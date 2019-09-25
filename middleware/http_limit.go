package middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/public"
	"github.com/didi/gatekeeper/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

//HTTPLimit http限流中间件
func HTTPLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取上游服务
		gws, ok := c.MustGet(MiddlewareServiceKey).(*service.GateWayService)
		if !ok {
			public.ResponseError(c, http.StatusBadRequest, errors.New("gateway_service not valid"))
			return
		}

		//入口流量统计
		currentModule := gws.CurrentModule()
		counter := public.FlowCounterHandler.GetRequestCounter(currentModule.Base.Name)
		counter.Increase(c.Request.Context(), c.Request.RemoteAddr)

		//客户端ip限流
		remoteIP := public.Substr(c.Request.RemoteAddr, 0, int64(strings.Index(c.Request.RemoteAddr, ":")))
		if currentModule.AccessControl.ClientFlowLimit > 0 {
			limiter := public.FlowLimiterHandler.GetModuleIPVisitor(currentModule.Base.Name+"_"+remoteIP, currentModule.AccessControl.ClientFlowLimit)
			if limiter.Allow() == false {
				errmsg := fmt.Sprintf("moduleName:%s remoteIP：%s, QPS limit : %d, %d", currentModule.Base.Name, remoteIP, int64(limiter.Limit()), limiter.Burst())
				public.ContextWarning(c.Request.Context(), service.DLTagAccessControlFailure, map[string]interface{}{
					"msg":        errmsg,
					"ip":         remoteIP,
					"moduleName": currentModule.Base.Name,
				})
				public.ResponseError(c, http.StatusBadRequest, errors.New(errmsg))
			}
		}

		//todo
		c.Next()
	}
}
