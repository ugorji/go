// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !notmono && !codec.notmono

package codec

import "io"

func callMake(v interface{}) {}

type encWriter interface{ encWriterI }
type decReader interface{ decReaderI }
type encDriver interface{ encDriverI }
type decDriver interface{ decDriverI }

func (h *SimpleHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverSimpleBytes{}.newEncoderBytes(out, h)
}

func (h *SimpleHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverSimpleIO{}.newEncoderIO(w, h)
}

func (h *SimpleHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverSimpleBytes{}.newDecoderBytes(in, h)
}

func (h *SimpleHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverSimpleIO{}.newDecoderIO(r, h)
}

func (h *JsonHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverJsonBytes{}.newEncoderBytes(out, h)
}

func (h *JsonHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverJsonIO{}.newEncoderIO(w, h)
}

func (h *JsonHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverJsonBytes{}.newDecoderBytes(in, h)
}

func (h *JsonHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverJsonIO{}.newDecoderIO(r, h)
}

func (h *MsgpackHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverMsgpackBytes{}.newEncoderBytes(out, h)
}

func (h *MsgpackHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverMsgpackIO{}.newEncoderIO(w, h)
}

func (h *MsgpackHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverMsgpackBytes{}.newDecoderBytes(in, h)
}

func (h *MsgpackHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverMsgpackIO{}.newDecoderIO(r, h)
}

func (h *BincHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverBincBytes{}.newEncoderBytes(out, h)
}

func (h *BincHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverBincIO{}.newEncoderIO(w, h)
}

func (h *BincHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverBincBytes{}.newDecoderBytes(in, h)
}

func (h *BincHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverBincIO{}.newDecoderIO(r, h)
}

func (h *CborHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverCborBytes{}.newEncoderBytes(out, h)
}

func (h *CborHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverCborIO{}.newEncoderIO(w, h)
}

func (h *CborHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverCborBytes{}.newDecoderBytes(in, h)
}

func (h *CborHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverCborIO{}.newDecoderIO(r, h)
}

var (
	bincFpEncIO    = helperEncDriverBincIO{}.fastpathEList()
	bincFpEncBytes = helperEncDriverBincBytes{}.fastpathEList()
	bincFpDecIO    = helperDecDriverBincIO{}.fastpathDList()
	bincFpDecBytes = helperDecDriverBincBytes{}.fastpathDList()
)

var (
	cborFpEncIO    = helperEncDriverCborIO{}.fastpathEList()
	cborFpEncBytes = helperEncDriverCborBytes{}.fastpathEList()
	cborFpDecIO    = helperDecDriverCborIO{}.fastpathDList()
	cborFpDecBytes = helperDecDriverCborBytes{}.fastpathDList()
)

var (
	jsonFpEncIO    = helperEncDriverJsonIO{}.fastpathEList()
	jsonFpEncBytes = helperEncDriverJsonBytes{}.fastpathEList()
	jsonFpDecIO    = helperDecDriverJsonIO{}.fastpathDList()
	jsonFpDecBytes = helperDecDriverJsonBytes{}.fastpathDList()
)

var (
	msgpackFpEncIO    = helperEncDriverMsgpackIO{}.fastpathEList()
	msgpackFpEncBytes = helperEncDriverMsgpackBytes{}.fastpathEList()
	msgpackFpDecIO    = helperDecDriverMsgpackIO{}.fastpathDList()
	msgpackFpDecBytes = helperDecDriverMsgpackBytes{}.fastpathDList()
)

var (
	simpleFpEncIO    = helperEncDriverSimpleIO{}.fastpathEList()
	simpleFpEncBytes = helperEncDriverSimpleBytes{}.fastpathEList()
	simpleFpDecIO    = helperDecDriverSimpleIO{}.fastpathDList()
	simpleFpDecBytes = helperDecDriverSimpleBytes{}.fastpathDList()
)

// // NewEncoder returns an Encoder for encoding into an io.Writer.
// //
// // For efficiency, Users are encouraged to configure WriterBufferSize on the handle
// // OR pass in a memory buffered writer (eg bufio.Writer, bytes.Buffer).
// func NewEncoder(w io.Writer, h Handle) *Encoder {
// 	var e encoderI
// 	switch h.(type) {
// 	case *SimpleHandle:
// 		var dh helperEncDriverSimpleIO
// 		e = dh.newEncoderIO(w, h)
// 	case *JsonHandle:
// 		var dh helperEncDriverJsonIO
// 		e = dh.newEncoderIO(w, h)
// 	case *CborHandle:
// 		var dh helperEncDriverCborIO
// 		e = dh.newEncoderIO(w, h)
// 	case *MsgpackHandle:
// 		var dh helperEncDriverMsgpackIO
// 		e = dh.newEncoderIO(w, h)
// 	case *BincHandle:
// 		var dh helperEncDriverBincIO
// 		e = dh.newEncoderIO(w, h)
// 	default:
// 		return nil
// 	}
// 	return &Encoder{e}
// }

// // NewEncoderBytes returns an encoder for encoding directly and efficiently
// // into a byte slice, using zero-copying to temporary slices.
// //
// // It will potentially replace the output byte slice pointed to.
// // After encoding, the out parameter contains the encoded contents.
// func NewEncoderBytes(out *[]byte, h Handle) *Encoder {
// 	var e encoderI
// 	switch h.(type) {
// 	case *SimpleHandle:
// 		var dh helperEncDriverSimpleBytes
// 		e = dh.newEncoderBytes(out, h)
// 	case *JsonHandle:
// 		var dh helperEncDriverJsonBytes
// 		e = dh.newEncoderBytes(out, h)
// 	case *CborHandle:
// 		var dh helperEncDriverCborBytes
// 		e = dh.newEncoderBytes(out, h)
// 	case *MsgpackHandle:
// 		var dh helperEncDriverMsgpackBytes
// 		e = dh.newEncoderBytes(out, h)
// 	case *BincHandle:
// 		var dh helperEncDriverBincBytes
// 		e = dh.newEncoderBytes(out, h)
// 	default:
// 		return nil
// 	}
// 	return &Encoder{e}
// }

// // NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
// //
// // For efficiency, Users are encouraged to configure ReaderBufferSize on the handle
// // OR pass in a memory buffered reader (eg bufio.Reader, bytes.Buffer).
// func NewDecoder(r io.Reader, h Handle) *Decoder {
// 	var d decoderI
// 	switch h.(type) {
// 	case *SimpleHandle:
// 		var dh helperDecDriverSimpleIO
// 		d = dh.newDecoderIO(r, h)
// 	case *JsonHandle:
// 		var dh helperDecDriverJsonIO
// 		d = dh.newDecoderIO(r, h)
// 	case *CborHandle:
// 		var dh helperDecDriverCborIO
// 		d = dh.newDecoderIO(r, h)
// 	case *MsgpackHandle:
// 		var dh helperDecDriverMsgpackIO
// 		d = dh.newDecoderIO(r, h)
// 	case *BincHandle:
// 		var dh helperDecDriverBincIO
// 		d = dh.newDecoderIO(r, h)
// 	default:
// 		return nil
// 	}
// 	return &Decoder{d}
// }

// // NewDecoderBytes returns a Decoder which efficiently decodes directly
// // from a byte slice with zero copying.
// func NewDecoderBytes(in []byte, h Handle) *Decoder {
// 	var d decoderI
// 	switch h.(type) {
// 	case *SimpleHandle:
// 		var dh helperDecDriverSimpleBytes
// 		d = dh.newDecoderBytes(in, h)
// 	case *JsonHandle:
// 		var dh helperDecDriverJsonBytes
// 		d = dh.newDecoderBytes(in, h)
// 	case *CborHandle:
// 		var dh helperDecDriverCborBytes
// 		d = dh.newDecoderBytes(in, h)
// 	case *MsgpackHandle:
// 		var dh helperDecDriverMsgpackBytes
// 		d = dh.newDecoderBytes(in, h)
// 	case *BincHandle:
// 		var dh helperDecDriverBincBytes
// 		d = dh.newDecoderBytes(in, h)
// 	default:
// 		return nil
// 	}
// 	return &Decoder{d}
// }
