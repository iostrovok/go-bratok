package History

import (
	"Config/Data"

	. "gopkg.in/check.v1"
	// "log"
	"fmt"
	"math/rand"
	"testing"
)

func TestConfigHistory(t *testing.T) {
	TestingT(t)
}

type ConfigHistoryTestsSuite struct{}

var _ = Suite(&ConfigHistoryTestsSuite{})

//  HELPERS ------>

func _init_ItemHistory_slice(history *History, prevID string, length int, startIs ...int) []ItemHistory {

	list := []ItemHistory{}
	startI := 0
	if len(startIs) > 0 {
		startI = startIs[0]

	}

	for i := 0; i < length; i++ {
		item := history.emptyItemHistory("replace")
		item.PrevID = prevID
		item.IsMaster = false
		item.ID = fmt.Sprintf("%d", i+startI)

		list = append(list, item)
	}

	return list
}

func _shake_ItemHistory_slice(list []ItemHistory) []ItemHistory {

	k := len(list)
	for i := 0; i < k; i++ {
		j := rand.Intn(k)
		list[j], list[i] = list[i], list[j]
	}

	return list
}

// <----- HELPERS

func (s *ConfigHistoryTestsSuite) TestConfigHistoryNew(c *C) {

	//c.Skip("Not now")
	history := New("workstation", true)

	c.Assert(history, NotNil)

}

func (s *ConfigHistoryTestsSuite) TestConfigHistorySetServerID(c *C) {
	//c.Skip("Not now")

	history := New("workstation", true)
	c.Assert(history.ServerID, Equals, "workstation")

	history.SetServerID("new_workstation")
	c.Assert(history.ServerID, Equals, "new_workstation")

	c.Assert(history.IsMaster, Equals, true)
	history.SetIsMaster(false)
	c.Assert(history.IsMaster, Equals, false)
}

func (s *ConfigHistoryTestsSuite) TestConfigHistoryPush(c *C) {
	//c.Skip("Not now")

	sr := &Data.Server{}
	sc := &Data.Script{}

	history := New("workstation", true)
	err := history.Push(sr, sc, "add")
	err = history.Push(sr, sc, "add")
	c.Assert(err, IsNil)

	history2 := New("", false)
	err = history2.Push(sr, sc, "add")
	c.Assert(err, NotNil)
}

func (s *ConfigHistoryTestsSuite) TestConfigFindCommonPrev(c *C) {
	history := New("workstation", true)

	ih1 := _init_ItemHistory_slice(history, "sssssssssss", 5, 0)
	ih2 := _init_ItemHistory_slice(history, "sssssssssss", 5, 10)

	brI := _shake_ItemHistory_slice(ih1)
	brJ := _shake_ItemHistory_slice(ih2)

	fmt.Printf("TestConfigFindCommonPrev brI.\n")
	for i, it := range brI {
		fmt.Printf("%d. ID: '%s', IsMaster: %t\n", i, it.ID, it.IsMaster)
	}

	fmt.Printf("TestConfigFindCommonPrev brJ.\n")
	for i, it := range brJ {
		fmt.Printf("%d. ID: '%s', IsMaster: %t\n", i, it.ID, it.IsMaster)
	}

	t := ListItemHistory{}

	i, j := t.FindCommonPrev(brI, brJ)

	c.Assert(i, Equals, -1)
	c.Assert(j, Equals, -1)
}

func (s *ConfigHistoryTestsSuite) TestConfigSortItemHistory(c *C) {
	c.Skip("Not now")

	history := New("workstation", true)

	ih := []ItemHistory{}

	prevID := "sssssssssss"
	for i := 0; i < 5; i++ {
		item := history.emptyItemHistory("replace")
		item.PrevID = prevID
		item.IsMaster = false
		item.ID = fmt.Sprintf("%d", i)

		ih = append(ih, item)
	}
	ih[0].IsMaster = true
	ih[2].IsMaster = true

	fmt.Printf("Before the sorting.\n")
	for i, it := range ih {
		fmt.Printf("%d. ID: '%s', IsMaster: %t\n", i, it.ID, it.IsMaster)

	}

	history.SortItemHistory(&ih)

	fmt.Printf("After the sorting.\n")
	for i, it := range ih {
		fmt.Printf("%d. ID: '%s', IsMaster: %t\n", i, it.ID, it.IsMaster)

	}

	c.Assert(ih, NotNil)

	c.Assert(ih[0].IsMaster, Equals, false)
	c.Assert(ih[1].IsMaster, Equals, false)
	c.Assert(ih[2].IsMaster, Equals, false)
	c.Assert(ih[3].IsMaster, Equals, true)
	c.Assert(ih[4].IsMaster, Equals, true)

	c.Assert(ih[0].ID, Equals, "1")
	c.Assert(ih[1].ID, Equals, "3")
	c.Assert(ih[2].ID, Equals, "4")
	c.Assert(ih[3].ID, Equals, "0")
	c.Assert(ih[4].ID, Equals, "2")
}

func (s *ConfigHistoryTestsSuite) TestConfigSortItemHistoryBig(c *C) {

	c.Skip("Not now")

	history := New("workstation", true)

	ih := _init_ItemHistory_slice(history, "sssssssssss", 5)

	ih[0].IsMaster = true
	ih[2].IsMaster = true

	ih = _shake_ItemHistory_slice(ih)

	fmt.Printf("Before the sorting.\n")
	for i, it := range ih {
		fmt.Printf("%d. ID: '%s', IsMaster: %t\n", i, it.ID, it.IsMaster)
	}

	history.SortItemHistory(&ih)

	fmt.Printf("After the sorting.\n")
	for i, it := range ih {
		fmt.Printf("%d. ID: '%s', IsMaster: %t\n", i, it.ID, it.IsMaster)

	}

	c.Assert(ih, NotNil)

	c.Assert(ih[0].IsMaster, Equals, false)
	c.Assert(ih[1].IsMaster, Equals, false)
	c.Assert(ih[2].IsMaster, Equals, false)
	c.Assert(ih[3].IsMaster, Equals, true)
	c.Assert(ih[4].IsMaster, Equals, true)

	c.Assert(ih[0].ID, Equals, "1")
	c.Assert(ih[1].ID, Equals, "3")
	c.Assert(ih[2].ID, Equals, "4")
	c.Assert(ih[3].ID, Equals, "0")
	//c.Assert(ih[4].ID, Equals, "2")
	c.Assert(ih[4].ID, Equals, "")
}
