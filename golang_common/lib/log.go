package lib

import (
	"context"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/trace"
	"strings"
)

var Log *Logger

type Logger struct {
}

func (l *Logger) TagInfo(traceContext *trace.Trace, dltag string, m map[string]interface{}) {
	ZLog.Infof(trace.SetCtxTrace(context.Background(), traceContext), dltag, parseParams(m))
}

func (l *Logger) TagWarn(traceContext *trace.Trace, dltag string, m map[string]interface{}) {
	ZLog.Warnf(trace.SetCtxTrace(context.Background(), traceContext), dltag, parseParams(m))
	//m[_dlTag] = checkDLTag(dltag)
	//m[_traceId] = trace.TraceId
	//m[_childSpanId] = trace.CSpanId
	//m[_spanId] = trace.SpanId
	//dlog.Warn(parseParams(m))
}

func (l *Logger) TagError(traceContext *trace.Trace, dltag string, m map[string]interface{}) {
	ZLog.Errorf(trace.SetCtxTrace(context.Background(), traceContext), dltag, parseParams(m))
}

func parseParams(m map[string]interface{}) string {
	var dltag string = "_undef"
	if _dltag, _have := m["dltag"]; _have {
		if __val, __ok := _dltag.(string); __ok {
			dltag = __val
		}
	}
	for _key, _val := range m {
		if _key == "dltag" {
			continue
		}
		dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _key, _val)
	}
	dltag = strings.Trim(fmt.Sprintf("%q", dltag), "\"")
	return dltag
}

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) string {
	return textColor(textBlack, str)
}

func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
