package quirk

var masks [32]uint64
var shifts [32]uint64
var keysched [32]byte

func hashb()

//go:noescape
//go:nosplit
func aeshash(s string) uint64
