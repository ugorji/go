// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
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

// type helperEncWriter[T encWriter] struct{}
// type helperDecReader[T decReader] struct{}
// func (helperDecReader[T]) decByteSlice(r T, clen, maxInitLen int, bs []byte) (bsOut []byte) {

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

type bincEncDriverM[T encWriter] struct {
	*bincEncDriver[T]
}

func (d *bincEncDriverM[T]) Make() {
	d.bincEncDriver = new(bincEncDriver[T])
}

type bincDecDriverM[T decReader] struct {
	*bincDecDriver[T]
}

func (d *bincDecDriverM[T]) Make() {
	d.bincDecDriver = new(bincDecDriver[T])
}

var (
	bincFpEncIO    = helperEncDriver[bincEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	bincFpEncBytes = helperEncDriver[bincEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	bincFpDecIO    = helperDecDriver[bincDecDriverM[ioDecReaderM]]{}.fastpathDList()
	bincFpDecBytes = helperDecDriver[bincDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (cbor.go)

type cborEncDriverM[T encWriter] struct {
	*cborEncDriver[T]
}

func (d *cborEncDriverM[T]) Make() {
	d.cborEncDriver = new(cborEncDriver[T])
}

type cborDecDriverM[T decReader] struct {
	*cborDecDriver[T]
}

func (d *cborDecDriverM[T]) Make() {
	d.cborDecDriver = new(cborDecDriver[T])
}

var (
	cborFpEncIO    = helperEncDriver[cborEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	cborFpEncBytes = helperEncDriver[cborEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	cborFpDecIO    = helperDecDriver[cborDecDriverM[ioDecReaderM]]{}.fastpathDList()
	cborFpDecBytes = helperDecDriver[cborDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (json.go)

type jsonEncDriverM[T encWriter] struct {
	*jsonEncDriver[T]
}

func (d *jsonEncDriverM[T]) Make() {
	d.jsonEncDriver = new(jsonEncDriver[T])
}

type jsonDecDriverM[T decReader] struct {
	*jsonDecDriver[T]
}

func (d *jsonDecDriverM[T]) Make() {
	d.jsonDecDriver = new(jsonDecDriver[T])
}

var (
	jsonFpEncIO    = helperEncDriver[jsonEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	jsonFpEncBytes = helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	jsonFpDecIO    = helperDecDriver[jsonDecDriverM[ioDecReaderM]]{}.fastpathDList()
	jsonFpDecBytes = helperDecDriver[jsonDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (msgpack.go)

type msgpackEncDriverM[T encWriter] struct {
	*msgpackEncDriver[T]
}

func (d *msgpackEncDriverM[T]) Make() {
	d.msgpackEncDriver = new(msgpackEncDriver[T])
}

type msgpackDecDriverM[T decReader] struct {
	*msgpackDecDriver[T]
}

func (d *msgpackDecDriverM[T]) Make() {
	d.msgpackDecDriver = new(msgpackDecDriver[T])
}

var (
	msgpackFpEncIO    = helperEncDriver[msgpackEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	msgpackFpEncBytes = helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	msgpackFpDecIO    = helperDecDriver[msgpackDecDriverM[ioDecReaderM]]{}.fastpathDList()
	msgpackFpDecBytes = helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (simple.go)

type simpleEncDriverM[T encWriter] struct {
	*simpleEncDriver[T]
}

func (d *simpleEncDriverM[T]) Make() {
	d.simpleEncDriver = new(simpleEncDriver[T])
}

type simpleDecDriverM[T decReader] struct {
	*simpleDecDriver[T]
}

func (d *simpleDecDriverM[T]) Make() {
	d.simpleDecDriver = new(simpleDecDriver[T])
}

var (
	simpleFpEncIO    = helperEncDriver[simpleEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	simpleFpEncBytes = helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	simpleFpDecIO    = helperDecDriver[simpleDecDriverM[ioDecReaderM]]{}.fastpathDList()
	simpleFpDecBytes = helperDecDriver[simpleDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

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
