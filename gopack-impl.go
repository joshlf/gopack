// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

type packer func(b []byte, v reflect.Value)
type unpacker func(b []byte, v reflect.Value)

func makePackerWrapper(strct reflect.Type) packer {
	p, bits := makePacker(0, strct)
	bytes := int(bits) / 8
	if bits%8 != 0 {
		bytes++
	}
	if bytes < 8 {
		// We need at least a full 64-bit word,
		// so if we aren't going to get that from
		// the user, we have to stack-allocate
		// and then copy over
		return func(b []byte, v reflect.Value) {
			if len(b) < bytes {
				panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
			}
			var buf [8]byte
			p(buf[:], v)
			copy(b, buf[:])
		}
	}
	return func(b []byte, v reflect.Value) {
		if len(b) < bytes {
			panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
		}
		p(b, v)
	}
}

func makeUnpackerWrapper(strct reflect.Type) unpacker {
	u, bits := makeUnpacker(0, strct)
	bytes := int(bits) / 8
	if bits%8 != 0 {
		bytes++
	}
	if bytes < 8 {
		// We need at least a full 64-bit word,
		// so if we aren't going to get that from
		// the user, we have to stack-allocate
		// and then copy over
		return func(b []byte, v reflect.Value) {
			if len(b) < bytes {
				panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
			}
			var buf [8]byte
			copy(buf[:], b)
			u(buf[:], v)
		}
	}
	return func(b []byte, v reflect.Value) {
		if len(b) < bytes {
			panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
		}
		u(b, v)
	}
}

// Returns the number of bits packed
func makePacker(lsb uint64, strct reflect.Type) (packer, uint64) {
	ptrType := strct.Kind() == reflect.Ptr
	if ptrType {
		strct = strct.Elem()
	}
	if strct.Kind() != reflect.Struct {
		return makePanicPacker(Error{fmt.Errorf("gopack: non-struct type %v", strct.String())}), 0
	}
	n := strct.NumField()
	packers := make([]packer, 0)
	var bitsPacked uint64
	for i := 0; i < n; i++ {
		field := strct.Field(i)
		if isExported(field) {
			f, bits := makeFieldPacker(lsb, field)
			lsb += bits
			bitsPacked += bits
			packers = append(packers, f)
		} else {
			packers = append(packers, noOpPacker)
		}
	}
	if ptrType {
		return makeDereferencePacker(makeCallAllPackers(packers)), bitsPacked
	}
	return makeCallAllPackers(packers), bitsPacked
}

// Returns the number of bits packed
func makeFieldPacker(lsb uint64, field reflect.StructField) (packer, uint64) {
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return makePanicPacker(err), 0
		}
		return makeSignedSinglePacker(uint8(lsb), uint8(bits)), bits
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return makePanicPacker(err), 0
		}
		return makeUnsignedSinglePacker(uint8(lsb), uint8(bits)), bits
	case reflect.Bool:
		return makeBoolSinglePacker(uint8(lsb)), 1
	case reflect.Struct:
		return makePacker(lsb, field.Type)
	default:
		return makePanicPacker(Error{fmt.Errorf("gopack: non-packable type %v", field.Type.String())}), 0
	}
}

func makeCallAllPackers(p []packer) packer {
	return func(b []byte, v reflect.Value) {
		for i, f := range p {
			f(b, v.Field(i))
		}
	}
}

func makeDereferencePacker(p packer) packer {
	return func(b []byte, v reflect.Value) {
		p(b, v.Elem())
	}
}

func noOpPacker(b []byte, v reflect.Value) {}

// Returns the number of bits packed
func makeUnpacker(lsb uint64, strct reflect.Type) (unpacker, uint64) {
	ptrType := strct.Kind() == reflect.Ptr
	if ptrType {
		strct = strct.Elem()
	}
	if strct.Kind() != reflect.Struct {
		return makePanicUnpacker(Error{fmt.Errorf("gopack: non-struct type %v", strct.String())}), 0
	}
	n := strct.NumField()
	unpackers := make([]unpacker, 0)
	var bitsUnpacked uint64
	for i := 0; i < n; i++ {
		field := strct.Field(i)
		if isExported(field) {
			f, bits := makeFieldUnpacker(lsb, field)
			lsb += bits
			bitsUnpacked += bits
			unpackers = append(unpackers, f)
		} else {
			unpackers = append(unpackers, noOpUnpacker)
		}
	}
	if ptrType {
		return makeDereferenceUnpacker(makeCallAllUnpackers(unpackers)), bitsUnpacked
	}
	return makeCallAllUnpackers(unpackers), bitsUnpacked
}

// Returns the number of bits unpacked
func makeFieldUnpacker(lsb uint64, field reflect.StructField) (unpacker, uint64) {
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return makePanicUnpacker(err), 0
		}
		return makeSignedSingleUnpacker(uint8(lsb), uint8(bits)), bits
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return makePanicUnpacker(err), 0
		}
		return makeUnsignedSingleUnpacker(uint8(lsb), uint8(bits)), bits
	case reflect.Bool:
		return makeBoolSingleUnpacker(uint8(lsb)), 1
	case reflect.Struct:
		return makeUnpacker(lsb, field.Type)
	default:
		return makePanicUnpacker(Error{fmt.Errorf("gopack: non-packable type %v", field.Type.String())}), 0
	}
}

func makeCallAllUnpackers(u []unpacker) unpacker {
	return func(b []byte, v reflect.Value) {
		for i, f := range u {
			f(b, v.Field(i))
		}
	}
}

func makeDereferenceUnpacker(u unpacker) unpacker {
	return func(b []byte, v reflect.Value) {
		u(b, v.Elem())
	}
}

func noOpUnpacker(b []byte, v reflect.Value) {}

/*
	Single packers and unpackers
	pack and unpack a single field
*/

func makeUnsignedSinglePacker(lsb, width uint8) packer {
	firstByte := lsb / 8
	lsb = lsb % 8
	if width < 64 {
		maxVal := (uint64(1) << width) - 1
		return func(b []byte, field reflect.Value) {
			u := field.Uint()
			if u > maxVal {
				panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
			}
			packUint64(b[firstByte:], u, lsb, width)
		}
	}
	return func(b []byte, field reflect.Value) {
		packUint64(b[firstByte:], field.Uint(), lsb, width)
	}
}

func makeSignedSinglePacker(lsb, width uint8) packer {
	firstByte := lsb / 8
	lsb = lsb % 8
	if width < 64 {
		minVal := int64(-1) << (width - 1)
		var maxVal int64
		// Place maxUval in a separate scope
		// to prevent the returned closure from
		// unnecessarily closing over it.
		{
			maxUval := uint64(math.MaxUint64) >> (65 - width)
			maxVal = *(*int64)(unsafe.Pointer(&maxUval))
		}
		return func(b []byte, field reflect.Value) {
			i := field.Int()
			if i < minVal || i > maxVal {
				panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, i)})
			}
			packInt64(b[firstByte:], i, lsb, width)
		}
	}
	return func(b []byte, field reflect.Value) {
		packInt64(b[firstByte:], field.Int(), lsb, width)
	}
}

func makeBoolSinglePacker(lsb uint8) packer {
	firstByte := lsb / 8
	lsb = lsb % 8
	return func(b []byte, field reflect.Value) {
		var val uint64
		if field.Bool() {
			val = 1
		}
		packUint64(b[firstByte:], val, lsb, 1)
	}
}

func makeUnsignedSingleUnpacker(lsb, width uint8) unpacker {
	firstByte := lsb / 8
	lsb = lsb % 8
	return func(b []byte, field reflect.Value) {
		field.SetUint(unpackUint64(b[firstByte:], lsb, width))
	}
}

func makeSignedSingleUnpacker(lsb, width uint8) unpacker {
	firstByte := lsb / 8
	lsb = lsb % 8
	return func(b []byte, field reflect.Value) {
		field.SetInt(unpackInt64(b[firstByte:], lsb, width))
	}
}

func makeBoolSingleUnpacker(lsb uint8) unpacker {
	firstByte := lsb / 8
	lsb = lsb % 8
	return func(b []byte, field reflect.Value) {
		field.SetBool(unpackUint64(b[firstByte:], lsb, 1) == 1)
	}
}

func makePanicPacker(err error) packer {
	return func(b []byte, v reflect.Value) {
		panic(err)
	}
}

func makePanicUnpacker(err error) unpacker {
	return func(b []byte, v reflect.Value) {
		panic(err)
	}
}

// Only call on uint and int types
func getFieldWidth(field reflect.StructField) (uint64, error) {
	bits := uint64(field.Type.Bits())
	str := field.Tag.Get("gopack")
	if str == "" {
		return bits, nil
	}

	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, Error{fmt.Errorf("gopack: struct tag on field \"%v\": %v", field.Name, err)}
	}
	if n > int(bits) {
		return 0, Error{fmt.Errorf("gopack: struct tag on field \"%v\" (type %v) too wide (%v)", field.Name, field.Type, n)}
	}
	if n < 1 {
		return 0, Error{fmt.Errorf("gopack: struct tag on field \"%v\" too small (%v)", field.Name, n)}
	}
	return uint64(n), nil
}

func isExported(field reflect.StructField) bool {
	// See http://golang.org/pkg/reflect/#StructField
	return field.PkgPath == ""
}
