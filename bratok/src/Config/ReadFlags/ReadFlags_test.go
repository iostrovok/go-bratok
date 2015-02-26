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

func (s *ReadFlagsTestsSuite) Test_ReadFlags_Check_ServerID(c *C) {

	f := NewFlagsDSL()

	err := f.ReadFlags()

	c.Assert(err, NotNil)
}
