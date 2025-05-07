// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// By default, this json support uses base64 encoding for bytes, because you cannot
// store and read any arbitrary string in json (only unicode).
// However, the user can configre how to encode/decode bytes.
//
// This library specifically supports UTF-8 for encoding and decoding only.
//
// Note that the library will happily encode/decode things which are not valid
// json e.g. a map[int64]string. We do it for consistency. With valid json,
// we will encode and decode appropriately.
// Users can specify their map type if necessary to force it.
//
// We cannot use strconv.(Q|Unq)uote because json quotes/unquotes differently.

import (
	"encoding/base64"
	"errors"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

//--------------------------------

// jsonLits and jsonLitb are defined at the package level,
// so they are guaranteed to be stored efficiently, making
// for better append/string comparison/etc.
//
// (anecdotal evidence from some benchmarking on go 1.20 devel in 20220104)
const jsonLits = `"true"false"null"`

var jsonLitb = []byte(jsonLits)

const (
	jsonLitT = 1
	jsonLitF = 6
	jsonLitN = 12
)

const jsonEncodeUintSmallsString = "" +
	"00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

var jsonEncodeUintSmallsStringBytes = (*[len(jsonEncodeUintSmallsString)]byte)([]byte(jsonEncodeUintSmallsString))

const (
	jsonU4Chk2 = '0'
	jsonU4Chk1 = 'a' - 10
	jsonU4Chk0 = 'A' - 10
)

const (
	// If !jsonValidateSymbols, decoding will be faster, by skipping some checks:
	//   - If we see first character of null, false or true,
	//     do not validate subsequent characters.
	//   - e.g. if we see a n, assume null and skip next 3 characters,
	//     and do not validate they are ull.
	// P.S. Do not expect a significant decoding boost from this.
	jsonValidateSymbols = true

	// jsonEscapeMultiByteUnicodeSep controls whether some unicode characters
	// that are valid json but may bomb in some contexts are escaped during encoeing.
	//
	// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
	// Both technically valid JSON, but bomb on JSONP, so fix here unconditionally.
	jsonEscapeMultiByteUnicodeSep = true

	// jsonRecognizeBoolNullInQuotedStr is used during decoding into a blank interface{}
	// to control whether we detect quoted values of bools and null where a map key is expected,
	// and treat as nil, true or false.
	jsonNakedBoolNullInQuotedStr = true

	jsonSpacesOrTabsLen = 128
)

var (
	// jsonTabs and jsonSpaces are used as caches for indents
	jsonTabs, jsonSpaces [jsonSpacesOrTabsLen]byte
)

func init() {
	for i := byte(0); i < jsonSpacesOrTabsLen; i++ {
		jsonSpaces[i] = ' '
		jsonTabs[i] = '\t'
	}

}

// ----------------

type jsonEncState struct {
	di int8   // indent per: if negative, use tabs
	d  bool   // indenting?
	dl uint16 // indent level
}

// func (x jsonEncState) captureState() interface{}   { return x }
// func (x *jsonEncState) restoreState(v interface{}) { *x = v.(jsonEncState) }

type jsonEncDriver[T encWriter] struct {
	noBuiltInTypes
	h *JsonHandle
	e *encoderBase
	s *bitset256 // safe set for characters (taking h.HTMLAsIs into consideration)

	w T
	// se interfaceExtWrapper

	enc encoderI
	// ---- cpu cache line boundary?
	jsonEncState

	ks bool // map key as string
	is byte // integer as string

	typical bool
	rawext  bool // rawext configured on the handle

	// buf *[]byte // used mostly for encoding []byte

	// scratch buffer for: encode time, numbers, etc
	//
	// RFC3339Nano uses 35 chars: 2006-01-02T15:04:05.999999999Z07:00
	// MaxUint64 uses 20 chars: 18446744073709551615
	// floats are encoded using: f/e fmt, and -1 precision, or 1 if no fractions.
	// This means we are limited by the number of characters for the
	// mantissa (up to 17), exponent (up to 3), signs (up to 3), dot (up to 1), E (up to 1)
	// for a total of 24 characters.
	//    -xxx.yyyyyyyyyyyye-zzz
	// Consequently, 35 characters should be sufficient for encoding time, integers or floats.
	// We use up all the remaining bytes to make this use full cache lines.
	b [48]byte
}

func (e *jsonEncDriver[T]) writeIndent() {
	e.w.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.di < 0 {
		x = -x
		for x > jsonSpacesOrTabsLen {
			e.w.writeb(jsonTabs[:])
			x -= jsonSpacesOrTabsLen
		}
		e.w.writeb(jsonTabs[:x])
	} else {
		for x > jsonSpacesOrTabsLen {
			e.w.writeb(jsonSpaces[:])
			x -= jsonSpacesOrTabsLen
		}
		e.w.writeb(jsonSpaces[:x])
	}
}

func (e *jsonEncDriver[T]) WriteArrayElem() {
	if e.e.c != containerArrayStart {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver[T]) WriteMapElemKey() {
	if e.e.c != containerMapStart {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver[T]) WriteMapElemValue() {
	if e.d {
		e.w.writen2(':', ' ')
	} else {
		e.w.writen1(':')
	}
}

func (e *jsonEncDriver[T]) EncodeNil() {
	// We always encode nil as just null (never in quotes)
	// so we can easily decode if a nil in the json stream ie if initial token is n.

	e.w.writestr(jsonLits[jsonLitN : jsonLitN+4])
}

func (e *jsonEncDriver[T]) EncodeTime(t time.Time) {
	// Do NOT use MarshalJSON, as it allocates internally.
	// instead, we call AppendFormat directly, using our scratch buffer (e.b)

	if t.IsZero() {
		e.EncodeNil()
	} else {
		e.b[0] = '"'
		b := t.AppendFormat(e.b[1:1], time.RFC3339Nano)
		e.b[len(b)+1] = '"'
		e.w.writeb(e.b[:len(b)+2])
	}
}

func (e *jsonEncDriver[T]) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.EncodeNil()
	} else {
		e.enc.encode(v)
	}
}

func (e *jsonEncDriver[T]) EncodeRawExt(re *RawExt) {
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encode(re.Value)
	} else {
		e.EncodeNil()
	}
}

var jsonEncBoolStrs = [2][2]string{
	{jsonLits[jsonLitF : jsonLitF+5], jsonLits[jsonLitT : jsonLitT+4]},
	{jsonLits[jsonLitF-1 : jsonLitF+6], jsonLits[jsonLitT-1 : jsonLitT+5]},
}

func (e *jsonEncDriver[T]) EncodeBool(b bool) {
	e.w.writestr(jsonEncBoolStrs[bool2int(e.ks && e.e.c == containerMapKey)%2][bool2int(b)%2])
}

// func (e *jsonEncDriver[T]) EncodeBool(b bool) {
// 	if e.ks && e.e.c == containerMapKey {
// 		if b {
// 			e.w.writestr(jsonLits[jsonLitT-1 : jsonLitT+5])
// 		} else {
// 			e.w.writestr(jsonLits[jsonLitF-1 : jsonLitF+6])
// 		}
// 	} else {
// 		if b {
// 			e.w.writestr(jsonLits[jsonLitT : jsonLitT+4])
// 		} else {
// 			e.w.writestr(jsonLits[jsonLitF : jsonLitF+5])
// 		}
// 	}
// }

func (e *jsonEncDriver[T]) encodeFloat(f float64, bitsize, fmt byte, prec int8) {
	var blen uint
	if e.ks && e.e.c == containerMapKey {
		blen = 2 + uint(len(strconv.AppendFloat(e.b[1:1], f, fmt, int(prec), int(bitsize))))
		// _ = e.b[:blen]
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
	} else {
		e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), int(bitsize)))
	}
}

func (e *jsonEncDriver[T]) EncodeFloat64(f float64) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.encodeFloat(f, 64, fmt, prec)
}

func (e *jsonEncDriver[T]) EncodeFloat32(f float32) {
	if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.encodeFloat(float64(f), 32, fmt, prec)
}

func jsonEncodeUint(neg, quotes bool, u uint64, b *[48]byte) []byte {
	// copied mostly from std library: strconv
	// this should only be called on 64bit OS.

	// const ss = jsonEncodeUintSmallsString
	var ss = jsonEncodeUintSmallsStringBytes[:]

	// typically, 19 or 20 bytes sufficient for decimal encoding a uint64
	// var a [24]byte
	// var a = (*[24]byte)(b[0:24])
	var a = b[:24]
	var i = uint(len(a))

	if quotes {
		i--
		setByteAt(a, i, '"')
		// a[i] = '"'
	}
	// u guaranteed to fit into a uint (as we are not 32bit OS)
	var is uint
	var us = uint(u)
	for us >= 100 {
		is = us % 100 * 2
		us /= 100
		i -= 2
		setByteAt(a, i+1, byteAt(ss, is+1))
		setByteAt(a, i, byteAt(ss, is))
		// a[i+1] = ss[is+1]
		// a[i] = ss[is]
	}

	// us < 100
	is = us * 2
	i--
	setByteAt(a, i, byteAt(ss, is+1))
	// a[i] = ss[is+1]
	if us >= 10 {
		i--
		setByteAt(a, i, byteAt(ss, is))
		// a[i] = ss[is]
	}
	if neg {
		i--
		setByteAt(a, i, '-')
		// a[i] = '-'
	}
	if quotes {
		i--
		setByteAt(a, i, '"')
		// a[i] = '"'
	}
	return a[i:]
}

func (e *jsonEncDriver[T]) encodeUint(neg bool, quotes bool, u uint64) {
	e.w.writeb(jsonEncodeUint(neg, quotes, u, &e.b))
}

func (e *jsonEncDriver[T]) EncodeInt(v int64) {
	quotes := e.is == 'A' || e.is == 'L' && (v > 1<<53 || v < -(1<<53)) ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		if quotes {
			blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
		}
		return
	}

	if v < 0 {
		e.encodeUint(true, quotes, uint64(-v))
	} else {
		e.encodeUint(false, quotes, uint64(v))
	}
}

func (e *jsonEncDriver[T]) EncodeUint(v uint64) {
	quotes := e.is == 'A' || e.is == 'L' && v > 1<<53 ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		// use strconv directly, as optimized encodeUint only works on 64-bit alone
		if quotes {
			blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
		}
		return
	}

	e.encodeUint(false, quotes, v)
}

func (e *jsonEncDriver[T]) EncodeString(v string) {
	if e.h.StringToRaw {
		e.EncodeStringBytesRaw(bytesView(v))
		return
	}
	e.quoteStr(v)
}

func (e *jsonEncDriver[T]) EncodeStringNoEscape4Json(v string) { e.w.writeqstr(v) }

func (e *jsonEncDriver[T]) EncodeStringBytesRaw(v []byte) {
	// if encoding raw bytes and RawBytesExt is configured, use it to encode
	if v == nil {
		e.EncodeNil()
		return
	}

	if e.rawext {
		// explicitly convert v to interface{} so that v doesn't escape to heap
		iv := e.h.RawBytesExt.ConvertExt(any(v))
		if iv == nil {
			e.EncodeNil()
		} else {
			e.enc.encode(iv)
		}
		return
	}

	slen := base64.StdEncoding.EncodedLen(len(v)) + 2

	// bs := e.e.blist.check(*e.buf, n)[:slen]
	// *e.buf = bs

	bs := e.e.blist.peek(slen, false)
	bs = bs[:slen]

	base64.StdEncoding.Encode(bs[1:], v)
	bs[len(bs)-1] = '"'
	bs[0] = '"'
	e.w.writeb(bs)
}

// indent is done as below:
//   - newline and indent are added before each mapKey or arrayElem
//   - newline and indent are added before each ending,
//     except there was no entry (so we can have {} or [])

func (e *jsonEncDriver[T]) WriteArrayStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('[')
}

func (e *jsonEncDriver[T]) WriteArrayEnd() {
	if e.d {
		e.dl--
		e.writeIndent()
	}
	e.w.writen1(']')
}

func (e *jsonEncDriver[T]) WriteMapStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('{')
}

func (e *jsonEncDriver[T]) WriteMapEnd() {
	if e.d {
		e.dl--
		if e.e.c != containerMapStart {
			e.writeIndent()
		}
	}
	e.w.writen1('}')
}

func (e *jsonEncDriver[T]) quoteStr(s string) {
	// adapted from std pkg encoding/json
	const hex = "0123456789abcdef"
	e.w.writen1('"')
	var i, start uint
	for i < uint(len(s)) {
		// encode all bytes < 0x20 (except \r, \n).
		// also encode < > & to prevent security holes when served to some browsers.

		// We optimize for ascii, by assuming that most characters are in the BMP
		// and natively consumed by json without much computation.

		// if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
		// if (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b) {
		b := s[i]
		if e.s.isset(b) {
			i++
			continue
		}
		if b < utf8.RuneSelf {
			if start < i {
				e.w.writestr(s[start:i])
			}

			// if b == '\\' || b == '"' {
			// 	e.w.writen2('\\', b)
			// } else if b == '\n' {
			// 	e.w.writen2('\\', 'n')
			// } else if b == '\t' {
			// 	e.w.writen2('\\', 't')
			// } else if b == '\r' {
			// 	e.w.writen2('\\', 'r')
			// } else if b == '\b' {
			// 	e.w.writen2('\\', 'b')
			// } else if b == '\f' {
			// 	e.w.writen2('\\', 'f')
			// } else {
			// 	e.w.writestr(`\u00`)
			// 	e.w.writen2(hex[b>>4], hex[b&0xF])
			// }

			switch b {
			case '\\':
				e.w.writen2('\\', '\\')
			case '"':
				e.w.writen2('\\', '"')
			case '\n':
				e.w.writen2('\\', 'n')
			case '\t':
				e.w.writen2('\\', 't')
			case '\r':
				e.w.writen2('\\', 'r')
			case '\b':
				e.w.writen2('\\', 'b')
			case '\f':
				e.w.writen2('\\', 'f')
			default:
				e.w.writestr(`\u00`)
				e.w.writen2(hex[b>>4], hex[b&0xF])
			}

			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 { // meaning invalid encoding (so output as-is)
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\uFFFD`)
			i++
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
		// Both technically valid JSON, but bomb on JSONP, so fix here *unconditionally*.
		if jsonEscapeMultiByteUnicodeSep && (c == '\u2028' || c == '\u2029') {
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\u202`)
			e.w.writen1(hex[c&0xF])
			i += uint(size)
			start = i
			continue
		}
		i += uint(size)
	}
	if start < uint(len(s)) {
		e.w.writestr(s[start:])
	}
	e.w.writen1('"')
}

func (e *jsonEncDriver[T]) atEndOfEncode() {
	if e.h.TermWhitespace {
		var c byte = ' ' // default is that scalar is written, so output space
		if e.e.c != 0 {
			c = '\n' // for containers (map/list), output a newline
		}
		e.w.writen1(c)
	}
}

// ----------

type jsonDecState struct {
	// scratch buffer used for base64 decoding (DecodeBytes in reuseBuf mode),
	// or reading doubleQuoted string (DecodeStringAsBytes, DecodeNaked)
	buf []byte

	rawext bool // rawext configured on the handle

	tok  uint8   // used to store the token read right after skipWhiteSpace
	_    bool    // found null
	_    byte    // padding
	bstr [4]byte // scratch used for string \UXXX parsing
}

// func (x jsonDecState) captureState() interface{}   { return x }
// func (x *jsonDecState) restoreState(v interface{}) { *x = v.(jsonDecState) }

type jsonDecDriver[T decReader] struct {
	noBuiltInTypes
	decDriverNoopNumberHelper
	h *JsonHandle
	d *decoderBase

	r T
	jsonDecState

	// se  interfaceExtWrapper

	// ---- cpu cache line boundary?

	// bytes bool

	dec decoderI
}

func (d *jsonDecDriver[T]) ReadMapStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '{' {
		halt.errorByte("read map - expect char '{' but got char: ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriver[T]) ReadArrayStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '[' {
		halt.errorByte("read array - expect char '[' but got char ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

// MARKER:
// We attempted making sure CheckBreak can be inlined, by moving the skipWhitespace
// call to an explicit (noinline) function call.
// However, this forces CheckBreak to always incur a function call if there was whitespace,
// with no clear benefit.

func (d *jsonDecDriver[T]) CheckBreak() bool {
	d.advance()
	return d.tok == '}' || d.tok == ']'
}

func (d *jsonDecDriver[T]) ReadArrayElem() {
	const xc uint8 = ','
	if d.d.c != containerArrayStart {
		d.advance()
		if d.tok != xc {
			d.readDelimError(xc)
		}
		d.tok = 0
	}
}

func (d *jsonDecDriver[T]) ReadArrayEnd() {
	const xc uint8 = ']'
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriver[T]) ReadMapElemKey() {
	const xc uint8 = ','
	if d.d.c != containerMapStart {
		d.advance()
		if d.tok != xc {
			d.readDelimError(xc)
		}
		d.tok = 0
	}
}

func (d *jsonDecDriver[T]) ReadMapElemValue() {
	const xc uint8 = ':'
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriver[T]) ReadMapEnd() {
	const xc uint8 = '}'
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

//go:inline
func (d *jsonDecDriver[T]) readDelimError(xc uint8) {
	halt.errorf("read json delimiter - expect char '%c' but got char '%c'", xc, d.tok)
}

// MARKER: checkLit takes the readn(3|4) result as a parameter so they can be inlined.
// We pass the array directly to errorf, as passing slice pushes past inlining threshold,
// and passing slice also might cause allocation of the bs array on the heap.

func (d *jsonDecDriver[T]) checkLit3(got, expect [3]byte) {
	d.tok = 0
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr3(got, expect)
	}
}

func (d *jsonDecDriver[T]) checkLit4(got, expect [4]byte) {
	d.tok = 0
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr4(got, expect)
	}
}

// MARKER: checkLitErr methods to prevent the got/expect parameters from escaping

//go:noinline
func jsonCheckLitErr3(got, expect [3]byte) {
	halt.errorf("expecting %s: got %s", expect, got)
}

//go:noinline
func jsonCheckLitErr4(got, expect [4]byte) {
	halt.errorf("expecting %s: got %s", expect, got)
}

func (d *jsonDecDriver[T]) skipWhitespace() {
	d.tok = d.r.skipWhitespace()
}

func (d *jsonDecDriver[T]) advance() {
	if d.tok == 0 {
		d.skipWhitespace()
	}
}

func (d *jsonDecDriver[T]) nextValueBytes() []byte {
	consumeString := func() {
	TOP:
		_, c := d.r.jsonReadAsisChars()
		if c == '\\' { // consume next one and try again
			d.r.readn1()
			goto TOP
		}
	}

	d.advance() // ignore leading whitespace
	d.r.startRecording()

	// cursor = d.d.rb.c - 1 // cursor starts just before non-whitespace token
	switch d.tok {
	default:
		d.r.jsonReadNum()
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
	case '"':
		consumeString()
	case '{', '[':
		var elem struct{}
		var stack []struct{}

		stack = append(stack, elem)

		for len(stack) != 0 {
			c := d.r.readn1()
			switch c {
			case '"':
				consumeString()
			case '{', '[':
				stack = append(stack, elem)
			case '}', ']':
				stack = stack[:len(stack)-1]
			}
		}
	}
	d.tok = 0

	return d.r.stopRecording()
}

func (d *jsonDecDriver[T]) TryNil() bool {
	d.advance()
	// we shouldn't try to see if quoted "null" was here, right?
	// only the plain string: `null` denotes a nil (ie not quotes)
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriver[T]) DecodeBool() (v bool) {
	d.advance()
	// bool can be in quotes if and only if it's a map key
	fquot := d.d.c == containerMapKey && d.tok == '"'
	if fquot {
		d.tok = d.r.readn1()
	}
	switch d.tok {
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		// v = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		v = true
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		// v = false
	default:
		halt.errorByte("decode bool: got first char: ", d.tok)
		// v = false // "unreachable"
	}
	if fquot {
		d.r.readn1()
	}
	return
}

func (d *jsonDecDriver[T]) DecodeTime() (t time.Time) {
	// read string, and pass the string into json.unmarshal
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	d.ensureReadingString()
	bs := d.readUnescapedString()
	t, err := time.Parse(time.RFC3339, stringView(bs))
	halt.onerror(err)
	return
}

func (d *jsonDecDriver[T]) ContainerType() (vt valueType) {
	// check container type by checking the first char
	d.advance()

	// optimize this, so we don't do 4 checks but do one computation.
	// return jsonContainerSet[d.tok]

	// ContainerType is mostly called for Map and Array,
	// so this conditional is good enough (max 2 checks typically)
	if d.tok == '{' {
		return valueTypeMap
	} else if d.tok == '[' {
		return valueTypeArray
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return valueTypeNil
	} else if d.tok == '"' {
		return valueTypeString
	}
	return valueTypeUnset
}

func (d *jsonDecDriver[T]) decNumBytes() (bs []byte) {
	d.advance()
	if d.tok == '"' {
		bs = d.r.jsonReadUntilDblQuote()
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	} else {
		bs = d.r.jsonReadNum()
	}
	d.tok = 0
	return
}

func (d *jsonDecDriver[T]) DecodeUint64() (u uint64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if neg {
		halt.errorStr("negative number cannot be decoded as uint64")
	}
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseUint"))
	}
	return
}

func (d *jsonDecDriver[T]) DecodeInt64() (v int64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseInt"))
	}
	if chkOvf.Uint2Int(u, neg) {
		halt.errorBytes("overflow decoding number from ", b)
	}
	if neg {
		v = -int64(u)
	} else {
		v = int64(u)
	}
	return
}

func (d *jsonDecDriver[T]) DecodeFloat64() (f float64) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat64(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriver[T]) DecodeFloat32() (f float32) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat32(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriver[T]) advanceNil() (ok bool) {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if d.advanceNil() {
		return
	}
	if ext == SelfExt {
		d.dec.decodeAs(rv, basetype, false)
	} else {
		d.dec.interfaceExtConvertAndDecode(rv, ext)
	}
}

func (d *jsonDecDriver[T]) DecodeRawExt(re *RawExt) {
	if d.advanceNil() {
		return
	}
	d.dec.decode(&re.Value)
}

func (d *jsonDecDriver[T]) decBytesFromArray(bs []byte) []byte {
	if bs != nil {
		bs = bs[:0]
	}
	d.tok = 0
	bs = append(bs, uint8(d.DecodeUint64()))
	d.tok = d.r.skipWhitespace() // skip(&whitespaceCharBitset)
	for d.tok != ']' {
		if d.tok != ',' {
			halt.errorByte("read array element - expect char ',' but got char: ", d.tok)
		}
		d.tok = 0
		bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		d.tok = d.r.skipWhitespace() // skip(&whitespaceCharBitset)
	}
	d.tok = 0
	return bs
}

func (d *jsonDecDriver[T]) DecodeBytes() (bs []byte, state dBytesAttachState) {
	d.advance()
	state = dBytesDetach
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	// if decoding into raw bytes, and the RawBytesExt is configured, use it to decode.
	if d.rawext {
		d.dec.interfaceExtConvertAndDecode(&d.buf, d.h.RawBytesExt)
		bs = d.buf
		state = dBytesAttachBuffer
		return
	}
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.tok == '[' {
		// bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		bs = d.decBytesFromArray(d.buf[:0])
		d.buf = bs
		state = dBytesAttachBuffer
		return
	}

	// base64 encodes []byte{} as "", and we encode nil []byte as null.
	// Consequently, base64 should decode null as a nil []byte, and "" as an empty []byte{}.

	state = dBytesAttachBuffer
	d.ensureReadingString()
	bs1 := d.readUnescapedString()
	slen := base64.StdEncoding.DecodedLen(len(bs1))
	if slen == 0 {
		bs = zeroByteSlice
		state = dBytesDetach
	} else if slen <= cap(d.buf) {
		bs = d.buf[:slen]
	} else {
		d.buf = d.d.blist.putGet(d.buf, slen)[:slen]
		bs = d.buf
	}
	slen2, err := base64.StdEncoding.Decode(bs, bs1)
	if err != nil {
		halt.errorf("error decoding base64 binary '%s': %v", any(bs1), err)
	}
	if slen != slen2 {
		bs = bs[:slen2]
	}
	return
}

func (d *jsonDecDriver[T]) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
	d.advance()

	var cond bool
	// common case - hoist outside the switch statement
	if d.tok == '"' {
		bs, cond = d.dblQuoteStringAsBytes()
		state = d.d.attachState(cond)
		return
	}

	state = dBytesDetach
	// handle non-string scalar: null, true, false or a number
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		// out = nil // []byte{}
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		bs = jsonLitb[jsonLitF : jsonLitF+5]
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		bs = jsonLitb[jsonLitT : jsonLitT+4]
	default:
		// try to parse a valid number
		d.tok = 0
		bs = d.r.jsonReadNum()
		state = d.d.attachState(!d.d.bytes)
	}
	return
}

func (d *jsonDecDriver[T]) ensureReadingString() {
	if d.tok != '"' {
		halt.errorByte("expecting string starting with '\"'; got ", d.tok)
	}
}

func (d *jsonDecDriver[T]) readUnescapedString() (bs []byte) {
	// d.ensureReadingString()
	bs = d.r.jsonReadUntilDblQuote()
	d.tok = 0
	return
}

func (d *jsonDecDriver[T]) dblQuoteStringAsBytes() (buf []byte, usingBuf bool) {
	d.tok = 0

	bs, c := d.r.jsonReadAsisChars()
	if c == '"' {
		return bs, false
	}

	checkUtf8 := d.h.ValidateUnicode
	usingBuf = true
	// use a local buf variable, so we don't do pointer chasing within loop
	buf = append(d.buf[:0], bs...)

	for {
		// c is now '\'
		c = d.r.readn1()

		switch c {
		case '"', '\\', '/', '\'':
			buf = append(buf, c)
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'u':
			rr := d.appendStringAsBytesSlashU()
			if checkUtf8 && rr == unicode.ReplacementChar {
				d.buf = buf
				halt.errorBytes("invalid UTF-8 character found after: ", buf)
			}
			buf = append(buf, d.bstr[:utf8.EncodeRune(d.bstr[:], rr)]...)
		default:
			d.buf = buf
			halt.errorByte("unsupported escaped value: ", c)
		}

		bs, c = d.r.jsonReadAsisChars()
		// _ = bs[0] // bounds check hint - slice must be > 0 elements
		buf = append(buf, bs...)
		if c == '"' {
			break
		}
	}
	d.buf = buf
	return
}

func (d *jsonDecDriver[T]) appendStringAsBytesSlashU() (r rune) {
	var rr uint32
	var csu [2]byte
	var cs [4]byte = d.r.readn4()
	if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
		return unicode.ReplacementChar
	}
	r = rune(rr)
	if utf16.IsSurrogate(r) {
		csu = d.r.readn2()
		cs = d.r.readn4()
		if csu[0] == '\\' && csu[1] == 'u' {
			if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
				return unicode.ReplacementChar
			}
			return utf16.DecodeRune(r, rune(rr))
		}
		return unicode.ReplacementChar
	}
	return
}

func jsonSlashURune(cs [4]byte) (rr uint32) {
	for _, c := range cs {
		// best to use explicit if-else
		// - not a table, etc which involve memory loads, array lookup with bounds checks, etc
		if c >= '0' && c <= '9' {
			rr = rr*16 + uint32(c-jsonU4Chk2)
		} else if c >= 'a' && c <= 'f' {
			rr = rr*16 + uint32(c-jsonU4Chk1)
		} else if c >= 'A' && c <= 'F' {
			rr = rr*16 + uint32(c-jsonU4Chk0)
		} else {
			return unicode.ReplacementChar
		}
	}
	return
}

func jsonNakedNum(z *fauxUnion, bs []byte, preferFloat, signedInt bool) (err error) {
	// Note: jsonNakedNum is NEVER called with a zero-length []byte
	if preferFloat {
		z.v = valueTypeFloat
		z.f, err = parseFloat64(bs)
	} else {
		err = parseNumber(bs, z, signedInt)
	}
	return
}

func (d *jsonDecDriver[T]) DecodeNaked() {
	z := d.d.naked()

	d.advance()
	var bs []byte
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		z.v = valueTypeNil
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		z.v = valueTypeBool
		z.b = true
	case '{':
		z.v = valueTypeMap // don't consume. kInterfaceNaked will call ReadMapStart
	case '[':
		z.v = valueTypeArray // don't consume. kInterfaceNaked will call ReadArrayStart
	case '"':
		// if a string, and MapKeyAsString, then try to decode it as a bool or number first
		bs, z.b = d.dblQuoteStringAsBytes()
		att := d.d.attachState(z.b)
		if jsonNakedBoolNullInQuotedStr &&
			d.h.MapKeyAsString && len(bs) > 0 && d.d.c == containerMapKey {
			switch string(bs) {
			// case "null": // nil is never quoted
			// 	z.v = valueTypeNil
			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default:
				// check if a number: float, int or uint
				if err := jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
					z.v = valueTypeString
					z.s = d.d.string(bs, att)
				}
			}
		} else {
			z.v = valueTypeString
			z.s = d.d.string(bs, att)
		}
	default: // number
		bs = d.r.jsonReadNum()
		d.tok = 0
		if len(bs) == 0 {
			halt.errorStr("decode number from empty string")
		}
		if err := jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
			halt.errorf("decode number from %s: %v", any(bs), err)
		}
	}
}

//----------------------

// JsonHandle is a handle for JSON encoding format.
//
// Json is comprehensively supported:
//   - decodes numbers into interface{} as int, uint or float64
//     based on how the number looks and some config parameters e.g. PreferFloat, SignedInt, etc.
//   - decode integers from float formatted numbers e.g. 1.27e+8
//   - decode any json value (numbers, bool, etc) from quoted strings
//   - configurable way to encode/decode []byte .
//     by default, encodes and decodes []byte using base64 Std Encoding
//   - UTF-8 support for encoding and decoding
//
// It has better performance than the json library in the standard library,
// by leveraging the performance improvements of the codec library.
//
// In addition, it doesn't read more bytes than necessary during a decode, which allows
// reading multiple values from a stream containing json and non-json content.
// For example, a user can read a json value, then a cbor value, then a msgpack value,
// all from the same stream in sequence.
//
// Note that, when decoding quoted strings, invalid UTF-8 or invalid UTF-16 surrogate pairs are
// not treated as an error. Instead, they are replaced by the Unicode replacement character U+FFFD.
//
// Note also that the float values for NaN, +Inf or -Inf are encoded as null,
// as suggested by NOTE 4 of the ECMA-262 ECMAScript Language Specification 5.1 edition.
// see http://www.ecma-international.org/publications/files/ECMA-ST/Ecma-262.pdf .
//
// Note the following behaviour differences vs std-library encoding/json package:
//   - struct field names matched in case-sensitive manner
type JsonHandle struct {
	textEncodingType
	BasicHandle

	// Indent indicates how a value is encoded.
	//   - If positive, indent by that number of spaces.
	//   - If negative, indent by that number of tabs.
	Indent int8

	// IntegerAsString controls how integers (signed and unsigned) are encoded.
	//
	// Per the JSON Spec, JSON numbers are 64-bit floating point numbers.
	// Consequently, integers > 2^53 cannot be represented as a JSON number without losing precision.
	// This can be mitigated by configuring how to encode integers.
	//
	// IntegerAsString interpretes the following values:
	//   - if 'L', then encode integers > 2^53 as a json string.
	//   - if 'A', then encode all integers as a json string
	//             containing the exact integer representation as a decimal.
	//   - else    encode all integers as a json number (default)
	IntegerAsString byte

	// HTMLCharsAsIs controls how to encode some special characters to html: < > &
	//
	// By default, we encode them as \uXXX
	// to prevent security holes when served from some browsers.
	HTMLCharsAsIs bool

	// PreferFloat says that we will default to decoding a number as a float.
	// If not set, we will examine the characters of the number and decode as an
	// integer type if it doesn't have any of the characters [.eE].
	PreferFloat bool

	// TermWhitespace says that we add a whitespace character
	// at the end of an encoding.
	//
	// The whitespace is important, especially if using numbers in a context
	// where multiple items are written to a stream.
	TermWhitespace bool

	// MapKeyAsString says to encode all map keys as strings.
	//
	// Use this to enforce strict json output.
	// The only caveat is that nil value is ALWAYS written as null (never as "null")
	MapKeyAsString bool

	// _ uint64 // padding (cache line)

	// Note: below, we store hardly-used items e.g. RawBytesExt.
	// These values below may straddle a cache line, but they are hardly-used,
	// so shouldn't contribute to false-sharing except in rare cases.

	// RawBytesExt, if configured, is used to encode and decode raw bytes in a custom way.
	// If not configured, raw bytes are encoded to/from base64 text.
	RawBytesExt InterfaceExt
}

func (h *JsonHandle) isJson() bool { return true }

// Name returns the name of the handle: json
func (h *JsonHandle) Name() string { return "json" }

// func (h *JsonHandle) desc(bd byte) string { return str4byte(bd) }
func (h *JsonHandle) desc(bd byte) string { return string(bd) }

func (h *JsonHandle) typical() bool {
	return h.Indent == 0 && !h.MapKeyAsString && h.IntegerAsString != 'A' && h.IntegerAsString != 'L'
}

// SetInterfaceExt sets an extension
func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}

// func (h *JsonHandle) newEncDriver() encDriver {
// 	var e = &jsonEncDriver{h: h}
// 	// var x []byte
// 	// e.buf = &x
// 	e.e.e = e
// 	e.e.js = true
// 	e.e.init(h)
// 	e.reset()
// 	return e
// }

// func (h *JsonHandle) newDecDriver() decDriver {
// 	var d = &jsonDecDriver{h: h}
// 	var x []byte
// 	d.buf = &x
// 	d.d.d = d
// 	d.d.js = true
// 	d.d.jsms = h.MapKeyAsString
// 	d.d.init(h)
// 	d.reset()
// 	return d
// }

// func (e *jsonEncDriver[T]) resetState() {
// 	e.dl = 0
// }

func (e *jsonEncDriver[T]) reset() {
	e.dl = 0
	// e.resetState()
	// (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b)
	// cache values from the handle
	e.typical = e.h.typical()
	if e.h.HTMLCharsAsIs {
		e.s = &jsonCharSafeBitset
	} else {
		e.s = &jsonCharHtmlSafeBitset
	}
	e.rawext = e.h.RawBytesExt != nil
	e.di = int8(e.h.Indent)
	e.d = e.h.Indent != 0
	e.ks = e.h.MapKeyAsString
	e.is = e.h.IntegerAsString
}

// func (d *jsonDecDriver[T]) resetState() {
// 	*d.buf = d.d.blist.check(*d.buf, 256)
// 	d.tok = 0
// }

func (d *jsonDecDriver[T]) reset() {
	d.buf = d.d.blist.check(d.buf, 256)
	d.tok = 0
	// d.resetState()
	d.rawext = d.h.RawBytesExt != nil
}

func jsonFloatStrconvFmtPrec64(f float64) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	fbits := math.Float64bits(f)
	abs := math.Float64frombits(fbits &^ (1 << 63))
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac64(fbits) {
		prec = 1
	}
	return
}

func jsonFloatStrconvFmtPrec32(f float32) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	// directly handle Modf (to get fractions) and Abs (to get absolute)
	fbits := math.Float32bits(f)
	abs := math.Float32frombits(fbits &^ (1 << 31))
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac32(fbits) {
		prec = 1
	}
	return
}

// var _ decDriverContainerTracker = (*jsonDecDriver[T])(nil)
// var _ encDriverContainerTracker = (*jsonEncDriver[T])(nil)
// var _ decDriver = (*jsonDecDriver[T])(nil)
// var _ encDriver = (*jsonEncDriver[T])(nil)

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *jsonEncDriver[T]) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*JsonHandle)
	d.e = shared
	if shared.bytes {
		fp = jsonFpEncBytes
	} else {
		fp = jsonFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *jsonEncDriver[T]) writeBytesAsis(b []byte) { e.w.writeb(b) }

// func (e *jsonEncDriver[T]) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *jsonEncDriver[T]) writerEnd() { e.w.end() }

func (e *jsonEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *jsonEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *jsonDecDriver[T]) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*JsonHandle)
	d.d = shared
	if shared.bytes {
		fp = jsonFpDecBytes
	} else {
		fp = jsonFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *jsonDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *jsonDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *jsonDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

// ---- (custom stanza)

var errJsonNoBd = errors.New("descBd unsupported in json")

func (d *jsonDecDriver[T]) descBd() (s string) {
	halt.onerror(errJsonNoBd)
	return
}

func (d *jsonEncDriver[T]) init2(enc encoderI) {
	d.enc = enc
	// d.e.js = true
}

func (d *jsonDecDriver[T]) init2(dec decoderI) {
	d.dec = dec
	// var x []byte
	// d.buf = &x
	// d.buf = new([]byte)
	d.buf = d.buf[:0]
	// d.d.js = true
	d.d.jsms = d.h.MapKeyAsString
}
