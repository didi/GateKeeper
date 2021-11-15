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
cd GateKeeper
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

### [自动化集成测试](https://github.com/didi/GateKeeper/blob/master/test_suites/README.md)

#### 插件化帮助文档
**插件机制是允许开发者进行二次开发的灵活机制。开发者可以通过开发自己的插件来实现自身所需的业务逻辑。**
![avatar](https://s3-gzpu.didistatic.com/stscar/didiGateway/E337CB04-EFAF-475B-B5EE-5363F0B6C85A_4_5005_c.jpeg)
* **服务端：**
    * 插件生产
        在http_proxy_middleware目录中，创建一个go文件，在文件中编写自身的业务逻辑即可。
        以http_flow_limit插件为例：
            1.在http_proxy_middleware目录下新建http_flow_limit.go文件，在文件中编写逻辑，通过调用serviceDetail下的PluginConf自带的GetPath方法获取存储在json中的值。
            ![avatar](https://s3-gzpu.didistatic.com/stscar/didiGateway/36801F68-EF4B-44FA-A926-FA6EA5430AC3.jpeg)
    * 插件注册
        插件分为配置型插件、负载型插件。这两种插件的注册方式不同，配置型插件我们采用在路由中注册的方式，负载型插件的注册方式为通过调用RegisterLoadBalanceStrategyHandler函数在内存中进行注册。
        以http_flow_limit配置插件为例：
        ![avatar](https://s3-gzpu.didistatic.com/stscar/didiGateway/C3A3DB2F-EDC7-4921-BB7A-90EB511BF4C0.jpeg)
        以consistent_hash_strategy负载插件为例：
        ![avatar](https://s3-gzpu.didistatic.com/stscar/didiGateway/CF92D3EE-5ABA-4D1E-9DB9-EB54546E3D19.jpeg)
    * 插件配置  
        插件功能完善之后，如果想在控制端进行显示，则需要进行配置文件的编写（这里基于前端自动化配置化的思路）。
        1. 配置方式：
            在conf目录下找到自己当前运行环境所调用的目录，选中plugin_config.toml文件。按照下图所示方式进行填写。本配置文件支持toml格式。
            ![avatar](https://s3-gzpu.didistatic.com/stscar/didiGateway/EFB8E1A3-15EE-442B-A9E4-82EF3A917548.jpeg)
        2. 配置参数
            | key | value | 注释 |
            | --- | --- | --- |
            | display_name | "限流中间件" | 中间件的名称 |
            | sort | 3 | 中间件的排序，默认值越大越靠前 |
            | postion | "normal" | 中间件的类型，目前分三种normal、auth、loadbalance。normal的则会再基本信息之后，鉴权之前展示，auth和loadbalance都是通过选择框进行选择展示 |
             | unique_name | "tcp_flow_limit"  | 中间件的唯一key，插件生效也是通过这个key去内存中获取信息 |
             | field_type | "input"  | 中间件的内部某一行的前端展示类型，目前支持input、select、textarea，radio、checkbox、switch这几种类型 |
             | field_display | "inline"  | 中间件的内部某一行的前端样式，目前包括inline、block。默认为block。block表示当前的一个items即为一行，后面的items换行展示 |
             | field_clear | "none"  | 中间件的内部某一行的前端辅助样式, 目前为left、none、right。当field_display为inline时生效。如果开头为left或者none时，表示后续的items仍然在后面进行展示，直到为right时才终止，下一个items才换行显示 |
             | field_placeholder | ""  | 中间件内部的某一行的前言，例如input的框内提示文字但不限于 |
             | field_option | ""  | 中间件的选项，例如select的下拉选择框中的内容但不限于 |
             | field_value | ""  | 中间件内部某一行的选中的value信息 |
             | field_default_value | "0"  | 中间件某一项的默认值 |
             | field_unique_name | "service_flow_limit_num"  | 中间件某一项唯一key，通过这个key后台进行获取当前的设置值 |
             | field_display_name | "服务限流数"  | 中间件某一项的显示名称 |
             | field_required | false  | 中间件某一项是否必填 |
             | field_valid_rule |  "/^[0-9]$/"    | 中间件某一项的正则校验 |
* **控制端：**
    * 插件渲染  
        服务端将配置好的插件信息传递到控制端，此时数据近似于树的主干，我们支持根据配置化来进行不同前端样式的展示，所以我们自己写了一个生成器，用于生成所需要的树。再将生成的树通过虚拟dom的方式渲染到页面中。
        ![avatar](https://s3-gzpu.didistatic.com/stscar/didiGateway/42E0EC5D-7A8E-47D0-9F75-EA0EC36D5F17.jpeg)
    * 插件使用  
        目前项目中支持创建HTTP服务，在创建的HTTP服务中，可以根据自己配置的插件，进行插件内容的填写，之后点击确定保存。所有经过新建的服务的请求都会受限于该插件。
