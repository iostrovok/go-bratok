package Config

import (
	"log"
)

// Update whole config, all data
func (config *Config) FullUpdate(data []byte) error {

	log.Printf("\n\nFullUpdat 0 data: %s\n", data)

	config.mu.Lock()
	defer config.mu.Unlock()

	// Update static file data
	if res := config.ConfigData.FromLine(data); res != nil {
		log.Printf("FullUpdat. config.ConfigData.FromLine. error: %s\n", res)
		return res
	}

	newScripts := map[string]bool{}
	for _, sc := range config.ConfigData.ListScript() {
		newScripts[sc.ID] = true
	}

	log.Printf("\n\nFullUpdat 2-0 newScripts: %v\n", newScripts)

	// Delete old scripts && update exists
	oldScripts := config._scriptsListNonLock()

	log.Printf("\n\nFullUpdat 2-1 oldScripts: %v\n", oldScripts)
	log.Printf("\n\nFullUpdat 2-2 oldScripts: %T\n", oldScripts)

	for _, sc := range oldScripts {
		log.Printf("\n\nFullUpdat 3-0 data. sc.ID: %s\n", sc.ID)
		if _, find := newScripts[sc.ID]; !find {
			// Delete old scripts
			config.RemoveCronScriptNonLock(sc.ID)
		} else {
			// Update exists
			if scriptJs, findSc := config.ConfigData.FindScript(sc.ID); findSc {
				sc := config.InitCronScript(scriptJs)
				config.UpdateCronScriptNonLock(sc)
			}
		}

		newScripts[sc.ID] = false
	}
	log.Printf("\n\nFullUpdat 3 data\n")

	// Insert new scripts
	for id, needAdd := range newScripts {
		if needAdd {
			if scriptJs, findSc := config.ConfigData.FindScript(id); findSc {
				sc := config.InitCronScript(scriptJs)
				config.AddCronScript(sc)
			}
		}
	}

	// RemoveCronScript

	log.Printf("FullUpdate last data: %s\n", data)

	return config.Store(true)
}
