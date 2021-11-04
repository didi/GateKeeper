package http_proxy_middleware

import (
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"strings"
)

func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		transferRule := serviceDetail.PluginConf.GetPath("header_transfer", "header_transfer_rule").MustString()
		for _, item := range strings.Split(transferRule, "\n") {
			if item == "" {
				continue
			}
			items := strings.Split(item, " ")
			if len(items) == 3 {
				if items[0] == "add" {
					c.Request.Header.Add(items[1], items[2])
				}
				if items[0] == "edit" {
					c.Request.Header.Set(items[1], items[2])
				}
				continue
			}
			if len(items) == 2 && items[0] == "del" {
				c.Request.Header.Del(items[1])
			}
		}
		c.Next()
	}
}
