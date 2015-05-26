package Servers

import (
	"BUtils"
	//"Config/Config"
	//"Config/File"
	"Web/Handlers/Common"

	//"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/* Handlers ====> */
// ConfigUpdate is handler
// Update config
func ConfigUpdate(res http.ResponseWriter, req *http.Request) {

	log.Printf("ConfigUpdate data: %s\n", req.Body)

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Common.ErrorPage(res, err)
		return
	}

	dataForManager := map[string]interface{}{"body": body}

	mes := Common.ToFromManager("", "config_update", dataForManager)
	Common.SendJsonMess(res, mes)
}

// Config returns config
func GetConfig(res http.ResponseWriter, req *http.Request) {
	conf := Common.GetConfigData()
	log.Printf("Servers. Config: %s\n ", conf)
	Common.SendJsonSuccess(res, conf)
}

// All is handler
// Returns all servers from config
func All(res http.ResponseWriter, req *http.Request) {

	server_id := string(req.URL.Query().Get("server_id"))

	dataForManager := map[string]interface{}{"script_id": server_id}

	mes := Common.ToFromManager("", "server_list", dataForManager)
	Common.SendJsonMess(res, mes)
}

// One is handler
// Returns info about the server
func One(res http.ResponseWriter, req *http.Request) {
	id := string(req.URL.Query().Get(":id"))

	if id == "" {
		err := fmt.Errorf("Server's %s is not found", id)
		Common.ErrorPage(res, err)
		return
	}

	dataForManager := map[string]interface{}{
		"server_id": id,
	}

	mes := Common.ToFromManager(id, "server_info", dataForManager)
	Common.SendJsonMess(res, mes)
}

// Save is handler
func Save(res http.ResponseWriter, req *http.Request) {

	form, err := Common.ParseJsonForm(req)
	if err != nil {
		Common.ErrorPage(res, err)
		return
	}

	id := BUtils.AnyToString(BUtils.GetPath(form, "id"))
	ip := BUtils.AnyToString(BUtils.GetPath(form, "ip"))
	host := BUtils.AnyToString(BUtils.GetPath(form, "host"))
	port := BUtils.AnyToInt(BUtils.GetPath(form, "port"))
	is_master := BUtils.AnyToBool(BUtils.GetPath(form, "is_master"))

	/*  Get only selected scripts from full list */
	scriptsFull := BUtils.AnyToInterfaceArray(BUtils.GetPath(form, "scriptsFull"))

	scripts := []string{}
	for _, s := range scriptsFull {
		if BUtils.AnyToBool(BUtils.GetPath(s, "selected")) {
			id := BUtils.AnyToString(BUtils.GetPath(s, "id"))
			if id != "" {
				scripts = append(scripts, id)
			}
		}
	}

	dataForManager := map[string]interface{}{
		"id":        id,
		"ip":        ip,
		"host":      host,
		"port":      port,
		"is_master": is_master,
		"scripts":   scripts,
	}

	mes := Common.ToFromManager(id, "save_server", dataForManager)
	Common.SendJsonMess(res, mes)
}

/* <===== Handlers */
