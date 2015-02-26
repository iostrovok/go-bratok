package CronMessage

import (
	. "gopkg.in/check.v1"
	//"log"
	"testing"
	"time"
)

const (
	layoutTest string = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

func TestCronMessage(t *testing.T) {
	TestingT(t)
}

type CronMessageTestsSuite struct{}

var _ = Suite(&CronMessageTestsSuite{})

func (s *CronMessageTestsSuite) TestCronMessageInit(c *C) {

	t := NewMess("1", "test", nil, nil, nil, nil)

	c.Assert(t, NotNil)
	c.Assert(t.Type, Equals, "test")
	c.Assert(t.ID, Equals, "1")

}

func (s *CronMessageTestsSuite) TestCronMessageSendMesst(c *C) {

	m := NewMess("1", "test", nil, nil, nil, nil)
	m.Time = nil
	ch := Channel()
	err := SendMess(ch, m)
	close(ch)
	c.Assert(err, IsNil)
}

func (s *CronMessageTestsSuite) TestCronMessageMakeChannelMake(c *C) {

	ch := Channel()

	c.Assert(ch, NotNil)
}

func (s *CronMessageTestsSuite) TestCronMessageSend(c *C) {

	ch := Channel()

	Send(ch, "1", "test", []byte{}, nil)

	t := <-ch

	c.Assert(t, NotNil)
	c.Assert(t.Type, Equals, "test")
	c.Assert(t.Error, IsNil)
}

func (s *CronMessageTestsSuite) TestCronMessageSendHash(c *C) {

	ch := Channel()

	SendHash(ch, "1", "test", map[string]interface{}{"s": 1}, nil)

	t := <-ch

	c.Assert(t, NotNil)
	c.Assert(t.Type, Equals, "test")
	c.Assert(t.Hash["s"], Equals, 1)
	c.Assert(t.Error, IsNil)
}

func (s *CronMessageTestsSuite) TestCronMessageSendTime(c *C) {

	ch := Channel()

	Send(ch, "1", "test", []byte{}, nil, time.Now())

	t := <-ch

	c.Assert(t.Time.Format(layoutTest), Equals, time.Now().Format(layoutTest))
}

func (s *CronMessageTestsSuite) TestCronMessageToLog(c *C) {

	t := NewMess("1", "test", []byte{}, nil, nil, nil)

	res := t.ToLog()

	c.Assert(res, NotNil)
	c.Assert(res["data"], Equals, "")
	c.Assert(res["error"], IsNil)
}

func (s *CronMessageTestsSuite) TestCronMessageCleanCh(c *C) {
	ch := Channel()
	t := NewMess("1", "test", nil, nil, nil, ch)
	c.Assert(t.BackChannel, NotNil)
	t.CleanCh()
	c.Assert(t.BackChannel, IsNil)
}
