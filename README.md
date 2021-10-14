# 快速安装（Quick Start）

- xxxx

# 手动安装服务

- 首先git clone 本项目

`git clone git@github.com:didi/gatekeeper.git`

- 确保本地环境安装了Go 1.12+版本

```
go version
go version go1.12.15 darwin/amd64
```

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

- 运行服务

运行管理面板

```
go run main.go run -c ./conf/dev/ -p control
```

运行代理服务

```
go run main.go run -c ./conf/dev/ -p proxy
```

# 自动化集成测试

请参照 test_suites/README.md