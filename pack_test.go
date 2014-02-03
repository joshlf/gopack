// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func TestPackUint64(t *testing.T) {
	testPackUint64(t, math.MaxUint64, 0, 64)
	testPackUint64(t, math.MaxUint64, 1, 64)
	testPackUint64(t, math.MaxUint64, 7, 64)

	// Make sure the test is
	// deterministic
	rand.Seed(15830)
	for i := 0; i < 1000*1000; i++ {
		width := 1 + (rand.Uint32() & 0x3F) // Restrict to range [1, 64]
		lsb := rand.Uint32() & 0x3F         // Restrict to range [0, 63]
		if width+lsb > 64 {
			continue
		}
		val := (uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)) & nOnes(int(width)) // Restrict to range [0, 2^width)
		testPackUint64(t, val, uint8(lsb), uint8(width))
	}

	// Test to make sure we're not pulling in other
	// values packed next to the target value
	b := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	val := unpackUint64(b, 8, 8)
	if val != 0xFF {
		t.Fatalf("Expected %v; got %v (width: %v, lsb: %v)", 0xFF, val, 8, 8)
	}
	val = unpackUint64(b, 63, 8)
	if val != 0xFF {
		t.Fatalf("Expected %v; got %v (width: %v, lsb: %v)", 0xFF, val, 8, 8)
	}
}

func testPackUint64(t *testing.T, val uint64, lsb, width uint8) {
	var b []byte
	if lsb+width > 64 {
		b = make([]byte, 9)
	} else {
		b = make([]byte, 8)
	}
	packUint64(b, val, lsb, width)
	val2 := unpackUint64(b, lsb, width)
	if val2 != val {
		fmt.Printf("%064b\n%064b\n", val, val2)
		fmt.Println(b)
		t.Fatalf("Expected %v; got %v (width: %v, lsb: %v)", val, val2, width, lsb)
	}
}

func TestPackInt64(t *testing.T) {
	testPackInt64(t, math.MaxInt64, 0, 64)
	testPackInt64(t, math.MaxInt64, 1, 64)
	testPackInt64(t, math.MaxInt64, 7, 64)

	// Make sure the test is
	// deterministic
	rand.Seed(14271)
	for i := 0; i < 1000*1000; i++ {
		width := 1 + (rand.Uint32() & 0x3F) // Restrict to range [1, 64]
		lsb := rand.Uint32() & 0x3F         // Restrict to range [0, 63]
		// Restrict width to [2, 63] because of
		// restrictions imposed by rand.Int63n
		if width+lsb > 64 || width == 1 || width == 64 {
			continue
		}

		val := rand.Int63n(1 << (width - 1))
		// Int63n only returns positive values
		if rand.Int()%2 == 0 {
			val *= -1
		}
		testPackInt64(t, val, uint8(lsb), uint8(width))
	}

	// The above code does not test for width = 1 or width = 64
	for lsb := 0; lsb < 64; lsb++ {
		testPackInt64(t, -1, uint8(lsb), 1)
	}
	for i := 0; i < 1000*1000; i++ {
		// Test every possible value (this is a poor man's rand.Uint64)
		uval := uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
		val := *(*int64)(unsafe.Pointer(&uval))
		testPackInt64(t, val, 0, 64)
	}

	// Test to make sure we're not pulling in other
	// values packed next to the target value
	b := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	val := unpackInt64(b, 8, 8)
	if val != -1 {
		t.Fatalf("Expected %v; got %v (width: %v, lsb: %v)", -1, val, 8, 8)
	}
	val = unpackInt64(b, 63, 8)
	if val != -1 {
		t.Fatalf("Expected %v; got %v (width: %v, lsb: %v)", -1, val, 8, 8)
	}
}

func testPackInt64(t *testing.T, val int64, lsb, width uint8) {
	var b []byte
	if lsb+width > 64 {
		b = make([]byte, 9)
	} else {
		b = make([]byte, 8)
	}
	packInt64(b, val, lsb, width)
	val2 := unpackInt64(b, lsb, width)
	if val2 != val {
		fmt.Printf("%064b\n%064b\n", val, val2)
		fmt.Println(b)
		t.Fatalf("Expected %v; got %v (width: %v, lsb: %v)", val, val2, width, lsb)
	}
}

func BenchmarkPackUint64(b *testing.B) {
	bytes := make([]byte, 9)
	rand.Seed(time.Now().UnixNano())

	// Poor man's rand.Uint64
	u := uint64(rand.Uint32()) | uint64(rand.Uint32())<<32
	width := 1 + uint8(rand.Uint32()&0x3F) // Restrict to range [0, 63]
	lsb := uint8(rand.Uint32() & 0x3F)     // Restrict to range [0, 63]
	for lsb+width > 64 {
		width = 1 + uint8(rand.Uint32()&0x3F)
		lsb = uint8(rand.Uint32() & 0x3F)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packUint64(bytes, u, lsb, width)
	}
}

// NOTE: This test takes about 2 minutes
// on a 2.4GHz i7.
func BenchmarkPackUint64Branch(b *testing.B) {
	bytes := make([]byte, 9)
	rand.Seed(time.Now().UnixNano())

	// Randomly generate values to throw off the branch
	// predictor. Precompute these values so that the
	// computation time doesn't affect the measured time.
	us := make([]uint64, b.N)
	widths := make([]uint8, b.N)
	lsbs := make([]uint8, b.N)
	for i := range us {
		// Poor man's rand.Uint64
		us[i] = uint64(rand.Uint32()) | uint64(rand.Uint32())<<32
		width := 1 + uint8(rand.Uint32()&0x3F) // Restrict to range [0, 63]
		lsb := uint8(rand.Uint32() & 0x3F)     // Restrict to range [0, 63]
		for lsb+width > 64 {
			width = 1 + uint8(rand.Uint32()&0x3F)
			lsb = uint8(rand.Uint32() & 0x3F)
		}
		widths[i] = width
		lsbs[i] = lsb
	}
	b.ResetTimer()
	for i := range us {
		packUint64(bytes, us[i], lsbs[i], widths[i])
	}
}

// This function uses the same values every time
// (so it doesn't test the branch predictor), but
// uses the same overhead of slices that the
// branch prediction test uses so that this test
// can be used as a baseline against which to compare
// the branch prediction test
func BenchmarkPackUint64Slice(b *testing.B) {
	bytes := make([]byte, 9)
	rand.Seed(time.Now().UnixNano())

	// Randomly generate values to throw off the branch
	// predictor. Precompute these values so that the
	// computation time doesn't affect the measured time.
	us := make([]uint64, b.N)
	widths := make([]uint8, b.N)
	lsbs := make([]uint8, b.N)
	// Poor man's rand.Uint64
	u := uint64(rand.Uint32()) | uint64(rand.Uint32())<<32
	width := 1 + uint8(rand.Uint32()&0x3F) // Restrict to range [0, 63]
	lsb := uint8(rand.Uint32() & 0x3F)     // Restrict to range [0, 63]
	for lsb+width > 64 {
		width = 1 + uint8(rand.Uint32()&0x3F)
		lsb = uint8(rand.Uint32() & 0x3F)
	}
	for i := range us {
		us[i] = u
		widths[i] = width
		lsbs[i] = lsb
	}
	b.ResetTimer()
	for i := range us {
		packUint64(bytes, us[i], lsbs[i], widths[i])
	}
}
