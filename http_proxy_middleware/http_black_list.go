package http_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//匹配接入方式 基于请求信息
func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		whiteListStr := serviceDetail.PluginConf.GetPath("http_whiteblacklist", "ip_white_list").MustString()
		blackListStr := serviceDetail.PluginConf.GetPath("http_whiteblacklist", "ip_black_list").MustString()
		if whiteListStr == "" && public.InIPSliceStr(c.ClientIP(), blackListStr) {
			public.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s in black ip list", c.ClientIP())))
			c.Abort()
			return
		}

		c.Next()
	}
}
