package CronTimer

import (
	"Cron/CronMessage"
	. "gopkg.in/check.v1"
	//"log"
	"sync"
	"testing"
	"time"
)

const (
	layoutTest string = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

func TestCronTimer(t *testing.T) {
	TestingT(t)
}

type CronTimerTestsSuite struct{}

var _ = Suite(&CronTimerTestsSuite{})

func (s *CronTimerTestsSuite) TestCronTimerInit(c *C) {

	var wg sync.WaitGroup

	chOut := CronMessage.Channel()

	timer := Start(2, wg, chOut)

	c.Assert(timer, NotNil)
}

func (s *CronTimerTestsSuite) TestCronTimerRead5Times(c *C) {

	var wg sync.WaitGroup

	chOut := CronMessage.Channel()

	timer := Start(1, wg, chOut)

	var wgTest sync.WaitGroup

	Count := 0
	wgTest.Add(1)
	go func() {
		defer wgTest.Done()
		for Count < 5 {
			Count++
			select {
			case m := <-chOut:
				c.Assert(m.Time.Format(layoutTest), Equals, time.Now().Format(layoutTest))
			}
		}
	}()

	wgTest.Wait()
	CronMessage.Send(timer.chIn, "", "stop", []byte{}, nil)
	wg.Wait()

	c.Assert(Count, Equals, 5)
}

func (s *CronTimerTestsSuite) TestCronTimerReset(c *C) {

	var wg sync.WaitGroup

	chOut := CronMessage.Channel()

	timer := Start(1, wg, chOut)

	Count := 0
	for Count < 5 {
		Count++
		select {
		case m := <-chOut:
			if Count == 2 {
				CronMessage.Send(timer.chIn, "", "reset", []byte("2"), nil)
			}
			c.Assert(m.Time.Format(layoutTest), Equals, time.Now().Format(layoutTest))
		}
	}

	close(timer.chIn)
	wg.Wait()

	c.Assert(Count, Equals, 5)
}
