package NotFound

import (
	"net/http"
)

func Handler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "NOT FOUND", 404)
}
