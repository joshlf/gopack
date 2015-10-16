// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
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

func TestLargeValue(t *testing.T) {
	type big struct {
		V1, V2, V3, V4, V5, V6, V7, V8 uint32 // 32 bytes

		B uint8
	}

	b := big{B: 42}
	bytes := make([]byte, 33)
	Pack(bytes, &b)
	if bytes[0] != 0 {
		t.Errorf("First byte should have been 0 but was %v", bytes[0])
	}

	b2 := big{}
	Unpack(bytes, &b2)
	if b != b2 {
		t.Fatalf("Expected %#v; got %#v", b, b2)
	}
}
