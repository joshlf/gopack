package benchmark

import (
	"testing"

	"github.com/joshlf13/gopack"
)

func BenchmarkPackNestingOverhead1Level(b *testing.B) {
	type typ struct {
		F1 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackNestingOverhead2Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type typ struct {
		F1 t1
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackNestingOverhead4Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 t1
	}
	type t3 struct {
		F1 t2
	}
	type typ struct {
		F1 t3
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackNestingOverhead8Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 t1
	}
	type t3 struct {
		F1 t2
	}
	type t4 struct {
		F1 t3
	}
	type t5 struct {
		F1 t4
	}
	type t6 struct {
		F1 t5
	}
	type t7 struct {
		F1 t6
	}
	type typ struct {
		F1 t7
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}
func BenchmarkUnpackNestingOverhead1Level(b *testing.B) {
	type typ struct {
		F1 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackNestingOverhead2Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type typ struct {
		F1 t1
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackNestingOverhead4Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 t1
	}
	type t3 struct {
		F1 t2
	}
	type typ struct {
		F1 t3
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackNestingOverhead8Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 t1
	}
	type t3 struct {
		F1 t2
	}
	type t4 struct {
		F1 t3
	}
	type t5 struct {
		F1 t4
	}
	type t6 struct {
		F1 t5
	}
	type t7 struct {
		F1 t6
	}
	type typ struct {
		F1 t7
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

/*
	Multiplicative
*/

func BenchmarkPackNestingMultiplicative1Level(b *testing.B) {
	type typ struct {
		F1 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackNestingMultiplicative2Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type typ struct {
		F1 uint8
		F2 t1
	}
	intface, bytes := benchmarkUtil(typ{}, 2)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackNestingMultiplicative4Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 uint8
		F2 t1
	}
	type t3 struct {
		F1 uint8
		F2 t2
	}
	type typ struct {
		F1 uint8
		F2 t3
	}
	intface, bytes := benchmarkUtil(typ{}, 4)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkPackNestingMultiplicative8Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 uint8
		F2 t1
	}
	type t3 struct {
		F1 uint8
		F2 t2
	}
	type t4 struct {
		F1 uint8
		F2 t3
	}
	type t5 struct {
		F1 uint8
		F2 t4
	}
	type t6 struct {
		F1 uint8
		F2 t5
	}
	type t7 struct {
		F1 uint8
		F2 t6
	}
	type typ struct {
		F1 uint8
		F2 t7
	}
	intface, bytes := benchmarkUtil(typ{}, 8)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Pack(bytes, val)
	}
}

func BenchmarkUnpackNestingMultiplicative1Level(b *testing.B) {
	type typ struct {
		F1 uint8
	}
	intface, bytes := benchmarkUtil(typ{}, 1)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackNestingMultiplicative2Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type typ struct {
		F1 uint8
		F2 t1
	}
	intface, bytes := benchmarkUtil(typ{}, 2)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackNestingMultiplicative4Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 uint8
		F2 t1
	}
	type t3 struct {
		F1 uint8
		F2 t2
	}
	type typ struct {
		F1 uint8
		F2 t3
	}
	intface, bytes := benchmarkUtil(typ{}, 4)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}

func BenchmarkUnpackNestingMultiplicative8Levels(b *testing.B) {
	type t1 struct {
		F1 uint8
	}
	type t2 struct {
		F1 uint8
		F2 t1
	}
	type t3 struct {
		F1 uint8
		F2 t2
	}
	type t4 struct {
		F1 uint8
		F2 t3
	}
	type t5 struct {
		F1 uint8
		F2 t4
	}
	type t6 struct {
		F1 uint8
		F2 t5
	}
	type t7 struct {
		F1 uint8
		F2 t6
	}
	type typ struct {
		F1 uint8
		F2 t7
	}
	intface, bytes := benchmarkUtil(typ{}, 8)
	val := intface.(typ)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopack.Unpack(bytes, &val)
	}
}
