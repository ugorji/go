// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"testing"
)

func init() {
	testPreInitFns = append(testPreInitFns, simpleTestInit)
}

func simpleTestInit() {
	// create legacy functions suitable for deprecated AddExt functionality,
	// and use on some places for testSimpleH e.g. for time.Time and wrapInt64
	var tUintToBytesExt testUintToBytesExt
	var tBytesExt wrapBytesExt
	var tI64Ext wrapInt64Ext
	var tTimeExt timeBytesExt

	timeExtEncFn := func(rv reflect.Value) (bs []byte, err error) { return basicTestExtEncFn(tTimeExt, rv) }
	timeExtDecFn := func(rv reflect.Value, bs []byte) (err error) { return basicTestExtDecFn(tTimeExt, rv, bs) }
	wrapInt64ExtEncFn := func(rv reflect.Value) (bs []byte, err error) { return basicTestExtEncFn(&tI64Ext, rv) }
	wrapInt64ExtDecFn := func(rv reflect.Value, bs []byte) (err error) { return basicTestExtDecFn(&tI64Ext, rv, bs) }

	halt.onerror(testSimpleH.AddExt(timeTyp, 1, timeExtEncFn, timeExtDecFn))
	halt.onerror(testSimpleH.SetBytesExt(testSelfExtTyp, 78, SelfExt))
	halt.onerror(testSimpleH.SetBytesExt(testSelfExt2Typ, 79, SelfExt))
	halt.onerror(testSimpleH.SetBytesExt(wrapBytesTyp, 32, &tBytesExt))
	halt.onerror(testSimpleH.SetBytesExt(testUintToBytesTyp, 33, &tUintToBytesExt))
	halt.onerror(testSimpleH.AddExt(wrapInt64Typ, 16, wrapInt64ExtEncFn, wrapInt64ExtDecFn))
}

func TestSimpleCodecsTable(t *testing.T) {
	doTestCodecTableOne(t, testSimpleH)
}

func TestSimpleCodecsMisc(t *testing.T) {
	doTestCodecMiscOne(t, testSimpleH)
}

func TestSimpleCodecsEmbeddedPointer(t *testing.T) {
	doTestCodecEmbeddedPointer(t, testSimpleH)
}

func TestSimpleStdEncIntf(t *testing.T) {
	doTestStdEncIntf(t, testSimpleH)
}

func TestSimpleRaw(t *testing.T) {
	doTestRawValue(t, testSimpleH)
}

func TestSimpleRpcGo(t *testing.T) {
	doTestCodecRpcOne(t, GoRpc, testSimpleH, true, 0)
}

func TestSimpleMapEncodeForCanonical(t *testing.T) {
	doTestMapEncodeForCanonical(t, testSimpleH)
}

func TestSimpleSwallowAndZero(t *testing.T) {
	doTestSwallowAndZero(t, testSimpleH)
}

func TestSimpleRawExt(t *testing.T) {
	doTestRawExt(t, testSimpleH)
}

func TestSimpleMapStructKey(t *testing.T) {
	doTestMapStructKey(t, testSimpleH)
}

func TestSimpleDecodeNilMapValue(t *testing.T) {
	doTestDecodeNilMapValue(t, testSimpleH)
}

func TestSimpleEmbeddedFieldPrecedence(t *testing.T) {
	doTestEmbeddedFieldPrecedence(t, testSimpleH)
}

func TestSimpleLargeContainerLen(t *testing.T) {
	doTestLargeContainerLen(t, testSimpleH)
}

func TestSimpleTime(t *testing.T) {
	doTestTime(t, testSimpleH)
}

func TestSimpleUintToInt(t *testing.T) {
	doTestUintToInt(t, testSimpleH)
}

func TestSimpleDifferentMapOrSliceType(t *testing.T) {
	doTestDifferentMapOrSliceType(t, testSimpleH)
}

func TestSimpleScalars(t *testing.T) {
	doTestScalars(t, testSimpleH)
}

func TestSimpleOmitempty(t *testing.T) {
	doTestOmitempty(t, testSimpleH)
}

func TestSimpleIntfMapping(t *testing.T) {
	doTestIntfMapping(t, testSimpleH)
}

func TestSimpleMissingFields(t *testing.T) {
	doTestMissingFields(t, testSimpleH)
}

func TestSimpleMaxDepth(t *testing.T) {
	doTestMaxDepth(t, testSimpleH)
}

func TestSimpleSelfExt(t *testing.T) {
	doTestSelfExt(t, testSimpleH)
}

func TestSimpleBytesEncodedAsArray(t *testing.T) {
	doTestBytesEncodedAsArray(t, testSimpleH)
}

func TestSimpleStrucEncDec(t *testing.T) {
	doTestStrucEncDec(t, testSimpleH)
}

func TestSimpleRawToStringToRawEtc(t *testing.T) {
	doTestRawToStringToRawEtc(t, testSimpleH)
}

func TestSimpleStructKeyType(t *testing.T) {
	doTestStructKeyType(t, testSimpleH)
}

func TestSimplePreferArrayOverSlice(t *testing.T) {
	doTestPreferArrayOverSlice(t, testSimpleH)
}

func TestSimpleZeroCopyBytes(t *testing.T) {
	doTestZeroCopyBytes(t, testSimpleH)
}

func TestSimpleNextValueBytes(t *testing.T) {
	doTestNextValueBytes(t, testSimpleH)
}

func TestSimpleNumbers(t *testing.T) {
	doTestNumbers(t, testSimpleH)
}

func TestSimpleDesc(t *testing.T) {
	doTestDesc(t, testSimpleH, simpledescNames)
}

func TestSimpleStructFieldInfoToArray(t *testing.T) {
	doTestStructFieldInfoToArray(t, testSimpleH)
}
