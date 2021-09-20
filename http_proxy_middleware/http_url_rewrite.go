package http_proxy_middleware

import (
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

//匹配接入方式 基于请求信息
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail,err:= model.GetServiceDetailFromGinContext(c)
		if err!=nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		rewriteUrl:=serviceDetail.PluginConf.GetPath("rewrite_rule","rewrite_rule").MustString()
		if rewriteUrl == "" {
			c.Next()
			return
		}
		for _, item := range strings.Split(rewriteUrl, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			regexp, err := regexp.Compile(items[0])
			if err != nil {
				continue
			}
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
		}
		c.Next()
	}
}
