// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"unsafe"
)

// lsb must be < 8
// b must be at least 8 bytes,
// and at least lsb+width bits
func packUint64(b []byte, val uint64, lsb, width uint8) {
	*(*uint64)(unsafe.Pointer(&b[0])) |= (val << lsb)
	lsbPlusWidth := lsb + width
	if lsbPlusWidth > 64 {
		*(*uint64)(unsafe.Pointer(&b[1])) |= ((val >> (64 - lsb)) << (128 - lsbPlusWidth))
	}
}

func unpackUint64(b []byte, lsb, width uint8) uint64 {
	u := *(*uint64)(unsafe.Pointer(&b[0]))
	u >>= lsb
	lsbPlusWidth := lsb + width
	if lsbPlusWidth > 64 {
		u |= (*(*uint64)(unsafe.Pointer(&b[1])) >> (128 - lsbPlusWidth)) << (64 - lsb)
	}
	// We could have pulled in bits from the subsequent
	// packed value in b, so zero out the high bits.
	return (u << (64 - width)) >> (64 - width)
}

// lsb must be < 8
// b must be at least 8 bytes,
// and at least lsb+width bits
func packInt64(b []byte, val int64, lsb, width uint8) {
	uval := *(*uint64)(unsafe.Pointer(&val))
	uval = (uval << (64 - width)) >> (64 - width)
	packUint64(b, uval, lsb, width)
}

func unpackInt64(b []byte, lsb, width uint8) int64 {
	uval := unpackUint64(b, lsb, width)
	ones := fillFirstBit(uval >> (width - 1))
	uval |= ones << width
	return *(*int64)(unsafe.Pointer(&uval))
}
