package File

import (
	"crypto/md5"
	"fmt"
	"time"
)

const layout = "2006-01-02T15:04:05.999999999Z07:00"

type ItemHistory struct {
	ID       string  `json:"id"`
	Time     string  `json:"time"`
	Script   *Script `json:"script"`
	Server   *Server `json:"server"`
	ServerId string  `json:"server_id"`
}

type History struct {
	serverID string        `json:"-"`
	line     []ItemHistory `json:"history"`
}

func NewHistory(server_id string) *History {

	h := &History{
		serverID: server_id,
		line:     []ItemHistory{},
	}

	return h
}

func (h *History) SetServerID(server_id string) {
	h.serverID = server_id
}

func (h *History) Push(sr *Server, sc *Script) error {

	if h.serverID == "" {
		return fmt.Errorf("Empty serverID")
	}

	t := time.Now().Format(layout)
	it := ItemHistory{
		ID:     fmt.Sprintf("%s", md5.Sum([]byte(h.serverID+t))),
		Time:   t,
		Script: sc,
		Server: sr,
	}

	h.line = append(h.line, it)
	return nil
}
