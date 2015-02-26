package CronTask

import (
	"errors"
)

/* One Task */
type Task struct {
	id     string
	exe    string
	params []string
	mark   string
}

func NewTask(id, exe string, p ...[]string) *Task {
	t := Task{
		id:     id,
		exe:    exe,
		params: []string{},
		mark:   "added",
	}

	if len(p) > 0 {
		t.params = p[0]
	}

	return &t
}

func (t *Task) Id() string {
	return t.id
}

func (t *Task) Exe() string {
	return t.exe
}

func (t *Task) Params() []string {
	return t.params
}

func (t *Task) Mark(key ...string) string {
	if len(key) > 0 {
		t.mark = key[0]
	}
	return t.mark
}

/* List of tasks */
type TaskList struct {
	mlist map[string]*Task
}

func New() *TaskList {
	return &TaskList{
		mlist: map[string]*Task{},
	}
}

func (l *TaskList) AddTask(t *Task) error {
	if _, f := l.mlist[t.id]; f {
		return errors.New("Duplicate task id")
	}
	l.mlist[t.id] = t
	return nil
}

func (l *TaskList) Add(id, exe string, p ...[]string) error {
	if _, f := l.mlist[id]; f {
		return errors.New("Duplicate task id")
	}
	l.mlist[id] = NewTask(id, exe, p...)
	return nil
}

func (l *TaskList) Mark(id, key string) error {
	if t, f := l.mlist[id]; f {
		t.Mark(key)
		return nil
	}
	return errors.New("Key not found")
}

func (l *TaskList) GetMark(id string) (string, error) {
	if t, f := l.mlist[id]; f {
		return t.Mark(), nil
	}
	return "", errors.New("Key not found")
}

func (l *TaskList) GetMarkTaskCount(key string) int {
	count := 0
	for _, t := range l.mlist {
		if t.Mark() == key {
			count++
		}
	}
	return count
}

func (l *TaskList) Len() int {
	return len(l.mlist)
}

func (l *TaskList) Get(id string) (*Task, bool) {
	t, f := l.mlist[id]
	return t, f
}
