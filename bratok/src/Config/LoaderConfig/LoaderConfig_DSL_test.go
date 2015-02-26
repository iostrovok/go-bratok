package LoaderConfig

import (
	"errors"
	"time"
)

/*
LoaderConfig is a structure with all command line flags
*/
type LoaderConfigDSL struct {
	URL      string
	body     []byte
	LastTime time.Time
}

// NewLoaderConfig is constructor
func NewLoaderConfigDSL() *LoaderConfigDSL {

	l := LoaderConfigDSL{
		body: []byte{},
	}

	return &l
}

// Load - loads config
func (f *LoaderConfigDSL) Load(url string) error {
	if url == "" {
		return errors.New("URL mustn't be empty")
	}

	return nil
}

// Body - return body as string
func (f *LoaderConfigDSL) Body() (string, error) {
	return "", nil
}
