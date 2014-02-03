// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

// Reflection reports Typ.f1 as exported
// when it's defined in the same package
// as the reflection is called from, so
// we need to define this type in another
// package.
type Typ struct {
	f1 string
	F1 uint8
}

func MakeTyp(f1 string, F1 uint8) Typ {
	return Typ{f1, F1}
}
