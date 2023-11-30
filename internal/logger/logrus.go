package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	log *logrus.Logger
}

func newLogrusLogger(conf Configuration) Logger {
	var writers []io.Writer
	logger := logrus.StandardLogger()

	if conf.EnableConsole {
		writers = append(writers, os.Stdout)
	}
	if conf.EnableFile && len(conf.FileLocation) > 0 {
		dir := filepath.Dir(conf.FileLocation)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			logger.Errorf("creating log file, %v", err)
		} else {
			file, err := os.OpenFile(conf.FileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o664)
			if err != nil {
				logger.Errorf("error opening file %v", err)
			} else {
				writers = append(writers, file)
			}
		}
	}

	var formatter logrus.Formatter
	formatter = &logrus.TextFormatter{}

	if conf.ConsoleJSONFormat {
		formatter = &logrus.JSONFormatter{}
	}

	logger = &logrus.Logger{
		Out:       io.MultiWriter(writers...),
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     getLogrusLevel(conf.ConsoleLevel),
	}

	return &logrusLogger{log: logger}
}

func getLogrusLevel(level Level) logrus.Level {
	switch level {
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case DebugLevel:
		return logrus.DebugLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case TraceLevel:
		return logrus.TraceLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logrusLogger) Tracef(format string, args ...interface{}) {
	l.log.Tracef(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l *logrusLogger) DebugWithFields(fields map[string]interface{}) {
	l.log.WithFields(fields).Debug()
}
