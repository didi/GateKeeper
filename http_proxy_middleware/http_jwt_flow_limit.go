package http_proxy_middleware

import (
	"bytes"
	"fmt"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		if serviceDetail.Info.AuthType != "jwt_auth" {
			c.Next()
			return
		}

		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*model.App)
		if appInfo.Qps > 0 {
			cLimiterBuffer := bytes.NewBufferString(public.FlowAppPrefix)
			cLimiterBuffer.WriteString(appInfo.AppID)
			cLimiterBuffer.WriteString("_")
			cLimiterBuffer.WriteString(c.ClientIP())
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(cLimiterBuffer.String(), float64(appInfo.Qps), 0, true)
			if err != nil {
				public.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				public.ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps), ))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
