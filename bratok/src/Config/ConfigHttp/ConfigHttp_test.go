package ConfigHttp

import (
	cronConf "Config/Config"
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

	cronConfig := cronConf.New(ReadFlags.New())
	config := New(ReadFlags.New(), cronConfig)

	c.Assert(config, NotNil)
}

func (s *ConfigHttpTestsSuite) TestConfigHttpPort(c *C) {

	cronConfig := cronConf.New(ReadFlags.New())
	config := New(ReadFlags.New(), cronConfig)

	port := config.Port()

	c.Assert(port, Equals, 21222)
}

// func (s *ConfigHttpTestsSuite) TestConfigHttpStaticDir(c *C) {

// 	config := New(ReadFlags.New())
// 	port := config.StaticDir()

// 	c.Assert(port, Equals, "/tmp")
// }
