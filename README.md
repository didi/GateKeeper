# 项目介绍
GateKeeper是一款golang编写的支持快速开发、插件化的高性能网关。使用二进制文件安装即刻体验网关功能。

功能特性：   
1、快速开发：插件化支持功能拓展+参数配置化，双管齐下。  
2、健康检查：支持主动+被动探活检测，还可拓展第三方配置中心。     
3、内置插件：分布式限流(秒/时/日三种粒度)、header头转换、白名单、租户鉴权、QPS统计等

系统架构：
![avatar](http://img-hxy021.didistatic.com/static/itstool_public/do1_QEfWPEJgafZ8aPjT83eG)

1、系统分为两个端：控制端(Control Plane)负责服务编辑配置和自定义参数管理；代理端(Proxy Plane)负责协议数据中间处理及转发。两端可以分开启动也同时启动。      
2、控制端与代理端交互：控制端编辑完服务信息会写入Mysql中，代理端定时从Mysql中拉取服务配置信息，然后平滑处理代理端服务启停。     
3、自定义化插件：支持最常见的业务需求快速定制化，主要包含业务中间件自定义、服务负载配置自定义、负载策略自定义等，如果插件需要参数配置则可以在控制端配置相应参数。

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
golang | \>=1.16 | [官网](https://golang.google.cn/dl/) |  [教程](https://www.runoob.com/go/go-environment.html)

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

# 插件化
网关中最常见的业务需求，包含业务中间件自定义、服务负载配置自定义、负载策略自定义等。  
这些常见功能我们都支持插件化定义，除此之外还支持了插件内部的动态参数配置及获取。

## 业务中间件插件
由于中间件内部采用AOP切面编程实现，所以业务中间件我们直接套用了 gin 中间件定义：
```
type HandlerFunc func(*Context)
```
示例demo，func直接返回定义的方法即可：
```
func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		//todo
		c.Next()
	}
}
```
最后记得完成中间件在 http_proxy_router/route.go 的内部注册。
```
router.Use(
    http_proxy_middleware.HTTPAccessModeMiddleware(),
    ...
    http_proxy_middleware.HTTPFlowLimitMiddleware(),
    http_proxy_middleware.HTTPReverseProxyMiddleware())
```
## 负载配置插件
系统默认支持兜底下游节点配置，但是有些服务是使用服务发现协调器实现的如 consul、zookeeper 这就需要自定义服务负载配置了。    
与上面相同首先是接口定义：
```
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
	CloseWatch()
}
```
实现完上述接口的结构后，还需要注册这个步骤。
```
func init() {
	RegisterCheckConfigHandler("upstream_config", NewLoadBalanceCheckConf)
}
```
## 负载策略插件
系统默认提供了4种负载均衡策略：轮询、权重轮询、基于客户端IP的一致性Hash、随机。如果有定义需求还可以使用插件进行功能拓展。
同样首先是接口定义
```
type LoadBalanceStrategy interface {
	Add(...string) error
	RemoveAll() error
	GetAll() ([]string, error)
	Get(string) (string, error)
}
```
实现完上述接口的结构后，还需要注册这个步骤。
```
func init() {
	RegisterLoadBalanceStrategyHandler("random", func() LoadBalanceStrategy {
		return &RandomStrategy{}
	})
}
```

## 动态参数配置及获取
插件要做到一处开发多处使用，就离不开动态参数配置功能。
GateKeeper要实现参数配置主要以下几个步骤：     
首先编辑 `conf/dev/plugin_config.toml` 增加参数配置，具体参数格式介绍 待补充，现举例如下：
```
[[http]]
  display_name = "url地址重写"
  sort = 6
  postion = "normal"
  unique_name = "url_rewrite"

  [[http.items]]
    field_type = "textarea"
    field_display = "block"
    field_clear = "none"
    field_placeholder = "格式：^/test_service(.*) $1\n多条请换行"
    field_option = ""
    field_value = ""
    field_default_value = ""
    field_unique_name = "rewrite_rule"
    field_display_name = "URL重写"
    field_required = false
    field_valid_rule = "/^[\\S]+ [\\S]+$/is"
```

其次，打开控制端修改响应服务信息。
![avatar](http://img-hxy021.didistatic.com/static/itstool_public/do1_d18zzT6DBk9zXHwq3AdN)

最后在响应中间件调用即可。
```
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		rewriteUrl := serviceDetail.PluginConf.GetPath("url_rewrite", "rewrite_rule").MustString()
		...
		c.Next()
	}
}
```

# 性能测试
GateKeeper 对比业界其他知名API网关 APISIX、Nginx、HAProxy、Kong、manba  
硬件配置：128G、40核 Intel(R) Xeon(R) Silver 4114 2.20GHz  
测试参数：wrk -t50 -c500 -d30s --latency  "http://xxxx/xxxx/"  

![avatar](http://img-hxy021.didistatic.com/static/itstool_public/do1_E82bzkP6K8qaiUhmgkdA)

![avatar](http://img-hxy021.didistatic.com/static/itstool_public/do1_pDPElgrgBdBpMyVHafkF)

![avatar](http://img-hxy021.didistatic.com/static/itstool_public/do1_y5gXsG6Kx9yhEk6FsAPE)

通过以上图片分析不难得出，GateKeeper性能的表现特点为：
- 高并发压力下并发性可达到主流网关需求。
- 高并发压力下比同类软件内存占用较高。
- 高并发压力下99分位比同类软件延迟最低。

注：以上服务均为默认安装未做调优、不同电脑配置测试结果可能不同。    
更多详细内容 待补充

# 自动化集成测试
[自动化集成测试](https://github.com/didi/GateKeeper/blob/master/test_suites/README.md)