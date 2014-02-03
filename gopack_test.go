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

func TestPackUnpack(t *testing.T) {
	// Do twice to test it once
	// the type is already in cache
	testPackUnpack(t)
	testPackUnpack(t)
}

func testPackUnpack(t *testing.T) {
	type typ struct {
		F1 uint8
	}

	t1 := typ{255}
	bytes := []byte{0}
	Pack(bytes, t1)

	t1 = typ{}
	Unpack(bytes, &t1)
	if t1 != (typ{255}) {
		t.Fatalf("Expected %v; got %v", typ{255}, t1)
	}
}

func BenchmarkLockContention(b *testing.B) {
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

	runtime.GOMAXPROCS(2)
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
	runtime.GOMAXPROCS(1)
}

func BenchmarkNoLockContention(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := typ{}
	bytes := []byte{0}
	Pack(bytes, t)
	b.ResetTimer()
	for i := 0; i < b.N; i += 2 {
		Pack(bytes, t)
	}
}

func BenchmarkPack1Field(b *testing.B) {
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
	p, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
	}
}

func BenchmarkPack2Fields(b *testing.B) {
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
	p, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
	}
}

func BenchmarkPack4Fields(b *testing.B) {
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
	p, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
	}
}

func BenchmarkPack8Fields(b *testing.B) {
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
	p, _ := makePacker(0, reflect.TypeOf(t))
	v := reflect.ValueOf(t)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, v)
	}
}
