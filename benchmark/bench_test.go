// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/synful/gopack"
)

/*
	NOTE: We could use a helper function
	for these benchmarks. However, that
	function would have to take an empty
	interface argument. The conversion
	to the interface argument when calling
	the helper function would be performed
	once ahead of time rather than in each
	iteration of the loop (as is the case
	with these benchmarks). Surprisingly,
	this conversion incurs a meaningful
	overhead. Since, in real life, the
	user would probably not perform interface
	conversion ahead of time, this overhead
	is realistic, and should be taken into
	account. Thus, we don't use such a
	helper function here.
*/

// benchmarkUtil takes an example of the type
// (always a value; never a pointer), and the
// length in bytes required to pack the type.
// It registers packers and unpackers for the
// type and generates and returns a random
// instance of the type itself, and of a
// properly-lengthed byte slice.
func benchmarkUtil(example interface{}, numBytes int) (interface{}, []byte) {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, numBytes)

	typ := reflect.TypeOf(example)
	gopack.Pack(bytes, example)
	gopack.Unpack(bytes, reflect.New(typ).Interface())

	randBytes(bytes)
	val := randInstance(typ).Interface()

	return val, bytes
}

func BenchmarkPackBool1Field(b *testing.B) {
	type typ struct {
		F1 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackBool2Fields(b *testing.B) {
	type typ struct {
		F1, F2 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackBool4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackBool8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackBool16Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13, F14, F15, F16 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 2)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkUnpackBool1Field(b *testing.B) {
	type typ struct {
		F1 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackBool2Fields(b *testing.B) {
	type typ struct {
		F1, F2 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackBool4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackBool8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, val)
	}
}

func BenchmarkUnpackBool16Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13, F14, F15, F16 bool
	}
	intface, bytes := benchmarkUtil(typ{}, 2)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkPackUint8_1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackUint8_2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 2)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackUint8_4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 4)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackUint8_8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 8)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackUint8_16Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13, F14, F15, F16 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 16)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkUnpackUint8_1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackUint8_2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 2)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackUint8_4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 4)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackUint8_8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 8)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackUint8_16Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13, F14, F15, F16 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 16)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}
