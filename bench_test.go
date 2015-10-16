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

func clearCaches() {
	packerCache.Lock()
	defer packerCache.Unlock()
	unpackerCache.Lock()
	defer unpackerCache.Unlock()

	packerCache.m = make(map[reflect.Type]packer)
	unpackerCache.m = make(map[reflect.Type]unpacker)
}

func BenchmarkBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkBool1Field(b *testing.B) {
	type typ struct {
		F1 bool
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkBool2Fields(b *testing.B) {
	type typ struct {
		F1, F2 bool
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkBool4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 bool
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkBool8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 bool
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 1)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkBool9Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9 bool
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 2)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkBool16Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13, F14, F15, F16 bool
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	bytes := make([]byte, 2)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkLockContention(b *testing.B) {
	clearCaches()
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
	clearCaches()
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
	p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
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
	p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
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
	p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
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
	p, _, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
	}
}

func BenchmarkMakeEmptyPacker(b *testing.B) {
	type typ struct {
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakeEmptyUnpacker(b *testing.B) {
	type typ struct {
	}

	t := reflect.TypeOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeUnpackerWrapper(t)
	}
}

func BenchmarkMakeUnpacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := reflect.TypeOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeUnpackerWrapper(t)
	}
}

func BenchmarkMakeUnpacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := reflect.TypeOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeUnpackerWrapper(t)
	}
}

func BenchmarkMakeUnpacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := reflect.TypeOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeUnpackerWrapper(t)
	}
}

func BenchmarkMakeUnpacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := reflect.TypeOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeUnpackerWrapper(t)
	}
}

func BenchmarkUseEmptyPacker(b *testing.B) {
	type typ struct {
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 1)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 2)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 4)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 8)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker8FieldsSigned(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 int8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 8)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUseEmptyUnpacker(b *testing.B) {
	type typ struct {
	}

	t := reflect.TypeOf(&typ{})
	u := makeUnpackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u(bytes, val)
	}
}

func BenchmarkUseUnpacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := reflect.TypeOf(&typ{})
	u := makeUnpackerWrapper(t)
	bytes := make([]byte, 1)
	val := reflect.ValueOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u(bytes, val)
	}
}

func BenchmarkUseUnpacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := reflect.TypeOf(&typ{})
	u := makeUnpackerWrapper(t)
	bytes := make([]byte, 2)
	val := reflect.ValueOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u(bytes, val)
	}
}

func BenchmarkUseUnpacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := reflect.TypeOf(&typ{})
	u := makeUnpackerWrapper(t)
	bytes := make([]byte, 4)
	val := reflect.ValueOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u(bytes, val)
	}
}

func BenchmarkUseUnpacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := reflect.TypeOf(&typ{})
	u := makeUnpackerWrapper(t)
	bytes := make([]byte, 8)
	val := reflect.ValueOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u(bytes, val)
	}
}

func BenchmarkUseUnpacker8FieldsSigned(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 int8
	}

	t := reflect.TypeOf(&typ{})
	u := makeUnpackerWrapper(t)
	bytes := make([]byte, 8)
	val := reflect.ValueOf(&typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u(bytes, val)
	}
}

func BenchmarkNested1Level(b *testing.B) {
	type typ struct {
		F1 struct {
		}
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkNested2Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
			}
		}
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
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

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
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

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}
