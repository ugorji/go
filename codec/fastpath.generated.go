//go:build !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from fast-path.go.tmpl - DO NOT EDIT.

package codec

// Fast path functions try to create a fast path encode or decode implementation
// for common maps and slices.
//
// We define the functions and register them in this single file
// so as not to pollute the encode.go and decode.go, and create a dependency in there.
// This file can be omitted without causing a build failure.
//
// The advantage of fast paths is:
//	  - Many calls bypass reflection altogether
//
// Currently support
//	  - slice of all builtin types (numeric, bool, string, []byte)
//    - maps of builtin types to builtin or interface{} type, EXCEPT FOR
//      keys of type uintptr, int8/16/32, uint16/32, float32/64, bool, interface{}
//      AND values of type type int8/16/32, uint16/32
// This should provide adequate "typical" implementations.
//
// Note that fast track decode functions must handle values for which an address cannot be obtained.
// For example:
//	 m2 := map[string]int{}
//	 p2 := []interface{}{m2}
//	 // decoding into p2 will bomb if fast track functions do not treat like unaddressable.
//

import (
	"reflect"
	"slices"
	"sort"
)

const fastpathEnabled = true

type fastpathE[T encDriver] struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoder[T], *encFnInfo, reflect.Value)
}
type fastpathD[T decDriver] struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoder[T], *decFnInfo, reflect.Value)
}
type fastpathEs[T encDriver] [56]fastpathE[T]
type fastpathDs[T decDriver] [56]fastpathD[T]

type fastpathET[T encDriver] struct{}
type fastpathDT[T decDriver] struct{}

type fastpathARtid [56]uintptr

type fastpathRtRtid struct {
	rtid uintptr
	rt   reflect.Type
}
type fastpathARtRtid [56]fastpathRtRtid

var fastpathAvRtid fastpathARtid
var fastpathAvRtRtid fastpathARtRtid

func fastpathAvIndex(rtid uintptr) (i uint, ok bool) {
	// use binary search to grab the index (adapted from sort/search.go)
	// Note: we use goto (instead of for loop) so this can be inlined.
	// h, i, j := 0, 0, 56
	var h uint
	var j uint = 56
LOOP:
	if i < j {
		h = (i + j) >> 1 // avoid overflow when computing h // h = i + (j-i)/2
		if fastpathAvRtid[h] < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	return i, i < 56 && fastpathAvRtid[i] == rtid
}

func init() {
	var i uint = 0
	fn := func(v interface{}) {
		xrt := reflect.TypeOf(v)
		xrtid := rt2id(xrt)
		fastpathAvRtid[i] = xrtid
		fastpathAvRtRtid[i] = fastpathRtRtid{rtid: xrtid, rt: xrt}
		i++
	}

	fn([]interface{}(nil))
	fn([]string(nil))
	fn([][]byte(nil))
	fn([]float32(nil))
	fn([]float64(nil))
	fn([]uint8(nil))
	fn([]uint64(nil))
	fn([]int(nil))
	fn([]int32(nil))
	fn([]int64(nil))
	fn([]bool(nil))

	fn(map[string]interface{}(nil))
	fn(map[string]string(nil))
	fn(map[string][]byte(nil))
	fn(map[string]uint8(nil))
	fn(map[string]uint64(nil))
	fn(map[string]int(nil))
	fn(map[string]int32(nil))
	fn(map[string]float64(nil))
	fn(map[string]bool(nil))
	fn(map[uint8]interface{}(nil))
	fn(map[uint8]string(nil))
	fn(map[uint8][]byte(nil))
	fn(map[uint8]uint8(nil))
	fn(map[uint8]uint64(nil))
	fn(map[uint8]int(nil))
	fn(map[uint8]int32(nil))
	fn(map[uint8]float64(nil))
	fn(map[uint8]bool(nil))
	fn(map[uint64]interface{}(nil))
	fn(map[uint64]string(nil))
	fn(map[uint64][]byte(nil))
	fn(map[uint64]uint8(nil))
	fn(map[uint64]uint64(nil))
	fn(map[uint64]int(nil))
	fn(map[uint64]int32(nil))
	fn(map[uint64]float64(nil))
	fn(map[uint64]bool(nil))
	fn(map[int]interface{}(nil))
	fn(map[int]string(nil))
	fn(map[int][]byte(nil))
	fn(map[int]uint8(nil))
	fn(map[int]uint64(nil))
	fn(map[int]int(nil))
	fn(map[int]int32(nil))
	fn(map[int]float64(nil))
	fn(map[int]bool(nil))
	fn(map[int32]interface{}(nil))
	fn(map[int32]string(nil))
	fn(map[int32][]byte(nil))
	fn(map[int32]uint8(nil))
	fn(map[int32]uint64(nil))
	fn(map[int32]int(nil))
	fn(map[int32]int32(nil))
	fn(map[int32]float64(nil))
	fn(map[int32]bool(nil))

	sort.Slice(fastpathAvRtid[:], func(i, j int) bool { return fastpathAvRtid[i] < fastpathAvRtid[j] })
	sort.Slice(fastpathAvRtRtid[:], func(i, j int) bool { return fastpathAvRtRtid[i].rtid < fastpathAvRtRtid[j].rtid })
}

func (helperEncDriver[T]) fastpathEList() *fastpathEs[T] {
	var i uint = 0
	var s fastpathEs[T]
	fn := func(v interface{}, fe func(*encoder[T], *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathE[T]{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoder[T]).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoder[T]).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoder[T]).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoder[T]).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoder[T]).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoder[T]).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoder[T]).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoder[T]).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoder[T]).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoder[T]).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoder[T]).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoder[T]).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoder[T]).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoder[T]).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoder[T]).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoder[T]).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoder[T]).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoder[T]).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoder[T]).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoder[T]).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoder[T]).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoder[T]).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoder[T]).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoder[T]).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoder[T]).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoder[T]).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoder[T]).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoder[T]).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoder[T]).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoder[T]).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoder[T]).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoder[T]).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoder[T]).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoder[T]).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoder[T]).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoder[T]).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoder[T]).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoder[T]).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoder[T]).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoder[T]).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoder[T]).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoder[T]).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoder[T]).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoder[T]).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoder[T]).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoder[T]).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoder[T]).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoder[T]).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoder[T]).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoder[T]).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoder[T]).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoder[T]).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoder[T]).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoder[T]).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoder[T]).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoder[T]).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriver[T]) fastpathDList() *fastpathDs[T] {
	var i uint = 0
	var s fastpathDs[T]
	fn := func(v interface{}, fd func(*decoder[T], *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathD[T]{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoder[T]).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoder[T]).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoder[T]).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoder[T]).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoder[T]).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoder[T]).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoder[T]).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoder[T]).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoder[T]).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoder[T]).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoder[T]).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoder[T]).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoder[T]).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoder[T]).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoder[T]).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoder[T]).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoder[T]).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoder[T]).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoder[T]).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoder[T]).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoder[T]).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoder[T]).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoder[T]).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoder[T]).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoder[T]).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoder[T]).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoder[T]).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoder[T]).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoder[T]).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoder[T]).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoder[T]).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoder[T]).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoder[T]).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoder[T]).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoder[T]).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoder[T]).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoder[T]).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoder[T]).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoder[T]).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoder[T]).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoder[T]).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoder[T]).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoder[T]).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoder[T]).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoder[T]).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoder[T]).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoder[T]).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoder[T]).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoder[T]).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoder[T]).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoder[T]).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoder[T]).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoder[T]).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoder[T]).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoder[T]).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoder[T]).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

// -- encode

// -- -- fast path type switch
func (helperEncDriver[T]) fastpathEncodeTypeSwitch(iv interface{}, e *encoder[T]) bool {
	var ft fastpathET[T]
	switch v := iv.(type) {
	case []interface{}:
		ft.EncSliceIntfV(v, e)
	case *[]interface{}:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceIntfV(*v, e)
		}
	case []string:
		ft.EncSliceStringV(v, e)
	case *[]string:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceStringV(*v, e)
		}
	case [][]byte:
		ft.EncSliceBytesV(v, e)
	case *[][]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceBytesV(*v, e)
		}
	case []float32:
		ft.EncSliceFloat32V(v, e)
	case *[]float32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceFloat32V(*v, e)
		}
	case []float64:
		ft.EncSliceFloat64V(v, e)
	case *[]float64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceFloat64V(*v, e)
		}
	case []uint8:
		ft.EncSliceUint8V(v, e)
	case *[]uint8:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceUint8V(*v, e)
		}
	case []uint64:
		ft.EncSliceUint64V(v, e)
	case *[]uint64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceUint64V(*v, e)
		}
	case []int:
		ft.EncSliceIntV(v, e)
	case *[]int:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceIntV(*v, e)
		}
	case []int32:
		ft.EncSliceInt32V(v, e)
	case *[]int32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceInt32V(*v, e)
		}
	case []int64:
		ft.EncSliceInt64V(v, e)
	case *[]int64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceInt64V(*v, e)
		}
	case []bool:
		ft.EncSliceBoolV(v, e)
	case *[]bool:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncSliceBoolV(*v, e)
		}
	case map[string]interface{}:
		ft.EncMapStringIntfV(v, e)
	case *map[string]interface{}:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringIntfV(*v, e)
		}
	case map[string]string:
		ft.EncMapStringStringV(v, e)
	case *map[string]string:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringStringV(*v, e)
		}
	case map[string][]byte:
		ft.EncMapStringBytesV(v, e)
	case *map[string][]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringBytesV(*v, e)
		}
	case map[string]uint8:
		ft.EncMapStringUint8V(v, e)
	case *map[string]uint8:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringUint8V(*v, e)
		}
	case map[string]uint64:
		ft.EncMapStringUint64V(v, e)
	case *map[string]uint64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringUint64V(*v, e)
		}
	case map[string]int:
		ft.EncMapStringIntV(v, e)
	case *map[string]int:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringIntV(*v, e)
		}
	case map[string]int32:
		ft.EncMapStringInt32V(v, e)
	case *map[string]int32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringInt32V(*v, e)
		}
	case map[string]float64:
		ft.EncMapStringFloat64V(v, e)
	case *map[string]float64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringFloat64V(*v, e)
		}
	case map[string]bool:
		ft.EncMapStringBoolV(v, e)
	case *map[string]bool:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapStringBoolV(*v, e)
		}
	case map[uint8]interface{}:
		ft.EncMapUint8IntfV(v, e)
	case *map[uint8]interface{}:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8IntfV(*v, e)
		}
	case map[uint8]string:
		ft.EncMapUint8StringV(v, e)
	case *map[uint8]string:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8StringV(*v, e)
		}
	case map[uint8][]byte:
		ft.EncMapUint8BytesV(v, e)
	case *map[uint8][]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8BytesV(*v, e)
		}
	case map[uint8]uint8:
		ft.EncMapUint8Uint8V(v, e)
	case *map[uint8]uint8:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8Uint8V(*v, e)
		}
	case map[uint8]uint64:
		ft.EncMapUint8Uint64V(v, e)
	case *map[uint8]uint64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8Uint64V(*v, e)
		}
	case map[uint8]int:
		ft.EncMapUint8IntV(v, e)
	case *map[uint8]int:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8IntV(*v, e)
		}
	case map[uint8]int32:
		ft.EncMapUint8Int32V(v, e)
	case *map[uint8]int32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8Int32V(*v, e)
		}
	case map[uint8]float64:
		ft.EncMapUint8Float64V(v, e)
	case *map[uint8]float64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8Float64V(*v, e)
		}
	case map[uint8]bool:
		ft.EncMapUint8BoolV(v, e)
	case *map[uint8]bool:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint8BoolV(*v, e)
		}
	case map[uint64]interface{}:
		ft.EncMapUint64IntfV(v, e)
	case *map[uint64]interface{}:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64IntfV(*v, e)
		}
	case map[uint64]string:
		ft.EncMapUint64StringV(v, e)
	case *map[uint64]string:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64StringV(*v, e)
		}
	case map[uint64][]byte:
		ft.EncMapUint64BytesV(v, e)
	case *map[uint64][]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64BytesV(*v, e)
		}
	case map[uint64]uint8:
		ft.EncMapUint64Uint8V(v, e)
	case *map[uint64]uint8:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64Uint8V(*v, e)
		}
	case map[uint64]uint64:
		ft.EncMapUint64Uint64V(v, e)
	case *map[uint64]uint64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64Uint64V(*v, e)
		}
	case map[uint64]int:
		ft.EncMapUint64IntV(v, e)
	case *map[uint64]int:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64IntV(*v, e)
		}
	case map[uint64]int32:
		ft.EncMapUint64Int32V(v, e)
	case *map[uint64]int32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64Int32V(*v, e)
		}
	case map[uint64]float64:
		ft.EncMapUint64Float64V(v, e)
	case *map[uint64]float64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64Float64V(*v, e)
		}
	case map[uint64]bool:
		ft.EncMapUint64BoolV(v, e)
	case *map[uint64]bool:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapUint64BoolV(*v, e)
		}
	case map[int]interface{}:
		ft.EncMapIntIntfV(v, e)
	case *map[int]interface{}:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntIntfV(*v, e)
		}
	case map[int]string:
		ft.EncMapIntStringV(v, e)
	case *map[int]string:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntStringV(*v, e)
		}
	case map[int][]byte:
		ft.EncMapIntBytesV(v, e)
	case *map[int][]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntBytesV(*v, e)
		}
	case map[int]uint8:
		ft.EncMapIntUint8V(v, e)
	case *map[int]uint8:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntUint8V(*v, e)
		}
	case map[int]uint64:
		ft.EncMapIntUint64V(v, e)
	case *map[int]uint64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntUint64V(*v, e)
		}
	case map[int]int:
		ft.EncMapIntIntV(v, e)
	case *map[int]int:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntIntV(*v, e)
		}
	case map[int]int32:
		ft.EncMapIntInt32V(v, e)
	case *map[int]int32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntInt32V(*v, e)
		}
	case map[int]float64:
		ft.EncMapIntFloat64V(v, e)
	case *map[int]float64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntFloat64V(*v, e)
		}
	case map[int]bool:
		ft.EncMapIntBoolV(v, e)
	case *map[int]bool:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapIntBoolV(*v, e)
		}
	case map[int32]interface{}:
		ft.EncMapInt32IntfV(v, e)
	case *map[int32]interface{}:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32IntfV(*v, e)
		}
	case map[int32]string:
		ft.EncMapInt32StringV(v, e)
	case *map[int32]string:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32StringV(*v, e)
		}
	case map[int32][]byte:
		ft.EncMapInt32BytesV(v, e)
	case *map[int32][]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32BytesV(*v, e)
		}
	case map[int32]uint8:
		ft.EncMapInt32Uint8V(v, e)
	case *map[int32]uint8:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32Uint8V(*v, e)
		}
	case map[int32]uint64:
		ft.EncMapInt32Uint64V(v, e)
	case *map[int32]uint64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32Uint64V(*v, e)
		}
	case map[int32]int:
		ft.EncMapInt32IntV(v, e)
	case *map[int32]int:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32IntV(*v, e)
		}
	case map[int32]int32:
		ft.EncMapInt32Int32V(v, e)
	case *map[int32]int32:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32Int32V(*v, e)
		}
	case map[int32]float64:
		ft.EncMapInt32Float64V(v, e)
	case *map[int32]float64:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32Float64V(*v, e)
		}
	case map[int32]bool:
		ft.EncMapInt32BoolV(v, e)
	case *map[int32]bool:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			ft.EncMapInt32BoolV(*v, e)
		}
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

// -- -- fast path functions
func (e *encoder[T]) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []interface{}
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]interface{})
	}
	if f.ti.mbs {
		ft.EncAsMapSliceIntfV(v, e)
	} else {
		ft.EncSliceIntfV(v, e)
	}
}
func (fastpathET[T]) EncSliceIntfV(v []interface{}, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceIntfV(v []interface{}, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.encode(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []string
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]string)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceStringV(v, e)
	} else {
		ft.EncSliceStringV(v, e)
	}
}
func (fastpathET[T]) EncSliceStringV(v []string, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceStringV(v []string, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeString(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v [][]byte
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([][]byte)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceBytesV(v, e)
	} else {
		ft.EncSliceBytesV(v, e)
	}
}
func (fastpathET[T]) EncSliceBytesV(v [][]byte, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceBytesV(v [][]byte, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []float32
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]float32)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceFloat32V(v, e)
	} else {
		ft.EncSliceFloat32V(v, e)
	}
}
func (fastpathET[T]) EncSliceFloat32V(v []float32, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceFloat32V(v []float32, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat32(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []float64
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]float64)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceFloat64V(v, e)
	} else {
		ft.EncSliceFloat64V(v, e)
	}
}
func (fastpathET[T]) EncSliceFloat64V(v []float64, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceFloat64V(v []float64, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat64(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []uint8
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]uint8)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceUint8V(v, e)
	} else {
		ft.EncSliceUint8V(v, e)
	}
}
func (fastpathET[T]) EncSliceUint8V(v []uint8, e *encoder[T]) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathET[T]) EncAsMapSliceUint8V(v []uint8, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(uint64(v[j]))
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []uint64
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]uint64)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceUint64V(v, e)
	} else {
		ft.EncSliceUint64V(v, e)
	}
}
func (fastpathET[T]) EncSliceUint64V(v []uint64, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceUint64V(v []uint64, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []int
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]int)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceIntV(v, e)
	} else {
		ft.EncSliceIntV(v, e)
	}
}
func (fastpathET[T]) EncSliceIntV(v []int, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceIntV(v []int, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []int32
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]int32)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceInt32V(v, e)
	} else {
		ft.EncSliceInt32V(v, e)
	}
}
func (fastpathET[T]) EncSliceInt32V(v []int32, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceInt32V(v []int32, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []int64
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]int64)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceInt64V(v, e)
	} else {
		ft.EncSliceInt64V(v, e)
	}
}
func (fastpathET[T]) EncSliceInt64V(v []int64, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceInt64V(v []int64, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []bool
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]bool)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceBoolV(v, e)
	} else {
		ft.EncSliceBoolV(v, e)
	}
}
func (fastpathET[T]) EncSliceBoolV(v []bool, e *encoder[T]) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathET[T]) EncAsMapSliceBoolV(v []bool, e *encoder[T]) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeBool(v[j])
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathET[T]) EncMapStringIntfV(v map[string]interface{}, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathET[T]) EncMapStringStringV(v map[string]string, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeString(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathET[T]) EncMapStringBytesV(v map[string][]byte, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathET[T]) EncMapStringUint8V(v map[string]uint8, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathET[T]) EncMapStringUint64V(v map[string]uint64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathET[T]) EncMapStringIntV(v map[string]int, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathET[T]) EncMapStringInt32V(v map[string]int32, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathET[T]) EncMapStringFloat64V(v map[string]float64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathET[T]) EncMapStringBoolV(v map[string]bool, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathET[T]) EncMapUint8IntfV(v map[uint8]interface{}, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathET[T]) EncMapUint8StringV(v map[uint8]string, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeString(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathET[T]) EncMapUint8BytesV(v map[uint8][]byte, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathET[T]) EncMapUint8Uint8V(v map[uint8]uint8, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathET[T]) EncMapUint8Uint64V(v map[uint8]uint64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathET[T]) EncMapUint8IntV(v map[uint8]int, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathET[T]) EncMapUint8Int32V(v map[uint8]int32, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathET[T]) EncMapUint8Float64V(v map[uint8]float64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathET[T]) EncMapUint8BoolV(v map[uint8]bool, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathET[T]) EncMapUint64IntfV(v map[uint64]interface{}, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathET[T]) EncMapUint64StringV(v map[uint64]string, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeString(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathET[T]) EncMapUint64BytesV(v map[uint64][]byte, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathET[T]) EncMapUint64Uint8V(v map[uint64]uint8, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathET[T]) EncMapUint64Uint64V(v map[uint64]uint64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathET[T]) EncMapUint64IntV(v map[uint64]int, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathET[T]) EncMapUint64Int32V(v map[uint64]int32, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathET[T]) EncMapUint64Float64V(v map[uint64]float64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathET[T]) EncMapUint64BoolV(v map[uint64]bool, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathET[T]) EncMapIntIntfV(v map[int]interface{}, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathET[T]) EncMapIntStringV(v map[int]string, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathET[T]) EncMapIntBytesV(v map[int][]byte, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathET[T]) EncMapIntUint8V(v map[int]uint8, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathET[T]) EncMapIntUint64V(v map[int]uint64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathET[T]) EncMapIntIntV(v map[int]int, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathET[T]) EncMapIntInt32V(v map[int]int32, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathET[T]) EncMapIntFloat64V(v map[int]float64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathET[T]) EncMapIntBoolV(v map[int]bool, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathET[T]) EncMapInt32IntfV(v map[int32]interface{}, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.encode(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.encode(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathET[T]) EncMapInt32StringV(v map[int32]string, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathET[T]) EncMapInt32BytesV(v map[int32][]byte, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeStringBytesRaw(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathET[T]) EncMapInt32Uint8V(v map[int32]uint8, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathET[T]) EncMapInt32Uint64V(v map[int32]uint64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathET[T]) EncMapInt32IntV(v map[int32]int, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathET[T]) EncMapInt32Int32V(v map[int32]int32, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathET[T]) EncMapInt32Float64V(v map[int32]float64, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
		}
	}
	e.mapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathET[T]) EncMapInt32BoolV(v map[int32]bool, e *encoder[T]) {
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		var i uint
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for _, k2 := range v2 {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		for k2, v2 := range v {
			e.mapElemKey()
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
		}
	}
	e.mapEnd()
}

// -- decode

// -- -- fast path type switch
func (helperDecDriver[T]) fastpathDecodeTypeSwitch(iv interface{}, d *decoder[T]) bool {
	var ft fastpathDT[T]
	var changed bool
	var containerLen int
	switch v := iv.(type) {
	case []interface{}:
		ft.DecSliceIntfN(v, d)
	case *[]interface{}:
		var v2 []interface{}
		if v2, changed = ft.DecSliceIntfY(*v, d); changed {
			*v = v2
		}
	case []string:
		ft.DecSliceStringN(v, d)
	case *[]string:
		var v2 []string
		if v2, changed = ft.DecSliceStringY(*v, d); changed {
			*v = v2
		}
	case [][]byte:
		ft.DecSliceBytesN(v, d)
	case *[][]byte:
		var v2 [][]byte
		if v2, changed = ft.DecSliceBytesY(*v, d); changed {
			*v = v2
		}
	case []float32:
		ft.DecSliceFloat32N(v, d)
	case *[]float32:
		var v2 []float32
		if v2, changed = ft.DecSliceFloat32Y(*v, d); changed {
			*v = v2
		}
	case []float64:
		ft.DecSliceFloat64N(v, d)
	case *[]float64:
		var v2 []float64
		if v2, changed = ft.DecSliceFloat64Y(*v, d); changed {
			*v = v2
		}
	case []uint8:
		ft.DecSliceUint8N(v, d)
	case *[]uint8:
		var v2 []uint8
		if v2, changed = ft.DecSliceUint8Y(*v, d); changed {
			*v = v2
		}
	case []uint64:
		ft.DecSliceUint64N(v, d)
	case *[]uint64:
		var v2 []uint64
		if v2, changed = ft.DecSliceUint64Y(*v, d); changed {
			*v = v2
		}
	case []int:
		ft.DecSliceIntN(v, d)
	case *[]int:
		var v2 []int
		if v2, changed = ft.DecSliceIntY(*v, d); changed {
			*v = v2
		}
	case []int32:
		ft.DecSliceInt32N(v, d)
	case *[]int32:
		var v2 []int32
		if v2, changed = ft.DecSliceInt32Y(*v, d); changed {
			*v = v2
		}
	case []int64:
		ft.DecSliceInt64N(v, d)
	case *[]int64:
		var v2 []int64
		if v2, changed = ft.DecSliceInt64Y(*v, d); changed {
			*v = v2
		}
	case []bool:
		ft.DecSliceBoolN(v, d)
	case *[]bool:
		var v2 []bool
		if v2, changed = ft.DecSliceBoolY(*v, d); changed {
			*v = v2
		}
	case map[string]interface{}:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringIntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]interface{}:
		ft.DecMapStringIntfX(v, d)
	case map[string]string:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringStringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]string:
		ft.DecMapStringStringX(v, d)
	case map[string][]byte:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringBytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string][]byte:
		ft.DecMapStringBytesX(v, d)
	case map[string]uint8:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringUint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]uint8:
		ft.DecMapStringUint8X(v, d)
	case map[string]uint64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringUint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]uint64:
		ft.DecMapStringUint64X(v, d)
	case map[string]int:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringIntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]int:
		ft.DecMapStringIntX(v, d)
	case map[string]int32:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringInt32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]int32:
		ft.DecMapStringInt32X(v, d)
	case map[string]float64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringFloat64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]float64:
		ft.DecMapStringFloat64X(v, d)
	case map[string]bool:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringBoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]bool:
		ft.DecMapStringBoolX(v, d)
	case map[uint8]interface{}:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8IntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]interface{}:
		ft.DecMapUint8IntfX(v, d)
	case map[uint8]string:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8StringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]string:
		ft.DecMapUint8StringX(v, d)
	case map[uint8][]byte:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8BytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8][]byte:
		ft.DecMapUint8BytesX(v, d)
	case map[uint8]uint8:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Uint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]uint8:
		ft.DecMapUint8Uint8X(v, d)
	case map[uint8]uint64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Uint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]uint64:
		ft.DecMapUint8Uint64X(v, d)
	case map[uint8]int:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8IntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]int:
		ft.DecMapUint8IntX(v, d)
	case map[uint8]int32:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Int32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]int32:
		ft.DecMapUint8Int32X(v, d)
	case map[uint8]float64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Float64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]float64:
		ft.DecMapUint8Float64X(v, d)
	case map[uint8]bool:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8BoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]bool:
		ft.DecMapUint8BoolX(v, d)
	case map[uint64]interface{}:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64IntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]interface{}:
		ft.DecMapUint64IntfX(v, d)
	case map[uint64]string:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64StringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]string:
		ft.DecMapUint64StringX(v, d)
	case map[uint64][]byte:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64BytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64][]byte:
		ft.DecMapUint64BytesX(v, d)
	case map[uint64]uint8:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Uint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]uint8:
		ft.DecMapUint64Uint8X(v, d)
	case map[uint64]uint64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Uint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]uint64:
		ft.DecMapUint64Uint64X(v, d)
	case map[uint64]int:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64IntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]int:
		ft.DecMapUint64IntX(v, d)
	case map[uint64]int32:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Int32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]int32:
		ft.DecMapUint64Int32X(v, d)
	case map[uint64]float64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Float64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]float64:
		ft.DecMapUint64Float64X(v, d)
	case map[uint64]bool:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64BoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]bool:
		ft.DecMapUint64BoolX(v, d)
	case map[int]interface{}:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntIntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]interface{}:
		ft.DecMapIntIntfX(v, d)
	case map[int]string:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntStringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]string:
		ft.DecMapIntStringX(v, d)
	case map[int][]byte:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntBytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int][]byte:
		ft.DecMapIntBytesX(v, d)
	case map[int]uint8:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntUint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]uint8:
		ft.DecMapIntUint8X(v, d)
	case map[int]uint64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntUint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]uint64:
		ft.DecMapIntUint64X(v, d)
	case map[int]int:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntIntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]int:
		ft.DecMapIntIntX(v, d)
	case map[int]int32:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntInt32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]int32:
		ft.DecMapIntInt32X(v, d)
	case map[int]float64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntFloat64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]float64:
		ft.DecMapIntFloat64X(v, d)
	case map[int]bool:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntBoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]bool:
		ft.DecMapIntBoolX(v, d)
	case map[int32]interface{}:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32IntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]interface{}:
		ft.DecMapInt32IntfX(v, d)
	case map[int32]string:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32StringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]string:
		ft.DecMapInt32StringX(v, d)
	case map[int32][]byte:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32BytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32][]byte:
		ft.DecMapInt32BytesX(v, d)
	case map[int32]uint8:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Uint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]uint8:
		ft.DecMapInt32Uint8X(v, d)
	case map[int32]uint64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Uint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]uint64:
		ft.DecMapInt32Uint64X(v, d)
	case map[int32]int:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32IntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]int:
		ft.DecMapInt32IntX(v, d)
	case map[int32]int32:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Int32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]int32:
		ft.DecMapInt32Int32X(v, d)
	case map[int32]float64:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Float64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]float64:
		ft.DecMapInt32Float64X(v, d)
	case map[int32]bool:
		containerLen = d.mapStart(d.d.ReadMapStart())
		if containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32BoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]bool:
		ft.DecMapInt32BoolX(v, d)
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

func fastpathDecodeSetZeroTypeSwitch(iv interface{}) bool {
	switch v := iv.(type) {
	case *[]interface{}:
		*v = nil
	case *[]string:
		*v = nil
	case *[][]byte:
		*v = nil
	case *[]float32:
		*v = nil
	case *[]float64:
		*v = nil
	case *[]uint8:
		*v = nil
	case *[]uint64:
		*v = nil
	case *[]int:
		*v = nil
	case *[]int32:
		*v = nil
	case *[]int64:
		*v = nil
	case *[]bool:
		*v = nil

	case *map[string]interface{}:
		*v = nil
	case *map[string]string:
		*v = nil
	case *map[string][]byte:
		*v = nil
	case *map[string]uint8:
		*v = nil
	case *map[string]uint64:
		*v = nil
	case *map[string]int:
		*v = nil
	case *map[string]int32:
		*v = nil
	case *map[string]float64:
		*v = nil
	case *map[string]bool:
		*v = nil
	case *map[uint8]interface{}:
		*v = nil
	case *map[uint8]string:
		*v = nil
	case *map[uint8][]byte:
		*v = nil
	case *map[uint8]uint8:
		*v = nil
	case *map[uint8]uint64:
		*v = nil
	case *map[uint8]int:
		*v = nil
	case *map[uint8]int32:
		*v = nil
	case *map[uint8]float64:
		*v = nil
	case *map[uint8]bool:
		*v = nil
	case *map[uint64]interface{}:
		*v = nil
	case *map[uint64]string:
		*v = nil
	case *map[uint64][]byte:
		*v = nil
	case *map[uint64]uint8:
		*v = nil
	case *map[uint64]uint64:
		*v = nil
	case *map[uint64]int:
		*v = nil
	case *map[uint64]int32:
		*v = nil
	case *map[uint64]float64:
		*v = nil
	case *map[uint64]bool:
		*v = nil
	case *map[int]interface{}:
		*v = nil
	case *map[int]string:
		*v = nil
	case *map[int][]byte:
		*v = nil
	case *map[int]uint8:
		*v = nil
	case *map[int]uint64:
		*v = nil
	case *map[int]int:
		*v = nil
	case *map[int]int32:
		*v = nil
	case *map[int]float64:
		*v = nil
	case *map[int]bool:
		*v = nil
	case *map[int32]interface{}:
		*v = nil
	case *map[int32]string:
		*v = nil
	case *map[int32][]byte:
		*v = nil
	case *map[int32]uint8:
		*v = nil
	case *map[int32]uint64:
		*v = nil
	case *map[int32]int:
		*v = nil
	case *map[int32]int32:
		*v = nil
	case *map[int32]float64:
		*v = nil
	case *map[int32]bool:
		*v = nil

	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

// -- -- fast path functions

func (d *decoder[T]) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []interface{}
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]interface{})
		var changed bool
		if v, changed = ft.DecSliceIntfY(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceIntfN(v, d)
	default:
		ft.DecSliceIntfN(rv2i(rv).([]interface{}), d)
	}
}
func (f fastpathDT[T]) DecSliceIntfX(vp *[]interface{}, d *decoder[T]) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceIntfY(v []interface{}, d *decoder[T]) (v2 []interface{}, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []interface{}{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]interface{}, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			v = make([]interface{}, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, nil)
			changed = true
		}
		slh.ElemContainerState(j)
		d.decode(&v[uint(j)])
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []interface{}{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceIntfN(v []interface{}, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		d.decode(&v[uint(j)])
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []string
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]string)
		var changed bool
		if v, changed = ft.DecSliceStringY(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceStringN(v, d)
	default:
		ft.DecSliceStringN(rv2i(rv).([]string), d)
	}
}
func (f fastpathDT[T]) DecSliceStringX(vp *[]string, d *decoder[T]) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceStringY(v []string, d *decoder[T]) (v2 []string, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []string{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]string, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			v = make([]string, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, "")
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.stringZC(d.d.DecodeStringAsBytes())
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []string{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceStringN(v []string, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.stringZC(d.d.DecodeStringAsBytes())
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v [][]byte
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[][]byte)
		var changed bool
		if v, changed = ft.DecSliceBytesY(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceBytesN(v, d)
	default:
		ft.DecSliceBytesN(rv2i(rv).([][]byte), d)
	}
}
func (f fastpathDT[T]) DecSliceBytesX(vp *[][]byte, d *decoder[T]) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceBytesY(v [][]byte, d *decoder[T]) (v2 [][]byte, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = [][]byte{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 24)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([][]byte, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 24)
			v = make([][]byte, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, nil)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeBytes(zeroByteSlice)
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = [][]byte{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceBytesN(v [][]byte, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeBytes(zeroByteSlice)
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []float32
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]float32)
		var changed bool
		if v, changed = ft.DecSliceFloat32Y(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceFloat32N(v, d)
	default:
		ft.DecSliceFloat32N(rv2i(rv).([]float32), d)
	}
}
func (f fastpathDT[T]) DecSliceFloat32X(vp *[]float32, d *decoder[T]) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceFloat32Y(v []float32, d *decoder[T]) (v2 []float32, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []float32{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]float32, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			v = make([]float32, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = float32(d.d.DecodeFloat32())
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []float32{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceFloat32N(v []float32, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = float32(d.d.DecodeFloat32())
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []float64
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]float64)
		var changed bool
		if v, changed = ft.DecSliceFloat64Y(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceFloat64N(v, d)
	default:
		ft.DecSliceFloat64N(rv2i(rv).([]float64), d)
	}
}
func (f fastpathDT[T]) DecSliceFloat64X(vp *[]float64, d *decoder[T]) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceFloat64Y(v []float64, d *decoder[T]) (v2 []float64, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []float64{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]float64, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			v = make([]float64, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeFloat64()
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []float64{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceFloat64N(v []float64, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeFloat64()
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []uint8
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]uint8)
		var changed bool
		if v, changed = ft.DecSliceUint8Y(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceUint8N(v, d)
	default:
		ft.DecSliceUint8N(rv2i(rv).([]uint8), d)
	}
}
func (f fastpathDT[T]) DecSliceUint8X(vp *[]uint8, d *decoder[T]) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceUint8Y(v []uint8, d *decoder[T]) (v2 []uint8, changed bool) {
	switch d.d.ContainerType() {
	case valueTypeNil, valueTypeMap:
		break
	default:
		v2 = d.decodeBytesInto(v[:len(v):len(v)])
		changed = !(len(v2) > 0 && len(v2) == len(v) && &v2[0] == &v[0]) // not same slice
		return
	}
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []uint8{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint8, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			v = make([]uint8, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []uint8{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceUint8N(v []uint8, d *decoder[T]) {
	switch d.d.ContainerType() {
	case valueTypeNil, valueTypeMap:
		break
	default:
		v2 := d.decodeBytesInto(v[:len(v):len(v)])
		if !(len(v2) > 0 && len(v2) == len(v) && &v2[0] == &v[0]) { // not same slice
			copy(v, v2)
		}
		return
	}
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []uint64
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]uint64)
		var changed bool
		if v, changed = ft.DecSliceUint64Y(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceUint64N(v, d)
	default:
		ft.DecSliceUint64N(rv2i(rv).([]uint64), d)
	}
}
func (f fastpathDT[T]) DecSliceUint64X(vp *[]uint64, d *decoder[T]) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceUint64Y(v []uint64, d *decoder[T]) (v2 []uint64, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []uint64{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]uint64, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			v = make([]uint64, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeUint64()
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []uint64{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceUint64N(v []uint64, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeUint64()
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []int
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]int)
		var changed bool
		if v, changed = ft.DecSliceIntY(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceIntN(v, d)
	default:
		ft.DecSliceIntN(rv2i(rv).([]int), d)
	}
}
func (f fastpathDT[T]) DecSliceIntX(vp *[]int, d *decoder[T]) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceIntY(v []int, d *decoder[T]) (v2 []int, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []int{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			v = make([]int, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []int{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceIntN(v []int, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []int32
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]int32)
		var changed bool
		if v, changed = ft.DecSliceInt32Y(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceInt32N(v, d)
	default:
		ft.DecSliceInt32N(rv2i(rv).([]int32), d)
	}
}
func (f fastpathDT[T]) DecSliceInt32X(vp *[]int32, d *decoder[T]) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceInt32Y(v []int32, d *decoder[T]) (v2 []int32, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []int32{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int32, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 4)
			v = make([]int32, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []int32{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceInt32N(v []int32, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []int64
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]int64)
		var changed bool
		if v, changed = ft.DecSliceInt64Y(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceInt64N(v, d)
	default:
		ft.DecSliceInt64N(rv2i(rv).([]int64), d)
	}
}
func (f fastpathDT[T]) DecSliceInt64X(vp *[]int64, d *decoder[T]) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceInt64Y(v []int64, d *decoder[T]) (v2 []int64, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []int64{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]int64, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 8)
			v = make([]int64, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, 0)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeInt64()
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []int64{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceInt64N(v []int64, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeInt64()
	}
	slh.End()
}

func (d *decoder[T]) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	var v []bool
	switch rv.Kind() {
	case reflect.Ptr:
		vp := rv2i(rv).(*[]bool)
		var changed bool
		if v, changed = ft.DecSliceBoolY(*vp, d); changed {
			*vp = v
		}
	case reflect.Array:
		rvGetSlice4Array(rv, &v)
		ft.DecSliceBoolN(v, d)
	default:
		ft.DecSliceBoolN(rv2i(rv).([]bool), d)
	}
}
func (f fastpathDT[T]) DecSliceBoolX(vp *[]bool, d *decoder[T]) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDT[T]) DecSliceBoolY(v []bool, d *decoder[T]) (v2 []bool, changed bool) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		if v == nil {
			return
		}
		return nil, true
	}
	if containerLenS == 0 {
		if v == nil {
			v = []bool{}
		} else if len(v) != 0 {
			v = v[:0]
		}
		slh.End()
		return v, true
	}
	hasLen := containerLenS > 0
	var xlen int
	if hasLen {
		if containerLenS > cap(v) {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			if xlen <= cap(v) {
				v = v[:uint(xlen)]
			} else {
				v = make([]bool, uint(xlen))
			}
			changed = true
		} else if containerLenS != len(v) {
			v = v[:containerLenS]
			changed = true
		}
	}
	var j int
	for j = 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 && len(v) == 0 { // means hasLen == false
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 1)
			v = make([]bool, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, false)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeBool()
	}
	if j < len(v) {
		v = v[:uint(j)]
		changed = true
	} else if j == 0 && v == nil {
		v = []bool{}
		changed = true
	}
	slh.End()
	return v, changed
}
func (fastpathDT[T]) DecSliceBoolN(v []bool, d *decoder[T]) {
	slh, containerLenS := d.decSliceHelperStart()
	if slh.IsNil {
		return
	}
	if containerLenS == 0 {
		slh.End()
		return
	}
	hasLen := containerLenS > 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j >= len(v) {
			slh.arrayCannotExpand(hasLen, len(v), j, containerLenS)
			return
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.d.DecodeBool()
	}
	slh.End()
}
func (d *decoder[T]) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]interface{})
		if *vp == nil {
			*vp = make(map[string]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			ft.DecMapStringIntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringIntfL(rv2i(rv).(map[string]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringIntfX(vp *map[string]interface{}, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 32))
	}
	if containerLen != 0 {
		f.DecMapStringIntfL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk string
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]string)
		if *vp == nil {
			*vp = make(map[string]string, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			ft.DecMapStringStringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringStringL(rv2i(rv).(map[string]string), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringStringX(vp *map[string]string, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]string, decInferLen(containerLen, d.h.MaxInitLen, 32))
	}
	if containerLen != 0 {
		f.DecMapStringStringL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringStringL(v map[string]string, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]string given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes())
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string][]byte)
		if *vp == nil {
			*vp = make(map[string][]byte, decInferLen(containerLen, d.h.MaxInitLen, 40))
		}
		if containerLen != 0 {
			ft.DecMapStringBytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringBytesL(rv2i(rv).(map[string][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringBytesX(vp *map[string][]byte, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string][]byte, decInferLen(containerLen, d.h.MaxInitLen, 40))
	}
	if containerLen != 0 {
		f.DecMapStringBytesL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk string
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.decodeBytesInto(mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]uint8)
		if *vp == nil {
			*vp = make(map[string]uint8, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			ft.DecMapStringUint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringUint8L(rv2i(rv).(map[string]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringUint8X(vp *map[string]uint8, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]uint8, decInferLen(containerLen, d.h.MaxInitLen, 17))
	}
	if containerLen != 0 {
		f.DecMapStringUint8L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]uint64)
		if *vp == nil {
			*vp = make(map[string]uint64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapStringUint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringUint64L(rv2i(rv).(map[string]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringUint64X(vp *map[string]uint64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]uint64, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapStringUint64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]int)
		if *vp == nil {
			*vp = make(map[string]int, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapStringIntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringIntL(rv2i(rv).(map[string]int), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringIntX(vp *map[string]int, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]int, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapStringIntL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringIntL(v map[string]int, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]int given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]int32)
		if *vp == nil {
			*vp = make(map[string]int32, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			ft.DecMapStringInt32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringInt32L(rv2i(rv).(map[string]int32), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringInt32X(vp *map[string]int32, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]int32, decInferLen(containerLen, d.h.MaxInitLen, 20))
	}
	if containerLen != 0 {
		f.DecMapStringInt32L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]float64)
		if *vp == nil {
			*vp = make(map[string]float64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapStringFloat64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringFloat64L(rv2i(rv).(map[string]float64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringFloat64X(vp *map[string]float64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]float64, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapStringFloat64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]bool)
		if *vp == nil {
			*vp = make(map[string]bool, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			ft.DecMapStringBoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringBoolL(rv2i(rv).(map[string]bool), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapStringBoolX(vp *map[string]bool, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[string]bool, decInferLen(containerLen, d.h.MaxInitLen, 17))
	}
	if containerLen != 0 {
		f.DecMapStringBoolL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]bool given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]interface{})
		if *vp == nil {
			*vp = make(map[uint8]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			ft.DecMapUint8IntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8IntfL(rv2i(rv).(map[uint8]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 17))
	}
	if containerLen != 0 {
		f.DecMapUint8IntfL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk uint8
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]string)
		if *vp == nil {
			*vp = make(map[uint8]string, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			ft.DecMapUint8StringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8StringL(rv2i(rv).(map[uint8]string), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8StringX(vp *map[uint8]string, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]string, decInferLen(containerLen, d.h.MaxInitLen, 17))
	}
	if containerLen != 0 {
		f.DecMapUint8StringL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]string given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes())
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8][]byte)
		if *vp == nil {
			*vp = make(map[uint8][]byte, decInferLen(containerLen, d.h.MaxInitLen, 25))
		}
		if containerLen != 0 {
			ft.DecMapUint8BytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8BytesL(rv2i(rv).(map[uint8][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8][]byte, decInferLen(containerLen, d.h.MaxInitLen, 25))
	}
	if containerLen != 0 {
		f.DecMapUint8BytesL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk uint8
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.decodeBytesInto(mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]uint8)
		if *vp == nil {
			*vp = make(map[uint8]uint8, decInferLen(containerLen, d.h.MaxInitLen, 2))
		}
		if containerLen != 0 {
			ft.DecMapUint8Uint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Uint8L(rv2i(rv).(map[uint8]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]uint8, decInferLen(containerLen, d.h.MaxInitLen, 2))
	}
	if containerLen != 0 {
		f.DecMapUint8Uint8L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]uint64)
		if *vp == nil {
			*vp = make(map[uint8]uint64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapUint8Uint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Uint64L(rv2i(rv).(map[uint8]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]uint64, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapUint8Uint64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]int)
		if *vp == nil {
			*vp = make(map[uint8]int, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapUint8IntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8IntL(rv2i(rv).(map[uint8]int), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8IntX(vp *map[uint8]int, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]int, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapUint8IntL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]int given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]int32)
		if *vp == nil {
			*vp = make(map[uint8]int32, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			ft.DecMapUint8Int32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Int32L(rv2i(rv).(map[uint8]int32), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8Int32X(vp *map[uint8]int32, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]int32, decInferLen(containerLen, d.h.MaxInitLen, 5))
	}
	if containerLen != 0 {
		f.DecMapUint8Int32L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]float64)
		if *vp == nil {
			*vp = make(map[uint8]float64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapUint8Float64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Float64L(rv2i(rv).(map[uint8]float64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8Float64X(vp *map[uint8]float64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]float64, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapUint8Float64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]bool)
		if *vp == nil {
			*vp = make(map[uint8]bool, decInferLen(containerLen, d.h.MaxInitLen, 2))
		}
		if containerLen != 0 {
			ft.DecMapUint8BoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8BoolL(rv2i(rv).(map[uint8]bool), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint8BoolX(vp *map[uint8]bool, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint8]bool, decInferLen(containerLen, d.h.MaxInitLen, 2))
	}
	if containerLen != 0 {
		f.DecMapUint8BoolL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]bool given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]interface{})
		if *vp == nil {
			*vp = make(map[uint64]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapUint64IntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64IntfL(rv2i(rv).(map[uint64]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapUint64IntfL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk uint64
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]string)
		if *vp == nil {
			*vp = make(map[uint64]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapUint64StringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64StringL(rv2i(rv).(map[uint64]string), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64StringX(vp *map[uint64]string, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapUint64StringL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]string given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes())
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64][]byte)
		if *vp == nil {
			*vp = make(map[uint64][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			ft.DecMapUint64BytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64BytesL(rv2i(rv).(map[uint64][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
	}
	if containerLen != 0 {
		f.DecMapUint64BytesL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk uint64
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.decodeBytesInto(mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]uint8)
		if *vp == nil {
			*vp = make(map[uint64]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapUint64Uint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Uint8L(rv2i(rv).(map[uint64]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapUint64Uint8L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]uint64)
		if *vp == nil {
			*vp = make(map[uint64]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			ft.DecMapUint64Uint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Uint64L(rv2i(rv).(map[uint64]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
	}
	if containerLen != 0 {
		f.DecMapUint64Uint64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]int)
		if *vp == nil {
			*vp = make(map[uint64]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			ft.DecMapUint64IntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64IntL(rv2i(rv).(map[uint64]int), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64IntX(vp *map[uint64]int, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
	}
	if containerLen != 0 {
		f.DecMapUint64IntL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]int given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]int32)
		if *vp == nil {
			*vp = make(map[uint64]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			ft.DecMapUint64Int32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Int32L(rv2i(rv).(map[uint64]int32), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64Int32X(vp *map[uint64]int32, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
	}
	if containerLen != 0 {
		f.DecMapUint64Int32L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]float64)
		if *vp == nil {
			*vp = make(map[uint64]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			ft.DecMapUint64Float64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Float64L(rv2i(rv).(map[uint64]float64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64Float64X(vp *map[uint64]float64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
	}
	if containerLen != 0 {
		f.DecMapUint64Float64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]bool)
		if *vp == nil {
			*vp = make(map[uint64]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapUint64BoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64BoolL(rv2i(rv).(map[uint64]bool), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapUint64BoolX(vp *map[uint64]bool, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[uint64]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapUint64BoolL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]bool given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]interface{})
		if *vp == nil {
			*vp = make(map[int]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapIntIntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntIntfL(rv2i(rv).(map[int]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntIntfX(vp *map[int]interface{}, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapIntIntfL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk int
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]string)
		if *vp == nil {
			*vp = make(map[int]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			ft.DecMapIntStringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntStringL(rv2i(rv).(map[int]string), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntStringX(vp *map[int]string, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
	}
	if containerLen != 0 {
		f.DecMapIntStringL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntStringL(v map[int]string, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]string given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes())
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int][]byte)
		if *vp == nil {
			*vp = make(map[int][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			ft.DecMapIntBytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntBytesL(rv2i(rv).(map[int][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntBytesX(vp *map[int][]byte, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
	}
	if containerLen != 0 {
		f.DecMapIntBytesL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk int
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.decodeBytesInto(mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]uint8)
		if *vp == nil {
			*vp = make(map[int]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapIntUint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntUint8L(rv2i(rv).(map[int]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntUint8X(vp *map[int]uint8, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapIntUint8L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]uint64)
		if *vp == nil {
			*vp = make(map[int]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			ft.DecMapIntUint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntUint64L(rv2i(rv).(map[int]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntUint64X(vp *map[int]uint64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
	}
	if containerLen != 0 {
		f.DecMapIntUint64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]int)
		if *vp == nil {
			*vp = make(map[int]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			ft.DecMapIntIntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntIntL(rv2i(rv).(map[int]int), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntIntX(vp *map[int]int, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
	}
	if containerLen != 0 {
		f.DecMapIntIntL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntIntL(v map[int]int, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]int given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]int32)
		if *vp == nil {
			*vp = make(map[int]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			ft.DecMapIntInt32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntInt32L(rv2i(rv).(map[int]int32), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntInt32X(vp *map[int]int32, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
	}
	if containerLen != 0 {
		f.DecMapIntInt32L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]float64)
		if *vp == nil {
			*vp = make(map[int]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			ft.DecMapIntFloat64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntFloat64L(rv2i(rv).(map[int]float64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntFloat64X(vp *map[int]float64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
	}
	if containerLen != 0 {
		f.DecMapIntFloat64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]bool)
		if *vp == nil {
			*vp = make(map[int]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			ft.DecMapIntBoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntBoolL(rv2i(rv).(map[int]bool), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapIntBoolX(vp *map[int]bool, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
	}
	if containerLen != 0 {
		f.DecMapIntBoolL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]bool given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]interface{})
		if *vp == nil {
			*vp = make(map[int32]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			ft.DecMapInt32IntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32IntfL(rv2i(rv).(map[int32]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 20))
	}
	if containerLen != 0 {
		f.DecMapInt32IntfL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk int32
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]string)
		if *vp == nil {
			*vp = make(map[int32]string, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			ft.DecMapInt32StringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32StringL(rv2i(rv).(map[int32]string), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32StringX(vp *map[int32]string, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]string, decInferLen(containerLen, d.h.MaxInitLen, 20))
	}
	if containerLen != 0 {
		f.DecMapInt32StringL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]string given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes())
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32][]byte)
		if *vp == nil {
			*vp = make(map[int32][]byte, decInferLen(containerLen, d.h.MaxInitLen, 28))
		}
		if containerLen != 0 {
			ft.DecMapInt32BytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32BytesL(rv2i(rv).(map[int32][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32BytesX(vp *map[int32][]byte, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32][]byte, decInferLen(containerLen, d.h.MaxInitLen, 28))
	}
	if containerLen != 0 {
		f.DecMapInt32BytesL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk int32
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		mv = d.decodeBytesInto(mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]uint8)
		if *vp == nil {
			*vp = make(map[int32]uint8, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			ft.DecMapInt32Uint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Uint8L(rv2i(rv).(map[int32]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]uint8, decInferLen(containerLen, d.h.MaxInitLen, 5))
	}
	if containerLen != 0 {
		f.DecMapInt32Uint8L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]uint64)
		if *vp == nil {
			*vp = make(map[int32]uint64, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			ft.DecMapInt32Uint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Uint64L(rv2i(rv).(map[int32]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]uint64, decInferLen(containerLen, d.h.MaxInitLen, 12))
	}
	if containerLen != 0 {
		f.DecMapInt32Uint64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]int)
		if *vp == nil {
			*vp = make(map[int32]int, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			ft.DecMapInt32IntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32IntL(rv2i(rv).(map[int32]int), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32IntX(vp *map[int32]int, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]int, decInferLen(containerLen, d.h.MaxInitLen, 12))
	}
	if containerLen != 0 {
		f.DecMapInt32IntL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]int given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]int32)
		if *vp == nil {
			*vp = make(map[int32]int32, decInferLen(containerLen, d.h.MaxInitLen, 8))
		}
		if containerLen != 0 {
			ft.DecMapInt32Int32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Int32L(rv2i(rv).(map[int32]int32), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32Int32X(vp *map[int32]int32, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]int32, decInferLen(containerLen, d.h.MaxInitLen, 8))
	}
	if containerLen != 0 {
		f.DecMapInt32Int32L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]float64)
		if *vp == nil {
			*vp = make(map[int32]float64, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			ft.DecMapInt32Float64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Float64L(rv2i(rv).(map[int32]float64), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32Float64X(vp *map[int32]float64, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]float64, decInferLen(containerLen, d.h.MaxInitLen, 12))
	}
	if containerLen != 0 {
		f.DecMapInt32Float64L(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]bool)
		if *vp == nil {
			*vp = make(map[int32]bool, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			ft.DecMapInt32BoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32BoolL(rv2i(rv).(map[int32]bool), containerLen, d)
	}
	d.mapEnd()
}
func (f fastpathDT[T]) DecMapInt32BoolX(vp *map[int32]bool, d *decoder[T]) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
		return
	}
	if *vp == nil {
		*vp = make(map[int32]bool, decInferLen(containerLen, d.h.MaxInitLen, 5))
	}
	if containerLen != 0 {
		f.DecMapInt32BoolL(*vp, containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]bool given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
