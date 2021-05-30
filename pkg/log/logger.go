package log

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})

	Info(args ...interface{})

	Error(args ...interface{})

	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Errorf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

func New() Logger {
	l, _ := zap.NewProduction()
	return newWithZap(l)
}

func newWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}
