package logger

import (
	"github.com/sirupsen/logrus"
)

// Implementation of grpclog.LoggerV2 contract
// for using logger.Logger instead of default gRPC logger.

// stub: l *logrusLogger grpclog.LoggerV2

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (l *LogrusLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *LogrusLogger) Infoln(args ...interface{}) {
	l.log.Infoln(args...)
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (l *LogrusLogger) Warning(args ...interface{}) {
	l.log.Warn(args...)
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (l *LogrusLogger) Warningln(args ...interface{}) {
	l.log.Warnln(args...)
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (l *LogrusLogger) Warningf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (l *LogrusLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (l *LogrusLogger) Errorln(args ...interface{}) {
	l.log.Errorln(args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *LogrusLogger) Fatalln(args ...interface{}) {
	l.log.Fatalln(args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (l *LogrusLogger) V(level int) bool {
	const (
		// infoLog indicates Info severity.
		infoLog int = iota
		// warningLog indicates Warning severity.
		warningLog
		// errorLog indicates Error severity.
		errorLog
		// fatalLog indicates Fatal severity.
		fatalLog
	)
	// logrus have levels from info(4) to fatal(1)
	// grpclog have levels from info(0) to fatal(3)
	var lev logrus.Level
	switch level {
	case infoLog: // info
		lev = logrus.InfoLevel
	case warningLog: // warn
		lev = logrus.WarnLevel
	case errorLog: // error
		lev = logrus.ErrorLevel
	case fatalLog: // fatal
		lev = logrus.FatalLevel
	default:
		return false
	}

	return l.log.Logger.IsLevelEnabled(lev)
}
