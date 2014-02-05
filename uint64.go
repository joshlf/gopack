// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"unsafe"
)

// Pack val into bits [lsb, lsb + width) in target,
// returning the new value of target. If val doesn't
// fit in width bits, the behavior of PackUnsigned
// is undefined.
func PackUnsigned(target, val uint64, lsb, width uint8) uint64 {
	// Zero out target region
	msk := (uint64(1) << width) - 1
	msk <<= lsb
	msk = ^msk
	target &= msk

	val <<= lsb
	return target | val
}

// Unpack the value stored in [lsb, lsb + width) in target.
func UnpackUnsigned(target uint64, lsb, width uint8) uint64 {
	return (target >> lsb) & ((uint64(1) << width) - 1)
}

// Pack val into bits [lsb, lsb + width) in target,
// returning the new value of target. If val doesn't
// fit in width bits, the behavior of PackSigned is
// undefined.
func PackSigned(target uint64, val int64, lsb, width uint8) uint64 {
	uval := *(*uint64)(unsafe.Pointer(&val))
	// If val is negative, there will
	// be 1's outside of the target range.
	msk := (uint64(1) << width) - 1
	uval &= msk

	msk <<= lsb
	msk = ^msk
	target &= msk

	uval <<= lsb
	return target | uval
}

// Unpack the value stored in [lsb, lsb + width) in target.
func UnpackSigned(target uint64, lsb, width uint8) int64 {
	uval := (target >> lsb) & ((uint64(1) << width) - 1)
	val := *(*int64)(unsafe.Pointer(&uval))
	return (val << (64 - width)) >> (64 - width)
}
