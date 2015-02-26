package Config

import (
	"Config/CronScript"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

func TestFileLine() []byte {
	return []byte(`{
			"logfile":"FILE","logdir":"DIR","staticdir":"STATIC-DIR","config_id":12312312,
			"scripts":[
				{"id":"ls1","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls2","time":["*/2"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]}
			],
			"servers":[
				{"id":"workstation","ip":"127.0.0.1","host":"","port":21222,"is_master":true},
				{"id":"somethere","ip":"192.168.0.10","host":"wks-l","port":21222,"is_master":false}
			]
	}`)
}

//var TestFileLine []byte = []byte(`{"logfile":"","logdir":"","scripts":[{"id":"ls11231","time":["*","*","*"],"exe":"sh","params":["ls","-a","./]}],"servers":[{"id":"workstation","ip":"127.0.0.1","host":"","port":21222}]}`)

type Script struct {
	ID     string   `json:"id"`
	Time   []string `json:"time"`
	Exe    string   `json:"exe"`
	Params []string `json:"params"`
	Evn    []string `json:"env"`
}

type Server struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	IsMaster bool   `json:"is_master"`
}

/* Data from config file */
type ConfigData struct {
	Scripts        *[]*Script `json:"scripts"`
	Servers        *[]*Server `json:"servers"`
	ConfigID       int64      `json:"config_id"`
	ScriptLogDir   string     `json:"logdir"`
	ScriptLogFile  string     `json:"logfile"`
	StaticFilesDir string     `json:"staticdir"`
}

func (config *Config) LoadConfigFile() error {
	// then config file settings

	if config.configFile == "" {
		return nil
	}

	log.Printf("config.configFile: %s\n", config.configFile)
	log.Printf("config.autoConfigFile: %s\n", config.autoConfigFile)

	file := config.configFile

	if config.useAutoConfig {
		file = config.autoConfigFile
	}

	log.Printf("\nfile: %s\n", file)

	data, errReadFile := ioutil.ReadFile(file)
	log.Printf("\ndata: %s\n", data)

	if errReadFile != nil {
		log.Printf("0. LoadConfigFile: error: %s\n", errReadFile)
		return errReadFile
	}

	if err := config._parseConfigData(data); err != nil {
		log.Printf("1. LoadConfigFile: error: %s\n", err)
		log.Printf("1. LoadConfigFile: data: %s\n", data)
		return err
	}

	config._checkConfigData()

	log.Printf("config.ConfigData: %+v\n", config.ConfigData)

	if err := config._loadConfigData(); err != nil {
		log.Printf("2. LoadConfigFile: error: %s\n", err)
		return err
	}

	return nil
}

func (config *Config) _checkConfigData() {
	if config.ConfigData == nil {
		config.ConfigData = &ConfigData{}
	}
	if config.ConfigData.Servers == nil {
		config.ConfigData.Servers = &[]*Server{}
	}
	if config.ConfigData.Scripts == nil {
		config.ConfigData.Scripts = &[]*Script{}
	}
}

func (config *Config) _parseConfigData(data []byte) error {
	config.ConfigData = &ConfigData{}
	return json.Unmarshal(data, config.ConfigData)
}

func (config *Config) _loadConfigData() error {
	cf := config.ConfigData

	if cf.ScriptLogDir != "" {
		config.scriptLogDir = cf.ScriptLogDir
	}
	if cf.ScriptLogFile != "" {
		config.scriptLogFile = cf.ScriptLogFile
	}
	if cf.StaticFilesDir != "" {
		config.staticDir = cf.StaticFilesDir
	}

	log.Printf("cf: %+v\n", cf)

	if cf.Scripts != nil {
		for _, scriptJs := range *cf.Scripts {
			script := CronScript.New(scriptJs.ID, scriptJs.Exe, scriptJs.Params...)
			script.SetEnv(scriptJs.Evn)
			for _, oneTime := range scriptJs.Time {
				//Evn    []string `json:"env"`

				times := strings.Split(oneTime, " ")
				log.Printf("times v: %+v\n", times)
				log.Printf("times T: %T\n", times)
				log.Printf("times s: %s\n", times)
				log.Printf("times len(times): %d\n", len(times))

				script.SetTime(times...)

			}

			log.Printf("script: %+v\n", script)
			config.AddScript(script)
		}
	}

	// type ConfigData struct {
	// 	Scripts       *[]Script `json:"scripts"`
	// 	Servers       *[]Server `json:"servers"`
	// 	ScriptLogDir  string    `json:"logdir"`
	// 	 string    `json:"logfile"`

	log.Printf("_loadConfigData: config 2: %+v\n", config)

	return nil
}
