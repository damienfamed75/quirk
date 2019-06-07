package quirk

import (
	"github.com/damienfamed75/quirk/logging"
	"go.uber.org/zap/zapcore"
)

func NewNilLogger() logging.Logger {
	return logging.NewNilLogger()
}

func NewCustomLogger(level []byte, config zapcore.EncoderConfig) logging.Logger {
	return logging.NewCustomLogger(level, config)
}

func NewDebugLogger() logging.Logger {
	return logging.NewDebugLogger()
}
