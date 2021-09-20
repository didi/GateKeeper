package ddlog

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/didi/gatekeeper/golang_common/zerolog"
)

const (
	//LogTypeStdout :Log Type Stdout
	LogTypeStdout = "stdout"
	//LogTypeFile :Log Type File
	LogTypeFile = "file"
)

type DiLogHandle struct {
	zerolog.Logger
	CtxFormatFunc func(ctx context.Context) string
}

func (zlog *DiLogHandle) Debug(args ...interface{}) {
	if zerolog.DebugLevel < zlog.Logger.GetLevel() {
		return
	}
	msg := fmt.Sprint(args...)
	zlog.Logger.Debug().Timestamp().CallerDefault("gatekeeper").Msg(msg)
}

func (zlog *DiLogHandle) Debugf(ctx context.Context, tag string, format string, args ...interface{}) {
	if zerolog.DebugLevel < zlog.Logger.GetLevel() {
		return
	}
	var ctxStr string
	if zlog.CtxFormatFunc != nil {
		ctxStr = zlog.CtxFormatFunc(ctx)
	}
	msg := fmt.Sprintf(format, args...)
	zlog.Logger.Debug().Timestamp().CallerDefault("gatekeeper").Tag(tag).Str("", ctxStr).Msg(msg)
}

func (zlog *DiLogHandle) Info(args ...interface{}) {
	if zerolog.InfoLevel < zlog.Logger.GetLevel() {
		return
	}
	msg := fmt.Sprint(args...)
	zlog.Logger.Info().Timestamp().CallerDefault("gatekeeper").Msg(msg)
}

func (zlog *DiLogHandle) Infof(ctx context.Context, tag string, format string, args ...interface{}) {
	if zerolog.InfoLevel < zlog.Logger.GetLevel() {
		return
	}
	var ctxStr string
	if zlog.CtxFormatFunc != nil {
		ctxStr = zlog.CtxFormatFunc(ctx)
	}
	msg := fmt.Sprintf(format, args...)
	zlog.Logger.Info().Timestamp().CallerDefault("gatekeeper").Tag(tag).Str("", ctxStr).Msg(msg)
}

func (zlog *DiLogHandle) Warn(args ...interface{}) {
	if zerolog.WarnLevel < zlog.Logger.GetLevel() {
		return
	}
	msg := fmt.Sprint(args...)
	zlog.Logger.Warn().Timestamp().CallerDefault("gatekeeper").Msg(msg)
}

func (zlog *DiLogHandle) Warnf(ctx context.Context, tag string, format string, args ...interface{}) {
	if zerolog.WarnLevel < zlog.Logger.GetLevel() {
		return
	}
	var ctxStr string
	if zlog.CtxFormatFunc != nil {
		ctxStr = zlog.CtxFormatFunc(ctx)
	}
	msg := fmt.Sprintf(format, args...)
	zlog.Logger.Warn().Timestamp().CallerDefault("gatekeeper").Tag(tag).Str("", ctxStr).Msg(msg)
}

func (zlog *DiLogHandle) Error(args ...interface{}) {
	if zerolog.ErrorLevel < zlog.Logger.GetLevel() {
		return
	}
	msg := fmt.Sprint(args...)
	zlog.Logger.Error().Timestamp().CallerDefault("gatekeeper").Msg(msg)
}

func (zlog *DiLogHandle) Errorf(ctx context.Context, tag string, format string, args ...interface{}) {
	if zerolog.ErrorLevel < zlog.Logger.GetLevel() {
		return
	}
	var ctxStr string
	if zlog.CtxFormatFunc != nil {
		ctxStr = zlog.CtxFormatFunc(ctx)
	}
	msg := fmt.Sprintf(format, args...)
	zlog.Logger.Error().Timestamp().CallerDefault("gatekeeper").Tag(tag).Str("", ctxStr).Msg(msg)
}

func (zlog *DiLogHandle) Fatal(args ...interface{}) {
	if zerolog.FatalLevel < zlog.Logger.GetLevel() {
		return
	}
	msg := fmt.Sprint(args...)
	zlog.Logger.Fatal().Timestamp().CallerDefault("gatekeeper").Msg(msg)
}

func (zlog *DiLogHandle) Fatalf(ctx context.Context, tag string, format string, args ...interface{}) {
	if zerolog.FatalLevel < zlog.Logger.GetLevel() {
		return
	}
	var ctxStr string
	if zlog.CtxFormatFunc != nil {
		ctxStr = zlog.CtxFormatFunc(ctx)
	}
	msg := fmt.Sprintf(format, args...)
	zlog.Logger.Fatal().Timestamp().CallerDefault("gatekeeper").Tag(tag).Str("", ctxStr).Msg(msg)
}

func (zlog *DiLogHandle) RegisterContextFormat(ctxFmt func(ctx context.Context) string) {
	zlog.CtxFormatFunc = ctxFmt
}

type ICtxKey interface{
	GetCtxKey() interface{}
}

type trace interface{
	IsPressureTraffic() bool
}

type PubLog struct {
	zerolog.Logger
	ICtxKey
}

func (plog *PubLog) Public(ctx context.Context, key string, pairs map[string]interface{}) {

	var opera string
	opera = "opera_stat_key=" + key
	if plog.ICtxKey !=nil{
		key := plog.ICtxKey.GetCtxKey()
		val:= ctx.Value(key)
		if val !=nil{
			if v , ok := val.(trace); ok{
				if v.IsPressureTraffic() {
					opera = opera +  "_shadow"
				}
			}
		}

	}
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString("||")
	buffer.WriteString("timestamp=")
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05"))

	plog.Info().Str("", buffer.String()).Fields(pairs).Msg(opera)
}

//PublicString will log  string directly
func (plog *PubLog) PublicString(public string) {
	plog.Logger.Info().Str("", public).Msg("")
}
