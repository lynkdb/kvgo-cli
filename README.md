# kvgo-cli

kvgo-cli is the command line tool to connect and manage the [kvgo-server](https://github.com/lynkdb/kvgo-server).


## Installing

Ensure that the golang development environment has been installed correctly

``` shell
go install github.com/lynkdb/kvgo-cli/cmd/kvgo-cli
```

## Usage


### 1. first instance setup

``` shell
[guest@VM-0-17-centos ~]$ kvgo-cli
no instance setup in /home/guest/.kvgo-cli.conf, try to use 'instance new' to create new connection to kvgo-server
kvgo-cli : instance new
input alias name of instance (ex: prod, demo, ...) : demo
input instace address (ex. 127.0.0.1:9200) : 10.0.0.1:9200
input access key id : 00000000
input access key secret:  : the-access-key-secret-of-server
Use Instance demo
kvgo-cli (demo):
```

### 2. use [table list] to check if the configuration is correct

``` shell
kvgo-cli (demo):  table list

ID  | Name   |  Keys |   Size | Log      | Incr      | Async                          | Desc   | Created   
10  | main   |     0 |      0 |          |           |                                |        | 2020-08-03
101 | zone   | 19165 |  24 MB | 37860944 |           |                                | zone   | 2020-09-03
102 | global |  1412 | 440 KB | 67427400 | meta 1100 | 10.0.0.10:1111:global 45229880 | global | 2020-09-03
    |        |       |        |          | role 1100 | 10.0.0.20:2222:main 50344506   |        |           
103 | inpack |   716 |   1 GB | 5684     |           | 10.0.0.10:1111:inpack 4889     | inpack | 2020-09-03
    |        |       |        |          |           | 10.0.0.20:2222:inpack 4889     |        |           
```

### 3. use [help] to list more methods

``` shell
kvgo-cli (demo):  help

kvgo-cli usage:
  instance list
  instance use <name>
  instance new
  table list
  table set
  access key list
  access key get
  access key set
  help
  quit
```

## Dependent or referenced

* kvgo [https://github.com/lynkdb/kvgo](https://github.com/lynkdb/kvgo)
* kvgo-server [https://github.com/lynkdb/kvgo-server](https://github.com/lynkdb/kvgo-server)
