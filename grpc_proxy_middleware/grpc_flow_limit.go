package grpc_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

func GrpcFlowLimitMiddleware(serviceDetail *model.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, sHandler grpc.StreamHandler) error {
		serviceFlowNum := serviceDetail.PluginConf.GetPath("grpc_flow_limit","service_flow_limit_num").MustInt()
		serviceFlowType := serviceDetail.PluginConf.GetPath("grpc_flow_limit","service_flow_limit_type").MustInt()
		if serviceFlowNum != 0 {
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(serviceFlowNum), serviceFlowType, true)
			if err != nil {
				return err
			}
			if !serviceLimiter.Allow() {
				return errors.New(fmt.Sprintf("service flow limit %v", serviceFlowNum), )
			}
		}
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		clientIpFlowNum := serviceDetail.PluginConf.GetPath("grpc_flow_limit","clientip_flow_limit_num").MustInt()
		clientIpFlowType := serviceDetail.PluginConf.GetPath("grpc_flow_limit","clientip_flow_limit_type").MustInt()
		if clientIpFlowNum > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP, float64(clientIpFlowNum), clientIpFlowType, true)
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf("%v flow limit %v", clientIP,clientIpFlowNum), )
			}
		}
		if err := sHandler(srv, ss); err != nil {
			log.Printf("GrpcFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
