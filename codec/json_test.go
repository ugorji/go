// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ----

func doTestJsonDecodeNonStringScalarInStringContext(t *testing.T, h Handle) {
	defer testSetup2(t, &h)()
	var b = `{"s.true": "true", "b.true": true, "s.false": "false", "b.false": false, "s.10": "10", "i.10": 10, "i.-10": -10}`
	var golden = map[string]string{"s.true": "true", "b.true": "true", "s.false": "false", "b.false": "false", "s.10": "10", "i.10": "10", "i.-10": "-10"}

	var m map[string]string
	d := NewDecoderBytes([]byte(b), h)
	d.MustDecode(&m)
	if err := testEqual(golden, m); err == nil {
		if testv.Verbose {
			t.Logf("++++ match: decoded: %#v", m)
		}
	} else {
		t.Logf("---- mismatch: %v ==> golden: %#v, decoded: %#v", err, golden, m)
		t.FailNow()
	}

	jh := h.(*JsonHandle)

	defer func(a, b, c, d bool, e reflect.Type) {
		jh.MapKeyAsString, jh.PreferFloat, jh.SignedInteger, jh.HTMLCharsAsIs, jh.MapType = a, b, c, d, e
	}(jh.MapKeyAsString, jh.PreferFloat, jh.SignedInteger, jh.HTMLCharsAsIs, jh.MapType)

	jh.MapType = mapIntfIntfTyp
	jh.HTMLCharsAsIs = false
	jh.MapKeyAsString = true
	jh.PreferFloat = false
	jh.SignedInteger = false

	// Also test out decoding string values into naked interface
	b = `{"true": true, "false": false, null: null, "7700000000000000000": 7700000000000000000}`
	const num = 7700000000000000000 // 77
	var golden2 = map[interface{}]interface{}{
		true:  true,
		false: false,
		nil:   nil,
		// uint64(num):uint64(num),
	}

	fn := func() {
		d.ResetBytes([]byte(b))
		var mf interface{}
		d.MustDecode(&mf)
		if err := testEqual(golden2, mf); err != nil {
			t.Logf("---- mismatch: %v ==> golden: %#v, decoded: %#v", err, golden2, mf)
			t.FailNow()
		}
	}

	golden2[uint64(num)] = uint64(num)
	fn()
	delete(golden2, uint64(num))

	jh.SignedInteger = true
	golden2[int64(num)] = int64(num)
	fn()
	delete(golden2, int64(num))
	jh.SignedInteger = false

	jh.PreferFloat = true
	golden2[float64(num)] = float64(num)
	fn()
	delete(golden2, float64(num))
	jh.PreferFloat = false
}

func doTestJsonEncodeIndent(t *testing.T, h Handle) {
	defer testSetup2(t, &h)()
	v := TestSimplish{
		Ii: -794,
		Ss: `A Man is
after the new line
	after new line and tab
`,
	}
	v2 := v
	v.Mm = make(map[string]*TestSimplish)
	for i := 0; i < len(v.Ar); i++ {
		v3 := v2
		v3.Ii += (i * 4)
		v3.Ss = fmt.Sprintf("%d - %s", v3.Ii, v3.Ss)
		if i%2 == 0 {
			v.Ar[i] = &v3
		}
		// v3 = v2
		v.Sl = append(v.Sl, &v3)
		v.Mm[strconv.FormatInt(int64(i), 10)] = &v3
	}
	jh := h.(*JsonHandle)

	defer func(a, b, c bool, d int8) {
		jh.Canonical, jh.StructToArray, jh.NilCollectionToZeroLength, jh.Indent = a, b, c, d
	}(jh.Canonical, jh.StructToArray, jh.NilCollectionToZeroLength, jh.Indent)
	jh.Canonical = true
	jh.Indent = -1
	jh.StructToArray = false
	jh.NilCollectionToZeroLength = false
	var bs []byte
	NewEncoderBytes(&bs, jh).MustEncode(&v)
	txt1Tab := string(bs)
	bs = nil
	jh.Indent = 120
	NewEncoderBytes(&bs, jh).MustEncode(&v)
	txtSpaces := string(bs)

	goldenResultTab := `{
	"Ar": [
		{
			"Ar": [
				null,
				null
			],
			"Ii": -794,
			"Mm": null,
			"Sl": null,
			"Ss": "-794 - A Man is\nafter the new line\n\tafter new line and tab\n"
		},
		null
	],
	"Ii": -794,
	"Mm": {
		"0": {
			"Ar": [
				null,
				null
			],
			"Ii": -794,
			"Mm": null,
			"Sl": null,
			"Ss": "-794 - A Man is\nafter the new line\n\tafter new line and tab\n"
		},
		"1": {
			"Ar": [
				null,
				null
			],
			"Ii": -790,
			"Mm": null,
			"Sl": null,
			"Ss": "-790 - A Man is\nafter the new line\n\tafter new line and tab\n"
		}
	},
	"Sl": [
		{
			"Ar": [
				null,
				null
			],
			"Ii": -794,
			"Mm": null,
			"Sl": null,
			"Ss": "-794 - A Man is\nafter the new line\n\tafter new line and tab\n"
		},
		{
			"Ar": [
				null,
				null
			],
			"Ii": -790,
			"Mm": null,
			"Sl": null,
			"Ss": "-790 - A Man is\nafter the new line\n\tafter new line and tab\n"
		}
	],
	"Ss": "A Man is\nafter the new line\n\tafter new line and tab\n"
}`

	if txt1Tab != goldenResultTab {
		t.Logf("decoded indented with tabs != expected: \nexpected: %s\nencoded: %s", goldenResultTab, txt1Tab)
		t.FailNow()
	}
	if txtSpaces != strings.Replace(goldenResultTab, "\t", strings.Repeat(" ", 120), -1) {
		t.Logf("decoded indented with spaces != expected: \nexpected: %s\nencoded: %s", goldenResultTab, txtSpaces)
		t.FailNow()
	}
}

func doTestJsonLargeInteger(t *testing.T, h Handle) {
	if !testRecoverPanicToErr {
		t.Skip(testSkipIfNotRecoverPanicToErrMsg)
	}
	defer testSetup2(t, &h)()
	jh := h.(*JsonHandle)
	for _, i := range []uint8{'L', 'A', 0} {
		for _, j := range []interface{}{
			int64(1 << 60),
			-int64(1 << 60),
			0,
			1 << 20,
			-(1 << 20),
			uint64(1 << 60),
			uint(0),
			uint(1 << 20),
			int64(1840000e-2),
			uint64(1840000e+2),
		} {
			__doTestJsonLargeInteger(t, j, i, jh)
		}
	}

	if oldIAS := jh.IntegerAsString; oldIAS != 0 {
		jh.IntegerAsString = 0
		defer func() { jh.IntegerAsString = oldIAS }()
	}

	type tt struct {
		s           string
		canI, canUi bool
		i           int64
		ui          uint64
	}

	var i int64
	var ui uint64
	var err error
	var d *Decoder = NewDecoderBytes(nil, jh)
	for _, v := range []tt{
		{"0", true, true, 0, 0},
		{"0000", false, false, 0, 0},
		{"0.00e+2", true, true, 0, 0},
		{"000e-2", false, false, 0, 0},
		{"0.00e-2", true, true, 0, 0},

		{"9223372036854775807", true, true, math.MaxInt64, math.MaxInt64},                             // maxint64
		{"92233720368547758.07e+2", true, true, math.MaxInt64, math.MaxInt64},                         // maxint64
		{"922337203685477580700e-2", true, true, math.MaxInt64, math.MaxInt64},                        // maxint64
		{"9223372.036854775807E+12", true, true, math.MaxInt64, math.MaxInt64},                        // maxint64
		{"9223372036854775807000000000000E-12", true, true, math.MaxInt64, math.MaxInt64},             // maxint64
		{"0.9223372036854775807E+19", true, true, math.MaxInt64, math.MaxInt64},                       // maxint64
		{"92233720368547758070000000000000000000E-19", true, true, math.MaxInt64, math.MaxInt64},      // maxint64
		{"0.000009223372036854775807E+24", true, true, math.MaxInt64, math.MaxInt64},                  // maxint64
		{"9223372036854775807000000000000000000000000E-24", true, true, math.MaxInt64, math.MaxInt64}, // maxint64

		{"-9223372036854775808", true, false, math.MinInt64, 0},                             // minint64
		{"-92233720368547758.08e+2", true, false, math.MinInt64, 0},                         // minint64
		{"-922337203685477580800E-2", true, false, math.MinInt64, 0},                        // minint64
		{"-9223372.036854775808e+12", true, false, math.MinInt64, 0},                        // minint64
		{"-9223372036854775808000000000000E-12", true, false, math.MinInt64, 0},             // minint64
		{"-0.9223372036854775808e+19", true, false, math.MinInt64, 0},                       // minint64
		{"-92233720368547758080000000000000000000E-19", true, false, math.MinInt64, 0},      // minint64
		{"-0.000009223372036854775808e+24", true, false, math.MinInt64, 0},                  // minint64
		{"-9223372036854775808000000000000000000000000E-24", true, false, math.MinInt64, 0}, // minint64

		{"18446744073709551615", false, true, 0, math.MaxUint64},                             // maxuint64
		{"18446744.073709551615E+12", false, true, 0, math.MaxUint64},                        // maxuint64
		{"18446744073709551615000000000000E-12", false, true, 0, math.MaxUint64},             // maxuint64
		{"0.000018446744073709551615E+24", false, true, 0, math.MaxUint64},                   // maxuint64
		{"18446744073709551615000000000000000000000000E-24", false, true, 0, math.MaxUint64}, // maxuint64

		// Add test for limit of uint64 where last digit is 0
		{"18446744073709551610", false, true, 0, math.MaxUint64 - 5},                             // maxuint64
		{"18446744.073709551610E+12", false, true, 0, math.MaxUint64 - 5},                        // maxuint64
		{"18446744073709551610000000000000E-12", false, true, 0, math.MaxUint64 - 5},             // maxuint64
		{"0.000018446744073709551610E+24", false, true, 0, math.MaxUint64 - 5},                   // maxuint64
		{"18446744073709551610000000000000000000000000E-24", false, true, 0, math.MaxUint64 - 5}, // maxuint64
		// {"", true, true},
	} {
		if v.s == "" {
			continue
		}
		d.ResetBytes([]byte(v.s))
		err = d.Decode(&ui)
		if (v.canUi && err != nil) || (!v.canUi && err == nil) || (v.canUi && err == nil && v.ui != ui) {
			t.Logf("Failing to decode %s (as unsigned): %v", v.s, err)
			t.FailNow()
		}

		d.ResetBytes([]byte(v.s))
		err = d.Decode(&i)
		if (v.canI && err != nil) || (!v.canI && err == nil) || (v.canI && err == nil && v.i != i) {
			t.Logf("Failing to decode %s (as signed): %v", v.s, err)
			t.FailNow()
		}
	}
}

func doTestJsonInvalidUnicode(t *testing.T, h Handle) {
	if !testRecoverPanicToErr {
		t.Skip(testSkipIfNotRecoverPanicToErrMsg)
	}
	defer testSetup(t, &h)()
	// t.Skipf("new json implementation does not handle bad unicode robustly")
	jh := h.(*JsonHandle)

	var err error

	// ---- test unmarshal ---
	var m = map[string]string{
		`"\udc49\u0430abc"`: "\uFFFDabc",
		`"\udc49\u0430"`:    "\uFFFD",
		`"\udc49abc"`:       "\uFFFDabc",
		`"\udc49"`:          "\uFFFD",
		`"\udZ49\u0430abc"`: "\uFFFD\u0430abc",
		`"\udcG9\u0430"`:    "\uFFFD\u0430",
		`"\uHc49abc"`:       "\uFFFDabc",
		`"\uKc49"`:          "\uFFFD",
	}

	for k, v := range m {
		// call testUnmarshal directly, so we can check for EOF
		// testUnmarshalErr(&s, []byte(k), jh, t, "-")
		// t.Logf("%s >> %s", k, v)
		var s string
		err = testUnmarshal(&s, []byte(k), jh)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				continue
			}
			t.Logf("%s: unmarshal failed: %v", "-", err)
			t.FailNow()
		}

		if s != v {
			t.Logf("unmarshal: not equal: %q, %q", v, s)
			t.FailNow()
		}
	}

	// test some valid edge cases
	m = map[string]string{
		`"az\uD834\udD1E"`: "az𝄞",
		`"n\ud834\uDD1en"`: "n\U0001D11En", // "\uf09d849e", // "\UD834DD1E" // U+1DD1E g clef

		`"a\\\"\/\"\b\f\n\r\"\tz"`: "a\\\"/\"\b\f\n\r\"\tz",
	}
	for k, v := range m {
		var s string
		err = testUnmarshal(&s, []byte(k), jh)
		if err != nil {
			t.Logf("%s: unmarshal failed: %v", "-", err)
			t.FailNow()
		}

		if s != v {
			t.Logf("unmarshal: not equal: %q, %q", v, s)
			t.FailNow()
		}
	}

	// ---- test marshal ---
	var b = []byte{'"', 0xef, 0xbf, 0xbd} // this is " and unicode.ReplacementChar (as bytes)
	var m2 = map[string][]byte{
		string([]byte{0xef, 0xbf, 0xbd}):           append(b, '"'),
		string([]byte{0xef, 0xbf, 0xbd, 0x0, 0x0}): append(b, `\u0000\u0000"`...),

		"a\\\"/\"\b\f\n\r\"\tz": []byte(`"a\\\"/\"\b\f\n\r\"\tz"`),

		// our encoder doesn't support encoding using only ascii ... so need to use the utf-8 version
		// "n\U0001D11En": []byte(`"n\uD834\uDD1En"`),
		// "az𝄞": []byte(`"az\uD834\uDD1E"`),
		"n\U0001D11En": []byte(`"n𝄞n"`),
		"az𝄞":          []byte(`"az𝄞"`),

		string([]byte{129, 129}): []byte(`"\uFFFD\uFFFD"`),
	}

	for k, v := range m2 {
		b, err = testMarshal(k, jh)
		if err != nil {
			t.Logf("%s: marshal failed: %v", "-", err)
			t.FailNow()
		}

		if !bytes.Equal(b, v) {
			t.Logf("marshal: not equal: %q, %q", v, b)
			t.FailNow()
		}
		testReleaseBytes(b)
	}
}

func doTestJsonNumberParsing(t *testing.T, h Handle) {
	defer testSetup(t, &h)()
	parseInt64_reader := func(r readFloatResult) (v int64, fail bool) {
		u, fail := parseUint64_reader(r)
		if fail {
			return
		}
		if r.neg {
			v = -int64(u)
		} else {
			v = int64(u)
		}
		return
	}
	for _, f64 := range testFloatsToParse {
		// using large prec might make a non-exact float ..., while small prec might make lose some precision.
		// However, we can do a check, and only check if the (u)int64 and float64 can be converted to one another
		// without losing any precision. Then we can use any precision for the tests.
		var precs = [...]int{32, -1}
		var r readFloatResult
		var fail bool
		var fs uint64
		var fsi int64
		var fint, ffrac float64
		_ = ffrac
		var bv []byte
		for _, prec := range precs {
			f := f64
			if math.IsNaN(f64) || math.IsInf(f64, 0) {
				goto F32
			}
			bv = strconv.AppendFloat(nil, f, 'E', prec, 64)
			if fs, err := parseFloat64_custom(bv); err != nil || fs != f {
				t.Logf("float64 -> float64 error (prec: %v) parsing '%s', got %v, expected %v: %v", prec, bv, fs, f, err)
				t.FailNow()
			}
			// try decoding it a uint64 or int64
			fint, ffrac = math.Modf(f)
			bv = strconv.AppendFloat(bv[:0], fint, 'E', prec, 64)
			if f < 0 || fint != float64(uint64(fint)) {
				goto F64i
			}
			r = readFloat(bv, fi64u)
			fail = !r.ok
			if r.ok {
				fs, fail = parseUint64_reader(r)
			}
			if fail || fs != uint64(fint) {
				t.Logf("float64 -> uint64 error (prec: %v, fail: %v) parsing '%s', got %v, expected %v", prec, fail, bv, fs, uint64(fint))
				t.FailNow()
			}
		F64i:
			if fint != float64(int64(fint)) {
				goto F32
			}
			r = readFloat(bv, fi64u)
			fail = !r.ok
			if r.ok {
				fsi, fail = parseInt64_reader(r)
			}
			if fail || fsi != int64(fint) {
				t.Logf("float64 -> int64 error (prec: %v, fail: %v) parsing '%s', got %v, expected %v", prec, fail, bv, fsi, int64(fint))
				t.FailNow()
			}
		F32:
			f32 := float32(f64)
			f64n := float64(f32)
			if !(math.IsNaN(f64n) || math.IsInf(f64n, 0)) {
				bv := strconv.AppendFloat(nil, f64n, 'E', prec, 32)
				if fs, err := parseFloat32_custom(bv); err != nil || fs != f32 {
					t.Logf("float32 -> float32 error (prec: %v) parsing '%s', got %v, expected %v: %v", prec, bv, fs, f32, err)
					t.FailNow()
				}
			}
		}
	}

	for _, v := range testUintsToParse {
		const prec int = 32    // test with precisions of 32 so formatting doesn't truncate what doesn't fit into float
		v = uint64(float64(v)) // use a value that when converted back and forth, remains the same
		bv := strconv.AppendFloat(nil, float64(v), 'E', prec, 64)
		r := readFloat(bv, fi64u)
		if r.ok {
			if fs, fail := parseUint64_reader(r); fail || fs != v {
				t.Logf("uint64 -> uint64 error (prec: %v, fail: %v) parsing '%s', got %v, expected %v", prec, fail, bv, fs, v)
				t.FailNow()
			}
		} else {
			t.Logf("uint64 -> uint64 not ok (prec: %v) parsing '%s', expected %v", prec, bv, v)
			t.FailNow()
		}
		vi := int64(v)
		bv = strconv.AppendFloat(nil, float64(vi), 'E', prec, 64)
		r = readFloat(bv, fi64u)
		if r.ok {
			if fs, fail := parseInt64_reader(r); fail || fs != vi {
				t.Logf("int64 -> int64 error (prec: %v, fail: %v) parsing '%s', got %v, expected %v", prec, fail, bv, fs, vi)
				t.FailNow()
			}
		} else {
			t.Logf("int64 -> int64 not ok (prec: %v) parsing '%s', expected %v", prec, bv, vi)
			t.FailNow()
		}
	}
	{
		var f64 float64 = 64.0
		var f32 float32 = 32.0
		var i int = 11
		var u64 uint64 = 128

		for _, s := range []string{`""`, `null`} {
			b := []byte(s)
			testUnmarshalErr(&f64, b, h, t, "")
			testDeepEqualErr(f64, float64(0), t, "")
			testUnmarshalErr(&f32, b, h, t, "")
			testDeepEqualErr(f32, float32(0), t, "")
			testUnmarshalErr(&i, b, h, t, "")
			testDeepEqualErr(i, int(0), t, "")
			testUnmarshalErr(&u64, b, h, t, "")
			testDeepEqualErr(u64, uint64(0), t, "")
			testUnmarshalErr(&u64, b, h, t, "")
			testDeepEqualErr(u64, uint64(0), t, "")
		}
	}
}

func __doTestJsonLargeInteger(t *testing.T, v interface{}, ias uint8, jh *JsonHandle) {
	if testv.Verbose {
		t.Logf("Running TestJsonLargeInteger: v: %#v, ias: %c", v, ias)
	}

	if oldIAS := jh.IntegerAsString; oldIAS != ias {
		jh.IntegerAsString = ias
		defer func() { jh.IntegerAsString = oldIAS }()
	}

	var vu uint
	var vi int
	var vb bool
	var b []byte
	e := NewEncoderBytes(&b, jh)
	e.MustEncode(v)
	e.MustEncode(true)
	d := NewDecoderBytes(b, jh)
	// below, we validate that the json string or number was encoded,
	// then decode, and validate that the correct value was decoded.
	fnStrChk := func() {
		// check that output started with '"', and ended with '"true'
		// if !(len(b) >= 7 && b[0] == '"' && string(b[len(b)-7:]) == `" true `) {
		if !(len(b) >= 5 && b[0] == '"' && string(b[len(b)-5:]) == `"true`) {
			t.Logf("Expecting a JSON string, got: '%s'", b)
			t.FailNow()
		}
	}

	switch ias {
	case 'L':
		switch v2 := v.(type) {
		case int:
			v2n := int64(v2) // done to work with 32-bit OS
			if v2n > 1<<53 || (v2n < 0 && -v2n > 1<<53) {
				fnStrChk()
			}
		case uint:
			v2n := uint64(v2) // done to work with 32-bit OS
			if v2n > 1<<53 {
				fnStrChk()
			}
		}
	case 'A':
		fnStrChk()
	default:
		// check that output doesn't contain " at all
		for _, i := range b {
			if i == '"' {
				t.Logf("Expecting a JSON Number without quotation: got: %s", b)
				t.FailNow()
			}
		}
	}
	switch v2 := v.(type) {
	case int:
		d.MustDecode(&vi)
		d.MustDecode(&vb)
		// check that vb = true, and vi == v2
		if !(vb && vi == v2) {
			t.Logf("Expecting equal values from %s: got golden: %v, decoded: %v", b, v2, vi)
			t.FailNow()
		}
	case uint:
		d.MustDecode(&vu)
		d.MustDecode(&vb)
		// check that vb = true, and vi == v2
		if !(vb && vu == v2) {
			t.Logf("Expecting equal values from %s: got golden: %v, decoded: %v", b, v2, vu)
			t.FailNow()
		}
	}
}

func doTestJsonTimeAndBytesOptions(t *testing.T, h Handle) {
	defer testSetup2(t, &h)()
	jh := h.(*JsonHandle)
	defer func(a, b []string, c InterfaceExt, d bool, e int8) {
		jh.TimeFormat, jh.BytesFormat, jh.RawBytesExt, jh.StructToArray, jh.Indent = a, b, c, d, e
	}(jh.TimeFormat, jh.BytesFormat, jh.RawBytesExt, jh.StructToArray, jh.Indent)

	jh.RawBytesExt = nil
	jh.Indent = 0
	jh.StructToArray = false

	type params struct {
		TimeFormat  []string
		BytesFormat []string
		Expected    string
	}
	var tt struct {
		Bytes []byte
		Time  time.Time
	}
	tt0 := tt

	tt.Bytes = bytesView(strRpt(1, testLongSentence))
	// func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
	// tt.Time = time.Now().Round(time.Second).UTC()
	tt.Time = time.Date(2025, 6, 19, 14, 32, 49, 0, time.UTC)

	var bs []byte

	table := []params{
		{[]string{}, []string{}, `{"Bytes":"c29tZSByZWFsbHkgcmVhbGx5IGNvb2wgbmFtZXMgdGhhdCBhcmUgbmlnZXJpYW4gYW5kIGFtZXJpY2FuIGxpa2UgInVnb3JqaSBtZWxvZHkgbndva2UiIC0gZ2V0IGl0PyA=","Time":"2025-06-19T14:32:49Z"}`},
		{[]string{time.RFC3339, time.RFC1123Z}, []string{"array"}, `{"Bytes":[115,111,109,101,32,114,101,97,108,108,121,32,114,101,97,108,108,121,32,99,111,111,108,32,110,97,109,101,115,32,116,104,97,116,32,97,114,101,32,110,105,103,101,114,105,97,110,32,97,110,100,32,97,109,101,114,105,99,97,110,32,108,105,107,101,32,34,117,103,111,114,106,105,32,109,101,108,111,100,121,32,110,119,111,107,101,34,32,45,32,103,101,116,32,105,116,63,32],"Time":"Thu, 19 Jun 2025 14:32:49 +0000"}`},
		{[]string{time.RFC1123Z, time.RFC850}, []string{"base32", "hex"}, `{"Bytes":"ONXW2ZJAOJSWC3DMPEQHEZLBNRWHSIDDN5XWYIDOMFWWK4ZAORUGC5BAMFZGKIDONFTWK4TJMFXCAYLOMQQGC3LFOJUWGYLOEBWGS23FEARHKZ3POJVGSIDNMVWG6ZDZEBXHO33LMURCALJAM5SXIIDJOQ7SA===","Time":"Thu, 19 Jun 2025 14:32:49 +0000"}`},
		{[]string{"unixmilli", time.RFC1123Z}, []string{"hex", "base32hex"}, `{"Bytes":"736f6d65207265616c6c79207265616c6c7920636f6f6c206e616d6573207468617420617265206e6967657269616e20616e6420616d65726963616e206c696b65202275676f726a69206d656c6f6479206e776f6b6522202d206765742069743f20","Time":1750343569000}`},
	}
	for i, v := range table {
		jh.TimeFormat = v.TimeFormat
		jh.BytesFormat = v.BytesFormat
		name := fmt.Sprintf("%d", i)
		bs = testMarshalErr(&tt, h, t, name)
		testDeepEqualErr(stringView(bs), v.Expected, t, name)
		tt1 := tt0
		testUnmarshalErr(&tt1, bs, h, t, name)
		tt1.Time = tt1.Time.Round(time.Second).UTC()
		testDeepEqualErr(tt1, tt, t, name)
	}
}

func TestJsonCodecsTable(t *testing.T) {
	doTestCodecTableOne(t, testJsonH)
}

func TestJsonCodecsMisc(t *testing.T) {
	doTestCodecMiscOne(t, testJsonH)
}

func TestJsonCodecsEmbeddedPointer(t *testing.T) {
	doTestCodecEmbeddedPointer(t, testJsonH)
}

func TestJsonCodecChan(t *testing.T) {
	doTestCodecChan(t, testJsonH)
}

func TestJsonStdEncIntf(t *testing.T) {
	doTestStdEncIntf(t, testJsonH)
}

// ----- Raw ---------
func TestJsonRaw(t *testing.T) {
	doTestRawValue(t, testJsonH)
}

func TestJsonRpcGo(t *testing.T) {
	doTestCodecRpcOne(t, GoRpc, testJsonH, true, 0)
}

func TestJsonMapEncodeForCanonical(t *testing.T) {
	doTestMapEncodeForCanonical(t, testJsonH)
}

func TestJsonSwallowAndZero(t *testing.T) {
	doTestSwallowAndZero(t, testJsonH)
}

func TestJsonRawExt(t *testing.T) {
	doTestRawExt(t, testJsonH)
}

func TestJsonMapStructKey(t *testing.T) {
	doTestMapStructKey(t, testJsonH)
}

func TestJsonDecodeNilMapValue(t *testing.T) {
	doTestDecodeNilMapValue(t, testJsonH)
}

func TestJsonEmbeddedFieldPrecedence(t *testing.T) {
	doTestEmbeddedFieldPrecedence(t, testJsonH)
}

func TestJsonLargeContainerLen(t *testing.T) {
	t.Skipf("skipping as json doesn't support prefixed lengths")
	doTestLargeContainerLen(t, testJsonH)
}
func TestJsonTime(t *testing.T) {
	doTestTime(t, testJsonH)
}

func TestJsonUintToInt(t *testing.T) {
	doTestUintToInt(t, testJsonH)
}

func TestJsonDifferentMapOrSliceType(t *testing.T) {
	doTestDifferentMapOrSliceType(t, testJsonH)
}

func TestJsonScalars(t *testing.T) {
	doTestScalars(t, testJsonH)
}

func TestJsonOmitempty(t *testing.T) {
	doTestOmitempty(t, testJsonH)
}

func TestJsonIntfMapping(t *testing.T) {
	doTestIntfMapping(t, testJsonH)
}

func TestJsonMissingFields(t *testing.T) {
	doTestMissingFields(t, testJsonH)
}

func TestJsonMaxDepth(t *testing.T) {
	doTestMaxDepth(t, testJsonH)
}

func TestJsonSelfExt(t *testing.T) {
	doTestSelfExt(t, testJsonH)
}

func TestJsonBytesEncodedAsArray(t *testing.T) {
	doTestBytesEncodedAsArray(t, testJsonH)
}

func TestJsonStrucEncDec(t *testing.T) {
	doTestStrucEncDec(t, testJsonH)
}

func TestJsonRawToStringToRawEtc(t *testing.T) {
	doTestRawToStringToRawEtc(t, testJsonH)
}

func TestJsonStructKeyType(t *testing.T) {
	doTestStructKeyType(t, testJsonH)
}

func TestJsonPreferArrayOverSlice(t *testing.T) {
	doTestPreferArrayOverSlice(t, testJsonH)
}

func TestJsonZeroCopyBytes(t *testing.T) {
	t.Skipf("skipping ... zero copy bytes not supported by json handle")
	doTestZeroCopyBytes(t, testJsonH)
}

func TestJsonNextValueBytes(t *testing.T) {
	doTestNextValueBytes(t, testJsonH)
}

func TestJsonNumbers(t *testing.T) {
	doTestNumbers(t, testJsonH)
}

func TestJsonDesc(t *testing.T) {
	doTestDesc(t, testJsonH, map[byte]string{'"': `"`, '{': `{`, '}': `}`, '[': `[`, ']': `]`})
}

func TestJsonStructFieldInfoToArray(t *testing.T) {
	doTestStructFieldInfoToArray(t, testJsonH)
}

func TestJsonEncodeIndent(t *testing.T) {
	doTestJsonEncodeIndent(t, testJsonH)
}

func TestJsonDecodeNonStringScalarInStringContext(t *testing.T) {
	doTestJsonDecodeNonStringScalarInStringContext(t, testJsonH)
}

func TestJsonLargeInteger(t *testing.T) {
	doTestJsonLargeInteger(t, testJsonH)
}

func TestJsonInvalidUnicode(t *testing.T) {
	doTestJsonInvalidUnicode(t, testJsonH)
}

func TestJsonNumberParsing(t *testing.T) {
	doTestJsonNumberParsing(t, testJsonH)
}

func TestJsonMultipleEncDec(t *testing.T) {
	doTestMultipleEncDec(t, testJsonH)
}

func TestJsonAllErrWriter(t *testing.T) {
	doTestAllErrWriter(t, testJsonH)
}

func TestJsonTimeAndBytesOptions(t *testing.T) {
	doTestJsonTimeAndBytesOptions(t, testJsonH)
}
