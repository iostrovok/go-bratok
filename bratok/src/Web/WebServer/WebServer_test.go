package WebServer

import (
	"Config/Config"
	"Config/ConfigHttp"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"testing"
)

func TestNewWebServer(t *testing.T) {
	TestingT(t)
}

type NewWebServerTestsSuite struct{}

var _ = Suite(&NewWebServerTestsSuite{})

func (s *NewWebServerTestsSuite) TestNewWebServerCanNewNewWebServer(c *C) {

	flags := ReadFlags.NewTest()
	cronConfig := Config.New(flags)
	config := ConfigHttp.New(flags, cronConfig)
	server, err := New(config)

	c.Assert(err, IsNil)
	c.Assert(server, NotNil)
}

func (s *NewWebServerTestsSuite) TestNewWebServerPort(c *C) {

	//init_handlers(m *pat.PatternServeMux, list map[string]interface{})
}
