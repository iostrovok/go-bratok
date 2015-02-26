package Cron

import (
	"Config/Config"
	"Cron/CronMessage"
	"Cron/CronTask"
	"Cron/CronTimer"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

type Info struct {
	TaskID string
	Proc   *os.Process
	StdIn  chan CronMessage.Mess
	WaitCh chan CronMessage.Mess
	Wait   sync.WaitGroup
}

// It is main structure
type Cron struct {
	config         *Config.Config
	tasks          *CronTask.TaskList
	ChCronFrom     chan CronMessage.Mess
	ChCronTo       chan CronMessage.Mess
	chProscessFrom chan CronMessage.Mess
	chProscessTo   chan CronMessage.Mess
	chTimerFrom    chan CronMessage.Mess
	chTimerTo      chan CronMessage.Mess
	timer          *CronTimer.Timer
	wgCron         sync.WaitGroup
	wgProc         sync.WaitGroup
}

func NewCron(config *Config.Config) *Cron {

	c := Cron{}
	c.tasks = CronTask.New()

	c.chProscessFrom = CronMessage.Channel()
	c.chTimerFrom = CronMessage.Channel()
	c.ChCronTo = CronMessage.Channel()
	c.ChCronFrom = CronMessage.Channel()
	c.timer = CronTimer.Start(1, c.wgCron, c.chTimerFrom)

	return &c
}

func (c *Cron) GetChannelProscessFrom() chan CronMessage.Mess {
	return c.chProscessFrom
}

func (c *Cron) AddTask(t *CronTask.Task) error {
	if err := c.tasks.AddTask(t); err != nil {
		return err
	}
	return c.tasks.Mark(t.Id(), "stop")
}

func (c *Cron) CountTask() int {
	return c.tasks.Len()
}

func (c *Cron) GetTask(id string) (*CronTask.Task, bool) {
	return c.tasks.Get(id)
}

func (c *Cron) MarkTask(id, key string) error {
	return c.tasks.Mark(id, key)
}

func (c *Cron) GetMarkTask(id string) (string, error) {
	return c.tasks.GetMark(id)
}

func (c *Cron) GetMarkTaskCount(key string) int {
	return c.tasks.GetMarkTaskCount(key)
}

func (c *Cron) Exe(id string) (*exec.Cmd, error) {

	t, find := c.tasks.Get(id)
	if !find {
		return nil, fmt.Errorf("Task [%s] not found", id)
	}

	c.tasks.Mark(id, "run")

	cmd := exec.Command(t.Exe(), t.Params()...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	err = cmd.Wait()
	c.tasks.Mark(id, "stop")
	return cmd, nil
}

// "Do" starts one outside process and makes channel in/to it
func (c *Cron) Do(id string) (*Info, error) {

	task, findTask := c.GetTask(id)
	if !findTask {
		return nil, fmt.Errorf("Task [%s] not found", id)
	}

	cmd := exec.Command(task.Exe(), task.Params()...)

	c.MarkTask(id, "run")

	stdIn, errStdinPipe := cmd.StdinPipe()
	if errStdinPipe != nil {
		c.MarkTask(id, "stop")
		return nil, errStdinPipe
	}

	var wg sync.WaitGroup
	out := &Info{
		TaskID: id,
		Proc:   cmd.Process,
		StdIn:  CronMessage.Channel(),
		Wait:   wg,
	}

	// Stdout buffer && Attach buffer to command
	cmdOutput := &bytes.Buffer{}
	cmdErrput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	cmd.Stderr = cmdErrput

	if err := cmd.Start(); err != nil {
		c.MarkTask(id, "stop")
		return nil, err
	}

	mu := &sync.Mutex{}
	chRes := CronMessage.Channel()

	mu.Lock()
	go func() {
		noFinish := true
		mu.Unlock()
		for noFinish {
			select {
			case m, isGood := <-chRes:
				if isGood {
					c.chProscessFrom <- m
				} else {
					noFinish = false
				}
			case m, isGood := <-out.StdIn:
				log.Printf("out.StdIn: %+v\n", m)
				if !isGood {
					cmd.Process.Kill()
					break
				}

				if m.Type == "kill" {
					cmd.Process.Kill()
					break
				}
				if m.Type == "stdin" {
					io.Copy(stdIn, bytes.NewBufferString(string(m.Data)))
				}
			}
		}

		CronMessage.Send(c.chProscessFrom, id, "stderr", cmdErrput.Bytes(), nil)
		CronMessage.Send(c.chProscessFrom, id, "stdout", cmdOutput.Bytes(), nil)
		CronMessage.Send(c.chProscessFrom, id, "finish", []byte{}, nil)

		c.MarkTask(id, "stop")
		wg.Wait()
	}()

	mu.Lock()
	go func() {
		mu.Unlock()
		err := cmd.Wait()
		CronMessage.Send(chRes, id, "result", []byte{}, err)
		close(chRes)
	}()

	return out, nil
}

/*--------------------------------*/
