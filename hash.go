package quirk

import "hash/fnv"

var masks [32]uint64
var shifts [32]uint64
var keysched [32]byte

func hashb()

//go:noescape
//go:nosplit
func aeshash(s string) uint64

func aeshashold(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
