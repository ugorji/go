// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	// using codec.XXX directly
	// . "github.com/ugorji/go/codec"
)

// variables that are not flags, but which can configure the handles
var (
	testEncodeOptions EncodeOptions
	testDecodeOptions DecodeOptions
	testRPCOptions    RPCOptions
)

type testHED struct {
	H   Handle
	Eb  *Encoder
	Db  *Decoder
	Eio *Encoder
	Dio *Decoder
}

var (
	// testNoopH    = NoopHandle(8)
	testBincH    = &BincHandle{}
	testMsgpackH = &MsgpackHandle{}
	testJsonH    = &JsonHandle{}
	testSimpleH  = &SimpleHandle{}
	testCborH    = &CborHandle{}

	testHandles []Handle

	testHEDs []testHED
)

func init() {
	// testHEDs = make([]testHED, 0, 32)

	// JSON should do HTMLCharsAsIs by default
	testJsonH.HTMLCharsAsIs = true
	// testJsonH.InternString = true

	testHandles = append(testHandles, testSimpleH, testJsonH, testCborH, testMsgpackH, testBincH)

	testReInitFns = append(testReInitFns, func() { testHEDs = nil })
}

func testHEDGet(h Handle) (d *testHED) {
	for i := range testHEDs {
		v := &testHEDs[i]
		if v.H == h {
			// fmt.Printf("testHEDGet: <cached> h: %v\n", h)
			return v
		}
	}
	testHEDs = append(testHEDs, testHED{
		H:   h,
		Eio: NewEncoder(nil, h),
		Dio: NewDecoder(nil, h),
		Eb:  NewEncoderBytes(nil, h),
		Db:  NewDecoderBytes(nil, h),
	})
	d = &testHEDs[len(testHEDs)-1]
	// fmt.Printf("testHEDGet: <new> h: %v\n", h)
	return
}

func testSharedCodecEncode(ts interface{}, bsIn []byte, fn func([]byte) *bytes.Buffer,
	h Handle, bh *BasicHandle, useMust bool) (bs []byte, err error) {
	// bs = make([]byte, 0, approxSize)
	var e *Encoder
	var buf *bytes.Buffer
	useIO := testUseIoEncDec >= 0
	if testUseReset && !testUseParallel {
		hed := testHEDGet(h)
		if useIO {
			e = hed.Eio
		} else {
			e = hed.Eb
		}
	} else if useIO {
		e = NewEncoder(nil, h)
	} else {
		e = NewEncoderBytes(nil, h)
	}

	var oldWriteBufferSize int
	if useIO {
		buf = fn(bsIn)
		// set the encode options for using a buffer
		oldWriteBufferSize = bh.WriterBufferSize
		bh.WriterBufferSize = testUseIoEncDec
		if testUseIoWrapper {
			e.Reset(ioWriterWrapper{buf})
		} else {
			e.Reset(buf)
		}
	} else {
		bs = bsIn
		e.ResetBytes(&bs)
	}
	if useMust {
		e.MustEncode(ts)
	} else {
		err = e.Encode(ts)
	}
	if testUseIoEncDec >= 0 {
		bs = buf.Bytes()
		bh.WriterBufferSize = oldWriteBufferSize
	}
	return
}

func testSharedCodecDecoder(bs []byte, h Handle, bh *BasicHandle) (d *Decoder, oldReadBufferSize int) {
	// var buf *bytes.Reader
	useIO := testUseIoEncDec >= 0
	if testUseReset && !testUseParallel {
		hed := testHEDGet(h)
		// fmt.Printf("testSharedCodecDecoder: <hed>: %v\n", hed)
		if useIO {
			d = hed.Dio
		} else {
			d = hed.Db
		}
	} else if useIO {
		d = NewDecoder(nil, h)
	} else {
		d = NewDecoderBytes(nil, h)
		// fmt.Printf("testSharedCodecDecoder: not use IO: d: %v\n", d)
	}
	// fmt.Printf("testSharedCodecDecoder: d: %v\n", d)
	if useIO {
		buf := bytes.NewReader(bs)
		oldReadBufferSize = bh.ReaderBufferSize
		bh.ReaderBufferSize = testUseIoEncDec
		if testUseIoWrapper {
			d.Reset(ioReaderWrapper{buf})
		} else {
			d.Reset(buf)
		}
	} else {
		d.ResetBytes(bs)
	}
	return
}

func testSharedCodecDecoderAfter(d *Decoder, oldReadBufferSize int, bh *BasicHandle) {
	if testUseIoEncDec >= 0 {
		bh.ReaderBufferSize = oldReadBufferSize
	}
}

func testSharedCodecDecode(bs []byte, ts interface{}, h Handle, bh *BasicHandle, useMust bool) (err error) {
	d, oldReadBufferSize := testSharedCodecDecoder(bs, h, bh)
	if useMust {
		d.MustDecode(ts)
	} else {
		err = d.Decode(ts)
	}
	testSharedCodecDecoderAfter(d, oldReadBufferSize, bh)
	return
}
