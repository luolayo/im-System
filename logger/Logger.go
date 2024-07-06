package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Level Define log level type
type Level = zapcore.Level

// Define specific log level constants
const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
)

// Logger struct, contains a SugaredLogger instance
type Logger struct {
	sugarLogger *zap.SugaredLogger
}

// Custom time format encoder
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// NewLogger function, creates and returns a Logger instance
func NewLogger(level Level) *Logger {

	// Set log level
	al := zap.NewAtomicLevelAt(level)

	// Custom EncoderConfig
	config := zapcore.EncoderConfig{
		LevelKey:       "level",                          // Key name for log level
		TimeKey:        "time",                           // Key name for time
		NameKey:        "logger",                         // Key name for logger
		CallerKey:      "caller",                         // Key name for caller
		MessageKey:     "msg",                            // Key name for log message
		StacktraceKey:  "stacktrace",                     // Key name for stack trace
		LineEnding:     zapcore.DefaultLineEnding,        // Line ending character
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Output log level in uppercase with color
		EncodeTime:     customTimeEncoder,                // Custom time format
		EncodeDuration: zapcore.StringDurationEncoder,    // Duration format
		EncodeCaller:   zapcore.ShortCallerEncoder,       // Short file path and line number
		EncodeName:     zapcore.FullNameEncoder,          // Full name encoder
	}

	if _, err := os.Stat("../log"); os.IsNotExist(err) {
		err := os.Mkdir("../log", 0755)
		if err != nil {
			return nil
		} // Create log directory
	}
	logFile, err := os.OpenFile("../log/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// Create the log core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),                                                 // Use ConsoleEncoder to output logs in a custom format
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logFile)), // Output to both console and specified output
		al,
	)

	// Create the Logger instance
	return &Logger{
		sugarLogger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar(),
	}
}

// Debug level log method
func (l *Logger) Debug(msg string, args ...any) {
	l.sugarLogger.Debugf(msg, args...)
}

// Info level log method
func (l *Logger) Info(msg string, args ...any) {
	l.sugarLogger.Infof(msg, args...)
}

// Warn level log method
func (l *Logger) Warn(msg string, args ...any) {
	l.sugarLogger.Warnf(msg, args...)
}

// Error level log method
func (l *Logger) Error(msg string, args ...any) {
	l.sugarLogger.Errorf(msg, args...)
}

// Sync method to flush the buffer
func (l *Logger) Sync() {
	_ = l.sugarLogger.Sync()
}
