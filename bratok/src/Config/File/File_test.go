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

	server, find := file.FindServer("workstation_asdsdfdasfsdfa")
	c.Assert(server, IsNil)
	c.Assert(find, Equals, false)

	server, find = file.FindServer("workstation")
	c.Assert(server, NotNil)
	c.Assert(find, Equals, true)

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

func (s *ConfigFileTestsSuite) TestConfigFileUpdateId(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	oldId := file.ConfigID()
	file.UpdateId()
	newId := file.ConfigID()

	c.Assert(oldId, Not(Equals), newId)

	if newId < oldId {
		c.Fatalf("Bad old-new id: %d - %d\n", oldId, newId)
	}
}

func (s *ConfigFileTestsSuite) TestConfigFileLoadHTTPLine(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")

	err := file.LoadHTTPLine(TestFileLineHttp())
	c.Assert(err, IsNil)

	c.Assert(file.ServerID(), Equals, "workstation")
	c.Assert(file.ConfigID(), Equals, int64(12312312))
	c.Assert(file.ScriptLogDir(), Equals, "workstation_DIR")
	c.Assert(file.StaticFilesDir(), Equals, "workstation_STATIC-DIR")
	c.Assert(file.ScriptLogFile(), Equals, "workstation_FILE")

}

func (s *ConfigFileTestsSuite) TestConfigFileLoadHTTPLNoServer(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")

	err := file.LoadHTTPLine(TestFileLineHttpNoServer())
	c.Assert(err, NotNil)

}

func (s *ConfigFileTestsSuite) TestConfigFile_checkDefault(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.Data = Data{}

	err := file._checkDefault()
	c.Assert(err, IsNil)

	c.Assert(file.Data.Servers, NotNil)
	c.Assert(len(file.Data.Servers), Equals, 0)

	c.Assert(file.Data.Scripts, NotNil)
	c.Assert(len(file.Data.Scripts), Equals, 0)

	c.Assert(file.Data.History, NotNil)
}

func (s *ConfigFileTestsSuite) TestConfigFileList(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	servers := file.ListServer()
	c.Assert(len(servers), Equals, 2)

	scripts := file.ListScript()
	c.Assert(len(scripts), Equals, 4)
}

func (s *ConfigFileTestsSuite) TestConfigFileSetServer(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	// Check data in start
	c.Assert(file.Data.Servers, NotNil)
	c.Assert(len(file.Data.Servers), Equals, 2)

	server, _ := file.FindServer("somethere")

	// Step 1 replace
	server2 := server.Clone()
	server2.IP = "222.222.222.111"
	find := file.SetServer(server2)
	c.Assert(find, Equals, true)

	// Check result
	c.Assert(file.Data.Servers, NotNil)
	c.Assert(len(file.Data.Servers), Equals, 2)
	server, _ = file.FindServer("somethere")
	c.Assert(server.IP, Equals, "222.222.222.111")

	// Step 2 insert
	server3 := server.Clone()
	server3.ID = "new_script_id"
	find = file.SetServer(server3)
	c.Assert(find, Equals, false)

	// Check result
	c.Assert(file.Data.Servers, NotNil)
	c.Assert(len(file.Data.Servers), Equals, 3)
	server, _ = file.FindServer("somethere")
	c.Assert(server.IP, Equals, "222.222.222.111")

}

//SetScript(script *CronScript.Script) bool
func (s *ConfigFileTestsSuite) TestConfigFileSetScript(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	c.Assert(len(file.Data.Scripts), Equals, 4)

	script, find := file.FindScript("ls1")
	c.Assert(find, Equals, true)
	c.Assert(script.Exe, Not(Equals), "SUPER-PUPER-SCRIPT")

	script.Exe = "SUPER-PUPER-SCRIPT"

	isReplace := file.SetScript(script)
	c.Assert(isReplace, Equals, true)
	c.Assert(len(file.Data.Scripts), Equals, 4)

	script, find = file.FindScript("ls1")
	c.Assert(find, Equals, true)
	c.Assert(script.Exe, Equals, "SUPER-PUPER-SCRIPT")

}

// func (s *ConfigFileTestsSuite) TestConfigFileScriptLogFile(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")

// 	out := file.ScriptLogFile("my.log")
// 	c.Assert(out, Equals, "/tmp/my.log")
// }

// func (s *ConfigFileTestsSuite) TestConfigFileCheckInitNew(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	err := file.InitNew(map[string]interface{}{})

// 	c.Assert(err, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileAddCronScript(c *C) {
// 	//c.Skip("Not now")

// 	script := CronScript.New("1", "ls", "-a", "-r")
// 	file := New("workstation", "configFile")

// 	err := file.AddCronScript(script)
// 	c.Assert(err, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetScript(c *C) {
// 	//c.Skip("Not now")

// 	script1 := CronScript.New("1", "ls", "-a", "-r")
// 	script2 := CronScript.New("2", "ls", "-a", "-r")
// 	file := New("workstation", "configFile")

// 	file.AddCronScript(script1)
// 	file.AddCronScript(script2)

// 	script, find := file.FindScript("1")

// 	c.Assert(find, Equals, true)
// 	c.Assert(script, NotNil)

// 	script, find = file.FindScript("34")

// 	c.Assert(find, Equals, false)
// 	c.Assert(script, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileStartNow(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
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

// 	file := New("workstation", "configFile")

// 	err := file._parseFileData(TestConfigFileLine())

// 	c.Assert(err, IsNil)
// 	c.Assert(len(*file.FileData.Scripts), Equals, 4)
// 	c.Assert(len(*file.FileData.Servers), Equals, 2)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileLoadFileData(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")

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

// 	file := New("workstation", "configFile")

// 	file._parseFileData(TestConfigFileLine())
// 	err := file._loadFileData()

// 	c.Assert(err, IsNil)

// 	c.Assert(file.ScriptLogDir(), Equals, "DIR")
// 	c.Assert(file.ScriptLogDir("DIR-NEW/"), Equals, "DIR-NEW/")

// 	c.Assert(file.ScriptLogFile(), Equals, "DIR-NEW/FILE")
// 	c.Assert(file.ScriptLogFile("FILE-NEW"), Equals, "DIR-NEW/FILE-NEW")
// }

// func (s *ConfigFileTestsSuite) TestConfigFileGetHttpData(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

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

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

// 	p := []*Server{}
// 	file.FileData.Servers = &p
// 	server, err := file.GetHttpData()

// 	c.Assert(err, NotNil)
// 	c.Assert(server, IsNil)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileFindServer(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())
// 	server, find := file.FindServer("somethere")

// 	c.Assert(find, Equals, true)
// 	c.Assert(server, NotNil)
// 	c.Assert(server.Host, Equals, "wks-l")
// 	c.Assert(server.Port, Equals, 21222)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileScriptsList_1(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())
// 	scripts := file.ScriptsList()

// 	c.Assert(len(scripts), Equals, 4)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileScriptsList_2(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

// 	scripts := file.ScriptsList("somethere")
// 	c.Assert(len(scripts), Equals, 1)

// 	scripts = file.ScriptsList("bad-server-id")
// 	c.Assert(len(scripts), Equals, 0)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileUpdateId(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

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

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

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

// 	file := New("workstation", "configFile")
// 	file._parseFileData(TestConfigFileLine())

// 	sList := file.ServersList()

// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 2)
// }

// func (s *ConfigFileTestsSuite) TestConfigFileSetServerADD(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

// 	// Check data in start
// 	sList := file.ServersList()
// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 2)

// 	server, _ := file.FindServer("somethere")

// 	server2 := server.Clone()
// 	server2.ID = "nothere"
// 	server2.IP = "222.222.222.111"

// 	find := file.SetServer(server2)
// 	c.Assert(find, Equals, false)

// 	// Check result
// 	sList = file.ServersList()
// 	c.Assert(sList, NotNil)
// 	c.Assert(len(sList), Equals, 3)

// 	server, _ = file.FindServer("nothere")
// 	c.Assert(server.IP, Equals, "222.222.222.111")

// }

// //SetScript(script *CronScript.Script) bool
// func (s *ConfigFileTestsSuite) TestConfigFileSetScriptADD(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

// 	list := *file.FileData.Scripts
// 	c.Assert(len(list), Equals, 4)

// 	scriptOld, findOld := file.FindScript("ls1")
// 	c.Assert(findOld, Equals, true)
// 	c.Assert(scriptOld.Exe, Not(Equals), "SUPER-PUPER-SCRIPT")
// 	c.Assert(scriptOld.ID, Not(Equals), "SUPER-PUPER-SCRIPT")

// 	scriptNew := scriptOld.Clone()
// 	scriptNew.ID = "SUPER-PUPER-SCRIPT"
// 	scriptNew.SetExe("SUPER-PUPER-SCRIPT")

// 	isReplace := file.SetScript(scriptNew)
// 	log.Printf("%+v\n", file.Scripts)
// 	c.Assert(isReplace, Equals, false)
// 	list = *file.FileData.Scripts
// 	c.Assert(len(list), Equals, 5)

// 	script, find := file.FindScript("SUPER-PUPER-SCRIPT")
// 	c.Assert(find, Equals, true)
// 	c.Assert(script.Exe, Equals, "SUPER-PUPER-SCRIPT")
// }
