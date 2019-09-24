package middleware

import (
	"errors"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

//AccessControl 权限控制中间件
func AccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		gws,ok:=c.MustGet(MiddlewareServiceKey).(*service.GateWayService)
		if !ok{
			public.ResponseError(c, http.StatusBadRequest, errors.New("gateway_service not valid"))
			return
		}
		if err := gws.AccessControl(); err != nil {
			public.ResponseError(c, http.StatusUnauthorized, err)
			return
		}
		c.Set(MiddlewareServiceKey,gws)
		c.Next()
	}
}
