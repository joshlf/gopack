// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

func makeUnsignedSinglePacker(typ reflect.Type, ilsb uint64, width uint8) packer {
	firstByte := ilsb / 8
	lsb := uint8(ilsb % 8)
	canOverflow := width != uint8(typ.Bits())
	maxVal := (uint64(1) << width) - 1
	switch {
	case lsb+width <= 8:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				b[firstByte] |= byte(u << lsb)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				b[firstByte] |= byte(field.Uint() << lsb)
			}
		}
	case lsb+width <= 16:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(field.Uint() << lsb)
			}
		}
	case lsb+width <= 24:
		shift := 16 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
				b[firstByte+2] |= byte(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
				b[firstByte+2] |= byte(u >> shift)
			}
		}
	case lsb+width <= 32:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(field.Uint() << lsb)
			}
		}
	case lsb+width <= 40:
		shift := 32 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				b[firstByte+4] |= byte(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				b[firstByte+4] |= byte(u >> shift)
			}
		}
	case lsb+width <= 48:
		shift := 32 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift)
			}
		}
	case lsb+width <= 56:
		shift1 := 32 - lsb
		shift2 := 48 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift1)
				b[firstByte+6] |= byte(u >> shift2)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift1)
				b[firstByte+6] |= byte(u >> shift2)
			}
		}
	case lsb+width <= 64:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
			}
		} else {
			return func(b []byte, field reflect.Value) {
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= field.Uint() << lsb
			}
		}
	default:
		// Assume lsb+width <= 72
		shift := 64 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				if u > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v; got %v", maxVal, u)})
				}
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
				b[firstByte+8] = byte(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				u := field.Uint()
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
				b[firstByte+8] = byte(u >> shift)
			}
		}
	}
}

func makeSignedSinglePacker(typ reflect.Type, ilsb uint64, width uint8) packer {
	firstByte := ilsb / 8
	lsb := uint8(ilsb % 8)
	canOverflow := width != uint8(typ.Bits())
	minVal := int64(-1) << (width - 1)
	maxUval := uint64(math.MaxUint64) >> (65 - width)
	maxVal := *(*int64)(unsafe.Pointer(&maxUval))
	switch {
	case lsb+width <= 8:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				b[firstByte] |= byte(u << lsb)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				b[firstByte] |= byte(u << lsb)
			}
		}
	case lsb+width <= 16:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
			}
		}
	case lsb+width <= 24:
		shift := 16 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
				b[firstByte+2] |= byte(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint16)(unsafe.Pointer(&b[firstByte])) |= uint16(u << lsb)
				b[firstByte+2] |= byte(u >> shift)
			}
		}
	case lsb+width <= 32:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
			}
		}
	case lsb+width <= 40:
		shift := 32 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				b[firstByte+4] |= byte(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				b[firstByte+4] |= byte(u >> shift)
			}
		}
	case lsb+width <= 48:
		shift := 32 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift)
			}
		}
	case lsb+width <= 56:
		shift1 := 32 - lsb
		shift2 := 48 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift1)
				b[firstByte+6] |= byte(u >> shift2)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint32)(unsafe.Pointer(&b[firstByte])) |= uint32(u << lsb)
				*(*uint16)(unsafe.Pointer(&b[firstByte+4])) |= uint16(u >> shift1)
				b[firstByte+6] |= byte(u >> shift2)
			}
		}
	case lsb+width <= 64:
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << (64 - width)) >> (64 - width)
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
			}
		}
	default:
		// Assume lsb+width <= 72
		// 64 - ((width + lsb) - 64)
		shift1 := 64 - width
		shift2 := 64 - lsb
		if canOverflow {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				if val < minVal || val > maxVal {
					panic(Error{fmt.Errorf("gopack: value out of range: max %v, min %v; got %v", maxVal, minVal, val)})
				}
				u := (uint64(val) << shift1) >> shift1
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
				b[firstByte+8] |= byte(u >> shift2)
			}
		} else {
			return func(b []byte, field reflect.Value) {
				val := field.Int()
				u := (uint64(val) << shift1) >> shift1
				*(*uint64)(unsafe.Pointer(&b[firstByte])) |= u << lsb
				b[firstByte+8] |= byte(u >> shift2)
			}
		}
	}
}

func makeUnsignedSingleUnpacker(typ reflect.Type, ilsb uint64, width uint8) unpacker {
	firstByte := ilsb / 8
	lsb := uint8(ilsb % 8)
	switch {
	case lsb+width <= 8:
		shift1 := 8 - (lsb + width)
		shift2 := 8 - width
		return func(b []byte, field reflect.Value) {
			field.SetUint(uint64((b[firstByte] << shift1) >> shift2))
		}
	case lsb+width <= 16:
		shift1 := 16 - (lsb + width)
		shift2 := 16 - width
		return func(b []byte, field reflect.Value) {
			field.SetUint(uint64((*(*uint16)(unsafe.Pointer(&b[firstByte])) << shift1) >> shift2))
		}
	case lsb+width <= 24:
		shift1 := 64 - ((lsb + width) - 16)
		shift2 := (shift1 + lsb) - 16
		return func(b []byte, field reflect.Value) {
			u := uint64(*(*uint16)(unsafe.Pointer(&b[firstByte]))) >> lsb
			field.SetUint(u | ((uint64(b[firstByte+2]) << shift1) >> shift2))
		}
	case lsb+width <= 32:
		shift1 := 32 - (lsb + width)
		shift2 := 32 - width
		return func(b []byte, field reflect.Value) {
			field.SetUint(uint64((*(*uint32)(unsafe.Pointer(&b[firstByte])) << shift1) >> shift2))
		}

	case lsb+width <= 40:
		shift1 := 64 - ((lsb + width) - 32)
		shift2 := (shift1 + lsb) - 32
		return func(b []byte, field reflect.Value) {
			u := uint64(*(*uint32)(unsafe.Pointer(&b[firstByte]))) >> lsb
			field.SetUint(u | ((uint64(b[firstByte+4]) << shift1) >> shift2))
		}
	case lsb+width <= 48:
		shift1 := 64 - ((lsb + width) - 32)
		shift2 := (shift1 + lsb) - 32
		return func(b []byte, field reflect.Value) {
			u := uint64(*(*uint32)(unsafe.Pointer(&b[firstByte]))) >> lsb
			field.SetUint(u | ((uint64(*(*uint16)(unsafe.Pointer(&b[firstByte+4]))) << shift1) >> shift2))
		}
	case lsb+width <= 56:
		shift1 := 32 - lsb
		shift2 := 64 - ((lsb + width) - 48)
		shift3 := (shift2 + lsb) - 48
		return func(b []byte, field reflect.Value) {
			u := uint64(*(*uint32)(unsafe.Pointer(&b[firstByte])) >> lsb)
			u |= uint64(*(*uint16)(unsafe.Pointer(&b[firstByte+4]))) >> shift1
			field.SetUint(u | (uint64(b[firstByte+6])<<shift2)>>shift3)
		}
	case lsb+width <= 64:
		shift1 := 64 - (lsb + width)
		shift2 := 64 - width
		return func(b []byte, field reflect.Value) {
			field.SetUint((*(*uint64)(unsafe.Pointer(&b[firstByte])) << shift1) >> shift2)
		}
	default:
		// Assume lsb+width <= 72
		shift1 := 128 - (lsb + width)
		shift2 := (shift1 + lsb) - 64
		return func(b []byte, field reflect.Value) {
			u := *(*uint64)(unsafe.Pointer(&b[firstByte])) >> lsb
			field.SetUint(u | ((uint64(b[firstByte+8]) << shift1) >> shift2))
		}
	}
}

func makeSignedSingleUnpacker(typ reflect.Type, ilsb uint64, width uint8) unpacker {
	firstByte := ilsb / 8
	lsb := uint8(ilsb % 8)
	switch {
	case lsb+width <= 8:
		shift1 := 64 - (lsb + width)
		shift2 := 64 - width
		return func(b []byte, field reflect.Value) {
			field.SetInt((int64(b[firstByte]) << shift1) >> shift2)
		}
	case lsb+width <= 16:
		shift1 := 64 - (lsb + width)
		shift2 := 64 - width
		return func(b []byte, field reflect.Value) {
			field.SetInt((int64(*(*uint16)(unsafe.Pointer(&b[firstByte]))) << shift1) >> shift2)
		}
	case lsb+width <= 24:
		shift1 := 64 - ((lsb + width) - 16)
		shift2 := (shift1 + lsb) - 16
		return func(b []byte, field reflect.Value) {
			i := int64(*(*uint16)(unsafe.Pointer(&b[firstByte]))) >> lsb
			field.SetInt(i | ((int64(b[firstByte+2]) << shift1) >> shift2))
		}
	case lsb+width <= 32:
		shift1 := 64 - (lsb + width)
		shift2 := 64 - width
		return func(b []byte, field reflect.Value) {
			field.SetInt((int64(*(*uint32)(unsafe.Pointer(&b[firstByte]))) << shift1) >> shift2)
		}

	case lsb+width <= 40:
		shift1 := 64 - ((lsb + width) - 32)
		shift2 := (shift1 + lsb) - 32
		return func(b []byte, field reflect.Value) {
			i := int64(*(*uint32)(unsafe.Pointer(&b[firstByte]))) >> lsb
			field.SetInt(i | ((int64(b[firstByte+4]) << shift1) >> shift2))
		}
	case lsb+width <= 48:
		shift1 := 64 - ((lsb + width) - 32)
		shift2 := (shift1 + lsb) - 32
		return func(b []byte, field reflect.Value) {
			i := int64(*(*uint32)(unsafe.Pointer(&b[firstByte]))) >> lsb
			field.SetInt(i | ((int64(*(*uint16)(unsafe.Pointer(&b[firstByte+4]))) << shift1) >> shift2))
		}
	case lsb+width <= 56:
		shift1 := 32 - lsb
		shift2 := 64 - ((lsb + width) - 48)
		shift3 := (shift2 + lsb) - 48
		return func(b []byte, field reflect.Value) {
			i := int64(*(*uint32)(unsafe.Pointer(&b[firstByte])) >> lsb)
			i |= int64(*(*uint16)(unsafe.Pointer(&b[firstByte+4]))) >> shift1
			field.SetInt(i | (int64(b[firstByte+6])<<shift2)>>shift3)
		}
	case lsb+width <= 64:
		shift1 := 64 - (lsb + width)
		shift2 := 64 - width
		return func(b []byte, field reflect.Value) {
			field.SetInt((*(*int64)(unsafe.Pointer(&b[firstByte])) << shift1) >> shift2)
		}
	default:
		// Assume lsb+width <= 72
		shift1 := 128 - (lsb + width)
		shift2 := (shift1 + lsb) - 64
		return func(b []byte, field reflect.Value) {
			i := int64(*(*uint64)(unsafe.Pointer(&b[firstByte])) >> lsb)
			field.SetInt(i | ((int64(b[firstByte+8]) << shift1) >> shift2))
		}
	}
}

func makeBoolSinglePacker(lsb uint64) packer {
	firstByte := lsb / 8
	tru := byte(1) << uint8(lsb%8)
	return func(b []byte, field reflect.Value) {
		if field.Bool() {
			b[firstByte] |= tru
		}
	}
}

func makeBoolSingleUnpacker(lsb uint64) unpacker {
	firstByte := lsb / 8
	tru := byte(1) << uint8(lsb%8)
	return func(b []byte, field reflect.Value) {
		field.SetBool(b[firstByte]&tru > 0)
	}
}
