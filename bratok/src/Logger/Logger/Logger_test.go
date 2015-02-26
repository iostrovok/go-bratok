package Logger

import (
	"Config/Config"
	"Config/ReadFlags"
	"Cron/CronMessage"
	"errors"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

const (
	layout string = "2006-01-02 15:04:05"
)

func TestLogger(t *testing.T) {
	TestingT(t)
}

type LoggerTestsSuite struct{}

var _ = Suite(&LoggerTestsSuite{})

func __init_test_message() *CronMessage.Mess {
	t, _ := time.Parse("2006-01-02T15:04:05.999999999", "2015-03-30T19:07:47.109782672")
	hash := map[string]interface{}{}
	return CronMessage.NewMess("id", "mtype", nil, hash, nil, nil, t)
}

func __init_test_logger() *Logger {

	cronConfig := Config.New(ReadFlags.New())
	return New(cronConfig)
}

func (s *LoggerTestsSuite) TestLoggerInit(c *C) {

	logger := __init_test_logger()
	defer logger.Close()

	c.Assert(logger, NotNil)
}

func (s *LoggerTestsSuite) TestLoggerWriteString(c *C) {

	logger := __init_test_logger()
	defer logger.Close()

	err := logger.WriteString("1")

	c.Assert(err, IsNil)
}

func (s *LoggerTestsSuite) TestLoggerToLogAct(c *C) {

	logger := __init_test_logger()
	defer logger.Close()

	logger.ToLogAct("sss", "string...string...string...string...string... string...string...string...")

	c.Assert(logger, NotNil)
}

func (s *LoggerTestsSuite) TestLoggerToLogString(c *C) {

	logger := __init_test_logger()
	defer logger.Close()

	logger.ToLog("string...string...string...string...string... string...string...string...")

	c.Assert(logger, NotNil)
}

func (s *LoggerTestsSuite) TestLoggerinterfaceToLogText(c *C) {

	hash := map[string]interface{}{
		"1": 1,
		"2": "2",
	}

	str, err := interfaceToLogText(hash)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, `{"1":1,"2":"2"}`)
}

func (s *LoggerTestsSuite) TestLoggerMessToLogText(c *C) {
	mess := __init_test_message()
	str := messToLogText(mess)
	c.Assert(str, Equals, "id	mtype	2015-03-30T19:07:47.109782672		has_hash: 0	has_channel: 0	")
}

func (s *LoggerTestsSuite) TestLoggerTexttoText(c *C) {
	logger := __init_test_logger()
	defer logger.Close()
	str, err := logger.toText("id	mtype	2015-03-30T19:07:47.109782672		has_hash: 0	has_channel: 0	")
	c.Assert(err, IsNil)
	c.Assert(str, Equals, "id	mtype	2015-03-30T19:07:47.109782672		has_hash: 0	has_channel: 0	")
}

func (s *LoggerTestsSuite) TestLoggerMapIntertoText(c *C) {
	logger := __init_test_logger()
	defer logger.Close()

	hash := map[string]interface{}{
		"1": 1,
		"2": "2",
	}

	str, err := logger.toText(hash)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, `{"1":1,"2":"2"}`)

	str, err = logger.toText(&hash)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, `{"1":1,"2":"2"}`)
}

func (s *LoggerTestsSuite) TestLoggertoTextRefMessage(c *C) {
	logger := __init_test_logger()
	defer logger.Close()

	mess := __init_test_message()

	mess.Hash = map[string]interface{}{
		"1": 1,
		"2": "2",
	}
	mess.BackChannel = CronMessage.Channel()
	mess.Error = errors.New("test")

	str, err := logger.toText(*mess)

	c.Assert(err, IsNil)
	c.Assert(str, Equals, "id	mtype	2015-03-30T19:07:47.109782672	test	has_hash: 1	has_channel: 1	")
}

func (s *LoggerTestsSuite) TestLoggertoTextMessage(c *C) {
	logger := __init_test_logger()
	defer logger.Close()

	mess := __init_test_message()

	str, err := logger.toText(mess)

	c.Assert(err, IsNil)
	c.Assert(str, Equals, "id	mtype	2015-03-30T19:07:47.109782672		has_hash: 0	has_channel: 0	")
}

func (s *LoggerTestsSuite) TestLoggerSaveResult(c *C) {
	logger := __init_test_logger()
	defer logger.Close()

	mess1 := __init_test_message()
	mess2 := __init_test_message()

	err := logger.SaveResult(mess1, mess2)

	c.Assert(err, IsNil)
}

func (s *LoggerTestsSuite) TestLoggerStrRefToText(c *C) {
	logger := __init_test_logger()
	defer logger.Close()

	send := "asdsadsad"

	str, err := logger.toText(&send)

	c.Assert(err, IsNil)
	c.Assert(str, Equals, `"asdsadsad"`)
}
