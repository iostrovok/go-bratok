package main

import (
	"Config/Config"
	"Config/ConfigHttp"
	"Config/CronScript"
	"Config/ReadFlags"
	"Logger/Logger"
	"Manager/Manager"
	"Web/Handlers/Common"
	"Web/WebClient"
	"Web/WebServer"
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {

	flags := ReadFlags.New()
	err := flags.Read()
	if err != nil {
		fmt.Printf("flags.Read: %s\n", err)
		os.Exit(1)
	}

	runtime.GOMAXPROCS(MaxParallelism())

	/* Do we have to load a remote config? Loading... */
	RemoutConfig := []byte{}
	if flags.ConfigHost != "" {
		RemoutConfig, err = WebClient.LoadRemoutConfig(flags.ConfigHost)
		if err != nil {
			fmt.Printf("We have load remout config [%s] error:: %s\n", flags.ConfigHost, err)
			os.Exit(1)
		}
		fmt.Printf("RemoutConfig: %s\n", RemoutConfig)
		os.Exit(0)
	}

	cronConfig := Config.New(flags)
	if cronConfig.ErrorLoad != nil {
		fmt.Printf("cronConfig.ErrorLoad: %s\n", cronConfig.ErrorLoad)
		os.Exit(1)
	}
	Common.SetCronConfig(cronConfig)

	logger := Logger.New(cronConfig)
	logger.ToLog("start")

	httpConfig := ConfigHttp.New(flags, cronConfig)

	http, errWebServer := WebServer.New(httpConfig)
	if errWebServer != nil {
		log.Panicf("%s", errWebServer)
	}
	http.Listen()

	httpClient, errWebClient := WebClient.New(cronConfig, httpConfig)
	if errWebClient != nil {
		log.Panicf("%s", errWebClient)
	}

	/* TEST POINT START */
	script := CronScript.New("fromRun", "ls", "-a", "-r", "/")
	script.SetTime("*/1", "*", "*", "*", "*")
	cronConfig.AddScript(script)
	/* TEST POINT FINISH */

	manager := Manager.New(cronConfig)
	manager.SetLogger(logger)
	manager.Live(http, httpClient)

	manager.Wg.Wait()
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
