package Config

import (
	"Config/CronScript"
	"Config/ReadFlags"
	"encoding/json"
	"errors"
	"io/ioutil"
	//"log"
	"os"
	"sort"
	"sync"
	"time"
)

/*
Config is a main structure
*/

const (
	formatConfID string = "20060102150405999999999"
)

type Config struct {
	ServerID       string
	Scripts        map[string]*CronScript.Script
	scriptLogDir   string
	scriptLogFile  string
	configFile     string
	autoConfigFile string
	useAutoConfig  bool
	staticDir      string
	flags          *ReadFlags.Flags
	ConfigData     *ConfigData
	ErrorLoad      error
	mu             *sync.Mutex
}

func New(flags *ReadFlags.Flags) *Config {
	config := Config{}
	config.mu = &sync.Mutex{}
	config.Scripts = map[string]*CronScript.Script{}

	config.ServerID = flags.ServerID

	config.scriptLogDir = "/tmp/"
	config.scriptLogFile = "bratok.scripts.log"

	if flags != nil {
		config.flags = flags

		if config.flags.ConfFile != "" {
			config.configFile = config.flags.ConfFile
		}
	}

	config.autoConfigFile = config.configFile + ".auto.js"
	if _, err := os.Stat(config.autoConfigFile); err == nil {
		config.useAutoConfig = true
	}

	config._checkConfigData()
	config.ErrorLoad = config.LoadConfigFile()

	return &config
}

func (config *Config) GetHttpData() (*Server, error) {

	for _, d := range *config.ConfigData.Servers {
		if config.ServerID == d.ID {
			return d, nil
		}
	}

	return nil, errors.New("Not found")
}

func (config *Config) GetConfigDataByte() ([]byte, error) {
	data, err := json.Marshal(config.ConfigData)

	if err == nil && len(data) == 0 {
		err = errors.New("Empty config data")
	}

	return data, err
}

func (config *Config) Store() error {

	n := time.Now()
	config.ConfigData.ConfigID = int64(n.Unix())*100000 + int64(n.Nanosecond()%100000)

	data, err := json.Marshal(config.ConfigData)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.autoConfigFile, data, 0644)
}

func (config *Config) ScriptStaticDir(d ...string) string {
	if len(d) > 0 {
		config.staticDir = d[0]
	}
	return config.staticDir
}

func (config *Config) ScriptLogDir(d ...string) string {
	if len(d) > 0 {
		config.scriptLogDir = d[0]
	}
	return config.scriptLogDir
}

func (config *Config) ScriptLogFile(f ...string) string {
	if len(f) > 0 {
		config.scriptLogFile = f[0]
	}
	return config.scriptLogDir + config.scriptLogFile
}

func (config *Config) InitNew(data map[string]interface{}) error {
	return nil
}

func (config *Config) AddScript(script *CronScript.Script) error {

	config.Scripts[script.ID] = script

	return nil
}

func (config *Config) GetScript(id string) (*CronScript.Script, bool) {

	for keyId, script := range config.Scripts {
		if keyId == id {
			return script, true
		}
	}

	return nil, false
}

func (config *Config) RaplaceScript(script *CronScript.Script) bool {
	config.mu.Lock()
	defer config.mu.Unlock()

	for keyId := range config.Scripts {
		if keyId == script.ID {
			config.Scripts[keyId] = script
			break
		}
	}

	find := false

	fdata := &Script{
		ID:     script.ID,
		Time:   script.TimeStr,
		Exe:    script.Exe,
		Params: script.Params,
		Evn:    script.Env,
	}

	list := *config.ConfigData.Scripts
	for i, sc := range list {
		if sc.ID == script.ID {
			list[i] = fdata
			find = true
			break
		}
	}

	if !find {
		list = append(list, fdata)
	}

	config.ConfigData.Scripts = &list

	return find
}

func (config *Config) StartNow(t time.Time) []*CronScript.Script {

	out := []*CronScript.Script{}

	for _, script := range config.Scripts {
		if script.StartNow(t) {
			out = append(out, script)
		}
	}

	return out
}

func (config *Config) ScriptsList(server_id string) []*CronScript.Script {
	out := []*CronScript.Script{}

	for _, script := range config.Scripts {
		out = append(out, script)
	}

	sort.Sort(CronScript.SortList(out))

	return out
}

func (config *Config) ServersList() []*Server {
	return *config.ConfigData.Servers
}
