package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func zapTimeFormat(format string) zapcore.TimeEncoder {
	if format == "" {
		format = _defaultTimeFormat
	}
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}

func getZapEncoder(isJSON bool, timeFormat string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapTimeFormat(timeFormat)
	encoderConfig.TimeKey = "time"
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(l string) zapcore.Level {
	ll := strings.ToLower(l)
	level, _ := zapcore.ParseLevel(ll)
	return level
}

func NewZapLogger(logLevel string, logJson bool, timeFormat string) Logger {
	level := getZapLevel(logLevel)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(getZapEncoder(logJson, timeFormat), writer, level)

	logger := zap.New(core,
		zap.AddCallerSkip(1),
		zap.AddCaller(),
	).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

func (l *zapLogger) WithError(err error) Logger {
	return l.WithFields(Fields{"error": err})
}

func (l *zapLogger) Logf(level Level, format string, args ...interface{}) {
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
