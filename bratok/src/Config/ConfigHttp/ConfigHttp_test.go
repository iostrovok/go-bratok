package ConfigHttp

import (
	ConfigCron "Config/Config"
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

	config := ConfigCron.New(ReadFlags.NewTest())
	c.Assert(config.ErrorLoad, IsNil)

	configHttp := New(ReadFlags.NewTest(), config)

	c.Assert(configHttp, NotNil)
}

func (s *ConfigHttpTestsSuite) TestConfigHttpPort(c *C) {

	config := ConfigCron.New(ReadFlags.NewTest())
	c.Assert(config.ErrorLoad, IsNil)

	configHttp := New(ReadFlags.NewTest(), config)

	port := configHttp.Port()

	c.Assert(port, Equals, 21222)
}

// func (s *ConfigHttpTestsSuite) TestConfigHttpStaticDir(c *C) {

//	config := ConfigCron.New(ReadFlags.NewTest())
// 	port := config.StaticDir()

// 	c.Assert(port, Equals, "/tmp")
// }
