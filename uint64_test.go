// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func TestUnsigned(t *testing.T) {
	// Make sure the test is
	// deterministic
	rand.Seed(27953)

	var u uint64
	for i := 0; i < 1000*1000; i++ {
		width, lsb := randWidthLSBPair()
		val := randUint64() & nOnes(int(width)) // Restrict to range [0, 2^width)
		u = PackUnsigned(u, val, uint8(lsb), uint8(width))
		val2 := UnpackUnsigned(u, uint8(lsb), uint8(width))
		if val2 != val {
			t.Fatalf("Expected 0x%X; got 0x%X (width: %v, lsb: %v)", val, val2, width, lsb)
		}
	}
}

func TestSigned(t *testing.T) {
	// Make sure the test is
	// deterministic
	rand.Seed(21459)

	var u uint64
	for i := 0; i < 1000*1000; i++ {
		width, lsb := randWidthLSBPair()
		// Restrict width to [2, 63] because of
		// restrictions imposed by rand.Int63n
		if width == 1 || width == 64 {
			continue
		}

		val := rand.Int63n(1 << (width - 1))
		// Int63n only returns positive values
		if rand.Int()%2 == 0 {
			val *= -1
		}
		u = PackSigned(u, val, uint8(lsb), uint8(width))
		val2 := UnpackSigned(u, uint8(lsb), uint8(width))
		if val2 != val {
			t.Fatalf("Expected 0x%X; got 0x%X (width: %v, lsb: %v)", val, val2, width, lsb)
		}
	}

	// The above code does not test for width = 1 or width = 64
	for lsb := 0; lsb < 64; lsb++ {
		u = PackSigned(u, -1, uint8(lsb), 1)
		val := UnpackSigned(u, uint8(lsb), 1)
		if val != -1 {
			t.Fatalf("Expected 0x%X; got 0x%X (width: %v, lsb: %v)", -1, val, 1, lsb)
		}
	}
	for i := 0; i < 1000*1000; i++ {
		uval := randUint64()
		val := *(*int64)(unsafe.Pointer(&uval))
		u = PackSigned(u, val, 0, 64)
		val2 := UnpackSigned(u, 0, 64)
		if val2 != val {
			t.Fatalf("Expected 0x%X; got 0x%X (width: %v, lsb: %v)", val, val2, 64, 0)
		}
	}
}

func BenchmarkPackUnsigned(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	width, lsb := randWidthLSBPair()
	val := randUint64()
	var u uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PackUnsigned(u, val, lsb, width)
	}
}

func BenchmarkUnpackUnsigned(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	width, lsb := randWidthLSBPair()
	u := randUint64()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnpackUnsigned(u, lsb, width)
	}
}

func BenchmarkPackSigned(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	width, lsb := randWidthLSBPair()
	val := randInt64()
	var u uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PackSigned(u, val, lsb, width)
	}
}

func BenchmarkUnpackSigned(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	width, lsb := randWidthLSBPair()
	u := randUint64()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnpackSigned(u, lsb, width)
	}
}
