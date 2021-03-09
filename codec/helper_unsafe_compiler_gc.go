// +build !safe
// +build !codec.safe
// +build !appengine
// +build go1.9
// +build gc

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"unsafe"
)

func rvType(rv reflect.Value) reflect.Type {
	return rvPtrToType(((*unsafeReflectValue)(unsafe.Pointer(&rv))).typ) // rv.Type()
}

// rcGrowSlice updates the slice to point to a new array with the cap incremented, and len set to the new cap value.
// It copies data from old slice to new slice.
// It returns set=true iff it updates it, else it just returns a new slice pointing to a newly made array.
func rvGrowSlice(rv reflect.Value, ti *typeInfo, xcap, incr int) (v reflect.Value, newcap int, set bool) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	ux := (*unsafeSlice)(urv.ptr)
	t := ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	// zz.Debugf("before growslice: ux:Data(%x) Len/Cap: %d/%d, xcap/incr/xcap+incr: %d/%d/%d; rv: kind: %v, type: %v",
	// ux.Data, ux.Len, ux.Cap, xcap, incr, xcap+incr, rv.Kind(), rv.Type())
	*ux = growslice(t, *ux, xcap+incr)
	// zz.Debug2f("after growslice:  ux:Data(%x) Len/Cap: %d/%d", ux.Data, ux.Len, ux.Cap)
	ux.Len = ux.Cap
	return rv, ux.Cap, true
}

//go:linkname unsafeZeroArr runtime.zeroVal
var unsafeZeroArr [1024]byte

//go:linkname rvPtrToType reflect.toType
//go:noescape
func rvPtrToType(typ unsafe.Pointer) reflect.Type

//go:linkname growslice runtime.growslice
//go:noescape
func growslice(typ unsafe.Pointer, old unsafeSlice, cap int) unsafeSlice
