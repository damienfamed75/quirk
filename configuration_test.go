package quirk

import (
	"testing"

	. "github.com/franela/goblin"
)

func BenchmarkWithLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = WithLogger(NewDebugLogger())
	}
}

func BenchmarkWithPredicateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = WithPredicateKey("key")
	}
}

func TestWithLogger(t *testing.T) {
	g := Goblin(t)

	g.Describe("Logger", func() {
		dl := NewDebugLogger()
		g.It("should equal a debug logger", func() {
			g.Assert(NewClient(WithLogger(dl)).logger).
				Equal(dl)
		})
	})
}

func TestWithPredicateKey(t *testing.T) {
	g := Goblin(t)

	g.Describe("Predicate key", func() {
		key := "customKey"
		g.It("should equal customKey", func() {
			g.Assert(NewClient(WithPredicateKey(key)).predicateKey).
				Equal(key)
		})
	})
}

func TestWithMaxWorkerCount(t *testing.T) {
	g := Goblin(t)

	g.Describe("Client max worker count", func() {
		count := 2
		g.It("should equal count", func() {
			g.Assert(NewClient(WithMaxWorkerCount(count)).maxWorkerCount).
				Equal(count)
		})
	})
}
