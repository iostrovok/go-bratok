package Config

import (
	"Config/CronScript"
	D "Config/Data"
	"Config/File"
	"Config/ReadFlags"

	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
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
	ServerID   string
	Scripts    map[string]*CronScript.Script
	logFile    string
	configFile string
	staticDir  string
	flags      *ReadFlags.Flags
	ConfigData *File.File
	ErrorLoad  error
	mu         *sync.Mutex
}

func New(flags *ReadFlags.Flags) *Config {

	config := _init(flags)
	if config.ErrorLoad != nil {
		return config
	}

	if config.flags.ConfFile != "" {
		log.Printf("config load from file %s\n", config.flags.ConfFile)
		config.configFile = config.flags.ConfFile
	}

	config._postInit()
	config.ErrorLoad = config.LoadConfigFile()

	return config
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
	err := config.LoadHTTPLine(RemoutConfig)
	log.Printf("Config.NewFromRemout - config.FromLine: %s\n", err)

	return config, err
}

func (config *Config) LoadHTTPLine(data []byte) error {

	log.Printf("config - LoadHTTPLine: %s\n-------------------------\n", data)

	err := config.ConfigData.LoadHTTPLine(data)
	if err != nil {
		return err
	}

	return config._loadConfigData()
}

func (config *Config) FromLine(RemoutConfig []byte) error {

	if err := config.ConfigData.FromLine(RemoutConfig); err != nil {
		return err
	}

	return config._loadConfigData()
}

func (config *Config) _loadConfigData() error {
	cd := config.ConfigData

	if cd.LogFile() != "" {
		config.logFile = cd.LogFile()
	}
	if cd.StaticFilesDir() != "" {
		config.staticDir = cd.StaticFilesDir()
	}

	if scripts := cd.ListScript(); len(scripts) > 0 {
		for _, scriptJs := range scripts {
			sc := config.InitCronScript(scriptJs)
			config.AddCronScript(sc)
		}
	}

	return nil
}

func (config *Config) LoadConfigFile() error {
	// then config file settings

	if res := config.ConfigData.LoadFile(); res != nil {
		return res
	}

	return config._loadConfigData()
}

func Upadte() {

}

func (config *Config) ID() int64 {
	return config.ConfigData.ConfigID()
}

func (config *Config) _postInit() {
	if config.ConfigData == nil {
		config.ConfigData = File.New(config.ServerID, config.configFile)
	}
}

func _init(flags *ReadFlags.Flags) *Config {
	config := Config{}
	config.mu = &sync.Mutex{}
	config.Scripts = map[string]*CronScript.Script{}

	config.logFile = "/tmp/bratok.scripts.log"

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

	if flags.StaticDir != "" {
		config.staticDir = flags.StaticDir
	}
	if flags.LogFile != "" {
		config.logFile = flags.LogFile
	}

	dir := path.Dir(config.logFile)
	exDir, errDir := existsDirFile(dir)
	if errDir != nil {
		config.ErrorLoad = errDir
	} else if !exDir {
		config.ErrorLoad = os.MkdirAll(dir, 0777)
	}

	return &config
}

// exists returns whether the given file or directory exists or not
func existsDirFile(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func (config *Config) NextConfigId() {
}

func (config *Config) GetHTTPData() (*D.Server, error) {
	if sever, find := config.GetServer(config.ServerID); find {
		return sever, nil
	}

	return nil, fmt.Errorf("GetHTTPData. Not found server for %s", config.ServerID)
}

func (config *Config) GetConfigDataByte() ([]byte, error) {
	return config.ConfigData.Byte()
}

func (config *Config) Store(noUpdateIds ...bool) error {

	log.Printf("Config.Store noUpdateIds: %f\n", noUpdateIds)

	return config.ConfigData.Store(noUpdateIds...)
}

func (config *Config) ScriptStaticDir(d ...string) string {
	// config.mu.Lock()
	// defer config.mu.Unlock()

	// if len(d) > 0 {
	// 	config.staticDir = d[0]
	// }
	return config.staticDir
}

func (config *Config) LogFile(f ...string) string {
	return config.logFile
}

func (config *Config) InitNew(data map[string]interface{}) error {
	return nil
}

func (config *Config) InitCronScript(scriptJs *D.Script) *CronScript.Script {
	script := CronScript.New(scriptJs.ID, scriptJs.Exe, scriptJs.Params...)
	script.SetEnv(scriptJs.Evn)
	for _, oneTime := range scriptJs.Time {
		times := strings.Split(oneTime, " ")
		script.SetTime(times...)
	}

	return script
}

func (config *Config) UpdateCronScript(script *CronScript.Script) error {
	config.mu.Lock()
	defer config.mu.Unlock()
	return config.UpdateCronScriptNonLock(script)
}

func (config *Config) UpdateCronScriptNonLock(script *CronScript.Script) error {
	if _, find := config.Scripts[script.ID]; !find {
		return fmt.Errorf("UpdateCronScriptNonLock.script %s not found\n", script.ID)
	}

	config.Scripts[script.ID].Update(script)

	return nil
}

func (config *Config) AddCronScript(script *CronScript.Script) error {
	config.mu.Lock()
	defer config.mu.Unlock()
	return config.AddCronScriptNonLock(script)
}

func (config *Config) AddCronScriptNonLock(script *CronScript.Script) error {

	if _, find := config.Scripts[script.ID]; !find {
		config.Scripts[script.ID] = script
	}

	return nil
}

func (config *Config) CheckDeletedScript(script_id string) {
	config.mu.Lock()
	defer config.mu.Unlock()

	script, find := config.Scripts[script_id]
	if !find {
		return
	}

	if !script.IsDeleted {
		return
	}

	if script.IsWork {
		return
	}

	delete(config.Scripts, script_id)
}

func (config *Config) RemoveCronScript(script_id string) {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.RemoveCronScriptNonLock(script_id)
}

func (config *Config) RemoveCronScriptNonLock(script_id string) {

	script, find := config.Scripts[script_id]
	if !find {
		return
	}
	if script.IsWork {
		script.IsDeleted = true
		return
	}

	delete(config.Scripts, script_id)
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

func (config *Config) GetServer(id string) (*D.Server, bool) {
	return config.ConfigData.FindServer(id)
}

func (config *Config) ReplaceScript(script *CronScript.Script) bool {
	config.mu.Lock()
	defer config.mu.Unlock()

	findExe := false

	// Update exists script
	for keyId, oldScript := range config.Scripts {
		if keyId == script.ID {

			oldScript.Update(script)
			config.Scripts[keyId] = oldScript

			findExe = true
			break
		}
	}

	if !findExe {
		config.AddCronScriptNonLock(script)
	}

	fdata := &D.Script{
		ID:     script.ID,
		Time:   script.TimeStr,
		Exe:    script.Exe,
		Params: script.Params,
		Evn:    script.Env,
	}

	return config.ConfigData.SetScript(fdata)
}

func (config *Config) ReplaceServer(server *D.Server) bool {
	config.mu.Lock()
	defer config.mu.Unlock()

	return config.ConfigData.SetServer(server)
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

	serverId := ""
	if len(server_ids) > 0 {
		serverId = server_ids[0]
	}

	if serverId == "" {
		return config._scriptsList()
	}

	out := []*CronScript.Script{}
	server, find := config.GetServer(serverId)
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
	return config._scriptsListNonLock()
}

func (config *Config) _scriptsListNonLock() []*CronScript.Script {
	out := []*CronScript.Script{}

	for _, script := range config.Scripts {
		out = append(out, script)
	}

	sort.Sort(CronScript.SortList(out))

	return out
}

func (config *Config) ServersList() []*D.Server {
	return config.ConfigData.ListServer()
}
