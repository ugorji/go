// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
	"time"
)

//BincHandle is a Handle for the Binc Schema-Free Encoding Format
//defined at http://www.ugorji.net/project/binc .
//
//BincHandle currently supports all Binc features with the following EXCEPTIONS:
//  - only integers up to 64 bits of precision (degree of precision <= 3) are supported.
//    big integers are unsupported.
//  - Only IEEE 754 binary32 and binary64 floats are supported (ie Go float32 and float64 types).
//    extended precision and decimal IEEE 754 floats are unsupported.
//  - Only UTF-8 strings supported. 
//    Unicode_Other Binc types (UTF16, UTF32) are currently unsupported.
//Note that these EXCEPTIONS are temporary and full support is possible and may happen soon.
type BincHandle struct {
	encdecHandle
	DecodeOptions
}

// vd as low 4 bits (there are 16 slots)
const (
	bincVdSpecial byte = iota
	bincVdUint
	bincVdInt
	bincVdFloat
	
	bincVdString
	bincVdByteArray
	bincVdArray
	bincVdMap
	
	bincVdTimestamp
	bincVdSmallInt
	bincVdUnicodeOther

	// 4 open slots left ...

	bincVdCustomExt = 0x0f
)

const (
	bincSpNil byte = iota
	bincSpFalse
	bincSpTrue
	bincSpNan
	bincSpPosInf
	bincSpNegInf
	bincSpZero
    bincSpNegOne
)

const (
	bincFlBin16 byte = iota
	bincFlBin32
	_ // bincFlBin32e
	bincFlBin64
	// others not currently supported
)

type bincEncoder struct { 
	w encWriter
}

type bincDecoder struct {
	r decReader
	bdRead bool
	bd byte 
	vd byte
	vs byte 
}

func (_ *BincHandle) newEncoder(w encWriter) encoder {
	return &bincEncoder{w: w}
}

func (_ *BincHandle) newDecoder(r decReader) decoder {
	return &bincDecoder{r: r}
}

func (_ *BincHandle) writeExt() bool {
	return true
}

func (e *bincEncoder) encodeBuiltinType(rt reflect.Type, rv reflect.Value) bool {
	switch rt {
	case timeTyp:
		bs := encodeTime(rv.Interface().(time.Time))
		e.w.writen1(bincVdTimestamp << 4 | uint8(len(bs)))
		e.w.writeb(bs)
		return true
	}
	return false
}

func (e *bincEncoder) encodeNil() { 
	e.w.writen1(bincVdSpecial << 4 | bincSpNil)
}

func (e *bincEncoder) encodeBool(b bool) { 
	if b {
		e.w.writen1(bincVdSpecial << 4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial << 4 | bincSpFalse)
	}		
}

func (e *bincEncoder) encodeFloat32(f float32) { 
	e.w.writen1(bincVdFloat << 4 | bincFlBin32)
	e.w.writeUint32(math.Float32bits(f))
}

func (e *bincEncoder) encodeFloat64(f float64) { 
	e.w.writen1(bincVdFloat << 4 | bincFlBin64)
	e.w.writeUint64(math.Float64bits(f))
}

func (e *bincEncoder) encodeInt(v int64) { 
	switch {
	case v == 0:
		e.w.writen1(bincVdSpecial << 4 | bincSpZero)
	case v == -1:
		e.w.writen1(bincVdSpecial << 4 | bincSpNegOne)
	case v >= 1 && v <= 16:
		e.w.writen1(bincVdSmallInt << 4 | byte(v-1))
	case v >= math.MinInt8 && v <= math.MaxInt8:
		e.w.writen2(bincVdInt << 4, byte(v))
	case v >= math.MinInt16 && v <= math.MaxInt16:
		e.w.writen1(bincVdInt << 4 | 0x01)
		e.w.writeUint16(uint16(v))
	case v >= math.MinInt32 && v <= math.MaxInt32:
		e.w.writen1(bincVdInt << 4 | 0x02)
		e.w.writeUint32(uint32(v))
	default:
		e.w.writen1(bincVdInt << 4 | 0x03)
		e.w.writeUint64(uint64(v))
	}
}

func (e *bincEncoder) encodeUint(v uint64) { 
	e.encNumber(bincVdUint << 4, v)
}

func (e *bincEncoder) encodeExtPreamble(xtag byte, length int)  {
	e.encLen(bincVdCustomExt << 4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncoder) encodeArrayPreamble(length int) {
	e.encLen(bincVdArray << 4, uint64(length))
}

func (e *bincEncoder) encodeMapPreamble(length int) { 
	e.encLen(bincVdMap << 4, uint64(length))
}

func (e *bincEncoder) encodeString(c charEncoding, v string) { 
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 {
		e.w.writestr(v)
	}	
}

func (e *bincEncoder) encodeStringBytes(c charEncoding, v []byte) { 
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 {
		e.w.writeb(v)
	}	
}

func (e *bincEncoder) encBytesLen(c charEncoding, length uint64) {
	//TODO: support bincUnicodeOther (for now, just use string or bytearray)
	if c == c_RAW {
		e.encLen(bincVdByteArray << 4, length)
	} else {
		e.encLen(bincVdString << 4, length)
	}
}

func (e *bincEncoder) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l + 4))
	} else {
		e.encNumber(bd, l)
	}
}
	
func (e *bincEncoder) encNumber(bd byte, v uint64) {
	switch {
	case v <= math.MaxUint8:
		e.w.writen2(bd, byte(v))
	case v <= math.MaxUint16:
		e.w.writen1(bd | 0x01)
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		e.w.writen1(bd | 0x02)
		e.w.writeUint32(uint32(v))
	default:
		e.w.writen1(bd | 0x03)
		e.w.writeUint64(uint64(v))
	}
}


//------------------------------------


func (d *bincDecoder) initReadNext() {
	if d.bdRead {
		return
	}
	d.bd = d.r.readUint8()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
}

func (d *bincDecoder) currentIsNil() bool {
	if d.bd == bincVdSpecial << 4 | bincSpNil {
		d.bdRead = false
		return true
	} 
	return false
}

func (d *bincDecoder) decodeBuiltinType(rt reflect.Type, rv reflect.Value) bool {
	switch rt {
	case timeTyp:
		if d.vd != bincVdTimestamp {
			decErr("Invalid d.vd. Expecting 0x%x. Received: 0x%x", bincVdTimestamp, d.vd)
		}
		tt, err := decodeTime(d.r.readn(int(d.vs)))
		if err != nil {
			panic(err)
		}
		rv.Set(reflect.ValueOf(tt))
		d.bdRead = false
		return true
	}
	return false
}

func (d *bincDecoder) decFloat() (f float64) {
	switch d.vs {
	case bincFlBin32:
		f = float64(math.Float32frombits(d.r.readUint32()))
	case bincFlBin64:
		f = math.Float64frombits(d.r.readUint64())
	default:
		decErr("only float32 and float64 are supported")
	}
	return
}

func (d *bincDecoder) decFloatI() (f interface{}) {
	switch d.vs {
	case bincFlBin32:
		f = math.Float32frombits(d.r.readUint32())
	case bincFlBin64:
		f = math.Float64frombits(d.r.readUint64())
	default:
		decErr("only float32 and float64 are supported")
	}
	return
}

func (d *bincDecoder) decInt() (v int64) {
	switch d.vs {
	case 0:
		v = int64(int8(d.r.readUint8()))
	case 1:
		v = int64(int16(d.r.readUint16()))
	case 2:
		v = int64(int32(d.r.readUint32()))
	case 3:
		v = int64(d.r.readUint64())
	default:
		decErr("integers with greater than 64 bits of precision not supported")
	}
	return
}

func (d *bincDecoder) decIntI() (v interface{}) {
	switch d.vs {
	case 0:
		v = int8(d.r.readUint8())
	case 1:
		v = int16(d.r.readUint16())
	case 2:
		v = int32(d.r.readUint32())
	case 3:
		v = int64(d.r.readUint64())
	default:
		decErr("integers with greater than 64 bits of precision not supported")
	}
	return
}

func (d *bincDecoder) decUint() (v uint64) {
	switch d.vs {
	case 0:
		v = uint64(d.r.readUint8())
	case 1:
		v = uint64(d.r.readUint16())
	case 2:
		v = uint64(d.r.readUint32())
	case 3:
		v = uint64(d.r.readUint64())
	default:
		decErr("integers with greater than 64 bits of precision not supported")
	}
	return
}

func (d *bincDecoder) decUintI() (v interface{}) {
	switch d.vs {
	case 0:
		v = d.r.readUint8()
	case 1:
		v = d.r.readUint16()
	case 2:
		v = d.r.readUint32()
	case 3:
		v = d.r.readUint64()
	default:
		decErr("integers with greater than 64 bits of precision not supported")
	}
	return
}

func (d *bincDecoder) decIntAny() (i int64) {
	switch d.vd {
	case bincVdInt: 
		i = d.decInt()
	case bincVdSmallInt:
		i = int64(d.vs) + 1
	case bincVdSpecial:
		switch d.vs {
		case bincSpZero:
			//i = 0
		case bincSpNegOne:
			i = -1
		default:
			decErr("numeric decode fails for special value: d.vs: 0x%x", d.vs)
		}
	default:
		decErr("number can only be decoded from uint or int values. d.bd: 0x%x, d.vd: 0x%x", d.bd, d.vd)
	}
	return
}
		
func (d *bincDecoder) decodeInt(bitsize uint8) (i int64) {
	switch d.vd {
	case bincVdUint: 
		i = int64(d.decUint())
	default:
		i = d.decIntAny()
	}
	// check overflow (logic adapted from std pkg reflect/value.go OverflowUint()
	if bitsize > 0 {
		if trunc := (i << (64 - bitsize)) >> (64 - bitsize); i != trunc {
			decErr("Overflow int value: %v", i)
		}
	}
	d.bdRead = false
	return
}

func (d *bincDecoder) decodeUint(bitsize uint8) (ui uint64)  {
	switch d.vd {
	case bincVdUint: 
		ui = d.decUint()
	default:
		if i := d.decIntAny(); i >= 0 {
			ui = uint64(i)
		} else {
			decErr("Assigning negative signed value: %v, to unsigned type", i)
		}
	}	
	// check overflow (logic adapted from std pkg reflect/value.go OverflowUint()
	if bitsize > 0 {
		if trunc := (ui << (64 - bitsize)) >> (64 - bitsize); ui != trunc {
			decErr("Overflow uint value: %v", ui) 
		}
	}
	d.bdRead = false
	return
}

func (d *bincDecoder) decodeFloat(chkOverflow32 bool) (f float64) {
	if d.vd == bincVdSpecial {
		switch d.vs {
		case bincSpNan:
			return math.NaN()
		case bincSpPosInf:
			return math.Inf(1)
		case bincSpNegInf:
			return math.Inf(-1)
		}
	}
	switch d.vd {
	case bincVdFloat: 
		f = d.decFloat()
	case bincVdUint: 
		f = float64(d.decUint())
	default:
		f = float64(d.decIntAny())
	}
	
	// check overflow (logic adapted from std pkg reflect/value.go OverflowFloat()
	if chkOverflow32 {
		f2 := f
		if f2 < 0 {
			f2 = -f
		}
		if math.MaxFloat32 < f2 && f2 <= math.MaxFloat64 {
			decErr("Overflow float32 value: %v", f2)
		}
	}
	d.bdRead = false
	return
}


// bool can be decoded from bool only (single byte).
func (d *bincDecoder) decodeBool() (b bool)  {
	switch d.bd {
	case (bincVdSpecial | bincSpFalse):
		// b = false
	case (bincVdSpecial | bincSpTrue):
		b = true
	default:
		decErr("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *bincDecoder) readMapLen() (length int) {
	if d.vd != bincVdMap {
		decErr("Invalid d.vd for map. Expecting 0x%x. Got: 0x%x", bincVdMap, d.vd)
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecoder) readArrayLen() (length int) {
	if d.vd != bincVdArray {
		decErr("Invalid d.vd for array. Expecting 0x%x. Got: 0x%x", bincVdArray, d.vd)
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecoder) decLen() int {
	if d.vs <= 3 {
		return int(d.decUint())
	} 
	return int(d.vs - 4)
}

func (d *bincDecoder) decBytesLen() int {
	//decode string from either bytearray or string
	if d.vd != bincVdString && d.vd != bincVdByteArray {
		decErr("Invalid d.vd for string. Expecting string: 0x%x or bytearray: 0x%x. Got: 0x%x", 
			bincVdString, bincVdByteArray, d.vd)
	}
	return d.decLen()
}

func (d *bincDecoder) decodeString() (s string)  {
	if length := d.decBytesLen(); length > 0 {
		s = string(d.r.readn(length))
	}
	d.bdRead = false
	return
}

func (d *bincDecoder) decodeStringBytes(bs []byte) (bsOut []byte, changed bool) {
	clen := d.decBytesLen()
	if clen > 0 {
		// if no contents in stream, don't update the passed byteslice	
		if len(bs) != clen {
			if len(bs) > clen {
				bs = bs[:clen]
			} else {
				bs = make([]byte, clen)
			}
			bsOut = bs
			changed = true
		}
		d.r.readb(bs)
	}
	d.bdRead = false
	return
}

func (d *bincDecoder) decodeExt(tag byte) (xbs []byte) {
	switch d.vd {
	case bincVdCustomExt:
		l := d.decLen()
		if xtag := d.r.readUint8(); xtag != tag {
			decErr("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		xbs = d.r.readn(l)
	case bincVdByteArray:
		xbs, _ = d.decodeStringBytes(nil)
	default:
		decErr("Invalid d.vd for extensions (Expecting extensions or byte array). Got: 0x%x", d.vd)
	}
	d.bdRead = false
	return	
}

func (d *bincDecoder) decodeNaked(h decodeHandleI) (rv reflect.Value, ctx decodeNakedContext) {
	d.initReadNext()
	var v interface{}

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil:
			ctx = dncNil
			d.bdRead = false
		case bincSpFalse:
			v = false
		case bincSpTrue:
			v = true
		case bincSpNan:
			v = math.NaN()
		case bincSpPosInf:
			v = math.Inf(1)
		case bincSpNegInf:
			v = math.Inf(-1)
		case bincSpZero:
			v = int8(0)
		case bincSpNegOne:
			v = int8(-1)
		default:
			decErr("Unrecognized special value 0x%x", d.vs)
		}
	case bincVdUint:
		v = d.decUintI()
	case bincVdInt:
		v = d.decIntI()
	case bincVdSmallInt:
		v = int8(d.vs) + 1
	case bincVdFloat:
		v = d.decFloatI()
	case bincVdString:
		v = d.decodeString()
	case bincVdByteArray:
		v, _ = d.decodeStringBytes(nil)
	case bincVdTimestamp:
		tt, err := decodeTime(d.r.readn(int(d.vs)))
		if err != nil {
			panic(err)
		}
		v = tt
	case bincVdCustomExt:
		//ctx = dncExt
		l := d.decLen()
		xtag := d.r.readUint8()
		opts := h.(*BincHandle)
		rt, bfn := opts.getDecodeExtForTag(xtag)
		if rt == nil {
			decErr("Unable to find type mapped to extension tag: %v", xtag)
		}
		if rt.Kind() == reflect.Ptr {
			rv = reflect.New(rt.Elem())
		} else {
			rv = reflect.New(rt).Elem()
		}
		if fnerr := bfn(rv, d.r.readn(l)); fnerr != nil {
			panic(fnerr)
		} 
	case bincVdArray:
		ctx = dncContainer
		opts := h.(*BincHandle)
		if opts.SliceType == nil {
			rv = reflect.New(intfSliceTyp).Elem()
		} else {
			rv = reflect.New(opts.SliceType).Elem()
		}
	case bincVdMap:
		ctx = dncContainer
		opts := h.(*BincHandle)
		if opts.MapType == nil {
			rv = reflect.MakeMap(mapIntfIntfTyp)
		} else {
			rv = reflect.MakeMap(opts.MapType)
		}
	default:
		decErr("Unrecognized d.vd: 0x%x", d.vd)
	}

	if ctx == dncHandled {
		d.bdRead = false
		if v != nil {
			rv = reflect.ValueOf(v)
		}
	}
	return
}


