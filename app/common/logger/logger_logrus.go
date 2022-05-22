package logger

import (
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	log *logrus.Entry
}

func NewLogrusLogger(logLevel string, logJson bool, timeFormat string, logLines bool) *LogrusLogger {
	if timeFormat == "" {
		timeFormat = _defaultTimeFormat
	}

	level, _ := logrus.ParseLevel(logLevel)

	logger := logrus.New()
	if logJson {
		formatter := &logrus.JSONFormatter{}
		formatter.TimestampFormat = timeFormat
		logger.SetFormatter(formatter)
	} else {
		formatter := &logrus.TextFormatter{}
		formatter.TimestampFormat = timeFormat
		formatter.FullTimestamp = true
		formatter.ForceColors = true
		logger.SetFormatter(formatter)
	}
	logger.SetLevel(level)
	logger.SetReportCaller(logLines)

	return &LogrusLogger{
		log: logrus.NewEntry(logger),
	}
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *LogrusLogger) WithFields(fields Fields) Logger {
	newLogger := l.log.WithFields(logrus.Fields(fields))
	return &LogrusLogger{newLogger}
}

func (l *LogrusLogger) WithError(err error) Logger {
	return l.WithFields(Fields{"error": err})
}

func (l *LogrusLogger) Logf(level Level, format string, args ...interface{}) {
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
