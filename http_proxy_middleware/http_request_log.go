package http_proxy_middleware

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/trace"
	"github.com/gin-gonic/gin"
	"time"
)

func HTTPRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !lib.ConfBase.Log.On {
			c.Next()
			return
		}
		t := time.Now()
		c.Request = c.Request.WithContext(trace.SetCtxTrace(c.Request.Context(), trace.New(c.Request)))
		lib.ZLog.Infof(c.Request.Context(), trace.DLTagRequestIn, "")
		c.Next()
		lib.ZLog.Infof(c.Request.Context(), trace.DLTagRequestOut, "proc_time=%f||status=%v", time.Since(t).Seconds(), c.Writer.Status())
		if traceCtx, ok := trace.GetCtxTrace(c.Request.Context()); ok {
			trace.PutTrace(traceCtx)
		}
	}
}