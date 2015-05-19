package ReadFlags

import (
	"errors"
	"flag"
)

func NewTest() *Flags {
	f := Flags{
		IP:         "127.0.0.1",
		Port:       21222,
		Host:       "localhost",
		ServerID:   "workstation",
		DebugMode:  true,
		ConfigHost: "",
		Path:       "./",
		//ConfFile:   "./bratok/conf/config.js",
		ConfFile: "",
	}

	return &f
}

/*
Flags is a structure with all command line flags
*/
type Flags struct {
	ConfigHost string
	DebugMode  bool
	Host       string
	IP         string
	Port       int
	ServerID   string
	Path       string
	ConfFile   string
}

func New() *Flags {

	f := Flags{
		IP:         "127.0.0.1",
		Port:       0,
		Host:       "localhost",
		ServerID:   "",
		DebugMode:  false,
		ConfigHost: "",
		Path:       "",
		ConfFile:   "",
	}

	return &f
}

func (f *Flags) Read() error {

	debug := flag.Bool("debug", false, "-d")
	port := flag.Int("port", 0, "-po")
	ip := flag.String("ip", "127.0.0.1", "-ip")
	host := flag.String("host", "localhost", "-h")
	serverID := flag.String("id", "", "-id")
	configHost := flag.String("master_host", "", "-ch")
	path := flag.String("path", "", "-pa")
	confFile := flag.String("config", "", "-cf")

	flag.Parse()

	f.IP = *ip
	f.Port = *port
	f.Host = *host
	f.ServerID = *serverID
	f.DebugMode = *debug
	f.ConfigHost = *configHost
	f.Path = *path
	f.ConfFile = *confFile

	return f.Check()
}

// CheckFlags checks wrong/true flags
func (f *Flags) Check() error {

	if f.ServerID == "" {
		return errors.New("No defined server ID")
	}

	return nil
}
