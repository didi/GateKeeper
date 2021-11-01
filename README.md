# 快速开始（Quick Start）

## 环境要求 >= go 1.12.0
_源码环境要求_


## 环境安装

[GO环境安装](https://golang.google.cn/dl/) ||  [教程](https://www.runoob.com/go/go-environment.html)

[Mysql安装](https://dev.mysql.com/downloads/mysql/) || [教程](https://www.runoob.com/mysql/mysql-install.html)

[Redis安装](https://github.com/tporadowski/redis/releases) || [教程](https://www.runoob.com/redis/redis-install.html)


## 网关服务安装方式(二选一)
1. 源码运行
2. 下载编译包: https://github.com/didi/GateKeeper/releases


### 服务初始化(二选一)

_首次运行时需要_

#### (1)自动初始化(二选一)

_①下载对应版本编译包自动初始化并安装_
- curl -L 'https://github.com/didi/GateKeeper/releases/download/v1.0.0/setup.sh' | sh

_②手动执行自动初始化脚本_
- cd install && sh install.sh 


#### (2)手动初始化

- 下载类库依赖

```
export GO111MODULE=on && export GOPROXY=https://goproxy.cn
cd gatekeeper
go mod tidy
```

- 创建 db 并导入数据

```
mysql -h 127.0.0.1 -u root -p -e "CREATE DATABASE gatekeeper DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;"
mysql -h 127.0.0.1 -u root -p gatekeeper < gatekeeper.sql --default-character-set=utf8
```

- 调整 mysql、redis 配置文件

修改 ./conf/dev/mysql.toml 和 ./conf/dev/redis.toml 为自己的环境配置。



## 运行服务

### 运行管理面板

```
sh control.sh start_control 
```


### 运行代理服务

```
sh control.sh start_proxy 
```

# 自动化集成测试

[自动化集成测试](https://github.com/didi/Gatekeeper/blob/master/test_suites/README.md)