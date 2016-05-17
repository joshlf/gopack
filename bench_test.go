// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func clearCache() {
	layoutCache.Lock()
	layoutCache.m = make(map[reflect.Type]cachedLayout)
	layoutCache.Unlock()
}

func BenchmarkBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkBool1Field(b *testing.B) {
	type typ struct {
		F1 bool
	}

	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _ := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, val)
		// p(bytes, val)
	}
}

func BenchmarkBool2Fields(b *testing.B) {
	type typ struct {
		F1, F2 bool
	}

	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _ := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, val)
		// p(bytes, val)
	}
}

func BenchmarkBool4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 bool
	}

	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _ := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, val)
		// p(bytes, val)
	}
}

func BenchmarkBool8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 bool
	}

	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _ := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, val)
		// p(bytes, val)
	}
}

func BenchmarkBool9Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9 bool
	}

	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _ := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 2)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, val)
		// p(bytes, val)
	}
}

func BenchmarkBool16Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13, F14, F15, F16 bool
	}

	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _ := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 2)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, val)
		// p(bytes, val)
	}
}

func BenchmarkLockContention(b *testing.B) {
	clearCache()
	type typ struct {
		F1 uint8
	}

	t := typ{}
	bytes1 := []byte{0}
	bytes2 := []byte{0}

	// Pre-compute handler
	Pack(bytes1, t)

	var wg sync.WaitGroup
	wg.Add(1)

	old := runtime.GOMAXPROCS(2)
	b.ResetTimer()
	go func(t typ, bytes []byte, b *testing.B) {
		for i := 0; i < b.N; i += 2 {
			Pack(bytes2, t)
		}
		wg.Done()
	}(t, bytes2, b)
	for i := 0; i < b.N; i += 2 {
		Pack(bytes1, t)
	}
	wg.Wait()
	runtime.GOMAXPROCS(old)
}

func BenchmarkNoLockContention(b *testing.B) {
	clearCache()
	type typ struct {
		F1 uint8
	}

	t := typ{}
	bytes := []byte{0}
	// Pre-compute handler
	Pack(bytes, t)
	b.ResetTimer()
	for i := 0; i < b.N; i += 2 {
		Pack(bytes, t)
	}
}

func BenchmarkPackExported1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := typ{}
	bytes := []byte{0}
	Pack(bytes, t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(bytes, t)
	}
}

func BenchmarkPackBackend1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := typ{}
	bytes := []byte{0}
	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(&t).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkPackExported2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := typ{}
	bytes := []byte{0, 0}
	Pack(bytes, t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(bytes, t)
	}
}

func BenchmarkPackBackend2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := typ{}
	bytes := []byte{0, 0}
	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(&t).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkPackExported4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := typ{}
	bytes := []byte{0, 0, 0, 0}
	Pack(bytes, t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(bytes, t)
	}
}

func BenchmarkPackBackend4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := typ{}
	bytes := []byte{0, 0, 0, 0}
	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(&t).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkPackExported8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := typ{}
	bytes := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	Pack(bytes, t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(bytes, t)
	}
}

func BenchmarkPackBackend8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := typ{}
	bytes := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	l, _, _ := layoutFor(reflect.ValueOf(typ{}))
	// p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(&t).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkMakeLayout0Fields(b *testing.B) {
	type typ struct {
	}

	v := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeLayout(v)
	}
}

func BenchmarkMakeLayout1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	v := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeLayout(v)
	}
}

func BenchmarkMakeLayout2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	v := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeLayout(v)
	}
}

func BenchmarkMakeLayout4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	v := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeLayout(v)
	}
}

func BenchmarkMakeLayout8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	v := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeLayout(v)
	}
}

func BenchmarkUseEmptyPacker(b *testing.B) {
	type typ struct {
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUsePacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUsePacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUsePacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUsePacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUsePacker8FieldsSigned(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 int8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUseEmptyUnpacker(b *testing.B) {
	type typ struct {
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unpack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUseUnpacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unpack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUseUnpacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unpack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUseUnpacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unpack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUseUnpacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unpack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkUseUnpacker8FieldsSigned(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 int8
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unpack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkNested1Level(b *testing.B) {
	type typ struct {
		F1 struct {
		}
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkNested2Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
			}
		}
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkNested4Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
				F3 struct {
					F4 struct {
					}
				}
			}
		}
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}

func BenchmarkNested8Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
				F3 struct {
					F4 struct {
						F5 struct {
							F6 struct {
								F7 struct {
									F8 struct {
									}
								}
							}
						}
					}
				}
			}
		}
	}

	v := reflect.ValueOf(&typ{}).Elem()
	l, _, _ := makeLayout(v)
	// p, _ := makePackerWrapper(t)
	bytes := make([]byte, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(bytes, l, v)
		// p(bytes, val)
	}
}
