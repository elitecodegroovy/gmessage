package simple_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

var s []string = []string{}
var junk map[int]string

func init() {
	junk = make(map[int]string)
	for i := 0; i < 10000; i++ {
		junk[i] = strconv.Itoa(i + 10000)
	}
}

// 定义获取一个字符
func nextNumString() func() string {
	n := 0
	// closure captures variable n
	return func() string {
		n += 1
		return strconv.Itoa(n)
	}
}

var global string

// benchmarkNaiveConcat provides a benchmark for basic built-in
// Go string concatenation. Because strings are immutable in Go,
// it performs the worst of the tested methods. The time taken to
// set up the array that is appended is not counted towards the
// time for naive concatenation.
func benchmarkStrConcat(b *testing.B, numConcat int) {
	// Reports memory allocations
	b.ReportAllocs()

	var ns string
	for i := 0; i < b.N; i++ {
		next := nextNumString()
		ns = ""
		for u := 0; u < numConcat; u++ {
			ns += next()
		}
	}
	// we assign to a global variable to make sure compiler
	// or runtime optimizations don't skip over the operations
	// we were benchmarking. This might be unnecessary, but it's
	// safe.
	global = ns
}

func BenchmarkStrConcat10(b *testing.B) {
	benchmarkStrConcat(b, 10)
}
func BenchmarkStrConcat100(b *testing.B) {
	benchmarkStrConcat(b, 100)
}
func BenchmarkStrConcat500(b *testing.B) {
	benchmarkStrConcat(b, 500)
}
func BenchmarkStrConcat1000(b *testing.B) {
	benchmarkStrConcat(b, 1000)
}
func BenchmarkStrConcat5000(b *testing.B) {
	benchmarkStrConcat(b, 5000)
}

// benchmarkByteSlice provides a benchmark for the time it takes
// to repeatedly append returned strings to a byte slice, and
// finally casting the byte slice to string type.
func benchmarkByteSliceSize(b *testing.B, numConcat int) {
	// Reports memory allocations
	b.ReportAllocs()

	var ns string
	for i := 0; i < b.N; i++ {
		next := nextNumString()
		b := make([]byte, 0, numConcat*10)
		for u := 0; u < numConcat; u++ {
			b = append(b, next()...)
		}
		ns = string(b)
	}
	global = ns
}

func BenchmarkByteSliceSize10(b *testing.B) {
	benchmarkByteSliceSize(b, 10)
}
func BenchmarkByteSliceSize100(b *testing.B) {
	benchmarkByteSliceSize(b, 100)
}
func BenchmarkByteSliceSize500(b *testing.B) {
	benchmarkByteSliceSize(b, 500)
}
func BenchmarkByteSliceSize1000(b *testing.B) {
	benchmarkByteSliceSize(b, 1000)
}
func BenchmarkByteSliceSize5000(b *testing.B) {
	benchmarkByteSliceSize(b, 5000)
}

// benchmarkJoinSize provides a benchmark for the time it takes to set
// up an array with strings, and calling strings.Join on that array
// to get a fully concatenated string – when the (approximate) number of
// strings is known in advance.
//
// This is identical to benchmarkJoin, except numConcat is used to size
// the []string slice's initial capacity to avoid needless reallocation.
func benchmarkJoinSize(b *testing.B, numConcat int) {
	// Reports memory allocations
	b.ReportAllocs()

	var ns string
	for i := 0; i < b.N; i++ {
		next := nextNumString()
		a := make([]string, 0, numConcat)
		for u := 0; u < numConcat; u++ {
			a = append(a, next())
		}
		ns = strings.Join(a, "")
	}
	global = ns
}

func BenchmarkJoinSize10(b *testing.B) {
	benchmarkJoinSize(b, 10)
}
func BenchmarkJoinSize100(b *testing.B) {
	benchmarkJoinSize(b, 100)
}
func BenchmarkJoinSize500(b *testing.B) {
	benchmarkJoinSize(b, 500)
}
func BenchmarkJoinSize1000(b *testing.B) {
	benchmarkJoinSize(b, 1000)
}
func BenchmarkJoinSize5000(b *testing.B) {
	benchmarkJoinSize(b, 5000)
}

func benchmarkBufferSize(b *testing.B, numConcat int) {
	// Reports memory allocations
	b.ReportAllocs()

	var ns string
	for i := 0; i < b.N; i++ {
		next := nextNumString()
		buffer := bytes.NewBuffer(make([]byte, 0, numConcat*10))
		for u := 0; u < numConcat; u++ {
			buffer.WriteString(next())
		}
		ns = buffer.String()
	}
	global = ns
}

func BenchmarkBufferSize10(b *testing.B) {
	benchmarkBufferSize(b, 10)
}
func BenchmarkBufferSize100(b *testing.B) {
	benchmarkBufferSize(b, 100)
}
func BenchmarkBufferSize500(b *testing.B) {
	benchmarkBufferSize(b, 500)
}
func BenchmarkBufferSize1000(b *testing.B) {
	benchmarkBufferSize(b, 1000)
}
func BenchmarkBufferSize5000(b *testing.B) {
	benchmarkBufferSize(b, 5000)
}
