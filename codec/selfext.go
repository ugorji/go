// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "sync"

// This extension expects that types registered with it implement SelfExt interface.
//
// This is used by libraries that support BytesExt e.g. cbor, json.
var GlobalSelfInterfaceExt InterfaceExt = selfInterfaceExt{}

// var selfExtEncPool = sync.Pool{
// 	New: func() interface{} { return new(Encoder) },
// }
// var selfExtDecPool = sync.Pool{
// 	New: func() interface{} { return new(Decoder) },
// }

// SelfExt is the interface that users implement so that their types
// can, with minimal effort, be able to be an extension while allowing the
// library handling the encoding/decoding needs easily.
//
// We now support the ability for an extension to define a tag,
// but allow itself to be encoded in a default way.
//
// Because we use the type to determine how to encode or decode a value,
// we cannot tell a value to just encode itself, as that will lead to an
// infinite recursion.
//
// Instead, a value can define how it is to be encoded using a delegate value.
//
// If your types implement SelfExt, you can use this to define an extension.
//
// At encode time, the interface will call CodecConvertExt, and encode that value.
// At decode time, the library will call CodecConvertExt, decode into that value, and
// call the primary object's CodecUpdateExt with that value.
//
// The easiest way to do this is via struct embedding:
//
//    type T struct {
//        tHelper
//    }
//    func (t *T) CodecConvertExt() { return &t.tHelper }
//    func (t *T) CodecUpdateExt(interface{}) {  } // no-op (as our delegate is interior pointer)
//    type tHelper struct {
//        // ... all t fields
//    }
//
// Usage model:
//
//    cborHandle.SetInterfaceExt(reflect.TypeOf(T), 122, codec.GlobalSelfInterfaceExt)
//
//    msgpackHandle.SetBytesExt(reflect.TypeOf(T), 122, codec.NewSelfBytesExt(msgpackHandle, 1024))
//
type SelfExt interface {
	CodecConvertExt() interface{}
	CodecUpdateExt(interface{})
}

type selfBytesExt struct {
	// For performance and memory utilization, use sync.Pools
	// They all have to be local to the Ext, as the Ext is bound to a Handle.
	p sync.Pool // pool of byte buffers
	e sync.Pool
	d sync.Pool
	h Handle
	// bufcap int     // cap for each byte buffer created
}

type selfInterfaceExt struct{}

// NewSelfBytesExt will return a BytesExt implementation,
// that will call an encoder to encode the value to a stream
// so it can be placed into the encoder stream, and use a decoder
// to do the same on the other end.
//
// Users can specify a buffer size, and we will initialize that
// buffer for encoding the type. This allows users manage
// how big the buffer is based on their knowledge of the type being
// registered.
//
// This extension expects that types registered with it implement SelfExt interface.
//
// This is used by libraries that support BytesExt e.g. msgpack, binc.
func NewSelfBytesExt(h Handle, bufcap int) *selfBytesExt {
	var v = selfBytesExt{h: h}
	v.p.New = func() interface{} {
		return make([]byte, 0, bufcap)
	}
	v.e.New = func() interface{} { return NewEncoderBytes(nil, v.h) }
	v.d.New = func() interface{} { return NewDecoderBytes(nil, v.h) }
	return &v
}

func (x *selfBytesExt) WriteExt(v interface{}) (s []byte) {
	ei := x.e.Get()
	bi := x.p.Get()
	defer func() {
		x.e.Put(ei)
		x.p.Put(bi)
	}()
	b := (bi.([]byte))[:0]
	e := ei.(*Encoder)
	e.ResetBytes(&b)
	e.MustEncode(v.(SelfExt).CodecConvertExt())
	if len(b) > 0 {
		s = make([]byte, len(b))
		copy(s, b)
	}
	return
}

func (x *selfBytesExt) ReadExt(dst interface{}, src []byte) {
	di := x.d.Get()
	d := di.(*Decoder)
	defer func() {
		d.Release()
		x.d.Put(di)
	}()
	d.ResetBytes(src)
	v := dst.(SelfExt).CodecConvertExt()
	d.MustDecode(v)
	dst.(SelfExt).CodecUpdateExt(v)
	return
}

func (x selfInterfaceExt) ConvertExt(v interface{}) interface{} {
	return v.(SelfExt).CodecConvertExt()
}

func (x selfInterfaceExt) UpdateExt(dst interface{}, src interface{}) {
	dst.(SelfExt).CodecUpdateExt(src)
}
