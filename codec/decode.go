// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"io"
	"reflect"
)

// Some tagging information for error messages.
var (
	msgTagDec  = "codec.decoder"
	msgBadDesc = "Unrecognized descriptor byte"
)

// when decoding without schema, the nakedContext tells us what 
// we decoded into, or if decoding has been handled.
type decodeNakedContext uint8

const (
	dncHandled decodeNakedContext = iota
	dncNil
	// dncExt
	dncContainer
)

// decodeEncodedType is the current type in the encoded stream
type decodeEncodedType uint8

const (
	detUnset decodeEncodedType = iota
	detNil
	detInt
	detUint
	detFloat
	detBool
	detString
	detBytes
	detMap
	detArray
	detTimestamp
	detExt
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface {
	readn(n int) []byte
	readb([]byte)
	readn1() uint8
	readUint16() uint16
	readUint32() uint32
	readUint64() uint64
}

type decDriver interface {
	initReadNext()
	tryDecodeAsNil() bool
	currentEncodedType() decodeEncodedType
	isBuiltinType(rt uintptr) bool
	decodeBuiltinType(rt uintptr, rv reflect.Value)
	//decodeNaked should completely handle extensions, builtins, primitives, etc.
	//Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	decodeNaked() (rv reflect.Value, ctx decodeNakedContext)
	decodeInt(bitsize uint8) (i int64)
	decodeUint(bitsize uint8) (ui uint64)
	decodeFloat(chkOverflow32 bool) (f float64)
	decodeBool() (b bool)
	// decodeString can also decode symbols
	decodeString() (s string)
	decodeBytes(bs []byte) (bsOut []byte, changed bool)
	decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)
	readMapLen() int
	readArrayLen() int
}

// decFnInfo has methods for registering handling decoding of a specific type
// based on some characteristics (builtin, extension, reflect Kind, etc)
type decFnInfo struct {
	ti *typeInfo
	d   *Decoder
	dd  decDriver
	xfFn  func(reflect.Value, []byte) error
	xfTag byte 
}

type decFn struct {
	i *decFnInfo
	f func(*decFnInfo, reflect.Value) 
}

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct {
	r decReader
	d decDriver
	h decodeHandleI
	f map[uintptr]decFn
	x []uintptr
	s []decFn
}

func (f *decFnInfo) builtin(rv reflect.Value) {
	f.dd.decodeBuiltinType(f.ti.rtid, rv)
}

func (f *decFnInfo) rawExt(rv reflect.Value) {
	xtag, xbs := f.dd.decodeExt(false, 0)
	rv.Field(0).SetUint(uint64(xtag))
	rv.Field(1).SetBytes(xbs)
}

func (f *decFnInfo) ext(rv reflect.Value) {
	_, xbs := f.dd.decodeExt(true, f.xfTag)
	if fnerr := f.xfFn(rv, xbs); fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) binaryMarshal(rv reflect.Value) {
	var bm binaryUnmarshaler
	if f.ti.unmIndir == -1 {
		bm = rv.Addr().Interface().(binaryUnmarshaler)
	} else if f.ti.unmIndir == 0 {
		bm = rv.Interface().(binaryUnmarshaler)
	} else {
		for j, k := int8(0), f.ti.unmIndir; j < k; j++ {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		bm = rv.Interface().(binaryUnmarshaler)
	}
	xbs, _ := f.dd.decodeBytes(nil)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) kErr(rv reflect.Value) {
	decErr("Unhandled value for kind: %v: %s", rv.Kind(), msgBadDesc)
}

func (f *decFnInfo) kString(rv reflect.Value) {
	rv.SetString(f.dd.decodeString())
}

func (f *decFnInfo) kBool(rv reflect.Value) {
	rv.SetBool(f.dd.decodeBool())
}

func (f *decFnInfo) kInt(rv reflect.Value) {
	rv.SetInt(f.dd.decodeInt(intBitsize))
}

func (f *decFnInfo) kInt64(rv reflect.Value) {
	rv.SetInt(f.dd.decodeInt(64))
}

func (f *decFnInfo) kInt32(rv reflect.Value) {
	rv.SetInt(f.dd.decodeInt(32))
}

func (f *decFnInfo) kInt8(rv reflect.Value) {
	rv.SetInt(f.dd.decodeInt(8))
}

func (f *decFnInfo) kInt16(rv reflect.Value) {
	rv.SetInt(f.dd.decodeInt(16))
}

func (f *decFnInfo) kFloat32(rv reflect.Value) {
	rv.SetFloat(f.dd.decodeFloat(true))
}

func (f *decFnInfo) kFloat64(rv reflect.Value) {
	rv.SetFloat(f.dd.decodeFloat(false))
}

func (f *decFnInfo) kUint8(rv reflect.Value) {
	rv.SetUint(f.dd.decodeUint(8))
}

func (f *decFnInfo) kUint64(rv reflect.Value) {
	rv.SetUint(f.dd.decodeUint(64))
}

func (f *decFnInfo) kUint(rv reflect.Value) {
	rv.SetUint(f.dd.decodeUint(uintBitsize))
}

func (f *decFnInfo) kUint32(rv reflect.Value) {
	rv.SetUint(f.dd.decodeUint(32))
}

func (f *decFnInfo) kUint16(rv reflect.Value) {
	rv.SetUint(f.dd.decodeUint(16))
}

// func (f *decFnInfo) kPtr(rv reflect.Value) {
// 	debugf(">>>>>>> ??? decode kPtr called - shouldn't get called")
// 	if rv.IsNil() {
// 		rv.Set(reflect.New(rv.Type().Elem()))
// 	}
// 	f.d.decodeValue(rv.Elem())
// }

func (f *decFnInfo) kInterface(rv reflect.Value) {
	if rv.IsNil() {
		// if nil interface, use some hieristics to set the nil interface to an
		// appropriate value based on the first byte read (byte descriptor bd)
		rv2, ndesc := f.dd.decodeNaked()
		if ndesc == dncNil {
			return
		}
		// Cannot decode into nil interface with methods (e.g. error, io.Reader, etc) 
		// if non-nil value in stream.
		if num := f.ti.rt.NumMethod(); num > 0 {
			decErr("decodeValue: Cannot decode non-nil codec value into nil %v (%v methods)", 
				f.ti.rt , num)
		} 
		if ndesc == dncHandled {
			rv.Set(rv2)
			return
		}
		f.d.decodeValue(rv2)
		rv.Set(rv2)
	} else {
		f.d.decodeValue(rv.Elem())
	}
}

func (f *decFnInfo) kStruct(rv reflect.Value) {
	fti := f.ti
	if currEncodedType := f.dd.currentEncodedType(); currEncodedType == detMap {
		containerLen := f.dd.readMapLen()
		if containerLen == 0 {
			return
		}
		tisfi := fti.sfi 
		for j := 0; j < containerLen; j++ {
			// var rvkencname string
			// ddecode(&rvkencname)
			f.dd.initReadNext()
			rvkencname := f.dd.decodeString()
			// rvksi := ti.getForEncName(rvkencname)
			if k := fti.indexForEncName(rvkencname); k > -1 {
				sfik := tisfi[k]
				if sfik.i != -1 {
					f.d.decodeValue(rv.Field(int(sfik.i)))
				} else {
					f.d.decodeValue(rv.FieldByIndex(sfik.is))
				}
				// f.d.decodeValue(ti.field(k, rv))
			} else {
				if f.d.h.errorIfNoField() {
					decErr("No matching struct field found when decoding stream map with key: %v", 
						rvkencname)
				} else {
					var nilintf0 interface{}
					f.d.decodeValue(reflect.ValueOf(&nilintf0).Elem())
				}
			}
		}
	} else if currEncodedType == detArray {
		containerLen := f.dd.readArrayLen()
		if containerLen == 0 {
			return
		}
		for j, si := range fti.sfip {
			if j == containerLen {
				break
			}
			if si.i != -1 {
				f.d.decodeValue(rv.Field(int(si.i)))
			} else {
				f.d.decodeValue(rv.FieldByIndex(si.is))
			}
		}
		if containerLen > len(fti.sfip) {
			// read remaining values and throw away
			for j := len(fti.sfip); j < containerLen; j++ {
				var nilintf0 interface{}
				f.d.decodeValue(reflect.ValueOf(&nilintf0).Elem())
			}
		}
	} else {
		decErr("Only encoded map or array can be decoded into a struct. (decodeEncodedType: %x)", 
			currEncodedType)
	}
}

func (f *decFnInfo) kSlice(rv reflect.Value) {
	// Be more careful calling Set() here, because a reflect.Value from an array
	// may have come in here (which may not be settable).
	// In places where the slice got from an array could be, we should guard with CanSet() calls.

	if f.ti.rtid == byteSliceTypId { // rawbytes
		if bs2, changed2 := f.dd.decodeBytes(rv.Bytes()); changed2 {
			rv.SetBytes(bs2)
		}
		return
	}

	containerLen := f.dd.readArrayLen()

	if rv.IsNil() {
		rv.Set(reflect.MakeSlice(f.ti.rt, containerLen, containerLen))
	} 
	if containerLen == 0 {
		return
	}

	// if we need to reset rv but it cannot be set, we should err out.
	// for example, if slice is got from unaddressable array, CanSet = false
	if rvcap, rvlen := rv.Len(), rv.Cap(); containerLen > rvcap {
		if rv.CanSet() {
			rvn := reflect.MakeSlice(f.ti.rt, containerLen, containerLen)
			if rvlen > 0 {
				reflect.Copy(rvn, rv)
			}
			rv.Set(rvn)
		} else {
			decErr("Cannot reset slice with less cap: %v than stream contents: %v", 
				rvcap, containerLen)
		}
	} else if containerLen > rvlen {
		rv.SetLen(containerLen)
	}
	for j := 0; j < containerLen; j++ {
		f.d.decodeValue(rv.Index(j))
	}
}

func (f *decFnInfo) kArray(rv reflect.Value) {
	f.d.decodeValue(rv.Slice(0, rv.Len()))
}

func (f *decFnInfo) kMap(rv reflect.Value) {
	containerLen := f.dd.readMapLen()

	if rv.IsNil() {
		rv.Set(reflect.MakeMap(f.ti.rt))
	}
	
	if containerLen == 0 {
		return
	}

	ktype, vtype := f.ti.rt.Key(), f.ti.rt.Elem()
	for j := 0; j < containerLen; j++ {
		rvk := reflect.New(ktype).Elem()
		f.d.decodeValue(rvk)

		if ktype == intfTyp {
			rvk = rvk.Elem()
			if rvk.Type() == byteSliceTyp {
				rvk = reflect.ValueOf(string(rvk.Bytes()))
			}
		}
		rvv := rv.MapIndex(rvk)
		if !rvv.IsValid() {
			rvv = reflect.New(vtype).Elem()
		}

		f.d.decodeValue(rvv)
		rv.SetMapIndex(rvk, rvv)
	}
}

// ioDecReader is a decReader that reads off an io.Reader
type ioDecReader struct {
	r io.Reader
	br io.ByteReader
	x [8]byte //temp byte array re-used internally for efficiency
}

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct {
	b []byte // data
	c int    // cursor
	a int    // available
}

type decodeHandleI interface {
	getDecodeExt(rt uintptr) (tag byte, fn func(reflect.Value, []byte) error)
	errorIfNoField() bool
}

type DecodeOptions struct {
	// An instance of MapType is used during schema-less decoding of a map in the stream.
	// If nil, we use map[interface{}]interface{}
	MapType reflect.Type
	// An instance of SliceType is used during schema-less decoding of an array in the stream.
	// If nil, we use []interface{}
	SliceType reflect.Type
	// ErrorIfNoField controls whether an error is returned when decoding a map
	// from a codec stream into a struct, and no matching struct field is found.
	ErrorIfNoField bool
}

func (o *DecodeOptions) errorIfNoField() bool {
	return o.ErrorIfNoField
}

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
// 
// For efficiency, Users are encouraged to pass in a memory buffered writer
// (eg bufio.Reader, bytes.Buffer). 
func NewDecoder(r io.Reader, h Handle) *Decoder {
	z := ioDecReader{
		r: r,
	}
	z.br, _ = r.(io.ByteReader)
	return &Decoder{r: &z, d: h.newDecDriver(&z), h: h}
}

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder {
	z := bytesDecReader{
		b: in,
		a: len(in),
	}
	return &Decoder{r: &z, d: h.newDecDriver(&z), h: h}
}

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
// 
// Sample usages:
//   // Decoding into a non-nil typed value
//   var f float32
//   err = codec.NewDecoder(r, handle).Decode(&f)
//
//   // Decoding into nil interface
//   var v interface{}
//   dec := codec.NewDecoder(r, handle)
//   err = dec.Decode(&v)
// 
// When decoding into a nil interface{}, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64. 
//   - Other values are decoded appropriately depending on the encoding: 
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
// Configurations exist on the Handle to override defaults 
// (e.g. for MapType, SliceType and how to decode raw bytes).
// 
// When decoding into a non-nil interface{} value, the mode of encoding is based on the 
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
// 
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container. 
//   - This means that for a struct or map, we just update matching fields or keys.
//   - For a slice/array, we just update the first n elements, where n is length of the stream.
//   - However, if decoding into a nil map/slice and the length of the stream is 0,
//     we reset the destination map/slice to be a zero-length non-nil map/slice.
//   - Also, if the encoded value is Nil in the stream, then we try to set
//     the container to its "zero" value (e.g. nil for slice/map).
//   - Note that a struct can be decoded from an array in the stream,
//     by updating fields as they occur in the struct.
func (d *Decoder) Decode(v interface{}) (err error) {
	defer panicToErr(&err)
	d.decode(v)
	return
}

func (d *Decoder) decode(iv interface{}) {
	d.d.initReadNext()

	// Fast path included for various pointer types which cannot be registered as extensions
	switch v := iv.(type) {
	case nil:
		decErr("Cannot decode into nil.")
	case reflect.Value:
		d.chkPtrValue(v)
		d.decodeValue(v)
	case *string:
		*v = d.d.decodeString()
	case *bool:
		*v = d.d.decodeBool()
	case *int:
		*v = int(d.d.decodeInt(intBitsize))
	case *int8:
		*v = int8(d.d.decodeInt(8))
	case *int16:
		*v = int16(d.d.decodeInt(16))
	case *int32:
		*v = int32(d.d.decodeInt(32))
	case *int64:
		*v = int64(d.d.decodeInt(64))
	case *uint:
		*v = uint(d.d.decodeUint(uintBitsize))
	case *uint8:
		*v = uint8(d.d.decodeUint(8))
	case *uint16:
		*v = uint16(d.d.decodeUint(16))
	case *uint32:
		*v = uint32(d.d.decodeUint(32))
	case *uint64:
		*v = uint64(d.d.decodeUint(64))
	case *float32:
		*v = float32(d.d.decodeFloat(true))
	case *float64:
		*v = d.d.decodeFloat(false)
	case *interface{}:
		d.decodeValue(reflect.ValueOf(iv).Elem())
	default:
		rv := reflect.ValueOf(iv)
		d.chkPtrValue(rv)
		d.decodeValue(rv)
	}
}

func (d *Decoder) decodeValue(rv reflect.Value) {
	d.d.initReadNext()

	if d.d.tryDecodeAsNil() {
		// If value in stream is nil, set the dereferenced value to its "zero" value (if settable).
		for rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.CanSet() {
			rv.Set(reflect.Zero(rv.Type()))
		}
		return
	}
	
	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	
	rt := rv.Type()
	rtid := reflect.ValueOf(rt).Pointer()
	
	// retrieve or register a focus'ed function for this type
	// to eliminate need to do the retrieval multiple times
	
	// if d.f == nil && d.s == nil {
	// 	// debugf("---->Creating new dec f map for type: %v\n", rt)
	// }
	var fn decFn 
	var ok bool
	if useMapForCodecCache {
		fn, ok = d.f[rtid]
	} else {
		for i, v := range d.x {
			if v == rtid {
				fn, ok = d.s[i], true
				break
			}
		}
	}
	if !ok {
		// debugf("\tCreating new dec fn for type: %v\n", rt)
		fi := decFnInfo { ti:getTypeInfo(rtid, rt), d:d, dd:d.d }
		fn.i = &fi 
		// An extension can be registered for any type, regardless of the Kind
		// (e.g. type BitSet int64, type MyStruct { / * unexported fields * / }, type X []int, etc.
		//
		// We can't check if it's an extension byte here first, because the user may have
		// registered a pointer or non-pointer type, meaning we may have to recurse first
		// before matching a mapped type, even though the extension byte is already detected.
		//
		// NOTE: if decoding into a nil interface{}, we return a non-nil
		// value except even if the container registers a length of 0.
		if rtid == rawExtTypId {
			fn.f = (*decFnInfo).rawExt 
		} else if d.d.isBuiltinType(rtid) {
			fn.f = (*decFnInfo).builtin 
		} else if xfTag, xfFn := d.h.getDecodeExt(rtid); xfFn != nil {
			fi.xfTag, fi.xfFn = xfTag, xfFn
			fn.f = (*decFnInfo).ext 
		} else if supportBinaryMarshal && fi.ti.unm {
			fn.f = (*decFnInfo).binaryMarshal 
		} else {
			switch rk := rt.Kind(); rk {
			case reflect.String:
				fn.f = (*decFnInfo).kString 
			case reflect.Bool:
				fn.f = (*decFnInfo).kBool 
			case reflect.Int:
				fn.f = (*decFnInfo).kInt 
			case reflect.Int64:
				fn.f = (*decFnInfo).kInt64 
			case reflect.Int32:
				fn.f = (*decFnInfo).kInt32 
			case reflect.Int8:
				fn.f = (*decFnInfo).kInt8 
			case reflect.Int16:
				fn.f = (*decFnInfo).kInt16 
			case reflect.Float32:
				fn.f = (*decFnInfo).kFloat32 
			case reflect.Float64:
				fn.f = (*decFnInfo).kFloat64 
			case reflect.Uint8:
				fn.f = (*decFnInfo).kUint8 
			case reflect.Uint64:
				fn.f = (*decFnInfo).kUint64 
			case reflect.Uint:
				fn.f = (*decFnInfo).kUint 
			case reflect.Uint32:
				fn.f = (*decFnInfo).kUint32 
			case reflect.Uint16:
				fn.f = (*decFnInfo).kUint16 
			// case reflect.Ptr:
			// 	fn.f = (*decFnInfo).kPtr 
			case reflect.Interface:
				fn.f = (*decFnInfo).kInterface 
			case reflect.Struct:
				fn.f = (*decFnInfo).kStruct 
			case reflect.Slice:
				fn.f = (*decFnInfo).kSlice 
			case reflect.Array:
				fn.f = (*decFnInfo).kArray 
			case reflect.Map:
				fn.f = (*decFnInfo).kMap 
			default:
				fn.f = (*decFnInfo).kErr 
			}
		}
		if useMapForCodecCache {
			if d.f == nil {
				d.f = make(map[uintptr]decFn, 16)
			}
			d.f[rtid] = fn
		} else {
			d.s = append(d.s, fn)
			d.x = append(d.x, rtid)
		}
	}
	
	fn.f(fn.i, rv)
	
	return
}

func (d *Decoder) chkPtrValue(rv reflect.Value) {
	// We cannot marshal into a non-pointer or a nil pointer
	// (at least pass a nil interface so we can marshal into it)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		var rvi interface{} = rv
		if rv.IsValid() && rv.CanInterface() {
			rvi = rv.Interface()
		}
		decErr("Decode: Expecting valid pointer to decode into. Got: %v, %T, %v",
			rv.Kind(), rvi, rvi)
	}
}

// ------------------------------------

func (z *ioDecReader) readn(n int) (bs []byte) {
	bs = make([]byte, n)
	if _, err := io.ReadAtLeast(z.r, bs, n); err != nil {
		panic(err)
	}
	return
}

func (z *ioDecReader) readb(bs []byte) {	
	if _, err := io.ReadAtLeast(z.r, bs, len(bs)); err != nil {
		panic(err)
	}
}

func (z *ioDecReader) readn1() uint8 {
	if z.br != nil {
		b, err := z.br.ReadByte()
		if err != nil {
			panic(err)
		}
		return b
	}
	z.readb(z.x[:1])
	return z.x[0]
}

func (z *ioDecReader) readUint16() uint16 {
	z.readb(z.x[:2])
	return bigen.Uint16(z.x[:2])
}

func (z *ioDecReader) readUint32() uint32 {
	z.readb(z.x[:4])
	return bigen.Uint32(z.x[:4])
}

func (z *ioDecReader) readUint64() uint64 {
	z.readb(z.x[:8])
	return bigen.Uint64(z.x[:8])
}

// ------------------------------------

func (z *bytesDecReader) consume(n int) (oldcursor int) {
	if z.a == 0 {
		panic(io.EOF)
	}
	if n > z.a {
		decErr("Trying to read %v bytes. Only %v available", n, z.a)
	}
	// z.checkAvailable(n)
	oldcursor = z.c
	z.c = oldcursor + n
	z.a = z.a - n
	return
}

func (z *bytesDecReader) readn(n int) (bs []byte) {
	c0 := z.consume(n)
	bs = z.b[c0:z.c]
	return
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readn(len(bs)))
}

func (z *bytesDecReader) readn1() uint8 {
	c0 := z.consume(1)
	return z.b[c0]
}

// Use binaryEncoding helper for 4 and 8 bits, but inline it for 2 bits
// creating temp slice variable and copying it to helper function is expensive
// for just 2 bits.

func (z *bytesDecReader) readUint16() uint16 {
	c0 := z.consume(2)
	return uint16(z.b[c0+1]) | uint16(z.b[c0])<<8
}

func (z *bytesDecReader) readUint32() uint32 {
	c0 := z.consume(4)
	return bigen.Uint32(z.b[c0:z.c])
}

func (z *bytesDecReader) readUint64() uint64 {
	c0 := z.consume(8)
	return bigen.Uint64(z.b[c0:z.c])
}

// ----------------------------------------

func decErr(format string, params ...interface{}) {
	doPanic(msgTagDec, format, params...)
}

