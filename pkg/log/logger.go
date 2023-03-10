/*
Package log contains the CLI logger, abstracting away any major detail about the underlying logger.
*/
package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger LoggerI = &Logger{zap.NewNop().Sugar()}

// LoggerI abstracts the logger
type LoggerI interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

// Logger is the actual logger
type Logger struct {
	*zap.SugaredLogger
}

// Debug prints the msg, plus any additional even number of arguments,
// only if the log level is higher than debug
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Debugw(msg, args...)
}

// Info prints the msg, plus any additional even number of arguments,
// only if the log level is higher than info
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Infow(msg, args...)
}

// Warn prints the msg, plus any additional even number of arguments,
// only if the log level is higher than warn
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Warnw(msg, args...)
}

// Error prints the msg, plus any additional even number of arguments,
// only if the log level is higher than error
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Errorw(msg, args...)
}

// Debug prints the msg, plus any additional even number of arguments,
// only if the log level is higher than debug
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

// Info prints the msg, plus any additional even number of arguments,
// only if the log level is higher than info
func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

// Warn prints the msg, plus any additional even number of arguments,
// only if the log level is higher than warn
func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

// Error prints the msg, plus any additional even number of arguments,
// only if the log level is higher than error
func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

// DefaultLogger configures the logger, setting the given level
func DefaultLogger(level Level, encoding Encoding) (*Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level.Level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: string(encoding),
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	prodLogger, err := config.Build(zap.AddCallerSkip(2), zap.WithCaller(true))
	if err != nil {
		return nil, fmt.Errorf("unable to create new logger: %w", err)
	}
	return &Logger{prodLogger.Sugar()}, nil
}

// SetLogger sets the global logger
func SetLogger(l LoggerI) {
	logger = l
}
