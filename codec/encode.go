// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	//"bufio"
	"io"
	"reflect"
	//"fmt"
)

//var _ = fmt.Printf
const (
	// Some tagging information for error messages.
	msgTagEnc         = "codec.encoder"
	defEncByteBufSize = 1 << 6 // 4:16, 6:64, 8:256, 10:1024
	// maxTimeSecs32 = math.MaxInt32 / 60 / 24 / 366
)

// encWriter abstracting writing to a byte array or to an io.Writer.
type encWriter interface {
	writeUint16(uint16)
	writeUint32(uint32)
	writeUint64(uint64)
	writeb([]byte)
	writestr(string)
	writen1(byte)
	writen2(byte, byte)
	atEndOfEncode()
}

// encDriver abstracts the actual codec (binc vs msgpack, etc)
type encDriver interface {
	isBuiltinType(rt uintptr) bool
	encodeBuiltinType(rt uintptr, rv reflect.Value)
	encodeNil()
	encodeInt(i int64)
	encodeUint(i uint64)
	encodeBool(b bool)
	encodeFloat32(f float32)
	encodeFloat64(f float64)
	encodeExtPreamble(xtag byte, length int)
	encodeArrayPreamble(length int)
	encodeMapPreamble(length int)
	encodeString(c charEncoding, v string)
	encodeSymbol(v string)
	encodeStringBytes(c charEncoding, v []byte)
	//TODO
	//encBignum(f *big.Int)
	//encStringRunes(c charEncoding, v []rune)
}

// encodeHandleI is the interface that the encode functions need.
type encodeHandleI interface {
	getEncodeExt(rt uintptr) (tag byte, fn func(reflect.Value) ([]byte, error))
	writeExt() bool
	structToArray() bool
}

type encFnInfo struct {
	sis   *typeInfo
	e     *Encoder
	ee    encDriver
	rt    reflect.Type
	rtid  uintptr
	xfFn  func(reflect.Value) ([]byte, error)
	xfTag byte 
}

// encFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the 
// code block that should be called (encapsulated in a function) 
// instead of executing the checks every time.
type encFn struct {
	i *encFnInfo
	f func(*encFnInfo, reflect.Value)
}

// An Encoder writes an object to an output stream in the codec format.
type Encoder struct {
	w encWriter
	e encDriver
	h encodeHandleI
	f map[uintptr]encFn
}

type ioEncWriterWriter interface {
	WriteByte(c byte) error
	WriteString(s string) (n int, err error)
	Write(p []byte) (n int, err error)
}

type ioEncStringWriter interface {
	WriteString(s string) (n int, err error)
}

type simpleIoEncWriterWriter struct {
	w io.Writer 
	bw io.ByteWriter
	sw ioEncStringWriter
}

// ioEncWriter implements encWriter and can write to an io.Writer implementation
type ioEncWriter struct {
	w ioEncWriterWriter
	x [8]byte // temp byte array re-used internally for efficiency
}

// bytesEncWriter implements encWriter and can write to an byte slice.
// It is used by Marshal function.
type bytesEncWriter struct {
	b   []byte
	c   int     // cursor
	out *[]byte // write out on atEndOfEncode
}

type EncodeOptions struct {
	// Encode a struct as an array, and not as a map.
	StructToArray bool
}

func (o *simpleIoEncWriterWriter) WriteByte(c byte) (err error) {
	if o.bw != nil {
		return o.bw.WriteByte(c)
	}
	_, err = o.w.Write([]byte{c})
	return
}

func (o *simpleIoEncWriterWriter) WriteString(s string) (n int, err error) {
	if o.sw != nil {
		return o.sw.WriteString(s)
	}
	return o.w.Write([]byte(s))
}

func (o *simpleIoEncWriterWriter) Write(p []byte) (n int, err error) {
	return o.w.Write(p)
}


func (o *EncodeOptions) structToArray() bool {
	return o.StructToArray
}

func (f *encFnInfo) builtin(rv reflect.Value) {
	baseRv := rv
	for j := int8(0); j < f.sis.baseIndir; j++ {
		baseRv = baseRv.Elem()
	}
	f.ee.encodeBuiltinType(f.sis.baseId, baseRv)
}

func (f *encFnInfo) ext(rv reflect.Value) {
	baseRv := rv
	for j := int8(0); j < f.sis.baseIndir; j++ {
		baseRv = baseRv.Elem()
	}
	bs, fnerr := f.xfFn(baseRv)
	if fnerr != nil {
		panic(fnerr)
	}
	if bs == nil {
		f.ee.encodeNil()
		return
	}
	if f.e.h.writeExt() {
		f.ee.encodeExtPreamble(f.xfTag, len(bs))
		f.e.w.writeb(bs)
	} else {
		f.ee.encodeStringBytes(c_RAW, bs)
	}
	
}

func (f *encFnInfo) binaryMarshal(rv reflect.Value) {
	var bm binaryMarshaler
	if f.sis.mIndir == 0 {
		bm = rv.Interface().(binaryMarshaler)
	} else if f.sis.mIndir == -1 {
		bm = rv.Addr().Interface().(binaryMarshaler)
	} else {
		rv2 := rv
		for j := int8(0); j < f.sis.mIndir; j++ {
			rv2 = rv.Elem()
		}
		bm = rv2.Interface().(binaryMarshaler)
	}
	// debugf(">>>> binaryMarshaler: %T", rv.Interface())
	bs, fnerr := bm.MarshalBinary()
	if fnerr != nil {
		panic(fnerr)
	}
	if bs == nil {
		f.ee.encodeNil()
	} else {
		f.ee.encodeStringBytes(c_RAW, bs)
	}

}

func (f *encFnInfo) kBool(rv reflect.Value) {
	f.ee.encodeBool(rv.Bool())
}

func (f *encFnInfo) kString(rv reflect.Value) {
	f.ee.encodeString(c_UTF8, rv.String())
}

func (f *encFnInfo) kFloat64(rv reflect.Value) {
	f.ee.encodeFloat64(rv.Float())
}

func (f *encFnInfo) kFloat32(rv reflect.Value) {
	f.ee.encodeFloat32(float32(rv.Float()))
}

func (f *encFnInfo) kInt(rv reflect.Value) {
	f.ee.encodeInt(rv.Int())
}

func (f *encFnInfo) kUint(rv reflect.Value) {
	f.ee.encodeUint(rv.Uint())
}

func (f *encFnInfo) kInvalid(rv reflect.Value) {
	f.ee.encodeNil()
}

func (f *encFnInfo) kErr(rv reflect.Value) {
	encErr("Unsupported kind: %s, for: %#v", rv.Kind(), rv)
}

func (f *encFnInfo) kSlice(rv reflect.Value) {
	if rv.IsNil() {
		f.ee.encodeNil()
		return
	}
	if f.rtid == byteSliceTypId {
		f.ee.encodeStringBytes(c_RAW, rv.Bytes())
		return
	}
	l := rv.Len()
	f.ee.encodeArrayPreamble(l)
	if l == 0 {
		return
	}
	for j := 0; j < l; j++ {
		f.e.encodeValue(rv.Index(j))
	}
}

func (f *encFnInfo) kArray(rv reflect.Value) {
	f.e.encodeValue(rv.Slice(0, rv.Len()))
}

func (f *encFnInfo) kStruct(rv reflect.Value) {
	newlen := len(f.sis.sis)
	rvals := make([]reflect.Value, newlen)
	var encnames []string
	e := f.e
	sissis := f.sis.sisp
	toMap := !(f.sis.toArray || e.h.structToArray())
	// if toMap, use the sorted array. If toArray, use unsorted array (to match sequence in struct)
	if toMap {
		sissis = f.sis.sis
		encnames = make([]string, newlen)
	}
	newlen = 0
	for _, si := range sissis {
		if si.i != -1 {
			rvals[newlen] = rv.Field(int(si.i))
		} else {
			rvals[newlen] = rv.FieldByIndex(si.is)
		}
		if toMap {
			if si.omitEmpty && isEmptyValue(rvals[newlen]) {
				continue
			}
			encnames[newlen] = si.encName
		} else {
			if si.omitEmpty && isEmptyValue(rvals[newlen]) {
				rvals[newlen] = reflect.Value{} //encode as nil
			}
		}
		newlen++
	}

	if toMap {
		ee := f.ee //don't dereference everytime
		ee.encodeMapPreamble(newlen)
		for j := 0; j < newlen; j++ {
			ee.encodeSymbol(encnames[j])
			e.encodeValue(rvals[j])
		}
	} else {
		f.ee.encodeArrayPreamble(newlen)
		for j := 0; j < newlen; j++ {
			e.encodeValue(rvals[j])
		}
	}
}

func (f *encFnInfo) kPtr(rv reflect.Value) {
	if rv.IsNil() {
		f.ee.encodeNil()
		return
	}
	f.e.encodeValue(rv.Elem())
}

func (f *encFnInfo) kInterface(rv reflect.Value) {
	if rv.IsNil() {
		f.ee.encodeNil()
		return
	}
	f.e.encodeValue(rv.Elem())
}

func (f *encFnInfo) kMap(rv reflect.Value) {
	if rv.IsNil() {
		f.ee.encodeNil()
		return
	}
	l := rv.Len()
	f.ee.encodeMapPreamble(l)
	if l == 0 {
		return
	}
	keyTypeIsString := f.rt.Key().Kind() == reflect.String
	mks := rv.MapKeys()
	// for j, lmks := 0, len(mks); j < lmks; j++ {
	for j := range mks {
		if keyTypeIsString {
			f.ee.encodeSymbol(mks[j].String())
		} else {
			f.e.encodeValue(mks[j])
		}
		f.e.encodeValue(rv.MapIndex(mks[j]))
	}

}



// NewEncoder returns an Encoder for encoding into an io.Writer.
// 
// For efficiency, Users are encouraged to pass in a memory buffered writer
// (eg bufio.Writer, bytes.Buffer). 
func NewEncoder(w io.Writer, h Handle) *Encoder {
	ww, ok := w.(ioEncWriterWriter)
	if !ok {
		sww := simpleIoEncWriterWriter{w: w}
		sww.bw, _ = w.(io.ByteWriter)
		sww.sw, _ = w.(ioEncStringWriter)
		ww = &sww
		//ww = bufio.NewWriterSize(w, defEncByteBufSize)
	}
	z := ioEncWriter{
		w: ww,
	}
	return &Encoder{w: &z, h: h, e: h.newEncDriver(&z) }
}

// NewEncoderBytes returns an encoder for encoding directly and efficiently
// into a byte slice, using zero-copying to temporary slices.
//
// It will potentially replace the output byte slice pointed to.
// After encoding, the out parameter contains the encoded contents.
func NewEncoderBytes(out *[]byte, h Handle) *Encoder {
	in := *out
	if in == nil {
		in = make([]byte, defEncByteBufSize)
	}
	z := bytesEncWriter{
		b:   in,
		out: out,
	}
	return &Encoder{w: &z, h: h, e: h.newEncDriver(&z) }
}

// Encode writes an object into a stream in the codec format.
//
// Encoding can be configured via the "codec" struct tag for the fields.
// 
// The "codec" key in struct field's tag value is the key name,
// followed by an optional comma and options.
// 
// To set an option on all fields (e.g. omitempty on all fields), you
// can create a field called _struct, and set flags on it. 
// 
// Struct values "usually" encode as maps. Each exported struct field is encoded unless:
//    - the field's codec tag is "-", OR
//    - the field is empty and its codec tag specifies the "omitempty" option.
// 
// When encoding as a map, the first string in the tag (before the comma)
// is the map key string to use when encoding.
// 
// However, struct values may encode as arrays. This happens when:
//    - StructToArray Encode option is set, OR
//    - the codec tag on the _struct field sets the "toarray" option
// 
// The empty values (for omitempty option) are false, 0, any nil pointer 
// or interface value, and any array, slice, map, or string of length zero.
//
// Anonymous fields are encoded inline if no struct tag is present.
// Else they are encoded as regular fields.
// 
// Examples:
//
//      type MyStruct struct {
//          _struct bool    `codec:",omitempty"`   //set omitempty for every field
//          Field1 string   `codec:"-"`            //skip this field
//          Field2 int      `codec:"myName"`       //Use key "myName" in encode stream
//          Field3 int32    `codec:",omitempty"`   //use key "Field3". Omit if empty.
//          Field4 bool     `codec:"f4,omitempty"` //use key "f4". Omit if empty.
//          ...
//      }
//      
//      type MyStruct struct {
//          _struct bool    `codec:",omitempty,toarray"`   //set omitempty for every field
//                                                         //and encode struct as an array
//      }   
//
// The mode of encoding is based on the type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryMarshaler, call its MarshalBinary() (data []byte, err error)
//   - Else encode it based on its reflect.Kind
// 
// Note that struct field names and keys in map[string]XXX will be treated as symbols.
// Some formats support symbols (e.g. binc) and will properly encode the string
// only once in the stream, and use a tag to refer to it thereafter. 
func (e *Encoder) Encode(v interface{}) (err error) {
	defer panicToErr(&err)
	e.encode(v)
	e.w.atEndOfEncode()
	return
}

func (e *Encoder) encode(iv interface{}) {
	switch v := iv.(type) {
	case nil:
		e.e.encodeNil()

	case reflect.Value:
		e.encodeValue(v)

	case string:
		e.e.encodeString(c_UTF8, v)
	case bool:
		e.e.encodeBool(v)
	case int:
		e.e.encodeInt(int64(v))
	case int8:
		e.e.encodeInt(int64(v))
	case int16:
		e.e.encodeInt(int64(v))
	case int32:
		e.e.encodeInt(int64(v))
	case int64:
		e.e.encodeInt(v)
	case uint:
		e.e.encodeUint(uint64(v))
	case uint8:
		e.e.encodeUint(uint64(v))
	case uint16:
		e.e.encodeUint(uint64(v))
	case uint32:
		e.e.encodeUint(uint64(v))
	case uint64:
		e.e.encodeUint(v)
	case float32:
		e.e.encodeFloat32(v)
	case float64:
		e.e.encodeFloat64(v)

	case *string:
		e.e.encodeString(c_UTF8, *v)
	case *bool:
		e.e.encodeBool(*v)
	case *int:
		e.e.encodeInt(int64(*v))
	case *int8:
		e.e.encodeInt(int64(*v))
	case *int16:
		e.e.encodeInt(int64(*v))
	case *int32:
		e.e.encodeInt(int64(*v))
	case *int64:
		e.e.encodeInt(*v)
	case *uint:
		e.e.encodeUint(uint64(*v))
	case *uint8:
		e.e.encodeUint(uint64(*v))
	case *uint16:
		e.e.encodeUint(uint64(*v))
	case *uint32:
		e.e.encodeUint(uint64(*v))
	case *uint64:
		e.e.encodeUint(*v)
	case *float32:
		e.e.encodeFloat32(*v)
	case *float64:
		e.e.encodeFloat64(*v)

	default:
		e.encodeValue(reflect.ValueOf(iv))
	}

}

func (e *Encoder) encodeValue(rv reflect.Value) {
	rt := rv.Type()
	rtid := reflect.ValueOf(rt).Pointer()
	
	if e.f == nil {
		// debugf("---->Creating new enc f map for type: %v\n", rt)
		e.f = make(map[uintptr]encFn, 16)
	}

	fn, ok := e.f[rtid]
	if !ok {
		// debugf("\tCreating new enc fn for type: %v\n", rt)
		fi := encFnInfo { sis:getTypeInfo(rtid, rt), e:e, ee:e.e, rt:rt, rtid:rtid }
		if e.e.isBuiltinType(fi.sis.baseId) {
			fn = encFn{ &fi, (*encFnInfo).builtin }
		} else if xfTag, xfFn := e.h.getEncodeExt(fi.sis.baseId); xfFn != nil {
			fi.xfTag, fi.xfFn = xfTag, xfFn
			fn = encFn{ &fi, (*encFnInfo).ext }
		} else if supportBinaryMarshal && fi.sis.m {
			fn = encFn{ &fi, (*encFnInfo).binaryMarshal }
		} else {
			switch rk := rt.Kind(); rk {
			case reflect.Bool:
				fn = encFn{ &fi, (*encFnInfo).kBool }
			case reflect.String:
				fn = encFn{ &fi, (*encFnInfo).kString }
			case reflect.Float64:
				fn = encFn{ &fi, (*encFnInfo).kFloat64 }
			case reflect.Float32:
				fn = encFn{ &fi, (*encFnInfo).kFloat32 }
			case reflect.Int, reflect.Int8, reflect.Int64, reflect.Int32, reflect.Int16:
				fn = encFn{ &fi, (*encFnInfo).kInt }
			case reflect.Uint8, reflect.Uint64, reflect.Uint, reflect.Uint32, reflect.Uint16:
				fn = encFn{ &fi, (*encFnInfo).kUint }
			case reflect.Invalid:
				fn = encFn{ &fi, (*encFnInfo).kInvalid }
			case reflect.Slice:
				fn = encFn{ &fi, (*encFnInfo).kSlice }
			case reflect.Array:
				fn = encFn{ &fi, (*encFnInfo).kArray }
			case reflect.Struct:
				fn = encFn{ &fi, (*encFnInfo).kStruct }
			case reflect.Ptr:
				fn = encFn{ &fi, (*encFnInfo).kPtr }
			case reflect.Interface:
				fn = encFn{ &fi, (*encFnInfo).kInterface }
			case reflect.Map:
				fn = encFn{ &fi, (*encFnInfo).kMap }
			default:
				fn = encFn{ &fi, (*encFnInfo).kErr }
			}
		}
		e.f[rtid] = fn
	}
	
	fn.f(fn.i, rv)

}

// ----------------------------------------

func (z *ioEncWriter) writeUint16(v uint16) {
	bigen.PutUint16(z.x[:2], v)
	z.writeb(z.x[:2])
}

func (z *ioEncWriter) writeUint32(v uint32) {
	bigen.PutUint32(z.x[:4], v)
	z.writeb(z.x[:4])
}

func (z *ioEncWriter) writeUint64(v uint64) {
	bigen.PutUint64(z.x[:8], v)
	z.writeb(z.x[:8])
}

func (z *ioEncWriter) writeb(bs []byte) {
	n, err := z.w.Write(bs)
	if err != nil {
		panic(err)
	}
	if n != len(bs) {
		doPanic(msgTagEnc, "write: Incorrect num bytes written. Expecting: %v, Wrote: %v", len(bs), n)
	}
}

func (z *ioEncWriter) writestr(s string) {
	n, err := z.w.WriteString(s)
	if err != nil {
		panic(err)
	}
	if n != len(s) {
		doPanic(msgTagEnc, "write: Incorrect num bytes written. Expecting: %v, Wrote: %v", len(s), n)
	}
}

func (z *ioEncWriter) writen1(b byte) {
	if err := z.w.WriteByte(b); err != nil {
		panic(err)
	}
}

func (z *ioEncWriter) writen2(b1 byte, b2 byte) {
	z.writen1(b1)
	z.writen1(b2)
}

func (z *ioEncWriter) atEndOfEncode() { }

// ----------------------------------------

func (z *bytesEncWriter) writeUint16(v uint16) {
	c := z.grow(2)
	z.b[c] = byte(v >> 8)
	z.b[c+1] = byte(v)
}

func (z *bytesEncWriter) writeUint32(v uint32) {
	c := z.grow(4)
	z.b[c] = byte(v >> 24)
	z.b[c+1] = byte(v >> 16)
	z.b[c+2] = byte(v >> 8)
	z.b[c+3] = byte(v)
}

func (z *bytesEncWriter) writeUint64(v uint64) {
	c := z.grow(8)
	z.b[c] = byte(v >> 56)
	z.b[c+1] = byte(v >> 48)
	z.b[c+2] = byte(v >> 40)
	z.b[c+3] = byte(v >> 32)
	z.b[c+4] = byte(v >> 24)
	z.b[c+5] = byte(v >> 16)
	z.b[c+6] = byte(v >> 8)
	z.b[c+7] = byte(v)
}

func (z *bytesEncWriter) writeb(s []byte) {
	c := z.grow(len(s))
	copy(z.b[c:], s)
}

func (z *bytesEncWriter) writestr(s string) {
	c := z.grow(len(s))
	copy(z.b[c:], s)
}

func (z *bytesEncWriter) writen1(b1 byte) {
	c := z.grow(1)
	z.b[c] = b1
}

func (z *bytesEncWriter) writen2(b1 byte, b2 byte) {
	c := z.grow(2)
	z.b[c] = b1
	z.b[c+1] = b2
}

func (z *bytesEncWriter) atEndOfEncode() {
	*(z.out) = z.b[:z.c]
}

func (z *bytesEncWriter) grow(n int) (oldcursor int) {
	oldcursor = z.c
	z.c = oldcursor + n
	if z.c > cap(z.b) {
		// It tried using appendslice logic: (if cap < 1024, *2, else *1.25).
		// However, it was too expensive, causing too many iterations of copy.
		// Using bytes.Buffer model was much better (2*cap + n)
		bs := make([]byte, 2*cap(z.b)+n)
		copy(bs, z.b[:oldcursor])
		z.b = bs
	} else if z.c > len(z.b) {
		z.b = z.b[:cap(z.b)]
	}
	return
}

// ----------------------------------------

func encErr(format string, params ...interface{}) {
	doPanic(msgTagEnc, format, params...)
}

