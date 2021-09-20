package tcp_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"strings"
)

func TCPFlowLimitMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*model.ServiceDetail)
		serviceFlowNum := serviceDetail.PluginConf.GetPath("http_flow_limit","service_flow_limit_num").MustInt()
		serviceFlowType := serviceDetail.PluginConf.GetPath("http_flow_limit","service_flow_limit_type").MustInt()
		if serviceFlowNum != 0 {
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(serviceFlowNum), serviceFlowType, true)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("service flow limit %v", serviceFlowNum)))
				c.Abort()
				return
			}
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		clientIpFlowNum := serviceDetail.PluginConf.GetPath("tcp_flow_limit","clientip_flow_limit_num").MustInt()
		clientIpFlowType := serviceDetail.PluginConf.GetPath("tcp_flow_limit","clientip_flow_limit_type").MustInt()
		if clientIpFlowNum > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP, float64(clientIpFlowNum), clientIpFlowType, true)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%v flow limit %v", clientIP, clientIpFlowNum)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
