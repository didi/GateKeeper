package dashboard_middleware

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/trace"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ginHands struct {
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
}

func RequestLogger(serName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !lib.ConfBase.Log.On {
			c.Next()
			return
		}
		c.Request = c.Request.WithContext(trace.SetCtxTrace(c.Request.Context(), trace.New(c.Request)))
		t := time.Now()
		path := c.Request.URL.Path
		logData := &ginHands{
			SerName:    serName,
			Path:       path,
			Latency:    time.Since(t),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ClientIP:   c.ClientIP(),
		}
		logRequestIn(c, logData)
		c.Next()
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}
		msg := c.Errors.String()
		if msg != "" {
			msg = strings.Replace(c.Errors.String(), "\n", "", -1)
		} else {
			msg = "Request"
		}
		logData.Latency = time.Since(t)
		logData.StatusCode = c.Writer.Status()
		logData.MsgStr = msg
		logRequestOut(c, logData)
	}
}

func logRequestIn(c *gin.Context, data *ginHands) {
	lib.ZLog.Infof(c.Request.Context(), trace.DLTagRequestIn, "ser_name=%v||method=%v||path=%v||client_ip=%v", data.SerName, data.Method, data.Path, data.ClientIP)
}

func logRequestOut(c *gin.Context, data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		lib.ZLog.Warnf(c.Request.Context(), trace.DLTagRequestOut, "ser_name=%v||method=%v||path=%v||proc_time=%v||status=%v||client_ip=%v||msg=%v", data.SerName, data.Method, data.Path, data.Latency, data.StatusCode, data.ClientIP, data.MsgStr)
	case data.StatusCode >= 500:
		lib.ZLog.Errorf(c.Request.Context(), trace.DLTagRequestOut, "ser_name=%v||method=%v||path=%v||proc_time=%v||status=%v||client_ip=%v||msg=%v", data.SerName, data.Method, data.Path, data.Latency, data.StatusCode, data.ClientIP, data.MsgStr)
	default:
		lib.ZLog.Infof(c.Request.Context(), trace.DLTagRequestOut, "ser_name=%v||method=%v||path=%v||proc_time=%v||status=%v||client_ip=%v||msg=%v", data.SerName, data.Method, data.Path, data.Latency, data.StatusCode, data.ClientIP, data.MsgStr)
	}
}
