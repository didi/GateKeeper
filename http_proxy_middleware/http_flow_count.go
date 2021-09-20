package http_proxy_middleware

import (
	"bytes"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
)

func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail,err:= model.GetServiceDetailFromGinContext(c)
		if err!=nil{
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		totalCounter, err := handler.ServiceCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			public.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		totalCounter.Increase()

		sCounterBuffer := bytes.NewBufferString(public.FlowServicePrefix)
		sCounterBuffer.WriteString(serviceDetail.Info.ServiceName)
		serviceCounter, err := handler.ServiceCounterHandler.GetCounter(sCounterBuffer.String())
		if err != nil {
			public.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}
}
