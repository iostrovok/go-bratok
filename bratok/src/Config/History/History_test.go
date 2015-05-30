package History

import (
	"Config/Data"

	. "gopkg.in/check.v1"
	// "log"
	// "strings"
	"testing"
	// "time"
)

func TestConfigHistory(t *testing.T) {
	TestingT(t)
}

type ConfigHistoryTestsSuite struct{}

var _ = Suite(&ConfigHistoryTestsSuite{})

func (s *ConfigHistoryTestsSuite) TestConfigHistoryNew(c *C) {

	//c.Skip("Not now")
	history := New("workstation")

	c.Assert(history, NotNil)

}

func (s *ConfigHistoryTestsSuite) TestConfigHistorySetServerID(c *C) {
	//c.Skip("Not now")

	history := New("workstation")
	c.Assert(history.ServerID, Equals, "workstation")

	history.SetServerID("new_workstation")
	c.Assert(history.ServerID, Equals, "new_workstation")
}

func (s *ConfigHistoryTestsSuite) TestConfigHistoryPush(c *C) {
	//c.Skip("Not now")

	sr := &Data.Server{}
	sc := &Data.Script{}

	history := New("workstation")
	err := history.Push(sr, sc, "add")
	err = history.Push(sr, sc, "add")
	c.Assert(err, IsNil)

	history2 := New("")
	err = history2.Push(sr, sc, "add")
	c.Assert(err, NotNil)
}
