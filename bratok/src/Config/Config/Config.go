package Config

import (
	"Config/CronScript"
	"Config/ReadFlags"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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

func NewFromRemout(flags *ReadFlags.Flags, RemoutConfig []byte) (*Config, error) {

	if len(RemoutConfig) == 0 {
		return nil, errors.New("RemoutConfig is empty")
	}

	config := _init(flags)
	if config.ErrorLoad != nil {
		return config, config.ErrorLoad
	}

	config._postInit()
	res := config.LoadConfigFileFromLine(RemoutConfig)

	return config, res
}

func New(flags *ReadFlags.Flags) *Config {
	log.Printf("NEW. flags: %+v\n", flags)

	config := _init(flags)
	if config.ErrorLoad != nil {
		return config
	}

	log.Printf("NEW. config.flags: %+v\n", config.flags)

	if config.flags.ConfFile != "" {
		config.configFile = config.flags.ConfFile
	}

	config._postInit()
	config.ErrorLoad = config.LoadConfigFile()

	return config
}

func (config *Config) _postInit() {

	config.autoConfigFile = config.configFile + ".auto.js"
	if _, err := os.Stat(config.autoConfigFile); err == nil {
		config.useAutoConfig = true
	}

}

func _init(flags *ReadFlags.Flags) *Config {
	config := Config{}
	config.mu = &sync.Mutex{}
	config.Scripts = map[string]*CronScript.Script{}

	config.scriptLogDir = "/tmp/"
	config.scriptLogFile = "bratok.scripts.log"

	if flags == nil {
		config.ErrorLoad = errors.New("No found ServerID")
		return &config
	}

	if flags.ServerID == "" {
		config.ErrorLoad = errors.New("No found ServerID")
		return &config
	}

	config.ServerID = flags.ServerID
	config.flags = flags

	config._checkConfigData()

	return &config
}

func (config *Config) NextConfigId() {
}

func (config *Config) GetHttpData() (*Server, error) {
	if sever, find := config.GetServer(config.ServerID); find {
		return sever, nil
	}

	return nil, errors.New("Not found")
}

func (config *Config) GetConfigDataByte() ([]byte, error) {
	config.mu.Lock()
	defer config.mu.Unlock()

	data, err := json.Marshal(config.ConfigData)

	if err == nil && len(data) == 0 {
		err = errors.New("Empty config data")
	}

	return data, err
}

func (config *Config) UpdateId() {
	config.mu.Lock()
	defer config.mu.Unlock()

	n := time.Now()
	config.ConfigData.ConfigID = int64(n.Unix())*100000 + int64(n.Nanosecond()%100000)
}

func (config *Config) Store() error {
	config.UpdateId()

	config.mu.Lock()
	defer config.mu.Unlock()

	data, err := json.Marshal(config.ConfigData)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.autoConfigFile, data, 0644)
}

func (config *Config) ScriptStaticDir(d ...string) string {
	config.mu.Lock()
	defer config.mu.Unlock()

	if len(d) > 0 {
		config.staticDir = d[0]
	}
	return config.staticDir
}

func (config *Config) ScriptLogDir(d ...string) string {
	config.mu.Lock()
	defer config.mu.Unlock()

	if len(d) > 0 {
		config.scriptLogDir = d[0]
	}
	return config.scriptLogDir
}

func (config *Config) ScriptLogFile(f ...string) string {
	config.mu.Lock()
	defer config.mu.Unlock()

	if len(f) > 0 {
		config.scriptLogFile = f[0]
	}
	return config.scriptLogDir + config.scriptLogFile
}

func (config *Config) InitNew(data map[string]interface{}) error {
	return nil
}

func (config *Config) AddScript(script *CronScript.Script) error {
	config.mu.Lock()
	defer config.mu.Unlock()

	config.Scripts[script.ID] = script

	return nil
}

func (config *Config) GetScript(id string) (*CronScript.Script, bool) {
	config.mu.Lock()
	defer config.mu.Unlock()

	for keyId, script := range config.Scripts {
		if keyId == id {
			return script, true
		}
	}

	return nil, false
}

func (config *Config) GetServer(id string) (*Server, bool) {
	config.mu.Lock()
	defer config.mu.Unlock()

	for _, server := range *config.ConfigData.Servers {
		if server.ID == id {
			return server, true
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

func (config *Config) RaplaceServer(server *Server) bool {
	config.mu.Lock()
	defer config.mu.Unlock()

	find := false
	list := *config.ConfigData.Servers
	for i, s := range list {
		if s.ID == server.ID {
			list[i] = server
			find = true
			break
		}
	}

	if !find {
		list = append(list, server)
	}

	config.ConfigData.Servers = &list

	return false
}

func (config *Config) StartNow(t time.Time) []*CronScript.Script {
	config.mu.Lock()
	defer config.mu.Unlock()

	out := []*CronScript.Script{}

	for _, script := range config.Scripts {
		if script.StartNow(t) {
			out = append(out, script)
		}
	}

	return out
}

func (config *Config) ScriptsList(server_ids ...string) []*CronScript.Script {

	if len(server_ids) == 0 {
		return config._scriptsList()
	}

	out := []*CronScript.Script{}
	server, find := config.GetServer(server_ids[0])
	if !find {
		return out
	}

	config.mu.Lock()
	defer config.mu.Unlock()

	for _, script := range config.Scripts {
		for _, script_id := range server.Scripts {
			if script_id == script.ID {
				out = append(out, script)
				break
			}
		}
	}

	sort.Sort(CronScript.SortList(out))

	return out
}

func (config *Config) _scriptsList() []*CronScript.Script {
	config.mu.Lock()
	defer config.mu.Unlock()

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
