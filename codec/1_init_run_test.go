// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"io"
	"reflect"
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

type ioReaderWrapper struct {
	r io.Reader
}

func (x ioReaderWrapper) Read(p []byte) (n int, err error) {
	return x.r.Read(p)
}

type ioWriterWrapper struct {
	w io.Writer
}

func (x ioWriterWrapper) Write(p []byte) (n int, err error) {
	return x.w.Write(p)
}

// the handles are declared here, and initialized during the init function.
//
// Note the following:
//   - if running parallel tests, then skip all tests which modify the Handle.
//     This prevents data races during the test execution.
var (
	// testNoopH    = NoopHandle(8)
	testBincH    *BincHandle
	testMsgpackH *MsgpackHandle
	testJsonH    *JsonHandle
	testSimpleH  *SimpleHandle
	testCborH    *CborHandle

	testHandles []Handle

	testHEDs []testHED
)

func init() {
	// doTestInit()
	testPreInitFns = append(testPreInitFns, doTestInit)
	testPostInitFns = append(testPostInitFns, doTestPostInit)
	// doTestInit MUST be the first function executed during a reinit
	// testReInitFns = slices.Insert(testReInitFns, 0, doTestInit)
	// testReInitFns = slices.Insert(testReInitFns, 0, doTestReinit)
	testReInitFns = append(testReInitFns, doTestInit, doTestPostInit)
}

func doTestInit() {
	testBincH = &BincHandle{}
	testMsgpackH = &MsgpackHandle{}
	testJsonH = &JsonHandle{}
	testSimpleH = &SimpleHandle{}
	testCborH = &CborHandle{}

	// testHEDs = make([]testHED, 0, 32)

	// JSON should do HTMLCharsAsIs by default
	testJsonH.HTMLCharsAsIs = true
	// testJsonH.InternString = true

	testHandles = nil
	testHandles = append(testHandles, testSimpleH, testJsonH,
		testCborH, testMsgpackH, testBincH)

	testHEDs = nil
}

func doTestPostInit() {
	testUpdateBasicHandleOptions(&testBincH.BasicHandle)
	testUpdateBasicHandleOptions(&testMsgpackH.BasicHandle)
	testUpdateBasicHandleOptions(&testJsonH.BasicHandle)
	testUpdateBasicHandleOptions(&testSimpleH.BasicHandle)
	testUpdateBasicHandleOptions(&testCborH.BasicHandle)
}

// func doTestReinit() {
// 	// doTestInit()
// 	// MARKER 2025 - instead, just reset them all
// 	for _, h := range testHandles {
// 		bh := testBasicHandle(h)
// 		bh.basicHandleRuntimeState = basicHandleRuntimeState{}
// 		atomic.StoreUint32(&bh.inited, 0)
// 		initHandle(h)
// 	}
// 	testHEDs = nil
// }

func testHEDGet(h Handle) (d *testHED) {
	for i := range testHEDs {
		v := &testHEDs[i]
		if v.H == h {
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
	return
}

func testSharedCodecEncode(ts interface{}, bsIn []byte,
	fn func([]byte) *bytes.Buffer,
	h Handle, useMust bool) (bs []byte, err error) {
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

	// var oldWriteBufferSize int
	if useIO {
		buf = fn(bsIn)
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
	}
	return
}

func testSharedCodecDecoder(bs []byte, h Handle) (d *Decoder) {
	// var buf *bytes.Reader
	useIO := testUseIoEncDec >= 0
	if testUseReset && !testUseParallel {
		hed := testHEDGet(h)
		if useIO {
			d = hed.Dio
		} else {
			d = hed.Db
		}
	} else if useIO {
		d = NewDecoder(nil, h)
	} else {
		d = NewDecoderBytes(nil, h)
	}
	if useIO {
		buf := bytes.NewReader(bs)
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

func testSharedCodecDecode(bs []byte, ts interface{}, h Handle, useMust bool) (err error) {
	d := testSharedCodecDecoder(bs, h)
	if useMust {
		d.MustDecode(ts)
	} else {
		err = d.Decode(ts)
	}
	return
}

func testUpdateBasicHandleOptions(bh *BasicHandle) {
	// bh.clearInited() // so it is reinitialized next time around // MARKER 2025 (may have to put it back)
	// pre-fill them first
	bh.EncodeOptions = testEncodeOptions
	bh.DecodeOptions = testDecodeOptions
	bh.RPCOptions = testRPCOptions
	// bh.InterfaceReset = true
	// bh.PreferArrayOverSlice = true
	// modify from flag'ish things
	// bh.MaxInitLen = testMaxInitLen
}

type testNameBasicHandle struct {
	n string
	h *BasicHandle
}

func testUpdateExts(nhs ...testNameBasicHandle) {
	var tI64Ext wrapInt64Ext
	var tUintToBytesExt testUintToBytesExt
	var tBytesExt wrapBytesExt
	var tTimeBytesExt timeBytesExt
	var tUnixTimeIntfExt testUnixNanoTimeInterfaceExt

	timeExtEncFn := func(rv reflect.Value) ([]byte, error) { return basicTestExtEncFn(tTimeBytesExt, rv) }
	timeExtDecFn := func(rv reflect.Value, bs []byte) error { return basicTestExtDecFn(tTimeBytesExt, rv, bs) }
	wrapInt64ExtEncFn := func(rv reflect.Value) ([]byte, error) { return basicTestExtEncFn(&tI64Ext, rv) }
	wrapInt64ExtDecFn := func(rv reflect.Value, bs []byte) error { return basicTestExtDecFn(&tI64Ext, rv, bs) }

	var bh *BasicHandle
	ix := func(rt reflect.Type, tag uint64, ext interface{}) {
		halt.onerror(bh.SetExt(rt, tag, makeExt(ext)))
	}

	for _, nh := range nhs {
		bh = nh.h
		ix(testSelfExtTyp, 78, SelfExt)
		ix(testSelfExt2Typ, 79, SelfExt)
		ix(wrapBytesTyp, 32, &tBytesExt)
		ix(testUintToBytesTyp, 33, &tUintToBytesExt)
		// Now, add extensions for the type wrapInt64 and wrapBytes,
		// so we can execute the Encode/Decode Ext paths.
		if nh.n == "simple" {
			halt.onerror(bh.AddExt(wrapInt64Typ, 16, wrapInt64ExtEncFn, wrapInt64ExtDecFn))
		} else {
			ix(wrapInt64Typ, 16, &tI64Ext)
		}

		// add extensions for time.Time
		switch nh.n {
		case "cbor":
			ix(timeTyp, 1, tUnixTimeIntfExt)
		case "binc":
			// ix(timeTyp, 1, timeExt{}) // time is builtin for binc
		case "msgpack":
			ix(timeTyp, 1, tTimeBytesExt)
		case "simple":
			halt.onerror(bh.AddExt(timeTyp, 1, timeExtEncFn, timeExtDecFn))
		}
	}
}
