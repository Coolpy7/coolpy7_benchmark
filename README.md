# coolpy7_benchmark

This is a simple MQTT benchmark tool written in Golang. The main purpose of the tool is to benchmark how many concurrent connections a MQTT broker could support.

## Sub Benchmark

```
$ ./go_build_cp7_bench_sub_go_linux -h
Usage of ./go_build_cp7_bench_sub_go_linux:
  -cid string
    	client id start with (default "testclient")
  -clear
    	clear session (default true)
  -keepalive string
    	keepalive (default "300s")
  -qos uint
    	sub qos level
  -topic string
    	the used topic (default "cp7sub%i")
  -url string
    	broker url (default "tcp://127.0.0.1:1883")
  -workers int
    	number of workers (default 100)

  -h                 help information
  -cid               client id start with this value profix + workers id [like: testclient0]
  -clear             clean session [default: true]
  -keepalive         keep alive in seconds [default: 300]
  -qos               subscribe qos [default: 0]
  -topic             topic subscribe, support %i variables
  -url               mqtt connect string [like: tcp://username:password@192.168.100.2:1883]
  -workers           mqtt connection clients count [default: 100]
```

## Pub Benchmark

```
$ ./go_build_cp7_bench_pub_go_linux -h
Usage of ./go_build_cp7_bench_pub_go_linux:
  -I int
    	interval of publishing message(ms) (default 1000)
  -cid string
    	client id start with (default "client")
  -clear
    	clear session (default true)
  -i int
    	interval of connecting to the broker(ms) (default 10)
  -keepalive string
    	keepalive (default "300s")
  -qos uint
    	pub qos level
  -s int
    	payload size (default 256)
  -topic string
    	pub topic (default "cp7sub%i")
  -url string
    	broker url (default "tcp://username:password@192.168.100.2:1883")
  -workers int
    	number of workers (default 200)

  -h                 help information
  -cid               client id start with this value profix + workers id [like: client0]
  -clear             clean session [default: true]
  -keepalive         keep alive in seconds [default: 300]
  -qos               subscribe qos [default: 0]
  -topic             topic subscribe, support %i variables
  -url               mqtt connect string [like: tcp://username:password@192.168.100.2:1883]
  -workers           mqtt connection clients count [default: 100]
  -s                 payload size [default: 256]
  -I                 interval of publishing message(ms) [default 1000]
  -i                 interval of connecting to the broker(ms) [default 10]
```
