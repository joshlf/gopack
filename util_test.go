// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"unsafe"
)

func randUint64() uint64 {
	return uint64(rand.Uint32()) | (uint64(rand.Uint32()) << 32)
}

func randUint64Bits(bits uint8) uint64 {
	return randUint64() >> (64 - bits)
}

func randInt64() int64 {
	u := randUint64()
	return *(*int64)(unsafe.Pointer(&u))
}

func randInt64Bits(bits uint8) int64 {
	u := randUint64Bits(bits)
	return (*(*int64)(unsafe.Pointer(&u)) << (64 - bits)) >> (64 - bits)
}

func randBool() bool {
	return rand.Int()%2 == 0
}

func randWidthLSBPair() (uint8, uint8) {
	width := uint8(1 + (rand.Uint32() & 0x3F)) // Restrict to range [1, 64]
	lsb := uint8(rand.Uint32() & 0x3F)         // Restrict to range [0, 63]
	for width+lsb > 64 {
		width = uint8(1 + (rand.Uint32() & 0x3F)) // Restrict to range [1, 64]
		lsb = uint8(rand.Uint32() & 0x3F)         // Restrict to range [0, 63]
	}
	return width, lsb
}

func nOnes(n int) uint64 {
	var u uint64
	for i := 0; i < n; i++ {
		u = (u << 1) | 1
	}
	return u
}

// func randInstance(typ reflect.Type) reflect.Value {
// 	val := reflect.New(typ).Elem()
// 	for i := 0; i < typ.NumField(); i++ {
// 		switch typ.Field(i).Type.Kind() {
// 		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 			n, _ := getFieldWidth(typ.Field(i))
// 			val.Field(i).SetInt(randInt64Bits(uint8(n)))
// 		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 			n, _ := getFieldWidth(typ.Field(i))
// 			val.Field(i).SetUint(randUint64Bits(uint8(n)))
// 		case reflect.Bool:
// 			val.Field(i).SetBool(randBool())
// 		default:
// 			panic(fmt.Sprint("Cannot generate type:", typ.Field(i).Type))
// 		}
// 	}
// 	return val
// }

func randInstance(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()
	for i := 0; i < typ.NumField(); i++ {
		switch typ.Field(i).Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, _ := getFieldWidth(typ.Field(i))
			val.Field(i).SetInt(randInt64Bits(uint8(n)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, _ := getFieldWidth(typ.Field(i))
			val.Field(i).SetUint(randUint64Bits(uint8(n)))
		case reflect.Bool:
			val.Field(i).SetBool(randBool())
		case reflect.Struct:
			val.Field(i).Set(randInstance(typ.Field(i).Type))
		default:
			panic(fmt.Sprint("Cannot generate type:", typ.Field(i).Type))
		}
	}
	return val
}

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

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(nsprintf(1, "%v", err))
	}
}

func mustError(t *testing.T, err error) {
	if err == nil {
		t.Fatal(nsprintf(1, "unexpected nil error"))
	}
	if err.Error() == "internal server error" {
		t.Fatal(nsprintf(1, "unexpected internal server error"))
	}
}

// nsprintf returns a formatted log line
// that includes the file/line number
// of the nth caller up the stack
func nsprintf(n int, format string, args ...interface{}) string {
	_, file, line, ok := runtime.Caller(n + 1)
	if !ok {
		return fmt.Sprintf("unknown file/line: "+format, args...)
	}
	file = filepath.Base(file)
	return fmt.Sprintf("%v:%v: "+format, append([]interface{}{file, line}, args...)...)
}
