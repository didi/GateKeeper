package main

import (
	"github.com/didi/gatekeeper/dashboard_router"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/grpc_proxy_router"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/http_proxy_router"
	"github.com/didi/gatekeeper/tcp_proxy_router"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := lib.CmdExecute(); err != nil || lib.GetCmdPanelType() == "" {
		os.Exit(1)
	}
	if lib.GetCmdPanelType() == "proxy" {
		startProxy()
	}
	if lib.GetCmdPanelType() == "control" {
		startControl()
	}
}

func startControl() {
	log.Info().Msg(lib.Purple("start controller application"))
	lib.InitConf(lib.GetCmdConfPath())
	defer lib.DestroyConf()
	handler.ServiceManagerHandler.LoadAndWatch()
	dashboard_router.HttpServerRun()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	dashboard_router.HttpServerStop()
}

func startProxy() {
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

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	tcp_proxy_router.TcpManagerHandler.TcpServerStop()
	grpc_proxy_router.GrpcManagerHandler.GrpcServerStop()
	http_proxy_router.HttpServerStop()
	http_proxy_router.HttpsServerStop()
}
