// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "testing"

func TestBincCodecsTable(t *testing.T) {
	doTestCodecTableOne(t, testBincH)
}

func TestBincCodecsMisc(t *testing.T) {
	doTestCodecMiscOne(t, testBincH)
}

func TestBincCodecsEmbeddedPointer(t *testing.T) {
	doTestCodecEmbeddedPointer(t, testBincH)
}

func TestBincStdEncIntf(t *testing.T) {
	doTestStdEncIntf(t, testBincH)
}

func TestBincRaw(t *testing.T) {
	doTestRawValue(t, testBincH)
}

func TestBincRpcGo(t *testing.T) {
	doTestCodecRpcOne(t, GoRpc, testBincH, true, 0)
}

func TestBincMapEncodeForCanonical(t *testing.T) {
	t.Skipf("skipping ... needs investigation") // MARKER: testing fails??? Need to research
	doTestMapEncodeForCanonical(t, testBincH)
}

func TestBincUnderlyingType(t *testing.T) {
	testCodecUnderlyingType(t, testBincH)
}

func TestBincSwallowAndZero(t *testing.T) {
	doTestSwallowAndZero(t, testBincH)
}

func TestBincRawExt(t *testing.T) {
	doTestRawExt(t, testBincH)
}

func TestBincMapStructKey(t *testing.T) {
	doTestMapStructKey(t, testBincH)
}

func TestBincDecodeNilMapValue(t *testing.T) {
	doTestDecodeNilMapValue(t, testBincH)
}

func TestBincEmbeddedFieldPrecedence(t *testing.T) {
	doTestEmbeddedFieldPrecedence(t, testBincH)
}

func TestBincLargeContainerLen(t *testing.T) {
	doTestLargeContainerLen(t, testBincH)
}

func TestBincTime(t *testing.T) {
	doTestTime(t, testBincH)
}

func TestBincUintToInt(t *testing.T) {
	doTestUintToInt(t, testBincH)
}

func TestBincDifferentMapOrSliceType(t *testing.T) {
	doTestDifferentMapOrSliceType(t, testBincH)
}

func TestBincScalars(t *testing.T) {
	doTestScalars(t, testBincH)
}

func TestBincOmitempty(t *testing.T) {
	doTestOmitempty(t, testBincH)
}

func TestBincIntfMapping(t *testing.T) {
	doTestIntfMapping(t, testBincH)
}

func TestBincMissingFields(t *testing.T) {
	doTestMissingFields(t, testBincH)
}

func TestBincMaxDepth(t *testing.T) {
	doTestMaxDepth(t, testBincH)
}

func TestBincSelfExt(t *testing.T) {
	doTestSelfExt(t, testBincH)
}

func TestBincBytesEncodedAsArray(t *testing.T) {
	doTestBytesEncodedAsArray(t, testBincH)
}

func TestBincStrucEncDec(t *testing.T) {
	doTestStrucEncDec(t, testBincH)
}

func TestBincRawToStringToRawEtc(t *testing.T) {
	doTestRawToStringToRawEtc(t, testBincH)
}

func TestBincStructKeyType(t *testing.T) {
	doTestStructKeyType(t, testBincH)
}

func TestBincPreferArrayOverSlice(t *testing.T) {
	doTestPreferArrayOverSlice(t, testBincH)
}

func TestBincZeroCopyBytes(t *testing.T) {
	doTestZeroCopyBytes(t, testBincH)
}

func TestBincNextValueBytes(t *testing.T) {
	doTestNextValueBytes(t, testBincH)
}

func TestBincNumbers(t *testing.T) {
	doTestNumbers(t, testBincH)
}

func TestBincDesc(t *testing.T) {
	m := make(map[byte]string)
	for k, v := range bincdescVdNames {
		m[k<<4] = v
	}
	for k, v := range bincdescSpecialVsNames {
		m[k] = v
	}
	delete(m, bincVdSpecial<<4)
	doTestDesc(t, testBincH, m)
}

func TestBincStructFieldInfoToArray(t *testing.T) {
	doTestStructFieldInfoToArray(t, testBincH)
}
