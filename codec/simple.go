// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
	"math"
	"reflect"
	"time"
)

// func init() {
// 	var t *SimpleHandle
// 	h := handleNewFn{
// 		typ: rt2id(reflect.TypeOf(t)),
// 		encBytes: func(v *[]byte, h Handle) encoderI {
// 			return newEncDriverBytes[simpleEncDriverM[bytesEncAppenderM]](v, h)
// 		},
// 		encIO:    newEncDriverIO[simpleEncDriverM[bufioEncAppenderM], simpleEncDriverM[bytesEncAppenderM]],
// 		decBytes: newDecDriverBytes[simpleDecDriverM[bytesDecReaderM]],
// 		decIO:    newDecDriverIO[simpleDecDriverM[ioDecReaderM], simpleDecDriverM[bytesDecReaderM]],
// 	}
// 	handleNewFns = append(handleNewFns, h)
// }

const (
	_               uint8 = iota
	simpleVdNil           = 1
	simpleVdFalse         = 2
	simpleVdTrue          = 3
	simpleVdFloat32       = 4
	simpleVdFloat64       = 5

	// each lasts for 4 (ie n, n+1, n+2, n+3)
	simpleVdPosInt = 8
	simpleVdNegInt = 12

	simpleVdTime = 24

	// containers: each lasts for 4 (ie n, n+1, n+2, ... n+7)
	simpleVdString    = 216
	simpleVdByteArray = 224
	simpleVdArray     = 232
	simpleVdMap       = 240
	simpleVdExt       = 248
)

var simpledescNames = map[byte]string{
	simpleVdNil:     "null",
	simpleVdFalse:   "false",
	simpleVdTrue:    "true",
	simpleVdFloat32: "float32",
	simpleVdFloat64: "float64",

	simpleVdPosInt: "+int",
	simpleVdNegInt: "-int",

	simpleVdTime: "time",

	simpleVdString:    "string",
	simpleVdByteArray: "binary",
	simpleVdArray:     "array",
	simpleVdMap:       "map",
	simpleVdExt:       "ext",
}

func simpledesc(bd byte) (s string) {
	s = simpledescNames[bd]
	if s == "" {
		s = "unknown"
	}
	return
}

type simpleEncDriver[T encWriter] struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *SimpleHandle
	e *encoderShared
	// b [8]byte
	w T

	// bytes bool

	// // we cannot reference a *encoder here, due to recursive limitations of go generics
	// // this MUST be *encoder[simpleEncDriverM[bytesEncAppenderM]]
	// // es encoderI // must-be *encoder[simpleEncDriverM[bytesEncAppenderM]]
	// // es *encoder[simpleEncDriverM[bytesEncAppenderM]]
	// es interface{}
}

func (e *simpleEncDriver[T]) EncodeNil() {
	e.w.writen1(simpleVdNil)
}

func (e *simpleEncDriver[T]) EncodeBool(b bool) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && !b {
		e.EncodeNil()
		return
	}
	if b {
		e.w.writen1(simpleVdTrue)
	} else {
		e.w.writen1(simpleVdFalse)
	}
}

func (e *simpleEncDriver[T]) EncodeFloat32(f float32) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat32)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *simpleEncDriver[T]) EncodeFloat64(f float64) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *simpleEncDriver[T]) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), simpleVdNegInt)
	} else {
		e.encUint(uint64(v), simpleVdPosInt)
	}
}

func (e *simpleEncDriver[T]) EncodeUint(v uint64) {
	e.encUint(v, simpleVdPosInt)
}

func (e *simpleEncDriver[T]) encUint(v uint64, bd uint8) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == 0 {
		e.EncodeNil()
		return
	}
	if v <= math.MaxUint8 {
		e.w.writen2(bd, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 1)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 2)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else { // if v <= math.MaxUint64 {
		e.w.writen1(bd + 3)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *simpleEncDriver[T]) encLen(bd byte, length int) {
	if length == 0 {
		e.w.writen1(bd)
	} else if length <= math.MaxUint8 {
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	} else if length <= math.MaxUint16 {
		e.w.writen1(bd + 2)
		e.w.writen2(bigen.PutUint16(uint16(length)))
	} else if int64(length) <= math.MaxUint32 {
		e.w.writen1(bd + 3)
		e.w.writen4(bigen.PutUint32(uint32(length)))
	} else {
		e.w.writen1(bd + 4)
		e.w.writen8(bigen.PutUint64(uint64(length)))
	}
}

func (e *simpleEncDriver[T]) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	var bs0, bs []byte
	if ext == SelfExt {
		bs0 = e.e.blist.get(1024)
		bs = bs0
		e.sideEncoder(&bs)
		e.sideEncode(v, basetype, 0)
	} else {
		bs = ext.WriteExt(v)
	}
	if bs == nil {
		e.EncodeNil()
		goto END
	}
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *simpleEncDriver[T]) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *simpleEncDriver[T]) encodeExtPreamble(xtag byte, length int) {
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
}

func (e *simpleEncDriver[T]) WriteArrayStart(length int) {
	e.encLen(simpleVdArray, length)
}

func (e *simpleEncDriver[T]) WriteMapStart(length int) {
	e.encLen(simpleVdMap, length)
}

func (e *simpleEncDriver[T]) EncodeString(v string) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == "" {
		e.EncodeNil()
		return
	}
	if e.h.StringToRaw {
		e.encLen(simpleVdByteArray, len(v))
	} else {
		e.encLen(simpleVdString, len(v))
	}
	e.w.writestr(v)
}

func (e *simpleEncDriver[T]) EncodeStringBytesRaw(v []byte) {
	// if e.h.EncZeroValuesAsNil && e.c != containerMapKey && v == nil {
	if v == nil {
		e.EncodeNil()
		return
	}
	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
}

func (e *simpleEncDriver[T]) EncodeTime(t time.Time) {
	// if e.h.EncZeroValuesAsNil && e.c != containerMapKey && t.IsZero() {
	if t.IsZero() {
		e.EncodeNil()
		return
	}
	v, err := t.MarshalBinary()
	halt.onerror(err)
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
}

//------------------------------------

type simpleDecDriver[T decReader] struct {
	h *SimpleHandle
	d *decoderShared
	r T

	rh helperDecReader[T]

	bdAndBdread
	bytes bool

	noBuiltInTypes
	decDriverNoopNumberHelper
	decDriverNoopContainerReader
	decInit2er

	// ds interface{} // must be *decoder[simpleDecDriverM[bytes...]]
}

func (d *simpleDecDriver[T]) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *simpleDecDriver[T]) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true // null = true
	}
	return
}

func (d *simpleDecDriver[T]) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	switch d.bd {
	case simpleVdNil:
		d.bdRead = false
		return valueTypeNil
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		return valueTypeBytes
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		return valueTypeString
	case simpleVdArray, simpleVdArray + 1,
		simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		return valueTypeArray
	case simpleVdMap, simpleVdMap + 1,
		simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *simpleDecDriver[T]) TryNil() bool {
	return d.advanceNil()
}

func (d *simpleDecDriver[T]) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case simpleVdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *simpleDecDriver[T]) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case simpleVdPosInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case simpleVdPosInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
		neg = true
	default:
		ok = false
		// halt.errorf("integer only valid from pos/neg integer1..8. Invalid descriptor: %v", d.bd)
	}
	// DO NOT do this check below, because callers may only want the unsigned value:
	//
	// if ui > math.MaxInt64 {
	// 	halt.errorf("decIntAny: Integer out of range for signed int64: %v", ui)
	//		return
	// }
	return
}

func (d *simpleDecDriver[T]) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

// bool can be decoded from bool only (single byte).
func (d *simpleDecDriver[T]) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == simpleVdFalse {
	} else if d.bd == simpleVdTrue {
		b = true
	} else {
		halt.errorf("cannot decode bool - %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriver[T]) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriver[T]) uint2Len(ui uint64) int {
	if chkOvf.Uint(ui, intBitsize) {
		halt.errorf("overflow integer: %v", ui)
	}
	return int(ui)
}

func (d *simpleDecDriver[T]) decLen() int {
	switch d.bd & 7 { // d.bd % 8 {
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(bigen.Uint16(d.r.readn2()))
	case 3:
		return d.uint2Len(uint64(bigen.Uint32(d.r.readn4())))
	case 4:
		return d.uint2Len(bigen.Uint64(d.r.readn8()))
	}
	halt.errorf("cannot read length: bd%%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
}

func (d *simpleDecDriver[T]) DecodeStringAsBytes() (s []byte) {
	return d.DecodeBytes(nil)
}

func (d *simpleDecDriver[T]) DecodeBytes(bs []byte) (bsOut []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.bd >= simpleVdArray && d.bd <= simpleVdMap+4 {
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			bs = d.d.b[:]
		}
		slen := d.ReadArrayStart()
		var changed bool
		if bs, changed = usableByteSlice(bs, slen); changed {
			d.d.decByteState = decByteStateNone
		}
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}

	clen := d.decLen()
	d.bdRead = false
	// if d.d.zerocopy() {
	if d.bytes && d.h.ZeroCopy {
		d.d.decByteState = decByteStateZerocopy
		return d.r.readx(uint(clen))
	}
	if bs == nil {
		d.d.decByteState = decByteStateReuseBuf
		bs = d.d.b[:]
	}
	return d.rh.decByteSlice(d.r, clen, d.h.MaxInitLen, bs)
}

func (d *simpleDecDriver[T]) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd != simpleVdTime {
		halt.errorf("invalid descriptor for time.Time - expect 0x%x, received 0x%x", simpleVdTime, d.bd)
	}
	d.bdRead = false
	clen := uint(d.r.readn1())
	b := d.r.readx(clen)
	halt.onerror((&t).UnmarshalBinary(b))
	return
}

func (d *simpleDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if xtag > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtag)
	}
	if d.advanceNil() {
		return
	}
	xbs, realxtag1, zerocopy := d.decodeExtV(ext != nil, uint8(xtag))
	realxtag := uint64(realxtag1)
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = realxtag
		re.setData(xbs, zerocopy)
	} else if ext == SelfExt {
		d.sideDecoder(xbs)
		d.sideDecode(rv, basetype)
	} else {
		ext.ReadExt(rv, xbs)
	}
}

func (d *simpleDecDriver[T]) decodeExtV(verifyTag bool, tag byte) (xbs []byte, xtag byte, zerocopy bool) {
	switch d.bd {
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		if d.d.bytes {
			xbs = d.r.readx(uint(l))
			zerocopy = true
		} else {
			xbs = d.rh.decByteSlice(d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs = d.DecodeBytes(nil)
	default:
		halt.errorf("ext - %s - expecting extensions/bytearray, got: 0x%x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd {
	case simpleVdNil:
		n.v = valueTypeNil
	case simpleVdFalse:
		n.v = valueTypeBool
		n.b = false
	case simpleVdTrue:
		n.v = valueTypeBool
		n.b = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case simpleVdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdTime:
		n.v = valueTypeTime
		n.t = d.DecodeTime()
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString, d.h.ZeroCopy)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		if d.d.bytes {
			n.l = d.r.readx(uint(l))
		} else {
			n.l = d.rh.decByteSlice(d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2,
		simpleVdArray + 3, simpleVdArray + 4:
		n.v = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, d.bd)
	}

	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *simpleDecDriver[T]) nextValueBytes(v0 []byte) (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	v = v0
	var h decNextValueBytesHelper
	// var h = decNextValueBytesHelper[T]{d}
	var cursor uint
	if d.bytes {
		cursor = d.r.numread() - 1
	}
	h.append1(&v, d.bytes, d.bd)
	v = d.nextValueBytesBdReadR(v)
	d.bdRead = false
	// h.bytesRdV(&v, d.bytes, cursor)
	if d.bytes {
		v = d.r.bytesReadFrom(cursor)
	}
	return
}

func (d *simpleDecDriver[T]) nextValueBytesR(v0 []byte) (v []byte) {
	d.readNextBd()
	v = v0
	var h decNextValueBytesHelper
	// var h = decNextValueBytesHelper{d: &d.d}
	h.append1(&v, d.bytes, d.bd)
	return d.nextValueBytesBdReadR(v)
}

func (d *simpleDecDriver[T]) nextValueBytesBdReadR(v0 []byte) (v []byte) {
	v = v0
	var h decNextValueBytesHelper
	// var h = decNextValueBytesHelper{d: &d.d}

	c := d.bd

	var length uint

	switch c {
	case simpleVdNil, simpleVdFalse, simpleVdTrue, simpleVdString, simpleVdByteArray:
		// pass
	case simpleVdPosInt, simpleVdNegInt:
		h.append1(&v, d.bytes, d.r.readn1())
	case simpleVdPosInt + 1, simpleVdNegInt + 1:
		h.appendN(&v, d.bytes, d.r.readx(2)...)
	case simpleVdPosInt + 2, simpleVdNegInt + 2, simpleVdFloat32:
		h.appendN(&v, d.bytes, d.r.readx(4)...)
	case simpleVdPosInt + 3, simpleVdNegInt + 3, simpleVdFloat64:
		h.appendN(&v, d.bytes, d.r.readx(8)...)
	case simpleVdTime:
		c = d.r.readn1()
		h.append1(&v, d.bytes, c)
		h.appendN(&v, d.bytes, d.r.readx(uint(c))...)

	default:
		switch c & 7 { // c % 8 {
		case 0:
			length = 0
		case 1:
			b := d.r.readn1()
			length = uint(b)
			h.append1(&v, d.bytes, b)
		case 2:
			x := d.r.readn2()
			length = uint(bigen.Uint16(x))
			h.appendN(&v, d.bytes, x[:]...)
		case 3:
			x := d.r.readn4()
			length = uint(bigen.Uint32(x))
			h.appendN(&v, d.bytes, x[:]...)
		case 4:
			x := d.r.readn8()
			length = uint(bigen.Uint64(x))
			h.appendN(&v, d.bytes, x[:]...)
		}

		bExt := c >= simpleVdExt && c <= simpleVdExt+7
		bStr := c >= simpleVdString && c <= simpleVdString+7
		bByteArray := c >= simpleVdByteArray && c <= simpleVdByteArray+7
		bArray := c >= simpleVdArray && c <= simpleVdArray+7
		bMap := c >= simpleVdMap && c <= simpleVdMap+7

		if !(bExt || bStr || bByteArray || bArray || bMap) {
			halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, c)
		}

		if bExt {
			h.append1(&v, d.bytes, d.r.readn1()) // tag
		}

		if length == 0 {
			break
		}

		if bArray {
			for i := uint(0); i < length; i++ {
				v = d.nextValueBytesR(v)
			}
		} else if bMap {
			for i := uint(0); i < length; i++ {
				v = d.nextValueBytesR(v)
				v = d.nextValueBytesR(v)
			}
		} else {
			h.appendN(&v, d.bytes, d.r.readx(length)...)
		}
	}
	return
}

//------------------------------------

// SimpleHandle is a Handle for a very simple encoding format.
//
// simple is a simplistic codec similar to binc, but not as compact.
//   - Encoding of a value is always preceded by the descriptor byte (bd)
//   - True, false, nil are encoded fully in 1 byte (the descriptor)
//   - Integers (intXXX, uintXXX) are encoded in 1, 2, 4 or 8 bytes (plus a descriptor byte).
//     There are positive (uintXXX and intXXX >= 0) and negative (intXXX < 0) integers.
//   - Floats are encoded in 4 or 8 bytes (plus a descriptor byte)
//   - Length of containers (strings, bytes, array, map, extensions)
//     are encoded in 0, 1, 2, 4 or 8 bytes.
//     Zero-length containers have no length encoded.
//     For others, the number of bytes is given by pow(2, bd%3)
//   - maps are encoded as [bd] [length] [[key][value]]...
//   - arrays are encoded as [bd] [length] [value]...
//   - extensions are encoded as [bd] [length] [tag] [byte]...
//   - strings/bytearrays are encoded as [bd] [length] [byte]...
//   - time.Time are encoded as [bd] [length] [byte]...
//
// The full spec will be published soon.
type SimpleHandle struct {
	binaryEncodingType
	BasicHandle

	// EncZeroValuesAsNil says to encode zero values for numbers, bool, string, etc as nil
	EncZeroValuesAsNil bool
}

// func (h *SimpleHandle) init() {
// 	// h.rt.init()
// }

// Name returns the name of the handle: simple
func (h *SimpleHandle) Name() string { return "simple" }

func (h *SimpleHandle) desc(bd byte) string { return simpledesc(bd) }

// func (h *SimpleHandle) clone() interface{} {
// 	h2 := *h
// 	return &h2
// }

// SetBytesExt sets an extension
func (h *SimpleHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *simpleEncDriver[T]) init(hh Handle, shared *encoderShared, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*SimpleHandle)
	d.e = shared
	if shared.bytes {
		fp = simpleFpEncBytes
	} else {
		fp = simpleFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *simpleEncDriver[T]) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *simpleEncDriver[T]) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *simpleEncDriver[T]) writerEnd()                        { e.w.end() }

func (e *simpleEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *simpleEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *simpleDecDriver[T]) init(hh Handle, shared *decoderShared, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*SimpleHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = simpleFpDecBytes
	} else {
		fp = simpleFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *simpleDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *simpleDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *simpleDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

// ---- (custom stanza)

func (d *simpleDecDriver[T]) descBd() string {
	return sprintf("%v (%s)", d.bd, simpledesc(d.bd))
}

func (d *simpleDecDriver[T]) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
