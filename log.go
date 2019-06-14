package quirk

import (
	"github.com/damienfamed75/quirk/logging"
	"go.uber.org/zap/zapcore"
)

// NewNilLogger returns the *logging.NilLogger logger which doesn't log anything.
func NewNilLogger() logging.Logger {
	return logging.NewNilLogger()
}

// NewCustomLogger returns a *logging.CustomLogger with the desired level
// and encoder configuration that may be passed in.
func NewCustomLogger(level []byte, config zapcore.EncoderConfig) logging.Logger {
	return logging.NewCustomLogger(level, config)
}

// NewDebugLogger is a debug level zap logger that can be used when testing.
func NewDebugLogger() logging.Logger {
	return logging.NewDebugLogger()
}
