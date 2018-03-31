package log

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	gaelog "google.golang.org/appengine/log"
)

var (
	ctx context.Context
)

func Init(r *http.Request) {
	ctx = appengine.NewContext(r)
}

func Debug(format string, args ...interface{}) {
	gaelog.Debugf(ctx, format, args...)
}

func Info(format string, args ...interface{}) {
	gaelog.Infof(ctx, format, args...)
}

func Warning(format string, args ...interface{}) {
	gaelog.Warningf(ctx, format, args...)
}

func Error(format string, args ...interface{}) {
	gaelog.Errorf(ctx, format, args...)
}

func Critical(format string, args ...interface{}) {
	gaelog.Criticalf(ctx, format, args...)
}
