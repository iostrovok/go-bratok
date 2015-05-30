package File

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

const layout = "2006-01-02T15:04:05.999999999Z07:00"

type ItemHistory struct {
	ID       string  `json:"id"`
	PrevID   string  `json:"prev_id"`
	Time     string  `json:"time"`
	Script   *Script `json:"script"`
	Server   *Server `json:"server"`
	ServerId string  `json:"server_id"`
	Act      string  `json:"act"`
}

type History struct {
	ServerID string        `json:"-"`
	Line     []ItemHistory `json:"line"`
}

func NewHistory(server_id string) *History {

	h := &History{
		ServerID: server_id,
		Line:     []ItemHistory{},
	}

	return h
}

func (h *History) SetServerID(server_id string) {
	h.ServerID = server_id
}

func (h *History) Push(sr *Server, sc *Script, act string) error {

	log.Printf("File.Push 0 file.Data.History: %+v\n", h)

	if h.ServerID == "" {
		return fmt.Errorf("Empty serverID")
	}

	t := time.Now().Format(layout)
	it := ItemHistory{
		ID:     fmt.Sprintf("%s", md5.Sum([]byte(h.ServerID+t))),
		PrevID: "",
		Time:   t,
		Script: sc,
		Server: sr,
		Act:    act,
	}

	if len(h.Line) > 0 {
		it.PrevID = fmt.Sprintf("%s", h.Line[len(h.Line)-1].ID)
	}

	h.Line = append(h.Line, it)
	return nil
}
