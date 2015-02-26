package Test

import (
	"Web/Handlers/Common"
	"net/http"
)

func Handler(res http.ResponseWriter, req *http.Request) {

	dataForManager := map[string]interface{}{
		"test_line": 1,
	}

	mes := Common.ToFromManager("1", "test", dataForManager)
	Common.SendJsonMess(res, mes)
}
