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
