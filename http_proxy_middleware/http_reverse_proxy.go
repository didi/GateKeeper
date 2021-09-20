package http_proxy_middleware

import (
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/didi/gatekeeper/reverse_proxy"
	"github.com/gin-gonic/gin"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		lb, err := handler.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			public.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		trans, err := handler.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			public.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}
