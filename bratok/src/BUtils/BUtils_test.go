package BUtils

import (
	"encoding/json"
	. "gopkg.in/check.v1"
	"testing"
)

func TestNew(t *testing.T) {
	TestingT(t)
}

type NewHandBUtilsTestsSuite struct{}

var _ = Suite(&NewHandBUtilsTestsSuite{})

func __json_read(c *C, s string) map[string]interface{} {
	out := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &out)
	c.Assert(err, IsNil)
	return out
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsAnyToInt(c *C) {

	var t1 string = "11"
	t2 := &t1

	t3 := []string{`1`, `1`}
	t4 := &t3

	t5 := [][]byte{[]byte(`11`)}
	t6 := &t5

	t7 := [][]uint8{[]uint8(`11`)}
	t8 := &t7

	t9 := []byte(`11`)
	t10 := &t9

	var t11 int = 11
	t12 := &t11

	var t13 int32 = 11
	t14 := &t13

	var t15 int64 = 11
	t16 := &t15

	var t17 int8 = 11
	t18 := &t17

	list := []interface{}{
		t1, t2, t3, t4, t5, t6, t7, t8, t9,
		t10, t11, t12, t13, t14,
		t15, t16, t17, t18,
	}

	for _, v := range list {
		str := AnyToInt(v)
		c.Assert(str, Equals, 11)
	}

	var tFloat float64 = 11

	str := AnyToInt(tFloat)
	c.Assert(str, Equals, 11)

	str = AnyToInt(&tFloat)
	c.Assert(str, Equals, 11)

	str = AnyToInt(nil)
	c.Assert(str, Equals, 0)

	str = AnyToInt("")
	c.Assert(str, Equals, 0)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsAnyToIntMinMax(c *C) {

	var t string = "11"

	str := AnyToInt(t, 8, 10)
	c.Assert(str, Equals, 10)

	str = AnyToInt(t, 12, 22)
	c.Assert(str, Equals, 12)

}

func (s *NewHandBUtilsTestsSuite) TestBUtilsAnyToString(c *C) {

	var t1 string = "11"
	t2 := &t1

	t3 := []string{`1`, `1`}
	t4 := &t3

	//t5 := []uint8{`1`, `1`}
	t6 := []byte(`11`)

	var t9 int = 11
	t10 := &t9

	var t11 int32 = 11
	t12 := &t11

	var t13 int64 = 11
	t14 := &t13

	var t15 int8 = 11
	t16 := &t15

	list := []interface{}{
		t1, t2, t3, t4, t6, t9,
		t10, t11, t12, t13, t14, t15, t16,
	}

	for _, v := range list {
		str := AnyToString(v)
		c.Assert(str, Equals, `11`)
	}

	var tFloat float64 = 11

	str := AnyToString(tFloat)
	c.Assert(str, Equals, `11.000000`)

	str = AnyToString(&tFloat)
	c.Assert(str, Equals, `11.000000`)

	str = AnyToString(nil)
	c.Assert(str, Equals, ``)

}

func (s *NewHandBUtilsTestsSuite) TestBUtilsCopyStringsList(c *C) {
	a := []string{"1", "2"}
	str := CopyStringsList(a)
	c.Assert(str[1], Equals, "2")
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsCopyIntsList(c *C) {
	a := []int{11, 22}
	str := CopyIntsList(a)
	c.Assert(str[1], Equals, 22)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsSplitTimeList(c *C) {
	a := []string{
		`*/5 90 re 43 , 5433 54545,ee `,
		`*/5 9 re 44 , 54545,ee `,
	}

	str := SplitTimeList(a)
	c.Assert(str[1][3], Equals, `44`)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsAnyToStringArray(c *C) {

	var t1 string = "11"
	t2 := &t1

	t3 := []string{`11`, `12`}
	t4 := &t3

	t5 := []byte(`11`)

	var t6 int = 11
	t7 := &t6

	var t8 int32 = 11
	t9 := &t8

	var t10 int64 = 11
	t11 := &t10

	var t12 int8 = 11
	t13 := &t12

	t14 := []int{11}
	t15 := &t14

	t16 := []interface{}{11, `12`, `13`}

	t17 := []*string{&t1}
	t18 := &t17

	list := []interface{}{
		t1, t2, t3, t4, t5, t6, t7, t8, t9, t10,
		t11, t12, t13, t14, t15, t16, t17, t18,
	}

	for _, v := range list {
		str := AnyToStringArray(v)
		c.Assert(str[0], Equals, `11`)
	}

	var tFloat float64 = 11

	str := AnyToStringArray(tFloat)
	c.Assert(str[0], Equals, `11.000000`)

	str = AnyToStringArray(&tFloat)
	c.Assert(str[0], Equals, `11.000000`)

	str = AnyToStringArray(nil)
	c.Assert(len(str), Equals, 0)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsAnyTo2StringArray(c *C) {

	var t1 string = "11"
	t2 := &t1

	t3 := []string{`11`, `12`}
	t4 := &t3

	t5 := []byte(`11`)

	var t6 int = 11
	t7 := &t6

	var t8 int32 = 11
	t9 := &t8

	var t10 int64 = 11
	t11 := &t10

	var t12 int8 = 11
	t13 := &t12

	t14 := [][]string{[]string{`11`, `12`}, []string{`13`, `14`}}
	t15 := &t14

	t16 := []*[]string{&t3, &t3}
	t17 := &t16

	t18 := []int{11}
	t19 := &t18

	t20 := []interface{}{11, `12`, `13`}

	list := []interface{}{
		t1, t2, t3, t4, t5, t6, t7, t8, t9, t10,
		t11, t12, t13, t14, t15, t16, t17, t18, t19,
		t20,
	}

	for _, v := range list {
		str := AnyTo2StringArray(v)
		c.Assert(str[0][0], Equals, `11`)
	}

	var tFloat float64 = 11

	str := AnyTo2StringArray(tFloat)
	c.Assert(str[0][0], Equals, `11.000000`)

	str = AnyTo2StringArray(&tFloat)
	c.Assert(str[0][0], Equals, `11.000000`)

	str = AnyTo2StringArray(nil)
	c.Assert(len(str), Equals, 0)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsInGetPath1(c *C) {

	data := __json_read(c, `{"first":"one"}`)
	paths := "first"

	str := AnyToString(GetPath(data, paths))
	c.Assert(str, Equals, `one`)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsInGetPath2(c *C) {

	data := __json_read(c, `{"first":{"second":"one"}}`)

	str := AnyToString(GetPath(data, "/first/second"))
	c.Assert(str, Equals, `one`)

	str = AnyToString(GetPath(data, "first/second"))
	c.Assert(str, Equals, `one`)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsInGetPath3(c *C) {

	data := __json_read(c, `{"first":{"second":["one"]}}`)

	str := AnyToString(GetPath(data, "/first/second/0"))
	c.Assert(str, Equals, `one`)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsInGetPath4(c *C) {

	data := __json_read(c, `{"first":{"second":{"next":"one"}}}`)

	str := GetPath(data, "/first/second/*")
	e := str.(map[string]interface{})
	c.Assert(AnyToString(e["next"]), Equals, `one`)
}

func (s *NewHandBUtilsTestsSuite) TestBUtilsInGetPath5(c *C) {

	data := __json_read(c, `{"first":{"second":{"next":["zero","one"]}}}`)

	str := GetPath(data, "/first/second/next/*")
	e := str.([]interface{})
	c.Assert(AnyToString(e[1]), Equals, `one`)
}
