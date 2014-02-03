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

type packerHandler struct {
	packer
	len int
}

type unpackerHandler struct {
	unpacker
	len int
}

var packerCache struct {
	sync.RWMutex
	m map[reflect.Type]packerHandler
}

var unpackerCache struct {
	sync.RWMutex
	m map[reflect.Type]unpackerHandler
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
// If b is not sufficiently long to hold all of the
// bits of strct, Pack will panic.
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
	typ := v.Type()
	packerCache.RLock()
	handler, ok := packerCache.m[typ]
	packerCache.RUnlock()
	if ok {
		if len(b) < handler.len {
			panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), handler.len)})
		}
		handler.packer(b, v)
		return
	}

	var l uint64
	handler.packer, l = makePacker(0, typ)
	handler.len = int(l) / 8
	if l%8 != 0 {
		handler.len++
	}
	packerCache.Lock()
	packerCache.m[typ] = handler
	packerCache.Unlock()
	if len(b) < handler.len {
		panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), handler.len)})
	}
	handler.packer(b, v)
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
	handler, ok := unpackerCache.m[typ]
	unpackerCache.RUnlock()
	if ok {
		if len(b) < handler.len {
			panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), handler.len)})
		}
		handler.unpacker(b, v)
		return
	}

	var l uint64
	handler.unpacker, l = makeUnpacker(0, typ)
	handler.len = int(l) / 8
	if l%8 != 0 {
		handler.len++
	}
	unpackerCache.Lock()
	unpackerCache.m[typ] = handler
	unpackerCache.Unlock()
	if len(b) < handler.len {
		panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), handler.len)})
	}
	handler.unpacker(b, v)
}

func init() {
	packerCache.m = make(map[reflect.Type]packerHandler)
	unpackerCache.m = make(map[reflect.Type]unpackerHandler)
}
