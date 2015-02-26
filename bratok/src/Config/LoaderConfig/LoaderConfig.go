package LoaderConfig

import (
	"Config/Config"
	"Net/HTTPLoader"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type LoaderConfigBase interface {
}

/*
LoaderConfig is a structure with all command line flags
*/
type LoaderConfig struct {
	c *Config.Config
	h *HTTPLoader.HTTPLoader
}

// NewLoaderConfig is constructor
func NewLoaderConfig(c *Config.Config, hl *HTTPLoader.HTTPLoader) *LoaderConfig {

	lc := LoaderConfig{
		c: c,
		h: hl,
	}

	return &lc
}

// Load - loads config
func (lc *LoaderConfig) LoadURL(hostPath string) error {

	js, errLoadJson := lc.h.LoadJson("GET", hostPath, map[string]interface{}{})
	if errLoadJson != nil {
		return errLoadJson
	}

	return lc.c.InitNew(js)
}

// Load - loads config
func (lc *LoaderConfig) LoadFile(file string) error {

	data, errReadFile := ioutil.ReadFile(file)

	if errReadFile != nil {
		return errReadFile
	}

	fmt.Print(string(data))

	js, errToJson := toJson(data)
	if errToJson != nil {
		return errToJson
	}

	return lc.c.InitNew(js)
}

func toJson(data []byte) (map[string]interface{}, error) {

	out := map[string]interface{}{}

	err := json.Unmarshal(data, &out)
	return out, err
}
