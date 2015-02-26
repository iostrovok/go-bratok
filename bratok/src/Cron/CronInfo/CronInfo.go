package CronInfo

import (
	"Config/CronScript"
	"Cron/Cron"
	"Cron/CronMessage"
)

type MyCronInfo struct {
	StdOut *CronMessage.Mess
	StdErr *CronMessage.Mess
	StdRes *CronMessage.Mess
	Info   *Cron.Info
	Script *CronScript.Script
}

// Help Gargabe Collector
func (info *MyCronInfo) Clean() {
	if info.StdRes != nil {
		info.StdRes.CleanCh()
	}
	if info.StdRes != nil {
		info.StdRes.CleanCh()
	}
	if info.StdOut != nil {
		info.StdOut.CleanCh()
	}
}
