package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func GrpcJwtFlowCountMiddleware(serviceDetail *model.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error{
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, grpcHandler grpc.StreamHandler) error{
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		if serviceDetail.Info.AuthType != "jwt_auth" {
			if err := grpcHandler(srv, ss); err != nil {
				log.Printf("GrpcJwtAuthTokenMiddleware failed with error %v\n", err)
				return err
			}
			return nil
		}

		appInfos := md.Get("app")
		if len(appInfos)==0 {
			if err := grpcHandler(srv, ss);err != nil {
				log.Printf("RPC failed with error %v\n", err)
				return err
			}
			return nil
		}
		appInfo := &model.App{}
		if err:=json.Unmarshal([]byte(appInfos[0]),appInfo);err!=nil{
			return err
		}
		appCounter, err := handler.AppCounterHandler.GetCounter(public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			return err
		}
		appCounter.Increase()
		if appInfo.Qpd>0 && appCounter.TotalCount>appInfo.Qpd{
			return errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v",appInfo.Qpd,appCounter.TotalCount))
		}
		if err := grpcHandler(srv, ss);err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
