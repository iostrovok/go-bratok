package Servers

import (
	// "BUtils"
	"Web/Handlers/Common"
	// "fmt"
	"log"
	"net/http"
)

/* Handlers ====> */
// ConfigUpdate is handler
// Update config
func ConfigUpdate(res http.ResponseWriter, req *http.Request) {

}

// Config returns config
func Config(res http.ResponseWriter, req *http.Request) {
	conf := Common.GetConfigData()
	log.Printf("Servers. Config: %s\n ", conf)
	Common.SendJsonSuccess(res, conf)
}

/* <===== Handlers */
