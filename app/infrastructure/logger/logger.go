package logger

const (
	_defaultTimeFormat = "2006-01-02T15:04:05.000Z0700"
)

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

type Level string

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = "debug"

	// InfoLevel is the default logging priority.
	InfoLevel Level = "info"

	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel Level = "warn"

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel Level = "error"

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = "fatal"

	// PanicLevel logs a message
	PanicLevel Level = "panic"
)

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger

	WithError(err error) Logger

	Logf(level Level, format string, args ...interface{})
}

func With(l Logger, a ...interface{}) Logger {
	if len(a)%2 != 0 {
		return l
	}
	fields := make(Fields)
	for i := 0; i < len(a); i += 2 {
		s, ok := a[i].(string)
		if !ok {
			return l
		}
		fields[s] = a[i+1]
	}
	return l.WithFields(fields)
}
