package History

import (
	D "Config/Data"

	"crypto/md5"
	"fmt"
	"log"
	"time"
)

const layout = "2006-01-02T15:04:05.999999999Z07:00"

type State int

const (
	Activ State = 1 + iota
	Delete
)

type ItemHistory struct {
	ID       string    `json:"id"`
	PrevID   string    `json:"prev_id"`
	Time     string    `json:"time"`
	Script   *D.Script `json:"script"`
	Server   *D.Server `json:"server"`
	ServerId string    `json:"server_id"`
	Act      string    `json:"act"`
	IsMaster bool      `json:"is_master"`
	State    State     `json:"state"`
	Order    int       `json:"order"`
}

type History struct {
	ServerID string        `json:"-"`
	IsMaster bool          `json:"-"`
	Scripts  []ItemHistory `json:"scripts"`
	Servers  []ItemHistory `json:"servers"`
}

func New(server_id string, is_master bool) *History {

	h := &History{
		ServerID: server_id,
		IsMaster: is_master,
		Scripts:  []ItemHistory{},
		Servers:  []ItemHistory{},
	}

	return h
}

func (h *History) Empty() bool {
	return len(h.Scripts) > 0 || len(h.Servers) > 0
}

func (h *History) SetIsMaster(is_master bool) {
	h.IsMaster = is_master
}

func (h *History) SetServerID(server_id string) {
	h.ServerID = server_id
}

func (h *History) Push(sr *D.Server, sc *D.Script, act string) error {

	log.Printf("File.Push 0 file.D.History: %+v\n", h)

	if h.ServerID == "" {
		return fmt.Errorf("Empty serverID")
	}

	if sr != nil {
		it := h.emptyItemHistory(act)
		it.Server = sr

		if len(h.Servers) > 0 {
			it.PrevID = fmt.Sprintf("%s", h.Servers[len(h.Servers)-1].ID)
		}

		h.Servers = append(h.Servers, it)
	}

	if sc != nil {
		it := h.emptyItemHistory(act)
		it.Script = sc

		if len(h.Scripts) > 0 {
			it.PrevID = fmt.Sprintf("%s", h.Scripts[len(h.Scripts)-1].ID)
		}

		h.Scripts = append(h.Scripts, it)
	}

	return nil
}

func (h *History) emptyItemHistory(act string) ItemHistory {
	t := time.Now().Format(layout)
	return ItemHistory{
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(h.ServerID+t+"script"))),
		PrevID:   "",
		Time:     t,
		Act:      act,
		State:    Activ,
		IsMaster: h.IsMaster,
	}
}

func (h *History) mergeInsertScript(sc ItemHistory) {
	for j, sc1 := range h.Scripts {
		if sc.PrevID == sc1.ID {
			tail := h.Scripts[j+1:]
			head := h.Scripts[:j]

			h.Scripts = append(head, sc)
			h.Scripts = append(h.Scripts, tail...)
			return
		}
	}

	h.Scripts = append(h.Scripts, sc)
}

func (h *History) mergeScript(outsideHistory *History) {
	checkSc := map[string]bool{}
	for _, sc := range h.Scripts {
		checkSc[sc.ID] = true
	}

	for _, sc0 := range outsideHistory.Scripts {
		if _, find := checkSc[sc0.ID]; find {
			continue
		}

		notFind := true
		for j, sc1 := range h.Scripts {
			if sc0.PrevID == sc1.ID {
				tail := h.Scripts[j+1:]
				head := h.Scripts[:j]
				h.Scripts = append(head, sc0)
				h.Scripts = append(h.Scripts, tail...)
				notFind = false
				break
			}
		}

		if notFind {
			h.Scripts = append(h.Scripts, sc0)
		}
	}
}
