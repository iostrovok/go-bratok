package HttpLoader

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HttpLoader is loader form HTTP 
type HttpLoader struct {
}

// NewHttpLoader is constructor
func NewHttpLoader() *HttpLoader {
	return &HttpLoader{}
}

// Load - loads config
func (f *HttpLoader) Load(kind, how, url string, v url.Values) ([]byte, error) {
	//fmt.Printf("H: %s %s\n", how, where)
	//fmt.Printf("jsonStr: %s\n", js))

	if how == "" {
		return nil, errors.New("Empty type of request [how]. Must be POST, GET, PUT, DELETE.")
	}

	if url == "" {
		return nil, errors.New("Empty url.")
	}

	req, errNewRequest := http.NewRequest(how, url, bytes.NewBuffer([]byte(v.Encode()))
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
