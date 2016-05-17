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

// Any error returned by this package
// will be of type Error.
type Error struct{ error }

// since we never return an Error value as
// an error (just call panic on them directly),
// make sure it actually implements error
var _ = error(Error{nil})

type cachedLayout struct {
	layout
	bytes int
}

var layoutCache = struct {
	sync.RWMutex
	m map[reflect.Type]cachedLayout
}{m: make(map[reflect.Type]cachedLayout)}

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
func Pack(b []byte, v interface{}) (err error) {
	rv, err := normalizeArgument(v, false, "Pack")
	if err != nil {
		return err
	}

	layout, bytes, err := layoutFor(rv)
	if err != nil {
		return errorf("gopack: Pack: %v", err)
	}

	if len(b) < bytes {
		return errorf("gopack: Pack: buffer too small (got %v; need %v)", len(b), bytes)
	}

	// clear it since pack uses |= operations
	for i := range b {
		b[i] = 0
	}
	err = pack(b, layout, rv)
	if err != nil {
		return errorf("gopack: Pack: %v", err)
	}
	return nil
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
func Unpack(b []byte, v interface{}) (err error) {
	rv, err := normalizeArgument(v, true, "Unpack")
	if err != nil {
		return err
	}

	layout, bytes, err := layoutFor(rv)
	if err != nil {
		return errorf("gopack: Unpack: %v", err)
	}

	if len(b) < bytes {
		return errorf("gopack: Unpack: buffer too small (got %v; need %v)", len(b), bytes)
	}

	unpack(b, layout, rv)
	return nil
}

// PackedSizeof returns the number of bytes needed to pack the given value.
func PackedSizeof(v interface{}) (bytes int, err error) {
	rv, err := normalizeArgument(v, false, "PackedSizeof")
	if err != nil {
		return 0, err
	}
	_, bytes, err = layoutFor(rv)
	if err != nil {
		return 0, errorf("gopack: PackedSizeof: %v", err)
	}
	return bytes, nil
}

func layoutFor(v reflect.Value) (layout, int, error) {
	typ := v.Type()
	layoutCache.RLock()
	entry, ok := layoutCache.m[typ]
	layoutCache.RUnlock()
	if ok {
		return entry.layout, entry.bytes, nil
	}

	l, bytes, err := makeLayout(v)
	if err != nil {
		return nil, 0, err
	}

	layoutCache.Lock()
	layoutCache.m[typ] = cachedLayout{layout: l, bytes: bytes}
	layoutCache.Unlock()
	return l, bytes, nil
}

// normalizeArgument normalizes v according to the following rules:
//  - if unpack is true, v must be a pointer, and is dereferenced
//    until a non-pointer is encountered, and that is returned (it
//    will be addressable)
//  - if pack is true:
//    - if v is a pointer, it is dereferenced until a non-pointer
//      is encountered, and that is returned (it will be addressable)
//    - if v is a non-pointer, a new addressable value is allocated,
//      and its contents are set to those of v; that is returned
func normalizeArgument(v interface{}, unpack bool, fname string) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	if v == nil {
		return reflect.Value{}, errorf("gopack: %v(nil)", fname)
	}
	if unpack && rv.Kind() != reflect.Ptr {
		return reflect.Value{}, errorf("gopack: %v(non-pointer %v)", fname, rv.Type())
	}
	if rv.Kind() != reflect.Ptr {
		// unpack must be false, or the previous
		// condition would have been true
		newrv := reflect.New(rv.Type()).Elem()
		newrv.Set(rv)
		rv = newrv
	}
	for {
		if rv.Kind() != reflect.Ptr {
			break
		}
		if rv.IsNil() {
			// TODO(joshlf)
			return reflect.Value{}, errorf("")
		}
		rv = rv.Elem()
	}
	return rv, nil
}

func errorf(format string, a ...interface{}) error {
	return Error{fmt.Errorf(format, a...)}
}
