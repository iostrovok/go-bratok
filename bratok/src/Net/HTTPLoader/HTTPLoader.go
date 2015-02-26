package HTTPLoader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPLoader is loader form HTTP
type HTTPLoaderBase interface {
	Load(kind, how, url string, v url.Values) ([]byte, error)
	LoadJson(kind, how, url string, v map[string]interface{}) (map[string]interface{}, error)
}

// HTTPLoader is loader form HTTP
type HTTPLoader struct {
}

// NewHTTPLoader is constructor
func NewHTTPLoader() *HTTPLoader {
	return &HTTPLoader{}
}

func (f *HTTPLoader) LoadJson(PGPD, hostPath string,
	data map[string]interface{}) (map[string]interface{}, error) {

	out := map[string]interface{}{}

	body, errLoad := f.Load("json", PGPD, hostPath, data)
	if errLoad != nil {
		return out, errLoad
	}

	err := json.Unmarshal(body, &out)
	return out, err
}

// Load - loads config
func (f *HTTPLoader) Load(kind, PGPD, hostPath string, data map[string]interface{}) ([]byte, error) {

	if err := checkParams(kind, PGPD, hostPath); err != nil {
		return nil, err
	}

	params := url.Values{}
	for k, s := range data {
		params.Set(k, iutils.AnyToString(s))
	}

	req, errNewRequest := http.NewRequest(PGPD, hostPath, bytes.NewBuffer([]byte(params.Encode())))
	if errNewRequest != nil {
		return nil, errNewRequest
	}

	if kind == "json" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "text/html")
	}

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func checkParams(kind, PGPD, hostPath string) error {
	if PGPD != "POST" && PGPD != "GET" && PGPD != "PUT" && PGPD != "DELETE" {
		return fmt.Errorf("Wrong type of request, how. Now is '%s'. Must be POST, GET, PUT, DELETE.", PGPD)
	}

	if hostPath == "" {
		return errors.New("Empty url/hostPath.")
	}

	return nil
}
