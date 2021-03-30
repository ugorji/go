// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !go1.9 || safe || codec.safe || appengine || !gc
// +build !go1.9 safe codec.safe appengine !gc

package codec

import "reflect"

// This files contains safe versions of the code where the unsafe versions are not supported
// in either gccgo or gollvm.
//
// - rvGrowSlice:
//   runtime.growslice does not work with gccgo,
//   failing with "growslice: cap out of range" error.
// - mapSet/mapGet:
//   runtime.{mapassign, mapaccess} are not supported in gollvm,
//   failing with "error: undefined reference" error.
// - rvType:
//   reflect.toType is not supported in gccgo, gollvm.
// - reflect.{unsafe_New, unsafe_NewArray} are not supported in gollvm,
//   failing with "error: undefined reference" error.
//   however, runtime.{mallocgc, newarray} are supported, to use that instead.
//
// When runtime.{mapassign, mapaccess} are supported in gollvm, then we will
// use mapSet and mapGet from helper_unsafe_compiler_gc.go, as that gives significant
// performance improvement and reduces allocation.

func rvType(rv reflect.Value) reflect.Type {
	return rv.Type()
}

func rvGrowSlice(rv reflect.Value, ti *typeInfo, xcap, incr int) (v reflect.Value, newcap int, set bool) {
	newcap = int(growCap(uint(xcap), uint(ti.elemsize), uint(incr)))
	v = reflect.MakeSlice(ti.rt, newcap, newcap)
	if rv.Len() > 0 {
		reflect.Copy(v, rv)
	}
	return
}

func mapSet(m, k, v reflect.Value, keyFastKind mapKeyFastKind, valIsIndirect, valIsRef bool) {
	m.SetMapIndex(k, v)
}

func mapGet(m, k, v reflect.Value, keyFastKind mapKeyFastKind, valIsIndirect, valIsRef bool) (vv reflect.Value) {
	return m.MapIndex(k)
}
