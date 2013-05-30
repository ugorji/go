// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"reflect"
	"math"
	"time"
	"fmt"
	"net/rpc"
	"io"
)

const (
	mpPosFixNumMin byte = 0x00
	mpPosFixNumMax = 0x7f
	mpFixMapMin = 0x80
	mpFixMapMax = 0x8f
	mpFixArrayMin =  0x90
	mpFixArrayMax =  0x9f
	mpFixRawMin  = 0xa0
	mpFixRawMax  = 0xbf
	mpNil = 0xc0
	mpFalse = 0xc2
	mpTrue = 0xc3
	mpFloat = 0xca
	mpDouble = 0xcb
	mpUint8 = 0xcc
	mpUint16 = 0xcd
	mpUint32 = 0xce
	mpUint64 = 0xcf
	mpInt8 = 0xd0
	mpInt16 = 0xd1
	mpInt32 = 0xd2
	mpInt64 = 0xd3
	mpRaw16 = 0xda
	mpRaw32 = 0xdb
	mpArray16 = 0xdc
	mpArray32 = 0xdd
	mpMap16 = 0xde
	mpMap32 = 0xdf
	mpNegFixNumMin = 0xe0
	mpNegFixNumMax = 0xff

	// extensions below
	// mpBin8 = 0xc4
	// mpBin16 = 0xc5
	// mpBin32 = 0xc6
	// mpExt8 = 0xc7
	// mpExt16 = 0xc8
	// mpExt32 = 0xc9
	// mpFixExt1 = 0xd4
	// mpFixExt2 = 0xd5
	// mpFixExt4 = 0xd6
	// mpFixExt8 = 0xd7
	// mpFixExt16 = 0xd8

	// extensions based off v4: https://gist.github.com/frsyuki/5235364
	mpXv4Fixext0 = 0xc4
	mpXv4Fixext1 = 0xc5
	mpXv4Fixext2 = 0xc6
	mpXv4Fixext3 = 0xc7
	mpXv4Fixext4 = 0xc8
	mpXv4Fixext5 = 0xc9

	mpXv4Ext8m = 0xd4
	mpXv4Ext16m = 0xd5
	mpXv4Ext32m = 0xd6
	mpXv4Ext8 = 0xd7
	mpXv4Ext16 = 0xd8
	mpXv4Ext32 = 0xd9
)	


// A MsgpackContainer type specifies the different types of msgpackContainers.
type msgpackContainerType struct {
	cutoff int8
	b0, b1, b2 byte
}

var (
	msgpackContainerRawBytes = msgpackContainerType{32, mpFixRawMin, mpRaw16, mpRaw32}
	msgpackContainerList = msgpackContainerType{16, mpFixArrayMin, mpArray16, mpArray32}
	msgpackContainerMap = msgpackContainerType{16, mpFixMapMin, mpMap16, mpMap32}
)

// MsgpackSpecRpc is the implementation of Rpc that uses custom communication protocol 
// as defined in the msgpack spec at http://wiki.msgpack.org/display/MSGPACK/RPC+specification
type MsgpackSpecRpc struct{}

type msgpackSpecRpcCodec struct {
	rpcCodec
}

//MsgpackHandle is a Handle for the Msgpack Schema-Free Encoding Format.
type MsgpackHandle struct {
	// RawToString controls how raw bytes are decoded into a nil interface{}.
	// Note that setting an extension func for []byte ensures that raw bytes 
	// are decoded as strings, regardless of this setting. 
	// This setting is used only if an extension func isn't defined for []byte.
	RawToString bool
	// WriteExt flag supports encoding configured extensions with extension tags.
	// 
	// With WriteExt=false, configured extensions are serialized as raw bytes.
	// 
	// They can still be decoded into a typed object, provided an appropriate one is 
	// provided, but the type cannot be inferred from the stream. If no appropriate
	// type is provided (e.g. decoding into a nil interface{}), you get back
	// a []byte or string based on the setting of RawToString.
	WriteExt bool	

	encdecHandle
	DecodeOptions
}

type msgpackEncoder struct { 
	w encWriter
}

type msgpackDecoder struct {
	r decReader
	bd byte
	bdRead bool
}

func (e *msgpackEncoder) encodeBuiltinType(rt reflect.Type, rv reflect.Value) bool {
	//no builtin types. All encodings are based on kinds. Types supported as extensions.
	return false
}

func (e *msgpackEncoder) encodeNil() {
	e.w.writen1(mpNil)
}

func (e *msgpackEncoder) encodeInt(i int64) {
	switch {
	case i >= -32 && i <= math.MaxInt8:
		e.w.writen1(byte(i))
	case i < -32 && i >= math.MinInt8:
		e.w.writen2(mpInt8, byte(i))
	case i >= math.MinInt16 && i <= math.MaxInt16:
		e.w.writen1(mpInt16)
		e.w.writeUint16(uint16(i))
	case i >= math.MinInt32 && i <= math.MaxInt32:
		e.w.writen1(mpInt32)
		e.w.writeUint32(uint32(i))
	case i >= math.MinInt64 && i <= math.MaxInt64:
		e.w.writen1(mpInt64)
		e.w.writeUint64(uint64(i))
	default:
		encErr("encInt64: Unreachable block")
	}
}

func (e *msgpackEncoder) encodeUint(i uint64) {
	// uints are not fixnums. fixnums are always signed.
	// case i <= math.MaxInt8:
	// 	e.w.writen1(byte(i))
	switch {
	case i <= math.MaxUint8:
		e.w.writen2(mpUint8, byte(i))
	case i <= math.MaxUint16:
		e.w.writen1(mpUint16)
		e.w.writeUint16(uint16(i))
	case i <= math.MaxUint32:
		e.w.writen1(mpUint32)
		e.w.writeUint32(uint32(i))
	default:
		e.w.writen1(mpUint64)
		e.w.writeUint64(uint64(i))
	}
}

func (e *msgpackEncoder) encodeBool(b bool) {
	if b {
		e.w.writen1(mpTrue)
	} else {
		e.w.writen1(mpFalse)
	}
}

func (e *msgpackEncoder) encodeFloat32(f float32) {
	e.w.writen1(mpFloat)
	e.w.writeUint32(math.Float32bits(f))
}

func (e *msgpackEncoder) encodeFloat64(f float64) {
	e.w.writen1(mpDouble)
	e.w.writeUint64(math.Float64bits(f))
}

func (e *msgpackEncoder) encodeExtPreamble(xtag byte, l int) {
	switch {
	case l <= 4:
		e.w.writen2(0xd4 | byte(l), xtag)
	case l <= 8:
		e.w.writen2(0xc0 | byte(l), xtag)
	case l < 256:
		e.w.writen3(mpXv4Fixext5, xtag, byte(l))
	case l < 65536:
		e.w.writen2(mpXv4Ext16, xtag)
		e.w.writeUint16(uint16(l))
	default:
		e.w.writen2(mpXv4Ext32, xtag)
		e.w.writeUint32(uint32(l))
	}
}

func (e *msgpackEncoder) encodeArrayPreamble(length int) {
	e.writeContainerLen(msgpackContainerList, length)
}


func (e *msgpackEncoder) encodeMapPreamble(length int) {
	e.writeContainerLen(msgpackContainerMap, length)	
}

func (e *msgpackEncoder) encodeString(c charEncoding, s string) {
	//ignore charEncoding. 
	e.writeContainerLen(msgpackContainerRawBytes, len(s))
	if len(s) > 0 {
		e.w.writestr(s)
	}
}

func (e *msgpackEncoder) encodeStringBytes(c charEncoding, bs []byte) {
	//ignore charEncoding. 
	e.writeContainerLen(msgpackContainerRawBytes, len(bs))
	if len(bs) > 0 {
		e.w.writeb(bs)
	}
}

func (e *msgpackEncoder) writeContainerLen(ct msgpackContainerType, l int) {
	switch {
	case l < int(ct.cutoff):
		e.w.writen1(ct.b0 | byte(l))
	case l < 65536:
		e.w.writen1(ct.b1)
		e.w.writeUint16(uint16(l))
	default:
		e.w.writen1(ct.b2)
		e.w.writeUint32(uint32(l))
	}
}

//---------------------------------------------

func (d *msgpackDecoder) decodeBuiltinType(rt reflect.Type, rv reflect.Value) bool { 
	return false
}

// Note: This returns either a primitive (int, bool, etc) for non-containers,
// or a containerType, or a specific type denoting nil or extension. 
// It is called when a nil interface{} is passed, leaving it up to the Decoder
// to introspect the stream and decide how best to decode.
// It deciphers the value by looking at the stream first.
func (d *msgpackDecoder) decodeNaked(h decodeHandleI) (rv reflect.Value, ctx decodeNakedContext) {
	d.initReadNext()
	bd := d.bd
	
	var v interface{}

	switch bd {
	case mpNil:
		ctx = dncNil
		d.bdRead = false
	case mpFalse:
		v = false
	case mpTrue:
		v = true

	case mpFloat:
		v = math.Float32frombits(d.r.readUint32())
	case mpDouble:
		v = math.Float64frombits(d.r.readUint64())
		
	case mpUint8:
		v = d.r.readUint8()
	case mpUint16:
		v = d.r.readUint16()
	case mpUint32:
		v = d.r.readUint32()
	case mpUint64:
		v = d.r.readUint64()
		
	case mpInt8:
		v = int8(d.r.readUint8())
	case mpInt16:
		v = int16(d.r.readUint16())
	case mpInt32:
		v = int32(d.r.readUint32())
	case mpInt64:
		v = int64(d.r.readUint64())
		
	default:
		switch {
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:
			// positive fixnum (always signed)
			v = int8(bd)
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
			// negative fixnum
			v = int8(bd)		
		case bd == mpRaw16, bd == mpRaw32, bd >= mpFixRawMin && bd <= mpFixRawMax:
			ctx = dncContainer
			// v = containerRawBytes
			opts := h.(*MsgpackHandle)
			if opts.rawToStringOverride || opts.RawToString {
				var rvm string
				rv = reflect.ValueOf(&rvm).Elem()
			} else {
				rv = reflect.New(byteSliceTyp).Elem() // Use New, not Zero, so it's settable
			}
		case bd == mpArray16, bd == mpArray32, bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			ctx = dncContainer
			// v = containerList
			opts := h.(*MsgpackHandle)
			if opts.SliceType == nil {
				rv = reflect.New(intfSliceTyp).Elem()
			} else {
				rv = reflect.New(opts.SliceType).Elem()
			}
		case bd == mpMap16, bd == mpMap32, bd >= mpFixMapMin && bd <= mpFixMapMax:
			ctx = dncContainer
			// v = containerMap
			opts := h.(*MsgpackHandle)
			if opts.MapType == nil {
				rv = reflect.MakeMap(mapIntfIntfTyp)
			} else {
				rv = reflect.MakeMap(opts.MapType)
			}
		case bd >= mpXv4Fixext0 && bd <= mpXv4Fixext5, bd >= mpXv4Ext8m && bd <= mpXv4Ext32:
			//ctx = dncExt
			xtag := d.r.readUint8()
			opts := h.(*MsgpackHandle)
			rt, bfn := opts.getDecodeExtForTag(xtag)
			if rt == nil {
				decErr("Unable to find type mapped to extension tag: %v", xtag)
			}
			if rt.Kind() == reflect.Ptr {
				rv = reflect.New(rt.Elem())
			} else {
				rv = reflect.New(rt).Elem()
			}
			if fnerr := bfn(rv, d.r.readn(d.readExtLen())); fnerr != nil {
				panic(fnerr)
			} 
		default:
			decErr("Nil-Deciphered DecodeValue: %s: hex: %x, dec: %d", msgBadDesc, bd, bd)
		}
	}
	if ctx == dncHandled {
		d.bdRead = false
		if v != nil {
			rv = reflect.ValueOf(v)
		}
	}
	return
}

// int can be decoded from msgpack type: intXXX or uintXXX 
func (d *msgpackDecoder) decodeInt(bitsize uint8) (i int64) {
	switch d.bd {
	case mpUint8:
		i = int64(uint64(d.r.readUint8()))
	case mpUint16:
		i = int64(uint64(d.r.readUint16()))
	case mpUint32:
		i = int64(uint64(d.r.readUint32()))
	case mpUint64:
		i = int64(d.r.readUint64())
	case mpInt8:
		i = int64(int8(d.r.readUint8()))
	case mpInt16:
		i = int64(int16(d.r.readUint16()))
	case mpInt32:
		i = int64(int32(d.r.readUint32()))
	case mpInt64:
		i = int64(d.r.readUint64())
	default:
		switch {
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			i = int64(int8(d.bd))
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			i = int64(int8(d.bd))
		default:
			decErr("Unhandled single-byte unsigned integer value: %s: %x", msgBadDesc, d.bd)
		}
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


// uint can be decoded from msgpack type: intXXX or uintXXX 
func (d *msgpackDecoder) decodeUint(bitsize uint8) (ui uint64) {
	switch d.bd {
	case mpUint8:
		ui = uint64(d.r.readUint8())
	case mpUint16:
		ui = uint64(d.r.readUint16())
	case mpUint32:
		ui = uint64(d.r.readUint32())
	case mpUint64:
		ui = d.r.readUint64()
	case mpInt8:
		if i := int64(int8(d.r.readUint8())); i >= 0 {
			ui = uint64(i)
		} else {
			decErr("Assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt16:
		if i := int64(int16(d.r.readUint16())); i >= 0 {
			ui = uint64(i)
		} else {
			decErr("Assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt32:
		if i := int64(int32(d.r.readUint32())); i >= 0 {
			ui = uint64(i)
		} else {
			decErr("Assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt64:
		if i := int64(d.r.readUint64()); i >= 0 {
			ui = uint64(i)
		} else {
			decErr("Assigning negative signed value: %v, to unsigned type", i)
		}
	default:
		switch {
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			ui = uint64(d.bd)
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			decErr("Assigning negative signed value: %v, to unsigned type", int(d.bd))
		default:
			decErr("Unhandled single-byte unsigned integer value: %s: %x", msgBadDesc, d.bd)
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

// float can either be decoded from msgpack type: float, double or intX
func (d *msgpackDecoder) decodeFloat(chkOverflow32 bool) (f float64) {
	switch d.bd {
	case mpFloat:
		f = float64(math.Float32frombits(d.r.readUint32()))
	case mpDouble:
		f = math.Float64frombits(d.r.readUint64())
	default:
		f = float64(d.decodeInt(0))
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

// bool can be decoded from bool, fixnum 0 or 1.
func (d *msgpackDecoder) decodeBool() (b bool) {
	switch d.bd {
	case mpFalse, 0:
		// b = false
	case mpTrue, 1:
		b = true
	default:
		decErr("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}
	
func (d *msgpackDecoder) decodeString() (s string) {
	clen := d.readContainerLen(msgpackContainerRawBytes)
	if clen > 0 {
		s = string(d.r.readn(clen))
	}
	d.bdRead = false
	return
}

// Callers must check if changed=true (to decide whether to replace the one they have)
func (d *msgpackDecoder) decodeStringBytes(bs []byte) (bsOut []byte, changed bool) {
	clen := d.readContainerLen(msgpackContainerRawBytes)
	// if clen < 0 {
	// 	changed = true
	// 	panic("length cannot be zero. this cannot be nil.")
	// } 
	if clen > 0 {
		// if no contents in stream, don't update the passed byteslice	
		if len(bs) != clen {
			// Return changed=true if length of passed slice is different from length of bytes in the stream.
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

// Every top-level decode funcs (i.e. decodeValue, decode) must call this first.
func (d *msgpackDecoder) initReadNext() {
	if d.bdRead {
		return
	}
	d.bd = d.r.readUint8()
	d.bdRead = true
}

func (d *msgpackDecoder) currentIsNil() bool {
	if d.bd == mpNil {
		d.bdRead = false
		return true
	} 
	return false
}

func (d *msgpackDecoder) readContainerLen(ct msgpackContainerType) (clen int) {
	switch {
	case d.bd == mpNil:
		clen = -1 // to represent nil
	case d.bd == ct.b1:
		clen = int(d.r.readUint16())
	case d.bd == ct.b2:
		clen = int(d.r.readUint32())
	case (ct.b0 & d.bd) == ct.b0:
		clen = int(ct.b0 ^ d.bd)
	default:
		decErr("readContainerLen: %s: hex: %x, dec: %d", msgBadDesc, d.bd, d.bd)
	}
	d.bdRead = false
	return	
}

func (d *msgpackDecoder) readMapLen() int {
	return d.readContainerLen(msgpackContainerMap)
}

func (d *msgpackDecoder) readArrayLen() int {
	return d.readContainerLen(msgpackContainerList)
}


func (d *msgpackDecoder) readExtLen() (clen int) {
	switch d.bd {
	case mpNil:
		clen = -1 // to represent nil
	case mpXv4Fixext5:
		clen = int(d.r.readUint8())
	case mpXv4Ext16:
		clen = int(d.r.readUint16())
	case mpXv4Ext32:
		clen = int(d.r.readUint32())
	default: 
		switch {
		case d.bd >= mpXv4Fixext0 && d.bd <= mpXv4Fixext4:
			clen = int(d.bd & 0x0f)
		case d.bd >= mpXv4Ext8m && d.bd <= mpXv4Ext8:
			clen = int(d.bd & 0x03)
		default:
			decErr("decoding ext bytes: found unexpected byte: %x", d.bd)
		}
	}
	return
}

func (d *msgpackDecoder) decodeExt(tag byte) (xbs []byte) {
	// if (d.bd >= mpXv4Fixext0 && d.bd <= mpXv4Fixext5) || (d.bd >= mpXv4Ext8m && d.bd <= mpXv4Ext32) {
	xbd := d.bd
	switch {
	case xbd >= mpXv4Fixext0 && xbd <= mpXv4Fixext5, xbd >= mpXv4Ext8m && xbd <= mpXv4Ext32:
		if xtag := d.r.readUint8(); xtag != tag {
			decErr("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		xbs = d.r.readn(d.readExtLen())
	case xbd == mpRaw16, xbd == mpRaw32, xbd >= mpFixRawMin && xbd <= mpFixRawMax:
		xbs, _ = d.decodeStringBytes(nil)
	default:
		decErr("Wrong byte descriptor (Expecting extensions or raw bytes). Got: 0x%x", xbd)
	}		
	d.bdRead = false
	return
}

//--------------------------------------------------

func (MsgpackSpecRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) (rpc.ServerCodec) {
	return &msgpackSpecRpcCodec{ newRPCCodec(conn, h) }
}

func (MsgpackSpecRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) (rpc.ClientCodec) {
	return &msgpackSpecRpcCodec{ newRPCCodec(conn, h) }
}

// /////////////// Spec RPC Codec ///////////////////
func (c msgpackSpecRpcCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	return c.writeCustomBody(0, r.Seq, r.ServiceMethod, body)
}

func (c msgpackSpecRpcCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	return c.writeCustomBody(1, r.Seq, r.Error, body)
}

func (c msgpackSpecRpcCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.parseCustomHeader(1, &r.Seq, &r.Error)
}

func (c msgpackSpecRpcCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.parseCustomHeader(0, &r.Seq, &r.ServiceMethod)
}

func (c msgpackSpecRpcCodec) parseCustomHeader(expectTypeByte byte, msgid *uint64, methodOrError *string) (err error) {

	// We read the response header by hand 
	// so that the body can be decoded on its own from the stream at a later time.

	bs := make([]byte, 1)
	n, err := c.rwc.Read(bs)
	if err != nil {
		return 
	}
	if n != 1 {
		err = fmt.Errorf("Couldn't read array descriptor: No bytes read")
		return
	}
	const fia byte = 0x94 //four item array descriptor value
	if bs[0] != fia {
		err = fmt.Errorf("Unexpected value for array descriptor: Expecting %v. Received %v", fia, bs[0])
		return
	}
	var b byte
	if err = c.read(&b, msgid, methodOrError); err != nil {
		return
	}
	if b != expectTypeByte {
		err = fmt.Errorf("Unexpected byte descriptor in header. Expecting %v. Received %v", expectTypeByte, b)
		return
	}
	return
}

func (c msgpackSpecRpcCodec) writeCustomBody(typeByte byte, msgid uint64, methodOrError string, body interface{}) (err error) {
	var moe interface{} = methodOrError
	// response needs nil error (not ""), and only one of error or body can be nil
	if typeByte == 1 {
		if methodOrError == "" {
			moe = nil
		}
		if moe != nil && body != nil {
			body = nil
		}
	}
	r2 := []interface{}{ typeByte, uint32(msgid), moe, body }
	return c.enc.Encode(r2)
}



//--------------------------------------------------

// EncodeBinaryExt returns the underlying bytes of this value AS-IS.
// Configure this to support the Binary Extension using tag 0.
func (_ *MsgpackHandle) BinaryEncodeExt(rv reflect.Value) ([]byte, error) {
	if rv.IsNil() {
		return nil, nil
	}
	return rv.Bytes(), nil
}

// DecodeBinaryExt sets passed byte array AS-IS into the reflect Value.
// Configure this to support the Binary Extension using tag 0.
func (_ *MsgpackHandle) BinaryDecodeExt(rv reflect.Value, bs []byte) (err error) {
	rv.SetBytes(bs)
	return
}

// EncodeBinaryExt returns the underlying bytes of this value AS-IS.
// Configure this to support the Binary Extension using tag 0.
func (_ *MsgpackHandle) TimeEncodeExt(rv reflect.Value) (bs []byte, err error) {
	bs = encodeTime(rv.Interface().(time.Time))
	return
}

func (_ *MsgpackHandle) TimeDecodeExt(rv reflect.Value, bs []byte) (err error) {
	tt, err := decodeTime(bs)
	if err == nil {
		rv.Set(reflect.ValueOf(tt))
	}
	return
}

func (_ *MsgpackHandle) newEncoder(w encWriter) encoder {
	return &msgpackEncoder{w: w}
}

func (_ *MsgpackHandle) newDecoder(r decReader) decoder {
	return &msgpackDecoder{r: r}
}

func (o *MsgpackHandle) writeExt() bool {
	return o.WriteExt
}

