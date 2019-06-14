package logging

import (
	"testing"

	. "github.com/franela/goblin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggers(t *testing.T) {
	g := Goblin(t)

	g.Describe("Debug logger", func() {
		l := NewDebugLogger()
		if l == nil {
			g.Fail("logger returned is nil")
		}

		g.It("should not panic", func() {
			l.Info("msg", zap.Skip())
			l.Debug("msg")
			l.Error("msg")
			l.Warn("msg")
		})
	})

	g.Describe("Custom logger", func() {
		l := NewCustomLogger([]byte("dpanic"), zapcore.EncoderConfig{})
		if l == nil {
			g.Fail("logger returned is nil")
		}

		g.It("should not panic", func() {
			l.Info("msg", zap.Skip())
			l.Debug("msg")
			l.Error("msg")
			l.Warn("msg")
		})
	})

	g.Describe("Nil logger", func() {
		g.Assert(NewNilLogger()).Equal(&NilLogger{})
	})
}
