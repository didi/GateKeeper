package http_proxy_middleware

import (
	"strings"

	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
)

func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := model.GetServiceDetailFromGinContext(c)
		if err != nil {
			public.ResponseError(c, 2001, err)
			c.Abort()
			return
		}
		if serviceDetail.Info.HttpStripPrefix == 1 {
			strArr := strings.Split(serviceDetail.Info.HTTPPaths,"\n")
			tmpStr := c.Request.URL.Path
			for _, replaceStr := range strArr {
				c.Request.URL.Path = strings.Replace(c.Request.URL.Path, replaceStr, "", 1)
				if c.Request.URL.Path != tmpStr {
					break
				}
			}

		}
		c.Next()
	}
}
