package logger

import "strings"

// Level represents the log level.
type Level string

const (
	// DebugLevel has verbose message
	DebugLevel Level = "debug"
	// InfoLevel is default log level
	InfoLevel Level = "info"
	// WarnLevel is for logging messages about possible issues
	WarnLevel Level = "warn"
	// TraceLevel is for finer-grained informational events than the Debug.
	TraceLevel Level = "trace"
	// ErrorLevel is for logging errors
	ErrorLevel Level = "error"
	// FatalLevel is for logging fatal messages. The system shutdown after logging the message.
	FatalLevel Level = "fatal"
)

var log Logger

func init() {
	defaultConf := Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      InfoLevel,
		EnableFile:        false,
	}
	// log = newZapLogger(defaultConf)
	log = newLogrusLogger(defaultConf)
}

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Tracef(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	DebugWithFields(fields map[string]interface{})
}

// Configuration contains all parameters for the logger.
type Configuration struct {
	// Console
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      Level

	// Files
	EnableFile     bool
	FileJSONFormat bool
	FileLevel      Level
	FileLocation   string
}

// GetLevel returns the log level by string.
func GetLevel(l string) Level {
	switch strings.ToLower(l) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "trace":
		return TraceLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// NewLogger returns an instance of logger
func NewLogger(config Configuration) {
	// log = newZapLogger(config)
	log = newLogrusLogger(config)
}

// Info logs a message at info level. The message includes any fields passed.
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn logs a message at warn level. The message includes any fields passed.
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Debug logs a message at debug level. The message includes any fields passed.
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Error logs a message at error level. The message includes any fields passed.
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Tracef logs a message at trace level. The message includes any fields passed.
func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args)
}

// Fatal logs a message at fatal level. The message includes any fields passed.
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panic logs a message at panic level. The message includes any fields passed.
func Panic(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// DebugWithFields adds a struct of fields to the log entry.
func DebugWithFields(fields map[string]interface{}) Logger {
	log.DebugWithFields(fields)
	return log
}
