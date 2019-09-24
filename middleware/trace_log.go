package middleware

import (
	"bytes"
	"context"
	"github.com/didichuxing/gatekeeper/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

//RequestTraceLog trace中间件
func RequestTraceLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequestInLog(c)
		defer RequestOutLog(c)
		c.Next()
	}
}

//RequestInLog 请求进入日志
func RequestInLog(c *gin.Context) {
	traceContext := lib.NewTrace()
	if traceID := c.Request.Header.Get("didi-header-rid"); traceID != "" {
		traceContext.TraceId = traceID
	}
	if spanID := c.Request.Header.Get("didi-header-spanid"); spanID != "" {
		traceContext.SpanId = spanID
	}
	bodyBytes, rerr := ioutil.ReadAll(c.Request.Body)
	if rerr != nil {
		public.TraceTagError(traceContext, "read_body_err", map[string]interface{}{
			"uri":    c.Request.RequestURI,
			"method": c.Request.Method,
			"args":   c.Request.PostForm,
			"body":   string(bodyBytes),
			"from":   c.ClientIP(),
		})
	}
	c.Set(MiddlewareRequestBodyKey, bodyBytes)
	public.TraceTagInfo(traceContext, "_com_request_in", map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"body":   string(bodyBytes),
		"from":   c.ClientIP(),
	})
	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)
	c.Request.Header.Set("didi-header-rid", traceContext.TraceId)

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), public.ContextKey("trace"), traceContext))
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), public.ContextKey("request_url"), c.Request.URL.Path))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}

//RequestOutLog 请求结束日志
func RequestOutLog(c *gin.Context) {
	endExecTime := time.Now()

	//response, _ := c.Get("response")
	st, stok := c.Get("startExecTime")
	if !stok {
		public.ComLogWarning(c, "gin_context_geterr", map[string]interface{}{
			"uri":    c.Request.RequestURI,
			"method": c.Request.Method,
			"args":   c.Request.PostForm,
			"from":   c.ClientIP(),
			//"response": response,
			"key": "startExecTime",
		})
	}
	startExecTime, ok := st.(time.Time)
	if !ok {
		public.ComLogWarning(c, "golang_assertion_err", map[string]interface{}{
			"uri":    c.Request.RequestURI,
			"method": c.Request.Method,
			"args":   c.Request.PostForm,
			"from":   c.ClientIP(),
			//"response": response,
			"key": "st",
			"st":  st,
		})
	}
	public.ComLogNotice(c, "_com_request_out", map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"from":   c.ClientIP(),
		//"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	})
}
