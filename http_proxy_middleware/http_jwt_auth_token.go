package http_proxy_middleware

import (
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
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
		appMatched := false
		claims, err := public.JwtDecode(strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", ""))
		if err != nil {
			public.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		appList := handler.AppManagerHandler.GetAppList()
		for _, appInfo := range appList {
			if appInfo.AppID == claims.Issuer {
				c.Set("app", appInfo)
				appMatched = true
				break
			}
		}
		if !appMatched {
			public.ResponseError(c, 2003, errors.New("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}
