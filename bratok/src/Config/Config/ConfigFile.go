package Config

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"log"
// )

// func TestFileLineHttp() []byte {
// 	out := `{"error":"","data":` + string(TestFileLine()) + `}`
// 	return []byte(out)
// }

// func TestFileLine() []byte {
// 	return []byte(`{
// 			"logfile":"FILE","logdir":"DIR","staticdir":"STATIC-DIR","config_id":12312312,
// 			"scripts":[
// 				{"id":"ls22","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
// 				{"id":"ls33","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
// 				{"id":"ls1","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
// 				{"id":"ls2","time":["*/2"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]}
// 			],
// 			"servers":[
// 				{"id":"workstation","ip":"127.0.0.1","host":"","port":21222,"is_master":true,"scripts": ["ls2","ls1"]},
// 				{"id":"somethere","ip":"192.168.0.10","host":"wks-l","port":21222,"is_master":false,"scripts": ["ls2"]}
// 			]
// 	}`)
// }

// type Script struct {
// 	ID     string   `json:"id"`
// 	Time   []string `json:"time"`
// 	Exe    string   `json:"exe"`
// 	Params []string `json:"params"`
// 	Evn    []string `json:"env"`
// }

// type Server struct {
// 	ID             string   `json:"id"`
// 	IP             string   `json:"ip"`
// 	Host           string   `json:"host"`
// 	Port           int      `json:"port"`
// 	IsMaster       bool     `json:"is_master"`
// 	Scripts        []string `json:"scripts"`
// 	StaticFilesDir string   `json:"staticdir"`
// 	ScriptLogFile  string   `json:"logfile"`
// 	ScriptLogDir   string   `json:"logdir"`
// }

// /* Data from config file */
// type ConfigDataHttp struct {
// 	Error string      `json:"error"`
// 	Data  *ConfigData `json:"data"`
// }

// /* Data from config file */
// type ConfigData struct {
// 	Scripts        *[]*Script `json:"scripts"`
// 	Servers        *[]*Server `json:"servers"`
// 	ConfigID       int64      `json:"config_id"`
// 	ScriptLogDir   string     `json:"logdir"`
// 	ScriptLogFile  string     `json:"logfile"`
// 	StaticFilesDir string     `json:"staticdir"`
// }

// func (config *Config) LoadConfigFileFromLine(data []byte) error {

// 	log.Printf("LoadConfigFileFromLine data: %s\n", data)

// 	confData := ConfigDataHttp{}
// 	err := json.Unmarshal(data, &confData)
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("LoadConfigFileFromLine err: %+v\n", err)
// 	log.Printf("LoadConfigFileFromLine confData: %+v\n", confData)

// 	config.ConfigData = confData.Data
// 	return config._postProceecingConfigFile()
// }

// func (config *Config) _postProceecingConfigFile() error {
// 	config._checkConfigData()
// 	return config._loadConfigData()
// }

// func (config *Config) LoadConfigFile() error {
// 	// then config file settings

// 	if config.configFile == "" {
// 		return nil
// 	}

// 	file := config.configFile

// 	if config.useAutoConfig {
// 		file = config.autoConfigFile
// 	}

// 	log.Printf("\nfile: %s\n", file)

// 	data, errReadFile := ioutil.ReadFile(file)
// 	log.Printf("\ndata: %s\n", data)

// 	if errReadFile != nil {
// 		return errReadFile
// 	}

// 	if err := config._parseConfigData(data); err != nil {
// 		return err
// 	}

// 	return config._postProceecingConfigFile()
// }

// func (config *Config) _checkConfigData() {
// 	if config.ConfigData == nil {
// 		config.ConfigData = &ConfigData{}
// 	}
// 	if config.ConfigData.Servers == nil {
// 		config.ConfigData.Servers = &[]*Server{}
// 	}
// 	if config.ConfigData.Scripts == nil {
// 		config.ConfigData.Scripts = &[]*Script{}
// 	}
// }

// func (config *Config) _parseConfigData(data []byte) error {
// 	config.ConfigData = &ConfigData{}
// 	return json.Unmarshal(data, config.ConfigData)
// }

// func (config *Config) _loadConfigData() error {
// 	cf := config.ConfigData

// 	log.Printf("_loadConfigData 1: %+v\n", cf)

// 	if cf.ScriptLogDir != "" {
// 		config.scriptLogDir = cf.ScriptLogDir
// 	}
// 	if cf.ScriptLogFile != "" {
// 		config.scriptLogFile = cf.ScriptLogFile
// 	}
// 	if cf.StaticFilesDir != "" {
// 		config.staticDir = cf.StaticFilesDir
// 	}

// 	log.Printf("_loadConfigData 2: %+v\n", cf)

// 	if cf.Scripts != nil {
// 		for _, scriptJs := range *cf.Scripts {
// 			sc := config.SetCronScript(scriptJs)
// 			config.AddCronScript(sc)
// 		}
// 	}

// 	return nil
// }

// func (server *Server) Clone() *Server {
// 	return &Server{
// 		ID:       server.ID,
// 		IP:       server.IP,
// 		Host:     server.Host,
// 		Port:     server.Port,
// 		IsMaster: server.IsMaster,
// 		Scripts:  server.Scripts,
// 	}
// }
