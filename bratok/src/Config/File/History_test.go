package File

import (
	. "gopkg.in/check.v1"
	// "log"
	// "strings"
	"testing"
	// "time"
)

func TestConfigFileHistory(t *testing.T) {
	TestingT(t)
}

type ConfigFileHistoryTestsSuite struct{}

var _ = Suite(&ConfigFileHistoryTestsSuite{})

func (s *ConfigFileHistoryTestsSuite) TestConfigFileHistoryNew(c *C) {

	//c.Skip("Not now")
	history := NewHistory("workstation")

	c.Assert(history, NotNil)

}

func (s *ConfigFileHistoryTestsSuite) TestConfigFileHistorySetServerID(c *C) {
	//c.Skip("Not now")

	history := NewHistory("workstation")
	c.Assert(history.serverID, Equals, "workstation")

	history.SetServerID("new_workstation")
	c.Assert(history.serverID, Equals, "new_workstation")
}

func (s *ConfigFileHistoryTestsSuite) TestConfigFileHistoryPush(c *C) {
	//c.Skip("Not now")

	sr := &Server{}
	sc := &Script{}

	history := NewHistory("workstation")
	err := history.Push(sr, sc)
	c.Assert(err, IsNil)

	history2 := NewHistory("")
	err = history2.Push(sr, sc)
	c.Assert(err, NotNil)
}
