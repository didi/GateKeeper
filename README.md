<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [代码帮助](#%E4%BB%A3%E7%A0%81%E5%B8%AE%E5%8A%A9)
    - [运行后端项目](#%E8%BF%90%E8%A1%8C%E5%90%8E%E7%AB%AF%E9%A1%B9%E7%9B%AE)
    - [运行前端项目](#%E8%BF%90%E8%A1%8C%E5%89%8D%E7%AB%AF%E9%A1%B9%E7%9B%AE)
    
## 当前2.0版本，老用户请查看 [[v1]](https://github.com/didi/Gatekeeper/tree/v1) 版本

## 代码帮助

### 运行后端项目

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
mysql -h localhost -u root -p -e "CREATE DATABASE gatekeeper DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;"
mysql -h localhost -u root -p gatekeeper < gatekeeper.sql --default-character-set=utf8
```

- 调整 mysql、redis 配置文件

修改 ./conf/dev/mysql.toml 和 ./conf/dev/redis.toml 为自己的环境配置。

- 运行面板、代理服务

运行管理面板配合前端项目 - 达成服务管理功能

```
go run main.go run -c ./conf/dev/ -p control
```

打开控制面板，账号密码默认

```
http://127.0.0.1:8880/dist/
```

运行代理服务

```
go run main.go run -c ./conf/dev/ -p proxy
```

### docker及k8s部署使用

待补充

### 构建ingress使用

待补充

### 测试

性能测试待补充 覆盖率测试待补充