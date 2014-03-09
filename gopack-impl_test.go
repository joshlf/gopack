// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"

	gopack_testing "github.com/joshlf13/gopack/testing"
)

func TestMakePacker(t *testing.T) {
	type typ struct {
		f1 uint8
	}

	var b [1]byte
	Pack(b[:], typ{127})

	// Since typ.f1 isn't exported,
	// it shouldn't be packed
	if b[0] != 0 {
		t.Fatalf("Expected %v; got %v", 127, b[0])
	}
}

func TestMakeUnpacker(t *testing.T) {
	type typ struct {
		F1 uint8
	}

	b := []byte{127}
	val := typ{}
	Unpack(b, &val)
	if (val != typ{127}) {
		t.Fatalf("Expected %v; got %v", typ{127}, val)
	}

	val = typ{}
	Unpack(b, val)
	if val != (typ{}) {
		t.Fatalf("Expected %v; got %v", typ{}, val)
	}
}

func TestMultipleFields(t *testing.T) {
	type typ struct {
		F1, F2 uint8
	}

	var b [2]byte
	val := typ{127, 255}
	Pack(b[:], val)
	if b != [...]byte{127, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{127, 255}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{127, 255}) {
		t.Fatalf("Expected %v; got %v", typ{127, 255}, val)
	}
}

func TestCoverUnsigned(t *testing.T) {
	rand.Seed(20421)
	testCover(t, struct{ F1 uint8 }{})
	testCover(t, struct {
		F1 uint8 `gopack:"4"`
	}{})

	testCover(t, struct {
		F1 uint16
	}{})
	testCover(t, struct {
		F1 bool
		F2 uint16 `gopack:"15"`
	}{})

	testCover(t, struct {
		F1 uint32 `gopack:"24"`
	}{})
	testCover(t, struct {
		F1 uint8 `gopack:"7"`
		F2 uint16
	}{})

	testCover(t, struct {
		F1 uint32 `gopack:"32"`
	}{})
	testCover(t, struct {
		F1 bool
		F2 uint32 `gopack:"31"`
	}{})

	testCover(t, struct {
		F1 uint64 `gopack:"40"`
	}{})
	testCover(t, struct {
		F1 uint8  `gopack:"7"`
		F2 uint32 `gopack:"32"`
	}{})

	testCover(t, struct {
		F1 uint64 `gopack:"48"`
	}{})

	testCover(t, struct {
		F1 uint64
	}{})

	testCover(t, struct {
		F1 bool
		F2 uint64 `gopack:"63"`
	}{})

	testCover(t, struct {
		F1 uint8  `gopack:"7"`
		F2 uint64 `gopack:"63"`
	}{})

	testCover(t, struct {
		F1 uint8 `gopack:"7"`
		F2 uint64
	}{})
}

func TestCoverSigned(t *testing.T) {
	rand.Seed(18140)
	testCover(t, struct{ F1 int8 }{})
	testCover(t, struct {
		F1 int8 `gopack:"4"`
	}{})

	testCover(t, struct {
		F1 int16
	}{})
	testCover(t, struct {
		F1 bool
		F2 int16 `gopack:"15"`
	}{})

	testCover(t, struct {
		F1 int32 `gopack:"24"`
	}{})
	testCover(t, struct {
		F1 int8 `gopack:"7"`
		F2 int16
	}{})

	testCover(t, struct {
		F1 int32 `gopack:"32"`
	}{})
	testCover(t, struct {
		F1 bool
		F2 int32 `gopack:"31"`
	}{})

	testCover(t, struct {
		F1 int64 `gopack:"40"`
	}{})
	testCover(t, struct {
		F1 int8  `gopack:"7"`
		F2 int32 `gopack:"32"`
	}{})

	testCover(t, struct {
		F1 int64 `gopack:"48"`
	}{})

	testCover(t, struct {
		F1 int64
	}{})

	testCover(t, struct {
		F1 bool
		F2 int64 `gopack:"63"`
	}{})

	testCover(t, struct {
		F1 int8  `gopack:"7"`
		F2 int64 `gopack:"63"`
	}{})

	testCover(t, struct {
		F1 int8 `gopack:"7"`
		F2 int64
	}{})
}

func testCover(t *testing.T, v interface{}) {
	typ := reflect.TypeOf(v)
	p := makePackerWrapper(typ)
	u := makeUnpackerWrapper(reflect.PtrTo(typ))
	_, n, _ := makePacker(0, typ)
	b := make([]byte, (int(n)/8)+1)

	// Note: increasing the iterations to 1000*1000
	// will cause the full test suite to take ~30s
	for i := 0; i < 1000*100; i++ {
		val1 := randInstance(typ)
		val2 := reflect.New(typ)
		p(b, val1)
		u(b, val2)
		if val2.Elem().Interface() != val1.Interface() {
			t.Fatalf("Expected \n%v; got \n%v\n(on type %v)", val1.Interface(), val2.Elem().Interface(), typ)
		}
	}
}

func TestByteBoundaries(t *testing.T) {
	rand.Seed(235)
	type typ struct {
		F1 uint8 `gopack:"1"`
		F2 uint8 `gopack:"5"`
		F3 uint8 `gopack:"2"` // Reallign

		F4 uint16 `gopack:"3"`
		F5 uint16 `gopack:"9"`
		F6 uint8  `gopack:"4"` // Reallign

		F7 uint32 `gopack:"5"`
		F8 uint32 `gopack:"18"`
		F9 uint32 `gopack:"1"` // Reallign

		F10 uint64 `gopack:"3"`
		F11 uint64 `gopack:"35"`
		F12 uint64 `gopack:"2"` // Reallign

		F13 uint64 `gopack:"3"`
		F14 uint64 `gopack:"43"`
		F15 uint64 `gopack:"2"` // Reallign

		F16 uint64 `gopack:"3"`
		F17 uint64 `gopack:"63"`
	}

	val := typ{}
	p := makePackerWrapper(reflect.TypeOf(val))
	u := makeUnpackerWrapper(reflect.TypeOf(&val))
	var b [32]byte

	for i := 0; i < 1000*1000; i++ {
		val = typ{
			uint8(randUint64Bits(1)),
			uint8(randUint64Bits(5)),
			uint8(randUint64Bits(2)),
			uint16(randUint64Bits(3)),
			uint16(randUint64Bits(9)),
			uint8(randUint64Bits(4)),
			uint32(randUint64Bits(5)),
			uint32(randUint64Bits(18)),
			uint32(randUint64Bits(1)),
			randUint64Bits(3),
			randUint64Bits(35),
			randUint64Bits(2),
			randUint64Bits(3),
			randUint64Bits(43),
			randUint64Bits(2),
			randUint64Bits(3),
			randUint64Bits(63),
		}
		val2 := typ{}
		p(b[:], reflect.ValueOf(val))
		u(b[:], reflect.ValueOf(&val2))
		if val2 != val {
			t.Fatalf("Expected \n%v; got \n%v", val, val2)
		}
	}
}

func TestNesting(t *testing.T) {
	type typ1 struct {
		F1 uint8
	}
	type typ struct {
		F1 uint8
		F2 typ1
	}

	var b [2]byte
	val := typ{127, typ1{255}}
	Pack(b[:], val)
	if b != [...]byte{127, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{127, 255}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{127, typ1{255}}) {
		t.Fatalf("Expected %v; got %v", typ{127, typ1{255}}, val)
	}
}

func TestTags(t *testing.T) {
	type typ struct {
		F1 uint8 `gopack:"5"`
		F2 uint8 `gopack:"4"`
	}

	var b [2]byte
	val := typ{21, 9}
	Pack(b[:], val)
	if b != [...]byte{53, 1} {
		t.Fatalf("Expected %v; got %v", [...]byte{53, 1}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{21, 9}) {
		t.Fatalf("Expected %v; got %v", typ{21, 9}, val)
	}
}

func TestBoolTag(t *testing.T) {
	type typ struct {
		F1 bool `gopack:"2"`
		F2 bool `gopack:"0"`
	}

	var b [1]byte
	val := typ{false, true}
	Pack(b[:], val)
	if b != [...]byte{2} {
		t.Fatalf("Expected %v; got %v", [...]byte{2}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{false, true}) {
		t.Fatalf("Expected %v; got %v", typ{false, true}, val)
	}
}

func TestUnexported(t *testing.T) {
	var b [1]byte
	val := gopack_testing.MakeTyp("hi", 255)
	Pack(b[:], val)
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = gopack_testing.Typ{}
	Unpack(b[:], &val)
	if val != gopack_testing.MakeTyp("", 255) {
		t.Fatalf("Expected %v; got %v", gopack_testing.MakeTyp("", 255), val)
	}
}

func TestPointer(t *testing.T) {
	type typ struct {
		F1 uint8
	}

	var b [1]byte
	val := typ{255}
	Pack(b[:], &val)
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{255}) {
		t.Fatalf("Expected %v; got %v", typ{255}, val)
	}

	val = typ{}
	Unpack(b[:], val)
	if val != (typ{}) {
		t.Fatalf("Expected %v; got %v", typ{}, val)
	}
}

func TestSignedPacking(t *testing.T) {
	type typ struct {
		F1 int8
	}

	var b [1]byte
	val := typ{-1}
	Pack(b[:], val)
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{-1}) {
		t.Fatalf("Expected %v; got %v", typ{-1}, val)
	}
}

func TestBoolPacking(t *testing.T) {
	type typ struct {
		F1 int8 `gopack:"7"`
		F2 bool
	}

	var b [1]byte
	val := typ{-1, false}
	Pack(b[:], val)
	if b != [...]byte{127} {
		t.Fatalf("Expected %v; got %v", [...]byte{127}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{-1, false}) {
		t.Fatalf("Expected %v; got %v", typ{-1, false}, val)
	}

	val = typ{-1, true}
	Pack(b[:], val)
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{-1, true}) {
		t.Fatalf("Expected %v; got %v", typ{-1, true}, val)
	}
}

func Test64Bit(t *testing.T) {
	type typ struct {
		F1 uint64
	}

	var b [8]byte
	val := typ{math.MaxUint64}
	Pack(b[:], &val)
	if b != [...]byte{255, 255, 255, 255, 255, 255, 255, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255, 255, 255, 255, 255, 255, 255, 255}, b)
	}

	val = typ{}
	Unpack(b[:], &val)
	if val != (typ{math.MaxUint64}) {
		t.Fatalf("Expected %v; got %v", typ{math.MaxUint64}, val)
	}

	type typ1 struct {
		F1 int64
	}

	val2 := typ1{-1}
	Pack(b[:], &val2)
	if b != [...]byte{255, 255, 255, 255, 255, 255, 255, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255, 255, 255, 255, 255, 255, 255, 255}, b)
	}

	val2 = typ1{}
	Unpack(b[:], &val2)
	if val2 != (typ1{-1}) {
		t.Fatalf("Expected %v; got %v", typ1{-1}, val)
	}
}

func TestErrors(t *testing.T) {
	i := 0
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		Pack(nil, 0)
	})
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		Pack(nil, &i)
	})
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		Unpack(nil, i)
	})
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		Pack(nil, i)
	})

	type typ struct {
		F1 string
	}
	testError(t, Error{fmt.Errorf("gopack: non-packable type string")}, func() {
		Pack(nil, typ{})
	})
	testError(t, Error{fmt.Errorf("gopack: non-packable type string")}, func() {
		Unpack(nil, typ{})
	})

	type typ1 struct {
		F1 uint8 `gopack:"numerals"`
	}
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\": strconv.ParseInt: parsing \"numerals\": invalid syntax")}, func() {
		Pack(nil, typ1{})
	})
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\": strconv.ParseInt: parsing \"numerals\": invalid syntax")}, func() {
		Unpack(nil, typ1{})
	})

	type typ2 struct {
		F1 int8 `gopack:"9"`
	}
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" (type int8) too wide (9)")}, func() {
		Pack(nil, typ2{})
	})
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" (type int8) too wide (9)")}, func() {
		Unpack(nil, typ2{})
	})

	type typ3 struct {
		F1 uint8 `gopack:"0"`
	}
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" too small (0)")}, func() {
		Pack(nil, typ3{})
	})
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" too small (0)")}, func() {
		Pack(nil, typ3{})
	})

	type typ4 struct {
		F1 uint8 `gopack:"4"`
		F2 int8  `gopack:"4"`
	}
	testError(t, Error{fmt.Errorf("gopack: value out of range: max 15; got 16")}, func() {
		Pack(make([]byte, 1), typ4{15, 0})
		Pack(make([]byte, 1), typ4{16, 0})
	})
	testError(t, Error{fmt.Errorf("gopack: value out of range: max 7, min -8; got 8")}, func() {
		Pack(make([]byte, 1), typ4{15, -8})
		Pack(make([]byte, 1), typ4{15, 8})
	})
	testError(t, Error{fmt.Errorf("gopack: value out of range: max 7, min -8; got -9")}, func() {
		Pack(make([]byte, 1), typ4{15, -8})
		Pack(make([]byte, 1), typ4{15, -9})
	})

	type typ5 struct {
		F1 uint8
	}
	type typ6 struct {
		F1 *typ5
	}
	testError(t, Error{fmt.Errorf("gopack: non-packable type *gopack.typ5")}, func() {
		Pack(nil, typ6{})
	})
	testError(t, Error{fmt.Errorf("gopack: non-packable type *gopack.typ5")}, func() {
		Unpack(nil, typ6{})
	})

	type typ7 struct {
		F1 uint8
		F2 uint8 `gopack:"1"`
	}

	t1 := typ7{255, 255}
	bytes := []byte{}
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 2)")}, func() {
		Pack(bytes, t1)
	})
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 2)")}, func() {
		Unpack(bytes, &t1)
	})

	type typ8 struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t2 := typ8{255, 255, 255, 255, 255, 255, 255, 255}
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 8)")}, func() {
		Pack(bytes, t2)
	})
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 8)")}, func() {
		Unpack(bytes, &t2)
	})

	// Make sure that Unpack reports errors
	// even for non-pointer types
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		Unpack(nil, 0)
	})
	testError(t, Error{fmt.Errorf("gopack: non-packable type *uint8")}, func() {
		type typ9 struct {
			F1 *uint8
		}
		Unpack(nil, typ9{})
	})
}

func testError(t *testing.T, err interface{}, f func()) {
	defer func() {
		r := recover()
		if r.(Error).Error() != err.(Error).Error() {
			t.Fatalf("Expected error \"%v\"; got \"%v\"", err, r)
		}
	}()
	f()
}
