package public

import (
	"context"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"runtime"
)

//TraceTagInfo trace日志
func TraceTagInfo(trace *lib.TraceContext, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		lib.Log.TagInfo(trace, dltag, m)
	}
}

//TraceTagWarn  trace报警
func TraceTagWarn(trace *lib.TraceContext, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		lib.Log.TagWarn(trace, dltag, m)
	}
}

//TraceTagError trace错误
func TraceTagError(trace *lib.TraceContext, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		lib.Log.TagError(trace, dltag, m)
	}
}

//ContextWarning 错误日志
func ContextWarning(c context.Context, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		v := c.Value(ContextKey("trace"))
		traceContext, ok := v.(*lib.TraceContext)
		if !ok {
			traceContext = lib.NewTrace()
		}
		lib.Log.TagWarn(traceContext, dltag, m)
	}
}

//ContextError 错误日志
func ContextError(c context.Context, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		v := c.Value(ContextKey("trace"))
		traceContext, ok := v.(*lib.TraceContext)
		if !ok {
			traceContext = lib.NewTrace()
		}
		lib.Log.TagError(traceContext, dltag, m)
	}
}

//ContextNotice 普通日志
func ContextNotice(c context.Context, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		v := c.Value(ContextKey("trace"))
		traceContext, ok := v.(*lib.TraceContext)
		if !ok {
			traceContext = lib.NewTrace()
		}
		lib.Log.TagInfo(traceContext, dltag, m)
	}
}

//ComLogWarning 错误日志
func ComLogWarning(c *gin.Context, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		traceContext := GetGinTraceContext(c)
		lib.Log.TagError(traceContext, dltag, m)
	}
}

//ComLogNotice 普通日志
func ComLogNotice(c *gin.Context, dltag string, m map[string]interface{}) {
	if TraceLoggerOn {
		traceContext := GetGinTraceContext(c)
		lib.Log.TagInfo(traceContext, dltag, m)
	}
}

//GetGinTraceContext 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *lib.TraceContext {
	// 防御
	if c == nil {
		return lib.NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*lib.TraceContext); ok {
			return tc
		}
	}
	return lib.NewTrace()
}

//GetTraceContext 从Context中获取数据
func GetTraceContext(c context.Context) *lib.TraceContext {
	if c == nil {
		return lib.NewTrace()
	}
	traceContext := c.Value(ContextKey("trace"))
	if tc, ok := traceContext.(*lib.TraceContext); ok {
		return tc
	}
	return lib.NewTrace()
}

//LogErrorf 将msg单独提出来打印
func LogErrorf(c context.Context, msg string, m Map) {
	_, f, l, _ := runtime.Caller(1)
	if m == nil {
		m = Map{}
	}
	m["caller"] = fmt.Sprintf("file: %s, line: %d", f, l)
	m["msg"] = msg
	SysLogger.Error("%v", m)
}

//LogWarnf 将msg单独提出来打印
func LogWarnf(c context.Context, msg string, m Map) {
	_, f, l, _ := runtime.Caller(1)
	if m == nil {
		m = Map{}
	}
	m["caller"] = fmt.Sprintf("file: %s, line: %d", f, l)
	m["msg"] = msg
	SysLogger.Warn("%v", m)
}

//LogInfof 将msg单独提出来打印
func LogInfof(c context.Context, msg string, m Map) {
	_, f, l, _ := runtime.Caller(1)
	if m == nil {
		m = Map{}
	}
	m["caller"] = fmt.Sprintf("file: %s, line: %d", f, l)
	m["msg"] = msg
	SysLogger.Info("%v", m)
	//ContextNotice(c, lib.DLTagUndefind, m)
	if !IsProductEnv() {
		fmt.Printf("[INFO]: %s. info: %+v\n", msg, m)
	}
}
