package File

import (
	D "Config/Data"
	. "gopkg.in/check.v1"
	"testing"
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

// func (s *ConfigFileTestsSuite) TestConfigFileSetAutoConfig(c *C) {
// 	//c.Skip("Not now")

// 	file := New("workstation", "configFile")
// 	file.FromLine(TestFileLine())

// 	err := file.SetAutoConfig(true, "line")
// 	c.Assert(err, IsNil)

// 	err = file.SetAutoConfig(false)
// 	c.Assert(err, IsNil)

// 	err = file.SetAutoConfig(false, "line")
// 	c.Assert(err, IsNil)

// 	err = file.SetAutoConfig(true, "")
// 	c.Assert(err, NotNil)

// 	err = file.SetAutoConfig(true)
// 	c.Assert(err, NotNil)
// }

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
	c.Assert(file.StaticFilesDir(), Equals, "STATIC-DIR")
	c.Assert(file.LogFile(), Equals, "FILE")
}

func (s *ConfigFileTestsSuite) TestConfigFileDirsB(c *C) {
	//c.Skip("Not now")

	file := New("workstation", "configFile")
	file.FromLine(TestFileLine())

	c.Assert(file.ServerID(), Equals, "workstation")
	c.Assert(file.ConfigID(), Equals, int64(12312312))
	c.Assert(file.StaticFilesDir(), Equals, "workstation_STATIC-DIR")
	c.Assert(file.LogFile(), Equals, "workstation_FILE")
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
	c.Assert(file.StaticFilesDir(), Equals, "workstation_STATIC-DIR")
	c.Assert(file.LogFile(), Equals, "workstation_FILE")

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
	file.Data = D.EmptyData()

	err := file._checkDefault()
	c.Assert(err, IsNil)

	c.Assert(file.Data.Servers, NotNil)
	c.Assert(len(file.Data.Servers), Equals, 0)

	c.Assert(file.Data.Scripts, NotNil)
	c.Assert(len(file.Data.Scripts), Equals, 0)

	c.Assert(file.History, NotNil)
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
