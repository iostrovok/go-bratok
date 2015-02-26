package CronScheduler

import (
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"regexp"
	"sort"
)

var (
	findDigitalReg      = regexp.MustCompile(`([0-9]+)`)
	findOnlyDigitalReg  = regexp.MustCompile(`^[0-9]+$`)
	findСommaReg        = regexp.MustCompile(`,`)
	findHyphenReg       = regexp.MustCompile(`-`)
	findSlashReg        = regexp.MustCompile(`/`)
	findSlashDigitalReg = regexp.MustCompile(`^\*/[0-9]+$`)

	cleanTimeLine = regexp.MustCompile(`[^\*/0-9]+`)
)

func CleanTimeLinePoint(a string) string {
	return cleanTimeLine.ReplaceAllString(a, "")
}

// TODO Add Jan...Dec & Sun...Sat
// Minute   Hour   Day of Month       Month          Day of Week
// (0-59)  (0-23)     (1-31)    (1-12 or Jan-Dec)  (0-6 or Sun-Sat)
//     Only numbers PLEASE!
func (sched *Scheduler) SetTime(mIn ...string) error {

	if len(mIn) == 0 {
		return fmt.Errorf("Empty any time")
	}
	//Min, Hour, DayM, Month, DayW
	m := []string{"*", "*", "*", "*", "*"}

	if len(mIn) > len(m) {
		return fmt.Errorf("Bad time data")
	}

	isEmptyTime := true
	for i := 0; i < len(mIn); i++ {
		if mIn[i] != "" {
			isEmptyTime = false
			m[i] = mIn[i]
		}
	}

	if isEmptyTime {
		return fmt.Errorf("Time is empty")
	}

	sched.Min = m[0]
	sched.Hour = m[1]
	sched.DayM = m[2]
	sched.Month = m[3]
	sched.DayW = m[4]

	//Min, Hour, DayM, Month, DayW
	listMin, isGoodMin := sched.prepareTime(m[0], 0, 59)
	if !isGoodMin {
		return fmt.Errorf("Bad minuts %sched. 0-59", m[0])
	}
	sched.MinInt = listMin

	listHours, isGoodHours := sched.prepareTime(m[1], 0, 23)
	if !isGoodHours {
		return fmt.Errorf("Bad Hours %sched. 0-23", m[1])
	}
	sched.HourInt = listHours

	listDays, isGoodDays := sched.prepareTime(m[2], 1, 31)
	if !isGoodDays {
		return fmt.Errorf("Bad day of month %sched. 1-31", m[2])
	}
	sched.DayMInt = listDays

	listMonths, isGoodMonths := sched.prepareTime(m[3], 1, 12)
	if !isGoodMonths {
		return fmt.Errorf("Bad month %sched. 1-12", m[3])
	}
	sched.MonthInt = listMonths

	listDWs, isGoodDWs := sched.prepareTime(m[4], 0, 6)
	if !isGoodDWs {
		return fmt.Errorf("Bad day of week %sched. 0-6, Mn-Su", m[4])
	}
	sched.DayWInt = listDWs

	return nil
}

func (sched *Scheduler) prepareTime(m string, min, max int) ([]int, bool) {

	list, isGood := sched._prepareTime(m, min, max)
	if !isGood {
		return nil, false
	}

	sort.Ints(list)
	out := []int{list[0]}
	last := list[0]
	for _, v := range list {
		if last != v {
			out = append(out, v)
			last = v
		}
	}
	return out, true
}

func (sched *Scheduler) _prepareTime(m string, min, max int) ([]int, bool) {

	if findСommaReg.FindString(m) != "" {
		list := findСommaReg.Split(m, -1)
		out := []int{}
		for _, m := range list {

			t, f := sched.prepareTime(m, min, max)
			if f {
				out = append(out, t...)
			} else {
				return []int{}, false
			}
		}
		return out, true
	}

	if out, f := sched.prepareTimeStar(m, min, max); f {
		return out, f
	}
	if out, f := sched.prepareTimeSlash(m, min, max); f {
		return out, f
	}
	if out, f := sched.prepareTimeHyphen(m, min, max); f {
		return out, f
	}
	if out, f := sched.prepareTimeDigital(m, min, max); f {
		return out, f
	}
	return []int{}, false
}

func (sched *Scheduler) prepareTimeStar(m string, min, max int) ([]int, bool) {

	out := []int{}
	f := false

	if m == "*" {
		f = true
		for i := min; i <= max; i++ {
			out = append(out, i)
		}
	}

	return out, f
}

func (sched *Scheduler) prepareTimeSlash(m string, min, max int) ([]int, bool) {

	if findSlashReg.FindString(m) == "" {
		return []int{}, false
	}

	if findSlashDigitalReg.FindString(m) == "" {
		return []int{}, false
	}

	t := findDigitalReg.FindAllString(m, -1)
	k := iutils.AnyToInt(t[0])

	if k == 0 {
		return []int{}, false
	}

	out := []int{}
	for i := min; i <= max; i++ {
		if i%k == 0 {
			out = append(out, i)
		}
	}
	return out, true
}

func (sched *Scheduler) prepareTimeHyphen(m string, min, max int) ([]int, bool) {

	if findHyphenReg.FindString(m) == "" {
		return []int{}, false
	}

	t := findDigitalReg.FindAllString(m, -1)
	if len(t) < 2 {
		return []int{}, false
	}

	from := iutils.AnyToInt(t[0])
	to := iutils.AnyToInt(t[1])

	if from < min || from > to || max < to {
		return []int{}, false
	}

	out := []int{}
	for i := from; i <= to; i++ {
		out = append(out, i)
	}

	return out, true
}

func (sched *Scheduler) prepareTimeDigital(m string, min, max int) ([]int, bool) {

	if findOnlyDigitalReg.FindString(m) == "" {
		return []int{}, false
	}

	i := iutils.AnyToInt(m)
	if i < min || max < i {
		return []int{}, false
	}

	return []int{i}, true
}
