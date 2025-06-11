// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding/hex"
	"math"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestCborIndefiniteLength(t *testing.T) {
	ch := testCborH
	var h Handle = ch
	defer testSetup2(t, &h)()
	bh := testBasicHandle(h)
	defer func(a reflect.Type, b, c, d bool) {
		bh.MapType, bh.RawToString, bh.SignedInteger, ch.IndefiniteLength = a, b, c, d
	}(bh.MapType, bh.RawToString, bh.SignedInteger, ch.IndefiniteLength)
	bh.MapType = testMapStrIntfTyp
	bh.RawToString = false
	bh.SignedInteger = false
	// since we're writing this manually using MustEncode, don't use IndefiniteLength
	ch.IndefiniteLength = false
	// var (
	// 	M1 map[string][]byte
	// 	M2 map[uint64]bool
	// 	L1 []interface{}
	// 	S1 []string
	// 	B1 []byte
	// )
	var v, vv interface{}
	// define it (v), encode it using indefinite lengths, decode it (vv), compare v to vv
	v = map[string]interface{}{
		"one-byte-key":   []byte{1, 2, 3, 4, 5, 6},
		"two-string-key": "two-value",
		"three-list-key": []interface{}{true, false, uint64(1), int64(-1)},
	}
	var buf bytes.Buffer
	// buf.Reset()
	e := NewEncoder(&buf, h)
	buf.WriteByte(cborBdIndefiniteMap)
	//----
	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("one-")
	e.MustEncode("byte-")
	e.MustEncode("key")
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdIndefiniteBytes)
	e.MustEncode([]byte{1, 2, 3})
	e.MustEncode([]byte{4, 5, 6})
	buf.WriteByte(cborBdBreak)

	//----
	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("two-")
	e.MustEncode("string-")
	e.MustEncode("key")
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("two-")
	e.MustEncode("value")
	buf.WriteByte(cborBdBreak)

	//----
	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("three-")
	e.MustEncode("list-")
	e.MustEncode("key")
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdIndefiniteArray)
	e.MustEncode(true)
	e.MustEncode(false)
	e.MustEncode(uint64(1))
	e.MustEncode(int64(-1))
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdBreak) // close map

	NewDecoderBytes(buf.Bytes(), h).MustDecode(&vv)
	if err := testEqual(v, vv); err != nil {
		t.Logf("-------- Before and After marshal do not match: Error: %v", err)
		if testv.Verbose {
			t.Logf("    ....... GOLDEN:  (%T) %#v", v, v)
			t.Logf("    ....... DECODED: (%T) %#v", vv, vv)
		}
		t.FailNow()
	}
}

// "If any item between the indefinite-length string indicator (0b010_11111 or 0b011_11111) and the
// "break" stop code is not a definite-length string item of the same major type, the string is not
// well-formed."
func TestCborIndefiniteLengthStringChunksCannotMixTypes(t *testing.T) {
	if !testRecoverPanicToErr {
		t.Skip(testSkipIfNotRecoverPanicToErrMsg)
	}
	var h Handle = testCborH
	defer testSetup(t, &h)()

	for _, in := range [][]byte{
		{cborBdIndefiniteString, 0x40, cborBdBreak}, // byte string chunk in indefinite length text string
		{cborBdIndefiniteBytes, 0x60, cborBdBreak},  // text string chunk in indefinite length byte string
	} {
		var out string
		err := NewDecoderBytes(in, h).Decode(&out)
		if err == nil {
			t.Errorf("expected error but decoded 0x%x to: %q", in, out)
		}
	}
}

// "If any definite-length text string inside an indefinite-length text string is invalid, the
// indefinite-length text string is invalid. Note that this implies that the UTF-8 bytes of a single
// Unicode code point (scalar value) cannot be spread between chunks: a new chunk of a text string
// can only be started at a code point boundary."
func TestCborIndefiniteLengthTextStringChunksAreUTF8(t *testing.T) {
	if !testRecoverPanicToErr {
		t.Skip(testSkipIfNotRecoverPanicToErrMsg)
	}
	var h Handle = testCborH
	defer testSetup2(t, &h)()

	bh := testBasicHandle(h)
	if !bh.ValidateUnicode {
		bh.ValidateUnicode = true
		defer func() { bh.ValidateUnicode = false }()
	}

	var out string
	in := []byte{cborBdIndefiniteString, 0x61, 0xc2, 0x61, 0xa3, cborBdBreak}
	err := NewDecoderBytes(in, h).Decode(&out)
	if err == nil {
		t.Errorf("expected error but decoded to: %q", out)
	}
}

type testCborGolden struct {
	Base64     string      `codec:"cbor"`
	Hex        string      `codec:"hex"`
	Roundtrip  bool        `codec:"roundtrip"`
	Decoded    interface{} `codec:"decoded"`
	Diagnostic string      `codec:"diagnostic"`
	Skip       bool        `codec:"skip"`
}

// Some tests are skipped because they include numbers outside the range of int64/uint64
func TestCborGoldens(t *testing.T) {
	ch := testCborH
	var h Handle = ch
	defer testSetup2(t, &h)()
	bh := testBasicHandle(h)
	defer func(a reflect.Type, b, c, d bool) {
		bh.MapType, bh.RawToString, bh.SignedInteger, ch.SkipUnexpectedTags = a, b, c, d
	}(bh.MapType, bh.RawToString, bh.SignedInteger, ch.SkipUnexpectedTags)
	bh.MapType = testMapStrIntfTyp
	bh.RawToString = false
	bh.SignedInteger = false
	ch.SkipUnexpectedTags = false

	// decode test-cbor-goldens.json into a list of []*testCborGolden
	// for each one,
	// - decode hex into []byte bs
	// - decode bs into interface{} v
	// - compare both using deepequal
	// - for any miss, record it
	var gs []*testCborGolden
	f, err := os.ReadFile("test-cbor-goldens.json")
	if err != nil {
		t.Logf("error opening test-cbor-goldens.json: %v", err)
		t.FailNow()
	}
	jh := testJsonH // new(JsonHandle)
	defer func(v reflect.Type, b bool) { jh.MapType, jh.SignedInteger = v, b }(jh.MapType, jh.SignedInteger)
	jh.MapType = testMapStrIntfTyp
	jh.SignedInteger = false

	// d := NewDecoder(f, jh)
	// d := NewDecoder(bufio.NewReader(f), jh)
	d := NewDecoderBytes(f, jh)
	// err = d.Decode(&gs)
	d.MustDecode(&gs)
	if err != nil {
		t.Logf("error json decoding test-cbor-goldens.json: %v", err)
		t.FailNow()
	}

	tagregex := regexp.MustCompile(`[\d]+\(.+?\)`)
	hexregex := regexp.MustCompile(`h'([0-9a-fA-F]*)'`)
	var bs []byte
	for i, g := range gs {
		// skip tags or simple or those with prefix, as we can't verify them.
		if g.Skip || strings.HasPrefix(g.Diagnostic, "simple(") || tagregex.MatchString(g.Diagnostic) {
			if testv.Verbose {
				t.Logf("[%v] skipping because skip=true OR unsupported simple value or Tag Value", i)
			}
			continue
		}
		if hexregex.MatchString(g.Diagnostic) {
			// println(i, "g.Diagnostic matched hex")
			if s2 := g.Diagnostic[2 : len(g.Diagnostic)-1]; s2 == "" {
				g.Decoded = zeroByteSlice
			} else if bs2, err2 := hex.DecodeString(s2); err2 == nil {
				g.Decoded = bs2
			}
		}
		bs, err = hex.DecodeString(g.Hex)
		if err != nil {
			t.Logf("[%v] error hex decoding %s [%v]: %v", i, g.Hex, g.Hex, err)
			t.FailNow()
		}
		var v interface{}
		NewDecoderBytes(bs, h).MustDecode(&v)
		switch v.(type) {
		case RawExt:
			continue
		case *RawExt:
			continue
		}
		// check the diagnostics to compare
		var b bool
		var v0 interface{}
		switch g.Diagnostic {
		case "Infinity":
			b = math.IsInf(v.(float64), 1)
			v0 = math.Inf(1)
		case "-Infinity":
			b = math.IsInf(v.(float64), -1)
			v0 = math.Inf(-1)
		case "NaN":
			// println(i, "checking NaN")
			b = math.IsNaN(v.(float64))
			v0 = math.NaN()
		case "undefined":
			b = v == nil
		default:
			v0 = g.Decoded
			err = testEqual(v0, v)
			b = true
			// testCborCoerceJsonNumber(reflect.ValueOf(&v0))
		}
		testCborError(t, g.Diagnostic, i, v0, v, err, b)
	}
}

func testCborError(t *testing.T, diagnostic string, i int, v0, v1 interface{}, err error, equal bool) {
	// only err or equal is set
	// if err == nil && equal {
	// 	return
	// }
	if err != nil {
		t.Logf("[%v] testEqual error: %v", i, err)
		if testv.Verbose {
			t.Logf("    ....... GOLDEN:  (%T) %#v", v0, v0)
			t.Logf("    ....... DECODED: (%T) %#v", v1, v1)
		}
		t.FailNow()
		return
	}
	if !equal {
		t.Logf("[%v] values not equal", i)
		if testv.Verbose {
			t.Logf("    ....... GOLDEN:  (%T) %#v", v0, v0)
			t.Logf("    ....... DECODED: (%T) %#v", v1, v1)
		}
		t.FailNow()
		return
	}
}

func TestCborHalfFloat(t *testing.T) {
	var h Handle = testCborH
	defer testSetup(t, &h)()
	m := map[uint16]float64{
		// using examples from
		// https://en.wikipedia.org/wiki/Half-precision_floating-point_format
		0x3c00: 1,
		0x3c01: 1 + math.Pow(2, -10),
		0xc000: -2,
		0x7bff: 65504,
		0x0400: math.Pow(2, -14),
		0x03ff: math.Pow(2, -14) - math.Pow(2, -24),
		0x0001: math.Pow(2, -24),
		0x0000: 0,
		0x8000: 0.0, // -0.0
	}
	var ba [3]byte
	ba[0] = cborBdFloat16
	var res float64
	for k, v := range m {
		res = 0
		bigenstd.PutUint16(ba[1:], k)
		testUnmarshalErr(&res, ba[:3], h, t, "-")
		if res == v {
			if testv.Verbose {
				t.Logf("equal floats: from %x %b, %v", k, k, v)
			}
		} else {
			t.Logf("unequal floats: from %x %b, %v != %v", k, k, res, v)
			t.FailNow()
		}
	}
}

func TestCborSkipTags(t *testing.T) {
	var hh Handle = testCborH
	defer testSetup2(t, &hh)()
	type Tcbortags struct {
		A string
		M map[string]interface{}
		// A []interface{}
	}
	var b8 [8]byte
	var w bytesEncAppender
	w.b = []byte{}

	// To make it easier,
	//    - use tags between math.MaxUint8 and math.MaxUint16 (incl SelfDesc)
	//    - use 1 char strings for key names
	//    - use 3-6 char strings for map keys
	//    - use integers that fit in 2 bytes (between 0x20 and 0xff)

	var tags = [...]uint64{math.MaxUint8 * 2, math.MaxUint8 * 8, 55799, math.MaxUint16 / 2}
	var tagIdx int
	var doAddTag bool
	addTagFn8To16 := func() {
		if !doAddTag {
			return
		}
		// writes a tag between MaxUint8 and MaxUint16 (culled from cborEncDriver.encUint)
		w.writen1(cborBaseTag + 0x19)
		// bigenHelper.writeUint16
		bigenstd.PutUint16(b8[:2], uint16(tags[tagIdx%len(tags)]))
		w.writeb(b8[:2])
		tagIdx++
	}

	var v Tcbortags
	v.A = "cbor"
	v.M = make(map[string]interface{})
	v.M["111"] = uint64(111)
	v.M["111.11"] = 111.11
	v.M["true"] = true
	// v.A = append(v.A, 222, 22.22, "true")

	// make stream manually (interspacing tags around it)
	// WriteMapStart - e.encLen(cborBaseMap, length) - encUint(length, bd)
	// EncodeStringEnc - e.encStringBytesS(cborBaseString, v)

	fnEncode := func() {
		w.b = w.b[:0]
		// addTagFn8To16()
		// write v (Tcbortags, with 2 fields = map with 2 entries)
		w.writen1(2 + cborBaseMap) // 2 fields = 2 entries
		// write v.A
		addTagFn8To16()
		var s = "A"
		w.writen1(byte(len(s)) + cborBaseString)
		w.writestr(s)
		w.writen1(byte(len(v.A)) + cborBaseString)
		w.writestr(v.A)
		//w.writen1(0)

		addTagFn8To16()
		s = "M"
		w.writen1(byte(len(s)) + cborBaseString)
		w.writestr(s)

		addTagFn8To16()
		w.writen1(byte(len(v.M)) + cborBaseMap)

		addTagFn8To16()
		s = "111"
		w.writen1(byte(len(s)) + cborBaseString)
		w.writestr(s)
		w.writen2(cborBaseUint+0x18, uint8(111))

		addTagFn8To16()
		s = "111.11"
		w.writen1(byte(len(s)) + cborBaseString)
		w.writestr(s)
		w.writen1(cborBdFloat64)
		bigenstd.PutUint64(b8[:8], math.Float64bits(111.11))
		w.writeb(b8[:8])

		addTagFn8To16()
		s = "true"
		w.writen1(byte(len(s)) + cborBaseString)
		w.writestr(s)
		w.writen1(cborBdTrue)
	}

	var h = testCborH // &CborHandle{}
	defer func(a, b, c, d, e, f bool) {
		h.SkipUnexpectedTags, h.Canonical, h.ZeroCopy = a, b, c
		h.SignedInteger, h.StructToArray, h.IndefiniteLength = d, e, f
	}(h.SkipUnexpectedTags, h.Canonical, h.ZeroCopy,
		h.SignedInteger, h.StructToArray, h.IndefiniteLength)
	h.SkipUnexpectedTags = true
	h.Canonical = true
	h.ZeroCopy = false // so we can reuse the w.b slice without worrying about sharedBytes
	h.SignedInteger = false
	h.StructToArray = false
	// manually encoding the value, so must not have IndefiniteLength=true
	h.IndefiniteLength = false

	var gold []byte
	NewEncoderBytes(&gold, h).MustEncode(v)

	// w.b is the encoded bytes
	var v2 Tcbortags
	doAddTag = false
	fnEncode()

	testDeepEqualErr(gold, w.b, t, "cbor-skip-tags--bytes---")
	NewDecoderBytes(w.b, h).MustDecode(&v2)
	testDeepEqualErr(v, v2, t, "cbor-skip-tags--no-tags-")

	var v3 Tcbortags
	doAddTag = true
	fnEncode()
	NewDecoderBytes(w.b, h).MustDecode(&v3)
	testDeepEqualErr(v, v3, t, "cbor-skip-tags--has-tags")
	// include this test below (against v2) to ensure that new writes to w.b does not affect v2
	testDeepEqualErr(v, v2, t, "cbor-skip-tags--has-tags")

	// Github 300 - tests naked path
	{
		expected := []interface{}{"x", uint64(0x0)}
		toDecode := []byte{0x82, 0x61, 0x78, 0x00}

		var raw interface{}

		NewDecoderBytes(toDecode, h).MustDecode(&raw)
		testDeepEqualErr(expected, raw, t, "cbor-skip-tags--gh-300---no-skips")

		toDecode = []byte{0xd9, 0xd9, 0xf7, 0x82, 0x61, 0x78, 0x00}
		raw = nil
		NewDecoderBytes(toDecode, h).MustDecode(&raw)
		testDeepEqualErr(expected, raw, t, "cbor-skip-tags--gh-300--has-skips")
	}
}

func TestCborMalformed(t *testing.T) {
	if !testRecoverPanicToErr {
		t.Skip(testSkipIfNotRecoverPanicToErrMsg)
	}
	var h Handle = testCborH
	defer testSetup(t, &h)()
	var bad = [][]byte{
		[]byte("\x9b\x00\x00000000"),
		[]byte("\x9b\x00\x00\x81112233"),
	}

	var out interface{}
	for _, v := range bad {
		out = nil
		err := testUnmarshal(&out, v, h)
		if err == nil {
			t.Logf("missing expected error decoding malformed cbor")
			t.FailNow()
		}
	}
}

// ----------
//
// cbor tests shared with other formats

func TestCborCodecsTable(t *testing.T) {
	doTestCodecTableOne(t, testCborH)
}

func TestCborCodecsMisc(t *testing.T) {
	doTestCodecMiscOne(t, testCborH)
}

func TestCborCodecsEmbeddedPointer(t *testing.T) {
	doTestCodecEmbeddedPointer(t, testCborH)
}

func TestCborCodecChan(t *testing.T) {
	doTestCodecChan(t, testCborH)
}

func TestCborStdEncIntf(t *testing.T) {
	doTestStdEncIntf(t, testCborH)
}

func TestCborRaw(t *testing.T) {
	doTestRawValue(t, testCborH)
}

// ----- RPC -----

func TestCborRpcGo(t *testing.T) {
	doTestCodecRpcOne(t, GoRpc, testCborH, true, 0)
}

// ----- OTHERS -----

func TestCborMapEncodeForCanonical(t *testing.T) {
	doTestMapEncodeForCanonical(t, testCborH)
}

func TestCborSwallowAndZero(t *testing.T) {
	doTestSwallowAndZero(t, testCborH)
}

func TestCborRawExt(t *testing.T) {
	doTestRawExt(t, testCborH)
}

func TestCborMapStructKey(t *testing.T) {
	doTestMapStructKey(t, testCborH)
}

func TestCborDecodeNilMapValue(t *testing.T) {
	doTestDecodeNilMapValue(t, testCborH)
}

func TestCborEmbeddedFieldPrecedence(t *testing.T) {
	doTestEmbeddedFieldPrecedence(t, testCborH)
}

func TestCborLargeContainerLen(t *testing.T) {
	if testCborH.IndefiniteLength {
		t.Skipf("skipping as cbor Indefinite Length doesn't use prefixed lengths")
	}
	doTestLargeContainerLen(t, testCborH)
}

func TestCborTime(t *testing.T) {
	doTestTime(t, testCborH)
}

func TestCborUintToInt(t *testing.T) {
	doTestUintToInt(t, testCborH)
}

func TestCborDifferentMapOrSliceType(t *testing.T) {
	doTestDifferentMapOrSliceType(t, testCborH)
}

func TestCborScalars(t *testing.T) {
	doTestScalars(t, testCborH)
}

func TestCborOmitempty(t *testing.T) {
	doTestOmitempty(t, testCborH)
}

func TestCborIntfMapping(t *testing.T) {
	doTestIntfMapping(t, testCborH)
}

func TestCborMissingFields(t *testing.T) {
	doTestMissingFields(t, testCborH)
}

func TestCborMaxDepth(t *testing.T) {
	doTestMaxDepth(t, testCborH)
}

func TestCborSelfExt(t *testing.T) {
	doTestSelfExt(t, testCborH)
}

func TestCborBytesEncodedAsArray(t *testing.T) {
	doTestBytesEncodedAsArray(t, testCborH)
}

func TestCborStrucEncDec(t *testing.T) {
	doTestStrucEncDec(t, testCborH)
}

func TestCborRawToStringToRawEtc(t *testing.T) {
	doTestRawToStringToRawEtc(t, testCborH)
}

func TestCborStructKeyType(t *testing.T) {
	doTestStructKeyType(t, testCborH)
}

func TestCborPreferArrayOverSlice(t *testing.T) {
	doTestPreferArrayOverSlice(t, testCborH)
}

func TestCborZeroCopyBytes(t *testing.T) {
	if testCborH.IndefiniteLength {
		t.Skipf("skipping ... zero copy bytes not supported by cbor handle with IndefiniteLength=true")
	}
	doTestZeroCopyBytes(t, testCborH)
}

func TestCborNextValueBytes(t *testing.T) {
	doTestNextValueBytes(t, testCborH)
}

func TestCborNumbers(t *testing.T) {
	doTestNumbers(t, testCborH)
}

func TestCborDesc(t *testing.T) {
	m := make(map[byte]string)
	for k, v := range cbordescMajorNames {
		// if k == cborMajorSimpleOrFloat { m[k<<5] = "nil" }
		m[k<<5] = v
	}
	for k, v := range cbordescSimpleNames {
		m[k] = v
	}
	delete(m, cborMajorSimpleOrFloat<<5)
	doTestDesc(t, testCborH, m)
}

func TestCborStructFieldInfoToArray(t *testing.T) {
	doTestStructFieldInfoToArray(t, testCborH)
}
func TestCborAllErrWriter(t *testing.T) {
	// doTestAllErrWriter(t, testCborH, testJsonH)
	doTestAllErrWriter(t, testCborH)
}

func TestCborAllEncCircularRef(t *testing.T) {
	doTestEncCircularRef(t, testCborH)
}

func TestCborAllAnonCycle(t *testing.T) {
	doTestAnonCycle(t, testCborH)
}
