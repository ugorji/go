// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// This file sets up the variables used, including testInitFns.
// Each file should add initialization that should be performed
// after flags are parsed.
//
// init is a multi-step process:
//   - setup vars (handled by init functions in each file)
//   - parse flags
//   - setup derived vars (handled by pre-init registered functions - registered in init function)
//   - post init (handled by post-init registered functions - registered in init function)
// This way, no one has to manage carefully control the initialization
// using file names, etc.
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
	"bytes"
	"flag"
	"io"
	"log"
	"sync"
)

type ioReaderWrapper struct {
	r io.Reader
}

func (x ioReaderWrapper) Read(p []byte) (n int, err error) {
	return x.r.Read(p)
}

type ioWriterWrapper struct {
	w io.Writer
}

func (x ioWriterWrapper) Write(p []byte) (n int, err error) {
	return x.w.Write(p)
}

var (
	testPreInitFns  []func()
	testPostInitFns []func()
	testReInitFns   []func()

	testOnce sync.Once
)

// flag variables used by tests (and bench)
var (
	testVerbose bool

	//depth of 0 maps to ~400bytes json-encoded string, 1 maps to ~1400 bytes, etc
	//For depth>1, we likely trigger stack growth for encoders, making benchmarking unreliable.
	testDepth int

	testMaxInitLen int

	testUseReset    bool
	testUseParallel bool

	testSkipIntf bool

	testUseIoEncDec  int
	testUseIoWrapper bool

	testNumRepeatString int

	testRpcBufsize       int
	testMapStringKeyOnly bool

	testBenchmarkNoConfig bool

	testBenchmarkWithRuntimeMetrics bool
)

func init() {
	log.SetOutput(io.Discard) // don't allow things log to standard out/err
	testInitFlags()
	benchInitFlags()
}

func testInitFlags() {
	var bIgnore bool
	// delete(testDecOpts.ExtFuncs, timeTyp)
	flag.BoolVar(&testVerbose, "tv", false, "Text Extra Verbose Logging if -v if set")
	flag.IntVar(&testUseIoEncDec, "ti", -1, "Use IO Reader/Writer for Marshal/Unmarshal ie >= 0")
	flag.BoolVar(&testUseIoWrapper, "tiw", false, "Wrap the IO Reader/Writer with a base pass-through reader/writer")

	flag.BoolVar(&testSkipIntf, "tf", false, "Skip Interfaces")
	flag.BoolVar(&testUseReset, "tr", false, "Use Reset")
	flag.BoolVar(&testUseParallel, "tp", false, "Run tests in parallel")
	flag.IntVar(&testNumRepeatString, "trs", 8, "Create string variables by repeating a string N times")
	flag.BoolVar(&bIgnore, "tm", true, "(Deprecated) Use Must(En|De)code")

	flag.IntVar(&testMaxInitLen, "tx", 0, "Max Init Len")

	flag.IntVar(&testDepth, "tsd", 0, "Test Struc Depth")
	flag.BoolVar(&testMapStringKeyOnly, "tsk", false, "use maps with string keys only")
}

func benchInitFlags() {
	flag.BoolVar(&testBenchmarkNoConfig, "bnc", false, "benchmarks: do not make configuration changes for fair benchmarking")
	flag.BoolVar(&testBenchmarkWithRuntimeMetrics, "brm", false, "benchmarks: include runtime metrics")
	// flags reproduced here for compatibility (duplicate some in testInitFlags)
	flag.BoolVar(&testMapStringKeyOnly, "bs", false, "benchmarks: use maps with string keys only")
	flag.IntVar(&testDepth, "bd", 1, "Benchmarks: Test Struc Depth")
}

func testReinit() {
	testOnce = sync.Once{}
	for _, f := range testReInitFns {
		f()
	}
}

func testInitAll() {
	// only parse it once.
	if !flag.Parsed() {
		flag.Parse()
	}
	for _, f := range testPreInitFns {
		f()
	}
	for _, f := range testPostInitFns {
		f()
	}
}

// --- functions below are used only by benchmarks alone

func fnBenchmarkByteBuf(bsIn []byte) (buf *bytes.Buffer) {
	// var buf bytes.Buffer
	// buf.Grow(approxSize)
	buf = bytes.NewBuffer(bsIn)
	buf.Truncate(0)
	return
}

// // --- functions below are used by both benchmarks and tests

// // log message only when testVerbose = true (ie go test ... -- -tv).
// //
// // These are for intormational messages that do not necessarily
// // help with diagnosing a failure, or which are too large.
// func logTv(x interface{}, format string, args ...interface{}) {
// 	if !testVerbose {
// 		return
// 	}
// 	if t, ok := x.(testing.TB); ok { // only available from go 1.9
// 		t.Helper()
// 	}
// 	logT(x, format, args...)
// }

// // logT logs messages when running as go test -v
// //
// // Use it for diagnostics messages that help diagnost failure,
// // and when the output is not too long ie shorter than like 100 characters.
// //
// // In general, any logT followed by failT should call this.
// func logT(x interface{}, format string, args ...interface{}) {
// 	if x == nil {
// 		if len(format) == 0 || format[len(format)-1] != '\n' {
// 			format = format + "\n"
// 		}
// 		fmt.Printf(format, args...)
// 		return
// 	}
// 	if t, ok := x.(testing.TB); ok { // only available from go 1.9
// 		t.Helper()
// 		t.Logf(format, args...)
// 	}
// }

// func failTv(x testing.TB, args ...interface{}) {
// 	x.Helper()
// 	if testVerbose {
// 		failTMsg(x, args...)
// 	}
// 	x.FailNow()
// }

// func failT(x testing.TB, args ...interface{}) {
// 	x.Helper()
// 	failTMsg(x, args...)
// 	x.FailNow()
// }

// func failTMsg(x testing.TB, args ...interface{}) {
// 	x.Helper()
// 	if len(args) > 0 {
// 		if format, ok := args[0].(string); ok {
// 			logT(x, format, args[1:]...)
// 		} else if len(args) == 1 {
// 			logT(x, "%v", args[0])
// 		} else {
// 			logT(x, "%v", args)
// 		}
// 	}
// }

// --- functions below are used only by benchmarks alone

// func benchFnCodecEncode(ts interface{}, bsIn []byte, h Handle) (bs []byte, err error) {
// 	return testCodecEncode(ts, bsIn, fnBenchmarkByteBuf, h)
// }

// func benchFnCodecDecode(bs []byte, ts interface{}, h Handle) (err error) {
// 	return testCodecDecode(bs, ts, h)
// }
