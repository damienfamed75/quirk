package quirk

import (
	"github.com/damienfamed75/yalp"
	"go.uber.org/zap/zapcore"
)

// NewNilLogger returns the *yalp.NilLogger logger which doesn't log anything.
func NewNilLogger() yalp.Logger {
	return yalp.NewNilLogger()
}

// NewCustomLogger returns a *yalp.CustomLogger with the desired level
// and encoder configuration that may be passed in.
func NewCustomLogger(level []byte, config zapcore.EncoderConfig) yalp.Logger {
	return yalp.NewCustomLogger(level, config)
}

// NewDebugLogger is a debug level zap logger that can be used when testing.
func NewDebugLogger() yalp.Logger {
	return yalp.NewDebugLogger()
}
