package HistoryMem

import (
	//"Config/Config"
	"Config/CronScript"
	"Cron/Cron"
	"Cron/CronInfo"
	"Cron/CronMessage"
	//"errors"
	. "gopkg.in/check.v1"
	"testing"
)

func TestHistoryMem(t *testing.T) {
	TestingT(t)
}

type HistoryMemTestsSuite struct{}

var _ = Suite(&HistoryMemTestsSuite{})

func (s *HistoryMemTestsSuite) TestHistoryMemInit(c *C) {

	hm := New()
	c.Assert(hm, NotNil)
}

func (s *HistoryMemTestsSuite) TestHistoryMemAdd(c *C) {

	hm := New()
	info := &CronInfo.MyCronInfo{
		StdOut: &CronMessage.Mess{},
		StdErr: &CronMessage.Mess{},
		StdRes: &CronMessage.Mess{},
		Info:   &Cron.Info{},
		Script: &CronScript.Script{},
	}

	hm.Add(info)
	c.Assert(len(hm.History), Equals, 1)
}

func (s *HistoryMemTestsSuite) TestHistoryMemAddLong(c *C) {

	hm := New()
	for i := 0; i < 3000; i++ {
		info := &CronInfo.MyCronInfo{
			StdOut: &CronMessage.Mess{},
			StdErr: &CronMessage.Mess{},
			StdRes: &CronMessage.Mess{},
			Info:   &Cron.Info{},
			Script: &CronScript.Script{},
		}
		hm.Add(info)
	}

	c.Assert(len(hm.History), Equals, 2000)
}

func (s *HistoryMemTestsSuite) TestHistoryMemEmpty(c *C) {

	hm := New()
	mes, find := hm.Last()

	c.Assert(find, Equals, false)
	c.Assert(mes, IsNil)
}

func (s *HistoryMemTestsSuite) TestHistoryMemLast(c *C) {

	hm := New()
	info := &CronInfo.MyCronInfo{
		StdOut: &CronMessage.Mess{},
		StdErr: &CronMessage.Mess{},
		StdRes: &CronMessage.Mess{},
		Info:   &Cron.Info{},
		Script: &CronScript.Script{},
	}

	hm.Add(info)
	mes, find := hm.Last()

	c.Assert(find, Equals, true)
	c.Assert(mes, NotNil)
}

func (s *HistoryMemTestsSuite) TestHistoryMemLastExe(c *C) {

	hm := New()
	info := &CronInfo.MyCronInfo{
		StdOut: &CronMessage.Mess{},
		StdErr: &CronMessage.Mess{},
		StdRes: &CronMessage.Mess{},
		Info:   &Cron.Info{},
		Script: &CronScript.Script{ID: "23"},
	}

	hm.Add(info)
	mes, find := hm.LastExe("23")

	// Check LastExe for 23
	c.Assert(find, Equals, true)
	// Check LastExe for 23
	c.Assert(mes, NotNil)
	// Check LastExe for 23
	c.Assert(len(mes), Equals, 1)

	mes, find = hm.LastExe("33")

	// Check LastExe for 33
	c.Assert(find, Equals, false)
	// Check LastExe for 33
	c.Assert(mes, NotNil)
	// Check LastExe for 33
	c.Assert(len(mes), Equals, 0)
}
