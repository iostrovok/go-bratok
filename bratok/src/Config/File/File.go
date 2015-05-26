package File

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func TestFileLineHttp() []byte {
	out := `{"error":"","data":` + string(TestFileLine()) + `}`
	return []byte(out)
}

func TestFileLine() []byte {
	return []byte(`{
			"logfile":"FILE","logdir":"DIR","staticdir":"STATIC-DIR","file_id":12312312,
			"scripts":[
				{"id":"ls22","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls33","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls1","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls2","time":["*/2"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]}
			],
			"servers":[
				{
					"id":"workstation","ip":"127.0.0.1","host":"","port":21222,"is_master":true,"scripts": ["ls2","ls1"],
					"logfile":"workstation_FILE","logdir":"workstation_DIR","staticdir":"workstation_STATIC-DIR"
				},
				{"id":"somethere","ip":"192.168.0.10","host":"wks-l","port":21223,"is_master":false,"scripts": ["ls2"]}
			]
	}`)
}

type Script struct {
	ID     string   `json:"id"`
	Time   []string `json:"time"`
	Exe    string   `json:"exe"`
	Params []string `json:"params"`
	Evn    []string `json:"env"`
}

type Server struct {
	ID             string   `json:"id"`
	IP             string   `json:"ip"`
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	IsMaster       bool     `json:"is_master"`
	Scripts        []string `json:"scripts"`
	StaticFilesDir string   `json:"staticdir"`
	ScriptLogFile  string   `json:"logfile"`
	ScriptLogDir   string   `json:"logdir"`
}

/* Data from file file */
type FileHttp struct {
	Error string `json:"error"`
	Data  Data   `json:"data"`
}

/* Data from file file */
type Data struct {
	Scripts        []*Script `json:"scripts"`
	Servers        []*Server `json:"servers"`
	ConfigID       int64     `json:"file_id"`
	ScriptLogDir   string    `json:"logdir"`
	ScriptLogFile  string    `json:"logfile"`
	StaticFilesDir string    `json:"staticdir"`
	History        *History  `json:"history,omitempty"`
}

/* Data from file file */
type File struct {
	autoConfigFile string
	useAutoConfig  bool

	configFile string
	Data       Data
	mu         *sync.Mutex
	serverID   string
}

func New(ServerID, configFile string) *File {

	data := Data{
		Scripts: []*Script{},
		Servers: []*Server{},
	}

	return &File{
		Data:       data,
		serverID:   ServerID,
		configFile: configFile,
		mu:         &sync.Mutex{},
	}
}

func (file *File) LoadHTTPLine(data []byte) error {

	confData := FileHttp{}
	err := json.Unmarshal(data, &confData)
	if err != nil {
		return err
	}

	file.Data = confData.Data
	return file._checkDefault()
}

func (file *File) SetAutoConfig(u bool, f ...string) error {

	if len(f) > 0 {
		file.autoConfigFile = f[0]
	}

	if u && file.autoConfigFile != "" {
		file.useAutoConfig = true
		return nil
	}

	if !u {
		file.useAutoConfig = false
		return nil
	}

	return errors.New("Bad option params for AutoConfig")
}

func (file *File) LoadFile() error {
	// then file file settings

	f := file.configFile
	if file.useAutoConfig {
		f = file.autoConfigFile
	}

	if f == "" {
		return nil
	}

	data, errReadFile := ioutil.ReadFile(f)

	if errReadFile != nil {
		return errReadFile
	}

	return file.FromLine(data)
}

func (file *File) FromLine(data []byte) error {

	if err := file._parseConfigData(data); err != nil {
		log.Printf("LoadConfigFileFromLine err: %+v\n", err)
		return err
	}

	return file._checkDefault()
}

func (file *File) _checkDefault() error {

	if file.Data.Servers == nil {
		file.Data.Servers = []*Server{}
	}
	if file.Data.Scripts == nil {
		file.Data.Scripts = []*Script{}
	}
	if file.Data.History == nil {
		file.Data.History = NewHistory(file.serverID)
	} else {
		file.Data.History.SetServerID(file.serverID)
	}

	return nil
}

func (file *File) _parseConfigData(data []byte) error {
	file.Data = Data{}
	return json.Unmarshal(data, &file.Data)
}

func (file *File) ServerID() string {
	return file.serverID
}

func (file *File) ConfigID() int64 {
	return file.Data.ConfigID
}

func (file *File) ScriptLogDir() string {
	server, find := file.FindServer(file.serverID)
	if find {
		if server.ScriptLogDir != "" {
			return server.ScriptLogDir
		}
	}

	return file.Data.ScriptLogDir
}

func (file *File) StaticFilesDir() string {
	if server, find := file.FindServer(file.serverID); find {
		if server.StaticFilesDir != "" {
			return server.StaticFilesDir
		}
	}

	return file.Data.StaticFilesDir
}

func (file *File) ScriptLogFile() string {
	if server, find := file.FindServer(file.serverID); find {
		if server.ScriptLogFile != "" {
			return server.ScriptLogFile
		}
	}

	return file.Data.ScriptLogFile
}

func (file *File) FindServer(id string) (*Server, bool) {

	for _, server := range file.Data.Servers {
		if server.ID == id {
			return server, true
		}
	}

	return nil, false
}

func (file *File) ListServer() []*Server {
	return file.Data.Servers
}

func (file *File) ListScript() []*Script {
	return file.Data.Scripts
}

func (file *File) FindScript(id string) (*Script, bool) {

	for _, script := range file.Data.Scripts {
		if script.ID == id {
			return script, true
		}
	}

	return nil, false
}

func (file *File) UpdateId() {
	file.mu.Lock()
	defer file.mu.Unlock()

	n := time.Now()

	file.Data.ConfigID = int64(n.Unix())*100000 + int64(n.Nanosecond()%100000)
}

func (file *File) Store(noUpdateIds ...bool) error {

	noUpdateId := false
	if len(noUpdateIds) > 0 {
		noUpdateId = noUpdateIds[0]
	}

	if !noUpdateId {
		file.UpdateId()
	}

	file.mu.Lock()
	defer file.mu.Unlock()

	data, err := json.Marshal(file.Data)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(file.autoConfigFile, data, 0644)
}

func (file *File) Byte() ([]byte, error) {
	file.mu.Lock()
	defer file.mu.Unlock()

	data, err := json.Marshal(file.Data)

	if err == nil && len(data) == 0 {
		err = errors.New("Empty config data")
	}

	return data, err
}

func (file *File) SetScript(script *Script) bool {

	find := false

	file.mu.Lock()
	defer file.mu.Unlock()

	list := file.Data.Scripts

	for i, sc := range list {
		if sc.ID == script.ID {
			list[i] = sc
			find = true
			break
		}
	}

	if !find {
		list = append(list, script)
	}

	file.Data.Scripts = list

	return find
}

func (file *File) SetServer(server *Server) bool {
	file.mu.Lock()
	defer file.mu.Unlock()

	find := false
	list := file.Data.Servers
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

	file.Data.Servers = list

	return find
}

func (server *Server) Clone() *Server {
	return &Server{
		ID:       server.ID,
		IP:       server.IP,
		Host:     server.Host,
		Port:     server.Port,
		IsMaster: server.IsMaster,
		Scripts:  server.Scripts,
	}
}