package public

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/trace"
	"github.com/gin-gonic/gin"
)

func GinLogNotice(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	lib.Log.TagInfo(traceContext, dltag, m)
}

func GinLogWarning(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	lib.Log.TagError(traceContext, dltag, m)
}

func GetGinTraceContext(c *gin.Context) *trace.Trace {
	traceContext, exists := trace.GetCtxTrace(c.Request.Context())
	if exists {
		return traceContext
	}
	return trace.New(c.Request)
}