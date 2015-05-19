package WebClient

import (
	"Config/Config"
	"Config/ConfigHttp"
	"Cron/CronMessage"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

/*

GET 	| 200 (OK), single customer. 404 (Not Found), if ID not found or invalid.
PUT 	| 200 (OK) or 204 (No Content). 404 (Not Found), if ID not found or invalid.
POST 	| 201 (Created), 'Location' header with link to /customers/{id} containing new ID. 	404 (Not Found).
DELETE 	| 200 (OK). 404 (Not Found), if ID not found or invalid.

*/

type Client struct {
	WG         sync.WaitGroup
	HttpConfig *ConfigHttp.Config
	Config     *Config.Config
	ChannelIn  chan CronMessage.Mess
	CronOut    chan CronMessage.Mess
	Transport  *http.Transport
}

func (client *Client) Run() {
	go func() {
		client._Run()
	}()
	log.Printf("WebClient is working\n")
}

func (client *Client) _Run() {

	for {
		select {
		// Read messages from WEB SERVER
		case mes, isGood := <-client.ChannelIn:

			log.Printf("client _Run: %s\n", mes.Data)

			if !isGood {
				break
			}

			switch mes.Type {
			case "stop":
				break
			default:
				go client.oneMessage(mes)
			}
		}
	}
}

func (client *Client) oneMessage(mes CronMessage.Mess) {
	switch mes.Type {
	case "all_server":
		servers := client.Config.ServersList()
		for _, s := range servers {
			if client.Config.ServerID != s.ID {
				continue
			}
			host := s.IP
			if s.Host != "" {
				host = s.Host
			}
			url := fmt.Sprintf("http://%s:%d/api/server/config", host, s.Port)
			go func(url string, data []byte) {
				client.SendPost(url, data)
			}(url, mes.Data)
		}
		//res, err = manager.saveScript(mes)
	default:
		//res["status"] = "INTERNAL ERROR"
		//err = fmt.Errorf("INTERNAL ERROR. NOT FOUND MESSAGE TYPR '%s'", mes.Type)
	}
}

func (client *Client) SendPost(url string, data []byte) {
	log.Printf("SendPost: %s => %s\n", url, data)
}

func LoadRemoutConfig(host string) ([]byte, error) {
	url := host + "/api/get/server/config"
	log.Printf("LoadUrl: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func New(config *Config.Config, httpConfig *ConfigHttp.Config) (*Client, error) {

	log.Printf("config: %+v\n", config)

	var transport *http.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &Client{
		Transport:  transport,
		Config:     config,
		HttpConfig: httpConfig,
		ChannelIn:  CronMessage.Channel(),
	}

	return client, nil
}
