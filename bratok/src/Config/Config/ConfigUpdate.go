package Config

import (
	"log"
)

// Update whole config, all data
func (config *Config) FullUpdate(data []byte) error {
	config.mu.Lock()
	defer config.mu.Unlock()

	if res := config.ConfigData.FromLine(data); res != nil {
		return res

	}
	// //
	// oldList := config._scriptsList()
	// for

	log.Printf("FullUpdate data: %s\n", data)

	return config.Store(true)
}
