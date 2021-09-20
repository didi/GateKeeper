package http_proxy_middleware

import (
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
)

//匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := handler.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			public.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		c.Set(public.ServiceDetailContextKey, service)
		c.Next()
	}
}
