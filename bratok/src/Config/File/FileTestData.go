package File

import (
	"fmt"
)

const (
	httpFormat string = `{"error":"","data":%s,"history":%s,"result":1}`
)

func TestFileLineHttp() []byte {
	//out := `{"error":"","data":{"Data":` + string(TestFileLine()) + `},"result":1}`
	out := fmt.Sprintf(httpFormat, TestFileLine(), TestHistoryLine())
	return []byte(out)
}

func TestFileLineHttpNoServer() []byte {
	out := fmt.Sprintf(httpFormat, TestFileLineNoServer())
	//out := `{"error":"","data":` + string(TestFileLineNoServer()) + `,"result":1}`
	return []byte(out)
}

func TestFileLine() []byte {
	return []byte(`{
			"logfile":"FILE","logdir":"DIR","staticdir":"STATIC-DIR","file_id":12312312,
			"scripts":[
				{"id":"ls22","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls33","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls1","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls2","time":["*/2"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]}
			],
			"servers":[
				{
					"id":"workstation","ip":"127.0.0.1","host":"","port":21222,"is_master":true,"scripts": ["ls2","ls1"],
					"logfile":"workstation_FILE","logdir":"workstation_DIR","staticdir":"workstation_STATIC-DIR"
				},
				{"id":"somethere","ip":"192.168.0.10","host":"wks-l","port":21223,"is_master":false,"scripts": ["ls2"]}
			]
	}`)
}

func TestHistoryLine() []byte {
	return []byte(`{
			"line":[{"id":"jEp\u0014o|��b\"�{�X�p","prev_id":"","time":"2015-05-26T18:15:05.750142396+03:00","script":{"id":"longperl","time":["*/20","*/5 * * * *"],"exe":"/Users/ostrovok/Work/go-bratok/test_long_time.pl","params":["10"],"env":[]},"server":null,"server_id":"","act":"replace"},{"id":"}\u0019g�-�ݻ�y�\r�P\u001b�","prev_id":"jEp\u0014o|��b\"�{�X�p","time":"2015-05-27T08:57:05.267868168+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\u000c\u0005�E.\u0001�4H��Դ��.","prev_id":"}\u0019g�-�ݻ�y�\r�P\u001b�","time":"2015-05-27T08:59:54.571068906+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"J\ufffd\tANA\ufffd:\ufffd\u0004\n\ufffd\u0001-v\u0002","prev_id":"\u000c\u0005�E.\u0001�4H��Դ��.","time":"2015-05-27T09:00:00.527895765+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffd@4\ufffd\ufffd$\u000f\ufffd\ufffd\u0026\ufffd%)fŧ","prev_id":"J\ufffd\tANA\ufffd:\ufffd\u0004\n\ufffd\u0001-v\u0002","time":"2015-05-27T09:00:17.319823344+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffd\ufffd\u0026^ْ\ufffd2|\ufffd\r\ufffd˸\ufffd\ufffd","prev_id":"\ufffd@4\ufffd\ufffd$\u000f\ufffd\ufffd\u0026\ufffd%)fŧ","time":"2015-05-27T09:00:30.877613022+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffd$K\ufffd\ufffd\ufffdW\ufffd\ufffd^u\ufffdu\ufffd\ufffdb","prev_id":"\ufffd\ufffd\u0026^ْ\ufffd2|\ufffd\r\ufffd˸\ufffd\ufffd","time":"2015-05-27T09:00:46.627640173+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"8e-5}4'd\ufffd\ufffdj\ufffd\ufffdP\ufffd\ufffd","prev_id":"\ufffd$K\ufffd\ufffd\ufffdW\ufffd\ufffd^u\ufffdu\ufffd\ufffdb","time":"2015-05-27T09:04:20.67246623+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21225,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffdH2z\ufffd^\ufffd\ufffd9\ufffd\ufffd\ufffd\ufffd\u0026E-","prev_id":"8e-5}4'd\ufffd\ufffdj\ufffd\ufffdP\ufffd\ufffd","time":"2015-05-27T09:05:28.022454916+03:00","script":null,"server":{"id":"third","ip":"127.0.0.2","host":"localhost","port":21226,"is_master":true,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"q\ufffdY^\ufffdc\ufffd\ufffdU\ufffd/\u0018m\ufffd+","prev_id":"\ufffdH2z\ufffd^\ufffd\ufffd9\ufffd\ufffd\ufffd\ufffd\u0026E-","time":"2015-05-27T09:05:39.718491451+03:00","script":null,"server":{"id":"third","ip":"127.0.0.2","host":"localhost","port":21226,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"e\ufffd\ufffd+ \ufffd\ufffd\ufffd\ufffdɣ\ufffdr\\\ufffd\ufffd","prev_id":"q\ufffdY^\ufffdc\ufffd\ufffdU\ufffd/\u0018m\ufffd+","time":"2015-05-27T09:05:44.23161871+03:00","script":null,"server":{"id":"third","ip":"127.0.0.2","host":"localhost","port":21226,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"}]
	}`)
}

func TestFileLineNoServer() []byte {
	return []byte(`{
			"logfile":"FILE","logdir":"DIR","staticdir":"STATIC-DIR","file_id":12312312,
			"scripts":[
				{"id":"ls22","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls33","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls1","time":["*/1 * * * *"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]},
				{"id":"ls2","time":["*/2"],"exe":"sh","params":["ls","-a","-r","./"],"env":[]}
			],
			"history":{"line":[{"id":"jEp\u0014o|��b\"�{�X�p","prev_id":"","time":"2015-05-26T18:15:05.750142396+03:00","script":{"id":"longperl","time":["*/20","*/5 * * * *"],"exe":"/Users/ostrovok/Work/go-bratok/test_long_time.pl","params":["10"],"env":[]},"server":null,"server_id":"","act":"replace"},{"id":"}\u0019g�-�ݻ�y�\r�P\u001b�","prev_id":"jEp\u0014o|��b\"�{�X�p","time":"2015-05-27T08:57:05.267868168+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\u000c\u0005�E.\u0001�4H��Դ��.","prev_id":"}\u0019g�-�ݻ�y�\r�P\u001b�","time":"2015-05-27T08:59:54.571068906+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"J\ufffd\tANA\ufffd:\ufffd\u0004\n\ufffd\u0001-v\u0002","prev_id":"\u000c\u0005�E.\u0001�4H��Դ��.","time":"2015-05-27T09:00:00.527895765+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffd@4\ufffd\ufffd$\u000f\ufffd\ufffd\u0026\ufffd%)fŧ","prev_id":"J\ufffd\tANA\ufffd:\ufffd\u0004\n\ufffd\u0001-v\u0002","time":"2015-05-27T09:00:17.319823344+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffd\ufffd\u0026^ْ\ufffd2|\ufffd\r\ufffd˸\ufffd\ufffd","prev_id":"\ufffd@4\ufffd\ufffd$\u000f\ufffd\ufffd\u0026\ufffd%)fŧ","time":"2015-05-27T09:00:30.877613022+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffd$K\ufffd\ufffd\ufffdW\ufffd\ufffd^u\ufffdu\ufffd\ufffdb","prev_id":"\ufffd\ufffd\u0026^ْ\ufffd2|\ufffd\r\ufffd˸\ufffd\ufffd","time":"2015-05-27T09:00:46.627640173+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21224,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"8e-5}4'd\ufffd\ufffdj\ufffd\ufffdP\ufffd\ufffd","prev_id":"\ufffd$K\ufffd\ufffd\ufffdW\ufffd\ufffd^u\ufffdu\ufffd\ufffdb","time":"2015-05-27T09:04:20.67246623+03:00","script":null,"server":{"id":"third","ip":"127.0.0.1","host":"","port":21225,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"\ufffdH2z\ufffd^\ufffd\ufffd9\ufffd\ufffd\ufffd\ufffd\u0026E-","prev_id":"8e-5}4'd\ufffd\ufffdj\ufffd\ufffdP\ufffd\ufffd","time":"2015-05-27T09:05:28.022454916+03:00","script":null,"server":{"id":"third","ip":"127.0.0.2","host":"localhost","port":21226,"is_master":true,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"q\ufffdY^\ufffdc\ufffd\ufffdU\ufffd/\u0018m\ufffd+","prev_id":"\ufffdH2z\ufffd^\ufffd\ufffd9\ufffd\ufffd\ufffd\ufffd\u0026E-","time":"2015-05-27T09:05:39.718491451+03:00","script":null,"server":{"id":"third","ip":"127.0.0.2","host":"localhost","port":21226,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"},{"id":"e\ufffd\ufffd+ \ufffd\ufffd\ufffd\ufffdɣ\ufffdr\\\ufffd\ufffd","prev_id":"q\ufffdY^\ufffdc\ufffd\ufffdU\ufffd/\u0018m\ufffd+","time":"2015-05-27T09:05:44.23161871+03:00","script":null,"server":{"id":"third","ip":"127.0.0.2","host":"localhost","port":21226,"is_master":false,"scripts":["ls1","ls2"],"staticdir":"","logfile":"","logdir":""},"server_id":"","act":"replace"}]}
	}`)
}
