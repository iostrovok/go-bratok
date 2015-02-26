package Common

import (
	"Config/Config"
	"Config/ConfigHttp"
	"Cron/CronMessage"
	"encoding/json"
	//"errors"
	"fmt"
	"net/http"
)

var myConfig *Config.Config
var myConfigHttp *ConfigHttp.Config
var webChannel chan CronMessage.Mess

func SetCronConfig(c *Config.Config) {
	myConfig = c
}
func GetCronConfig() *Config.Config {
	return myConfig
}
func GetConfigData() *Config.ConfigData {
	return myConfig.ConfigData
}

func ToHttpLog(str string) {
	fmt.Printf("ToHttpLog: %s\n", str)
}

func SetWebChannel(c chan CronMessage.Mess) {
	webChannel = c
}

func GetWebChannel() chan CronMessage.Mess {
	return webChannel
}

func SetConfig(Config *ConfigHttp.Config) {
	myConfigHttp = Config
}
func GetStaticDir() string {
	return myConfigHttp.StaticDir()
}
func GetStaticUrl() string {
	return myConfigHttp.StaticUrl()
}

const (
	InternalJsonError string = `{"error":"Internal Erroror. Internal code 1.","data":null,"result":1}`
)

type jsResult struct {
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
	Result int         `json:"result"`
}

func NewResultEmpty() *jsResult {
	return &jsResult{
		Error:  "",
		Data:   nil,
		Result: 1,
	}
}

func NewResult(err error, d ...interface{}) *jsResult {

	var data interface{} = nil
	if len(d) > 0 {
		data = d[0]
	}

	out := &jsResult{
		Error:  "",
		Data:   data,
		Result: 1,
	}

	if err != nil {
		out.Error = fmt.Sprintf("%s", err)
		out.Result = 0
	}

	return out
}

func JsonResult(err error, d interface{}) []byte {
	if d, err := NewResult(err, d).ToJson(); err == nil {
		return d
	}

	// TODO log error
	return []byte(InternalJsonError)
}

func (result *jsResult) ToJsonStr() string {
	if d, err := result.ToJson(); err == nil {
		return string(d)
	}

	return InternalJsonError
}

func (result *jsResult) ToJson() ([]byte, error) {
	return json.Marshal(result)
}

func ErrorPage(res http.ResponseWriter, err error) {
	data := map[string]interface{}{
		"error": err,
	}
	b := JsonResult(err, data)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(b)
}

// SendJson - Silple REST API with body
func SendJsonMess(res http.ResponseWriter, mes *CronMessage.Mess) {

	// Not
	if mes == nil {
		SendJsonError(500, res, nil)
		return
	}

	if mes.Error != nil {
		SendJsonError(404, res, mes.Error)
		return
	}

	b := mes.Hash
	if len(b) == 0 {
		b = map[string]interface{}{"data": mes.Data}
	}
	b["time"] = mes.Time

	SendJsonSuccess(res, b)
}

func SendJsonSuccess(res http.ResponseWriter, d ...interface{}) {
	var data interface{} = nil
	if len(d) > 0 {
		data = d[0]
	}
	b := JsonResult(nil, data)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(b)
}

func SendJsonError(code int, res http.ResponseWriter, err error) {
	b := JsonResult(err, map[string]interface{}{})
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	http.Error(res, string(b), code)
	res.WriteHeader(code)
	res.Write(b)
	//	fmt.Fprintln(res, error)
}

func ToFromManager(id, mtype string, dataForManager map[string]interface{}) *CronMessage.Mess {

	myCh := CronMessage.Channel()
	mess := CronMessage.NewMess(id, mtype, nil, dataForManager, nil, myCh)

	GetWebChannel() <- *mess

	select {
	case mes, isGood := <-myCh:
		if isGood {
			return &mes
		}
	}
	return nil
}
