package Server

import (
	"Config/Config"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"testing"
)

func TestServer(t *testing.T) {
	TestingT(t)
}

type ServerTestsSuite struct{}

var _ = Suite(&ServerTestsSuite{})

func (s *ServerTestsSuite) TestServerInit(c *C) {

	server := NewServer()

	c.Assert(server, NotNil)
}

func (s *ServerTestsSuite) TestServerCanStart(c *C) {

	config := Config.New(ReadFlags.New())
	server := NewServer()

	err := server.Start(config)

	c.Assert(err, IsNil)
}
