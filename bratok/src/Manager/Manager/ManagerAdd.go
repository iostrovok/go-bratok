package Manager

import (
	"BUtils"
	"Config/CronScheduler"
	"Cron/CronMessage"
	"fmt"
	"log"
)

func (manager *Manager) StoreConfig() error {
	// save to file & set new ConfigId
	if err := manager.config.Store(); err != nil {
		return err
	}

	conf, err := manager.config.GetConfigDataByte()
	if err == nil {
		// send new config to all servers
		CronMessage.Send(manager.httpClientChannel, "1", "all_server", conf, nil)
	}
	return err
}

func (manager *Manager) updateConfig(mes CronMessage.Mess) (map[string]interface{}, error) {
	body := BUtils.GetPath(mes.Hash, "body")
	log.Printf("updateConfig body: %s\n", body)

	var err error
	if s, f := body.([]byte); f {
		err = manager.config.FullUpdate(s)
	} else {
		err = fmt.Errorf("cannot use body (type interface {}) as type []byte in argument to manager.config.FullUpdate: need type assertion")
	}
	return map[string]interface{}{
		"error": err,
	}, nil
}

func (manager *Manager) saveServer(mes CronMessage.Mess) (map[string]interface{}, error) {
	_, err := manager.infoServerById(mes.ID)
	if err != nil {
		return map[string]interface{}{}, err
	}

	serverOld, find := manager.config.GetServer(mes.ID)
	if !find {
		return map[string]interface{}{},
			fmt.Errorf("saveServer. Server [id =  %s] is not found", mes.ID)
	}

	server := serverOld.Clone()

	server.IP = BUtils.AnyToString(BUtils.GetPath(mes.Hash, "ip"))
	server.Host = BUtils.AnyToString(BUtils.GetPath(mes.Hash, "host"))
	server.Port = BUtils.AnyToInt(BUtils.GetPath(mes.Hash, "port"))
	server.Scripts = BUtils.AnyToStringArray(BUtils.GetPath(mes.Hash, "scripts"))
	server.IsMaster = BUtils.AnyToBool(BUtils.GetPath(mes.Hash, "is_master"))

	server.StaticFilesDir = BUtils.AnyToString(BUtils.GetPath(mes.Hash, "staticdir"))
	server.LogFile = BUtils.AnyToString(BUtils.GetPath(mes.Hash, "logfile"))

	log.Printf("Save/Add server: %s\n", server)

	manager.config.ReplaceServer(server)

	errorStore := manager.StoreConfig()
	return map[string]interface{}{}, errorStore
}

func (manager *Manager) saveScript(mes CronMessage.Mess) (map[string]interface{}, error) {
	_, err := manager.infoScriptById(mes.ID)
	if err != nil {
		return map[string]interface{}{}, err
	}

	scriptOld, find := manager.config.GetScript(mes.ID)
	if !find {
		return map[string]interface{}{},
			fmt.Errorf("saveScript. Script [id =  %s] is not found", mes.ID)
	}

	script := scriptOld.Clone()

	times := BUtils.GetPath(mes.Hash, "times")

	script.CleanTime()

	for _, t := range BUtils.AnyTo2StringArray(times) {
		for i := range t {
			t[i] = CronScheduler.CleanTimeLinePoint(t[i])
		}
		if err := script.SetTime(t...); err != nil {
			return map[string]interface{}{}, err
		}
	}

	script.SetParams(BUtils.AnyToStringArray(BUtils.GetPath(mes.Hash, "params")))
	script.SetEnv(BUtils.AnyToStringArray(BUtils.GetPath(mes.Hash, "env")))
	script.SetExe(BUtils.AnyToString(BUtils.GetPath(mes.Hash, "exe")))

	manager.config.ReplaceScript(script)

	errorStore := manager.StoreConfig()
	return map[string]interface{}{}, errorStore
}

func (manager *Manager) scriptsList(mes CronMessage.Mess) (map[string]interface{}, error) {

	// mes.ID - server id
	list := manager.config.ScriptsList(mes.ID)

	return map[string]interface{}{
		"list": list,
	}, nil
}

func (manager *Manager) serversList(mes CronMessage.Mess) (map[string]interface{}, error) {

	// mes.ID - server id
	list := manager.config.ServersList()

	return map[string]interface{}{
		"list": list,
	}, nil
}

func (manager *Manager) infoServerById(server_id string) (map[string]interface{}, error) {

	if server_id == "" {
		return nil, fmt.Errorf("infoServerById. Server id is empty")
	}

	server, find := manager.config.GetServer(server_id)
	if find {
		return map[string]interface{}{
				"data": server,
			},
			nil
	}

	return nil, fmt.Errorf("infoServerById. Server [id =  %s] is not found", server_id)

}

func (manager *Manager) infoScriptById(script_id string) (map[string]interface{}, error) {

	if script_id == "" {
		return nil, fmt.Errorf("infoScriptById. Script id is empty")
	}

	script, find := manager.config.GetScript(script_id)
	if find {
		return map[string]interface{}{
				"data": script,
			},
			nil
	}

	return nil, fmt.Errorf("infoScriptById. Script [id =  %s] is not found", script_id)

}
