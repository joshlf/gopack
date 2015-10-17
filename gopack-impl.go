// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	"reflect"
	"strconv"
)

type packer func(b []byte, v reflect.Value)
type unpacker func(b []byte, v reflect.Value)

func makePackerWrapper(strct reflect.Type) packer {
	p, bits, err := makePacker(0, strct)
	if err != nil {
		return func(b []byte, v reflect.Value) {
			panic(err)
		}
	}
	bytes := int(bits) / 8
	if bits%8 != 0 {
		bytes++
	}
	return func(b []byte, v reflect.Value) {
		if len(b) < bytes {
			panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
		}
		for i := 0; i < bytes; i++ {
			b[i] = 0
		}
		p(b, v)
	}
}

func makeUnpackerWrapper(strct reflect.Type) unpacker {
	u, bits, err := makeUnpacker(0, strct)
	if err != nil {
		return func(b []byte, v reflect.Value) {
			panic(err)
		}
	}
	// Check for non-pointers after
	// checking for errors so that
	// passing a non-pointer value
	// with an invalid type panics
	// (as opposed to being a no-op)
	if strct.Kind() != reflect.Ptr {
		return noOpUnpacker
	}
	bytes := int(bits) / 8
	if bits%8 != 0 {
		bytes++
	}
	return func(b []byte, v reflect.Value) {
		if len(b) < bytes {
			panic(Error{fmt.Errorf("gopack: buffer too small (%v; need %v)", len(b), bytes)})
		}
		u(b, v)
	}
}

// Returns the number of bits packed
// as the second return value
func makePacker(lsb uint64, strct reflect.Type) (packer, uint64, error) {
	ptrType := strct.Kind() == reflect.Ptr
	if ptrType {
		strct = strct.Elem()
	}
	if strct.Kind() != reflect.Struct {
		return nil, 0, Error{fmt.Errorf("gopack: non-struct type %v", strct.String())}
	}
	n := strct.NumField()
	packers := make([]packer, 0)
	var bitsPacked uint64
	for i := 0; i < n; i++ {
		field := strct.Field(i)
		if isExported(field) {
			f, bits, err := makeFieldPacker(lsb, field)
			if err != nil {
				return nil, 0, err
			}
			lsb += bits
			bitsPacked += bits
			packers = append(packers, f)
		} else {
			packers = append(packers, noOpPacker)
		}
	}
	return makeCallAllPackers(packers, ptrType), bitsPacked, nil
}

// Returns the number of bits packed
// as the second return value
func makeFieldPacker(lsb uint64, field reflect.StructField) (packer, uint64, error) {
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return nil, 0, err
		}
		return makeSignedSinglePacker(field.Type, lsb, uint8(bits)), bits, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return nil, 0, err
		}
		return makeUnsignedSinglePacker(field.Type, lsb, uint8(bits)), bits, nil
	case reflect.Bool:
		return makeBoolSinglePacker(lsb), 1, nil
	case reflect.Struct:
		return makePacker(lsb, field.Type)
	default:
		return nil, 0, Error{fmt.Errorf("gopack: non-packable type %v", field.Type.String())}
	}
}

func makeCallAllPackers(p []packer, ptrType bool) packer {
	if ptrType {
		return func(b []byte, v reflect.Value) {
			v = v.Elem()
			for i, f := range p {
				f(b, v.Field(i))
			}
		}
	} else {
		return func(b []byte, v reflect.Value) {
			for i, f := range p {
				f(b, v.Field(i))
			}
		}
	}
}

func noOpPacker(b []byte, v reflect.Value) {}

// Returns the number of bits unpacked
// as the second return value
func makeUnpacker(lsb uint64, strct reflect.Type) (unpacker, uint64, error) {
	ptrType := strct.Kind() == reflect.Ptr
	if ptrType {
		strct = strct.Elem()
	}
	if strct.Kind() != reflect.Struct {
		return nil, 0, Error{fmt.Errorf("gopack: non-struct type %v", strct.String())}
	}
	n := strct.NumField()
	unpackers := make([]unpacker, 0)
	var bitsUnpacked uint64
	for i := 0; i < n; i++ {
		field := strct.Field(i)
		if isExported(field) {
			f, bits, err := makeFieldUnpacker(lsb, field)
			if err != nil {
				return nil, 0, err
			}
			lsb += bits
			bitsUnpacked += bits
			unpackers = append(unpackers, f)
		} else {
			unpackers = append(unpackers, noOpUnpacker)
		}
	}
	return makeCallAllUnpackers(unpackers, ptrType), bitsUnpacked, nil
}

// Returns the number of bits unpacked
// as the second return value
func makeFieldUnpacker(lsb uint64, field reflect.StructField) (unpacker, uint64, error) {
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return nil, 0, err
		}
		return makeSignedSingleUnpacker(field.Type, lsb, uint8(bits)), bits, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits, err := getFieldWidth(field)
		if err != nil {
			return nil, 0, err
		}
		return makeUnsignedSingleUnpacker(field.Type, lsb, uint8(bits)), bits, nil
	case reflect.Bool:
		return makeBoolSingleUnpacker(lsb), 1, nil
	case reflect.Struct:
		return makeUnpacker(lsb, field.Type)
	default:
		return nil, 0, Error{fmt.Errorf("gopack: non-packable type %v", field.Type.String())}
	}
}

func makeCallAllUnpackers(u []unpacker, ptrType bool) unpacker {
	if ptrType {
		return func(b []byte, v reflect.Value) {
			v = v.Elem()
			for i, f := range u {
				f(b, v.Field(i))
			}
		}
	} else {
		return func(b []byte, v reflect.Value) {
			for i, f := range u {
				f(b, v.Field(i))
			}
		}
	}
}

func noOpUnpacker(b []byte, v reflect.Value) {}

// Only call on uint and int types
func getFieldWidth(field reflect.StructField) (uint64, error) {
	bits := uint64(field.Type.Bits())
	str := field.Tag.Get("gopack")
	if str == "" {
		return bits, nil
	}

	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, Error{fmt.Errorf("gopack: struct tag on field %q: %s",
			field.Name, err)}
	} else if n > int(bits) {
		return 0, Error{fmt.Errorf("gopack: struct tag on field %q (type %s) too wide (%d)",
			field.Name, field.Type, n)}
	} else if n < 1 {
		return 0, Error{fmt.Errorf("gopack: struct tag on field %q too small (%d)",
			field.Name, n)}
	}
	return uint64(n), nil
}

func isExported(field reflect.StructField) bool {
	// See http://golang.org/pkg/reflect/#StructField
	return field.PkgPath == ""
}
