// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
)

const (
	cborMajorUint byte = iota
	cborMajorNegInt
	cborMajorBytes
	cborMajorText
	cborMajorArray
	cborMajorMap
	cborMajorTag
	cborMajorOther
)

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

const (
	cborBdIndefiniteBytes  byte = 0x5f
	cborBdIndefiniteString      = 0x7f
	cborBdIndefiniteArray       = 0x9f
	cborBdIndefiniteMap         = 0xbf
	cborBdBreak                 = 0xff
)

const (
	CborStreamBytes  byte = 0x5f
	CborStreamString      = 0x7f
	CborStreamArray       = 0x9f
	CborStreamMap         = 0xbf
	cborStreamBreak       = 0xff
)

const (
	cborBaseUint   byte = 0x00
	cborBaseNegInt      = 0x20
	cborBaseBytes       = 0x40
	cborBaseString      = 0x60
	cborBaseArray       = 0x80
	cborBaseMap         = 0xa0
	cborBaseTag         = 0xc0
	cborBaseSimple      = 0xe0
)

// -------------------

type cborEncDriver struct {
	w encWriter
	h *CborHandle
	noBuiltInTypes
}

func (e *cborEncDriver) encodeNil() {
	e.w.writen1(cborBdNil)
}

func (e *cborEncDriver) encodeBool(b bool) {
	if b {
		e.w.writen1(cborBdTrue)
	} else {
		e.w.writen1(cborBdFalse)
	}
}

func (e *cborEncDriver) encodeFloat32(f float32) {
	e.w.writen1(cborBdFloat32)
	e.w.writeUint32(math.Float32bits(f))
}

func (e *cborEncDriver) encodeFloat64(f float64) {
	e.w.writen1(cborBdFloat64)
	e.w.writeUint64(math.Float64bits(f))
}

func (e *cborEncDriver) encUint(v uint64, bd byte) {
	switch {
	case v <= 0x17:
		e.w.writen1(byte(v) + bd)
	case v <= math.MaxUint8:
		e.w.writen2(bd+0x18, uint8(v))
	case v <= math.MaxUint16:
		e.w.writen1(bd + 0x19)
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		e.w.writen1(bd + 0x1a)
		e.w.writeUint32(uint32(v))
	case v <= math.MaxUint64:
		e.w.writen1(bd + 0x1b)
		e.w.writeUint64(v)
	}
}

func (e *cborEncDriver) encodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-1-v), cborBaseNegInt)
	} else {
		e.encUint(uint64(v), cborBaseUint)
	}
}

func (e *cborEncDriver) encodeUint(v uint64) {
	e.encUint(v, cborBaseUint)
}

func (e *cborEncDriver) encLen(bd byte, length int) {
	e.encUint(uint64(length), bd)
}

func (e *cborEncDriver) encodeExt(rv reflect.Value, xtag uint64, ext Ext, en *Encoder) {
	e.encUint(uint64(xtag), cborBaseTag)
	if v := ext.ConvertExt(rv); v == nil {
		e.encodeNil()
	} else {
		en.encode(v)
	}
}

func (e *cborEncDriver) encodeRawExt(re *RawExt, en *Encoder) {
	e.encUint(uint64(re.Tag), cborBaseTag)
	if re.Data != nil {
		en.encode(re.Data)
	} else if re.Value == nil {
		e.encodeNil()
	} else {
		en.encode(re.Value)
	}
}

func (e *cborEncDriver) encodeArrayPreamble(length int) {
	e.encLen(cborBaseArray, length)
}

func (e *cborEncDriver) encodeMapPreamble(length int) {
	e.encLen(cborBaseMap, length)
}

func (e *cborEncDriver) encodeString(c charEncoding, v string) {
	e.encLen(cborBaseString, len(v))
	e.w.writestr(v)
}

func (e *cborEncDriver) encodeSymbol(v string) {
	e.encodeString(c_UTF8, v)
}

func (e *cborEncDriver) encodeStringBytes(c charEncoding, v []byte) {
	e.encLen(cborBaseBytes, len(v))
	e.w.writeb(v)
}

// ----------------------

type cborDecDriver struct {
	h      *CborHandle
	r      decReader
	bdRead bool
	bdType valueType
	bd     byte
	noBuiltInTypes
}

func (d *cborDecDriver) initReadNext() {
	if d.bdRead {
		return
	}
	d.bd = d.r.readn1()
	d.bdRead = true
	d.bdType = valueTypeUnset
}

func (d *cborDecDriver) currentEncodedType() valueType {
	if d.bdType == valueTypeUnset {
		switch d.bd {
		case cborBdNil:
			d.bdType = valueTypeNil
		case cborBdFalse, cborBdTrue:
			d.bdType = valueTypeBool
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			d.bdType = valueTypeFloat
		case cborBdIndefiniteBytes:
			d.bdType = valueTypeBytes
		case cborBdIndefiniteString:
			d.bdType = valueTypeString
		case cborBdIndefiniteArray:
			d.bdType = valueTypeArray
		case cborBdIndefiniteMap:
			d.bdType = valueTypeMap
		default:
			switch {
			case d.bd >= cborBaseUint && d.bd < cborBaseNegInt:
				if d.h.SignedInteger {
					d.bdType = valueTypeInt
				} else {
					d.bdType = valueTypeUint
				}
			case d.bd >= cborBaseNegInt && d.bd < cborBaseBytes:
				d.bdType = valueTypeInt
			case d.bd >= cborBaseBytes && d.bd < cborBaseString:
				d.bdType = valueTypeBytes
			case d.bd >= cborBaseString && d.bd < cborBaseArray:
				d.bdType = valueTypeString
			case d.bd >= cborBaseArray && d.bd < cborBaseMap:
				d.bdType = valueTypeArray
			case d.bd >= cborBaseMap && d.bd < cborBaseTag:
				d.bdType = valueTypeMap
			case d.bd >= cborBaseTag && d.bd < cborBaseSimple:
				d.bdType = valueTypeExt
			default:
				decErr("currentEncodedType: Unrecognized d.bd: 0x%x", d.bd)
			}
		}
	}
	return d.bdType
}

func (d *cborDecDriver) tryDecodeAsNil() bool {
	// treat Nil and Undefined as nil values
	if d.bd == cborBdNil || d.bd == cborBdUndefined {
		d.bdRead = false
		return true
	}
	return false
}

func (d *cborDecDriver) checkBreak() bool {
	d.initReadNext()
	if d.bd == cborBdBreak {
		d.bdRead = false
		return true
	}
	return false
}

func (d *cborDecDriver) decUint() (ui uint64) {
	v := d.bd & 0x1f
	if v <= 0x17 {
		ui = uint64(v)
	} else {
		switch v {
		case 0x18:
			ui = uint64(d.r.readn1())
		case 0x19:
			ui = uint64(d.r.readUint16())
		case 0x1a:
			ui = uint64(d.r.readUint32())
		case 0x1b:
			ui = uint64(d.r.readUint64())
		default:
			decErr("decUint: Invalid descriptor: %v", d.bd)
		}
	}
	return
}

func (d *cborDecDriver) decCheckInteger() (neg bool) {
	switch major := d.bd >> 5; major {
	case cborMajorUint:
	case cborMajorNegInt:
		neg = true
	default:
		decErr("invalid major: %v (bd: %v)", major, d.bd)
	}
	return
}

func (d *cborDecDriver) decodeInt(bitsize uint8) (i int64) {
	neg := d.decCheckInteger()
	ui := d.decUint()
	// check if this number can be converted to an int without overflow
	if neg {
		i = -checkOverflowUint64ToInt64(ui + 1)
	} else {
		i = checkOverflowUint64ToInt64(ui)
	}
	checkOverflow(0, i, bitsize)
	d.bdRead = false
	return
}

func (d *cborDecDriver) decodeUint(bitsize uint8) (ui uint64) {
	if d.decCheckInteger() {
		decErr("Assigning negative signed value to unsigned type")
	}
	ui = d.decUint()
	checkOverflow(ui, 0, bitsize)
	d.bdRead = false
	return
}

func (d *cborDecDriver) decodeFloat(chkOverflow32 bool) (f float64) {
	switch d.bd {
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(d.r.readUint16())))
	case cborBdFloat32:
		f = float64(math.Float32frombits(d.r.readUint32()))
	case cborBdFloat64:
		f = math.Float64frombits(d.r.readUint64())
	default:
		if d.bd >= cborBaseUint && d.bd < cborBaseBytes {
			f = float64(d.decodeInt(64))
		} else {
			decErr("Float only valid from float16/32/64: Invalid descriptor: %v", d.bd)
		}
	}
	checkOverflowFloat32(f, chkOverflow32)
	d.bdRead = false
	return
}

// bool can be decoded from bool only (single byte).
func (d *cborDecDriver) decodeBool() (b bool) {
	switch d.bd {
	case cborBdTrue:
		b = true
	case cborBdFalse:
	default:
		decErr("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *cborDecDriver) readMapLen() (length int) {
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap {
		return -1
	}
	return d.decLen()
}

func (d *cborDecDriver) readArrayLen() (length int) {
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray {
		return -1
	}
	return d.decLen()
}

func (d *cborDecDriver) decLen() int {
	return int(d.decUint())
}

func (d *cborDecDriver) decIndefiniteBytes(bs []byte) []byte {
	d.bdRead = false
	for {
		if d.checkBreak() {
			break
		}
		if major := d.bd >> 5; major != cborMajorBytes && major != cborMajorText {
			decErr("cbor: expect bytes or string major type in indefinite string/bytes; got: %v, byte: %v", major, d.bd)
		}
		bs = append(bs, d.r.readn(d.decLen())...)
		d.bdRead = false
	}
	d.bdRead = false
	return bs
}

func (d *cborDecDriver) decodeString() (s string) {
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		s = string(d.decIndefiniteBytes(nil))
		return
	}
	s = string(d.r.readn(d.decLen()))
	d.bdRead = false
	return
}

func (d *cborDecDriver) decodeBytes(bs []byte) (bsOut []byte, changed bool) {
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		if bs == nil {
			bs = []byte{}
			changed = true
		}
		bsOut = d.decIndefiniteBytes(bs[:0])
		if len(bsOut) != len(bs) {
			changed = true
		}
		return
	}
	clen := d.decLen()
	if clen >= 0 {
		if bs == nil {
			bs = []byte{}
			bsOut = bs
			changed = true
		}
		if len(bs) != clen {
			if len(bs) > clen {
				bs = bs[:clen]
			} else {
				bs = make([]byte, clen)
			}
			bsOut = bs
			changed = true
		}
		if len(bs) > 0 {
			d.r.readb(bs)
		}
	}
	d.bdRead = false
	return
}

func (d *cborDecDriver) decodeExt(rv reflect.Value, xtag uint64, ext Ext, de *Decoder) (realxtag uint64) {
	u := d.decUint()
	d.bdRead = false
	realxtag = u
	if ext == nil {
		re := rv.Interface().(*RawExt)
		re.Tag = realxtag
		de.decode(&re.Value)
	} else if xtag != realxtag {
		decErr("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
	} else {
		var v interface{}
		de.decode(&v)
		ext.UpdateExt(rv, v)
	}
	d.bdRead = false
	return
}

func (d *cborDecDriver) decodeNaked(de *Decoder) (v interface{}, vt valueType, decodeFurther bool) {
	d.initReadNext()

	switch d.bd {
	case cborBdNil:
		vt = valueTypeNil
	case cborBdFalse:
		vt = valueTypeBool
		v = false
	case cborBdTrue:
		vt = valueTypeBool
		v = true
	case cborBdFloat16, cborBdFloat32:
		vt = valueTypeFloat
		v = d.decodeFloat(true)
	case cborBdFloat64:
		vt = valueTypeFloat
		v = d.decodeFloat(false)
	case cborBdIndefiniteBytes:
		vt = valueTypeBytes
		v, _ = d.decodeBytes(nil)
	case cborBdIndefiniteString:
		vt = valueTypeString
		v = d.decodeString()
	case cborBdIndefiniteArray:
		vt = valueTypeArray
		decodeFurther = true
	case cborBdIndefiniteMap:
		vt = valueTypeMap
		decodeFurther = true
	default:
		switch {
		case d.bd >= cborBaseUint && d.bd < cborBaseNegInt:
			if d.h.SignedInteger {
				vt = valueTypeInt
				v = d.decodeInt(64)
			} else {
				vt = valueTypeUint
				v = d.decodeUint(64)
			}
		case d.bd >= cborBaseNegInt && d.bd < cborBaseBytes:
			vt = valueTypeInt
			v = d.decodeInt(64)
		case d.bd >= cborBaseBytes && d.bd < cborBaseString:
			vt = valueTypeBytes
			v, _ = d.decodeBytes(nil)
		case d.bd >= cborBaseString && d.bd < cborBaseArray:
			vt = valueTypeString
			v = d.decodeString()
		case d.bd >= cborBaseArray && d.bd < cborBaseMap:
			vt = valueTypeArray
			decodeFurther = true
		case d.bd >= cborBaseMap && d.bd < cborBaseTag:
			vt = valueTypeMap
			decodeFurther = true
		case d.bd >= cborBaseTag && d.bd < cborBaseSimple:
			vt = valueTypeExt
			var re RawExt
			ui := d.decUint()
			d.bdRead = false
			re.Tag = ui
			de.decode(&re.Value)
			v = &re
			// decodeFurther = true
		default:
			decErr("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		}
	}

	if !decodeFurther {
		d.bdRead = false
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
//   - timestamp, BigNum, BigFloat, Decimals, Encoded Text (e.g. URL, regexp, base64, MIME Message), etc.
//
// To encode with indefinite lengths (streaming), users will use
// (Must)Write and (Must)Encode methods of *Encoder, along with CborStreamXXX constants.
//
// For example, to encode "one-byte" as an indefinite length string:
//     var buf bytes.Buffer
//     e := NewEncoder(&buf, new(CborHandle))
//     e.MustWrite([]byte{CborStreamString})
//     e.MustEncode("one-")
//     e.MustEncode("byte")
//     e.MustWrite([]byte{CborStreamBreak})
//     encodedBytes := buf.Bytes()
//     var vv interface{}
//     NewDecoderBytes(buf.Bytes(), new(CborHandle)).MustDecode(&vv)
//     // Now, vv contains the same string "one-byte"
//
type CborHandle struct {
	BasicHandle
}

func (h *CborHandle) newEncDriver(w encWriter) encDriver {
	return &cborEncDriver{w: w, h: h}
}

func (h *CborHandle) newDecDriver(r decReader) decDriver {
	return &cborDecDriver{r: r, h: h}
}

var _ decDriver = (*cborDecDriver)(nil)
var _ encDriver = (*cborEncDriver)(nil)
