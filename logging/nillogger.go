package logging

// NilLogger is the default logger used for the Quirk client.
// This logger will not print anything out.
type NilLogger struct{}

var _ Logger = &NilLogger{}

// NewNilLogger returns a nil logging
// object for the Quirk client to use.
func NewNilLogger() *NilLogger {
	return &NilLogger{}
}

// Info does nothing.
func (*NilLogger) Info(string, ...interface{}) {}

// Debug logs nothing.
func (*NilLogger) Debug(string, ...interface{}) {}

// Error logs nothing.
func (*NilLogger) Error(string, ...interface{}) {}

// Warn does nothing.
func (*NilLogger) Warn(string, ...interface{}) {}

// Fatal logs nothing.
func (*NilLogger) Fatal(string, ...interface{}) {}
