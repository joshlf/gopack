// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"testing"

	"github.com/joshlf/gopack"
)

func BenchmarkPackOverflowCheckBaseline(b *testing.B) {
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

func BenchmarkPackOverflowCheckOverflow(b *testing.B) {
	type typ struct {
		F1 uint8 `gopack:"7"`
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkUnpackOverflowCheckBaseline(b *testing.B) {
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

func BenchmarkUnpackOverflowCheckOverflow(b *testing.B) {
	type typ struct {
		F1 uint8 `gopack:"7"`
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}
