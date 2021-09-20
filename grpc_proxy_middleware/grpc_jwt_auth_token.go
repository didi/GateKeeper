package grpc_proxy_middleware

import (
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

//jwt auth token
func GrpcJwtAuthTokenMiddleware(serviceDetail *model.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, grpcHandler grpc.StreamHandler) error {
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

		authToken := ""
		auths := md.Get("authorization")
		if len(auths) > 0 {
			authToken = auths[0]
		}
		appMatched := false
		claims, err := public.JwtDecode(strings.ReplaceAll(authToken, "Bearer ", ""))
		if err != nil {
			return errors.WithMessage(err, "JwtDecode")
		}
		appList := handler.AppManagerHandler.GetAppList()
		for _, appInfo := range appList {
			if appInfo.AppID == claims.Issuer {
				md.Set("app", public.Obj2Json(appInfo))
				appMatched = true
				break
			}
		}
		if !appMatched {
			return errors.New("not match valid app")
		}
		if err := grpcHandler(srv, ss); err != nil {
			log.Printf("GrpcJwtAuthTokenMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
