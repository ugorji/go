// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
	"sync/atomic"
	"time"
	//"fmt"
)

//var _ = fmt.Printf

//BincHandle is a Handle for the Binc Schema-Free Encoding Format
//defined at https://github.com/ugorji/binc .
//
//BincHandle currently supports all Binc features with the following EXCEPTIONS:
//  - only integers up to 64 bits of precision are supported.
//    big integers are unsupported.
//  - Only IEEE 754 binary32 and binary64 floats are supported (ie Go float32 and float64 types).
//    extended precision and decimal IEEE 754 floats are unsupported.
//  - Only UTF-8 strings supported.
//    Unicode_Other Binc types (UTF16, UTF32) are currently unsupported.
//Note that these EXCEPTIONS are temporary and full support is possible and may happen soon.
type BincHandle struct {
	extHandle
	EncodeOptions
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
	bincVdSymbol

	bincVdDecimal
	_               // open slot
	_               // open slot
	bincVdCustomExt = 0x0f
)

const (
	bincSpNil byte = iota
	bincSpFalse
	bincSpTrue
	bincSpNan
	bincSpPosInf
	bincSpNegInf
	bincSpZeroFloat
	bincSpZero
	bincSpNegOne
)

const (
	bincFlBin16 byte = iota
	bincFlBin32
	_ // bincFlBin32e
	bincFlBin64
	_ // bincFlBin64e
	// others not currently supported
)

type bincEncDriver struct {
	w encWriter
	m map[string]uint16 // symbols
	s uint32            // symbols sequencer
	b [8]byte
}

type bincDecDriver struct {
	r      decReader
	h      *BincHandle
	bdRead bool
	bdType decodeEncodedType
	bd     byte
	vd     byte
	vs     byte
	b      [8]byte
	m      map[uint32]string // symbols (use uint32 as key, as map optimizes for it)
}

func (_ *BincHandle) newEncDriver(w encWriter) encDriver {
	return &bincEncDriver{w: w}
}

func (h *BincHandle) newDecDriver(r decReader) decDriver {
	return &bincDecDriver{r: r, h: h}
}

func (_ *BincHandle) writeExt() bool {
	return true
}

func (e *bincEncDriver) isBuiltinType(rt uintptr) bool {
	return rt == timeTypId
}

func (e *bincEncDriver) encodeBuiltinType(rt uintptr, rv reflect.Value) {
	switch rt {
	case timeTypId:
		bs := encodeTime(rv.Interface().(time.Time))
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	}
}

func (e *bincEncDriver) encodeNil() {
	e.w.writen1(bincVdSpecial<<4 | bincSpNil)
}

func (e *bincEncDriver) encodeBool(b bool) {
	if b {
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	}
}

func (e *bincEncDriver) encodeFloat32(f float32) {
	if f == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	}
	e.w.writen1(bincVdFloat<<4 | bincFlBin32)
	e.w.writeUint32(math.Float32bits(f))
}

func (e *bincEncDriver) encodeFloat64(f float64) {
	//if true { e.w.writen1(bincVdFloat << 4 | bincFlBin64); e.w.writeUint64(math.Float64bits(f)); return; }
	if f == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	}
	bigen.PutUint64(e.b[:], math.Float64bits(f))
	var i int = 7
	for ; i >= 0 && (e.b[i] == 0); i-- {
	}
	i++
	if i > 6 { // 7 or 8 ie < 2 trailing zeros
		e.w.writen1(bincVdFloat<<4 | bincFlBin64)
		e.w.writeb(e.b[:])
	} else {
		e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
		e.w.writen1(byte(i))
		e.w.writeb(e.b[:i])
	}
}

func (e *bincEncDriver) encInteger4(bd byte, v uint32) {
	const lim int = 4
	eb := e.b[:lim]
	bigen.PutUint32(eb, v)
	i := pruneSignExt(eb)
	e.w.writen1(bd | byte(lim-1-i))
	e.w.writeb(e.b[i:lim])
}

func (e *bincEncDriver) encInteger8(bd byte, v uint64) {
	const lim int = 8
	eb := e.b[:lim]
	bigen.PutUint64(eb, v)
	i := pruneSignExt(eb)
	e.w.writen1(bd | byte(lim-1-i))
	e.w.writeb(e.b[i:lim])
}

func (e *bincEncDriver) encodeInt(v int64) {
	const bd byte = bincVdInt << 4
	switch {
	case v == 0:
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	case v == -1:
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	case v >= 1 && v <= 16:
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	case v >= math.MinInt8 && v <= math.MaxInt8:
		e.w.writen2(bd|0x0, byte(v))
	case v >= math.MinInt16 && v <= math.MaxInt16:
		e.w.writen1(bd | 0x1)
		e.w.writeUint16(uint16(v))
	case v >= math.MinInt32 && v <= math.MaxInt32:
		e.encInteger4(bd, uint32(v))
	default:
		e.encInteger8(bd, uint64(v))
	}
}

func (e *bincEncDriver) encodeUint(v uint64) {
	const bd byte = bincVdUint << 4
	switch {
	case v <= math.MaxUint8:
		e.w.writen2(bd|0x0, byte(v))
	case v <= math.MaxUint16:
		e.w.writen1(bd | 0x01)
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		e.encInteger4(bd, uint32(v))
	default:
		e.encInteger8(bd, v)
	}
}

func (e *bincEncDriver) encodeExtPreamble(xtag byte, length int) {
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncDriver) encodeArrayPreamble(length int) {
	e.encLen(bincVdArray<<4, uint64(length))
}

func (e *bincEncDriver) encodeMapPreamble(length int) {
	e.encLen(bincVdMap<<4, uint64(length))
}

func (e *bincEncDriver) encodeString(c charEncoding, v string) {
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 {
		e.w.writestr(v)
	}
}

func (e *bincEncDriver) encodeSymbol(v string) {
	//if true { e.encodeString(c_UTF8, v); return; }

	//symbols only offer benefit when string length > 1.
	//This is because strings with length 1 take only 2 bytes to store
	//(bd with embedded length, and single byte for string val).

	l := len(v)
	switch l {
	case 0:
		e.encBytesLen(c_UTF8, 0)
		return
	case 1:
		e.encBytesLen(c_UTF8, 1)
		e.w.writen1(v[0])
		return
	}
	if e.m == nil {
		e.m = make(map[string]uint16, 16)
	}
	ui, ok := e.m[v]
	if ok {
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			e.w.writeUint16(ui)
		}
	} else {
		//e.s++
		//ui = uint16(e.s)
		ui = uint16(atomic.AddUint32(&e.s, 1))
		e.m[v] = ui
		var lenprec uint8
		switch {
		case l <= math.MaxUint8:
			// lenprec = 0
		case l <= math.MaxUint16:
			lenprec = 1
		case int64(l) <= math.MaxUint32:
			lenprec = 2
		default:
			lenprec = 3
		}
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4|0x0|0x4|lenprec, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			e.w.writeUint16(ui)
		}
		switch lenprec {
		case 0:
			e.w.writen1(byte(l))
		case 1:
			e.w.writeUint16(uint16(l))
		case 2:
			e.w.writeUint32(uint32(l))
		default:
			e.w.writeUint64(uint64(l))
		}
		e.w.writestr(v)
	}
}

func (e *bincEncDriver) encodeStringBytes(c charEncoding, v []byte) {
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 {
		e.w.writeb(v)
	}
}

func (e *bincEncDriver) encBytesLen(c charEncoding, length uint64) {
	//TODO: support bincUnicodeOther (for now, just use string or bytearray)
	if c == c_RAW {
		e.encLen(bincVdByteArray<<4, length)
	} else {
		e.encLen(bincVdString<<4, length)
	}
}

func (e *bincEncDriver) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l+4))
	} else {
		e.encLenNumber(bd, l)
	}
}

func (e *bincEncDriver) encLenNumber(bd byte, v uint64) {
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

func (d *bincDecDriver) initReadNext() {
	if d.bdRead {
		return
	}
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
	d.bdType = detUnset
}

func (d *bincDecDriver) currentEncodedType() decodeEncodedType {
	if d.bdType == detUnset {
		switch d.vd {
		case bincVdSpecial:
			switch d.vs {
			case bincSpNil:
				d.bdType = detNil
			case bincSpFalse, bincSpTrue:
				d.bdType = detBool
			case bincSpNan, bincSpNegInf, bincSpPosInf, bincSpZeroFloat:
				d.bdType = detFloat
			case bincSpZero, bincSpNegOne:
				d.bdType = detInt
			default:
				decErr("currentEncodedType: Unrecognized special value 0x%x", d.vs)
			}
		case bincVdSmallInt:
			d.bdType = detInt
		case bincVdUint:
			d.bdType = detUint
		case bincVdInt:
			d.bdType = detInt
		case bincVdFloat:
			d.bdType = detFloat
		case bincVdSymbol, bincVdString:
			d.bdType = detString
		case bincVdByteArray:
			d.bdType = detBytes
		case bincVdTimestamp:
			d.bdType = detTimestamp
		case bincVdCustomExt:
			d.bdType = detExt
		case bincVdArray:
			d.bdType = detArray
		case bincVdMap:
			d.bdType = detMap
		default:
			decErr("currentEncodedType: Unrecognized d.vd: 0x%x", d.vd)
		}		
	}
	return d.bdType
}

func (d *bincDecDriver) tryDecodeAsNil() bool {
	if d.bd == bincVdSpecial<<4|bincSpNil {
		d.bdRead = false
		return true
	}
	return false
}

func (d *bincDecDriver) isBuiltinType(rt uintptr) bool {
	return rt == timeTypId
}

func (d *bincDecDriver) decodeBuiltinType(rt uintptr, rv reflect.Value) {
	switch rt {
	case timeTypId:
		if d.vd != bincVdTimestamp {
			decErr("Invalid d.vd. Expecting 0x%x. Received: 0x%x", bincVdTimestamp, d.vd)
		}
		tt, err := decodeTime(d.r.readn(int(d.vs)))
		if err != nil {
			panic(err)
		}
		rv.Set(reflect.ValueOf(tt))
		d.bdRead = false
	}
}

func (d *bincDecDriver) decFloatPre(vs, defaultLen byte) {
	if vs&0x8 == 0 {
		d.r.readb(d.b[0:defaultLen])
	} else {
		l := d.r.readn1()
		if l > 8 {
			decErr("At most 8 bytes used to represent float. Received: %v bytes", l)
		}
		for i := l; i < 8; i++ {
			d.b[i] = 0
		}
		d.r.readb(d.b[0:l])
	}
}

func (d *bincDecDriver) decFloat() (f float64) {
	//if true { f = math.Float64frombits(d.r.readUint64()); break; }
	switch vs := d.vs; vs & 0x7 {
	case bincFlBin32:
		d.decFloatPre(vs, 4)
		f = float64(math.Float32frombits(bigen.Uint32(d.b[0:4])))
	case bincFlBin64:
		d.decFloatPre(vs, 8)
		f = math.Float64frombits(bigen.Uint64(d.b[0:8]))
	default:
		decErr("only float32 and float64 are supported. d.vd: 0x%x, d.vs: 0x%x", d.vd, d.vs)
	}
	return
}

func (d *bincDecDriver) decInt() (v int64) {
	// need to inline the code (interface conversion and type assertion expensive)
	switch d.vs {
	case 0:
		v = int64(int8(d.r.readn1()))
	case 1:
		d.r.readb(d.b[6:])
		v = int64(int16(bigen.Uint16(d.b[6:])))
	case 2:
		d.r.readb(d.b[5:])
		if d.b[5]&0x80 == 0 {
			d.b[4] = 0
		} else {
			d.b[4] = 0xff
		}
		v = int64(int32(bigen.Uint32(d.b[4:])))
	case 3:
		d.r.readb(d.b[4:])
		v = int64(int32(bigen.Uint32(d.b[4:])))
	case 4, 5, 6:
		lim := int(7 - d.vs)
		d.r.readb(d.b[lim:])
		var fillval byte = 0
		if d.b[lim]&0x80 != 0 {
			fillval = 0xff
		}
		for i := 0; i < lim; i++ {
			d.b[i] = fillval
		}
		v = int64(bigen.Uint64(d.b[:]))
	case 7:
		d.r.readb(d.b[:])
		v = int64(bigen.Uint64(d.b[:]))
	default:
		decErr("integers with greater than 64 bits of precision not supported")
	}
	return
}

func (d *bincDecDriver) decUint() (v uint64) {
	// need to inline the code (interface conversion and type assertion expensive)
	switch d.vs {
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		d.r.readb(d.b[6:])
		v = uint64(bigen.Uint16(d.b[6:]))
	case 2:
		d.b[4] = 0
		d.r.readb(d.b[5:])
		v = uint64(bigen.Uint32(d.b[4:]))
	case 3:
		d.r.readb(d.b[4:])
		v = uint64(bigen.Uint32(d.b[4:]))
	case 4, 5, 6:
		lim := int(7 - d.vs)
		d.r.readb(d.b[lim:])
		for i := 0; i < lim; i++ {
			d.b[i] = 0
		}
		v = uint64(bigen.Uint64(d.b[:]))
	case 7:
		d.r.readb(d.b[:])
		v = uint64(bigen.Uint64(d.b[:]))
	default:
		decErr("unsigned integers with greater than 64 bits of precision not supported")
	}
	return
}

func (d *bincDecDriver) decIntAny() (i int64) {
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

func (d *bincDecDriver) decodeInt(bitsize uint8) (i int64) {
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

func (d *bincDecDriver) decodeUint(bitsize uint8) (ui uint64) {
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

func (d *bincDecDriver) decodeFloat(chkOverflow32 bool) (f float64) {
	switch d.vd {
	case bincVdSpecial:
		d.bdRead = false
		switch d.vs {
		case bincSpNan:
			return math.NaN()
		case bincSpPosInf:
			return math.Inf(1)
		case bincSpZeroFloat, bincSpZero:
			return
		case bincSpNegInf:
			return math.Inf(-1)
		default:
			decErr("Invalid d.vs decoding float where d.vd=bincVdSpecial: %v", d.vs)
		}		
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
func (d *bincDecDriver) decodeBool() (b bool) {
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

func (d *bincDecDriver) readMapLen() (length int) {
	if d.vd != bincVdMap {
		decErr("Invalid d.vd for map. Expecting 0x%x. Got: 0x%x", bincVdMap, d.vd)
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriver) readArrayLen() (length int) {
	if d.vd != bincVdArray {
		decErr("Invalid d.vd for array. Expecting 0x%x. Got: 0x%x", bincVdArray, d.vd)
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriver) decLen() int {
	if d.vs <= 3 {
		return int(d.decUint())
	}
	return int(d.vs - 4)
}

func (d *bincDecDriver) decodeString() (s string) {
	switch d.vd {
	case bincVdString, bincVdByteArray:
		if length := d.decLen(); length > 0 {
			s = string(d.r.readn(length))
		}
	case bincVdSymbol:
		//from vs: extract numSymbolBytes, containsStringVal, strLenPrecision,
		//extract symbol
		//if containsStringVal, read it and put in map
		//else look in map for string value
		var symbol uint32
		vs := d.vs
		//fmt.Printf(">>>> d.vs: 0b%b, & 0x8: %v, & 0x4: %v\n", d.vs, vs & 0x8, vs & 0x4)
		if vs&0x8 == 0 {
			symbol = uint32(d.r.readn1())
		} else {
			symbol = uint32(d.r.readUint16())
		}
		if d.m == nil {
			d.m = make(map[uint32]string, 16)
		}

		if vs&0x4 == 0 {
			s = d.m[symbol]
		} else {
			var slen int
			switch vs & 0x3 {
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(d.r.readUint16())
			case 2:
				slen = int(d.r.readUint32())
			case 3:
				slen = int(d.r.readUint64())
			}
			s = string(d.r.readn(slen))
			d.m[symbol] = s
		}
	default:
		decErr("Invalid d.vd for string. Expecting string:0x%x, bytearray:0x%x or symbol: 0x%x. Got: 0x%x",
			bincVdString, bincVdByteArray, bincVdSymbol, d.vd)
	}
	d.bdRead = false
	return
}

func (d *bincDecDriver) decodeBytes(bs []byte) (bsOut []byte, changed bool) {
	var clen int
	switch d.vd {
	case bincVdString, bincVdByteArray:
		clen = d.decLen()
	default:
		decErr("Invalid d.vd for bytes. Expecting string:0x%x or bytearray:0x%x. Got: 0x%x",
			bincVdString, bincVdByteArray, d.vd)
	}
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

func (d *bincDecDriver) decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte) {
	switch d.vd {
	case bincVdCustomExt:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			decErr("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		xbs = d.r.readn(l)
	case bincVdByteArray:
		xbs, _ = d.decodeBytes(nil)
	default:
		decErr("Invalid d.vd for extensions (Expecting extensions or byte array). Got: 0x%x", d.vd)
	}
	d.bdRead = false
	return
}

func (d *bincDecDriver) decodeNaked() (rv reflect.Value, ctx decodeNakedContext) {
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
		case bincSpZeroFloat:
			v = float64(0)
		case bincSpZero:
			v = int64(0) // int8(0)
		case bincSpNegOne:
			v = int64(-1) // int8(-1)
		default:
			decErr("decodeNaked: Unrecognized special value 0x%x", d.vs)
		}
	case bincVdSmallInt:
		v = int64(int8(d.vs)) + 1 // int8(d.vs) + 1
	case bincVdUint:
		v = d.decUint()
	case bincVdInt:
		v = d.decInt()
	case bincVdFloat:
		v = d.decFloat()
	case bincVdSymbol:
		v = d.decodeString()
	case bincVdString:
		v = d.decodeString()
	case bincVdByteArray:
		v, _ = d.decodeBytes(nil)
	case bincVdTimestamp:
		tt, err := decodeTime(d.r.readn(int(d.vs)))
		if err != nil {
			panic(err)
		}
		v = tt
	case bincVdCustomExt:
		//ctx = dncExt
		l := d.decLen()
		xtag := d.r.readn1()
		xbs := d.r.readn(l)
		var bfn func(reflect.Value, []byte) error
		rv, bfn = d.h.getDecodeExtForTag(xtag)
		if bfn == nil {
			// decErr("decodeNaked: Unable to find type mapped to extension tag: %v", xtag)
			re := RawExt { xtag, xbs }
			rv = reflect.ValueOf(&re).Elem()
		} else if fnerr := bfn(rv, xbs); fnerr != nil {
			panic(fnerr)
		}
	case bincVdArray:
		ctx = dncContainer
		if d.h.SliceType == nil {
			rv = reflect.New(intfSliceTyp).Elem()
		} else {
			rv = reflect.New(d.h.SliceType).Elem()
		}
	case bincVdMap:
		ctx = dncContainer
		if d.h.MapType == nil {
			rv = reflect.MakeMap(mapIntfIntfTyp)
		} else {
			rv = reflect.MakeMap(d.h.MapType)
		}
	default:
		decErr("decodeNaked: Unrecognized d.vd: 0x%x", d.vd)
	}

	if ctx == dncHandled {
		d.bdRead = false
		if v != nil {
			rv = reflect.ValueOf(v)
		}
	}
	return
}

