package StaticFiles

import (
	"Config/ConfigHttp"
	"Config/ReadFlags"
	. "gopkg.in/check.v1"
	"testing"
)

func TestNewNotFound(t *testing.T) {
	TestingT(t)
}

type NewNotFoundTestsSuite struct{}

var _ = Suite(&NewNotFoundTestsSuite{})

func (s *NewNotFoundTestsSuite) TestNewNotFoundCanNewNewNotFound(c *C) {

	config := ConfigHttp.New(ReadFlags.New())
	server, err := NewNotFound(config)

	c.Assert(err, IsNil)
	c.Assert(server, NotNil)
}
