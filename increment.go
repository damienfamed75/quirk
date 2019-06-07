package quirk

import "sync/atomic"

var (
	charIncrementor uint64 = 'a'
	uidIncrementor  uint64
)

// GetLastUID will retrieve the newest value set to the
// uidIncrementor. This can be used when swapping Quirk clients.
func GetLastUID() uint64 {
	return atomic.LoadUint64(&uidIncrementor)
}

// setStartUID swaps the value of the default starting value
// when incrementing the UID to the new given value.
func setStartUID(i uint64) {
	atomic.SwapUint64(&uidIncrementor, i)
}

func resetIncrementor(incrementor *uint64) {
	atomic.SwapUint64(incrementor, 'a')
}

func increment(mode rdfMode) uint64 {
	var incrementor = &uidIncrementor
	if mode == auto {
		incrementor = &charIncrementor
	}

	atomic.AddUint64(incrementor, 1)

	return atomic.LoadUint64(incrementor)
}
