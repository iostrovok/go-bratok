package WebServer

import (
	"Config/ConfigHttp"
	"Cron/CronMessage"
	"Web/Handlers/Common"
	"Web/Handlers/NotFound"
	"Web/Handlers/ScriptsEdit"
	"Web/Handlers/Servers"
	"Web/Handlers/StartScript"
	"Web/Handlers/StaticFiles"
	"Web/Handlers/Test"
	"fmt"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
	_ "net/http/pprof"
)

/*

GET 	| 200 (OK), single customer. 404 (Not Found), if ID not found or invalid.
PUT 	| 200 (OK) or 204 (No Content). 404 (Not Found), if ID not found or invalid.
POST 	| 201 (Created), 'Location' header with link to /customers/{id} containing new ID. 	404 (Not Found).
DELETE 	| 200 (OK). 404 (Not Found), if ID not found or invalid.

*/

type Server struct {
	Pat     *pat.PatternServeMux
	Config  *ConfigHttp.Config
	Channel chan CronMessage.Mess
}

func (server *Server) Listen() {
	go func() {
		log.Printf("%s", http.ListenAndServe(fmt.Sprintf(":%d", server.Config.Port()), nil))
	}()

	log.Printf("Listen on %d", server.Config.Port())
}

func New(config *ConfigHttp.Config) (*Server, error) {

	log.Printf("config: %+v\n", config)

	//http.Handle("/static", http.FileServer(http.Dir(config.StaticDir())))
	//http.Handle("/static/", http.FileServer(http.Dir(config.StaticDir())))

	webChannel := CronMessage.Channel()

	Common.SetConfig(config)
	Common.SetWebChannel(webChannel)

	_pat := pat.New()

	server := &Server{
		Pat:     _pat,
		Config:  config,
		Channel: webChannel,
	}

	//staticPath := config.StaticUrl()

	POST_listID := map[string]interface{}{}

	POST_list := map[string]interface{}{
		"/404":               NotFound.Handler,
		"/api/test":          Test.Handler,
		"/api/script/save":   ScriptsEdit.Save,
		"/api/server/config": Servers.ConfigUpdate,
	}

	GET_list := map[string]interface{}{
		"/404":                   NotFound.Handler,
		"/api/test":              Test.Handler,
		"/api/script_list":       ScriptsEdit.All,
		"/api/script/info":       ScriptsEdit.One,
		"/api/get/server/config": Servers.Config,
		"/api/server_list":       Servers.All,
	}

	GET_listID := map[string]interface{}{
		"/api/start":       StartScript.Handler,
		"/api/script/log":  ScriptsEdit.Log,
		"/api/script/info": ScriptsEdit.One,
		"/api/server/info": Servers.One,
		"/api/script/save": ScriptsEdit.Save,
	}

	PUT_listID := map[string]interface{}{}
	DEL_listID := map[string]interface{}{}

	for _, staticPath := range []string{"/html", "/css", "/javascript"} {
		for _, path := range []string{"", "/*", "/*/*", "/*/*/*", "/*/*/*/*", "/*/*/*/*/*"} {
			GET_list[staticPath+path] = StaticFiles.Handler
			GET_listID[staticPath+path] = StaticFiles.Handler
		}
	}

	http.Handle("/", _pat)

	if err := server.init_handlers("PUT", _pat, PUT_listID, true); err != nil {
		return nil, err
	}
	if err := server.init_handlers("DEL", _pat, DEL_listID, true); err != nil {
		return nil, err
	}

	if err := server.init_handlers("GET", _pat, GET_listID, true); err != nil {
		return nil, err
	}
	if err := server.init_handlers("POST", _pat, POST_listID, true); err != nil {
		return nil, err
	}

	if err := server.init_handlers("GET", _pat, GET_list, false); err != nil {
		return nil, err
	}
	if err := server.init_handlers("POST", _pat, POST_list, false); err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) init_handlers(t string, m *pat.PatternServeMux, list map[string]interface{}, isID bool) error {
	for path, fun := range list {

		if isID {
			path += "/:id"
		}

		log.Printf("Init [%s] %s\n", t, path)
		switch t {
		case "POST":
			m.Post(path, http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
			m.Post(path+"/", http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
		case "GET":
			m.Get(path, http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
			m.Get(path+"/", http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
		case "PUT":
			m.Put(path, http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
			m.Put(path+"/", http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
		case "DEL":
			m.Del(path, http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
			m.Del(path+"/", http.HandlerFunc(fun.(func(http.ResponseWriter, *http.Request))))
		default:
			return fmt.Errorf("Bad HTTP request type: %s", t)
		}
	}
	return nil
}
