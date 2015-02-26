package Manager

import (
	"Config/Config"
	"Config/CronScript"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	TestingT(t)
}

type ManagerTestsSuite struct{}

var _ = Suite(&ManagerTestsSuite{})

func (s *ManagerTestsSuite) TestManagerInit(c *C) {

	config := Config.New(ReadFlags.New())
	cron := New(config)

	c.Assert(cron, NotNil)
}

func (s *ManagerTestsSuite) TestManagerReadData(c *C) {

	script := CronScript.New("1", "ls", "-a", "-r")
	script.SetTime("*/1", "*", "*", "*", "*")

	config := Config.New(ReadFlags.New())
	config.AddScript(script)

	manager := New(config)

	c.Assert(manager, NotNil)
}

func (s *ManagerTestsSuite) TestManagerTimerAction(c *C) {

	script := CronScript.New("1", "ls", "-a", "-r")
	script.SetTime("*/1", "*", "*", "*", "*")

	script2 := CronScript.New("2", "ls", "-a", "-r")
	script2.SetTime("*/1", "*", "*", "*", "*")

	script3 := CronScript.New("3", "ls", "-a", "-r")
	script3.SetTime("1", "1", "1", "1", "1")

	config := Config.New(ReadFlags.New())
	config.AddScript(script)
	config.AddScript(script2)
	config.AddScript(script3)

	d, err := time.ParseDuration("5m")
	c.Assert(err, IsNil)

	manager := New(config)
	manager.TimerAction(time.Now().Add(d))

	c.Assert(manager.Currently, NotNil)
	c.Assert(len(manager.Currently), Equals, 2)
}
