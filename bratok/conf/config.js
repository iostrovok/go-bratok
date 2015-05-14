{
    "config_id": 123123,
    "staticdir": "/Users/ostrovok/Work/bratok-js/public-test/",
    "logfile": "bratok.scripts.log",
    "logdir": "/tmp/",
    "scripts": [{
        "id": "ls1",
        "time": ["*/1 * * * *", "*/5 * * * *"],
        "exe": "ls",
        "params": ["-a", "-r", "./"],
        "env": ["export DBUSER=boomboom", "export DBPASS=boomboom", "export DBBASE=somethere"]
    }, {
        "id": "ls2",
        "time": ["*/2"],
        "exe": "sh",
        "params": ["ls", "-a", "-r", "./"],
        "env": []
    }, {
        "id": "longperl",
        "time": [],
        "exe": "/Users/ostrovok/Work/go-bratok/test_long_time.pl",
        "params": ["10"],
        "env": []
    }],
    "servers": [{
        "id": "first",
        "ip": "127.0.0.1",
        "host": "",
        "port": 21222,
        "is_master": true,
        "scripts": ["ls1", "longperl", "ls2"]
    }, {
        "id": "second",
        "ip": "127.0.0.1",
        "host": "",
        "port": 21223,
        "is_master": false,
        "scripts": ["ls1" ]
    }, {
        "id": "third",
        "ip": "127.0.0.1",
        "host": "",
        "port": 21224,
        "is_master": false,
        "scripts": ["ls2"]
    }]
}
