// +build !safe
// +build !appengine
// +build go1.7

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"
)

// This file has unsafe variants of some helper methods.
// NOTE: See helper_not_unsafe.go for the usage information.

// For reflect.Value code, we decided to do the following:
//    - if we know the kind, we can elide conditional checks for
//      - SetXXX (Int, Uint, String, Bool, etc)
//      - SetLen
//
// We can also optimize
//      - IsNil
//
// We cannot do the same for Cap, Len if we still have to do conditional.

// var zeroRTv [4]uintptr

const safeMode = false

// keep in sync with GO_ROOT/src/reflect/value.go
const (
	unsafeFlagIndir    = 1 << 7
	unsafeFlagKindMask = (1 << 5) - 1 // 5 bits for 27 kinds (up to 31)
	// unsafeTypeKindDirectIface = 1 << 5
)

type unsafeString struct {
	Data unsafe.Pointer
	Len  int
}

type unsafeSlice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type unsafeIntf struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

type unsafeReflectValue struct {
	typ  unsafe.Pointer
	ptr  unsafe.Pointer
	flag uintptr
}

func stringView(v []byte) string {
	if len(v) == 0 {
		return ""
	}
	bx := (*unsafeSlice)(unsafe.Pointer(&v))
	return *(*string)(unsafe.Pointer(&unsafeString{bx.Data, bx.Len}))
}

func bytesView(v string) []byte {
	if len(v) == 0 {
		return zeroByteSlice
	}
	sx := (*unsafeString)(unsafe.Pointer(&v))
	return *(*[]byte)(unsafe.Pointer(&unsafeSlice{sx.Data, sx.Len, sx.Len}))
}

// // isNilRef says whether the interface is a nil reference or not.
// //
// // A reference here is a pointer-sized reference i.e. map, ptr, chan, func, unsafepointer.
// // It is optional to extend this to also check if slices or interfaces are nil also.
// //
// // NOTE: There is no global way of checking if an interface is nil.
// // For true references (map, ptr, func, chan), you can just look
// // at the word of the interface.
// // However, for slices, you have to dereference
// // the word, and get a pointer to the 3-word interface value.
// func isNilRef(v interface{}) (rv reflect.Value, isnil bool) {
// 	isnil = ((*unsafeIntf)(unsafe.Pointer(&v))).word == nil
// 	return
// }

func isNil(v interface{}) (rv reflect.Value, isnil bool) {
	var ui *unsafeIntf = (*unsafeIntf)(unsafe.Pointer(&v))
	if ui.word == nil {
		isnil = true
		return
	}
	rv = reflect.ValueOf(v) // reflect.value is cheap and inline'able
	tk := rv.Kind()
	isnil = (tk == reflect.Interface || tk == reflect.Slice) && *(*unsafe.Pointer)(ui.word) == nil
	return
	// fmt.Printf(">>>> definitely nil: isnil: %v, TYPE: \t%T, word: %v, *word: %v, type: %v, nil: %v\n",
	// 	v == nil, v, word, *((*unsafe.Pointer)(word)), ui.typ, nil)
}

func rv2ptr(urv *unsafeReflectValue) (ptr unsafe.Pointer) {
	// true references (map, func, chan, ptr - NOT slice) may be double-referenced? as flagIndir
	// rv := *((*reflect.Value)(unsafe.Pointer(urv)))
	if refBitset.isset(byte(urv.flag&unsafeFlagKindMask)) && urv.flag&unsafeFlagIndir != 0 {
		ptr = *(*unsafe.Pointer)(urv.ptr)
	} else {
		ptr = urv.ptr
	}
	return
}

func rv2i(rv reflect.Value) interface{} {
	// TODO: consider a more generally-known optimization for reflect.Value ==> Interface
	//
	// Currently, we use this fragile method that taps into implememtation details from
	// the source go stdlib reflect/value.go, and trims the implementation.

	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*interface{})(unsafe.Pointer(&unsafeIntf{typ: urv.typ, word: rv2ptr(urv)}))
}

func rvisnil(rv reflect.Value) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if urv.flag&unsafeFlagIndir != 0 {
		return *(*unsafe.Pointer)(urv.ptr) == nil
	}
	return urv.ptr == nil
}

func rvssetlen(rv reflect.Value, length int) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	(*unsafeString)(urv.ptr).Len = length
}

// func rvisnilref(rv reflect.Value) bool {
// 	return (*unsafeReflectValue)(unsafe.Pointer(&rv)).ptr == nil
// }

// func rvslen(rv reflect.Value) int {
// 	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	return (*unsafeString)(urv.ptr).Len
// }

// func rv2rtid(rv reflect.Value) uintptr {
// 	return uintptr((*unsafeReflectValue)(unsafe.Pointer(&rv)).typ)
// }

func rt2id(rt reflect.Type) uintptr {
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&rt))).word)
}

func i2rtid(i interface{}) uintptr {
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&i))).typ)
}

// --------------------------

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, deref, checkStruct bool) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if urv.flag == 0 {
		return true
	}
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String:
		return (*unsafeString)(urv.ptr).Len == 0
	case reflect.Slice:
		return (*unsafeSlice)(urv.ptr).Len == 0
	case reflect.Bool:
		return !*(*bool)(urv.ptr)
	case reflect.Int:
		return *(*int)(urv.ptr) == 0
	case reflect.Int8:
		return *(*int8)(urv.ptr) == 0
	case reflect.Int16:
		return *(*int16)(urv.ptr) == 0
	case reflect.Int32:
		return *(*int32)(urv.ptr) == 0
	case reflect.Int64:
		return *(*int64)(urv.ptr) == 0
	case reflect.Uint:
		return *(*uint)(urv.ptr) == 0
	case reflect.Uint8:
		return *(*uint8)(urv.ptr) == 0
	case reflect.Uint16:
		return *(*uint16)(urv.ptr) == 0
	case reflect.Uint32:
		return *(*uint32)(urv.ptr) == 0
	case reflect.Uint64:
		return *(*uint64)(urv.ptr) == 0
	case reflect.Uintptr:
		return *(*uintptr)(urv.ptr) == 0
	case reflect.Float32:
		return *(*float32)(urv.ptr) == 0
	case reflect.Float64:
		return *(*float64)(urv.ptr) == 0
	case reflect.Interface:
		isnil := urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
		if deref {
			if isnil {
				return true
			}
			return isEmptyValue(v.Elem(), tinfos, deref, checkStruct)
		}
		return isnil
	case reflect.Ptr:
		// isnil := urv.ptr == nil (not sufficient, as a pointer value encodes the type)
		isnil := urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
		if deref {
			if isnil {
				return true
			}
			return isEmptyValue(v.Elem(), tinfos, deref, checkStruct)
		}
		return isnil
	case reflect.Struct:
		return isEmptyStruct(v, tinfos, deref, checkStruct)
	case reflect.Map, reflect.Array, reflect.Chan:
		return v.Len() == 0
	}
	return false
}

// --------------------------

// atomicXXX is expected to be 2 words (for symmetry with atomic.Value)
//
// Note that we do not atomically load/store length and data pointer separately,
// as this could lead to some races. Instead, we atomically load/store cappedSlice.
//
// Note: with atomic.(Load|Store)Pointer, we MUST work with an unsafe.Pointer directly.

// ----------------------
type atomicTypeInfoSlice struct {
	v unsafe.Pointer // *[]rtid2ti
	_ uint64         // padding (atomicXXX expected to be 2 words)
}

func (x *atomicTypeInfoSlice) load() (s []rtid2ti) {
	x2 := atomic.LoadPointer(&x.v)
	if x2 != nil {
		s = *(*[]rtid2ti)(x2)
	}
	return
}

func (x *atomicTypeInfoSlice) store(p []rtid2ti) {
	atomic.StorePointer(&x.v, unsafe.Pointer(&p))
}

// --------------------------
type atomicRtidFnSlice struct {
	v unsafe.Pointer // *[]codecRtidFn
	_ uint64         // padding (atomicXXX expected to be 2 words) (make 1 word so JsonHandle fits)
}

func (x *atomicRtidFnSlice) load() (s []codecRtidFn) {
	x2 := atomic.LoadPointer(&x.v)
	if x2 != nil {
		s = *(*[]codecRtidFn)(x2)
	}
	return
}

func (x *atomicRtidFnSlice) store(p []codecRtidFn) {
	atomic.StorePointer(&x.v, unsafe.Pointer(&p))
}

// --------------------------
type atomicClsErr struct {
	v unsafe.Pointer // *clsErr
	_ uint64         // padding (atomicXXX expected to be 2 words)
}

func (x *atomicClsErr) load() (e clsErr) {
	x2 := (*clsErr)(atomic.LoadPointer(&x.v))
	if x2 != nil {
		e = *x2
	}
	return
}

func (x *atomicClsErr) store(p clsErr) {
	atomic.StorePointer(&x.v, unsafe.Pointer(&p))
}

// --------------------------

// to create a reflect.Value for each member field of decNaked,
// we first create a global decNaked, and create reflect.Value
// for them all.
// This way, we have the flags and type in the reflect.Value.
// Then, when a reflect.Value is called, we just copy it,
// update the ptr to the decNaked's, and return it.

type unsafeDecNakedWrapper struct {
	decNaked
	ru, ri, rf, rl, rs, rb, rt reflect.Value // mapping to the primitives above
}

func (n *unsafeDecNakedWrapper) init() {
	n.ru = reflect.ValueOf(&n.u).Elem()
	n.ri = reflect.ValueOf(&n.i).Elem()
	n.rf = reflect.ValueOf(&n.f).Elem()
	n.rl = reflect.ValueOf(&n.l).Elem()
	n.rs = reflect.ValueOf(&n.s).Elem()
	n.rt = reflect.ValueOf(&n.t).Elem()
	n.rb = reflect.ValueOf(&n.b).Elem()
	// n.rr[] = reflect.ValueOf(&n.)
}

var defUnsafeDecNakedWrapper unsafeDecNakedWrapper

func init() {
	defUnsafeDecNakedWrapper.init()
}

func (n *decNaked) ru() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.ru
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.u)
	return
}
func (n *decNaked) ri() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.ri
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.i)
	return
}
func (n *decNaked) rf() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rf
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.f)
	return
}
func (n *decNaked) rl() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rl
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.l)
	return
}
func (n *decNaked) rs() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rs
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.s)
	return
}
func (n *decNaked) rt() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rt
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.t)
	return
}
func (n *decNaked) rb() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rb
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.b)
	return
}

// --------------------------
func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*[]byte)(urv.ptr) = d.rawBytes()
}

func (d *Decoder) kString(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*string)(urv.ptr) = string(d.d.DecodeStringAsBytes())
}

func (d *Decoder) kBool(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*bool)(urv.ptr) = d.d.DecodeBool()
}

func (d *Decoder) kTime(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*time.Time)(urv.ptr) = d.d.DecodeTime()
}

func (d *Decoder) kFloat32(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*float32)(urv.ptr) = d.decodeFloat32()
}

func (d *Decoder) kFloat64(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*float64)(urv.ptr) = d.d.DecodeFloat64()
}

func (d *Decoder) kInt(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int)(urv.ptr) = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
}

func (d *Decoder) kInt8(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int8)(urv.ptr) = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
}

func (d *Decoder) kInt16(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int16)(urv.ptr) = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
}

func (d *Decoder) kInt32(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int32)(urv.ptr) = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
}

func (d *Decoder) kInt64(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int64)(urv.ptr) = d.d.DecodeInt64()
}

func (d *Decoder) kUint(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint)(urv.ptr) = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
}

func (d *Decoder) kUintptr(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uintptr)(urv.ptr) = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
}

func (d *Decoder) kUint8(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint8)(urv.ptr) = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
}

func (d *Decoder) kUint16(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint16)(urv.ptr) = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
}

func (d *Decoder) kUint32(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint32)(urv.ptr) = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
}

func (d *Decoder) kUint64(f *codecFnInfo, rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint64)(urv.ptr) = d.d.DecodeUint64()
}

// ------------

func (e *Encoder) kBool(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeBool(*(*bool)(v.ptr))
}

func (e *Encoder) kTime(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeTime(*(*time.Time)(v.ptr))
}

func (e *Encoder) kString(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if e.h.StringToRaw {
		e.e.EncodeStringBytesRaw(bytesView(*(*string)(v.ptr)))
	} else {
		e.e.EncodeStringEnc(cUTF8, *(*string)(v.ptr))
	}
}

// func (e *Encoder) kStringToRaw(f *codecFnInfo, rv reflect.Value) {
// 	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	s := *(*string)(v.ptr)
// 	e.e.EncodeStringBytesRaw(bytesView(s))
// }

// func (e *Encoder) kStringEnc(f *codecFnInfo, rv reflect.Value) {
// 	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	s := *(*string)(v.ptr)
// 	e.e.EncodeStringEnc(cUTF8, s)
// }

func (e *Encoder) kFloat64(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeFloat64(*(*float64)(v.ptr))
}

func (e *Encoder) kFloat32(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeFloat32(*(*float32)(v.ptr))
}

func (e *Encoder) kInt(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int)(v.ptr)))
}

func (e *Encoder) kInt8(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int8)(v.ptr)))
}

func (e *Encoder) kInt16(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int16)(v.ptr)))
}

func (e *Encoder) kInt32(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int32)(v.ptr)))
}

func (e *Encoder) kInt64(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int64)(v.ptr)))
}

func (e *Encoder) kUint(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint)(v.ptr)))
}

func (e *Encoder) kUint8(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint8)(v.ptr)))
}

func (e *Encoder) kUint16(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint16)(v.ptr)))
}

func (e *Encoder) kUint32(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint32)(v.ptr)))
}

func (e *Encoder) kUint64(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint64)(v.ptr)))
}

func (e *Encoder) kUintptr(f *codecFnInfo, rv reflect.Value) {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uintptr)(v.ptr)))
}

// ------------ map range and map indexing ----------

// regular calls to map via reflection: MapKeys, MapIndex, MapRange/MapIter etc
// will always allocate for each map key or value.
//
// It is more performant to provide a value that the map entry is set into,
// and that elides the allocation.

// unsafeMapHashIter
//
// go 1.4+ has runtime/hashmap.go or runtime/map.go which has a
// hIter struct with the first 2 values being key and value
// of the current iteration.
//
// This *hIter is passed to mapiterinit, mapiternext, mapiterkey, mapiterelem.
// We bypass the reflect wrapper functions and just use the *hIter directly.
//
// Though *hIter has many fields, we only care about the first 2.
type unsafeMapHashIter struct {
	key, value unsafe.Pointer
	// other fields are ignored
}

// type unsafeReflectMapIter struct {
// 	m  unsafeReflectValue
// 	it unsafe.Pointer
// }

type unsafeMapIter struct {
	it               *unsafeMapHashIter
	k, v             reflect.Value
	mtyp, ktyp, vtyp unsafe.Pointer
	mptr, kptr, vptr unsafe.Pointer
	kisref, visref   bool
	mapvalues        bool
	done             bool

	// _ [2]uint64 // padding (cache-aligned)
}

func (t *unsafeMapIter) Next() (r bool) {
	if t.done {
		return
	}
	if t.it == nil {
		t.it = (*unsafeMapHashIter)(mapiterinit(t.mtyp, t.mptr))
	} else {
		mapiternext((unsafe.Pointer)(t.it))
	}
	t.done = t.it.key == nil
	if t.done {
		return
	}
	unsafeSet(t.kptr, t.ktyp, t.it.key, t.kisref)
	if t.mapvalues {
		unsafeSet(t.vptr, t.vtyp, t.it.value, t.visref)
	}
	return true
}

func (t *unsafeMapIter) Key() reflect.Value {
	return t.k
}

func (t *unsafeMapIter) Value() (r reflect.Value) {
	if t.mapvalues {
		return t.v
	}
	return
}

func unsafeSet(p, ptyp, p2 unsafe.Pointer, isref bool) {
	if isref {
		*(*unsafe.Pointer)(p) = *(*unsafe.Pointer)(p2) // p2
	} else {
		typedmemmove(ptyp, p, p2) // *(*unsafe.Pointer)(p2)) // p2)
	}
}

func unsafeMapKVPtr(urv *unsafeReflectValue) unsafe.Pointer {
	if urv.flag&unsafeFlagIndir == 0 {
		return unsafe.Pointer(&urv.ptr)
	}
	return urv.ptr
}

func mapRange(m, k, v reflect.Value, mapvalues bool) *unsafeMapIter {
	if m.IsNil() {
		return &unsafeMapIter{done: true}
	}
	t := &unsafeMapIter{
		k: k, v: v,
		mapvalues: mapvalues,
	}

	var urv *unsafeReflectValue

	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	t.mtyp = urv.typ
	t.mptr = rv2ptr(urv)

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

func mapGet(m, k, v reflect.Value) (vv reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)

	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))

	vvptr := mapaccess(urv.typ, rv2ptr(urv), kptr)
	if vvptr == nil {
		return
	}
	// vvptr = *(*unsafe.Pointer)(vvptr)

	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))

	unsafeSet(urv.ptr, urv.typ, vvptr, refBitset.isset(byte(v.Kind())))
	return v
}

func mapSet(m, k, v reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))
	var vptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mapassign(urv.typ, rv2ptr(urv), kptr, vptr)
}

func mapDelete(m, k reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mapdelete(urv.typ, rv2ptr(urv), kptr)
}

// return an addressable reflect value that can be used in mapRange and mapGet operations.
//
// all calls to mapGet or mapRange will call here to get an addressable reflect.Value.
func mapAddressableRV(t reflect.Type) (r reflect.Value) {
	return reflect.New(t).Elem()
}

//go:linkname mapiterinit reflect.mapiterinit
//go:noescape
func mapiterinit(typ unsafe.Pointer, it unsafe.Pointer) (key unsafe.Pointer)

//go:linkname mapiternext reflect.mapiternext
//go:noescape
func mapiternext(it unsafe.Pointer) (key unsafe.Pointer)

//go:linkname mapaccess reflect.mapaccess
//go:noescape
func mapaccess(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) (val unsafe.Pointer)

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(typ unsafe.Pointer, m unsafe.Pointer, key, val unsafe.Pointer)

//go:linkname mapdelete reflect.mapdelete
//go:noescape
func mapdelete(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer)

//go:linkname typedmemmove reflect.typedmemmove
//go:noescape
func typedmemmove(typ unsafe.Pointer, dst, src unsafe.Pointer)
