package grpc_proxy_router

import (
	"fmt"
	"github.com/didi/gatekeeper/grpc_proxy_middleware"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/didi/gatekeeper/reverse_proxy"
	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type GrpcManager struct {
	ServerList []*warpGrpcServer
}

func init() {
	GrpcManagerHandler = NewGrpcManager()
}

func NewGrpcManager() *GrpcManager {
	return &GrpcManager{}
}

var GrpcManagerHandler *GrpcManager

type warpGrpcServer struct {
	Addr        string
	ServiceName string
	UpdateAt    time.Time
	*grpc.Server
}

func (g *GrpcManager) grpcOneServerRun(service *model.ServiceDetail) {
	addr := fmt.Sprintf(":%d", service.Info.Port)
	rb, err := handler.LoadBalancerHandler.GetLoadBalancer(service)
	if err != nil {
		log.Printf(" [ERROR] GetTcpLoadBalancer %v err:%v\n", addr, err)
		return
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf(" [ERROR] GrpcListen %v err:%v\n", addr, err)
		return
	}
	grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(rb)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
		grpc.ChainStreamInterceptor(
			grpc_proxy_middleware.GrpcFlowCountMiddleware(service),
			grpc_proxy_middleware.GrpcFlowLimitMiddleware(service),
			grpc_proxy_middleware.GrpcJwtAuthTokenMiddleware(service),
			grpc_proxy_middleware.GrpcJwtFlowCountMiddleware(service),
			grpc_proxy_middleware.GrpcJwtFlowLimitMiddleware(service),
			grpc_proxy_middleware.GrpcWhiteListMiddleware(service),
			grpc_proxy_middleware.GrpcBlackListMiddleware(service),
			grpc_proxy_middleware.GrpcMetadataTransferMiddleware(service),
		),
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(grpcHandler))

	GrpcManagerHandler.ServerList = append(GrpcManagerHandler.ServerList, &warpGrpcServer{
		Addr:        addr,
		ServiceName: service.Info.ServiceName,
		UpdateAt:    service.Info.UpdatedAt,
		Server:      s,
	})
	log.Printf(" [INFO] grpc_proxy_run %v\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Printf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
	}
}

func (g *GrpcManager) GrpcServerRun() {
	serviceList := handler.ServiceManagerHandler.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go g.grpcOneServerRun(tempItem)
	}
	handler.ServiceManagerHandler.Regist(g)
}

func (g *GrpcManager) Update(e *handler.ServiceEvent) {
	log.Printf("GrpcManager.Update")
	delList := e.DeleteService
	for _, delService := range delList {
		if delService.Info.LoadType == public.LoadTypeGRPC {
			continue
		}
		for _, tcpServer := range GrpcManagerHandler.ServerList {
			if delService.Info.ServiceName != tcpServer.ServiceName {
				continue
			}
			tcpServer.GracefulStop()
			log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", tcpServer.Addr)
		}
	}
	addList := e.AddService
	for _, addService := range addList {
		if addService.Info.LoadType != public.LoadTypeGRPC {
			continue
		}
		go g.grpcOneServerRun(addService)
	}
}

func (g *GrpcManager) GrpcServerStop() {
	for _, grpcServer := range GrpcManagerHandler.ServerList {
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			defer func() {
				wait.Done()
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			grpcServer.GracefulStop()
		}()
		wait.Wait()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
