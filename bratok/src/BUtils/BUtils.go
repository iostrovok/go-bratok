package BUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// construct a regexp to extract values:
var (
	splitTimeListRegExp = regexp.MustCompile(`[ ,\.]+`)
	NonDigitalRE        = regexp.MustCompile(`[^-0-9,\.]+`)
)

func CopyStringsList(a []string) []string {
	out := make([]string, len(a))
	for i, v := range a {
		out[i] = v
	}
	return out
}

func CopyIntsList(a []int) []int {
	out := make([]int, len(a))
	for i, v := range a {
		out[i] = v
	}
	return out
}

func ReadJson(s string) (map[string]interface{}, error) {

	out := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &out)

	return out, err
}

func SplitTimeList(list []string) [][]string {
	out := make([][]string, len(list))
	for i := range list {
		out[i] = splitTimeListRegExp.Split(list[i], -1)
	}
	return out
}

func AnyToString(s interface{}) string {
	switch v := s.(type) {
	case string:
		return s.(string)
	case *string:
		return *s.(*string)

	case []string:
		return strings.Join(s.([]string), "")
	case *[]string:
		return strings.Join(*s.(*[]string), "")

	case []uint8:
		raw := s.([]uint8)
		return string(raw)

	case float64:
		return fmt.Sprintf("%f", s.(float64))
	case *float64:
		return fmt.Sprintf("%f", *s.(*float64))

	case int:
		return strconv.Itoa(s.(int))
	case *int:
		return strconv.Itoa(*s.(*int))

	case int8:
		return strconv.Itoa(int(s.(int8)))
	case *int8:
		return strconv.Itoa(int(*s.(*int8)))

	case int32:
		return strconv.FormatUint(uint64(s.(int32)), 10)
	case *int32:
		return strconv.FormatUint(uint64(*s.(*int32)), 10)

	case int64:
		return strconv.FormatUint(uint64(s.(int64)), 10)
	case *int64:
		return strconv.FormatUint(uint64(*s.(*int64)), 10)

	case nil:
		return ""

	default:
		log.Fatalf("AnyToString. unknown type '%t' => '%s'\v", s, v)

	}

	return ""
}

func GetPath(data interface{}, path string) interface{} {
	paths := strings.Split(path, "/")
	if paths[0] == "" {
		paths = paths[1:]
	}
	return inGetPath(data, paths)
}

func inGetPath(data interface{}, paths []string) interface{} {

	if len(paths) == 0 {
		return data
	}

	path := paths[0]
	paths = paths[1:]
	switch data.(type) {
	case map[string]interface{}:
		if path == "*" {
			out := map[string]interface{}{}
			for k, v := range data.(map[string]interface{}) {
				out[k] = inGetPath(v, paths)
			}
			return out
		}
		if val, find := data.(map[string]interface{})[path]; find {
			return inGetPath(val, paths)
		}
	case []interface{}:
		d := data.([]interface{})
		if path == "*" {
			out := make([]interface{}, len(d))
			for k, v := range d {
				out[k] = inGetPath(v, paths)
			}
			return out
		}

		i := AnyToInt(path)
		if i < len(d) {
			return inGetPath(d[i], paths)
		}
	default:
		return data
	}
	return nil
}

func AnyToInt(s interface{}, minmax ...int) int {
	i := _AnyToInt(s)
	if len(minmax) == 0 {
		return i
	}
	if len(minmax) > 0 && i < minmax[0] {
		return minmax[0]
	}
	if len(minmax) > 1 && i > minmax[1] {
		return minmax[1]
	}
	return i
}

func _AnyToInt(s interface{}) int {

	if s == nil {
		return 0
	}

	switch v := s.(type) {
	case *string:
		return _AnyToInt(*s.(*string))
	case string:
		st := NonDigitalRE.ReplaceAllString(s.(string), "")
		if st == "" {
			return 0
		}
		i, err := strconv.Atoi(st)
		if err != nil {
			log.Println(err)
			return 0
		}
		return int(i)

	case []byte: // && []uint8
		return _AnyToInt(string(s.([]byte)))
	case *[]byte: // && *[]uint8
		return _AnyToInt(string(*s.(*[]byte)))

	case [][]byte:
		return _AnyToInt(bytes.Join(s.([][]byte), []byte{}))
	case *[][]byte:
		return _AnyToInt(bytes.Join(*s.(*[][]byte), []byte{}))

	case []string:
		return _AnyToInt(strings.Join(s.([]string), ""))
	case *[]string:
		return _AnyToInt(strings.Join(*s.(*[]string), ""))

	case int:
		return int(s.(int))
	case int8:
		return int(s.(int8))
	case int32:
		return int(s.(int32))
	case int64:
		return int(s.(int64))
	case *int32:
		return int(*s.(*int32))
	case *int64:
		return int(*s.(*int64))
	case *int:
		return int(*s.(*int))
	case *int8:
		return int(*s.(*int8))

	case float64:
		return int(s.(float64))
	case *float64:
		return int(*s.(*float64))

	default:
		log.Fatalf("AnyToInt. unknown type %s\v", v)
	}
	return 0
}

func AnyTo2StringArray(s interface{}) [][]string {

	switch s.(type) {
	case [][]string:
		return s.([][]string)
	case *[][]string:
		return *s.(*[][]string)
	case []*[]string:
		out := [][]string{}
		list := s.([]*[]string)
		for _, a := range list {
			out = append(out, *a)
		}
		return out

	case *[]*[]string:
		out := [][]string{}
		list := *s.(*[]*[]string)
		for _, a := range list {
			out = append(out, *a)
		}
		return out

	case string:
		return [][]string{[]string{s.(string)}}

	case []string:
		return [][]string{s.([]string)}

	case *[]string:
		return [][]string{*s.(*[]string)}

	case []interface{}:
		list := s.([]interface{})
		out := [][]string{}
		for i := range list {
			out = append(out, AnyToStringArray(list[i]))
		}
		return out

	case []int:
		list := s.([]int)
		out := [][]string{}

		for i := range list {
			out = append(out, []string{AnyToString(list[i])})
		}
		return out

	case *[]int:
		return AnyTo2StringArray(*s.(*[]int))

	case nil:
		return [][]string{}
	default:
		return [][]string{AnyToStringArray(s)}
	}
	return [][]string{}
}

func AnyToStringArray(s interface{}) []string {

	log.Printf("AnyToStringArray: %T\n", s)

	switch s.(type) {
	case string:
		return []string{s.(string)}
	case []string:
		return s.([]string)

	case *[]string:
		return *s.(*[]string)

	case *[]*string:
		return AnyToStringArray(*s.(*[]*string))

	case []*string:
		list := s.([]*string)
		out := make([]string, len(list))
		for i := range list {
			out[i] = *list[i]
		}
		return out

	case []interface{}:
		list := s.([]interface{})
		out := make([]string, len(list))
		for i := range list {
			out[i] = AnyToString(list[i])
		}
		return out

	case *[]int:
		return AnyToStringArray(*s.(*[]int))

	case []int:
		list := s.([]int)
		out := make([]string, len(list))
		for i := range list {
			out[i] = AnyToString(list[i])
		}
		return out
	case nil:
		return []string{}
	default:
		return []string{AnyToString(s)}
	}
	return []string{}
}
