package ReadFlags

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestReadFlags(t *testing.T) {
	TestingT(t)
}

type ReadFlagsTestsSuite struct{}

var _ = Suite(&ReadFlagsTestsSuite{})

func (s *ReadFlagsTestsSuite) TestReadFlagsInit(c *C) {
	f := New()
	c.Assert(f, NotNil)
}

func (s *ReadFlagsTestsSuite) TestReadFlagsRead(c *C) {
	f := New()
	err := f.Read()
	c.Assert(err, NotNil)
}

func (s *ReadFlagsTestsSuite) TestReadFlagsCheck(c *C) {

	f := New()

	err := f.Check()
	c.Assert(err, NotNil)

	f.ServerID = "ssss"
	err = f.Check()
	c.Assert(err, IsNil)
}
