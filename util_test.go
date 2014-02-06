package gopack

import (
	"math/rand"
	"unsafe"
)

func randUint64() uint64 {
	return uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
}

func randUint64Bits(bits uint8) uint64 {
	return randUint64() >> (64 - bits)
}

func randInt64() int64 {
	u := randUint64()
	return *(*int64)(unsafe.Pointer(&u))
}

func randInt64Bits(bits uint8) int64 {
	u := randUint64Bits(bits)
	return (*(*int64)(unsafe.Pointer(&u)) << (64 - bits)) >> (64 - bits)
}

func randBool() bool {
	return rand.Int()%2 == 0
}

func randWidthLSBPair() (uint8, uint8) {
	width := uint8(1 + (rand.Uint32() & 0x3F)) // Restrict to range [1, 64]
	lsb := uint8(rand.Uint32() & 0x3F)         // Restrict to range [0, 63]
	for width+lsb > 64 {
		width = uint8(1 + (rand.Uint32() & 0x3F)) // Restrict to range [1, 64]
		lsb = uint8(rand.Uint32() & 0x3F)         // Restrict to range [0, 63]
	}
	return width, lsb
}

func nOnes(n int) uint64 {
	var u uint64
	for i := 0; i < n; i++ {
		u = (u << 1) | 1
	}
	return u
}
