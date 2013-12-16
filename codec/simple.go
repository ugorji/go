// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

/*
simple is a simplistic codec which does not pack integers, floats, etc into smaller chunks.

All numbers are encoded with 8 bytes (int64, uint64, float64).

Beyond this, it looks very similar to binc or msgpack.

- nil, false, true are represented with single bytes
- positive/negative ints, and floats are distinct types, all represented by 1,2,4 or 8 bytes
- lengths of containers (strings, arrays, ext, map) occupy 2 bytes and
  exist right after the bd
- maps are encoded as [bd] [length] [[key][value]]...
- arrays are encoded as [bd] [length] [value]...
- extensions are encoded as [bd] [length] [tag] [byte]...
- strings/bytearrays are encoded as [bd] [length] [tag] [byte]...

*/

package codec

import "math"

const (
	_ byte = iota
	simpleVdNil
	simpleVdFalse
	simpleVdTrue
	simpleVdPosInt1
	simpleVdPosInt2
	simpleVdPosInt4
	simpleVdPosInt8
	simpleVdNegInt1
	simpleVdNegInt2
	simpleVdNegInt4
	simpleVdNegInt8
	simpleVdFloat32
	simpleVdFloat64

	simpleVdString
	simpleVdByteArray
	simpleVdArray
	simpleVdMap

	simpleVdExt
)

type simpleEncDriver struct {
	w encWriter
	b [8]byte
}

func (e *simpleEncDriver) isBuiltinType(rt uintptr) bool {
	return false
}

func (e *simpleEncDriver) encodeBuiltin(rt uintptr, v interface{}) {
}

func (e *simpleEncDriver) encodeNil() {
	e.w.writen1(simpleVdNil)
	// e.w.writen1(0)
}

func (e *simpleEncDriver) encodeBool(b bool) {
	if b {
		e.w.writen1(simpleVdTrue)
	} else {
		e.w.writen1(simpleVdFalse)
	}
}

func (e *simpleEncDriver) encodeFloat32(f float32) {
	e.w.writen1(simpleVdFloat32)
	e.w.writeUint32(math.Float32bits(f))
}

func (e *simpleEncDriver) encodeFloat64(f float64) {
	e.w.writen1(simpleVdFloat64)
	e.w.writeUint64(math.Float64bits(f))
}

func (e *simpleEncDriver) encodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), false)
	} else {
		e.encUint(uint64(v), true)
	}
}

func (e *simpleEncDriver) encodeUint(v uint64) {
	e.encUint(v, true)
}

func (e *simpleEncDriver) encUint(v uint64, pos bool) {
	switch {
	case v <= math.MaxUint8:
		if pos {
			e.w.writen1(simpleVdPosInt1)
		} else {
			e.w.writen1(simpleVdNegInt1)
		}
		e.w.writen1(uint8(v))
	case v <= math.MaxUint16:
		if pos {
			e.w.writen1(simpleVdPosInt2)
		} else {
			e.w.writen1(simpleVdNegInt2)
		}
		e.w.writeUint16(uint16(v))
	case v <= math.MaxUint32:
		if pos {
			e.w.writen1(simpleVdPosInt4)
		} else {
			e.w.writen1(simpleVdNegInt4)
		}
		e.w.writeUint32(uint32(v))
	case v <= math.MaxUint64:
		if pos {
			e.w.writen1(simpleVdPosInt8)
		} else {
			e.w.writen1(simpleVdNegInt8)
		}
		e.w.writeUint64(v)
	}
}

func (e *simpleEncDriver) encodeExtPreamble(xtag byte, length int) {
	e.w.writen1(simpleVdExt)
	e.w.writeUint32(uint32(length))
	e.w.writen1(xtag)
}

func (e *simpleEncDriver) encodeArrayPreamble(length int) {
	e.w.writen1(simpleVdArray)
	e.w.writeUint32(uint32(length))
}

func (e *simpleEncDriver) encodeMapPreamble(length int) {
	e.w.writen1(simpleVdMap)
	e.w.writeUint32(uint32(length))
}

func (e *simpleEncDriver) encodeString(c charEncoding, v string) {
	e.w.writen1(simpleVdString)
	e.w.writeUint32(uint32(len(v)))
	e.w.writestr(v)
}

func (e *simpleEncDriver) encodeSymbol(v string) {
	e.encodeString(c_UTF8, v)
}

func (e *simpleEncDriver) encodeStringBytes(c charEncoding, v []byte) {
	e.w.writen1(simpleVdByteArray)
	e.w.writeUint32(uint32(len(v)))
	e.w.writeb(v)
}

//------------------------------------

type simpleDecDriver struct {
	r      decReader
	bdRead bool
	bdType valueType
	bd     byte
	b      [8]byte
}

func (d *simpleDecDriver) initReadNext() {
	if d.bdRead {
		return
	}
	d.bd = d.r.readn1()
	d.bdRead = true
	d.bdType = valueTypeUnset
}

func (d *simpleDecDriver) currentEncodedType() valueType {
	if d.bdType == valueTypeUnset {
		switch d.bd {
		case simpleVdNil:
			d.bdType = valueTypeNil
		case simpleVdTrue, simpleVdFalse:
			d.bdType = valueTypeBool
		case simpleVdPosInt1, simpleVdPosInt2, simpleVdPosInt4, simpleVdPosInt8:
			d.bdType = valueTypeUint
		case simpleVdNegInt1, simpleVdNegInt2, simpleVdNegInt4, simpleVdNegInt8:
			d.bdType = valueTypeInt
		case simpleVdFloat32, simpleVdFloat64:
			d.bdType = valueTypeFloat
		case simpleVdString:
			d.bdType = valueTypeString
		case simpleVdByteArray:
			d.bdType = valueTypeBytes
		case simpleVdExt:
			d.bdType = valueTypeExt
		case simpleVdArray:
			d.bdType = valueTypeArray
		case simpleVdMap:
			d.bdType = valueTypeMap
		default:
			decErr("currentEncodedType: Unrecognized d.vd: 0x%x", d.bd)
		}
	}
	return d.bdType
}

func (d *simpleDecDriver) tryDecodeAsNil() bool {
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true
	}
	return false
}

func (d *simpleDecDriver) isBuiltinType(rt uintptr) bool {
	return false
}

func (d *simpleDecDriver) decodeBuiltin(rt uintptr, v interface{}) {
}

func (d *simpleDecDriver) decIntAny() (ui uint64, i int64, neg bool) {
	switch d.bd {
	case simpleVdPosInt1:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt2:
		ui = uint64(d.r.readUint16())
	case simpleVdPosInt4:
		ui = uint64(d.r.readUint32())
	case simpleVdPosInt8:
		ui = uint64(d.r.readUint64())
	case simpleVdNegInt1:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt2:
		ui = uint64(d.r.readUint16())
		neg = true
	case simpleVdNegInt4:
		ui = uint64(d.r.readUint32())
		neg = true
	case simpleVdNegInt8:
		ui = uint64(d.r.readUint64())
		neg = true
	default:
		decErr("Integer only valid from pos/neg integer1..8. Invalid descriptor: %v", d.bd)
	}
	if neg {
		i = -(int64(ui))
	} else {
		i = int64(ui)
	}
	return
}

func (d *simpleDecDriver) decodeInt(bitsize uint8) (i int64) {
	_, i, _ = d.decIntAny()
	checkOverflow(0, i, bitsize)
	d.bdRead = false
	return
}

func (d *simpleDecDriver) decodeUint(bitsize uint8) (ui uint64) {
	ui, i, neg := d.decIntAny()
	if neg {
		decErr("Assigning negative signed value: %v, to unsigned type", i)
	}
	checkOverflow(ui, 0, bitsize)
	d.bdRead = false
	return
}

func (d *simpleDecDriver) decodeFloat(chkOverflow32 bool) (f float64) {
	switch d.bd {
	case simpleVdFloat32:
		f = float64(math.Float32frombits(d.r.readUint32()))
	case simpleVdFloat64:
		f = math.Float64frombits(d.r.readUint64())
	default:
		if d.bd >= simpleVdPosInt1 && d.bd <= simpleVdNegInt8 {
			_, i, _ := d.decIntAny()
			f = float64(i)
		} else {
			decErr("Float only valid from float32/64: Invalid descriptor: %v", d.bd)
		}
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
	return
}

// bool can be decoded from bool only (single byte).
func (d *simpleDecDriver) decodeBool() (b bool) {
	b = d.r.readn1() != 0
	d.bdRead = false
	return
}

func (d *simpleDecDriver) readMapLen() (length int) {
	return d.decLen()
}

func (d *simpleDecDriver) readArrayLen() (length int) {
	return d.decLen()
}

func (d *simpleDecDriver) decLen() int {
	d.r.readb(d.b[:4])
	return int(bigen.Uint32(d.b[:4]))
}

func (d *simpleDecDriver) decodeString() (s string) {
	s = string(d.r.readn(d.decLen()))
	d.bdRead = false
	return
}

func (d *simpleDecDriver) decodeBytes(bs []byte) (bsOut []byte, changed bool) {
	if clen := d.decLen(); clen > 0 {
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

func (d *simpleDecDriver) decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte) {
	switch d.bd {
	case simpleVdExt:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			decErr("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		xbs = d.r.readn(l)
	case simpleVdByteArray:
		xbs, _ = d.decodeBytes(nil)
	default:
		decErr("Invalid d.vd for extensions (Expecting extensions or byte array). Got: 0x%x", d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriver) decodeNaked() (v interface{}, vt valueType, decodeFurther bool) {
	d.initReadNext()

	switch d.bd {
	case simpleVdNil:
		vt = valueTypeNil
	case simpleVdFalse:
		vt = valueTypeBool
		v = false
	case simpleVdTrue:
		vt = valueTypeBool
		v = true
	case simpleVdPosInt1, simpleVdPosInt2, simpleVdPosInt4, simpleVdPosInt8:
		vt = valueTypeUint
		ui, _, _ := d.decIntAny()
		v = ui
	case simpleVdNegInt1, simpleVdNegInt2, simpleVdNegInt4, simpleVdNegInt8:
		vt = valueTypeInt
		_, i, _ := d.decIntAny()
		v = i
	case simpleVdFloat32:
		vt = valueTypeFloat
		v = d.decodeFloat(true)
	case simpleVdFloat64:
		vt = valueTypeFloat
		v = d.decodeFloat(false)
	case simpleVdString:
		vt = valueTypeString
		v = d.decodeString()
	case simpleVdByteArray:
		vt = valueTypeBytes
		v, _ = d.decodeBytes(nil)
	case simpleVdExt:
		vt = valueTypeExt
		l := d.decLen()
		var re RawExt
		re.Tag = d.r.readn1()
		re.Data = d.r.readn(l)
		v = &re
		vt = valueTypeExt
	case simpleVdArray:
		vt = valueTypeArray
		decodeFurther = true
	case simpleVdMap:
		vt = valueTypeMap
		decodeFurther = true
	default:
		decErr("decodeNaked: Unrecognized d.vd: 0x%x", d.bd)
	}

	if !decodeFurther {
		d.bdRead = false
	}
	return
}

//------------------------------------

//SimpleHandle is a Handle for a very simple encoding format.
//
//The simple format is similar to binc:
//  - All numbers are represented with 8 bytes (int64, uint64, float64)
//  - Strings, []byte, arrays and maps have the length pre-pended as a uint64
//  - Thus, there isn't much packing, but the format is extremely simple.
//    and easy to create different implementations of (e.g. in C).
type SimpleHandle struct {
	BasicHandle
}

func (h *SimpleHandle) newEncDriver(w encWriter) encDriver {
	return &simpleEncDriver{w: w}
}

func (h *SimpleHandle) newDecDriver(r decReader) decDriver {
	return &simpleDecDriver{r: r}
}

func (_ *SimpleHandle) writeExt() bool {
	return true
}

func (h *SimpleHandle) getBasicHandle() *BasicHandle {
	return &h.BasicHandle
}

var _ decDriver = (*simpleDecDriver)(nil)
var _ encDriver = (*simpleEncDriver)(nil)
