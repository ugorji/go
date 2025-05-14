//go:build !codec.nobench && go1.24

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// bench_test is the "helper" file for all benchmarking tests.
//
// There are also benchmarks which depend on just codec and the stdlib,
// and benchmarks which depend on external libraries.
// It is an explicit goal that you can run benchmarks without external
// dependencies (which is why the 'x' build tag was explicitly introduced).
//
// There are 2 ways of running tests:
//    - generated
//    - not generated
//
// Consequently, we have 4 groups:
//    - codec_bench   (gen, !gen)
//    - stdlib_bench  (!gen only)
//    - x_bench       (!gen only)
//    - x_bench_gen   (gen only)
//
// We also have 4 matching suite files.
//    - z_all_bench (rename later to z_all_codec_bench???)
//    - z_all_stdlib_bench
//    - z_all_x_bench
//    - z_all_x_bench_gen
//
// Finally, we have a single test (TestBenchInit) that
// will log information about whether each format can
// encode or not, how long to encode (unscientifically),
// and the encode size.
//
// This test MUST be run always, as it calls init() internally

import (
	"reflect"
	"runtime"
	"runtime/metrics"
	"strings"
	"testing"
	"time"
)

// Sample way to run:
// go test -bi -bv -bd=1 -benchmem -bench=.

const (
	benchUnscientificRes = true
	benchVerify          = true
	benchRecover         = true
	benchShowJsonOnError = true
)

var (
	benchTs       *TestStruc
	approxSize    int
	benchCheckers []benchChecker
)

type benchEncFn func(interface{}, []byte) ([]byte, error)
type benchDecFn func([]byte, interface{}) error
type benchIntfFn func() interface{}

type benchChecker struct {
	name     string
	encodefn benchEncFn
	decodefn benchDecFn
}

func init() {
	// testPreInitFns = append(testPreInitFns, benchPreInit)
	// testPostInitFns = append(testPostInitFns, codecbenchPostInit)
	testPostInitFns = append(testPostInitFns, benchInit)
	testReInitFns = append(testReInitFns, benchReinit)
}

func benchInit() {
	benchTs = newTestStruc(testDepth, testNumRepeatString, true, !testSkipIntf, testMapStringKeyOnly)
	approxSize = approxDataSize(reflect.ValueOf(benchTs)) * 2 // multiply by 1.5 or 2 to appease msgp, and prevent alloc
	// bytesLen := 1024 * 4 * (testDepth + 1) * (testDepth + 1)
	// if bytesLen < approxSize {
	// 	bytesLen = approxSize
	// }
	if !testBenchmarkNoConfig {
		// benchmark comparisons use zerocopy (easyjson, json-iterator, etc).
		// use zerocopy for the benchmarks, for best performance, and better comparison to others.
		testJsonH.ZeroCopy = true
		testCborH.ZeroCopy = true
		testMsgpackH.ZeroCopy = true
		testSimpleH.ZeroCopy = true
		testBincH.ZeroCopy = true

		// std-lib encoging/json does the following, which we will do for codec in these benchmarks.
		// - sets into a map without getting what's there first.
		// - sets into an interface value regardless of what was in there
		// - sets slice to zero len first, then appends (equivalent to ignoring slice contents)

		testJsonH.MapValueReset = true
		testJsonH.InterfaceReset = true
		testJsonH.SliceElementReset = true
	}
}

func benchReinit() {
	benchCheckers = nil
}

func benchmarkDivider() {
	// logTv(nil, "-------------------------------\n")
	println()
}

// func Test0(t *testing.T) {
// 	testOnce.Do(testInitAll)
// }

func TestBenchInit(t *testing.T) {
	// testOnce.Do(testInitAll)
	if !testing.Verbose() {
		return
	}
	// t.Logf("..............................................")
	t.Logf("BENCHMARK INIT: %v", time.Now())
	// t.Logf("To run full benchmark comparing encodings, use: \"go test -bench=.\"")
	t.Logf("Benchmark: ")
	t.Logf("\tStruct recursive Depth:             %d", testDepth)
	if approxSize > 0 {
		t.Logf("\tApproxDeepSize Of benchmark Struct: %d bytes", approxSize)
	}
	if benchUnscientificRes {
		t.Logf("Benchmark One-Pass Run (with Unscientific Encode/Decode times): ")
	} else {
		t.Logf("Benchmark One-Pass Run:")
	}
	for _, bc := range benchCheckers {
		doBenchCheck(t, bc.name, bc.encodefn, bc.decodefn)
	}
	if testVerbose {
		t.Logf("..............................................")
		t.Logf("<<<<====>>>> depth: %v, ts: %#v\n", testDepth, benchTs)
	}
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
}

var vBenchTs = TestStruc{}

func fnBenchNewTs() interface{} {
	vBenchTs = TestStruc{}
	return &vBenchTs
	// return new(TestStruc)
}

// const benchCheckDoDeepEqual = false

func benchRecoverPanic(t *testing.B) {
	if benchRecover {
		if r := recover(); r != nil {
			t.Logf("(recovered) panic: %v\n", r)
			t.FailNow()
		}
	}
}

func benchRecoverPanicT(t *testing.T) {
	if benchRecover {
		if r := recover(); r != nil {
			t.Logf("(recovered) panic: %v\n", r)
			t.FailNow()
		}
	}
}

func doBenchCheck(t *testing.T, name string, encfn benchEncFn, decfn benchDecFn) {
	// if benchUnscientificRes {
	// 	t.Logf("-------------- %s ----------------", name)
	// }
	defer benchRecoverPanicT(t)
	runtime.GC()
	tnow := time.Now()
	buf, err := encfn(benchTs, nil)
	if err != nil {
		t.Logf("\t%10s: **** Error encoding benchTs: %v", name, err)
		return
	}
	encDur := time.Since(tnow)
	encLen := len(buf)
	runtime.GC()
	if !benchUnscientificRes {
		t.Logf("\t%10s: len: %d bytes\n", name, encLen)
		return
	}
	tnow = time.Now()
	var ts2 TestStruc
	if err = decfn(buf, &ts2); err != nil {
		t.Logf("\t%10s: **** Error decoding into new TestStruc: %v", name, err)
		return
	}
	decDur := time.Since(tnow)
	// if benchCheckDoDeepEqual {
	if benchVerify {
		err = deepEqual(benchTs, &ts2)
		if err == nil {
			t.Logf("\t%10s: len: %d bytes,\t encode: %v,\t decode: %v,\tencoded == decoded", name, encLen, encDur, decDur)
		} else {
			t.Logf("\t%10s: len: %d bytes,\t encode: %v,\t decode: %v,\tencoded != decoded: %v", name, encLen, encDur, decDur, err)
		}
	} else {
		t.Logf("\t%10s: len: %d bytes,\t encode: %v,\t decode: %v", name, encLen, encDur, decDur)
	}
}

func fnBenchmarkEncode(b *testing.B, encName string, ts interface{}, encfn benchEncFn) {
	defer benchRecoverPanic(b)
	// testOnce.Do(testInitAll)
	// ignore method params: ts, and work on benchTs directly
	ts = benchTs
	// do initial warm up by running encode one time
	bs, err := encfn(ts, make([]byte, 0, approxSize))
	// var err error
	// bs := make([]byte, 0, approxSize)
	fnRun := func() {
		if _, err = encfn(ts, bs); err != nil {
			b.Logf("Error encoding benchTs: %s: %v", encName, err)
			b.FailNow()
		}
	}
	fnRun()
	fnBenchmarkRun(b, fnRun)
}

func fnBenchmarkDecode(b *testing.B, encName string, ts interface{},
	encfn benchEncFn, decfn benchDecFn, newfn benchIntfFn,
) {
	defer benchRecoverPanic(b)
	// testOnce.Do(testInitAll)

	// MARKER: to ensure same sequence of bytes to be decoded, always encode using codec encoder.
	//
	// ignore method params:
	// - ts: use benchTs instead
	// - newfn: use TestStruc instead
	// - benchEncFn: use codec's encfn instead (based on name/format)

	ts = benchTs

	if strings.Contains(encName, "json") {
		encfn = fnJsonEncodeFn
	} else if strings.Contains(encName, "cbor") {
		encfn = fnCborEncodeFn
	} else if strings.Contains(encName, "msgpack") {
		encfn = fnMsgpackEncodeFn
	}

	buf := make([]byte, 0, approxSize)
	buf, err := encfn(ts, buf)
	if err != nil {
		b.Logf("Error encoding benchTs: %s: %v", encName, err)
		b.FailNow()
	}

	// do initial warm up by running decode one time
	locTs := new(TestStruc)
	ts = locTs

	fnRun := func() {
		*locTs = TestStruc{}
		if err = decfn(buf, ts); err != nil {
			b.Logf("Error decoding into new TestStruc: %s: %v", encName, err)
			b.FailNow()
		}
	}

	fnRun()
	fnBenchmarkRun(b, fnRun)

	// if false && benchVerify { // do not do benchVerify during decode
	// 	// ts2 := newfn()
	// 	ts1 := ts.(*TestStruc)
	// 	ts2 := new(TestStruc)
	// 	if err = decfn(buf, ts2); err != nil {
	// 		failT(b, "BenchVerify: Error decoding benchTs: %s: %v", encName, err)
	// 	}
	// 	if err = deepEqual(ts1, ts2); err != nil {
	// 		failT(b, "BenchVerify: Error comparing benchTs: %s: %v", encName, err)
	// 	}
	// }
}

func fnBenchmarkRun(b *testing.B, fn func()) {
	fn() // run one time first - to init things
	if testBenchmarkWithRuntimeMetrics {
		fnBenchmarkRunWithMetrics(b, fn)
	} else {
		fnBenchmarkRunNoMetrics(b, fn)
	}
}

func fnBenchmarkRunNoMetrics(b *testing.B, fn func()) {
	runtime.GC()
	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	for b.Loop() {
		fn()
	}
}

func fnBenchmarkRunWithMetrics(b *testing.B, fn func()) {
	var names = [...]string{
		`/gc/cycles/automatic:gc-cycles`,
		`/gc/scan/heap:bytes`,
		`/cpu/classes/gc/total:cpu-seconds`,
		//`/cpu/classes/total:cpu-seconds`,
		`/cpu/classes/idle:cpu-seconds`,
		`/cpu/classes/user:cpu-seconds`,
	}
	var cols = [...]string{
		`gcRuns`,
		`gcScanBytes`,
		`gcCpuSec`,
		//`userRtCpuSec`,
		`idleCpuSec`,
		`userCpuSec`,
	}
	var s1, s2 [len(names)]metrics.Sample
	for i, s := range names {
		s1[i].Name = s
		s2[i].Name = s
	}

	fnRM := func(i int) {
		var fv float64
		if strings.HasSuffix(cols[i], "CpuSec") {
			fv = s2[i].Value.Float64() - s1[i].Value.Float64()
		} else {
			i1, i2 := s1[i].Value.Uint64(), s2[i].Value.Uint64()
			// println("fnBenchmarkRunWithMetrics: i1: %d, i2: %d", i1, i2)
			if i2 >= i1 {
				fv = float64(i2 - i1)
			} else {
				fv = -float64(i1 - i2)
			}
		}
		b.ReportMetric(fv, cols[i])
	}

	runtime.GC()
	metrics.Read(s1[:])
	// runtime.ReadMemStats(&m0)

	for b.Loop() {
		fn()
	}

	// runtime.ReadMemStats(&m1)
	// b.ReportMetric(float64(m1.NumGC-m0.NumGC), "gcRuns")
	// b.ReportMetric(float64(m1.Lookups-m0.Lookups), "ptrLookups")
	// // b.ReportMetric(float64(gcRuns)/float64(b.N), "gc_runs/op")

	runtime.GC()
	metrics.Read(s2[:])
	for i := range len(names) {
		fnRM(i)
	}
}
