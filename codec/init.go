// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
	"reflect"
)

// This contains all the iniatializations of generics.
// Putting it into one file, ensures that we can go generics or not.

// ---- (writer.go)

type encWriter interface {
	bufioEncWriterM | bytesEncAppenderM
	encWriterI
}

// ---- reader.go

type decReader interface {
	bytesDecReaderM | ioDecReaderM

	decReaderI
}

// ---- (encode.go)

type encDriver interface {
	simpleEncDriverM[bufioEncWriterM] |
		simpleEncDriverM[bytesEncAppenderM] |
		jsonEncDriverM[bufioEncWriterM] |
		jsonEncDriverM[bytesEncAppenderM] |
		cborEncDriverM[bufioEncWriterM] |
		cborEncDriverM[bytesEncAppenderM] |
		msgpackEncDriverM[bufioEncWriterM] |
		msgpackEncDriverM[bytesEncAppenderM] |
		bincEncDriverM[bufioEncWriterM] |
		bincEncDriverM[bytesEncAppenderM]

	encDriverI
}

// NewEncoder returns an Encoder for encoding into an io.Writer.
//
// For efficiency, Users are encouraged to configure WriterBufferSize on the handle
// OR pass in a memory buffered writer (eg bufio.Writer, bytes.Buffer).
func NewEncoder(w io.Writer, h Handle) *Encoder {
	var e encoderI
	switch h.(type) {
	case *SimpleHandle:
		var dh helperEncDriver[simpleEncDriverM[bufioEncWriterM]]
		e = dh.newEncDriverIO(w, h)
	case *JsonHandle:
		var dh helperEncDriver[jsonEncDriverM[bufioEncWriterM]]
		e = dh.newEncDriverIO(w, h)
	case *CborHandle:
		var dh helperEncDriver[cborEncDriverM[bufioEncWriterM]]
		e = dh.newEncDriverIO(w, h)
	case *MsgpackHandle:
		var dh helperEncDriver[msgpackEncDriverM[bufioEncWriterM]]
		e = dh.newEncDriverIO(w, h)
	case *BincHandle:
		var dh helperEncDriver[bincEncDriverM[bufioEncWriterM]]
		e = dh.newEncDriverIO(w, h)
	default:
		return nil
	}
	return &Encoder{e}
}

// NewEncoderBytes returns an encoder for encoding directly and efficiently
// into a byte slice, using zero-copying to temporary slices.
//
// It will potentially replace the output byte slice pointed to.
// After encoding, the out parameter contains the encoded contents.
func NewEncoderBytes(out *[]byte, h Handle) *Encoder {
	var e encoderI
	switch h.(type) {
	case *SimpleHandle:
		var dh helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]
		e = dh.newEncDriverBytes(out, h)
	case *JsonHandle:
		var dh helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]
		e = dh.newEncDriverBytes(out, h)
	case *CborHandle:
		var dh helperEncDriver[cborEncDriverM[bytesEncAppenderM]]
		e = dh.newEncDriverBytes(out, h)
	case *MsgpackHandle:
		var dh helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]
		e = dh.newEncDriverBytes(out, h)
	case *BincHandle:
		var dh helperEncDriver[bincEncDriverM[bytesEncAppenderM]]
		e = dh.newEncDriverBytes(out, h)
	default:
		return nil
	}
	return &Encoder{e}
}

// ---- (decode.go)

type decDriver interface {
	simpleDecDriverM[bytesDecReaderM] |
		simpleDecDriverM[ioDecReaderM] |
		jsonDecDriverM[bytesDecReaderM] |
		jsonDecDriverM[ioDecReaderM] |
		cborDecDriverM[bytesDecReaderM] |
		cborDecDriverM[ioDecReaderM] |
		msgpackDecDriverM[bytesDecReaderM] |
		msgpackDecDriverM[ioDecReaderM] |
		bincDecDriverM[bytesDecReaderM] |
		bincDecDriverM[ioDecReaderM]

	decDriverI
}

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to configure ReaderBufferSize on the handle
// OR pass in a memory buffered reader (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder {
	var d decoderI
	switch h.(type) {
	case *SimpleHandle:
		var dh helperDecDriver[simpleDecDriverM[ioDecReaderM]]
		d = dh.newDecDriverIO(r, h)
	case *JsonHandle:
		var dh helperDecDriver[jsonDecDriverM[ioDecReaderM]]
		d = dh.newDecDriverIO(r, h)
	case *CborHandle:
		var dh helperDecDriver[cborDecDriverM[ioDecReaderM]]
		d = dh.newDecDriverIO(r, h)
	case *MsgpackHandle:
		var dh helperDecDriver[msgpackDecDriverM[ioDecReaderM]]
		d = dh.newDecDriverIO(r, h)
	case *BincHandle:
		var dh helperDecDriver[bincDecDriverM[ioDecReaderM]]
		d = dh.newDecDriverIO(r, h)
	default:
		return nil
	}
	return &Decoder{d}
}

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder {
	var d decoderI
	switch h.(type) {
	case *SimpleHandle:
		var dh helperDecDriver[simpleDecDriverM[bytesDecReaderM]]
		d = dh.newDecDriverBytes(in, h)
	case *JsonHandle:
		var dh helperDecDriver[jsonDecDriverM[bytesDecReaderM]]
		d = dh.newDecDriverBytes(in, h)
	case *CborHandle:
		var dh helperDecDriver[cborDecDriverM[bytesDecReaderM]]
		d = dh.newDecDriverBytes(in, h)
	case *MsgpackHandle:
		var dh helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]
		d = dh.newDecDriverBytes(in, h)
	case *BincHandle:
		var dh helperDecDriver[bincDecDriverM[bytesDecReaderM]]
		d = dh.newDecDriverBytes(in, h)
	default:
		return nil
	}
	return &Decoder{d}
}

// Below: <format>.go files

// ---- (binc.go)

var (
	bincFpEncIO    = helperEncDriver[bincEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	bincFpEncBytes = helperEncDriver[bincEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	bincFpDecIO    = helperDecDriver[bincDecDriverM[ioDecReaderM]]{}.fastpathDList()
	bincFpDecBytes = helperDecDriver[bincDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

func (e *bincEncDriver[T]) sideEncoder(out *[]byte) {
	var dh helperEncDriver[bincEncDriverM[bytesEncAppenderM]]
	dh.sideEncoder(out, e.e, e.h)
}

func (e *bincEncDriver[T]) sideEncode(v interface{}, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[bincEncDriverM[bytesEncAppenderM]]
	dh.sideEncode(e.e.se.(*encoder[bincEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (e *bincEncDriver[T]) sideEncodeRV(v reflect.Value, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[bincEncDriverM[bytesEncAppenderM]]
	dh.sideEncodeRV(e.e.se.(*encoder[bincEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (d *bincDecDriver[T]) sideDecoder(in []byte) {
	var dh helperDecDriver[bincDecDriverM[bytesDecReaderM]]
	dh.sideDecoder(in, d.d, d.h)
}

func (d *bincDecDriver[T]) sideDecode(v interface{}, basetype reflect.Type) {
	var dh helperDecDriver[bincDecDriverM[bytesDecReaderM]]
	dh.sideDecode(d.d.sd.(*decoder[bincDecDriverM[bytesDecReaderM]]), v, basetype)
}

func (d *bincDecDriver[T]) sideDecodeRV(v reflect.Value, basetype reflect.Type) {
	var dh helperDecDriver[bincDecDriverM[bytesDecReaderM]]
	dh.sideDecodeRV(d.d.sd.(*decoder[bincDecDriverM[bytesDecReaderM]]), v, basetype)
}

// ---- (cbor.go)

var (
	cborFpEncIO    = helperEncDriver[cborEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	cborFpEncBytes = helperEncDriver[cborEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	cborFpDecIO    = helperDecDriver[cborDecDriverM[ioDecReaderM]]{}.fastpathDList()
	cborFpDecBytes = helperDecDriver[cborDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

func (e *cborEncDriver[T]) sideEncoder(out *[]byte) {
	var dh helperEncDriver[cborEncDriverM[bytesEncAppenderM]]
	dh.sideEncoder(out, e.e, e.h)
}

func (e *cborEncDriver[T]) sideEncode(v interface{}, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[cborEncDriverM[bytesEncAppenderM]]
	dh.sideEncode(e.e.se.(*encoder[cborEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (e *cborEncDriver[T]) sideEncodeRV(v reflect.Value, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[cborEncDriverM[bytesEncAppenderM]]
	dh.sideEncodeRV(e.e.se.(*encoder[cborEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (d *cborDecDriver[T]) sideDecoder(in []byte) {
	var dh helperDecDriver[cborDecDriverM[bytesDecReaderM]]
	dh.sideDecoder(in, d.d, d.h)
}

func (d *cborDecDriver[T]) sideDecode(v interface{}, basetype reflect.Type) {
	var dh helperDecDriver[cborDecDriverM[bytesDecReaderM]]
	dh.sideDecode(d.d.sd.(*decoder[cborDecDriverM[bytesDecReaderM]]), v, basetype)
}

func (d *cborDecDriver[T]) sideDecodeRV(v reflect.Value, basetype reflect.Type) {
	var dh helperDecDriver[cborDecDriverM[bytesDecReaderM]]
	dh.sideDecodeRV(d.d.sd.(*decoder[cborDecDriverM[bytesDecReaderM]]), v, basetype)
}

// ---- (json.go)

var (
	jsonFpEncIO    = helperEncDriver[jsonEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	jsonFpEncBytes = helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	jsonFpDecIO    = helperDecDriver[jsonDecDriverM[ioDecReaderM]]{}.fastpathDList()
	jsonFpDecBytes = helperDecDriver[jsonDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

func (e *jsonEncDriver[T]) sideEncoder(out *[]byte) {
	var dh helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]
	dh.sideEncoder(out, e.e, e.h)
}

func (e *jsonEncDriver[T]) sideEncode(v interface{}, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]
	dh.sideEncode(e.e.se.(*encoder[jsonEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (e *jsonEncDriver[T]) sideEncodeRV(v reflect.Value, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]
	dh.sideEncodeRV(e.e.se.(*encoder[jsonEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (d *jsonDecDriver[T]) sideDecoder(in []byte) {
	var dh helperDecDriver[jsonDecDriverM[bytesDecReaderM]]
	dh.sideDecoder(in, d.d, d.h)
}

func (d *jsonDecDriver[T]) sideDecode(v interface{}, basetype reflect.Type) {
	var dh helperDecDriver[jsonDecDriverM[bytesDecReaderM]]
	dh.sideDecode(d.d.sd.(*decoder[jsonDecDriverM[bytesDecReaderM]]), v, basetype)
}

func (d *jsonDecDriver[T]) sideDecodeRV(v reflect.Value, basetype reflect.Type) {
	var dh helperDecDriver[jsonDecDriverM[bytesDecReaderM]]
	dh.sideDecodeRV(d.d.sd.(*decoder[jsonDecDriverM[bytesDecReaderM]]), v, basetype)
}

// ---- (msgpack.go)

var (
	msgpackFpEncIO    = helperEncDriver[msgpackEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	msgpackFpEncBytes = helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	msgpackFpDecIO    = helperDecDriver[msgpackDecDriverM[ioDecReaderM]]{}.fastpathDList()
	msgpackFpDecBytes = helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

func (e *msgpackEncDriver[T]) sideEncoder(out *[]byte) {
	var dh helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]
	dh.sideEncoder(out, e.e, e.h)
}

func (e *msgpackEncDriver[T]) sideEncode(v interface{}, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]
	dh.sideEncode(e.e.se.(*encoder[msgpackEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (e *msgpackEncDriver[T]) sideEncodeRV(v reflect.Value, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]
	dh.sideEncodeRV(e.e.se.(*encoder[msgpackEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (d *msgpackDecDriver[T]) sideDecoder(in []byte) {
	var dh helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]
	dh.sideDecoder(in, d.d, d.h)
}

func (d *msgpackDecDriver[T]) sideDecode(v interface{}, basetype reflect.Type) {
	var dh helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]
	dh.sideDecode(d.d.sd.(*decoder[msgpackDecDriverM[bytesDecReaderM]]), v, basetype)
}

func (d *msgpackDecDriver[T]) sideDecodeRV(v reflect.Value, basetype reflect.Type) {
	var dh helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]
	dh.sideDecodeRV(d.d.sd.(*decoder[msgpackDecDriverM[bytesDecReaderM]]), v, basetype)
}

// ---- (simple.go)

var (
	simpleFpEncIO    = helperEncDriver[simpleEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	simpleFpEncBytes = helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	simpleFpDecIO    = helperDecDriver[simpleDecDriverM[ioDecReaderM]]{}.fastpathDList()
	simpleFpDecBytes = helperDecDriver[simpleDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

func (e *simpleEncDriver[T]) sideEncoder(out *[]byte) {
	var dh helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]
	dh.sideEncoder(out, e.e, e.h)
}

func (e *simpleEncDriver[T]) sideEncode(v any, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]
	dh.sideEncode(e.e.se.(*encoder[simpleEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (e *simpleEncDriver[T]) sideEncodeRV(v reflect.Value, basetype reflect.Type, cs containerState) {
	var dh helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]
	dh.sideEncodeRV(e.e.se.(*encoder[simpleEncDriverM[bytesEncAppenderM]]), v, basetype, cs)
}

func (d *simpleDecDriver[T]) sideDecoder(in []byte) {
	var dh helperDecDriver[simpleDecDriverM[bytesDecReaderM]]
	dh.sideDecoder(in, d.d, d.h)
}

func (d *simpleDecDriver[T]) sideDecode(v any, basetype reflect.Type) {
	var dh helperDecDriver[simpleDecDriverM[bytesDecReaderM]]
	dh.sideDecode(d.d.sd.(*decoder[simpleDecDriverM[bytesDecReaderM]]), v, basetype)
}

func (d *simpleDecDriver[T]) sideDecodeRV(v reflect.Value, basetype reflect.Type) {
	var dh helperDecDriver[simpleDecDriverM[bytesDecReaderM]]
	dh.sideDecodeRV(d.d.sd.(*decoder[simpleDecDriverM[bytesDecReaderM]]), v, basetype)
}

// ---- commented out stuff

// func encResetBytes[T encWriter](w T, out *[]byte) (ok bool) {
// 	v, ok := any(w).(bytesEncAppenderM)
// 	if ok {
// 		v.resetBytes(*out, out)
// 	}
// 	return
// }

// func encResetIO[T encWriter](w T, out io.Writer, bufsize int, blist *bytesFreelist) (ok bool) {
// 	v, ok := any(w).(bufioEncWriterM)
// 	if ok {
// 		v.resetIO(out, bufsize, blist)
// 	}
// 	return
// }

// func decResetBytes[T decReader](r T, in []byte) (ok bool) {
// 	v, ok := any(r).(bytesDecReaderM)
// 	if ok {
// 		v.resetBytes(in)
// 	}
// 	return
// }

// func decResetIO[T decReader](r T, in io.Reader, bufsize int, blist *bytesFreelist) (ok bool) {
// 	v, ok := any(r).(ioDecReaderM)
// 	if ok {
// 		v.resetIO(in, bufsize, blist)
// 	}
// 	return
// }
