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

/* Data from file file */
type Data struct {
	Scripts        []*Script `json:"scripts"`
	Servers        []*Server `json:"servers"`
	ConfigID       int64     `json:"file_id"`
	LogFile        string    `json:"logfile"`
	StaticFilesDir string    `json:"staticdir"`
}

// IData
type IData interface {
	IsEmpty() bool
	GetId() string
}

type Script struct {
	IData
	ID     string   `json:"id"`
	Time   []string `json:"time"`
	Exe    string   `json:"exe"`
	Params []string `json:"params"`
	Evn    []string `json:"env"`
}

type Server struct {
	IData
	ID             string   `json:"id"`
	IP             string   `json:"ip"`
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	IsMaster       bool     `json:"is_master"`
	Scripts        []string `json:"scripts"`
	StaticFilesDir string   `json:"staticdir"`
	LogFile        string   `json:"logfile"`
}

var _ IData = &Script{}
var _ IData = &Server{}

func EmptyData() Data {

	d := Data{
		Scripts: []*Script{},
		Servers: []*Server{},
	}

	return d
}

func (script Script) GetId() string {
	return script.ID
}

func (server Server) GetId() string {
	return server.ID
}

func (server Script) IsEmpty() bool {
	return false
}

func (script Server) IsEmpty() bool {
	return false
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
