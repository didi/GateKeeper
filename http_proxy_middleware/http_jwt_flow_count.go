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

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
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
		countBuffer := bytes.NewBufferString(public.FlowAppPrefix)
		countBuffer.WriteString(appInfo.AppID)
		appCounter, err := handler.ServiceCounterHandler.GetCounter(countBuffer.String())
		if err != nil {
			public.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			public.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount)))
			c.Abort()
			return
		}
		c.Next()
	}
}
