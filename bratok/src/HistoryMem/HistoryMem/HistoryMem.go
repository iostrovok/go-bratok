package HistoryMem

import (
	"Cron/CronInfo"
	"Cron/CronMessage"
	"Logger/Logger"
	"time"
)

type Item struct {
	StdOut   *CronMessage.Mess `json:"stdout"`
	StdErr   *CronMessage.Mess `json:"error"`
	StdRes   *CronMessage.Mess `json:"result"`
	Start    time.Time         `json:"start"`
	Finish   time.Time         `json:"finish"`
	TaskID   string            `json:"task_id"`
	ScriptID string            `json:"script_id"`
}

// It is main structure
type HistoryMem struct {
	History []*Item
}

func New() *HistoryMem {
	return &HistoryMem{
		History: []*Item{},
	}
}

func (h *HistoryMem) Add(info *CronInfo.MyCronInfo) {
	i := &Item{
		StdRes:   info.StdRes,
		StdErr:   info.StdErr,
		StdOut:   info.StdOut,
		Start:    info.Script.LastStart,
		Finish:   info.Script.LastFinish,
		TaskID:   info.Info.TaskID,
		ScriptID: info.Script.ID,
	}

	h.History = append(h.History, i)
	if len(h.History) > 2000 {
		h.History = h.History[:1999]
	}
}

func (last *Item) forPrint() map[string]interface{} {
	return map[string]interface{}{
		"stdout":    last.StdOut.ToStringLog(),
		"result":    last.StdRes.ToStringLog(),
		"error":     last.StdErr.ToStringLog(),
		"start":     last.Start.Format(Logger.DateLayOut),
		"finish":    last.Finish.Format(Logger.DateLayOut),
		"task_id":   last.TaskID,
		"script_id": last.ScriptID,
	}
}

func (h *HistoryMem) LastExe(id string) ([]map[string]interface{}, bool) {

	out := []map[string]interface{}{}

	find := false
	for i := len(h.History) - 1; i >= 0; i-- {
		if h.History[i].ScriptID == id {
			out = append(out, h.History[i].forPrint())
			find = true
		}
	}

	return out, find
}

func (h *HistoryMem) Last() (map[string]interface{}, bool) {

	if len(h.History) == 0 {
		return nil, false

	}

	last := h.History[len(h.History)-1]

	return last.forPrint(), true
}
