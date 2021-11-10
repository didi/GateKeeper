package suites

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/grpc_proxy_router"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/http_proxy_router"
	"github.com/didi/gatekeeper/tcp_proxy_router"
	testsqlhandler "github.com/didi/gatekeeper/test_suites/SqlHandler"
	"github.com/didi/gatekeeper/test_suites/testhttp"
	"github.com/didi/gatekeeper/test_suites/testrpc/thriftserver"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

var (
	testTCP  *thriftserver.TestTCPDestServer
	testHTTP *testhttp.HTTPDestServer
)

func TestRunSuite(t *testing.T) {
	SetUp()
	defer TearDown()
	time.Sleep(500 * time.Millisecond)
	runCase(t, TestGoConvey)
	//runCase(t, TestHostServiceVisit)
	runCase(t, TestStripPrefix)
	//xxx
}

func runCase(t *testing.T, testCase func(*testing.T)) {
	FuncBefore()
	defer FuncAfter()
	testCase(t)
}

func SetUp() {
	gin.SetMode(gin.ReleaseMode)
	lib.SetCmdConfPath("./conf/")
	log.Info().Msg(lib.Purple("start proxy application"))
	lib.InitConf(lib.GetCmdConfPath())
	defer lib.DestroyConf()
	handler.ServiceManagerHandler.LoadAndWatch()
	handler.AppManagerHandler.LoadAndWatch()
	go func() {
		http_proxy_router.HttpServerRun()
	}()
	go func() {
		http_proxy_router.HttpsServerRun()
	}()
	go func() {
		tcp_proxy_router.TcpManagerHandler.TcpServerRun()
	}()
	go func() {
		grpc_proxy_router.GrpcManagerHandler.GrpcServerRun()
	}()

	testHTTP = testhttp.NewTestHTTPDestServer()
	httpAddrSlice := []string{":8881", ":8882"}
	for _, addr := range httpAddrSlice {
		testHTTP.Run(addr)
	}
	testsqlhandler.InitGORMHandler()
	time.Sleep(500 * time.Millisecond)
	//testTCP = thriftserver.NewTestTCPDestServer()
	//tcpAddrSlice := lib.GetStringSliceConf("test_dest.tcp_dest.addrs")
	//for _, addr := range tcpAddrSlice {
	//	testTCP.Run(addr)
	//}
}

func FuncBefore() {
}

func TearDown() {
	tcp_proxy_router.TcpManagerHandler.TcpServerStop()
	grpc_proxy_router.GrpcManagerHandler.GrpcServerStop()
	http_proxy_router.HttpServerStop()
	http_proxy_router.HttpsServerStop()
}

func FuncAfter() {

}
