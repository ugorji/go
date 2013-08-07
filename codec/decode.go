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

type decodeNakedContext uint8

const (
	dncHandled decodeNakedContext = iota
	dncNil
	dncExt
	dncContainer
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
	currentIsNil() bool
	decodeBuiltinType(rt reflect.Type, rv reflect.Value) bool
	//decodeNaked should completely handle extensions, builtins, primitives, etc.
	//Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	decodeNaked(h decodeHandleI) (rv reflect.Value, ctx decodeNakedContext)
	decodeInt(bitsize uint8) (i int64)
	decodeUint(bitsize uint8) (ui uint64)
	decodeFloat(chkOverflow32 bool) (f float64)
	decodeBool() (b bool)
	// decodeString can also decode symbols
	decodeString() (s string)
	decodeBytes(bs []byte) (bsOut []byte, changed bool)
	decodeExt(tag byte) []byte
	readMapLen() int
	readArrayLen() int
}

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct {
	r decReader
	d decDriver
	h decodeHandleI
}

// ioDecReader is a decReader that reads off an io.Reader
type ioDecReader struct {
	r io.Reader
	x [8]byte //temp byte array re-used internally for efficiency
}

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct {
	b []byte // data
	c int    // cursor
	a int    // available
}

type decExtTagFn struct {
	fn  func(reflect.Value, []byte) error
	tag byte
}

type decExtTypeTagFn struct {
	rt reflect.Type
	decExtTagFn
}

type decodeHandleI interface {
	getDecodeExt(rt reflect.Type) (tag byte, fn func(reflect.Value, []byte) error)
	errorIfNoField() bool
}

type decHandle struct {
	// put word-aligned fields first (before bools, etc)
	exts     []decExtTypeTagFn
	extFuncs map[reflect.Type]decExtTagFn
	// if an extension for byte slice is defined, then always decode Raw as strings
	rawToStringOverride bool
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

// addDecodeExt registers a function to handle decoding into a given type when an
// extension type and specific tag byte is detected in the codec stream.
// To remove an extension, pass fn=nil.
func (o *decHandle) addDecodeExt(rt reflect.Type, tag byte, fn func(reflect.Value, []byte) error) {
	if o.exts == nil {
		o.exts = make([]decExtTypeTagFn, 0, 2)
		o.extFuncs = make(map[reflect.Type]decExtTagFn, 2)
	}
	if _, ok := o.extFuncs[rt]; ok {
		delete(o.extFuncs, rt)
		if rt == byteSliceTyp {
			o.rawToStringOverride = false
		}
	}
	if fn != nil {
		o.extFuncs[rt] = decExtTagFn{fn, tag}
		if rt == byteSliceTyp {
			o.rawToStringOverride = true
		}
	}

	if leno := len(o.extFuncs); leno > cap(o.exts) {
		o.exts = make([]decExtTypeTagFn, leno, (leno * 3 / 2))
	} else {
		o.exts = o.exts[0:leno]
	}
	var i int
	for k, v := range o.extFuncs {
		o.exts[i] = decExtTypeTagFn{k, v}
		i++
	}
}

func (o *decHandle) getDecodeExtForTag(tag byte) (rt reflect.Type, fn func(reflect.Value, []byte) error) {
	for i, l := 0, len(o.exts); i < l; i++ {
		if o.exts[i].tag == tag {
			return o.exts[i].rt, o.exts[i].fn
		}
	}
	return
}

func (o *decHandle) getDecodeExt(rt reflect.Type) (tag byte, fn func(reflect.Value, []byte) error) {
	if l := len(o.exts); l == 0 {
		return
	} else if l < mapAccessThreshold {
		for i := 0; i < l; i++ {
			if o.exts[i].rt == rt {
				x := o.exts[i].decExtTagFn
				return x.tag, x.fn
			}
		}
	} else {
		x := o.extFuncs[rt]
		return x.tag, x.fn
	}
	return
}

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
// 
// For efficiency, Users are encouraged to pass in a memory buffered writer
// (eg bufio.Reader, bytes.Buffer). 
func NewDecoder(r io.Reader, h Handle) *Decoder {
	z := ioDecReader{
		r: r,
	}
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
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container. 
//   - This means that for a struct or map, we just update matching fields or keys.
//   - For a slice/array, we just update the first n elements, where n is length of the stream.
//   - However, if decoding into a nil map/slice and the length of the stream is 0,
//     we reset the destination map/slice to be a zero-length non-nil map/slice.
//   - Also, if the encoded value is Nil in the stream, then we try to set
//     the container to its "zero" value (e.g. nil for slice/map).
// 
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
	// Note: if stream is set to nil, we set the corresponding value to its "zero" value

	// var ctr int (define this above the  function if trying to do this run)
	// ctr++
	// log(".. [%v] enter decode: rv: %v <==> %T <==> %v", ctr, rv, rv.Interface(), rv.Interface())
	// defer func(ctr2 int) {
	// 	log(".... [%v] exit decode: rv: %v <==> %T <==> %v", ctr2, rv, rv.Interface(), rv.Interface())
	// }(ctr)
	dd := d.d //so we don't dereference constantly
	dd.initReadNext()

	rvOrig := rv
	wasNilIntf := rv.Kind() == reflect.Interface && rv.IsNil()
	rt := rv.Type()

	var ndesc decodeNakedContext
	//if nil interface, use some hieristics to set the nil interface to an
	//appropriate value based on the first byte read (byte descriptor bd)
	if wasNilIntf {
		if dd.currentIsNil() {
			return
		}
		//Prevent from decoding into e.g. error, io.Reader, etc if it's nil and non-nil value in stream.
		//We can only decode into interface{} (0 methods). Else reflect.Set fails later.
		if num := rt.NumMethod(); num > 0 {
			decErr("decodeValue: Cannot decode non-nil codec value into nil %v (%v methods)", rt, num)
		} else {
			rv, ndesc = dd.decodeNaked(d.h)
			if ndesc == dncHandled {
				rvOrig.Set(rv)
				return
			}
			rt = rv.Type()
		}
	} else if dd.currentIsNil() {
		// Note: if stream is set to nil, we set the dereferenced value to its "zero" value (if settable).
		for rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.CanSet() {
			rv.Set(reflect.Zero(rv.Type()))
		}
		return
	}

	// An extension can be registered for any type, regardless of the Kind
	// (e.g. type BitSet int64, type MyStruct { / * unexported fields * / }, type X []int, etc.
	//
	// We can't check if it's an extension byte here first, because the user may have
	// registered a pointer or non-pointer type, meaning we may have to recurse first
	// before matching a mapped type, even though the extension byte is already detected.
	//
	// If we are checking for builtin or ext type here, it means we didn't go through decodeNaked,
	// Because decodeNaked would have handled it. It also means wasNilIntf = false.
	if dd.decodeBuiltinType(rt, rv) {
		return
	}
	if bfnTag, bfnFn := d.h.getDecodeExt(rt); bfnFn != nil {
		xbs := dd.decodeExt(bfnTag)
		if fnerr := bfnFn(rv, xbs); fnerr != nil {
			panic(fnerr)
		}
		return
	}

	// NOTE: if decoding into a nil interface{}, we return a non-nil
	// value except even if the container registers a length of 0.
	//
	// NOTE: Do not make blocks for struct, slice, map, etc individual methods.
	// It ends up being more expensive, because they recursively calls decodeValue
	//
	// (Mar 7, 2013. DON'T REARRANGE ... code clarity)
	// tried arranging in sequence of most probable ones.
	// string, bool, integer, float, struct, ptr, slice, array, map, interface, uint.
	switch rk := rv.Kind(); rk {
	case reflect.String:
		rv.SetString(dd.decodeString())
	case reflect.Bool:
		rv.SetBool(dd.decodeBool())
	case reflect.Int:
		rv.SetInt(dd.decodeInt(intBitsize))
	case reflect.Int64:
		rv.SetInt(dd.decodeInt(64))
	case reflect.Int32:
		rv.SetInt(dd.decodeInt(32))
	case reflect.Int8:
		rv.SetInt(dd.decodeInt(8))
	case reflect.Int16:
		rv.SetInt(dd.decodeInt(16))
	case reflect.Float32:
		rv.SetFloat(dd.decodeFloat(true))
	case reflect.Float64:
		rv.SetFloat(dd.decodeFloat(false))
	case reflect.Uint8:
		rv.SetUint(dd.decodeUint(8))
	case reflect.Uint64:
		rv.SetUint(dd.decodeUint(64))
	case reflect.Uint:
		rv.SetUint(dd.decodeUint(uintBitsize))
	case reflect.Uint32:
		rv.SetUint(dd.decodeUint(32))
	case reflect.Uint16:
		rv.SetUint(dd.decodeUint(16))
	case reflect.Ptr:
		if rv.IsNil() {
			if wasNilIntf {
				rv = reflect.New(rt.Elem())
			} else {
				rv.Set(reflect.New(rt.Elem()))
			}
		}
		d.decodeValue(rv.Elem())
	case reflect.Interface:
		d.decodeValue(rv.Elem())
	case reflect.Struct:
		containerLen := dd.readMapLen()

		if containerLen == 0 {
			break
		}

		sfi := getStructFieldInfos(rt)
		for j := 0; j < containerLen; j++ {
			// var rvkencname string
			// ddecode(&rvkencname)
			dd.initReadNext()
			rvkencname := dd.decodeString()
			// rvksi := sfi.getForEncName(rvkencname)
			if k := sfi.indexForEncName(rvkencname); k > -1 {
				sfik := sfi[k]
				if sfik.i > -1 {
					d.decodeValue(rv.Field(int(sfik.i)))
				} else {
					d.decodeValue(rv.FieldByIndex(sfik.is))
				}
				// d.decodeValue(sfi.field(k, rv))
			} else {
				if d.h.errorIfNoField() {
					decErr("No matching struct field found when decoding stream map with key: %v", rvkencname)
				} else {
					var nilintf0 interface{}
					d.decodeValue(reflect.ValueOf(&nilintf0).Elem())
				}
			}
		}
	case reflect.Slice:
		// Be more careful calling Set() here, because a reflect.Value from an array
		// may have come in here (which may not be settable).
		// In places where the slice got from an array could be, we should guard with CanSet() calls.

		if rt == byteSliceTyp { // rawbytes
			if bs2, changed2 := dd.decodeBytes(rv.Bytes()); changed2 {
				rv.SetBytes(bs2)
			}
			if wasNilIntf && rv.IsNil() {
				rv.SetBytes([]byte{})
			}
			break
		}

		containerLen := dd.readArrayLen()

		if wasNilIntf {
			rv = reflect.MakeSlice(rt, containerLen, containerLen)
		}
		if containerLen == 0 {
			if rv.IsNil() {
				rv.Set(reflect.MakeSlice(rt, containerLen, containerLen))
			} 
			break
		}

		if rv.IsNil() {
			// wasNilIntf only applies if rv is nil (since that's what we did earlier)
			rv.Set(reflect.MakeSlice(rt, containerLen, containerLen))
		} else {
			// if we need to reset rv but it cannot be set, we should err out.
			// for example, if slice is got from unaddressable array, CanSet = false
			if rvcap, rvlen := rv.Len(), rv.Cap(); containerLen > rvcap {
				if rv.CanSet() {
					rvn := reflect.MakeSlice(rt, containerLen, containerLen)
					if rvlen > 0 {
						reflect.Copy(rvn, rv)
					}
					rv.Set(rvn)
				} else {
					decErr("Cannot reset slice with less cap: %v that stream contents: %v", rvcap, containerLen)
				}
			} else if containerLen > rvlen {
				rv.SetLen(containerLen)
			}
		}
		for j := 0; j < containerLen; j++ {
			d.decodeValue(rv.Index(j))
		}
	case reflect.Array:
		d.decodeValue(rv.Slice(0, rv.Len()))
	case reflect.Map:
		containerLen := dd.readMapLen()

		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rt))
		}
		
		if containerLen == 0 {
			break
		}

		ktype, vtype := rt.Key(), rt.Elem()
		for j := 0; j < containerLen; j++ {
			rvk := reflect.New(ktype).Elem()
			d.decodeValue(rvk)

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

			d.decodeValue(rvv)
			rv.SetMapIndex(rvk, rvv)
		}
	default:
		decErr("Unhandled value for kind: %v: %s", rk, msgBadDesc)
	}

	if wasNilIntf {
		rvOrig.Set(rv)
	}
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
	if _, err := io.ReadFull(z.r, bs); err != nil {
		panic(err)
	}
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if _, err := io.ReadFull(z.r, bs); err != nil {
		panic(err)
	}
}

func (z *ioDecReader) readn1() uint8 {
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
		doPanic(msgTagDec, "Trying to read %v bytes. Only %v available", n, z.a)
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

