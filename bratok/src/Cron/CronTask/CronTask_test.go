package CronTask

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestCronTask(t *testing.T) {
	TestingT(t)
}

type CronTaskTestsSuite struct{}

var _ = Suite(&CronTaskTestsSuite{})

func (s *CronTaskTestsSuite) TestCronTaskNew(c *C) {

	ch := NewTask("12", "ls")

	c.Assert(ch, NotNil)
}

func (s *CronTaskTestsSuite) TestCronTaskGetId(c *C) {

	ch := NewTask("12", "ls")

	c.Assert(ch.Id(), Equals, "12")
	c.Assert(ch.Exe(), Equals, "ls")
}

func (s *CronTaskTestsSuite) TestCronTaskMark(c *C) {

	t := NewTask("12", "ls")

	t.Mark("done")

	c.Assert(t.Mark(), Equals, "done")
}

func (s *CronTaskTestsSuite) TestCronTaskNewMultyParams(c *C) {

	t := NewTask("12", "ls", []string{"asdasd"})
	c.Assert(t, NotNil)

	p := t.Params()

	c.Assert(p[0], Equals, "asdasd")
}

func (s *CronTaskTestsSuite) TestCronTaskListNew(c *C) {

	tl := New()

	c.Assert(tl, NotNil)
}

func (s *CronTaskTestsSuite) TestCronTaskListAdd(c *C) {

	tl := New()
	err := tl.Add("12", "ls")

	c.Assert(err, IsNil)
}

func (s *CronTaskTestsSuite) TestCronTaskLen(c *C) {

	tl := New()
	tl.Add("12", "ls")

	c.Assert(tl.Len(), Equals, 1)
}

func (s *CronTaskTestsSuite) TestCronTaskListAddTask(c *C) {

	tl := New()
	t := NewTask("12", "ls", []string{"asdasd"})

	err := tl.AddTask(t)

	c.Assert(err, IsNil)

	task, f := tl.Get("12")
	c.Assert(task, NotNil)
	c.Assert(f, Equals, true)
}

func (s *CronTaskTestsSuite) TestCronTaskListAddTaskBad(c *C) {

	tl := New()
	t := NewTask("12", "ls", []string{"asdasd"})

	tl.AddTask(t)
	err := tl.AddTask(t)

	c.Assert(err, NotNil)
}

func (s *CronTaskTestsSuite) TestCronTaskListAddBad(c *C) {

	tl := New()
	tl.Add("12", "ls")
	err := tl.Add("12", "ls")

	c.Assert(err, NotNil)
}

func (s *CronTaskTestsSuite) TestCronTaskListMark(c *C) {

	tl := New()
	tl.Add("12", "ls")
	err := tl.Mark("12", "new_mark")
	c.Assert(err, IsNil)

	markGood, errGoodFind := tl.GetMark("12")
	c.Assert(errGoodFind, IsNil)
	c.Assert(markGood, Equals, "new_mark")

	markBad, errBadFind := tl.GetMark("13")
	c.Assert(errBadFind, NotNil)
	c.Assert(markBad, Equals, "")
}

func (s *CronTaskTestsSuite) TestCronTaskListMarkCount(c *C) {

	tl := New()

	tl.Add("12", "ls")
	tl.Mark("12", "new_mark")

	tl.Add("13", "ls")
	tl.Mark("13", "new_mark")

	count := tl.GetMarkTaskCount("new_mark")
	c.Assert(count, Equals, 2)
}

func (s *CronTaskTestsSuite) TestCronTaskListMarkBad(c *C) {

	tl := New()
	tl.Add("12", "ls")
	err := tl.Mark("13", "new_key")

	c.Assert(err, NotNil)
}
