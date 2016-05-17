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

type layout []field

type field struct {
	// offset from address of the top-level value
	valByteOffset uint64

	// offset from the beginning of the byte slice
	bufByteOffset uint64

	// LSB within the first target byte in the byte slice
	bufLSB uint8

	// width of the packed field
	bits uint8

	// how many bytes of the packed slice this field
	// covers (e.g., an 8-bit value starting at LSB
	// 4 would cover two bytes)
	bytesSpanned uint8

	// whether values in this field can overflow
	// their packed representations (i.e., whether
	// the packed width is smaller than the field
	// type's width)
	canOverflow bool

	kind reflect.Kind

	signed bool
}

func makeLayout(v reflect.Value) (l layout, bytes int, err error) {
	if !v.CanAddr() {
		v = reflect.New(v.Type()).Elem()
	}
	bits, err := processType(&l, v, 0, v.UnsafeAddr(), 0, "")
	if err != nil {
		return nil, 0, err
	}
	bytestmp := bits / 8
	if bits%8 != 0 {
		bytestmp++
	}
	bytestmpInt := int(bytestmp)
	if bytestmpInt < 0 {
		bytestmpInt = 0
	}
	if uint64(bytestmpInt) != bytestmp {
		// check for int in particular because
		// slice indices are ints, so this will
		// cause our algorithm to fail
		return nil, 0, fmt.Errorf("packed byte length overflows int: %v", bytestmp)
	}
	return l, int(bytestmp), nil
}

// processType adds entries to l for v or v's fields
// or elements. It returns the number of bits that
// were used in the packed buffer, and any error
// encountered.
//
// bufLsb is the absolute LSB from the beginning of the
// packing buffer. valByteOffset is the offset of this
// field from the address of the top-lvel value.
// fieldBits is only set if a custom bit width was set
// for this field. fieldName is only set if this is
// a struct field.
func processType(l *layout, v reflect.Value, bufLsb uint64, baseAddr uintptr, fieldBits uint64, fieldName string) (bits uint64, err error) {
	t := v.Type()

	if fieldBits > 0 {
		// fieldBits > 0 means that fieldBits was intentionally set,
		// which means that it has already been verified that it's
		// valid to have a struct field tag on this field.

		if int(fieldBits) > t.Bits() {
			return 0, fmt.Errorf("struct tag on field %q too big", fieldName)
		}
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var signed bool
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			signed = true
		}
		bits = uint64(t.Bits())
		if fieldBits > 0 {
			bits = fieldBits
		}
		bytesSpanned := ((bufLsb % 8) + bits) / 8
		if ((bufLsb%8)+bits)%8 != 0 {
			bytesSpanned++
		}
		f := field{
			valByteOffset: uint64(v.UnsafeAddr() - baseAddr),
			bufByteOffset: bufLsb / 8,
			bufLSB:        uint8(bufLsb % 8),
			bits:          uint8(bits),
			bytesSpanned:  uint8(bytesSpanned),
			canOverflow:   bits != uint64(t.Bits()),
			kind:          t.Kind(),
			signed:        signed,
		}
		*l = append(*l, f)
		return bits, nil
	case reflect.Bool:
		f := field{
			valByteOffset: uint64(v.UnsafeAddr() - baseAddr),
			bufByteOffset: bufLsb / 8,
			bufLSB:        uint8(bufLsb % 8),
			bits:          1,
			bytesSpanned:  1,
			canOverflow:   false,
			kind:          reflect.Bool,
			signed:        false,
		}
		*l = append(*l, f)
		return 1, nil
	case reflect.Array:
		bits = 0
		for i := 0; i < v.Len(); i++ {
			vv := v.Index(i)
			// Pass most arguments through from the top level. baseAddr is constant,
			// and fieldBits and fieldName, if they were set by the parent, should
			// be interpreted by the child.
			b, err := processType(l, vv, bufLsb+bits, baseAddr, fieldBits, fieldName)
			if err != nil {
				return 0, err
			}
			bits += b
		}
		return bits, nil
	case reflect.Struct:
		bits = 0
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			vv := v.Field(i)

			if field.PkgPath != "" {
				// It's not exported; see
				// http://golang.org/pkg/reflect/#StructField
				continue
			}

			fieldBits = 0
			str := field.Tag.Get("gopack")

			if str != "" {
				k := vv.Kind()
				if k == reflect.Array {
					k = vv.Elem().Type().Kind()
				}
				switch k {
				case reflect.Uint, reflect.Int8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Int, reflect.Uint8, reflect.Int16, reflect.Int32, reflect.Int64:
				default:
					return 0, fmt.Errorf("struct tag not allowed on field %q", field.Name)
				}

				n, err := strconv.Atoi(str)
				if err != nil {
					return 0, fmt.Errorf("struct tag on field %q: %s", field.Name, err)
				}
				if n < 1 {
					return 0, fmt.Errorf("struct tag on field %q too small", field.Name)
				}

				fieldBits = uint64(n)
			}

			b, err := processType(l, vv, bufLsb+bits, baseAddr, fieldBits, field.Name)
			if err != nil {
				return 0, err
			}
			bits += b
		}
		return bits, nil
	default:
		return 0, fmt.Errorf("cannot pack type %v", t)
	}
}

// TODO(joshlf): Add name to field type so we can give
// more helpful error messages

// pack packs the value in v into buf. It assumes that v
// is addressable, and is described by l, and that buf
// contains entirely zeroes.
func pack(buf []byte, l layout, v reflect.Value) error {
	// TODO(joshlf): Is it safe to make any function
	// calls after getting the address like this?
	// Could there be an issue if a function call
	// overflows the stack, causes a new stack to
	// be allocated, but addr isn't recognized as
	// a pointer, so it's not adjusted? This would
	// only happen if v were stack-allocated, which
	// would at the very least require the escape
	// analysis to detect that it won't escape, which
	// seems exceedingly unlikely, but maybe...?
	addr := uintptr(v.UnsafeAddr())
	for _, field := range l {
		if field.kind == reflect.Bool {
			if *(*bool)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) {
				buf[field.bufByteOffset] |= byte(1) << uint8(field.bufLSB%8)
			}
		} else if field.signed {
			var val int64
			switch field.kind {
			case reflect.Int8:
				val = int64(*(*int8)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			case reflect.Int16:
				val = int64(*(*int16)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			case reflect.Int32:
				val = int64(*(*int32)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			case reflect.Int64:
				val = int64(*(*int64)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			}
			minVal := int64(-1) << (field.bits - 1)
			maxUval := uint64(math.MaxUint64) >> (65 - field.bits)
			maxVal := *(*int64)(unsafe.Pointer(&maxUval))
			switch field.bytesSpanned {
			case 1:
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {

						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					buf[field.bufByteOffset] |= byte(u << field.bufLSB)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					buf[field.bufByteOffset] |= byte(u << field.bufLSB)
				}
			case 2:
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(u << field.bufLSB)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(u << field.bufLSB)
				}
			case 3:
				shift := 16 - field.bufLSB
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(u << field.bufLSB)
					buf[field.bufByteOffset+2] |= byte(u >> shift)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(u << field.bufLSB)
					buf[field.bufByteOffset+2] |= byte(u >> shift)
				}
			case 4:
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
				}
			case 5:
				shift := 32 - field.bufLSB
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
					buf[field.bufByteOffset+4] |= byte(u >> shift)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
					buf[field.bufByteOffset+4] |= byte(u >> shift)
				}
			case 6:
				shift := 32 - field.bufLSB
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(u >> shift)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(u >> shift)
				}
			case 7:
				shift1 := 32 - field.bufLSB
				shift2 := 48 - field.bufLSB
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(u >> shift1)
					buf[field.bufByteOffset+6] |= byte(u >> shift2)
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(u << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(u >> shift1)
					buf[field.bufByteOffset+6] |= byte(u >> shift2)
				}
			case 8:
				if field.canOverflow {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= u << field.bufLSB
				} else {
					u := (uint64(val) << (64 - field.bits)) >> (64 - field.bits)
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= u << field.bufLSB
				}
			case 9:
				// 64 - ((bits + lsb) - 64)
				shift1 := 64 - field.bits
				shift2 := 64 - field.bufLSB
				if field.canOverflow {
					if val < minVal || val > maxVal {
						return fmt.Errorf("value out of range: max %v, min %v; got %v", maxVal, minVal, val)
					}
					u := (uint64(val) << shift1) >> shift1
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= u << field.bufLSB
					buf[field.bufByteOffset+8] |= byte(u >> shift2)
				} else {
					u := (uint64(val) << shift1) >> shift1
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= u << field.bufLSB
					buf[field.bufByteOffset+8] |= byte(u >> shift2)
				}
			default:
				panic("unreachable")
			}
		} else {
			var val uint64
			switch field.kind {
			case reflect.Uint8:
				val = uint64(*(*uint8)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			case reflect.Uint16:
				val = uint64(*(*uint16)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			case reflect.Uint32:
				val = uint64(*(*uint32)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			case reflect.Uint64:
				val = uint64(*(*uint64)(unsafe.Pointer(addr + uintptr(field.valByteOffset))))
			}
			maxVal := (uint64(1) << field.bits) - 1
			switch field.bytesSpanned {
			case 1:
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					buf[field.bufByteOffset] |= byte(val << field.bufLSB)
				} else {
					buf[field.bufByteOffset] |= byte(val << field.bufLSB)
				}
			case 2:
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(val << field.bufLSB)
				} else {
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(val << field.bufLSB)
				}
			case 3:
				shift := 16 - field.bufLSB
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(val << field.bufLSB)
					buf[field.bufByteOffset+2] |= byte(val >> shift)
				} else {
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint16(val << field.bufLSB)
					buf[field.bufByteOffset+2] |= byte(val >> shift)
				}
			case 4:
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
				} else {
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
				}
			case 5:
				shift := 32 - field.bufLSB
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
					buf[field.bufByteOffset+4] |= byte(val >> shift)
				} else {
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
					buf[field.bufByteOffset+4] |= byte(val >> shift)
				}
			case 6:
				shift := 32 - field.bufLSB
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(val >> shift)
				} else {
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(val >> shift)
				}
			case 7:
				shift1 := 32 - field.bufLSB
				shift2 := 48 - field.bufLSB
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(val >> shift1)
					buf[field.bufByteOffset+6] |= byte(val >> shift2)
				} else {
					*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) |= uint32(val << field.bufLSB)
					*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4])) |= uint16(val >> shift1)
					buf[field.bufByteOffset+6] |= byte(val >> shift2)
				}
			case 8:
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= val << field.bufLSB
				} else {
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= val << field.bufLSB
				}
			case 9:
				shift := 64 - field.bufLSB
				if field.canOverflow {
					if val > maxVal {
						return fmt.Errorf("value out of range: max %v; got %v", maxVal, val)
					}
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= val << field.bufLSB
					buf[field.bufByteOffset+8] = byte(val >> shift)
				} else {
					*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) |= val << field.bufLSB
					buf[field.bufByteOffset+8] = byte(val >> shift)
				}
			default:
				panic("unreachable")
			}
		}
	}

	return nil
}

// unpack unpacks the value in buf into v. It assumes that v
// is addressable and is described by l.
func unpack(buf []byte, l layout, v reflect.Value) {
	// TODO(joshlf): Is it safe to make any function
	// calls after getting the address like this?
	// Could there be an issue if a function call
	// overflows the stack, causes a new stack to
	// be allocated, but addr isn't recognized as
	// a pointer, so it's not adjusted? This would
	// only happen if v were stack-allocated, which
	// would at the very least require the escape
	// analysis to detect that it won't escape, which
	// seems exceedingly unlikely, but maybe...?
	addr := uintptr(v.UnsafeAddr())
	for _, field := range l {
		if field.kind == reflect.Bool {
			*(*bool)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) =
				buf[field.bufByteOffset]&(byte(1)<<uint8(field.bufLSB%8)) > 0
		} else if field.signed {
			var val int64
			switch field.bytesSpanned {
			case 1:
				shift1 := 64 - (field.bufLSB + field.bits)
				shift2 := 64 - field.bits
				val = (int64(buf[field.bufByteOffset]) << shift1) >> shift2
			case 2:
				shift1 := 64 - (field.bufLSB + field.bits)
				shift2 := 64 - field.bits
				val = (int64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset]))) << shift1) >> shift2
			case 3:
				shift1 := 64 - ((field.bufLSB + field.bits) - 16)
				shift2 := (shift1 + field.bufLSB) - 16
				i := int64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset]))) >> field.bufLSB
				val = i | ((int64(buf[field.bufByteOffset+2]) << shift1) >> shift2)
			case 4:
				shift1 := 64 - (field.bufLSB + field.bits)
				shift2 := 64 - field.bits
				val = (int64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset]))) << shift1) >> shift2
			case 5:
				shift1 := 64 - ((field.bufLSB + field.bits) - 32)
				shift2 := (shift1 + field.bufLSB) - 32
				i := int64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset]))) >> field.bufLSB
				val = i | ((int64(buf[field.bufByteOffset+4]) << shift1) >> shift2)
			case 6:
				shift1 := 64 - ((field.bufLSB + field.bits) - 32)
				shift2 := (shift1 + field.bufLSB) - 32
				i := int64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset]))) >> field.bufLSB
				val = i | ((int64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4]))) << shift1) >> shift2)
			case 7:
				shift1 := 32 - field.bufLSB
				shift2 := 64 - ((field.bufLSB + field.bits) - 48)
				shift3 := (shift2 + field.bufLSB) - 48
				i := int64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) >> field.bufLSB)
				i |= int64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4]))) >> shift1
				val = i | (int64(buf[field.bufByteOffset+6])<<shift2)>>shift3
			case 8:
				shift1 := 64 - (field.bufLSB + field.bits)
				shift2 := 64 - field.bits
				val = (*(*int64)(unsafe.Pointer(&buf[field.bufByteOffset])) << shift1) >> shift2
			case 9:
				shift1 := 128 - (field.bufLSB + field.bits)
				shift2 := (shift1 + field.bufLSB) - 64
				i := int64(*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) >> field.bufLSB)
				val = i | ((int64(buf[field.bufByteOffset+8]) << shift1) >> shift2)
			default:
				panic("unreachable")
			}
			switch field.kind {
			case reflect.Int8:
				*(*int8)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = int8(val)
			case reflect.Int16:
				*(*int16)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = int16(val)
			case reflect.Int32:
				*(*int32)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = int32(val)
			case reflect.Int64:
				*(*int64)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = int64(val)
			}
		} else {
			var val uint64
			switch field.bytesSpanned {
			case 1:
				shift1 := 8 - (field.bufLSB + field.bits)
				shift2 := 8 - field.bits
				val = uint64((buf[field.bufByteOffset] << shift1) >> shift2)
			case 2:
				shift1 := 16 - (field.bufLSB + field.bits)
				shift2 := 16 - field.bits
				val = uint64((*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset])) << shift1) >> shift2)
			case 3:
				shift1 := 64 - ((field.bufLSB + field.bits) - 16)
				shift2 := (shift1 + field.bufLSB) - 16
				u := uint64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset]))) >> field.bufLSB
				val = u | ((uint64(buf[field.bufByteOffset+2]) << shift1) >> shift2)
			case 4:
				shift1 := 32 - (field.bufLSB + field.bits)
				shift2 := 32 - field.bits
				val = uint64((*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) << shift1) >> shift2)
			case 5:
				shift1 := 64 - ((field.bufLSB + field.bits) - 32)
				shift2 := (shift1 + field.bufLSB) - 32
				u := uint64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset]))) >> field.bufLSB
				val = u | ((uint64(buf[field.bufByteOffset+4]) << shift1) >> shift2)
			case 6:
				shift1 := 64 - ((field.bufLSB + field.bits) - 32)
				shift2 := (shift1 + field.bufLSB) - 32
				u := uint64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset]))) >> field.bufLSB
				val = u | ((uint64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4]))) << shift1) >> shift2)
			case 7:
				shift1 := 32 - field.bufLSB
				shift2 := 64 - ((field.bufLSB + field.bits) - 48)
				shift3 := (shift2 + field.bufLSB) - 48
				u := uint64(*(*uint32)(unsafe.Pointer(&buf[field.bufByteOffset])) >> field.bufLSB)
				u |= uint64(*(*uint16)(unsafe.Pointer(&buf[field.bufByteOffset+4]))) >> shift1
				val = u | (uint64(buf[field.bufByteOffset+6])<<shift2)>>shift3
			case 8:
				shift1 := 64 - (field.bufLSB + field.bits)
				shift2 := 64 - field.bits
				val = (*(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) << shift1) >> shift2
			case 9:
				shift1 := 128 - (field.bufLSB + field.bits)
				shift2 := (shift1 + field.bufLSB) - 64
				u := *(*uint64)(unsafe.Pointer(&buf[field.bufByteOffset])) >> field.bufLSB
				val = u | ((uint64(buf[field.bufByteOffset+8]) << shift1) >> shift2)
			default:
				panic("unreachable")
			}
			switch field.kind {
			case reflect.Uint8:
				*(*uint8)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = uint8(val)
			case reflect.Uint16:
				*(*uint16)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = uint16(val)
			case reflect.Uint32:
				*(*uint32)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = uint32(val)
			case reflect.Uint64:
				*(*uint64)(unsafe.Pointer(addr + uintptr(field.valByteOffset))) = uint64(val)
			}
		}
	}
}
