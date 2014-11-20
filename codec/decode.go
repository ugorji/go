// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"encoding"
	"io"
	"reflect"
	// "runtime/debug"
)

// Some tagging information for error messages.
const (
	msgTagDec             = "codec.decoder"
	msgBadDesc            = "Unrecognized descriptor byte"
	msgDecCannotExpandArr = "cannot expand go array from %v to stream length: %v"
)

// fastpathsDec holds the rtid (reflect.Type Pointer) to fast decode function for a selected slice/map type.
var fastpathsDec = make(map[uintptr]func(*decFnInfo, reflect.Value))

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface {
	unreadn1()
	readn(n int) []byte
	readb([]byte)
	readn1() uint8
	readUint16() uint16
	readUint32() uint32
	readUint64() uint64
}

type decDriver interface {
	initReadNext()
	// this will call initReadNext implicitly, and then check if the next token is a break.
	checkBreak() bool
	tryDecodeAsNil() bool
	// check if a container type: vt is one of: Bytes, String, Nil, Slice or Map
	isContainerType(vt valueType) bool
	isBuiltinType(rt uintptr) bool
	decodeBuiltin(rt uintptr, v interface{})
	//decodeNaked: Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	//for extensions, decodeNaked must completely decode them as a *RawExt.
	decodeNaked(*Decoder) (v interface{}, vt valueType, decodeFurther bool)
	decodeInt(bitsize uint8) (i int64)
	decodeUint(bitsize uint8) (ui uint64)
	decodeFloat(chkOverflow32 bool) (f float64)
	decodeBool() (b bool)
	// decodeString can also decode symbols
	decodeString() (s string)
	decodeBytes(bs []byte) (bsOut []byte, changed bool)
	// decodeExt will decode into a *RawExt or into an extension.
	decodeExt(rv reflect.Value, xtag uint64, ext Ext, d *Decoder) (realxtag uint64)
	// decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)
	readMapStart() int
	readArrayStart() int
	readMapEnd()
	readArrayEnd()
	readArrayEntrySeparator()
	readMapEntrySeparator()
	readMapKVSeparator()
}

// decDrivers may implement this interface to bypass allocation
type decDriverStringAsBytes interface {
	decStringAsBytes(bs []byte) []byte
}

type decNoMapArrayEnd struct{}

func (_ decNoMapArrayEnd) readMapEnd()   {}
func (_ decNoMapArrayEnd) readArrayEnd() {}

type decNoMapArraySeparator struct{}

func (_ decNoMapArraySeparator) readArrayEntrySeparator() {}
func (_ decNoMapArraySeparator) readMapEntrySeparator()   {}
func (_ decNoMapArraySeparator) readMapKVSeparator()      {}

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
	// If true, use the int64 during schema-less decoding of unsigned values (not uint64).
	SignedInteger bool
}

// ------------------------------------

// ioDecReader is a decReader that reads off an io.Reader
type ioDecReader struct {
	r  io.Reader
	br io.ByteScanner
	x  [8]byte //temp byte array re-used internally for efficiency
	l  byte    // last byte
	ls uint8   // last byte status: 0: unset, 1: read, 2: unread
}

func (z *ioDecReader) readn(n int) (bs []byte) {
	if n <= 0 {
		return
	}
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
		if b, err := z.br.ReadByte(); err == nil {
			return b
		} else {
			panic(err)
		}
	}
	if z.ls == 2 {
		z.ls = 0
	} else {
		z.readb(z.x[:1])
		z.l = z.x[0]
		z.ls = 1
	}
	return z.l
}

func (z *ioDecReader) unreadn1() {
	if z.br != nil {
		if err := z.br.UnreadByte(); err != nil {
			panic(err)
		}
		return
	}
	if z.ls == 2 {
		decErr("cannot unread when last byte has been unread")
	}
	if z.ls == 1 {
		z.ls = 2
		return
	}
	decErr("cannot unread when no byte has been read")
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

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct {
	b []byte // data
	c int    // cursor
	a int    // available
}

func (z *bytesDecReader) consume(n int) (oldcursor int) {
	if z.a == 0 {
		panic(io.EOF)
	}
	if n > z.a {
		decErr("cannot read %v bytes, when only %v available", n, z.a)
	}
	// z.checkAvailable(n)
	oldcursor = z.c
	z.c = oldcursor + n
	z.a = z.a - n
	return
}

func (z *bytesDecReader) unreadn1() {
	if z.c == 0 || len(z.b) == 0 {
		decErr("cannot unread last byte read")
	}
	z.c--
	z.a++
	return
}

func (z *bytesDecReader) readn(n int) (bs []byte) {
	if n <= 0 {
		return
	}
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

// ------------------------------------

// decFnInfo has methods for registering handling decoding of a specific type
// based on some characteristics (builtin, extension, reflect Kind, etc)
type decFnInfo struct {
	ti    *typeInfo
	d     *Decoder
	dd    decDriver
	xfFn  Ext
	xfTag uint64
	array bool
}

// ----------------------------------------

type decFn struct {
	i *decFnInfo
	f func(*decFnInfo, reflect.Value)
}

func (f *decFnInfo) builtin(rv reflect.Value) {
	f.dd.decodeBuiltin(f.ti.rtid, rv.Addr().Interface())
}

func (f *decFnInfo) rawExt(rv reflect.Value) {
	f.dd.decodeExt(rv, 0, nil, f.d)
}

func (f *decFnInfo) ext(rv reflect.Value) {
	f.dd.decodeExt(rv, f.xfTag, f.xfFn, f.d)
}

func (f *decFnInfo) binaryUnmarshal(rv reflect.Value) {
	var bm encoding.BinaryUnmarshaler
	if f.ti.bunmIndir == -1 {
		bm = rv.Addr().Interface().(encoding.BinaryUnmarshaler)
	} else if f.ti.bunmIndir == 0 {
		bm = rv.Interface().(encoding.BinaryUnmarshaler)
	} else {
		for j, k := int8(0), f.ti.bunmIndir; j < k; j++ {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		bm = rv.Interface().(encoding.BinaryUnmarshaler)
	}
	xbs, _ := f.dd.decodeBytes(nil)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) textUnmarshal(rv reflect.Value) {
	var tm encoding.TextUnmarshaler
	if f.ti.tunmIndir == -1 {
		tm = rv.Addr().Interface().(encoding.TextUnmarshaler)
	} else if f.ti.tunmIndir == 0 {
		tm = rv.Interface().(encoding.TextUnmarshaler)
	} else {
		for j, k := int8(0), f.ti.tunmIndir; j < k; j++ {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		tm = rv.Interface().(encoding.TextUnmarshaler)
	}
	var fnerr error
	if sb, sbok := f.dd.(decDriverStringAsBytes); sbok {
		fnerr = tm.UnmarshalText(sb.decStringAsBytes(f.d.b[:0]))
	} else {
		fnerr = tm.UnmarshalText([]byte(f.dd.decodeString()))
	}
	if fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) kErr(rv reflect.Value) {
	decErr("no decoding function defined for kind %v", rv.Kind())
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
	// debugf("\t===> kInterface")
	if !rv.IsNil() {
		f.d.decodeValue(rv.Elem(), decFn{})
		return
	}
	// nil interface:
	// use some hieristics to set the nil interface to an
	// appropriate value based on the first byte read (byte descriptor bd)
	v, vt, decodeFurther := f.dd.decodeNaked(f.d)
	if vt == valueTypeNil {
		return
	}
	// Cannot decode into nil interface with methods (e.g. error, io.Reader, etc)
	// if non-nil value in stream.
	if num := f.ti.rt.NumMethod(); num > 0 {
		decErr("cannot decode non-nil codec value into nil %v (%v methods)",
			f.ti.rt, num)
	}
	var rvn reflect.Value
	var useRvn bool
	switch vt {
	case valueTypeMap:
		if f.d.h.MapType == nil {
			var m2 map[interface{}]interface{}
			v = &m2
		} else {
			rvn = reflect.New(f.d.h.MapType).Elem()
			useRvn = true
		}
	case valueTypeArray:
		if f.d.h.SliceType == nil {
			var m2 []interface{}
			v = &m2
		} else {
			rvn = reflect.New(f.d.h.SliceType).Elem()
			useRvn = true
		}
	case valueTypeExt:
		re := v.(*RawExt)
		bfn := f.d.h.getExtForTag(re.Tag)
		if bfn == nil {
			rvn = reflect.ValueOf(*re)
		} else {
			rvn = reflect.New(bfn.rt).Elem()
			if re.Data != nil {
				bfn.ext.ReadExt(rvn, re.Data)
			} else {
				bfn.ext.UpdateExt(rvn, re.Value)
			}
		}
		rv.Set(rvn)
		return
	}
	if decodeFurther {
		if useRvn {
			f.d.decodeValue(rvn, decFn{})
		} else if v != nil {
			// this v is a pointer, so we need to dereference it when done
			f.d.decode(v)
			rvn = reflect.ValueOf(v).Elem()
			useRvn = true
		}
	}
	if useRvn {
		rv.Set(rvn)
	} else if v != nil {
		rv.Set(reflect.ValueOf(v))
	}
}

func (f *decFnInfo) kStruct(rv reflect.Value) {
	fti := f.ti
	if f.dd.isContainerType(valueTypeMap) {
		containerLen := f.dd.readMapStart()
		if containerLen == 0 {
			f.dd.readMapEnd()
			return
		}
		tisfi := fti.sfi
		for j := 0; ; j++ {
			if containerLen >= 0 {
				if j >= containerLen {
					break
				}
			} else if f.dd.checkBreak() {
				break
			}
			if j > 0 {
				f.dd.readMapEntrySeparator()
			}
			// var rvkencname string
			// ddecode(&rvkencname)
			f.dd.initReadNext()
			rvkencname := f.dd.decodeString()
			f.dd.readMapKVSeparator()
			// rvksi := ti.getForEncName(rvkencname)
			if k := fti.indexForEncName(rvkencname); k > -1 {
				sfik := tisfi[k]
				if sfik.i != -1 {
					f.d.decodeValue(rv.Field(int(sfik.i)), decFn{})
				} else {
					f.d.decEmbeddedField(rv, sfik.is)
				}
				// f.d.decodeValue(ti.field(k, rv))
			} else {
				if f.d.h.ErrorIfNoField {
					decErr("no matching struct field found when decoding stream map with key %v",
						rvkencname)
				} else {
					var nilintf0 interface{}
					f.d.decodeValue(reflect.ValueOf(&nilintf0).Elem(), decFn{})
				}
			}
		}
		f.dd.readMapEnd()
	} else if f.dd.isContainerType(valueTypeArray) {
		containerLen := f.dd.readArrayStart()
		if containerLen == 0 {
			f.dd.readArrayEnd()
			return
		}
		for j, si := range fti.sfip {
			if containerLen >= 0 {
				if j == containerLen {
					break
				}
			} else if f.dd.checkBreak() {
				break
			}
			if j > 0 {
				f.dd.readArrayEntrySeparator()
			}
			if si.i != -1 {
				f.d.decodeValue(rv.Field(int(si.i)), decFn{})
			} else {
				f.d.decEmbeddedField(rv, si.is)
			}
		}
		if containerLen > len(fti.sfip) {
			// read remaining values and throw away
			for j := len(fti.sfip); j < containerLen; j++ {
				var nilintf0 interface{}
				if j > 0 {
					f.dd.readArrayEntrySeparator()
				}
				f.d.decodeValue(reflect.ValueOf(&nilintf0).Elem(), decFn{})
			}
		}
		f.dd.readArrayEnd()
	} else {
		decErr("only encoded map or array can be decoded into a struct")
	}
}

func (f *decFnInfo) kSlice(rv reflect.Value) {
	// A slice can be set from a map or array in stream. This way, the order can be kept (as order is lost with map).
	if f.dd.isContainerType(valueTypeBytes) || f.dd.isContainerType(valueTypeString) {
		if f.ti.rtid == uint8SliceTypId || f.ti.rt.Elem().Kind() == reflect.Uint8 {
			rvbs := rv.Bytes()
			if bs2, changed2 := f.dd.decodeBytes(rvbs); changed2 {
				if rv.CanSet() {
					rv.SetBytes(bs2)
				} else {
					copy(rvbs, bs2)
				}
			}
			return
		}
	}

	slh := decSliceHelper{dd: f.dd}
	containerLenS := slh.start()

	// an array can never return a nil slice. so no need to check f.array here.
	if rv.IsNil() {
		if containerLenS < 0 {
			rv.Set(reflect.MakeSlice(f.ti.rt, 0, 0))
		} else {
			rv.Set(reflect.MakeSlice(f.ti.rt, containerLenS, containerLenS))
		}
	}

	if containerLenS == 0 {
		rv.SetLen(0)
		f.dd.readArrayEnd()
		return
	}

	rv0 := rv
	rvChanged := false

	if rvcap, rvlen := rv.Len(), rv.Cap(); containerLenS > rvcap {
		if f.array { // !rv.CanSet()
			decErr(msgDecCannotExpandArr, rvcap, containerLenS)
		}
		rv = reflect.MakeSlice(f.ti.rt, containerLenS, containerLenS)
		if rvlen > 0 {
			reflect.Copy(rv, rv0)
		}
		rvChanged = true
	} else if containerLenS > rvlen {
		rv.SetLen(containerLenS)
	}

	rtelem0 := f.ti.rt.Elem()
	rtelem := rtelem0
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}
	fn := f.d.getDecFn(rtelem)
	// for j := 0; j < containerLenS; j++ {

	for j := 0; ; j++ {
		if containerLenS >= 0 {
			if j >= containerLenS {
				break
			}
		} else if f.dd.checkBreak() {
			break
		}
		// if indefinite, etc, then expand the slice if necessary
		if j >= rv.Len() {
			// TODO: Need to update this appropriately
			rv = reflect.Append(rv, reflect.Zero(rtelem0))
			rvChanged = true
		}
		if j > 0 {
			slh.sep(j)
		}
		f.d.decodeValue(rv.Index(j), fn)
	}
	slh.end()
	if rvChanged {
		rv0.Set(rv)
	}
}

func (f *decFnInfo) kArray(rv reflect.Value) {
	// f.d.decodeValue(rv.Slice(0, rv.Len()))
	f.kSlice(rv.Slice(0, rv.Len()))
}

func (f *decFnInfo) kMap(rv reflect.Value) {
	containerLen := f.dd.readMapStart()

	if rv.IsNil() {
		rv.Set(reflect.MakeMap(f.ti.rt))
	}

	if containerLen == 0 {
		f.dd.readMapEnd()
		return
	}

	ktype, vtype := f.ti.rt.Key(), f.ti.rt.Elem()
	ktypeId := reflect.ValueOf(ktype).Pointer()
	var keyFn, valFn decFn
	var xtyp reflect.Type
	for xtyp = ktype; xtyp.Kind() == reflect.Ptr; xtyp = xtyp.Elem() {
	}
	keyFn = f.d.getDecFn(xtyp)
	for xtyp = vtype; xtyp.Kind() == reflect.Ptr; xtyp = xtyp.Elem() {
	}
	valFn = f.d.getDecFn(xtyp)
	// for j := 0; j < containerLen; j++ {
	for j := 0; ; j++ {
		if containerLen >= 0 {
			if j >= containerLen {
				break
			}
		} else if f.dd.checkBreak() {
			break
		}
		if j > 0 {
			f.dd.readMapEntrySeparator()
		}
		rvk := reflect.New(ktype).Elem()
		f.d.decodeValue(rvk, keyFn)

		// special case if a byte array.
		if ktypeId == intfTypId {
			rvk = rvk.Elem()
			if rvk.Type() == uint8SliceTyp {
				rvk = reflect.ValueOf(string(rvk.Bytes()))
			}
		}
		rvv := rv.MapIndex(rvk)
		if !rvv.IsValid() {
			rvv = reflect.New(vtype).Elem()
		}
		f.dd.readMapKVSeparator()
		f.d.decodeValue(rvv, valFn)
		rv.SetMapIndex(rvk, rvv)
	}
	f.dd.readMapEnd()
}

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct {
	r  decReader
	d  decDriver
	h  *BasicHandle
	hh Handle
	f  map[uintptr]decFn
	x  []uintptr
	s  []decFn
	b  [32]byte
}

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to pass in a memory buffered reader
// (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder {
	z := ioDecReader{
		r: r,
	}
	z.br, _ = r.(io.ByteScanner)
	return &Decoder{r: &z, hh: h, h: h.getBasicHandle(), d: h.newDecDriver(&z)}
}

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder {
	z := bytesDecReader{
		b: in,
		a: len(in),
	}
	return &Decoder{r: &z, hh: h, h: h.getBasicHandle(), d: h.newDecDriver(&z)}
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
//   - Other values are decoded appropriately depending on the type:
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
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// When decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
func (d *Decoder) Decode(v interface{}) (err error) {
	defer panicToErr(&err)
	d.decode(v)
	return
}

// MustDecode is like Decode, but panics if unable to Decode.
// This provides insight to the code location that triggered the error.
func (d *Decoder) MustDecode(v interface{}) {
	d.decode(v)
}

func (d *Decoder) decode(iv interface{}) {
	d.d.initReadNext()

	switch v := iv.(type) {
	case nil:
		decErr("cannot decode into nil.")

	case reflect.Value:
		d.chkPtrValue(v)
		d.decodeValue(v.Elem(), decFn{})

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
		*v = d.d.decodeInt(64)
	case *uint:
		*v = uint(d.d.decodeUint(uintBitsize))
	case *uint8:
		*v = uint8(d.d.decodeUint(8))
	case *uint16:
		*v = uint16(d.d.decodeUint(16))
	case *uint32:
		*v = uint32(d.d.decodeUint(32))
	case *uint64:
		*v = d.d.decodeUint(64)
	case *float32:
		*v = float32(d.d.decodeFloat(true))
	case *float64:
		*v = d.d.decodeFloat(false)
	case *[]uint8:
		*v, _ = d.d.decodeBytes(*v)

	case *interface{}:
		d.decodeValue(reflect.ValueOf(iv).Elem(), decFn{})

	default:
		rv := reflect.ValueOf(iv)
		d.chkPtrValue(rv)
		d.decodeValue(rv.Elem(), decFn{})
	}
}

func (d *Decoder) decodeValue(rv reflect.Value, fn decFn) {
	d.d.initReadNext()

	if d.d.tryDecodeAsNil() {
		// If value in stream is nil, set the dereferenced value to its "zero" value (if settable)
		if rv.Kind() == reflect.Ptr {
			if !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return
		}
		// for rv.Kind() == reflect.Ptr {
		// 	rv = rv.Elem()
		// }
		if rv.IsValid() { // rv.CanSet() // always settable, except it's invalid
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

	if fn.i == nil {
		fn = d.getDecFn(rv.Type())
	}
	fn.f(fn.i, rv)
	return
}

func (d *Decoder) getDecFn(rt reflect.Type) (fn decFn) {
	rtid := reflect.ValueOf(rt).Pointer()

	// retrieve or register a focus'ed function for this type
	// to eliminate need to do the retrieval multiple times

	// if d.f == nil && d.s == nil { debugf("---->Creating new dec f map for type: %v\n", rt) }
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
		fi := decFnInfo{ti: getTypeInfo(rtid, rt), d: d, dd: d.d}
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
		} else if xfFn := d.h.getExt(rtid); xfFn != nil {
			fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
			fn.f = (*decFnInfo).ext
		} else if supportMarshalInterfaces && d.hh.isBinaryEncoding() && fi.ti.bunm {
			fn.f = (*decFnInfo).binaryUnmarshal
		} else if supportMarshalInterfaces && !d.hh.isBinaryEncoding() && fi.ti.tunm {
			fn.f = (*decFnInfo).textUnmarshal
		} else {
			rk := rt.Kind()
			if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice) {
				if fn.f, ok = fastpathsDec[rtid]; !ok && rt.PkgPath() != "" {
					// use mapping for underlying type if there
					var rtu reflect.Type
					if rk == reflect.Map {
						rtu = reflect.MapOf(rt.Key(), rt.Elem())
					} else {
						rtu = reflect.SliceOf(rt.Elem())
					}
					rtuid := reflect.ValueOf(rtu).Pointer()
					if fn.f, ok = fastpathsDec[rtuid]; ok {
						xfnf := fn.f
						xrt := fastpathsTyp[rtuid]
						fn.f = func(xf *decFnInfo, xrv reflect.Value) {
							// xfnf(xf, xrv.Convert(xrt))
							xfnf(xf, xrv.Addr().Convert(reflect.PtrTo(xrt)).Elem())
						}
					}
				}
			}
			if fn.f == nil {
				switch rk {
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
					fi.array = true
					fn.f = (*decFnInfo).kArray
				case reflect.Map:
					fn.f = (*decFnInfo).kMap
				default:
					fn.f = (*decFnInfo).kErr
				}
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

	return
}

func (d *Decoder) chkPtrValue(rv reflect.Value) {
	// We can only decode into a non-nil pointer
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return
	}
	if !rv.IsValid() {
		decErr("cannot decode into a zero (ie invalid) reflect.Value")
	}
	if !rv.CanInterface() {
		decErr("cannot decode into a value without an interface: %v", rv)
	}
	rvi := rv.Interface()
	decErr("cannot decode into non-pointer or nil pointer. Got: %v, %T, %v",
		rv.Kind(), rvi, rvi)
}

func (d *Decoder) decEmbeddedField(rv reflect.Value, index []int) {
	// d.decodeValue(rv.FieldByIndex(index))
	// nil pointers may be here; so reproduce FieldByIndex logic + enhancements
	for _, j := range index {
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			// If a pointer, it must be a pointer to struct (based on typeInfo contract)
			rv = rv.Elem()
		}
		rv = rv.Field(j)
	}
	d.decodeValue(rv, decFn{})
}

// --------------------------------------------------

// decSliceHelper assists when decoding into a slice, from a map or an array in the stream.
// A slice can be set from a map or array in stream. This supports the MapBySlice interface.
type decSliceHelper struct {
	dd decDriver
	ct valueType
}

func (x *decSliceHelper) start() (sliceLen int) {
	if x.dd.isContainerType(valueTypeArray) {
		x.ct = valueTypeArray
		return x.dd.readArrayStart()
	}
	if x.dd.isContainerType(valueTypeMap) {
		x.ct = valueTypeMap
		return x.dd.readMapStart() * 2
	}
	decErr("only encoded map or array can be decoded into a slice")
	panic("unreachable")
}

func (x *decSliceHelper) sep(index int) {
	if x.ct == valueTypeArray {
		x.dd.readArrayEntrySeparator()
	} else {
		if index%2 == 0 {
			x.dd.readMapEntrySeparator()
		} else {
			x.dd.readMapKVSeparator()
		}
	}
}

func (x *decSliceHelper) end() {
	if x.ct == valueTypeArray {
		x.dd.readArrayEnd()
	} else {
		x.dd.readMapEnd()
	}
}

func decErr(format string, params ...interface{}) {
	doPanic(msgTagDec, format, params...)
}
