package Data

import (
// "Config/History"
// "encoding/json"
// "errors"
// "io/ioutil"
// "log"
// "os"
// "sync"
// "time"
)

const (
	DefaultCofigFile string = "./bratok/conf/config.js"
)

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
