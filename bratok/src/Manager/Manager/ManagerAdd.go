package Manager

import (
	"BUtils"
	"Config/CronScheduler"
	"Cron/CronMessage"
	"fmt"
)
import (
	"log"
)

func (manager *Manager) updateConfigAllServers() error {
	conf, err := manager.config.GetConfigDataByte()
	if err == nil {
		CronMessage.Send(manager.httpClientChannel, "1", "all_server", conf, nil)
	}
	return err
}

func (manager *Manager) saveScript(mes CronMessage.Mess) (map[string]interface{}, error) {

	log.Printf("saveScript. mes: %+v\n", mes)

	// mes.ID - script id
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

	log.Printf("saveScript. times. v: %+v\n", times)

	script.CleanTime()

	log.Println("CleanTime. Finish")

	for _, t := range BUtils.AnyTo2StringArray(times) {
		for i := range t {
			log.Printf("CleanTime 0: %s\n", t[i])

			t[i] = CronScheduler.CleanTimeLinePoint(t[i])
			log.Printf("CleanTime 1: %s\n", t[i])

		}
		if err := script.SetTime(t...); err != nil {
			return map[string]interface{}{}, err
		}
	}

	script.SetParams(BUtils.AnyToStringArray(BUtils.GetPath(mes.Hash, "params")))
	script.SetEnv(BUtils.AnyToStringArray(BUtils.GetPath(mes.Hash, "env")))
	script.SetExe(BUtils.AnyToString(BUtils.GetPath(mes.Hash, "exe")))

	log.Printf("Set new script: %s\n", script)

	manager.config.RaplaceScript(script)
	manager.config.Store()

	return map[string]interface{}{}, manager.updateConfigAllServers()
}

func (manager *Manager) scriptsList(mes CronMessage.Mess) (map[string]interface{}, error) {

	// mes.ID - server id
	list := manager.config.ScriptsList(mes.ID)

	return map[string]interface{}{
		"list": list,
	}, nil
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
