package quirk

import (
	"fmt"
	"testing"
)

func BenchmarkAeshash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = aeshash(string(15 * i))
	}
}

func BenchmarkAeshashOld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = aeshashold(string(15 * i))
	}
}

func BenchmarkAeshashVerbose(b *testing.B) {
	var hashes = make(map[uint64]bool)
	var dupes int

	for i := 0; i < b.N; i++ {
		v := aeshash(string(15 * i))
		if _, ok := hashes[v]; ok {
			dupes++
		} else {
			hashes[v] = true
		}
	}

	fmt.Println(dupes)
}
