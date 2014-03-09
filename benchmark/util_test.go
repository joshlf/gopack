package benchmark

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
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

func randByte() byte {
	return byte(randUint64Bits(8))
}

func randBytes(b []byte) {
	for i := range b {
		b[i] = randByte()
	}
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

func randInstance(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()
	for i := 0; i < typ.NumField(); i++ {
		switch typ.Field(i).Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n := getFieldWidth(typ.Field(i))
			val.Field(i).SetInt(randInt64Bits(uint8(n)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n := getFieldWidth(typ.Field(i))
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

func getFieldWidth(field reflect.StructField) uint64 {
	str := field.Tag.Get("gopack")
	if str == "" {
		return uint64(field.Type.Bits())
	}
	n, _ := strconv.Atoi(str)
	return uint64(n)
}
