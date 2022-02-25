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

// implementation: l *logrusLogger grpclog.LoggerV2

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (l *logrusLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *logrusLogger) Infoln(args ...interface{}) {
	l.log.Infoln(args...)
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (l *logrusLogger) Warning(args ...interface{}) {
	l.log.Warn(args...)
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (l *logrusLogger) Warningln(args ...interface{}) {
	l.log.Warnln(args...)
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (l *logrusLogger) Warningf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (l *logrusLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (l *logrusLogger) Errorln(args ...interface{}) {
	l.log.Errorln(args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *logrusLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.log.Fatalln(args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (l *logrusLogger) V(level int) bool {
	// logrus have levels from info(4) to fatal(1)
	// grpclog have levels from info(0) to fatal(3)
	var lev logrus.Level
	switch level {
	case 0: // info
		lev = logrus.InfoLevel
	case 1: // warn
		lev = logrus.WarnLevel
	case 2: // error
		lev = logrus.ErrorLevel
	case 3: // fatal
		lev = logrus.FatalLevel
	default:
		return false
	}

	return l.log.Logger.IsLevelEnabled(lev)
}
