// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	// using codec.XXX directly
	// . "github.com/ugorji/go/codec"
)

type testHED struct {
	H Handle
	E *Encoder
	D *Decoder
}

var (
	// testNoopH    = NoopHandle(8)
	testMsgpackH = &MsgpackHandle{}
	testBincH    = &BincHandle{}
	testSimpleH  = &SimpleHandle{}
	testCborH    = &CborHandle{}
	testJsonH    = &JsonHandle{}

	testHandles []Handle

	testHEDs []testHED
)

// variables that are not flags, but which can configure the handles
var (
	testEncodeOptions EncodeOptions
	testDecodeOptions DecodeOptions
	testRPCOptions    RPCOptions
)

func init() {
	testHEDs = make([]testHED, 0, 32)
	testHandles = append(testHandles,
		// testNoopH,
		testMsgpackH, testBincH, testSimpleH, testCborH, testJsonH)
	// JSON should do HTMLCharsAsIs by default
	testJsonH.HTMLCharsAsIs = true
	// testJsonH.InternString = true
	testReInitFns = append(testReInitFns, func() { testHEDs = nil })
}

func testHEDGet(h Handle) *testHED {
	for i := range testHEDs {
		v := &testHEDs[i]
		if v.H == h {
			return v
		}
	}
	testHEDs = append(testHEDs, testHED{h, NewEncoder(nil, h), NewDecoder(nil, h)})
	return &testHEDs[len(testHEDs)-1]
}

func testSharedCodecEncode(ts interface{}, bsIn []byte, fn func([]byte) *bytes.Buffer,
	h Handle, bh *BasicHandle, useMust bool) (bs []byte, err error) {
	// bs = make([]byte, 0, approxSize)
	var e *Encoder
	var buf *bytes.Buffer
	if testUseReset && !testUseParallel {
		e = testHEDGet(h).E
	} else {
		e = NewEncoder(nil, h)
	}
	var oldWriteBufferSize int
	if testUseIoEncDec >= 0 {
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
	if testUseReset && !testUseParallel {
		d = testHEDGet(h).D
	} else {
		d = NewDecoder(nil, h)
	}
	if testUseIoEncDec >= 0 {
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
