package Config

import (
	"Config/CronScript"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"log"
	"strings"
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

	config := New(ReadFlags.NewTest())

	c.Assert(config, NotNil)
}

func (s *ConfigTestsSuite) TestConfigScriptLogDir(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	out := config.ScriptLogDir("/wqeqweq/")
	c.Assert(out, Equals, "/wqeqweq/")
}

func (s *ConfigTestsSuite) TestConfigScriptLogFile(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	out := config.ScriptLogFile("my.log")
	c.Assert(out, Equals, "/tmp/my.log")
}

func (s *ConfigTestsSuite) TestConfigCheckInitNew(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	err := config.InitNew(map[string]interface{}{})

	c.Assert(err, IsNil)
}

func (s *ConfigTestsSuite) TestConfigAddScript(c *C) {
	//c.Skip("Not now")

	script := CronScript.New("1", "ls", "-a", "-r")
	config := New(ReadFlags.NewTest())

	err := config.AddScript(script)
	c.Assert(err, IsNil)
}

func (s *ConfigTestsSuite) TestConfigGetScript(c *C) {
	//c.Skip("Not now")

	script1 := CronScript.New("1", "ls", "-a", "-r")
	script2 := CronScript.New("2", "ls", "-a", "-r")
	config := New(ReadFlags.NewTest())

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

	config := New(ReadFlags.NewTest())
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

	config := New(ReadFlags.NewTest())

	err := config._parseConfigData(TestFileLine())

	c.Assert(err, IsNil)
	c.Assert(len(*config.ConfigData.Scripts), Equals, 4)
	c.Assert(len(*config.ConfigData.Servers), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigLoadConfigData(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	config._parseConfigData(TestFileLine())
	err := config._loadConfigData()

	c.Assert(err, IsNil)
	c.Assert(config.scriptLogDir, Equals, "DIR")
	c.Assert(config.scriptLogFile, Equals, "FILE")
	c.Assert(config.staticDir, Equals, "STATIC-DIR")
	c.Assert(len(*config.ConfigData.Scripts), Equals, 4)
	c.Assert(len(*config.ConfigData.Servers), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigLoadConfigData2(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	config._parseConfigData(TestFileLine())
	err := config._loadConfigData()

	c.Assert(err, IsNil)

	c.Assert(config.ScriptLogDir(), Equals, "DIR")
	c.Assert(config.ScriptLogDir("DIR-NEW/"), Equals, "DIR-NEW/")

	c.Assert(config.ScriptLogFile(), Equals, "DIR-NEW/FILE")
	c.Assert(config.ScriptLogFile("FILE-NEW"), Equals, "DIR-NEW/FILE-NEW")
}

func (s *ConfigTestsSuite) TestConfigLoadConfigFileFromLine(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	err := config.LoadConfigFileFromLine(TestFileLineHttp())
	c.Assert(err, IsNil)

	c.Assert(config.ScriptLogDir(), Equals, "DIR")
	c.Assert(config.ScriptLogFile(), Equals, "DIRFILE")
}

func (s *ConfigTestsSuite) TestConfigGetHttpData(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())

	p := *config.ConfigData.Servers
	c.Assert(len(p), Equals, 2)

	server, err := config.GetHttpData()
	c.Assert(err, IsNil)
	c.Assert(server, NotNil)

	c.Assert(server.IP, Equals, "127.0.0.1")
	c.Assert(server.Port, Equals, 21222)
}

func (s *ConfigTestsSuite) TestConfigGetHttpDataError(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())

	p := []*Server{}
	config.ConfigData.Servers = &p
	server, err := config.GetHttpData()

	c.Assert(err, NotNil)
	c.Assert(server, IsNil)
}

func (s *ConfigTestsSuite) TestConfigGetServer(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())
	server, find := config.GetServer("somethere")

	c.Assert(find, Equals, true)
	c.Assert(server, NotNil)
	c.Assert(server.Host, Equals, "wks-l")
	c.Assert(server.Port, Equals, 21222)
}

func (s *ConfigTestsSuite) TestConfigScriptsList_1(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())
	scripts := config.ScriptsList()

	c.Assert(len(scripts), Equals, 4)
}

func (s *ConfigTestsSuite) TestConfigScriptsList_2(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())

	scripts := config.ScriptsList("somethere")
	c.Assert(len(scripts), Equals, 1)

	scripts = config.ScriptsList("bad-server-id")
	c.Assert(len(scripts), Equals, 0)
}

func (s *ConfigTestsSuite) TestConfigUpdateId(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())

	oldId := config.ConfigData.ConfigID
	config.UpdateId()
	newId := config.ConfigData.ConfigID

	c.Assert(oldId, Not(Equals), newId)

	if newId < oldId {
		c.Fatalf("Bad old-new id: %d - %d\n", oldId, newId)
	}
}

func (s *ConfigTestsSuite) TestConfigGetConfigDataByte(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadConfigFileFromLine(TestFileLineHttp())

	b, err := config.GetConfigDataByte()
	c.Assert(b, NotNil)
	c.Assert(err, IsNil)

	str := string(b)
	log.Printf("TestConfigGetConfigDataByte: %s\n", str)

	c.Assert(len(str), Not(Equals), 0)
	c.Assert(strings.Index(str, "scripts"), Not(Equals), 0)
	c.Assert(strings.Index(str, "is_master"), Not(Equals), 0)

}
