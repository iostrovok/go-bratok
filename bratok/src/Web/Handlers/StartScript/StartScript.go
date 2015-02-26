package StartScript

import (
	"Web/Handlers/Common"
	"fmt"
	"log"
	"net/http"
)

func Handler(res http.ResponseWriter, req *http.Request) {

	id := string(req.URL.Query().Get(":id"))

	log.Printf("StartScript.Handler. ID: %s", id)

	if id == "" {
		err := fmt.Errorf("Script's %s is not found", id)
		Common.ErrorPage(res, err)
		return
	}

	dataForManager := map[string]interface{}{
		"script_id": id,
	}

	mes := Common.ToFromManager(id, "start_script", dataForManager)
	Common.SendJsonMess(res, mes)
}
