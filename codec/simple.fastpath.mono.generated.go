//go:build !codec.notmono  && !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"slices"

	"sort"
)

type fastpathESimpleBytes struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderSimpleBytes, *encFnInfo, reflect.Value)
}
type fastpathDSimpleBytes struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderSimpleBytes, *decFnInfo, reflect.Value)
}
type fastpathEsSimpleBytes [56]fastpathESimpleBytes
type fastpathDsSimpleBytes [56]fastpathDSimpleBytes
type fastpathETSimpleBytes struct{}
type fastpathDTSimpleBytes struct{}

func (helperEncDriverSimpleBytes) fastpathEList() *fastpathEsSimpleBytes {
	var i uint = 0
	var s fastpathEsSimpleBytes
	fn := func(v interface{}, fe func(*encoderSimpleBytes, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathESimpleBytes{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderSimpleBytes).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderSimpleBytes).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderSimpleBytes).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderSimpleBytes).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderSimpleBytes).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderSimpleBytes).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderSimpleBytes).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderSimpleBytes).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderSimpleBytes).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderSimpleBytes).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderSimpleBytes).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderSimpleBytes).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderSimpleBytes).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderSimpleBytes).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderSimpleBytes).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderSimpleBytes).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderSimpleBytes).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderSimpleBytes).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderSimpleBytes).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderSimpleBytes).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderSimpleBytes).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderSimpleBytes).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderSimpleBytes).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderSimpleBytes).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderSimpleBytes).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderSimpleBytes).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderSimpleBytes).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderSimpleBytes).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderSimpleBytes).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderSimpleBytes).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderSimpleBytes).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderSimpleBytes).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderSimpleBytes).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderSimpleBytes).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderSimpleBytes).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderSimpleBytes).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderSimpleBytes).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderSimpleBytes).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderSimpleBytes).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderSimpleBytes).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderSimpleBytes).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderSimpleBytes).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderSimpleBytes).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderSimpleBytes).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderSimpleBytes).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderSimpleBytes).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderSimpleBytes).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderSimpleBytes).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderSimpleBytes).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderSimpleBytes).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderSimpleBytes).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderSimpleBytes).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderSimpleBytes).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderSimpleBytes).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderSimpleBytes).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderSimpleBytes).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverSimpleBytes) fastpathDList() *fastpathDsSimpleBytes {
	var i uint = 0
	var s fastpathDsSimpleBytes
	fn := func(v interface{}, fd func(*decoderSimpleBytes, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDSimpleBytes{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderSimpleBytes).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderSimpleBytes).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderSimpleBytes).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderSimpleBytes).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderSimpleBytes).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderSimpleBytes).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderSimpleBytes).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderSimpleBytes).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderSimpleBytes).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderSimpleBytes).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderSimpleBytes).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderSimpleBytes).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderSimpleBytes).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderSimpleBytes).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderSimpleBytes).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderSimpleBytes).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderSimpleBytes).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderSimpleBytes).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderSimpleBytes).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderSimpleBytes).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderSimpleBytes).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderSimpleBytes).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderSimpleBytes).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderSimpleBytes).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderSimpleBytes).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderSimpleBytes).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderSimpleBytes).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderSimpleBytes).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderSimpleBytes).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderSimpleBytes).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderSimpleBytes).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderSimpleBytes).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderSimpleBytes).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderSimpleBytes).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderSimpleBytes).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderSimpleBytes).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderSimpleBytes).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderSimpleBytes).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderSimpleBytes).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderSimpleBytes).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderSimpleBytes).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderSimpleBytes).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderSimpleBytes).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderSimpleBytes).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderSimpleBytes).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderSimpleBytes).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderSimpleBytes).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderSimpleBytes).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderSimpleBytes).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderSimpleBytes).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderSimpleBytes).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderSimpleBytes).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderSimpleBytes).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderSimpleBytes).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderSimpleBytes).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderSimpleBytes).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverSimpleBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderSimpleBytes) bool {
	var ft fastpathETSimpleBytes
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
		_ = v
		return false
	}
	return true
}

func (e *encoderSimpleBytes) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceIntfV(v []interface{}, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceIntfV(v []interface{}, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.encode(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceStringV(v []string, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceStringV(v []string, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeString(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceBytesV(v [][]byte, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceBytesV(v [][]byte, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceFloat32V(v []float32, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceFloat32V(v []float32, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat32(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceFloat64V(v []float64, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceFloat64V(v []float64, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat64(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceUint8V(v []uint8, e *encoderSimpleBytes) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETSimpleBytes) EncAsMapSliceUint8V(v []uint8, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(uint64(v[j]))
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceUint64V(v []uint64, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceUint64V(v []uint64, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceIntV(v []int, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceIntV(v []int, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceInt32V(v []int32, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceInt32V(v []int32, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceInt64V(v []int64, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceInt64V(v []int64, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleBytes
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
func (fastpathETSimpleBytes) EncSliceBoolV(v []bool, e *encoderSimpleBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleBytes) EncAsMapSliceBoolV(v []bool, e *encoderSimpleBytes) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeBool(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleBytes) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETSimpleBytes) EncMapStringIntfV(v map[string]interface{}, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETSimpleBytes) EncMapStringStringV(v map[string]string, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETSimpleBytes) EncMapStringBytesV(v map[string][]byte, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETSimpleBytes) EncMapStringUint8V(v map[string]uint8, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETSimpleBytes) EncMapStringUint64V(v map[string]uint64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETSimpleBytes) EncMapStringIntV(v map[string]int, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETSimpleBytes) EncMapStringInt32V(v map[string]int32, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETSimpleBytes) EncMapStringFloat64V(v map[string]float64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETSimpleBytes) EncMapStringBoolV(v map[string]bool, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETSimpleBytes) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETSimpleBytes) EncMapUint8StringV(v map[uint8]string, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETSimpleBytes) EncMapUint8BytesV(v map[uint8][]byte, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETSimpleBytes) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETSimpleBytes) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETSimpleBytes) EncMapUint8IntV(v map[uint8]int, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETSimpleBytes) EncMapUint8Int32V(v map[uint8]int32, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETSimpleBytes) EncMapUint8Float64V(v map[uint8]float64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETSimpleBytes) EncMapUint8BoolV(v map[uint8]bool, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETSimpleBytes) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETSimpleBytes) EncMapUint64StringV(v map[uint64]string, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETSimpleBytes) EncMapUint64BytesV(v map[uint64][]byte, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETSimpleBytes) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETSimpleBytes) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETSimpleBytes) EncMapUint64IntV(v map[uint64]int, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETSimpleBytes) EncMapUint64Int32V(v map[uint64]int32, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETSimpleBytes) EncMapUint64Float64V(v map[uint64]float64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETSimpleBytes) EncMapUint64BoolV(v map[uint64]bool, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETSimpleBytes) EncMapIntIntfV(v map[int]interface{}, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETSimpleBytes) EncMapIntStringV(v map[int]string, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETSimpleBytes) EncMapIntBytesV(v map[int][]byte, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETSimpleBytes) EncMapIntUint8V(v map[int]uint8, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETSimpleBytes) EncMapIntUint64V(v map[int]uint64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETSimpleBytes) EncMapIntIntV(v map[int]int, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETSimpleBytes) EncMapIntInt32V(v map[int]int32, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETSimpleBytes) EncMapIntFloat64V(v map[int]float64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETSimpleBytes) EncMapIntBoolV(v map[int]bool, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETSimpleBytes) EncMapInt32IntfV(v map[int32]interface{}, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETSimpleBytes) EncMapInt32StringV(v map[int32]string, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETSimpleBytes) EncMapInt32BytesV(v map[int32][]byte, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETSimpleBytes) EncMapInt32Uint8V(v map[int32]uint8, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETSimpleBytes) EncMapInt32Uint64V(v map[int32]uint64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETSimpleBytes) EncMapInt32IntV(v map[int32]int, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETSimpleBytes) EncMapInt32Int32V(v map[int32]int32, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETSimpleBytes) EncMapInt32Float64V(v map[int32]float64, e *encoderSimpleBytes) {
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
func (e *encoderSimpleBytes) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleBytes{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETSimpleBytes) EncMapInt32BoolV(v map[int32]bool, e *encoderSimpleBytes) {
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

func (helperDecDriverSimpleBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderSimpleBytes) bool {
	var ft fastpathDTSimpleBytes
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
		_ = v
		return false
	}
	return true
}

func (d *decoderSimpleBytes) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceIntfX(vp *[]interface{}, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceIntfY(v []interface{}, d *decoderSimpleBytes) (v2 []interface{}, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceIntfN(v []interface{}, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceStringX(vp *[]string, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceStringY(v []string, d *decoderSimpleBytes) (v2 []string, changed bool) {
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
		if j == 0 && len(v) == 0 {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			v = make([]string, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, "")
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.stringZC(d.d.DecodeStringAsBytes(nil))
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
func (fastpathDTSimpleBytes) DecSliceStringN(v []string, d *decoderSimpleBytes) {
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
		v[uint(j)] = d.stringZC(d.d.DecodeStringAsBytes(nil))
	}
	slh.End()
}

func (d *decoderSimpleBytes) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceBytesX(vp *[][]byte, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceBytesY(v [][]byte, d *decoderSimpleBytes) (v2 [][]byte, changed bool) {
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
		if j == 0 && len(v) == 0 {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 24)
			v = make([][]byte, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, nil)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = bytesOk(d.d.DecodeBytes(zeroByteSlice))
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
func (fastpathDTSimpleBytes) DecSliceBytesN(v [][]byte, d *decoderSimpleBytes) {
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
		v[uint(j)] = bytesOk(d.d.DecodeBytes(zeroByteSlice))
	}
	slh.End()
}

func (d *decoderSimpleBytes) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceFloat32X(vp *[]float32, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceFloat32Y(v []float32, d *decoderSimpleBytes) (v2 []float32, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceFloat32N(v []float32, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceFloat64X(vp *[]float64, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceFloat64Y(v []float64, d *decoderSimpleBytes) (v2 []float64, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceFloat64N(v []float64, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceUint8X(vp *[]uint8, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceUint8Y(v []uint8, d *decoderSimpleBytes) (v2 []uint8, changed bool) {
	switch d.d.ContainerType() {
	case valueTypeNil, valueTypeMap:
		break
	default:
		v2 = d.decodeBytesInto(v[:len(v):len(v)])
		changed = !(len(v2) > 0 && len(v2) == len(v) && &v2[0] == &v[0])
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceUint8N(v []uint8, d *decoderSimpleBytes) {
	switch d.d.ContainerType() {
	case valueTypeNil, valueTypeMap:
		break
	default:
		v2 := d.decodeBytesInto(v[:len(v):len(v)])
		if !(len(v2) > 0 && len(v2) == len(v) && &v2[0] == &v[0]) {
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

func (d *decoderSimpleBytes) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceUint64X(vp *[]uint64, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceUint64Y(v []uint64, d *decoderSimpleBytes) (v2 []uint64, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceUint64N(v []uint64, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceIntX(vp *[]int, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceIntY(v []int, d *decoderSimpleBytes) (v2 []int, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceIntN(v []int, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceInt32X(vp *[]int32, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceInt32Y(v []int32, d *decoderSimpleBytes) (v2 []int32, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceInt32N(v []int32, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceInt64X(vp *[]int64, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceInt64Y(v []int64, d *decoderSimpleBytes) (v2 []int64, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceInt64N(v []int64, d *decoderSimpleBytes) {
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

func (d *decoderSimpleBytes) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecSliceBoolX(vp *[]bool, d *decoderSimpleBytes) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleBytes) DecSliceBoolY(v []bool, d *decoderSimpleBytes) (v2 []bool, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleBytes) DecSliceBoolN(v []bool, d *decoderSimpleBytes) {
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
func (d *decoderSimpleBytes) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringIntfX(vp *map[string]interface{}, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk string
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
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
func (d *decoderSimpleBytes) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringStringX(vp *map[string]string, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringStringL(v map[string]string, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]string given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringBytesX(vp *map[string][]byte, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk string
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
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
func (d *decoderSimpleBytes) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringUint8X(vp *map[string]uint8, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringUint64X(vp *map[string]uint64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringIntX(vp *map[string]int, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringIntL(v map[string]int, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]int given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringInt32X(vp *map[string]int32, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringFloat64X(vp *map[string]float64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapStringBoolX(vp *map[string]bool, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]bool given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8StringX(vp *map[uint8]string, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]string given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8IntX(vp *map[uint8]int, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8Int32X(vp *map[uint8]int32, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8Float64X(vp *map[uint8]float64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint8BoolX(vp *map[uint8]bool, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]bool given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64StringX(vp *map[uint64]string, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]string given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64IntX(vp *map[uint64]int, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64Int32X(vp *map[uint64]int32, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64Float64X(vp *map[uint64]float64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapUint64BoolX(vp *map[uint64]bool, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]bool given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntIntfX(vp *map[int]interface{}, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntStringX(vp *map[int]string, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntStringL(v map[int]string, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]string given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntBytesX(vp *map[int][]byte, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntUint8X(vp *map[int]uint8, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntUint64X(vp *map[int]uint64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntIntX(vp *map[int]int, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntIntL(v map[int]int, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntInt32X(vp *map[int]int32, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntFloat64X(vp *map[int]float64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapIntBoolX(vp *map[int]bool, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]bool given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32StringX(vp *map[int32]string, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]string given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleBytes) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32BytesX(vp *map[int32][]byte, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32IntX(vp *map[int32]int, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32Int32X(vp *map[int32]int32, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32Float64X(vp *map[int32]float64, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleBytes) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleBytes
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
func (f fastpathDTSimpleBytes) DecMapInt32BoolX(vp *map[int32]bool, d *decoderSimpleBytes) {
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
func (fastpathDTSimpleBytes) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderSimpleBytes) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]bool given stream length: ", int64(containerLen))
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

type fastpathESimpleIO struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderSimpleIO, *encFnInfo, reflect.Value)
}
type fastpathDSimpleIO struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderSimpleIO, *decFnInfo, reflect.Value)
}
type fastpathEsSimpleIO [56]fastpathESimpleIO
type fastpathDsSimpleIO [56]fastpathDSimpleIO
type fastpathETSimpleIO struct{}
type fastpathDTSimpleIO struct{}

func (helperEncDriverSimpleIO) fastpathEList() *fastpathEsSimpleIO {
	var i uint = 0
	var s fastpathEsSimpleIO
	fn := func(v interface{}, fe func(*encoderSimpleIO, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathESimpleIO{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderSimpleIO).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderSimpleIO).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderSimpleIO).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderSimpleIO).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderSimpleIO).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderSimpleIO).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderSimpleIO).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderSimpleIO).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderSimpleIO).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderSimpleIO).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderSimpleIO).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderSimpleIO).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderSimpleIO).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderSimpleIO).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderSimpleIO).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderSimpleIO).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderSimpleIO).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderSimpleIO).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderSimpleIO).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderSimpleIO).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderSimpleIO).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderSimpleIO).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderSimpleIO).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderSimpleIO).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderSimpleIO).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderSimpleIO).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderSimpleIO).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderSimpleIO).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderSimpleIO).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderSimpleIO).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderSimpleIO).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderSimpleIO).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderSimpleIO).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderSimpleIO).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderSimpleIO).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderSimpleIO).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderSimpleIO).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderSimpleIO).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderSimpleIO).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderSimpleIO).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderSimpleIO).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderSimpleIO).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderSimpleIO).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderSimpleIO).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderSimpleIO).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderSimpleIO).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderSimpleIO).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderSimpleIO).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderSimpleIO).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderSimpleIO).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderSimpleIO).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderSimpleIO).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderSimpleIO).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderSimpleIO).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderSimpleIO).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderSimpleIO).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverSimpleIO) fastpathDList() *fastpathDsSimpleIO {
	var i uint = 0
	var s fastpathDsSimpleIO
	fn := func(v interface{}, fd func(*decoderSimpleIO, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDSimpleIO{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderSimpleIO).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderSimpleIO).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderSimpleIO).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderSimpleIO).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderSimpleIO).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderSimpleIO).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderSimpleIO).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderSimpleIO).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderSimpleIO).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderSimpleIO).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderSimpleIO).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderSimpleIO).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderSimpleIO).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderSimpleIO).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderSimpleIO).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderSimpleIO).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderSimpleIO).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderSimpleIO).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderSimpleIO).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderSimpleIO).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderSimpleIO).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderSimpleIO).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderSimpleIO).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderSimpleIO).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderSimpleIO).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderSimpleIO).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderSimpleIO).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderSimpleIO).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderSimpleIO).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderSimpleIO).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderSimpleIO).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderSimpleIO).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderSimpleIO).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderSimpleIO).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderSimpleIO).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderSimpleIO).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderSimpleIO).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderSimpleIO).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderSimpleIO).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderSimpleIO).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderSimpleIO).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderSimpleIO).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderSimpleIO).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderSimpleIO).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderSimpleIO).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderSimpleIO).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderSimpleIO).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderSimpleIO).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderSimpleIO).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderSimpleIO).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderSimpleIO).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderSimpleIO).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderSimpleIO).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderSimpleIO).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderSimpleIO).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderSimpleIO).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverSimpleIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderSimpleIO) bool {
	var ft fastpathETSimpleIO
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
		_ = v
		return false
	}
	return true
}

func (e *encoderSimpleIO) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceIntfV(v []interface{}, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceIntfV(v []interface{}, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.encode(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceStringV(v []string, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceStringV(v []string, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeString(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceBytesV(v [][]byte, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceBytesV(v [][]byte, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceFloat32V(v []float32, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceFloat32V(v []float32, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat32(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceFloat64V(v []float64, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceFloat64V(v []float64, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat64(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceUint8V(v []uint8, e *encoderSimpleIO) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETSimpleIO) EncAsMapSliceUint8V(v []uint8, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(uint64(v[j]))
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceUint64V(v []uint64, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceUint64V(v []uint64, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceIntV(v []int, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceIntV(v []int, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceInt32V(v []int32, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceInt32V(v []int32, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceInt64V(v []int64, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceInt64V(v []int64, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETSimpleIO
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
func (fastpathETSimpleIO) EncSliceBoolV(v []bool, e *encoderSimpleIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathETSimpleIO) EncAsMapSliceBoolV(v []bool, e *encoderSimpleIO) {
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		e.e.EncodeBool(v[j])
	}
	e.mapEnd()
}
func (e *encoderSimpleIO) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETSimpleIO) EncMapStringIntfV(v map[string]interface{}, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETSimpleIO) EncMapStringStringV(v map[string]string, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETSimpleIO) EncMapStringBytesV(v map[string][]byte, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETSimpleIO) EncMapStringUint8V(v map[string]uint8, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETSimpleIO) EncMapStringUint64V(v map[string]uint64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETSimpleIO) EncMapStringIntV(v map[string]int, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETSimpleIO) EncMapStringInt32V(v map[string]int32, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETSimpleIO) EncMapStringFloat64V(v map[string]float64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETSimpleIO) EncMapStringBoolV(v map[string]bool, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETSimpleIO) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETSimpleIO) EncMapUint8StringV(v map[uint8]string, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETSimpleIO) EncMapUint8BytesV(v map[uint8][]byte, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETSimpleIO) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETSimpleIO) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETSimpleIO) EncMapUint8IntV(v map[uint8]int, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETSimpleIO) EncMapUint8Int32V(v map[uint8]int32, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETSimpleIO) EncMapUint8Float64V(v map[uint8]float64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETSimpleIO) EncMapUint8BoolV(v map[uint8]bool, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETSimpleIO) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETSimpleIO) EncMapUint64StringV(v map[uint64]string, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETSimpleIO) EncMapUint64BytesV(v map[uint64][]byte, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETSimpleIO) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETSimpleIO) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETSimpleIO) EncMapUint64IntV(v map[uint64]int, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETSimpleIO) EncMapUint64Int32V(v map[uint64]int32, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETSimpleIO) EncMapUint64Float64V(v map[uint64]float64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETSimpleIO) EncMapUint64BoolV(v map[uint64]bool, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETSimpleIO) EncMapIntIntfV(v map[int]interface{}, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETSimpleIO) EncMapIntStringV(v map[int]string, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETSimpleIO) EncMapIntBytesV(v map[int][]byte, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETSimpleIO) EncMapIntUint8V(v map[int]uint8, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETSimpleIO) EncMapIntUint64V(v map[int]uint64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETSimpleIO) EncMapIntIntV(v map[int]int, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETSimpleIO) EncMapIntInt32V(v map[int]int32, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETSimpleIO) EncMapIntFloat64V(v map[int]float64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETSimpleIO) EncMapIntBoolV(v map[int]bool, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETSimpleIO) EncMapInt32IntfV(v map[int32]interface{}, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETSimpleIO) EncMapInt32StringV(v map[int32]string, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETSimpleIO) EncMapInt32BytesV(v map[int32][]byte, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETSimpleIO) EncMapInt32Uint8V(v map[int32]uint8, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETSimpleIO) EncMapInt32Uint64V(v map[int32]uint64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETSimpleIO) EncMapInt32IntV(v map[int32]int, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETSimpleIO) EncMapInt32Int32V(v map[int32]int32, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETSimpleIO) EncMapInt32Float64V(v map[int32]float64, e *encoderSimpleIO) {
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
func (e *encoderSimpleIO) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETSimpleIO{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETSimpleIO) EncMapInt32BoolV(v map[int32]bool, e *encoderSimpleIO) {
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

func (helperDecDriverSimpleIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderSimpleIO) bool {
	var ft fastpathDTSimpleIO
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
		_ = v
		return false
	}
	return true
}

func (d *decoderSimpleIO) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceIntfX(vp *[]interface{}, d *decoderSimpleIO) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceIntfY(v []interface{}, d *decoderSimpleIO) (v2 []interface{}, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceIntfN(v []interface{}, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceStringX(vp *[]string, d *decoderSimpleIO) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceStringY(v []string, d *decoderSimpleIO) (v2 []string, changed bool) {
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
		if j == 0 && len(v) == 0 {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 16)
			v = make([]string, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, "")
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = d.stringZC(d.d.DecodeStringAsBytes(nil))
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
func (fastpathDTSimpleIO) DecSliceStringN(v []string, d *decoderSimpleIO) {
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
		v[uint(j)] = d.stringZC(d.d.DecodeStringAsBytes(nil))
	}
	slh.End()
}

func (d *decoderSimpleIO) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceBytesX(vp *[][]byte, d *decoderSimpleIO) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceBytesY(v [][]byte, d *decoderSimpleIO) (v2 [][]byte, changed bool) {
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
		if j == 0 && len(v) == 0 {
			xlen = decInferLen(containerLenS, d.h.MaxInitLen, 24)
			v = make([][]byte, uint(xlen))
			changed = true
		}
		if j >= len(v) {
			v = append(v, nil)
			changed = true
		}
		slh.ElemContainerState(j)
		v[uint(j)] = bytesOk(d.d.DecodeBytes(zeroByteSlice))
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
func (fastpathDTSimpleIO) DecSliceBytesN(v [][]byte, d *decoderSimpleIO) {
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
		v[uint(j)] = bytesOk(d.d.DecodeBytes(zeroByteSlice))
	}
	slh.End()
}

func (d *decoderSimpleIO) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceFloat32X(vp *[]float32, d *decoderSimpleIO) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceFloat32Y(v []float32, d *decoderSimpleIO) (v2 []float32, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceFloat32N(v []float32, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceFloat64X(vp *[]float64, d *decoderSimpleIO) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceFloat64Y(v []float64, d *decoderSimpleIO) (v2 []float64, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceFloat64N(v []float64, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceUint8X(vp *[]uint8, d *decoderSimpleIO) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceUint8Y(v []uint8, d *decoderSimpleIO) (v2 []uint8, changed bool) {
	switch d.d.ContainerType() {
	case valueTypeNil, valueTypeMap:
		break
	default:
		v2 = d.decodeBytesInto(v[:len(v):len(v)])
		changed = !(len(v2) > 0 && len(v2) == len(v) && &v2[0] == &v[0])
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceUint8N(v []uint8, d *decoderSimpleIO) {
	switch d.d.ContainerType() {
	case valueTypeNil, valueTypeMap:
		break
	default:
		v2 := d.decodeBytesInto(v[:len(v):len(v)])
		if !(len(v2) > 0 && len(v2) == len(v) && &v2[0] == &v[0]) {
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

func (d *decoderSimpleIO) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceUint64X(vp *[]uint64, d *decoderSimpleIO) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceUint64Y(v []uint64, d *decoderSimpleIO) (v2 []uint64, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceUint64N(v []uint64, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceIntX(vp *[]int, d *decoderSimpleIO) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceIntY(v []int, d *decoderSimpleIO) (v2 []int, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceIntN(v []int, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceInt32X(vp *[]int32, d *decoderSimpleIO) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceInt32Y(v []int32, d *decoderSimpleIO) (v2 []int32, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceInt32N(v []int32, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceInt64X(vp *[]int64, d *decoderSimpleIO) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceInt64Y(v []int64, d *decoderSimpleIO) (v2 []int64, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceInt64N(v []int64, d *decoderSimpleIO) {
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

func (d *decoderSimpleIO) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecSliceBoolX(vp *[]bool, d *decoderSimpleIO) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTSimpleIO) DecSliceBoolY(v []bool, d *decoderSimpleIO) (v2 []bool, changed bool) {
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
		if j == 0 && len(v) == 0 {
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
func (fastpathDTSimpleIO) DecSliceBoolN(v []bool, d *decoderSimpleIO) {
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
func (d *decoderSimpleIO) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringIntfX(vp *map[string]interface{}, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	var mk string
	var mv interface{}
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
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
func (d *decoderSimpleIO) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringStringX(vp *map[string]string, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringStringL(v map[string]string, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]string given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringBytesX(vp *map[string][]byte, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := !d.h.MapValueReset
	var mk string
	var mv []byte
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
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
func (d *decoderSimpleIO) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringUint8X(vp *map[string]uint8, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]uint8 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv uint8
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringUint64X(vp *map[string]uint64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]uint64 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv uint64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.d.DecodeUint64()
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringIntX(vp *map[string]int, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringIntL(v map[string]int, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]int given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv int
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringInt32X(vp *map[string]int32, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]int32 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv int32
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringFloat64X(vp *map[string]float64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]float64 given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv float64
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.d.DecodeFloat64()
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapStringBoolX(vp *map[string]bool, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]bool given stream length: ", int64(containerLen))
		return
	}
	var mk string
	var mv bool
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.stringZC(d.d.DecodeStringAsBytes(nil))
		d.mapElemValue()
		mv = d.d.DecodeBool()
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8StringX(vp *map[uint8]string, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]string given stream length: ", int64(containerLen))
		return
	}
	var mk uint8
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8IntX(vp *map[uint8]int, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8Int32X(vp *map[uint8]int32, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8Float64X(vp *map[uint8]float64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint8BoolX(vp *map[uint8]bool, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]bool given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64StringX(vp *map[uint64]string, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]string given stream length: ", int64(containerLen))
		return
	}
	var mk uint64
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = d.d.DecodeUint64()
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64IntX(vp *map[uint64]int, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64Int32X(vp *map[uint64]int32, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64Float64X(vp *map[uint64]float64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapUint64BoolX(vp *map[uint64]bool, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]bool given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntIntfX(vp *map[int]interface{}, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntStringX(vp *map[int]string, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntStringL(v map[int]string, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]string given stream length: ", int64(containerLen))
		return
	}
	var mk int
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntBytesX(vp *map[int][]byte, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntUint8X(vp *map[int]uint8, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntUint64X(vp *map[int]uint64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntIntX(vp *map[int]int, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntIntL(v map[int]int, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntInt32X(vp *map[int]int32, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntFloat64X(vp *map[int]float64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapIntBoolX(vp *map[int]bool, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]bool given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]interface{} given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32StringX(vp *map[int32]string, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]string given stream length: ", int64(containerLen))
		return
	}
	var mk int32
	var mv string
	hasLen := containerLen > 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey()
		mk = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		mv = d.stringZC(d.d.DecodeStringAsBytes(nil))
		v[mk] = mv
	}
}
func (d *decoderSimpleIO) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32BytesX(vp *map[int32][]byte, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32][]byte given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]uint8 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]uint64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32IntX(vp *map[int32]int, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]int given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32Int32X(vp *map[int32]int32, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]int32 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32Float64X(vp *map[int32]float64, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]float64 given stream length: ", int64(containerLen))
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
func (d *decoderSimpleIO) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTSimpleIO
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
func (f fastpathDTSimpleIO) DecMapInt32BoolX(vp *map[int32]bool, d *decoderSimpleIO) {
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
func (fastpathDTSimpleIO) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderSimpleIO) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]bool given stream length: ", int64(containerLen))
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
