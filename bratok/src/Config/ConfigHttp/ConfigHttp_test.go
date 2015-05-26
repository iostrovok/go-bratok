package ConfigHttp

import (
	ConfigCron "Config/Config"
	"Config/File"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"testing"
)

func TestConfigHttp(t *testing.T) {
	TestingT(t)
}

type ConfigHttpTestsSuite struct{}

var _ = Suite(&ConfigHttpTestsSuite{})

func (s *ConfigHttpTestsSuite) TestConfigHttpCanNewHttp(c *C) {

	config := ConfigCron.New(ReadFlags.NewTest("somethere"))
	config.FromLine(File.TestFileLine())

	c.Assert(config.ErrorLoad, IsNil)

	configHttp := New(ReadFlags.NewTest("somethere"), config)

	c.Assert(configHttp, NotNil)
}

func (s *ConfigHttpTestsSuite) TestConfigHttpPort(c *C) {

	flag := ReadFlags.NewTest("somethere")
	flag.Port = 0

	config := ConfigCron.New(flag)
	config.FromLine(File.TestFileLine())

	c.Assert(config.ErrorLoad, IsNil)

	configHttp := New(flag, config)

	port := configHttp.Port()

	c.Assert(port, Equals, 21223)
}

// func (s *ConfigHttpTestsSuite) TestConfigHttpStaticDir(c *C) {

//	config := ConfigCron.New(ReadFlags.NewTest())
// 	port := config.StaticDir()

// 	c.Assert(port, Equals, "/tmp")
// }
