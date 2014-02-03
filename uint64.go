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
	msk := mask(width)
	msk <<= lsb
	msk = ^msk
	target &= msk

	val <<= lsb
	return target | val
}

// Unpack the value stored in [lsb, lsb + width) in target.
func UnpackUnsigned(target uint64, lsb, width uint8) uint64 {
	return (target >> lsb) & mask(width)
}

// Pack val into bits [lsb, lsb + width) in target,
// returning the new value of target. If val doesn't
// fit in width bits, the behavior of PackSigned is
// undefined.
func PackSigned(target uint64, val int64, lsb, width uint8) uint64 {
	uval := *(*uint64)(unsafe.Pointer(&val))
	// If val is negative, there will
	// be 1's outside of the target range.
	uval &= mask(width)
	return PackUnsigned(target, uval, lsb, width)
}

// Unpack the value stored in [lsb, lsb + width) in target.
func UnpackSigned(target uint64, lsb, width uint8) int64 {
	target >>= lsb
	msk := mask(width)
	target &= msk

	// The return value of fillFirstBit
	// should be either all 0s or all 1s
	// depending on the value of the msb
	// of the target range.
	target |= (fillFirstBit((target>>(width-1))&1) & ^msk)
	return *(*int64)(unsafe.Pointer(&target))
}
