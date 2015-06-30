package Logger

import (
	"Config/Config"
	"Cron/CronMessage"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DateLayOut      string = "2006-01-02T15:04:05.999999999"
	DateLayOutInner string = "2006/01/02 15:04:05"
)

type Logger struct {
	config *Config.Config
	file   *os.File
}

func New(config *Config.Config) *Logger {

	logFile, err := os.OpenFile(config.LogFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	logger := &Logger{
		config: config,
		file:   logFile,
	}

	return logger
}

func (logger *Logger) Close() {
	logger.file.Close()
}

func (logger *Logger) ToLogAct(act string, s ...interface{}) {

	for _, t := range s {
		text, err := logger.toText(t)
		if err != nil {
			log.Panicln(err)
		}

		logger.WriteString(act + "	" + text)
	}

}

func (logger *Logger) ToLog(s ...interface{}) {

	for _, t := range s {
		text, err := logger.toText(t)
		if err != nil {
			log.Panicln(err)
		}

		logger.WriteString(text)
	}

}

func (logger *Logger) WriteString(text string) error {
	if logger.file == nil {
		return errors.New("No log file")
	}
	_, err := logger.file.WriteString(time.Now().Format(DateLayOutInner) + "	" + text + "\n")
	return err
}

func (logger *Logger) toText(t interface{}) (string, error) {

	switch t.(type) {
	case string:
		return t.(string), nil

	case CronMessage.Mess:
		m := t.(CronMessage.Mess)
		return messToLogText(&m), nil
	case *CronMessage.Mess:
		m := t.(*CronMessage.Mess)
		return messToLogText(m), nil

	case *map[string]interface{}:
		m := t.(*map[string]interface{})
		return interfaceToLogText(*m)
	case map[string]interface{}:
		m := t.(map[string]interface{})
		return interfaceToLogText(m)
	}

	return interfaceToLogText(t)
}

func interfaceToLogText(t interface{}) (string, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func messToLogText(m *CronMessage.Mess) string {

	hasBackChannel := 0
	if m.BackChannel != nil {
		hasBackChannel = 1
	}

	hasHash := 0
	if m.Hash != nil && len(m.Hash) > 0 {
		hasHash = 1
	}

	err := ""
	if m.Error != nil {
		err = m.Error.Error()
	}

	return fmt.Sprintf("%s	%s	%s.%09d	%s	has_hash: %d	has_channel: %d	%s",
		m.ID, m.Type, m.Time.Format("2006-01-02T15:04:05"), m.Time.Nanosecond(),
		err, hasHash, hasBackChannel, m.Data)
}

func (logger *Logger) SaveResult(m ...*CronMessage.Mess) error {

	text, err := interfaceToLogText(map[string]interface{}{
		"messages": m,
	})

	if err != nil {
		return err
	}

	logger.ToLog(text)

	return nil
}
