// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED from mammoth-test.go.tmpl
// ************************************************************

package codec

import "testing"
import "fmt"

// TestMammoth has all the different paths optimized in fast-path
// It has all the primitives, slices and maps.
//
// For each of those types, it has a pointer and a non-pointer field.

func init() { _ = fmt.Printf } // so we can include fmt as needed

type TestMammoth struct {
	FIntf       interface{}
	FptrIntf    *interface{}
	FString     string
	FptrString  *string
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

	FSliceIntf       []interface{}
	FptrSliceIntf    *[]interface{}
	FSliceString     []string
	FptrSliceString  *[]string
	FSliceFloat32    []float32
	FptrSliceFloat32 *[]float32
	FSliceFloat64    []float64
	FptrSliceFloat64 *[]float64
	FSliceUint       []uint
	FptrSliceUint    *[]uint
	FSliceUint16     []uint16
	FptrSliceUint16  *[]uint16
	FSliceUint32     []uint32
	FptrSliceUint32  *[]uint32
	FSliceUint64     []uint64
	FptrSliceUint64  *[]uint64
	FSliceUintptr    []uintptr
	FptrSliceUintptr *[]uintptr
	FSliceInt        []int
	FptrSliceInt     *[]int
	FSliceInt8       []int8
	FptrSliceInt8    *[]int8
	FSliceInt16      []int16
	FptrSliceInt16   *[]int16
	FSliceInt32      []int32
	FptrSliceInt32   *[]int32
	FSliceInt64      []int64
	FptrSliceInt64   *[]int64
	FSliceBool       []bool
	FptrSliceBool    *[]bool

	FMapIntfIntf          map[interface{}]interface{}
	FptrMapIntfIntf       *map[interface{}]interface{}
	FMapIntfString        map[interface{}]string
	FptrMapIntfString     *map[interface{}]string
	FMapIntfUint          map[interface{}]uint
	FptrMapIntfUint       *map[interface{}]uint
	FMapIntfUint8         map[interface{}]uint8
	FptrMapIntfUint8      *map[interface{}]uint8
	FMapIntfUint16        map[interface{}]uint16
	FptrMapIntfUint16     *map[interface{}]uint16
	FMapIntfUint32        map[interface{}]uint32
	FptrMapIntfUint32     *map[interface{}]uint32
	FMapIntfUint64        map[interface{}]uint64
	FptrMapIntfUint64     *map[interface{}]uint64
	FMapIntfUintptr       map[interface{}]uintptr
	FptrMapIntfUintptr    *map[interface{}]uintptr
	FMapIntfInt           map[interface{}]int
	FptrMapIntfInt        *map[interface{}]int
	FMapIntfInt8          map[interface{}]int8
	FptrMapIntfInt8       *map[interface{}]int8
	FMapIntfInt16         map[interface{}]int16
	FptrMapIntfInt16      *map[interface{}]int16
	FMapIntfInt32         map[interface{}]int32
	FptrMapIntfInt32      *map[interface{}]int32
	FMapIntfInt64         map[interface{}]int64
	FptrMapIntfInt64      *map[interface{}]int64
	FMapIntfFloat32       map[interface{}]float32
	FptrMapIntfFloat32    *map[interface{}]float32
	FMapIntfFloat64       map[interface{}]float64
	FptrMapIntfFloat64    *map[interface{}]float64
	FMapIntfBool          map[interface{}]bool
	FptrMapIntfBool       *map[interface{}]bool
	FMapStringIntf        map[string]interface{}
	FptrMapStringIntf     *map[string]interface{}
	FMapStringString      map[string]string
	FptrMapStringString   *map[string]string
	FMapStringUint        map[string]uint
	FptrMapStringUint     *map[string]uint
	FMapStringUint8       map[string]uint8
	FptrMapStringUint8    *map[string]uint8
	FMapStringUint16      map[string]uint16
	FptrMapStringUint16   *map[string]uint16
	FMapStringUint32      map[string]uint32
	FptrMapStringUint32   *map[string]uint32
	FMapStringUint64      map[string]uint64
	FptrMapStringUint64   *map[string]uint64
	FMapStringUintptr     map[string]uintptr
	FptrMapStringUintptr  *map[string]uintptr
	FMapStringInt         map[string]int
	FptrMapStringInt      *map[string]int
	FMapStringInt8        map[string]int8
	FptrMapStringInt8     *map[string]int8
	FMapStringInt16       map[string]int16
	FptrMapStringInt16    *map[string]int16
	FMapStringInt32       map[string]int32
	FptrMapStringInt32    *map[string]int32
	FMapStringInt64       map[string]int64
	FptrMapStringInt64    *map[string]int64
	FMapStringFloat32     map[string]float32
	FptrMapStringFloat32  *map[string]float32
	FMapStringFloat64     map[string]float64
	FptrMapStringFloat64  *map[string]float64
	FMapStringBool        map[string]bool
	FptrMapStringBool     *map[string]bool
	FMapFloat32Intf       map[float32]interface{}
	FptrMapFloat32Intf    *map[float32]interface{}
	FMapFloat32String     map[float32]string
	FptrMapFloat32String  *map[float32]string
	FMapFloat32Uint       map[float32]uint
	FptrMapFloat32Uint    *map[float32]uint
	FMapFloat32Uint8      map[float32]uint8
	FptrMapFloat32Uint8   *map[float32]uint8
	FMapFloat32Uint16     map[float32]uint16
	FptrMapFloat32Uint16  *map[float32]uint16
	FMapFloat32Uint32     map[float32]uint32
	FptrMapFloat32Uint32  *map[float32]uint32
	FMapFloat32Uint64     map[float32]uint64
	FptrMapFloat32Uint64  *map[float32]uint64
	FMapFloat32Uintptr    map[float32]uintptr
	FptrMapFloat32Uintptr *map[float32]uintptr
	FMapFloat32Int        map[float32]int
	FptrMapFloat32Int     *map[float32]int
	FMapFloat32Int8       map[float32]int8
	FptrMapFloat32Int8    *map[float32]int8
	FMapFloat32Int16      map[float32]int16
	FptrMapFloat32Int16   *map[float32]int16
	FMapFloat32Int32      map[float32]int32
	FptrMapFloat32Int32   *map[float32]int32
	FMapFloat32Int64      map[float32]int64
	FptrMapFloat32Int64   *map[float32]int64
	FMapFloat32Float32    map[float32]float32
	FptrMapFloat32Float32 *map[float32]float32
	FMapFloat32Float64    map[float32]float64
	FptrMapFloat32Float64 *map[float32]float64
	FMapFloat32Bool       map[float32]bool
	FptrMapFloat32Bool    *map[float32]bool
	FMapFloat64Intf       map[float64]interface{}
	FptrMapFloat64Intf    *map[float64]interface{}
	FMapFloat64String     map[float64]string
	FptrMapFloat64String  *map[float64]string
	FMapFloat64Uint       map[float64]uint
	FptrMapFloat64Uint    *map[float64]uint
	FMapFloat64Uint8      map[float64]uint8
	FptrMapFloat64Uint8   *map[float64]uint8
	FMapFloat64Uint16     map[float64]uint16
	FptrMapFloat64Uint16  *map[float64]uint16
	FMapFloat64Uint32     map[float64]uint32
	FptrMapFloat64Uint32  *map[float64]uint32
	FMapFloat64Uint64     map[float64]uint64
	FptrMapFloat64Uint64  *map[float64]uint64
	FMapFloat64Uintptr    map[float64]uintptr
	FptrMapFloat64Uintptr *map[float64]uintptr
	FMapFloat64Int        map[float64]int
	FptrMapFloat64Int     *map[float64]int
	FMapFloat64Int8       map[float64]int8
	FptrMapFloat64Int8    *map[float64]int8
	FMapFloat64Int16      map[float64]int16
	FptrMapFloat64Int16   *map[float64]int16
	FMapFloat64Int32      map[float64]int32
	FptrMapFloat64Int32   *map[float64]int32
	FMapFloat64Int64      map[float64]int64
	FptrMapFloat64Int64   *map[float64]int64
	FMapFloat64Float32    map[float64]float32
	FptrMapFloat64Float32 *map[float64]float32
	FMapFloat64Float64    map[float64]float64
	FptrMapFloat64Float64 *map[float64]float64
	FMapFloat64Bool       map[float64]bool
	FptrMapFloat64Bool    *map[float64]bool
	FMapUintIntf          map[uint]interface{}
	FptrMapUintIntf       *map[uint]interface{}
	FMapUintString        map[uint]string
	FptrMapUintString     *map[uint]string
	FMapUintUint          map[uint]uint
	FptrMapUintUint       *map[uint]uint
	FMapUintUint8         map[uint]uint8
	FptrMapUintUint8      *map[uint]uint8
	FMapUintUint16        map[uint]uint16
	FptrMapUintUint16     *map[uint]uint16
	FMapUintUint32        map[uint]uint32
	FptrMapUintUint32     *map[uint]uint32
	FMapUintUint64        map[uint]uint64
	FptrMapUintUint64     *map[uint]uint64
	FMapUintUintptr       map[uint]uintptr
	FptrMapUintUintptr    *map[uint]uintptr
	FMapUintInt           map[uint]int
	FptrMapUintInt        *map[uint]int
	FMapUintInt8          map[uint]int8
	FptrMapUintInt8       *map[uint]int8
	FMapUintInt16         map[uint]int16
	FptrMapUintInt16      *map[uint]int16
	FMapUintInt32         map[uint]int32
	FptrMapUintInt32      *map[uint]int32
	FMapUintInt64         map[uint]int64
	FptrMapUintInt64      *map[uint]int64
	FMapUintFloat32       map[uint]float32
	FptrMapUintFloat32    *map[uint]float32
	FMapUintFloat64       map[uint]float64
	FptrMapUintFloat64    *map[uint]float64
	FMapUintBool          map[uint]bool
	FptrMapUintBool       *map[uint]bool
	FMapUint8Intf         map[uint8]interface{}
	FptrMapUint8Intf      *map[uint8]interface{}
	FMapUint8String       map[uint8]string
	FptrMapUint8String    *map[uint8]string
	FMapUint8Uint         map[uint8]uint
	FptrMapUint8Uint      *map[uint8]uint
	FMapUint8Uint8        map[uint8]uint8
	FptrMapUint8Uint8     *map[uint8]uint8
	FMapUint8Uint16       map[uint8]uint16
	FptrMapUint8Uint16    *map[uint8]uint16
	FMapUint8Uint32       map[uint8]uint32
	FptrMapUint8Uint32    *map[uint8]uint32
	FMapUint8Uint64       map[uint8]uint64
	FptrMapUint8Uint64    *map[uint8]uint64
	FMapUint8Uintptr      map[uint8]uintptr
	FptrMapUint8Uintptr   *map[uint8]uintptr
	FMapUint8Int          map[uint8]int
	FptrMapUint8Int       *map[uint8]int
	FMapUint8Int8         map[uint8]int8
	FptrMapUint8Int8      *map[uint8]int8
	FMapUint8Int16        map[uint8]int16
	FptrMapUint8Int16     *map[uint8]int16
	FMapUint8Int32        map[uint8]int32
	FptrMapUint8Int32     *map[uint8]int32
	FMapUint8Int64        map[uint8]int64
	FptrMapUint8Int64     *map[uint8]int64
	FMapUint8Float32      map[uint8]float32
	FptrMapUint8Float32   *map[uint8]float32
	FMapUint8Float64      map[uint8]float64
	FptrMapUint8Float64   *map[uint8]float64
	FMapUint8Bool         map[uint8]bool
	FptrMapUint8Bool      *map[uint8]bool
	FMapUint16Intf        map[uint16]interface{}
	FptrMapUint16Intf     *map[uint16]interface{}
	FMapUint16String      map[uint16]string
	FptrMapUint16String   *map[uint16]string
	FMapUint16Uint        map[uint16]uint
	FptrMapUint16Uint     *map[uint16]uint
	FMapUint16Uint8       map[uint16]uint8
	FptrMapUint16Uint8    *map[uint16]uint8
	FMapUint16Uint16      map[uint16]uint16
	FptrMapUint16Uint16   *map[uint16]uint16
	FMapUint16Uint32      map[uint16]uint32
	FptrMapUint16Uint32   *map[uint16]uint32
	FMapUint16Uint64      map[uint16]uint64
	FptrMapUint16Uint64   *map[uint16]uint64
	FMapUint16Uintptr     map[uint16]uintptr
	FptrMapUint16Uintptr  *map[uint16]uintptr
	FMapUint16Int         map[uint16]int
	FptrMapUint16Int      *map[uint16]int
	FMapUint16Int8        map[uint16]int8
	FptrMapUint16Int8     *map[uint16]int8
	FMapUint16Int16       map[uint16]int16
	FptrMapUint16Int16    *map[uint16]int16
	FMapUint16Int32       map[uint16]int32
	FptrMapUint16Int32    *map[uint16]int32
	FMapUint16Int64       map[uint16]int64
	FptrMapUint16Int64    *map[uint16]int64
	FMapUint16Float32     map[uint16]float32
	FptrMapUint16Float32  *map[uint16]float32
	FMapUint16Float64     map[uint16]float64
	FptrMapUint16Float64  *map[uint16]float64
	FMapUint16Bool        map[uint16]bool
	FptrMapUint16Bool     *map[uint16]bool
	FMapUint32Intf        map[uint32]interface{}
	FptrMapUint32Intf     *map[uint32]interface{}
	FMapUint32String      map[uint32]string
	FptrMapUint32String   *map[uint32]string
	FMapUint32Uint        map[uint32]uint
	FptrMapUint32Uint     *map[uint32]uint
	FMapUint32Uint8       map[uint32]uint8
	FptrMapUint32Uint8    *map[uint32]uint8
	FMapUint32Uint16      map[uint32]uint16
	FptrMapUint32Uint16   *map[uint32]uint16
	FMapUint32Uint32      map[uint32]uint32
	FptrMapUint32Uint32   *map[uint32]uint32
	FMapUint32Uint64      map[uint32]uint64
	FptrMapUint32Uint64   *map[uint32]uint64
	FMapUint32Uintptr     map[uint32]uintptr
	FptrMapUint32Uintptr  *map[uint32]uintptr
	FMapUint32Int         map[uint32]int
	FptrMapUint32Int      *map[uint32]int
	FMapUint32Int8        map[uint32]int8
	FptrMapUint32Int8     *map[uint32]int8
	FMapUint32Int16       map[uint32]int16
	FptrMapUint32Int16    *map[uint32]int16
	FMapUint32Int32       map[uint32]int32
	FptrMapUint32Int32    *map[uint32]int32
	FMapUint32Int64       map[uint32]int64
	FptrMapUint32Int64    *map[uint32]int64
	FMapUint32Float32     map[uint32]float32
	FptrMapUint32Float32  *map[uint32]float32
	FMapUint32Float64     map[uint32]float64
	FptrMapUint32Float64  *map[uint32]float64
	FMapUint32Bool        map[uint32]bool
	FptrMapUint32Bool     *map[uint32]bool
	FMapUint64Intf        map[uint64]interface{}
	FptrMapUint64Intf     *map[uint64]interface{}
	FMapUint64String      map[uint64]string
	FptrMapUint64String   *map[uint64]string
	FMapUint64Uint        map[uint64]uint
	FptrMapUint64Uint     *map[uint64]uint
	FMapUint64Uint8       map[uint64]uint8
	FptrMapUint64Uint8    *map[uint64]uint8
	FMapUint64Uint16      map[uint64]uint16
	FptrMapUint64Uint16   *map[uint64]uint16
	FMapUint64Uint32      map[uint64]uint32
	FptrMapUint64Uint32   *map[uint64]uint32
	FMapUint64Uint64      map[uint64]uint64
	FptrMapUint64Uint64   *map[uint64]uint64
	FMapUint64Uintptr     map[uint64]uintptr
	FptrMapUint64Uintptr  *map[uint64]uintptr
	FMapUint64Int         map[uint64]int
	FptrMapUint64Int      *map[uint64]int
	FMapUint64Int8        map[uint64]int8
	FptrMapUint64Int8     *map[uint64]int8
	FMapUint64Int16       map[uint64]int16
	FptrMapUint64Int16    *map[uint64]int16
	FMapUint64Int32       map[uint64]int32
	FptrMapUint64Int32    *map[uint64]int32
	FMapUint64Int64       map[uint64]int64
	FptrMapUint64Int64    *map[uint64]int64
	FMapUint64Float32     map[uint64]float32
	FptrMapUint64Float32  *map[uint64]float32
	FMapUint64Float64     map[uint64]float64
	FptrMapUint64Float64  *map[uint64]float64
	FMapUint64Bool        map[uint64]bool
	FptrMapUint64Bool     *map[uint64]bool
	FMapUintptrIntf       map[uintptr]interface{}
	FptrMapUintptrIntf    *map[uintptr]interface{}
	FMapUintptrString     map[uintptr]string
	FptrMapUintptrString  *map[uintptr]string
	FMapUintptrUint       map[uintptr]uint
	FptrMapUintptrUint    *map[uintptr]uint
	FMapUintptrUint8      map[uintptr]uint8
	FptrMapUintptrUint8   *map[uintptr]uint8
	FMapUintptrUint16     map[uintptr]uint16
	FptrMapUintptrUint16  *map[uintptr]uint16
	FMapUintptrUint32     map[uintptr]uint32
	FptrMapUintptrUint32  *map[uintptr]uint32
	FMapUintptrUint64     map[uintptr]uint64
	FptrMapUintptrUint64  *map[uintptr]uint64
	FMapUintptrUintptr    map[uintptr]uintptr
	FptrMapUintptrUintptr *map[uintptr]uintptr
	FMapUintptrInt        map[uintptr]int
	FptrMapUintptrInt     *map[uintptr]int
	FMapUintptrInt8       map[uintptr]int8
	FptrMapUintptrInt8    *map[uintptr]int8
	FMapUintptrInt16      map[uintptr]int16
	FptrMapUintptrInt16   *map[uintptr]int16
	FMapUintptrInt32      map[uintptr]int32
	FptrMapUintptrInt32   *map[uintptr]int32
	FMapUintptrInt64      map[uintptr]int64
	FptrMapUintptrInt64   *map[uintptr]int64
	FMapUintptrFloat32    map[uintptr]float32
	FptrMapUintptrFloat32 *map[uintptr]float32
	FMapUintptrFloat64    map[uintptr]float64
	FptrMapUintptrFloat64 *map[uintptr]float64
	FMapUintptrBool       map[uintptr]bool
	FptrMapUintptrBool    *map[uintptr]bool
	FMapIntIntf           map[int]interface{}
	FptrMapIntIntf        *map[int]interface{}
	FMapIntString         map[int]string
	FptrMapIntString      *map[int]string
	FMapIntUint           map[int]uint
	FptrMapIntUint        *map[int]uint
	FMapIntUint8          map[int]uint8
	FptrMapIntUint8       *map[int]uint8
	FMapIntUint16         map[int]uint16
	FptrMapIntUint16      *map[int]uint16
	FMapIntUint32         map[int]uint32
	FptrMapIntUint32      *map[int]uint32
	FMapIntUint64         map[int]uint64
	FptrMapIntUint64      *map[int]uint64
	FMapIntUintptr        map[int]uintptr
	FptrMapIntUintptr     *map[int]uintptr
	FMapIntInt            map[int]int
	FptrMapIntInt         *map[int]int
	FMapIntInt8           map[int]int8
	FptrMapIntInt8        *map[int]int8
	FMapIntInt16          map[int]int16
	FptrMapIntInt16       *map[int]int16
	FMapIntInt32          map[int]int32
	FptrMapIntInt32       *map[int]int32
	FMapIntInt64          map[int]int64
	FptrMapIntInt64       *map[int]int64
	FMapIntFloat32        map[int]float32
	FptrMapIntFloat32     *map[int]float32
	FMapIntFloat64        map[int]float64
	FptrMapIntFloat64     *map[int]float64
	FMapIntBool           map[int]bool
	FptrMapIntBool        *map[int]bool
	FMapInt8Intf          map[int8]interface{}
	FptrMapInt8Intf       *map[int8]interface{}
	FMapInt8String        map[int8]string
	FptrMapInt8String     *map[int8]string
	FMapInt8Uint          map[int8]uint
	FptrMapInt8Uint       *map[int8]uint
	FMapInt8Uint8         map[int8]uint8
	FptrMapInt8Uint8      *map[int8]uint8
	FMapInt8Uint16        map[int8]uint16
	FptrMapInt8Uint16     *map[int8]uint16
	FMapInt8Uint32        map[int8]uint32
	FptrMapInt8Uint32     *map[int8]uint32
	FMapInt8Uint64        map[int8]uint64
	FptrMapInt8Uint64     *map[int8]uint64
	FMapInt8Uintptr       map[int8]uintptr
	FptrMapInt8Uintptr    *map[int8]uintptr
	FMapInt8Int           map[int8]int
	FptrMapInt8Int        *map[int8]int
	FMapInt8Int8          map[int8]int8
	FptrMapInt8Int8       *map[int8]int8
	FMapInt8Int16         map[int8]int16
	FptrMapInt8Int16      *map[int8]int16
	FMapInt8Int32         map[int8]int32
	FptrMapInt8Int32      *map[int8]int32
	FMapInt8Int64         map[int8]int64
	FptrMapInt8Int64      *map[int8]int64
	FMapInt8Float32       map[int8]float32
	FptrMapInt8Float32    *map[int8]float32
	FMapInt8Float64       map[int8]float64
	FptrMapInt8Float64    *map[int8]float64
	FMapInt8Bool          map[int8]bool
	FptrMapInt8Bool       *map[int8]bool
	FMapInt16Intf         map[int16]interface{}
	FptrMapInt16Intf      *map[int16]interface{}
	FMapInt16String       map[int16]string
	FptrMapInt16String    *map[int16]string
	FMapInt16Uint         map[int16]uint
	FptrMapInt16Uint      *map[int16]uint
	FMapInt16Uint8        map[int16]uint8
	FptrMapInt16Uint8     *map[int16]uint8
	FMapInt16Uint16       map[int16]uint16
	FptrMapInt16Uint16    *map[int16]uint16
	FMapInt16Uint32       map[int16]uint32
	FptrMapInt16Uint32    *map[int16]uint32
	FMapInt16Uint64       map[int16]uint64
	FptrMapInt16Uint64    *map[int16]uint64
	FMapInt16Uintptr      map[int16]uintptr
	FptrMapInt16Uintptr   *map[int16]uintptr
	FMapInt16Int          map[int16]int
	FptrMapInt16Int       *map[int16]int
	FMapInt16Int8         map[int16]int8
	FptrMapInt16Int8      *map[int16]int8
	FMapInt16Int16        map[int16]int16
	FptrMapInt16Int16     *map[int16]int16
	FMapInt16Int32        map[int16]int32
	FptrMapInt16Int32     *map[int16]int32
	FMapInt16Int64        map[int16]int64
	FptrMapInt16Int64     *map[int16]int64
	FMapInt16Float32      map[int16]float32
	FptrMapInt16Float32   *map[int16]float32
	FMapInt16Float64      map[int16]float64
	FptrMapInt16Float64   *map[int16]float64
	FMapInt16Bool         map[int16]bool
	FptrMapInt16Bool      *map[int16]bool
	FMapInt32Intf         map[int32]interface{}
	FptrMapInt32Intf      *map[int32]interface{}
	FMapInt32String       map[int32]string
	FptrMapInt32String    *map[int32]string
	FMapInt32Uint         map[int32]uint
	FptrMapInt32Uint      *map[int32]uint
	FMapInt32Uint8        map[int32]uint8
	FptrMapInt32Uint8     *map[int32]uint8
	FMapInt32Uint16       map[int32]uint16
	FptrMapInt32Uint16    *map[int32]uint16
	FMapInt32Uint32       map[int32]uint32
	FptrMapInt32Uint32    *map[int32]uint32
	FMapInt32Uint64       map[int32]uint64
	FptrMapInt32Uint64    *map[int32]uint64
	FMapInt32Uintptr      map[int32]uintptr
	FptrMapInt32Uintptr   *map[int32]uintptr
	FMapInt32Int          map[int32]int
	FptrMapInt32Int       *map[int32]int
	FMapInt32Int8         map[int32]int8
	FptrMapInt32Int8      *map[int32]int8
	FMapInt32Int16        map[int32]int16
	FptrMapInt32Int16     *map[int32]int16
	FMapInt32Int32        map[int32]int32
	FptrMapInt32Int32     *map[int32]int32
	FMapInt32Int64        map[int32]int64
	FptrMapInt32Int64     *map[int32]int64
	FMapInt32Float32      map[int32]float32
	FptrMapInt32Float32   *map[int32]float32
	FMapInt32Float64      map[int32]float64
	FptrMapInt32Float64   *map[int32]float64
	FMapInt32Bool         map[int32]bool
	FptrMapInt32Bool      *map[int32]bool
	FMapInt64Intf         map[int64]interface{}
	FptrMapInt64Intf      *map[int64]interface{}
	FMapInt64String       map[int64]string
	FptrMapInt64String    *map[int64]string
	FMapInt64Uint         map[int64]uint
	FptrMapInt64Uint      *map[int64]uint
	FMapInt64Uint8        map[int64]uint8
	FptrMapInt64Uint8     *map[int64]uint8
	FMapInt64Uint16       map[int64]uint16
	FptrMapInt64Uint16    *map[int64]uint16
	FMapInt64Uint32       map[int64]uint32
	FptrMapInt64Uint32    *map[int64]uint32
	FMapInt64Uint64       map[int64]uint64
	FptrMapInt64Uint64    *map[int64]uint64
	FMapInt64Uintptr      map[int64]uintptr
	FptrMapInt64Uintptr   *map[int64]uintptr
	FMapInt64Int          map[int64]int
	FptrMapInt64Int       *map[int64]int
	FMapInt64Int8         map[int64]int8
	FptrMapInt64Int8      *map[int64]int8
	FMapInt64Int16        map[int64]int16
	FptrMapInt64Int16     *map[int64]int16
	FMapInt64Int32        map[int64]int32
	FptrMapInt64Int32     *map[int64]int32
	FMapInt64Int64        map[int64]int64
	FptrMapInt64Int64     *map[int64]int64
	FMapInt64Float32      map[int64]float32
	FptrMapInt64Float32   *map[int64]float32
	FMapInt64Float64      map[int64]float64
	FptrMapInt64Float64   *map[int64]float64
	FMapInt64Bool         map[int64]bool
	FptrMapInt64Bool      *map[int64]bool
	FMapBoolIntf          map[bool]interface{}
	FptrMapBoolIntf       *map[bool]interface{}
	FMapBoolString        map[bool]string
	FptrMapBoolString     *map[bool]string
	FMapBoolUint          map[bool]uint
	FptrMapBoolUint       *map[bool]uint
	FMapBoolUint8         map[bool]uint8
	FptrMapBoolUint8      *map[bool]uint8
	FMapBoolUint16        map[bool]uint16
	FptrMapBoolUint16     *map[bool]uint16
	FMapBoolUint32        map[bool]uint32
	FptrMapBoolUint32     *map[bool]uint32
	FMapBoolUint64        map[bool]uint64
	FptrMapBoolUint64     *map[bool]uint64
	FMapBoolUintptr       map[bool]uintptr
	FptrMapBoolUintptr    *map[bool]uintptr
	FMapBoolInt           map[bool]int
	FptrMapBoolInt        *map[bool]int
	FMapBoolInt8          map[bool]int8
	FptrMapBoolInt8       *map[bool]int8
	FMapBoolInt16         map[bool]int16
	FptrMapBoolInt16      *map[bool]int16
	FMapBoolInt32         map[bool]int32
	FptrMapBoolInt32      *map[bool]int32
	FMapBoolInt64         map[bool]int64
	FptrMapBoolInt64      *map[bool]int64
	FMapBoolFloat32       map[bool]float32
	FptrMapBoolFloat32    *map[bool]float32
	FMapBoolFloat64       map[bool]float64
	FptrMapBoolFloat64    *map[bool]float64
	FMapBoolBool          map[bool]bool
	FptrMapBoolBool       *map[bool]bool
}

type typeSliceIntf []interface{}

func (_ typeSliceIntf) MapBySlice() {}

type typeSliceString []string

func (_ typeSliceString) MapBySlice() {}

type typeSliceFloat32 []float32

func (_ typeSliceFloat32) MapBySlice() {}

type typeSliceFloat64 []float64

func (_ typeSliceFloat64) MapBySlice() {}

type typeSliceUint []uint

func (_ typeSliceUint) MapBySlice() {}

type typeSliceUint16 []uint16

func (_ typeSliceUint16) MapBySlice() {}

type typeSliceUint32 []uint32

func (_ typeSliceUint32) MapBySlice() {}

type typeSliceUint64 []uint64

func (_ typeSliceUint64) MapBySlice() {}

type typeSliceUintptr []uintptr

func (_ typeSliceUintptr) MapBySlice() {}

type typeSliceInt []int

func (_ typeSliceInt) MapBySlice() {}

type typeSliceInt8 []int8

func (_ typeSliceInt8) MapBySlice() {}

type typeSliceInt16 []int16

func (_ typeSliceInt16) MapBySlice() {}

type typeSliceInt32 []int32

func (_ typeSliceInt32) MapBySlice() {}

type typeSliceInt64 []int64

func (_ typeSliceInt64) MapBySlice() {}

type typeSliceBool []bool

func (_ typeSliceBool) MapBySlice() {}

func doTestMammothSlices(t *testing.T, h Handle) {

	for _, v := range [][]interface{}{nil, []interface{}{}, []interface{}{"string-is-an-interface", "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth slice v1: %v\n", v)
		var v1v1, v1v2, v1v3, v1v4 []interface{}
		v1v1 = v
		bs1 := testMarshalErr(v1v1, h, t, "enc-slice-v1")
		if v != nil {
			v1v2 = make([]interface{}, len(v))
		}
		testUnmarshalErr(v1v2, bs1, h, t, "dec-slice-v1")
		testDeepEqualErr(v1v1, v1v2, t, "equal-slice-v1")
		bs1 = testMarshalErr(&v1v1, h, t, "enc-slice-v1-p")
		v1v2 = nil
		testUnmarshalErr(&v1v2, bs1, h, t, "dec-slice-v1-p")
		testDeepEqualErr(v1v1, v1v2, t, "equal-slice-v1-p")
		// ...
		v1v2 = nil
		if v != nil {
			v1v2 = make([]interface{}, len(v))
		}
		v1v3 = typeSliceIntf(v1v1)
		bs1 = testMarshalErr(v1v3, h, t, "enc-slice-v1-custom")
		v1v4 = typeSliceIntf(v1v2)
		testUnmarshalErr(v1v4, bs1, h, t, "dec-slice-v1-custom")
		testDeepEqualErr(v1v3, v1v4, t, "equal-slice-v1-custom")
		v1v2 = nil
		bs1 = testMarshalErr(&v1v3, h, t, "enc-slice-v1-custom-p")
		v1v4 = typeSliceIntf(v1v2)
		testUnmarshalErr(&v1v4, bs1, h, t, "dec-slice-v1-custom-p")
		testDeepEqualErr(v1v3, v1v4, t, "equal-slice-v1-custom-p")
	}

	for _, v := range [][]string{nil, []string{}, []string{"some-string", "some-string"}} {
		// fmt.Printf(">>>> running mammoth slice v19: %v\n", v)
		var v19v1, v19v2, v19v3, v19v4 []string
		v19v1 = v
		bs19 := testMarshalErr(v19v1, h, t, "enc-slice-v19")
		if v != nil {
			v19v2 = make([]string, len(v))
		}
		testUnmarshalErr(v19v2, bs19, h, t, "dec-slice-v19")
		testDeepEqualErr(v19v1, v19v2, t, "equal-slice-v19")
		bs19 = testMarshalErr(&v19v1, h, t, "enc-slice-v19-p")
		v19v2 = nil
		testUnmarshalErr(&v19v2, bs19, h, t, "dec-slice-v19-p")
		testDeepEqualErr(v19v1, v19v2, t, "equal-slice-v19-p")
		// ...
		v19v2 = nil
		if v != nil {
			v19v2 = make([]string, len(v))
		}
		v19v3 = typeSliceString(v19v1)
		bs19 = testMarshalErr(v19v3, h, t, "enc-slice-v19-custom")
		v19v4 = typeSliceString(v19v2)
		testUnmarshalErr(v19v4, bs19, h, t, "dec-slice-v19-custom")
		testDeepEqualErr(v19v3, v19v4, t, "equal-slice-v19-custom")
		v19v2 = nil
		bs19 = testMarshalErr(&v19v3, h, t, "enc-slice-v19-custom-p")
		v19v4 = typeSliceString(v19v2)
		testUnmarshalErr(&v19v4, bs19, h, t, "dec-slice-v19-custom-p")
		testDeepEqualErr(v19v3, v19v4, t, "equal-slice-v19-custom-p")
	}

	for _, v := range [][]float32{nil, []float32{}, []float32{10.1, 10.1}} {
		// fmt.Printf(">>>> running mammoth slice v37: %v\n", v)
		var v37v1, v37v2, v37v3, v37v4 []float32
		v37v1 = v
		bs37 := testMarshalErr(v37v1, h, t, "enc-slice-v37")
		if v != nil {
			v37v2 = make([]float32, len(v))
		}
		testUnmarshalErr(v37v2, bs37, h, t, "dec-slice-v37")
		testDeepEqualErr(v37v1, v37v2, t, "equal-slice-v37")
		bs37 = testMarshalErr(&v37v1, h, t, "enc-slice-v37-p")
		v37v2 = nil
		testUnmarshalErr(&v37v2, bs37, h, t, "dec-slice-v37-p")
		testDeepEqualErr(v37v1, v37v2, t, "equal-slice-v37-p")
		// ...
		v37v2 = nil
		if v != nil {
			v37v2 = make([]float32, len(v))
		}
		v37v3 = typeSliceFloat32(v37v1)
		bs37 = testMarshalErr(v37v3, h, t, "enc-slice-v37-custom")
		v37v4 = typeSliceFloat32(v37v2)
		testUnmarshalErr(v37v4, bs37, h, t, "dec-slice-v37-custom")
		testDeepEqualErr(v37v3, v37v4, t, "equal-slice-v37-custom")
		v37v2 = nil
		bs37 = testMarshalErr(&v37v3, h, t, "enc-slice-v37-custom-p")
		v37v4 = typeSliceFloat32(v37v2)
		testUnmarshalErr(&v37v4, bs37, h, t, "dec-slice-v37-custom-p")
		testDeepEqualErr(v37v3, v37v4, t, "equal-slice-v37-custom-p")
	}

	for _, v := range [][]float64{nil, []float64{}, []float64{10.1, 10.1}} {
		// fmt.Printf(">>>> running mammoth slice v55: %v\n", v)
		var v55v1, v55v2, v55v3, v55v4 []float64
		v55v1 = v
		bs55 := testMarshalErr(v55v1, h, t, "enc-slice-v55")
		if v != nil {
			v55v2 = make([]float64, len(v))
		}
		testUnmarshalErr(v55v2, bs55, h, t, "dec-slice-v55")
		testDeepEqualErr(v55v1, v55v2, t, "equal-slice-v55")
		bs55 = testMarshalErr(&v55v1, h, t, "enc-slice-v55-p")
		v55v2 = nil
		testUnmarshalErr(&v55v2, bs55, h, t, "dec-slice-v55-p")
		testDeepEqualErr(v55v1, v55v2, t, "equal-slice-v55-p")
		// ...
		v55v2 = nil
		if v != nil {
			v55v2 = make([]float64, len(v))
		}
		v55v3 = typeSliceFloat64(v55v1)
		bs55 = testMarshalErr(v55v3, h, t, "enc-slice-v55-custom")
		v55v4 = typeSliceFloat64(v55v2)
		testUnmarshalErr(v55v4, bs55, h, t, "dec-slice-v55-custom")
		testDeepEqualErr(v55v3, v55v4, t, "equal-slice-v55-custom")
		v55v2 = nil
		bs55 = testMarshalErr(&v55v3, h, t, "enc-slice-v55-custom-p")
		v55v4 = typeSliceFloat64(v55v2)
		testUnmarshalErr(&v55v4, bs55, h, t, "dec-slice-v55-custom-p")
		testDeepEqualErr(v55v3, v55v4, t, "equal-slice-v55-custom-p")
	}

	for _, v := range [][]uint{nil, []uint{}, []uint{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v73: %v\n", v)
		var v73v1, v73v2, v73v3, v73v4 []uint
		v73v1 = v
		bs73 := testMarshalErr(v73v1, h, t, "enc-slice-v73")
		if v != nil {
			v73v2 = make([]uint, len(v))
		}
		testUnmarshalErr(v73v2, bs73, h, t, "dec-slice-v73")
		testDeepEqualErr(v73v1, v73v2, t, "equal-slice-v73")
		bs73 = testMarshalErr(&v73v1, h, t, "enc-slice-v73-p")
		v73v2 = nil
		testUnmarshalErr(&v73v2, bs73, h, t, "dec-slice-v73-p")
		testDeepEqualErr(v73v1, v73v2, t, "equal-slice-v73-p")
		// ...
		v73v2 = nil
		if v != nil {
			v73v2 = make([]uint, len(v))
		}
		v73v3 = typeSliceUint(v73v1)
		bs73 = testMarshalErr(v73v3, h, t, "enc-slice-v73-custom")
		v73v4 = typeSliceUint(v73v2)
		testUnmarshalErr(v73v4, bs73, h, t, "dec-slice-v73-custom")
		testDeepEqualErr(v73v3, v73v4, t, "equal-slice-v73-custom")
		v73v2 = nil
		bs73 = testMarshalErr(&v73v3, h, t, "enc-slice-v73-custom-p")
		v73v4 = typeSliceUint(v73v2)
		testUnmarshalErr(&v73v4, bs73, h, t, "dec-slice-v73-custom-p")
		testDeepEqualErr(v73v3, v73v4, t, "equal-slice-v73-custom-p")
	}

	for _, v := range [][]uint16{nil, []uint16{}, []uint16{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v108: %v\n", v)
		var v108v1, v108v2, v108v3, v108v4 []uint16
		v108v1 = v
		bs108 := testMarshalErr(v108v1, h, t, "enc-slice-v108")
		if v != nil {
			v108v2 = make([]uint16, len(v))
		}
		testUnmarshalErr(v108v2, bs108, h, t, "dec-slice-v108")
		testDeepEqualErr(v108v1, v108v2, t, "equal-slice-v108")
		bs108 = testMarshalErr(&v108v1, h, t, "enc-slice-v108-p")
		v108v2 = nil
		testUnmarshalErr(&v108v2, bs108, h, t, "dec-slice-v108-p")
		testDeepEqualErr(v108v1, v108v2, t, "equal-slice-v108-p")
		// ...
		v108v2 = nil
		if v != nil {
			v108v2 = make([]uint16, len(v))
		}
		v108v3 = typeSliceUint16(v108v1)
		bs108 = testMarshalErr(v108v3, h, t, "enc-slice-v108-custom")
		v108v4 = typeSliceUint16(v108v2)
		testUnmarshalErr(v108v4, bs108, h, t, "dec-slice-v108-custom")
		testDeepEqualErr(v108v3, v108v4, t, "equal-slice-v108-custom")
		v108v2 = nil
		bs108 = testMarshalErr(&v108v3, h, t, "enc-slice-v108-custom-p")
		v108v4 = typeSliceUint16(v108v2)
		testUnmarshalErr(&v108v4, bs108, h, t, "dec-slice-v108-custom-p")
		testDeepEqualErr(v108v3, v108v4, t, "equal-slice-v108-custom-p")
	}

	for _, v := range [][]uint32{nil, []uint32{}, []uint32{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v126: %v\n", v)
		var v126v1, v126v2, v126v3, v126v4 []uint32
		v126v1 = v
		bs126 := testMarshalErr(v126v1, h, t, "enc-slice-v126")
		if v != nil {
			v126v2 = make([]uint32, len(v))
		}
		testUnmarshalErr(v126v2, bs126, h, t, "dec-slice-v126")
		testDeepEqualErr(v126v1, v126v2, t, "equal-slice-v126")
		bs126 = testMarshalErr(&v126v1, h, t, "enc-slice-v126-p")
		v126v2 = nil
		testUnmarshalErr(&v126v2, bs126, h, t, "dec-slice-v126-p")
		testDeepEqualErr(v126v1, v126v2, t, "equal-slice-v126-p")
		// ...
		v126v2 = nil
		if v != nil {
			v126v2 = make([]uint32, len(v))
		}
		v126v3 = typeSliceUint32(v126v1)
		bs126 = testMarshalErr(v126v3, h, t, "enc-slice-v126-custom")
		v126v4 = typeSliceUint32(v126v2)
		testUnmarshalErr(v126v4, bs126, h, t, "dec-slice-v126-custom")
		testDeepEqualErr(v126v3, v126v4, t, "equal-slice-v126-custom")
		v126v2 = nil
		bs126 = testMarshalErr(&v126v3, h, t, "enc-slice-v126-custom-p")
		v126v4 = typeSliceUint32(v126v2)
		testUnmarshalErr(&v126v4, bs126, h, t, "dec-slice-v126-custom-p")
		testDeepEqualErr(v126v3, v126v4, t, "equal-slice-v126-custom-p")
	}

	for _, v := range [][]uint64{nil, []uint64{}, []uint64{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v144: %v\n", v)
		var v144v1, v144v2, v144v3, v144v4 []uint64
		v144v1 = v
		bs144 := testMarshalErr(v144v1, h, t, "enc-slice-v144")
		if v != nil {
			v144v2 = make([]uint64, len(v))
		}
		testUnmarshalErr(v144v2, bs144, h, t, "dec-slice-v144")
		testDeepEqualErr(v144v1, v144v2, t, "equal-slice-v144")
		bs144 = testMarshalErr(&v144v1, h, t, "enc-slice-v144-p")
		v144v2 = nil
		testUnmarshalErr(&v144v2, bs144, h, t, "dec-slice-v144-p")
		testDeepEqualErr(v144v1, v144v2, t, "equal-slice-v144-p")
		// ...
		v144v2 = nil
		if v != nil {
			v144v2 = make([]uint64, len(v))
		}
		v144v3 = typeSliceUint64(v144v1)
		bs144 = testMarshalErr(v144v3, h, t, "enc-slice-v144-custom")
		v144v4 = typeSliceUint64(v144v2)
		testUnmarshalErr(v144v4, bs144, h, t, "dec-slice-v144-custom")
		testDeepEqualErr(v144v3, v144v4, t, "equal-slice-v144-custom")
		v144v2 = nil
		bs144 = testMarshalErr(&v144v3, h, t, "enc-slice-v144-custom-p")
		v144v4 = typeSliceUint64(v144v2)
		testUnmarshalErr(&v144v4, bs144, h, t, "dec-slice-v144-custom-p")
		testDeepEqualErr(v144v3, v144v4, t, "equal-slice-v144-custom-p")
	}

	for _, v := range [][]uintptr{nil, []uintptr{}, []uintptr{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v162: %v\n", v)
		var v162v1, v162v2, v162v3, v162v4 []uintptr
		v162v1 = v
		bs162 := testMarshalErr(v162v1, h, t, "enc-slice-v162")
		if v != nil {
			v162v2 = make([]uintptr, len(v))
		}
		testUnmarshalErr(v162v2, bs162, h, t, "dec-slice-v162")
		testDeepEqualErr(v162v1, v162v2, t, "equal-slice-v162")
		bs162 = testMarshalErr(&v162v1, h, t, "enc-slice-v162-p")
		v162v2 = nil
		testUnmarshalErr(&v162v2, bs162, h, t, "dec-slice-v162-p")
		testDeepEqualErr(v162v1, v162v2, t, "equal-slice-v162-p")
		// ...
		v162v2 = nil
		if v != nil {
			v162v2 = make([]uintptr, len(v))
		}
		v162v3 = typeSliceUintptr(v162v1)
		bs162 = testMarshalErr(v162v3, h, t, "enc-slice-v162-custom")
		v162v4 = typeSliceUintptr(v162v2)
		testUnmarshalErr(v162v4, bs162, h, t, "dec-slice-v162-custom")
		testDeepEqualErr(v162v3, v162v4, t, "equal-slice-v162-custom")
		v162v2 = nil
		bs162 = testMarshalErr(&v162v3, h, t, "enc-slice-v162-custom-p")
		v162v4 = typeSliceUintptr(v162v2)
		testUnmarshalErr(&v162v4, bs162, h, t, "dec-slice-v162-custom-p")
		testDeepEqualErr(v162v3, v162v4, t, "equal-slice-v162-custom-p")
	}

	for _, v := range [][]int{nil, []int{}, []int{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v180: %v\n", v)
		var v180v1, v180v2, v180v3, v180v4 []int
		v180v1 = v
		bs180 := testMarshalErr(v180v1, h, t, "enc-slice-v180")
		if v != nil {
			v180v2 = make([]int, len(v))
		}
		testUnmarshalErr(v180v2, bs180, h, t, "dec-slice-v180")
		testDeepEqualErr(v180v1, v180v2, t, "equal-slice-v180")
		bs180 = testMarshalErr(&v180v1, h, t, "enc-slice-v180-p")
		v180v2 = nil
		testUnmarshalErr(&v180v2, bs180, h, t, "dec-slice-v180-p")
		testDeepEqualErr(v180v1, v180v2, t, "equal-slice-v180-p")
		// ...
		v180v2 = nil
		if v != nil {
			v180v2 = make([]int, len(v))
		}
		v180v3 = typeSliceInt(v180v1)
		bs180 = testMarshalErr(v180v3, h, t, "enc-slice-v180-custom")
		v180v4 = typeSliceInt(v180v2)
		testUnmarshalErr(v180v4, bs180, h, t, "dec-slice-v180-custom")
		testDeepEqualErr(v180v3, v180v4, t, "equal-slice-v180-custom")
		v180v2 = nil
		bs180 = testMarshalErr(&v180v3, h, t, "enc-slice-v180-custom-p")
		v180v4 = typeSliceInt(v180v2)
		testUnmarshalErr(&v180v4, bs180, h, t, "dec-slice-v180-custom-p")
		testDeepEqualErr(v180v3, v180v4, t, "equal-slice-v180-custom-p")
	}

	for _, v := range [][]int8{nil, []int8{}, []int8{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v198: %v\n", v)
		var v198v1, v198v2, v198v3, v198v4 []int8
		v198v1 = v
		bs198 := testMarshalErr(v198v1, h, t, "enc-slice-v198")
		if v != nil {
			v198v2 = make([]int8, len(v))
		}
		testUnmarshalErr(v198v2, bs198, h, t, "dec-slice-v198")
		testDeepEqualErr(v198v1, v198v2, t, "equal-slice-v198")
		bs198 = testMarshalErr(&v198v1, h, t, "enc-slice-v198-p")
		v198v2 = nil
		testUnmarshalErr(&v198v2, bs198, h, t, "dec-slice-v198-p")
		testDeepEqualErr(v198v1, v198v2, t, "equal-slice-v198-p")
		// ...
		v198v2 = nil
		if v != nil {
			v198v2 = make([]int8, len(v))
		}
		v198v3 = typeSliceInt8(v198v1)
		bs198 = testMarshalErr(v198v3, h, t, "enc-slice-v198-custom")
		v198v4 = typeSliceInt8(v198v2)
		testUnmarshalErr(v198v4, bs198, h, t, "dec-slice-v198-custom")
		testDeepEqualErr(v198v3, v198v4, t, "equal-slice-v198-custom")
		v198v2 = nil
		bs198 = testMarshalErr(&v198v3, h, t, "enc-slice-v198-custom-p")
		v198v4 = typeSliceInt8(v198v2)
		testUnmarshalErr(&v198v4, bs198, h, t, "dec-slice-v198-custom-p")
		testDeepEqualErr(v198v3, v198v4, t, "equal-slice-v198-custom-p")
	}

	for _, v := range [][]int16{nil, []int16{}, []int16{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v216: %v\n", v)
		var v216v1, v216v2, v216v3, v216v4 []int16
		v216v1 = v
		bs216 := testMarshalErr(v216v1, h, t, "enc-slice-v216")
		if v != nil {
			v216v2 = make([]int16, len(v))
		}
		testUnmarshalErr(v216v2, bs216, h, t, "dec-slice-v216")
		testDeepEqualErr(v216v1, v216v2, t, "equal-slice-v216")
		bs216 = testMarshalErr(&v216v1, h, t, "enc-slice-v216-p")
		v216v2 = nil
		testUnmarshalErr(&v216v2, bs216, h, t, "dec-slice-v216-p")
		testDeepEqualErr(v216v1, v216v2, t, "equal-slice-v216-p")
		// ...
		v216v2 = nil
		if v != nil {
			v216v2 = make([]int16, len(v))
		}
		v216v3 = typeSliceInt16(v216v1)
		bs216 = testMarshalErr(v216v3, h, t, "enc-slice-v216-custom")
		v216v4 = typeSliceInt16(v216v2)
		testUnmarshalErr(v216v4, bs216, h, t, "dec-slice-v216-custom")
		testDeepEqualErr(v216v3, v216v4, t, "equal-slice-v216-custom")
		v216v2 = nil
		bs216 = testMarshalErr(&v216v3, h, t, "enc-slice-v216-custom-p")
		v216v4 = typeSliceInt16(v216v2)
		testUnmarshalErr(&v216v4, bs216, h, t, "dec-slice-v216-custom-p")
		testDeepEqualErr(v216v3, v216v4, t, "equal-slice-v216-custom-p")
	}

	for _, v := range [][]int32{nil, []int32{}, []int32{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v234: %v\n", v)
		var v234v1, v234v2, v234v3, v234v4 []int32
		v234v1 = v
		bs234 := testMarshalErr(v234v1, h, t, "enc-slice-v234")
		if v != nil {
			v234v2 = make([]int32, len(v))
		}
		testUnmarshalErr(v234v2, bs234, h, t, "dec-slice-v234")
		testDeepEqualErr(v234v1, v234v2, t, "equal-slice-v234")
		bs234 = testMarshalErr(&v234v1, h, t, "enc-slice-v234-p")
		v234v2 = nil
		testUnmarshalErr(&v234v2, bs234, h, t, "dec-slice-v234-p")
		testDeepEqualErr(v234v1, v234v2, t, "equal-slice-v234-p")
		// ...
		v234v2 = nil
		if v != nil {
			v234v2 = make([]int32, len(v))
		}
		v234v3 = typeSliceInt32(v234v1)
		bs234 = testMarshalErr(v234v3, h, t, "enc-slice-v234-custom")
		v234v4 = typeSliceInt32(v234v2)
		testUnmarshalErr(v234v4, bs234, h, t, "dec-slice-v234-custom")
		testDeepEqualErr(v234v3, v234v4, t, "equal-slice-v234-custom")
		v234v2 = nil
		bs234 = testMarshalErr(&v234v3, h, t, "enc-slice-v234-custom-p")
		v234v4 = typeSliceInt32(v234v2)
		testUnmarshalErr(&v234v4, bs234, h, t, "dec-slice-v234-custom-p")
		testDeepEqualErr(v234v3, v234v4, t, "equal-slice-v234-custom-p")
	}

	for _, v := range [][]int64{nil, []int64{}, []int64{10, 10}} {
		// fmt.Printf(">>>> running mammoth slice v252: %v\n", v)
		var v252v1, v252v2, v252v3, v252v4 []int64
		v252v1 = v
		bs252 := testMarshalErr(v252v1, h, t, "enc-slice-v252")
		if v != nil {
			v252v2 = make([]int64, len(v))
		}
		testUnmarshalErr(v252v2, bs252, h, t, "dec-slice-v252")
		testDeepEqualErr(v252v1, v252v2, t, "equal-slice-v252")
		bs252 = testMarshalErr(&v252v1, h, t, "enc-slice-v252-p")
		v252v2 = nil
		testUnmarshalErr(&v252v2, bs252, h, t, "dec-slice-v252-p")
		testDeepEqualErr(v252v1, v252v2, t, "equal-slice-v252-p")
		// ...
		v252v2 = nil
		if v != nil {
			v252v2 = make([]int64, len(v))
		}
		v252v3 = typeSliceInt64(v252v1)
		bs252 = testMarshalErr(v252v3, h, t, "enc-slice-v252-custom")
		v252v4 = typeSliceInt64(v252v2)
		testUnmarshalErr(v252v4, bs252, h, t, "dec-slice-v252-custom")
		testDeepEqualErr(v252v3, v252v4, t, "equal-slice-v252-custom")
		v252v2 = nil
		bs252 = testMarshalErr(&v252v3, h, t, "enc-slice-v252-custom-p")
		v252v4 = typeSliceInt64(v252v2)
		testUnmarshalErr(&v252v4, bs252, h, t, "dec-slice-v252-custom-p")
		testDeepEqualErr(v252v3, v252v4, t, "equal-slice-v252-custom-p")
	}

	for _, v := range [][]bool{nil, []bool{}, []bool{true, true}} {
		// fmt.Printf(">>>> running mammoth slice v270: %v\n", v)
		var v270v1, v270v2, v270v3, v270v4 []bool
		v270v1 = v
		bs270 := testMarshalErr(v270v1, h, t, "enc-slice-v270")
		if v != nil {
			v270v2 = make([]bool, len(v))
		}
		testUnmarshalErr(v270v2, bs270, h, t, "dec-slice-v270")
		testDeepEqualErr(v270v1, v270v2, t, "equal-slice-v270")
		bs270 = testMarshalErr(&v270v1, h, t, "enc-slice-v270-p")
		v270v2 = nil
		testUnmarshalErr(&v270v2, bs270, h, t, "dec-slice-v270-p")
		testDeepEqualErr(v270v1, v270v2, t, "equal-slice-v270-p")
		// ...
		v270v2 = nil
		if v != nil {
			v270v2 = make([]bool, len(v))
		}
		v270v3 = typeSliceBool(v270v1)
		bs270 = testMarshalErr(v270v3, h, t, "enc-slice-v270-custom")
		v270v4 = typeSliceBool(v270v2)
		testUnmarshalErr(v270v4, bs270, h, t, "dec-slice-v270-custom")
		testDeepEqualErr(v270v3, v270v4, t, "equal-slice-v270-custom")
		v270v2 = nil
		bs270 = testMarshalErr(&v270v3, h, t, "enc-slice-v270-custom-p")
		v270v4 = typeSliceBool(v270v2)
		testUnmarshalErr(&v270v4, bs270, h, t, "dec-slice-v270-custom-p")
		testDeepEqualErr(v270v3, v270v4, t, "equal-slice-v270-custom-p")
	}

}

func doTestMammothMaps(t *testing.T, h Handle) {

	for _, v := range []map[interface{}]interface{}{nil, map[interface{}]interface{}{}, map[interface{}]interface{}{"string-is-an-interface": "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v2: %v\n", v)
		var v2v1, v2v2 map[interface{}]interface{}
		v2v1 = v
		bs2 := testMarshalErr(v2v1, h, t, "enc-map-v2")
		if v != nil {
			v2v2 = make(map[interface{}]interface{}, len(v))
		}
		testUnmarshalErr(v2v2, bs2, h, t, "dec-map-v2")
		testDeepEqualErr(v2v1, v2v2, t, "equal-map-v2")
		bs2 = testMarshalErr(&v2v1, h, t, "enc-map-v2-p")
		v2v2 = nil
		testUnmarshalErr(&v2v2, bs2, h, t, "dec-map-v2-p")
		testDeepEqualErr(v2v1, v2v2, t, "equal-map-v2-p")
	}

	for _, v := range []map[interface{}]string{nil, map[interface{}]string{}, map[interface{}]string{"string-is-an-interface": "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v3: %v\n", v)
		var v3v1, v3v2 map[interface{}]string
		v3v1 = v
		bs3 := testMarshalErr(v3v1, h, t, "enc-map-v3")
		if v != nil {
			v3v2 = make(map[interface{}]string, len(v))
		}
		testUnmarshalErr(v3v2, bs3, h, t, "dec-map-v3")
		testDeepEqualErr(v3v1, v3v2, t, "equal-map-v3")
		bs3 = testMarshalErr(&v3v1, h, t, "enc-map-v3-p")
		v3v2 = nil
		testUnmarshalErr(&v3v2, bs3, h, t, "dec-map-v3-p")
		testDeepEqualErr(v3v1, v3v2, t, "equal-map-v3-p")
	}

	for _, v := range []map[interface{}]uint{nil, map[interface{}]uint{}, map[interface{}]uint{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v4: %v\n", v)
		var v4v1, v4v2 map[interface{}]uint
		v4v1 = v
		bs4 := testMarshalErr(v4v1, h, t, "enc-map-v4")
		if v != nil {
			v4v2 = make(map[interface{}]uint, len(v))
		}
		testUnmarshalErr(v4v2, bs4, h, t, "dec-map-v4")
		testDeepEqualErr(v4v1, v4v2, t, "equal-map-v4")
		bs4 = testMarshalErr(&v4v1, h, t, "enc-map-v4-p")
		v4v2 = nil
		testUnmarshalErr(&v4v2, bs4, h, t, "dec-map-v4-p")
		testDeepEqualErr(v4v1, v4v2, t, "equal-map-v4-p")
	}

	for _, v := range []map[interface{}]uint8{nil, map[interface{}]uint8{}, map[interface{}]uint8{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v5: %v\n", v)
		var v5v1, v5v2 map[interface{}]uint8
		v5v1 = v
		bs5 := testMarshalErr(v5v1, h, t, "enc-map-v5")
		if v != nil {
			v5v2 = make(map[interface{}]uint8, len(v))
		}
		testUnmarshalErr(v5v2, bs5, h, t, "dec-map-v5")
		testDeepEqualErr(v5v1, v5v2, t, "equal-map-v5")
		bs5 = testMarshalErr(&v5v1, h, t, "enc-map-v5-p")
		v5v2 = nil
		testUnmarshalErr(&v5v2, bs5, h, t, "dec-map-v5-p")
		testDeepEqualErr(v5v1, v5v2, t, "equal-map-v5-p")
	}

	for _, v := range []map[interface{}]uint16{nil, map[interface{}]uint16{}, map[interface{}]uint16{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v6: %v\n", v)
		var v6v1, v6v2 map[interface{}]uint16
		v6v1 = v
		bs6 := testMarshalErr(v6v1, h, t, "enc-map-v6")
		if v != nil {
			v6v2 = make(map[interface{}]uint16, len(v))
		}
		testUnmarshalErr(v6v2, bs6, h, t, "dec-map-v6")
		testDeepEqualErr(v6v1, v6v2, t, "equal-map-v6")
		bs6 = testMarshalErr(&v6v1, h, t, "enc-map-v6-p")
		v6v2 = nil
		testUnmarshalErr(&v6v2, bs6, h, t, "dec-map-v6-p")
		testDeepEqualErr(v6v1, v6v2, t, "equal-map-v6-p")
	}

	for _, v := range []map[interface{}]uint32{nil, map[interface{}]uint32{}, map[interface{}]uint32{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v7: %v\n", v)
		var v7v1, v7v2 map[interface{}]uint32
		v7v1 = v
		bs7 := testMarshalErr(v7v1, h, t, "enc-map-v7")
		if v != nil {
			v7v2 = make(map[interface{}]uint32, len(v))
		}
		testUnmarshalErr(v7v2, bs7, h, t, "dec-map-v7")
		testDeepEqualErr(v7v1, v7v2, t, "equal-map-v7")
		bs7 = testMarshalErr(&v7v1, h, t, "enc-map-v7-p")
		v7v2 = nil
		testUnmarshalErr(&v7v2, bs7, h, t, "dec-map-v7-p")
		testDeepEqualErr(v7v1, v7v2, t, "equal-map-v7-p")
	}

	for _, v := range []map[interface{}]uint64{nil, map[interface{}]uint64{}, map[interface{}]uint64{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v8: %v\n", v)
		var v8v1, v8v2 map[interface{}]uint64
		v8v1 = v
		bs8 := testMarshalErr(v8v1, h, t, "enc-map-v8")
		if v != nil {
			v8v2 = make(map[interface{}]uint64, len(v))
		}
		testUnmarshalErr(v8v2, bs8, h, t, "dec-map-v8")
		testDeepEqualErr(v8v1, v8v2, t, "equal-map-v8")
		bs8 = testMarshalErr(&v8v1, h, t, "enc-map-v8-p")
		v8v2 = nil
		testUnmarshalErr(&v8v2, bs8, h, t, "dec-map-v8-p")
		testDeepEqualErr(v8v1, v8v2, t, "equal-map-v8-p")
	}

	for _, v := range []map[interface{}]uintptr{nil, map[interface{}]uintptr{}, map[interface{}]uintptr{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v9: %v\n", v)
		var v9v1, v9v2 map[interface{}]uintptr
		v9v1 = v
		bs9 := testMarshalErr(v9v1, h, t, "enc-map-v9")
		if v != nil {
			v9v2 = make(map[interface{}]uintptr, len(v))
		}
		testUnmarshalErr(v9v2, bs9, h, t, "dec-map-v9")
		testDeepEqualErr(v9v1, v9v2, t, "equal-map-v9")
		bs9 = testMarshalErr(&v9v1, h, t, "enc-map-v9-p")
		v9v2 = nil
		testUnmarshalErr(&v9v2, bs9, h, t, "dec-map-v9-p")
		testDeepEqualErr(v9v1, v9v2, t, "equal-map-v9-p")
	}

	for _, v := range []map[interface{}]int{nil, map[interface{}]int{}, map[interface{}]int{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v10: %v\n", v)
		var v10v1, v10v2 map[interface{}]int
		v10v1 = v
		bs10 := testMarshalErr(v10v1, h, t, "enc-map-v10")
		if v != nil {
			v10v2 = make(map[interface{}]int, len(v))
		}
		testUnmarshalErr(v10v2, bs10, h, t, "dec-map-v10")
		testDeepEqualErr(v10v1, v10v2, t, "equal-map-v10")
		bs10 = testMarshalErr(&v10v1, h, t, "enc-map-v10-p")
		v10v2 = nil
		testUnmarshalErr(&v10v2, bs10, h, t, "dec-map-v10-p")
		testDeepEqualErr(v10v1, v10v2, t, "equal-map-v10-p")
	}

	for _, v := range []map[interface{}]int8{nil, map[interface{}]int8{}, map[interface{}]int8{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v11: %v\n", v)
		var v11v1, v11v2 map[interface{}]int8
		v11v1 = v
		bs11 := testMarshalErr(v11v1, h, t, "enc-map-v11")
		if v != nil {
			v11v2 = make(map[interface{}]int8, len(v))
		}
		testUnmarshalErr(v11v2, bs11, h, t, "dec-map-v11")
		testDeepEqualErr(v11v1, v11v2, t, "equal-map-v11")
		bs11 = testMarshalErr(&v11v1, h, t, "enc-map-v11-p")
		v11v2 = nil
		testUnmarshalErr(&v11v2, bs11, h, t, "dec-map-v11-p")
		testDeepEqualErr(v11v1, v11v2, t, "equal-map-v11-p")
	}

	for _, v := range []map[interface{}]int16{nil, map[interface{}]int16{}, map[interface{}]int16{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v12: %v\n", v)
		var v12v1, v12v2 map[interface{}]int16
		v12v1 = v
		bs12 := testMarshalErr(v12v1, h, t, "enc-map-v12")
		if v != nil {
			v12v2 = make(map[interface{}]int16, len(v))
		}
		testUnmarshalErr(v12v2, bs12, h, t, "dec-map-v12")
		testDeepEqualErr(v12v1, v12v2, t, "equal-map-v12")
		bs12 = testMarshalErr(&v12v1, h, t, "enc-map-v12-p")
		v12v2 = nil
		testUnmarshalErr(&v12v2, bs12, h, t, "dec-map-v12-p")
		testDeepEqualErr(v12v1, v12v2, t, "equal-map-v12-p")
	}

	for _, v := range []map[interface{}]int32{nil, map[interface{}]int32{}, map[interface{}]int32{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v13: %v\n", v)
		var v13v1, v13v2 map[interface{}]int32
		v13v1 = v
		bs13 := testMarshalErr(v13v1, h, t, "enc-map-v13")
		if v != nil {
			v13v2 = make(map[interface{}]int32, len(v))
		}
		testUnmarshalErr(v13v2, bs13, h, t, "dec-map-v13")
		testDeepEqualErr(v13v1, v13v2, t, "equal-map-v13")
		bs13 = testMarshalErr(&v13v1, h, t, "enc-map-v13-p")
		v13v2 = nil
		testUnmarshalErr(&v13v2, bs13, h, t, "dec-map-v13-p")
		testDeepEqualErr(v13v1, v13v2, t, "equal-map-v13-p")
	}

	for _, v := range []map[interface{}]int64{nil, map[interface{}]int64{}, map[interface{}]int64{"string-is-an-interface": 10}} {
		// fmt.Printf(">>>> running mammoth map v14: %v\n", v)
		var v14v1, v14v2 map[interface{}]int64
		v14v1 = v
		bs14 := testMarshalErr(v14v1, h, t, "enc-map-v14")
		if v != nil {
			v14v2 = make(map[interface{}]int64, len(v))
		}
		testUnmarshalErr(v14v2, bs14, h, t, "dec-map-v14")
		testDeepEqualErr(v14v1, v14v2, t, "equal-map-v14")
		bs14 = testMarshalErr(&v14v1, h, t, "enc-map-v14-p")
		v14v2 = nil
		testUnmarshalErr(&v14v2, bs14, h, t, "dec-map-v14-p")
		testDeepEqualErr(v14v1, v14v2, t, "equal-map-v14-p")
	}

	for _, v := range []map[interface{}]float32{nil, map[interface{}]float32{}, map[interface{}]float32{"string-is-an-interface": 10.1}} {
		// fmt.Printf(">>>> running mammoth map v15: %v\n", v)
		var v15v1, v15v2 map[interface{}]float32
		v15v1 = v
		bs15 := testMarshalErr(v15v1, h, t, "enc-map-v15")
		if v != nil {
			v15v2 = make(map[interface{}]float32, len(v))
		}
		testUnmarshalErr(v15v2, bs15, h, t, "dec-map-v15")
		testDeepEqualErr(v15v1, v15v2, t, "equal-map-v15")
		bs15 = testMarshalErr(&v15v1, h, t, "enc-map-v15-p")
		v15v2 = nil
		testUnmarshalErr(&v15v2, bs15, h, t, "dec-map-v15-p")
		testDeepEqualErr(v15v1, v15v2, t, "equal-map-v15-p")
	}

	for _, v := range []map[interface{}]float64{nil, map[interface{}]float64{}, map[interface{}]float64{"string-is-an-interface": 10.1}} {
		// fmt.Printf(">>>> running mammoth map v16: %v\n", v)
		var v16v1, v16v2 map[interface{}]float64
		v16v1 = v
		bs16 := testMarshalErr(v16v1, h, t, "enc-map-v16")
		if v != nil {
			v16v2 = make(map[interface{}]float64, len(v))
		}
		testUnmarshalErr(v16v2, bs16, h, t, "dec-map-v16")
		testDeepEqualErr(v16v1, v16v2, t, "equal-map-v16")
		bs16 = testMarshalErr(&v16v1, h, t, "enc-map-v16-p")
		v16v2 = nil
		testUnmarshalErr(&v16v2, bs16, h, t, "dec-map-v16-p")
		testDeepEqualErr(v16v1, v16v2, t, "equal-map-v16-p")
	}

	for _, v := range []map[interface{}]bool{nil, map[interface{}]bool{}, map[interface{}]bool{"string-is-an-interface": true}} {
		// fmt.Printf(">>>> running mammoth map v17: %v\n", v)
		var v17v1, v17v2 map[interface{}]bool
		v17v1 = v
		bs17 := testMarshalErr(v17v1, h, t, "enc-map-v17")
		if v != nil {
			v17v2 = make(map[interface{}]bool, len(v))
		}
		testUnmarshalErr(v17v2, bs17, h, t, "dec-map-v17")
		testDeepEqualErr(v17v1, v17v2, t, "equal-map-v17")
		bs17 = testMarshalErr(&v17v1, h, t, "enc-map-v17-p")
		v17v2 = nil
		testUnmarshalErr(&v17v2, bs17, h, t, "dec-map-v17-p")
		testDeepEqualErr(v17v1, v17v2, t, "equal-map-v17-p")
	}

	for _, v := range []map[string]interface{}{nil, map[string]interface{}{}, map[string]interface{}{"some-string": "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v20: %v\n", v)
		var v20v1, v20v2 map[string]interface{}
		v20v1 = v
		bs20 := testMarshalErr(v20v1, h, t, "enc-map-v20")
		if v != nil {
			v20v2 = make(map[string]interface{}, len(v))
		}
		testUnmarshalErr(v20v2, bs20, h, t, "dec-map-v20")
		testDeepEqualErr(v20v1, v20v2, t, "equal-map-v20")
		bs20 = testMarshalErr(&v20v1, h, t, "enc-map-v20-p")
		v20v2 = nil
		testUnmarshalErr(&v20v2, bs20, h, t, "dec-map-v20-p")
		testDeepEqualErr(v20v1, v20v2, t, "equal-map-v20-p")
	}

	for _, v := range []map[string]string{nil, map[string]string{}, map[string]string{"some-string": "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v21: %v\n", v)
		var v21v1, v21v2 map[string]string
		v21v1 = v
		bs21 := testMarshalErr(v21v1, h, t, "enc-map-v21")
		if v != nil {
			v21v2 = make(map[string]string, len(v))
		}
		testUnmarshalErr(v21v2, bs21, h, t, "dec-map-v21")
		testDeepEqualErr(v21v1, v21v2, t, "equal-map-v21")
		bs21 = testMarshalErr(&v21v1, h, t, "enc-map-v21-p")
		v21v2 = nil
		testUnmarshalErr(&v21v2, bs21, h, t, "dec-map-v21-p")
		testDeepEqualErr(v21v1, v21v2, t, "equal-map-v21-p")
	}

	for _, v := range []map[string]uint{nil, map[string]uint{}, map[string]uint{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v22: %v\n", v)
		var v22v1, v22v2 map[string]uint
		v22v1 = v
		bs22 := testMarshalErr(v22v1, h, t, "enc-map-v22")
		if v != nil {
			v22v2 = make(map[string]uint, len(v))
		}
		testUnmarshalErr(v22v2, bs22, h, t, "dec-map-v22")
		testDeepEqualErr(v22v1, v22v2, t, "equal-map-v22")
		bs22 = testMarshalErr(&v22v1, h, t, "enc-map-v22-p")
		v22v2 = nil
		testUnmarshalErr(&v22v2, bs22, h, t, "dec-map-v22-p")
		testDeepEqualErr(v22v1, v22v2, t, "equal-map-v22-p")
	}

	for _, v := range []map[string]uint8{nil, map[string]uint8{}, map[string]uint8{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v23: %v\n", v)
		var v23v1, v23v2 map[string]uint8
		v23v1 = v
		bs23 := testMarshalErr(v23v1, h, t, "enc-map-v23")
		if v != nil {
			v23v2 = make(map[string]uint8, len(v))
		}
		testUnmarshalErr(v23v2, bs23, h, t, "dec-map-v23")
		testDeepEqualErr(v23v1, v23v2, t, "equal-map-v23")
		bs23 = testMarshalErr(&v23v1, h, t, "enc-map-v23-p")
		v23v2 = nil
		testUnmarshalErr(&v23v2, bs23, h, t, "dec-map-v23-p")
		testDeepEqualErr(v23v1, v23v2, t, "equal-map-v23-p")
	}

	for _, v := range []map[string]uint16{nil, map[string]uint16{}, map[string]uint16{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v24: %v\n", v)
		var v24v1, v24v2 map[string]uint16
		v24v1 = v
		bs24 := testMarshalErr(v24v1, h, t, "enc-map-v24")
		if v != nil {
			v24v2 = make(map[string]uint16, len(v))
		}
		testUnmarshalErr(v24v2, bs24, h, t, "dec-map-v24")
		testDeepEqualErr(v24v1, v24v2, t, "equal-map-v24")
		bs24 = testMarshalErr(&v24v1, h, t, "enc-map-v24-p")
		v24v2 = nil
		testUnmarshalErr(&v24v2, bs24, h, t, "dec-map-v24-p")
		testDeepEqualErr(v24v1, v24v2, t, "equal-map-v24-p")
	}

	for _, v := range []map[string]uint32{nil, map[string]uint32{}, map[string]uint32{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v25: %v\n", v)
		var v25v1, v25v2 map[string]uint32
		v25v1 = v
		bs25 := testMarshalErr(v25v1, h, t, "enc-map-v25")
		if v != nil {
			v25v2 = make(map[string]uint32, len(v))
		}
		testUnmarshalErr(v25v2, bs25, h, t, "dec-map-v25")
		testDeepEqualErr(v25v1, v25v2, t, "equal-map-v25")
		bs25 = testMarshalErr(&v25v1, h, t, "enc-map-v25-p")
		v25v2 = nil
		testUnmarshalErr(&v25v2, bs25, h, t, "dec-map-v25-p")
		testDeepEqualErr(v25v1, v25v2, t, "equal-map-v25-p")
	}

	for _, v := range []map[string]uint64{nil, map[string]uint64{}, map[string]uint64{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v26: %v\n", v)
		var v26v1, v26v2 map[string]uint64
		v26v1 = v
		bs26 := testMarshalErr(v26v1, h, t, "enc-map-v26")
		if v != nil {
			v26v2 = make(map[string]uint64, len(v))
		}
		testUnmarshalErr(v26v2, bs26, h, t, "dec-map-v26")
		testDeepEqualErr(v26v1, v26v2, t, "equal-map-v26")
		bs26 = testMarshalErr(&v26v1, h, t, "enc-map-v26-p")
		v26v2 = nil
		testUnmarshalErr(&v26v2, bs26, h, t, "dec-map-v26-p")
		testDeepEqualErr(v26v1, v26v2, t, "equal-map-v26-p")
	}

	for _, v := range []map[string]uintptr{nil, map[string]uintptr{}, map[string]uintptr{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v27: %v\n", v)
		var v27v1, v27v2 map[string]uintptr
		v27v1 = v
		bs27 := testMarshalErr(v27v1, h, t, "enc-map-v27")
		if v != nil {
			v27v2 = make(map[string]uintptr, len(v))
		}
		testUnmarshalErr(v27v2, bs27, h, t, "dec-map-v27")
		testDeepEqualErr(v27v1, v27v2, t, "equal-map-v27")
		bs27 = testMarshalErr(&v27v1, h, t, "enc-map-v27-p")
		v27v2 = nil
		testUnmarshalErr(&v27v2, bs27, h, t, "dec-map-v27-p")
		testDeepEqualErr(v27v1, v27v2, t, "equal-map-v27-p")
	}

	for _, v := range []map[string]int{nil, map[string]int{}, map[string]int{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v28: %v\n", v)
		var v28v1, v28v2 map[string]int
		v28v1 = v
		bs28 := testMarshalErr(v28v1, h, t, "enc-map-v28")
		if v != nil {
			v28v2 = make(map[string]int, len(v))
		}
		testUnmarshalErr(v28v2, bs28, h, t, "dec-map-v28")
		testDeepEqualErr(v28v1, v28v2, t, "equal-map-v28")
		bs28 = testMarshalErr(&v28v1, h, t, "enc-map-v28-p")
		v28v2 = nil
		testUnmarshalErr(&v28v2, bs28, h, t, "dec-map-v28-p")
		testDeepEqualErr(v28v1, v28v2, t, "equal-map-v28-p")
	}

	for _, v := range []map[string]int8{nil, map[string]int8{}, map[string]int8{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v29: %v\n", v)
		var v29v1, v29v2 map[string]int8
		v29v1 = v
		bs29 := testMarshalErr(v29v1, h, t, "enc-map-v29")
		if v != nil {
			v29v2 = make(map[string]int8, len(v))
		}
		testUnmarshalErr(v29v2, bs29, h, t, "dec-map-v29")
		testDeepEqualErr(v29v1, v29v2, t, "equal-map-v29")
		bs29 = testMarshalErr(&v29v1, h, t, "enc-map-v29-p")
		v29v2 = nil
		testUnmarshalErr(&v29v2, bs29, h, t, "dec-map-v29-p")
		testDeepEqualErr(v29v1, v29v2, t, "equal-map-v29-p")
	}

	for _, v := range []map[string]int16{nil, map[string]int16{}, map[string]int16{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v30: %v\n", v)
		var v30v1, v30v2 map[string]int16
		v30v1 = v
		bs30 := testMarshalErr(v30v1, h, t, "enc-map-v30")
		if v != nil {
			v30v2 = make(map[string]int16, len(v))
		}
		testUnmarshalErr(v30v2, bs30, h, t, "dec-map-v30")
		testDeepEqualErr(v30v1, v30v2, t, "equal-map-v30")
		bs30 = testMarshalErr(&v30v1, h, t, "enc-map-v30-p")
		v30v2 = nil
		testUnmarshalErr(&v30v2, bs30, h, t, "dec-map-v30-p")
		testDeepEqualErr(v30v1, v30v2, t, "equal-map-v30-p")
	}

	for _, v := range []map[string]int32{nil, map[string]int32{}, map[string]int32{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v31: %v\n", v)
		var v31v1, v31v2 map[string]int32
		v31v1 = v
		bs31 := testMarshalErr(v31v1, h, t, "enc-map-v31")
		if v != nil {
			v31v2 = make(map[string]int32, len(v))
		}
		testUnmarshalErr(v31v2, bs31, h, t, "dec-map-v31")
		testDeepEqualErr(v31v1, v31v2, t, "equal-map-v31")
		bs31 = testMarshalErr(&v31v1, h, t, "enc-map-v31-p")
		v31v2 = nil
		testUnmarshalErr(&v31v2, bs31, h, t, "dec-map-v31-p")
		testDeepEqualErr(v31v1, v31v2, t, "equal-map-v31-p")
	}

	for _, v := range []map[string]int64{nil, map[string]int64{}, map[string]int64{"some-string": 10}} {
		// fmt.Printf(">>>> running mammoth map v32: %v\n", v)
		var v32v1, v32v2 map[string]int64
		v32v1 = v
		bs32 := testMarshalErr(v32v1, h, t, "enc-map-v32")
		if v != nil {
			v32v2 = make(map[string]int64, len(v))
		}
		testUnmarshalErr(v32v2, bs32, h, t, "dec-map-v32")
		testDeepEqualErr(v32v1, v32v2, t, "equal-map-v32")
		bs32 = testMarshalErr(&v32v1, h, t, "enc-map-v32-p")
		v32v2 = nil
		testUnmarshalErr(&v32v2, bs32, h, t, "dec-map-v32-p")
		testDeepEqualErr(v32v1, v32v2, t, "equal-map-v32-p")
	}

	for _, v := range []map[string]float32{nil, map[string]float32{}, map[string]float32{"some-string": 10.1}} {
		// fmt.Printf(">>>> running mammoth map v33: %v\n", v)
		var v33v1, v33v2 map[string]float32
		v33v1 = v
		bs33 := testMarshalErr(v33v1, h, t, "enc-map-v33")
		if v != nil {
			v33v2 = make(map[string]float32, len(v))
		}
		testUnmarshalErr(v33v2, bs33, h, t, "dec-map-v33")
		testDeepEqualErr(v33v1, v33v2, t, "equal-map-v33")
		bs33 = testMarshalErr(&v33v1, h, t, "enc-map-v33-p")
		v33v2 = nil
		testUnmarshalErr(&v33v2, bs33, h, t, "dec-map-v33-p")
		testDeepEqualErr(v33v1, v33v2, t, "equal-map-v33-p")
	}

	for _, v := range []map[string]float64{nil, map[string]float64{}, map[string]float64{"some-string": 10.1}} {
		// fmt.Printf(">>>> running mammoth map v34: %v\n", v)
		var v34v1, v34v2 map[string]float64
		v34v1 = v
		bs34 := testMarshalErr(v34v1, h, t, "enc-map-v34")
		if v != nil {
			v34v2 = make(map[string]float64, len(v))
		}
		testUnmarshalErr(v34v2, bs34, h, t, "dec-map-v34")
		testDeepEqualErr(v34v1, v34v2, t, "equal-map-v34")
		bs34 = testMarshalErr(&v34v1, h, t, "enc-map-v34-p")
		v34v2 = nil
		testUnmarshalErr(&v34v2, bs34, h, t, "dec-map-v34-p")
		testDeepEqualErr(v34v1, v34v2, t, "equal-map-v34-p")
	}

	for _, v := range []map[string]bool{nil, map[string]bool{}, map[string]bool{"some-string": true}} {
		// fmt.Printf(">>>> running mammoth map v35: %v\n", v)
		var v35v1, v35v2 map[string]bool
		v35v1 = v
		bs35 := testMarshalErr(v35v1, h, t, "enc-map-v35")
		if v != nil {
			v35v2 = make(map[string]bool, len(v))
		}
		testUnmarshalErr(v35v2, bs35, h, t, "dec-map-v35")
		testDeepEqualErr(v35v1, v35v2, t, "equal-map-v35")
		bs35 = testMarshalErr(&v35v1, h, t, "enc-map-v35-p")
		v35v2 = nil
		testUnmarshalErr(&v35v2, bs35, h, t, "dec-map-v35-p")
		testDeepEqualErr(v35v1, v35v2, t, "equal-map-v35-p")
	}

	for _, v := range []map[float32]interface{}{nil, map[float32]interface{}{}, map[float32]interface{}{10.1: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v38: %v\n", v)
		var v38v1, v38v2 map[float32]interface{}
		v38v1 = v
		bs38 := testMarshalErr(v38v1, h, t, "enc-map-v38")
		if v != nil {
			v38v2 = make(map[float32]interface{}, len(v))
		}
		testUnmarshalErr(v38v2, bs38, h, t, "dec-map-v38")
		testDeepEqualErr(v38v1, v38v2, t, "equal-map-v38")
		bs38 = testMarshalErr(&v38v1, h, t, "enc-map-v38-p")
		v38v2 = nil
		testUnmarshalErr(&v38v2, bs38, h, t, "dec-map-v38-p")
		testDeepEqualErr(v38v1, v38v2, t, "equal-map-v38-p")
	}

	for _, v := range []map[float32]string{nil, map[float32]string{}, map[float32]string{10.1: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v39: %v\n", v)
		var v39v1, v39v2 map[float32]string
		v39v1 = v
		bs39 := testMarshalErr(v39v1, h, t, "enc-map-v39")
		if v != nil {
			v39v2 = make(map[float32]string, len(v))
		}
		testUnmarshalErr(v39v2, bs39, h, t, "dec-map-v39")
		testDeepEqualErr(v39v1, v39v2, t, "equal-map-v39")
		bs39 = testMarshalErr(&v39v1, h, t, "enc-map-v39-p")
		v39v2 = nil
		testUnmarshalErr(&v39v2, bs39, h, t, "dec-map-v39-p")
		testDeepEqualErr(v39v1, v39v2, t, "equal-map-v39-p")
	}

	for _, v := range []map[float32]uint{nil, map[float32]uint{}, map[float32]uint{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v40: %v\n", v)
		var v40v1, v40v2 map[float32]uint
		v40v1 = v
		bs40 := testMarshalErr(v40v1, h, t, "enc-map-v40")
		if v != nil {
			v40v2 = make(map[float32]uint, len(v))
		}
		testUnmarshalErr(v40v2, bs40, h, t, "dec-map-v40")
		testDeepEqualErr(v40v1, v40v2, t, "equal-map-v40")
		bs40 = testMarshalErr(&v40v1, h, t, "enc-map-v40-p")
		v40v2 = nil
		testUnmarshalErr(&v40v2, bs40, h, t, "dec-map-v40-p")
		testDeepEqualErr(v40v1, v40v2, t, "equal-map-v40-p")
	}

	for _, v := range []map[float32]uint8{nil, map[float32]uint8{}, map[float32]uint8{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v41: %v\n", v)
		var v41v1, v41v2 map[float32]uint8
		v41v1 = v
		bs41 := testMarshalErr(v41v1, h, t, "enc-map-v41")
		if v != nil {
			v41v2 = make(map[float32]uint8, len(v))
		}
		testUnmarshalErr(v41v2, bs41, h, t, "dec-map-v41")
		testDeepEqualErr(v41v1, v41v2, t, "equal-map-v41")
		bs41 = testMarshalErr(&v41v1, h, t, "enc-map-v41-p")
		v41v2 = nil
		testUnmarshalErr(&v41v2, bs41, h, t, "dec-map-v41-p")
		testDeepEqualErr(v41v1, v41v2, t, "equal-map-v41-p")
	}

	for _, v := range []map[float32]uint16{nil, map[float32]uint16{}, map[float32]uint16{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v42: %v\n", v)
		var v42v1, v42v2 map[float32]uint16
		v42v1 = v
		bs42 := testMarshalErr(v42v1, h, t, "enc-map-v42")
		if v != nil {
			v42v2 = make(map[float32]uint16, len(v))
		}
		testUnmarshalErr(v42v2, bs42, h, t, "dec-map-v42")
		testDeepEqualErr(v42v1, v42v2, t, "equal-map-v42")
		bs42 = testMarshalErr(&v42v1, h, t, "enc-map-v42-p")
		v42v2 = nil
		testUnmarshalErr(&v42v2, bs42, h, t, "dec-map-v42-p")
		testDeepEqualErr(v42v1, v42v2, t, "equal-map-v42-p")
	}

	for _, v := range []map[float32]uint32{nil, map[float32]uint32{}, map[float32]uint32{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v43: %v\n", v)
		var v43v1, v43v2 map[float32]uint32
		v43v1 = v
		bs43 := testMarshalErr(v43v1, h, t, "enc-map-v43")
		if v != nil {
			v43v2 = make(map[float32]uint32, len(v))
		}
		testUnmarshalErr(v43v2, bs43, h, t, "dec-map-v43")
		testDeepEqualErr(v43v1, v43v2, t, "equal-map-v43")
		bs43 = testMarshalErr(&v43v1, h, t, "enc-map-v43-p")
		v43v2 = nil
		testUnmarshalErr(&v43v2, bs43, h, t, "dec-map-v43-p")
		testDeepEqualErr(v43v1, v43v2, t, "equal-map-v43-p")
	}

	for _, v := range []map[float32]uint64{nil, map[float32]uint64{}, map[float32]uint64{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v44: %v\n", v)
		var v44v1, v44v2 map[float32]uint64
		v44v1 = v
		bs44 := testMarshalErr(v44v1, h, t, "enc-map-v44")
		if v != nil {
			v44v2 = make(map[float32]uint64, len(v))
		}
		testUnmarshalErr(v44v2, bs44, h, t, "dec-map-v44")
		testDeepEqualErr(v44v1, v44v2, t, "equal-map-v44")
		bs44 = testMarshalErr(&v44v1, h, t, "enc-map-v44-p")
		v44v2 = nil
		testUnmarshalErr(&v44v2, bs44, h, t, "dec-map-v44-p")
		testDeepEqualErr(v44v1, v44v2, t, "equal-map-v44-p")
	}

	for _, v := range []map[float32]uintptr{nil, map[float32]uintptr{}, map[float32]uintptr{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v45: %v\n", v)
		var v45v1, v45v2 map[float32]uintptr
		v45v1 = v
		bs45 := testMarshalErr(v45v1, h, t, "enc-map-v45")
		if v != nil {
			v45v2 = make(map[float32]uintptr, len(v))
		}
		testUnmarshalErr(v45v2, bs45, h, t, "dec-map-v45")
		testDeepEqualErr(v45v1, v45v2, t, "equal-map-v45")
		bs45 = testMarshalErr(&v45v1, h, t, "enc-map-v45-p")
		v45v2 = nil
		testUnmarshalErr(&v45v2, bs45, h, t, "dec-map-v45-p")
		testDeepEqualErr(v45v1, v45v2, t, "equal-map-v45-p")
	}

	for _, v := range []map[float32]int{nil, map[float32]int{}, map[float32]int{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v46: %v\n", v)
		var v46v1, v46v2 map[float32]int
		v46v1 = v
		bs46 := testMarshalErr(v46v1, h, t, "enc-map-v46")
		if v != nil {
			v46v2 = make(map[float32]int, len(v))
		}
		testUnmarshalErr(v46v2, bs46, h, t, "dec-map-v46")
		testDeepEqualErr(v46v1, v46v2, t, "equal-map-v46")
		bs46 = testMarshalErr(&v46v1, h, t, "enc-map-v46-p")
		v46v2 = nil
		testUnmarshalErr(&v46v2, bs46, h, t, "dec-map-v46-p")
		testDeepEqualErr(v46v1, v46v2, t, "equal-map-v46-p")
	}

	for _, v := range []map[float32]int8{nil, map[float32]int8{}, map[float32]int8{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v47: %v\n", v)
		var v47v1, v47v2 map[float32]int8
		v47v1 = v
		bs47 := testMarshalErr(v47v1, h, t, "enc-map-v47")
		if v != nil {
			v47v2 = make(map[float32]int8, len(v))
		}
		testUnmarshalErr(v47v2, bs47, h, t, "dec-map-v47")
		testDeepEqualErr(v47v1, v47v2, t, "equal-map-v47")
		bs47 = testMarshalErr(&v47v1, h, t, "enc-map-v47-p")
		v47v2 = nil
		testUnmarshalErr(&v47v2, bs47, h, t, "dec-map-v47-p")
		testDeepEqualErr(v47v1, v47v2, t, "equal-map-v47-p")
	}

	for _, v := range []map[float32]int16{nil, map[float32]int16{}, map[float32]int16{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v48: %v\n", v)
		var v48v1, v48v2 map[float32]int16
		v48v1 = v
		bs48 := testMarshalErr(v48v1, h, t, "enc-map-v48")
		if v != nil {
			v48v2 = make(map[float32]int16, len(v))
		}
		testUnmarshalErr(v48v2, bs48, h, t, "dec-map-v48")
		testDeepEqualErr(v48v1, v48v2, t, "equal-map-v48")
		bs48 = testMarshalErr(&v48v1, h, t, "enc-map-v48-p")
		v48v2 = nil
		testUnmarshalErr(&v48v2, bs48, h, t, "dec-map-v48-p")
		testDeepEqualErr(v48v1, v48v2, t, "equal-map-v48-p")
	}

	for _, v := range []map[float32]int32{nil, map[float32]int32{}, map[float32]int32{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v49: %v\n", v)
		var v49v1, v49v2 map[float32]int32
		v49v1 = v
		bs49 := testMarshalErr(v49v1, h, t, "enc-map-v49")
		if v != nil {
			v49v2 = make(map[float32]int32, len(v))
		}
		testUnmarshalErr(v49v2, bs49, h, t, "dec-map-v49")
		testDeepEqualErr(v49v1, v49v2, t, "equal-map-v49")
		bs49 = testMarshalErr(&v49v1, h, t, "enc-map-v49-p")
		v49v2 = nil
		testUnmarshalErr(&v49v2, bs49, h, t, "dec-map-v49-p")
		testDeepEqualErr(v49v1, v49v2, t, "equal-map-v49-p")
	}

	for _, v := range []map[float32]int64{nil, map[float32]int64{}, map[float32]int64{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v50: %v\n", v)
		var v50v1, v50v2 map[float32]int64
		v50v1 = v
		bs50 := testMarshalErr(v50v1, h, t, "enc-map-v50")
		if v != nil {
			v50v2 = make(map[float32]int64, len(v))
		}
		testUnmarshalErr(v50v2, bs50, h, t, "dec-map-v50")
		testDeepEqualErr(v50v1, v50v2, t, "equal-map-v50")
		bs50 = testMarshalErr(&v50v1, h, t, "enc-map-v50-p")
		v50v2 = nil
		testUnmarshalErr(&v50v2, bs50, h, t, "dec-map-v50-p")
		testDeepEqualErr(v50v1, v50v2, t, "equal-map-v50-p")
	}

	for _, v := range []map[float32]float32{nil, map[float32]float32{}, map[float32]float32{10.1: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v51: %v\n", v)
		var v51v1, v51v2 map[float32]float32
		v51v1 = v
		bs51 := testMarshalErr(v51v1, h, t, "enc-map-v51")
		if v != nil {
			v51v2 = make(map[float32]float32, len(v))
		}
		testUnmarshalErr(v51v2, bs51, h, t, "dec-map-v51")
		testDeepEqualErr(v51v1, v51v2, t, "equal-map-v51")
		bs51 = testMarshalErr(&v51v1, h, t, "enc-map-v51-p")
		v51v2 = nil
		testUnmarshalErr(&v51v2, bs51, h, t, "dec-map-v51-p")
		testDeepEqualErr(v51v1, v51v2, t, "equal-map-v51-p")
	}

	for _, v := range []map[float32]float64{nil, map[float32]float64{}, map[float32]float64{10.1: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v52: %v\n", v)
		var v52v1, v52v2 map[float32]float64
		v52v1 = v
		bs52 := testMarshalErr(v52v1, h, t, "enc-map-v52")
		if v != nil {
			v52v2 = make(map[float32]float64, len(v))
		}
		testUnmarshalErr(v52v2, bs52, h, t, "dec-map-v52")
		testDeepEqualErr(v52v1, v52v2, t, "equal-map-v52")
		bs52 = testMarshalErr(&v52v1, h, t, "enc-map-v52-p")
		v52v2 = nil
		testUnmarshalErr(&v52v2, bs52, h, t, "dec-map-v52-p")
		testDeepEqualErr(v52v1, v52v2, t, "equal-map-v52-p")
	}

	for _, v := range []map[float32]bool{nil, map[float32]bool{}, map[float32]bool{10.1: true}} {
		// fmt.Printf(">>>> running mammoth map v53: %v\n", v)
		var v53v1, v53v2 map[float32]bool
		v53v1 = v
		bs53 := testMarshalErr(v53v1, h, t, "enc-map-v53")
		if v != nil {
			v53v2 = make(map[float32]bool, len(v))
		}
		testUnmarshalErr(v53v2, bs53, h, t, "dec-map-v53")
		testDeepEqualErr(v53v1, v53v2, t, "equal-map-v53")
		bs53 = testMarshalErr(&v53v1, h, t, "enc-map-v53-p")
		v53v2 = nil
		testUnmarshalErr(&v53v2, bs53, h, t, "dec-map-v53-p")
		testDeepEqualErr(v53v1, v53v2, t, "equal-map-v53-p")
	}

	for _, v := range []map[float64]interface{}{nil, map[float64]interface{}{}, map[float64]interface{}{10.1: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v56: %v\n", v)
		var v56v1, v56v2 map[float64]interface{}
		v56v1 = v
		bs56 := testMarshalErr(v56v1, h, t, "enc-map-v56")
		if v != nil {
			v56v2 = make(map[float64]interface{}, len(v))
		}
		testUnmarshalErr(v56v2, bs56, h, t, "dec-map-v56")
		testDeepEqualErr(v56v1, v56v2, t, "equal-map-v56")
		bs56 = testMarshalErr(&v56v1, h, t, "enc-map-v56-p")
		v56v2 = nil
		testUnmarshalErr(&v56v2, bs56, h, t, "dec-map-v56-p")
		testDeepEqualErr(v56v1, v56v2, t, "equal-map-v56-p")
	}

	for _, v := range []map[float64]string{nil, map[float64]string{}, map[float64]string{10.1: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v57: %v\n", v)
		var v57v1, v57v2 map[float64]string
		v57v1 = v
		bs57 := testMarshalErr(v57v1, h, t, "enc-map-v57")
		if v != nil {
			v57v2 = make(map[float64]string, len(v))
		}
		testUnmarshalErr(v57v2, bs57, h, t, "dec-map-v57")
		testDeepEqualErr(v57v1, v57v2, t, "equal-map-v57")
		bs57 = testMarshalErr(&v57v1, h, t, "enc-map-v57-p")
		v57v2 = nil
		testUnmarshalErr(&v57v2, bs57, h, t, "dec-map-v57-p")
		testDeepEqualErr(v57v1, v57v2, t, "equal-map-v57-p")
	}

	for _, v := range []map[float64]uint{nil, map[float64]uint{}, map[float64]uint{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v58: %v\n", v)
		var v58v1, v58v2 map[float64]uint
		v58v1 = v
		bs58 := testMarshalErr(v58v1, h, t, "enc-map-v58")
		if v != nil {
			v58v2 = make(map[float64]uint, len(v))
		}
		testUnmarshalErr(v58v2, bs58, h, t, "dec-map-v58")
		testDeepEqualErr(v58v1, v58v2, t, "equal-map-v58")
		bs58 = testMarshalErr(&v58v1, h, t, "enc-map-v58-p")
		v58v2 = nil
		testUnmarshalErr(&v58v2, bs58, h, t, "dec-map-v58-p")
		testDeepEqualErr(v58v1, v58v2, t, "equal-map-v58-p")
	}

	for _, v := range []map[float64]uint8{nil, map[float64]uint8{}, map[float64]uint8{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v59: %v\n", v)
		var v59v1, v59v2 map[float64]uint8
		v59v1 = v
		bs59 := testMarshalErr(v59v1, h, t, "enc-map-v59")
		if v != nil {
			v59v2 = make(map[float64]uint8, len(v))
		}
		testUnmarshalErr(v59v2, bs59, h, t, "dec-map-v59")
		testDeepEqualErr(v59v1, v59v2, t, "equal-map-v59")
		bs59 = testMarshalErr(&v59v1, h, t, "enc-map-v59-p")
		v59v2 = nil
		testUnmarshalErr(&v59v2, bs59, h, t, "dec-map-v59-p")
		testDeepEqualErr(v59v1, v59v2, t, "equal-map-v59-p")
	}

	for _, v := range []map[float64]uint16{nil, map[float64]uint16{}, map[float64]uint16{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v60: %v\n", v)
		var v60v1, v60v2 map[float64]uint16
		v60v1 = v
		bs60 := testMarshalErr(v60v1, h, t, "enc-map-v60")
		if v != nil {
			v60v2 = make(map[float64]uint16, len(v))
		}
		testUnmarshalErr(v60v2, bs60, h, t, "dec-map-v60")
		testDeepEqualErr(v60v1, v60v2, t, "equal-map-v60")
		bs60 = testMarshalErr(&v60v1, h, t, "enc-map-v60-p")
		v60v2 = nil
		testUnmarshalErr(&v60v2, bs60, h, t, "dec-map-v60-p")
		testDeepEqualErr(v60v1, v60v2, t, "equal-map-v60-p")
	}

	for _, v := range []map[float64]uint32{nil, map[float64]uint32{}, map[float64]uint32{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v61: %v\n", v)
		var v61v1, v61v2 map[float64]uint32
		v61v1 = v
		bs61 := testMarshalErr(v61v1, h, t, "enc-map-v61")
		if v != nil {
			v61v2 = make(map[float64]uint32, len(v))
		}
		testUnmarshalErr(v61v2, bs61, h, t, "dec-map-v61")
		testDeepEqualErr(v61v1, v61v2, t, "equal-map-v61")
		bs61 = testMarshalErr(&v61v1, h, t, "enc-map-v61-p")
		v61v2 = nil
		testUnmarshalErr(&v61v2, bs61, h, t, "dec-map-v61-p")
		testDeepEqualErr(v61v1, v61v2, t, "equal-map-v61-p")
	}

	for _, v := range []map[float64]uint64{nil, map[float64]uint64{}, map[float64]uint64{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v62: %v\n", v)
		var v62v1, v62v2 map[float64]uint64
		v62v1 = v
		bs62 := testMarshalErr(v62v1, h, t, "enc-map-v62")
		if v != nil {
			v62v2 = make(map[float64]uint64, len(v))
		}
		testUnmarshalErr(v62v2, bs62, h, t, "dec-map-v62")
		testDeepEqualErr(v62v1, v62v2, t, "equal-map-v62")
		bs62 = testMarshalErr(&v62v1, h, t, "enc-map-v62-p")
		v62v2 = nil
		testUnmarshalErr(&v62v2, bs62, h, t, "dec-map-v62-p")
		testDeepEqualErr(v62v1, v62v2, t, "equal-map-v62-p")
	}

	for _, v := range []map[float64]uintptr{nil, map[float64]uintptr{}, map[float64]uintptr{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v63: %v\n", v)
		var v63v1, v63v2 map[float64]uintptr
		v63v1 = v
		bs63 := testMarshalErr(v63v1, h, t, "enc-map-v63")
		if v != nil {
			v63v2 = make(map[float64]uintptr, len(v))
		}
		testUnmarshalErr(v63v2, bs63, h, t, "dec-map-v63")
		testDeepEqualErr(v63v1, v63v2, t, "equal-map-v63")
		bs63 = testMarshalErr(&v63v1, h, t, "enc-map-v63-p")
		v63v2 = nil
		testUnmarshalErr(&v63v2, bs63, h, t, "dec-map-v63-p")
		testDeepEqualErr(v63v1, v63v2, t, "equal-map-v63-p")
	}

	for _, v := range []map[float64]int{nil, map[float64]int{}, map[float64]int{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v64: %v\n", v)
		var v64v1, v64v2 map[float64]int
		v64v1 = v
		bs64 := testMarshalErr(v64v1, h, t, "enc-map-v64")
		if v != nil {
			v64v2 = make(map[float64]int, len(v))
		}
		testUnmarshalErr(v64v2, bs64, h, t, "dec-map-v64")
		testDeepEqualErr(v64v1, v64v2, t, "equal-map-v64")
		bs64 = testMarshalErr(&v64v1, h, t, "enc-map-v64-p")
		v64v2 = nil
		testUnmarshalErr(&v64v2, bs64, h, t, "dec-map-v64-p")
		testDeepEqualErr(v64v1, v64v2, t, "equal-map-v64-p")
	}

	for _, v := range []map[float64]int8{nil, map[float64]int8{}, map[float64]int8{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v65: %v\n", v)
		var v65v1, v65v2 map[float64]int8
		v65v1 = v
		bs65 := testMarshalErr(v65v1, h, t, "enc-map-v65")
		if v != nil {
			v65v2 = make(map[float64]int8, len(v))
		}
		testUnmarshalErr(v65v2, bs65, h, t, "dec-map-v65")
		testDeepEqualErr(v65v1, v65v2, t, "equal-map-v65")
		bs65 = testMarshalErr(&v65v1, h, t, "enc-map-v65-p")
		v65v2 = nil
		testUnmarshalErr(&v65v2, bs65, h, t, "dec-map-v65-p")
		testDeepEqualErr(v65v1, v65v2, t, "equal-map-v65-p")
	}

	for _, v := range []map[float64]int16{nil, map[float64]int16{}, map[float64]int16{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v66: %v\n", v)
		var v66v1, v66v2 map[float64]int16
		v66v1 = v
		bs66 := testMarshalErr(v66v1, h, t, "enc-map-v66")
		if v != nil {
			v66v2 = make(map[float64]int16, len(v))
		}
		testUnmarshalErr(v66v2, bs66, h, t, "dec-map-v66")
		testDeepEqualErr(v66v1, v66v2, t, "equal-map-v66")
		bs66 = testMarshalErr(&v66v1, h, t, "enc-map-v66-p")
		v66v2 = nil
		testUnmarshalErr(&v66v2, bs66, h, t, "dec-map-v66-p")
		testDeepEqualErr(v66v1, v66v2, t, "equal-map-v66-p")
	}

	for _, v := range []map[float64]int32{nil, map[float64]int32{}, map[float64]int32{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v67: %v\n", v)
		var v67v1, v67v2 map[float64]int32
		v67v1 = v
		bs67 := testMarshalErr(v67v1, h, t, "enc-map-v67")
		if v != nil {
			v67v2 = make(map[float64]int32, len(v))
		}
		testUnmarshalErr(v67v2, bs67, h, t, "dec-map-v67")
		testDeepEqualErr(v67v1, v67v2, t, "equal-map-v67")
		bs67 = testMarshalErr(&v67v1, h, t, "enc-map-v67-p")
		v67v2 = nil
		testUnmarshalErr(&v67v2, bs67, h, t, "dec-map-v67-p")
		testDeepEqualErr(v67v1, v67v2, t, "equal-map-v67-p")
	}

	for _, v := range []map[float64]int64{nil, map[float64]int64{}, map[float64]int64{10.1: 10}} {
		// fmt.Printf(">>>> running mammoth map v68: %v\n", v)
		var v68v1, v68v2 map[float64]int64
		v68v1 = v
		bs68 := testMarshalErr(v68v1, h, t, "enc-map-v68")
		if v != nil {
			v68v2 = make(map[float64]int64, len(v))
		}
		testUnmarshalErr(v68v2, bs68, h, t, "dec-map-v68")
		testDeepEqualErr(v68v1, v68v2, t, "equal-map-v68")
		bs68 = testMarshalErr(&v68v1, h, t, "enc-map-v68-p")
		v68v2 = nil
		testUnmarshalErr(&v68v2, bs68, h, t, "dec-map-v68-p")
		testDeepEqualErr(v68v1, v68v2, t, "equal-map-v68-p")
	}

	for _, v := range []map[float64]float32{nil, map[float64]float32{}, map[float64]float32{10.1: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v69: %v\n", v)
		var v69v1, v69v2 map[float64]float32
		v69v1 = v
		bs69 := testMarshalErr(v69v1, h, t, "enc-map-v69")
		if v != nil {
			v69v2 = make(map[float64]float32, len(v))
		}
		testUnmarshalErr(v69v2, bs69, h, t, "dec-map-v69")
		testDeepEqualErr(v69v1, v69v2, t, "equal-map-v69")
		bs69 = testMarshalErr(&v69v1, h, t, "enc-map-v69-p")
		v69v2 = nil
		testUnmarshalErr(&v69v2, bs69, h, t, "dec-map-v69-p")
		testDeepEqualErr(v69v1, v69v2, t, "equal-map-v69-p")
	}

	for _, v := range []map[float64]float64{nil, map[float64]float64{}, map[float64]float64{10.1: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v70: %v\n", v)
		var v70v1, v70v2 map[float64]float64
		v70v1 = v
		bs70 := testMarshalErr(v70v1, h, t, "enc-map-v70")
		if v != nil {
			v70v2 = make(map[float64]float64, len(v))
		}
		testUnmarshalErr(v70v2, bs70, h, t, "dec-map-v70")
		testDeepEqualErr(v70v1, v70v2, t, "equal-map-v70")
		bs70 = testMarshalErr(&v70v1, h, t, "enc-map-v70-p")
		v70v2 = nil
		testUnmarshalErr(&v70v2, bs70, h, t, "dec-map-v70-p")
		testDeepEqualErr(v70v1, v70v2, t, "equal-map-v70-p")
	}

	for _, v := range []map[float64]bool{nil, map[float64]bool{}, map[float64]bool{10.1: true}} {
		// fmt.Printf(">>>> running mammoth map v71: %v\n", v)
		var v71v1, v71v2 map[float64]bool
		v71v1 = v
		bs71 := testMarshalErr(v71v1, h, t, "enc-map-v71")
		if v != nil {
			v71v2 = make(map[float64]bool, len(v))
		}
		testUnmarshalErr(v71v2, bs71, h, t, "dec-map-v71")
		testDeepEqualErr(v71v1, v71v2, t, "equal-map-v71")
		bs71 = testMarshalErr(&v71v1, h, t, "enc-map-v71-p")
		v71v2 = nil
		testUnmarshalErr(&v71v2, bs71, h, t, "dec-map-v71-p")
		testDeepEqualErr(v71v1, v71v2, t, "equal-map-v71-p")
	}

	for _, v := range []map[uint]interface{}{nil, map[uint]interface{}{}, map[uint]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v74: %v\n", v)
		var v74v1, v74v2 map[uint]interface{}
		v74v1 = v
		bs74 := testMarshalErr(v74v1, h, t, "enc-map-v74")
		if v != nil {
			v74v2 = make(map[uint]interface{}, len(v))
		}
		testUnmarshalErr(v74v2, bs74, h, t, "dec-map-v74")
		testDeepEqualErr(v74v1, v74v2, t, "equal-map-v74")
		bs74 = testMarshalErr(&v74v1, h, t, "enc-map-v74-p")
		v74v2 = nil
		testUnmarshalErr(&v74v2, bs74, h, t, "dec-map-v74-p")
		testDeepEqualErr(v74v1, v74v2, t, "equal-map-v74-p")
	}

	for _, v := range []map[uint]string{nil, map[uint]string{}, map[uint]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v75: %v\n", v)
		var v75v1, v75v2 map[uint]string
		v75v1 = v
		bs75 := testMarshalErr(v75v1, h, t, "enc-map-v75")
		if v != nil {
			v75v2 = make(map[uint]string, len(v))
		}
		testUnmarshalErr(v75v2, bs75, h, t, "dec-map-v75")
		testDeepEqualErr(v75v1, v75v2, t, "equal-map-v75")
		bs75 = testMarshalErr(&v75v1, h, t, "enc-map-v75-p")
		v75v2 = nil
		testUnmarshalErr(&v75v2, bs75, h, t, "dec-map-v75-p")
		testDeepEqualErr(v75v1, v75v2, t, "equal-map-v75-p")
	}

	for _, v := range []map[uint]uint{nil, map[uint]uint{}, map[uint]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v76: %v\n", v)
		var v76v1, v76v2 map[uint]uint
		v76v1 = v
		bs76 := testMarshalErr(v76v1, h, t, "enc-map-v76")
		if v != nil {
			v76v2 = make(map[uint]uint, len(v))
		}
		testUnmarshalErr(v76v2, bs76, h, t, "dec-map-v76")
		testDeepEqualErr(v76v1, v76v2, t, "equal-map-v76")
		bs76 = testMarshalErr(&v76v1, h, t, "enc-map-v76-p")
		v76v2 = nil
		testUnmarshalErr(&v76v2, bs76, h, t, "dec-map-v76-p")
		testDeepEqualErr(v76v1, v76v2, t, "equal-map-v76-p")
	}

	for _, v := range []map[uint]uint8{nil, map[uint]uint8{}, map[uint]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v77: %v\n", v)
		var v77v1, v77v2 map[uint]uint8
		v77v1 = v
		bs77 := testMarshalErr(v77v1, h, t, "enc-map-v77")
		if v != nil {
			v77v2 = make(map[uint]uint8, len(v))
		}
		testUnmarshalErr(v77v2, bs77, h, t, "dec-map-v77")
		testDeepEqualErr(v77v1, v77v2, t, "equal-map-v77")
		bs77 = testMarshalErr(&v77v1, h, t, "enc-map-v77-p")
		v77v2 = nil
		testUnmarshalErr(&v77v2, bs77, h, t, "dec-map-v77-p")
		testDeepEqualErr(v77v1, v77v2, t, "equal-map-v77-p")
	}

	for _, v := range []map[uint]uint16{nil, map[uint]uint16{}, map[uint]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v78: %v\n", v)
		var v78v1, v78v2 map[uint]uint16
		v78v1 = v
		bs78 := testMarshalErr(v78v1, h, t, "enc-map-v78")
		if v != nil {
			v78v2 = make(map[uint]uint16, len(v))
		}
		testUnmarshalErr(v78v2, bs78, h, t, "dec-map-v78")
		testDeepEqualErr(v78v1, v78v2, t, "equal-map-v78")
		bs78 = testMarshalErr(&v78v1, h, t, "enc-map-v78-p")
		v78v2 = nil
		testUnmarshalErr(&v78v2, bs78, h, t, "dec-map-v78-p")
		testDeepEqualErr(v78v1, v78v2, t, "equal-map-v78-p")
	}

	for _, v := range []map[uint]uint32{nil, map[uint]uint32{}, map[uint]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v79: %v\n", v)
		var v79v1, v79v2 map[uint]uint32
		v79v1 = v
		bs79 := testMarshalErr(v79v1, h, t, "enc-map-v79")
		if v != nil {
			v79v2 = make(map[uint]uint32, len(v))
		}
		testUnmarshalErr(v79v2, bs79, h, t, "dec-map-v79")
		testDeepEqualErr(v79v1, v79v2, t, "equal-map-v79")
		bs79 = testMarshalErr(&v79v1, h, t, "enc-map-v79-p")
		v79v2 = nil
		testUnmarshalErr(&v79v2, bs79, h, t, "dec-map-v79-p")
		testDeepEqualErr(v79v1, v79v2, t, "equal-map-v79-p")
	}

	for _, v := range []map[uint]uint64{nil, map[uint]uint64{}, map[uint]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v80: %v\n", v)
		var v80v1, v80v2 map[uint]uint64
		v80v1 = v
		bs80 := testMarshalErr(v80v1, h, t, "enc-map-v80")
		if v != nil {
			v80v2 = make(map[uint]uint64, len(v))
		}
		testUnmarshalErr(v80v2, bs80, h, t, "dec-map-v80")
		testDeepEqualErr(v80v1, v80v2, t, "equal-map-v80")
		bs80 = testMarshalErr(&v80v1, h, t, "enc-map-v80-p")
		v80v2 = nil
		testUnmarshalErr(&v80v2, bs80, h, t, "dec-map-v80-p")
		testDeepEqualErr(v80v1, v80v2, t, "equal-map-v80-p")
	}

	for _, v := range []map[uint]uintptr{nil, map[uint]uintptr{}, map[uint]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v81: %v\n", v)
		var v81v1, v81v2 map[uint]uintptr
		v81v1 = v
		bs81 := testMarshalErr(v81v1, h, t, "enc-map-v81")
		if v != nil {
			v81v2 = make(map[uint]uintptr, len(v))
		}
		testUnmarshalErr(v81v2, bs81, h, t, "dec-map-v81")
		testDeepEqualErr(v81v1, v81v2, t, "equal-map-v81")
		bs81 = testMarshalErr(&v81v1, h, t, "enc-map-v81-p")
		v81v2 = nil
		testUnmarshalErr(&v81v2, bs81, h, t, "dec-map-v81-p")
		testDeepEqualErr(v81v1, v81v2, t, "equal-map-v81-p")
	}

	for _, v := range []map[uint]int{nil, map[uint]int{}, map[uint]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v82: %v\n", v)
		var v82v1, v82v2 map[uint]int
		v82v1 = v
		bs82 := testMarshalErr(v82v1, h, t, "enc-map-v82")
		if v != nil {
			v82v2 = make(map[uint]int, len(v))
		}
		testUnmarshalErr(v82v2, bs82, h, t, "dec-map-v82")
		testDeepEqualErr(v82v1, v82v2, t, "equal-map-v82")
		bs82 = testMarshalErr(&v82v1, h, t, "enc-map-v82-p")
		v82v2 = nil
		testUnmarshalErr(&v82v2, bs82, h, t, "dec-map-v82-p")
		testDeepEqualErr(v82v1, v82v2, t, "equal-map-v82-p")
	}

	for _, v := range []map[uint]int8{nil, map[uint]int8{}, map[uint]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v83: %v\n", v)
		var v83v1, v83v2 map[uint]int8
		v83v1 = v
		bs83 := testMarshalErr(v83v1, h, t, "enc-map-v83")
		if v != nil {
			v83v2 = make(map[uint]int8, len(v))
		}
		testUnmarshalErr(v83v2, bs83, h, t, "dec-map-v83")
		testDeepEqualErr(v83v1, v83v2, t, "equal-map-v83")
		bs83 = testMarshalErr(&v83v1, h, t, "enc-map-v83-p")
		v83v2 = nil
		testUnmarshalErr(&v83v2, bs83, h, t, "dec-map-v83-p")
		testDeepEqualErr(v83v1, v83v2, t, "equal-map-v83-p")
	}

	for _, v := range []map[uint]int16{nil, map[uint]int16{}, map[uint]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v84: %v\n", v)
		var v84v1, v84v2 map[uint]int16
		v84v1 = v
		bs84 := testMarshalErr(v84v1, h, t, "enc-map-v84")
		if v != nil {
			v84v2 = make(map[uint]int16, len(v))
		}
		testUnmarshalErr(v84v2, bs84, h, t, "dec-map-v84")
		testDeepEqualErr(v84v1, v84v2, t, "equal-map-v84")
		bs84 = testMarshalErr(&v84v1, h, t, "enc-map-v84-p")
		v84v2 = nil
		testUnmarshalErr(&v84v2, bs84, h, t, "dec-map-v84-p")
		testDeepEqualErr(v84v1, v84v2, t, "equal-map-v84-p")
	}

	for _, v := range []map[uint]int32{nil, map[uint]int32{}, map[uint]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v85: %v\n", v)
		var v85v1, v85v2 map[uint]int32
		v85v1 = v
		bs85 := testMarshalErr(v85v1, h, t, "enc-map-v85")
		if v != nil {
			v85v2 = make(map[uint]int32, len(v))
		}
		testUnmarshalErr(v85v2, bs85, h, t, "dec-map-v85")
		testDeepEqualErr(v85v1, v85v2, t, "equal-map-v85")
		bs85 = testMarshalErr(&v85v1, h, t, "enc-map-v85-p")
		v85v2 = nil
		testUnmarshalErr(&v85v2, bs85, h, t, "dec-map-v85-p")
		testDeepEqualErr(v85v1, v85v2, t, "equal-map-v85-p")
	}

	for _, v := range []map[uint]int64{nil, map[uint]int64{}, map[uint]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v86: %v\n", v)
		var v86v1, v86v2 map[uint]int64
		v86v1 = v
		bs86 := testMarshalErr(v86v1, h, t, "enc-map-v86")
		if v != nil {
			v86v2 = make(map[uint]int64, len(v))
		}
		testUnmarshalErr(v86v2, bs86, h, t, "dec-map-v86")
		testDeepEqualErr(v86v1, v86v2, t, "equal-map-v86")
		bs86 = testMarshalErr(&v86v1, h, t, "enc-map-v86-p")
		v86v2 = nil
		testUnmarshalErr(&v86v2, bs86, h, t, "dec-map-v86-p")
		testDeepEqualErr(v86v1, v86v2, t, "equal-map-v86-p")
	}

	for _, v := range []map[uint]float32{nil, map[uint]float32{}, map[uint]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v87: %v\n", v)
		var v87v1, v87v2 map[uint]float32
		v87v1 = v
		bs87 := testMarshalErr(v87v1, h, t, "enc-map-v87")
		if v != nil {
			v87v2 = make(map[uint]float32, len(v))
		}
		testUnmarshalErr(v87v2, bs87, h, t, "dec-map-v87")
		testDeepEqualErr(v87v1, v87v2, t, "equal-map-v87")
		bs87 = testMarshalErr(&v87v1, h, t, "enc-map-v87-p")
		v87v2 = nil
		testUnmarshalErr(&v87v2, bs87, h, t, "dec-map-v87-p")
		testDeepEqualErr(v87v1, v87v2, t, "equal-map-v87-p")
	}

	for _, v := range []map[uint]float64{nil, map[uint]float64{}, map[uint]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v88: %v\n", v)
		var v88v1, v88v2 map[uint]float64
		v88v1 = v
		bs88 := testMarshalErr(v88v1, h, t, "enc-map-v88")
		if v != nil {
			v88v2 = make(map[uint]float64, len(v))
		}
		testUnmarshalErr(v88v2, bs88, h, t, "dec-map-v88")
		testDeepEqualErr(v88v1, v88v2, t, "equal-map-v88")
		bs88 = testMarshalErr(&v88v1, h, t, "enc-map-v88-p")
		v88v2 = nil
		testUnmarshalErr(&v88v2, bs88, h, t, "dec-map-v88-p")
		testDeepEqualErr(v88v1, v88v2, t, "equal-map-v88-p")
	}

	for _, v := range []map[uint]bool{nil, map[uint]bool{}, map[uint]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v89: %v\n", v)
		var v89v1, v89v2 map[uint]bool
		v89v1 = v
		bs89 := testMarshalErr(v89v1, h, t, "enc-map-v89")
		if v != nil {
			v89v2 = make(map[uint]bool, len(v))
		}
		testUnmarshalErr(v89v2, bs89, h, t, "dec-map-v89")
		testDeepEqualErr(v89v1, v89v2, t, "equal-map-v89")
		bs89 = testMarshalErr(&v89v1, h, t, "enc-map-v89-p")
		v89v2 = nil
		testUnmarshalErr(&v89v2, bs89, h, t, "dec-map-v89-p")
		testDeepEqualErr(v89v1, v89v2, t, "equal-map-v89-p")
	}

	for _, v := range []map[uint8]interface{}{nil, map[uint8]interface{}{}, map[uint8]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v91: %v\n", v)
		var v91v1, v91v2 map[uint8]interface{}
		v91v1 = v
		bs91 := testMarshalErr(v91v1, h, t, "enc-map-v91")
		if v != nil {
			v91v2 = make(map[uint8]interface{}, len(v))
		}
		testUnmarshalErr(v91v2, bs91, h, t, "dec-map-v91")
		testDeepEqualErr(v91v1, v91v2, t, "equal-map-v91")
		bs91 = testMarshalErr(&v91v1, h, t, "enc-map-v91-p")
		v91v2 = nil
		testUnmarshalErr(&v91v2, bs91, h, t, "dec-map-v91-p")
		testDeepEqualErr(v91v1, v91v2, t, "equal-map-v91-p")
	}

	for _, v := range []map[uint8]string{nil, map[uint8]string{}, map[uint8]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v92: %v\n", v)
		var v92v1, v92v2 map[uint8]string
		v92v1 = v
		bs92 := testMarshalErr(v92v1, h, t, "enc-map-v92")
		if v != nil {
			v92v2 = make(map[uint8]string, len(v))
		}
		testUnmarshalErr(v92v2, bs92, h, t, "dec-map-v92")
		testDeepEqualErr(v92v1, v92v2, t, "equal-map-v92")
		bs92 = testMarshalErr(&v92v1, h, t, "enc-map-v92-p")
		v92v2 = nil
		testUnmarshalErr(&v92v2, bs92, h, t, "dec-map-v92-p")
		testDeepEqualErr(v92v1, v92v2, t, "equal-map-v92-p")
	}

	for _, v := range []map[uint8]uint{nil, map[uint8]uint{}, map[uint8]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v93: %v\n", v)
		var v93v1, v93v2 map[uint8]uint
		v93v1 = v
		bs93 := testMarshalErr(v93v1, h, t, "enc-map-v93")
		if v != nil {
			v93v2 = make(map[uint8]uint, len(v))
		}
		testUnmarshalErr(v93v2, bs93, h, t, "dec-map-v93")
		testDeepEqualErr(v93v1, v93v2, t, "equal-map-v93")
		bs93 = testMarshalErr(&v93v1, h, t, "enc-map-v93-p")
		v93v2 = nil
		testUnmarshalErr(&v93v2, bs93, h, t, "dec-map-v93-p")
		testDeepEqualErr(v93v1, v93v2, t, "equal-map-v93-p")
	}

	for _, v := range []map[uint8]uint8{nil, map[uint8]uint8{}, map[uint8]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v94: %v\n", v)
		var v94v1, v94v2 map[uint8]uint8
		v94v1 = v
		bs94 := testMarshalErr(v94v1, h, t, "enc-map-v94")
		if v != nil {
			v94v2 = make(map[uint8]uint8, len(v))
		}
		testUnmarshalErr(v94v2, bs94, h, t, "dec-map-v94")
		testDeepEqualErr(v94v1, v94v2, t, "equal-map-v94")
		bs94 = testMarshalErr(&v94v1, h, t, "enc-map-v94-p")
		v94v2 = nil
		testUnmarshalErr(&v94v2, bs94, h, t, "dec-map-v94-p")
		testDeepEqualErr(v94v1, v94v2, t, "equal-map-v94-p")
	}

	for _, v := range []map[uint8]uint16{nil, map[uint8]uint16{}, map[uint8]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v95: %v\n", v)
		var v95v1, v95v2 map[uint8]uint16
		v95v1 = v
		bs95 := testMarshalErr(v95v1, h, t, "enc-map-v95")
		if v != nil {
			v95v2 = make(map[uint8]uint16, len(v))
		}
		testUnmarshalErr(v95v2, bs95, h, t, "dec-map-v95")
		testDeepEqualErr(v95v1, v95v2, t, "equal-map-v95")
		bs95 = testMarshalErr(&v95v1, h, t, "enc-map-v95-p")
		v95v2 = nil
		testUnmarshalErr(&v95v2, bs95, h, t, "dec-map-v95-p")
		testDeepEqualErr(v95v1, v95v2, t, "equal-map-v95-p")
	}

	for _, v := range []map[uint8]uint32{nil, map[uint8]uint32{}, map[uint8]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v96: %v\n", v)
		var v96v1, v96v2 map[uint8]uint32
		v96v1 = v
		bs96 := testMarshalErr(v96v1, h, t, "enc-map-v96")
		if v != nil {
			v96v2 = make(map[uint8]uint32, len(v))
		}
		testUnmarshalErr(v96v2, bs96, h, t, "dec-map-v96")
		testDeepEqualErr(v96v1, v96v2, t, "equal-map-v96")
		bs96 = testMarshalErr(&v96v1, h, t, "enc-map-v96-p")
		v96v2 = nil
		testUnmarshalErr(&v96v2, bs96, h, t, "dec-map-v96-p")
		testDeepEqualErr(v96v1, v96v2, t, "equal-map-v96-p")
	}

	for _, v := range []map[uint8]uint64{nil, map[uint8]uint64{}, map[uint8]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v97: %v\n", v)
		var v97v1, v97v2 map[uint8]uint64
		v97v1 = v
		bs97 := testMarshalErr(v97v1, h, t, "enc-map-v97")
		if v != nil {
			v97v2 = make(map[uint8]uint64, len(v))
		}
		testUnmarshalErr(v97v2, bs97, h, t, "dec-map-v97")
		testDeepEqualErr(v97v1, v97v2, t, "equal-map-v97")
		bs97 = testMarshalErr(&v97v1, h, t, "enc-map-v97-p")
		v97v2 = nil
		testUnmarshalErr(&v97v2, bs97, h, t, "dec-map-v97-p")
		testDeepEqualErr(v97v1, v97v2, t, "equal-map-v97-p")
	}

	for _, v := range []map[uint8]uintptr{nil, map[uint8]uintptr{}, map[uint8]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v98: %v\n", v)
		var v98v1, v98v2 map[uint8]uintptr
		v98v1 = v
		bs98 := testMarshalErr(v98v1, h, t, "enc-map-v98")
		if v != nil {
			v98v2 = make(map[uint8]uintptr, len(v))
		}
		testUnmarshalErr(v98v2, bs98, h, t, "dec-map-v98")
		testDeepEqualErr(v98v1, v98v2, t, "equal-map-v98")
		bs98 = testMarshalErr(&v98v1, h, t, "enc-map-v98-p")
		v98v2 = nil
		testUnmarshalErr(&v98v2, bs98, h, t, "dec-map-v98-p")
		testDeepEqualErr(v98v1, v98v2, t, "equal-map-v98-p")
	}

	for _, v := range []map[uint8]int{nil, map[uint8]int{}, map[uint8]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v99: %v\n", v)
		var v99v1, v99v2 map[uint8]int
		v99v1 = v
		bs99 := testMarshalErr(v99v1, h, t, "enc-map-v99")
		if v != nil {
			v99v2 = make(map[uint8]int, len(v))
		}
		testUnmarshalErr(v99v2, bs99, h, t, "dec-map-v99")
		testDeepEqualErr(v99v1, v99v2, t, "equal-map-v99")
		bs99 = testMarshalErr(&v99v1, h, t, "enc-map-v99-p")
		v99v2 = nil
		testUnmarshalErr(&v99v2, bs99, h, t, "dec-map-v99-p")
		testDeepEqualErr(v99v1, v99v2, t, "equal-map-v99-p")
	}

	for _, v := range []map[uint8]int8{nil, map[uint8]int8{}, map[uint8]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v100: %v\n", v)
		var v100v1, v100v2 map[uint8]int8
		v100v1 = v
		bs100 := testMarshalErr(v100v1, h, t, "enc-map-v100")
		if v != nil {
			v100v2 = make(map[uint8]int8, len(v))
		}
		testUnmarshalErr(v100v2, bs100, h, t, "dec-map-v100")
		testDeepEqualErr(v100v1, v100v2, t, "equal-map-v100")
		bs100 = testMarshalErr(&v100v1, h, t, "enc-map-v100-p")
		v100v2 = nil
		testUnmarshalErr(&v100v2, bs100, h, t, "dec-map-v100-p")
		testDeepEqualErr(v100v1, v100v2, t, "equal-map-v100-p")
	}

	for _, v := range []map[uint8]int16{nil, map[uint8]int16{}, map[uint8]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v101: %v\n", v)
		var v101v1, v101v2 map[uint8]int16
		v101v1 = v
		bs101 := testMarshalErr(v101v1, h, t, "enc-map-v101")
		if v != nil {
			v101v2 = make(map[uint8]int16, len(v))
		}
		testUnmarshalErr(v101v2, bs101, h, t, "dec-map-v101")
		testDeepEqualErr(v101v1, v101v2, t, "equal-map-v101")
		bs101 = testMarshalErr(&v101v1, h, t, "enc-map-v101-p")
		v101v2 = nil
		testUnmarshalErr(&v101v2, bs101, h, t, "dec-map-v101-p")
		testDeepEqualErr(v101v1, v101v2, t, "equal-map-v101-p")
	}

	for _, v := range []map[uint8]int32{nil, map[uint8]int32{}, map[uint8]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v102: %v\n", v)
		var v102v1, v102v2 map[uint8]int32
		v102v1 = v
		bs102 := testMarshalErr(v102v1, h, t, "enc-map-v102")
		if v != nil {
			v102v2 = make(map[uint8]int32, len(v))
		}
		testUnmarshalErr(v102v2, bs102, h, t, "dec-map-v102")
		testDeepEqualErr(v102v1, v102v2, t, "equal-map-v102")
		bs102 = testMarshalErr(&v102v1, h, t, "enc-map-v102-p")
		v102v2 = nil
		testUnmarshalErr(&v102v2, bs102, h, t, "dec-map-v102-p")
		testDeepEqualErr(v102v1, v102v2, t, "equal-map-v102-p")
	}

	for _, v := range []map[uint8]int64{nil, map[uint8]int64{}, map[uint8]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v103: %v\n", v)
		var v103v1, v103v2 map[uint8]int64
		v103v1 = v
		bs103 := testMarshalErr(v103v1, h, t, "enc-map-v103")
		if v != nil {
			v103v2 = make(map[uint8]int64, len(v))
		}
		testUnmarshalErr(v103v2, bs103, h, t, "dec-map-v103")
		testDeepEqualErr(v103v1, v103v2, t, "equal-map-v103")
		bs103 = testMarshalErr(&v103v1, h, t, "enc-map-v103-p")
		v103v2 = nil
		testUnmarshalErr(&v103v2, bs103, h, t, "dec-map-v103-p")
		testDeepEqualErr(v103v1, v103v2, t, "equal-map-v103-p")
	}

	for _, v := range []map[uint8]float32{nil, map[uint8]float32{}, map[uint8]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v104: %v\n", v)
		var v104v1, v104v2 map[uint8]float32
		v104v1 = v
		bs104 := testMarshalErr(v104v1, h, t, "enc-map-v104")
		if v != nil {
			v104v2 = make(map[uint8]float32, len(v))
		}
		testUnmarshalErr(v104v2, bs104, h, t, "dec-map-v104")
		testDeepEqualErr(v104v1, v104v2, t, "equal-map-v104")
		bs104 = testMarshalErr(&v104v1, h, t, "enc-map-v104-p")
		v104v2 = nil
		testUnmarshalErr(&v104v2, bs104, h, t, "dec-map-v104-p")
		testDeepEqualErr(v104v1, v104v2, t, "equal-map-v104-p")
	}

	for _, v := range []map[uint8]float64{nil, map[uint8]float64{}, map[uint8]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v105: %v\n", v)
		var v105v1, v105v2 map[uint8]float64
		v105v1 = v
		bs105 := testMarshalErr(v105v1, h, t, "enc-map-v105")
		if v != nil {
			v105v2 = make(map[uint8]float64, len(v))
		}
		testUnmarshalErr(v105v2, bs105, h, t, "dec-map-v105")
		testDeepEqualErr(v105v1, v105v2, t, "equal-map-v105")
		bs105 = testMarshalErr(&v105v1, h, t, "enc-map-v105-p")
		v105v2 = nil
		testUnmarshalErr(&v105v2, bs105, h, t, "dec-map-v105-p")
		testDeepEqualErr(v105v1, v105v2, t, "equal-map-v105-p")
	}

	for _, v := range []map[uint8]bool{nil, map[uint8]bool{}, map[uint8]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v106: %v\n", v)
		var v106v1, v106v2 map[uint8]bool
		v106v1 = v
		bs106 := testMarshalErr(v106v1, h, t, "enc-map-v106")
		if v != nil {
			v106v2 = make(map[uint8]bool, len(v))
		}
		testUnmarshalErr(v106v2, bs106, h, t, "dec-map-v106")
		testDeepEqualErr(v106v1, v106v2, t, "equal-map-v106")
		bs106 = testMarshalErr(&v106v1, h, t, "enc-map-v106-p")
		v106v2 = nil
		testUnmarshalErr(&v106v2, bs106, h, t, "dec-map-v106-p")
		testDeepEqualErr(v106v1, v106v2, t, "equal-map-v106-p")
	}

	for _, v := range []map[uint16]interface{}{nil, map[uint16]interface{}{}, map[uint16]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v109: %v\n", v)
		var v109v1, v109v2 map[uint16]interface{}
		v109v1 = v
		bs109 := testMarshalErr(v109v1, h, t, "enc-map-v109")
		if v != nil {
			v109v2 = make(map[uint16]interface{}, len(v))
		}
		testUnmarshalErr(v109v2, bs109, h, t, "dec-map-v109")
		testDeepEqualErr(v109v1, v109v2, t, "equal-map-v109")
		bs109 = testMarshalErr(&v109v1, h, t, "enc-map-v109-p")
		v109v2 = nil
		testUnmarshalErr(&v109v2, bs109, h, t, "dec-map-v109-p")
		testDeepEqualErr(v109v1, v109v2, t, "equal-map-v109-p")
	}

	for _, v := range []map[uint16]string{nil, map[uint16]string{}, map[uint16]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v110: %v\n", v)
		var v110v1, v110v2 map[uint16]string
		v110v1 = v
		bs110 := testMarshalErr(v110v1, h, t, "enc-map-v110")
		if v != nil {
			v110v2 = make(map[uint16]string, len(v))
		}
		testUnmarshalErr(v110v2, bs110, h, t, "dec-map-v110")
		testDeepEqualErr(v110v1, v110v2, t, "equal-map-v110")
		bs110 = testMarshalErr(&v110v1, h, t, "enc-map-v110-p")
		v110v2 = nil
		testUnmarshalErr(&v110v2, bs110, h, t, "dec-map-v110-p")
		testDeepEqualErr(v110v1, v110v2, t, "equal-map-v110-p")
	}

	for _, v := range []map[uint16]uint{nil, map[uint16]uint{}, map[uint16]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v111: %v\n", v)
		var v111v1, v111v2 map[uint16]uint
		v111v1 = v
		bs111 := testMarshalErr(v111v1, h, t, "enc-map-v111")
		if v != nil {
			v111v2 = make(map[uint16]uint, len(v))
		}
		testUnmarshalErr(v111v2, bs111, h, t, "dec-map-v111")
		testDeepEqualErr(v111v1, v111v2, t, "equal-map-v111")
		bs111 = testMarshalErr(&v111v1, h, t, "enc-map-v111-p")
		v111v2 = nil
		testUnmarshalErr(&v111v2, bs111, h, t, "dec-map-v111-p")
		testDeepEqualErr(v111v1, v111v2, t, "equal-map-v111-p")
	}

	for _, v := range []map[uint16]uint8{nil, map[uint16]uint8{}, map[uint16]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v112: %v\n", v)
		var v112v1, v112v2 map[uint16]uint8
		v112v1 = v
		bs112 := testMarshalErr(v112v1, h, t, "enc-map-v112")
		if v != nil {
			v112v2 = make(map[uint16]uint8, len(v))
		}
		testUnmarshalErr(v112v2, bs112, h, t, "dec-map-v112")
		testDeepEqualErr(v112v1, v112v2, t, "equal-map-v112")
		bs112 = testMarshalErr(&v112v1, h, t, "enc-map-v112-p")
		v112v2 = nil
		testUnmarshalErr(&v112v2, bs112, h, t, "dec-map-v112-p")
		testDeepEqualErr(v112v1, v112v2, t, "equal-map-v112-p")
	}

	for _, v := range []map[uint16]uint16{nil, map[uint16]uint16{}, map[uint16]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v113: %v\n", v)
		var v113v1, v113v2 map[uint16]uint16
		v113v1 = v
		bs113 := testMarshalErr(v113v1, h, t, "enc-map-v113")
		if v != nil {
			v113v2 = make(map[uint16]uint16, len(v))
		}
		testUnmarshalErr(v113v2, bs113, h, t, "dec-map-v113")
		testDeepEqualErr(v113v1, v113v2, t, "equal-map-v113")
		bs113 = testMarshalErr(&v113v1, h, t, "enc-map-v113-p")
		v113v2 = nil
		testUnmarshalErr(&v113v2, bs113, h, t, "dec-map-v113-p")
		testDeepEqualErr(v113v1, v113v2, t, "equal-map-v113-p")
	}

	for _, v := range []map[uint16]uint32{nil, map[uint16]uint32{}, map[uint16]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v114: %v\n", v)
		var v114v1, v114v2 map[uint16]uint32
		v114v1 = v
		bs114 := testMarshalErr(v114v1, h, t, "enc-map-v114")
		if v != nil {
			v114v2 = make(map[uint16]uint32, len(v))
		}
		testUnmarshalErr(v114v2, bs114, h, t, "dec-map-v114")
		testDeepEqualErr(v114v1, v114v2, t, "equal-map-v114")
		bs114 = testMarshalErr(&v114v1, h, t, "enc-map-v114-p")
		v114v2 = nil
		testUnmarshalErr(&v114v2, bs114, h, t, "dec-map-v114-p")
		testDeepEqualErr(v114v1, v114v2, t, "equal-map-v114-p")
	}

	for _, v := range []map[uint16]uint64{nil, map[uint16]uint64{}, map[uint16]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v115: %v\n", v)
		var v115v1, v115v2 map[uint16]uint64
		v115v1 = v
		bs115 := testMarshalErr(v115v1, h, t, "enc-map-v115")
		if v != nil {
			v115v2 = make(map[uint16]uint64, len(v))
		}
		testUnmarshalErr(v115v2, bs115, h, t, "dec-map-v115")
		testDeepEqualErr(v115v1, v115v2, t, "equal-map-v115")
		bs115 = testMarshalErr(&v115v1, h, t, "enc-map-v115-p")
		v115v2 = nil
		testUnmarshalErr(&v115v2, bs115, h, t, "dec-map-v115-p")
		testDeepEqualErr(v115v1, v115v2, t, "equal-map-v115-p")
	}

	for _, v := range []map[uint16]uintptr{nil, map[uint16]uintptr{}, map[uint16]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v116: %v\n", v)
		var v116v1, v116v2 map[uint16]uintptr
		v116v1 = v
		bs116 := testMarshalErr(v116v1, h, t, "enc-map-v116")
		if v != nil {
			v116v2 = make(map[uint16]uintptr, len(v))
		}
		testUnmarshalErr(v116v2, bs116, h, t, "dec-map-v116")
		testDeepEqualErr(v116v1, v116v2, t, "equal-map-v116")
		bs116 = testMarshalErr(&v116v1, h, t, "enc-map-v116-p")
		v116v2 = nil
		testUnmarshalErr(&v116v2, bs116, h, t, "dec-map-v116-p")
		testDeepEqualErr(v116v1, v116v2, t, "equal-map-v116-p")
	}

	for _, v := range []map[uint16]int{nil, map[uint16]int{}, map[uint16]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v117: %v\n", v)
		var v117v1, v117v2 map[uint16]int
		v117v1 = v
		bs117 := testMarshalErr(v117v1, h, t, "enc-map-v117")
		if v != nil {
			v117v2 = make(map[uint16]int, len(v))
		}
		testUnmarshalErr(v117v2, bs117, h, t, "dec-map-v117")
		testDeepEqualErr(v117v1, v117v2, t, "equal-map-v117")
		bs117 = testMarshalErr(&v117v1, h, t, "enc-map-v117-p")
		v117v2 = nil
		testUnmarshalErr(&v117v2, bs117, h, t, "dec-map-v117-p")
		testDeepEqualErr(v117v1, v117v2, t, "equal-map-v117-p")
	}

	for _, v := range []map[uint16]int8{nil, map[uint16]int8{}, map[uint16]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v118: %v\n", v)
		var v118v1, v118v2 map[uint16]int8
		v118v1 = v
		bs118 := testMarshalErr(v118v1, h, t, "enc-map-v118")
		if v != nil {
			v118v2 = make(map[uint16]int8, len(v))
		}
		testUnmarshalErr(v118v2, bs118, h, t, "dec-map-v118")
		testDeepEqualErr(v118v1, v118v2, t, "equal-map-v118")
		bs118 = testMarshalErr(&v118v1, h, t, "enc-map-v118-p")
		v118v2 = nil
		testUnmarshalErr(&v118v2, bs118, h, t, "dec-map-v118-p")
		testDeepEqualErr(v118v1, v118v2, t, "equal-map-v118-p")
	}

	for _, v := range []map[uint16]int16{nil, map[uint16]int16{}, map[uint16]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v119: %v\n", v)
		var v119v1, v119v2 map[uint16]int16
		v119v1 = v
		bs119 := testMarshalErr(v119v1, h, t, "enc-map-v119")
		if v != nil {
			v119v2 = make(map[uint16]int16, len(v))
		}
		testUnmarshalErr(v119v2, bs119, h, t, "dec-map-v119")
		testDeepEqualErr(v119v1, v119v2, t, "equal-map-v119")
		bs119 = testMarshalErr(&v119v1, h, t, "enc-map-v119-p")
		v119v2 = nil
		testUnmarshalErr(&v119v2, bs119, h, t, "dec-map-v119-p")
		testDeepEqualErr(v119v1, v119v2, t, "equal-map-v119-p")
	}

	for _, v := range []map[uint16]int32{nil, map[uint16]int32{}, map[uint16]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v120: %v\n", v)
		var v120v1, v120v2 map[uint16]int32
		v120v1 = v
		bs120 := testMarshalErr(v120v1, h, t, "enc-map-v120")
		if v != nil {
			v120v2 = make(map[uint16]int32, len(v))
		}
		testUnmarshalErr(v120v2, bs120, h, t, "dec-map-v120")
		testDeepEqualErr(v120v1, v120v2, t, "equal-map-v120")
		bs120 = testMarshalErr(&v120v1, h, t, "enc-map-v120-p")
		v120v2 = nil
		testUnmarshalErr(&v120v2, bs120, h, t, "dec-map-v120-p")
		testDeepEqualErr(v120v1, v120v2, t, "equal-map-v120-p")
	}

	for _, v := range []map[uint16]int64{nil, map[uint16]int64{}, map[uint16]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v121: %v\n", v)
		var v121v1, v121v2 map[uint16]int64
		v121v1 = v
		bs121 := testMarshalErr(v121v1, h, t, "enc-map-v121")
		if v != nil {
			v121v2 = make(map[uint16]int64, len(v))
		}
		testUnmarshalErr(v121v2, bs121, h, t, "dec-map-v121")
		testDeepEqualErr(v121v1, v121v2, t, "equal-map-v121")
		bs121 = testMarshalErr(&v121v1, h, t, "enc-map-v121-p")
		v121v2 = nil
		testUnmarshalErr(&v121v2, bs121, h, t, "dec-map-v121-p")
		testDeepEqualErr(v121v1, v121v2, t, "equal-map-v121-p")
	}

	for _, v := range []map[uint16]float32{nil, map[uint16]float32{}, map[uint16]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v122: %v\n", v)
		var v122v1, v122v2 map[uint16]float32
		v122v1 = v
		bs122 := testMarshalErr(v122v1, h, t, "enc-map-v122")
		if v != nil {
			v122v2 = make(map[uint16]float32, len(v))
		}
		testUnmarshalErr(v122v2, bs122, h, t, "dec-map-v122")
		testDeepEqualErr(v122v1, v122v2, t, "equal-map-v122")
		bs122 = testMarshalErr(&v122v1, h, t, "enc-map-v122-p")
		v122v2 = nil
		testUnmarshalErr(&v122v2, bs122, h, t, "dec-map-v122-p")
		testDeepEqualErr(v122v1, v122v2, t, "equal-map-v122-p")
	}

	for _, v := range []map[uint16]float64{nil, map[uint16]float64{}, map[uint16]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v123: %v\n", v)
		var v123v1, v123v2 map[uint16]float64
		v123v1 = v
		bs123 := testMarshalErr(v123v1, h, t, "enc-map-v123")
		if v != nil {
			v123v2 = make(map[uint16]float64, len(v))
		}
		testUnmarshalErr(v123v2, bs123, h, t, "dec-map-v123")
		testDeepEqualErr(v123v1, v123v2, t, "equal-map-v123")
		bs123 = testMarshalErr(&v123v1, h, t, "enc-map-v123-p")
		v123v2 = nil
		testUnmarshalErr(&v123v2, bs123, h, t, "dec-map-v123-p")
		testDeepEqualErr(v123v1, v123v2, t, "equal-map-v123-p")
	}

	for _, v := range []map[uint16]bool{nil, map[uint16]bool{}, map[uint16]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v124: %v\n", v)
		var v124v1, v124v2 map[uint16]bool
		v124v1 = v
		bs124 := testMarshalErr(v124v1, h, t, "enc-map-v124")
		if v != nil {
			v124v2 = make(map[uint16]bool, len(v))
		}
		testUnmarshalErr(v124v2, bs124, h, t, "dec-map-v124")
		testDeepEqualErr(v124v1, v124v2, t, "equal-map-v124")
		bs124 = testMarshalErr(&v124v1, h, t, "enc-map-v124-p")
		v124v2 = nil
		testUnmarshalErr(&v124v2, bs124, h, t, "dec-map-v124-p")
		testDeepEqualErr(v124v1, v124v2, t, "equal-map-v124-p")
	}

	for _, v := range []map[uint32]interface{}{nil, map[uint32]interface{}{}, map[uint32]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v127: %v\n", v)
		var v127v1, v127v2 map[uint32]interface{}
		v127v1 = v
		bs127 := testMarshalErr(v127v1, h, t, "enc-map-v127")
		if v != nil {
			v127v2 = make(map[uint32]interface{}, len(v))
		}
		testUnmarshalErr(v127v2, bs127, h, t, "dec-map-v127")
		testDeepEqualErr(v127v1, v127v2, t, "equal-map-v127")
		bs127 = testMarshalErr(&v127v1, h, t, "enc-map-v127-p")
		v127v2 = nil
		testUnmarshalErr(&v127v2, bs127, h, t, "dec-map-v127-p")
		testDeepEqualErr(v127v1, v127v2, t, "equal-map-v127-p")
	}

	for _, v := range []map[uint32]string{nil, map[uint32]string{}, map[uint32]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v128: %v\n", v)
		var v128v1, v128v2 map[uint32]string
		v128v1 = v
		bs128 := testMarshalErr(v128v1, h, t, "enc-map-v128")
		if v != nil {
			v128v2 = make(map[uint32]string, len(v))
		}
		testUnmarshalErr(v128v2, bs128, h, t, "dec-map-v128")
		testDeepEqualErr(v128v1, v128v2, t, "equal-map-v128")
		bs128 = testMarshalErr(&v128v1, h, t, "enc-map-v128-p")
		v128v2 = nil
		testUnmarshalErr(&v128v2, bs128, h, t, "dec-map-v128-p")
		testDeepEqualErr(v128v1, v128v2, t, "equal-map-v128-p")
	}

	for _, v := range []map[uint32]uint{nil, map[uint32]uint{}, map[uint32]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v129: %v\n", v)
		var v129v1, v129v2 map[uint32]uint
		v129v1 = v
		bs129 := testMarshalErr(v129v1, h, t, "enc-map-v129")
		if v != nil {
			v129v2 = make(map[uint32]uint, len(v))
		}
		testUnmarshalErr(v129v2, bs129, h, t, "dec-map-v129")
		testDeepEqualErr(v129v1, v129v2, t, "equal-map-v129")
		bs129 = testMarshalErr(&v129v1, h, t, "enc-map-v129-p")
		v129v2 = nil
		testUnmarshalErr(&v129v2, bs129, h, t, "dec-map-v129-p")
		testDeepEqualErr(v129v1, v129v2, t, "equal-map-v129-p")
	}

	for _, v := range []map[uint32]uint8{nil, map[uint32]uint8{}, map[uint32]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v130: %v\n", v)
		var v130v1, v130v2 map[uint32]uint8
		v130v1 = v
		bs130 := testMarshalErr(v130v1, h, t, "enc-map-v130")
		if v != nil {
			v130v2 = make(map[uint32]uint8, len(v))
		}
		testUnmarshalErr(v130v2, bs130, h, t, "dec-map-v130")
		testDeepEqualErr(v130v1, v130v2, t, "equal-map-v130")
		bs130 = testMarshalErr(&v130v1, h, t, "enc-map-v130-p")
		v130v2 = nil
		testUnmarshalErr(&v130v2, bs130, h, t, "dec-map-v130-p")
		testDeepEqualErr(v130v1, v130v2, t, "equal-map-v130-p")
	}

	for _, v := range []map[uint32]uint16{nil, map[uint32]uint16{}, map[uint32]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v131: %v\n", v)
		var v131v1, v131v2 map[uint32]uint16
		v131v1 = v
		bs131 := testMarshalErr(v131v1, h, t, "enc-map-v131")
		if v != nil {
			v131v2 = make(map[uint32]uint16, len(v))
		}
		testUnmarshalErr(v131v2, bs131, h, t, "dec-map-v131")
		testDeepEqualErr(v131v1, v131v2, t, "equal-map-v131")
		bs131 = testMarshalErr(&v131v1, h, t, "enc-map-v131-p")
		v131v2 = nil
		testUnmarshalErr(&v131v2, bs131, h, t, "dec-map-v131-p")
		testDeepEqualErr(v131v1, v131v2, t, "equal-map-v131-p")
	}

	for _, v := range []map[uint32]uint32{nil, map[uint32]uint32{}, map[uint32]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v132: %v\n", v)
		var v132v1, v132v2 map[uint32]uint32
		v132v1 = v
		bs132 := testMarshalErr(v132v1, h, t, "enc-map-v132")
		if v != nil {
			v132v2 = make(map[uint32]uint32, len(v))
		}
		testUnmarshalErr(v132v2, bs132, h, t, "dec-map-v132")
		testDeepEqualErr(v132v1, v132v2, t, "equal-map-v132")
		bs132 = testMarshalErr(&v132v1, h, t, "enc-map-v132-p")
		v132v2 = nil
		testUnmarshalErr(&v132v2, bs132, h, t, "dec-map-v132-p")
		testDeepEqualErr(v132v1, v132v2, t, "equal-map-v132-p")
	}

	for _, v := range []map[uint32]uint64{nil, map[uint32]uint64{}, map[uint32]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v133: %v\n", v)
		var v133v1, v133v2 map[uint32]uint64
		v133v1 = v
		bs133 := testMarshalErr(v133v1, h, t, "enc-map-v133")
		if v != nil {
			v133v2 = make(map[uint32]uint64, len(v))
		}
		testUnmarshalErr(v133v2, bs133, h, t, "dec-map-v133")
		testDeepEqualErr(v133v1, v133v2, t, "equal-map-v133")
		bs133 = testMarshalErr(&v133v1, h, t, "enc-map-v133-p")
		v133v2 = nil
		testUnmarshalErr(&v133v2, bs133, h, t, "dec-map-v133-p")
		testDeepEqualErr(v133v1, v133v2, t, "equal-map-v133-p")
	}

	for _, v := range []map[uint32]uintptr{nil, map[uint32]uintptr{}, map[uint32]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v134: %v\n", v)
		var v134v1, v134v2 map[uint32]uintptr
		v134v1 = v
		bs134 := testMarshalErr(v134v1, h, t, "enc-map-v134")
		if v != nil {
			v134v2 = make(map[uint32]uintptr, len(v))
		}
		testUnmarshalErr(v134v2, bs134, h, t, "dec-map-v134")
		testDeepEqualErr(v134v1, v134v2, t, "equal-map-v134")
		bs134 = testMarshalErr(&v134v1, h, t, "enc-map-v134-p")
		v134v2 = nil
		testUnmarshalErr(&v134v2, bs134, h, t, "dec-map-v134-p")
		testDeepEqualErr(v134v1, v134v2, t, "equal-map-v134-p")
	}

	for _, v := range []map[uint32]int{nil, map[uint32]int{}, map[uint32]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v135: %v\n", v)
		var v135v1, v135v2 map[uint32]int
		v135v1 = v
		bs135 := testMarshalErr(v135v1, h, t, "enc-map-v135")
		if v != nil {
			v135v2 = make(map[uint32]int, len(v))
		}
		testUnmarshalErr(v135v2, bs135, h, t, "dec-map-v135")
		testDeepEqualErr(v135v1, v135v2, t, "equal-map-v135")
		bs135 = testMarshalErr(&v135v1, h, t, "enc-map-v135-p")
		v135v2 = nil
		testUnmarshalErr(&v135v2, bs135, h, t, "dec-map-v135-p")
		testDeepEqualErr(v135v1, v135v2, t, "equal-map-v135-p")
	}

	for _, v := range []map[uint32]int8{nil, map[uint32]int8{}, map[uint32]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v136: %v\n", v)
		var v136v1, v136v2 map[uint32]int8
		v136v1 = v
		bs136 := testMarshalErr(v136v1, h, t, "enc-map-v136")
		if v != nil {
			v136v2 = make(map[uint32]int8, len(v))
		}
		testUnmarshalErr(v136v2, bs136, h, t, "dec-map-v136")
		testDeepEqualErr(v136v1, v136v2, t, "equal-map-v136")
		bs136 = testMarshalErr(&v136v1, h, t, "enc-map-v136-p")
		v136v2 = nil
		testUnmarshalErr(&v136v2, bs136, h, t, "dec-map-v136-p")
		testDeepEqualErr(v136v1, v136v2, t, "equal-map-v136-p")
	}

	for _, v := range []map[uint32]int16{nil, map[uint32]int16{}, map[uint32]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v137: %v\n", v)
		var v137v1, v137v2 map[uint32]int16
		v137v1 = v
		bs137 := testMarshalErr(v137v1, h, t, "enc-map-v137")
		if v != nil {
			v137v2 = make(map[uint32]int16, len(v))
		}
		testUnmarshalErr(v137v2, bs137, h, t, "dec-map-v137")
		testDeepEqualErr(v137v1, v137v2, t, "equal-map-v137")
		bs137 = testMarshalErr(&v137v1, h, t, "enc-map-v137-p")
		v137v2 = nil
		testUnmarshalErr(&v137v2, bs137, h, t, "dec-map-v137-p")
		testDeepEqualErr(v137v1, v137v2, t, "equal-map-v137-p")
	}

	for _, v := range []map[uint32]int32{nil, map[uint32]int32{}, map[uint32]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v138: %v\n", v)
		var v138v1, v138v2 map[uint32]int32
		v138v1 = v
		bs138 := testMarshalErr(v138v1, h, t, "enc-map-v138")
		if v != nil {
			v138v2 = make(map[uint32]int32, len(v))
		}
		testUnmarshalErr(v138v2, bs138, h, t, "dec-map-v138")
		testDeepEqualErr(v138v1, v138v2, t, "equal-map-v138")
		bs138 = testMarshalErr(&v138v1, h, t, "enc-map-v138-p")
		v138v2 = nil
		testUnmarshalErr(&v138v2, bs138, h, t, "dec-map-v138-p")
		testDeepEqualErr(v138v1, v138v2, t, "equal-map-v138-p")
	}

	for _, v := range []map[uint32]int64{nil, map[uint32]int64{}, map[uint32]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v139: %v\n", v)
		var v139v1, v139v2 map[uint32]int64
		v139v1 = v
		bs139 := testMarshalErr(v139v1, h, t, "enc-map-v139")
		if v != nil {
			v139v2 = make(map[uint32]int64, len(v))
		}
		testUnmarshalErr(v139v2, bs139, h, t, "dec-map-v139")
		testDeepEqualErr(v139v1, v139v2, t, "equal-map-v139")
		bs139 = testMarshalErr(&v139v1, h, t, "enc-map-v139-p")
		v139v2 = nil
		testUnmarshalErr(&v139v2, bs139, h, t, "dec-map-v139-p")
		testDeepEqualErr(v139v1, v139v2, t, "equal-map-v139-p")
	}

	for _, v := range []map[uint32]float32{nil, map[uint32]float32{}, map[uint32]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v140: %v\n", v)
		var v140v1, v140v2 map[uint32]float32
		v140v1 = v
		bs140 := testMarshalErr(v140v1, h, t, "enc-map-v140")
		if v != nil {
			v140v2 = make(map[uint32]float32, len(v))
		}
		testUnmarshalErr(v140v2, bs140, h, t, "dec-map-v140")
		testDeepEqualErr(v140v1, v140v2, t, "equal-map-v140")
		bs140 = testMarshalErr(&v140v1, h, t, "enc-map-v140-p")
		v140v2 = nil
		testUnmarshalErr(&v140v2, bs140, h, t, "dec-map-v140-p")
		testDeepEqualErr(v140v1, v140v2, t, "equal-map-v140-p")
	}

	for _, v := range []map[uint32]float64{nil, map[uint32]float64{}, map[uint32]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v141: %v\n", v)
		var v141v1, v141v2 map[uint32]float64
		v141v1 = v
		bs141 := testMarshalErr(v141v1, h, t, "enc-map-v141")
		if v != nil {
			v141v2 = make(map[uint32]float64, len(v))
		}
		testUnmarshalErr(v141v2, bs141, h, t, "dec-map-v141")
		testDeepEqualErr(v141v1, v141v2, t, "equal-map-v141")
		bs141 = testMarshalErr(&v141v1, h, t, "enc-map-v141-p")
		v141v2 = nil
		testUnmarshalErr(&v141v2, bs141, h, t, "dec-map-v141-p")
		testDeepEqualErr(v141v1, v141v2, t, "equal-map-v141-p")
	}

	for _, v := range []map[uint32]bool{nil, map[uint32]bool{}, map[uint32]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v142: %v\n", v)
		var v142v1, v142v2 map[uint32]bool
		v142v1 = v
		bs142 := testMarshalErr(v142v1, h, t, "enc-map-v142")
		if v != nil {
			v142v2 = make(map[uint32]bool, len(v))
		}
		testUnmarshalErr(v142v2, bs142, h, t, "dec-map-v142")
		testDeepEqualErr(v142v1, v142v2, t, "equal-map-v142")
		bs142 = testMarshalErr(&v142v1, h, t, "enc-map-v142-p")
		v142v2 = nil
		testUnmarshalErr(&v142v2, bs142, h, t, "dec-map-v142-p")
		testDeepEqualErr(v142v1, v142v2, t, "equal-map-v142-p")
	}

	for _, v := range []map[uint64]interface{}{nil, map[uint64]interface{}{}, map[uint64]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v145: %v\n", v)
		var v145v1, v145v2 map[uint64]interface{}
		v145v1 = v
		bs145 := testMarshalErr(v145v1, h, t, "enc-map-v145")
		if v != nil {
			v145v2 = make(map[uint64]interface{}, len(v))
		}
		testUnmarshalErr(v145v2, bs145, h, t, "dec-map-v145")
		testDeepEqualErr(v145v1, v145v2, t, "equal-map-v145")
		bs145 = testMarshalErr(&v145v1, h, t, "enc-map-v145-p")
		v145v2 = nil
		testUnmarshalErr(&v145v2, bs145, h, t, "dec-map-v145-p")
		testDeepEqualErr(v145v1, v145v2, t, "equal-map-v145-p")
	}

	for _, v := range []map[uint64]string{nil, map[uint64]string{}, map[uint64]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v146: %v\n", v)
		var v146v1, v146v2 map[uint64]string
		v146v1 = v
		bs146 := testMarshalErr(v146v1, h, t, "enc-map-v146")
		if v != nil {
			v146v2 = make(map[uint64]string, len(v))
		}
		testUnmarshalErr(v146v2, bs146, h, t, "dec-map-v146")
		testDeepEqualErr(v146v1, v146v2, t, "equal-map-v146")
		bs146 = testMarshalErr(&v146v1, h, t, "enc-map-v146-p")
		v146v2 = nil
		testUnmarshalErr(&v146v2, bs146, h, t, "dec-map-v146-p")
		testDeepEqualErr(v146v1, v146v2, t, "equal-map-v146-p")
	}

	for _, v := range []map[uint64]uint{nil, map[uint64]uint{}, map[uint64]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v147: %v\n", v)
		var v147v1, v147v2 map[uint64]uint
		v147v1 = v
		bs147 := testMarshalErr(v147v1, h, t, "enc-map-v147")
		if v != nil {
			v147v2 = make(map[uint64]uint, len(v))
		}
		testUnmarshalErr(v147v2, bs147, h, t, "dec-map-v147")
		testDeepEqualErr(v147v1, v147v2, t, "equal-map-v147")
		bs147 = testMarshalErr(&v147v1, h, t, "enc-map-v147-p")
		v147v2 = nil
		testUnmarshalErr(&v147v2, bs147, h, t, "dec-map-v147-p")
		testDeepEqualErr(v147v1, v147v2, t, "equal-map-v147-p")
	}

	for _, v := range []map[uint64]uint8{nil, map[uint64]uint8{}, map[uint64]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v148: %v\n", v)
		var v148v1, v148v2 map[uint64]uint8
		v148v1 = v
		bs148 := testMarshalErr(v148v1, h, t, "enc-map-v148")
		if v != nil {
			v148v2 = make(map[uint64]uint8, len(v))
		}
		testUnmarshalErr(v148v2, bs148, h, t, "dec-map-v148")
		testDeepEqualErr(v148v1, v148v2, t, "equal-map-v148")
		bs148 = testMarshalErr(&v148v1, h, t, "enc-map-v148-p")
		v148v2 = nil
		testUnmarshalErr(&v148v2, bs148, h, t, "dec-map-v148-p")
		testDeepEqualErr(v148v1, v148v2, t, "equal-map-v148-p")
	}

	for _, v := range []map[uint64]uint16{nil, map[uint64]uint16{}, map[uint64]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v149: %v\n", v)
		var v149v1, v149v2 map[uint64]uint16
		v149v1 = v
		bs149 := testMarshalErr(v149v1, h, t, "enc-map-v149")
		if v != nil {
			v149v2 = make(map[uint64]uint16, len(v))
		}
		testUnmarshalErr(v149v2, bs149, h, t, "dec-map-v149")
		testDeepEqualErr(v149v1, v149v2, t, "equal-map-v149")
		bs149 = testMarshalErr(&v149v1, h, t, "enc-map-v149-p")
		v149v2 = nil
		testUnmarshalErr(&v149v2, bs149, h, t, "dec-map-v149-p")
		testDeepEqualErr(v149v1, v149v2, t, "equal-map-v149-p")
	}

	for _, v := range []map[uint64]uint32{nil, map[uint64]uint32{}, map[uint64]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v150: %v\n", v)
		var v150v1, v150v2 map[uint64]uint32
		v150v1 = v
		bs150 := testMarshalErr(v150v1, h, t, "enc-map-v150")
		if v != nil {
			v150v2 = make(map[uint64]uint32, len(v))
		}
		testUnmarshalErr(v150v2, bs150, h, t, "dec-map-v150")
		testDeepEqualErr(v150v1, v150v2, t, "equal-map-v150")
		bs150 = testMarshalErr(&v150v1, h, t, "enc-map-v150-p")
		v150v2 = nil
		testUnmarshalErr(&v150v2, bs150, h, t, "dec-map-v150-p")
		testDeepEqualErr(v150v1, v150v2, t, "equal-map-v150-p")
	}

	for _, v := range []map[uint64]uint64{nil, map[uint64]uint64{}, map[uint64]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v151: %v\n", v)
		var v151v1, v151v2 map[uint64]uint64
		v151v1 = v
		bs151 := testMarshalErr(v151v1, h, t, "enc-map-v151")
		if v != nil {
			v151v2 = make(map[uint64]uint64, len(v))
		}
		testUnmarshalErr(v151v2, bs151, h, t, "dec-map-v151")
		testDeepEqualErr(v151v1, v151v2, t, "equal-map-v151")
		bs151 = testMarshalErr(&v151v1, h, t, "enc-map-v151-p")
		v151v2 = nil
		testUnmarshalErr(&v151v2, bs151, h, t, "dec-map-v151-p")
		testDeepEqualErr(v151v1, v151v2, t, "equal-map-v151-p")
	}

	for _, v := range []map[uint64]uintptr{nil, map[uint64]uintptr{}, map[uint64]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v152: %v\n", v)
		var v152v1, v152v2 map[uint64]uintptr
		v152v1 = v
		bs152 := testMarshalErr(v152v1, h, t, "enc-map-v152")
		if v != nil {
			v152v2 = make(map[uint64]uintptr, len(v))
		}
		testUnmarshalErr(v152v2, bs152, h, t, "dec-map-v152")
		testDeepEqualErr(v152v1, v152v2, t, "equal-map-v152")
		bs152 = testMarshalErr(&v152v1, h, t, "enc-map-v152-p")
		v152v2 = nil
		testUnmarshalErr(&v152v2, bs152, h, t, "dec-map-v152-p")
		testDeepEqualErr(v152v1, v152v2, t, "equal-map-v152-p")
	}

	for _, v := range []map[uint64]int{nil, map[uint64]int{}, map[uint64]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v153: %v\n", v)
		var v153v1, v153v2 map[uint64]int
		v153v1 = v
		bs153 := testMarshalErr(v153v1, h, t, "enc-map-v153")
		if v != nil {
			v153v2 = make(map[uint64]int, len(v))
		}
		testUnmarshalErr(v153v2, bs153, h, t, "dec-map-v153")
		testDeepEqualErr(v153v1, v153v2, t, "equal-map-v153")
		bs153 = testMarshalErr(&v153v1, h, t, "enc-map-v153-p")
		v153v2 = nil
		testUnmarshalErr(&v153v2, bs153, h, t, "dec-map-v153-p")
		testDeepEqualErr(v153v1, v153v2, t, "equal-map-v153-p")
	}

	for _, v := range []map[uint64]int8{nil, map[uint64]int8{}, map[uint64]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v154: %v\n", v)
		var v154v1, v154v2 map[uint64]int8
		v154v1 = v
		bs154 := testMarshalErr(v154v1, h, t, "enc-map-v154")
		if v != nil {
			v154v2 = make(map[uint64]int8, len(v))
		}
		testUnmarshalErr(v154v2, bs154, h, t, "dec-map-v154")
		testDeepEqualErr(v154v1, v154v2, t, "equal-map-v154")
		bs154 = testMarshalErr(&v154v1, h, t, "enc-map-v154-p")
		v154v2 = nil
		testUnmarshalErr(&v154v2, bs154, h, t, "dec-map-v154-p")
		testDeepEqualErr(v154v1, v154v2, t, "equal-map-v154-p")
	}

	for _, v := range []map[uint64]int16{nil, map[uint64]int16{}, map[uint64]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v155: %v\n", v)
		var v155v1, v155v2 map[uint64]int16
		v155v1 = v
		bs155 := testMarshalErr(v155v1, h, t, "enc-map-v155")
		if v != nil {
			v155v2 = make(map[uint64]int16, len(v))
		}
		testUnmarshalErr(v155v2, bs155, h, t, "dec-map-v155")
		testDeepEqualErr(v155v1, v155v2, t, "equal-map-v155")
		bs155 = testMarshalErr(&v155v1, h, t, "enc-map-v155-p")
		v155v2 = nil
		testUnmarshalErr(&v155v2, bs155, h, t, "dec-map-v155-p")
		testDeepEqualErr(v155v1, v155v2, t, "equal-map-v155-p")
	}

	for _, v := range []map[uint64]int32{nil, map[uint64]int32{}, map[uint64]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v156: %v\n", v)
		var v156v1, v156v2 map[uint64]int32
		v156v1 = v
		bs156 := testMarshalErr(v156v1, h, t, "enc-map-v156")
		if v != nil {
			v156v2 = make(map[uint64]int32, len(v))
		}
		testUnmarshalErr(v156v2, bs156, h, t, "dec-map-v156")
		testDeepEqualErr(v156v1, v156v2, t, "equal-map-v156")
		bs156 = testMarshalErr(&v156v1, h, t, "enc-map-v156-p")
		v156v2 = nil
		testUnmarshalErr(&v156v2, bs156, h, t, "dec-map-v156-p")
		testDeepEqualErr(v156v1, v156v2, t, "equal-map-v156-p")
	}

	for _, v := range []map[uint64]int64{nil, map[uint64]int64{}, map[uint64]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v157: %v\n", v)
		var v157v1, v157v2 map[uint64]int64
		v157v1 = v
		bs157 := testMarshalErr(v157v1, h, t, "enc-map-v157")
		if v != nil {
			v157v2 = make(map[uint64]int64, len(v))
		}
		testUnmarshalErr(v157v2, bs157, h, t, "dec-map-v157")
		testDeepEqualErr(v157v1, v157v2, t, "equal-map-v157")
		bs157 = testMarshalErr(&v157v1, h, t, "enc-map-v157-p")
		v157v2 = nil
		testUnmarshalErr(&v157v2, bs157, h, t, "dec-map-v157-p")
		testDeepEqualErr(v157v1, v157v2, t, "equal-map-v157-p")
	}

	for _, v := range []map[uint64]float32{nil, map[uint64]float32{}, map[uint64]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v158: %v\n", v)
		var v158v1, v158v2 map[uint64]float32
		v158v1 = v
		bs158 := testMarshalErr(v158v1, h, t, "enc-map-v158")
		if v != nil {
			v158v2 = make(map[uint64]float32, len(v))
		}
		testUnmarshalErr(v158v2, bs158, h, t, "dec-map-v158")
		testDeepEqualErr(v158v1, v158v2, t, "equal-map-v158")
		bs158 = testMarshalErr(&v158v1, h, t, "enc-map-v158-p")
		v158v2 = nil
		testUnmarshalErr(&v158v2, bs158, h, t, "dec-map-v158-p")
		testDeepEqualErr(v158v1, v158v2, t, "equal-map-v158-p")
	}

	for _, v := range []map[uint64]float64{nil, map[uint64]float64{}, map[uint64]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v159: %v\n", v)
		var v159v1, v159v2 map[uint64]float64
		v159v1 = v
		bs159 := testMarshalErr(v159v1, h, t, "enc-map-v159")
		if v != nil {
			v159v2 = make(map[uint64]float64, len(v))
		}
		testUnmarshalErr(v159v2, bs159, h, t, "dec-map-v159")
		testDeepEqualErr(v159v1, v159v2, t, "equal-map-v159")
		bs159 = testMarshalErr(&v159v1, h, t, "enc-map-v159-p")
		v159v2 = nil
		testUnmarshalErr(&v159v2, bs159, h, t, "dec-map-v159-p")
		testDeepEqualErr(v159v1, v159v2, t, "equal-map-v159-p")
	}

	for _, v := range []map[uint64]bool{nil, map[uint64]bool{}, map[uint64]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v160: %v\n", v)
		var v160v1, v160v2 map[uint64]bool
		v160v1 = v
		bs160 := testMarshalErr(v160v1, h, t, "enc-map-v160")
		if v != nil {
			v160v2 = make(map[uint64]bool, len(v))
		}
		testUnmarshalErr(v160v2, bs160, h, t, "dec-map-v160")
		testDeepEqualErr(v160v1, v160v2, t, "equal-map-v160")
		bs160 = testMarshalErr(&v160v1, h, t, "enc-map-v160-p")
		v160v2 = nil
		testUnmarshalErr(&v160v2, bs160, h, t, "dec-map-v160-p")
		testDeepEqualErr(v160v1, v160v2, t, "equal-map-v160-p")
	}

	for _, v := range []map[uintptr]interface{}{nil, map[uintptr]interface{}{}, map[uintptr]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v163: %v\n", v)
		var v163v1, v163v2 map[uintptr]interface{}
		v163v1 = v
		bs163 := testMarshalErr(v163v1, h, t, "enc-map-v163")
		if v != nil {
			v163v2 = make(map[uintptr]interface{}, len(v))
		}
		testUnmarshalErr(v163v2, bs163, h, t, "dec-map-v163")
		testDeepEqualErr(v163v1, v163v2, t, "equal-map-v163")
		bs163 = testMarshalErr(&v163v1, h, t, "enc-map-v163-p")
		v163v2 = nil
		testUnmarshalErr(&v163v2, bs163, h, t, "dec-map-v163-p")
		testDeepEqualErr(v163v1, v163v2, t, "equal-map-v163-p")
	}

	for _, v := range []map[uintptr]string{nil, map[uintptr]string{}, map[uintptr]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v164: %v\n", v)
		var v164v1, v164v2 map[uintptr]string
		v164v1 = v
		bs164 := testMarshalErr(v164v1, h, t, "enc-map-v164")
		if v != nil {
			v164v2 = make(map[uintptr]string, len(v))
		}
		testUnmarshalErr(v164v2, bs164, h, t, "dec-map-v164")
		testDeepEqualErr(v164v1, v164v2, t, "equal-map-v164")
		bs164 = testMarshalErr(&v164v1, h, t, "enc-map-v164-p")
		v164v2 = nil
		testUnmarshalErr(&v164v2, bs164, h, t, "dec-map-v164-p")
		testDeepEqualErr(v164v1, v164v2, t, "equal-map-v164-p")
	}

	for _, v := range []map[uintptr]uint{nil, map[uintptr]uint{}, map[uintptr]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v165: %v\n", v)
		var v165v1, v165v2 map[uintptr]uint
		v165v1 = v
		bs165 := testMarshalErr(v165v1, h, t, "enc-map-v165")
		if v != nil {
			v165v2 = make(map[uintptr]uint, len(v))
		}
		testUnmarshalErr(v165v2, bs165, h, t, "dec-map-v165")
		testDeepEqualErr(v165v1, v165v2, t, "equal-map-v165")
		bs165 = testMarshalErr(&v165v1, h, t, "enc-map-v165-p")
		v165v2 = nil
		testUnmarshalErr(&v165v2, bs165, h, t, "dec-map-v165-p")
		testDeepEqualErr(v165v1, v165v2, t, "equal-map-v165-p")
	}

	for _, v := range []map[uintptr]uint8{nil, map[uintptr]uint8{}, map[uintptr]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v166: %v\n", v)
		var v166v1, v166v2 map[uintptr]uint8
		v166v1 = v
		bs166 := testMarshalErr(v166v1, h, t, "enc-map-v166")
		if v != nil {
			v166v2 = make(map[uintptr]uint8, len(v))
		}
		testUnmarshalErr(v166v2, bs166, h, t, "dec-map-v166")
		testDeepEqualErr(v166v1, v166v2, t, "equal-map-v166")
		bs166 = testMarshalErr(&v166v1, h, t, "enc-map-v166-p")
		v166v2 = nil
		testUnmarshalErr(&v166v2, bs166, h, t, "dec-map-v166-p")
		testDeepEqualErr(v166v1, v166v2, t, "equal-map-v166-p")
	}

	for _, v := range []map[uintptr]uint16{nil, map[uintptr]uint16{}, map[uintptr]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v167: %v\n", v)
		var v167v1, v167v2 map[uintptr]uint16
		v167v1 = v
		bs167 := testMarshalErr(v167v1, h, t, "enc-map-v167")
		if v != nil {
			v167v2 = make(map[uintptr]uint16, len(v))
		}
		testUnmarshalErr(v167v2, bs167, h, t, "dec-map-v167")
		testDeepEqualErr(v167v1, v167v2, t, "equal-map-v167")
		bs167 = testMarshalErr(&v167v1, h, t, "enc-map-v167-p")
		v167v2 = nil
		testUnmarshalErr(&v167v2, bs167, h, t, "dec-map-v167-p")
		testDeepEqualErr(v167v1, v167v2, t, "equal-map-v167-p")
	}

	for _, v := range []map[uintptr]uint32{nil, map[uintptr]uint32{}, map[uintptr]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v168: %v\n", v)
		var v168v1, v168v2 map[uintptr]uint32
		v168v1 = v
		bs168 := testMarshalErr(v168v1, h, t, "enc-map-v168")
		if v != nil {
			v168v2 = make(map[uintptr]uint32, len(v))
		}
		testUnmarshalErr(v168v2, bs168, h, t, "dec-map-v168")
		testDeepEqualErr(v168v1, v168v2, t, "equal-map-v168")
		bs168 = testMarshalErr(&v168v1, h, t, "enc-map-v168-p")
		v168v2 = nil
		testUnmarshalErr(&v168v2, bs168, h, t, "dec-map-v168-p")
		testDeepEqualErr(v168v1, v168v2, t, "equal-map-v168-p")
	}

	for _, v := range []map[uintptr]uint64{nil, map[uintptr]uint64{}, map[uintptr]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v169: %v\n", v)
		var v169v1, v169v2 map[uintptr]uint64
		v169v1 = v
		bs169 := testMarshalErr(v169v1, h, t, "enc-map-v169")
		if v != nil {
			v169v2 = make(map[uintptr]uint64, len(v))
		}
		testUnmarshalErr(v169v2, bs169, h, t, "dec-map-v169")
		testDeepEqualErr(v169v1, v169v2, t, "equal-map-v169")
		bs169 = testMarshalErr(&v169v1, h, t, "enc-map-v169-p")
		v169v2 = nil
		testUnmarshalErr(&v169v2, bs169, h, t, "dec-map-v169-p")
		testDeepEqualErr(v169v1, v169v2, t, "equal-map-v169-p")
	}

	for _, v := range []map[uintptr]uintptr{nil, map[uintptr]uintptr{}, map[uintptr]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v170: %v\n", v)
		var v170v1, v170v2 map[uintptr]uintptr
		v170v1 = v
		bs170 := testMarshalErr(v170v1, h, t, "enc-map-v170")
		if v != nil {
			v170v2 = make(map[uintptr]uintptr, len(v))
		}
		testUnmarshalErr(v170v2, bs170, h, t, "dec-map-v170")
		testDeepEqualErr(v170v1, v170v2, t, "equal-map-v170")
		bs170 = testMarshalErr(&v170v1, h, t, "enc-map-v170-p")
		v170v2 = nil
		testUnmarshalErr(&v170v2, bs170, h, t, "dec-map-v170-p")
		testDeepEqualErr(v170v1, v170v2, t, "equal-map-v170-p")
	}

	for _, v := range []map[uintptr]int{nil, map[uintptr]int{}, map[uintptr]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v171: %v\n", v)
		var v171v1, v171v2 map[uintptr]int
		v171v1 = v
		bs171 := testMarshalErr(v171v1, h, t, "enc-map-v171")
		if v != nil {
			v171v2 = make(map[uintptr]int, len(v))
		}
		testUnmarshalErr(v171v2, bs171, h, t, "dec-map-v171")
		testDeepEqualErr(v171v1, v171v2, t, "equal-map-v171")
		bs171 = testMarshalErr(&v171v1, h, t, "enc-map-v171-p")
		v171v2 = nil
		testUnmarshalErr(&v171v2, bs171, h, t, "dec-map-v171-p")
		testDeepEqualErr(v171v1, v171v2, t, "equal-map-v171-p")
	}

	for _, v := range []map[uintptr]int8{nil, map[uintptr]int8{}, map[uintptr]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v172: %v\n", v)
		var v172v1, v172v2 map[uintptr]int8
		v172v1 = v
		bs172 := testMarshalErr(v172v1, h, t, "enc-map-v172")
		if v != nil {
			v172v2 = make(map[uintptr]int8, len(v))
		}
		testUnmarshalErr(v172v2, bs172, h, t, "dec-map-v172")
		testDeepEqualErr(v172v1, v172v2, t, "equal-map-v172")
		bs172 = testMarshalErr(&v172v1, h, t, "enc-map-v172-p")
		v172v2 = nil
		testUnmarshalErr(&v172v2, bs172, h, t, "dec-map-v172-p")
		testDeepEqualErr(v172v1, v172v2, t, "equal-map-v172-p")
	}

	for _, v := range []map[uintptr]int16{nil, map[uintptr]int16{}, map[uintptr]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v173: %v\n", v)
		var v173v1, v173v2 map[uintptr]int16
		v173v1 = v
		bs173 := testMarshalErr(v173v1, h, t, "enc-map-v173")
		if v != nil {
			v173v2 = make(map[uintptr]int16, len(v))
		}
		testUnmarshalErr(v173v2, bs173, h, t, "dec-map-v173")
		testDeepEqualErr(v173v1, v173v2, t, "equal-map-v173")
		bs173 = testMarshalErr(&v173v1, h, t, "enc-map-v173-p")
		v173v2 = nil
		testUnmarshalErr(&v173v2, bs173, h, t, "dec-map-v173-p")
		testDeepEqualErr(v173v1, v173v2, t, "equal-map-v173-p")
	}

	for _, v := range []map[uintptr]int32{nil, map[uintptr]int32{}, map[uintptr]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v174: %v\n", v)
		var v174v1, v174v2 map[uintptr]int32
		v174v1 = v
		bs174 := testMarshalErr(v174v1, h, t, "enc-map-v174")
		if v != nil {
			v174v2 = make(map[uintptr]int32, len(v))
		}
		testUnmarshalErr(v174v2, bs174, h, t, "dec-map-v174")
		testDeepEqualErr(v174v1, v174v2, t, "equal-map-v174")
		bs174 = testMarshalErr(&v174v1, h, t, "enc-map-v174-p")
		v174v2 = nil
		testUnmarshalErr(&v174v2, bs174, h, t, "dec-map-v174-p")
		testDeepEqualErr(v174v1, v174v2, t, "equal-map-v174-p")
	}

	for _, v := range []map[uintptr]int64{nil, map[uintptr]int64{}, map[uintptr]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v175: %v\n", v)
		var v175v1, v175v2 map[uintptr]int64
		v175v1 = v
		bs175 := testMarshalErr(v175v1, h, t, "enc-map-v175")
		if v != nil {
			v175v2 = make(map[uintptr]int64, len(v))
		}
		testUnmarshalErr(v175v2, bs175, h, t, "dec-map-v175")
		testDeepEqualErr(v175v1, v175v2, t, "equal-map-v175")
		bs175 = testMarshalErr(&v175v1, h, t, "enc-map-v175-p")
		v175v2 = nil
		testUnmarshalErr(&v175v2, bs175, h, t, "dec-map-v175-p")
		testDeepEqualErr(v175v1, v175v2, t, "equal-map-v175-p")
	}

	for _, v := range []map[uintptr]float32{nil, map[uintptr]float32{}, map[uintptr]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v176: %v\n", v)
		var v176v1, v176v2 map[uintptr]float32
		v176v1 = v
		bs176 := testMarshalErr(v176v1, h, t, "enc-map-v176")
		if v != nil {
			v176v2 = make(map[uintptr]float32, len(v))
		}
		testUnmarshalErr(v176v2, bs176, h, t, "dec-map-v176")
		testDeepEqualErr(v176v1, v176v2, t, "equal-map-v176")
		bs176 = testMarshalErr(&v176v1, h, t, "enc-map-v176-p")
		v176v2 = nil
		testUnmarshalErr(&v176v2, bs176, h, t, "dec-map-v176-p")
		testDeepEqualErr(v176v1, v176v2, t, "equal-map-v176-p")
	}

	for _, v := range []map[uintptr]float64{nil, map[uintptr]float64{}, map[uintptr]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v177: %v\n", v)
		var v177v1, v177v2 map[uintptr]float64
		v177v1 = v
		bs177 := testMarshalErr(v177v1, h, t, "enc-map-v177")
		if v != nil {
			v177v2 = make(map[uintptr]float64, len(v))
		}
		testUnmarshalErr(v177v2, bs177, h, t, "dec-map-v177")
		testDeepEqualErr(v177v1, v177v2, t, "equal-map-v177")
		bs177 = testMarshalErr(&v177v1, h, t, "enc-map-v177-p")
		v177v2 = nil
		testUnmarshalErr(&v177v2, bs177, h, t, "dec-map-v177-p")
		testDeepEqualErr(v177v1, v177v2, t, "equal-map-v177-p")
	}

	for _, v := range []map[uintptr]bool{nil, map[uintptr]bool{}, map[uintptr]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v178: %v\n", v)
		var v178v1, v178v2 map[uintptr]bool
		v178v1 = v
		bs178 := testMarshalErr(v178v1, h, t, "enc-map-v178")
		if v != nil {
			v178v2 = make(map[uintptr]bool, len(v))
		}
		testUnmarshalErr(v178v2, bs178, h, t, "dec-map-v178")
		testDeepEqualErr(v178v1, v178v2, t, "equal-map-v178")
		bs178 = testMarshalErr(&v178v1, h, t, "enc-map-v178-p")
		v178v2 = nil
		testUnmarshalErr(&v178v2, bs178, h, t, "dec-map-v178-p")
		testDeepEqualErr(v178v1, v178v2, t, "equal-map-v178-p")
	}

	for _, v := range []map[int]interface{}{nil, map[int]interface{}{}, map[int]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v181: %v\n", v)
		var v181v1, v181v2 map[int]interface{}
		v181v1 = v
		bs181 := testMarshalErr(v181v1, h, t, "enc-map-v181")
		if v != nil {
			v181v2 = make(map[int]interface{}, len(v))
		}
		testUnmarshalErr(v181v2, bs181, h, t, "dec-map-v181")
		testDeepEqualErr(v181v1, v181v2, t, "equal-map-v181")
		bs181 = testMarshalErr(&v181v1, h, t, "enc-map-v181-p")
		v181v2 = nil
		testUnmarshalErr(&v181v2, bs181, h, t, "dec-map-v181-p")
		testDeepEqualErr(v181v1, v181v2, t, "equal-map-v181-p")
	}

	for _, v := range []map[int]string{nil, map[int]string{}, map[int]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v182: %v\n", v)
		var v182v1, v182v2 map[int]string
		v182v1 = v
		bs182 := testMarshalErr(v182v1, h, t, "enc-map-v182")
		if v != nil {
			v182v2 = make(map[int]string, len(v))
		}
		testUnmarshalErr(v182v2, bs182, h, t, "dec-map-v182")
		testDeepEqualErr(v182v1, v182v2, t, "equal-map-v182")
		bs182 = testMarshalErr(&v182v1, h, t, "enc-map-v182-p")
		v182v2 = nil
		testUnmarshalErr(&v182v2, bs182, h, t, "dec-map-v182-p")
		testDeepEqualErr(v182v1, v182v2, t, "equal-map-v182-p")
	}

	for _, v := range []map[int]uint{nil, map[int]uint{}, map[int]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v183: %v\n", v)
		var v183v1, v183v2 map[int]uint
		v183v1 = v
		bs183 := testMarshalErr(v183v1, h, t, "enc-map-v183")
		if v != nil {
			v183v2 = make(map[int]uint, len(v))
		}
		testUnmarshalErr(v183v2, bs183, h, t, "dec-map-v183")
		testDeepEqualErr(v183v1, v183v2, t, "equal-map-v183")
		bs183 = testMarshalErr(&v183v1, h, t, "enc-map-v183-p")
		v183v2 = nil
		testUnmarshalErr(&v183v2, bs183, h, t, "dec-map-v183-p")
		testDeepEqualErr(v183v1, v183v2, t, "equal-map-v183-p")
	}

	for _, v := range []map[int]uint8{nil, map[int]uint8{}, map[int]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v184: %v\n", v)
		var v184v1, v184v2 map[int]uint8
		v184v1 = v
		bs184 := testMarshalErr(v184v1, h, t, "enc-map-v184")
		if v != nil {
			v184v2 = make(map[int]uint8, len(v))
		}
		testUnmarshalErr(v184v2, bs184, h, t, "dec-map-v184")
		testDeepEqualErr(v184v1, v184v2, t, "equal-map-v184")
		bs184 = testMarshalErr(&v184v1, h, t, "enc-map-v184-p")
		v184v2 = nil
		testUnmarshalErr(&v184v2, bs184, h, t, "dec-map-v184-p")
		testDeepEqualErr(v184v1, v184v2, t, "equal-map-v184-p")
	}

	for _, v := range []map[int]uint16{nil, map[int]uint16{}, map[int]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v185: %v\n", v)
		var v185v1, v185v2 map[int]uint16
		v185v1 = v
		bs185 := testMarshalErr(v185v1, h, t, "enc-map-v185")
		if v != nil {
			v185v2 = make(map[int]uint16, len(v))
		}
		testUnmarshalErr(v185v2, bs185, h, t, "dec-map-v185")
		testDeepEqualErr(v185v1, v185v2, t, "equal-map-v185")
		bs185 = testMarshalErr(&v185v1, h, t, "enc-map-v185-p")
		v185v2 = nil
		testUnmarshalErr(&v185v2, bs185, h, t, "dec-map-v185-p")
		testDeepEqualErr(v185v1, v185v2, t, "equal-map-v185-p")
	}

	for _, v := range []map[int]uint32{nil, map[int]uint32{}, map[int]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v186: %v\n", v)
		var v186v1, v186v2 map[int]uint32
		v186v1 = v
		bs186 := testMarshalErr(v186v1, h, t, "enc-map-v186")
		if v != nil {
			v186v2 = make(map[int]uint32, len(v))
		}
		testUnmarshalErr(v186v2, bs186, h, t, "dec-map-v186")
		testDeepEqualErr(v186v1, v186v2, t, "equal-map-v186")
		bs186 = testMarshalErr(&v186v1, h, t, "enc-map-v186-p")
		v186v2 = nil
		testUnmarshalErr(&v186v2, bs186, h, t, "dec-map-v186-p")
		testDeepEqualErr(v186v1, v186v2, t, "equal-map-v186-p")
	}

	for _, v := range []map[int]uint64{nil, map[int]uint64{}, map[int]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v187: %v\n", v)
		var v187v1, v187v2 map[int]uint64
		v187v1 = v
		bs187 := testMarshalErr(v187v1, h, t, "enc-map-v187")
		if v != nil {
			v187v2 = make(map[int]uint64, len(v))
		}
		testUnmarshalErr(v187v2, bs187, h, t, "dec-map-v187")
		testDeepEqualErr(v187v1, v187v2, t, "equal-map-v187")
		bs187 = testMarshalErr(&v187v1, h, t, "enc-map-v187-p")
		v187v2 = nil
		testUnmarshalErr(&v187v2, bs187, h, t, "dec-map-v187-p")
		testDeepEqualErr(v187v1, v187v2, t, "equal-map-v187-p")
	}

	for _, v := range []map[int]uintptr{nil, map[int]uintptr{}, map[int]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v188: %v\n", v)
		var v188v1, v188v2 map[int]uintptr
		v188v1 = v
		bs188 := testMarshalErr(v188v1, h, t, "enc-map-v188")
		if v != nil {
			v188v2 = make(map[int]uintptr, len(v))
		}
		testUnmarshalErr(v188v2, bs188, h, t, "dec-map-v188")
		testDeepEqualErr(v188v1, v188v2, t, "equal-map-v188")
		bs188 = testMarshalErr(&v188v1, h, t, "enc-map-v188-p")
		v188v2 = nil
		testUnmarshalErr(&v188v2, bs188, h, t, "dec-map-v188-p")
		testDeepEqualErr(v188v1, v188v2, t, "equal-map-v188-p")
	}

	for _, v := range []map[int]int{nil, map[int]int{}, map[int]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v189: %v\n", v)
		var v189v1, v189v2 map[int]int
		v189v1 = v
		bs189 := testMarshalErr(v189v1, h, t, "enc-map-v189")
		if v != nil {
			v189v2 = make(map[int]int, len(v))
		}
		testUnmarshalErr(v189v2, bs189, h, t, "dec-map-v189")
		testDeepEqualErr(v189v1, v189v2, t, "equal-map-v189")
		bs189 = testMarshalErr(&v189v1, h, t, "enc-map-v189-p")
		v189v2 = nil
		testUnmarshalErr(&v189v2, bs189, h, t, "dec-map-v189-p")
		testDeepEqualErr(v189v1, v189v2, t, "equal-map-v189-p")
	}

	for _, v := range []map[int]int8{nil, map[int]int8{}, map[int]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v190: %v\n", v)
		var v190v1, v190v2 map[int]int8
		v190v1 = v
		bs190 := testMarshalErr(v190v1, h, t, "enc-map-v190")
		if v != nil {
			v190v2 = make(map[int]int8, len(v))
		}
		testUnmarshalErr(v190v2, bs190, h, t, "dec-map-v190")
		testDeepEqualErr(v190v1, v190v2, t, "equal-map-v190")
		bs190 = testMarshalErr(&v190v1, h, t, "enc-map-v190-p")
		v190v2 = nil
		testUnmarshalErr(&v190v2, bs190, h, t, "dec-map-v190-p")
		testDeepEqualErr(v190v1, v190v2, t, "equal-map-v190-p")
	}

	for _, v := range []map[int]int16{nil, map[int]int16{}, map[int]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v191: %v\n", v)
		var v191v1, v191v2 map[int]int16
		v191v1 = v
		bs191 := testMarshalErr(v191v1, h, t, "enc-map-v191")
		if v != nil {
			v191v2 = make(map[int]int16, len(v))
		}
		testUnmarshalErr(v191v2, bs191, h, t, "dec-map-v191")
		testDeepEqualErr(v191v1, v191v2, t, "equal-map-v191")
		bs191 = testMarshalErr(&v191v1, h, t, "enc-map-v191-p")
		v191v2 = nil
		testUnmarshalErr(&v191v2, bs191, h, t, "dec-map-v191-p")
		testDeepEqualErr(v191v1, v191v2, t, "equal-map-v191-p")
	}

	for _, v := range []map[int]int32{nil, map[int]int32{}, map[int]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v192: %v\n", v)
		var v192v1, v192v2 map[int]int32
		v192v1 = v
		bs192 := testMarshalErr(v192v1, h, t, "enc-map-v192")
		if v != nil {
			v192v2 = make(map[int]int32, len(v))
		}
		testUnmarshalErr(v192v2, bs192, h, t, "dec-map-v192")
		testDeepEqualErr(v192v1, v192v2, t, "equal-map-v192")
		bs192 = testMarshalErr(&v192v1, h, t, "enc-map-v192-p")
		v192v2 = nil
		testUnmarshalErr(&v192v2, bs192, h, t, "dec-map-v192-p")
		testDeepEqualErr(v192v1, v192v2, t, "equal-map-v192-p")
	}

	for _, v := range []map[int]int64{nil, map[int]int64{}, map[int]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v193: %v\n", v)
		var v193v1, v193v2 map[int]int64
		v193v1 = v
		bs193 := testMarshalErr(v193v1, h, t, "enc-map-v193")
		if v != nil {
			v193v2 = make(map[int]int64, len(v))
		}
		testUnmarshalErr(v193v2, bs193, h, t, "dec-map-v193")
		testDeepEqualErr(v193v1, v193v2, t, "equal-map-v193")
		bs193 = testMarshalErr(&v193v1, h, t, "enc-map-v193-p")
		v193v2 = nil
		testUnmarshalErr(&v193v2, bs193, h, t, "dec-map-v193-p")
		testDeepEqualErr(v193v1, v193v2, t, "equal-map-v193-p")
	}

	for _, v := range []map[int]float32{nil, map[int]float32{}, map[int]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v194: %v\n", v)
		var v194v1, v194v2 map[int]float32
		v194v1 = v
		bs194 := testMarshalErr(v194v1, h, t, "enc-map-v194")
		if v != nil {
			v194v2 = make(map[int]float32, len(v))
		}
		testUnmarshalErr(v194v2, bs194, h, t, "dec-map-v194")
		testDeepEqualErr(v194v1, v194v2, t, "equal-map-v194")
		bs194 = testMarshalErr(&v194v1, h, t, "enc-map-v194-p")
		v194v2 = nil
		testUnmarshalErr(&v194v2, bs194, h, t, "dec-map-v194-p")
		testDeepEqualErr(v194v1, v194v2, t, "equal-map-v194-p")
	}

	for _, v := range []map[int]float64{nil, map[int]float64{}, map[int]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v195: %v\n", v)
		var v195v1, v195v2 map[int]float64
		v195v1 = v
		bs195 := testMarshalErr(v195v1, h, t, "enc-map-v195")
		if v != nil {
			v195v2 = make(map[int]float64, len(v))
		}
		testUnmarshalErr(v195v2, bs195, h, t, "dec-map-v195")
		testDeepEqualErr(v195v1, v195v2, t, "equal-map-v195")
		bs195 = testMarshalErr(&v195v1, h, t, "enc-map-v195-p")
		v195v2 = nil
		testUnmarshalErr(&v195v2, bs195, h, t, "dec-map-v195-p")
		testDeepEqualErr(v195v1, v195v2, t, "equal-map-v195-p")
	}

	for _, v := range []map[int]bool{nil, map[int]bool{}, map[int]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v196: %v\n", v)
		var v196v1, v196v2 map[int]bool
		v196v1 = v
		bs196 := testMarshalErr(v196v1, h, t, "enc-map-v196")
		if v != nil {
			v196v2 = make(map[int]bool, len(v))
		}
		testUnmarshalErr(v196v2, bs196, h, t, "dec-map-v196")
		testDeepEqualErr(v196v1, v196v2, t, "equal-map-v196")
		bs196 = testMarshalErr(&v196v1, h, t, "enc-map-v196-p")
		v196v2 = nil
		testUnmarshalErr(&v196v2, bs196, h, t, "dec-map-v196-p")
		testDeepEqualErr(v196v1, v196v2, t, "equal-map-v196-p")
	}

	for _, v := range []map[int8]interface{}{nil, map[int8]interface{}{}, map[int8]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v199: %v\n", v)
		var v199v1, v199v2 map[int8]interface{}
		v199v1 = v
		bs199 := testMarshalErr(v199v1, h, t, "enc-map-v199")
		if v != nil {
			v199v2 = make(map[int8]interface{}, len(v))
		}
		testUnmarshalErr(v199v2, bs199, h, t, "dec-map-v199")
		testDeepEqualErr(v199v1, v199v2, t, "equal-map-v199")
		bs199 = testMarshalErr(&v199v1, h, t, "enc-map-v199-p")
		v199v2 = nil
		testUnmarshalErr(&v199v2, bs199, h, t, "dec-map-v199-p")
		testDeepEqualErr(v199v1, v199v2, t, "equal-map-v199-p")
	}

	for _, v := range []map[int8]string{nil, map[int8]string{}, map[int8]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v200: %v\n", v)
		var v200v1, v200v2 map[int8]string
		v200v1 = v
		bs200 := testMarshalErr(v200v1, h, t, "enc-map-v200")
		if v != nil {
			v200v2 = make(map[int8]string, len(v))
		}
		testUnmarshalErr(v200v2, bs200, h, t, "dec-map-v200")
		testDeepEqualErr(v200v1, v200v2, t, "equal-map-v200")
		bs200 = testMarshalErr(&v200v1, h, t, "enc-map-v200-p")
		v200v2 = nil
		testUnmarshalErr(&v200v2, bs200, h, t, "dec-map-v200-p")
		testDeepEqualErr(v200v1, v200v2, t, "equal-map-v200-p")
	}

	for _, v := range []map[int8]uint{nil, map[int8]uint{}, map[int8]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v201: %v\n", v)
		var v201v1, v201v2 map[int8]uint
		v201v1 = v
		bs201 := testMarshalErr(v201v1, h, t, "enc-map-v201")
		if v != nil {
			v201v2 = make(map[int8]uint, len(v))
		}
		testUnmarshalErr(v201v2, bs201, h, t, "dec-map-v201")
		testDeepEqualErr(v201v1, v201v2, t, "equal-map-v201")
		bs201 = testMarshalErr(&v201v1, h, t, "enc-map-v201-p")
		v201v2 = nil
		testUnmarshalErr(&v201v2, bs201, h, t, "dec-map-v201-p")
		testDeepEqualErr(v201v1, v201v2, t, "equal-map-v201-p")
	}

	for _, v := range []map[int8]uint8{nil, map[int8]uint8{}, map[int8]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v202: %v\n", v)
		var v202v1, v202v2 map[int8]uint8
		v202v1 = v
		bs202 := testMarshalErr(v202v1, h, t, "enc-map-v202")
		if v != nil {
			v202v2 = make(map[int8]uint8, len(v))
		}
		testUnmarshalErr(v202v2, bs202, h, t, "dec-map-v202")
		testDeepEqualErr(v202v1, v202v2, t, "equal-map-v202")
		bs202 = testMarshalErr(&v202v1, h, t, "enc-map-v202-p")
		v202v2 = nil
		testUnmarshalErr(&v202v2, bs202, h, t, "dec-map-v202-p")
		testDeepEqualErr(v202v1, v202v2, t, "equal-map-v202-p")
	}

	for _, v := range []map[int8]uint16{nil, map[int8]uint16{}, map[int8]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v203: %v\n", v)
		var v203v1, v203v2 map[int8]uint16
		v203v1 = v
		bs203 := testMarshalErr(v203v1, h, t, "enc-map-v203")
		if v != nil {
			v203v2 = make(map[int8]uint16, len(v))
		}
		testUnmarshalErr(v203v2, bs203, h, t, "dec-map-v203")
		testDeepEqualErr(v203v1, v203v2, t, "equal-map-v203")
		bs203 = testMarshalErr(&v203v1, h, t, "enc-map-v203-p")
		v203v2 = nil
		testUnmarshalErr(&v203v2, bs203, h, t, "dec-map-v203-p")
		testDeepEqualErr(v203v1, v203v2, t, "equal-map-v203-p")
	}

	for _, v := range []map[int8]uint32{nil, map[int8]uint32{}, map[int8]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v204: %v\n", v)
		var v204v1, v204v2 map[int8]uint32
		v204v1 = v
		bs204 := testMarshalErr(v204v1, h, t, "enc-map-v204")
		if v != nil {
			v204v2 = make(map[int8]uint32, len(v))
		}
		testUnmarshalErr(v204v2, bs204, h, t, "dec-map-v204")
		testDeepEqualErr(v204v1, v204v2, t, "equal-map-v204")
		bs204 = testMarshalErr(&v204v1, h, t, "enc-map-v204-p")
		v204v2 = nil
		testUnmarshalErr(&v204v2, bs204, h, t, "dec-map-v204-p")
		testDeepEqualErr(v204v1, v204v2, t, "equal-map-v204-p")
	}

	for _, v := range []map[int8]uint64{nil, map[int8]uint64{}, map[int8]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v205: %v\n", v)
		var v205v1, v205v2 map[int8]uint64
		v205v1 = v
		bs205 := testMarshalErr(v205v1, h, t, "enc-map-v205")
		if v != nil {
			v205v2 = make(map[int8]uint64, len(v))
		}
		testUnmarshalErr(v205v2, bs205, h, t, "dec-map-v205")
		testDeepEqualErr(v205v1, v205v2, t, "equal-map-v205")
		bs205 = testMarshalErr(&v205v1, h, t, "enc-map-v205-p")
		v205v2 = nil
		testUnmarshalErr(&v205v2, bs205, h, t, "dec-map-v205-p")
		testDeepEqualErr(v205v1, v205v2, t, "equal-map-v205-p")
	}

	for _, v := range []map[int8]uintptr{nil, map[int8]uintptr{}, map[int8]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v206: %v\n", v)
		var v206v1, v206v2 map[int8]uintptr
		v206v1 = v
		bs206 := testMarshalErr(v206v1, h, t, "enc-map-v206")
		if v != nil {
			v206v2 = make(map[int8]uintptr, len(v))
		}
		testUnmarshalErr(v206v2, bs206, h, t, "dec-map-v206")
		testDeepEqualErr(v206v1, v206v2, t, "equal-map-v206")
		bs206 = testMarshalErr(&v206v1, h, t, "enc-map-v206-p")
		v206v2 = nil
		testUnmarshalErr(&v206v2, bs206, h, t, "dec-map-v206-p")
		testDeepEqualErr(v206v1, v206v2, t, "equal-map-v206-p")
	}

	for _, v := range []map[int8]int{nil, map[int8]int{}, map[int8]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v207: %v\n", v)
		var v207v1, v207v2 map[int8]int
		v207v1 = v
		bs207 := testMarshalErr(v207v1, h, t, "enc-map-v207")
		if v != nil {
			v207v2 = make(map[int8]int, len(v))
		}
		testUnmarshalErr(v207v2, bs207, h, t, "dec-map-v207")
		testDeepEqualErr(v207v1, v207v2, t, "equal-map-v207")
		bs207 = testMarshalErr(&v207v1, h, t, "enc-map-v207-p")
		v207v2 = nil
		testUnmarshalErr(&v207v2, bs207, h, t, "dec-map-v207-p")
		testDeepEqualErr(v207v1, v207v2, t, "equal-map-v207-p")
	}

	for _, v := range []map[int8]int8{nil, map[int8]int8{}, map[int8]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v208: %v\n", v)
		var v208v1, v208v2 map[int8]int8
		v208v1 = v
		bs208 := testMarshalErr(v208v1, h, t, "enc-map-v208")
		if v != nil {
			v208v2 = make(map[int8]int8, len(v))
		}
		testUnmarshalErr(v208v2, bs208, h, t, "dec-map-v208")
		testDeepEqualErr(v208v1, v208v2, t, "equal-map-v208")
		bs208 = testMarshalErr(&v208v1, h, t, "enc-map-v208-p")
		v208v2 = nil
		testUnmarshalErr(&v208v2, bs208, h, t, "dec-map-v208-p")
		testDeepEqualErr(v208v1, v208v2, t, "equal-map-v208-p")
	}

	for _, v := range []map[int8]int16{nil, map[int8]int16{}, map[int8]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v209: %v\n", v)
		var v209v1, v209v2 map[int8]int16
		v209v1 = v
		bs209 := testMarshalErr(v209v1, h, t, "enc-map-v209")
		if v != nil {
			v209v2 = make(map[int8]int16, len(v))
		}
		testUnmarshalErr(v209v2, bs209, h, t, "dec-map-v209")
		testDeepEqualErr(v209v1, v209v2, t, "equal-map-v209")
		bs209 = testMarshalErr(&v209v1, h, t, "enc-map-v209-p")
		v209v2 = nil
		testUnmarshalErr(&v209v2, bs209, h, t, "dec-map-v209-p")
		testDeepEqualErr(v209v1, v209v2, t, "equal-map-v209-p")
	}

	for _, v := range []map[int8]int32{nil, map[int8]int32{}, map[int8]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v210: %v\n", v)
		var v210v1, v210v2 map[int8]int32
		v210v1 = v
		bs210 := testMarshalErr(v210v1, h, t, "enc-map-v210")
		if v != nil {
			v210v2 = make(map[int8]int32, len(v))
		}
		testUnmarshalErr(v210v2, bs210, h, t, "dec-map-v210")
		testDeepEqualErr(v210v1, v210v2, t, "equal-map-v210")
		bs210 = testMarshalErr(&v210v1, h, t, "enc-map-v210-p")
		v210v2 = nil
		testUnmarshalErr(&v210v2, bs210, h, t, "dec-map-v210-p")
		testDeepEqualErr(v210v1, v210v2, t, "equal-map-v210-p")
	}

	for _, v := range []map[int8]int64{nil, map[int8]int64{}, map[int8]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v211: %v\n", v)
		var v211v1, v211v2 map[int8]int64
		v211v1 = v
		bs211 := testMarshalErr(v211v1, h, t, "enc-map-v211")
		if v != nil {
			v211v2 = make(map[int8]int64, len(v))
		}
		testUnmarshalErr(v211v2, bs211, h, t, "dec-map-v211")
		testDeepEqualErr(v211v1, v211v2, t, "equal-map-v211")
		bs211 = testMarshalErr(&v211v1, h, t, "enc-map-v211-p")
		v211v2 = nil
		testUnmarshalErr(&v211v2, bs211, h, t, "dec-map-v211-p")
		testDeepEqualErr(v211v1, v211v2, t, "equal-map-v211-p")
	}

	for _, v := range []map[int8]float32{nil, map[int8]float32{}, map[int8]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v212: %v\n", v)
		var v212v1, v212v2 map[int8]float32
		v212v1 = v
		bs212 := testMarshalErr(v212v1, h, t, "enc-map-v212")
		if v != nil {
			v212v2 = make(map[int8]float32, len(v))
		}
		testUnmarshalErr(v212v2, bs212, h, t, "dec-map-v212")
		testDeepEqualErr(v212v1, v212v2, t, "equal-map-v212")
		bs212 = testMarshalErr(&v212v1, h, t, "enc-map-v212-p")
		v212v2 = nil
		testUnmarshalErr(&v212v2, bs212, h, t, "dec-map-v212-p")
		testDeepEqualErr(v212v1, v212v2, t, "equal-map-v212-p")
	}

	for _, v := range []map[int8]float64{nil, map[int8]float64{}, map[int8]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v213: %v\n", v)
		var v213v1, v213v2 map[int8]float64
		v213v1 = v
		bs213 := testMarshalErr(v213v1, h, t, "enc-map-v213")
		if v != nil {
			v213v2 = make(map[int8]float64, len(v))
		}
		testUnmarshalErr(v213v2, bs213, h, t, "dec-map-v213")
		testDeepEqualErr(v213v1, v213v2, t, "equal-map-v213")
		bs213 = testMarshalErr(&v213v1, h, t, "enc-map-v213-p")
		v213v2 = nil
		testUnmarshalErr(&v213v2, bs213, h, t, "dec-map-v213-p")
		testDeepEqualErr(v213v1, v213v2, t, "equal-map-v213-p")
	}

	for _, v := range []map[int8]bool{nil, map[int8]bool{}, map[int8]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v214: %v\n", v)
		var v214v1, v214v2 map[int8]bool
		v214v1 = v
		bs214 := testMarshalErr(v214v1, h, t, "enc-map-v214")
		if v != nil {
			v214v2 = make(map[int8]bool, len(v))
		}
		testUnmarshalErr(v214v2, bs214, h, t, "dec-map-v214")
		testDeepEqualErr(v214v1, v214v2, t, "equal-map-v214")
		bs214 = testMarshalErr(&v214v1, h, t, "enc-map-v214-p")
		v214v2 = nil
		testUnmarshalErr(&v214v2, bs214, h, t, "dec-map-v214-p")
		testDeepEqualErr(v214v1, v214v2, t, "equal-map-v214-p")
	}

	for _, v := range []map[int16]interface{}{nil, map[int16]interface{}{}, map[int16]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v217: %v\n", v)
		var v217v1, v217v2 map[int16]interface{}
		v217v1 = v
		bs217 := testMarshalErr(v217v1, h, t, "enc-map-v217")
		if v != nil {
			v217v2 = make(map[int16]interface{}, len(v))
		}
		testUnmarshalErr(v217v2, bs217, h, t, "dec-map-v217")
		testDeepEqualErr(v217v1, v217v2, t, "equal-map-v217")
		bs217 = testMarshalErr(&v217v1, h, t, "enc-map-v217-p")
		v217v2 = nil
		testUnmarshalErr(&v217v2, bs217, h, t, "dec-map-v217-p")
		testDeepEqualErr(v217v1, v217v2, t, "equal-map-v217-p")
	}

	for _, v := range []map[int16]string{nil, map[int16]string{}, map[int16]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v218: %v\n", v)
		var v218v1, v218v2 map[int16]string
		v218v1 = v
		bs218 := testMarshalErr(v218v1, h, t, "enc-map-v218")
		if v != nil {
			v218v2 = make(map[int16]string, len(v))
		}
		testUnmarshalErr(v218v2, bs218, h, t, "dec-map-v218")
		testDeepEqualErr(v218v1, v218v2, t, "equal-map-v218")
		bs218 = testMarshalErr(&v218v1, h, t, "enc-map-v218-p")
		v218v2 = nil
		testUnmarshalErr(&v218v2, bs218, h, t, "dec-map-v218-p")
		testDeepEqualErr(v218v1, v218v2, t, "equal-map-v218-p")
	}

	for _, v := range []map[int16]uint{nil, map[int16]uint{}, map[int16]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v219: %v\n", v)
		var v219v1, v219v2 map[int16]uint
		v219v1 = v
		bs219 := testMarshalErr(v219v1, h, t, "enc-map-v219")
		if v != nil {
			v219v2 = make(map[int16]uint, len(v))
		}
		testUnmarshalErr(v219v2, bs219, h, t, "dec-map-v219")
		testDeepEqualErr(v219v1, v219v2, t, "equal-map-v219")
		bs219 = testMarshalErr(&v219v1, h, t, "enc-map-v219-p")
		v219v2 = nil
		testUnmarshalErr(&v219v2, bs219, h, t, "dec-map-v219-p")
		testDeepEqualErr(v219v1, v219v2, t, "equal-map-v219-p")
	}

	for _, v := range []map[int16]uint8{nil, map[int16]uint8{}, map[int16]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v220: %v\n", v)
		var v220v1, v220v2 map[int16]uint8
		v220v1 = v
		bs220 := testMarshalErr(v220v1, h, t, "enc-map-v220")
		if v != nil {
			v220v2 = make(map[int16]uint8, len(v))
		}
		testUnmarshalErr(v220v2, bs220, h, t, "dec-map-v220")
		testDeepEqualErr(v220v1, v220v2, t, "equal-map-v220")
		bs220 = testMarshalErr(&v220v1, h, t, "enc-map-v220-p")
		v220v2 = nil
		testUnmarshalErr(&v220v2, bs220, h, t, "dec-map-v220-p")
		testDeepEqualErr(v220v1, v220v2, t, "equal-map-v220-p")
	}

	for _, v := range []map[int16]uint16{nil, map[int16]uint16{}, map[int16]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v221: %v\n", v)
		var v221v1, v221v2 map[int16]uint16
		v221v1 = v
		bs221 := testMarshalErr(v221v1, h, t, "enc-map-v221")
		if v != nil {
			v221v2 = make(map[int16]uint16, len(v))
		}
		testUnmarshalErr(v221v2, bs221, h, t, "dec-map-v221")
		testDeepEqualErr(v221v1, v221v2, t, "equal-map-v221")
		bs221 = testMarshalErr(&v221v1, h, t, "enc-map-v221-p")
		v221v2 = nil
		testUnmarshalErr(&v221v2, bs221, h, t, "dec-map-v221-p")
		testDeepEqualErr(v221v1, v221v2, t, "equal-map-v221-p")
	}

	for _, v := range []map[int16]uint32{nil, map[int16]uint32{}, map[int16]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v222: %v\n", v)
		var v222v1, v222v2 map[int16]uint32
		v222v1 = v
		bs222 := testMarshalErr(v222v1, h, t, "enc-map-v222")
		if v != nil {
			v222v2 = make(map[int16]uint32, len(v))
		}
		testUnmarshalErr(v222v2, bs222, h, t, "dec-map-v222")
		testDeepEqualErr(v222v1, v222v2, t, "equal-map-v222")
		bs222 = testMarshalErr(&v222v1, h, t, "enc-map-v222-p")
		v222v2 = nil
		testUnmarshalErr(&v222v2, bs222, h, t, "dec-map-v222-p")
		testDeepEqualErr(v222v1, v222v2, t, "equal-map-v222-p")
	}

	for _, v := range []map[int16]uint64{nil, map[int16]uint64{}, map[int16]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v223: %v\n", v)
		var v223v1, v223v2 map[int16]uint64
		v223v1 = v
		bs223 := testMarshalErr(v223v1, h, t, "enc-map-v223")
		if v != nil {
			v223v2 = make(map[int16]uint64, len(v))
		}
		testUnmarshalErr(v223v2, bs223, h, t, "dec-map-v223")
		testDeepEqualErr(v223v1, v223v2, t, "equal-map-v223")
		bs223 = testMarshalErr(&v223v1, h, t, "enc-map-v223-p")
		v223v2 = nil
		testUnmarshalErr(&v223v2, bs223, h, t, "dec-map-v223-p")
		testDeepEqualErr(v223v1, v223v2, t, "equal-map-v223-p")
	}

	for _, v := range []map[int16]uintptr{nil, map[int16]uintptr{}, map[int16]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v224: %v\n", v)
		var v224v1, v224v2 map[int16]uintptr
		v224v1 = v
		bs224 := testMarshalErr(v224v1, h, t, "enc-map-v224")
		if v != nil {
			v224v2 = make(map[int16]uintptr, len(v))
		}
		testUnmarshalErr(v224v2, bs224, h, t, "dec-map-v224")
		testDeepEqualErr(v224v1, v224v2, t, "equal-map-v224")
		bs224 = testMarshalErr(&v224v1, h, t, "enc-map-v224-p")
		v224v2 = nil
		testUnmarshalErr(&v224v2, bs224, h, t, "dec-map-v224-p")
		testDeepEqualErr(v224v1, v224v2, t, "equal-map-v224-p")
	}

	for _, v := range []map[int16]int{nil, map[int16]int{}, map[int16]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v225: %v\n", v)
		var v225v1, v225v2 map[int16]int
		v225v1 = v
		bs225 := testMarshalErr(v225v1, h, t, "enc-map-v225")
		if v != nil {
			v225v2 = make(map[int16]int, len(v))
		}
		testUnmarshalErr(v225v2, bs225, h, t, "dec-map-v225")
		testDeepEqualErr(v225v1, v225v2, t, "equal-map-v225")
		bs225 = testMarshalErr(&v225v1, h, t, "enc-map-v225-p")
		v225v2 = nil
		testUnmarshalErr(&v225v2, bs225, h, t, "dec-map-v225-p")
		testDeepEqualErr(v225v1, v225v2, t, "equal-map-v225-p")
	}

	for _, v := range []map[int16]int8{nil, map[int16]int8{}, map[int16]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v226: %v\n", v)
		var v226v1, v226v2 map[int16]int8
		v226v1 = v
		bs226 := testMarshalErr(v226v1, h, t, "enc-map-v226")
		if v != nil {
			v226v2 = make(map[int16]int8, len(v))
		}
		testUnmarshalErr(v226v2, bs226, h, t, "dec-map-v226")
		testDeepEqualErr(v226v1, v226v2, t, "equal-map-v226")
		bs226 = testMarshalErr(&v226v1, h, t, "enc-map-v226-p")
		v226v2 = nil
		testUnmarshalErr(&v226v2, bs226, h, t, "dec-map-v226-p")
		testDeepEqualErr(v226v1, v226v2, t, "equal-map-v226-p")
	}

	for _, v := range []map[int16]int16{nil, map[int16]int16{}, map[int16]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v227: %v\n", v)
		var v227v1, v227v2 map[int16]int16
		v227v1 = v
		bs227 := testMarshalErr(v227v1, h, t, "enc-map-v227")
		if v != nil {
			v227v2 = make(map[int16]int16, len(v))
		}
		testUnmarshalErr(v227v2, bs227, h, t, "dec-map-v227")
		testDeepEqualErr(v227v1, v227v2, t, "equal-map-v227")
		bs227 = testMarshalErr(&v227v1, h, t, "enc-map-v227-p")
		v227v2 = nil
		testUnmarshalErr(&v227v2, bs227, h, t, "dec-map-v227-p")
		testDeepEqualErr(v227v1, v227v2, t, "equal-map-v227-p")
	}

	for _, v := range []map[int16]int32{nil, map[int16]int32{}, map[int16]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v228: %v\n", v)
		var v228v1, v228v2 map[int16]int32
		v228v1 = v
		bs228 := testMarshalErr(v228v1, h, t, "enc-map-v228")
		if v != nil {
			v228v2 = make(map[int16]int32, len(v))
		}
		testUnmarshalErr(v228v2, bs228, h, t, "dec-map-v228")
		testDeepEqualErr(v228v1, v228v2, t, "equal-map-v228")
		bs228 = testMarshalErr(&v228v1, h, t, "enc-map-v228-p")
		v228v2 = nil
		testUnmarshalErr(&v228v2, bs228, h, t, "dec-map-v228-p")
		testDeepEqualErr(v228v1, v228v2, t, "equal-map-v228-p")
	}

	for _, v := range []map[int16]int64{nil, map[int16]int64{}, map[int16]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v229: %v\n", v)
		var v229v1, v229v2 map[int16]int64
		v229v1 = v
		bs229 := testMarshalErr(v229v1, h, t, "enc-map-v229")
		if v != nil {
			v229v2 = make(map[int16]int64, len(v))
		}
		testUnmarshalErr(v229v2, bs229, h, t, "dec-map-v229")
		testDeepEqualErr(v229v1, v229v2, t, "equal-map-v229")
		bs229 = testMarshalErr(&v229v1, h, t, "enc-map-v229-p")
		v229v2 = nil
		testUnmarshalErr(&v229v2, bs229, h, t, "dec-map-v229-p")
		testDeepEqualErr(v229v1, v229v2, t, "equal-map-v229-p")
	}

	for _, v := range []map[int16]float32{nil, map[int16]float32{}, map[int16]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v230: %v\n", v)
		var v230v1, v230v2 map[int16]float32
		v230v1 = v
		bs230 := testMarshalErr(v230v1, h, t, "enc-map-v230")
		if v != nil {
			v230v2 = make(map[int16]float32, len(v))
		}
		testUnmarshalErr(v230v2, bs230, h, t, "dec-map-v230")
		testDeepEqualErr(v230v1, v230v2, t, "equal-map-v230")
		bs230 = testMarshalErr(&v230v1, h, t, "enc-map-v230-p")
		v230v2 = nil
		testUnmarshalErr(&v230v2, bs230, h, t, "dec-map-v230-p")
		testDeepEqualErr(v230v1, v230v2, t, "equal-map-v230-p")
	}

	for _, v := range []map[int16]float64{nil, map[int16]float64{}, map[int16]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v231: %v\n", v)
		var v231v1, v231v2 map[int16]float64
		v231v1 = v
		bs231 := testMarshalErr(v231v1, h, t, "enc-map-v231")
		if v != nil {
			v231v2 = make(map[int16]float64, len(v))
		}
		testUnmarshalErr(v231v2, bs231, h, t, "dec-map-v231")
		testDeepEqualErr(v231v1, v231v2, t, "equal-map-v231")
		bs231 = testMarshalErr(&v231v1, h, t, "enc-map-v231-p")
		v231v2 = nil
		testUnmarshalErr(&v231v2, bs231, h, t, "dec-map-v231-p")
		testDeepEqualErr(v231v1, v231v2, t, "equal-map-v231-p")
	}

	for _, v := range []map[int16]bool{nil, map[int16]bool{}, map[int16]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v232: %v\n", v)
		var v232v1, v232v2 map[int16]bool
		v232v1 = v
		bs232 := testMarshalErr(v232v1, h, t, "enc-map-v232")
		if v != nil {
			v232v2 = make(map[int16]bool, len(v))
		}
		testUnmarshalErr(v232v2, bs232, h, t, "dec-map-v232")
		testDeepEqualErr(v232v1, v232v2, t, "equal-map-v232")
		bs232 = testMarshalErr(&v232v1, h, t, "enc-map-v232-p")
		v232v2 = nil
		testUnmarshalErr(&v232v2, bs232, h, t, "dec-map-v232-p")
		testDeepEqualErr(v232v1, v232v2, t, "equal-map-v232-p")
	}

	for _, v := range []map[int32]interface{}{nil, map[int32]interface{}{}, map[int32]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v235: %v\n", v)
		var v235v1, v235v2 map[int32]interface{}
		v235v1 = v
		bs235 := testMarshalErr(v235v1, h, t, "enc-map-v235")
		if v != nil {
			v235v2 = make(map[int32]interface{}, len(v))
		}
		testUnmarshalErr(v235v2, bs235, h, t, "dec-map-v235")
		testDeepEqualErr(v235v1, v235v2, t, "equal-map-v235")
		bs235 = testMarshalErr(&v235v1, h, t, "enc-map-v235-p")
		v235v2 = nil
		testUnmarshalErr(&v235v2, bs235, h, t, "dec-map-v235-p")
		testDeepEqualErr(v235v1, v235v2, t, "equal-map-v235-p")
	}

	for _, v := range []map[int32]string{nil, map[int32]string{}, map[int32]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v236: %v\n", v)
		var v236v1, v236v2 map[int32]string
		v236v1 = v
		bs236 := testMarshalErr(v236v1, h, t, "enc-map-v236")
		if v != nil {
			v236v2 = make(map[int32]string, len(v))
		}
		testUnmarshalErr(v236v2, bs236, h, t, "dec-map-v236")
		testDeepEqualErr(v236v1, v236v2, t, "equal-map-v236")
		bs236 = testMarshalErr(&v236v1, h, t, "enc-map-v236-p")
		v236v2 = nil
		testUnmarshalErr(&v236v2, bs236, h, t, "dec-map-v236-p")
		testDeepEqualErr(v236v1, v236v2, t, "equal-map-v236-p")
	}

	for _, v := range []map[int32]uint{nil, map[int32]uint{}, map[int32]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v237: %v\n", v)
		var v237v1, v237v2 map[int32]uint
		v237v1 = v
		bs237 := testMarshalErr(v237v1, h, t, "enc-map-v237")
		if v != nil {
			v237v2 = make(map[int32]uint, len(v))
		}
		testUnmarshalErr(v237v2, bs237, h, t, "dec-map-v237")
		testDeepEqualErr(v237v1, v237v2, t, "equal-map-v237")
		bs237 = testMarshalErr(&v237v1, h, t, "enc-map-v237-p")
		v237v2 = nil
		testUnmarshalErr(&v237v2, bs237, h, t, "dec-map-v237-p")
		testDeepEqualErr(v237v1, v237v2, t, "equal-map-v237-p")
	}

	for _, v := range []map[int32]uint8{nil, map[int32]uint8{}, map[int32]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v238: %v\n", v)
		var v238v1, v238v2 map[int32]uint8
		v238v1 = v
		bs238 := testMarshalErr(v238v1, h, t, "enc-map-v238")
		if v != nil {
			v238v2 = make(map[int32]uint8, len(v))
		}
		testUnmarshalErr(v238v2, bs238, h, t, "dec-map-v238")
		testDeepEqualErr(v238v1, v238v2, t, "equal-map-v238")
		bs238 = testMarshalErr(&v238v1, h, t, "enc-map-v238-p")
		v238v2 = nil
		testUnmarshalErr(&v238v2, bs238, h, t, "dec-map-v238-p")
		testDeepEqualErr(v238v1, v238v2, t, "equal-map-v238-p")
	}

	for _, v := range []map[int32]uint16{nil, map[int32]uint16{}, map[int32]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v239: %v\n", v)
		var v239v1, v239v2 map[int32]uint16
		v239v1 = v
		bs239 := testMarshalErr(v239v1, h, t, "enc-map-v239")
		if v != nil {
			v239v2 = make(map[int32]uint16, len(v))
		}
		testUnmarshalErr(v239v2, bs239, h, t, "dec-map-v239")
		testDeepEqualErr(v239v1, v239v2, t, "equal-map-v239")
		bs239 = testMarshalErr(&v239v1, h, t, "enc-map-v239-p")
		v239v2 = nil
		testUnmarshalErr(&v239v2, bs239, h, t, "dec-map-v239-p")
		testDeepEqualErr(v239v1, v239v2, t, "equal-map-v239-p")
	}

	for _, v := range []map[int32]uint32{nil, map[int32]uint32{}, map[int32]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v240: %v\n", v)
		var v240v1, v240v2 map[int32]uint32
		v240v1 = v
		bs240 := testMarshalErr(v240v1, h, t, "enc-map-v240")
		if v != nil {
			v240v2 = make(map[int32]uint32, len(v))
		}
		testUnmarshalErr(v240v2, bs240, h, t, "dec-map-v240")
		testDeepEqualErr(v240v1, v240v2, t, "equal-map-v240")
		bs240 = testMarshalErr(&v240v1, h, t, "enc-map-v240-p")
		v240v2 = nil
		testUnmarshalErr(&v240v2, bs240, h, t, "dec-map-v240-p")
		testDeepEqualErr(v240v1, v240v2, t, "equal-map-v240-p")
	}

	for _, v := range []map[int32]uint64{nil, map[int32]uint64{}, map[int32]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v241: %v\n", v)
		var v241v1, v241v2 map[int32]uint64
		v241v1 = v
		bs241 := testMarshalErr(v241v1, h, t, "enc-map-v241")
		if v != nil {
			v241v2 = make(map[int32]uint64, len(v))
		}
		testUnmarshalErr(v241v2, bs241, h, t, "dec-map-v241")
		testDeepEqualErr(v241v1, v241v2, t, "equal-map-v241")
		bs241 = testMarshalErr(&v241v1, h, t, "enc-map-v241-p")
		v241v2 = nil
		testUnmarshalErr(&v241v2, bs241, h, t, "dec-map-v241-p")
		testDeepEqualErr(v241v1, v241v2, t, "equal-map-v241-p")
	}

	for _, v := range []map[int32]uintptr{nil, map[int32]uintptr{}, map[int32]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v242: %v\n", v)
		var v242v1, v242v2 map[int32]uintptr
		v242v1 = v
		bs242 := testMarshalErr(v242v1, h, t, "enc-map-v242")
		if v != nil {
			v242v2 = make(map[int32]uintptr, len(v))
		}
		testUnmarshalErr(v242v2, bs242, h, t, "dec-map-v242")
		testDeepEqualErr(v242v1, v242v2, t, "equal-map-v242")
		bs242 = testMarshalErr(&v242v1, h, t, "enc-map-v242-p")
		v242v2 = nil
		testUnmarshalErr(&v242v2, bs242, h, t, "dec-map-v242-p")
		testDeepEqualErr(v242v1, v242v2, t, "equal-map-v242-p")
	}

	for _, v := range []map[int32]int{nil, map[int32]int{}, map[int32]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v243: %v\n", v)
		var v243v1, v243v2 map[int32]int
		v243v1 = v
		bs243 := testMarshalErr(v243v1, h, t, "enc-map-v243")
		if v != nil {
			v243v2 = make(map[int32]int, len(v))
		}
		testUnmarshalErr(v243v2, bs243, h, t, "dec-map-v243")
		testDeepEqualErr(v243v1, v243v2, t, "equal-map-v243")
		bs243 = testMarshalErr(&v243v1, h, t, "enc-map-v243-p")
		v243v2 = nil
		testUnmarshalErr(&v243v2, bs243, h, t, "dec-map-v243-p")
		testDeepEqualErr(v243v1, v243v2, t, "equal-map-v243-p")
	}

	for _, v := range []map[int32]int8{nil, map[int32]int8{}, map[int32]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v244: %v\n", v)
		var v244v1, v244v2 map[int32]int8
		v244v1 = v
		bs244 := testMarshalErr(v244v1, h, t, "enc-map-v244")
		if v != nil {
			v244v2 = make(map[int32]int8, len(v))
		}
		testUnmarshalErr(v244v2, bs244, h, t, "dec-map-v244")
		testDeepEqualErr(v244v1, v244v2, t, "equal-map-v244")
		bs244 = testMarshalErr(&v244v1, h, t, "enc-map-v244-p")
		v244v2 = nil
		testUnmarshalErr(&v244v2, bs244, h, t, "dec-map-v244-p")
		testDeepEqualErr(v244v1, v244v2, t, "equal-map-v244-p")
	}

	for _, v := range []map[int32]int16{nil, map[int32]int16{}, map[int32]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v245: %v\n", v)
		var v245v1, v245v2 map[int32]int16
		v245v1 = v
		bs245 := testMarshalErr(v245v1, h, t, "enc-map-v245")
		if v != nil {
			v245v2 = make(map[int32]int16, len(v))
		}
		testUnmarshalErr(v245v2, bs245, h, t, "dec-map-v245")
		testDeepEqualErr(v245v1, v245v2, t, "equal-map-v245")
		bs245 = testMarshalErr(&v245v1, h, t, "enc-map-v245-p")
		v245v2 = nil
		testUnmarshalErr(&v245v2, bs245, h, t, "dec-map-v245-p")
		testDeepEqualErr(v245v1, v245v2, t, "equal-map-v245-p")
	}

	for _, v := range []map[int32]int32{nil, map[int32]int32{}, map[int32]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v246: %v\n", v)
		var v246v1, v246v2 map[int32]int32
		v246v1 = v
		bs246 := testMarshalErr(v246v1, h, t, "enc-map-v246")
		if v != nil {
			v246v2 = make(map[int32]int32, len(v))
		}
		testUnmarshalErr(v246v2, bs246, h, t, "dec-map-v246")
		testDeepEqualErr(v246v1, v246v2, t, "equal-map-v246")
		bs246 = testMarshalErr(&v246v1, h, t, "enc-map-v246-p")
		v246v2 = nil
		testUnmarshalErr(&v246v2, bs246, h, t, "dec-map-v246-p")
		testDeepEqualErr(v246v1, v246v2, t, "equal-map-v246-p")
	}

	for _, v := range []map[int32]int64{nil, map[int32]int64{}, map[int32]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v247: %v\n", v)
		var v247v1, v247v2 map[int32]int64
		v247v1 = v
		bs247 := testMarshalErr(v247v1, h, t, "enc-map-v247")
		if v != nil {
			v247v2 = make(map[int32]int64, len(v))
		}
		testUnmarshalErr(v247v2, bs247, h, t, "dec-map-v247")
		testDeepEqualErr(v247v1, v247v2, t, "equal-map-v247")
		bs247 = testMarshalErr(&v247v1, h, t, "enc-map-v247-p")
		v247v2 = nil
		testUnmarshalErr(&v247v2, bs247, h, t, "dec-map-v247-p")
		testDeepEqualErr(v247v1, v247v2, t, "equal-map-v247-p")
	}

	for _, v := range []map[int32]float32{nil, map[int32]float32{}, map[int32]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v248: %v\n", v)
		var v248v1, v248v2 map[int32]float32
		v248v1 = v
		bs248 := testMarshalErr(v248v1, h, t, "enc-map-v248")
		if v != nil {
			v248v2 = make(map[int32]float32, len(v))
		}
		testUnmarshalErr(v248v2, bs248, h, t, "dec-map-v248")
		testDeepEqualErr(v248v1, v248v2, t, "equal-map-v248")
		bs248 = testMarshalErr(&v248v1, h, t, "enc-map-v248-p")
		v248v2 = nil
		testUnmarshalErr(&v248v2, bs248, h, t, "dec-map-v248-p")
		testDeepEqualErr(v248v1, v248v2, t, "equal-map-v248-p")
	}

	for _, v := range []map[int32]float64{nil, map[int32]float64{}, map[int32]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v249: %v\n", v)
		var v249v1, v249v2 map[int32]float64
		v249v1 = v
		bs249 := testMarshalErr(v249v1, h, t, "enc-map-v249")
		if v != nil {
			v249v2 = make(map[int32]float64, len(v))
		}
		testUnmarshalErr(v249v2, bs249, h, t, "dec-map-v249")
		testDeepEqualErr(v249v1, v249v2, t, "equal-map-v249")
		bs249 = testMarshalErr(&v249v1, h, t, "enc-map-v249-p")
		v249v2 = nil
		testUnmarshalErr(&v249v2, bs249, h, t, "dec-map-v249-p")
		testDeepEqualErr(v249v1, v249v2, t, "equal-map-v249-p")
	}

	for _, v := range []map[int32]bool{nil, map[int32]bool{}, map[int32]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v250: %v\n", v)
		var v250v1, v250v2 map[int32]bool
		v250v1 = v
		bs250 := testMarshalErr(v250v1, h, t, "enc-map-v250")
		if v != nil {
			v250v2 = make(map[int32]bool, len(v))
		}
		testUnmarshalErr(v250v2, bs250, h, t, "dec-map-v250")
		testDeepEqualErr(v250v1, v250v2, t, "equal-map-v250")
		bs250 = testMarshalErr(&v250v1, h, t, "enc-map-v250-p")
		v250v2 = nil
		testUnmarshalErr(&v250v2, bs250, h, t, "dec-map-v250-p")
		testDeepEqualErr(v250v1, v250v2, t, "equal-map-v250-p")
	}

	for _, v := range []map[int64]interface{}{nil, map[int64]interface{}{}, map[int64]interface{}{10: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v253: %v\n", v)
		var v253v1, v253v2 map[int64]interface{}
		v253v1 = v
		bs253 := testMarshalErr(v253v1, h, t, "enc-map-v253")
		if v != nil {
			v253v2 = make(map[int64]interface{}, len(v))
		}
		testUnmarshalErr(v253v2, bs253, h, t, "dec-map-v253")
		testDeepEqualErr(v253v1, v253v2, t, "equal-map-v253")
		bs253 = testMarshalErr(&v253v1, h, t, "enc-map-v253-p")
		v253v2 = nil
		testUnmarshalErr(&v253v2, bs253, h, t, "dec-map-v253-p")
		testDeepEqualErr(v253v1, v253v2, t, "equal-map-v253-p")
	}

	for _, v := range []map[int64]string{nil, map[int64]string{}, map[int64]string{10: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v254: %v\n", v)
		var v254v1, v254v2 map[int64]string
		v254v1 = v
		bs254 := testMarshalErr(v254v1, h, t, "enc-map-v254")
		if v != nil {
			v254v2 = make(map[int64]string, len(v))
		}
		testUnmarshalErr(v254v2, bs254, h, t, "dec-map-v254")
		testDeepEqualErr(v254v1, v254v2, t, "equal-map-v254")
		bs254 = testMarshalErr(&v254v1, h, t, "enc-map-v254-p")
		v254v2 = nil
		testUnmarshalErr(&v254v2, bs254, h, t, "dec-map-v254-p")
		testDeepEqualErr(v254v1, v254v2, t, "equal-map-v254-p")
	}

	for _, v := range []map[int64]uint{nil, map[int64]uint{}, map[int64]uint{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v255: %v\n", v)
		var v255v1, v255v2 map[int64]uint
		v255v1 = v
		bs255 := testMarshalErr(v255v1, h, t, "enc-map-v255")
		if v != nil {
			v255v2 = make(map[int64]uint, len(v))
		}
		testUnmarshalErr(v255v2, bs255, h, t, "dec-map-v255")
		testDeepEqualErr(v255v1, v255v2, t, "equal-map-v255")
		bs255 = testMarshalErr(&v255v1, h, t, "enc-map-v255-p")
		v255v2 = nil
		testUnmarshalErr(&v255v2, bs255, h, t, "dec-map-v255-p")
		testDeepEqualErr(v255v1, v255v2, t, "equal-map-v255-p")
	}

	for _, v := range []map[int64]uint8{nil, map[int64]uint8{}, map[int64]uint8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v256: %v\n", v)
		var v256v1, v256v2 map[int64]uint8
		v256v1 = v
		bs256 := testMarshalErr(v256v1, h, t, "enc-map-v256")
		if v != nil {
			v256v2 = make(map[int64]uint8, len(v))
		}
		testUnmarshalErr(v256v2, bs256, h, t, "dec-map-v256")
		testDeepEqualErr(v256v1, v256v2, t, "equal-map-v256")
		bs256 = testMarshalErr(&v256v1, h, t, "enc-map-v256-p")
		v256v2 = nil
		testUnmarshalErr(&v256v2, bs256, h, t, "dec-map-v256-p")
		testDeepEqualErr(v256v1, v256v2, t, "equal-map-v256-p")
	}

	for _, v := range []map[int64]uint16{nil, map[int64]uint16{}, map[int64]uint16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v257: %v\n", v)
		var v257v1, v257v2 map[int64]uint16
		v257v1 = v
		bs257 := testMarshalErr(v257v1, h, t, "enc-map-v257")
		if v != nil {
			v257v2 = make(map[int64]uint16, len(v))
		}
		testUnmarshalErr(v257v2, bs257, h, t, "dec-map-v257")
		testDeepEqualErr(v257v1, v257v2, t, "equal-map-v257")
		bs257 = testMarshalErr(&v257v1, h, t, "enc-map-v257-p")
		v257v2 = nil
		testUnmarshalErr(&v257v2, bs257, h, t, "dec-map-v257-p")
		testDeepEqualErr(v257v1, v257v2, t, "equal-map-v257-p")
	}

	for _, v := range []map[int64]uint32{nil, map[int64]uint32{}, map[int64]uint32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v258: %v\n", v)
		var v258v1, v258v2 map[int64]uint32
		v258v1 = v
		bs258 := testMarshalErr(v258v1, h, t, "enc-map-v258")
		if v != nil {
			v258v2 = make(map[int64]uint32, len(v))
		}
		testUnmarshalErr(v258v2, bs258, h, t, "dec-map-v258")
		testDeepEqualErr(v258v1, v258v2, t, "equal-map-v258")
		bs258 = testMarshalErr(&v258v1, h, t, "enc-map-v258-p")
		v258v2 = nil
		testUnmarshalErr(&v258v2, bs258, h, t, "dec-map-v258-p")
		testDeepEqualErr(v258v1, v258v2, t, "equal-map-v258-p")
	}

	for _, v := range []map[int64]uint64{nil, map[int64]uint64{}, map[int64]uint64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v259: %v\n", v)
		var v259v1, v259v2 map[int64]uint64
		v259v1 = v
		bs259 := testMarshalErr(v259v1, h, t, "enc-map-v259")
		if v != nil {
			v259v2 = make(map[int64]uint64, len(v))
		}
		testUnmarshalErr(v259v2, bs259, h, t, "dec-map-v259")
		testDeepEqualErr(v259v1, v259v2, t, "equal-map-v259")
		bs259 = testMarshalErr(&v259v1, h, t, "enc-map-v259-p")
		v259v2 = nil
		testUnmarshalErr(&v259v2, bs259, h, t, "dec-map-v259-p")
		testDeepEqualErr(v259v1, v259v2, t, "equal-map-v259-p")
	}

	for _, v := range []map[int64]uintptr{nil, map[int64]uintptr{}, map[int64]uintptr{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v260: %v\n", v)
		var v260v1, v260v2 map[int64]uintptr
		v260v1 = v
		bs260 := testMarshalErr(v260v1, h, t, "enc-map-v260")
		if v != nil {
			v260v2 = make(map[int64]uintptr, len(v))
		}
		testUnmarshalErr(v260v2, bs260, h, t, "dec-map-v260")
		testDeepEqualErr(v260v1, v260v2, t, "equal-map-v260")
		bs260 = testMarshalErr(&v260v1, h, t, "enc-map-v260-p")
		v260v2 = nil
		testUnmarshalErr(&v260v2, bs260, h, t, "dec-map-v260-p")
		testDeepEqualErr(v260v1, v260v2, t, "equal-map-v260-p")
	}

	for _, v := range []map[int64]int{nil, map[int64]int{}, map[int64]int{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v261: %v\n", v)
		var v261v1, v261v2 map[int64]int
		v261v1 = v
		bs261 := testMarshalErr(v261v1, h, t, "enc-map-v261")
		if v != nil {
			v261v2 = make(map[int64]int, len(v))
		}
		testUnmarshalErr(v261v2, bs261, h, t, "dec-map-v261")
		testDeepEqualErr(v261v1, v261v2, t, "equal-map-v261")
		bs261 = testMarshalErr(&v261v1, h, t, "enc-map-v261-p")
		v261v2 = nil
		testUnmarshalErr(&v261v2, bs261, h, t, "dec-map-v261-p")
		testDeepEqualErr(v261v1, v261v2, t, "equal-map-v261-p")
	}

	for _, v := range []map[int64]int8{nil, map[int64]int8{}, map[int64]int8{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v262: %v\n", v)
		var v262v1, v262v2 map[int64]int8
		v262v1 = v
		bs262 := testMarshalErr(v262v1, h, t, "enc-map-v262")
		if v != nil {
			v262v2 = make(map[int64]int8, len(v))
		}
		testUnmarshalErr(v262v2, bs262, h, t, "dec-map-v262")
		testDeepEqualErr(v262v1, v262v2, t, "equal-map-v262")
		bs262 = testMarshalErr(&v262v1, h, t, "enc-map-v262-p")
		v262v2 = nil
		testUnmarshalErr(&v262v2, bs262, h, t, "dec-map-v262-p")
		testDeepEqualErr(v262v1, v262v2, t, "equal-map-v262-p")
	}

	for _, v := range []map[int64]int16{nil, map[int64]int16{}, map[int64]int16{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v263: %v\n", v)
		var v263v1, v263v2 map[int64]int16
		v263v1 = v
		bs263 := testMarshalErr(v263v1, h, t, "enc-map-v263")
		if v != nil {
			v263v2 = make(map[int64]int16, len(v))
		}
		testUnmarshalErr(v263v2, bs263, h, t, "dec-map-v263")
		testDeepEqualErr(v263v1, v263v2, t, "equal-map-v263")
		bs263 = testMarshalErr(&v263v1, h, t, "enc-map-v263-p")
		v263v2 = nil
		testUnmarshalErr(&v263v2, bs263, h, t, "dec-map-v263-p")
		testDeepEqualErr(v263v1, v263v2, t, "equal-map-v263-p")
	}

	for _, v := range []map[int64]int32{nil, map[int64]int32{}, map[int64]int32{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v264: %v\n", v)
		var v264v1, v264v2 map[int64]int32
		v264v1 = v
		bs264 := testMarshalErr(v264v1, h, t, "enc-map-v264")
		if v != nil {
			v264v2 = make(map[int64]int32, len(v))
		}
		testUnmarshalErr(v264v2, bs264, h, t, "dec-map-v264")
		testDeepEqualErr(v264v1, v264v2, t, "equal-map-v264")
		bs264 = testMarshalErr(&v264v1, h, t, "enc-map-v264-p")
		v264v2 = nil
		testUnmarshalErr(&v264v2, bs264, h, t, "dec-map-v264-p")
		testDeepEqualErr(v264v1, v264v2, t, "equal-map-v264-p")
	}

	for _, v := range []map[int64]int64{nil, map[int64]int64{}, map[int64]int64{10: 10}} {
		// fmt.Printf(">>>> running mammoth map v265: %v\n", v)
		var v265v1, v265v2 map[int64]int64
		v265v1 = v
		bs265 := testMarshalErr(v265v1, h, t, "enc-map-v265")
		if v != nil {
			v265v2 = make(map[int64]int64, len(v))
		}
		testUnmarshalErr(v265v2, bs265, h, t, "dec-map-v265")
		testDeepEqualErr(v265v1, v265v2, t, "equal-map-v265")
		bs265 = testMarshalErr(&v265v1, h, t, "enc-map-v265-p")
		v265v2 = nil
		testUnmarshalErr(&v265v2, bs265, h, t, "dec-map-v265-p")
		testDeepEqualErr(v265v1, v265v2, t, "equal-map-v265-p")
	}

	for _, v := range []map[int64]float32{nil, map[int64]float32{}, map[int64]float32{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v266: %v\n", v)
		var v266v1, v266v2 map[int64]float32
		v266v1 = v
		bs266 := testMarshalErr(v266v1, h, t, "enc-map-v266")
		if v != nil {
			v266v2 = make(map[int64]float32, len(v))
		}
		testUnmarshalErr(v266v2, bs266, h, t, "dec-map-v266")
		testDeepEqualErr(v266v1, v266v2, t, "equal-map-v266")
		bs266 = testMarshalErr(&v266v1, h, t, "enc-map-v266-p")
		v266v2 = nil
		testUnmarshalErr(&v266v2, bs266, h, t, "dec-map-v266-p")
		testDeepEqualErr(v266v1, v266v2, t, "equal-map-v266-p")
	}

	for _, v := range []map[int64]float64{nil, map[int64]float64{}, map[int64]float64{10: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v267: %v\n", v)
		var v267v1, v267v2 map[int64]float64
		v267v1 = v
		bs267 := testMarshalErr(v267v1, h, t, "enc-map-v267")
		if v != nil {
			v267v2 = make(map[int64]float64, len(v))
		}
		testUnmarshalErr(v267v2, bs267, h, t, "dec-map-v267")
		testDeepEqualErr(v267v1, v267v2, t, "equal-map-v267")
		bs267 = testMarshalErr(&v267v1, h, t, "enc-map-v267-p")
		v267v2 = nil
		testUnmarshalErr(&v267v2, bs267, h, t, "dec-map-v267-p")
		testDeepEqualErr(v267v1, v267v2, t, "equal-map-v267-p")
	}

	for _, v := range []map[int64]bool{nil, map[int64]bool{}, map[int64]bool{10: true}} {
		// fmt.Printf(">>>> running mammoth map v268: %v\n", v)
		var v268v1, v268v2 map[int64]bool
		v268v1 = v
		bs268 := testMarshalErr(v268v1, h, t, "enc-map-v268")
		if v != nil {
			v268v2 = make(map[int64]bool, len(v))
		}
		testUnmarshalErr(v268v2, bs268, h, t, "dec-map-v268")
		testDeepEqualErr(v268v1, v268v2, t, "equal-map-v268")
		bs268 = testMarshalErr(&v268v1, h, t, "enc-map-v268-p")
		v268v2 = nil
		testUnmarshalErr(&v268v2, bs268, h, t, "dec-map-v268-p")
		testDeepEqualErr(v268v1, v268v2, t, "equal-map-v268-p")
	}

	for _, v := range []map[bool]interface{}{nil, map[bool]interface{}{}, map[bool]interface{}{true: "string-is-an-interface"}} {
		// fmt.Printf(">>>> running mammoth map v271: %v\n", v)
		var v271v1, v271v2 map[bool]interface{}
		v271v1 = v
		bs271 := testMarshalErr(v271v1, h, t, "enc-map-v271")
		if v != nil {
			v271v2 = make(map[bool]interface{}, len(v))
		}
		testUnmarshalErr(v271v2, bs271, h, t, "dec-map-v271")
		testDeepEqualErr(v271v1, v271v2, t, "equal-map-v271")
		bs271 = testMarshalErr(&v271v1, h, t, "enc-map-v271-p")
		v271v2 = nil
		testUnmarshalErr(&v271v2, bs271, h, t, "dec-map-v271-p")
		testDeepEqualErr(v271v1, v271v2, t, "equal-map-v271-p")
	}

	for _, v := range []map[bool]string{nil, map[bool]string{}, map[bool]string{true: "some-string"}} {
		// fmt.Printf(">>>> running mammoth map v272: %v\n", v)
		var v272v1, v272v2 map[bool]string
		v272v1 = v
		bs272 := testMarshalErr(v272v1, h, t, "enc-map-v272")
		if v != nil {
			v272v2 = make(map[bool]string, len(v))
		}
		testUnmarshalErr(v272v2, bs272, h, t, "dec-map-v272")
		testDeepEqualErr(v272v1, v272v2, t, "equal-map-v272")
		bs272 = testMarshalErr(&v272v1, h, t, "enc-map-v272-p")
		v272v2 = nil
		testUnmarshalErr(&v272v2, bs272, h, t, "dec-map-v272-p")
		testDeepEqualErr(v272v1, v272v2, t, "equal-map-v272-p")
	}

	for _, v := range []map[bool]uint{nil, map[bool]uint{}, map[bool]uint{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v273: %v\n", v)
		var v273v1, v273v2 map[bool]uint
		v273v1 = v
		bs273 := testMarshalErr(v273v1, h, t, "enc-map-v273")
		if v != nil {
			v273v2 = make(map[bool]uint, len(v))
		}
		testUnmarshalErr(v273v2, bs273, h, t, "dec-map-v273")
		testDeepEqualErr(v273v1, v273v2, t, "equal-map-v273")
		bs273 = testMarshalErr(&v273v1, h, t, "enc-map-v273-p")
		v273v2 = nil
		testUnmarshalErr(&v273v2, bs273, h, t, "dec-map-v273-p")
		testDeepEqualErr(v273v1, v273v2, t, "equal-map-v273-p")
	}

	for _, v := range []map[bool]uint8{nil, map[bool]uint8{}, map[bool]uint8{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v274: %v\n", v)
		var v274v1, v274v2 map[bool]uint8
		v274v1 = v
		bs274 := testMarshalErr(v274v1, h, t, "enc-map-v274")
		if v != nil {
			v274v2 = make(map[bool]uint8, len(v))
		}
		testUnmarshalErr(v274v2, bs274, h, t, "dec-map-v274")
		testDeepEqualErr(v274v1, v274v2, t, "equal-map-v274")
		bs274 = testMarshalErr(&v274v1, h, t, "enc-map-v274-p")
		v274v2 = nil
		testUnmarshalErr(&v274v2, bs274, h, t, "dec-map-v274-p")
		testDeepEqualErr(v274v1, v274v2, t, "equal-map-v274-p")
	}

	for _, v := range []map[bool]uint16{nil, map[bool]uint16{}, map[bool]uint16{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v275: %v\n", v)
		var v275v1, v275v2 map[bool]uint16
		v275v1 = v
		bs275 := testMarshalErr(v275v1, h, t, "enc-map-v275")
		if v != nil {
			v275v2 = make(map[bool]uint16, len(v))
		}
		testUnmarshalErr(v275v2, bs275, h, t, "dec-map-v275")
		testDeepEqualErr(v275v1, v275v2, t, "equal-map-v275")
		bs275 = testMarshalErr(&v275v1, h, t, "enc-map-v275-p")
		v275v2 = nil
		testUnmarshalErr(&v275v2, bs275, h, t, "dec-map-v275-p")
		testDeepEqualErr(v275v1, v275v2, t, "equal-map-v275-p")
	}

	for _, v := range []map[bool]uint32{nil, map[bool]uint32{}, map[bool]uint32{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v276: %v\n", v)
		var v276v1, v276v2 map[bool]uint32
		v276v1 = v
		bs276 := testMarshalErr(v276v1, h, t, "enc-map-v276")
		if v != nil {
			v276v2 = make(map[bool]uint32, len(v))
		}
		testUnmarshalErr(v276v2, bs276, h, t, "dec-map-v276")
		testDeepEqualErr(v276v1, v276v2, t, "equal-map-v276")
		bs276 = testMarshalErr(&v276v1, h, t, "enc-map-v276-p")
		v276v2 = nil
		testUnmarshalErr(&v276v2, bs276, h, t, "dec-map-v276-p")
		testDeepEqualErr(v276v1, v276v2, t, "equal-map-v276-p")
	}

	for _, v := range []map[bool]uint64{nil, map[bool]uint64{}, map[bool]uint64{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v277: %v\n", v)
		var v277v1, v277v2 map[bool]uint64
		v277v1 = v
		bs277 := testMarshalErr(v277v1, h, t, "enc-map-v277")
		if v != nil {
			v277v2 = make(map[bool]uint64, len(v))
		}
		testUnmarshalErr(v277v2, bs277, h, t, "dec-map-v277")
		testDeepEqualErr(v277v1, v277v2, t, "equal-map-v277")
		bs277 = testMarshalErr(&v277v1, h, t, "enc-map-v277-p")
		v277v2 = nil
		testUnmarshalErr(&v277v2, bs277, h, t, "dec-map-v277-p")
		testDeepEqualErr(v277v1, v277v2, t, "equal-map-v277-p")
	}

	for _, v := range []map[bool]uintptr{nil, map[bool]uintptr{}, map[bool]uintptr{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v278: %v\n", v)
		var v278v1, v278v2 map[bool]uintptr
		v278v1 = v
		bs278 := testMarshalErr(v278v1, h, t, "enc-map-v278")
		if v != nil {
			v278v2 = make(map[bool]uintptr, len(v))
		}
		testUnmarshalErr(v278v2, bs278, h, t, "dec-map-v278")
		testDeepEqualErr(v278v1, v278v2, t, "equal-map-v278")
		bs278 = testMarshalErr(&v278v1, h, t, "enc-map-v278-p")
		v278v2 = nil
		testUnmarshalErr(&v278v2, bs278, h, t, "dec-map-v278-p")
		testDeepEqualErr(v278v1, v278v2, t, "equal-map-v278-p")
	}

	for _, v := range []map[bool]int{nil, map[bool]int{}, map[bool]int{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v279: %v\n", v)
		var v279v1, v279v2 map[bool]int
		v279v1 = v
		bs279 := testMarshalErr(v279v1, h, t, "enc-map-v279")
		if v != nil {
			v279v2 = make(map[bool]int, len(v))
		}
		testUnmarshalErr(v279v2, bs279, h, t, "dec-map-v279")
		testDeepEqualErr(v279v1, v279v2, t, "equal-map-v279")
		bs279 = testMarshalErr(&v279v1, h, t, "enc-map-v279-p")
		v279v2 = nil
		testUnmarshalErr(&v279v2, bs279, h, t, "dec-map-v279-p")
		testDeepEqualErr(v279v1, v279v2, t, "equal-map-v279-p")
	}

	for _, v := range []map[bool]int8{nil, map[bool]int8{}, map[bool]int8{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v280: %v\n", v)
		var v280v1, v280v2 map[bool]int8
		v280v1 = v
		bs280 := testMarshalErr(v280v1, h, t, "enc-map-v280")
		if v != nil {
			v280v2 = make(map[bool]int8, len(v))
		}
		testUnmarshalErr(v280v2, bs280, h, t, "dec-map-v280")
		testDeepEqualErr(v280v1, v280v2, t, "equal-map-v280")
		bs280 = testMarshalErr(&v280v1, h, t, "enc-map-v280-p")
		v280v2 = nil
		testUnmarshalErr(&v280v2, bs280, h, t, "dec-map-v280-p")
		testDeepEqualErr(v280v1, v280v2, t, "equal-map-v280-p")
	}

	for _, v := range []map[bool]int16{nil, map[bool]int16{}, map[bool]int16{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v281: %v\n", v)
		var v281v1, v281v2 map[bool]int16
		v281v1 = v
		bs281 := testMarshalErr(v281v1, h, t, "enc-map-v281")
		if v != nil {
			v281v2 = make(map[bool]int16, len(v))
		}
		testUnmarshalErr(v281v2, bs281, h, t, "dec-map-v281")
		testDeepEqualErr(v281v1, v281v2, t, "equal-map-v281")
		bs281 = testMarshalErr(&v281v1, h, t, "enc-map-v281-p")
		v281v2 = nil
		testUnmarshalErr(&v281v2, bs281, h, t, "dec-map-v281-p")
		testDeepEqualErr(v281v1, v281v2, t, "equal-map-v281-p")
	}

	for _, v := range []map[bool]int32{nil, map[bool]int32{}, map[bool]int32{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v282: %v\n", v)
		var v282v1, v282v2 map[bool]int32
		v282v1 = v
		bs282 := testMarshalErr(v282v1, h, t, "enc-map-v282")
		if v != nil {
			v282v2 = make(map[bool]int32, len(v))
		}
		testUnmarshalErr(v282v2, bs282, h, t, "dec-map-v282")
		testDeepEqualErr(v282v1, v282v2, t, "equal-map-v282")
		bs282 = testMarshalErr(&v282v1, h, t, "enc-map-v282-p")
		v282v2 = nil
		testUnmarshalErr(&v282v2, bs282, h, t, "dec-map-v282-p")
		testDeepEqualErr(v282v1, v282v2, t, "equal-map-v282-p")
	}

	for _, v := range []map[bool]int64{nil, map[bool]int64{}, map[bool]int64{true: 10}} {
		// fmt.Printf(">>>> running mammoth map v283: %v\n", v)
		var v283v1, v283v2 map[bool]int64
		v283v1 = v
		bs283 := testMarshalErr(v283v1, h, t, "enc-map-v283")
		if v != nil {
			v283v2 = make(map[bool]int64, len(v))
		}
		testUnmarshalErr(v283v2, bs283, h, t, "dec-map-v283")
		testDeepEqualErr(v283v1, v283v2, t, "equal-map-v283")
		bs283 = testMarshalErr(&v283v1, h, t, "enc-map-v283-p")
		v283v2 = nil
		testUnmarshalErr(&v283v2, bs283, h, t, "dec-map-v283-p")
		testDeepEqualErr(v283v1, v283v2, t, "equal-map-v283-p")
	}

	for _, v := range []map[bool]float32{nil, map[bool]float32{}, map[bool]float32{true: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v284: %v\n", v)
		var v284v1, v284v2 map[bool]float32
		v284v1 = v
		bs284 := testMarshalErr(v284v1, h, t, "enc-map-v284")
		if v != nil {
			v284v2 = make(map[bool]float32, len(v))
		}
		testUnmarshalErr(v284v2, bs284, h, t, "dec-map-v284")
		testDeepEqualErr(v284v1, v284v2, t, "equal-map-v284")
		bs284 = testMarshalErr(&v284v1, h, t, "enc-map-v284-p")
		v284v2 = nil
		testUnmarshalErr(&v284v2, bs284, h, t, "dec-map-v284-p")
		testDeepEqualErr(v284v1, v284v2, t, "equal-map-v284-p")
	}

	for _, v := range []map[bool]float64{nil, map[bool]float64{}, map[bool]float64{true: 10.1}} {
		// fmt.Printf(">>>> running mammoth map v285: %v\n", v)
		var v285v1, v285v2 map[bool]float64
		v285v1 = v
		bs285 := testMarshalErr(v285v1, h, t, "enc-map-v285")
		if v != nil {
			v285v2 = make(map[bool]float64, len(v))
		}
		testUnmarshalErr(v285v2, bs285, h, t, "dec-map-v285")
		testDeepEqualErr(v285v1, v285v2, t, "equal-map-v285")
		bs285 = testMarshalErr(&v285v1, h, t, "enc-map-v285-p")
		v285v2 = nil
		testUnmarshalErr(&v285v2, bs285, h, t, "dec-map-v285-p")
		testDeepEqualErr(v285v1, v285v2, t, "equal-map-v285-p")
	}

	for _, v := range []map[bool]bool{nil, map[bool]bool{}, map[bool]bool{true: true}} {
		// fmt.Printf(">>>> running mammoth map v286: %v\n", v)
		var v286v1, v286v2 map[bool]bool
		v286v1 = v
		bs286 := testMarshalErr(v286v1, h, t, "enc-map-v286")
		if v != nil {
			v286v2 = make(map[bool]bool, len(v))
		}
		testUnmarshalErr(v286v2, bs286, h, t, "dec-map-v286")
		testDeepEqualErr(v286v1, v286v2, t, "equal-map-v286")
		bs286 = testMarshalErr(&v286v1, h, t, "enc-map-v286-p")
		v286v2 = nil
		testUnmarshalErr(&v286v2, bs286, h, t, "dec-map-v286-p")
		testDeepEqualErr(v286v1, v286v2, t, "equal-map-v286-p")
	}

}

func doTestMammothMapsAndSlices(t *testing.T, h Handle) {
	doTestMammothSlices(t, h)
	doTestMammothMaps(t, h)
}
