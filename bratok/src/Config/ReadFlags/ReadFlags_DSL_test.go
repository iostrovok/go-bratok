package ReadFlags

import (
	"errors"
	"flag"
)

/*
Flags is a structure with all command line flags
*/
type FlagsDSL struct {
	ConfigHost string
	DebugMode  bool
	Host       string
	IP         string
	Port       int
	ServerID   string
	Path       string
}

func NewFlagsDSL() *FlagsDSL {

	f := FlagsDSL{
		IP:         "127.0.0.1",
		Port:       10200,
		Host:       "localhost",
		ServerID:   "",
		DebugMode:  false,
		ConfigHost: "",
		Path:       "",
	}

	return &f
}

func (f *FlagsDSL) ReadFlags() error {

	debug := flag.Bool("debug", false, "-d")
	port := flag.Int("port", 10200, "-p")
	ip := flag.String("ip", "127.0.0.1", "-ip")
	host := flag.String("host", "localhost", "-h")
	serverID := flag.String("id", "", "-id")
	configHost := flag.String("config_host", "", "-ch")
	path := flag.String("path", "", "-pa")

	flag.Parse()

	f.IP = *ip
	f.Port = *port
	f.Host = *host
	f.ServerID = *serverID
	f.DebugMode = *debug
	f.ConfigHost = *configHost
	f.Path = *path

	return f.CheckFlags()
}

// CheckFlags checks wrong/true flags
func (f *FlagsDSL) CheckFlags() error {

	if f.ServerID == "" {
		return errors.New("No defined server ID")
	}

	return nil
}
