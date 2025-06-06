//go:build !codec.notmammoth

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from mammoth_test.go.tmpl - DO NOT EDIT.

package codec

import "testing"
import "fmt"
import "reflect"

// TestMammoth has all the different paths optimized in fastpath
// It has all the primitives, slices and maps.
//
// For each of those types, it has a pointer and a non-pointer field.

func init() { _ = fmt.Printf } // so we can include fmt as needed

type TestMammoth struct {
	FIntf       interface{}
	FptrIntf    *interface{}
	FString     string
	FptrString  *string
	FBytes      []byte
	FptrBytes   *[]byte
	FFloat32    float32
	FptrFloat32 *float32
	FFloat64    float64
	FptrFloat64 *float64
	FUint       uint
	FptrUint    *uint
	FUint8      uint8
	FptrUint8   *uint8
	FUint16     uint16
	FptrUint16  *uint16
	FUint32     uint32
	FptrUint32  *uint32
	FUint64     uint64
	FptrUint64  *uint64
	FUintptr    uintptr
	FptrUintptr *uintptr
	FInt        int
	FptrInt     *int
	FInt8       int8
	FptrInt8    *int8
	FInt16      int16
	FptrInt16   *int16
	FInt32      int32
	FptrInt32   *int32
	FInt64      int64
	FptrInt64   *int64
	FBool       bool
	FptrBool    *bool

	FSliceIntf        []interface{}
	FptrSliceIntf     *[]interface{}
	Farr4SliceIntf    [4]interface{}
	FSliceString      []string
	FptrSliceString   *[]string
	Farr4SliceString  [4]string
	FSliceBytes       [][]byte
	FptrSliceBytes    *[][]byte
	Farr4SliceBytes   [4][]byte
	FSliceFloat32     []float32
	FptrSliceFloat32  *[]float32
	Farr4SliceFloat32 [4]float32
	FSliceFloat64     []float64
	FptrSliceFloat64  *[]float64
	Farr4SliceFloat64 [4]float64
	FSliceUint8       []uint8
	FptrSliceUint8    *[]uint8
	Farr4SliceUint8   [4]uint8
	FSliceUint64      []uint64
	FptrSliceUint64   *[]uint64
	Farr4SliceUint64  [4]uint64
	FSliceInt         []int
	FptrSliceInt      *[]int
	Farr4SliceInt     [4]int
	FSliceInt32       []int32
	FptrSliceInt32    *[]int32
	Farr4SliceInt32   [4]int32
	FSliceInt64       []int64
	FptrSliceInt64    *[]int64
	Farr4SliceInt64   [4]int64
	FSliceBool        []bool
	FptrSliceBool     *[]bool
	Farr4SliceBool    [4]bool

	FMapStringIntf       map[string]interface{}
	FptrMapStringIntf    *map[string]interface{}
	FMapStringString     map[string]string
	FptrMapStringString  *map[string]string
	FMapStringBytes      map[string][]byte
	FptrMapStringBytes   *map[string][]byte
	FMapStringUint8      map[string]uint8
	FptrMapStringUint8   *map[string]uint8
	FMapStringUint64     map[string]uint64
	FptrMapStringUint64  *map[string]uint64
	FMapStringInt        map[string]int
	FptrMapStringInt     *map[string]int
	FMapStringInt32      map[string]int32
	FptrMapStringInt32   *map[string]int32
	FMapStringFloat64    map[string]float64
	FptrMapStringFloat64 *map[string]float64
	FMapStringBool       map[string]bool
	FptrMapStringBool    *map[string]bool
	FMapUint8Intf        map[uint8]interface{}
	FptrMapUint8Intf     *map[uint8]interface{}
	FMapUint8String      map[uint8]string
	FptrMapUint8String   *map[uint8]string
	FMapUint8Bytes       map[uint8][]byte
	FptrMapUint8Bytes    *map[uint8][]byte
	FMapUint8Uint8       map[uint8]uint8
	FptrMapUint8Uint8    *map[uint8]uint8
	FMapUint8Uint64      map[uint8]uint64
	FptrMapUint8Uint64   *map[uint8]uint64
	FMapUint8Int         map[uint8]int
	FptrMapUint8Int      *map[uint8]int
	FMapUint8Int32       map[uint8]int32
	FptrMapUint8Int32    *map[uint8]int32
	FMapUint8Float64     map[uint8]float64
	FptrMapUint8Float64  *map[uint8]float64
	FMapUint8Bool        map[uint8]bool
	FptrMapUint8Bool     *map[uint8]bool
	FMapUint64Intf       map[uint64]interface{}
	FptrMapUint64Intf    *map[uint64]interface{}
	FMapUint64String     map[uint64]string
	FptrMapUint64String  *map[uint64]string
	FMapUint64Bytes      map[uint64][]byte
	FptrMapUint64Bytes   *map[uint64][]byte
	FMapUint64Uint8      map[uint64]uint8
	FptrMapUint64Uint8   *map[uint64]uint8
	FMapUint64Uint64     map[uint64]uint64
	FptrMapUint64Uint64  *map[uint64]uint64
	FMapUint64Int        map[uint64]int
	FptrMapUint64Int     *map[uint64]int
	FMapUint64Int32      map[uint64]int32
	FptrMapUint64Int32   *map[uint64]int32
	FMapUint64Float64    map[uint64]float64
	FptrMapUint64Float64 *map[uint64]float64
	FMapUint64Bool       map[uint64]bool
	FptrMapUint64Bool    *map[uint64]bool
	FMapIntIntf          map[int]interface{}
	FptrMapIntIntf       *map[int]interface{}
	FMapIntString        map[int]string
	FptrMapIntString     *map[int]string
	FMapIntBytes         map[int][]byte
	FptrMapIntBytes      *map[int][]byte
	FMapIntUint8         map[int]uint8
	FptrMapIntUint8      *map[int]uint8
	FMapIntUint64        map[int]uint64
	FptrMapIntUint64     *map[int]uint64
	FMapIntInt           map[int]int
	FptrMapIntInt        *map[int]int
	FMapIntInt32         map[int]int32
	FptrMapIntInt32      *map[int]int32
	FMapIntFloat64       map[int]float64
	FptrMapIntFloat64    *map[int]float64
	FMapIntBool          map[int]bool
	FptrMapIntBool       *map[int]bool
	FMapInt32Intf        map[int32]interface{}
	FptrMapInt32Intf     *map[int32]interface{}
	FMapInt32String      map[int32]string
	FptrMapInt32String   *map[int32]string
	FMapInt32Bytes       map[int32][]byte
	FptrMapInt32Bytes    *map[int32][]byte
	FMapInt32Uint8       map[int32]uint8
	FptrMapInt32Uint8    *map[int32]uint8
	FMapInt32Uint64      map[int32]uint64
	FptrMapInt32Uint64   *map[int32]uint64
	FMapInt32Int         map[int32]int
	FptrMapInt32Int      *map[int32]int
	FMapInt32Int32       map[int32]int32
	FptrMapInt32Int32    *map[int32]int32
	FMapInt32Float64     map[int32]float64
	FptrMapInt32Float64  *map[int32]float64
	FMapInt32Bool        map[int32]bool
	FptrMapInt32Bool     *map[int32]bool
}

// -----------

// Increase codecoverage by covering all the codecgen paths, in fastpath ....
//
// Note: even though this is built based on fastpath, we will run these tests
// in all modes, including notfastpath, etc.
//
// Add test file for creating a mammoth generated file as _mammoth_generated.go
//
// Now, add some types:
//  - some that implement BinaryMarshal, TextMarshal, JSONMarshal, and one that implements none of it
//  - create a wrapper type that includes TestMammoth2, with it in slices, and maps, and the custom types
//  - this wrapper object is what we work encode/decode (so that the codecgen methods are called)

type testMammoth2Binary uint64

func (x testMammoth2Binary) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 8)
	bigenstd.PutUint64(data, uint64(x))
	return
}
func (x *testMammoth2Binary) UnmarshalBinary(data []byte) (err error) {
	*x = testMammoth2Binary(bigenstd.Uint64(data))
	return
}

type testMammoth2Text uint64

func (x testMammoth2Text) MarshalText() (data []byte, err error) {
	data = []byte(fmt.Sprintf("%b", uint64(x)))
	return
}
func (x *testMammoth2Text) UnmarshalText(data []byte) (err error) {
	_, err = fmt.Sscanf(string(data), "%b", (*uint64)(x))
	return
}

type testMammoth2Json uint64

func (x testMammoth2Json) MarshalJSON() (data []byte, err error) {
	data = []byte(fmt.Sprintf("%v", uint64(x)))
	return
}
func (x *testMammoth2Json) UnmarshalJSON(data []byte) (err error) {
	_, err = fmt.Sscanf(string(data), "%v", (*uint64)(x))
	return
}

type testMammoth2Basic [4]uint64

type TestMammoth2Wrapper struct {
	V TestMammoth
	T testMammoth2Text
	B testMammoth2Binary
	J testMammoth2Json
	C testMammoth2Basic
	M map[testMammoth2Basic]TestMammoth
	L []TestMammoth
	A [4]int64

	Tcomplex128 complex128
	Tcomplex64  complex64
	Tbytes      []uint8
	Tpbytes     *[]uint8
}

// -----------

type typMbsSliceIntf []interface{}

func (_ typMbsSliceIntf) MapBySlice() {}

type typMbsSliceString []string

func (_ typMbsSliceString) MapBySlice() {}

type typMbsSliceBytes [][]byte

func (_ typMbsSliceBytes) MapBySlice() {}

type typMbsSliceFloat32 []float32

func (_ typMbsSliceFloat32) MapBySlice() {}

type typMbsSliceFloat64 []float64

func (_ typMbsSliceFloat64) MapBySlice() {}

type typMbsSliceUint8 []uint8

func (_ typMbsSliceUint8) MapBySlice() {}

type typMbsSliceUint64 []uint64

func (_ typMbsSliceUint64) MapBySlice() {}

type typMbsSliceInt []int

func (_ typMbsSliceInt) MapBySlice() {}

type typMbsSliceInt32 []int32

func (_ typMbsSliceInt32) MapBySlice() {}

type typMbsSliceInt64 []int64

func (_ typMbsSliceInt64) MapBySlice() {}

type typMbsSliceBool []bool

func (_ typMbsSliceBool) MapBySlice() {}

type typMapMapStringIntf map[string]interface{}
type typMapMapStringString map[string]string
type typMapMapStringBytes map[string][]byte
type typMapMapStringUint8 map[string]uint8
type typMapMapStringUint64 map[string]uint64
type typMapMapStringInt map[string]int
type typMapMapStringInt32 map[string]int32
type typMapMapStringFloat64 map[string]float64
type typMapMapStringBool map[string]bool
type typMapMapUint8Intf map[uint8]interface{}
type typMapMapUint8String map[uint8]string
type typMapMapUint8Bytes map[uint8][]byte
type typMapMapUint8Uint8 map[uint8]uint8
type typMapMapUint8Uint64 map[uint8]uint64
type typMapMapUint8Int map[uint8]int
type typMapMapUint8Int32 map[uint8]int32
type typMapMapUint8Float64 map[uint8]float64
type typMapMapUint8Bool map[uint8]bool
type typMapMapUint64Intf map[uint64]interface{}
type typMapMapUint64String map[uint64]string
type typMapMapUint64Bytes map[uint64][]byte
type typMapMapUint64Uint8 map[uint64]uint8
type typMapMapUint64Uint64 map[uint64]uint64
type typMapMapUint64Int map[uint64]int
type typMapMapUint64Int32 map[uint64]int32
type typMapMapUint64Float64 map[uint64]float64
type typMapMapUint64Bool map[uint64]bool
type typMapMapIntIntf map[int]interface{}
type typMapMapIntString map[int]string
type typMapMapIntBytes map[int][]byte
type typMapMapIntUint8 map[int]uint8
type typMapMapIntUint64 map[int]uint64
type typMapMapIntInt map[int]int
type typMapMapIntInt32 map[int]int32
type typMapMapIntFloat64 map[int]float64
type typMapMapIntBool map[int]bool
type typMapMapInt32Intf map[int32]interface{}
type typMapMapInt32String map[int32]string
type typMapMapInt32Bytes map[int32][]byte
type typMapMapInt32Uint8 map[int32]uint8
type typMapMapInt32Uint64 map[int32]uint64
type typMapMapInt32Int map[int32]int
type typMapMapInt32Int32 map[int32]int32
type typMapMapInt32Float64 map[int32]float64
type typMapMapInt32Bool map[int32]bool

func __doTestMammothSlices(t *testing.T, h Handle) {
	var v17va [8]interface{}
	for _, v := range [][]interface{}{nil, {}, {"string-is-an-interface-2", nil, nil, "string-is-an-interface-3"}} {
		var v17v1, v17v2 []interface{}
		var bs17 []byte
		v17v1 = v
		bs17 = testMarshalErr(v17v1, h, t, "enc-slice-v17")
		if v == nil {
			v17v2 = make([]interface{}, 2)
			testUnmarshalErr(v17v2, bs17, h, t, "dec-slice-v17")
			testDeepEqualErr(v17v2[0], v17v2[1], t, "equal-slice-v17") // should not change
			testDeepEqualErr(len(v17v2), 2, t, "equal-slice-v17")      // should not change
			v17v2 = make([]interface{}, 2)
			testUnmarshalErr(reflect.ValueOf(v17v2), bs17, h, t, "dec-slice-v17-noaddr") // non-addressable value
			testDeepEqualErr(v17v2[0], v17v2[1], t, "equal-slice-v17-noaddr")            // should not change
			testDeepEqualErr(len(v17v2), 2, t, "equal-slice-v17")                        // should not change
		} else {
			v17v2 = make([]interface{}, len(v))
			testUnmarshalErr(v17v2, bs17, h, t, "dec-slice-v17")
			testDeepEqualErrHandle(v17v1, v17v2, h, t, "equal-slice-v17")
			v17v2 = make([]interface{}, len(v))
			testUnmarshalErr(reflect.ValueOf(v17v2), bs17, h, t, "dec-slice-v17-noaddr") // non-addressable value
			testDeepEqualErrHandle(v17v1, v17v2, h, t, "equal-slice-v17-noaddr")
		}
		testReleaseBytes(bs17)
		// ...
		bs17 = testMarshalErr(&v17v1, h, t, "enc-slice-v17-p")
		v17v2 = nil
		testUnmarshalErr(&v17v2, bs17, h, t, "dec-slice-v17-p")
		testDeepEqualErrHandle(v17v1, v17v2, h, t, "equal-slice-v17-p")
		v17va = [8]interface{}{} // clear the array
		testUnmarshalErr(&v17va, bs17, h, t, "dec-array-v17-p-1")
		if v17v1 == nil && v17v2 == nil {
			v17v2 = []interface{}{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v17va[:len(v17v2)], v17v2, h, t, "equal-array-v17-p-1")
		v17va = [8]interface{}{} // clear the array
		v17v2 = v17va[:1:1]
		testUnmarshalErr(&v17v2, bs17, h, t, "dec-slice-v17-p-1")
		testDeepEqualErrHandle(v17v1, v17v2, h, t, "equal-slice-v17-p-1")
		v17va = [8]interface{}{} // clear the array
		v17v2 = v17va[:len(v17v1):len(v17v1)]
		testUnmarshalErr(&v17v2, bs17, h, t, "dec-slice-v17-p-len")
		testDeepEqualErrHandle(v17v1, v17v2, h, t, "equal-slice-v17-p-len")
		v17va = [8]interface{}{} // clear the array
		v17v2 = v17va[:]
		testUnmarshalErr(&v17v2, bs17, h, t, "dec-slice-v17-p-cap")
		testDeepEqualErrHandle(v17v1, v17v2, h, t, "equal-slice-v17-p-cap")
		if len(v17v1) > 1 {
			v17va = [8]interface{}{} // clear the array
			testUnmarshalErr((&v17va)[:len(v17v1)], bs17, h, t, "dec-slice-v17-p-len-noaddr")
			testDeepEqualErrHandle(v17v1, v17va[:len(v17v1)], h, t, "equal-slice-v17-p-len-noaddr")
			v17va = [8]interface{}{} // clear the array
			testUnmarshalErr((&v17va)[:], bs17, h, t, "dec-slice-v17-p-cap-noaddr")
			testDeepEqualErrHandle(v17v1, v17va[:len(v17v1)], h, t, "equal-slice-v17-p-cap-noaddr")
		}
		testReleaseBytes(bs17)
		// ...
		var v17v3, v17v4 typMbsSliceIntf
		v17v2 = nil
		if v != nil {
			v17v2 = make([]interface{}, len(v))
		}
		v17v3 = typMbsSliceIntf(v17v1)
		v17v4 = typMbsSliceIntf(v17v2)
		if v != nil {
			bs17 = testMarshalErr(v17v3, h, t, "enc-slice-v17-custom")
			testUnmarshalErr(v17v4, bs17, h, t, "dec-slice-v17-custom")
			testDeepEqualErrHandle(v17v3, v17v4, h, t, "equal-slice-v17-custom")
			testReleaseBytes(bs17)
		}
		bs17 = testMarshalErr(&v17v3, h, t, "enc-slice-v17-custom-p")
		v17v2 = nil
		v17v4 = typMbsSliceIntf(v17v2)
		testUnmarshalErr(&v17v4, bs17, h, t, "dec-slice-v17-custom-p")
		testDeepEqualErrHandle(v17v3, v17v4, h, t, "equal-slice-v17-custom-p")
		testReleaseBytes(bs17)
	}
	var v18va [8]string
	for _, v := range [][]string{nil, {}, {"some-string-2", "", "", "some-string-3"}} {
		var v18v1, v18v2 []string
		var bs18 []byte
		v18v1 = v
		bs18 = testMarshalErr(v18v1, h, t, "enc-slice-v18")
		if v == nil {
			v18v2 = make([]string, 2)
			testUnmarshalErr(v18v2, bs18, h, t, "dec-slice-v18")
			testDeepEqualErr(v18v2[0], v18v2[1], t, "equal-slice-v18") // should not change
			testDeepEqualErr(len(v18v2), 2, t, "equal-slice-v18")      // should not change
			v18v2 = make([]string, 2)
			testUnmarshalErr(reflect.ValueOf(v18v2), bs18, h, t, "dec-slice-v18-noaddr") // non-addressable value
			testDeepEqualErr(v18v2[0], v18v2[1], t, "equal-slice-v18-noaddr")            // should not change
			testDeepEqualErr(len(v18v2), 2, t, "equal-slice-v18")                        // should not change
		} else {
			v18v2 = make([]string, len(v))
			testUnmarshalErr(v18v2, bs18, h, t, "dec-slice-v18")
			testDeepEqualErrHandle(v18v1, v18v2, h, t, "equal-slice-v18")
			v18v2 = make([]string, len(v))
			testUnmarshalErr(reflect.ValueOf(v18v2), bs18, h, t, "dec-slice-v18-noaddr") // non-addressable value
			testDeepEqualErrHandle(v18v1, v18v2, h, t, "equal-slice-v18-noaddr")
		}
		testReleaseBytes(bs18)
		// ...
		bs18 = testMarshalErr(&v18v1, h, t, "enc-slice-v18-p")
		v18v2 = nil
		testUnmarshalErr(&v18v2, bs18, h, t, "dec-slice-v18-p")
		testDeepEqualErrHandle(v18v1, v18v2, h, t, "equal-slice-v18-p")
		v18va = [8]string{} // clear the array
		testUnmarshalErr(&v18va, bs18, h, t, "dec-array-v18-p-1")
		if v18v1 == nil && v18v2 == nil {
			v18v2 = []string{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v18va[:len(v18v2)], v18v2, h, t, "equal-array-v18-p-1")
		v18va = [8]string{} // clear the array
		v18v2 = v18va[:1:1]
		testUnmarshalErr(&v18v2, bs18, h, t, "dec-slice-v18-p-1")
		testDeepEqualErrHandle(v18v1, v18v2, h, t, "equal-slice-v18-p-1")
		v18va = [8]string{} // clear the array
		v18v2 = v18va[:len(v18v1):len(v18v1)]
		testUnmarshalErr(&v18v2, bs18, h, t, "dec-slice-v18-p-len")
		testDeepEqualErrHandle(v18v1, v18v2, h, t, "equal-slice-v18-p-len")
		v18va = [8]string{} // clear the array
		v18v2 = v18va[:]
		testUnmarshalErr(&v18v2, bs18, h, t, "dec-slice-v18-p-cap")
		testDeepEqualErrHandle(v18v1, v18v2, h, t, "equal-slice-v18-p-cap")
		if len(v18v1) > 1 {
			v18va = [8]string{} // clear the array
			testUnmarshalErr((&v18va)[:len(v18v1)], bs18, h, t, "dec-slice-v18-p-len-noaddr")
			testDeepEqualErrHandle(v18v1, v18va[:len(v18v1)], h, t, "equal-slice-v18-p-len-noaddr")
			v18va = [8]string{} // clear the array
			testUnmarshalErr((&v18va)[:], bs18, h, t, "dec-slice-v18-p-cap-noaddr")
			testDeepEqualErrHandle(v18v1, v18va[:len(v18v1)], h, t, "equal-slice-v18-p-cap-noaddr")
		}
		testReleaseBytes(bs18)
		// ...
		var v18v3, v18v4 typMbsSliceString
		v18v2 = nil
		if v != nil {
			v18v2 = make([]string, len(v))
		}
		v18v3 = typMbsSliceString(v18v1)
		v18v4 = typMbsSliceString(v18v2)
		if v != nil {
			bs18 = testMarshalErr(v18v3, h, t, "enc-slice-v18-custom")
			testUnmarshalErr(v18v4, bs18, h, t, "dec-slice-v18-custom")
			testDeepEqualErrHandle(v18v3, v18v4, h, t, "equal-slice-v18-custom")
			testReleaseBytes(bs18)
		}
		bs18 = testMarshalErr(&v18v3, h, t, "enc-slice-v18-custom-p")
		v18v2 = nil
		v18v4 = typMbsSliceString(v18v2)
		testUnmarshalErr(&v18v4, bs18, h, t, "dec-slice-v18-custom-p")
		testDeepEqualErrHandle(v18v3, v18v4, h, t, "equal-slice-v18-custom-p")
		testReleaseBytes(bs18)
	}
	var v19va [8][]byte
	for _, v := range [][][]byte{nil, {}, {[]byte("some-string-2"), nil, nil, []byte("some-string-3")}} {
		var v19v1, v19v2 [][]byte
		var bs19 []byte
		v19v1 = v
		bs19 = testMarshalErr(v19v1, h, t, "enc-slice-v19")
		if v == nil {
			v19v2 = make([][]byte, 2)
			testUnmarshalErr(v19v2, bs19, h, t, "dec-slice-v19")
			testDeepEqualErr(v19v2[0], v19v2[1], t, "equal-slice-v19") // should not change
			testDeepEqualErr(len(v19v2), 2, t, "equal-slice-v19")      // should not change
			v19v2 = make([][]byte, 2)
			testUnmarshalErr(reflect.ValueOf(v19v2), bs19, h, t, "dec-slice-v19-noaddr") // non-addressable value
			testDeepEqualErr(v19v2[0], v19v2[1], t, "equal-slice-v19-noaddr")            // should not change
			testDeepEqualErr(len(v19v2), 2, t, "equal-slice-v19")                        // should not change
		} else {
			v19v2 = make([][]byte, len(v))
			testUnmarshalErr(v19v2, bs19, h, t, "dec-slice-v19")
			testDeepEqualErrHandle(v19v1, v19v2, h, t, "equal-slice-v19")
			v19v2 = make([][]byte, len(v))
			testUnmarshalErr(reflect.ValueOf(v19v2), bs19, h, t, "dec-slice-v19-noaddr") // non-addressable value
			testDeepEqualErrHandle(v19v1, v19v2, h, t, "equal-slice-v19-noaddr")
		}
		testReleaseBytes(bs19)
		// ...
		bs19 = testMarshalErr(&v19v1, h, t, "enc-slice-v19-p")
		v19v2 = nil
		testUnmarshalErr(&v19v2, bs19, h, t, "dec-slice-v19-p")
		testDeepEqualErrHandle(v19v1, v19v2, h, t, "equal-slice-v19-p")
		v19va = [8][]byte{} // clear the array
		testUnmarshalErr(&v19va, bs19, h, t, "dec-array-v19-p-1")
		if v19v1 == nil && v19v2 == nil {
			v19v2 = [][]byte{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v19va[:len(v19v2)], v19v2, h, t, "equal-array-v19-p-1")
		v19va = [8][]byte{} // clear the array
		v19v2 = v19va[:1:1]
		testUnmarshalErr(&v19v2, bs19, h, t, "dec-slice-v19-p-1")
		testDeepEqualErrHandle(v19v1, v19v2, h, t, "equal-slice-v19-p-1")
		v19va = [8][]byte{} // clear the array
		v19v2 = v19va[:len(v19v1):len(v19v1)]
		testUnmarshalErr(&v19v2, bs19, h, t, "dec-slice-v19-p-len")
		testDeepEqualErrHandle(v19v1, v19v2, h, t, "equal-slice-v19-p-len")
		v19va = [8][]byte{} // clear the array
		v19v2 = v19va[:]
		testUnmarshalErr(&v19v2, bs19, h, t, "dec-slice-v19-p-cap")
		testDeepEqualErrHandle(v19v1, v19v2, h, t, "equal-slice-v19-p-cap")
		if len(v19v1) > 1 {
			v19va = [8][]byte{} // clear the array
			testUnmarshalErr((&v19va)[:len(v19v1)], bs19, h, t, "dec-slice-v19-p-len-noaddr")
			testDeepEqualErrHandle(v19v1, v19va[:len(v19v1)], h, t, "equal-slice-v19-p-len-noaddr")
			v19va = [8][]byte{} // clear the array
			testUnmarshalErr((&v19va)[:], bs19, h, t, "dec-slice-v19-p-cap-noaddr")
			testDeepEqualErrHandle(v19v1, v19va[:len(v19v1)], h, t, "equal-slice-v19-p-cap-noaddr")
		}
		testReleaseBytes(bs19)
		// ...
		var v19v3, v19v4 typMbsSliceBytes
		v19v2 = nil
		if v != nil {
			v19v2 = make([][]byte, len(v))
		}
		v19v3 = typMbsSliceBytes(v19v1)
		v19v4 = typMbsSliceBytes(v19v2)
		if v != nil {
			bs19 = testMarshalErr(v19v3, h, t, "enc-slice-v19-custom")
			testUnmarshalErr(v19v4, bs19, h, t, "dec-slice-v19-custom")
			testDeepEqualErrHandle(v19v3, v19v4, h, t, "equal-slice-v19-custom")
			testReleaseBytes(bs19)
		}
		bs19 = testMarshalErr(&v19v3, h, t, "enc-slice-v19-custom-p")
		v19v2 = nil
		v19v4 = typMbsSliceBytes(v19v2)
		testUnmarshalErr(&v19v4, bs19, h, t, "dec-slice-v19-custom-p")
		testDeepEqualErrHandle(v19v3, v19v4, h, t, "equal-slice-v19-custom-p")
		testReleaseBytes(bs19)
	}
	var v20va [8]float32
	for _, v := range [][]float32{nil, {}, {22.2, 0, 0, 33.3e3}} {
		var v20v1, v20v2 []float32
		var bs20 []byte
		v20v1 = v
		bs20 = testMarshalErr(v20v1, h, t, "enc-slice-v20")
		if v == nil {
			v20v2 = make([]float32, 2)
			testUnmarshalErr(v20v2, bs20, h, t, "dec-slice-v20")
			testDeepEqualErr(v20v2[0], v20v2[1], t, "equal-slice-v20") // should not change
			testDeepEqualErr(len(v20v2), 2, t, "equal-slice-v20")      // should not change
			v20v2 = make([]float32, 2)
			testUnmarshalErr(reflect.ValueOf(v20v2), bs20, h, t, "dec-slice-v20-noaddr") // non-addressable value
			testDeepEqualErr(v20v2[0], v20v2[1], t, "equal-slice-v20-noaddr")            // should not change
			testDeepEqualErr(len(v20v2), 2, t, "equal-slice-v20")                        // should not change
		} else {
			v20v2 = make([]float32, len(v))
			testUnmarshalErr(v20v2, bs20, h, t, "dec-slice-v20")
			testDeepEqualErrHandle(v20v1, v20v2, h, t, "equal-slice-v20")
			v20v2 = make([]float32, len(v))
			testUnmarshalErr(reflect.ValueOf(v20v2), bs20, h, t, "dec-slice-v20-noaddr") // non-addressable value
			testDeepEqualErrHandle(v20v1, v20v2, h, t, "equal-slice-v20-noaddr")
		}
		testReleaseBytes(bs20)
		// ...
		bs20 = testMarshalErr(&v20v1, h, t, "enc-slice-v20-p")
		v20v2 = nil
		testUnmarshalErr(&v20v2, bs20, h, t, "dec-slice-v20-p")
		testDeepEqualErrHandle(v20v1, v20v2, h, t, "equal-slice-v20-p")
		v20va = [8]float32{} // clear the array
		testUnmarshalErr(&v20va, bs20, h, t, "dec-array-v20-p-1")
		if v20v1 == nil && v20v2 == nil {
			v20v2 = []float32{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v20va[:len(v20v2)], v20v2, h, t, "equal-array-v20-p-1")
		v20va = [8]float32{} // clear the array
		v20v2 = v20va[:1:1]
		testUnmarshalErr(&v20v2, bs20, h, t, "dec-slice-v20-p-1")
		testDeepEqualErrHandle(v20v1, v20v2, h, t, "equal-slice-v20-p-1")
		v20va = [8]float32{} // clear the array
		v20v2 = v20va[:len(v20v1):len(v20v1)]
		testUnmarshalErr(&v20v2, bs20, h, t, "dec-slice-v20-p-len")
		testDeepEqualErrHandle(v20v1, v20v2, h, t, "equal-slice-v20-p-len")
		v20va = [8]float32{} // clear the array
		v20v2 = v20va[:]
		testUnmarshalErr(&v20v2, bs20, h, t, "dec-slice-v20-p-cap")
		testDeepEqualErrHandle(v20v1, v20v2, h, t, "equal-slice-v20-p-cap")
		if len(v20v1) > 1 {
			v20va = [8]float32{} // clear the array
			testUnmarshalErr((&v20va)[:len(v20v1)], bs20, h, t, "dec-slice-v20-p-len-noaddr")
			testDeepEqualErrHandle(v20v1, v20va[:len(v20v1)], h, t, "equal-slice-v20-p-len-noaddr")
			v20va = [8]float32{} // clear the array
			testUnmarshalErr((&v20va)[:], bs20, h, t, "dec-slice-v20-p-cap-noaddr")
			testDeepEqualErrHandle(v20v1, v20va[:len(v20v1)], h, t, "equal-slice-v20-p-cap-noaddr")
		}
		testReleaseBytes(bs20)
		// ...
		var v20v3, v20v4 typMbsSliceFloat32
		v20v2 = nil
		if v != nil {
			v20v2 = make([]float32, len(v))
		}
		v20v3 = typMbsSliceFloat32(v20v1)
		v20v4 = typMbsSliceFloat32(v20v2)
		if v != nil {
			bs20 = testMarshalErr(v20v3, h, t, "enc-slice-v20-custom")
			testUnmarshalErr(v20v4, bs20, h, t, "dec-slice-v20-custom")
			testDeepEqualErrHandle(v20v3, v20v4, h, t, "equal-slice-v20-custom")
			testReleaseBytes(bs20)
		}
		bs20 = testMarshalErr(&v20v3, h, t, "enc-slice-v20-custom-p")
		v20v2 = nil
		v20v4 = typMbsSliceFloat32(v20v2)
		testUnmarshalErr(&v20v4, bs20, h, t, "dec-slice-v20-custom-p")
		testDeepEqualErrHandle(v20v3, v20v4, h, t, "equal-slice-v20-custom-p")
		testReleaseBytes(bs20)
	}
	var v21va [8]float64
	for _, v := range [][]float64{nil, {}, {11.1, 0, 0, 22.2}} {
		var v21v1, v21v2 []float64
		var bs21 []byte
		v21v1 = v
		bs21 = testMarshalErr(v21v1, h, t, "enc-slice-v21")
		if v == nil {
			v21v2 = make([]float64, 2)
			testUnmarshalErr(v21v2, bs21, h, t, "dec-slice-v21")
			testDeepEqualErr(v21v2[0], v21v2[1], t, "equal-slice-v21") // should not change
			testDeepEqualErr(len(v21v2), 2, t, "equal-slice-v21")      // should not change
			v21v2 = make([]float64, 2)
			testUnmarshalErr(reflect.ValueOf(v21v2), bs21, h, t, "dec-slice-v21-noaddr") // non-addressable value
			testDeepEqualErr(v21v2[0], v21v2[1], t, "equal-slice-v21-noaddr")            // should not change
			testDeepEqualErr(len(v21v2), 2, t, "equal-slice-v21")                        // should not change
		} else {
			v21v2 = make([]float64, len(v))
			testUnmarshalErr(v21v2, bs21, h, t, "dec-slice-v21")
			testDeepEqualErrHandle(v21v1, v21v2, h, t, "equal-slice-v21")
			v21v2 = make([]float64, len(v))
			testUnmarshalErr(reflect.ValueOf(v21v2), bs21, h, t, "dec-slice-v21-noaddr") // non-addressable value
			testDeepEqualErrHandle(v21v1, v21v2, h, t, "equal-slice-v21-noaddr")
		}
		testReleaseBytes(bs21)
		// ...
		bs21 = testMarshalErr(&v21v1, h, t, "enc-slice-v21-p")
		v21v2 = nil
		testUnmarshalErr(&v21v2, bs21, h, t, "dec-slice-v21-p")
		testDeepEqualErrHandle(v21v1, v21v2, h, t, "equal-slice-v21-p")
		v21va = [8]float64{} // clear the array
		testUnmarshalErr(&v21va, bs21, h, t, "dec-array-v21-p-1")
		if v21v1 == nil && v21v2 == nil {
			v21v2 = []float64{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v21va[:len(v21v2)], v21v2, h, t, "equal-array-v21-p-1")
		v21va = [8]float64{} // clear the array
		v21v2 = v21va[:1:1]
		testUnmarshalErr(&v21v2, bs21, h, t, "dec-slice-v21-p-1")
		testDeepEqualErrHandle(v21v1, v21v2, h, t, "equal-slice-v21-p-1")
		v21va = [8]float64{} // clear the array
		v21v2 = v21va[:len(v21v1):len(v21v1)]
		testUnmarshalErr(&v21v2, bs21, h, t, "dec-slice-v21-p-len")
		testDeepEqualErrHandle(v21v1, v21v2, h, t, "equal-slice-v21-p-len")
		v21va = [8]float64{} // clear the array
		v21v2 = v21va[:]
		testUnmarshalErr(&v21v2, bs21, h, t, "dec-slice-v21-p-cap")
		testDeepEqualErrHandle(v21v1, v21v2, h, t, "equal-slice-v21-p-cap")
		if len(v21v1) > 1 {
			v21va = [8]float64{} // clear the array
			testUnmarshalErr((&v21va)[:len(v21v1)], bs21, h, t, "dec-slice-v21-p-len-noaddr")
			testDeepEqualErrHandle(v21v1, v21va[:len(v21v1)], h, t, "equal-slice-v21-p-len-noaddr")
			v21va = [8]float64{} // clear the array
			testUnmarshalErr((&v21va)[:], bs21, h, t, "dec-slice-v21-p-cap-noaddr")
			testDeepEqualErrHandle(v21v1, v21va[:len(v21v1)], h, t, "equal-slice-v21-p-cap-noaddr")
		}
		testReleaseBytes(bs21)
		// ...
		var v21v3, v21v4 typMbsSliceFloat64
		v21v2 = nil
		if v != nil {
			v21v2 = make([]float64, len(v))
		}
		v21v3 = typMbsSliceFloat64(v21v1)
		v21v4 = typMbsSliceFloat64(v21v2)
		if v != nil {
			bs21 = testMarshalErr(v21v3, h, t, "enc-slice-v21-custom")
			testUnmarshalErr(v21v4, bs21, h, t, "dec-slice-v21-custom")
			testDeepEqualErrHandle(v21v3, v21v4, h, t, "equal-slice-v21-custom")
			testReleaseBytes(bs21)
		}
		bs21 = testMarshalErr(&v21v3, h, t, "enc-slice-v21-custom-p")
		v21v2 = nil
		v21v4 = typMbsSliceFloat64(v21v2)
		testUnmarshalErr(&v21v4, bs21, h, t, "dec-slice-v21-custom-p")
		testDeepEqualErrHandle(v21v3, v21v4, h, t, "equal-slice-v21-custom-p")
		testReleaseBytes(bs21)
	}
	var v22va [8]uint8
	for _, v := range [][]uint8{nil, {}, {77, 0, 0, 127}} {
		var v22v1, v22v2 []uint8
		var bs22 []byte
		v22v1 = v
		bs22 = testMarshalErr(v22v1, h, t, "enc-slice-v22")
		if v == nil {
			v22v2 = make([]uint8, 2)
			testUnmarshalErr(v22v2, bs22, h, t, "dec-slice-v22")
			testDeepEqualErr(v22v2[0], v22v2[1], t, "equal-slice-v22") // should not change
			testDeepEqualErr(len(v22v2), 2, t, "equal-slice-v22")      // should not change
			v22v2 = make([]uint8, 2)
			testUnmarshalErr(reflect.ValueOf(v22v2), bs22, h, t, "dec-slice-v22-noaddr") // non-addressable value
			testDeepEqualErr(v22v2[0], v22v2[1], t, "equal-slice-v22-noaddr")            // should not change
			testDeepEqualErr(len(v22v2), 2, t, "equal-slice-v22")                        // should not change
		} else {
			v22v2 = make([]uint8, len(v))
			testUnmarshalErr(v22v2, bs22, h, t, "dec-slice-v22")
			testDeepEqualErrHandle(v22v1, v22v2, h, t, "equal-slice-v22")
			v22v2 = make([]uint8, len(v))
			testUnmarshalErr(reflect.ValueOf(v22v2), bs22, h, t, "dec-slice-v22-noaddr") // non-addressable value
			testDeepEqualErrHandle(v22v1, v22v2, h, t, "equal-slice-v22-noaddr")
		}
		testReleaseBytes(bs22)
		// ...
		bs22 = testMarshalErr(&v22v1, h, t, "enc-slice-v22-p")
		v22v2 = nil
		testUnmarshalErr(&v22v2, bs22, h, t, "dec-slice-v22-p")
		testDeepEqualErrHandle(v22v1, v22v2, h, t, "equal-slice-v22-p")
		v22va = [8]uint8{} // clear the array
		testUnmarshalErr(&v22va, bs22, h, t, "dec-array-v22-p-1")
		if v22v1 == nil && v22v2 == nil {
			v22v2 = []uint8{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v22va[:len(v22v2)], v22v2, h, t, "equal-array-v22-p-1")
		v22va = [8]uint8{} // clear the array
		v22v2 = v22va[:1:1]
		testUnmarshalErr(&v22v2, bs22, h, t, "dec-slice-v22-p-1")
		testDeepEqualErrHandle(v22v1, v22v2, h, t, "equal-slice-v22-p-1")
		v22va = [8]uint8{} // clear the array
		v22v2 = v22va[:len(v22v1):len(v22v1)]
		testUnmarshalErr(&v22v2, bs22, h, t, "dec-slice-v22-p-len")
		testDeepEqualErrHandle(v22v1, v22v2, h, t, "equal-slice-v22-p-len")
		v22va = [8]uint8{} // clear the array
		v22v2 = v22va[:]
		testUnmarshalErr(&v22v2, bs22, h, t, "dec-slice-v22-p-cap")
		testDeepEqualErrHandle(v22v1, v22v2, h, t, "equal-slice-v22-p-cap")
		if len(v22v1) > 1 {
			v22va = [8]uint8{} // clear the array
			testUnmarshalErr((&v22va)[:len(v22v1)], bs22, h, t, "dec-slice-v22-p-len-noaddr")
			testDeepEqualErrHandle(v22v1, v22va[:len(v22v1)], h, t, "equal-slice-v22-p-len-noaddr")
			v22va = [8]uint8{} // clear the array
			testUnmarshalErr((&v22va)[:], bs22, h, t, "dec-slice-v22-p-cap-noaddr")
			testDeepEqualErrHandle(v22v1, v22va[:len(v22v1)], h, t, "equal-slice-v22-p-cap-noaddr")
		}
		testReleaseBytes(bs22)
		// ...
		var v22v3, v22v4 typMbsSliceUint8
		v22v2 = nil
		if v != nil {
			v22v2 = make([]uint8, len(v))
		}
		v22v3 = typMbsSliceUint8(v22v1)
		v22v4 = typMbsSliceUint8(v22v2)
		if v != nil {
			bs22 = testMarshalErr(v22v3, h, t, "enc-slice-v22-custom")
			testUnmarshalErr(v22v4, bs22, h, t, "dec-slice-v22-custom")
			testDeepEqualErrHandle(v22v3, v22v4, h, t, "equal-slice-v22-custom")
			testReleaseBytes(bs22)
		}
		bs22 = testMarshalErr(&v22v3, h, t, "enc-slice-v22-custom-p")
		v22v2 = nil
		v22v4 = typMbsSliceUint8(v22v2)
		testUnmarshalErr(&v22v4, bs22, h, t, "dec-slice-v22-custom-p")
		testDeepEqualErrHandle(v22v3, v22v4, h, t, "equal-slice-v22-custom-p")
		testReleaseBytes(bs22)
	}
	var v23va [8]uint64
	for _, v := range [][]uint64{nil, {}, {111, 0, 0, 77}} {
		var v23v1, v23v2 []uint64
		var bs23 []byte
		v23v1 = v
		bs23 = testMarshalErr(v23v1, h, t, "enc-slice-v23")
		if v == nil {
			v23v2 = make([]uint64, 2)
			testUnmarshalErr(v23v2, bs23, h, t, "dec-slice-v23")
			testDeepEqualErr(v23v2[0], v23v2[1], t, "equal-slice-v23") // should not change
			testDeepEqualErr(len(v23v2), 2, t, "equal-slice-v23")      // should not change
			v23v2 = make([]uint64, 2)
			testUnmarshalErr(reflect.ValueOf(v23v2), bs23, h, t, "dec-slice-v23-noaddr") // non-addressable value
			testDeepEqualErr(v23v2[0], v23v2[1], t, "equal-slice-v23-noaddr")            // should not change
			testDeepEqualErr(len(v23v2), 2, t, "equal-slice-v23")                        // should not change
		} else {
			v23v2 = make([]uint64, len(v))
			testUnmarshalErr(v23v2, bs23, h, t, "dec-slice-v23")
			testDeepEqualErrHandle(v23v1, v23v2, h, t, "equal-slice-v23")
			v23v2 = make([]uint64, len(v))
			testUnmarshalErr(reflect.ValueOf(v23v2), bs23, h, t, "dec-slice-v23-noaddr") // non-addressable value
			testDeepEqualErrHandle(v23v1, v23v2, h, t, "equal-slice-v23-noaddr")
		}
		testReleaseBytes(bs23)
		// ...
		bs23 = testMarshalErr(&v23v1, h, t, "enc-slice-v23-p")
		v23v2 = nil
		testUnmarshalErr(&v23v2, bs23, h, t, "dec-slice-v23-p")
		testDeepEqualErrHandle(v23v1, v23v2, h, t, "equal-slice-v23-p")
		v23va = [8]uint64{} // clear the array
		testUnmarshalErr(&v23va, bs23, h, t, "dec-array-v23-p-1")
		if v23v1 == nil && v23v2 == nil {
			v23v2 = []uint64{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v23va[:len(v23v2)], v23v2, h, t, "equal-array-v23-p-1")
		v23va = [8]uint64{} // clear the array
		v23v2 = v23va[:1:1]
		testUnmarshalErr(&v23v2, bs23, h, t, "dec-slice-v23-p-1")
		testDeepEqualErrHandle(v23v1, v23v2, h, t, "equal-slice-v23-p-1")
		v23va = [8]uint64{} // clear the array
		v23v2 = v23va[:len(v23v1):len(v23v1)]
		testUnmarshalErr(&v23v2, bs23, h, t, "dec-slice-v23-p-len")
		testDeepEqualErrHandle(v23v1, v23v2, h, t, "equal-slice-v23-p-len")
		v23va = [8]uint64{} // clear the array
		v23v2 = v23va[:]
		testUnmarshalErr(&v23v2, bs23, h, t, "dec-slice-v23-p-cap")
		testDeepEqualErrHandle(v23v1, v23v2, h, t, "equal-slice-v23-p-cap")
		if len(v23v1) > 1 {
			v23va = [8]uint64{} // clear the array
			testUnmarshalErr((&v23va)[:len(v23v1)], bs23, h, t, "dec-slice-v23-p-len-noaddr")
			testDeepEqualErrHandle(v23v1, v23va[:len(v23v1)], h, t, "equal-slice-v23-p-len-noaddr")
			v23va = [8]uint64{} // clear the array
			testUnmarshalErr((&v23va)[:], bs23, h, t, "dec-slice-v23-p-cap-noaddr")
			testDeepEqualErrHandle(v23v1, v23va[:len(v23v1)], h, t, "equal-slice-v23-p-cap-noaddr")
		}
		testReleaseBytes(bs23)
		// ...
		var v23v3, v23v4 typMbsSliceUint64
		v23v2 = nil
		if v != nil {
			v23v2 = make([]uint64, len(v))
		}
		v23v3 = typMbsSliceUint64(v23v1)
		v23v4 = typMbsSliceUint64(v23v2)
		if v != nil {
			bs23 = testMarshalErr(v23v3, h, t, "enc-slice-v23-custom")
			testUnmarshalErr(v23v4, bs23, h, t, "dec-slice-v23-custom")
			testDeepEqualErrHandle(v23v3, v23v4, h, t, "equal-slice-v23-custom")
			testReleaseBytes(bs23)
		}
		bs23 = testMarshalErr(&v23v3, h, t, "enc-slice-v23-custom-p")
		v23v2 = nil
		v23v4 = typMbsSliceUint64(v23v2)
		testUnmarshalErr(&v23v4, bs23, h, t, "dec-slice-v23-custom-p")
		testDeepEqualErrHandle(v23v3, v23v4, h, t, "equal-slice-v23-custom-p")
		testReleaseBytes(bs23)
	}
	var v24va [8]int
	for _, v := range [][]int{nil, {}, {127, 0, 0, 111}} {
		var v24v1, v24v2 []int
		var bs24 []byte
		v24v1 = v
		bs24 = testMarshalErr(v24v1, h, t, "enc-slice-v24")
		if v == nil {
			v24v2 = make([]int, 2)
			testUnmarshalErr(v24v2, bs24, h, t, "dec-slice-v24")
			testDeepEqualErr(v24v2[0], v24v2[1], t, "equal-slice-v24") // should not change
			testDeepEqualErr(len(v24v2), 2, t, "equal-slice-v24")      // should not change
			v24v2 = make([]int, 2)
			testUnmarshalErr(reflect.ValueOf(v24v2), bs24, h, t, "dec-slice-v24-noaddr") // non-addressable value
			testDeepEqualErr(v24v2[0], v24v2[1], t, "equal-slice-v24-noaddr")            // should not change
			testDeepEqualErr(len(v24v2), 2, t, "equal-slice-v24")                        // should not change
		} else {
			v24v2 = make([]int, len(v))
			testUnmarshalErr(v24v2, bs24, h, t, "dec-slice-v24")
			testDeepEqualErrHandle(v24v1, v24v2, h, t, "equal-slice-v24")
			v24v2 = make([]int, len(v))
			testUnmarshalErr(reflect.ValueOf(v24v2), bs24, h, t, "dec-slice-v24-noaddr") // non-addressable value
			testDeepEqualErrHandle(v24v1, v24v2, h, t, "equal-slice-v24-noaddr")
		}
		testReleaseBytes(bs24)
		// ...
		bs24 = testMarshalErr(&v24v1, h, t, "enc-slice-v24-p")
		v24v2 = nil
		testUnmarshalErr(&v24v2, bs24, h, t, "dec-slice-v24-p")
		testDeepEqualErrHandle(v24v1, v24v2, h, t, "equal-slice-v24-p")
		v24va = [8]int{} // clear the array
		testUnmarshalErr(&v24va, bs24, h, t, "dec-array-v24-p-1")
		if v24v1 == nil && v24v2 == nil {
			v24v2 = []int{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v24va[:len(v24v2)], v24v2, h, t, "equal-array-v24-p-1")
		v24va = [8]int{} // clear the array
		v24v2 = v24va[:1:1]
		testUnmarshalErr(&v24v2, bs24, h, t, "dec-slice-v24-p-1")
		testDeepEqualErrHandle(v24v1, v24v2, h, t, "equal-slice-v24-p-1")
		v24va = [8]int{} // clear the array
		v24v2 = v24va[:len(v24v1):len(v24v1)]
		testUnmarshalErr(&v24v2, bs24, h, t, "dec-slice-v24-p-len")
		testDeepEqualErrHandle(v24v1, v24v2, h, t, "equal-slice-v24-p-len")
		v24va = [8]int{} // clear the array
		v24v2 = v24va[:]
		testUnmarshalErr(&v24v2, bs24, h, t, "dec-slice-v24-p-cap")
		testDeepEqualErrHandle(v24v1, v24v2, h, t, "equal-slice-v24-p-cap")
		if len(v24v1) > 1 {
			v24va = [8]int{} // clear the array
			testUnmarshalErr((&v24va)[:len(v24v1)], bs24, h, t, "dec-slice-v24-p-len-noaddr")
			testDeepEqualErrHandle(v24v1, v24va[:len(v24v1)], h, t, "equal-slice-v24-p-len-noaddr")
			v24va = [8]int{} // clear the array
			testUnmarshalErr((&v24va)[:], bs24, h, t, "dec-slice-v24-p-cap-noaddr")
			testDeepEqualErrHandle(v24v1, v24va[:len(v24v1)], h, t, "equal-slice-v24-p-cap-noaddr")
		}
		testReleaseBytes(bs24)
		// ...
		var v24v3, v24v4 typMbsSliceInt
		v24v2 = nil
		if v != nil {
			v24v2 = make([]int, len(v))
		}
		v24v3 = typMbsSliceInt(v24v1)
		v24v4 = typMbsSliceInt(v24v2)
		if v != nil {
			bs24 = testMarshalErr(v24v3, h, t, "enc-slice-v24-custom")
			testUnmarshalErr(v24v4, bs24, h, t, "dec-slice-v24-custom")
			testDeepEqualErrHandle(v24v3, v24v4, h, t, "equal-slice-v24-custom")
			testReleaseBytes(bs24)
		}
		bs24 = testMarshalErr(&v24v3, h, t, "enc-slice-v24-custom-p")
		v24v2 = nil
		v24v4 = typMbsSliceInt(v24v2)
		testUnmarshalErr(&v24v4, bs24, h, t, "dec-slice-v24-custom-p")
		testDeepEqualErrHandle(v24v3, v24v4, h, t, "equal-slice-v24-custom-p")
		testReleaseBytes(bs24)
	}
	var v25va [8]int32
	for _, v := range [][]int32{nil, {}, {77, 0, 0, 127}} {
		var v25v1, v25v2 []int32
		var bs25 []byte
		v25v1 = v
		bs25 = testMarshalErr(v25v1, h, t, "enc-slice-v25")
		if v == nil {
			v25v2 = make([]int32, 2)
			testUnmarshalErr(v25v2, bs25, h, t, "dec-slice-v25")
			testDeepEqualErr(v25v2[0], v25v2[1], t, "equal-slice-v25") // should not change
			testDeepEqualErr(len(v25v2), 2, t, "equal-slice-v25")      // should not change
			v25v2 = make([]int32, 2)
			testUnmarshalErr(reflect.ValueOf(v25v2), bs25, h, t, "dec-slice-v25-noaddr") // non-addressable value
			testDeepEqualErr(v25v2[0], v25v2[1], t, "equal-slice-v25-noaddr")            // should not change
			testDeepEqualErr(len(v25v2), 2, t, "equal-slice-v25")                        // should not change
		} else {
			v25v2 = make([]int32, len(v))
			testUnmarshalErr(v25v2, bs25, h, t, "dec-slice-v25")
			testDeepEqualErrHandle(v25v1, v25v2, h, t, "equal-slice-v25")
			v25v2 = make([]int32, len(v))
			testUnmarshalErr(reflect.ValueOf(v25v2), bs25, h, t, "dec-slice-v25-noaddr") // non-addressable value
			testDeepEqualErrHandle(v25v1, v25v2, h, t, "equal-slice-v25-noaddr")
		}
		testReleaseBytes(bs25)
		// ...
		bs25 = testMarshalErr(&v25v1, h, t, "enc-slice-v25-p")
		v25v2 = nil
		testUnmarshalErr(&v25v2, bs25, h, t, "dec-slice-v25-p")
		testDeepEqualErrHandle(v25v1, v25v2, h, t, "equal-slice-v25-p")
		v25va = [8]int32{} // clear the array
		testUnmarshalErr(&v25va, bs25, h, t, "dec-array-v25-p-1")
		if v25v1 == nil && v25v2 == nil {
			v25v2 = []int32{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v25va[:len(v25v2)], v25v2, h, t, "equal-array-v25-p-1")
		v25va = [8]int32{} // clear the array
		v25v2 = v25va[:1:1]
		testUnmarshalErr(&v25v2, bs25, h, t, "dec-slice-v25-p-1")
		testDeepEqualErrHandle(v25v1, v25v2, h, t, "equal-slice-v25-p-1")
		v25va = [8]int32{} // clear the array
		v25v2 = v25va[:len(v25v1):len(v25v1)]
		testUnmarshalErr(&v25v2, bs25, h, t, "dec-slice-v25-p-len")
		testDeepEqualErrHandle(v25v1, v25v2, h, t, "equal-slice-v25-p-len")
		v25va = [8]int32{} // clear the array
		v25v2 = v25va[:]
		testUnmarshalErr(&v25v2, bs25, h, t, "dec-slice-v25-p-cap")
		testDeepEqualErrHandle(v25v1, v25v2, h, t, "equal-slice-v25-p-cap")
		if len(v25v1) > 1 {
			v25va = [8]int32{} // clear the array
			testUnmarshalErr((&v25va)[:len(v25v1)], bs25, h, t, "dec-slice-v25-p-len-noaddr")
			testDeepEqualErrHandle(v25v1, v25va[:len(v25v1)], h, t, "equal-slice-v25-p-len-noaddr")
			v25va = [8]int32{} // clear the array
			testUnmarshalErr((&v25va)[:], bs25, h, t, "dec-slice-v25-p-cap-noaddr")
			testDeepEqualErrHandle(v25v1, v25va[:len(v25v1)], h, t, "equal-slice-v25-p-cap-noaddr")
		}
		testReleaseBytes(bs25)
		// ...
		var v25v3, v25v4 typMbsSliceInt32
		v25v2 = nil
		if v != nil {
			v25v2 = make([]int32, len(v))
		}
		v25v3 = typMbsSliceInt32(v25v1)
		v25v4 = typMbsSliceInt32(v25v2)
		if v != nil {
			bs25 = testMarshalErr(v25v3, h, t, "enc-slice-v25-custom")
			testUnmarshalErr(v25v4, bs25, h, t, "dec-slice-v25-custom")
			testDeepEqualErrHandle(v25v3, v25v4, h, t, "equal-slice-v25-custom")
			testReleaseBytes(bs25)
		}
		bs25 = testMarshalErr(&v25v3, h, t, "enc-slice-v25-custom-p")
		v25v2 = nil
		v25v4 = typMbsSliceInt32(v25v2)
		testUnmarshalErr(&v25v4, bs25, h, t, "dec-slice-v25-custom-p")
		testDeepEqualErrHandle(v25v3, v25v4, h, t, "equal-slice-v25-custom-p")
		testReleaseBytes(bs25)
	}
	var v26va [8]int64
	for _, v := range [][]int64{nil, {}, {111, 0, 0, 77}} {
		var v26v1, v26v2 []int64
		var bs26 []byte
		v26v1 = v
		bs26 = testMarshalErr(v26v1, h, t, "enc-slice-v26")
		if v == nil {
			v26v2 = make([]int64, 2)
			testUnmarshalErr(v26v2, bs26, h, t, "dec-slice-v26")
			testDeepEqualErr(v26v2[0], v26v2[1], t, "equal-slice-v26") // should not change
			testDeepEqualErr(len(v26v2), 2, t, "equal-slice-v26")      // should not change
			v26v2 = make([]int64, 2)
			testUnmarshalErr(reflect.ValueOf(v26v2), bs26, h, t, "dec-slice-v26-noaddr") // non-addressable value
			testDeepEqualErr(v26v2[0], v26v2[1], t, "equal-slice-v26-noaddr")            // should not change
			testDeepEqualErr(len(v26v2), 2, t, "equal-slice-v26")                        // should not change
		} else {
			v26v2 = make([]int64, len(v))
			testUnmarshalErr(v26v2, bs26, h, t, "dec-slice-v26")
			testDeepEqualErrHandle(v26v1, v26v2, h, t, "equal-slice-v26")
			v26v2 = make([]int64, len(v))
			testUnmarshalErr(reflect.ValueOf(v26v2), bs26, h, t, "dec-slice-v26-noaddr") // non-addressable value
			testDeepEqualErrHandle(v26v1, v26v2, h, t, "equal-slice-v26-noaddr")
		}
		testReleaseBytes(bs26)
		// ...
		bs26 = testMarshalErr(&v26v1, h, t, "enc-slice-v26-p")
		v26v2 = nil
		testUnmarshalErr(&v26v2, bs26, h, t, "dec-slice-v26-p")
		testDeepEqualErrHandle(v26v1, v26v2, h, t, "equal-slice-v26-p")
		v26va = [8]int64{} // clear the array
		testUnmarshalErr(&v26va, bs26, h, t, "dec-array-v26-p-1")
		if v26v1 == nil && v26v2 == nil {
			v26v2 = []int64{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v26va[:len(v26v2)], v26v2, h, t, "equal-array-v26-p-1")
		v26va = [8]int64{} // clear the array
		v26v2 = v26va[:1:1]
		testUnmarshalErr(&v26v2, bs26, h, t, "dec-slice-v26-p-1")
		testDeepEqualErrHandle(v26v1, v26v2, h, t, "equal-slice-v26-p-1")
		v26va = [8]int64{} // clear the array
		v26v2 = v26va[:len(v26v1):len(v26v1)]
		testUnmarshalErr(&v26v2, bs26, h, t, "dec-slice-v26-p-len")
		testDeepEqualErrHandle(v26v1, v26v2, h, t, "equal-slice-v26-p-len")
		v26va = [8]int64{} // clear the array
		v26v2 = v26va[:]
		testUnmarshalErr(&v26v2, bs26, h, t, "dec-slice-v26-p-cap")
		testDeepEqualErrHandle(v26v1, v26v2, h, t, "equal-slice-v26-p-cap")
		if len(v26v1) > 1 {
			v26va = [8]int64{} // clear the array
			testUnmarshalErr((&v26va)[:len(v26v1)], bs26, h, t, "dec-slice-v26-p-len-noaddr")
			testDeepEqualErrHandle(v26v1, v26va[:len(v26v1)], h, t, "equal-slice-v26-p-len-noaddr")
			v26va = [8]int64{} // clear the array
			testUnmarshalErr((&v26va)[:], bs26, h, t, "dec-slice-v26-p-cap-noaddr")
			testDeepEqualErrHandle(v26v1, v26va[:len(v26v1)], h, t, "equal-slice-v26-p-cap-noaddr")
		}
		testReleaseBytes(bs26)
		// ...
		var v26v3, v26v4 typMbsSliceInt64
		v26v2 = nil
		if v != nil {
			v26v2 = make([]int64, len(v))
		}
		v26v3 = typMbsSliceInt64(v26v1)
		v26v4 = typMbsSliceInt64(v26v2)
		if v != nil {
			bs26 = testMarshalErr(v26v3, h, t, "enc-slice-v26-custom")
			testUnmarshalErr(v26v4, bs26, h, t, "dec-slice-v26-custom")
			testDeepEqualErrHandle(v26v3, v26v4, h, t, "equal-slice-v26-custom")
			testReleaseBytes(bs26)
		}
		bs26 = testMarshalErr(&v26v3, h, t, "enc-slice-v26-custom-p")
		v26v2 = nil
		v26v4 = typMbsSliceInt64(v26v2)
		testUnmarshalErr(&v26v4, bs26, h, t, "dec-slice-v26-custom-p")
		testDeepEqualErrHandle(v26v3, v26v4, h, t, "equal-slice-v26-custom-p")
		testReleaseBytes(bs26)
	}
	var v27va [8]bool
	for _, v := range [][]bool{nil, {}, {false, false, false, true}} {
		var v27v1, v27v2 []bool
		var bs27 []byte
		v27v1 = v
		bs27 = testMarshalErr(v27v1, h, t, "enc-slice-v27")
		if v == nil {
			v27v2 = make([]bool, 2)
			testUnmarshalErr(v27v2, bs27, h, t, "dec-slice-v27")
			testDeepEqualErr(v27v2[0], v27v2[1], t, "equal-slice-v27") // should not change
			testDeepEqualErr(len(v27v2), 2, t, "equal-slice-v27")      // should not change
			v27v2 = make([]bool, 2)
			testUnmarshalErr(reflect.ValueOf(v27v2), bs27, h, t, "dec-slice-v27-noaddr") // non-addressable value
			testDeepEqualErr(v27v2[0], v27v2[1], t, "equal-slice-v27-noaddr")            // should not change
			testDeepEqualErr(len(v27v2), 2, t, "equal-slice-v27")                        // should not change
		} else {
			v27v2 = make([]bool, len(v))
			testUnmarshalErr(v27v2, bs27, h, t, "dec-slice-v27")
			testDeepEqualErrHandle(v27v1, v27v2, h, t, "equal-slice-v27")
			v27v2 = make([]bool, len(v))
			testUnmarshalErr(reflect.ValueOf(v27v2), bs27, h, t, "dec-slice-v27-noaddr") // non-addressable value
			testDeepEqualErrHandle(v27v1, v27v2, h, t, "equal-slice-v27-noaddr")
		}
		testReleaseBytes(bs27)
		// ...
		bs27 = testMarshalErr(&v27v1, h, t, "enc-slice-v27-p")
		v27v2 = nil
		testUnmarshalErr(&v27v2, bs27, h, t, "dec-slice-v27-p")
		testDeepEqualErrHandle(v27v1, v27v2, h, t, "equal-slice-v27-p")
		v27va = [8]bool{} // clear the array
		testUnmarshalErr(&v27va, bs27, h, t, "dec-array-v27-p-1")
		if v27v1 == nil && v27v2 == nil {
			v27v2 = []bool{}
		} // so we can compare to zero len slice below
		testDeepEqualErrHandle(v27va[:len(v27v2)], v27v2, h, t, "equal-array-v27-p-1")
		v27va = [8]bool{} // clear the array
		v27v2 = v27va[:1:1]
		testUnmarshalErr(&v27v2, bs27, h, t, "dec-slice-v27-p-1")
		testDeepEqualErrHandle(v27v1, v27v2, h, t, "equal-slice-v27-p-1")
		v27va = [8]bool{} // clear the array
		v27v2 = v27va[:len(v27v1):len(v27v1)]
		testUnmarshalErr(&v27v2, bs27, h, t, "dec-slice-v27-p-len")
		testDeepEqualErrHandle(v27v1, v27v2, h, t, "equal-slice-v27-p-len")
		v27va = [8]bool{} // clear the array
		v27v2 = v27va[:]
		testUnmarshalErr(&v27v2, bs27, h, t, "dec-slice-v27-p-cap")
		testDeepEqualErrHandle(v27v1, v27v2, h, t, "equal-slice-v27-p-cap")
		if len(v27v1) > 1 {
			v27va = [8]bool{} // clear the array
			testUnmarshalErr((&v27va)[:len(v27v1)], bs27, h, t, "dec-slice-v27-p-len-noaddr")
			testDeepEqualErrHandle(v27v1, v27va[:len(v27v1)], h, t, "equal-slice-v27-p-len-noaddr")
			v27va = [8]bool{} // clear the array
			testUnmarshalErr((&v27va)[:], bs27, h, t, "dec-slice-v27-p-cap-noaddr")
			testDeepEqualErrHandle(v27v1, v27va[:len(v27v1)], h, t, "equal-slice-v27-p-cap-noaddr")
		}
		testReleaseBytes(bs27)
		// ...
		var v27v3, v27v4 typMbsSliceBool
		v27v2 = nil
		if v != nil {
			v27v2 = make([]bool, len(v))
		}
		v27v3 = typMbsSliceBool(v27v1)
		v27v4 = typMbsSliceBool(v27v2)
		if v != nil {
			bs27 = testMarshalErr(v27v3, h, t, "enc-slice-v27-custom")
			testUnmarshalErr(v27v4, bs27, h, t, "dec-slice-v27-custom")
			testDeepEqualErrHandle(v27v3, v27v4, h, t, "equal-slice-v27-custom")
			testReleaseBytes(bs27)
		}
		bs27 = testMarshalErr(&v27v3, h, t, "enc-slice-v27-custom-p")
		v27v2 = nil
		v27v4 = typMbsSliceBool(v27v2)
		testUnmarshalErr(&v27v4, bs27, h, t, "dec-slice-v27-custom-p")
		testDeepEqualErrHandle(v27v3, v27v4, h, t, "equal-slice-v27-custom-p")
		testReleaseBytes(bs27)
	}

}

func __doTestMammothMaps(t *testing.T, h Handle) {
	for _, v := range []map[string]interface{}{nil, {}, {"some-string-1": nil, "some-string-2": "string-is-an-interface-1"}} {
		var v28v1, v28v2 map[string]interface{}
		var bs28 []byte
		v28v1 = v
		bs28 = testMarshalErr(v28v1, h, t, "enc-map-v28")
		if v != nil {
			v28v2 = make(map[string]interface{}, len(v)) // reset map
			testUnmarshalErr(v28v2, bs28, h, t, "dec-map-v28")
			testDeepEqualErrHandle(v28v1, v28v2, h, t, "equal-map-v28")
			v28v2 = make(map[string]interface{}, len(v))                               // reset map
			testUnmarshalErr(reflect.ValueOf(v28v2), bs28, h, t, "dec-map-v28-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v28v1, v28v2, h, t, "equal-map-v28-noaddr")
		}
		if v == nil {
			v28v2 = nil
		} else {
			v28v2 = make(map[string]interface{}, len(v))
		} // reset map
		testUnmarshalErr(&v28v2, bs28, h, t, "dec-map-v28-p-len")
		testDeepEqualErrHandle(v28v1, v28v2, h, t, "equal-map-v28-p-len")
		testReleaseBytes(bs28)
		bs28 = testMarshalErr(&v28v1, h, t, "enc-map-v28-p")
		v28v2 = nil
		testUnmarshalErr(&v28v2, bs28, h, t, "dec-map-v28-p-nil")
		testDeepEqualErrHandle(v28v1, v28v2, h, t, "equal-map-v28-p-nil")
		testReleaseBytes(bs28)
		// ...
		if v == nil {
			v28v2 = nil
		} else {
			v28v2 = make(map[string]interface{}, len(v))
		} // reset map
		var v28v3, v28v4 typMapMapStringIntf
		v28v3 = typMapMapStringIntf(v28v1)
		v28v4 = typMapMapStringIntf(v28v2)
		if v != nil {
			bs28 = testMarshalErr(v28v3, h, t, "enc-map-v28-custom")
			testUnmarshalErr(v28v4, bs28, h, t, "dec-map-v28-p-len")
			testDeepEqualErrHandle(v28v3, v28v4, h, t, "equal-map-v28-p-len")
			testReleaseBytes(bs28)
		}
		type s28T struct {
			M  map[string]interface{}
			Mp *map[string]interface{}
		}
		var m28v99 = map[string]interface{}{
			"":              nil,
			"some-string-3": "string-is-an-interface-2",
		}
		var s28v1, s28v2 s28T
		bs28 = testMarshalErr(s28v1, h, t, "enc-map-v28-custom")
		testUnmarshalErr(&s28v2, bs28, h, t, "dec-map-v28-p-len")
		testDeepEqualErrHandle(s28v1, s28v2, h, t, "equal-map-v28-p-len")
		testReleaseBytes(bs28)
		s28v2 = s28T{}
		s28v1.M = m28v99
		bs28 = testMarshalErr(s28v1, h, t, "enc-map-v28-custom")
		testUnmarshalErr(&s28v2, bs28, h, t, "dec-map-v28-p-len")
		testDeepEqualErrHandle(s28v1, s28v2, h, t, "equal-map-v28-p-len")
		testReleaseBytes(bs28)
		s28v2 = s28T{}
		s28v1.Mp = &m28v99
		bs28 = testMarshalErr(s28v1, h, t, "enc-map-v28-custom")
		testUnmarshalErr(&s28v2, bs28, h, t, "dec-map-v28-p-len")
		testDeepEqualErrHandle(s28v1, s28v2, h, t, "equal-map-v28-p-len")
		testReleaseBytes(bs28)
	}
	for _, v := range []map[string]string{nil, {}, {"some-string-1": "", "some-string-2": "some-string-3"}} {
		var v29v1, v29v2 map[string]string
		var bs29 []byte
		v29v1 = v
		bs29 = testMarshalErr(v29v1, h, t, "enc-map-v29")
		if v != nil {
			v29v2 = make(map[string]string, len(v)) // reset map
			testUnmarshalErr(v29v2, bs29, h, t, "dec-map-v29")
			testDeepEqualErrHandle(v29v1, v29v2, h, t, "equal-map-v29")
			v29v2 = make(map[string]string, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v29v2), bs29, h, t, "dec-map-v29-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v29v1, v29v2, h, t, "equal-map-v29-noaddr")
		}
		if v == nil {
			v29v2 = nil
		} else {
			v29v2 = make(map[string]string, len(v))
		} // reset map
		testUnmarshalErr(&v29v2, bs29, h, t, "dec-map-v29-p-len")
		testDeepEqualErrHandle(v29v1, v29v2, h, t, "equal-map-v29-p-len")
		testReleaseBytes(bs29)
		bs29 = testMarshalErr(&v29v1, h, t, "enc-map-v29-p")
		v29v2 = nil
		testUnmarshalErr(&v29v2, bs29, h, t, "dec-map-v29-p-nil")
		testDeepEqualErrHandle(v29v1, v29v2, h, t, "equal-map-v29-p-nil")
		testReleaseBytes(bs29)
		// ...
		if v == nil {
			v29v2 = nil
		} else {
			v29v2 = make(map[string]string, len(v))
		} // reset map
		var v29v3, v29v4 typMapMapStringString
		v29v3 = typMapMapStringString(v29v1)
		v29v4 = typMapMapStringString(v29v2)
		if v != nil {
			bs29 = testMarshalErr(v29v3, h, t, "enc-map-v29-custom")
			testUnmarshalErr(v29v4, bs29, h, t, "dec-map-v29-p-len")
			testDeepEqualErrHandle(v29v3, v29v4, h, t, "equal-map-v29-p-len")
			testReleaseBytes(bs29)
		}
		type s29T struct {
			M  map[string]string
			Mp *map[string]string
		}
		var m29v99 = map[string]string{
			"":              "",
			"some-string-1": "some-string-2",
		}
		var s29v1, s29v2 s29T
		bs29 = testMarshalErr(s29v1, h, t, "enc-map-v29-custom")
		testUnmarshalErr(&s29v2, bs29, h, t, "dec-map-v29-p-len")
		testDeepEqualErrHandle(s29v1, s29v2, h, t, "equal-map-v29-p-len")
		testReleaseBytes(bs29)
		s29v2 = s29T{}
		s29v1.M = m29v99
		bs29 = testMarshalErr(s29v1, h, t, "enc-map-v29-custom")
		testUnmarshalErr(&s29v2, bs29, h, t, "dec-map-v29-p-len")
		testDeepEqualErrHandle(s29v1, s29v2, h, t, "equal-map-v29-p-len")
		testReleaseBytes(bs29)
		s29v2 = s29T{}
		s29v1.Mp = &m29v99
		bs29 = testMarshalErr(s29v1, h, t, "enc-map-v29-custom")
		testUnmarshalErr(&s29v2, bs29, h, t, "dec-map-v29-p-len")
		testDeepEqualErrHandle(s29v1, s29v2, h, t, "equal-map-v29-p-len")
		testReleaseBytes(bs29)
	}
	for _, v := range []map[string][]byte{nil, {}, {"some-string-3": nil, "some-string-1": []byte("some-string-1")}} {
		var v30v1, v30v2 map[string][]byte
		var bs30 []byte
		v30v1 = v
		bs30 = testMarshalErr(v30v1, h, t, "enc-map-v30")
		if v != nil {
			v30v2 = make(map[string][]byte, len(v)) // reset map
			testUnmarshalErr(v30v2, bs30, h, t, "dec-map-v30")
			testDeepEqualErrHandle(v30v1, v30v2, h, t, "equal-map-v30")
			v30v2 = make(map[string][]byte, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v30v2), bs30, h, t, "dec-map-v30-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v30v1, v30v2, h, t, "equal-map-v30-noaddr")
		}
		if v == nil {
			v30v2 = nil
		} else {
			v30v2 = make(map[string][]byte, len(v))
		} // reset map
		testUnmarshalErr(&v30v2, bs30, h, t, "dec-map-v30-p-len")
		testDeepEqualErrHandle(v30v1, v30v2, h, t, "equal-map-v30-p-len")
		testReleaseBytes(bs30)
		bs30 = testMarshalErr(&v30v1, h, t, "enc-map-v30-p")
		v30v2 = nil
		testUnmarshalErr(&v30v2, bs30, h, t, "dec-map-v30-p-nil")
		testDeepEqualErrHandle(v30v1, v30v2, h, t, "equal-map-v30-p-nil")
		testReleaseBytes(bs30)
		// ...
		if v == nil {
			v30v2 = nil
		} else {
			v30v2 = make(map[string][]byte, len(v))
		} // reset map
		var v30v3, v30v4 typMapMapStringBytes
		v30v3 = typMapMapStringBytes(v30v1)
		v30v4 = typMapMapStringBytes(v30v2)
		if v != nil {
			bs30 = testMarshalErr(v30v3, h, t, "enc-map-v30-custom")
			testUnmarshalErr(v30v4, bs30, h, t, "dec-map-v30-p-len")
			testDeepEqualErrHandle(v30v3, v30v4, h, t, "equal-map-v30-p-len")
			testReleaseBytes(bs30)
		}
		type s30T struct {
			M  map[string][]byte
			Mp *map[string][]byte
		}
		var m30v99 = map[string][]byte{
			"":              nil,
			"some-string-2": []byte("some-string-2"),
		}
		var s30v1, s30v2 s30T
		bs30 = testMarshalErr(s30v1, h, t, "enc-map-v30-custom")
		testUnmarshalErr(&s30v2, bs30, h, t, "dec-map-v30-p-len")
		testDeepEqualErrHandle(s30v1, s30v2, h, t, "equal-map-v30-p-len")
		testReleaseBytes(bs30)
		s30v2 = s30T{}
		s30v1.M = m30v99
		bs30 = testMarshalErr(s30v1, h, t, "enc-map-v30-custom")
		testUnmarshalErr(&s30v2, bs30, h, t, "dec-map-v30-p-len")
		testDeepEqualErrHandle(s30v1, s30v2, h, t, "equal-map-v30-p-len")
		testReleaseBytes(bs30)
		s30v2 = s30T{}
		s30v1.Mp = &m30v99
		bs30 = testMarshalErr(s30v1, h, t, "enc-map-v30-custom")
		testUnmarshalErr(&s30v2, bs30, h, t, "dec-map-v30-p-len")
		testDeepEqualErrHandle(s30v1, s30v2, h, t, "equal-map-v30-p-len")
		testReleaseBytes(bs30)
	}
	for _, v := range []map[string]uint8{nil, {}, {"some-string-3": 0, "some-string-1": 127}} {
		var v31v1, v31v2 map[string]uint8
		var bs31 []byte
		v31v1 = v
		bs31 = testMarshalErr(v31v1, h, t, "enc-map-v31")
		if v != nil {
			v31v2 = make(map[string]uint8, len(v)) // reset map
			testUnmarshalErr(v31v2, bs31, h, t, "dec-map-v31")
			testDeepEqualErrHandle(v31v1, v31v2, h, t, "equal-map-v31")
			v31v2 = make(map[string]uint8, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v31v2), bs31, h, t, "dec-map-v31-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v31v1, v31v2, h, t, "equal-map-v31-noaddr")
		}
		if v == nil {
			v31v2 = nil
		} else {
			v31v2 = make(map[string]uint8, len(v))
		} // reset map
		testUnmarshalErr(&v31v2, bs31, h, t, "dec-map-v31-p-len")
		testDeepEqualErrHandle(v31v1, v31v2, h, t, "equal-map-v31-p-len")
		testReleaseBytes(bs31)
		bs31 = testMarshalErr(&v31v1, h, t, "enc-map-v31-p")
		v31v2 = nil
		testUnmarshalErr(&v31v2, bs31, h, t, "dec-map-v31-p-nil")
		testDeepEqualErrHandle(v31v1, v31v2, h, t, "equal-map-v31-p-nil")
		testReleaseBytes(bs31)
		// ...
		if v == nil {
			v31v2 = nil
		} else {
			v31v2 = make(map[string]uint8, len(v))
		} // reset map
		var v31v3, v31v4 typMapMapStringUint8
		v31v3 = typMapMapStringUint8(v31v1)
		v31v4 = typMapMapStringUint8(v31v2)
		if v != nil {
			bs31 = testMarshalErr(v31v3, h, t, "enc-map-v31-custom")
			testUnmarshalErr(v31v4, bs31, h, t, "dec-map-v31-p-len")
			testDeepEqualErrHandle(v31v3, v31v4, h, t, "equal-map-v31-p-len")
			testReleaseBytes(bs31)
		}
		type s31T struct {
			M  map[string]uint8
			Mp *map[string]uint8
		}
		var m31v99 = map[string]uint8{
			"":              0,
			"some-string-2": 111,
		}
		var s31v1, s31v2 s31T
		bs31 = testMarshalErr(s31v1, h, t, "enc-map-v31-custom")
		testUnmarshalErr(&s31v2, bs31, h, t, "dec-map-v31-p-len")
		testDeepEqualErrHandle(s31v1, s31v2, h, t, "equal-map-v31-p-len")
		testReleaseBytes(bs31)
		s31v2 = s31T{}
		s31v1.M = m31v99
		bs31 = testMarshalErr(s31v1, h, t, "enc-map-v31-custom")
		testUnmarshalErr(&s31v2, bs31, h, t, "dec-map-v31-p-len")
		testDeepEqualErrHandle(s31v1, s31v2, h, t, "equal-map-v31-p-len")
		testReleaseBytes(bs31)
		s31v2 = s31T{}
		s31v1.Mp = &m31v99
		bs31 = testMarshalErr(s31v1, h, t, "enc-map-v31-custom")
		testUnmarshalErr(&s31v2, bs31, h, t, "dec-map-v31-p-len")
		testDeepEqualErrHandle(s31v1, s31v2, h, t, "equal-map-v31-p-len")
		testReleaseBytes(bs31)
	}
	for _, v := range []map[string]uint64{nil, {}, {"some-string-3": 0, "some-string-1": 77}} {
		var v32v1, v32v2 map[string]uint64
		var bs32 []byte
		v32v1 = v
		bs32 = testMarshalErr(v32v1, h, t, "enc-map-v32")
		if v != nil {
			v32v2 = make(map[string]uint64, len(v)) // reset map
			testUnmarshalErr(v32v2, bs32, h, t, "dec-map-v32")
			testDeepEqualErrHandle(v32v1, v32v2, h, t, "equal-map-v32")
			v32v2 = make(map[string]uint64, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v32v2), bs32, h, t, "dec-map-v32-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v32v1, v32v2, h, t, "equal-map-v32-noaddr")
		}
		if v == nil {
			v32v2 = nil
		} else {
			v32v2 = make(map[string]uint64, len(v))
		} // reset map
		testUnmarshalErr(&v32v2, bs32, h, t, "dec-map-v32-p-len")
		testDeepEqualErrHandle(v32v1, v32v2, h, t, "equal-map-v32-p-len")
		testReleaseBytes(bs32)
		bs32 = testMarshalErr(&v32v1, h, t, "enc-map-v32-p")
		v32v2 = nil
		testUnmarshalErr(&v32v2, bs32, h, t, "dec-map-v32-p-nil")
		testDeepEqualErrHandle(v32v1, v32v2, h, t, "equal-map-v32-p-nil")
		testReleaseBytes(bs32)
		// ...
		if v == nil {
			v32v2 = nil
		} else {
			v32v2 = make(map[string]uint64, len(v))
		} // reset map
		var v32v3, v32v4 typMapMapStringUint64
		v32v3 = typMapMapStringUint64(v32v1)
		v32v4 = typMapMapStringUint64(v32v2)
		if v != nil {
			bs32 = testMarshalErr(v32v3, h, t, "enc-map-v32-custom")
			testUnmarshalErr(v32v4, bs32, h, t, "dec-map-v32-p-len")
			testDeepEqualErrHandle(v32v3, v32v4, h, t, "equal-map-v32-p-len")
			testReleaseBytes(bs32)
		}
		type s32T struct {
			M  map[string]uint64
			Mp *map[string]uint64
		}
		var m32v99 = map[string]uint64{
			"":              0,
			"some-string-2": 127,
		}
		var s32v1, s32v2 s32T
		bs32 = testMarshalErr(s32v1, h, t, "enc-map-v32-custom")
		testUnmarshalErr(&s32v2, bs32, h, t, "dec-map-v32-p-len")
		testDeepEqualErrHandle(s32v1, s32v2, h, t, "equal-map-v32-p-len")
		testReleaseBytes(bs32)
		s32v2 = s32T{}
		s32v1.M = m32v99
		bs32 = testMarshalErr(s32v1, h, t, "enc-map-v32-custom")
		testUnmarshalErr(&s32v2, bs32, h, t, "dec-map-v32-p-len")
		testDeepEqualErrHandle(s32v1, s32v2, h, t, "equal-map-v32-p-len")
		testReleaseBytes(bs32)
		s32v2 = s32T{}
		s32v1.Mp = &m32v99
		bs32 = testMarshalErr(s32v1, h, t, "enc-map-v32-custom")
		testUnmarshalErr(&s32v2, bs32, h, t, "dec-map-v32-p-len")
		testDeepEqualErrHandle(s32v1, s32v2, h, t, "equal-map-v32-p-len")
		testReleaseBytes(bs32)
	}
	for _, v := range []map[string]int{nil, {}, {"some-string-3": 0, "some-string-1": 111}} {
		var v33v1, v33v2 map[string]int
		var bs33 []byte
		v33v1 = v
		bs33 = testMarshalErr(v33v1, h, t, "enc-map-v33")
		if v != nil {
			v33v2 = make(map[string]int, len(v)) // reset map
			testUnmarshalErr(v33v2, bs33, h, t, "dec-map-v33")
			testDeepEqualErrHandle(v33v1, v33v2, h, t, "equal-map-v33")
			v33v2 = make(map[string]int, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v33v2), bs33, h, t, "dec-map-v33-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v33v1, v33v2, h, t, "equal-map-v33-noaddr")
		}
		if v == nil {
			v33v2 = nil
		} else {
			v33v2 = make(map[string]int, len(v))
		} // reset map
		testUnmarshalErr(&v33v2, bs33, h, t, "dec-map-v33-p-len")
		testDeepEqualErrHandle(v33v1, v33v2, h, t, "equal-map-v33-p-len")
		testReleaseBytes(bs33)
		bs33 = testMarshalErr(&v33v1, h, t, "enc-map-v33-p")
		v33v2 = nil
		testUnmarshalErr(&v33v2, bs33, h, t, "dec-map-v33-p-nil")
		testDeepEqualErrHandle(v33v1, v33v2, h, t, "equal-map-v33-p-nil")
		testReleaseBytes(bs33)
		// ...
		if v == nil {
			v33v2 = nil
		} else {
			v33v2 = make(map[string]int, len(v))
		} // reset map
		var v33v3, v33v4 typMapMapStringInt
		v33v3 = typMapMapStringInt(v33v1)
		v33v4 = typMapMapStringInt(v33v2)
		if v != nil {
			bs33 = testMarshalErr(v33v3, h, t, "enc-map-v33-custom")
			testUnmarshalErr(v33v4, bs33, h, t, "dec-map-v33-p-len")
			testDeepEqualErrHandle(v33v3, v33v4, h, t, "equal-map-v33-p-len")
			testReleaseBytes(bs33)
		}
		type s33T struct {
			M  map[string]int
			Mp *map[string]int
		}
		var m33v99 = map[string]int{
			"":              0,
			"some-string-2": 77,
		}
		var s33v1, s33v2 s33T
		bs33 = testMarshalErr(s33v1, h, t, "enc-map-v33-custom")
		testUnmarshalErr(&s33v2, bs33, h, t, "dec-map-v33-p-len")
		testDeepEqualErrHandle(s33v1, s33v2, h, t, "equal-map-v33-p-len")
		testReleaseBytes(bs33)
		s33v2 = s33T{}
		s33v1.M = m33v99
		bs33 = testMarshalErr(s33v1, h, t, "enc-map-v33-custom")
		testUnmarshalErr(&s33v2, bs33, h, t, "dec-map-v33-p-len")
		testDeepEqualErrHandle(s33v1, s33v2, h, t, "equal-map-v33-p-len")
		testReleaseBytes(bs33)
		s33v2 = s33T{}
		s33v1.Mp = &m33v99
		bs33 = testMarshalErr(s33v1, h, t, "enc-map-v33-custom")
		testUnmarshalErr(&s33v2, bs33, h, t, "dec-map-v33-p-len")
		testDeepEqualErrHandle(s33v1, s33v2, h, t, "equal-map-v33-p-len")
		testReleaseBytes(bs33)
	}
	for _, v := range []map[string]int32{nil, {}, {"some-string-3": 0, "some-string-1": 127}} {
		var v34v1, v34v2 map[string]int32
		var bs34 []byte
		v34v1 = v
		bs34 = testMarshalErr(v34v1, h, t, "enc-map-v34")
		if v != nil {
			v34v2 = make(map[string]int32, len(v)) // reset map
			testUnmarshalErr(v34v2, bs34, h, t, "dec-map-v34")
			testDeepEqualErrHandle(v34v1, v34v2, h, t, "equal-map-v34")
			v34v2 = make(map[string]int32, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v34v2), bs34, h, t, "dec-map-v34-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v34v1, v34v2, h, t, "equal-map-v34-noaddr")
		}
		if v == nil {
			v34v2 = nil
		} else {
			v34v2 = make(map[string]int32, len(v))
		} // reset map
		testUnmarshalErr(&v34v2, bs34, h, t, "dec-map-v34-p-len")
		testDeepEqualErrHandle(v34v1, v34v2, h, t, "equal-map-v34-p-len")
		testReleaseBytes(bs34)
		bs34 = testMarshalErr(&v34v1, h, t, "enc-map-v34-p")
		v34v2 = nil
		testUnmarshalErr(&v34v2, bs34, h, t, "dec-map-v34-p-nil")
		testDeepEqualErrHandle(v34v1, v34v2, h, t, "equal-map-v34-p-nil")
		testReleaseBytes(bs34)
		// ...
		if v == nil {
			v34v2 = nil
		} else {
			v34v2 = make(map[string]int32, len(v))
		} // reset map
		var v34v3, v34v4 typMapMapStringInt32
		v34v3 = typMapMapStringInt32(v34v1)
		v34v4 = typMapMapStringInt32(v34v2)
		if v != nil {
			bs34 = testMarshalErr(v34v3, h, t, "enc-map-v34-custom")
			testUnmarshalErr(v34v4, bs34, h, t, "dec-map-v34-p-len")
			testDeepEqualErrHandle(v34v3, v34v4, h, t, "equal-map-v34-p-len")
			testReleaseBytes(bs34)
		}
		type s34T struct {
			M  map[string]int32
			Mp *map[string]int32
		}
		var m34v99 = map[string]int32{
			"":              0,
			"some-string-2": 111,
		}
		var s34v1, s34v2 s34T
		bs34 = testMarshalErr(s34v1, h, t, "enc-map-v34-custom")
		testUnmarshalErr(&s34v2, bs34, h, t, "dec-map-v34-p-len")
		testDeepEqualErrHandle(s34v1, s34v2, h, t, "equal-map-v34-p-len")
		testReleaseBytes(bs34)
		s34v2 = s34T{}
		s34v1.M = m34v99
		bs34 = testMarshalErr(s34v1, h, t, "enc-map-v34-custom")
		testUnmarshalErr(&s34v2, bs34, h, t, "dec-map-v34-p-len")
		testDeepEqualErrHandle(s34v1, s34v2, h, t, "equal-map-v34-p-len")
		testReleaseBytes(bs34)
		s34v2 = s34T{}
		s34v1.Mp = &m34v99
		bs34 = testMarshalErr(s34v1, h, t, "enc-map-v34-custom")
		testUnmarshalErr(&s34v2, bs34, h, t, "dec-map-v34-p-len")
		testDeepEqualErrHandle(s34v1, s34v2, h, t, "equal-map-v34-p-len")
		testReleaseBytes(bs34)
	}
	for _, v := range []map[string]float64{nil, {}, {"some-string-3": 0, "some-string-1": 33.3e3}} {
		var v35v1, v35v2 map[string]float64
		var bs35 []byte
		v35v1 = v
		bs35 = testMarshalErr(v35v1, h, t, "enc-map-v35")
		if v != nil {
			v35v2 = make(map[string]float64, len(v)) // reset map
			testUnmarshalErr(v35v2, bs35, h, t, "dec-map-v35")
			testDeepEqualErrHandle(v35v1, v35v2, h, t, "equal-map-v35")
			v35v2 = make(map[string]float64, len(v))                                   // reset map
			testUnmarshalErr(reflect.ValueOf(v35v2), bs35, h, t, "dec-map-v35-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v35v1, v35v2, h, t, "equal-map-v35-noaddr")
		}
		if v == nil {
			v35v2 = nil
		} else {
			v35v2 = make(map[string]float64, len(v))
		} // reset map
		testUnmarshalErr(&v35v2, bs35, h, t, "dec-map-v35-p-len")
		testDeepEqualErrHandle(v35v1, v35v2, h, t, "equal-map-v35-p-len")
		testReleaseBytes(bs35)
		bs35 = testMarshalErr(&v35v1, h, t, "enc-map-v35-p")
		v35v2 = nil
		testUnmarshalErr(&v35v2, bs35, h, t, "dec-map-v35-p-nil")
		testDeepEqualErrHandle(v35v1, v35v2, h, t, "equal-map-v35-p-nil")
		testReleaseBytes(bs35)
		// ...
		if v == nil {
			v35v2 = nil
		} else {
			v35v2 = make(map[string]float64, len(v))
		} // reset map
		var v35v3, v35v4 typMapMapStringFloat64
		v35v3 = typMapMapStringFloat64(v35v1)
		v35v4 = typMapMapStringFloat64(v35v2)
		if v != nil {
			bs35 = testMarshalErr(v35v3, h, t, "enc-map-v35-custom")
			testUnmarshalErr(v35v4, bs35, h, t, "dec-map-v35-p-len")
			testDeepEqualErrHandle(v35v3, v35v4, h, t, "equal-map-v35-p-len")
			testReleaseBytes(bs35)
		}
		type s35T struct {
			M  map[string]float64
			Mp *map[string]float64
		}
		var m35v99 = map[string]float64{
			"":              0,
			"some-string-2": 11.1,
		}
		var s35v1, s35v2 s35T
		bs35 = testMarshalErr(s35v1, h, t, "enc-map-v35-custom")
		testUnmarshalErr(&s35v2, bs35, h, t, "dec-map-v35-p-len")
		testDeepEqualErrHandle(s35v1, s35v2, h, t, "equal-map-v35-p-len")
		testReleaseBytes(bs35)
		s35v2 = s35T{}
		s35v1.M = m35v99
		bs35 = testMarshalErr(s35v1, h, t, "enc-map-v35-custom")
		testUnmarshalErr(&s35v2, bs35, h, t, "dec-map-v35-p-len")
		testDeepEqualErrHandle(s35v1, s35v2, h, t, "equal-map-v35-p-len")
		testReleaseBytes(bs35)
		s35v2 = s35T{}
		s35v1.Mp = &m35v99
		bs35 = testMarshalErr(s35v1, h, t, "enc-map-v35-custom")
		testUnmarshalErr(&s35v2, bs35, h, t, "dec-map-v35-p-len")
		testDeepEqualErrHandle(s35v1, s35v2, h, t, "equal-map-v35-p-len")
		testReleaseBytes(bs35)
	}
	for _, v := range []map[string]bool{nil, {}, {"some-string-3": false, "some-string-1": true}} {
		var v36v1, v36v2 map[string]bool
		var bs36 []byte
		v36v1 = v
		bs36 = testMarshalErr(v36v1, h, t, "enc-map-v36")
		if v != nil {
			v36v2 = make(map[string]bool, len(v)) // reset map
			testUnmarshalErr(v36v2, bs36, h, t, "dec-map-v36")
			testDeepEqualErrHandle(v36v1, v36v2, h, t, "equal-map-v36")
			v36v2 = make(map[string]bool, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v36v2), bs36, h, t, "dec-map-v36-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v36v1, v36v2, h, t, "equal-map-v36-noaddr")
		}
		if v == nil {
			v36v2 = nil
		} else {
			v36v2 = make(map[string]bool, len(v))
		} // reset map
		testUnmarshalErr(&v36v2, bs36, h, t, "dec-map-v36-p-len")
		testDeepEqualErrHandle(v36v1, v36v2, h, t, "equal-map-v36-p-len")
		testReleaseBytes(bs36)
		bs36 = testMarshalErr(&v36v1, h, t, "enc-map-v36-p")
		v36v2 = nil
		testUnmarshalErr(&v36v2, bs36, h, t, "dec-map-v36-p-nil")
		testDeepEqualErrHandle(v36v1, v36v2, h, t, "equal-map-v36-p-nil")
		testReleaseBytes(bs36)
		// ...
		if v == nil {
			v36v2 = nil
		} else {
			v36v2 = make(map[string]bool, len(v))
		} // reset map
		var v36v3, v36v4 typMapMapStringBool
		v36v3 = typMapMapStringBool(v36v1)
		v36v4 = typMapMapStringBool(v36v2)
		if v != nil {
			bs36 = testMarshalErr(v36v3, h, t, "enc-map-v36-custom")
			testUnmarshalErr(v36v4, bs36, h, t, "dec-map-v36-p-len")
			testDeepEqualErrHandle(v36v3, v36v4, h, t, "equal-map-v36-p-len")
			testReleaseBytes(bs36)
		}
		type s36T struct {
			M  map[string]bool
			Mp *map[string]bool
		}
		var m36v99 = map[string]bool{
			"":              false,
			"some-string-2": false,
		}
		var s36v1, s36v2 s36T
		bs36 = testMarshalErr(s36v1, h, t, "enc-map-v36-custom")
		testUnmarshalErr(&s36v2, bs36, h, t, "dec-map-v36-p-len")
		testDeepEqualErrHandle(s36v1, s36v2, h, t, "equal-map-v36-p-len")
		testReleaseBytes(bs36)
		s36v2 = s36T{}
		s36v1.M = m36v99
		bs36 = testMarshalErr(s36v1, h, t, "enc-map-v36-custom")
		testUnmarshalErr(&s36v2, bs36, h, t, "dec-map-v36-p-len")
		testDeepEqualErrHandle(s36v1, s36v2, h, t, "equal-map-v36-p-len")
		testReleaseBytes(bs36)
		s36v2 = s36T{}
		s36v1.Mp = &m36v99
		bs36 = testMarshalErr(s36v1, h, t, "enc-map-v36-custom")
		testUnmarshalErr(&s36v2, bs36, h, t, "dec-map-v36-p-len")
		testDeepEqualErrHandle(s36v1, s36v2, h, t, "equal-map-v36-p-len")
		testReleaseBytes(bs36)
	}
	for _, v := range []map[uint8]interface{}{nil, {}, {77: nil, 127: "string-is-an-interface-3"}} {
		var v37v1, v37v2 map[uint8]interface{}
		var bs37 []byte
		v37v1 = v
		bs37 = testMarshalErr(v37v1, h, t, "enc-map-v37")
		if v != nil {
			v37v2 = make(map[uint8]interface{}, len(v)) // reset map
			testUnmarshalErr(v37v2, bs37, h, t, "dec-map-v37")
			testDeepEqualErrHandle(v37v1, v37v2, h, t, "equal-map-v37")
			v37v2 = make(map[uint8]interface{}, len(v))                                // reset map
			testUnmarshalErr(reflect.ValueOf(v37v2), bs37, h, t, "dec-map-v37-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v37v1, v37v2, h, t, "equal-map-v37-noaddr")
		}
		if v == nil {
			v37v2 = nil
		} else {
			v37v2 = make(map[uint8]interface{}, len(v))
		} // reset map
		testUnmarshalErr(&v37v2, bs37, h, t, "dec-map-v37-p-len")
		testDeepEqualErrHandle(v37v1, v37v2, h, t, "equal-map-v37-p-len")
		testReleaseBytes(bs37)
		bs37 = testMarshalErr(&v37v1, h, t, "enc-map-v37-p")
		v37v2 = nil
		testUnmarshalErr(&v37v2, bs37, h, t, "dec-map-v37-p-nil")
		testDeepEqualErrHandle(v37v1, v37v2, h, t, "equal-map-v37-p-nil")
		testReleaseBytes(bs37)
		// ...
		if v == nil {
			v37v2 = nil
		} else {
			v37v2 = make(map[uint8]interface{}, len(v))
		} // reset map
		var v37v3, v37v4 typMapMapUint8Intf
		v37v3 = typMapMapUint8Intf(v37v1)
		v37v4 = typMapMapUint8Intf(v37v2)
		if v != nil {
			bs37 = testMarshalErr(v37v3, h, t, "enc-map-v37-custom")
			testUnmarshalErr(v37v4, bs37, h, t, "dec-map-v37-p-len")
			testDeepEqualErrHandle(v37v3, v37v4, h, t, "equal-map-v37-p-len")
			testReleaseBytes(bs37)
		}
		type s37T struct {
			M  map[uint8]interface{}
			Mp *map[uint8]interface{}
		}
		var m37v99 = map[uint8]interface{}{
			0:   nil,
			111: "string-is-an-interface-1",
		}
		var s37v1, s37v2 s37T
		bs37 = testMarshalErr(s37v1, h, t, "enc-map-v37-custom")
		testUnmarshalErr(&s37v2, bs37, h, t, "dec-map-v37-p-len")
		testDeepEqualErrHandle(s37v1, s37v2, h, t, "equal-map-v37-p-len")
		testReleaseBytes(bs37)
		s37v2 = s37T{}
		s37v1.M = m37v99
		bs37 = testMarshalErr(s37v1, h, t, "enc-map-v37-custom")
		testUnmarshalErr(&s37v2, bs37, h, t, "dec-map-v37-p-len")
		testDeepEqualErrHandle(s37v1, s37v2, h, t, "equal-map-v37-p-len")
		testReleaseBytes(bs37)
		s37v2 = s37T{}
		s37v1.Mp = &m37v99
		bs37 = testMarshalErr(s37v1, h, t, "enc-map-v37-custom")
		testUnmarshalErr(&s37v2, bs37, h, t, "dec-map-v37-p-len")
		testDeepEqualErrHandle(s37v1, s37v2, h, t, "equal-map-v37-p-len")
		testReleaseBytes(bs37)
	}
	for _, v := range []map[uint8]string{nil, {}, {77: "", 127: "some-string-3"}} {
		var v38v1, v38v2 map[uint8]string
		var bs38 []byte
		v38v1 = v
		bs38 = testMarshalErr(v38v1, h, t, "enc-map-v38")
		if v != nil {
			v38v2 = make(map[uint8]string, len(v)) // reset map
			testUnmarshalErr(v38v2, bs38, h, t, "dec-map-v38")
			testDeepEqualErrHandle(v38v1, v38v2, h, t, "equal-map-v38")
			v38v2 = make(map[uint8]string, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v38v2), bs38, h, t, "dec-map-v38-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v38v1, v38v2, h, t, "equal-map-v38-noaddr")
		}
		if v == nil {
			v38v2 = nil
		} else {
			v38v2 = make(map[uint8]string, len(v))
		} // reset map
		testUnmarshalErr(&v38v2, bs38, h, t, "dec-map-v38-p-len")
		testDeepEqualErrHandle(v38v1, v38v2, h, t, "equal-map-v38-p-len")
		testReleaseBytes(bs38)
		bs38 = testMarshalErr(&v38v1, h, t, "enc-map-v38-p")
		v38v2 = nil
		testUnmarshalErr(&v38v2, bs38, h, t, "dec-map-v38-p-nil")
		testDeepEqualErrHandle(v38v1, v38v2, h, t, "equal-map-v38-p-nil")
		testReleaseBytes(bs38)
		// ...
		if v == nil {
			v38v2 = nil
		} else {
			v38v2 = make(map[uint8]string, len(v))
		} // reset map
		var v38v3, v38v4 typMapMapUint8String
		v38v3 = typMapMapUint8String(v38v1)
		v38v4 = typMapMapUint8String(v38v2)
		if v != nil {
			bs38 = testMarshalErr(v38v3, h, t, "enc-map-v38-custom")
			testUnmarshalErr(v38v4, bs38, h, t, "dec-map-v38-p-len")
			testDeepEqualErrHandle(v38v3, v38v4, h, t, "equal-map-v38-p-len")
			testReleaseBytes(bs38)
		}
		type s38T struct {
			M  map[uint8]string
			Mp *map[uint8]string
		}
		var m38v99 = map[uint8]string{
			0:   "",
			111: "some-string-1",
		}
		var s38v1, s38v2 s38T
		bs38 = testMarshalErr(s38v1, h, t, "enc-map-v38-custom")
		testUnmarshalErr(&s38v2, bs38, h, t, "dec-map-v38-p-len")
		testDeepEqualErrHandle(s38v1, s38v2, h, t, "equal-map-v38-p-len")
		testReleaseBytes(bs38)
		s38v2 = s38T{}
		s38v1.M = m38v99
		bs38 = testMarshalErr(s38v1, h, t, "enc-map-v38-custom")
		testUnmarshalErr(&s38v2, bs38, h, t, "dec-map-v38-p-len")
		testDeepEqualErrHandle(s38v1, s38v2, h, t, "equal-map-v38-p-len")
		testReleaseBytes(bs38)
		s38v2 = s38T{}
		s38v1.Mp = &m38v99
		bs38 = testMarshalErr(s38v1, h, t, "enc-map-v38-custom")
		testUnmarshalErr(&s38v2, bs38, h, t, "dec-map-v38-p-len")
		testDeepEqualErrHandle(s38v1, s38v2, h, t, "equal-map-v38-p-len")
		testReleaseBytes(bs38)
	}
	for _, v := range []map[uint8][]byte{nil, {}, {77: nil, 127: []byte("some-string-3")}} {
		var v39v1, v39v2 map[uint8][]byte
		var bs39 []byte
		v39v1 = v
		bs39 = testMarshalErr(v39v1, h, t, "enc-map-v39")
		if v != nil {
			v39v2 = make(map[uint8][]byte, len(v)) // reset map
			testUnmarshalErr(v39v2, bs39, h, t, "dec-map-v39")
			testDeepEqualErrHandle(v39v1, v39v2, h, t, "equal-map-v39")
			v39v2 = make(map[uint8][]byte, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v39v2), bs39, h, t, "dec-map-v39-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v39v1, v39v2, h, t, "equal-map-v39-noaddr")
		}
		if v == nil {
			v39v2 = nil
		} else {
			v39v2 = make(map[uint8][]byte, len(v))
		} // reset map
		testUnmarshalErr(&v39v2, bs39, h, t, "dec-map-v39-p-len")
		testDeepEqualErrHandle(v39v1, v39v2, h, t, "equal-map-v39-p-len")
		testReleaseBytes(bs39)
		bs39 = testMarshalErr(&v39v1, h, t, "enc-map-v39-p")
		v39v2 = nil
		testUnmarshalErr(&v39v2, bs39, h, t, "dec-map-v39-p-nil")
		testDeepEqualErrHandle(v39v1, v39v2, h, t, "equal-map-v39-p-nil")
		testReleaseBytes(bs39)
		// ...
		if v == nil {
			v39v2 = nil
		} else {
			v39v2 = make(map[uint8][]byte, len(v))
		} // reset map
		var v39v3, v39v4 typMapMapUint8Bytes
		v39v3 = typMapMapUint8Bytes(v39v1)
		v39v4 = typMapMapUint8Bytes(v39v2)
		if v != nil {
			bs39 = testMarshalErr(v39v3, h, t, "enc-map-v39-custom")
			testUnmarshalErr(v39v4, bs39, h, t, "dec-map-v39-p-len")
			testDeepEqualErrHandle(v39v3, v39v4, h, t, "equal-map-v39-p-len")
			testReleaseBytes(bs39)
		}
		type s39T struct {
			M  map[uint8][]byte
			Mp *map[uint8][]byte
		}
		var m39v99 = map[uint8][]byte{
			0:   nil,
			111: []byte("some-string-1"),
		}
		var s39v1, s39v2 s39T
		bs39 = testMarshalErr(s39v1, h, t, "enc-map-v39-custom")
		testUnmarshalErr(&s39v2, bs39, h, t, "dec-map-v39-p-len")
		testDeepEqualErrHandle(s39v1, s39v2, h, t, "equal-map-v39-p-len")
		testReleaseBytes(bs39)
		s39v2 = s39T{}
		s39v1.M = m39v99
		bs39 = testMarshalErr(s39v1, h, t, "enc-map-v39-custom")
		testUnmarshalErr(&s39v2, bs39, h, t, "dec-map-v39-p-len")
		testDeepEqualErrHandle(s39v1, s39v2, h, t, "equal-map-v39-p-len")
		testReleaseBytes(bs39)
		s39v2 = s39T{}
		s39v1.Mp = &m39v99
		bs39 = testMarshalErr(s39v1, h, t, "enc-map-v39-custom")
		testUnmarshalErr(&s39v2, bs39, h, t, "dec-map-v39-p-len")
		testDeepEqualErrHandle(s39v1, s39v2, h, t, "equal-map-v39-p-len")
		testReleaseBytes(bs39)
	}
	for _, v := range []map[uint8]uint8{nil, {}, {77: 0, 127: 111}} {
		var v40v1, v40v2 map[uint8]uint8
		var bs40 []byte
		v40v1 = v
		bs40 = testMarshalErr(v40v1, h, t, "enc-map-v40")
		if v != nil {
			v40v2 = make(map[uint8]uint8, len(v)) // reset map
			testUnmarshalErr(v40v2, bs40, h, t, "dec-map-v40")
			testDeepEqualErrHandle(v40v1, v40v2, h, t, "equal-map-v40")
			v40v2 = make(map[uint8]uint8, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v40v2), bs40, h, t, "dec-map-v40-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v40v1, v40v2, h, t, "equal-map-v40-noaddr")
		}
		if v == nil {
			v40v2 = nil
		} else {
			v40v2 = make(map[uint8]uint8, len(v))
		} // reset map
		testUnmarshalErr(&v40v2, bs40, h, t, "dec-map-v40-p-len")
		testDeepEqualErrHandle(v40v1, v40v2, h, t, "equal-map-v40-p-len")
		testReleaseBytes(bs40)
		bs40 = testMarshalErr(&v40v1, h, t, "enc-map-v40-p")
		v40v2 = nil
		testUnmarshalErr(&v40v2, bs40, h, t, "dec-map-v40-p-nil")
		testDeepEqualErrHandle(v40v1, v40v2, h, t, "equal-map-v40-p-nil")
		testReleaseBytes(bs40)
		// ...
		if v == nil {
			v40v2 = nil
		} else {
			v40v2 = make(map[uint8]uint8, len(v))
		} // reset map
		var v40v3, v40v4 typMapMapUint8Uint8
		v40v3 = typMapMapUint8Uint8(v40v1)
		v40v4 = typMapMapUint8Uint8(v40v2)
		if v != nil {
			bs40 = testMarshalErr(v40v3, h, t, "enc-map-v40-custom")
			testUnmarshalErr(v40v4, bs40, h, t, "dec-map-v40-p-len")
			testDeepEqualErrHandle(v40v3, v40v4, h, t, "equal-map-v40-p-len")
			testReleaseBytes(bs40)
		}
		type s40T struct {
			M  map[uint8]uint8
			Mp *map[uint8]uint8
		}
		var m40v99 = map[uint8]uint8{
			0:  0,
			77: 127,
		}
		var s40v1, s40v2 s40T
		bs40 = testMarshalErr(s40v1, h, t, "enc-map-v40-custom")
		testUnmarshalErr(&s40v2, bs40, h, t, "dec-map-v40-p-len")
		testDeepEqualErrHandle(s40v1, s40v2, h, t, "equal-map-v40-p-len")
		testReleaseBytes(bs40)
		s40v2 = s40T{}
		s40v1.M = m40v99
		bs40 = testMarshalErr(s40v1, h, t, "enc-map-v40-custom")
		testUnmarshalErr(&s40v2, bs40, h, t, "dec-map-v40-p-len")
		testDeepEqualErrHandle(s40v1, s40v2, h, t, "equal-map-v40-p-len")
		testReleaseBytes(bs40)
		s40v2 = s40T{}
		s40v1.Mp = &m40v99
		bs40 = testMarshalErr(s40v1, h, t, "enc-map-v40-custom")
		testUnmarshalErr(&s40v2, bs40, h, t, "dec-map-v40-p-len")
		testDeepEqualErrHandle(s40v1, s40v2, h, t, "equal-map-v40-p-len")
		testReleaseBytes(bs40)
	}
	for _, v := range []map[uint8]uint64{nil, {}, {111: 0, 77: 127}} {
		var v41v1, v41v2 map[uint8]uint64
		var bs41 []byte
		v41v1 = v
		bs41 = testMarshalErr(v41v1, h, t, "enc-map-v41")
		if v != nil {
			v41v2 = make(map[uint8]uint64, len(v)) // reset map
			testUnmarshalErr(v41v2, bs41, h, t, "dec-map-v41")
			testDeepEqualErrHandle(v41v1, v41v2, h, t, "equal-map-v41")
			v41v2 = make(map[uint8]uint64, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v41v2), bs41, h, t, "dec-map-v41-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v41v1, v41v2, h, t, "equal-map-v41-noaddr")
		}
		if v == nil {
			v41v2 = nil
		} else {
			v41v2 = make(map[uint8]uint64, len(v))
		} // reset map
		testUnmarshalErr(&v41v2, bs41, h, t, "dec-map-v41-p-len")
		testDeepEqualErrHandle(v41v1, v41v2, h, t, "equal-map-v41-p-len")
		testReleaseBytes(bs41)
		bs41 = testMarshalErr(&v41v1, h, t, "enc-map-v41-p")
		v41v2 = nil
		testUnmarshalErr(&v41v2, bs41, h, t, "dec-map-v41-p-nil")
		testDeepEqualErrHandle(v41v1, v41v2, h, t, "equal-map-v41-p-nil")
		testReleaseBytes(bs41)
		// ...
		if v == nil {
			v41v2 = nil
		} else {
			v41v2 = make(map[uint8]uint64, len(v))
		} // reset map
		var v41v3, v41v4 typMapMapUint8Uint64
		v41v3 = typMapMapUint8Uint64(v41v1)
		v41v4 = typMapMapUint8Uint64(v41v2)
		if v != nil {
			bs41 = testMarshalErr(v41v3, h, t, "enc-map-v41-custom")
			testUnmarshalErr(v41v4, bs41, h, t, "dec-map-v41-p-len")
			testDeepEqualErrHandle(v41v3, v41v4, h, t, "equal-map-v41-p-len")
			testReleaseBytes(bs41)
		}
		type s41T struct {
			M  map[uint8]uint64
			Mp *map[uint8]uint64
		}
		var m41v99 = map[uint8]uint64{
			0:   0,
			111: 77,
		}
		var s41v1, s41v2 s41T
		bs41 = testMarshalErr(s41v1, h, t, "enc-map-v41-custom")
		testUnmarshalErr(&s41v2, bs41, h, t, "dec-map-v41-p-len")
		testDeepEqualErrHandle(s41v1, s41v2, h, t, "equal-map-v41-p-len")
		testReleaseBytes(bs41)
		s41v2 = s41T{}
		s41v1.M = m41v99
		bs41 = testMarshalErr(s41v1, h, t, "enc-map-v41-custom")
		testUnmarshalErr(&s41v2, bs41, h, t, "dec-map-v41-p-len")
		testDeepEqualErrHandle(s41v1, s41v2, h, t, "equal-map-v41-p-len")
		testReleaseBytes(bs41)
		s41v2 = s41T{}
		s41v1.Mp = &m41v99
		bs41 = testMarshalErr(s41v1, h, t, "enc-map-v41-custom")
		testUnmarshalErr(&s41v2, bs41, h, t, "dec-map-v41-p-len")
		testDeepEqualErrHandle(s41v1, s41v2, h, t, "equal-map-v41-p-len")
		testReleaseBytes(bs41)
	}
	for _, v := range []map[uint8]int{nil, {}, {127: 0, 111: 77}} {
		var v42v1, v42v2 map[uint8]int
		var bs42 []byte
		v42v1 = v
		bs42 = testMarshalErr(v42v1, h, t, "enc-map-v42")
		if v != nil {
			v42v2 = make(map[uint8]int, len(v)) // reset map
			testUnmarshalErr(v42v2, bs42, h, t, "dec-map-v42")
			testDeepEqualErrHandle(v42v1, v42v2, h, t, "equal-map-v42")
			v42v2 = make(map[uint8]int, len(v))                                        // reset map
			testUnmarshalErr(reflect.ValueOf(v42v2), bs42, h, t, "dec-map-v42-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v42v1, v42v2, h, t, "equal-map-v42-noaddr")
		}
		if v == nil {
			v42v2 = nil
		} else {
			v42v2 = make(map[uint8]int, len(v))
		} // reset map
		testUnmarshalErr(&v42v2, bs42, h, t, "dec-map-v42-p-len")
		testDeepEqualErrHandle(v42v1, v42v2, h, t, "equal-map-v42-p-len")
		testReleaseBytes(bs42)
		bs42 = testMarshalErr(&v42v1, h, t, "enc-map-v42-p")
		v42v2 = nil
		testUnmarshalErr(&v42v2, bs42, h, t, "dec-map-v42-p-nil")
		testDeepEqualErrHandle(v42v1, v42v2, h, t, "equal-map-v42-p-nil")
		testReleaseBytes(bs42)
		// ...
		if v == nil {
			v42v2 = nil
		} else {
			v42v2 = make(map[uint8]int, len(v))
		} // reset map
		var v42v3, v42v4 typMapMapUint8Int
		v42v3 = typMapMapUint8Int(v42v1)
		v42v4 = typMapMapUint8Int(v42v2)
		if v != nil {
			bs42 = testMarshalErr(v42v3, h, t, "enc-map-v42-custom")
			testUnmarshalErr(v42v4, bs42, h, t, "dec-map-v42-p-len")
			testDeepEqualErrHandle(v42v3, v42v4, h, t, "equal-map-v42-p-len")
			testReleaseBytes(bs42)
		}
		type s42T struct {
			M  map[uint8]int
			Mp *map[uint8]int
		}
		var m42v99 = map[uint8]int{
			0:   0,
			127: 111,
		}
		var s42v1, s42v2 s42T
		bs42 = testMarshalErr(s42v1, h, t, "enc-map-v42-custom")
		testUnmarshalErr(&s42v2, bs42, h, t, "dec-map-v42-p-len")
		testDeepEqualErrHandle(s42v1, s42v2, h, t, "equal-map-v42-p-len")
		testReleaseBytes(bs42)
		s42v2 = s42T{}
		s42v1.M = m42v99
		bs42 = testMarshalErr(s42v1, h, t, "enc-map-v42-custom")
		testUnmarshalErr(&s42v2, bs42, h, t, "dec-map-v42-p-len")
		testDeepEqualErrHandle(s42v1, s42v2, h, t, "equal-map-v42-p-len")
		testReleaseBytes(bs42)
		s42v2 = s42T{}
		s42v1.Mp = &m42v99
		bs42 = testMarshalErr(s42v1, h, t, "enc-map-v42-custom")
		testUnmarshalErr(&s42v2, bs42, h, t, "dec-map-v42-p-len")
		testDeepEqualErrHandle(s42v1, s42v2, h, t, "equal-map-v42-p-len")
		testReleaseBytes(bs42)
	}
	for _, v := range []map[uint8]int32{nil, {}, {77: 0, 127: 111}} {
		var v43v1, v43v2 map[uint8]int32
		var bs43 []byte
		v43v1 = v
		bs43 = testMarshalErr(v43v1, h, t, "enc-map-v43")
		if v != nil {
			v43v2 = make(map[uint8]int32, len(v)) // reset map
			testUnmarshalErr(v43v2, bs43, h, t, "dec-map-v43")
			testDeepEqualErrHandle(v43v1, v43v2, h, t, "equal-map-v43")
			v43v2 = make(map[uint8]int32, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v43v2), bs43, h, t, "dec-map-v43-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v43v1, v43v2, h, t, "equal-map-v43-noaddr")
		}
		if v == nil {
			v43v2 = nil
		} else {
			v43v2 = make(map[uint8]int32, len(v))
		} // reset map
		testUnmarshalErr(&v43v2, bs43, h, t, "dec-map-v43-p-len")
		testDeepEqualErrHandle(v43v1, v43v2, h, t, "equal-map-v43-p-len")
		testReleaseBytes(bs43)
		bs43 = testMarshalErr(&v43v1, h, t, "enc-map-v43-p")
		v43v2 = nil
		testUnmarshalErr(&v43v2, bs43, h, t, "dec-map-v43-p-nil")
		testDeepEqualErrHandle(v43v1, v43v2, h, t, "equal-map-v43-p-nil")
		testReleaseBytes(bs43)
		// ...
		if v == nil {
			v43v2 = nil
		} else {
			v43v2 = make(map[uint8]int32, len(v))
		} // reset map
		var v43v3, v43v4 typMapMapUint8Int32
		v43v3 = typMapMapUint8Int32(v43v1)
		v43v4 = typMapMapUint8Int32(v43v2)
		if v != nil {
			bs43 = testMarshalErr(v43v3, h, t, "enc-map-v43-custom")
			testUnmarshalErr(v43v4, bs43, h, t, "dec-map-v43-p-len")
			testDeepEqualErrHandle(v43v3, v43v4, h, t, "equal-map-v43-p-len")
			testReleaseBytes(bs43)
		}
		type s43T struct {
			M  map[uint8]int32
			Mp *map[uint8]int32
		}
		var m43v99 = map[uint8]int32{
			0:  0,
			77: 127,
		}
		var s43v1, s43v2 s43T
		bs43 = testMarshalErr(s43v1, h, t, "enc-map-v43-custom")
		testUnmarshalErr(&s43v2, bs43, h, t, "dec-map-v43-p-len")
		testDeepEqualErrHandle(s43v1, s43v2, h, t, "equal-map-v43-p-len")
		testReleaseBytes(bs43)
		s43v2 = s43T{}
		s43v1.M = m43v99
		bs43 = testMarshalErr(s43v1, h, t, "enc-map-v43-custom")
		testUnmarshalErr(&s43v2, bs43, h, t, "dec-map-v43-p-len")
		testDeepEqualErrHandle(s43v1, s43v2, h, t, "equal-map-v43-p-len")
		testReleaseBytes(bs43)
		s43v2 = s43T{}
		s43v1.Mp = &m43v99
		bs43 = testMarshalErr(s43v1, h, t, "enc-map-v43-custom")
		testUnmarshalErr(&s43v2, bs43, h, t, "dec-map-v43-p-len")
		testDeepEqualErrHandle(s43v1, s43v2, h, t, "equal-map-v43-p-len")
		testReleaseBytes(bs43)
	}
	for _, v := range []map[uint8]float64{nil, {}, {111: 0, 77: 22.2}} {
		var v44v1, v44v2 map[uint8]float64
		var bs44 []byte
		v44v1 = v
		bs44 = testMarshalErr(v44v1, h, t, "enc-map-v44")
		if v != nil {
			v44v2 = make(map[uint8]float64, len(v)) // reset map
			testUnmarshalErr(v44v2, bs44, h, t, "dec-map-v44")
			testDeepEqualErrHandle(v44v1, v44v2, h, t, "equal-map-v44")
			v44v2 = make(map[uint8]float64, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v44v2), bs44, h, t, "dec-map-v44-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v44v1, v44v2, h, t, "equal-map-v44-noaddr")
		}
		if v == nil {
			v44v2 = nil
		} else {
			v44v2 = make(map[uint8]float64, len(v))
		} // reset map
		testUnmarshalErr(&v44v2, bs44, h, t, "dec-map-v44-p-len")
		testDeepEqualErrHandle(v44v1, v44v2, h, t, "equal-map-v44-p-len")
		testReleaseBytes(bs44)
		bs44 = testMarshalErr(&v44v1, h, t, "enc-map-v44-p")
		v44v2 = nil
		testUnmarshalErr(&v44v2, bs44, h, t, "dec-map-v44-p-nil")
		testDeepEqualErrHandle(v44v1, v44v2, h, t, "equal-map-v44-p-nil")
		testReleaseBytes(bs44)
		// ...
		if v == nil {
			v44v2 = nil
		} else {
			v44v2 = make(map[uint8]float64, len(v))
		} // reset map
		var v44v3, v44v4 typMapMapUint8Float64
		v44v3 = typMapMapUint8Float64(v44v1)
		v44v4 = typMapMapUint8Float64(v44v2)
		if v != nil {
			bs44 = testMarshalErr(v44v3, h, t, "enc-map-v44-custom")
			testUnmarshalErr(v44v4, bs44, h, t, "dec-map-v44-p-len")
			testDeepEqualErrHandle(v44v3, v44v4, h, t, "equal-map-v44-p-len")
			testReleaseBytes(bs44)
		}
		type s44T struct {
			M  map[uint8]float64
			Mp *map[uint8]float64
		}
		var m44v99 = map[uint8]float64{
			0:   0,
			127: 33.3e3,
		}
		var s44v1, s44v2 s44T
		bs44 = testMarshalErr(s44v1, h, t, "enc-map-v44-custom")
		testUnmarshalErr(&s44v2, bs44, h, t, "dec-map-v44-p-len")
		testDeepEqualErrHandle(s44v1, s44v2, h, t, "equal-map-v44-p-len")
		testReleaseBytes(bs44)
		s44v2 = s44T{}
		s44v1.M = m44v99
		bs44 = testMarshalErr(s44v1, h, t, "enc-map-v44-custom")
		testUnmarshalErr(&s44v2, bs44, h, t, "dec-map-v44-p-len")
		testDeepEqualErrHandle(s44v1, s44v2, h, t, "equal-map-v44-p-len")
		testReleaseBytes(bs44)
		s44v2 = s44T{}
		s44v1.Mp = &m44v99
		bs44 = testMarshalErr(s44v1, h, t, "enc-map-v44-custom")
		testUnmarshalErr(&s44v2, bs44, h, t, "dec-map-v44-p-len")
		testDeepEqualErrHandle(s44v1, s44v2, h, t, "equal-map-v44-p-len")
		testReleaseBytes(bs44)
	}
	for _, v := range []map[uint8]bool{nil, {}, {111: false, 77: true}} {
		var v45v1, v45v2 map[uint8]bool
		var bs45 []byte
		v45v1 = v
		bs45 = testMarshalErr(v45v1, h, t, "enc-map-v45")
		if v != nil {
			v45v2 = make(map[uint8]bool, len(v)) // reset map
			testUnmarshalErr(v45v2, bs45, h, t, "dec-map-v45")
			testDeepEqualErrHandle(v45v1, v45v2, h, t, "equal-map-v45")
			v45v2 = make(map[uint8]bool, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v45v2), bs45, h, t, "dec-map-v45-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v45v1, v45v2, h, t, "equal-map-v45-noaddr")
		}
		if v == nil {
			v45v2 = nil
		} else {
			v45v2 = make(map[uint8]bool, len(v))
		} // reset map
		testUnmarshalErr(&v45v2, bs45, h, t, "dec-map-v45-p-len")
		testDeepEqualErrHandle(v45v1, v45v2, h, t, "equal-map-v45-p-len")
		testReleaseBytes(bs45)
		bs45 = testMarshalErr(&v45v1, h, t, "enc-map-v45-p")
		v45v2 = nil
		testUnmarshalErr(&v45v2, bs45, h, t, "dec-map-v45-p-nil")
		testDeepEqualErrHandle(v45v1, v45v2, h, t, "equal-map-v45-p-nil")
		testReleaseBytes(bs45)
		// ...
		if v == nil {
			v45v2 = nil
		} else {
			v45v2 = make(map[uint8]bool, len(v))
		} // reset map
		var v45v3, v45v4 typMapMapUint8Bool
		v45v3 = typMapMapUint8Bool(v45v1)
		v45v4 = typMapMapUint8Bool(v45v2)
		if v != nil {
			bs45 = testMarshalErr(v45v3, h, t, "enc-map-v45-custom")
			testUnmarshalErr(v45v4, bs45, h, t, "dec-map-v45-p-len")
			testDeepEqualErrHandle(v45v3, v45v4, h, t, "equal-map-v45-p-len")
			testReleaseBytes(bs45)
		}
		type s45T struct {
			M  map[uint8]bool
			Mp *map[uint8]bool
		}
		var m45v99 = map[uint8]bool{
			0:   false,
			127: true,
		}
		var s45v1, s45v2 s45T
		bs45 = testMarshalErr(s45v1, h, t, "enc-map-v45-custom")
		testUnmarshalErr(&s45v2, bs45, h, t, "dec-map-v45-p-len")
		testDeepEqualErrHandle(s45v1, s45v2, h, t, "equal-map-v45-p-len")
		testReleaseBytes(bs45)
		s45v2 = s45T{}
		s45v1.M = m45v99
		bs45 = testMarshalErr(s45v1, h, t, "enc-map-v45-custom")
		testUnmarshalErr(&s45v2, bs45, h, t, "dec-map-v45-p-len")
		testDeepEqualErrHandle(s45v1, s45v2, h, t, "equal-map-v45-p-len")
		testReleaseBytes(bs45)
		s45v2 = s45T{}
		s45v1.Mp = &m45v99
		bs45 = testMarshalErr(s45v1, h, t, "enc-map-v45-custom")
		testUnmarshalErr(&s45v2, bs45, h, t, "dec-map-v45-p-len")
		testDeepEqualErrHandle(s45v1, s45v2, h, t, "equal-map-v45-p-len")
		testReleaseBytes(bs45)
	}
	for _, v := range []map[uint64]interface{}{nil, {}, {111: nil, 77: "string-is-an-interface-2"}} {
		var v46v1, v46v2 map[uint64]interface{}
		var bs46 []byte
		v46v1 = v
		bs46 = testMarshalErr(v46v1, h, t, "enc-map-v46")
		if v != nil {
			v46v2 = make(map[uint64]interface{}, len(v)) // reset map
			testUnmarshalErr(v46v2, bs46, h, t, "dec-map-v46")
			testDeepEqualErrHandle(v46v1, v46v2, h, t, "equal-map-v46")
			v46v2 = make(map[uint64]interface{}, len(v))                               // reset map
			testUnmarshalErr(reflect.ValueOf(v46v2), bs46, h, t, "dec-map-v46-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v46v1, v46v2, h, t, "equal-map-v46-noaddr")
		}
		if v == nil {
			v46v2 = nil
		} else {
			v46v2 = make(map[uint64]interface{}, len(v))
		} // reset map
		testUnmarshalErr(&v46v2, bs46, h, t, "dec-map-v46-p-len")
		testDeepEqualErrHandle(v46v1, v46v2, h, t, "equal-map-v46-p-len")
		testReleaseBytes(bs46)
		bs46 = testMarshalErr(&v46v1, h, t, "enc-map-v46-p")
		v46v2 = nil
		testUnmarshalErr(&v46v2, bs46, h, t, "dec-map-v46-p-nil")
		testDeepEqualErrHandle(v46v1, v46v2, h, t, "equal-map-v46-p-nil")
		testReleaseBytes(bs46)
		// ...
		if v == nil {
			v46v2 = nil
		} else {
			v46v2 = make(map[uint64]interface{}, len(v))
		} // reset map
		var v46v3, v46v4 typMapMapUint64Intf
		v46v3 = typMapMapUint64Intf(v46v1)
		v46v4 = typMapMapUint64Intf(v46v2)
		if v != nil {
			bs46 = testMarshalErr(v46v3, h, t, "enc-map-v46-custom")
			testUnmarshalErr(v46v4, bs46, h, t, "dec-map-v46-p-len")
			testDeepEqualErrHandle(v46v3, v46v4, h, t, "equal-map-v46-p-len")
			testReleaseBytes(bs46)
		}
		type s46T struct {
			M  map[uint64]interface{}
			Mp *map[uint64]interface{}
		}
		var m46v99 = map[uint64]interface{}{
			0:   nil,
			127: "string-is-an-interface-3",
		}
		var s46v1, s46v2 s46T
		bs46 = testMarshalErr(s46v1, h, t, "enc-map-v46-custom")
		testUnmarshalErr(&s46v2, bs46, h, t, "dec-map-v46-p-len")
		testDeepEqualErrHandle(s46v1, s46v2, h, t, "equal-map-v46-p-len")
		testReleaseBytes(bs46)
		s46v2 = s46T{}
		s46v1.M = m46v99
		bs46 = testMarshalErr(s46v1, h, t, "enc-map-v46-custom")
		testUnmarshalErr(&s46v2, bs46, h, t, "dec-map-v46-p-len")
		testDeepEqualErrHandle(s46v1, s46v2, h, t, "equal-map-v46-p-len")
		testReleaseBytes(bs46)
		s46v2 = s46T{}
		s46v1.Mp = &m46v99
		bs46 = testMarshalErr(s46v1, h, t, "enc-map-v46-custom")
		testUnmarshalErr(&s46v2, bs46, h, t, "dec-map-v46-p-len")
		testDeepEqualErrHandle(s46v1, s46v2, h, t, "equal-map-v46-p-len")
		testReleaseBytes(bs46)
	}
	for _, v := range []map[uint64]string{nil, {}, {111: "", 77: "some-string-2"}} {
		var v47v1, v47v2 map[uint64]string
		var bs47 []byte
		v47v1 = v
		bs47 = testMarshalErr(v47v1, h, t, "enc-map-v47")
		if v != nil {
			v47v2 = make(map[uint64]string, len(v)) // reset map
			testUnmarshalErr(v47v2, bs47, h, t, "dec-map-v47")
			testDeepEqualErrHandle(v47v1, v47v2, h, t, "equal-map-v47")
			v47v2 = make(map[uint64]string, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v47v2), bs47, h, t, "dec-map-v47-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v47v1, v47v2, h, t, "equal-map-v47-noaddr")
		}
		if v == nil {
			v47v2 = nil
		} else {
			v47v2 = make(map[uint64]string, len(v))
		} // reset map
		testUnmarshalErr(&v47v2, bs47, h, t, "dec-map-v47-p-len")
		testDeepEqualErrHandle(v47v1, v47v2, h, t, "equal-map-v47-p-len")
		testReleaseBytes(bs47)
		bs47 = testMarshalErr(&v47v1, h, t, "enc-map-v47-p")
		v47v2 = nil
		testUnmarshalErr(&v47v2, bs47, h, t, "dec-map-v47-p-nil")
		testDeepEqualErrHandle(v47v1, v47v2, h, t, "equal-map-v47-p-nil")
		testReleaseBytes(bs47)
		// ...
		if v == nil {
			v47v2 = nil
		} else {
			v47v2 = make(map[uint64]string, len(v))
		} // reset map
		var v47v3, v47v4 typMapMapUint64String
		v47v3 = typMapMapUint64String(v47v1)
		v47v4 = typMapMapUint64String(v47v2)
		if v != nil {
			bs47 = testMarshalErr(v47v3, h, t, "enc-map-v47-custom")
			testUnmarshalErr(v47v4, bs47, h, t, "dec-map-v47-p-len")
			testDeepEqualErrHandle(v47v3, v47v4, h, t, "equal-map-v47-p-len")
			testReleaseBytes(bs47)
		}
		type s47T struct {
			M  map[uint64]string
			Mp *map[uint64]string
		}
		var m47v99 = map[uint64]string{
			0:   "",
			127: "some-string-3",
		}
		var s47v1, s47v2 s47T
		bs47 = testMarshalErr(s47v1, h, t, "enc-map-v47-custom")
		testUnmarshalErr(&s47v2, bs47, h, t, "dec-map-v47-p-len")
		testDeepEqualErrHandle(s47v1, s47v2, h, t, "equal-map-v47-p-len")
		testReleaseBytes(bs47)
		s47v2 = s47T{}
		s47v1.M = m47v99
		bs47 = testMarshalErr(s47v1, h, t, "enc-map-v47-custom")
		testUnmarshalErr(&s47v2, bs47, h, t, "dec-map-v47-p-len")
		testDeepEqualErrHandle(s47v1, s47v2, h, t, "equal-map-v47-p-len")
		testReleaseBytes(bs47)
		s47v2 = s47T{}
		s47v1.Mp = &m47v99
		bs47 = testMarshalErr(s47v1, h, t, "enc-map-v47-custom")
		testUnmarshalErr(&s47v2, bs47, h, t, "dec-map-v47-p-len")
		testDeepEqualErrHandle(s47v1, s47v2, h, t, "equal-map-v47-p-len")
		testReleaseBytes(bs47)
	}
	for _, v := range []map[uint64][]byte{nil, {}, {111: nil, 77: []byte("some-string-2")}} {
		var v48v1, v48v2 map[uint64][]byte
		var bs48 []byte
		v48v1 = v
		bs48 = testMarshalErr(v48v1, h, t, "enc-map-v48")
		if v != nil {
			v48v2 = make(map[uint64][]byte, len(v)) // reset map
			testUnmarshalErr(v48v2, bs48, h, t, "dec-map-v48")
			testDeepEqualErrHandle(v48v1, v48v2, h, t, "equal-map-v48")
			v48v2 = make(map[uint64][]byte, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v48v2), bs48, h, t, "dec-map-v48-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v48v1, v48v2, h, t, "equal-map-v48-noaddr")
		}
		if v == nil {
			v48v2 = nil
		} else {
			v48v2 = make(map[uint64][]byte, len(v))
		} // reset map
		testUnmarshalErr(&v48v2, bs48, h, t, "dec-map-v48-p-len")
		testDeepEqualErrHandle(v48v1, v48v2, h, t, "equal-map-v48-p-len")
		testReleaseBytes(bs48)
		bs48 = testMarshalErr(&v48v1, h, t, "enc-map-v48-p")
		v48v2 = nil
		testUnmarshalErr(&v48v2, bs48, h, t, "dec-map-v48-p-nil")
		testDeepEqualErrHandle(v48v1, v48v2, h, t, "equal-map-v48-p-nil")
		testReleaseBytes(bs48)
		// ...
		if v == nil {
			v48v2 = nil
		} else {
			v48v2 = make(map[uint64][]byte, len(v))
		} // reset map
		var v48v3, v48v4 typMapMapUint64Bytes
		v48v3 = typMapMapUint64Bytes(v48v1)
		v48v4 = typMapMapUint64Bytes(v48v2)
		if v != nil {
			bs48 = testMarshalErr(v48v3, h, t, "enc-map-v48-custom")
			testUnmarshalErr(v48v4, bs48, h, t, "dec-map-v48-p-len")
			testDeepEqualErrHandle(v48v3, v48v4, h, t, "equal-map-v48-p-len")
			testReleaseBytes(bs48)
		}
		type s48T struct {
			M  map[uint64][]byte
			Mp *map[uint64][]byte
		}
		var m48v99 = map[uint64][]byte{
			0:   nil,
			127: []byte("some-string-3"),
		}
		var s48v1, s48v2 s48T
		bs48 = testMarshalErr(s48v1, h, t, "enc-map-v48-custom")
		testUnmarshalErr(&s48v2, bs48, h, t, "dec-map-v48-p-len")
		testDeepEqualErrHandle(s48v1, s48v2, h, t, "equal-map-v48-p-len")
		testReleaseBytes(bs48)
		s48v2 = s48T{}
		s48v1.M = m48v99
		bs48 = testMarshalErr(s48v1, h, t, "enc-map-v48-custom")
		testUnmarshalErr(&s48v2, bs48, h, t, "dec-map-v48-p-len")
		testDeepEqualErrHandle(s48v1, s48v2, h, t, "equal-map-v48-p-len")
		testReleaseBytes(bs48)
		s48v2 = s48T{}
		s48v1.Mp = &m48v99
		bs48 = testMarshalErr(s48v1, h, t, "enc-map-v48-custom")
		testUnmarshalErr(&s48v2, bs48, h, t, "dec-map-v48-p-len")
		testDeepEqualErrHandle(s48v1, s48v2, h, t, "equal-map-v48-p-len")
		testReleaseBytes(bs48)
	}
	for _, v := range []map[uint64]uint8{nil, {}, {111: 0, 77: 127}} {
		var v49v1, v49v2 map[uint64]uint8
		var bs49 []byte
		v49v1 = v
		bs49 = testMarshalErr(v49v1, h, t, "enc-map-v49")
		if v != nil {
			v49v2 = make(map[uint64]uint8, len(v)) // reset map
			testUnmarshalErr(v49v2, bs49, h, t, "dec-map-v49")
			testDeepEqualErrHandle(v49v1, v49v2, h, t, "equal-map-v49")
			v49v2 = make(map[uint64]uint8, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v49v2), bs49, h, t, "dec-map-v49-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v49v1, v49v2, h, t, "equal-map-v49-noaddr")
		}
		if v == nil {
			v49v2 = nil
		} else {
			v49v2 = make(map[uint64]uint8, len(v))
		} // reset map
		testUnmarshalErr(&v49v2, bs49, h, t, "dec-map-v49-p-len")
		testDeepEqualErrHandle(v49v1, v49v2, h, t, "equal-map-v49-p-len")
		testReleaseBytes(bs49)
		bs49 = testMarshalErr(&v49v1, h, t, "enc-map-v49-p")
		v49v2 = nil
		testUnmarshalErr(&v49v2, bs49, h, t, "dec-map-v49-p-nil")
		testDeepEqualErrHandle(v49v1, v49v2, h, t, "equal-map-v49-p-nil")
		testReleaseBytes(bs49)
		// ...
		if v == nil {
			v49v2 = nil
		} else {
			v49v2 = make(map[uint64]uint8, len(v))
		} // reset map
		var v49v3, v49v4 typMapMapUint64Uint8
		v49v3 = typMapMapUint64Uint8(v49v1)
		v49v4 = typMapMapUint64Uint8(v49v2)
		if v != nil {
			bs49 = testMarshalErr(v49v3, h, t, "enc-map-v49-custom")
			testUnmarshalErr(v49v4, bs49, h, t, "dec-map-v49-p-len")
			testDeepEqualErrHandle(v49v3, v49v4, h, t, "equal-map-v49-p-len")
			testReleaseBytes(bs49)
		}
		type s49T struct {
			M  map[uint64]uint8
			Mp *map[uint64]uint8
		}
		var m49v99 = map[uint64]uint8{
			0:   0,
			111: 77,
		}
		var s49v1, s49v2 s49T
		bs49 = testMarshalErr(s49v1, h, t, "enc-map-v49-custom")
		testUnmarshalErr(&s49v2, bs49, h, t, "dec-map-v49-p-len")
		testDeepEqualErrHandle(s49v1, s49v2, h, t, "equal-map-v49-p-len")
		testReleaseBytes(bs49)
		s49v2 = s49T{}
		s49v1.M = m49v99
		bs49 = testMarshalErr(s49v1, h, t, "enc-map-v49-custom")
		testUnmarshalErr(&s49v2, bs49, h, t, "dec-map-v49-p-len")
		testDeepEqualErrHandle(s49v1, s49v2, h, t, "equal-map-v49-p-len")
		testReleaseBytes(bs49)
		s49v2 = s49T{}
		s49v1.Mp = &m49v99
		bs49 = testMarshalErr(s49v1, h, t, "enc-map-v49-custom")
		testUnmarshalErr(&s49v2, bs49, h, t, "dec-map-v49-p-len")
		testDeepEqualErrHandle(s49v1, s49v2, h, t, "equal-map-v49-p-len")
		testReleaseBytes(bs49)
	}
	for _, v := range []map[uint64]uint64{nil, {}, {127: 0, 111: 77}} {
		var v50v1, v50v2 map[uint64]uint64
		var bs50 []byte
		v50v1 = v
		bs50 = testMarshalErr(v50v1, h, t, "enc-map-v50")
		if v != nil {
			v50v2 = make(map[uint64]uint64, len(v)) // reset map
			testUnmarshalErr(v50v2, bs50, h, t, "dec-map-v50")
			testDeepEqualErrHandle(v50v1, v50v2, h, t, "equal-map-v50")
			v50v2 = make(map[uint64]uint64, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v50v2), bs50, h, t, "dec-map-v50-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v50v1, v50v2, h, t, "equal-map-v50-noaddr")
		}
		if v == nil {
			v50v2 = nil
		} else {
			v50v2 = make(map[uint64]uint64, len(v))
		} // reset map
		testUnmarshalErr(&v50v2, bs50, h, t, "dec-map-v50-p-len")
		testDeepEqualErrHandle(v50v1, v50v2, h, t, "equal-map-v50-p-len")
		testReleaseBytes(bs50)
		bs50 = testMarshalErr(&v50v1, h, t, "enc-map-v50-p")
		v50v2 = nil
		testUnmarshalErr(&v50v2, bs50, h, t, "dec-map-v50-p-nil")
		testDeepEqualErrHandle(v50v1, v50v2, h, t, "equal-map-v50-p-nil")
		testReleaseBytes(bs50)
		// ...
		if v == nil {
			v50v2 = nil
		} else {
			v50v2 = make(map[uint64]uint64, len(v))
		} // reset map
		var v50v3, v50v4 typMapMapUint64Uint64
		v50v3 = typMapMapUint64Uint64(v50v1)
		v50v4 = typMapMapUint64Uint64(v50v2)
		if v != nil {
			bs50 = testMarshalErr(v50v3, h, t, "enc-map-v50-custom")
			testUnmarshalErr(v50v4, bs50, h, t, "dec-map-v50-p-len")
			testDeepEqualErrHandle(v50v3, v50v4, h, t, "equal-map-v50-p-len")
			testReleaseBytes(bs50)
		}
		type s50T struct {
			M  map[uint64]uint64
			Mp *map[uint64]uint64
		}
		var m50v99 = map[uint64]uint64{
			0:   0,
			127: 111,
		}
		var s50v1, s50v2 s50T
		bs50 = testMarshalErr(s50v1, h, t, "enc-map-v50-custom")
		testUnmarshalErr(&s50v2, bs50, h, t, "dec-map-v50-p-len")
		testDeepEqualErrHandle(s50v1, s50v2, h, t, "equal-map-v50-p-len")
		testReleaseBytes(bs50)
		s50v2 = s50T{}
		s50v1.M = m50v99
		bs50 = testMarshalErr(s50v1, h, t, "enc-map-v50-custom")
		testUnmarshalErr(&s50v2, bs50, h, t, "dec-map-v50-p-len")
		testDeepEqualErrHandle(s50v1, s50v2, h, t, "equal-map-v50-p-len")
		testReleaseBytes(bs50)
		s50v2 = s50T{}
		s50v1.Mp = &m50v99
		bs50 = testMarshalErr(s50v1, h, t, "enc-map-v50-custom")
		testUnmarshalErr(&s50v2, bs50, h, t, "dec-map-v50-p-len")
		testDeepEqualErrHandle(s50v1, s50v2, h, t, "equal-map-v50-p-len")
		testReleaseBytes(bs50)
	}
	for _, v := range []map[uint64]int{nil, {}, {77: 0, 127: 111}} {
		var v51v1, v51v2 map[uint64]int
		var bs51 []byte
		v51v1 = v
		bs51 = testMarshalErr(v51v1, h, t, "enc-map-v51")
		if v != nil {
			v51v2 = make(map[uint64]int, len(v)) // reset map
			testUnmarshalErr(v51v2, bs51, h, t, "dec-map-v51")
			testDeepEqualErrHandle(v51v1, v51v2, h, t, "equal-map-v51")
			v51v2 = make(map[uint64]int, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v51v2), bs51, h, t, "dec-map-v51-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v51v1, v51v2, h, t, "equal-map-v51-noaddr")
		}
		if v == nil {
			v51v2 = nil
		} else {
			v51v2 = make(map[uint64]int, len(v))
		} // reset map
		testUnmarshalErr(&v51v2, bs51, h, t, "dec-map-v51-p-len")
		testDeepEqualErrHandle(v51v1, v51v2, h, t, "equal-map-v51-p-len")
		testReleaseBytes(bs51)
		bs51 = testMarshalErr(&v51v1, h, t, "enc-map-v51-p")
		v51v2 = nil
		testUnmarshalErr(&v51v2, bs51, h, t, "dec-map-v51-p-nil")
		testDeepEqualErrHandle(v51v1, v51v2, h, t, "equal-map-v51-p-nil")
		testReleaseBytes(bs51)
		// ...
		if v == nil {
			v51v2 = nil
		} else {
			v51v2 = make(map[uint64]int, len(v))
		} // reset map
		var v51v3, v51v4 typMapMapUint64Int
		v51v3 = typMapMapUint64Int(v51v1)
		v51v4 = typMapMapUint64Int(v51v2)
		if v != nil {
			bs51 = testMarshalErr(v51v3, h, t, "enc-map-v51-custom")
			testUnmarshalErr(v51v4, bs51, h, t, "dec-map-v51-p-len")
			testDeepEqualErrHandle(v51v3, v51v4, h, t, "equal-map-v51-p-len")
			testReleaseBytes(bs51)
		}
		type s51T struct {
			M  map[uint64]int
			Mp *map[uint64]int
		}
		var m51v99 = map[uint64]int{
			0:  0,
			77: 127,
		}
		var s51v1, s51v2 s51T
		bs51 = testMarshalErr(s51v1, h, t, "enc-map-v51-custom")
		testUnmarshalErr(&s51v2, bs51, h, t, "dec-map-v51-p-len")
		testDeepEqualErrHandle(s51v1, s51v2, h, t, "equal-map-v51-p-len")
		testReleaseBytes(bs51)
		s51v2 = s51T{}
		s51v1.M = m51v99
		bs51 = testMarshalErr(s51v1, h, t, "enc-map-v51-custom")
		testUnmarshalErr(&s51v2, bs51, h, t, "dec-map-v51-p-len")
		testDeepEqualErrHandle(s51v1, s51v2, h, t, "equal-map-v51-p-len")
		testReleaseBytes(bs51)
		s51v2 = s51T{}
		s51v1.Mp = &m51v99
		bs51 = testMarshalErr(s51v1, h, t, "enc-map-v51-custom")
		testUnmarshalErr(&s51v2, bs51, h, t, "dec-map-v51-p-len")
		testDeepEqualErrHandle(s51v1, s51v2, h, t, "equal-map-v51-p-len")
		testReleaseBytes(bs51)
	}
	for _, v := range []map[uint64]int32{nil, {}, {111: 0, 77: 127}} {
		var v52v1, v52v2 map[uint64]int32
		var bs52 []byte
		v52v1 = v
		bs52 = testMarshalErr(v52v1, h, t, "enc-map-v52")
		if v != nil {
			v52v2 = make(map[uint64]int32, len(v)) // reset map
			testUnmarshalErr(v52v2, bs52, h, t, "dec-map-v52")
			testDeepEqualErrHandle(v52v1, v52v2, h, t, "equal-map-v52")
			v52v2 = make(map[uint64]int32, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v52v2), bs52, h, t, "dec-map-v52-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v52v1, v52v2, h, t, "equal-map-v52-noaddr")
		}
		if v == nil {
			v52v2 = nil
		} else {
			v52v2 = make(map[uint64]int32, len(v))
		} // reset map
		testUnmarshalErr(&v52v2, bs52, h, t, "dec-map-v52-p-len")
		testDeepEqualErrHandle(v52v1, v52v2, h, t, "equal-map-v52-p-len")
		testReleaseBytes(bs52)
		bs52 = testMarshalErr(&v52v1, h, t, "enc-map-v52-p")
		v52v2 = nil
		testUnmarshalErr(&v52v2, bs52, h, t, "dec-map-v52-p-nil")
		testDeepEqualErrHandle(v52v1, v52v2, h, t, "equal-map-v52-p-nil")
		testReleaseBytes(bs52)
		// ...
		if v == nil {
			v52v2 = nil
		} else {
			v52v2 = make(map[uint64]int32, len(v))
		} // reset map
		var v52v3, v52v4 typMapMapUint64Int32
		v52v3 = typMapMapUint64Int32(v52v1)
		v52v4 = typMapMapUint64Int32(v52v2)
		if v != nil {
			bs52 = testMarshalErr(v52v3, h, t, "enc-map-v52-custom")
			testUnmarshalErr(v52v4, bs52, h, t, "dec-map-v52-p-len")
			testDeepEqualErrHandle(v52v3, v52v4, h, t, "equal-map-v52-p-len")
			testReleaseBytes(bs52)
		}
		type s52T struct {
			M  map[uint64]int32
			Mp *map[uint64]int32
		}
		var m52v99 = map[uint64]int32{
			0:   0,
			111: 77,
		}
		var s52v1, s52v2 s52T
		bs52 = testMarshalErr(s52v1, h, t, "enc-map-v52-custom")
		testUnmarshalErr(&s52v2, bs52, h, t, "dec-map-v52-p-len")
		testDeepEqualErrHandle(s52v1, s52v2, h, t, "equal-map-v52-p-len")
		testReleaseBytes(bs52)
		s52v2 = s52T{}
		s52v1.M = m52v99
		bs52 = testMarshalErr(s52v1, h, t, "enc-map-v52-custom")
		testUnmarshalErr(&s52v2, bs52, h, t, "dec-map-v52-p-len")
		testDeepEqualErrHandle(s52v1, s52v2, h, t, "equal-map-v52-p-len")
		testReleaseBytes(bs52)
		s52v2 = s52T{}
		s52v1.Mp = &m52v99
		bs52 = testMarshalErr(s52v1, h, t, "enc-map-v52-custom")
		testUnmarshalErr(&s52v2, bs52, h, t, "dec-map-v52-p-len")
		testDeepEqualErrHandle(s52v1, s52v2, h, t, "equal-map-v52-p-len")
		testReleaseBytes(bs52)
	}
	for _, v := range []map[uint64]float64{nil, {}, {127: 0, 111: 11.1}} {
		var v53v1, v53v2 map[uint64]float64
		var bs53 []byte
		v53v1 = v
		bs53 = testMarshalErr(v53v1, h, t, "enc-map-v53")
		if v != nil {
			v53v2 = make(map[uint64]float64, len(v)) // reset map
			testUnmarshalErr(v53v2, bs53, h, t, "dec-map-v53")
			testDeepEqualErrHandle(v53v1, v53v2, h, t, "equal-map-v53")
			v53v2 = make(map[uint64]float64, len(v))                                   // reset map
			testUnmarshalErr(reflect.ValueOf(v53v2), bs53, h, t, "dec-map-v53-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v53v1, v53v2, h, t, "equal-map-v53-noaddr")
		}
		if v == nil {
			v53v2 = nil
		} else {
			v53v2 = make(map[uint64]float64, len(v))
		} // reset map
		testUnmarshalErr(&v53v2, bs53, h, t, "dec-map-v53-p-len")
		testDeepEqualErrHandle(v53v1, v53v2, h, t, "equal-map-v53-p-len")
		testReleaseBytes(bs53)
		bs53 = testMarshalErr(&v53v1, h, t, "enc-map-v53-p")
		v53v2 = nil
		testUnmarshalErr(&v53v2, bs53, h, t, "dec-map-v53-p-nil")
		testDeepEqualErrHandle(v53v1, v53v2, h, t, "equal-map-v53-p-nil")
		testReleaseBytes(bs53)
		// ...
		if v == nil {
			v53v2 = nil
		} else {
			v53v2 = make(map[uint64]float64, len(v))
		} // reset map
		var v53v3, v53v4 typMapMapUint64Float64
		v53v3 = typMapMapUint64Float64(v53v1)
		v53v4 = typMapMapUint64Float64(v53v2)
		if v != nil {
			bs53 = testMarshalErr(v53v3, h, t, "enc-map-v53-custom")
			testUnmarshalErr(v53v4, bs53, h, t, "dec-map-v53-p-len")
			testDeepEqualErrHandle(v53v3, v53v4, h, t, "equal-map-v53-p-len")
			testReleaseBytes(bs53)
		}
		type s53T struct {
			M  map[uint64]float64
			Mp *map[uint64]float64
		}
		var m53v99 = map[uint64]float64{
			0:  0,
			77: 22.2,
		}
		var s53v1, s53v2 s53T
		bs53 = testMarshalErr(s53v1, h, t, "enc-map-v53-custom")
		testUnmarshalErr(&s53v2, bs53, h, t, "dec-map-v53-p-len")
		testDeepEqualErrHandle(s53v1, s53v2, h, t, "equal-map-v53-p-len")
		testReleaseBytes(bs53)
		s53v2 = s53T{}
		s53v1.M = m53v99
		bs53 = testMarshalErr(s53v1, h, t, "enc-map-v53-custom")
		testUnmarshalErr(&s53v2, bs53, h, t, "dec-map-v53-p-len")
		testDeepEqualErrHandle(s53v1, s53v2, h, t, "equal-map-v53-p-len")
		testReleaseBytes(bs53)
		s53v2 = s53T{}
		s53v1.Mp = &m53v99
		bs53 = testMarshalErr(s53v1, h, t, "enc-map-v53-custom")
		testUnmarshalErr(&s53v2, bs53, h, t, "dec-map-v53-p-len")
		testDeepEqualErrHandle(s53v1, s53v2, h, t, "equal-map-v53-p-len")
		testReleaseBytes(bs53)
	}
	for _, v := range []map[uint64]bool{nil, {}, {127: false, 111: false}} {
		var v54v1, v54v2 map[uint64]bool
		var bs54 []byte
		v54v1 = v
		bs54 = testMarshalErr(v54v1, h, t, "enc-map-v54")
		if v != nil {
			v54v2 = make(map[uint64]bool, len(v)) // reset map
			testUnmarshalErr(v54v2, bs54, h, t, "dec-map-v54")
			testDeepEqualErrHandle(v54v1, v54v2, h, t, "equal-map-v54")
			v54v2 = make(map[uint64]bool, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v54v2), bs54, h, t, "dec-map-v54-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v54v1, v54v2, h, t, "equal-map-v54-noaddr")
		}
		if v == nil {
			v54v2 = nil
		} else {
			v54v2 = make(map[uint64]bool, len(v))
		} // reset map
		testUnmarshalErr(&v54v2, bs54, h, t, "dec-map-v54-p-len")
		testDeepEqualErrHandle(v54v1, v54v2, h, t, "equal-map-v54-p-len")
		testReleaseBytes(bs54)
		bs54 = testMarshalErr(&v54v1, h, t, "enc-map-v54-p")
		v54v2 = nil
		testUnmarshalErr(&v54v2, bs54, h, t, "dec-map-v54-p-nil")
		testDeepEqualErrHandle(v54v1, v54v2, h, t, "equal-map-v54-p-nil")
		testReleaseBytes(bs54)
		// ...
		if v == nil {
			v54v2 = nil
		} else {
			v54v2 = make(map[uint64]bool, len(v))
		} // reset map
		var v54v3, v54v4 typMapMapUint64Bool
		v54v3 = typMapMapUint64Bool(v54v1)
		v54v4 = typMapMapUint64Bool(v54v2)
		if v != nil {
			bs54 = testMarshalErr(v54v3, h, t, "enc-map-v54-custom")
			testUnmarshalErr(v54v4, bs54, h, t, "dec-map-v54-p-len")
			testDeepEqualErrHandle(v54v3, v54v4, h, t, "equal-map-v54-p-len")
			testReleaseBytes(bs54)
		}
		type s54T struct {
			M  map[uint64]bool
			Mp *map[uint64]bool
		}
		var m54v99 = map[uint64]bool{
			0:  false,
			77: true,
		}
		var s54v1, s54v2 s54T
		bs54 = testMarshalErr(s54v1, h, t, "enc-map-v54-custom")
		testUnmarshalErr(&s54v2, bs54, h, t, "dec-map-v54-p-len")
		testDeepEqualErrHandle(s54v1, s54v2, h, t, "equal-map-v54-p-len")
		testReleaseBytes(bs54)
		s54v2 = s54T{}
		s54v1.M = m54v99
		bs54 = testMarshalErr(s54v1, h, t, "enc-map-v54-custom")
		testUnmarshalErr(&s54v2, bs54, h, t, "dec-map-v54-p-len")
		testDeepEqualErrHandle(s54v1, s54v2, h, t, "equal-map-v54-p-len")
		testReleaseBytes(bs54)
		s54v2 = s54T{}
		s54v1.Mp = &m54v99
		bs54 = testMarshalErr(s54v1, h, t, "enc-map-v54-custom")
		testUnmarshalErr(&s54v2, bs54, h, t, "dec-map-v54-p-len")
		testDeepEqualErrHandle(s54v1, s54v2, h, t, "equal-map-v54-p-len")
		testReleaseBytes(bs54)
	}
	for _, v := range []map[int]interface{}{nil, {}, {127: nil, 111: "string-is-an-interface-1"}} {
		var v55v1, v55v2 map[int]interface{}
		var bs55 []byte
		v55v1 = v
		bs55 = testMarshalErr(v55v1, h, t, "enc-map-v55")
		if v != nil {
			v55v2 = make(map[int]interface{}, len(v)) // reset map
			testUnmarshalErr(v55v2, bs55, h, t, "dec-map-v55")
			testDeepEqualErrHandle(v55v1, v55v2, h, t, "equal-map-v55")
			v55v2 = make(map[int]interface{}, len(v))                                  // reset map
			testUnmarshalErr(reflect.ValueOf(v55v2), bs55, h, t, "dec-map-v55-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v55v1, v55v2, h, t, "equal-map-v55-noaddr")
		}
		if v == nil {
			v55v2 = nil
		} else {
			v55v2 = make(map[int]interface{}, len(v))
		} // reset map
		testUnmarshalErr(&v55v2, bs55, h, t, "dec-map-v55-p-len")
		testDeepEqualErrHandle(v55v1, v55v2, h, t, "equal-map-v55-p-len")
		testReleaseBytes(bs55)
		bs55 = testMarshalErr(&v55v1, h, t, "enc-map-v55-p")
		v55v2 = nil
		testUnmarshalErr(&v55v2, bs55, h, t, "dec-map-v55-p-nil")
		testDeepEqualErrHandle(v55v1, v55v2, h, t, "equal-map-v55-p-nil")
		testReleaseBytes(bs55)
		// ...
		if v == nil {
			v55v2 = nil
		} else {
			v55v2 = make(map[int]interface{}, len(v))
		} // reset map
		var v55v3, v55v4 typMapMapIntIntf
		v55v3 = typMapMapIntIntf(v55v1)
		v55v4 = typMapMapIntIntf(v55v2)
		if v != nil {
			bs55 = testMarshalErr(v55v3, h, t, "enc-map-v55-custom")
			testUnmarshalErr(v55v4, bs55, h, t, "dec-map-v55-p-len")
			testDeepEqualErrHandle(v55v3, v55v4, h, t, "equal-map-v55-p-len")
			testReleaseBytes(bs55)
		}
		type s55T struct {
			M  map[int]interface{}
			Mp *map[int]interface{}
		}
		var m55v99 = map[int]interface{}{
			0:  nil,
			77: "string-is-an-interface-2",
		}
		var s55v1, s55v2 s55T
		bs55 = testMarshalErr(s55v1, h, t, "enc-map-v55-custom")
		testUnmarshalErr(&s55v2, bs55, h, t, "dec-map-v55-p-len")
		testDeepEqualErrHandle(s55v1, s55v2, h, t, "equal-map-v55-p-len")
		testReleaseBytes(bs55)
		s55v2 = s55T{}
		s55v1.M = m55v99
		bs55 = testMarshalErr(s55v1, h, t, "enc-map-v55-custom")
		testUnmarshalErr(&s55v2, bs55, h, t, "dec-map-v55-p-len")
		testDeepEqualErrHandle(s55v1, s55v2, h, t, "equal-map-v55-p-len")
		testReleaseBytes(bs55)
		s55v2 = s55T{}
		s55v1.Mp = &m55v99
		bs55 = testMarshalErr(s55v1, h, t, "enc-map-v55-custom")
		testUnmarshalErr(&s55v2, bs55, h, t, "dec-map-v55-p-len")
		testDeepEqualErrHandle(s55v1, s55v2, h, t, "equal-map-v55-p-len")
		testReleaseBytes(bs55)
	}
	for _, v := range []map[int]string{nil, {}, {127: "", 111: "some-string-1"}} {
		var v56v1, v56v2 map[int]string
		var bs56 []byte
		v56v1 = v
		bs56 = testMarshalErr(v56v1, h, t, "enc-map-v56")
		if v != nil {
			v56v2 = make(map[int]string, len(v)) // reset map
			testUnmarshalErr(v56v2, bs56, h, t, "dec-map-v56")
			testDeepEqualErrHandle(v56v1, v56v2, h, t, "equal-map-v56")
			v56v2 = make(map[int]string, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v56v2), bs56, h, t, "dec-map-v56-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v56v1, v56v2, h, t, "equal-map-v56-noaddr")
		}
		if v == nil {
			v56v2 = nil
		} else {
			v56v2 = make(map[int]string, len(v))
		} // reset map
		testUnmarshalErr(&v56v2, bs56, h, t, "dec-map-v56-p-len")
		testDeepEqualErrHandle(v56v1, v56v2, h, t, "equal-map-v56-p-len")
		testReleaseBytes(bs56)
		bs56 = testMarshalErr(&v56v1, h, t, "enc-map-v56-p")
		v56v2 = nil
		testUnmarshalErr(&v56v2, bs56, h, t, "dec-map-v56-p-nil")
		testDeepEqualErrHandle(v56v1, v56v2, h, t, "equal-map-v56-p-nil")
		testReleaseBytes(bs56)
		// ...
		if v == nil {
			v56v2 = nil
		} else {
			v56v2 = make(map[int]string, len(v))
		} // reset map
		var v56v3, v56v4 typMapMapIntString
		v56v3 = typMapMapIntString(v56v1)
		v56v4 = typMapMapIntString(v56v2)
		if v != nil {
			bs56 = testMarshalErr(v56v3, h, t, "enc-map-v56-custom")
			testUnmarshalErr(v56v4, bs56, h, t, "dec-map-v56-p-len")
			testDeepEqualErrHandle(v56v3, v56v4, h, t, "equal-map-v56-p-len")
			testReleaseBytes(bs56)
		}
		type s56T struct {
			M  map[int]string
			Mp *map[int]string
		}
		var m56v99 = map[int]string{
			0:  "",
			77: "some-string-2",
		}
		var s56v1, s56v2 s56T
		bs56 = testMarshalErr(s56v1, h, t, "enc-map-v56-custom")
		testUnmarshalErr(&s56v2, bs56, h, t, "dec-map-v56-p-len")
		testDeepEqualErrHandle(s56v1, s56v2, h, t, "equal-map-v56-p-len")
		testReleaseBytes(bs56)
		s56v2 = s56T{}
		s56v1.M = m56v99
		bs56 = testMarshalErr(s56v1, h, t, "enc-map-v56-custom")
		testUnmarshalErr(&s56v2, bs56, h, t, "dec-map-v56-p-len")
		testDeepEqualErrHandle(s56v1, s56v2, h, t, "equal-map-v56-p-len")
		testReleaseBytes(bs56)
		s56v2 = s56T{}
		s56v1.Mp = &m56v99
		bs56 = testMarshalErr(s56v1, h, t, "enc-map-v56-custom")
		testUnmarshalErr(&s56v2, bs56, h, t, "dec-map-v56-p-len")
		testDeepEqualErrHandle(s56v1, s56v2, h, t, "equal-map-v56-p-len")
		testReleaseBytes(bs56)
	}
	for _, v := range []map[int][]byte{nil, {}, {127: nil, 111: []byte("some-string-1")}} {
		var v57v1, v57v2 map[int][]byte
		var bs57 []byte
		v57v1 = v
		bs57 = testMarshalErr(v57v1, h, t, "enc-map-v57")
		if v != nil {
			v57v2 = make(map[int][]byte, len(v)) // reset map
			testUnmarshalErr(v57v2, bs57, h, t, "dec-map-v57")
			testDeepEqualErrHandle(v57v1, v57v2, h, t, "equal-map-v57")
			v57v2 = make(map[int][]byte, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v57v2), bs57, h, t, "dec-map-v57-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v57v1, v57v2, h, t, "equal-map-v57-noaddr")
		}
		if v == nil {
			v57v2 = nil
		} else {
			v57v2 = make(map[int][]byte, len(v))
		} // reset map
		testUnmarshalErr(&v57v2, bs57, h, t, "dec-map-v57-p-len")
		testDeepEqualErrHandle(v57v1, v57v2, h, t, "equal-map-v57-p-len")
		testReleaseBytes(bs57)
		bs57 = testMarshalErr(&v57v1, h, t, "enc-map-v57-p")
		v57v2 = nil
		testUnmarshalErr(&v57v2, bs57, h, t, "dec-map-v57-p-nil")
		testDeepEqualErrHandle(v57v1, v57v2, h, t, "equal-map-v57-p-nil")
		testReleaseBytes(bs57)
		// ...
		if v == nil {
			v57v2 = nil
		} else {
			v57v2 = make(map[int][]byte, len(v))
		} // reset map
		var v57v3, v57v4 typMapMapIntBytes
		v57v3 = typMapMapIntBytes(v57v1)
		v57v4 = typMapMapIntBytes(v57v2)
		if v != nil {
			bs57 = testMarshalErr(v57v3, h, t, "enc-map-v57-custom")
			testUnmarshalErr(v57v4, bs57, h, t, "dec-map-v57-p-len")
			testDeepEqualErrHandle(v57v3, v57v4, h, t, "equal-map-v57-p-len")
			testReleaseBytes(bs57)
		}
		type s57T struct {
			M  map[int][]byte
			Mp *map[int][]byte
		}
		var m57v99 = map[int][]byte{
			0:  nil,
			77: []byte("some-string-2"),
		}
		var s57v1, s57v2 s57T
		bs57 = testMarshalErr(s57v1, h, t, "enc-map-v57-custom")
		testUnmarshalErr(&s57v2, bs57, h, t, "dec-map-v57-p-len")
		testDeepEqualErrHandle(s57v1, s57v2, h, t, "equal-map-v57-p-len")
		testReleaseBytes(bs57)
		s57v2 = s57T{}
		s57v1.M = m57v99
		bs57 = testMarshalErr(s57v1, h, t, "enc-map-v57-custom")
		testUnmarshalErr(&s57v2, bs57, h, t, "dec-map-v57-p-len")
		testDeepEqualErrHandle(s57v1, s57v2, h, t, "equal-map-v57-p-len")
		testReleaseBytes(bs57)
		s57v2 = s57T{}
		s57v1.Mp = &m57v99
		bs57 = testMarshalErr(s57v1, h, t, "enc-map-v57-custom")
		testUnmarshalErr(&s57v2, bs57, h, t, "dec-map-v57-p-len")
		testDeepEqualErrHandle(s57v1, s57v2, h, t, "equal-map-v57-p-len")
		testReleaseBytes(bs57)
	}
	for _, v := range []map[int]uint8{nil, {}, {127: 0, 111: 77}} {
		var v58v1, v58v2 map[int]uint8
		var bs58 []byte
		v58v1 = v
		bs58 = testMarshalErr(v58v1, h, t, "enc-map-v58")
		if v != nil {
			v58v2 = make(map[int]uint8, len(v)) // reset map
			testUnmarshalErr(v58v2, bs58, h, t, "dec-map-v58")
			testDeepEqualErrHandle(v58v1, v58v2, h, t, "equal-map-v58")
			v58v2 = make(map[int]uint8, len(v))                                        // reset map
			testUnmarshalErr(reflect.ValueOf(v58v2), bs58, h, t, "dec-map-v58-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v58v1, v58v2, h, t, "equal-map-v58-noaddr")
		}
		if v == nil {
			v58v2 = nil
		} else {
			v58v2 = make(map[int]uint8, len(v))
		} // reset map
		testUnmarshalErr(&v58v2, bs58, h, t, "dec-map-v58-p-len")
		testDeepEqualErrHandle(v58v1, v58v2, h, t, "equal-map-v58-p-len")
		testReleaseBytes(bs58)
		bs58 = testMarshalErr(&v58v1, h, t, "enc-map-v58-p")
		v58v2 = nil
		testUnmarshalErr(&v58v2, bs58, h, t, "dec-map-v58-p-nil")
		testDeepEqualErrHandle(v58v1, v58v2, h, t, "equal-map-v58-p-nil")
		testReleaseBytes(bs58)
		// ...
		if v == nil {
			v58v2 = nil
		} else {
			v58v2 = make(map[int]uint8, len(v))
		} // reset map
		var v58v3, v58v4 typMapMapIntUint8
		v58v3 = typMapMapIntUint8(v58v1)
		v58v4 = typMapMapIntUint8(v58v2)
		if v != nil {
			bs58 = testMarshalErr(v58v3, h, t, "enc-map-v58-custom")
			testUnmarshalErr(v58v4, bs58, h, t, "dec-map-v58-p-len")
			testDeepEqualErrHandle(v58v3, v58v4, h, t, "equal-map-v58-p-len")
			testReleaseBytes(bs58)
		}
		type s58T struct {
			M  map[int]uint8
			Mp *map[int]uint8
		}
		var m58v99 = map[int]uint8{
			0:   0,
			127: 111,
		}
		var s58v1, s58v2 s58T
		bs58 = testMarshalErr(s58v1, h, t, "enc-map-v58-custom")
		testUnmarshalErr(&s58v2, bs58, h, t, "dec-map-v58-p-len")
		testDeepEqualErrHandle(s58v1, s58v2, h, t, "equal-map-v58-p-len")
		testReleaseBytes(bs58)
		s58v2 = s58T{}
		s58v1.M = m58v99
		bs58 = testMarshalErr(s58v1, h, t, "enc-map-v58-custom")
		testUnmarshalErr(&s58v2, bs58, h, t, "dec-map-v58-p-len")
		testDeepEqualErrHandle(s58v1, s58v2, h, t, "equal-map-v58-p-len")
		testReleaseBytes(bs58)
		s58v2 = s58T{}
		s58v1.Mp = &m58v99
		bs58 = testMarshalErr(s58v1, h, t, "enc-map-v58-custom")
		testUnmarshalErr(&s58v2, bs58, h, t, "dec-map-v58-p-len")
		testDeepEqualErrHandle(s58v1, s58v2, h, t, "equal-map-v58-p-len")
		testReleaseBytes(bs58)
	}
	for _, v := range []map[int]uint64{nil, {}, {77: 0, 127: 111}} {
		var v59v1, v59v2 map[int]uint64
		var bs59 []byte
		v59v1 = v
		bs59 = testMarshalErr(v59v1, h, t, "enc-map-v59")
		if v != nil {
			v59v2 = make(map[int]uint64, len(v)) // reset map
			testUnmarshalErr(v59v2, bs59, h, t, "dec-map-v59")
			testDeepEqualErrHandle(v59v1, v59v2, h, t, "equal-map-v59")
			v59v2 = make(map[int]uint64, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v59v2), bs59, h, t, "dec-map-v59-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v59v1, v59v2, h, t, "equal-map-v59-noaddr")
		}
		if v == nil {
			v59v2 = nil
		} else {
			v59v2 = make(map[int]uint64, len(v))
		} // reset map
		testUnmarshalErr(&v59v2, bs59, h, t, "dec-map-v59-p-len")
		testDeepEqualErrHandle(v59v1, v59v2, h, t, "equal-map-v59-p-len")
		testReleaseBytes(bs59)
		bs59 = testMarshalErr(&v59v1, h, t, "enc-map-v59-p")
		v59v2 = nil
		testUnmarshalErr(&v59v2, bs59, h, t, "dec-map-v59-p-nil")
		testDeepEqualErrHandle(v59v1, v59v2, h, t, "equal-map-v59-p-nil")
		testReleaseBytes(bs59)
		// ...
		if v == nil {
			v59v2 = nil
		} else {
			v59v2 = make(map[int]uint64, len(v))
		} // reset map
		var v59v3, v59v4 typMapMapIntUint64
		v59v3 = typMapMapIntUint64(v59v1)
		v59v4 = typMapMapIntUint64(v59v2)
		if v != nil {
			bs59 = testMarshalErr(v59v3, h, t, "enc-map-v59-custom")
			testUnmarshalErr(v59v4, bs59, h, t, "dec-map-v59-p-len")
			testDeepEqualErrHandle(v59v3, v59v4, h, t, "equal-map-v59-p-len")
			testReleaseBytes(bs59)
		}
		type s59T struct {
			M  map[int]uint64
			Mp *map[int]uint64
		}
		var m59v99 = map[int]uint64{
			0:  0,
			77: 127,
		}
		var s59v1, s59v2 s59T
		bs59 = testMarshalErr(s59v1, h, t, "enc-map-v59-custom")
		testUnmarshalErr(&s59v2, bs59, h, t, "dec-map-v59-p-len")
		testDeepEqualErrHandle(s59v1, s59v2, h, t, "equal-map-v59-p-len")
		testReleaseBytes(bs59)
		s59v2 = s59T{}
		s59v1.M = m59v99
		bs59 = testMarshalErr(s59v1, h, t, "enc-map-v59-custom")
		testUnmarshalErr(&s59v2, bs59, h, t, "dec-map-v59-p-len")
		testDeepEqualErrHandle(s59v1, s59v2, h, t, "equal-map-v59-p-len")
		testReleaseBytes(bs59)
		s59v2 = s59T{}
		s59v1.Mp = &m59v99
		bs59 = testMarshalErr(s59v1, h, t, "enc-map-v59-custom")
		testUnmarshalErr(&s59v2, bs59, h, t, "dec-map-v59-p-len")
		testDeepEqualErrHandle(s59v1, s59v2, h, t, "equal-map-v59-p-len")
		testReleaseBytes(bs59)
	}
	for _, v := range []map[int]int{nil, {}, {111: 0, 77: 127}} {
		var v60v1, v60v2 map[int]int
		var bs60 []byte
		v60v1 = v
		bs60 = testMarshalErr(v60v1, h, t, "enc-map-v60")
		if v != nil {
			v60v2 = make(map[int]int, len(v)) // reset map
			testUnmarshalErr(v60v2, bs60, h, t, "dec-map-v60")
			testDeepEqualErrHandle(v60v1, v60v2, h, t, "equal-map-v60")
			v60v2 = make(map[int]int, len(v))                                          // reset map
			testUnmarshalErr(reflect.ValueOf(v60v2), bs60, h, t, "dec-map-v60-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v60v1, v60v2, h, t, "equal-map-v60-noaddr")
		}
		if v == nil {
			v60v2 = nil
		} else {
			v60v2 = make(map[int]int, len(v))
		} // reset map
		testUnmarshalErr(&v60v2, bs60, h, t, "dec-map-v60-p-len")
		testDeepEqualErrHandle(v60v1, v60v2, h, t, "equal-map-v60-p-len")
		testReleaseBytes(bs60)
		bs60 = testMarshalErr(&v60v1, h, t, "enc-map-v60-p")
		v60v2 = nil
		testUnmarshalErr(&v60v2, bs60, h, t, "dec-map-v60-p-nil")
		testDeepEqualErrHandle(v60v1, v60v2, h, t, "equal-map-v60-p-nil")
		testReleaseBytes(bs60)
		// ...
		if v == nil {
			v60v2 = nil
		} else {
			v60v2 = make(map[int]int, len(v))
		} // reset map
		var v60v3, v60v4 typMapMapIntInt
		v60v3 = typMapMapIntInt(v60v1)
		v60v4 = typMapMapIntInt(v60v2)
		if v != nil {
			bs60 = testMarshalErr(v60v3, h, t, "enc-map-v60-custom")
			testUnmarshalErr(v60v4, bs60, h, t, "dec-map-v60-p-len")
			testDeepEqualErrHandle(v60v3, v60v4, h, t, "equal-map-v60-p-len")
			testReleaseBytes(bs60)
		}
		type s60T struct {
			M  map[int]int
			Mp *map[int]int
		}
		var m60v99 = map[int]int{
			0:   0,
			111: 77,
		}
		var s60v1, s60v2 s60T
		bs60 = testMarshalErr(s60v1, h, t, "enc-map-v60-custom")
		testUnmarshalErr(&s60v2, bs60, h, t, "dec-map-v60-p-len")
		testDeepEqualErrHandle(s60v1, s60v2, h, t, "equal-map-v60-p-len")
		testReleaseBytes(bs60)
		s60v2 = s60T{}
		s60v1.M = m60v99
		bs60 = testMarshalErr(s60v1, h, t, "enc-map-v60-custom")
		testUnmarshalErr(&s60v2, bs60, h, t, "dec-map-v60-p-len")
		testDeepEqualErrHandle(s60v1, s60v2, h, t, "equal-map-v60-p-len")
		testReleaseBytes(bs60)
		s60v2 = s60T{}
		s60v1.Mp = &m60v99
		bs60 = testMarshalErr(s60v1, h, t, "enc-map-v60-custom")
		testUnmarshalErr(&s60v2, bs60, h, t, "dec-map-v60-p-len")
		testDeepEqualErrHandle(s60v1, s60v2, h, t, "equal-map-v60-p-len")
		testReleaseBytes(bs60)
	}
	for _, v := range []map[int]int32{nil, {}, {127: 0, 111: 77}} {
		var v61v1, v61v2 map[int]int32
		var bs61 []byte
		v61v1 = v
		bs61 = testMarshalErr(v61v1, h, t, "enc-map-v61")
		if v != nil {
			v61v2 = make(map[int]int32, len(v)) // reset map
			testUnmarshalErr(v61v2, bs61, h, t, "dec-map-v61")
			testDeepEqualErrHandle(v61v1, v61v2, h, t, "equal-map-v61")
			v61v2 = make(map[int]int32, len(v))                                        // reset map
			testUnmarshalErr(reflect.ValueOf(v61v2), bs61, h, t, "dec-map-v61-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v61v1, v61v2, h, t, "equal-map-v61-noaddr")
		}
		if v == nil {
			v61v2 = nil
		} else {
			v61v2 = make(map[int]int32, len(v))
		} // reset map
		testUnmarshalErr(&v61v2, bs61, h, t, "dec-map-v61-p-len")
		testDeepEqualErrHandle(v61v1, v61v2, h, t, "equal-map-v61-p-len")
		testReleaseBytes(bs61)
		bs61 = testMarshalErr(&v61v1, h, t, "enc-map-v61-p")
		v61v2 = nil
		testUnmarshalErr(&v61v2, bs61, h, t, "dec-map-v61-p-nil")
		testDeepEqualErrHandle(v61v1, v61v2, h, t, "equal-map-v61-p-nil")
		testReleaseBytes(bs61)
		// ...
		if v == nil {
			v61v2 = nil
		} else {
			v61v2 = make(map[int]int32, len(v))
		} // reset map
		var v61v3, v61v4 typMapMapIntInt32
		v61v3 = typMapMapIntInt32(v61v1)
		v61v4 = typMapMapIntInt32(v61v2)
		if v != nil {
			bs61 = testMarshalErr(v61v3, h, t, "enc-map-v61-custom")
			testUnmarshalErr(v61v4, bs61, h, t, "dec-map-v61-p-len")
			testDeepEqualErrHandle(v61v3, v61v4, h, t, "equal-map-v61-p-len")
			testReleaseBytes(bs61)
		}
		type s61T struct {
			M  map[int]int32
			Mp *map[int]int32
		}
		var m61v99 = map[int]int32{
			0:   0,
			127: 111,
		}
		var s61v1, s61v2 s61T
		bs61 = testMarshalErr(s61v1, h, t, "enc-map-v61-custom")
		testUnmarshalErr(&s61v2, bs61, h, t, "dec-map-v61-p-len")
		testDeepEqualErrHandle(s61v1, s61v2, h, t, "equal-map-v61-p-len")
		testReleaseBytes(bs61)
		s61v2 = s61T{}
		s61v1.M = m61v99
		bs61 = testMarshalErr(s61v1, h, t, "enc-map-v61-custom")
		testUnmarshalErr(&s61v2, bs61, h, t, "dec-map-v61-p-len")
		testDeepEqualErrHandle(s61v1, s61v2, h, t, "equal-map-v61-p-len")
		testReleaseBytes(bs61)
		s61v2 = s61T{}
		s61v1.Mp = &m61v99
		bs61 = testMarshalErr(s61v1, h, t, "enc-map-v61-custom")
		testUnmarshalErr(&s61v2, bs61, h, t, "dec-map-v61-p-len")
		testDeepEqualErrHandle(s61v1, s61v2, h, t, "equal-map-v61-p-len")
		testReleaseBytes(bs61)
	}
	for _, v := range []map[int]float64{nil, {}, {77: 0, 127: 33.3e3}} {
		var v62v1, v62v2 map[int]float64
		var bs62 []byte
		v62v1 = v
		bs62 = testMarshalErr(v62v1, h, t, "enc-map-v62")
		if v != nil {
			v62v2 = make(map[int]float64, len(v)) // reset map
			testUnmarshalErr(v62v2, bs62, h, t, "dec-map-v62")
			testDeepEqualErrHandle(v62v1, v62v2, h, t, "equal-map-v62")
			v62v2 = make(map[int]float64, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v62v2), bs62, h, t, "dec-map-v62-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v62v1, v62v2, h, t, "equal-map-v62-noaddr")
		}
		if v == nil {
			v62v2 = nil
		} else {
			v62v2 = make(map[int]float64, len(v))
		} // reset map
		testUnmarshalErr(&v62v2, bs62, h, t, "dec-map-v62-p-len")
		testDeepEqualErrHandle(v62v1, v62v2, h, t, "equal-map-v62-p-len")
		testReleaseBytes(bs62)
		bs62 = testMarshalErr(&v62v1, h, t, "enc-map-v62-p")
		v62v2 = nil
		testUnmarshalErr(&v62v2, bs62, h, t, "dec-map-v62-p-nil")
		testDeepEqualErrHandle(v62v1, v62v2, h, t, "equal-map-v62-p-nil")
		testReleaseBytes(bs62)
		// ...
		if v == nil {
			v62v2 = nil
		} else {
			v62v2 = make(map[int]float64, len(v))
		} // reset map
		var v62v3, v62v4 typMapMapIntFloat64
		v62v3 = typMapMapIntFloat64(v62v1)
		v62v4 = typMapMapIntFloat64(v62v2)
		if v != nil {
			bs62 = testMarshalErr(v62v3, h, t, "enc-map-v62-custom")
			testUnmarshalErr(v62v4, bs62, h, t, "dec-map-v62-p-len")
			testDeepEqualErrHandle(v62v3, v62v4, h, t, "equal-map-v62-p-len")
			testReleaseBytes(bs62)
		}
		type s62T struct {
			M  map[int]float64
			Mp *map[int]float64
		}
		var m62v99 = map[int]float64{
			0:   0,
			111: 11.1,
		}
		var s62v1, s62v2 s62T
		bs62 = testMarshalErr(s62v1, h, t, "enc-map-v62-custom")
		testUnmarshalErr(&s62v2, bs62, h, t, "dec-map-v62-p-len")
		testDeepEqualErrHandle(s62v1, s62v2, h, t, "equal-map-v62-p-len")
		testReleaseBytes(bs62)
		s62v2 = s62T{}
		s62v1.M = m62v99
		bs62 = testMarshalErr(s62v1, h, t, "enc-map-v62-custom")
		testUnmarshalErr(&s62v2, bs62, h, t, "dec-map-v62-p-len")
		testDeepEqualErrHandle(s62v1, s62v2, h, t, "equal-map-v62-p-len")
		testReleaseBytes(bs62)
		s62v2 = s62T{}
		s62v1.Mp = &m62v99
		bs62 = testMarshalErr(s62v1, h, t, "enc-map-v62-custom")
		testUnmarshalErr(&s62v2, bs62, h, t, "dec-map-v62-p-len")
		testDeepEqualErrHandle(s62v1, s62v2, h, t, "equal-map-v62-p-len")
		testReleaseBytes(bs62)
	}
	for _, v := range []map[int]bool{nil, {}, {77: false, 127: true}} {
		var v63v1, v63v2 map[int]bool
		var bs63 []byte
		v63v1 = v
		bs63 = testMarshalErr(v63v1, h, t, "enc-map-v63")
		if v != nil {
			v63v2 = make(map[int]bool, len(v)) // reset map
			testUnmarshalErr(v63v2, bs63, h, t, "dec-map-v63")
			testDeepEqualErrHandle(v63v1, v63v2, h, t, "equal-map-v63")
			v63v2 = make(map[int]bool, len(v))                                         // reset map
			testUnmarshalErr(reflect.ValueOf(v63v2), bs63, h, t, "dec-map-v63-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v63v1, v63v2, h, t, "equal-map-v63-noaddr")
		}
		if v == nil {
			v63v2 = nil
		} else {
			v63v2 = make(map[int]bool, len(v))
		} // reset map
		testUnmarshalErr(&v63v2, bs63, h, t, "dec-map-v63-p-len")
		testDeepEqualErrHandle(v63v1, v63v2, h, t, "equal-map-v63-p-len")
		testReleaseBytes(bs63)
		bs63 = testMarshalErr(&v63v1, h, t, "enc-map-v63-p")
		v63v2 = nil
		testUnmarshalErr(&v63v2, bs63, h, t, "dec-map-v63-p-nil")
		testDeepEqualErrHandle(v63v1, v63v2, h, t, "equal-map-v63-p-nil")
		testReleaseBytes(bs63)
		// ...
		if v == nil {
			v63v2 = nil
		} else {
			v63v2 = make(map[int]bool, len(v))
		} // reset map
		var v63v3, v63v4 typMapMapIntBool
		v63v3 = typMapMapIntBool(v63v1)
		v63v4 = typMapMapIntBool(v63v2)
		if v != nil {
			bs63 = testMarshalErr(v63v3, h, t, "enc-map-v63-custom")
			testUnmarshalErr(v63v4, bs63, h, t, "dec-map-v63-p-len")
			testDeepEqualErrHandle(v63v3, v63v4, h, t, "equal-map-v63-p-len")
			testReleaseBytes(bs63)
		}
		type s63T struct {
			M  map[int]bool
			Mp *map[int]bool
		}
		var m63v99 = map[int]bool{
			0:   false,
			111: false,
		}
		var s63v1, s63v2 s63T
		bs63 = testMarshalErr(s63v1, h, t, "enc-map-v63-custom")
		testUnmarshalErr(&s63v2, bs63, h, t, "dec-map-v63-p-len")
		testDeepEqualErrHandle(s63v1, s63v2, h, t, "equal-map-v63-p-len")
		testReleaseBytes(bs63)
		s63v2 = s63T{}
		s63v1.M = m63v99
		bs63 = testMarshalErr(s63v1, h, t, "enc-map-v63-custom")
		testUnmarshalErr(&s63v2, bs63, h, t, "dec-map-v63-p-len")
		testDeepEqualErrHandle(s63v1, s63v2, h, t, "equal-map-v63-p-len")
		testReleaseBytes(bs63)
		s63v2 = s63T{}
		s63v1.Mp = &m63v99
		bs63 = testMarshalErr(s63v1, h, t, "enc-map-v63-custom")
		testUnmarshalErr(&s63v2, bs63, h, t, "dec-map-v63-p-len")
		testDeepEqualErrHandle(s63v1, s63v2, h, t, "equal-map-v63-p-len")
		testReleaseBytes(bs63)
	}
	for _, v := range []map[int32]interface{}{nil, {}, {77: nil, 127: "string-is-an-interface-3"}} {
		var v64v1, v64v2 map[int32]interface{}
		var bs64 []byte
		v64v1 = v
		bs64 = testMarshalErr(v64v1, h, t, "enc-map-v64")
		if v != nil {
			v64v2 = make(map[int32]interface{}, len(v)) // reset map
			testUnmarshalErr(v64v2, bs64, h, t, "dec-map-v64")
			testDeepEqualErrHandle(v64v1, v64v2, h, t, "equal-map-v64")
			v64v2 = make(map[int32]interface{}, len(v))                                // reset map
			testUnmarshalErr(reflect.ValueOf(v64v2), bs64, h, t, "dec-map-v64-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v64v1, v64v2, h, t, "equal-map-v64-noaddr")
		}
		if v == nil {
			v64v2 = nil
		} else {
			v64v2 = make(map[int32]interface{}, len(v))
		} // reset map
		testUnmarshalErr(&v64v2, bs64, h, t, "dec-map-v64-p-len")
		testDeepEqualErrHandle(v64v1, v64v2, h, t, "equal-map-v64-p-len")
		testReleaseBytes(bs64)
		bs64 = testMarshalErr(&v64v1, h, t, "enc-map-v64-p")
		v64v2 = nil
		testUnmarshalErr(&v64v2, bs64, h, t, "dec-map-v64-p-nil")
		testDeepEqualErrHandle(v64v1, v64v2, h, t, "equal-map-v64-p-nil")
		testReleaseBytes(bs64)
		// ...
		if v == nil {
			v64v2 = nil
		} else {
			v64v2 = make(map[int32]interface{}, len(v))
		} // reset map
		var v64v3, v64v4 typMapMapInt32Intf
		v64v3 = typMapMapInt32Intf(v64v1)
		v64v4 = typMapMapInt32Intf(v64v2)
		if v != nil {
			bs64 = testMarshalErr(v64v3, h, t, "enc-map-v64-custom")
			testUnmarshalErr(v64v4, bs64, h, t, "dec-map-v64-p-len")
			testDeepEqualErrHandle(v64v3, v64v4, h, t, "equal-map-v64-p-len")
			testReleaseBytes(bs64)
		}
		type s64T struct {
			M  map[int32]interface{}
			Mp *map[int32]interface{}
		}
		var m64v99 = map[int32]interface{}{
			0:   nil,
			111: "string-is-an-interface-1",
		}
		var s64v1, s64v2 s64T
		bs64 = testMarshalErr(s64v1, h, t, "enc-map-v64-custom")
		testUnmarshalErr(&s64v2, bs64, h, t, "dec-map-v64-p-len")
		testDeepEqualErrHandle(s64v1, s64v2, h, t, "equal-map-v64-p-len")
		testReleaseBytes(bs64)
		s64v2 = s64T{}
		s64v1.M = m64v99
		bs64 = testMarshalErr(s64v1, h, t, "enc-map-v64-custom")
		testUnmarshalErr(&s64v2, bs64, h, t, "dec-map-v64-p-len")
		testDeepEqualErrHandle(s64v1, s64v2, h, t, "equal-map-v64-p-len")
		testReleaseBytes(bs64)
		s64v2 = s64T{}
		s64v1.Mp = &m64v99
		bs64 = testMarshalErr(s64v1, h, t, "enc-map-v64-custom")
		testUnmarshalErr(&s64v2, bs64, h, t, "dec-map-v64-p-len")
		testDeepEqualErrHandle(s64v1, s64v2, h, t, "equal-map-v64-p-len")
		testReleaseBytes(bs64)
	}
	for _, v := range []map[int32]string{nil, {}, {77: "", 127: "some-string-3"}} {
		var v65v1, v65v2 map[int32]string
		var bs65 []byte
		v65v1 = v
		bs65 = testMarshalErr(v65v1, h, t, "enc-map-v65")
		if v != nil {
			v65v2 = make(map[int32]string, len(v)) // reset map
			testUnmarshalErr(v65v2, bs65, h, t, "dec-map-v65")
			testDeepEqualErrHandle(v65v1, v65v2, h, t, "equal-map-v65")
			v65v2 = make(map[int32]string, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v65v2), bs65, h, t, "dec-map-v65-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v65v1, v65v2, h, t, "equal-map-v65-noaddr")
		}
		if v == nil {
			v65v2 = nil
		} else {
			v65v2 = make(map[int32]string, len(v))
		} // reset map
		testUnmarshalErr(&v65v2, bs65, h, t, "dec-map-v65-p-len")
		testDeepEqualErrHandle(v65v1, v65v2, h, t, "equal-map-v65-p-len")
		testReleaseBytes(bs65)
		bs65 = testMarshalErr(&v65v1, h, t, "enc-map-v65-p")
		v65v2 = nil
		testUnmarshalErr(&v65v2, bs65, h, t, "dec-map-v65-p-nil")
		testDeepEqualErrHandle(v65v1, v65v2, h, t, "equal-map-v65-p-nil")
		testReleaseBytes(bs65)
		// ...
		if v == nil {
			v65v2 = nil
		} else {
			v65v2 = make(map[int32]string, len(v))
		} // reset map
		var v65v3, v65v4 typMapMapInt32String
		v65v3 = typMapMapInt32String(v65v1)
		v65v4 = typMapMapInt32String(v65v2)
		if v != nil {
			bs65 = testMarshalErr(v65v3, h, t, "enc-map-v65-custom")
			testUnmarshalErr(v65v4, bs65, h, t, "dec-map-v65-p-len")
			testDeepEqualErrHandle(v65v3, v65v4, h, t, "equal-map-v65-p-len")
			testReleaseBytes(bs65)
		}
		type s65T struct {
			M  map[int32]string
			Mp *map[int32]string
		}
		var m65v99 = map[int32]string{
			0:   "",
			111: "some-string-1",
		}
		var s65v1, s65v2 s65T
		bs65 = testMarshalErr(s65v1, h, t, "enc-map-v65-custom")
		testUnmarshalErr(&s65v2, bs65, h, t, "dec-map-v65-p-len")
		testDeepEqualErrHandle(s65v1, s65v2, h, t, "equal-map-v65-p-len")
		testReleaseBytes(bs65)
		s65v2 = s65T{}
		s65v1.M = m65v99
		bs65 = testMarshalErr(s65v1, h, t, "enc-map-v65-custom")
		testUnmarshalErr(&s65v2, bs65, h, t, "dec-map-v65-p-len")
		testDeepEqualErrHandle(s65v1, s65v2, h, t, "equal-map-v65-p-len")
		testReleaseBytes(bs65)
		s65v2 = s65T{}
		s65v1.Mp = &m65v99
		bs65 = testMarshalErr(s65v1, h, t, "enc-map-v65-custom")
		testUnmarshalErr(&s65v2, bs65, h, t, "dec-map-v65-p-len")
		testDeepEqualErrHandle(s65v1, s65v2, h, t, "equal-map-v65-p-len")
		testReleaseBytes(bs65)
	}
	for _, v := range []map[int32][]byte{nil, {}, {77: nil, 127: []byte("some-string-3")}} {
		var v66v1, v66v2 map[int32][]byte
		var bs66 []byte
		v66v1 = v
		bs66 = testMarshalErr(v66v1, h, t, "enc-map-v66")
		if v != nil {
			v66v2 = make(map[int32][]byte, len(v)) // reset map
			testUnmarshalErr(v66v2, bs66, h, t, "dec-map-v66")
			testDeepEqualErrHandle(v66v1, v66v2, h, t, "equal-map-v66")
			v66v2 = make(map[int32][]byte, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v66v2), bs66, h, t, "dec-map-v66-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v66v1, v66v2, h, t, "equal-map-v66-noaddr")
		}
		if v == nil {
			v66v2 = nil
		} else {
			v66v2 = make(map[int32][]byte, len(v))
		} // reset map
		testUnmarshalErr(&v66v2, bs66, h, t, "dec-map-v66-p-len")
		testDeepEqualErrHandle(v66v1, v66v2, h, t, "equal-map-v66-p-len")
		testReleaseBytes(bs66)
		bs66 = testMarshalErr(&v66v1, h, t, "enc-map-v66-p")
		v66v2 = nil
		testUnmarshalErr(&v66v2, bs66, h, t, "dec-map-v66-p-nil")
		testDeepEqualErrHandle(v66v1, v66v2, h, t, "equal-map-v66-p-nil")
		testReleaseBytes(bs66)
		// ...
		if v == nil {
			v66v2 = nil
		} else {
			v66v2 = make(map[int32][]byte, len(v))
		} // reset map
		var v66v3, v66v4 typMapMapInt32Bytes
		v66v3 = typMapMapInt32Bytes(v66v1)
		v66v4 = typMapMapInt32Bytes(v66v2)
		if v != nil {
			bs66 = testMarshalErr(v66v3, h, t, "enc-map-v66-custom")
			testUnmarshalErr(v66v4, bs66, h, t, "dec-map-v66-p-len")
			testDeepEqualErrHandle(v66v3, v66v4, h, t, "equal-map-v66-p-len")
			testReleaseBytes(bs66)
		}
		type s66T struct {
			M  map[int32][]byte
			Mp *map[int32][]byte
		}
		var m66v99 = map[int32][]byte{
			0:   nil,
			111: []byte("some-string-1"),
		}
		var s66v1, s66v2 s66T
		bs66 = testMarshalErr(s66v1, h, t, "enc-map-v66-custom")
		testUnmarshalErr(&s66v2, bs66, h, t, "dec-map-v66-p-len")
		testDeepEqualErrHandle(s66v1, s66v2, h, t, "equal-map-v66-p-len")
		testReleaseBytes(bs66)
		s66v2 = s66T{}
		s66v1.M = m66v99
		bs66 = testMarshalErr(s66v1, h, t, "enc-map-v66-custom")
		testUnmarshalErr(&s66v2, bs66, h, t, "dec-map-v66-p-len")
		testDeepEqualErrHandle(s66v1, s66v2, h, t, "equal-map-v66-p-len")
		testReleaseBytes(bs66)
		s66v2 = s66T{}
		s66v1.Mp = &m66v99
		bs66 = testMarshalErr(s66v1, h, t, "enc-map-v66-custom")
		testUnmarshalErr(&s66v2, bs66, h, t, "dec-map-v66-p-len")
		testDeepEqualErrHandle(s66v1, s66v2, h, t, "equal-map-v66-p-len")
		testReleaseBytes(bs66)
	}
	for _, v := range []map[int32]uint8{nil, {}, {77: 0, 127: 111}} {
		var v67v1, v67v2 map[int32]uint8
		var bs67 []byte
		v67v1 = v
		bs67 = testMarshalErr(v67v1, h, t, "enc-map-v67")
		if v != nil {
			v67v2 = make(map[int32]uint8, len(v)) // reset map
			testUnmarshalErr(v67v2, bs67, h, t, "dec-map-v67")
			testDeepEqualErrHandle(v67v1, v67v2, h, t, "equal-map-v67")
			v67v2 = make(map[int32]uint8, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v67v2), bs67, h, t, "dec-map-v67-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v67v1, v67v2, h, t, "equal-map-v67-noaddr")
		}
		if v == nil {
			v67v2 = nil
		} else {
			v67v2 = make(map[int32]uint8, len(v))
		} // reset map
		testUnmarshalErr(&v67v2, bs67, h, t, "dec-map-v67-p-len")
		testDeepEqualErrHandle(v67v1, v67v2, h, t, "equal-map-v67-p-len")
		testReleaseBytes(bs67)
		bs67 = testMarshalErr(&v67v1, h, t, "enc-map-v67-p")
		v67v2 = nil
		testUnmarshalErr(&v67v2, bs67, h, t, "dec-map-v67-p-nil")
		testDeepEqualErrHandle(v67v1, v67v2, h, t, "equal-map-v67-p-nil")
		testReleaseBytes(bs67)
		// ...
		if v == nil {
			v67v2 = nil
		} else {
			v67v2 = make(map[int32]uint8, len(v))
		} // reset map
		var v67v3, v67v4 typMapMapInt32Uint8
		v67v3 = typMapMapInt32Uint8(v67v1)
		v67v4 = typMapMapInt32Uint8(v67v2)
		if v != nil {
			bs67 = testMarshalErr(v67v3, h, t, "enc-map-v67-custom")
			testUnmarshalErr(v67v4, bs67, h, t, "dec-map-v67-p-len")
			testDeepEqualErrHandle(v67v3, v67v4, h, t, "equal-map-v67-p-len")
			testReleaseBytes(bs67)
		}
		type s67T struct {
			M  map[int32]uint8
			Mp *map[int32]uint8
		}
		var m67v99 = map[int32]uint8{
			0:  0,
			77: 127,
		}
		var s67v1, s67v2 s67T
		bs67 = testMarshalErr(s67v1, h, t, "enc-map-v67-custom")
		testUnmarshalErr(&s67v2, bs67, h, t, "dec-map-v67-p-len")
		testDeepEqualErrHandle(s67v1, s67v2, h, t, "equal-map-v67-p-len")
		testReleaseBytes(bs67)
		s67v2 = s67T{}
		s67v1.M = m67v99
		bs67 = testMarshalErr(s67v1, h, t, "enc-map-v67-custom")
		testUnmarshalErr(&s67v2, bs67, h, t, "dec-map-v67-p-len")
		testDeepEqualErrHandle(s67v1, s67v2, h, t, "equal-map-v67-p-len")
		testReleaseBytes(bs67)
		s67v2 = s67T{}
		s67v1.Mp = &m67v99
		bs67 = testMarshalErr(s67v1, h, t, "enc-map-v67-custom")
		testUnmarshalErr(&s67v2, bs67, h, t, "dec-map-v67-p-len")
		testDeepEqualErrHandle(s67v1, s67v2, h, t, "equal-map-v67-p-len")
		testReleaseBytes(bs67)
	}
	for _, v := range []map[int32]uint64{nil, {}, {111: 0, 77: 127}} {
		var v68v1, v68v2 map[int32]uint64
		var bs68 []byte
		v68v1 = v
		bs68 = testMarshalErr(v68v1, h, t, "enc-map-v68")
		if v != nil {
			v68v2 = make(map[int32]uint64, len(v)) // reset map
			testUnmarshalErr(v68v2, bs68, h, t, "dec-map-v68")
			testDeepEqualErrHandle(v68v1, v68v2, h, t, "equal-map-v68")
			v68v2 = make(map[int32]uint64, len(v))                                     // reset map
			testUnmarshalErr(reflect.ValueOf(v68v2), bs68, h, t, "dec-map-v68-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v68v1, v68v2, h, t, "equal-map-v68-noaddr")
		}
		if v == nil {
			v68v2 = nil
		} else {
			v68v2 = make(map[int32]uint64, len(v))
		} // reset map
		testUnmarshalErr(&v68v2, bs68, h, t, "dec-map-v68-p-len")
		testDeepEqualErrHandle(v68v1, v68v2, h, t, "equal-map-v68-p-len")
		testReleaseBytes(bs68)
		bs68 = testMarshalErr(&v68v1, h, t, "enc-map-v68-p")
		v68v2 = nil
		testUnmarshalErr(&v68v2, bs68, h, t, "dec-map-v68-p-nil")
		testDeepEqualErrHandle(v68v1, v68v2, h, t, "equal-map-v68-p-nil")
		testReleaseBytes(bs68)
		// ...
		if v == nil {
			v68v2 = nil
		} else {
			v68v2 = make(map[int32]uint64, len(v))
		} // reset map
		var v68v3, v68v4 typMapMapInt32Uint64
		v68v3 = typMapMapInt32Uint64(v68v1)
		v68v4 = typMapMapInt32Uint64(v68v2)
		if v != nil {
			bs68 = testMarshalErr(v68v3, h, t, "enc-map-v68-custom")
			testUnmarshalErr(v68v4, bs68, h, t, "dec-map-v68-p-len")
			testDeepEqualErrHandle(v68v3, v68v4, h, t, "equal-map-v68-p-len")
			testReleaseBytes(bs68)
		}
		type s68T struct {
			M  map[int32]uint64
			Mp *map[int32]uint64
		}
		var m68v99 = map[int32]uint64{
			0:   0,
			111: 77,
		}
		var s68v1, s68v2 s68T
		bs68 = testMarshalErr(s68v1, h, t, "enc-map-v68-custom")
		testUnmarshalErr(&s68v2, bs68, h, t, "dec-map-v68-p-len")
		testDeepEqualErrHandle(s68v1, s68v2, h, t, "equal-map-v68-p-len")
		testReleaseBytes(bs68)
		s68v2 = s68T{}
		s68v1.M = m68v99
		bs68 = testMarshalErr(s68v1, h, t, "enc-map-v68-custom")
		testUnmarshalErr(&s68v2, bs68, h, t, "dec-map-v68-p-len")
		testDeepEqualErrHandle(s68v1, s68v2, h, t, "equal-map-v68-p-len")
		testReleaseBytes(bs68)
		s68v2 = s68T{}
		s68v1.Mp = &m68v99
		bs68 = testMarshalErr(s68v1, h, t, "enc-map-v68-custom")
		testUnmarshalErr(&s68v2, bs68, h, t, "dec-map-v68-p-len")
		testDeepEqualErrHandle(s68v1, s68v2, h, t, "equal-map-v68-p-len")
		testReleaseBytes(bs68)
	}
	for _, v := range []map[int32]int{nil, {}, {127: 0, 111: 77}} {
		var v69v1, v69v2 map[int32]int
		var bs69 []byte
		v69v1 = v
		bs69 = testMarshalErr(v69v1, h, t, "enc-map-v69")
		if v != nil {
			v69v2 = make(map[int32]int, len(v)) // reset map
			testUnmarshalErr(v69v2, bs69, h, t, "dec-map-v69")
			testDeepEqualErrHandle(v69v1, v69v2, h, t, "equal-map-v69")
			v69v2 = make(map[int32]int, len(v))                                        // reset map
			testUnmarshalErr(reflect.ValueOf(v69v2), bs69, h, t, "dec-map-v69-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v69v1, v69v2, h, t, "equal-map-v69-noaddr")
		}
		if v == nil {
			v69v2 = nil
		} else {
			v69v2 = make(map[int32]int, len(v))
		} // reset map
		testUnmarshalErr(&v69v2, bs69, h, t, "dec-map-v69-p-len")
		testDeepEqualErrHandle(v69v1, v69v2, h, t, "equal-map-v69-p-len")
		testReleaseBytes(bs69)
		bs69 = testMarshalErr(&v69v1, h, t, "enc-map-v69-p")
		v69v2 = nil
		testUnmarshalErr(&v69v2, bs69, h, t, "dec-map-v69-p-nil")
		testDeepEqualErrHandle(v69v1, v69v2, h, t, "equal-map-v69-p-nil")
		testReleaseBytes(bs69)
		// ...
		if v == nil {
			v69v2 = nil
		} else {
			v69v2 = make(map[int32]int, len(v))
		} // reset map
		var v69v3, v69v4 typMapMapInt32Int
		v69v3 = typMapMapInt32Int(v69v1)
		v69v4 = typMapMapInt32Int(v69v2)
		if v != nil {
			bs69 = testMarshalErr(v69v3, h, t, "enc-map-v69-custom")
			testUnmarshalErr(v69v4, bs69, h, t, "dec-map-v69-p-len")
			testDeepEqualErrHandle(v69v3, v69v4, h, t, "equal-map-v69-p-len")
			testReleaseBytes(bs69)
		}
		type s69T struct {
			M  map[int32]int
			Mp *map[int32]int
		}
		var m69v99 = map[int32]int{
			0:   0,
			127: 111,
		}
		var s69v1, s69v2 s69T
		bs69 = testMarshalErr(s69v1, h, t, "enc-map-v69-custom")
		testUnmarshalErr(&s69v2, bs69, h, t, "dec-map-v69-p-len")
		testDeepEqualErrHandle(s69v1, s69v2, h, t, "equal-map-v69-p-len")
		testReleaseBytes(bs69)
		s69v2 = s69T{}
		s69v1.M = m69v99
		bs69 = testMarshalErr(s69v1, h, t, "enc-map-v69-custom")
		testUnmarshalErr(&s69v2, bs69, h, t, "dec-map-v69-p-len")
		testDeepEqualErrHandle(s69v1, s69v2, h, t, "equal-map-v69-p-len")
		testReleaseBytes(bs69)
		s69v2 = s69T{}
		s69v1.Mp = &m69v99
		bs69 = testMarshalErr(s69v1, h, t, "enc-map-v69-custom")
		testUnmarshalErr(&s69v2, bs69, h, t, "dec-map-v69-p-len")
		testDeepEqualErrHandle(s69v1, s69v2, h, t, "equal-map-v69-p-len")
		testReleaseBytes(bs69)
	}
	for _, v := range []map[int32]int32{nil, {}, {77: 0, 127: 111}} {
		var v70v1, v70v2 map[int32]int32
		var bs70 []byte
		v70v1 = v
		bs70 = testMarshalErr(v70v1, h, t, "enc-map-v70")
		if v != nil {
			v70v2 = make(map[int32]int32, len(v)) // reset map
			testUnmarshalErr(v70v2, bs70, h, t, "dec-map-v70")
			testDeepEqualErrHandle(v70v1, v70v2, h, t, "equal-map-v70")
			v70v2 = make(map[int32]int32, len(v))                                      // reset map
			testUnmarshalErr(reflect.ValueOf(v70v2), bs70, h, t, "dec-map-v70-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v70v1, v70v2, h, t, "equal-map-v70-noaddr")
		}
		if v == nil {
			v70v2 = nil
		} else {
			v70v2 = make(map[int32]int32, len(v))
		} // reset map
		testUnmarshalErr(&v70v2, bs70, h, t, "dec-map-v70-p-len")
		testDeepEqualErrHandle(v70v1, v70v2, h, t, "equal-map-v70-p-len")
		testReleaseBytes(bs70)
		bs70 = testMarshalErr(&v70v1, h, t, "enc-map-v70-p")
		v70v2 = nil
		testUnmarshalErr(&v70v2, bs70, h, t, "dec-map-v70-p-nil")
		testDeepEqualErrHandle(v70v1, v70v2, h, t, "equal-map-v70-p-nil")
		testReleaseBytes(bs70)
		// ...
		if v == nil {
			v70v2 = nil
		} else {
			v70v2 = make(map[int32]int32, len(v))
		} // reset map
		var v70v3, v70v4 typMapMapInt32Int32
		v70v3 = typMapMapInt32Int32(v70v1)
		v70v4 = typMapMapInt32Int32(v70v2)
		if v != nil {
			bs70 = testMarshalErr(v70v3, h, t, "enc-map-v70-custom")
			testUnmarshalErr(v70v4, bs70, h, t, "dec-map-v70-p-len")
			testDeepEqualErrHandle(v70v3, v70v4, h, t, "equal-map-v70-p-len")
			testReleaseBytes(bs70)
		}
		type s70T struct {
			M  map[int32]int32
			Mp *map[int32]int32
		}
		var m70v99 = map[int32]int32{
			0:  0,
			77: 127,
		}
		var s70v1, s70v2 s70T
		bs70 = testMarshalErr(s70v1, h, t, "enc-map-v70-custom")
		testUnmarshalErr(&s70v2, bs70, h, t, "dec-map-v70-p-len")
		testDeepEqualErrHandle(s70v1, s70v2, h, t, "equal-map-v70-p-len")
		testReleaseBytes(bs70)
		s70v2 = s70T{}
		s70v1.M = m70v99
		bs70 = testMarshalErr(s70v1, h, t, "enc-map-v70-custom")
		testUnmarshalErr(&s70v2, bs70, h, t, "dec-map-v70-p-len")
		testDeepEqualErrHandle(s70v1, s70v2, h, t, "equal-map-v70-p-len")
		testReleaseBytes(bs70)
		s70v2 = s70T{}
		s70v1.Mp = &m70v99
		bs70 = testMarshalErr(s70v1, h, t, "enc-map-v70-custom")
		testUnmarshalErr(&s70v2, bs70, h, t, "dec-map-v70-p-len")
		testDeepEqualErrHandle(s70v1, s70v2, h, t, "equal-map-v70-p-len")
		testReleaseBytes(bs70)
	}
	for _, v := range []map[int32]float64{nil, {}, {111: 0, 77: 22.2}} {
		var v71v1, v71v2 map[int32]float64
		var bs71 []byte
		v71v1 = v
		bs71 = testMarshalErr(v71v1, h, t, "enc-map-v71")
		if v != nil {
			v71v2 = make(map[int32]float64, len(v)) // reset map
			testUnmarshalErr(v71v2, bs71, h, t, "dec-map-v71")
			testDeepEqualErrHandle(v71v1, v71v2, h, t, "equal-map-v71")
			v71v2 = make(map[int32]float64, len(v))                                    // reset map
			testUnmarshalErr(reflect.ValueOf(v71v2), bs71, h, t, "dec-map-v71-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v71v1, v71v2, h, t, "equal-map-v71-noaddr")
		}
		if v == nil {
			v71v2 = nil
		} else {
			v71v2 = make(map[int32]float64, len(v))
		} // reset map
		testUnmarshalErr(&v71v2, bs71, h, t, "dec-map-v71-p-len")
		testDeepEqualErrHandle(v71v1, v71v2, h, t, "equal-map-v71-p-len")
		testReleaseBytes(bs71)
		bs71 = testMarshalErr(&v71v1, h, t, "enc-map-v71-p")
		v71v2 = nil
		testUnmarshalErr(&v71v2, bs71, h, t, "dec-map-v71-p-nil")
		testDeepEqualErrHandle(v71v1, v71v2, h, t, "equal-map-v71-p-nil")
		testReleaseBytes(bs71)
		// ...
		if v == nil {
			v71v2 = nil
		} else {
			v71v2 = make(map[int32]float64, len(v))
		} // reset map
		var v71v3, v71v4 typMapMapInt32Float64
		v71v3 = typMapMapInt32Float64(v71v1)
		v71v4 = typMapMapInt32Float64(v71v2)
		if v != nil {
			bs71 = testMarshalErr(v71v3, h, t, "enc-map-v71-custom")
			testUnmarshalErr(v71v4, bs71, h, t, "dec-map-v71-p-len")
			testDeepEqualErrHandle(v71v3, v71v4, h, t, "equal-map-v71-p-len")
			testReleaseBytes(bs71)
		}
		type s71T struct {
			M  map[int32]float64
			Mp *map[int32]float64
		}
		var m71v99 = map[int32]float64{
			0:   0,
			127: 33.3e3,
		}
		var s71v1, s71v2 s71T
		bs71 = testMarshalErr(s71v1, h, t, "enc-map-v71-custom")
		testUnmarshalErr(&s71v2, bs71, h, t, "dec-map-v71-p-len")
		testDeepEqualErrHandle(s71v1, s71v2, h, t, "equal-map-v71-p-len")
		testReleaseBytes(bs71)
		s71v2 = s71T{}
		s71v1.M = m71v99
		bs71 = testMarshalErr(s71v1, h, t, "enc-map-v71-custom")
		testUnmarshalErr(&s71v2, bs71, h, t, "dec-map-v71-p-len")
		testDeepEqualErrHandle(s71v1, s71v2, h, t, "equal-map-v71-p-len")
		testReleaseBytes(bs71)
		s71v2 = s71T{}
		s71v1.Mp = &m71v99
		bs71 = testMarshalErr(s71v1, h, t, "enc-map-v71-custom")
		testUnmarshalErr(&s71v2, bs71, h, t, "dec-map-v71-p-len")
		testDeepEqualErrHandle(s71v1, s71v2, h, t, "equal-map-v71-p-len")
		testReleaseBytes(bs71)
	}
	for _, v := range []map[int32]bool{nil, {}, {111: false, 77: true}} {
		var v72v1, v72v2 map[int32]bool
		var bs72 []byte
		v72v1 = v
		bs72 = testMarshalErr(v72v1, h, t, "enc-map-v72")
		if v != nil {
			v72v2 = make(map[int32]bool, len(v)) // reset map
			testUnmarshalErr(v72v2, bs72, h, t, "dec-map-v72")
			testDeepEqualErrHandle(v72v1, v72v2, h, t, "equal-map-v72")
			v72v2 = make(map[int32]bool, len(v))                                       // reset map
			testUnmarshalErr(reflect.ValueOf(v72v2), bs72, h, t, "dec-map-v72-noaddr") // decode into non-addressable map value
			testDeepEqualErrHandle(v72v1, v72v2, h, t, "equal-map-v72-noaddr")
		}
		if v == nil {
			v72v2 = nil
		} else {
			v72v2 = make(map[int32]bool, len(v))
		} // reset map
		testUnmarshalErr(&v72v2, bs72, h, t, "dec-map-v72-p-len")
		testDeepEqualErrHandle(v72v1, v72v2, h, t, "equal-map-v72-p-len")
		testReleaseBytes(bs72)
		bs72 = testMarshalErr(&v72v1, h, t, "enc-map-v72-p")
		v72v2 = nil
		testUnmarshalErr(&v72v2, bs72, h, t, "dec-map-v72-p-nil")
		testDeepEqualErrHandle(v72v1, v72v2, h, t, "equal-map-v72-p-nil")
		testReleaseBytes(bs72)
		// ...
		if v == nil {
			v72v2 = nil
		} else {
			v72v2 = make(map[int32]bool, len(v))
		} // reset map
		var v72v3, v72v4 typMapMapInt32Bool
		v72v3 = typMapMapInt32Bool(v72v1)
		v72v4 = typMapMapInt32Bool(v72v2)
		if v != nil {
			bs72 = testMarshalErr(v72v3, h, t, "enc-map-v72-custom")
			testUnmarshalErr(v72v4, bs72, h, t, "dec-map-v72-p-len")
			testDeepEqualErrHandle(v72v3, v72v4, h, t, "equal-map-v72-p-len")
			testReleaseBytes(bs72)
		}
		type s72T struct {
			M  map[int32]bool
			Mp *map[int32]bool
		}
		var m72v99 = map[int32]bool{
			0:   false,
			127: true,
		}
		var s72v1, s72v2 s72T
		bs72 = testMarshalErr(s72v1, h, t, "enc-map-v72-custom")
		testUnmarshalErr(&s72v2, bs72, h, t, "dec-map-v72-p-len")
		testDeepEqualErrHandle(s72v1, s72v2, h, t, "equal-map-v72-p-len")
		testReleaseBytes(bs72)
		s72v2 = s72T{}
		s72v1.M = m72v99
		bs72 = testMarshalErr(s72v1, h, t, "enc-map-v72-custom")
		testUnmarshalErr(&s72v2, bs72, h, t, "dec-map-v72-p-len")
		testDeepEqualErrHandle(s72v1, s72v2, h, t, "equal-map-v72-p-len")
		testReleaseBytes(bs72)
		s72v2 = s72T{}
		s72v1.Mp = &m72v99
		bs72 = testMarshalErr(s72v1, h, t, "enc-map-v72-custom")
		testUnmarshalErr(&s72v2, bs72, h, t, "dec-map-v72-p-len")
		testDeepEqualErrHandle(s72v1, s72v2, h, t, "equal-map-v72-p-len")
		testReleaseBytes(bs72)
	}

}

func doTestMammothMapsAndSlices(t *testing.T, h Handle) {
	defer testSetup(t, &h)()
	if mh, ok := h.(*MsgpackHandle); ok {
		defer func(b bool) { mh.RawToString = b }(mh.RawToString)
		mh.RawToString = true
	}
	__doTestMammothSlices(t, h)
	__doTestMammothMaps(t, h)
}

func doTestMammoth(t *testing.T, h Handle) {
	defer testSetup(t, &h)()
	if mh, ok := h.(*MsgpackHandle); ok {
		defer func(b bool) { mh.RawToString = b }(mh.RawToString)
		mh.RawToString = true
	}

	name := h.Name()
	var b []byte

	var m, m2 TestMammoth
	testRandomFillRV(reflect.ValueOf(&m).Elem())
	b = testMarshalErr(&m, h, t, "mammoth-"+name)

	testUnmarshalErr(&m2, b, h, t, "mammoth-"+name)
	testDeepEqualErrHandle(&m, &m2, h, t, "mammoth-"+name)
	testReleaseBytes(b)

	if testing.Short() {
		t.Skipf("skipping rest of mammoth test in -short mode")
	}

	var mm, mm2 TestMammoth2Wrapper
	testRandomFillRV(reflect.ValueOf(&mm).Elem())
	b = testMarshalErr(&mm, h, t, "mammoth2-"+name)
	// os.Stderr.Write([]byte("\n\n\n\n" + string(b) + "\n\n\n\n"))
	testUnmarshalErr(&mm2, b, h, t, "mammoth2-"+name)
	testDeepEqualErrHandle(&mm, &mm2, h, t, "mammoth2-"+name)
	// testMammoth2(t, name, h)
	testReleaseBytes(b)
}

func TestJsonMammoth(t *testing.T) {
	doTestMammoth(t, testJsonH)
}
func TestCborMammoth(t *testing.T) {
	doTestMammoth(t, testCborH)
}
func TestMsgpackMammoth(t *testing.T) {
	doTestMammoth(t, testMsgpackH)
}
func TestBincMammoth(t *testing.T) {
	doTestMammoth(t, testBincH)
}
func TestSimpleMammoth(t *testing.T) {
	doTestMammoth(t, testSimpleH)
}

func TestJsonMammothMapsAndSlices(t *testing.T) {
	doTestMammothMapsAndSlices(t, testJsonH)
}
func TestCborMammothMapsAndSlices(t *testing.T) {
	doTestMammothMapsAndSlices(t, testCborH)
}
func TestMsgpackMammothMapsAndSlices(t *testing.T) {
	doTestMammothMapsAndSlices(t, testMsgpackH)
}
func TestBincMammothMapsAndSlices(t *testing.T) {
	doTestMammothMapsAndSlices(t, testBincH)
}
func TestSimpleMammothMapsAndSlices(t *testing.T) {
	doTestMammothMapsAndSlices(t, testSimpleH)
}
