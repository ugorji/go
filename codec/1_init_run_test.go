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
	var preferFloat, zeroAsNil, mapKeyAsStr bool
	_, _, _ = preferFloat, zeroAsNil, mapKeyAsStr
	// var nilColAsZeroLen, str2Raw, intAsStr bool
	bh := testBasicHandle(h)
	switch x := h.(type) {
	case *SimpleHandle:
		zeroAsNil = x.EncZeroValuesAsNil
	case *JsonHandle:
		mapKeyAsStr = x.MapKeyAsString
		preferFloat = x.PreferFloat
	}
	// var o gocmp.Option
	// var p []gocmp.Option
	// if zeroAsNil {
	// }
	// if mapKeyAsStr {
	// }
	// if bh.NilCollectionToZeroLength {
	var rc func(rv reflect.Value)
	rc = func(rv reflect.Value) {
		switch rv.Kind() {
		case reflect.Slice:
			if bh.NilCollectionToZeroLength && rv.IsNil() {
				rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
				// debugf("slice: %T: %v", hlBLUE, rv.Interface(), rv.Interface())
			}
			for i := 0; i < rv.Len(); i++ {
				rc(rv.Index(i))
			}
		case reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				rc(rv.Index(i))
			}
		case reflect.Map:
			if bh.NilCollectionToZeroLength && rv.IsNil() {
				rv.Set(reflect.MakeMapWithSize(rv.Type(), 0))
				// debugf("map: %T: %v", hlBLUE, rv.Interface(), rv.Interface())
			}
			for iter := rv.MapRange(); iter.Next(); {
				rv1 := iter.Key()
				rc(rv1) // MARKER - is this necessary?
				// if value is not addressable, we should clone it on the heap, modify it, and update it
				rv2 := iter.Value()
				if rv2.CanAddr() {
					rc(rv2)
				} else {
					rv3 := reflect.New(rv2.Type()).Elem()
					deepCopyRecursiveInto(rv2, rv3, make(map[uintptr]reflect.Value))
					rv2 = rv3
					rc(rv2)
					rv.SetMapIndex(rv1, rv2)
				}
			}
		case reflect.Chan:
			if bh.NilCollectionToZeroLength && rv.IsNil() {
				rv.Set(reflect.MakeChan(rv.Type(), 0))
			}
		case reflect.Pointer, reflect.Interface:
			if !rv.IsNil() {
				rc(rv.Elem())
			}
		case reflect.Struct:
			if !rv.CanAddr() {
				// debugf("struct not addressable: %v", hlRED, rv.Type())
			}
			tt := rv.Type()
			for i, n := 0, rv.NumField(); i < n; i++ {
				sf := tt.Field(i)
				if sf.IsExported() {
					rv2 := rv.Field(i)
					if kk := rv2.Kind(); (kk == reflect.Slice || kk == reflect.Map) && rv2.IsNil() {
						// debugf("struct field: [%v/%s] (%v) nil: %v", hlYELLOW, tt, sf.Name, rv2.Type(), rv2.IsNil())
					}
					rc(rv2)
				}
			}
			// }
			// tt := rv.Type()
			// for i := 0; i < tt.NumField(); i++ {
			// 	if tt.Field(i).IsExported() {
			// 		if rv2 := rv.Field(i); rv2.CanSet() { // in case an unaddressable struct
			// 			rc(rv2)
			// 		}
			// 	}
			// }
		}
	}

	f := func(in interface{}) (out interface{}) {
		// debugf("start transform: %T(%p) (isnil: %v)", hlGREEN, in, in, in == nil)
		v := deepCopyValue(in)
		if v.IsValid() {
			rc(v)
			out = v.Interface()
		}
		// debugf("transforming %T(%p) --> %T(%p) (isnil: %v)", hlGREEN, in, in, out, out, out == nil)
		return
	}

	// o = gocmpopts.AcyclicTransformer("T1", f)
	// p = append(p, o)
	// if !gocmp.Equal(v1, v2, p...) {
	// 	if testv.UseDiff {
	// 		err = errors.New(gocmp.Diff(v1, v2, p...))
	// 	} else {
	// 		err = errDeepEqualNotMatch
	// 	}
	// }
	// return

	v11 := f(v1)
	if !reflect.DeepEqual(v11, v2) {
		if testv.UseDiff {
			err = errors.New(gocmp.Diff(v11, v2))
		} else {
			err = errDeepEqualNotMatch
		}
	}
	return

	// o = cmp.Transformer("T1", func(in interface{}) (out interface{}) {
	// 	out.Real, out.Imag = real(in), imag(in)
	// 	return out
	// })
}

func deepCopyValue(val interface{}) (rv reflect.Value) {
	if val == nil {
		return
	}
	orig := reflect.ValueOf(val)
	visited := make(map[uintptr]reflect.Value)

	if orig.Kind() == reflect.Ptr {
		if orig.IsNil() {
			return reflect.Zero(orig.Type()) // or nil for typed nil
		}
		rv = reflect.New(orig.Type().Elem())
		visited[orig.Pointer()] = rv
		deepCopyRecursiveInto(orig.Elem(), rv.Elem(), visited)
	} else {
		rv = reflect.New(orig.Type()).Elem()
		deepCopyRecursiveInto(orig, rv, visited)
	}
	return
}

// deepCopyRecursiveInto copies from 'original' into 'target'. Both are concrete values.
func deepCopyRecursiveInto(orig, target reflect.Value, visited map[uintptr]reflect.Value) {
	for orig.Kind() == reflect.Interface { // Ptrs inside structs handled by checking visited
		if orig.IsNil() {
			return
		}
		orig = orig.Elem()
	}
	switch orig.Kind() {
	case reflect.Ptr:
		if orig.IsNil() {
			return // Target remains zero, which is correct for a nil pointer
		}
		addr := orig.Pointer()
		if copiedPtr, ok := visited[addr]; ok {
			if target.Type() == copiedPtr.Type() {
				target.Set(copiedPtr) // Set target to the already copied pointer
			} else if target.Type() == copiedPtr.Elem().Type() && copiedPtr.Kind() == reflect.Ptr {
				target.Set(copiedPtr.Elem()) // If target is value, set to copied pointer's element
			}
			return
		}
		// New pointer needed for the field/element
		elemCopy := reflect.New(orig.Elem().Type())
		visited[addr] = elemCopy // Store the new pointer
		deepCopyRecursiveInto(orig.Elem(), elemCopy.Elem(), visited)
		target.Set(elemCopy)
	case reflect.Struct:
		for i := 0; i < orig.NumField(); i++ {
			deepCopyRecursiveInto(orig.Field(i), target.Field(i), visited)
			// if target.Field(i).CanSet() {
			// 	deepCopyRecursiveInto(orig.Field(i), target.Field(i), visited)
			// }
		}
	case reflect.Slice:
		if orig.IsNil() {
			// target.Set(reflect.Zero(orig.Type())) // target becomes a nil slice
			return // target is already a zero slice if newly created.
		}
		if target.IsZero() || target.Cap() < orig.Cap() { // Check if target needs to be (re)made
			target.Set(reflect.MakeSlice(orig.Type(), orig.Len(), orig.Cap()))
		} else {
			target.SetLen(orig.Len()) // Ensure length matches
		}
		for i := 0; i < orig.Len(); i++ {
			deepCopyRecursiveInto(orig.Index(i), target.Index(i), visited)
		}
	case reflect.Map:
		if orig.IsNil() {
			// target.Set(reflect.Zero(orig.Type())) // target becomes a nil map
			return
		}
		if target.IsZero() {
			target.Set(reflect.MakeMap(orig.Type()))
		}

		iter := orig.MapRange()
		for iter.Next() {
			kOrig := iter.Key()
			vOrig := iter.Value()
			kTarget := reflect.New(kOrig.Type()).Elem()
			vTarget := reflect.New(vOrig.Type()).Elem()
			deepCopyRecursiveInto(kOrig, kTarget, visited)
			deepCopyRecursiveInto(vOrig, vTarget, visited)
			target.SetMapIndex(kTarget, vTarget)
		}
	default: // Basic types: int, string, bool, etc.
		if orig.Type().AssignableTo(target.Type()) {
			target.Set(orig)
		}
	}
}
