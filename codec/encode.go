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

type encDriver interface {
	encodeBuiltinType(rt uintptr, rv reflect.Value) bool
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

type encodeHandleI interface {
	getEncodeExt(rt uintptr) (tag byte, fn func(reflect.Value) ([]byte, error))
	writeExt() bool
}

// An Encoder writes an object to an output stream in the codec format.
type Encoder struct {
	w encWriter
	e encDriver
	h encodeHandleI
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

type encExtTagFn struct {
	fn  func(reflect.Value) ([]byte, error)
	tag byte
}

type encExtTypeTagFn struct {
	rt uintptr
	encExtTagFn
}

// EncoderOptions contain options for the encoder, e.g. registered extension functions.
type encHandle struct {
	extFuncs map[uintptr]encExtTagFn
	exts     []encExtTypeTagFn
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


// addEncodeExt registers a function to handle encoding a given type as an extension
// with a specific specific tag byte.
// To remove an extension, pass fn=nil.
func (o *encHandle) addEncodeExt(rt uintptr, tag byte, fn func(reflect.Value) ([]byte, error)) {
	if o.exts == nil {
		o.exts = make([]encExtTypeTagFn, 0, 8)
		o.extFuncs = make(map[uintptr]encExtTagFn, 8)
	} else {
		if _, ok := o.extFuncs[rt]; ok {
			delete(o.extFuncs, rt)
			for i := 0; i < len(o.exts); i++ {
				if o.exts[i].rt == rt {
					o.exts = append(o.exts[:i], o.exts[i+1:]...)
					break
				}
			}
		}
	}
	if fn != nil {
		o.extFuncs[rt] = encExtTagFn{fn, tag}
		o.exts = append(o.exts, encExtTypeTagFn{rt, encExtTagFn{fn, tag}})
	}
}

func (o *encHandle) getEncodeExt(rt uintptr) (tag byte, fn func(reflect.Value) ([]byte, error)) {
	if l := len(o.exts); l == 0 {
		return
	} else if l < mapAccessThreshold {
		for i := 0; i < l; i++ {
			if o.exts[i].rt == rt {
				x := o.exts[i].encExtTagFn
				return x.tag, x.fn
			}
		}
	} else {
		x := o.extFuncs[rt]
		return x.tag, x.fn
	}
	return
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
	return &Encoder{w: &z, h: h, e: h.newEncDriver(&z)}
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
	return &Encoder{w: &z, h: h, e: h.newEncDriver(&z)}
}

// Encode writes an object into a stream in the codec format.
//
// Struct values encode as maps. Each exported struct field is encoded unless:
//    - the field's tag is "-", or
//    - the field is empty and its tag specifies the "omitempty" option.
//
// The empty values are false, 0, any nil pointer or interface value,
// and any array, slice, map, or string of length zero.
//
// Anonymous fields are encoded inline if no struct tag is present.
// Else they are encoded as regular fields.
//
// The object's default key string is the struct field name but can be
// specified in the struct field's tag value.
// The "codec" key in struct field's tag value is the key name,
// followed by an optional comma and options.
//
// To set an option on all fields (e.g. omitempty on all fields), you
// can create a field called _struct, and set flags on it.
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
// Note:
//   - Encode will treat struct field names and keys in map[string]XXX as symbols.
//     Some formats support symbols (e.g. binc) and will properly encode the string
//     only once in the stream, and use a tag to refer to it thereafter.
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
	
	//encode based on type first, since over-rides are based on type.
	ee := e.e //don't dereference everytime
	if ee.encodeBuiltinType(rtid, rv) {
		return
	}

	//Note: tagFn must handle returning nil if value should be encoded as a nil.
	if xfTag, xfFn := e.h.getEncodeExt(rtid); xfFn != nil {
		bs, fnerr := xfFn(rv)
		if fnerr != nil {
			panic(fnerr)
		}
		if bs == nil {
			ee.encodeNil()
			return
		}
		if e.h.writeExt() {
			ee.encodeExtPreamble(xfTag, len(bs))
			e.w.writeb(bs)
		} else {
			ee.encodeStringBytes(c_RAW, bs)
		}
		return
	}
	
	// TODO: Encode if type is an encoding.BinaryMarshaler: MarshalBinary() (data []byte, err error)
	// There is a cost, as we need to change the rv to an interface{} first.

	// ensure more common cases appear early in switch.
	rk := rv.Kind()
	switch rk {
	case reflect.Bool:
		ee.encodeBool(rv.Bool())
	case reflect.String:
		ee.encodeString(c_UTF8, rv.String())
	case reflect.Float64:
		ee.encodeFloat64(rv.Float())
	case reflect.Float32:
		ee.encodeFloat32(float32(rv.Float()))
	case reflect.Slice:
		if rv.IsNil() {
			ee.encodeNil()
			break
		}
		if rt == byteSliceTyp {
			ee.encodeStringBytes(c_RAW, rv.Bytes())
			break
		}
		l := rv.Len()
		ee.encodeArrayPreamble(l)
		if l == 0 {
			break
		}
		for j := 0; j < l; j++ {
			e.encodeValue(rv.Index(j))
		}
	case reflect.Array:
		e.encodeValue(rv.Slice(0, rv.Len()))
	case reflect.Map:
		if rv.IsNil() {
			ee.encodeNil()
			break
		}
		l := rv.Len()
		ee.encodeMapPreamble(l)
		if l == 0 {
			break
		}
		keyTypeIsString := rt.Key().Kind() == reflect.String
		mks := rv.MapKeys()
		// for j, lmks := 0, len(mks); j < lmks; j++ {
		for j := range mks {
			if keyTypeIsString {
				ee.encodeSymbol(mks[j].String())
			} else {
				e.encodeValue(mks[j])
			}
			e.encodeValue(rv.MapIndex(mks[j]))
		}
	case reflect.Struct:
		sis := getStructFieldInfos(rtid, rt)
		e.encStruct(sis, rv)
	case reflect.Ptr:
		if rv.IsNil() {
			ee.encodeNil()
			break
		}
		e.encodeValue(rv.Elem())
	case reflect.Interface:
		if rv.IsNil() {
			ee.encodeNil()
			break
		}
		e.encodeValue(rv.Elem())
	case reflect.Int, reflect.Int8, reflect.Int64, reflect.Int32, reflect.Int16:
		ee.encodeInt(rv.Int())
	case reflect.Uint8, reflect.Uint64, reflect.Uint, reflect.Uint32, reflect.Uint16:
		ee.encodeUint(rv.Uint())
	case reflect.Invalid:
		ee.encodeNil()
	default:
		encErr("Unsupported kind: %s, for: %#v", rk, rv)
	}
	return
}

func (e *Encoder) encStruct(sis structFieldInfos, rv reflect.Value) {
	newlen := len(sis)
	rvals := make([]reflect.Value, newlen)
	encnames := make([]string, newlen)
	newlen = 0
	// var rv0 reflect.Value
	// for i := 0; i < l; i++ {
	// 	si := sis[i]
	for _, si := range sis {
		if si.i > -1 {
			rvals[newlen] = rv.Field(int(si.i))
		} else {
			rvals[newlen] = rv.FieldByIndex(si.is)
		}
		if si.omitEmpty && isEmptyValue(rvals[newlen]) {
			continue
		}
		// sivals[newlen] = i
		encnames[newlen] = si.encName
		newlen++
	}
	ee := e.e //don't dereference everytime
	ee.encodeMapPreamble(newlen)
	for j := 0; j < newlen; j++ {
		//e.encString(sis[sivals[j]].encName)
		ee.encodeSymbol(encnames[j])
		e.encodeValue(rvals[j])
	}
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

