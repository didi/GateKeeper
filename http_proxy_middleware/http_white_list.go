package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}

		ipWhiteListString := serviceDetail.PluginConf.GetPath("http_whiteblacklist", "ip_white_list").MustString()
		if ipWhiteListString != "" {
			if !public.InIPSliceStr(c.ClientIP(), ipWhiteListString) {
				public.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}

		urlWhiteUrlString := serviceDetail.PluginConf.GetPath("http_whiteblacklist", "url_white_list").MustString()
		if urlWhiteUrlString != "" {
			if !public.InURLSliceStr(c.Request.URL.Path, urlWhiteUrlString) {
				public.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white url list", c.Request.URL.Path)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
