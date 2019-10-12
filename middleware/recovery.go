package middleware

import (
	"errors"
	"fmt"
	"github.com/didi/gatekeeper/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

//Recovery 错误捕获中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				public.ComLogWarning(c, "_panic", map[string]interface{}{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				})
				if lib.ConfBase.DebugMode != "debug" {
					public.ResponseError(c, 500, errors.New("内部错误"))
					return
				}
				public.ResponseError(c, 500, errors.New(fmt.Sprint(err)))
				return
			}
		}()
		c.Next()
	}
}
