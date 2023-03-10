package log

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	// JSONEncoding should be used to set the log output to json format
	JSONEncoding Encoding = "json"
	// ConsoleEncoding should be used to set the log output to console format
	ConsoleEncoding Encoding = "console"
)

// Encoding represents the possible formats for logs
type Encoding string

// String is needed to be able to parse it correctly as flag
func (e *Encoding) String() string {
	return string(*e)
}

// Set is needed to be able to parse it correctly as flag
func (e *Encoding) Set(s string) error {
	switch s {
	case string(JSONEncoding):
		*e = JSONEncoding
	case string(ConsoleEncoding):
		*e = ConsoleEncoding
	default:
		return fmt.Errorf("unsupported encoding format: %s", s)
	}
	return nil
}

// Type is needed to be able to parse it correctly as flag
func (e Encoding) Type() string {
	return "string"
}

// Level encodes the log level, can not be a type alias because we need to implement the same interfaces
type Level struct {
	zapcore.Level
}

// Type returns "string" as type of the Level, required to correctly parse it
func (l Level) Type() string {
	return "string"
}

// GetDefaultLevel returns the default log level
func GetDefaultLevel() Level {
	return Level{zapcore.InfoLevel}
}
