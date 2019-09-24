package middleware

import (
	"bytes"
	"errors"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//LoadBalance 负载均衡中间件
func LoadBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		gws,ok:=c.MustGet(MiddlewareServiceKey).(*service.GateWayService)
		if !ok{
			public.ResponseError(c, http.StatusBadRequest, errors.New("gateway_service not valid"))
			return
		}
		proxy, err := gws.LoadBalance()
		if err != nil {
			public.ResponseError(c, http.StatusProxyAuthRequired, err)
			return
		}
		requestBody,ok:=c.MustGet(MiddlewareRequestBodyKey).([]byte)
		if !ok{
			public.ResponseError(c, http.StatusBadRequest, errors.New("request_body not valid"))
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
