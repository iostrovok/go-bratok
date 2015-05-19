package Manager

import (
	"Cron/CronMessage"
	"Web/WebServer"
	"fmt"
)
import (
	"log"
)

func (manager *Manager) Rules(server *WebServer.Server) {

	for {
		select {
		// Read messages from WEB SERVER
		case mes, isGood := <-server.Channel:
			manager.TxtLog("Manager-live. *WebServer.Server isGood %t, mes: %+v\n", isGood, mes)

			if !isGood {
				break
			}

			res := map[string]interface{}{}
			var err error = nil

			switch mes.Type {
			case "config_update":
				res, err = manager.updateConfig(mes)
			case "save_script":
				res, err = manager.saveScript(mes)
			case "start_script":
				err = manager.startScriptById(mes.ID)
				log.Printf("Live. start_script. %s\n", err)
				res["status"] = "Try to Start"
			case "script_log":
				res = manager.logScriptById(mes.ID)
			case "script_info":
				res, err = manager.infoScriptById(mes.ID)
			case "server_info":
				res, err = manager.infoServerById(mes.ID)
			case "save_server":
				res, err = manager.saveServer(mes)
			case "scripts_list":
				res, err = manager.scriptsList(mes)
			case "server_list":
				res, err = manager.serversList(mes)
			case "test":
				res["status"] = "DATA FROM MENEGER"
			default:
				res["status"] = "INTERNAL ERROR"
				err = fmt.Errorf("INTERNAL ERROR. NOT FOUND MESSAGE TYPR '%s'", mes.Type)
			}
			if err != nil {
				log.Printf("MANAGER ERROR: %s\n", err)
			}
			if mes.BackChannel != nil {
				CronMessage.SendHash(mes.BackChannel, mes.ID, mes.Type, res, err)
				close(mes.BackChannel)
			}
		// Timer action
		case _, isGood := <-manager.chTimerFrom:

			if !isGood {
				break
			}

			manager.TimerAction()

		case mes, isGood := <-manager.chCronFrom:

			if !isGood {
				manager.TxtLog("manager.chCronFrom is close")
				break
			}

			manager.FixCronMessage(mes)

		case mes, isGood := <-manager.chMainIn:
			manager.TxtLog("Manager-live. chMainIn m: %+v\n", mes)
			manager.TxtLog("Manager-live. chMainIn isGood: %+v\n", isGood)
		}
	}
}
