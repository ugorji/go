// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
	"math"
	"reflect"
	"time"
	"unicode/utf8"
)

// major
const (
	cborMajorUint byte = iota
	cborMajorNegInt
	cborMajorBytes
	cborMajorString
	cborMajorArray
	cborMajorMap
	cborMajorTag
	cborMajorSimpleOrFloat
)

// simple
const (
	cborBdFalse byte = 0xf4 + iota
	cborBdTrue
	cborBdNil
	cborBdUndefined
	cborBdExt
	cborBdFloat16
	cborBdFloat32
	cborBdFloat64
)

// indefinite
const (
	cborBdIndefiniteBytes  byte = 0x5f
	cborBdIndefiniteString byte = 0x7f
	cborBdIndefiniteArray  byte = 0x9f
	cborBdIndefiniteMap    byte = 0xbf
	cborBdBreak            byte = 0xff
)

// These define some in-stream descriptors for
// manual encoding e.g. when doing explicit indefinite-length
const (
	CborStreamBytes  byte = 0x5f
	CborStreamString byte = 0x7f
	CborStreamArray  byte = 0x9f
	CborStreamMap    byte = 0xbf
	CborStreamBreak  byte = 0xff
)

// base values
const (
	cborBaseUint   byte = 0x00
	cborBaseNegInt byte = 0x20
	cborBaseBytes  byte = 0x40
	cborBaseString byte = 0x60
	cborBaseArray  byte = 0x80
	cborBaseMap    byte = 0xa0
	cborBaseTag    byte = 0xc0
	cborBaseSimple byte = 0xe0
)

// const (
// 	cborSelfDesrTag  byte = 0xd9
// 	cborSelfDesrTag2 byte = 0xd9
// 	cborSelfDesrTag3 byte = 0xf7
// )

var (
	cbordescSimpleNames = map[byte]string{
		cborBdNil:     "nil",
		cborBdFalse:   "false",
		cborBdTrue:    "true",
		cborBdFloat16: "float",
		cborBdFloat32: "float",
		cborBdFloat64: "float",
		cborBdBreak:   "break",
	}
	cbordescIndefNames = map[byte]string{
		cborBdIndefiniteBytes:  "bytes*",
		cborBdIndefiniteString: "string*",
		cborBdIndefiniteArray:  "array*",
		cborBdIndefiniteMap:    "map*",
	}
	cbordescMajorNames = map[byte]string{
		cborMajorUint:          "(u)int",
		cborMajorNegInt:        "int",
		cborMajorBytes:         "bytes",
		cborMajorString:        "string",
		cborMajorArray:         "array",
		cborMajorMap:           "map",
		cborMajorTag:           "tag",
		cborMajorSimpleOrFloat: "simple",
	}
)

func cbordesc(bd byte) (s string) {
	bm := bd >> 5
	if bm == cborMajorSimpleOrFloat {
		s = cbordescSimpleNames[bd]
	} else {
		s = cbordescMajorNames[bm]
		if s == "" {
			s = cbordescIndefNames[bd]
		}
	}
	if s == "" {
		s = "unknown"
	}
	return
}

// -------------------

type cborEncDriver[T encWriter] struct {
	noBuiltInTypes
	encDriverNoState
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT

	h   *CborHandle
	e   *encoderBase
	w   T
	enc encoderI

	// scratch buffer for: encode time, numbers, etc
	//
	// RFC3339Nano uses 35 chars: 2006-01-02T15:04:05.999999999Z07:00
	b [40]byte
}

func (e *cborEncDriver[T]) EncodeNil() {
	e.w.writen1(cborBdNil)
}

func (e *cborEncDriver[T]) EncodeBool(b bool) {
	if b {
		e.w.writen1(cborBdTrue)
	} else {
		e.w.writen1(cborBdFalse)
	}
}

func (e *cborEncDriver[T]) EncodeFloat32(f float32) {
	b := math.Float32bits(f)
	if e.h.OptimumSize {
		if h := floatToHalfFloatBits(b); halfFloatToFloatBits(h) == b {
			e.w.writen1(cborBdFloat16)
			e.w.writen2(bigen.PutUint16(h))
			return
		}
	}
	e.w.writen1(cborBdFloat32)
	e.w.writen4(bigen.PutUint32(b))
}

func (e *cborEncDriver[T]) EncodeFloat64(f float64) {
	if e.h.OptimumSize {
		if f32 := float32(f); float64(f32) == f {
			e.EncodeFloat32(f32)
			return
		}
	}
	e.w.writen1(cborBdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *cborEncDriver[T]) encUint(v uint64, bd byte) {
	if v <= 0x17 {
		e.w.writen1(byte(v) + bd)
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd+0x18, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 0x19)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 0x1a)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else { // if v <= math.MaxUint64 {
		e.w.writen1(bd + 0x1b)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *cborEncDriver[T]) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-1-v), cborBaseNegInt)
	} else {
		e.encUint(uint64(v), cborBaseUint)
	}
}

func (e *cborEncDriver[T]) EncodeUint(v uint64) {
	e.encUint(v, cborBaseUint)
}

func (e *cborEncDriver[T]) encLen(bd byte, length int) {
	e.encUint(uint64(length), bd)
}

func (e *cborEncDriver[T]) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else if e.h.TimeRFC3339 {
		e.encUint(0, cborBaseTag)
		e.encStringBytesS(cborBaseString, stringView(t.AppendFormat(e.b[:0], time.RFC3339Nano)))
	} else {
		e.encUint(1, cborBaseTag)
		t = t.UTC().Round(time.Microsecond)
		sec, nsec := t.Unix(), uint64(t.Nanosecond())
		if nsec == 0 {
			e.EncodeInt(sec)
		} else {
			e.EncodeFloat64(float64(sec) + float64(nsec)/1e9)
		}
	}
}

func (e *cborEncDriver[T]) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	e.encUint(uint64(xtag), cborBaseTag)
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.EncodeNil()
	} else {
		e.enc.encode(v)
	}
}

func (e *cborEncDriver[T]) EncodeRawExt(re *RawExt) {
	e.encUint(uint64(re.Tag), cborBaseTag)
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encode(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *cborEncDriver[T]) WriteArrayEmpty() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteArray)
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(cborBaseArray, 0)
	}
}

func (e *cborEncDriver[T]) WriteMapEmpty() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteMap)
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(cborBaseMap, 0)
	}
}

func (e *cborEncDriver[T]) WriteArrayStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteArray)
	} else {
		e.encLen(cborBaseArray, length)
	}
}

func (e *cborEncDriver[T]) WriteMapStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteMap)
	} else {
		e.encLen(cborBaseMap, length)
	}
}

func (e *cborEncDriver[T]) WriteMapEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriver[T]) WriteArrayEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriver[T]) EncodeString(v string) {
	bb := cborBaseString
	if e.h.StringToRaw {
		bb = cborBaseBytes
	}
	e.encStringBytesS(bb, v)
}

func (e *cborEncDriver[T]) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *cborEncDriver[T]) EncodeStringBytesRaw(v []byte) {
	if v == nil {
		e.EncodeNil()
	} else {
		e.encStringBytesS(cborBaseBytes, stringView(v))
	}
}

func (e *cborEncDriver[T]) encStringBytesS(bb byte, v string) {
	if e.h.IndefiniteLength {
		if bb == cborBaseBytes {
			e.w.writen1(cborBdIndefiniteBytes)
		} else {
			e.w.writen1(cborBdIndefiniteString)
		}
		var vlen uint = uint(len(v))
		blen := vlen / 4
		if blen == 0 {
			blen = 64
		} else if blen > 1024 {
			blen = 1024
		}
		for i := uint(0); i < vlen; {
			var v2 string
			i2 := i + blen
			if i2 >= i && i2 < vlen {
				v2 = v[i:i2]
			} else {
				v2 = v[i:]
			}
			e.encLen(bb, len(v2))
			e.w.writestr(v2)
			i = i2
		}
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(bb, len(v))
		e.w.writestr(v)
	}
}

// ----------------------

type cborDecDriver[T decReader] struct {
	decDriverNoopContainerReader
	// decDriverNoopNumberHelper
	noBuiltInTypes

	h   *CborHandle
	d   *decoderBase
	r   T
	dec decoderI
	bdAndBdread
	st bool // skip tags
	// bytes bool
}

func (d *cborDecDriver[T]) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *cborDecDriver[T]) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdNil || d.bd == cborBdUndefined {
		d.bdRead = false
		return true // null = true
	}
	return
}

func (d *cborDecDriver[T]) TryNil() bool {
	return d.advanceNil()
}

// skipTags is called to skip any tags in the stream.
//
// Since any value can be tagged, then we should call skipTags
// before any value is decoded.
//
// By definition, skipTags should not be called before
// checking for break, or nil or undefined.
func (d *cborDecDriver[T]) skipTags() {
	for d.bd>>5 == cborMajorTag {
		d.decUint()
		d.bd = d.r.readn1()
	}
}

func (d *cborDecDriver[T]) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdNil {
		d.bdRead = false // always consume nil after seeing it in container type
		return valueTypeNil
	}
	major := d.bd >> 5
	if major == cborMajorBytes {
		return valueTypeBytes
	} else if major == cborMajorString {
		return valueTypeString
	} else if major == cborMajorArray {
		return valueTypeArray
	} else if major == cborMajorMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *cborDecDriver[T]) CheckBreak() (v bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdBreak {
		d.bdRead = false
		v = true
	}
	return
}

func (d *cborDecDriver[T]) decUint() (ui uint64) {
	v := d.bd & 0x1f
	if v <= 0x17 {
		ui = uint64(v)
	} else if v == 0x18 {
		ui = uint64(d.r.readn1())
	} else if v == 0x19 {
		ui = uint64(bigen.Uint16(d.r.readn2()))
	} else if v == 0x1a {
		ui = uint64(bigen.Uint32(d.r.readn4()))
	} else if v == 0x1b {
		ui = uint64(bigen.Uint64(d.r.readn8()))
	} else {
		halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
	}
	return
}

func (d *cborDecDriver[T]) decLen() int {
	return int(d.decUint())
}

func (d *cborDecDriver[T]) decAppendIndefiniteBytes(bs []byte, major byte) []byte {
	d.bdRead = false
	for !d.CheckBreak() {
		chunkMajor := d.bd >> 5
		if chunkMajor != major {
			halt.errorf("malformed indefinite string/bytes %x (%s); contains chunk with major type %v, expected %v",
				d.bd, cbordesc(d.bd), chunkMajor, major)
		}
		n := uint(d.decLen())
		oldLen := uint(len(bs))
		newLen := oldLen + n
		if newLen > uint(cap(bs)) {
			bs2 := make([]byte, newLen, 2*uint(cap(bs))+n)
			copy(bs2, bs)
			bs = bs2
		} else {
			bs = bs[:newLen]
		}
		d.r.readb(bs[oldLen:newLen])
		if d.h.ValidateUnicode && major == cborMajorString && !utf8.Valid(bs[oldLen:newLen]) {
			halt.errorf("indefinite-length text string contains chunk that is not a valid utf-8 sequence: 0x%x", bs[oldLen:newLen])
		}
		d.bdRead = false
	}
	d.bdRead = false
	return bs
}

func (d *cborDecDriver[T]) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.r.readn2()))))
	case cborBdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case cborBdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *cborDecDriver[T]) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd >> 5 {
	case cborMajorUint:
		ui = d.decUint()
	case cborMajorNegInt:
		ui = d.decUint()
		neg = true
	default:
		ok = false
	}
	return
}

func (d *cborDecDriver[T]) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, true)
	d.bdRead = false
	return
}

func (d *cborDecDriver[T]) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *cborDecDriver[T]) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, true)
	d.bdRead = false
	return
}

// bool can be decoded from bool only (single byte).
func (d *cborDecDriver[T]) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdTrue {
		b = true
	} else if d.bd == cborBdFalse {
	} else {
		halt.errorf("not bool - %s %x/%s", msgBadDesc, d.bd, cbordesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *cborDecDriver[T]) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.st {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorMap {
		halt.errorf("error reading map; got major type: %x, expected %x/%s", d.bd>>5, cborMajorMap, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriver[T]) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.st {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorArray {
		halt.errorf("invalid array; got major type: %x, expect: %x/%s", d.bd>>5, cborMajorArray, cbordesc(d.bd))
	}
	return d.decLen()
}

// MARKER d.d.buf is ONLY used within DecodeBytes.
// Safe to use freely here only.

func (d *cborDecDriver[T]) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		d.bdRead = false
		bs = d.decAppendIndefiniteBytes(d.d.buf[:0], d.bd>>5)
		d.d.buf = bs
		state = dBytesAttachBuffer
		return
	}
	if d.bd == cborBdIndefiniteArray {
		d.bdRead = false
		bs = d.d.buf[:0]
		for !d.CheckBreak() {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		d.d.buf = bs
		state = dBytesAttachBuffer
		return
	}
	var cond bool
	if d.bd>>5 == cborMajorArray {
		d.bdRead = false
		slen := d.decLen()
		bs, cond = usableByteSlice(d.d.buf, slen)
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		d.d.buf = bs
		state = dBytesAttachBuffer
		return
	}
	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *cborDecDriver[T]) DecodeStringAsBytes() (out []byte, state dBytesAttachState) {
	out, state = d.DecodeBytes()
	if d.h.ValidateUnicode && !utf8.Valid(out) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", out)
	}
	return
}

func (d *cborDecDriver[T]) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
}

func (d *cborDecDriver[T]) decodeTime(xtag uint64) (t time.Time) {
	switch xtag {
	case 0:
		var err error
		t, err = time.Parse(time.RFC3339, stringView(bytesOk(d.DecodeStringAsBytes())))
		halt.onerror(err)
	case 1:
		f1, f2 := math.Modf(d.DecodeFloat64())
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		halt.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	}
	t = t.UTC().Round(time.Microsecond)
	return
}

func (d *cborDecDriver[T]) preDecodeExt(checkTag bool, xtag uint64) (realxtag uint64, ok bool) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	realxtag = d.decUint()
	d.bdRead = false
	if checkTag && xtag != realxtag {
		halt.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
	}
	ok = true
	return
}

func (d *cborDecDriver[T]) DecodeRawExt(re *RawExt) {
	if realxtag, ok := d.preDecodeExt(false, 0); ok {
		re.Tag = realxtag
		d.dec.decode(&re.Value)
		d.bdRead = false
	}
}

func (d *cborDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if _, ok := d.preDecodeExt(true, xtag); ok {
		if ext == SelfExt {
			d.dec.decodeAs(rv, basetype, false)
		} else {
			d.dec.interfaceExtConvertAndDecode(rv, ext)
		}
		d.bdRead = false
	}
}

func (d *cborDecDriver[T]) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd >> 5 {
	case cborMajorUint:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case cborMajorNegInt:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case cborMajorBytes:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString) //, d.h.ZeroCopy)
	case cborMajorString:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case cborMajorArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborMajorMap:
		n.v = valueTypeMap
		decodeFurther = true
	case cborMajorTag:
		n.v = valueTypeExt
		n.u = d.decUint()
		n.l = nil
		if n.u == 0 || n.u == 1 {
			d.bdRead = false
			n.v = valueTypeTime
			n.t = d.decodeTime(n.u)
		} else if d.st && d.h.getExtForTag(n.u) == nil {
			// d.skipTags() // no need to call this - tags already skipped
			d.bdRead = false
			d.DecodeNaked()
			return // return when done (as true recursive function)
		}
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined:
			n.v = valueTypeNil
		case cborBdFalse:
			n.v = valueTypeBool
			n.b = false
		case cborBdTrue:
			n.v = valueTypeBool
			n.b = true
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			n.v = valueTypeFloat
			n.f = d.DecodeFloat64()
		default:
			halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		}
	default: // should never happen
		halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	}
	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *cborDecDriver[T]) uintBytes() (v []byte, ui uint64) {
	// this is only used by nextValueBytes, so it's ok to
	// use readx and bigenstd here.
	switch vv := d.bd & 0x1f; vv {
	case 0x18:
		v = d.r.readx(1)
		ui = uint64(v[0])
	case 0x19:
		v = d.r.readx(2)
		ui = uint64(bigenstd.Uint16(v))
	case 0x1a:
		v = d.r.readx(4)
		ui = uint64(bigenstd.Uint32(v))
	case 0x1b:
		v = d.r.readx(8)
		ui = uint64(bigenstd.Uint64(v))
	default:
		if vv > 0x1b {
			halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
		}
		ui = uint64(vv)
	}
	return
}

func (d *cborDecDriver[T]) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

// func (d *cborDecDriver[T]) nextValueBytesR(v0 []byte) (v []byte) {
// 	d.readNextBd()
// 	v0 = append(v0, d.bd)
// 	d.r.startRecording(v0)
// 	d.nextValueBytesBdReadR()
// 	v = d.r.stopRecording()
// 	return
// }

func (d *cborDecDriver[T]) nextValueBytesBdReadR() {
	// var bs []byte
	var ui uint64

	switch d.bd >> 5 {
	case cborMajorUint, cborMajorNegInt:
		d.uintBytes()
	case cborMajorString, cborMajorBytes:
		if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				_, ui = d.uintBytes()
				d.r.skip(uint(ui))
			}
		} else {
			_, ui = d.uintBytes()
			d.r.skip(uint(ui))
		}
	case cborMajorArray:
		if d.bd == cborBdIndefiniteArray {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				d.nextValueBytesBdReadR()
			}
		} else {
			_, ui = d.uintBytes()
			for i := uint64(0); i < ui; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		}
	case cborMajorMap:
		if d.bd == cborBdIndefiniteMap {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			_, ui = d.uintBytes()
			for i := uint64(0); i < ui; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		}
	case cborMajorTag:
		d.uintBytes()
		d.readNextBd()
		d.nextValueBytesBdReadR()
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined, cborBdFalse, cborBdTrue: // pass
		case cborBdFloat16:
			d.r.skip(2)
		case cborBdFloat32:
			d.r.skip(4)
		case cborBdFloat64:
			d.r.skip(8)
		default:
			halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
		}
	default: // should never happen
		halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
	}
	return
}

// -------------------------

// CborHandle is a Handle for the CBOR encoding format,
// defined at http://tools.ietf.org/html/rfc7049 and documented further at http://cbor.io .
//
// CBOR is comprehensively supported, including support for:
//   - indefinite-length arrays/maps/bytes/strings
//   - (extension) tags in range 0..0xffff (0 .. 65535)
//   - half, single and double-precision floats
//   - all numbers (1, 2, 4 and 8-byte signed and unsigned integers)
//   - nil, true, false, ...
//   - arrays and maps, bytes and text strings
//
// None of the optional extensions (with tags) defined in the spec are supported out-of-the-box.
// Users can implement them as needed (using SetExt), including spec-documented ones:
//   - timestamp, BigNum, BigFloat, Decimals,
//   - Encoded Text (e.g. URL, regexp, base64, MIME Message), etc.
type CborHandle struct {
	binaryEncodingType
	// noElemSeparators
	BasicHandle

	// IndefiniteLength=true, means that we encode using indefinitelength
	IndefiniteLength bool

	// TimeRFC3339 says to encode time.Time using RFC3339 format.
	// If unset, we encode time.Time using seconds past epoch.
	TimeRFC3339 bool

	// SkipUnexpectedTags says to skip over any tags for which extensions are
	// not defined. This is in keeping with the cbor spec on "Optional Tagging of Items".
	//
	// Furthermore, this allows the skipping over of the Self Describing Tag 0xd9d9f7.
	SkipUnexpectedTags bool
}

// Name returns the name of the handle: cbor
func (h *CborHandle) Name() string { return "cbor" }

func (h *CborHandle) desc(bd byte) string { return cbordesc(bd) }

// SetInterfaceExt sets an extension
func (h *CborHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}

func (d *cborDecDriver[T]) reset() {
	d.bdAndBdread.reset()
	d.st = d.h.SkipUnexpectedTags
}

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *cborEncDriver[T]) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*CborHandle)
	d.e = shared
	if shared.bytes {
		fp = cborFpEncBytes
	} else {
		fp = cborFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *cborEncDriver[T]) writeBytesAsis(b []byte) { e.w.writeb(b) }

// func (e *cborEncDriver[T]) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }

func (e *cborEncDriver[T]) writerEnd() { e.w.end() }

func (e *cborEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *cborEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *cborDecDriver[T]) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*CborHandle)
	d.d = shared
	if shared.bytes {
		fp = cborFpDecBytes
	} else {
		fp = cborFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *cborDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *cborDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *cborDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

// ---- (custom stanza)

func (d *cborDecDriver[T]) descBd() string {
	return sprintf("%v (%s)", d.bd, cbordesc(d.bd))
}

func (d *cborDecDriver[T]) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

func (d *cborEncDriver[T]) init2(enc encoderI) {
	d.enc = enc
}

func (d *cborDecDriver[T]) init2(dec decoderI) {
	d.dec = dec
	// d.d.cbor = true
}
