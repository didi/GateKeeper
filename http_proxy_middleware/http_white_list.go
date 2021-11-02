package http_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		whiteListString := serviceDetail.PluginConf.GetPath("http_whiteblacklist", "ip_white_list").MustString()
		if whiteListString != "" {
			if !public.InIPSliceStr(c.ClientIP(), whiteListString) {
				public.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
