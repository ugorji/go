// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"testing"
)

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

func TestSimpleMultipleEncDec(t *testing.T) {
	doTestMultipleEncDec(t, testSimpleH)
}

func TestSimpleAllErrWriter(t *testing.T) {
	doTestAllErrWriter(t, testSimpleH)
}
