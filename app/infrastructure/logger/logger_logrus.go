package logger

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	log *logrus.Entry
}

func NewLogrusLogger(LogLevel string, LogJson bool, TimeFormat string, LogLines bool) Logger {
	if TimeFormat == "" {
		TimeFormat = _defaultTimeFormat
	}

	level, _ := logrus.ParseLevel(LogLevel)

	logger := logrus.New()
	if LogJson {
		formatter := &logrus.JSONFormatter{}
		formatter.TimestampFormat = TimeFormat
		logger.SetFormatter(formatter)
	} else {
		formatter := &logrus.TextFormatter{}
		formatter.TimestampFormat = TimeFormat
		formatter.FullTimestamp = true
		formatter.ForceColors = true
		logger.SetFormatter(formatter)
	}
	logger.SetLevel(level)
	logger.SetReportCaller(LogLines)

	return &logrusLogger{
		log: logrus.NewEntry(logger),
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

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	newLogger := l.log.WithFields(logrus.Fields(fields))
	return &logrusLogger{newLogger}
}

func (l *logrusLogger) WithError(err error) Logger {
	return l.WithFields(Fields{"error": err})
}

func (l *logrusLogger) Logf(level Level, format string, args ...interface{}) {
	switch level {
	case DebugLevel:
		l.Debugf(format, args...)
	case InfoLevel:
		l.Infof(format, args...)
	case WarnLevel:
		l.Warnf(format, args...)
	case ErrorLevel:
		l.Errorf(format, args...)
	case FatalLevel:
		l.Fatalf(format, args...)
	case PanicLevel:
		l.Panicf(format, args...)
	}
}
