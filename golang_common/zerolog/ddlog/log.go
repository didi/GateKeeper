package ddlog

import (
	"context"

	"github.com/didi/gatekeeper/golang_common/zerolog"
)

var (
	//Logger is defualt log ,and print to stdout, if want format context, use RegisterContextFormat
	Logger *DiLogHandle
)

func init() {
	diLogger := zerolog.New(zerolog.NewStdoutWriter())
	diLogger = diLogger.Level(zerolog.DebugLevel)
	Logger = &DiLogHandle{Logger: diLogger}
}

func Debugf(ctx context.Context, tag string, format string, args ...interface{}) {
	Logger.Debugf(ctx, tag, format, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Infof(ctx context.Context, tag string, format string, args ...interface{}) {
	Logger.Infof(ctx, tag, format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Warnf(ctx context.Context, tag string, format string, args ...interface{}) {
	Logger.Warnf(ctx, tag, format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Errorf(ctx context.Context, tag string, format string, args ...interface{}) {
	Logger.Errorf(ctx, tag, format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Fatalf(ctx context.Context, tag string, format string, args ...interface{}) {
	Logger.Fatalf(ctx, tag, format, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}
