package Common

import (
	//"BUtils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func FirstNotNil(list ...error) error {
	for _, err := range list {
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseJsonForm(req *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(req.Body)
	js := map[string]interface{}{}
	err := decoder.Decode(&js)

	ID := string(req.URL.Query().Get(":id"))
	log.Printf("ID: %s\n", ID)
	if ID != "" {
		js["id"] = ID
	}

	return js, err
}

func ParseForm(req *http.Request) (map[string][]string, error) {
	if err := req.ParseForm(); err != nil {
		log.Printf("ParseForm err: %s\n", err)
		return nil, err
	}

	log.Printf("ParseForm req.Form: %+v\n", req.Form)

	ID := string(req.URL.Query().Get(":id"))
	log.Printf("ID: %s\n", ID)
	if ID != "" {
		req.Form.Add("id", ID)
	}

	return req.Form, nil
}

func IntParam(form map[string][]string, key string, def ...int) (int, error) {

	list, find := form[key]

	if find && 1 == len(list) {
		return strconv.Atoi(list[0])
	}

	if len(def) > 0 {
		return def[0], nil
	}

	return 0, fmt.Errorf("Integer param %s not found", key)
}

func StringParam(form map[string][]string, key string, def ...string) (string, error) {

	list, find := form[key]

	if find && 1 == len(list) {
		return list[0], nil
	}

	if len(def) > 0 {
		return def[0], nil
	}

	return "", fmt.Errorf("String param %s not found", key)
}

func StringSliceParam(form map[string][]string, key string) ([]string, error) {

	list, find := form[key]

	if find && 0 < len(list) {
		return list, nil
	}

	return []string{}, fmt.Errorf("StringSlice param %s not found", key)
}

func FloatParam(req *http.Request, key string, def ...float64) (float64, error) {

	str := string(req.FormValue(key))

	if str == "" {
		if len(def) > 0 {
			return def[0], nil
		}
		return 0.0, errors.New("Not found")
	}

	return 0, fmt.Errorf("Float param %s not found", key)
}
