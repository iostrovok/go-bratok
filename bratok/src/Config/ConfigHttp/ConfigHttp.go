package ConfigHttp

import (
	cronConf "Config/Config"
	"Config/ReadFlags"
	"log"
	"os"
)

/*
Config is a main structure
*/
type Config struct {
	_ID       string
	_IP       string
	_Host     string
	_Port     int
	_IsMaster bool
	staticDir string
	staticUrl string
	flags     *ReadFlags.Flags
	cronConf  *cronConf.Config
}

func New(flags *ReadFlags.Flags, cronConf *cronConf.Config) *Config {

	config := Config{
		_Port:     0,
		staticDir: "/tmp",
		staticUrl: "/static",
		flags:     flags,
		cronConf:  cronConf,
	}

	if config.cronConf.ScriptStaticDir() != "" {
		config.staticDir = config.cronConf.ScriptStaticDir()
	} else {
		dir, err := os.Getwd()
		if err == nil {
			config.staticDir = dir + "/static"
		}
		log.Printf("DIR: %s\n", dir)
	}

	/* Reads data from config file */
	if data, err := cronConf.GetHttpData(); err == nil {
		log.Printf("cronConf.GetHttpData. data: %+v\n", data)

		config._Port = data.Port
		config._ID = data.ID
		config._IP = data.IP
		config._Host = data.Host
		config._IsMaster = data.IsMaster
	}

	log.Printf("cronConf.GetHttpData. config: %+v\n", config)

	/* Re-reads data from command line */
	if flags.Port > 0 {
		config._Port = flags.Port
	}

	return &config
}

func (config *Config) Port() int {
	return config._Port
}

func (config *Config) StaticDir(dir ...string) string {
	if len(dir) > 0 {
		config.staticDir = dir[0]
	}
	return config.staticDir
}
func (config *Config) StaticUrl() string {
	return config.staticUrl
}
