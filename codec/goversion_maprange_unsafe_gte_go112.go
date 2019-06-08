// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build !safe
// +build !appengine
// +build go1.12

package codec

import (
	"reflect"
	"unsafe"
)

type unsafeReflectMapIter struct {
	m  unsafeReflectValue
	it unsafe.Pointer
}

type mapIter struct {
	t              *reflect.MapIter
	m, k, v        reflect.Value
	ktyp, vtyp     unsafe.Pointer
	kptr, vptr     unsafe.Pointer
	kisref, visref bool
	mapvalues      bool
	done           bool
	_              uint64 // padding (cache-aligned)
}

func (t *mapIter) Next() (r bool) {
	if t.done {
		return
	}
	r = t.t.Next()
	if r {
		it := (*unsafeReflectMapIter)(unsafe.Pointer(t.t))
		mapSet(t.kptr, t.ktyp, mapiterkey(it.it), t.kisref)
		if t.mapvalues {
			mapSet(t.vptr, t.vtyp, mapiterelem(it.it), t.visref)
		}
	} else {
		t.done = true
	}
	return
}

func (t *mapIter) Key() reflect.Value {
	return t.k
}

func (t *mapIter) Value() (r reflect.Value) {
	if t.mapvalues {
		return t.v
	}
	return
}

func mapSet(p, ptyp, p2 unsafe.Pointer, isref bool) {
	if isref {
		*(*unsafe.Pointer)(p) = *(*unsafe.Pointer)(p2) // p2
	} else {
		typedmemmove(ptyp, p, p2) // *(*unsafe.Pointer)(p2)) // p2)
	}
}

func mapRange(m, k, v reflect.Value, mapvalues bool) *mapIter {
	if m.IsNil() {
		return &mapIter{done: true}
	}
	t := &mapIter{
		m: m, k: k, v: v,
		t: m.MapRange(), mapvalues: mapvalues,
	}

	var urv *unsafeReflectValue
	urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	t.ktyp = urv.typ
	t.kptr = urv.ptr
	t.kisref = refBitset.isset(byte(k.Kind()))
	if mapvalues {
		urv = (*unsafeReflectValue)(unsafe.Pointer(&v))
		t.vtyp = urv.typ
		t.vptr = urv.ptr
		t.visref = refBitset.isset(byte(v.Kind()))
	}

	return t
}

func mapIndex(m, k, v reflect.Value) (vv reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr unsafe.Pointer
	if urv.flag&unsafeFlagIndir == 0 {
		kptr = unsafe.Pointer(&urv.ptr)
	} else {
		kptr = urv.ptr
	}

	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))

	vvptr := mapaccess(urv.typ, rv2ptr(urv), kptr)
	if vvptr == nil {
		return
	}
	// vvptr = *(*unsafe.Pointer)(vvptr)

	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))

	mapSet(urv.ptr, urv.typ, vvptr, refBitset.isset(byte(v.Kind())))
	return v
}

// return an addressable reflect value that can be used in mapRange and mapIndex operations.
//
// all calls to mapIndex or mapRange will call here to get an addressable reflect.Value.
func mapAddressableRV(t reflect.Type) (r reflect.Value) {
	return reflect.New(t).Elem()
}

//go:linkname mapiterkey reflect.mapiterkey
//go:noescape
func mapiterkey(it unsafe.Pointer) (key unsafe.Pointer)

//go:linkname mapiterelem reflect.mapiterelem
//go:noescape
func mapiterelem(it unsafe.Pointer) (elem unsafe.Pointer)

//go:linkname mapaccess reflect.mapaccess
//go:noescape
func mapaccess(rtype unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) (val unsafe.Pointer)

//go:linkname typedmemmove reflect.typedmemmove
//go:noescape
func typedmemmove(typ unsafe.Pointer, dst, src unsafe.Pointer)
