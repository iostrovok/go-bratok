GP := $(shell dirname $(realpath $(lastword $(GOPATH))))
ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH := ${ROOT}/bratok/:${GOPATH}

all:
	touch /Users/ostrovok/Work/go-bratok/bratok/cover.html 
	rm /Users/ostrovok/Work/go-bratok/bratok/cover.html
#export GOROOT := /usr/local/go/

#.PHONY: all test build index import run

# test:
# 	echo ${GOPATH}
# 	cd ./bratok; go test ./src/*/*/*/
# 	cd ./bratok; go test ./src/*/*/ 
# 	cd ./bratok; go test ./src/*/

run:
	 go run ./bratok/src/scripts/run.go -id=first -config="./bratok/conf/config.js"

run2:
	 go run ./bratok/src/scripts/run.go -id=second -master_host="http://127.0.0.1:21222"

run3:
	 go run ./bratok/src/scripts/run.go -id=third -master_host="http://127.0.0.1:21222"

log:
	tail -f /tmp/bratok.scripts.log

test: all test-butils test-flags test-hmem test-conf-file test-conf-http test-conf-load test-conf test-cron test-cronscheduler \
	test-cronscript test-h-common test-manager test-message test-net test-net test-server test-task \
	test-webserver test-timer


test-butils: all
	cd ./bratok; go test ./src/Butils/

test-butils-cover:
	cd ./bratok; go test ./src/Butils/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-manager: all
	cd ./bratok; go test ./src/Manager/Manager/

test-manager-cover:
	cd ./bratok; go test ./src/Manager/Manager/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out

test-server: all
	cd ./bratok; go test ./src/Net/Server/

test-server-cover:
	cd ./bratok; go test ./src/Net/Server/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-cron: all
	cd ./bratok; go test ./src/Cron/Cron/

test-cron-cover:
	cd ./bratok; go test ./src/Cron/Cron/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-task: all
	cd ./bratok; go test ./src/Cron/CronTask/

test-task-cover:
	cd ./bratok; go test ./src/Cron/CronTask/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-timer: all
	cd ./bratok; go test ./src/Cron/CronTimer/

test-timer-cover:
	cd ./bratok; go test ./src/Cron/CronTimer/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-message: all
	cd ./bratok; go test ./src/Cron/CronMessage/

test-message-cover:
	cd ./bratok; go test ./src/Cron/CronMessage/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out



test-net: all
	cd ./bratok; go test ./src/Net/HTTPLoader/

test-net-cover:
	cd ./bratok; go test ./src/Net/HTTPLoader/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out

test-conf: all
	cd ./bratok; go test ./src/Config/Config/
	
test-conf-cover:
	cd ./bratok; go test ./src/Config/Config/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out



test-cronscheduler: all
	cd ./bratok; go test ./src/Config/CronScheduler/
	
test-cronscheduler-cover:
	cd ./bratok; go test ./src/Config/CronScheduler/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-cronscript: all
	cd ./bratok; go test ./src/Config/CronScript/
	
test-cronscript-cover:
	cd ./bratok; go test ./src/Config/CronScript/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-conf-http: all
	cd ./bratok; go test ./src/Config/ConfigHttp/
	
test-conf-http-cover:
	cd ./bratok; go test ./src/Config/ConfigHttp/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-webserver: all
	cd ./bratok; go test ./src/Web/WebServer/
	
test-webserver-cover:
	cd ./bratok; go test ./src/Web/WebServer/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-h-common: all
	cd ./bratok; go test ./src/Web/Handlers/Common/
	
test-h-common-cover:
	cd ./bratok; go test ./src/Web/Handlers/Common/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out

# test-h-staticfiles:
# 	cd ./bratok; go test ./src/Web/Handlers/StaticFiles/
	
# test-h-staticfiles-cover:
# 	cd ./bratok; go test ./src/Web/Handlers/StaticFiles/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out

test-flags: all
	cd ./bratok; go test ./src/Config/ReadFlags/
	
test-flags-cover:
	cd ./bratok; go test ./src/Config/ReadFlags/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-hmem: all
	cd ./bratok; go test ./src/HistoryMem/HistoryMem/
	
test-hmem-cover:
	cd ./bratok; go test ./src/HistoryMem/HistoryMem/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-logger: all
	cd ./bratok; go test ./src/Logger/Logger/
	
test-logger-cover:
	cd ./bratok; go test ./src/Logger/Logger/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-conf-load: all
	cd ./bratok; go test ./src/Config/LoaderConfig/
	
test-conf-load-cover:
	cd ./bratok; go test ./src/Config/LoaderConfig/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out

test-conf-file: all
	cd ./bratok; go test ./src/Config/File/
	
test-conf-file-cover:
	cd ./bratok; go test ./src/Config/File/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out


test-conf-hist: all
	cd ./bratok; go test ./src/Config/History/
	
test-conf-hist-cover:
	cd ./bratok; go test ./src/Config/History/ -cover -coverprofile ./tmp.out; go tool cover -html=./tmp.out -o cover.html; rm ./tmp.out



