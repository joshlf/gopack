// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gopack provides utilities for performing bit packing.
package gopack

import (
	"fmt"
	"reflect"
	"sync"
)

// Any panic originating from this package
// will be of type Error.
type Error struct {
	error
}

type cachedPacker struct {
	packer
	bytes int
}

var packerCache struct {
	sync.RWMutex
	m map[reflect.Type]cachedPacker
}

var unpackerCache struct {
	sync.RWMutex
	m map[reflect.Type]unpacker
}

// Pack the fields of strct into b. Fields must be
// of an int or bool type, or must be of a struct
// type whose fields are properly typed (structs
// may be nested arbitrarily deep). If strct is not
// a struct or a pointer to a struct, or if any
// of the fields are not of an allowed type, Pack
// will panic.
//
// Fields may optionally be tagged with a size in
// bits using the key "gopack".
//
//	type unixMode struct {
//		User, Group, Other uint8 `gopack:"3"`
//	}
//
// Fields without a tag are taken to be their native
// sizes (for example, 8 bits for uint8, 16 for int16,
// etc). If a field is tagged with an impossible size
// (less than 1, or larger than the native size of the
// field), or if the tag value cannot be parsed as an
// integer, Pack will panic. If a field holds a value
// which cannot be packed in the specified number of
// bits, Pack will panic (for example, if unixMode.User
// from the above example were set to 8, which cannot
// be stored in 3 bits).
//
// bool-typed fields always take up 1 bit, and any field
// tags are ignored.
//
// If there are bits in the last used byte of b which
// are beyond the end of the packed data (for example,
// the last four bits of the second byte when packing
// 12 bits), those bits will be zeroed. If b is not
// sufficiently long to hold all of the bits of strct,
// Pack will panic.
//
// Fields which are not exported are ignored, and
// take up no space in the packing.
//
//	// Only requires 2 bytes to store
//	type person struct {
//		name string
//		Age, Height uint8
//	}
func Pack(b []byte, strct interface{}) {
	v := reflect.ValueOf(strct)
	p, bytes := packerFor(v)
	if len(b) < bytes {
		panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
	}
	for i := 0; i < bytes; i++ {
		b[i] = 0
	}
	p(b, v)
}

// PackedSizeof returns the number of bytes needed to pack the given value.
func PackedSizeof(strct interface{}) int {
	_, bytes := packerFor(reflect.ValueOf(strct))
	return bytes
}

// Returns the packer and number of bytes needed by this packer.
func packerFor(v reflect.Value) (packer, int) {
	typ := v.Type()
	packerCache.RLock()
	entry, ok := packerCache.m[typ]
	packerCache.RUnlock()
	if ok {
		return entry.packer, entry.bytes
	}

	p, bytes := makePackerWrapper(typ)
	packerCache.Lock()
	packerCache.m[typ] = cachedPacker{packer: p, bytes: bytes}
	packerCache.Unlock()
	return p, bytes
}

// Unpack the data in b into the fields of strct.
// strct must be either a struct or a pointer to
// a struct, or else Unpack will panic. However,
// if strct is not a pointer, all values extracted
// from b will be discarded since strct is passed by
// value.
//
// All of the restrictions on the type of strct
// documented for Pack apply to Unpack.
//
// If b is not sufficiently long to hold all of
// the bits of strct, Unpack will panic.
func Unpack(b []byte, strct interface{}) {
	v := reflect.ValueOf(strct)
	typ := v.Type()
	unpackerCache.RLock()
	u, ok := unpackerCache.m[typ]
	unpackerCache.RUnlock()
	if ok {
		u(b, v)
		return
	}

	u = makeUnpackerWrapper(typ)
	unpackerCache.Lock()
	unpackerCache.m[typ] = u
	unpackerCache.Unlock()
	u(b, v)
}

func init() {
	packerCache.m = make(map[reflect.Type]cachedPacker)
	unpackerCache.m = make(map[reflect.Type]unpacker)
}
