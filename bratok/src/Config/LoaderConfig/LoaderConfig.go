package LoaderConfig

import (
	"errors"
	"time"
)

type LoaderConfigBase interface {
}

/*
LoaderConfig is a structure with all command line flags
*/
type LoaderConfig struct {
	URL      string
	body     []byte
	LastTime time.Time
}

// NewLoaderConfig is constructor
func NewLoaderConfig() *LoaderConfig {

	l := LoaderConfig{
		body: []byte{},
	}

	return &l
}

// Load - loads config
func (f *LoaderConfig) Load(url string) error {
	if url == "" {
		return errors.New("URL mustn't be empty")
	}

	return nil
}

// Body - return body as string
func (f *LoaderConfig) Body() (string, error) {
	return "", nil
}
