// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

// If the lsb of u is 1, return all 1's,
// otherwise return all 0's
func fillFirstBit(u uint64) uint64 {
	u |= u << 1
	u |= u << 2
	u |= u << 4
	u |= u << 8
	u |= u << 16
	return u | (u << 32)
}

// Make a mask consisting of all 0's
// followed by width 1's
func mask(width uint8) uint64 {
	msk := uint64(1)
	msk <<= width
	return msk - 1
}
