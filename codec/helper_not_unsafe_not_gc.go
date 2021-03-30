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
//   runtime.{mapassign_fastXXX, mapaccess2_fastXXX} are not supported in gollvm,
//   failing with "error: undefined reference" error.
//   however, runtime.{mapassign, mapaccess2} are supported, so use that instead.
// - rvType:
//   reflect.toType is not supported in gccgo, gollvm.
// - reflect.{unsafe_New, unsafe_NewArray} are not supported in gollvm,
//   failing with "error: undefined reference" error.
//   however, runtime.{mallocgc, newarray} are supported, so use that instead.

func rvType(rv reflect.Value) reflect.Type {
	return rv.Type()
}
