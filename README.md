# 快速开始（Quick Start）

## 环境要求

软件 | 版本 | 官网 | 安装教程
---|---|---|---
mysql | \>=5.1 | [官网](https://dev.mysql.com/downloads/mysql/) | [教程](https://www.runoob.com/mysql/mysql-install.html)
redis | \>=3.0 | [官网](https://redis.io/download) | [教程](https://www.runoob.com/redis/redis-install.html)

## 服务安装(二选一)

类型 | 优点
---|---
二进制文件安装 | 适合无golang环境且无插件开发需求用户.
源码安装 | 支持插件定制化用户.


### 1. 二进制文件安装
目前支持64位windows/linux/mac 操作系统下用户，运行以下命令可自动化进行提示安装。  
windows用户需要启动powershell运行以下脚本。
```
bash <(curl -s -S -L 'https://download.fastgit.org/didi/GateKeeper/releases/download/v1.0.0/setup.sh')  
```
执行完毕后，会提示启动服务的命令，按照操作运行即可，执行二进制安装无需再进行源码安装。

### 2. 源码安装

使用源码安装，需依赖golang环境，要求如下：

软件 | 版本 | 官网 | 安装教程
---|---|---|---
golang | \>=1.12 | [官网](https://golang.google.cn/dl/) |  [教程](https://www.runoob.com/go/go-environment.html)

#### 2.1 git clone code

```
git clone https://github.com/didi/GateKeeper.git
```

#### 2.2 set go proxy & download golang dependency

```
export GO111MODULE=on && export GOPROXY=https://goproxy.cn
cd gatekeeper
go mod tidy
```

#### 2.3 edit mysql connect info

```
vim ./conf/dev/mysql_map.toml
# set data_source_name = "{dbUser}:{dbPwd}@tcp({dbIp}:{dbPort})/{dbDatabase}?charset=utf8&parseTime=true&loc=Asia%2FChongqin" 
```

#### 2.4 edit redis connect info


```
vim ./conf/dev/redis_map.toml
# set proxy_list = "redisIP:redsiPort"
# set password = "redisPwd"

vim ./conf/dev/base.toml
# set session redis_service = "redsiIp:redisPort"
# set session redis_password = "redisPwd"
```

#### 2.5 create database & import data

```
mysql -h 127.0.0.1 -u root -p -e "CREATE DATABASE gatekeeper DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;"
mysql -h 127.0.0.1 -u root -p gatekeeper < gatekeeper.sql --default-character-set=utf8
```

#### 2.6 run gatekeeper
- 启动控制端

```
sh control.sh start_control 
```
- 启动代理端
```
sh control.sh start_proxy 
```
- 同时双启动
```
sh control.sh start_both 
```
---

# 自动化集成测试

### [自动化集成测试](https://github.com/didi/Gatekeeper/blob/master/test_suites/README.md)