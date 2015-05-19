package ScriptsEdit

import (
	"BUtils"
	"Web/Handlers/Common"
	"fmt"
	"log"
	"net/http"
)

/* Handlers ====> */
// OneScript is handler
// Returns info about the script
func One(res http.ResponseWriter, req *http.Request) {

	id := string(req.URL.Query().Get(":id"))

	if id == "" {
		err := fmt.Errorf("Script's %s is not found", id)
		Common.ErrorPage(res, err)
		return
	}

	dataForManager := map[string]interface{}{
		"script_id": id,
	}

	mes := Common.ToFromManager(id, "script_info", dataForManager)
	Common.SendJsonMess(res, mes)
}

// OneScript is handler
// Returns info about the script
func Log(res http.ResponseWriter, req *http.Request) {

	id := string(req.URL.Query().Get(":id"))

	if id == "" {
		err := fmt.Errorf("Script's %s is not found", id)
		Common.ErrorPage(res, err)
		return
	}

	dataForManager := map[string]interface{}{
		"script_id": id,
	}

	mes := Common.ToFromManager(id, "script_log", dataForManager)
	Common.SendJsonMess(res, mes)
}

// All is handler
// Returns all script for one or all servers
func All(res http.ResponseWriter, req *http.Request) {

	server_id := string(req.URL.Query().Get("server_id"))

	dataForManager := map[string]interface{}{"script_id": server_id}

	mes := Common.ToFromManager(server_id, "scripts_list", dataForManager)
	Common.SendJsonMess(res, mes)
}

// Save is handler
func Save(res http.ResponseWriter, req *http.Request) {

	form, err := Common.ParseJsonForm(req)
	if err != nil {
		Common.ErrorPage(res, err)
		return
	}

	for k, v := range form {
		log.Printf("ScriptsEdit.Save. form: %s := %v\n", k, v)
	}

	log.Printf("ScriptsEdit.Save. form: %+v\n", form)

	id := BUtils.AnyToString(BUtils.GetPath(form, "id"))
	exe := BUtils.AnyToString(BUtils.GetPath(form, "exe"))
	times := BUtils.SplitTimeList(BUtils.AnyToStringArray(BUtils.GetPath(form, "times")))
	env := BUtils.AnyToStringArray(BUtils.GetPath(form, "env"))
	params := BUtils.AnyToStringArray(BUtils.GetPath(form, "params"))

	dataForManager := map[string]interface{}{
		"id":     id,
		"exe":    exe,
		"times":  times,
		"env":    env,
		"params": params,
	}

	log.Printf("ScriptsEdit.Save dataForManager: %+v\n", dataForManager)

	mes := Common.ToFromManager(id, "save_script", dataForManager)
	Common.SendJsonMess(res, mes)
}

/* <===== Handlers */
