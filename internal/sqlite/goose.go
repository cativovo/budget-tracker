package sqlite

import (
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type gooseLogger struct {
	base *zap.SugaredLogger
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	g.base.Infof(format, v...)
}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	g.base.Fatalf(format, v...)
}

var _ goose.Logger = (*gooseLogger)(nil)
