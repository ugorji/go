// +build !go1.7 safe appengine

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"sync/atomic"
	"time"
)

const safeMode = true

// stringView returns a view of the []byte as a string.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
//
// Usage: Always maintain a reference to v while result of this call is in use,
//        and call keepAlive4BytesView(v) at point where done with view.
func stringView(v []byte) string {
	return string(v)
}

// bytesView returns a view of the string as a []byte.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
//
// Usage: Always maintain a reference to v while result of this call is in use,
//        and call keepAlive4BytesView(v) at point where done with view.
func bytesView(v string) []byte {
	return []byte(v)
}

// isNil says whether the value v is nil.
// This applies to references like map/ptr/unsafepointer/chan/func,
// and non-reference values like interface/slice.
func isNil(v interface{}) (rv reflect.Value, isnil bool) {
	rv = reflect.ValueOf(v)
	if isnilBitset.isset(byte(rv.Kind())) {
		isnil = rv.IsNil()
	}
	return
}

func rv2i(rv reflect.Value) interface{} {
	return rv.Interface()
}

func rvisnil(rv reflect.Value) bool {
	return rv.IsNil()
}

func rvssetlen(rv reflect.Value, length int) {
	rv.SetLen(length)
}

// func rvisnilref(rv reflect.Value) bool {
// 	return rv.IsNil()
// }

// func rvslen(rv reflect.Value) int {
// 	return rv.Len()
// }

// func rv2rtid(rv reflect.Value) uintptr {
// 	return reflect.ValueOf(rv.Type()).Pointer()
// }

func rt2id(rt reflect.Type) uintptr {
	return reflect.ValueOf(rt).Pointer()
}

func i2rtid(i interface{}) uintptr {
	return reflect.ValueOf(reflect.TypeOf(i)).Pointer()
}

// --------------------------

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, deref, checkStruct bool) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		if deref {
			if v.IsNil() {
				return true
			}
			return isEmptyValue(v.Elem(), tinfos, deref, checkStruct)
		}
		return v.IsNil()
	case reflect.Struct:
		return isEmptyStruct(v, tinfos, deref, checkStruct)
	}
	return false
}

// --------------------------
// type ptrToRvMap struct{}

// func (*ptrToRvMap) init() {}
// func (*ptrToRvMap) get(i interface{}) reflect.Value {
// 	return reflect.ValueOf(i).Elem()
// }

// --------------------------
type atomicClsErr struct {
	v atomic.Value
}

func (x *atomicClsErr) load() (e clsErr) {
	if i := x.v.Load(); i != nil {
		e = i.(clsErr)
	}
	return
}

func (x *atomicClsErr) store(p clsErr) {
	x.v.Store(p)
}

// --------------------------
type atomicTypeInfoSlice struct { // expected to be 2 words
	v atomic.Value
}

func (x *atomicTypeInfoSlice) load() (e []rtid2ti) {
	if i := x.v.Load(); i != nil {
		e = i.([]rtid2ti)
	}
	return
}

func (x *atomicTypeInfoSlice) store(p []rtid2ti) {
	x.v.Store(p)
}

// --------------------------
type atomicRtidFnSlice struct { // expected to be 2 words
	v atomic.Value
}

func (x *atomicRtidFnSlice) load() (e []codecRtidFn) {
	if i := x.v.Load(); i != nil {
		e = i.([]codecRtidFn)
	}
	return
}

func (x *atomicRtidFnSlice) store(p []codecRtidFn) {
	x.v.Store(p)
}

// --------------------------
func (n *decNaked) ru() reflect.Value {
	return reflect.ValueOf(&n.u).Elem()
}
func (n *decNaked) ri() reflect.Value {
	return reflect.ValueOf(&n.i).Elem()
}
func (n *decNaked) rf() reflect.Value {
	return reflect.ValueOf(&n.f).Elem()
}
func (n *decNaked) rl() reflect.Value {
	return reflect.ValueOf(&n.l).Elem()
}
func (n *decNaked) rs() reflect.Value {
	return reflect.ValueOf(&n.s).Elem()
}
func (n *decNaked) rt() reflect.Value {
	return reflect.ValueOf(&n.t).Elem()
}
func (n *decNaked) rb() reflect.Value {
	return reflect.ValueOf(&n.b).Elem()
}

// --------------------------
func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) {
	rv.SetBytes(d.rawBytes())
}

func (d *Decoder) kString(f *codecFnInfo, rv reflect.Value) {
	rv.SetString(string(d.d.DecodeStringAsBytes()))
}

func (d *Decoder) kBool(f *codecFnInfo, rv reflect.Value) {
	rv.SetBool(d.d.DecodeBool())
}

func (d *Decoder) kTime(f *codecFnInfo, rv reflect.Value) {
	rv.Set(reflect.ValueOf(d.d.DecodeTime()))
}

func (d *Decoder) kFloat32(f *codecFnInfo, rv reflect.Value) {
	rv.SetFloat(float64(d.decodeFloat32()))
}

func (d *Decoder) kFloat64(f *codecFnInfo, rv reflect.Value) {
	rv.SetFloat(d.d.DecodeFloat64())
}

func (d *Decoder) kInt(f *codecFnInfo, rv reflect.Value) {
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
}

func (d *Decoder) kInt8(f *codecFnInfo, rv reflect.Value) {
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), 8))
}

func (d *Decoder) kInt16(f *codecFnInfo, rv reflect.Value) {
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), 16))
}

func (d *Decoder) kInt32(f *codecFnInfo, rv reflect.Value) {
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), 32))
}

func (d *Decoder) kInt64(f *codecFnInfo, rv reflect.Value) {
	rv.SetInt(d.d.DecodeInt64())
}

func (d *Decoder) kUint(f *codecFnInfo, rv reflect.Value) {
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
}

func (d *Decoder) kUintptr(f *codecFnInfo, rv reflect.Value) {
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
}

func (d *Decoder) kUint8(f *codecFnInfo, rv reflect.Value) {
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), 8))
}

func (d *Decoder) kUint16(f *codecFnInfo, rv reflect.Value) {
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), 16))
}

func (d *Decoder) kUint32(f *codecFnInfo, rv reflect.Value) {
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), 32))
}

func (d *Decoder) kUint64(f *codecFnInfo, rv reflect.Value) {
	rv.SetUint(d.d.DecodeUint64())
}

// ----------------

func (e *Encoder) kBool(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rv.Bool())
}

func (e *Encoder) kTime(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rv2i(rv).(time.Time))
}

func (e *Encoder) kString(f *codecFnInfo, rv reflect.Value) {
	if e.h.StringToRaw {
		e.e.EncodeStringBytesRaw(bytesView(rv.String()))
	} else {
		e.e.EncodeStringEnc(cUTF8, rv.String())
	}
}

// func (e *Encoder) kStringToRaw(f *codecFnInfo, rv reflect.Value) {
// 	e.e.EncodeStringBytesRaw(bytesView(rv.String()))
// }

// func (e *Encoder) kStringEnc(f *codecFnInfo, rv reflect.Value) {
// 	e.e.EncodeStringEnc(cUTF8, rv.String())
// }

func (e *Encoder) kFloat64(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rv.Float())
}

func (e *Encoder) kFloat32(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(float32(rv.Float()))
}

func (e *Encoder) kInt(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(rv.Int())
}

func (e *Encoder) kInt8(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(rv.Int())
}

func (e *Encoder) kInt16(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(rv.Int())
}

func (e *Encoder) kInt32(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(rv.Int())
}

func (e *Encoder) kInt64(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(rv.Int())
}

func (e *Encoder) kUint(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(rv.Uint())
}

func (e *Encoder) kUint8(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(rv.Uint())
}

func (e *Encoder) kUint16(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(rv.Uint())
}

func (e *Encoder) kUint32(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(rv.Uint())
}

func (e *Encoder) kUint64(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(rv.Uint())
}

func (e *Encoder) kUintptr(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(rv.Uint())
}

// ------------ map range and map indexing ----------

func mapGet(m, k, v reflect.Value) (vv reflect.Value) {
	return m.MapIndex(k)
}

func mapSet(m, k, v reflect.Value) {
	m.SetMapIndex(k, v)
}

func mapDelete(m, k reflect.Value) {
	m.SetMapIndex(k, reflect.Value{})
}

// return an addressable reflect value that can be used in mapRange and mapGet operations.
//
// all calls to mapGet or mapRange will call here to get an addressable reflect.Value.
func mapAddressableRV(t reflect.Type) (r reflect.Value) {
	return // reflect.New(t).Elem()
}
