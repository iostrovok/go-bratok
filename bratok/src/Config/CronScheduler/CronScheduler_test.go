package CronScheduler

import (
	. "gopkg.in/check.v1"
	"log"
	"testing"
	"time"
)

const (
	layout string = "2006-01-02 15:04:05"
)

func TestCronSchedulerScript(t *testing.T) {
	TestingT(t)
}

type CronSchedulerScriptTestsSuite struct{}

var _ = Suite(&CronSchedulerScriptTestsSuite{})

func (s *CronSchedulerScriptTestsSuite) TestCronSchedulerScriptInit(c *C) {

	sched := New()

	c.Assert(sched, NotNil)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeDigitalBad(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	_, isGood := sched.prepareTimeDigital("", 1, 12)
	// WRONG 1
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeDigital("1001", 0, 59)
	// WRONG 2
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeDigital("0", 1, 12)
	// WRONG 3
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeDigital("-1", 1, 12)
	// WRONG 4
	c.Assert(isGood, Equals, false)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeHyphen(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	t, isGood := sched.prepareTimeHyphen("1-12", 1, 59)
	c.Assert(isGood, Equals, true)

	c.Assert(len(t), Equals, 12)
	c.Assert(t[0], Equals, 1)
	c.Assert(t[10], Equals, 11)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeHyphenBad(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	_, isGood := sched.prepareTimeHyphen("-12", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeHyphen("1-101", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeHyphen("0-13", 1, 12)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeHyphen("0-1", 1, 12)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeHyphen("1-", 0, 59)
	c.Assert(isGood, Equals, false)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeStar(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	t, isGood := sched.prepareTimeStar("*", 0, 59)
	c.Assert(isGood, Equals, true)

	// Each minute
	c.Assert(len(t), Equals, 60)
	c.Assert(t[0], Equals, 0)
	c.Assert(t[56], Equals, 56)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeStarBad(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	_, isGood := sched.prepareTimeStar("**", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeStar("", 0, 59)
	c.Assert(isGood, Equals, false)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeSlash(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	t, isGood := sched.prepareTimeSlash("*/3", 0, 59)
	c.Assert(isGood, Equals, true)

	c.Assert(len(t), Equals, 20)
	c.Assert(t[0], Equals, 0)
	c.Assert(t[19], Equals, 57)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeSlashBad(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	_, isGood := sched.prepareTimeSlash("**/3", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeSlash("*/", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeSlash("*/-1", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeSlash("*/0", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTimeSlash("/", 0, 59)
	c.Assert(isGood, Equals, false)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeComplex(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()
	t, isGood := sched.prepareTime(`*/3,11,44,5-10`, 0, 59)
	c.Assert(isGood, Equals, true)

	c.Assert(len(t), Equals, 26)
	c.Assert(t[0], Equals, 0)
	c.Assert(t[19], Equals, 42)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptCheckPrepareTimeComplexBad(c *C) {
	sched := New()
	_, isGood := sched.prepareTime("*/3,11,44,5-101", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTime("*/3,11,1000,5-10", 0, 59)
	c.Assert(isGood, Equals, false)

	_, isGood = sched.prepareTime("*,/3,11,1000,5-10", 0, 59)
	c.Assert(isGood, Equals, false)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTimeBad(c *C) {
	sched := New()

	err := sched.SetTime("s", "*", "*", "*", "*")
	// Check Bad Min
	c.Assert(err, NotNil)

	err = sched.SetTime("*", "s", "*", "*", "*")
	// Check Bad Hour
	c.Assert(err, NotNil)

	err = sched.SetTime("*", "*", "s", "*", "*")
	// Check Bad DayM
	c.Assert(err, NotNil)

	err = sched.SetTime("*", "*", "*", "s", "*")
	// Check Bad Month
	c.Assert(err, NotNil)

	err = sched.SetTime("*", "*", "*", "*", "s")
	// Check Bad DayW
	c.Assert(err, NotNil)

	err = sched.SetTime("")
	// Check empty
	c.Assert(err, NotNil)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTimeEmpty(c *C) {
	sched := New()
	err := sched.SetTime()
	c.Assert(err, NotNil)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTime(c *C) {

	sched := New()
	err := sched.SetTime("*", "*", "*", "*", "*")

	// Check EMPTY TIME
	c.Assert(err, IsNil)
	c.Assert(len(sched.DayMInt), Equals, 31)
	c.Assert(len(sched.DayWInt), Equals, 7)
	c.Assert(len(sched.MinInt), Equals, 60)
	c.Assert(len(sched.MonthInt), Equals, 12)
	c.Assert(len(sched.HourInt), Equals, 24)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTimeShot(c *C) {

	sched := New()
	err := sched.SetTime("*")

	// Check EMPTY TIME
	c.Assert(err, IsNil)
	c.Assert(len(sched.DayMInt), Equals, 31)
	c.Assert(len(sched.DayWInt), Equals, 7)
	c.Assert(len(sched.MinInt), Equals, 60)
	c.Assert(len(sched.MonthInt), Equals, 12)
	c.Assert(len(sched.HourInt), Equals, 24)

}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptgetNextDay_29_Feb_2024(c *C) {

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("33", "12", "29", "2", "4")

	t := time.Date(2024, time.Month(1), 29, 23, 59, 1, 1, time.Now().Location())
	day, month, year := sched.getNextDay(t, 0)

	log.Printf("day: %d, month: %d, year: %d\n", day, month, year)

	// Check length array
	c.Assert(day, Equals, 29)
	c.Assert(month, Equals, 2)
	c.Assert(year, Equals, 2024)

	t = time.Date(2020, time.Month(1), 29, 23, 59, 1, 1, time.Now().Location())
	day, month, year = sched.getNextDay(t, 0)

	log.Printf("day: %d, month: %d, year: %d\n", day, month, year)

	// Check length array
	c.Assert(day, Equals, 29)
	c.Assert(month, Equals, 2)
	c.Assert(year, Equals, 2024)

}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptgetNextDayThisDay(c *C) {

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("*", "*", "*", "*", "*")

	t := time.Date(2024, time.Month(2), 29, 23, 59, 1, 1, time.Now().Location())
	day, month, year := sched.getNextDay(t, 0)

	log.Printf("day: %d, month: %d, year: %d\n", day, month, year)

	// Check length array
	c.Assert(day, Equals, 29)
	c.Assert(month, Equals, 2)
	c.Assert(year, Equals, 2024)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptgetNextHourMinNow(c *C) {

	//c.Skip("NO TEST NOW")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("*", "*", "*", "*", "*")

	t := time.Date(2024, time.Month(2), 29, 2, 5, 1, 1, time.Now().Location())
	add, hour, min := sched.getNextHourMin(t, 0)

	c.Assert(add, Equals, false)
	c.Assert(hour, Equals, 2)
	c.Assert(min, Equals, 5)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptgetNextHourMinNextDate(c *C) {

	////c.Skip("NO TEST NOW")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("3,45", "*", "*", "*", "*")

	// Check next hour
	t, _ := time.Parse(layout, "2024-02-29 2:59:01")
	add, hour, min := sched.getNextHourMin(t, 0)
	c.Assert(add, Equals, false)
	c.Assert(hour, Equals, 3)
	c.Assert(min, Equals, 3)

	// Check next hour
	t, _ = time.Parse(layout, "2024-02-29 5:59:55")
	add, hour, min = sched.getNextHourMin(t, 0)
	c.Assert(add, Equals, false)
	c.Assert(hour, Equals, 6)
	c.Assert(min, Equals, 3)

	// Check next day
	t, _ = time.Parse(layout, "2024-02-29 23:59:55")
	add, hour, min = sched.getNextHourMin(t, 0)
	c.Assert(add, Equals, true)
	c.Assert(hour, Equals, 0)
	c.Assert(min, Equals, 0)

	// Check next year
	t, _ = time.Parse(layout, "2024-12-31 23:59:55")
	add, hour, min = sched.getNextHourMin(t, 0)
	c.Assert(add, Equals, true)
	c.Assert(hour, Equals, 0)
	c.Assert(min, Equals, 0)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptGetNextTime(c *C) {

	//c.Skip("NO TEST NOW")

	if time.Now().Minute() == 59 {
		c.Skip("We don't test inconsistent situations now.")
	}

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("*", "*", "*", "*", "*")
	nt := sched.GetNextTime(time.Now())

	c.Assert(int(nt.Day()), Equals, int(time.Now().Day()))
	c.Assert(int(nt.Month()), Equals, int(time.Now().Month()))
	c.Assert(int(nt.Year()), Equals, int(time.Now().Year()))
	c.Assert(int(nt.Hour()), Equals, int(time.Now().Hour()))
	c.Assert(int(nt.Minute()), Equals, 1+int(time.Now().Minute()))
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptGetNextTime_29_Feb_2024(c *C) {

	//c.Skip("NO TEST NOW")
	sched := New()

	//Min, Hour, DayM, Month, DayW
	// 29 Feb 2020
	sched.SetTime("33", "12", "29", "2", "4")
	nt := sched.GetNextTime(time.Now())

	// Check 29 Feb 2020
	c.Assert(int(nt.Minute()), Equals, 33)
	c.Assert(int(nt.Hour()), Equals, 12)
	c.Assert(int(nt.Day()), Equals, 29)
	c.Assert(int(nt.Month()), Equals, 2)
	c.Assert(int(nt.Year()), Equals, 2024)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTime2(c *C) {

	c.Skip("OLD")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	err := sched.SetTime("*", "*", "*", "*", "6")
	c.Assert(err, IsNil)
	start, _ := time.Parse(layout, "2015-01-25 15:04:05")
	finish, _ := time.Parse(layout, "2015-01-25 15:04:05")

	isStart := sched.StartTime(start, finish)
	c.Assert(isStart, Equals, true)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTime3(c *C) {

	c.Skip("OLD")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	err := sched.SetTime("*/5", "*", "*", "*", "*")
	// Check returning values
	c.Assert(err, IsNil)

	// Check length array
	//c.Assert(len(sched.Time), Equals, 1)

	start, _ := time.Parse(layout, "2015-01-25 15:04:00")
	finish, _ := time.Parse(layout, "2015-01-25 15:05:01")

	isStart := sched.StartTime(start, finish)
	//c.Assert(isStart, Equals, false)
	c.Assert(isStart, Equals, true)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTime4(c *C) {

	c.Skip("OLD")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	err := sched.SetTime("*/5", "*/4", "*/3", "*/3", "*")
	// Check returning values
	c.Assert(err, IsNil)

	// Check length array
	//c.Assert(len(sched.Time), Equals, 1)

	start, _ := time.Parse(layout, "2014-01-21 15:05:00")
	finish, _ := time.Parse(layout, "2015-01-25 15:05:00")

	isStart := sched.StartTime(start, finish)
	//c.Assert(isStart, Equals, false)
	c.Assert(isStart, Equals, true)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTime5(c *C) {

	c.Skip("OLD")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("*", "*", "*", "*/5", "*/5")
	sched.SetTime("*", "*", "*/5", "*", "*/5")
	sched.SetTime("*", "*/5", "*", "*", "*/5")
	sched.SetTime("*/5", "*", "*", "*", "*/5")

	// Check length array
	//c.Assert(len(sched.Time), Equals, 4)

	start, _ := time.Parse(layout, "2014-01-21 16:16:00")
	finish, _ := time.Parse(layout, "2014-01-21 16:16:00")

	isStart := sched.StartTime(start, finish)
	c.Assert(isStart, Equals, false)
}

func (s *CronSchedulerScriptTestsSuite) TestConfigScriptAddTime6(c *C) {

	c.Skip("OLD")

	sched := New()

	//Min, Hour, DayM, Month, DayW
	sched.SetTime("*", "*/7", "*", "*", "*")
	sched.SetTime("*/7", "*", "*", "*", "*")

	start, _ := time.Parse(layout, "2014-01-21 17:17:00")
	finish, _ := time.Parse(layout, "2014-01-21 17:17:00")

	isStart := sched.StartTime(start, finish)
	c.Assert(isStart, Equals, false)
}
