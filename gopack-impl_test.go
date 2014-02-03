// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopack

import (
	"fmt"
	gopack_testing "github.com/joshlf13/gopack/testing"
	"math"
	"reflect"
	"testing"
)

func TestMakePacker(t *testing.T) {
	type typ struct {
		f1 uint8
	}

	p := makePackerWrapper(reflect.TypeOf(typ{}))
	b := []byte{0}
	p(b, reflect.ValueOf(typ{127}))

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
	u := makeUnpackerWrapper(reflect.TypeOf(val))
	u(b, reflect.ValueOf(&val).Elem())
	if (val != typ{127}) {
		t.Fatalf("Expected %v; got %v", typ{127}, val)
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
	p := makePackerWrapper(reflect.TypeOf(val))
	p(b[:], reflect.ValueOf(val))
	if b != [...]byte{127, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{127, 255}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(val))
	u(b[:], reflect.ValueOf(&val).Elem())
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
	p := makePackerWrapper(reflect.TypeOf(val))
	p(b[:], reflect.ValueOf(val))
	if b != [...]byte{53, 1} {
		t.Fatalf("Expected %v; got %v", [...]byte{53, 1}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(val))
	u(b[:], reflect.ValueOf(&val).Elem())
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
	p := makePackerWrapper(reflect.TypeOf(val))
	p(b[:], reflect.ValueOf(val))
	if b != [...]byte{2} {
		t.Fatalf("Expected %v; got %v", [...]byte{2}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(val))
	u(b[:], reflect.ValueOf(&val).Elem())
	if val != (typ{false, true}) {
		t.Fatalf("Expected %v; got %v", typ{false, true}, val)
	}
}

func TestUnexported(t *testing.T) {
	var b [1]byte
	val := gopack_testing.MakeTyp("hi", 255)
	p := makePackerWrapper(reflect.TypeOf(val))
	p(b[:], reflect.ValueOf(val))
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = gopack_testing.Typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(val))
	u(b[:], reflect.ValueOf(&val).Elem())
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
	p := makePackerWrapper(reflect.TypeOf(&val))
	p(b[:], reflect.ValueOf(&val))
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(&val))
	u(b[:], reflect.ValueOf(&val))
	if val != (typ{255}) {
		t.Fatalf("Expected %v; got %v", typ{255}, val)
	}
}

func TestSignedPacking(t *testing.T) {
	type typ struct {
		F1 int8
	}

	var b [1]byte
	val := typ{-1}
	p := makePackerWrapper(reflect.TypeOf(&val))
	p(b[:], reflect.ValueOf(&val))
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(&val))
	u(b[:], reflect.ValueOf(&val))
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
	p := makePackerWrapper(reflect.TypeOf(&val))
	p(b[:], reflect.ValueOf(&val))
	if b != [...]byte{127} {
		t.Fatalf("Expected %v; got %v", [...]byte{127}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(&val))
	u(b[:], reflect.ValueOf(&val))
	if val != (typ{-1, false}) {
		t.Fatalf("Expected %v; got %v", typ{-1, false}, val)
	}

	val = typ{-1, true}
	p(b[:], reflect.ValueOf(&val))
	if b != [...]byte{255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255}, b)
	}

	val = typ{}
	u(b[:], reflect.ValueOf(&val))
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
	p := makePackerWrapper(reflect.TypeOf(&val))
	p(b[:], reflect.ValueOf(&val))
	if b != [...]byte{255, 255, 255, 255, 255, 255, 255, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255, 255, 255, 255, 255, 255, 255, 255}, b)
	}

	val = typ{}
	u := makeUnpackerWrapper(reflect.TypeOf(&val))
	u(b[:], reflect.ValueOf(&val))
	if val != (typ{math.MaxUint64}) {
		t.Fatalf("Expected %v; got %v", typ{math.MaxUint64}, val)
	}

	type typ1 struct {
		F1 int64
	}

	val2 := typ1{-1}
	p = makePackerWrapper(reflect.TypeOf(&val2))
	p(b[:], reflect.ValueOf(&val2))
	if b != [...]byte{255, 255, 255, 255, 255, 255, 255, 255} {
		t.Fatalf("Expected %v; got %v", [...]byte{255, 255, 255, 255, 255, 255, 255, 255}, b)
	}

	val2 = typ1{}
	u = makeUnpackerWrapper(reflect.TypeOf(&val2))
	u(b[:], reflect.ValueOf(&val2))
	if val2 != (typ1{-1}) {
		t.Fatalf("Expected %v; got %v", typ1{-1}, val)
	}
}

func TestErrors(t *testing.T) {
	i := 0
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		p := makePackerWrapper(reflect.TypeOf(i))
		p(nil, reflect.Value{})
	})
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		p := makePackerWrapper(reflect.TypeOf(&i))
		p(nil, reflect.Value{})
	})
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(i))
		u(nil, reflect.Value{})
	})
	testError(t, Error{fmt.Errorf("gopack: non-struct type int")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(&i))
		u(nil, reflect.Value{})
	})

	type typ struct {
		F1 string
	}
	testError(t, Error{fmt.Errorf("gopack: non-packable type string")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ{}))
		p(nil, reflect.ValueOf(typ{}))
	})
	testError(t, Error{fmt.Errorf("gopack: non-packable type string")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(typ{}))
		u(nil, reflect.ValueOf(typ{}))
	})

	type typ1 struct {
		F1 uint8 `gopack:"numerals"`
	}
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\": strconv.ParseInt: parsing \"numerals\": invalid syntax")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ1{}))
		p(nil, reflect.ValueOf(typ1{}))
	})
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\": strconv.ParseInt: parsing \"numerals\": invalid syntax")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(typ1{}))
		u(nil, reflect.ValueOf(typ1{}))
	})

	type typ2 struct {
		F1 int8 `gopack:"9"`
	}
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" (type int8) too wide (9)")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ2{}))
		p(nil, reflect.ValueOf(typ2{}))
	})
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" (type int8) too wide (9)")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(typ2{}))
		u(nil, reflect.ValueOf(typ2{}))
	})

	type typ3 struct {
		F1 uint8 `gopack:"0"`
	}
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" too small (0)")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ3{}))
		p(nil, reflect.ValueOf(typ3{}))
	})
	testError(t, Error{fmt.Errorf("gopack: struct tag on field \"F1\" too small (0)")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(typ3{}))
		u(nil, reflect.ValueOf(typ3{}))
	})

	type typ4 struct {
		F1 uint8 `gopack:"4"`
		F2 int8  `gopack:"4"`
	}
	testError(t, Error{fmt.Errorf("gopack: value out of range: max 15; got 16")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ4{}))
		p(make([]byte, 1), reflect.ValueOf(typ4{15, 0}))
		p(make([]byte, 1), reflect.ValueOf(typ4{16, 0}))
	})
	testError(t, Error{fmt.Errorf("gopack: value out of range: max 7, min -8; got 8")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ4{}))
		p(make([]byte, 1), reflect.ValueOf(typ4{15, -8}))
		p(make([]byte, 1), reflect.ValueOf(typ4{15, 8}))
	})
	testError(t, Error{fmt.Errorf("gopack: value out of range: max 7, min -8; got -9")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ4{}))
		p(make([]byte, 1), reflect.ValueOf(typ4{15, -9}))
	})

	type typ5 struct {
		F1 uint8
	}
	type typ6 struct {
		F1 *typ5
	}
	testError(t, Error{fmt.Errorf("gopack: non-packable type *gopack.typ5")}, func() {
		p := makePackerWrapper(reflect.TypeOf(typ6{}))
		p(nil, reflect.ValueOf(typ6{}))
	})
	testError(t, Error{fmt.Errorf("gopack: non-packable type *gopack.typ5")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(typ6{}))
		u(nil, reflect.ValueOf(typ6{}))
	})

	type typ7 struct {
		F1 uint8
		F2 uint8 `gopack:"1"`
	}

	t1 := typ7{255, 255}
	bytes := []byte{}
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 2)")}, func() {
		p := makePackerWrapper(reflect.TypeOf(t1))
		p(bytes, reflect.ValueOf(t1))
	})
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 2)")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(&t1))
		u(bytes, reflect.ValueOf(&t1))
	})

	type typ8 struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t2 := typ8{255, 255, 255, 255, 255, 255, 255, 255}
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 8)")}, func() {
		p := makePackerWrapper(reflect.TypeOf(t2))
		p(bytes, reflect.ValueOf(t2))
	})
	testError(t, Error{fmt.Errorf("gopack: buffer too small (0; need 8)")}, func() {
		u := makeUnpackerWrapper(reflect.TypeOf(&t2))
		u(bytes, reflect.ValueOf(&t2))
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

func BenchmarkMakeEmptyPacker(b *testing.B) {
	type typ struct {
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkMakePacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := reflect.TypeOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makePackerWrapper(t)
	}
}

func BenchmarkUseEmptyPacker(b *testing.B) {
	type typ struct {
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(typ{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker1Field(b *testing.B) {
	type typ struct {
		F1 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 1)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker2Fields(b *testing.B) {
	type typ struct {
		F1, F2 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 2)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker4Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 4)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker8Fields(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 uint8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 8)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkUsePacker8FieldsSigned(b *testing.B) {
	type typ struct {
		F1, F2, F3, F4, F5, F6, F7, F8 int8
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 8)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkNested1Level(b *testing.B) {
	type typ struct {
		F1 struct {
		}
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkNested2Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
			}
		}
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkNested4Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
				F3 struct {
					F4 struct {
					}
				}
			}
		}
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}

func BenchmarkNested8Levels(b *testing.B) {
	type typ struct {
		F1 struct {
			F2 struct {
				F3 struct {
					F4 struct {
						F5 struct {
							F6 struct {
								F7 struct {
									F8 struct {
									}
								}
							}
						}
					}
				}
			}
		}
	}

	t := reflect.TypeOf(typ{})
	p := makePackerWrapper(t)
	bytes := make([]byte, 0)
	val := reflect.ValueOf(&typ{}).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(bytes, val)
	}
}
