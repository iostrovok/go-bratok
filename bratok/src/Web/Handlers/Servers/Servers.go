package Servers

import (
	// "BUtils"
	"Web/Handlers/Common"
	"fmt"
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

// All is handler
// Returns all servers from config
func All(res http.ResponseWriter, req *http.Request) {

	server_id := string(req.URL.Query().Get("server_id"))

	log.Printf("StatusScript.All. ID: %s", server_id)

	dataForManager := map[string]interface{}{"script_id": server_id}

	mes := Common.ToFromManager("", "server_list", dataForManager)
	Common.SendJsonMess(res, mes)
}

// One is handler
// Returns info about the server
func One(res http.ResponseWriter, req *http.Request) {

	id := string(req.URL.Query().Get(":id"))

	log.Printf("ServersEdit.One ID: %s", id)

	if id == "" {
		err := fmt.Errorf("Server's %s is not found", id)
		Common.ErrorPage(res, err)
		return
	}

	dataForManager := map[string]interface{}{
		"server_id": id,
	}

	mes := Common.ToFromManager(id, "server/info", dataForManager)
	Common.SendJsonMess(res, mes)
}

/* <===== Handlers */
