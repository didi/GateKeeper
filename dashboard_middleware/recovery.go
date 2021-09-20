package dashboard_middleware

import (
	"errors"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/trace"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if lib.ConfBase.Log.On {
					public.GinLogWarning(c, trace.DLTagUndefined, map[string]interface{}{
						"error": fmt.Sprint(err),
						"stack": string(debug.Stack()),
					})
				}
				if lib.ConfBase.Base.DebugMode != "debug" {
					ResponseError(c, 500, errors.New("内部错误"))
					return
				} else {
					ResponseError(c, 500, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}