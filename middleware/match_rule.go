package middleware

import (
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

//中间件常量
const (
	MiddlewareServiceKey     = "gateway_service"
	MiddlewareRequestBodyKey = "request_body"
)

//MatchRule 匹配模块中间件
func MatchRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		gws := service.NewGateWayService(c.Writer, c.Request)
		if err := gws.MatchRule(); err != nil {
			public.ResponseError(c, http.StatusBadRequest, err)
			return
		}
		c.Set(MiddlewareServiceKey,gws)
	}
}