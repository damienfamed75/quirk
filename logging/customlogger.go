package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CustomLogger
type CustomLogger struct {
	logger *zap.Logger
}

var _ Logger = &CustomLogger{}

// NewCustomLogger returns a custom logging
// object for the Classy service to use.
func NewCustomLogger(level []byte, config zapcore.EncoderConfig) *CustomLogger {
	logLevel := zap.NewAtomicLevel()
	logLevel.UnmarshalText(level)

	return &CustomLogger{
		logger: zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.Lock(os.Stdout), logLevel)),
	}
}

func (l *CustomLogger) Info(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Info(msg, fields...)
}

// Debug logs nothing.
func (l *CustomLogger) Debug(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Debug(msg, fields...)
}

// Warn warns the client.
func (l *CustomLogger) Warn(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Warn(msg, fields...)
}

// Error logs nothing.
func (l *CustomLogger) Error(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Error(msg, fields...)
}

// Fatal logs nothing.
func (l *CustomLogger) Fatal(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Fatal(msg, fields...)
}

func interfaceToZapField(iFields ...interface{}) (fields []zapcore.Field) {
	for i := 0; i < len(iFields); i++ {
		fields = append(fields, iFields[i].(zapcore.Field))
	}
	return
}
