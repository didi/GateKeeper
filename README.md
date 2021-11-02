# 快速开始（Quick Start）

## 环境要求 >= go 1.12.0
_源码环境要求_


## 环境安装

[GO环境安装](https://golang.google.cn/dl/) ||  [教程](https://www.runoob.com/go/go-environment.html) *源码环境需要*

[Mysql安装](https://dev.mysql.com/downloads/mysql/) || [教程](https://www.runoob.com/mysql/mysql-install.html) *必须*

[Redis安装](https://github.com/tporadowski/redis/releases) || [教程](https://www.runoob.com/redis/redis-install.html) *必须*

---

> ## Gatekeeper安装方式(二选一)

### 1. 源码安装
```
# step 1  get gatekeeper server 
git clone https://github.com/didi/GateKeeper.git

# or download zip
# wget https://hub.fastgit.org/didi/GateKeeper/archive/refs/heads/master.zip

# step 2  set go proxy & download library dependency
export GO111MODULE=on && export GOPROXY=https://goproxy.cn
cd gatekeeper
go mod tidy

# step 3 init config
# step 4 run gatekeeper
```

### 2.二进制文件安装

```
# step 1  get release gatekeeper server 
curl -L 'https://github.com/didi/GateKeeper/releases/download/v1.0.0/setup.sh' | sh
# step 2  run gatekeeper
# setup already automatic init config
```

---

> ## 服务配置初始化(二选一)

_源码环境&&首次运行时需要_

[config详解](https://github.com/didi/Gatekeeper/blob/master/doc/config/README.md)

### 1. 自动初始化
```
# run auto init shell 
cd install && sh install.sh

# input redis connect info 
# input mysql connect info
# auto create gatekeeper config
```

### 2.手动初始化
```
# step 1 create database & input data
mysql -h 127.0.0.1 -u root -p -e "CREATE DATABASE gatekeeper DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;"
mysql -h 127.0.0.1 -u root -p gatekeeper < gatekeeper.sql --default-character-set=utf8

# step 2 edit mysql connect info
vim ./conf/dev/mysql_map.toml
# set data_source_name = "{dbUser}:{dbPwd}@tcp({dbIp}:{dbPort})/{dbDatabase}?charset=utf8&parseTime=true&loc=Asia%2FChongqin" 

# step 3 edit redis connect info
vim ./conf/dev/redis_map.toml
# set proxy_list = "redisIP:redsiPort"
# set password = "redisPwd"

# step 4 edit base
vim ./conf/dev/base.toml
# set session redis_service = "redsiIp:redisPort"
# set session redis_password = "redisPwd"
```
---
```math
E = mc^2
```


> ## 运行服务

### 运行管理面板

```
sh control.sh start_control 
```

### 运行代理服务

```
sh control.sh start_proxy 
```

---

# 自动化集成测试

### [自动化集成测试](https://github.com/didi/Gatekeeper/blob/master/test_suites/README.md)