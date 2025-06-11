// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build (alltests || codec.alltests) && go1.9

package codec

// Run this using:
//   go test -tags=alltests -run=Suite -coverprofile=cov.out
//   go tool cover -html=cov.out
//
// Because build tags are a build time parameter, we will have to test out the
// different tags separately.
// Tags: x codecgen codec.safe codec.notfastpath appengine
//
// These tags should be added to alltests, e.g.
//   go test '-tags=alltests x codecgen' -run=Suite -coverprofile=cov.out
//
// To run all tests before submitting code, run:
//    a=( "" "codec.safe" "codecgen" "codec.notfastpath" "codecgen codec.notfastpath" "codecgen codec.safe" "codec.safe codec.notfastpath" )
//    for i in "${a[@]}"; do echo ">>>> TAGS: $i"; go test "-tags=alltests $i" -run=Suite; done
//
// This suite of tests requires support for subtests and suites,
// and conseqneutly will only run on go1.7 and above.

// Easy script to re-generate the calls below:
// -------------------------------------------
/*
# Old model when all tests were in codec_test.go
find . -name "codec_test.go" | xargs grep -e '^func Test' | \
    cut -d '(' -f 1 | cut -d ' ' -f 2 | \
    while read f; do echo "t.Run(\"$f\", $f)"; done

# New model when all tests are in format specific test files e.g. json_test.go
for a in binc cbor json msgpack simple; do
  echo "============ $a ============"
  grep "func Test" ${a}_test.go  | while IFS= read -r line; do
    if [[ "$line" =~ func\ ([A-Za-z0-9_]+).* ]]; then
      a="${BASH_REMATCH[1]}"
      echo -e "\tt.Run(\"${a}\", ${a})"
    fi
  done
done

# For breaking up the testJsonGroup vs testJsonGroupV style:
for a in binc cbor json msgpack simple; do
  echo "============ $a ============"
  grep "func Test" ${a}_test.go | grep -v -E '(CodecsTable|CodecsMisc|SwallowAndZero|NextValueBytes|StrucEncDec)' | while IFS= read -r line; do
    if [[ "$line" =~ func\ ([A-Za-z0-9_]+).* ]]; then
      a="${BASH_REMATCH[1]}"
      echo -e "\tt.Run(\"${a}\", ${a})"
    fi
  done
done

for a in binc cbor json msgpack simple; do
  echo "============ $a ============"
  grep -E 'func Test.*(CodecsTable|CodecsMisc|SwallowAndZero|NextValueBytes|StrucEncDec)' ${a}_test.go | while IFS= read -r line; do
    if [[ "$line" =~ func\ ([A-Za-z0-9_]+).* ]]; then
      a="${BASH_REMATCH[1]}"
      echo -e "\tt.Run(\"${a}\", ${a})"
    fi
  done
done

# testNonHandlesGroup
grep "func Test" codec_test.go | grep -v -E '(Cbor|Json|Simple|Msgpack|Binc)'

*/

import (
	"testing"
	"time"
)

// func TestMain(m *testing.M) {
// 	println("calling TestMain")
// 	// set some parameters
// 	exitcode := m.Run()
// 	os.Exit(exitcode)
// }

type testTimeTracker struct {
	t time.Time
}

func (tt *testTimeTracker) Elapsed() (d time.Duration) {
	if !tt.t.IsZero() {
		d = time.Since(tt.t)
	}
	tt.t = time.Now()
	return
}

func testCodecGroup(t *testing.T) {
	// <setup code>
	testJsonGroup(t)
	testBincGroup(t)
	testCborGroup(t)
	testMsgpackGroup(t)
	testSimpleGroup(t)
	// testSimpleMammothGroup(t)
	// testRpcGroup(t)
	testNonHandlesGroup(t)

	testCodecGroupV(t)
	// <tear-down code>
}

func testCodecGroupV(t *testing.T) {
	testJsonGroupV(t)
	testBincGroupV(t)
	testCborGroupV(t)
	testMsgpackGroupV(t)
	testSimpleGroupV(t)
}

func testJsonGroup(t *testing.T) {
	t.Run("TestJsonCodecsEmbeddedPointer", TestJsonCodecsEmbeddedPointer)
	t.Run("TestJsonCodecChan", TestJsonCodecChan)
	t.Run("TestJsonStdEncIntf", TestJsonStdEncIntf)
	t.Run("TestJsonRaw", TestJsonRaw)
	t.Run("TestJsonRpcGo", TestJsonRpcGo)
	t.Run("TestJsonMapEncodeForCanonical", TestJsonMapEncodeForCanonical)
	t.Run("TestJsonRawExt", TestJsonRawExt)
	t.Run("TestJsonMapStructKey", TestJsonMapStructKey)
	t.Run("TestJsonDecodeNilMapValue", TestJsonDecodeNilMapValue)
	t.Run("TestJsonEmbeddedFieldPrecedence", TestJsonEmbeddedFieldPrecedence)
	t.Run("TestJsonLargeContainerLen", TestJsonLargeContainerLen)
	t.Run("TestJsonTime", TestJsonTime)
	t.Run("TestJsonUintToInt", TestJsonUintToInt)
	t.Run("TestJsonDifferentMapOrSliceType", TestJsonDifferentMapOrSliceType)
	t.Run("TestJsonScalars", TestJsonScalars)
	t.Run("TestJsonOmitempty", TestJsonOmitempty)
	t.Run("TestJsonIntfMapping", TestJsonIntfMapping)
	t.Run("TestJsonMissingFields", TestJsonMissingFields)
	t.Run("TestJsonMaxDepth", TestJsonMaxDepth)
	t.Run("TestJsonSelfExt", TestJsonSelfExt)
	t.Run("TestJsonBytesEncodedAsArray", TestJsonBytesEncodedAsArray)
	t.Run("TestJsonRawToStringToRawEtc", TestJsonRawToStringToRawEtc)
	t.Run("TestJsonStructKeyType", TestJsonStructKeyType)
	t.Run("TestJsonPreferArrayOverSlice", TestJsonPreferArrayOverSlice)
	t.Run("TestJsonZeroCopyBytes", TestJsonZeroCopyBytes)
	t.Run("TestJsonNumbers", TestJsonNumbers)
	t.Run("TestJsonDesc", TestJsonDesc)
	t.Run("TestJsonStructFieldInfoToArray", TestJsonStructFieldInfoToArray)
	t.Run("TestJsonEncodeIndent", TestJsonEncodeIndent)
	t.Run("TestJsonDecodeNonStringScalarInStringContext", TestJsonDecodeNonStringScalarInStringContext)
	t.Run("TestJsonLargeInteger", TestJsonLargeInteger)
	t.Run("TestJsonInvalidUnicode", TestJsonInvalidUnicode)
	t.Run("TestJsonNumberParsing", TestJsonNumberParsing)
	t.Run("TestJsonMultipleEncDec", TestJsonMultipleEncDec)
	t.Run("TestJsonAllErrWriter", TestJsonAllErrWriter)
	t.Run("TestJsonMammoth", TestJsonMammoth)
	t.Run("TestJsonMammothMapsAndSlices", TestJsonMammothMapsAndSlices)

	t.Run("TestJsonTimeAndBytesOptions", TestJsonTimeAndBytesOptions)
}

func testJsonGroupV(t *testing.T) {
	t.Run("TestJsonCodecsTable", TestJsonCodecsTable)
	t.Run("TestJsonCodecsMisc", TestJsonCodecsMisc)
	t.Run("TestJsonSwallowAndZero", TestJsonSwallowAndZero)
	t.Run("TestJsonNextValueBytes", TestJsonNextValueBytes)
	t.Run("TestJsonStrucEncDec", TestJsonStrucEncDec)
	t.Run("TestJsonTimeAndBytesOptions", TestJsonTimeAndBytesOptions)
}

func testBincGroup(t *testing.T) {
	t.Run("TestBincCodecsEmbeddedPointer", TestBincCodecsEmbeddedPointer)
	t.Run("TestBincStdEncIntf", TestBincStdEncIntf)
	t.Run("TestBincRaw", TestBincRaw)
	t.Run("TestBincRpcGo", TestBincRpcGo)
	t.Run("TestBincMapEncodeForCanonical", TestBincMapEncodeForCanonical)
	t.Run("TestBincUnderlyingType", TestBincUnderlyingType)
	t.Run("TestBincRawExt", TestBincRawExt)
	t.Run("TestBincMapStructKey", TestBincMapStructKey)
	t.Run("TestBincDecodeNilMapValue", TestBincDecodeNilMapValue)
	t.Run("TestBincEmbeddedFieldPrecedence", TestBincEmbeddedFieldPrecedence)
	t.Run("TestBincLargeContainerLen", TestBincLargeContainerLen)
	t.Run("TestBincTime", TestBincTime)
	t.Run("TestBincUintToInt", TestBincUintToInt)
	t.Run("TestBincDifferentMapOrSliceType", TestBincDifferentMapOrSliceType)
	t.Run("TestBincScalars", TestBincScalars)
	t.Run("TestBincOmitempty", TestBincOmitempty)
	t.Run("TestBincIntfMapping", TestBincIntfMapping)
	t.Run("TestBincMissingFields", TestBincMissingFields)
	t.Run("TestBincMaxDepth", TestBincMaxDepth)
	t.Run("TestBincSelfExt", TestBincSelfExt)
	t.Run("TestBincBytesEncodedAsArray", TestBincBytesEncodedAsArray)
	t.Run("TestBincRawToStringToRawEtc", TestBincRawToStringToRawEtc)
	t.Run("TestBincStructKeyType", TestBincStructKeyType)
	t.Run("TestBincPreferArrayOverSlice", TestBincPreferArrayOverSlice)
	t.Run("TestBincZeroCopyBytes", TestBincZeroCopyBytes)
	t.Run("TestBincNumbers", TestBincNumbers)
	t.Run("TestBincDesc", TestBincDesc)
	t.Run("TestBincStructFieldInfoToArray", TestBincStructFieldInfoToArray)
	t.Run("TestBincMammoth", TestBincMammoth)
	t.Run("TestBincMammothMapsAndSlices", TestBincMammothMapsAndSlices)
}

func testBincGroupV(t *testing.T) {
	t.Run("TestBincCodecsTable", TestBincCodecsTable)
	t.Run("TestBincCodecsMisc", TestBincCodecsMisc)
	t.Run("TestBincSwallowAndZero", TestBincSwallowAndZero)
	t.Run("TestBincNextValueBytes", TestBincNextValueBytes)
	t.Run("TestBincStrucEncDec", TestBincStrucEncDec)
}

func testCborGroup(t *testing.T) {
	t.Run("TestCborIndefiniteLength", TestCborIndefiniteLength)
	t.Run("TestCborIndefiniteLengthStringChunksCannotMixTypes", TestCborIndefiniteLengthStringChunksCannotMixTypes)
	t.Run("TestCborIndefiniteLengthTextStringChunksAreUTF8", TestCborIndefiniteLengthTextStringChunksAreUTF8)
	t.Run("TestCborGoldens", TestCborGoldens)
	t.Run("TestCborHalfFloat", TestCborHalfFloat)
	t.Run("TestCborSkipTags", TestCborSkipTags)
	t.Run("TestCborMalformed", TestCborMalformed)
	t.Run("TestCborCodecsEmbeddedPointer", TestCborCodecsEmbeddedPointer)
	t.Run("TestCborCodecChan", TestCborCodecChan)
	t.Run("TestCborStdEncIntf", TestCborStdEncIntf)
	t.Run("TestCborRaw", TestCborRaw)
	t.Run("TestCborRpcGo", TestCborRpcGo)
	t.Run("TestCborMapEncodeForCanonical", TestCborMapEncodeForCanonical)
	t.Run("TestCborRawExt", TestCborRawExt)
	t.Run("TestCborMapStructKey", TestCborMapStructKey)
	t.Run("TestCborDecodeNilMapValue", TestCborDecodeNilMapValue)
	t.Run("TestCborEmbeddedFieldPrecedence", TestCborEmbeddedFieldPrecedence)
	t.Run("TestCborLargeContainerLen", TestCborLargeContainerLen)
	t.Run("TestCborTime", TestCborTime)
	t.Run("TestCborUintToInt", TestCborUintToInt)
	t.Run("TestCborDifferentMapOrSliceType", TestCborDifferentMapOrSliceType)
	t.Run("TestCborScalars", TestCborScalars)
	t.Run("TestCborOmitempty", TestCborOmitempty)
	t.Run("TestCborIntfMapping", TestCborIntfMapping)
	t.Run("TestCborMissingFields", TestCborMissingFields)
	t.Run("TestCborMaxDepth", TestCborMaxDepth)
	t.Run("TestCborSelfExt", TestCborSelfExt)
	t.Run("TestCborBytesEncodedAsArray", TestCborBytesEncodedAsArray)
	t.Run("TestCborRawToStringToRawEtc", TestCborRawToStringToRawEtc)
	t.Run("TestCborStructKeyType", TestCborStructKeyType)
	t.Run("TestCborPreferArrayOverSlice", TestCborPreferArrayOverSlice)
	t.Run("TestCborZeroCopyBytes", TestCborZeroCopyBytes)
	t.Run("TestCborNumbers", TestCborNumbers)
	t.Run("TestCborDesc", TestCborDesc)
	t.Run("TestCborStructFieldInfoToArray", TestCborStructFieldInfoToArray)
	t.Run("TestCborAllErrWriter", TestCborAllErrWriter)
	t.Run("TestCborAllEncCircularRef", TestCborAllEncCircularRef)
	t.Run("TestCborAllAnonCycle", TestCborAllAnonCycle)
	t.Run("TestCborMammoth", TestCborMammoth)
	t.Run("TestCborMammothMapsAndSlices", TestCborMammothMapsAndSlices)
}

func testCborGroupV(t *testing.T) {
	t.Run("TestCborCodecsTable", TestCborCodecsTable)
	t.Run("TestCborCodecsMisc", TestCborCodecsMisc)
	t.Run("TestCborSwallowAndZero", TestCborSwallowAndZero)
	t.Run("TestCborNextValueBytes", TestCborNextValueBytes)
	t.Run("TestCborStrucEncDec", TestCborStrucEncDec)
}

func testMsgpackGroup(t *testing.T) {
	t.Run("TestMsgpackCodecsEmbeddedPointer", TestMsgpackCodecsEmbeddedPointer)
	t.Run("TestMsgpackStdEncIntf", TestMsgpackStdEncIntf)
	t.Run("TestMsgpackRaw", TestMsgpackRaw)
	t.Run("TestMsgpackRpcSpec", TestMsgpackRpcSpec)
	t.Run("TestMsgpackRpcGo", TestMsgpackRpcGo)
	t.Run("TestMsgpackMapEncodeForCanonical", TestMsgpackMapEncodeForCanonical)
	t.Run("TestMsgpackRawExt", TestMsgpackRawExt)
	t.Run("TestMsgpackMapStructKey", TestMsgpackMapStructKey)
	t.Run("TestMsgpackDecodeNilMapValue", TestMsgpackDecodeNilMapValue)
	t.Run("TestMsgpackEmbeddedFieldPrecedence", TestMsgpackEmbeddedFieldPrecedence)
	t.Run("TestMsgpackLargeContainerLen", TestMsgpackLargeContainerLen)
	t.Run("TestMsgpackTime", TestMsgpackTime)
	t.Run("TestMsgpackUintToInt", TestMsgpackUintToInt)
	t.Run("TestMsgpackDifferentMapOrSliceType", TestMsgpackDifferentMapOrSliceType)
	t.Run("TestMsgpackScalars", TestMsgpackScalars)
	t.Run("TestMsgpackOmitempty", TestMsgpackOmitempty)
	t.Run("TestMsgpackIntfMapping", TestMsgpackIntfMapping)
	t.Run("TestMsgpackMissingFields", TestMsgpackMissingFields)
	t.Run("TestMsgpackMaxDepth", TestMsgpackMaxDepth)
	t.Run("TestMsgpackSelfExt", TestMsgpackSelfExt)
	t.Run("TestMsgpackBytesEncodedAsArray", TestMsgpackBytesEncodedAsArray)
	t.Run("TestMsgpackRawToStringToRawEtc", TestMsgpackRawToStringToRawEtc)
	t.Run("TestMsgpackStructKeyType", TestMsgpackStructKeyType)
	t.Run("TestMsgpackPreferArrayOverSlice", TestMsgpackPreferArrayOverSlice)
	t.Run("TestMsgpackZeroCopyBytes", TestMsgpackZeroCopyBytes)
	t.Run("TestMsgpackNumbers", TestMsgpackNumbers)
	t.Run("TestMsgpackDesc", TestMsgpackDesc)
	t.Run("TestMsgpackStructFieldInfoToArray", TestMsgpackStructFieldInfoToArray)
	t.Run("TestMsgpackDecodeMapAndExtSizeMismatch", TestMsgpackDecodeMapAndExtSizeMismatch)
	t.Run("TestMsgpackMammoth", TestMsgpackMammoth)
	t.Run("TestMsgpackMammothMapsAndSlices", TestMsgpackMammothMapsAndSlices)
}

func testMsgpackGroupV(t *testing.T) {
	t.Run("TestMsgpackCodecsTable", TestMsgpackCodecsTable)
	t.Run("TestMsgpackCodecsMisc", TestMsgpackCodecsMisc)
	t.Run("TestMsgpackSwallowAndZero", TestMsgpackSwallowAndZero)
	t.Run("TestMsgpackNextValueBytes", TestMsgpackNextValueBytes)
	t.Run("TestMsgpackStrucEncDec", TestMsgpackStrucEncDec)
}

func testSimpleGroup(t *testing.T) {
	t.Run("TestSimpleCodecsEmbeddedPointer", TestSimpleCodecsEmbeddedPointer)
	t.Run("TestSimpleStdEncIntf", TestSimpleStdEncIntf)
	t.Run("TestSimpleRaw", TestSimpleRaw)
	t.Run("TestSimpleRpcGo", TestSimpleRpcGo)
	t.Run("TestSimpleMapEncodeForCanonical", TestSimpleMapEncodeForCanonical)
	t.Run("TestSimpleRawExt", TestSimpleRawExt)
	t.Run("TestSimpleMapStructKey", TestSimpleMapStructKey)
	t.Run("TestSimpleDecodeNilMapValue", TestSimpleDecodeNilMapValue)
	t.Run("TestSimpleEmbeddedFieldPrecedence", TestSimpleEmbeddedFieldPrecedence)
	t.Run("TestSimpleLargeContainerLen", TestSimpleLargeContainerLen)
	t.Run("TestSimpleTime", TestSimpleTime)
	t.Run("TestSimpleUintToInt", TestSimpleUintToInt)
	t.Run("TestSimpleDifferentMapOrSliceType", TestSimpleDifferentMapOrSliceType)
	t.Run("TestSimpleScalars", TestSimpleScalars)
	t.Run("TestSimpleOmitempty", TestSimpleOmitempty)
	t.Run("TestSimpleIntfMapping", TestSimpleIntfMapping)
	t.Run("TestSimpleMissingFields", TestSimpleMissingFields)
	t.Run("TestSimpleMaxDepth", TestSimpleMaxDepth)
	t.Run("TestSimpleSelfExt", TestSimpleSelfExt)
	t.Run("TestSimpleBytesEncodedAsArray", TestSimpleBytesEncodedAsArray)
	t.Run("TestSimpleRawToStringToRawEtc", TestSimpleRawToStringToRawEtc)
	t.Run("TestSimpleStructKeyType", TestSimpleStructKeyType)
	t.Run("TestSimplePreferArrayOverSlice", TestSimplePreferArrayOverSlice)
	t.Run("TestSimpleZeroCopyBytes", TestSimpleZeroCopyBytes)
	t.Run("TestSimpleNumbers", TestSimpleNumbers)
	t.Run("TestSimpleDesc", TestSimpleDesc)
	t.Run("TestSimpleStructFieldInfoToArray", TestSimpleStructFieldInfoToArray)
	t.Run("TestSimpleMultipleEncDec", TestSimpleMultipleEncDec)
	t.Run("TestSimpleAllErrWriter", TestSimpleAllErrWriter)
	t.Run("TestSimpleMammoth", TestSimpleMammoth)
	t.Run("TestSimpleMammothMapsAndSlices", TestSimpleMammothMapsAndSlices)
}

func testSimpleGroupV(t *testing.T) {
	t.Run("TestSimpleCodecsTable", TestSimpleCodecsTable)
	t.Run("TestSimpleCodecsMisc", TestSimpleCodecsMisc)
	t.Run("TestSimpleSwallowAndZero", TestSimpleSwallowAndZero)
	t.Run("TestSimpleNextValueBytes", TestSimpleNextValueBytes)
	t.Run("TestSimpleStrucEncDec", TestSimpleStrucEncDec)
}

func testSimpleMammothGroup(t *testing.T) {
	t.Run("TestSimpleMammothMapsAndSlices", TestSimpleMammothMapsAndSlices)
}

func testRpcGroup(t *testing.T) {
	t.Run("TestBincRpcGo", TestBincRpcGo)
	t.Run("TestSimpleRpcGo", TestSimpleRpcGo)
	t.Run("TestMsgpackRpcGo", TestMsgpackRpcGo)
	t.Run("TestCborRpcGo", TestCborRpcGo)
	t.Run("TestJsonRpcGo", TestJsonRpcGo)
	t.Run("TestMsgpackRpcSpec", TestMsgpackRpcSpec)
}

func testNonHandlesGroup(t *testing.T) {
	// grep "func Test" codec_test.go | grep -v -E '(Cbor|Json|Simple|Msgpack|Binc)'
	t.Run("TestMapRangeIndex", TestMapRangeIndex)
}

func TestCodecSuite(t *testing.T) {
	// steps for each run:
	// - testGroupResetBase     - default reset of testv
	// - update testv           - (run-specific) updates to testv (which will apply to all handles within testReinit)
	// - testReinit             - recreate handles and apply testv updates
	// - testGroupResetHandles  - default updates to handles
	// - update handles         - (run-specific) updates to handles for this specific Run operation
	// - RUN                    - (run-specific) execution
	//
	// Most of these steps are optional

	var tt testTimeTracker
	tt.Elapsed()

	fnRun := func(s string, f func(t *testing.T)) {
		t.Run(s, f)
	}
	fnb2i := func(b bool, i, j int) int {
		if b {
			return i
		}
		return j
	}

	testGroupResetBase := func() {
		tbvars.D = DecodeOptions{}
		tbvars.E = EncodeOptions{}
		tbvars.R = RPCOptions{}
		tbvars.updateHandleOptions()

		testv.RpcBufsize = 2048
		testv.UseReset = false
		testv.UseParallel = false
		testv.UseIoWrapper = false
		testv.NumRepeatString = 8
		testv.Depth = 0
	}

	testGroupResetHandles := func() {
		testJsonH.Indent = 0 // -1?
		testJsonH.HTMLCharsAsIs = false
		testJsonH.MapKeyAsString = false
		testJsonH.PreferFloat = false

		testCborH.IndefiniteLength = false
		testCborH.TimeRFC3339 = false
		testCborH.SkipUnexpectedTags = false
		testBincH.AsSymbols = 0 // 2? AsSymbolNone

		testMsgpackH.WriteExt = false
		testMsgpackH.NoFixedNum = false
		testMsgpackH.PositiveIntUnsigned = false

		testSimpleH.EncZeroValuesAsNil = false
	}

	fnTrueV := func() {
		testGroupResetBase()
		tbvars.setBufsize(0)
		testv.UseReset = true
		testv.UseParallel = true

		tbvars.D.ZeroCopy = true
		tbvars.D.InternString = true
		tbvars.D.MapValueReset = true

		// tbvars.D.ErrorIfNoField = true // error, as expected fields not there
		// tbvars.D.ErrorIfNoArrayExpand = true // no error, but no error case either
		// tbvars.D.PreferArrayOverSlice = true // error??? because slice != array.
		tbvars.D.SignedInteger = true // error as deepEqual compares int64 to uint64
		tbvars.D.SliceElementReset = true
		tbvars.D.InterfaceReset = true
		tbvars.D.RawToString = true
		tbvars.D.PreferPointerForStructOrArray = true

		tbvars.E.StructToArray = true
		tbvars.E.Canonical = true
		tbvars.E.CheckCircularRef = true
		tbvars.E.RecursiveEmptyCheck = true
		tbvars.E.OptimumSize = true
		tbvars.E.NilCollectionToZeroLength = true
		// MARKER: we cannot test these below, as they will not encode as expected
		// meaning a decoded value will look different than expected.
		// e.g. encode nil slice, and get a decoded stream with zero-length array
		//
		// Consequently, we don't modify these here.
		// Standalone unit tests will test these out in codec_run_test.go
		//
		// tbvars.E.Raw = true
		// tbvars.E.StringToRaw = true
	}

	fnTrueH := func() {
		// MARKER: same as above - do not modify fields which cause the
		// encoded stream to contain unexpected values than what a decode would expect.
		//
		// testJsonH.MapKeyAsString = true
		// testJsonH.PreferFloat = true
		// testSimpleH.EncZeroValuesAsNil = true
		testJsonH.HTMLCharsAsIs = true
		testCborH.TimeRFC3339 = true
		testMsgpackH.WriteExt = true
		testMsgpackH.NoFixedNum = true
		testMsgpackH.PositiveIntUnsigned = true
		testCborH.IndefiniteLength = true
		testCborH.SkipUnexpectedTags = true
	}
	// --------------
	testGroupResetBase()
	tbvars.setBufsize(-1)
	testReinit()
	testGroupResetHandles()
	fnRun("optionsFalse", testCodecGroup)

	// --------------
	tbvars.setBufsize(0)
	testReinit()
	testGroupResetHandles()
	fnRun("optionsFalse-io", testCodecGroup)

	// --------------
	fnTrueV()
	testReinit()
	testGroupResetHandles()
	fnTrueH()
	fnRun("optionsTrue", testCodecGroup)

	// --------------
	fnTrueV()
	testv.Depth = fnb2i(testing.Short(), 2, 4)
	testReinit()
	testGroupResetHandles()
	fnRun("optionsTrue-deepstruct", testCodecGroupV)

	// --------------
	fnTrueV()
	testv.Depth = 0
	// tbvars.E.AsSymbols = AsSymbolAll
	testv.UseIoWrapper = true
	testReinit()
	testGroupResetHandles()
	fnTrueH()
	fnRun("optionsTrue-ioWrapper", testCodecGroupV)

	// --------------
	fnTrueV()
	testv.UseIoWrapper = false
	// testv.UseIoEncDec = -1
	// ---
	// make buffer small enough so that we have to re-fill multiple times
	// and also such that writing a quoted struct name e.g. "LongFieldNameXYZ"
	// will require a re-fill, and test out bufioEncWriter.writeqstr well.
	// Due to last requirement, we prefer 16 to 128.
	testv.SkipRPCTests = true
	tbvars.setBufsize(16)
	// tbvars.D.ReaderBufferSize = 128
	// tbvars.E.WriterBufferSize = 128
	testReinit()
	testGroupResetHandles()
	fnTrueH()
	fnRun("optionsTrue-bufio", testCodecGroupV)

	// --------------
	fnTrueV()
	// tbvars.D.ReaderBufferSize = 0
	// tbvars.E.WriterBufferSize = 0
	testv.SkipRPCTests = false
	tbvars.setBufsize(-1)
	// ---
	testv.NumRepeatString = 32
	testReinit()
	testGroupResetHandles()
	fnTrueH()
	fnRun("optionsTrue-largestrings", testCodecGroupV)

	// --------------
	testGroupResetBase()
	testv.NumRepeatString = 8
	defer func(ml int, d int8, hca, mkas bool) {
		tbvars.D.MaxInitLen = ml
		testJsonH.Indent = d
		testJsonH.HTMLCharsAsIs = hca
		testJsonH.MapKeyAsString = mkas
	}(tbvars.D.MaxInitLen, testJsonH.Indent, testJsonH.HTMLCharsAsIs, testJsonH.MapKeyAsString)
	tbvars.D.MaxInitLen = 10
	testReinit()
	testGroupResetHandles()
	testJsonH.MapKeyAsString = true
	testJsonH.Indent = 8
	testJsonH.HTMLCharsAsIs = true
	fnRun("json-spaces-htmlcharsasis-initLen10", testJsonGroup)

	testJsonH.Indent = -1
	testJsonH.HTMLCharsAsIs = false
	fnRun("json-tabs-initLen10", testJsonGroup)

	defer func(v uint8) { testBincH.AsSymbols = v }(testBincH.AsSymbols)
	testBincH.AsSymbols = 2 // AsSymbolNone
	fnRun("binc-no-symbols", testBincGroup)

	testBincH.AsSymbols = 1 // AsSymbolAll
	fnRun("binc-all-symbols", testBincGroup)

	defer func(v bool) { testSimpleH.EncZeroValuesAsNil = v }(testSimpleH.EncZeroValuesAsNil)
	testSimpleH.EncZeroValuesAsNil = !testSimpleH.EncZeroValuesAsNil
	fnRun("simple-enczeroasnil", testSimpleMammothGroup) // testSimpleGroup

	// --------------
	defer tbvars.setBufsize((int)(testv.bufsize))
	defer func(b bool) { tbvars.R.RPCNoBuffer = b }(tbvars.R.RPCNoBuffer)
	tbvars.setBufsize(16)
	tbvars.R.RPCNoBuffer = false
	testReinit()
	testGroupResetHandles()
	testv.RpcBufsize = 0
	fnRun("rpc-buf-0", testRpcGroup)
	testv.RpcBufsize = 0
	fnRun("rpc-buf-00", testRpcGroup)
	testv.RpcBufsize = 0
	fnRun("rpc-buf-000", testRpcGroup)
	testv.RpcBufsize = 16
	fnRun("rpc-buf-16", testRpcGroup)
	testv.RpcBufsize = 2048
	fnRun("rpc-buf-2048", testRpcGroup)
	tbvars.R.RPCNoBuffer = true
	testv.RpcBufsize = 0
	fnRun("rpc-buf-0-rpcNoBuffer", testRpcGroup)
	testv.RpcBufsize = 0
	fnRun("rpc-buf-00-rpcNoBuffer", testRpcGroup)
	testv.RpcBufsize = 2048
	fnRun("rpc-buf-2048-rpcNoBuffer", testRpcGroup)

	testGroupResetBase()
	testReinit()
	testGroupResetHandles()
}
