// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"errors"
	"io"
	"reflect"

	// using codec.XXX directly
	// . "github.com/ugorji/go/codec"

	gocmp "github.com/google/go-cmp/cmp"
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
// 	// instead, just reset them all
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
	useIO := testv.E.WriterBufferSize >= 0
	if testv.UseReset && !testv.UseParallel {
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
		if testv.UseIoWrapper {
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
	if testv.E.WriterBufferSize >= 0 {
		bs = buf.Bytes()
	}
	return
}

func testSharedCodecDecoder(bs []byte, h Handle) (d *Decoder) {
	// var buf *bytes.Reader
	useIO := testv.D.ReaderBufferSize >= 0
	if testv.UseReset && !testv.UseParallel {
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
		if testv.UseIoWrapper {
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
	// cleanInited() not needed, as we re-create the Handles on each reinit
	// bh.clearInited() // so it is reinitialized next time around
	// pre-fill them first
	bh.EncodeOptions = testv.E
	bh.DecodeOptions = testv.D
	bh.RPCOptions = testv.R
	// bh.InterfaceReset = true
	// bh.PreferArrayOverSlice = true
	// modify from flag'ish things
	// bh.MaxInitLen = testMaxInitLen
}

var errDeepEqualNotMatch = errors.New("not match")

// var testCmpOpts []gocmp.Option
//
// var testCmpOpts = []cmp.Option{
// 	cmpopts.EquateNaNs(),
// 	cmpopts.EquateApprox(0.001, 0.001),
// 	cmpopts.SortMaps(func(a, b float32) bool { return a < b }),
// 	cmpopts.SortMaps(func(a, b float64) bool { return a < b }),
// }

// perform a simple DeepEqual expecting same values
func testEqual(v1, v2 interface{}) (err error) {
	if !reflect.DeepEqual(v1, v2) {
		if testv.UseDiff {
			err = errors.New(gocmp.Diff(v1, v2))
		} else {
			err = errDeepEqualNotMatch
		}
	}
	return
}

// perform a comparison taking Handle fields into consideration
func testEqualH(v1, v2 interface{}, h Handle) (err error) {
	var preferFloat, zeroAsNil, mapKeyAsStr, isJson bool
	_, _, _, _ = preferFloat, zeroAsNil, mapKeyAsStr, isJson
	// var nilColAsZeroLen, str2Raw, intAsStr bool
	bh := testBasicHandle(h)
	switch x := h.(type) {
	case *SimpleHandle:
		zeroAsNil = x.EncZeroValuesAsNil
	case *JsonHandle:
		mapKeyAsStr = x.MapKeyAsString
		preferFloat = x.PreferFloat
		isJson = true
	}

	// create a clone that honors the Handle options
	// then compare that using reflect.DeepEqual

	const structExportedFieldsOnly = false

	visited := make(map[uintptr]reflect.Value)

	var rcopy func(src, target reflect.Value)
	rcopy = func(src, target reflect.Value) {
	TOP:
		switch src.Kind() {
		case reflect.Interface:
			if src.IsNil() {
				return
			}
			src = src.Elem()
			goto TOP
		case reflect.Ptr:
			if src.IsNil() {
				return // target remains zero, which is correct for a nil pointer
			}
			addr := src.Pointer()
			if copiedPtr, ok := visited[addr]; ok {
				if target.Type() == copiedPtr.Type() {
					target.Set(copiedPtr) // Set target to the already copied pointer
				} else if target.Type() == copiedPtr.Elem().Type() && copiedPtr.Kind() == reflect.Ptr {
					target.Set(copiedPtr.Elem()) // If target is value, set to copied pointer's element
				}
				return
			}
			// New pointer needed for the field/element
			elemCopy := reflect.New(src.Elem().Type())
			visited[addr] = elemCopy // Store the new pointer
			rcopy(src.Elem(), elemCopy.Elem())
			target.Set(elemCopy)
		case reflect.Array:
			target.Set(src)
			for i, slen := 0, src.Len(); i < slen; i++ {
				rcopy(src.Index(i), target.Index(i))
			}
		case reflect.Slice:
			if src.IsNil() {
				if bh.NilCollectionToZeroLength {
					target.Set(reflect.MakeSlice(src.Type(), 0, 0))
				}
				return
				// debugf("slice: %T: %v", hlBLUE, rv.Interface(), rv.Interface())
			}
			slen := src.Len()
			target.Set(reflect.MakeSlice(src.Type(), slen, src.Cap()))
			for i := 0; i < slen; i++ {
				rcopy(src.Index(i), target.Index(i))
			}
		case reflect.Map:
			if src.IsNil() {
				if bh.NilCollectionToZeroLength {
					target.Set(reflect.MakeMapWithSize(src.Type(), 0))
				}
				return
				// debugf("map: %T: %v", hlBLUE, rv.Interface(), rv.Interface())
			}
			target.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
			iter := src.MapRange()
			for iter.Next() {
				kSrc := iter.Key()
				vSrc := iter.Value()
				kTarget := reflect.New(kSrc.Type()).Elem()
				vTarget := reflect.New(vSrc.Type()).Elem()
				rcopy(kSrc, kTarget)
				rcopy(vSrc, vTarget)
				target.SetMapIndex(kTarget, vTarget)
			}
		case reflect.Chan:
			if src.IsNil() {
				if bh.NilCollectionToZeroLength {
					target.Set(reflect.MakeChan(src.Type(), 0))
				}
				return
			}
		case reflect.Struct:
			tt := src.Type()
			for i, n := 0, src.NumField(); i < n; i++ {
				if !structExportedFieldsOnly || tt.Field(i).IsExported() {
					// rv2 := rv.Field(i)
					// if kk := rv2.Kind(); (kk == reflect.Slice || kk == reflect.Map) && rv2.IsNil() {
					// 	// debugf("struct field: [%v/%s] (%v) nil: %v", hlYELLOW, tt, sf.Name, rv2.Type(), rv2.IsNil())
					// }
					rcopy(src.Field(i), target.Field(i))
				}
			}
		case reflect.String:
			// MARKER 2025 - need to all use same functions to compare
			// s := src.String()
			// if isJson && bh.StringToRaw {
			// 	s = base64.StdEncoding.EncodeToString(bytesView(s))
			// 	// dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
			// 	// n, err := base64.StdEncoding.Decode(dbuf, bytesView(s))
			// 	// if err == nil {
			// 	// 	s = stringView(dbuf[:n])
			// 	// }
			// }
			// target.SetString(s)
			target.Set(src)
		default: // Basic types: int, string, bool, etc.
			target.Set(src)
			// if src.Type().AssignableTo(target.Type()) {
			// 	target.Set(src)
			// }
		}
	}

	deepcopy := func(in interface{}) (rv reflect.Value) {
		// debugf("start transform: %T(%p) (isnil: %v)", hlGREEN, in, in, in == nil)
		clear(visited)
		src := reflect.ValueOf(in)
		if src.Kind() == reflect.Ptr {
			if src.IsNil() {
				rv = reflect.Zero(src.Type()) // or nil for typed nil
			} else {
				rv = reflect.New(src.Type().Elem())
				visited[src.Pointer()] = rv
				rcopy(src.Elem(), rv.Elem())
			}
		} else {
			rv = reflect.New(src.Type()).Elem()
			rcopy(src, rv)
		}
		// debugf("transforming %T(%p) --> %T(%p) (isnil: %v)", hlGREEN, in, in, out, out, out == nil)
		return
	}

	if v1 != nil {
		v1 = deepcopy(v1).Interface()
	}
	if !reflect.DeepEqual(v1, v2) {
		if testv.UseDiff {
			err = errors.New(gocmp.Diff(v1, v2))
		} else {
			err = errDeepEqualNotMatch
		}
	}
	return
}

// func testEqualH(v1, v2 interface{}, h Handle) (err error) {
// 	// ...
// 	//
// 	// var o gocmp.Option
// 	// var p []gocmp.Option
// 	// if zeroAsNil {
// 	// }
// 	// if mapKeyAsStr {
// 	// }
// 	// if bh.NilCollectionToZeroLength {
// 	//
// 	// ...
// 	//
// 	// o = gocmpopts.AcyclicTransformer("T1", f)
// 	// p = append(p, o)
// 	// if !gocmp.Equal(v1, v2, p...) {
// 	// 	if testv.UseDiff {
// 	// 		err = errors.New(gocmp.Diff(v1, v2, p...))
// 	// 	} else {
// 	// 		err = errDeepEqualNotMatch
// 	// 	}
// 	// }
// 	// return
// 	//
// 	// ...
// 	//
// 	// o = cmp.Transformer("T1", func(in interface{}) (out interface{}) {
// 	// 	out.Real, out.Imag = real(in), imag(in)
// 	// 	return out
// 	// })
// }
