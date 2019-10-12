package common

import (
	"github.com/didi/gatekeeper/router"
	"github.com/didi/gatekeeper/tester/testrpc/thriftserver"
	"github.com/e421083458/golang_common/lib"
)

//TearDown 测试套件退出函数
func TearDown() {
	router.HTTPServerStop()
	router.TCPServerStop()
	httpAddrSlice := lib.GetStringSliceConf("test_dest.http_dest.addrs")
	for _, addr := range httpAddrSlice {
		testHTTP.Stop(addr)
	}
	testTCP := thriftserver.NewTestTCPDestServer()
	tcpAddrSlice := lib.GetStringSliceConf("test_dest.tcp_dest.addrs")
	for _, addr := range tcpAddrSlice {
		testTCP.Stop(addr)
	}
}

//After 测试用例退出函数
func After() {

}
