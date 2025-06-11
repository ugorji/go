// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// This file sets up the variables used, including testInitFns.
// Each file should add initialization that should be performed
// after flags are parsed.
//
// init is a multi-step process:
//   - pre-init:
//     - each file will fully define and pre-configure their variables
//       (open to everyone in the package)
//     - preloaded with configure test/benchmark flags and parsing cmdline flags.
//   - post-init:
//     - configure variables which depend on fully configured state from other files
//
// Each file will append init functions of the form func() to the
// pre-init, post-init or re-init slices.
//
// Suite tests may require each of the pre-defined handles to modify
// their configuration and reinit. This is what the testReinit function handles.
//
// This way, no one has to manage carefully control the initialization
// using file names, etc.
//
// Note: during preInit: update testEncodeOptions and testDecodeOptions
// with values from the flags. Each handle will be updated by basically
// copying these into the EncodeOptions and DecodeOptions of the Handle's BasicHandle
// embedded value.
//
// TestMain will call the full init (via testInitAll) as a one-time initialization
// for the full test run. Consequently, no need to use sync.Once to manage that.
//
// Tests which require external dependencies need the -tag=x parameter.
// They should be run as:
//    go test -tags=x -run=. <other parameters ...>
// Benchmarks should also take this parameter, to include the sereal, xdr, etc.
// To run against codecgen, etc, make sure you pass extra parameters.
// Example usage:
//    go test "-tags=x codecgen" -bench=. <other parameters ...>
//
// To fully test everything:
//    go test -tags=x -benchtime=100ms -tv -bg -bi  -brw -bu -v -run=. -bench=.
//
// Handling flags
// define a set of global flags for testing, including:
//   - Use Reset
//   - Use IO reader/writer (vs direct bytes)
//   - Set Canonical
//   - Set InternStrings
//   - Use Symbols
//
// This way, we can test them all by running same set of tests with a different
// set of flags.
//
// Following this, all the benchmarks will utilize flags set by codec_test.go
// and will not redefine these "global" flags.

import (
	"flag"
	"reflect"
	"strconv"
	"testing"
)

func init() {
	// log.SetOutput(io.Discard) // don't allow things log to standard out/err
	testPreInitFns = append(testPreInitFns, testInitFlags, benchInitFlags, testParseFlags)
	// testPostInitFns = append(testPostInitFns, testUpdateOptionsFromFlags)
	// testReInitFns = append(testReInitFns, testUpdateOptionsFromFlags)
}

func TestMain(m *testing.M) {
	testInitAll() // with this, we can remove testOnce.Do(testInitAll) everywhere
	m.Run()
}

var (
	testPreInitFns  []func()
	testPostInitFns []func()
	testReInitFns   []func()
	// testOnce sync.Once
)

// flag variables used by tests (and bench)
type testVars struct {
	Verbose bool
	//depth of 0 maps to ~400bytes json-encoded string, 1 maps to ~1400 bytes, etc
	//For depth>1, we likely trigger stack growth for encoders, making benchmarking unreliable.
	Depth int

	UseDiff bool

	UseReset    bool
	UseParallel bool

	SkipIntf     bool
	SkipRPCTests bool

	UseIoWrapper bool

	NumRepeatString int

	RpcBufsize       int
	MapStringKeyOnly bool

	BenchmarkNoConfig bool

	BenchmarkWithRuntimeMetrics bool

	bufsize    testBufioSizeFlag
	maxInitLen int
	zeroCopy   bool

	// MaxInitLen int
	// ZeroCopy         bool
	// UseIoEncDec  int
}

var testv = testVars{
	bufsize:         -1,
	maxInitLen:      1024,
	zeroCopy:        true,
	NumRepeatString: 8,
}

type testBufioSizeFlag int

func (x *testBufioSizeFlag) String() string { return strconv.Itoa(int(*x)) }
func (x *testBufioSizeFlag) Set(s string) (err error) {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		v = -1
	}
	*x = testBufioSizeFlag(v)
	return
}
func (x *testBufioSizeFlag) Get() interface{} { return int(*x) }

func testInitFlags() {
	var bIgnore bool
	// delete(testDecOpts.ExtFuncs, timeTyp)
	flag.Var(&testv.bufsize, "ti", "Use IO Reader/Writer for Marshal/Unmarshal ie >= 0")
	flag.BoolVar(&testv.Verbose, "tv", false, "Text Extra Verbose Logging if -v if set")
	flag.BoolVar(&testv.UseIoWrapper, "tiw", false, "Wrap the IO Reader/Writer with a base pass-through reader/writer")

	flag.BoolVar(&testv.SkipIntf, "tf", false, "Skip Interfaces")
	flag.BoolVar(&testv.UseReset, "tr", false, "Use Reset")
	flag.BoolVar(&testv.UseParallel, "tp", false, "Run tests in parallel")
	flag.IntVar(&testv.NumRepeatString, "trs", 8, "Create string variables by repeating a string N times")
	flag.BoolVar(&testv.UseDiff, "tdiff", false, "Use Diff")
	flag.BoolVar(&testv.zeroCopy, "tzc", false, "Use Zero copy mode")

	flag.IntVar(&testv.maxInitLen, "tx", 0, "Max Init Len")

	flag.IntVar(&testv.Depth, "tsd", 0, "Test Struc Depth")
	flag.BoolVar(&testv.MapStringKeyOnly, "tsk", false, "use maps with string keys only")

	flag.BoolVar(&bIgnore, "tm", false, "(Deprecated) Use Must(En|De)code")
}

func benchInitFlags() {
	flag.BoolVar(&testv.BenchmarkNoConfig, "bnc", false, "benchmarks: do not make configuration changes for fair benchmarking")
	flag.BoolVar(&testv.BenchmarkWithRuntimeMetrics, "brm", false, "benchmarks: include runtime metrics")
	// flags reproduced here for compatibility (duplicate some in testInitFlags)
	flag.BoolVar(&testv.MapStringKeyOnly, "bs", false, "benchmarks: use maps with string keys only")
	flag.IntVar(&testv.Depth, "bd", 1, "Benchmarks: Test Struc Depth")
}

func testParseFlags() {
	// only parse it once.
	if !flag.Parsed() {
		flag.Parse()
	}
}

func testReinit() {
	// testOnce = sync.Once{}
	for _, f := range testReInitFns {
		f()
	}
}

func testInitAll() {
	for _, f := range testPreInitFns {
		f()
	}
	for _, f := range testPostInitFns {
		f()
	}
}

func approxDataSize(rv reflect.Value) (sum int) {
	switch rk := rv.Kind(); rk {
	case reflect.Invalid:
	case reflect.Ptr, reflect.Interface:
		sum += int(rv.Type().Size())
		sum += approxDataSize(rv.Elem())
	case reflect.Slice:
		sum += int(rv.Type().Size())
		for j := 0; j < rv.Len(); j++ {
			sum += approxDataSize(rv.Index(j))
		}
	case reflect.String:
		sum += int(rv.Type().Size())
		sum += rv.Len()
	case reflect.Map:
		sum += int(rv.Type().Size())
		for _, mk := range rv.MapKeys() {
			sum += approxDataSize(mk)
			sum += approxDataSize(rv.MapIndex(mk))
		}
	case reflect.Struct:
		//struct size already includes the full data size.
		//sum += int(rv.Type().Size())
		for j := 0; j < rv.NumField(); j++ {
			sum += approxDataSize(rv.Field(j))
		}
	default:
		//pure value types
		sum += int(rv.Type().Size())
	}
	return
}
