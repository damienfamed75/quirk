package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CustomLogger is a customizable logging.Logger where you can choose the level
// and the zapcore encoder configuration.
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

// Info logs at an info level.
func (l *CustomLogger) Info(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Info(msg, fields...)
}

// Debug logs at an debug level.
func (l *CustomLogger) Debug(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Debug(msg, fields...)
}

// Warn warns the client.
func (l *CustomLogger) Warn(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Warn(msg, fields...)
}

// Error logs at an error level.
func (l *CustomLogger) Error(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Error(msg, fields...)
}

// Fatal logs at a fatal level and exits.
func (l *CustomLogger) Fatal(msg string, iFields ...interface{}) {
	fields := interfaceToZapField(iFields...)
	l.logger.Fatal(msg, fields...)
}

// interfaceToZapField takes the interfaces passed in and type asserts them
// into a zap.Field and returns a slice.
func interfaceToZapField(iFields ...interface{}) (fields []zapcore.Field) {
	for i := 0; i < len(iFields); i++ {
		fields = append(fields, iFields[i].(zapcore.Field))
	}
	return
}
