package tcp_proxy_router

import (
	"context"
	"fmt"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/didi/gatekeeper/reverse_proxy"
	"github.com/didi/gatekeeper/tcp_proxy_middleware"
	"github.com/didi/gatekeeper/tcp_server"
	"log"
	"sync"
)

type TcpManager struct {
	ServerList []*tcp_server.TcpServer
}

func init() {
	TcpManagerHandler = NewTcpManager()
}

func NewTcpManager() *TcpManager {
	return &TcpManager{}
}

var TcpManagerHandler *TcpManager

func (t *TcpManager) tcpOneServerRun(service *model.ServiceDetail) {
	addr := fmt.Sprintf(":%d", service.Info.Port)
	rb, err := handler.LoadBalancerHandler.GetLoadBalancer(service)
	if err != nil {
		log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
		return
	}

	//构建路由及设置中间件
	router := tcp_proxy_middleware.NewTcpSliceRouter()
	router.Group("/").Use(
		tcp_proxy_middleware.TCPFlowCountMiddleware(),
		tcp_proxy_middleware.TCPFlowLimitMiddleware(),
		tcp_proxy_middleware.TCPWhiteListMiddleware(),
		tcp_proxy_middleware.TCPBlackListMiddleware(),
	)

	//构建回调handler
	routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
		func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
			return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
		}, router)

	baseCtx := context.WithValue(context.Background(), "service", service)
	tcpServer := &tcp_server.TcpServer{
		Addr:     addr,
		Handler:  routerHandler,
		BaseCtx:  baseCtx,
		UpdateAt: service.Info.UpdatedAt,
	}
	t.ServerList = append(t.ServerList, tcpServer)
	log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
	if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
		log.Printf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
	}
}

func (t *TcpManager) TcpServerRun() {
	serviceList := handler.ServiceManagerHandler.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(service *model.ServiceDetail) {
			t.tcpOneServerRun(service)
		}(tempItem)
	}
	handler.ServiceManagerHandler.Regist(t)
}

func (t *TcpManager) Update(e *handler.ServiceEvent) {
	log.Printf("TcpManager.Update")
	delList := e.DeleteService
	for _, delService := range delList {
		if delService.Info.LoadType == public.LoadTypeTCP {
			continue
		}
		for _, tcpServer := range TcpManagerHandler.ServerList {
			if delService.Info.ServiceName != tcpServer.ServiceName {
				continue
			}
			tcpServer.Close()
			log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
		}
	}
	addList := e.AddService
	for _, addService := range addList {
		if addService.Info.LoadType != public.LoadTypeTCP {
			continue
		}
		go t.tcpOneServerRun(addService)
	}
}

func (t *TcpManager) TcpServerStop() {
	for _, tcpServer := range t.ServerList {
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			defer func() {
				wait.Done()
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			tcpServer.Close()
		}()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}
