package grpc_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

func GrpcBlackListMiddleware(serviceDetail *model.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		whiteListStr := serviceDetail.PluginConf.GetPath("grpc_whiteblacklist", "ip_white_list").MustString()
		blackListStr := serviceDetail.PluginConf.GetPath("grpc_whiteblacklist", "ip_black_list").MustString()

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		if whiteListStr == "" && blackListStr != "" {
			if public.InIPSliceStr(clientIP, blackListStr) {
				return errors.New(fmt.Sprintf("%s in black ip list", clientIP))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
