package CronScheduler

import (
	"BUtils"
	"time"
)

type Scheduler struct {
	//Debug bool

	Min   string
	Hour  string
	DayM  string
	DayW  string
	Month string

	MinInt   []int
	HourInt  []int
	DayMInt  []int
	DayWInt  []int
	MonthInt []int
}

func New() *Scheduler {
	return &Scheduler{
		Min:   "*",
		Hour:  "*",
		DayM:  "*",
		DayW:  "*",
		Month: "*",

		MinInt:   []int{},
		HourInt:  []int{},
		DayMInt:  []int{},
		DayWInt:  []int{},
		MonthInt: []int{},
	}
}

func (sched *Scheduler) Clone() *Scheduler {
	return &Scheduler{
		Min:   sched.Min,
		Hour:  sched.Hour,
		DayM:  sched.DayM,
		DayW:  sched.DayW,
		Month: sched.Month,

		MinInt:   BUtils.CopyIntsList(sched.MinInt),
		HourInt:  BUtils.CopyIntsList(sched.HourInt),
		DayMInt:  BUtils.CopyIntsList(sched.DayMInt),
		DayWInt:  BUtils.CopyIntsList(sched.DayWInt),
		MonthInt: BUtils.CopyIntsList(sched.MonthInt),
	}
}

func (sched *Scheduler) StartTime(prevTime, finishTime time.Time) bool {
	return true
}

func (sched *Scheduler) GetNextTime(t time.Time) time.Time {

	Dur, _ := time.ParseDuration("1m")
	t = t.Add(Dur)

	day, month, year := sched.getNextDay(t, 0)

	if day > int(t.Day()) || month > int(t.Month()) || year > int(t.Year()) {
		return time.Date(year, time.Month(month), day, sched.HourInt[0],
			sched.MinInt[0], 0, 0, t.Location())
	}

	addDay, hour, min := sched.getNextHourMin(t, 0)
	if addDay {
		day, month, year := sched.getNextDay(t, 1)
		return time.Date(year, time.Month(month), day, sched.HourInt[0],
			sched.MinInt[0], 0, 0, t.Location())
	}

	return time.Date(year, time.Month(month), day, hour, min, 0, 0, t.Location())
}

func (sched *Scheduler) getNextHourMin(t time.Time, add int) (bool, int, int) {

	min := int(t.Minute())
	hour := int(t.Hour())

	minI := getNextI(min, sched.MinInt)
	hourI := getNextI(hour, sched.HourInt)

	if sched.HourInt[hourI] >= hour && sched.MinInt[minI] >= min {
		return false, sched.HourInt[hourI], sched.MinInt[minI]
	}

	if sched.HourInt[hourI] < hour {
		return true, 0, 0
	}

	minI = (minI + 1) % len(sched.MinInt)

	if sched.MinInt[minI] < min {
		hourI = (hourI + 1) % len(sched.HourInt)
		minI = 0
	}

	if sched.HourInt[hourI] >= hour {
		return false, sched.HourInt[hourI], sched.MinInt[minI]
	}

	return true, 0, 0
}

func (sched *Scheduler) getNextDay(t time.Time, add int) (int, int, int) {

	D, _ := time.ParseDuration("24h")

	day := int(t.Day())
	dayW := int(t.Weekday())
	month := int(t.Month())
	year := int(t.Year())

	if sched.checkDay(day, dayW, month) {
		return day, month, year
	}

	// 28 years * 365 day + 7 (of 29/02) = 10227 It's full cyle for "day of week", "day of month" and "month"
	for i := 0; i < 10230; i++ {
		t = t.Add(D)
		day = int(t.Day())
		dayW = int(t.Weekday())
		month = int(t.Month())
		year = int(t.Year())

		if sched.checkDay(day, dayW, month) {
			return day, month, year
		}
	}

	return 0, 0, 0
}

func (sched *Scheduler) checkDay(day, dayW, month int) bool {

	if !checkList(dayW, sched.DayWInt) {
		return false
	}

	if !checkList(month, sched.MonthInt) {
		return false
	}

	if !checkList(day, sched.DayMInt) {
		return false
	}

	return true
}

func checkList(v int, list []int) bool {
	for _, d := range list {
		if d == v {
			return true
		}
	}
	return false
}

func getNextI(v int, list []int) int {

	if v <= list[0] || list[len(list)-1] < v {
		return 0
	}

	for i := 1; i < len(list); i++ {
		if list[i-1] < v && v <= list[i] {
			return i
		}
	}
	return -1
}
