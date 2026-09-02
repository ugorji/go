// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build notmono || codec.notmono

package codec

import (
	"io"
)

// // This contains all the iniatializations of generics.
// // Putting it into one file, ensures that we can go generics or not.

func callMake(v any) {
	// we couldn't find an easy way to initializa these without resorting to reflection.
	// however, a type switch does it elegantly, since we have a finite set of types to support.
	switch x := v.(type) {
	case **bytesEncAppender:
		*x = new(bytesEncAppender)
	case **bufioEncWriter:
		*x = new(bufioEncWriter)
	case **bytesDecReader:
		*x = new(bytesDecReader)
	case **ioDecReader:
		*x = new(ioDecReader)
	case **simpleEncDriver[*bufioEncWriter]:
		*x = new(simpleEncDriver[*bufioEncWriter])
	case **simpleEncDriver[*bytesEncAppender]:
		*x = new(simpleEncDriver[*bytesEncAppender])
	case **jsonEncDriver[*bufioEncWriter]:
		*x = new(jsonEncDriver[*bufioEncWriter])
	case **jsonEncDriver[*bytesEncAppender]:
		*x = new(jsonEncDriver[*bytesEncAppender])
	case **cborEncDriver[*bufioEncWriter]:
		*x = new(cborEncDriver[*bufioEncWriter])
	case **cborEncDriver[*bytesEncAppender]:
		*x = new(cborEncDriver[*bytesEncAppender])
	case **msgpackEncDriver[*bufioEncWriter]:
		*x = new(msgpackEncDriver[*bufioEncWriter])
	case **msgpackEncDriver[*bytesEncAppender]:
		*x = new(msgpackEncDriver[*bytesEncAppender])
	case **bincEncDriver[*bufioEncWriter]:
		*x = new(bincEncDriver[*bufioEncWriter])
	case **bincEncDriver[*bytesEncAppender]:
		*x = new(bincEncDriver[*bytesEncAppender])
	case **simpleDecDriver[*bytesDecReader]:
		*x = new(simpleDecDriver[*bytesDecReader])
	case **simpleDecDriver[*ioDecReader]:
		*x = new(simpleDecDriver[*ioDecReader])
	case **jsonDecDriver[*bytesDecReader]:
		*x = new(jsonDecDriver[*bytesDecReader])
	case **jsonDecDriver[*ioDecReader]:
		*x = new(jsonDecDriver[*ioDecReader])
	case **cborDecDriver[*bytesDecReader]:
		*x = new(cborDecDriver[*bytesDecReader])
	case **cborDecDriver[*ioDecReader]:
		*x = new(cborDecDriver[*ioDecReader])
	case **msgpackDecDriver[*bytesDecReader]:
		*x = new(msgpackDecDriver[*bytesDecReader])
	case **msgpackDecDriver[*ioDecReader]:
		*x = new(msgpackDecDriver[*ioDecReader])
	case **bincDecDriver[*bytesDecReader]:
		*x = new(bincDecDriver[*bytesDecReader])
	case **bincDecDriver[*ioDecReader]:
		*x = new(bincDecDriver[*ioDecReader])
	}
}

// ---- (writer.go)

type encWriter interface {
	*bufioEncWriter | *bytesEncAppender
	encWriterI
}

// ---- reader.go

type decReader interface {
	*bytesDecReader | *ioDecReader
	decReaderI
}

// type helperEncWriter[T encWriter] struct{}
// type helperDecReader[T decReader] struct{}
// func (helperDecReader[T]) decByteSlice(r T, clen, maxInitLen int, bs []byte) (bsOut []byte) {

// ---- (encode.go)

type encDriver interface {
	*simpleEncDriver[*bufioEncWriter] |
		*simpleEncDriver[*bytesEncAppender] |
		*jsonEncDriver[*bufioEncWriter] |
		*jsonEncDriver[*bytesEncAppender] |
		*cborEncDriver[*bufioEncWriter] |
		*cborEncDriver[*bytesEncAppender] |
		*msgpackEncDriver[*bufioEncWriter] |
		*msgpackEncDriver[*bytesEncAppender] |
		*bincEncDriver[*bufioEncWriter] |
		*bincEncDriver[*bytesEncAppender]

	encDriverI
}

// ---- (decode.go)

type decDriver interface {
	*simpleDecDriver[*bytesDecReader] |
		*simpleDecDriver[*ioDecReader] |
		*jsonDecDriver[*bytesDecReader] |
		*jsonDecDriver[*ioDecReader] |
		*cborDecDriver[*bytesDecReader] |
		*cborDecDriver[*ioDecReader] |
		*msgpackDecDriver[*bytesDecReader] |
		*msgpackDecDriver[*ioDecReader] |
		*bincDecDriver[*bytesDecReader] |
		*bincDecDriver[*ioDecReader]

	decDriverI
}

// Below: <format>.go files

// ---- (binc.go)

var (
	bincFpEncIO    = helperEncDriver[*bincEncDriver[*bufioEncWriter]]{}.fastpathEList()
	bincFpEncBytes = helperEncDriver[*bincEncDriver[*bytesEncAppender]]{}.fastpathEList()
	bincFpDecIO    = helperDecDriver[*bincDecDriver[*ioDecReader]]{}.fastpathDList()
	bincFpDecBytes = helperDecDriver[*bincDecDriver[*bytesDecReader]]{}.fastpathDList()
)

// ---- (cbor.go)

var (
	cborFpEncIO    = helperEncDriver[*cborEncDriver[*bufioEncWriter]]{}.fastpathEList()
	cborFpEncBytes = helperEncDriver[*cborEncDriver[*bytesEncAppender]]{}.fastpathEList()
	cborFpDecIO    = helperDecDriver[*cborDecDriver[*ioDecReader]]{}.fastpathDList()
	cborFpDecBytes = helperDecDriver[*cborDecDriver[*bytesDecReader]]{}.fastpathDList()
)

// ---- (json.go)

var (
	jsonFpEncIO    = helperEncDriver[*jsonEncDriver[*bufioEncWriter]]{}.fastpathEList()
	jsonFpEncBytes = helperEncDriver[*jsonEncDriver[*bytesEncAppender]]{}.fastpathEList()
	jsonFpDecIO    = helperDecDriver[*jsonDecDriver[*ioDecReader]]{}.fastpathDList()
	jsonFpDecBytes = helperDecDriver[*jsonDecDriver[*bytesDecReader]]{}.fastpathDList()
)

// ---- (msgpack.go)

var (
	msgpackFpEncIO    = helperEncDriver[*msgpackEncDriver[*bufioEncWriter]]{}.fastpathEList()
	msgpackFpEncBytes = helperEncDriver[*msgpackEncDriver[*bytesEncAppender]]{}.fastpathEList()
	msgpackFpDecIO    = helperDecDriver[*msgpackDecDriver[*ioDecReader]]{}.fastpathDList()
	msgpackFpDecBytes = helperDecDriver[*msgpackDecDriver[*bytesDecReader]]{}.fastpathDList()
)

// ---- (simple.go)

var (
	simpleFpEncIO    = helperEncDriver[*simpleEncDriver[*bufioEncWriter]]{}.fastpathEList()
	simpleFpEncBytes = helperEncDriver[*simpleEncDriver[*bytesEncAppender]]{}.fastpathEList()
	simpleFpDecIO    = helperDecDriver[*simpleDecDriver[*ioDecReader]]{}.fastpathDList()
	simpleFpDecBytes = helperDecDriver[*simpleDecDriver[*bytesDecReader]]{}.fastpathDList()
)

func (h *SimpleHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[*simpleEncDriver[*bytesEncAppender]]{}.newEncoderBytes(out, h)
}

func (h *SimpleHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[*simpleEncDriver[*bufioEncWriter]]{}.newEncoderIO(w, h)
}

func (h *SimpleHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[*simpleDecDriver[*bytesDecReader]]{}.newDecoderBytes(in, h)
}

func (h *SimpleHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[*simpleDecDriver[*ioDecReader]]{}.newDecoderIO(r, h)
}

func (h *JsonHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[*jsonEncDriver[*bytesEncAppender]]{}.newEncoderBytes(out, h)
}

func (h *JsonHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[*jsonEncDriver[*bufioEncWriter]]{}.newEncoderIO(w, h)
}

func (h *JsonHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[*jsonDecDriver[*bytesDecReader]]{}.newDecoderBytes(in, h)
}

func (h *JsonHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[*jsonDecDriver[*ioDecReader]]{}.newDecoderIO(r, h)
}

func (h *MsgpackHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[*msgpackEncDriver[*bytesEncAppender]]{}.newEncoderBytes(out, h)
}

func (h *MsgpackHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[*msgpackEncDriver[*bufioEncWriter]]{}.newEncoderIO(w, h)
}

func (h *MsgpackHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[*msgpackDecDriver[*bytesDecReader]]{}.newDecoderBytes(in, h)
}

func (h *MsgpackHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[*msgpackDecDriver[*ioDecReader]]{}.newDecoderIO(r, h)
}

func (h *CborHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[*cborEncDriver[*bytesEncAppender]]{}.newEncoderBytes(out, h)
}

func (h *CborHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[*cborEncDriver[*bufioEncWriter]]{}.newEncoderIO(w, h)
}

func (h *CborHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[*cborDecDriver[*bytesDecReader]]{}.newDecoderBytes(in, h)
}

func (h *CborHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[*cborDecDriver[*ioDecReader]]{}.newDecoderIO(r, h)
}

func (h *BincHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[*bincEncDriver[*bytesEncAppender]]{}.newEncoderBytes(out, h)
}

func (h *BincHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[*bincEncDriver[*bufioEncWriter]]{}.newEncoderIO(w, h)
}

func (h *BincHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[*bincDecDriver[*bytesDecReader]]{}.newDecoderBytes(in, h)
}

func (h *BincHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[*bincDecDriver[*ioDecReader]]{}.newDecoderIO(r, h)
}
