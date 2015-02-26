package Cron

import (
	"Config/Config"
	"Config/ReadFlags"
	"Cron/CronMessage"
	"Cron/CronTask"
	. "gopkg.in/check.v1"
	"log"
	"testing"
)

func TestCron(t *testing.T) {
	TestingT(t)
}

type CronTestsSuite struct{}

var _ = Suite(&CronTestsSuite{})

func (s *CronTestsSuite) TestCronInit(c *C) {

	config := Config.New(ReadFlags.New())
	cron := NewCron(config)

	c.Assert(cron, NotNil)
}

func (s *CronTestsSuite) TestCronAddTask(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "ls", []string{"-l"})

	err := cron.AddTask(t)

	c.Assert(err, IsNil)
}

func (s *CronTestsSuite) TestCronMarkTask(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "ls", []string{"-l"})

	cron.AddTask(t)
	err := cron.MarkTask("1", "new_mark")
	mark, err2 := cron.GetMarkTask("1")

	c.Assert(err, IsNil)
	c.Assert(err2, IsNil)
	c.Assert(mark, Equals, "new_mark")

}

func (s *CronTestsSuite) TestCronCountTask(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "ls", []string{"-l"})

	cron.AddTask(t)

	l := cron.CountTask()

	c.Assert(l, Equals, 1)
}

func (s *CronTestsSuite) TestCronCountMarkTask(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))

	t := CronTask.NewTask("1", "ls", []string{"-l"})
	cron.AddTask(t)
	cron.MarkTask("1", "new_mark")

	t2 := CronTask.NewTask("2", "ls", []string{"-l"})
	cron.AddTask(t2)
	cron.MarkTask("2", "new_mark")

	l := cron.GetMarkTaskCount("new_mark")

	c.Assert(l, Equals, 2)
}

func (s *CronTestsSuite) TestCronGetTask(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "ls", []string{"-l"})

	cron.AddTask(t)

	_, f1 := cron.GetTask("1")
	c.Assert(f1, Equals, true)

	_, f2 := cron.GetTask("2")
	c.Assert(f2, Equals, false)

}

func (s *CronTestsSuite) TestCronExeOne(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "perl", []string{"-e 'sleep(10);'"})

	cron.AddTask(t)

	cmd, err := cron.Exe("1")
	c.Assert(err, IsNil)
	c.Assert(cmd, NotNil)
}

func (s *CronTestsSuite) TestCronExeNotFound(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "perl", []string{"-e 'sleep(10);'"})

	cron.AddTask(t)

	cmd, err := cron.Exe("2")
	c.Assert(err, NotNil)
	c.Assert(cmd, IsNil)
}

func (s *CronTestsSuite) TestCronExe(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "lsweqwrqr-rewrwer", []string{"-l"})

	cron.AddTask(t)

	cmd, err := cron.Exe("1")
	c.Assert(err, NotNil)
	c.Assert(cmd, IsNil)
}

func (s *CronTestsSuite) TestCronExeBad(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", "lsweqwrqr-rewrwer", []string{"-l"})

	cron.AddTask(t)

	cmd, err := cron.Exe("2")
	c.Assert(err, NotNil)
	c.Assert(cmd, IsNil)
}

func (s *CronTestsSuite) TestCronDo(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test.pl`)
	cron.AddTask(t)
	info, err := cron.Do("1")

	ResOut, ResErr := readAllFromCrom(cron.GetChannelProscessFrom())

	log.Printf("%s\n%s\n", ResOut, ResErr)

	c.Assert(err, IsNil)
	c.Assert(info, NotNil)
}

func (s *CronTestsSuite) TestCronReadLogData(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test_log_data.pl`)

	cron.AddTask(t)

	cron.Do("1")

	ResOut, ResErr := readAllFromCrom(cron.GetChannelProscessFrom())

	c.Assert(len(ResErr), Equals, 5172)
	c.Assert(len(ResOut), Equals, 5172)
}

func (s *CronTestsSuite) TestCronDoBad(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test_log_data.pl`)

	cron.AddTask(t)

	info, err := cron.Do("2")

	c.Assert(err, NotNil)
	c.Assert(info, IsNil)
}

func (s *CronTestsSuite) TestCronCanPassToStdIN(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test2.pl`)

	cron.AddTask(t)

	info, err := cron.Do("1")

	c.Assert(err, IsNil)

	CronMessage.Send(info.StdIn, "", "stdin", []byte("qwerty\n"), nil)

	ResOut, ResErr := readAllFromCrom(cron.GetChannelProscessFrom())

	c.Assert(string(ResOut), Equals, "ytrewq\n")
	c.Assert(string(ResErr), Equals, "qwerty\n")
}

func (s *CronTestsSuite) TestCronKillProccess(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))

	countRunInit := cron.GetMarkTaskCount("run")
	countStoppedInit := cron.GetMarkTaskCount("stop")

	t := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test2.pl`)

	cron.AddTask(t)

	info, err := cron.Do("1")

	c.Assert(err, IsNil)

	countRunBefor := cron.GetMarkTaskCount("run")
	countStoppedBefor := cron.GetMarkTaskCount("stop")

	CronMessage.Send(info.StdIn, "", "kill", []byte{}, nil)

	<-cron.GetChannelProscessFrom()

	countRunAfter := cron.GetMarkTaskCount("run")
	countStoppedAfter := cron.GetMarkTaskCount("stop")

	c.Assert(countRunInit, Equals, 0)
	c.Assert(countStoppedInit, Equals, 0)

	c.Assert(countRunBefor, Equals, 1)
	c.Assert(countStoppedBefor, Equals, 0)

	c.Assert(countRunAfter, Equals, 0)
	c.Assert(countStoppedAfter, Equals, 1)

}

func (s *CronTestsSuite) TestCronKillProccessWithCloseInCh(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))
	t := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test2.pl`)

	cron.AddTask(t)

	info, _ := cron.Do("1")

	countRunBefor := cron.GetMarkTaskCount("run")
	countStoppedBefor := cron.GetMarkTaskCount("stop")

	close(info.StdIn)

	<-cron.GetChannelProscessFrom()

	countRunAfter := cron.GetMarkTaskCount("run")
	countStoppedAfter := cron.GetMarkTaskCount("stop")

	c.Assert(countRunBefor, Equals, 1)
	c.Assert(countStoppedBefor, Equals, 0)

	c.Assert(countRunAfter, Equals, 0)
	c.Assert(countStoppedAfter, Equals, 1)
}

func (s *CronTestsSuite) TestCronStartManyTaskOneTime(c *C) {

	cron := NewCron(Config.New(ReadFlags.New()))

	t1 := CronTask.NewTask("1", `/Users/ostrovok/Work/go-bratok/test2.pl`)
	cron.AddTask(t1)
	t2 := CronTask.NewTask("2", `/Users/ostrovok/Work/go-bratok/test2.pl`)
	cron.AddTask(t2)
	t3 := CronTask.NewTask("3", `/Users/ostrovok/Work/go-bratok/test2.pl`)
	cron.AddTask(t3)

	info1, _ := cron.Do("1")
	info2, _ := cron.Do("2")
	CronMessage.Send(info1.StdIn, "", "stdin", []byte("123456\n"), nil)
	info3, _ := cron.Do("3")
	CronMessage.Send(info2.StdIn, "", "stdin", []byte("12345678\n"), nil)
	CronMessage.Send(info3.StdIn, "", "stdin", []byte("1234567891\n"), nil)

	ResOut, ResErr := readAllFromCromManyTask(3, cron.GetChannelProscessFrom())

	log.Printf("%+v\n%+v\n", ResOut, ResErr)

	//c.Assert(err, IsNil)
	c.Assert(len(ResOut["1"]), Equals, 7)
	c.Assert(len(ResOut["2"]), Equals, 9)
	c.Assert(len(ResOut["3"]), Equals, 11)
}

func readAllFromCrom(waitCh chan CronMessage.Mess) (string, string) {

	ResErr := ""
	ResOut := ""
	kk := true
	for kk {
		select {
		case m, isGood := <-waitCh:

			if !isGood {
				kk = false
				return ResOut, ResErr
			}

			log.Printf("m.Type: %s\n", m.Type)

			if m.Type == "finish" {
				kk = false
				return ResOut, ResErr
			}

			if m.Type == "stderr" {
				ResErr += string(m.Data)
			}

			if m.Type == "stdout" {
				ResOut += string(m.Data)
			}

		}
	}
	return ResOut, ResErr
}

func readAllFromCromManyTask(count int, waitCh chan CronMessage.Mess) (map[string]string, map[string]string) {

	ResErr := map[string]string{}
	ResOut := map[string]string{}
	NotFinish := map[string]bool{}

	kk := true
	for kk {
		select {
		case m, isGood := <-waitCh:

			if !isGood {
				return ResOut, ResErr
			}

			log.Printf("m.ID: %s, m.Type: %s\n", m.ID, m.Type)

			if m.ID == "" {
				continue
			}

			if _, f := ResErr[m.ID]; !f {
				ResErr[m.ID] = ""
			}

			if _, f := ResOut[m.ID]; !f {
				ResOut[m.ID] = ""
			}

			if m.Type == "finish" {
				NotFinish[m.ID] = false
			} else {
				NotFinish[m.ID] = true
			}

			if m.Type == "stderr" {
				ResErr[m.ID] += string(m.Data)
			}

			if m.Type == "stdout" {
				ResOut[m.ID] += string(m.Data)
			}

			if len(NotFinish) >= count {
				kk = false
				for _, v := range NotFinish {
					if v {
						kk = true
					}
				}
			}
		}
	}
	return ResOut, ResErr
}
