package quirk

import "testing"

func BenchmarkIncrementAuto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		increment(auto)
	}
}

func BenchmarkIncrementIncrementor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		increment(incrementor)
	}
}

func BenchmarkIncrementHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		increment(hash)
	}
}
