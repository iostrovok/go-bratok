package Config

import (
	"Config/CronScript"
	"Config/File"
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

func (s *ConfigTestsSuite) TestConfigLogFile(c *C) {
	c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	out := config.LogFile("my.log")
	c.Assert(out, Equals, "/tmp/my.log")
}

func (s *ConfigTestsSuite) TestConfigCheckInitNew(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	err := config.InitNew(map[string]interface{}{})

	c.Assert(err, IsNil)
}

func (s *ConfigTestsSuite) TestConfigAddCronScript(c *C) {
	//c.Skip("Not now")

	script := CronScript.New("1", "ls", "-a", "-r")
	config := New(ReadFlags.NewTest())

	err := config.AddCronScript(script)
	c.Assert(err, IsNil)
}

func (s *ConfigTestsSuite) TestConfigGetScript(c *C) {
	//c.Skip("Not now")

	script1 := CronScript.New("1", "ls", "-a", "-r")
	script2 := CronScript.New("2", "ls", "-a", "-r")
	config := New(ReadFlags.NewTest())

	config.AddCronScript(script1)
	config.AddCronScript(script2)

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
	config.AddCronScript(script)

	d := time.Now().Add(1 * time.Minute)

	list := config.StartNow(d)

	c.Assert(list, NotNil)
	c.Assert(len(list), Equals, 1)
}

func (s *ConfigTestsSuite) TestConfigParseConfigData(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	err := config.FromLine(File.TestFileLine())

	c.Assert(err, IsNil)
	c.Assert(len(config.ScriptsList()), Equals, 4)
	c.Assert(len(config.ServersList()), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigLoadConfigData(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	config.FromLine(File.TestFileLine())

	log.Printf("TestConfigLoadConfigData config: %+v\n", config)

	c.Assert(config.staticDir, Equals, "workstation_STATIC-DIR")
	c.Assert(len(config.ScriptsList()), Equals, 4)
	c.Assert(len(config.ServersList()), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigLoadConfigData2(c *C) {
	// We don't check the set dir/file now
	c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	config.FromLine(File.TestFileLine())

	c.Assert(config.LogFile(), Equals, "DIR-NEW/FILE")
	c.Assert(config.LogFile("FILE-NEW"), Equals, "DIR-NEW/FILE-NEW")
}

func (s *ConfigTestsSuite) TestConfigLoadHTTPLine(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())

	err := config.LoadHTTPLine(File.TestFileLineHttp())
	c.Assert(err, IsNil)

	c.Assert(config.LogFile(), Equals, "workstation_DIR/workstation_FILE")
}

func (s *ConfigTestsSuite) TestConfigGetHTTPData(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	p := config.ServersList()
	c.Assert(len(p), Equals, 2)

	server, err := config.GetHTTPData()
	c.Assert(err, IsNil)
	c.Assert(server, NotNil)

	c.Assert(server.IP, Equals, "127.0.0.1")
	c.Assert(server.Port, Equals, 21222)
}

func (s *ConfigTestsSuite) TestConfigGetHTTPDataError(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	server, err := config.GetHTTPData()

	c.Assert(err, IsNil)
	c.Assert(server, NotNil)
	c.Assert(server.IP, Equals, "127.0.0.1")
}

func (s *ConfigTestsSuite) TestConfigGetServer(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())
	server, find := config.GetServer("somethere")

	c.Assert(find, Equals, true)
	c.Assert(server, NotNil)
	c.Assert(server.Host, Equals, "wks-l")
	c.Assert(server.Port, Equals, 21223)
}

func (s *ConfigTestsSuite) TestConfigScriptsList_1(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())
	scripts := config.ScriptsList()

	c.Assert(len(scripts), Equals, 4)
}

func (s *ConfigTestsSuite) TestConfigScriptsList_2(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	scripts := config.ScriptsList("somethere")
	c.Assert(len(scripts), Equals, 1)

	scripts = config.ScriptsList("bad-server-id")
	c.Assert(len(scripts), Equals, 0)
}

func (s *ConfigTestsSuite) TestConfigGetConfigDataByte(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	b, err := config.GetConfigDataByte()
	c.Assert(b, NotNil)
	c.Assert(err, IsNil)

	str := string(b)
	log.Printf("TestConfigGetConfigDataByte: %s\n", str)

	c.Assert(len(str), Not(Equals), 0)
	c.Assert(strings.Index(str, "scripts"), Not(Equals), 0)
	c.Assert(strings.Index(str, "is_master"), Not(Equals), 0)

}

func (s *ConfigTestsSuite) TestConfigNewFromRemout(c *C) {
	//c.Skip("Not now")

	config, err := NewFromRemout(ReadFlags.NewTest(), File.TestFileLineHttp())

	c.Assert(config, NotNil)
	c.Assert(err, IsNil)

}

func (s *ConfigTestsSuite) TestConfigServersList(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.FromLine(File.TestFileLine())

	sList := config.ServersList()

	c.Assert(sList, NotNil)
	c.Assert(len(sList), Equals, 2)
}

func (s *ConfigTestsSuite) TestConfigReplaceServer(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	// Check data in start
	sList := config.ServersList()
	c.Assert(sList, NotNil)
	c.Assert(len(sList), Equals, 2)

	server, _ := config.GetServer("somethere")

	server2 := server.Clone()
	server2.IP = "222.222.222.111"

	find := config.ReplaceServer(server2)
	c.Assert(find, Equals, true)

	// Check result
	sList = config.ServersList()
	c.Assert(sList, NotNil)
	c.Assert(len(sList), Equals, 2)

	server, _ = config.GetServer("somethere")
	c.Assert(server.IP, Equals, "222.222.222.111")

}

func (s *ConfigTestsSuite) TestConfigReplaceServerADD(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	// Check data in start
	sList := config.ServersList()
	c.Assert(sList, NotNil)
	c.Assert(len(sList), Equals, 2)

	server, _ := config.GetServer("somethere")

	server2 := server.Clone()
	server2.ID = "nothere"
	server2.IP = "222.222.222.111"

	find := config.ReplaceServer(server2)
	c.Assert(find, Equals, false)

	// Check result
	sList = config.ServersList()
	c.Assert(sList, NotNil)
	c.Assert(len(sList), Equals, 3)

	server, _ = config.GetServer("nothere")
	c.Assert(server.IP, Equals, "222.222.222.111")

}

//ReplaceScript(script *CronScript.Script) bool
func (s *ConfigTestsSuite) TestConfigReplaceScript(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	list := config.ScriptsList()
	c.Assert(len(list), Equals, 4)

	script, find := config.GetScript("ls1")
	c.Assert(find, Equals, true)
	c.Assert(script.Exe, Not(Equals), "SUPER-PUPER-SCRIPT")

	script.SetExe("SUPER-PUPER-SCRIPT")

	isReplace := config.ReplaceScript(script)
	c.Assert(isReplace, Equals, true)
	list = config.ScriptsList()
	c.Assert(len(list), Equals, 4)

	script, find = config.GetScript("ls1")
	c.Assert(find, Equals, true)
	c.Assert(script.Exe, Equals, "SUPER-PUPER-SCRIPT")

}

//ReplaceScript(script *CronScript.Script) bool
func (s *ConfigTestsSuite) TestConfigReplaceScriptADD(c *C) {
	//c.Skip("Not now")

	config := New(ReadFlags.NewTest())
	config.LoadHTTPLine(File.TestFileLineHttp())

	list := config.ScriptsList()
	c.Assert(len(list), Equals, 4)

	scriptOld, findOld := config.GetScript("ls1")
	c.Assert(findOld, Equals, true)
	c.Assert(scriptOld.Exe, Not(Equals), "SUPER-PUPER-SCRIPT")
	c.Assert(scriptOld.ID, Not(Equals), "SUPER-PUPER-SCRIPT")

	scriptNew := scriptOld.Clone()
	scriptNew.ID = "SUPER-PUPER-SCRIPT"
	scriptNew.SetExe("SUPER-PUPER-SCRIPT")

	isReplace := config.ReplaceScript(scriptNew)
	log.Printf("%+v\n", config.Scripts)
	c.Assert(isReplace, Equals, false)
	list = config.ScriptsList()
	c.Assert(len(list), Equals, 5)

	script, find := config.GetScript("SUPER-PUPER-SCRIPT")
	c.Assert(find, Equals, true)
	c.Assert(script.Exe, Equals, "SUPER-PUPER-SCRIPT")
}
