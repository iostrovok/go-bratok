package StaticFiles

import (
	"Web/Handlers/Common"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}

	Common.ToHttpLog(upath)

	name := Common.GetStaticDir() + strings.TrimPrefix(upath, Common.GetStaticUrl())
	Common.ToHttpLog(name)

	http.ServeFile(w, r, name)
}
