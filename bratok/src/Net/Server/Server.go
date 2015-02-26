package Server

import (
	"Config/Config"
)

/*
Flags is a structure with all command line flags
*/
type ServerBase interface {
	Start(config *Config.Config) error
}

/*
Flags is a structure with all command line flags
*/
type Server struct {
	config *Config.Config
}

func NewServer() *Server {

	c := Server{}

	return &c
}

func (c *Server) Start(config *Config.Config) error {
	c.config = config
	return nil
}
