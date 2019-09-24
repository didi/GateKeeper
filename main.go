package main

import (
	"flag"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/router"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/e421083458/golang_common/lib"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var (
	conf *string
)

func main() {
	conf = flag.String("config", "./conf/dev/", "input config file like ./conf/dev/")
	flag.Parse()

	lib.InitModule(*conf,[]string{"base","mysql","redis",})
	defer lib.Destroy()
	public.InitMysql()
	public.InitConf()

	//配置管理
	service.SysConfMgr = service.NewSysConfigManage()
	service.SysConfMgr.InitConfig()
	service.SysConfMgr.MonitorConfig()

	//注册请求前验证request方法
	service.RegisterBeforeRequestAuthFunc(service.AuthAppToken)

	//注册请求后更改response方法
	service.RegisterModifyResponseFunc(service.FilterCityData([]string{"/gatekeeper/tester_filter/goods_list"}))

	router.HTTPServerRun()
	router.TCPServerRun()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.TCPServerStop()
	router.HTTPServerStop()
	signal.Stop(quit)
}