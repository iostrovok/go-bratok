package CronScript

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

const (
	layout string = "2006-01-02 15:04:05"
)

func TestConfigCronScript(t *testing.T) {
	TestingT(t)
}

type ConfigCronScriptTestsSuite struct{}

var _ = Suite(&ConfigCronScriptTestsSuite{})

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptInit(c *C) {

	script := New("1", "ls", "-a", "-r")

	c.Assert(script, NotNil)
	c.Assert(script.Exe, Equals, "ls")
	c.Assert(script.ID, Equals, "1")
	c.Assert(script.Params[0], Equals, "-a")
	c.Assert(script.Params[1], Equals, "-r")
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptCheckSetTime(c *C) {
	script := New("1", "ls", "-a", "-r")
	err := script.SetTime("*", "*", "*", "*", "*")
	c.Assert(err, IsNil)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptCheckSetTimeBad(c *C) {
	script := New("1", "ls", "-a", "-r")
	err := script.SetTime("s")
	c.Assert(err, NotNil)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptCheckSetTimeTwice(c *C) {
	script := New("1", "ls", "-a", "-r")
	script.SetTime("*/7", "*/6", "*", "*", "*")
	err := script.SetTime("*/3", "*/5", "*", "*", "*")
	c.Assert(err, IsNil)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkStartTime(c *C) {
	D, _ := time.ParseDuration("2s")

	script := New("1", "ls", "-a", "-r")
	c.Assert(script.IsWork, Equals, false)

	script.MarkStart(time.Now())

	c.Assert(time.Now().Before(script.LastStart.Add(D)), Equals, true)
	c.Assert(script.LastStart.Before(time.Now().Add(D)), Equals, true)
	c.Assert(script.IsWork, Equals, true)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkStart(c *C) {
	D, _ := time.ParseDuration("2s")

	script := New("1", "ls", "-a", "-r")
	c.Assert(script.IsWork, Equals, false)

	script.MarkStart()

	c.Assert(time.Now().Before(script.LastStart.Add(D)), Equals, true)
	c.Assert(script.LastStart.Before(time.Now().Add(D)), Equals, true)
	c.Assert(script.IsWork, Equals, true)

}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkFinishTime(c *C) {
	D, _ := time.ParseDuration("2s")

	script := New("1", "ls", "-a", "-r")
	c.Assert(script.IsWork, Equals, false)
	script.MarkFinish(time.Now())

	c.Assert(time.Now().Before(script.LastFinish.Add(D)), Equals, true)
	c.Assert(script.LastFinish.Before(time.Now().Add(D)), Equals, true)
	c.Assert(script.IsWork, Equals, false)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkFinish(c *C) {
	D, _ := time.ParseDuration("2s")

	script := New("1", "ls", "-a", "-r")
	c.Assert(script.IsWork, Equals, false)
	script.MarkFinish()

	c.Assert(time.Now().Before(script.LastFinish.Add(D)), Equals, true)
	c.Assert(script.LastFinish.Before(time.Now().Add(D)), Equals, true)
	c.Assert(script.IsWork, Equals, false)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkStartMarkFinish(c *C) {
	D, _ := time.ParseDuration("2s")

	script := New("1", "ls", "-a", "-r")
	c.Assert(script.IsWork, Equals, false)

	script.MarkStart()

	c.Assert(time.Now().Before(script.LastStart.Add(D)), Equals, true)
	c.Assert(script.LastStart.Before(time.Now().Add(D)), Equals, true)
	c.Assert(script.IsWork, Equals, true)

	script.MarkFinish()

	c.Assert(time.Now().Before(script.LastFinish.Add(D)), Equals, true)
	c.Assert(script.LastFinish.Before(time.Now().Add(D)), Equals, true)
	c.Assert(script.IsWork, Equals, false)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkStartStartNow(c *C) {
	d, _ := time.ParseDuration("1s")

	script := New("1", "ls", "-a", "-r")
	script.SetTime("*", "*", "*", "*", "*")

	isStart := script.StartNow(time.Now().Add(d))

	// Our precision equals 1 min
	c.Assert(isStart, Equals, false)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkStartStartNow2(c *C) {
	d, _ := time.ParseDuration("1m")

	script := New("1", "ls", "-a", "-r")
	script.SetTime("*", "*", "*", "*", "*")

	isStart := script.StartNow(time.Now().Add(d))

	c.Assert(isStart, Equals, true)
}

func (s *ConfigCronScriptTestsSuite) TestConfigCronScriptMarkStartStartNowBad(c *C) {
	d, err := time.ParseDuration("5m")
	c.Assert(err, IsNil)

	script := New("1", "ls", "-a", "-r")
	script.SetTime("*/3", "*/5", "*", "*", "*")
	script.SetTime("*/7", "*/6", "*", "*", "*")
	script.SetTime("*", "*", "*", "*", "*")

	script.MarkStart()

	// If script is working now we cann't start it
	isStart := script.StartNow(time.Now().Add(d))

	c.Assert(isStart, Equals, false)
}
