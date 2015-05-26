package File

import (
	. "gopkg.in/check.v1"
	// "log"
	// "strings"
	"testing"
	// "time"
)

func TestConfigFile(t *testing.T) {
	TestingT(t)
}

type ConfigFileTestsSuite struct{}

var _ = Suite(&ConfigFileTestsSuite{})

func (s *ConfigFileTestsSuite) TestConfigFileNew(c *C) {

	//c.Skip("Not now")

	// file := New("ServerID", "configFile")
	file := New("workstation", "configFile")

	c.Assert(file, NotNil)
	c.Assert(file.Data, NotNil)
	c.Assert(file.Data.Scripts, NotNil)
	c.Assert(file.Data.Servers, NotNil)
	c.Assert(file.mu, NotNil)
	c.Assert(file.serverID, Equals, "workstation")
	c.Assert(file.configFile, Equals, "configFile")

}

func (s *ConfigFileTestsSuite) TestConfigFileSetAutoConfig(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	err := file.SetAutoConfig(true, "line")
	c.Assert(err, IsNil)

	err = file.SetAutoConfig(false)
	c.Assert(err, IsNil)

	err = file.SetAutoConfig(false, "line")
	c.Assert(err, IsNil)

	err = file.SetAutoConfig(true, "")
	c.Assert(err, NotNil)

	err = file.SetAutoConfig(true)
	c.Assert(err, NotNil)
}

func (s *ConfigFileTestsSuite) TestConfigFileFindServer(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	server := file.FindServer("workstation_asdsdfdasfsdfa")
	c.Assert(server, IsNil)

	server = file.FindServer("workstation")
	c.Assert(server, NotNil)

}

func (s *ConfigFileTestsSuite) TestConfigFileFromLine(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	err := file.FromLine(TestFileLine())
	c.Assert(err, IsNil)
}

func (s *ConfigFileTestsSuite) TestConfigFileDirsA(c *C) {
	//c.Skip("Not now")

	file := New("somethere", "configFile")
	file.FromLine(TestFileLine())

	c.Assert(file.ServerID(), Equals, "somethere")
	c.Assert(file.ConfigID(), Equals, int64(12312312))
	c.Assert(file.ScriptLogDir(), Equals, "DIR")
	c.Assert(file.StaticFilesDir(), Equals, "STATIC-DIR")
	c.Assert(file.ScriptLogFile(), Equals, "FILE")
}

func (s *ConfigFileTestsSuite) TestConfigFileDirsB(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	c.Assert(file.ServerID(), Equals, "workstation")
	c.Assert(file.ConfigID(), Equals, int64(12312312))
	c.Assert(file.ScriptLogDir(), Equals, "workstation_DIR")
	c.Assert(file.StaticFilesDir(), Equals, "workstation_STATIC-DIR")
	c.Assert(file.ScriptLogFile(), Equals, "workstation_FILE")
}

func (s *ConfigTestsSuite) TestConfigFileUpdateId(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	oldId := file.ConfigID
	file.UpdateId()
	newId := file.ConfigID

	c.Assert(oldId, Not(Equals), newId)

	if newId < oldId {
		c.Fatalf("Bad old-new id: %d - %d\n", oldId, newId)
	}
}

// func (s *ConfigFileTestsSuite) TestConfigFileScriptLogFile(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())

// 	out := file.ScriptLogFile("my.log")
// 	c.Assert(out, Equals, "/tmp/my.log")
// }

// func (s *ConfigFileTestsSuite) TestConfigFileCheckInitNew(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	err := file.InitNew(map[string]interface{}{})

// 	c.Assert(err, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileAddCronScript(c *C) {
// 	//c.Skip("Not now")

// 	script := CronScript.New("1", "ls", "-a", "-r")
// 	file := New(ReadFlags.NewTest())

// 	err := file.AddCronScript(script)
// 	c.Assert(err, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetScript(c *C) {
// 	//c.Skip("Not now")

// 	script1 := CronScript.New("1", "ls", "-a", "-r")
// 	script2 := CronScript.New("2", "ls", "-a", "-r")
// 	file := New(ReadFlags.NewTest())

// 	file.AddCronScript(script1)
// 	file.AddCronScript(script2)

// 	script, find := file.GetScript("1")

// 	c.Assert(find, Equals, true)
// 	c.Assert(script, NotNil)

// 	script, find = file.GetScript("34")

// 	c.Assert(find, Equals, false)
// 	c.Assert(script, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileStartNow(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	script := CronScript.New("1", "ls", "-a", "-r")
// 	script.SetTime("*/1", "*", "*", "*", "*")
// 	file.AddCronScript(script)

// 	d := time.Now().Add(1 * time.Minute)

// 	list := file.StartNow(d)

// 	c.Assert(list, NotNil)
// 	c.Assert(len(list), Equals, 1)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileParseFileData(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())

// 	err := file._parseFileData(TestConfigFileLine())

// 	c.Assert(err, IsNil)
// 	c.Assert(len(*file.FileData.Scripts), Equals, 4)
// 	c.Assert(len(*file.FileData.Servers), Equals, 2)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileLoadFileData(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())

// 	file._parseFileData(TestConfigFileLine())
// 	err := file._loadFileData()

// 	c.Assert(err, IsNil)
// 	c.Assert(file.scriptLogDir, Equals, "DIR")
// 	c.Assert(file.scriptLogFile, Equals, "FILE")
// 	c.Assert(file.staticDir, Equals, "STATIC-DIR")
// 	c.Assert(len(*file.FileData.Scripts), Equals, 4)
// 	c.Assert(len(*file.FileData.Servers), Equals, 2)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileLoadFileData2(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())

// 	file._parseFileData(TestConfigFileLine())
// 	err := file._loadFileData()

// 	c.Assert(err, IsNil)

// 	c.Assert(file.ScriptLogDir(), Equals, "DIR")
// 	c.Assert(file.ScriptLogDir("DIR-NEW/"), Equals, "DIR-NEW/")

// 	c.Assert(file.ScriptLogFile(), Equals, "DIR-NEW/FILE")
// 	c.Assert(file.ScriptLogFile("FILE-NEW"), Equals, "DIR-NEW/FILE-NEW")
// }

// func (s *ConfigFileTestsSuite) TestConfigFileLoadFileFileFromLine(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())

// 	err := file.LoadFileFileFromLine(TestConfigFileLineHttp())
// 	c.Assert(err, IsNil)

// 	c.Assert(file.ScriptLogDir(), Equals, "DIR")
// 	c.Assert(file.ScriptLogFile(), Equals, "DIRFILE")
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetHttpData(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	p := *file.FileData.Servers
// 	c.Assert(len(p), Equals, 2)

// 	server, err := file.GetHttpData()
// 	c.Assert(err, IsNil)
// 	c.Assert(server, NotNil)

// 	c.Assert(server.IP, Equals, "127.0.0.1")
// 	c.Assert(server.Port, Equals, 21222)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetHttpDataError(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	p := []*Server{}
// 	file.FileData.Servers = &p
// 	server, err := file.GetHttpData()

// 	c.Assert(err, NotNil)
// 	c.Assert(server, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetServer(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())
// 	server, find := file.GetServer("somethere")

// 	c.Assert(find, Equals, true)
// 	c.Assert(server, NotNil)
// 	c.Assert(server.Host, Equals, "wks-l")
// 	c.Assert(server.Port, Equals, 21222)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileScriptsList_1(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())
// 	scripts := file.ScriptsList()

// 	c.Assert(len(scripts), Equals, 4)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileScriptsList_2(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	scripts := file.ScriptsList("somethere")
// 	c.Assert(len(scripts), Equals, 1)

// 	scripts = file.ScriptsList("bad-server-id")
// 	c.Assert(len(scripts), Equals, 0)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileUpdateId(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	oldId := file.FileData.FileID
// 	file.UpdateId()
// 	newId := file.FileData.FileID

// 	c.Assert(oldId, Not(Equals), newId)

// 	if newId < oldId {
// 		c.Fatalf("Bad old-new id: %d - %d\n", oldId, newId)
// 	}
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetFileDataByte(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	b, err := file.GetFileDataByte()
// 	c.Assert(b, NotNil)
// 	c.Assert(err, IsNil)

// 	str := string(b)
// 	log.Printf("TestConfigFileGetFileDataByte: %s\n", str)

// 	c.Assert(len(str), Not(Equals), 0)
// 	c.Assert(strings.Index(str, "scripts"), Not(Equals), 0)
// 	c.Assert(strings.Index(str, "is_master"), Not(Equals), 0)

// }

// func (s *ConfigFileTestsSuite) TestConfigFileNewFromRemout(c *C) {
// 	//c.Skip("Not now")

// 	file, err := NewFromRemout(ReadFlags.NewTest(), TestConfigFileLineHttp())

// 	c.Assert(file, NotNil)
// 	c.Assert(err, IsNil)

// }

// func (s *ConfigFileTestsSuite) TestConfigFileServersList(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file._parseFileData(TestConfigFileLine())

// 	sList := file.ServersList()

// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 2)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileReplaceServer(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	// Check data in start
// 	sList := file.ServersList()
// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 2)

// 	server, _ := file.GetServer("somethere")

// 	server2 := server.Clone()
// 	server2.IP = "222.222.222.111"

// 	find := file.ReplaceServer(server2)
// 	c.Assert(find, Equals, true)

// 	// Check result
// 	sList = file.ServersList()
// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 2)

// 	server, _ = file.GetServer("somethere")
// 	c.Assert(server.IP, Equals, "222.222.222.111")

// }

// func (s *ConfigFileTestsSuite) TestConfigFileReplaceServerADD(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	// Check data in start
// 	sList := file.ServersList()
// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 2)

// 	server, _ := file.GetServer("somethere")

// 	server2 := server.Clone()
// 	server2.ID = "nothere"
// 	server2.IP = "222.222.222.111"

// 	find := file.ReplaceServer(server2)
// 	c.Assert(find, Equals, false)

// 	// Check result
// 	sList = file.ServersList()
// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 3)

// 	server, _ = file.GetServer("nothere")
// 	c.Assert(server.IP, Equals, "222.222.222.111")

// }

// //ReplaceScript(script *CronScript.Script) bool
// func (s *ConfigFileTestsSuite) TestConfigFileReplaceScript(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	list := *file.FileData.Scripts
// 	c.Assert(len(list), Equals, 4)

// 	script, find := file.GetScript("ls1")
// 	c.Assert(find, Equals, true)
// 	c.Assert(script.Exe, Not(Equals), "SUPER-PUPER-SCRIPT")

// 	script.SetExe("SUPER-PUPER-SCRIPT")

// 	isReplace := file.ReplaceScript(script)
// 	c.Assert(isReplace, Equals, true)
// 	list = *file.FileData.Scripts
// 	c.Assert(len(list), Equals, 4)

// 	script, find = file.GetScript("ls1")
// 	c.Assert(find, Equals, true)
// 	c.Assert(script.Exe, Equals, "SUPER-PUPER-SCRIPT")

// }

// //ReplaceScript(script *CronScript.Script) bool
// func (s *ConfigFileTestsSuite) TestConfigFileReplaceScriptADD(c *C) {
// 	//c.Skip("Not now")

// 	file := New(ReadFlags.NewTest())
// 	file.LoadFileFileFromLine(TestConfigFileLineHttp())

// 	list := *file.FileData.Scripts
// 	c.Assert(len(list), Equals, 4)

// 	scriptOld, findOld := file.GetScript("ls1")
// 	c.Assert(findOld, Equals, true)
// 	c.Assert(scriptOld.Exe, Not(Equals), "SUPER-PUPER-SCRIPT")
// 	c.Assert(scriptOld.ID, Not(Equals), "SUPER-PUPER-SCRIPT")

// 	scriptNew := scriptOld.Clone()
// 	scriptNew.ID = "SUPER-PUPER-SCRIPT"
// 	scriptNew.SetExe("SUPER-PUPER-SCRIPT")

// 	isReplace := file.ReplaceScript(scriptNew)
// 	log.Printf("%+v\n", file.Scripts)
// 	c.Assert(isReplace, Equals, false)
// 	list = *file.FileData.Scripts
// 	c.Assert(len(list), Equals, 5)

// 	script, find := file.GetScript("SUPER-PUPER-SCRIPT")
// 	c.Assert(find, Equals, true)
// 	c.Assert(script.Exe, Equals, "SUPER-PUPER-SCRIPT")
// }
