package CronMessage

import (
	"time"
)

/*
Mess is interface for message
*/
type Mess struct {
	ID          string                 `json:"id,string"`
	Error       error                  `json:"error"`
	Data        []byte                 `json:"data,string"`
	Type        string                 `json:"type"`
	Time        *time.Time             `json:"time"`
	Hash        map[string]interface{} `json:"hash"`
	BackChannel chan Mess              `json:"-"`
}

// MakeChannel creates Channel of Messes
func Channel() chan Mess {
	return make(chan Mess, 100)
}

func checkTimeSlice(tis []time.Time) *time.Time {
	if len(tis) == 0 {
		tis = append(tis, time.Now())
	}
	return &tis[0]
}

func NewMess(id, tp string, data []byte, hash map[string]interface{}, err error, bCh chan Mess, tis ...time.Time) *Mess {
	return &Mess{
		ID:          id,
		Type:        tp,
		Error:       err,
		Data:        data,
		Time:        checkTimeSlice(tis),
		Hash:        hash,
		BackChannel: bCh,
	}
}

func (m *Mess) ToLog() map[string]interface{} {
	return map[string]interface{}{
		"error": m.Error,
		"data":  string(m.Data),
		"hash":  m.Hash,
	}
}

func (m *Mess) ToStringLog() string {
	return string(m.Data)
}

func (m *Mess) CleanCh() {
	if m.BackChannel != nil {
		m.BackChannel = nil
	}
}

func SendMess(ch chan Mess, m *Mess) error {

	if m.Time == nil {
		t := time.Now()
		m.Time = &t
	}

	ch <- *m
	return nil
}

func SendHash(ch chan Mess, id, tp string, hash map[string]interface{}, err error, tis ...time.Time) {
	ch <- Mess{
		ID:    id,
		Type:  tp,
		Error: err,
		Data:  []byte{},
		Time:  checkTimeSlice(tis),
		Hash:  hash,
	}
}

func Send(ch chan Mess, id, tp string, data []byte, err error, tis ...time.Time) {
	ch <- Mess{
		ID:    id,
		Type:  tp,
		Error: err,
		Data:  data,
		Time:  checkTimeSlice(tis),
		Hash:  map[string]interface{}{},
	}
}
