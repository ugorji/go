// +build !safe
// +build !codec.safe
// +build !appengine
// +build go1.9
// +build !gc

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "reflect"

var unsafeZeroArr [1024]byte

func rvType(rv reflect.Value) reflect.Type {
	return rv.Type()
}

// rvGrowSlice in helper_unsafe_compiler_gc.go does not work with gccgo,
// failing with a "growslice: cap out of range" error.
// Seems there's something going wrong in growslice.
// Until we find out why, use standard reflection.
func rvGrowSlice(rv reflect.Value, ti *typeInfo, xcap, incr int) (v reflect.Value, newcap int, set bool) {
	// zz.Debugf("rv: kind: %v, type: %v", rv.Kind(), rv.Type())
	newcap = int(growCap(uint(xcap), uint(ti.elemsize), uint(incr)))
	v = reflect.MakeSlice(ti.rt, newcap, newcap)
	if rv.Len() > 0 {
		reflect.Copy(v, rv)
	}
	return
}
