package CronScript

import (
	"BUtils"
	"Config/CronScheduler"
	"fmt"
	"strings"
	"time"
)

const (
	HourDurSecond  float64 = 60 * 60
	DayDurSecond   float64 = 60 * 60 * 24
	WeekDurSecond  float64 = 60 * 60 * 24 * 7
	MonthDurSecond float64 = 60 * 60 * 24 * 31
	YearDurSecond  float64 = 60 * 60 * 24 * 366
)

type Script struct {
	ID         string                     `json:"id"`
	Time       []*CronScheduler.Scheduler `json:"-"`
	TimeStr    []string                   `json:"times"`
	Exe        string                     `json:"exe"`
	Params     []string                   `json:"params"`
	LastStart  time.Time                  `json:"last_time"`
	LastFinish time.Time                  `json:"last_finish"`
	IsWork     bool                       `json:"is_work"`
	NextStart  *time.Time                 `json:"next_start"`
	Env        []string                   `json:"env"`
}

func New(id, exe string, params ...string) *Script {
	s := Script{
		ID:         id,
		Exe:        exe,
		Time:       []*CronScheduler.Scheduler{},
		TimeStr:    []string{},
		Params:     params,
		LastStart:  time.Now(),
		LastFinish: time.Now(),
		IsWork:     false,
		Env:        []string{},
	}

	return &s
}

func (script *Script) Clone() *Script {
	s := Script{
		ID:         script.ID,
		Exe:        script.Exe,
		Time:       script.CloneScheduler(),
		TimeStr:    BUtils.CopyStringsList(script.TimeStr),
		Params:     BUtils.CopyStringsList(script.Params),
		LastStart:  script.LastStart,
		LastFinish: script.LastFinish,
		IsWork:     script.IsWork,
		Env:        BUtils.CopyStringsList(script.Env),
	}

	return &s
}

func (script *Script) CloneScheduler() []*CronScheduler.Scheduler {
	out := make([]*CronScheduler.Scheduler, len(script.Time))
	for i, v := range script.Time {
		out[i] = v.Clone()
	}
	return out
}

func (script *Script) initNextTime(t time.Time) error {

	script.NextStart = nil

	if len(script.Time) == 0 {
		return fmt.Errorf("Not found scheduler")
	}

	nextTime := script.Time[0].GetNextTime(t)
	for i := 1; i < len(script.Time); i++ {
		nt := script.Time[i].GetNextTime(t)
		if nt.Before(nextTime) {
			nextTime = nt
		}
	}

	script.NextStart = &nextTime

	return nil
}

func (script *Script) SetEnv(m []string) error {
	script.Env = m
	return nil
}

func (script *Script) SetParams(m []string) error {
	script.Params = m
	return nil
}

func (script *Script) SetExe(m string) error {
	script.Exe = m
	return nil
}

func (script *Script) CleanTime() {
	script.Time = []*CronScheduler.Scheduler{}
	script.TimeStr = []string{}
}

func (script *Script) SetTime(m ...string) error {
	sched := CronScheduler.New()
	if err := sched.SetTime(m...); err != nil {
		return err
	}

	str := strings.Join(m, " ")
	script.TimeStr = append(script.TimeStr, str)
	script.Time = append(script.Time, sched)
	script.initNextTime(time.Now())
	return nil
}

func (script *Script) MarkStart(t ...time.Time) {
	if len(t) > 0 {
		script.LastStart = t[0]
	} else {
		script.LastStart = time.Now()
	}
	script.IsWork = true
}

func (script *Script) MarkFinish(t ...time.Time) {
	if len(t) > 0 {
		script.LastFinish = t[0]
	} else {
		script.LastFinish = time.Now()
	}
	script.initNextTime(script.LastFinish)
	script.IsWork = false
}

func (script *Script) StartNow(t time.Time) bool {

	if script.IsWork || script.NextStart == nil {
		return false
	}

	return t.After(*script.NextStart)
}

/* SORT DATA ------->

USING:
out := []*CronScript.Script{}
...
sort.Sort(CronScript.SortList(out))

*/
type SortList []*Script

func (d SortList) Len() int {
	return len(d)
}

func (d SortList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d SortList) Less(i, j int) bool {
	return d[i].ID < d[j].ID
}

/* <------- SORT DATA */
