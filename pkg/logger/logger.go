package logger

import (
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=logger.go -destination=mocks/logger_mock.go
type Ilogger interface {
	Info(message string)
	Warn(message string)
	Error(err error)
}

type lLogger struct {
	Log *logrus.Logger
}

func NewLogrusLogger() Ilogger {
	log := logrus.New()
	log.SetFormatter(
		&logrus.TextFormatter{
			FullTimestamp: true,
		},
	)
	return &lLogger{
		Log: log,
	}
}

func (l *lLogger) Info(message string) {
	l.Log.Info(message)
}
func (l *lLogger) Warn(message string) {
	l.Log.Warn(message)
}
func (l *lLogger) Error(err error) {
	l.Log.Error(err)
}
