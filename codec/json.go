// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// This json support uses base64 encoding for bytes, because you cannot
// store and read any arbitrary string in json (only unicode).
//
// This library specifically supports UTF-8 for encoding and decoding only.
//
// Note:
//   - we cannot use strconv.Quote and strconv.Unquote because json quotes/unquotes differently.
//   - encode does not beautify. There is no whitespace when encoding.
//   - rpc calls which take single integer arguments or write single numeric arguments will not
//     work well, as it may not be possible to know when a number ends (unlike string, etc with terminator char).
//     Luckily, rpc support in this package mitigates that via rpcEncodeTerminator interface.

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

//--------------------------------

var jsonLiterals = [...]byte{'t', 'r', 'u', 'e', 'f', 'a', 'l', 's', 'e', 'n', 'u', 'l', 'l'}

type jsonEncDriver struct {
	w  encWriter
	h  *JsonHandle
	b  [64]byte // scratch
	bs []byte   // scratch
	noBuiltInTypes
}

func (e *jsonEncDriver) encodeNil() {
	e.w.writeb(jsonLiterals[9:]) // null
}

func (e *jsonEncDriver) encodeBool(b bool) {
	if b {
		e.w.writeb(jsonLiterals[:4]) // true
	} else {
		e.w.writeb(jsonLiterals[4:9]) // false
	}
}

func (e *jsonEncDriver) encodeFloat32(f float32) {
	e.w.writeb(strconv.AppendFloat(e.b[:0], float64(f), 'E', -1, 32))
}

func (e *jsonEncDriver) encodeFloat64(f float64) {
	// e.w.writestr(strconv.FormatFloat(f, 'E', -1, 64))
	e.w.writeb(strconv.AppendFloat(e.b[:0], f, 'E', -1, 64))
}

func (e *jsonEncDriver) encodeInt(v int64) {
	e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
}

func (e *jsonEncDriver) encodeUint(v uint64) {
	e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
}

func (e *jsonEncDriver) encodeExt(rv reflect.Value, xtag uint64, ext Ext, en *Encoder) {
	if v := ext.ConvertExt(rv); v == nil {
		e.encodeNil()
	} else {
		en.encode(v)
	}
}

func (e *jsonEncDriver) encodeRawExt(re *RawExt, en *Encoder) {
	// only encodes re.Value (never re.Data)
	if re.Value == nil {
		e.encodeNil()
	} else {
		en.encode(re.Value)
	}
}

func (e *jsonEncDriver) encodeArrayStart(length int) {
	e.w.writen1('[')
}

func (e *jsonEncDriver) encodeArrayEntrySeparator() {
	e.w.writen1(',')
}

func (e *jsonEncDriver) encodeArrayEnd() {
	e.w.writen1(']')
}

func (e *jsonEncDriver) encodeMapStart(length int) {
	e.w.writen1('{')
}

func (e *jsonEncDriver) encodeMapEntrySeparator() {
	e.w.writen1(',')
}

func (e *jsonEncDriver) encodeMapKVSeparator() {
	e.w.writen1(':')
}

func (e *jsonEncDriver) encodeMapEnd() {
	e.w.writen1('}')
}

func (e *jsonEncDriver) encodeString(c charEncoding, v string) {
	// e.w.writestr(strconv.Quote(v))
	e.quoteStr(v)
}

func (e *jsonEncDriver) encodeSymbol(v string) {
	e.encodeString(c_UTF8, v)
}

func (e *jsonEncDriver) encodeStringBytes(c charEncoding, v []byte) {
	if c == c_RAW {
		slen := base64.StdEncoding.EncodedLen(len(v))
		if e.bs == nil {
			e.bs = e.b[:]
		}
		if cap(e.bs) >= slen {
			e.bs = e.bs[:slen]
		} else {
			e.bs = make([]byte, slen)
		}
		base64.StdEncoding.Encode(e.bs, v)
		e.w.writen1('"')
		e.w.writeb(e.bs)
		e.w.writen1('"')
	} else {
		e.encodeString(c, string(v))
	}
}

func (e *jsonEncDriver) quoteStr(s string) {
	// adapted from std pkg encoding/json
	const hex = "0123456789abcdef"
	w := e.w
	w.writen1('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				w.writestr(s[start:i])
			}
			switch b {
			case '\\', '"':
				w.writen2('\\', b)
			case '\n':
				w.writen2('\\', 'n')
			case '\r':
				w.writen2('\\', 'r')
			case '\b':
				w.writen2('\\', 'b')
			case '\f':
				w.writen2('\\', 'f')
			case '\t':
				w.writen2('\\', 't')
			default:
				w.writestr(`\u00`)
				w.writen1(hex[b>>4])
				w.writen1(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w.writestr(s[start:i])
			}
			w.writestr(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
		// Both technically valid JSON, but bomb on JSONP, so fix here.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w.writestr(s[start:i])
			}
			w.writestr(`\u202`)
			w.writen1(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w.writestr(s[start:])
	}
	w.writen1('"')
}

//--------------------------------

type jsonDecDriver struct {
	h  *JsonHandle
	r  decReader
	ct valueType // container type. one of unset, array or map.
	b  [64]byte  // scratch
	b2 [8]byte   // scratch
	noBuiltInTypes
}

// This will skip whitespace characters and return the next byte to read.
// The next byte determines what the value will be one of.
func (d *jsonDecDriver) skipWhitespace(unread bool) (b byte) {
	for b = d.r.readn1(); b == ' ' || b == '\t' || b == '\r' || b == '\n'; b = d.r.readn1() {
	}
	if unread {
		d.r.unreadn1()
	}
	return b
}

func (d *jsonDecDriver) initReadNext() {
	d.skipWhitespace(true)
	d.ct = valueTypeUnset
}

func (d *jsonDecDriver) checkBreak() bool {
	b := d.skipWhitespace(true)
	return b == '}' || b == ']'
}

func (d *jsonDecDriver) readStr(s []byte) {
	bs := d.b[:len(s)]
	d.r.readb(bs)
	if !bytes.Equal(bs, s) {
		decErr("json: expecting null: got %s", bs)
	}
}

func (d *jsonDecDriver) tryDecodeAsNil() bool {
	b := d.skipWhitespace(true)
	if b == 'n' {
		d.readStr(jsonLiterals[9:]) // null
		d.ct = valueTypeNil
		return true
	}
	return false
}

func (d *jsonDecDriver) decodeBool() bool {
	b := d.skipWhitespace(false)
	if b == 'f' {
		d.readStr(jsonLiterals[5:9]) // alse
		return false
	}
	if b == 't' {
		d.readStr(jsonLiterals[1:4]) // rue
		return true
	}
	decErr("json: decode bool: got first char %c", b)
	panic("unreachable")
}

func (d *jsonDecDriver) readMapStart() int {
	d.expectChar('{')
	d.ct = valueTypeMap
	return -1
}

func (d *jsonDecDriver) readArrayStart() int {
	d.expectChar('[')
	d.ct = valueTypeArray
	return -1
}
func (d *jsonDecDriver) readMapEnd() {
	d.expectChar('}')
}
func (d *jsonDecDriver) readArrayEnd() {
	d.expectChar(']')
}
func (d *jsonDecDriver) readArrayEntrySeparator() {
	d.expectChar(',')
}
func (d *jsonDecDriver) readMapEntrySeparator() {
	d.expectChar(',')
}
func (d *jsonDecDriver) readMapKVSeparator() {
	d.expectChar(':')
}
func (d *jsonDecDriver) expectChar(c uint8) {
	b := d.skipWhitespace(false)
	if b != c {
		decErr("json: expect char %c but got char %c", c, b)
	}
}

func (d *jsonDecDriver) isContainerType(vt valueType) bool {
	// check container type by checking the first char
	if d.ct == valueTypeUnset {
		b := d.skipWhitespace(true)
		switch b {
		case '{':
			d.ct = valueTypeMap
		case '[':
			d.ct = valueTypeArray
		case 'n':
			d.ct = valueTypeNil
		case '"':
			d.ct = valueTypeString
		}
	}
	switch vt {
	case valueTypeNil, valueTypeBytes, valueTypeString, valueTypeArray, valueTypeMap:
		return d.ct == vt
	}
	decErr("isContainerType: unsupported parameter: %v", vt)
	panic("unreachable")
}

func (d *jsonDecDriver) decNum() (i int64, f float64, isFloat bool) {
	// If it is has a . or an e|E, decode as a float; else decode as an int.
	b := d.skipWhitespace(false)
	if !(b == '+' || b == '-' || b == '.' || (b >= '0' && b <= '9')) {
		decErr("json: decNum: got first char %c", b)
	}
	var eof bool
	bs := d.b[:0]
	for {
		if b == '.' || b == 'e' || b == 'E' {
			isFloat = true
		} else if b == '+' || b == '-' || (b >= '0' && b <= '9') {
		} else {
			d.r.unreadn1()
			break
		}
		bs = append(bs, b)
		b, eof = d.readn1eof()
		if eof {
			break
		}
	}
	var err error
	if isFloat {
		f, err = strconv.ParseFloat(string(bs), 64)
	} else {
		i, err = strconv.ParseInt(string(bs), 10, 64)
	}
	if err != nil {
		decErr("decNum: %v", err)
	}
	return
}

func (d *jsonDecDriver) decodeInt(bitsize uint8) (i int64) {
	i, xf, xisFloat := d.decNum()
	if xisFloat {
		i = int64(xf)
	}
	checkOverflow(0, i, bitsize)
	return
}

func (d *jsonDecDriver) decodeUint(bitsize uint8) (ui uint64) {
	xi, xf, xisFloat := d.decNum()
	if (xisFloat && xf < 0) || (!xisFloat && xi < 0) {
		decErr("received negative number decoding number into uint: %v, %v", xi, xf)
	}
	if xisFloat {
		ui = uint64(xf)
	} else {
		ui = uint64(xi)
	}
	checkOverflow(ui, 0, bitsize)
	return
}

func (d *jsonDecDriver) decodeFloat(chkOverflow32 bool) (f float64) {
	xi, f, xisFloat := d.decNum()
	if !xisFloat {
		f = float64(xi)
	}
	checkOverflowFloat32(f, chkOverflow32)
	return
}

func (d *jsonDecDriver) decodeBytes(bs []byte) (bsOut []byte, changed bool) {
	s := d.decStringAsBytes(nil)
	slen := base64.StdEncoding.DecodedLen(len(s))
	if len(bs) < slen {
		changed = true
		bsOut = make([]byte, slen)
		bs = bsOut
	} else if len(bs) > slen {
		changed = true
		bsOut = bs[:slen]
		bs = bsOut
	}
	base64.StdEncoding.Decode(bs, s)
	return
}

func (d *jsonDecDriver) decodeExt(rv reflect.Value, xtag uint64, ext Ext, de *Decoder) (realxtag uint64) {
	if ext == nil {
		re := rv.Interface().(*RawExt)
		re.Tag = xtag
		de.decode(&re.Value)
	} else {
		var v interface{}
		de.decode(&v)
		ext.UpdateExt(rv, v)
	}
	return
}

func (d *jsonDecDriver) decodeString() (s string) {
	return string(d.decStringAsBytes(nil))
}

func (d *jsonDecDriver) decStringAsBytes(v []byte) []byte {
	d.expectChar('"')
	if v == nil {
		v = d.b[:0]
	}
	for {
		c := d.r.readn1()
		if c == '"' {
			break
		} else if c == '\\' {
			c = d.r.readn1()
			switch c {
			case '"', '\\', '/', '\'':
				v = append(v, c)
			case 'b':
				v = append(v, '\b')
			case 'f':
				v = append(v, '\f')
			case 'n':
				v = append(v, '\n')
			case 'r':
				v = append(v, '\r')
			case 't':
				v = append(v, '\t')
			case 'u':
				rr := d.jsonU4(false)
				fmt.Printf("$$$$$$$$$: is surrogate: %v\n", utf16.IsSurrogate(rr))
				if utf16.IsSurrogate(rr) {
					rr = utf16.DecodeRune(rr, d.jsonU4(true))
				}
				w2 := utf8.EncodeRune(d.b2[:], rr)
				v = append(v, d.b2[:w2]...)
			default:
				decErr("json: unsupported escaped value: %c", c)
			}
		} else {
			v = append(v, c)
		}
	}
	return v
}

func (d *jsonDecDriver) jsonU4(checkSlashU bool) rune {
	if checkSlashU && !(d.r.readn1() == '\\' && d.r.readn1() == 'u') {
		decErr(`json: unquoteStr: invalid unicode sequence. Expecting \u`)
	}
	d.r.readb(d.b2[:4])
	ui4, err4 := strconv.ParseUint(string(d.b2[:4]), 16, 64)
	if err4 != nil {
		decErr("json: unquoteStr: %v", err4)
	}
	return rune(ui4)
}

func (d *jsonDecDriver) decodeNaked(de *Decoder) (v interface{}, vt valueType, decodeFurther bool) {
	n := d.skipWhitespace(true)
	switch n {
	case 'n':
		d.readStr(jsonLiterals[9:]) // null
		vt = valueTypeNil
	case 'f':
		d.readStr(jsonLiterals[4:9]) // false
		vt = valueTypeBool
		v = false
	case 't':
		d.readStr(jsonLiterals[:4]) // true
		vt = valueTypeBool
		v = true
	case '{':
		vt = valueTypeMap
		decodeFurther = true
	case '[':
		vt = valueTypeArray
		decodeFurther = true
	case '"':
		vt = valueTypeString
		v = d.decodeString()
	default:
		xi, xf, xisFloat := d.decNum()
		if xisFloat {
			vt = valueTypeFloat
			v = xf
		} else if xi < 0 || d.h.SignedInteger {
			vt = valueTypeInt
			v = xi
		} else {
			vt = valueTypeUint
			v = uint64(xi)
		}
	}
	return
}

func (d *jsonDecDriver) readn1eof() (v uint8, eof bool) {
	defer func() {
		if x := recover(); x != nil {
			if x != io.EOF {
				panic(x)
			}
			eof = true
		}
	}()
	v = d.r.readn1()
	return
}

//----------------------

// JsonHandle is a handle for JSON encoding format.
//
// Json is comprehensively supported:
//    - decodes numbers into interface{} as int, uint or float64
//    - encodes and decodes []byte using base64 Std Encoding
//    - UTF-8 support for encoding and decoding
//
// It has better performance than the json library in the standard library,
// by leveraging the performance improvements of the codec library and minimizing allocations.
//
// In addition, it doesn't read more bytes than necessary during a decode, which allows
// the use of this to read multiple values from a stream containing json and non-json content.
type JsonHandle struct {
	BasicHandle
}

func (h *JsonHandle) newEncDriver(w encWriter) encDriver {
	return &jsonEncDriver{w: w, h: h}
}

func (h *JsonHandle) newDecDriver(r decReader) decDriver {
	return &jsonDecDriver{r: r, h: h}
}

func (h *JsonHandle) isBinaryEncoding() bool {
	return false
}

var jsonEncodeTerminate = []byte{' '}

func (h *JsonHandle) rpcEncodeTerminate() []byte {
	return jsonEncodeTerminate
}

var _ decDriver = (*jsonDecDriver)(nil)
var _ encDriver = (*jsonEncDriver)(nil)
