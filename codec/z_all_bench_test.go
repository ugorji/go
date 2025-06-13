//go:build (alltests || codec.alltests) && !codec.nobench && !nobench && go1.24

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// see notes in z_all_test.go

import (
	"fmt"
	"testing"
	"time"
)

func init() {
	// testPostInitFns = append(testPostInitFns, benchmarkGroupInitAll)
}

// func benchmarkGroupInitAll() {}

func benchmarkOneFn(fns []func(*testing.B)) func(*testing.B) {
	switch len(fns) {
	case 0:
		return nil
	case 1:
		return fns[0]
	default:
		return func(t *testing.B) {
			for _, f := range fns {
				f(t)
			}
		}
	}
}

func benchmarkSuiteNoop(b *testing.B) {
	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	for b.Loop() {
		time.Sleep(10 * time.Millisecond)
	}
}

func benchmarkSuite(t *testing.B, fns ...func(t *testing.B)) {
	defer tbvars.setBufsize((int)(testv.bufsize))

	f := benchmarkOneFn(fns)
	// find . -name "*_test.go" | xargs grep -e 'flag.' | cut -d '&' -f 2 | cut -d ',' -f 1 | grep -e '^bench'

	testReinit()
	tbvars.setBufsize(-1)
	testReinit()
	t.Run("use-bytes.......", f)

	tbvars.setBufsize(1024)
	testReinit()
	t.Run("use-io-1024-....", f)

	tbvars.setBufsize(0)
	testReinit()
	t.Run("use-io-0-.......", f)
}

func benchmarkVeryQuickSuite(t *testing.B, name string, fns ...func(t *testing.B)) {
	defer tbvars.setBufsize((int)(testv.bufsize))
	benchmarkDivider()

	tbvars.setBufsize(-1)
	testReinit()

	t.Run(fmt.Sprintf("%s-bd%d........", name, testv.Depth), benchmarkOneFn(fns))
}

func benchmarkQuickSuite(t *testing.B, name string, fns ...func(t *testing.B)) {
	defer tbvars.setBufsize((int)(testv.bufsize))
	benchmarkVeryQuickSuite(t, name, fns...)

	// encoded size of TestStruc is between 20K and 30K for bd=1 // consider buffer=1024 * 16 * testv.Depth
	tbvars.setBufsize(1024) // (value of byteBufSize): use smaller buffer, and more flushes - it's ok.
	testReinit()
	t.Run(fmt.Sprintf("%s-bd%d-buf%d", name, testv.Depth, 1024), benchmarkOneFn(fns))

	tbvars.setBufsize(0)
	testReinit()
	t.Run(fmt.Sprintf("%s-bd%d-io.....", name, testv.Depth), benchmarkOneFn(fns))
}

/*
z='bench_test.go'
find . -name "$z" | xargs grep -e '^func Benchmark.*Encode' | \
    cut -d '(' -f 1 | cut -d ' ' -f 2 | \
    while read f; do echo "t.Run(\"$f\", $f)"; done &&
echo &&
find . -name "$z" | xargs grep -e '^func Benchmark.*Decode' | \
    cut -d '(' -f 1 | cut -d ' ' -f 2 | \
    while read f; do echo "t.Run(\"$f\", $f)"; done
*/

func benchmarkCodecGroup(t *testing.B) {
	benchmarkDivider()
	t.Run("Benchmark__Msgpack____Encode", Benchmark__Msgpack____Encode)
	t.Run("Benchmark__Binc_______Encode", Benchmark__Binc_______Encode)
	t.Run("Benchmark__Simple_____Encode", Benchmark__Simple_____Encode)
	t.Run("Benchmark__Cbor_______Encode", Benchmark__Cbor_______Encode)
	t.Run("Benchmark__Json_______Encode", Benchmark__Json_______Encode)
	benchmarkDivider()
	t.Run("Benchmark__Msgpack____Decode", Benchmark__Msgpack____Decode)
	t.Run("Benchmark__Binc_______Decode", Benchmark__Binc_______Decode)
	t.Run("Benchmark__Simple_____Decode", Benchmark__Simple_____Decode)
	t.Run("Benchmark__Cbor_______Decode", Benchmark__Cbor_______Decode)
	t.Run("Benchmark__Json_______Decode", Benchmark__Json_______Decode)
}

func BenchmarkCodecSuite(t *testing.B) { benchmarkSuite(t, benchmarkCodecGroup) }

func benchmarkJsonEncodeGroup(t *testing.B) {
	t.Run("Benchmark__Json_______Encode", Benchmark__Json_______Encode)
}

func benchmarkJsonDecodeGroup(t *testing.B) {
	t.Run("Benchmark__Json_______Decode", Benchmark__Json_______Decode)
}

func benchmarkCborEncodeGroup(t *testing.B) {
	t.Run("Benchmark__Cbor_______Encode", Benchmark__Cbor_______Encode)
}

func benchmarkCborDecodeGroup(t *testing.B) {
	t.Run("Benchmark__Cbor_______Decode", Benchmark__Cbor_______Decode)
}

func BenchmarkCodecQuickSuite(t *testing.B) {
	benchmarkQuickSuite(t, "cbor", benchmarkCborEncodeGroup)
	benchmarkQuickSuite(t, "cbor", benchmarkCborDecodeGroup)
	benchmarkQuickSuite(t, "json", benchmarkJsonEncodeGroup)
	benchmarkQuickSuite(t, "json", benchmarkJsonDecodeGroup)

	// depths := [...]int{1, 4}
	// for _, d := range depths {
	// 	benchmarkQuickSuite(t, d, benchmarkJsonEncodeGroup)
	// 	benchmarkQuickSuite(t, d, benchmarkJsonDecodeGroup)
	// }

	// benchmarkQuickSuite(t, 1, benchmarkJsonEncodeGroup)
	// benchmarkQuickSuite(t, 4, benchmarkJsonEncodeGroup)
	// benchmarkQuickSuite(t, 1, benchmarkJsonDecodeGroup)
	// benchmarkQuickSuite(t, 4, benchmarkJsonDecodeGroup)

	// benchmarkQuickSuite(t, 1, benchmarkJsonEncodeGroup, benchmarkJsonDecodeGroup)
	// benchmarkQuickSuite(t, 4, benchmarkJsonEncodeGroup, benchmarkJsonDecodeGroup)
	// benchmarkQuickSuite(t, benchmarkJsonEncodeGroup)
	// benchmarkQuickSuite(t, benchmarkJsonDecodeGroup)
}
