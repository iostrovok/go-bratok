package Config

import (
	"Config/CronScript"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"log"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	TestingT(t)
}

type ConfigTestsSuite struct{}

var _ = Suite(&ConfigTestsSuite{})

func (s *ConfigTestsSuite) TestConfigCanNew(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())

	c.Assert(config, NotNil)
}

func (s *ConfigTestsSuite) TestConfigScriptLogDir(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())

	out := config.ScriptLogDir("/wqeqweq/")
	c.Assert(out, Equals, "/wqeqweq/")
}

func (s *ConfigTestsSuite) TestConfigScriptLogFile(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())

	out := config.ScriptLogFile("my.log")
	c.Assert(out, Equals, "/tmp/my.log")
}

func (s *ConfigTestsSuite) TestConfigCheckInitNew(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())
	err := config.InitNew(map[string]interface{}{})

	c.Assert(err, IsNil)
}

func (s *ConfigTestsSuite) TestConfigAddScript(c *C) {

	//c.Skip("Not now")

	script := CronScript.New("1", "ls", "-a", "-r")
	config := New(ReadFlags.New())

	log.Printf("script: %+v\nconfig: %+v\n", script, config)

	err := config.AddScript(script)
	c.Assert(err, IsNil)
}

func (s *ConfigTestsSuite) TestConfigGetScript(c *C) {

	//c.Skip("Not now")

	script1 := CronScript.New("1", "ls", "-a", "-r")
	script2 := CronScript.New("2", "ls", "-a", "-r")
	config := New(ReadFlags.New())

	log.Printf("script: %+v\nconfig: %+v\n", script1, config)

	config.AddScript(script1)
	config.AddScript(script2)

	script, find := config.GetScript("1")

	c.Assert(find, Equals, true)
	c.Assert(script, NotNil)

	script, find = config.GetScript("34")

	c.Assert(find, Equals, false)
	c.Assert(script, IsNil)
}

func (s *ConfigTestsSuite) TestConfigStartNow(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())
	script := CronScript.New("1", "ls", "-a", "-r")
	script.SetTime("*/1", "*", "*", "*", "*")
	config.AddScript(script)

	d := time.Now().Add(1 * time.Minute)

	list := config.StartNow(d)

	c.Assert(list, NotNil)
	c.Assert(len(list), Equals, 1)
}

func (s *ConfigTestsSuite) TestConfigParseConfigData(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())

	err := config._parseConfigData(TestFileLine())

	c.Assert(err, IsNil)
	c.Assert(len(*config.ConfigData.Scripts), Equals, 2)
	c.Assert(len(*config.ConfigData.Servers), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigLoadConfigData(c *C) {

	//c.Skip("Not now")

	config := New(ReadFlags.New())

	log.Printf("\nTestFileLine: %s\n\n", TestFileLine())

	config._parseConfigData(TestFileLine())
	err := config._loadConfigData()

	c.Assert(err, IsNil)
	c.Assert(config.scriptLogDir, Equals, "DIR")
	c.Assert(config.scriptLogFile, Equals, "FILE")
	c.Assert(config.staticDir, Equals, "STATIC-DIR")
	c.Assert(len(*config.ConfigData.Scripts), Equals, 2)
	c.Assert(len(*config.ConfigData.Servers), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigLoadConfigData2(c *C) {

	config := New(ReadFlags.New())

	config._parseConfigData(TestFileLine())
	err := config._loadConfigData()

	c.Assert(err, IsNil)

	c.Assert(config.ScriptLogDir(), Equals, "DIR")
	c.Assert(config.ScriptLogDir("DIR-NEW/"), Equals, "DIR-NEW/")

	c.Assert(config.ScriptLogFile(), Equals, "DIR-NEW/FILE")
	c.Assert(config.ScriptLogFile("FILE-NEW"), Equals, "DIR-NEW/FILE-NEW")

}
