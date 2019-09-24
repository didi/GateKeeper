package common

import (
	"flag"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/router"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/didichuxing/gatekeeper/tester/testhttp"
	"github.com/didichuxing/gatekeeper/tester/testrpc/thriftserver"
	"github.com/e421083458/golang_common/lib"
	"time"
)

var(
	testTCP  *thriftserver.TestTCPDestServer
	testHTTP *testhttp.HTTPDestServer
)

//SetUp 套件创建执行函数
func SetUp()  {
	conf := flag.String(
		"config",
		"../../conf/test/",	//默认配置文件
		"input config file like ../../conf/dev/")
	flag.Parse()
	lib.InitModule(*conf,[]string{"base","mysql","redis","test_dest",})
	defer lib.Destroy()
	public.InitMysql()
	public.InitConf()

	//配置管理
	service.SysConfMgr = service.NewSysConfigManage()
	service.SysConfMgr.InitConfig()
	service.SysConfMgr.MonitorConfig()

	//todo
	//注册验证方法
	service.RegisterBeforeRequestAuthFunc(service.AuthAppToken)

	//注册过滤方法
	service.RegisterModifyResponseFunc(service.FilterCityData([]string{"/gatekeeper/tester_filter/goods_list"}))

	//目标http服务器
	testHTTP = testhttp.NewTestHTTPDestServer()
	httpAddrSlice := lib.GetStringSliceConf("test_dest.http_dest.addrs")
	for _,addr:=range httpAddrSlice {
		testHTTP.Run(addr)
	}

	//目标tcp服务器
	testTCP = thriftserver.NewTestTCPDestServer()
	tcpAddrSlice := lib.GetStringSliceConf("test_dest.tcp_dest.addrs")
	for _,addr:=range tcpAddrSlice {
		testTCP.Run(addr)
	}
	time.Sleep(5*time.Second)
	router.HTTPServerRun()
	router.TCPServerRun()
	time.Sleep(1*time.Second)
}

//Before 用例创建执行函数
func Before()  {
}