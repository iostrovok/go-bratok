package Common

import (
	"errors"
	. "gopkg.in/check.v1"
	"testing"
)

func TestNew(t *testing.T) {
	TestingT(t)
}

type NewHandCommonTestsSuite struct{}

var _ = Suite(&NewHandCommonTestsSuite{})

func (s *NewHandCommonTestsSuite) TestHandCommonInit(c *C) {

	res := NewResultEmpty()

	c.Assert(res, NotNil)
	c.Assert(res.Error, Equals, ``)
	c.Assert(res.Data, Equals, nil)
	c.Assert(res.Result, Equals, 1)
}

func (s *NewHandCommonTestsSuite) TestHandCommonInit2(c *C) {

	res := NewResult(nil, map[string]interface{}{"super-puper": 1})

	c.Assert(res, NotNil)
	c.Assert(res.Error, Equals, ``)
	c.Assert(res.Data, NotNil)
	c.Assert(res.Result, Equals, 1)
}

func (s *NewHandCommonTestsSuite) TestHandCommonInit3(c *C) {

	res := NewResult(errors.New("Testerror"))

	c.Assert(res, NotNil)
	c.Assert(res.Error, Equals, `Testerror`)
	c.Assert(res.Data, IsNil)
	c.Assert(res.Result, Equals, 0)
}

func (s *NewHandCommonTestsSuite) TestHandCommonJsonResult(c *C) {

	res := JsonResult(nil, map[string]interface{}{"super-puper": 1})

	c.Assert(res, NotNil)
	c.Assert(string(res), Equals, `{"error":"","data":{"super-puper":1},"result":1}`)
}

func (s *NewHandCommonTestsSuite) TestHandCommonJsonResultError(c *C) {

	res := JsonResult(errors.New("test error"), map[string]interface{}{"super-puper": 1})

	c.Assert(res, NotNil)
	c.Assert(string(res), Equals, `{"error":"test error","data":{"super-puper":1},"result":0}`)
}

func (s *NewHandCommonTestsSuite) TestHandCommonToJsonStr(c *C) {
	str := NewResult(nil, map[string]interface{}{"super-puper": 1}).ToJsonStr()
	c.Assert(str, Equals, `{"error":"","data":{"super-puper":1},"result":1}`)
}
