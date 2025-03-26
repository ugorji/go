// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !codec.generics

package codec

import (
	"encoding"
	"io"
	"math"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

func (e *encoderCborBytes) setContainerState(cs containerState) {
	if cs != 0 {
		e.c = cs
	}
}

func (e *encoderCborBytes) rawExt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *encoderCborBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderCborBytes) selferMarshal(f *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderCborBytes) binaryMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderCborBytes) textMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderCborBytes) jsonMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderCborBytes) raw(f *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderCborBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderCborBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderCborBytes) kBool(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderCborBytes) kTime(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderCborBytes) kString(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderCborBytes) kFloat32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderCborBytes) kFloat64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderCborBytes) kComplex64(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderCborBytes) kComplex128(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderCborBytes) kInt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderCborBytes) kInt8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderCborBytes) kInt16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderCborBytes) kInt32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderCborBytes) kInt64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderCborBytes) kUint(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderCborBytes) kUint8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderCborBytes) kUint16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderCborBytes) kUint32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderCborBytes) kUint64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderCborBytes) kUintptr(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderCborBytes) kErr(f *encFnInfo, rv reflect.Value) {
	e.errorf("unsupported encoding kind %s, for %#v", rv.Kind(), any(rv))
}

func (e *encoderCborBytes) kSeqFn(rtelem reflect.Type) (fn *encFnCborBytes) {
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}

	if rtelem.Kind() != reflect.Interface {
		fn = e.fn(rtelem)
	}
	return
}

func (e *encoderCborBytes) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	if l == 0 {
		e.mapStart(0)
	} else {
		e.haltOnMbsOddLen(l)
		e.mapStart(l >> 1)
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			if j&1 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
	e.mapEnd()
}

func (e *encoderCborBytes) kSliceW(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	e.arrayStart(l)
	if l > 0 {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
	e.arrayEnd()
}

func (e *encoderCborBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	if l == 0 {
		e.mapStart(0)
	} else {
		e.haltOnMbsOddLen(l)
		e.mapStart(l >> 1)
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			if j&1 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encodeValue(rv.Index(j), fn)
		}
	}
	e.mapEnd()
}

func (e *encoderCborBytes) kArrayW(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	e.arrayStart(l)
	if l > 0 {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rv.Index(j), fn)
		}
	}
	e.arrayEnd()
}

func (e *encoderCborBytes) kChan(f *encFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		e.errorStr("send-only channel cannot be encoded")
	}
	if !f.ti.mbs && uint8TypId == rt2id(f.ti.elem) {
		e.kSliceBytesChan(rv)
		return
	}
	rtslice := reflect.SliceOf(f.ti.elem)
	rv = chanToSlice(rv, rtslice, e.h.ChanRecvTimeout)
	ti := e.h.getTypeInfo(rt2id(rtslice), rtslice)
	if f.ti.mbs {
		e.kSliceWMbs(rv, ti)
	} else {
		e.kSliceW(rv, ti)
	}
}

func (e *encoderCborBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoderCborBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoderCborBytes) kSliceBytesChan(rv reflect.Value) {

	bs0 := e.blist.peek(32, true)
	bs := bs0

	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok {
		ch = irv.(chan byte)
	}

L1:
	switch timeout := e.h.ChanRecvTimeout; {
	case timeout == 0:
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			}
		}
	case timeout > 0:
		tt := time.NewTimer(timeout)
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:

				break L1
			}
		}
	default:
		for b := range ch {
			bs = append(bs, b)
		}
	}

	e.e.EncodeStringBytesRaw(bs)
	e.blist.put(bs)
	if !byteSliceSameData(bs0, bs) {
		e.blist.put(bs0)
	}
}

func (e *encoderCborBytes) kStructSfi(f *encFnInfo) []*structFieldInfo {
	if e.h.Canonical {
		return f.ti.sfi.sorted()
	}
	return f.ti.sfi.source()
}

func (e *encoderCborBytes) kStructNoOmitempty(f *encFnInfo, rv reflect.Value) {
	var tisfi []*structFieldInfo
	if f.ti.toArray || e.h.StructToArray {
		tisfi = f.ti.sfi.source()
		e.arrayStart(len(tisfi))
		for _, si := range tisfi {
			e.arrayElem()
			e.encodeValue(si.path.field(rv), nil)
		}
		e.arrayEnd()
	} else {
		tisfi = e.kStructSfi(f)
		e.mapStart(len(tisfi))
		keytyp := f.ti.keyType
		for _, si := range tisfi {
			e.mapElemKey()
			e.kStructFieldKey(keytyp, si.path.encNameAsciiAlphaNum, si.encName)
			e.mapElemValue()
			e.encodeValue(si.path.field(rv), nil)
		}
		e.mapEnd()
	}
}

func (e *encoderCborBytes) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
	if keyType == valueTypeString {
		if e.js && encNameAsciiAlphaNum {
			e.e.writeStringAsisDblQuoted(encName)
		} else {
			e.e.EncodeString(encName)
		}
	} else if keyType == valueTypeInt {
		e.e.EncodeInt(must.Int(strconv.ParseInt(encName, 10, 64)))
	} else if keyType == valueTypeUint {
		e.e.EncodeUint(must.Uint(strconv.ParseUint(encName, 10, 64)))
	} else if keyType == valueTypeFloat {
		e.e.EncodeFloat64(must.Float(strconv.ParseFloat(encName, 64)))
	} else {
		halt.errorStr2("invalid struct key type: ", keyType.String())
	}

}

func (e *encoderCborBytes) kStruct(f *encFnInfo, rv reflect.Value) {
	var newlen int
	ti := f.ti
	toMap := !(ti.toArray || e.h.StructToArray)
	var mf map[string]interface{}
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	} else if ti.flagMissingFielderPtr {
		rv2 := e.addrRV(rv, ti.rt, ti.ptr)
		mf = rv2i(rv2).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	}
	tisfi := ti.sfi.source()
	newlen += len(tisfi)

	var fkvs = e.slist.get(newlen)[:newlen]

	recur := e.h.RecursiveEmptyCheck

	var kv sfiRv
	var j int
	if toMap {
		newlen = 0
		for _, si := range e.kStructSfi(f) {
			kv.r = si.path.field(rv)
			if si.path.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
				continue
			}
			kv.v = si
			fkvs[newlen] = kv
			newlen++
		}

		var mf2s []stringIntf
		if len(mf) > 0 {
			mf2s = make([]stringIntf, 0, len(mf))
			for k, v := range mf {
				if k == "" {
					continue
				}
				if ti.infoFieldOmitempty && isEmptyValue(reflect.ValueOf(v), e.h.TypeInfos, recur) {
					continue
				}
				mf2s = append(mf2s, stringIntf{k, v})
			}
		}

		e.mapStart(newlen + len(mf2s))

		if len(mf2s) > 0 && e.h.Canonical {
			mf2w := make([]encStructFieldObj, newlen+len(mf2s))
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				mf2w[j] = encStructFieldObj{kv.v.encName, kv.r, nil, kv.v.path.encNameAsciiAlphaNum, true}
			}
			for _, v := range mf2s {
				mf2w[j] = encStructFieldObj{v.v, reflect.Value{}, v.i, false, false}
				j++
			}
			sort.Sort((encStructFieldObjSlice)(mf2w))
			for _, v := range mf2w {
				e.mapElemKey()
				e.kStructFieldKey(ti.keyType, v.ascii, v.key)
				e.mapElemValue()
				if v.isRv {
					e.encodeValue(v.rv, nil)
				} else {
					e.encode(v.intf)
				}
			}
		} else {
			keytyp := ti.keyType
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				e.mapElemKey()
				e.kStructFieldKey(keytyp, kv.v.path.encNameAsciiAlphaNum, kv.v.encName)
				e.mapElemValue()
				e.encodeValue(kv.r, nil)
			}
			for _, v := range mf2s {
				e.mapElemKey()
				e.kStructFieldKey(keytyp, false, v.v)
				e.mapElemValue()
				e.encode(v.i)
			}
		}

		e.mapEnd()
	} else {
		newlen = len(tisfi)
		for i, si := range tisfi {
			kv.r = si.path.field(rv)

			if si.path.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
				switch kv.r.Kind() {
				case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
					kv.r = reflect.Value{}
				}
			}
			fkvs[i] = kv
		}

		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ {
			e.arrayElem()
			e.encodeValue(fkvs[j].r, nil)
		}
		e.arrayEnd()
	}

	e.slist.put(fkvs)
}

func (e *encoderCborBytes) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	var keyFn, valFn *encFnCborBytes

	ktypeKind := reflect.Kind(f.ti.keykind)
	vtypeKind := reflect.Kind(f.ti.elemkind)

	rtval := f.ti.elem
	rtvalkind := vtypeKind
	for rtvalkind == reflect.Ptr {
		rtval = rtval.Elem()
		rtvalkind = rtval.Kind()
	}
	if rtvalkind != reflect.Interface {
		valFn = e.fn(rtval)
	}

	var rvv = mapAddrLoopvarRV(f.ti.elem, vtypeKind)

	rtkey := f.ti.key
	var keyTypeIsString = stringTypId == rt2id(rtkey)
	if keyTypeIsString {
		keyFn = e.fn(rtkey)
	} else {
		for rtkey.Kind() == reflect.Ptr {
			rtkey = rtkey.Elem()
		}
		if rtkey.Kind() != reflect.Interface {
			keyFn = e.fn(rtkey)
		}
	}

	if e.h.Canonical {
		e.kMapCanonical(f.ti, rv, rvv, keyFn, valFn)
		e.mapEnd()
		return
	}

	var rvk = mapAddrLoopvarRV(f.ti.key, ktypeKind)

	var it mapIter
	mapRange(&it, rv, rvk, rvv, true)

	for it.Next() {
		e.mapElemKey()
		if keyTypeIsString {
			e.e.EncodeString(it.Key().String())
		} else {
			e.encodeValue(it.Key(), keyFn)
		}
		e.mapElemValue()
		e.encodeValue(it.Value(), valFn)
	}
	it.Done()

	e.mapEnd()
}

func (e *encoderCborBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnCborBytes) {

	rtkey := ti.key
	rtkeydecl := rtkey.PkgPath() == "" && rtkey.Name() != ""

	mks := rv.MapKeys()
	rtkeyKind := rtkey.Kind()
	kfast := mapKeyFastKindFor(rtkeyKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	switch rtkeyKind {
	case reflect.Bool:

		if len(mks) == 2 && mks[0].Bool() {
			mks[0], mks[1] = mks[1], mks[0]
		}
		for i := range mks {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeBool(mks[i].Bool())
			} else {
				e.encodeValueNonNil(mks[i], keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mks[i], rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.String:
		mksv := make([]orderedRv[string], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.String()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeString(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		mksv := make([]orderedRv[uint64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Uint()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeUint(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		mksv := make([]orderedRv[int64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Int()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeInt(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Float32:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeFloat32(float32(mksv[i].v))
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Float64:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeFloat64(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	default:
		if rtkey == timeTyp {
			mksv := make([]timeRv, len(mks))
			for i, k := range mks {
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			}
			slices.SortFunc(mksv, cmpTimeRv)
			for i := range mksv {
				e.mapElemKey()
				e.e.EncodeTime(mksv[i].v)
				e.mapElemValue()
				e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
			}
			break
		}

		bs0 := e.blist.get(len(mks) * 16)
		mksv := bs0
		mksbv := make([]bytesRv, len(mks))

		func() {
			se := e.h.sideEncPool.Get().(encoderI)
			defer e.h.sideEncPool.Put(se)
			se.ResetBytes(&mksv)
			for i, k := range mks {
				v := &mksbv[i]
				l := len(mksv)
				se.setContainerState(containerMapKey)
				se.encode(rv2i(baseRVRV(k)))
				se.atEndOfEncode()
				se.writerEnd()
				v.r = k
				v.v = mksv[l:]
			}
		}()

		slices.SortFunc(mksbv, cmpBytesRv)
		for j := range mksbv {
			e.mapElemKey()
			e.e.writeBytesAsis(mksbv[j].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksbv[j].r, rvv, kfast, visindirect, visref), valFn)
		}
		e.blist.put(mksv)
		if !byteSliceSameData(bs0, mksv) {
			e.blist.put(bs0)
		}
	}
}

type encoderCborBytes struct {
	panicHdl
	perType encPerType

	dh helperEncDriverCborBytes

	fp *fastpathEsCborBytes

	h *BasicHandle

	e cborEncDriverBytes

	encoderShared

	hh Handle

	ci []interface{}

	slist sfiRvFreelist
}

func (e *encoderCborBytes) HandleName() string {
	return e.hh.Name()
}

func (e *encoderCborBytes) Release() {
}

func (e *encoderCborBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderShared, e).(*fastpathEsCborBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderCborBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderCborBytes) MustEncode(v interface{}) {
	halt.onerror(e.err)
	if e.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	e.calls++
	e.encode(v)
	e.calls--
	if e.calls == 0 {
		e.e.atEndOfEncode()
		e.e.writerEnd()
	}
}

func (e *encoderCborBytes) Encode(v interface{}) (err error) {

	if !debugging {
		defer func() {

			if x := recover(); x != nil {
				panicValToErr(e, x, &e.err)
				err = e.err
			}
		}()
	}

	e.MustEncode(v)
	return
}

func (e *encoderCborBytes) encode(iv interface{}) {

	if iv == nil {
		e.e.EncodeNil()
		return
	}

	rv, ok := isNil(iv)
	if ok {
		e.e.EncodeNil()
		return
	}

	switch v := iv.(type) {

	case Raw:
		e.rawBytes(v)
	case reflect.Value:
		e.encodeValue(v, nil)

	case string:
		e.e.EncodeString(v)
	case bool:
		e.e.EncodeBool(v)
	case int:
		e.e.EncodeInt(int64(v))
	case int8:
		e.e.EncodeInt(int64(v))
	case int16:
		e.e.EncodeInt(int64(v))
	case int32:
		e.e.EncodeInt(int64(v))
	case int64:
		e.e.EncodeInt(v)
	case uint:
		e.e.EncodeUint(uint64(v))
	case uint8:
		e.e.EncodeUint(uint64(v))
	case uint16:
		e.e.EncodeUint(uint64(v))
	case uint32:
		e.e.EncodeUint(uint64(v))
	case uint64:
		e.e.EncodeUint(v)
	case uintptr:
		e.e.EncodeUint(uint64(v))
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)
	case complex64:
		e.encodeComplex64(v)
	case complex128:
		e.encodeComplex128(v)
	case time.Time:
		e.e.EncodeTime(v)
	case []byte:
		e.e.EncodeStringBytesRaw(v)
	case *Raw:
		e.rawBytes(*v)
	case *string:
		e.e.EncodeString(*v)
	case *bool:
		e.e.EncodeBool(*v)
	case *int:
		e.e.EncodeInt(int64(*v))
	case *int8:
		e.e.EncodeInt(int64(*v))
	case *int16:
		e.e.EncodeInt(int64(*v))
	case *int32:
		e.e.EncodeInt(int64(*v))
	case *int64:
		e.e.EncodeInt(*v)
	case *uint:
		e.e.EncodeUint(uint64(*v))
	case *uint8:
		e.e.EncodeUint(uint64(*v))
	case *uint16:
		e.e.EncodeUint(uint64(*v))
	case *uint32:
		e.e.EncodeUint(uint64(*v))
	case *uint64:
		e.e.EncodeUint(*v)
	case *uintptr:
		e.e.EncodeUint(uint64(*v))
	case *float32:
		e.e.EncodeFloat32(*v)
	case *float64:
		e.e.EncodeFloat64(*v)
	case *complex64:
		e.encodeComplex64(*v)
	case *complex128:
		e.encodeComplex128(*v)
	case *time.Time:
		e.e.EncodeTime(*v)
	case *[]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			e.e.EncodeStringBytesRaw(*v)
		}
	default:

		if skipFastpathTypeSwitchInDirectCall || !e.dh.fastpathEncodeTypeSwitch(iv, e) {
			e.encodeValue(rv, nil)
		}
	}
}

func (e *encoderCborBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderCborBytes) encodeValue(rv reflect.Value, fn *encFnCborBytes) {

	var sptr interface{}
	var rvp reflect.Value
	var rvpValid bool
TOP:
	switch rv.Kind() {
	case reflect.Ptr:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
		rvpValid = true
		rvp = rv
		rv = rv.Elem()
		goto TOP
	case reflect.Interface:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
		rvpValid = false
		rvp = reflect.Value{}
		rv = rv.Elem()
		goto TOP
	case reflect.Struct:
		if rvpValid && e.h.CheckCircularRef {
			sptr = rv2i(rvp)
			for _, vv := range e.ci {
				if eq4i(sptr, vv) {
					e.errorf("circular reference found: %p, %T", sptr, sptr)
				}
			}
			e.ci = append(e.ci, sptr)
		}
	case reflect.Slice, reflect.Map, reflect.Chan:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		return
	}

	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if !fn.i.addrE {

	} else if rvpValid {
		rv = rvp
	} else {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)

	if sptr != nil {
		e.ci = e.ci[:len(e.ci)-1]
	}
}

func (e *encoderCborBytes) encodeValueNonNil(rv reflect.Value, fn *encFnCborBytes) {
	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if fn.i.addrE {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderCborBytes) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
	if rv.CanAddr() {
		return rvAddr(rv, ptrType)
	}
	if e.h.NoAddressableReadonly {
		rva = reflect.New(typ)
		rvSetDirect(rva.Elem(), rv)
		return
	}
	return rvAddr(e.perType.AddressableRO(rv), ptrType)
}

func (e *encoderCborBytes) marshalUtf8(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderCborBytes) marshalAsis(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderCborBytes) marshalRaw(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoderCborBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		e.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderCborBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}

func (e *encoderCborBytes) fn(t reflect.Type) *encFnCborBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoderCborBytes) fnNoExt(t reflect.Type) *encFnCborBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

func (e *encoderCborBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderCborBytes) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoderCborBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderCborBytes) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoderCborBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderCborBytes) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoderCborBytes) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

func (e *encoderCborBytes) haltOnMbsOddLen(length int) {
	if length&1 != 0 {
		e.errorInt("mapBySlice requires even slice length, but got ", int64(length))
	}
}

func (e *encoderCborBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderCborBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderCborBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderCborBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderCborBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverCborBytes) newEncoderBytes(out *[]byte, h Handle) *encoderCborBytes {
	var c1 encoderCborBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverCborBytes) newEncoderIO(out io.Writer, h Handle) *encoderCborBytes {
	var c1 encoderCborBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverCborBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsCborBytes) (f *fastpathECborBytes, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

type encFnCborBytes struct {
	i  encFnInfo
	fe func(*encoderCborBytes, *encFnInfo, reflect.Value)
}
type encRtidFnCborBytes struct {
	rtid uintptr
	fn   *encFnCborBytes
}

func (helperEncDriverCborBytes) encFindRtidFn(s []encRtidFnCborBytes, rtid uintptr) (i uint, fn *encFnCborBytes) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperEncDriverCborBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnCborBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnCborBytes](v))
	}
	return
}

func (dh helperEncDriverCborBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsCborBytes, checkExt bool) (fn *encFnCborBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverCborBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnCborBytes
	sp = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverCborBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnCborBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnCborBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnCborBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnCborBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverCborBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborBytes) {
	fn = new(encFnCborBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderCborBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderCborBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderCborBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderCborBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fe = (*encoderCborBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderCborBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderCborBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderCborBytes).textMarshal
		fi.addrE = ti.flagTextMarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {

			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fe = fp[idx].encfn
				}
			} else {

				xfe, xrt := dh.encFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf := xfe.encfn
					fn.fe = func(e *encoderCborBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderCborBytes).kBool
			case reflect.String:

				fn.fe = (*encoderCborBytes).kString
			case reflect.Int:
				fn.fe = (*encoderCborBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderCborBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderCborBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderCborBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderCborBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderCborBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderCborBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderCborBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderCborBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderCborBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderCborBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderCborBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderCborBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderCborBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderCborBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderCborBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderCborBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderCborBytes).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*encoderCborBytes).kStruct
				} else {
					fn.fe = (*encoderCborBytes).kStructNoOmitempty
				}
			case reflect.Map:
				fn.fe = (*encoderCborBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderCborBytes).kErr
			default:

				fn.fe = (*encoderCborBytes).kErr
			}
		}
	}
	return
}

type helperEncDriverCborBytes struct{}

func (d *decoderCborBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *decoderCborBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderCborBytes) selferUnmarshal(f *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderCborBytes) binaryUnmarshal(f *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	d.onerror(fnerr)
}

func (d *decoderCborBytes) textUnmarshal(f *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	d.onerror(fnerr)
}

func (d *decoderCborBytes) jsonUnmarshal(f *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderCborBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

	var bs0 = zeroByteSlice
	if !d.bytes {
		bs0 = d.blist.get(256)
	}
	bs := d.d.nextValueBytes(bs0)
	fnerr := tm.UnmarshalJSON(bs)
	if !d.bytes {
		d.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			d.blist.put(bs0)
		}
	}
	d.onerror(fnerr)
}

func (d *decoderCborBytes) kErr(f *decFnInfo, rv reflect.Value) {
	d.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderCborBytes) raw(f *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderCborBytes) kString(f *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes()))
}

func (d *decoderCborBytes) kBool(f *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderCborBytes) kTime(f *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderCborBytes) kFloat32(f *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderCborBytes) kFloat64(f *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderCborBytes) kComplex64(f *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderCborBytes) kComplex128(f *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderCborBytes) kInt(f *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderCborBytes) kInt8(f *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderCborBytes) kInt16(f *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderCborBytes) kInt32(f *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderCborBytes) kInt64(f *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderCborBytes) kUint(f *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborBytes) kUintptr(f *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborBytes) kUint8(f *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderCborBytes) kUint16(f *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderCborBytes) kUint32(f *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderCborBytes) kUint64(f *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderCborBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

	n := d.naked()
	d.d.DecodeNaked()

	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
	}
	switch n.v {
	case valueTypeMap:
		mtid := d.mtid
		if mtid == 0 {
			if d.jsms {
				mtid = mapStrIntfTypId
			} else {
				mtid = mapIntfIntfTypId
			}
		}
		if mtid == mapStrIntfTypId {
			var v2 map[string]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if mtid == mapIntfIntfTypId {
			var v2 map[interface{}]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.mtr {
			rvn = reflect.New(d.h.MapType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.MapType, reflect.Map)
			d.decodeValue(rvn, nil)
		}
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId {
			var v2 []interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.str {
			rvn = reflect.New(d.h.SliceType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.SliceType, reflect.Slice)
			d.decodeValue(rvn, nil)
		}
		if d.h.PreferArrayOverSlice {
			rvn = rvGetArray4Slice(rvn)
		}
	case valueTypeExt:
		tag, bytes := n.u, n.l
		bfn := d.h.getExtForTag(tag)
		var re = RawExt{Tag: tag}
		if bytes == nil {

			if bfn == nil {
				d.decode(&re.Value)
				rvn = rv4iptr(&re).Elem()
			} else {
				if bfn.ext == SelfExt {
					rvn = rvZeroAddrK(bfn.rt, bfn.rt.Kind())
					d.decodeValue(rvn, d.fnNoExt(bfn.rt))
				} else {
					rvn = reflect.New(bfn.rt)
					d.interfaceExtConvertAndDecode(rv2i(rvn), bfn.ext)
					rvn = rvn.Elem()
				}
			}
		} else {

			if bfn == nil {
				re.setData(bytes, false)
				rvn = rv4iptr(&re).Elem()
			} else {
				rvn = reflect.New(bfn.rt)
				if bfn.ext == SelfExt {
					sideDecode(d.h, rv2i(rvn), bytes, bfn.rt, true)
				} else {
					bfn.ext.ReadExt(rv2i(rvn), bytes)
				}
				rvn = rvn.Elem()
			}
		}

		if d.h.PreferPointerForStructOrArray && rvn.CanAddr() {
			if rk := rvn.Kind(); rk == reflect.Array || rk == reflect.Struct {
				rvn = rvn.Addr()
			}
		}
	case valueTypeNil:

	case valueTypeInt:
		rvn = n.ri()
	case valueTypeUint:
		rvn = n.ru()
	case valueTypeFloat:
		rvn = n.rf()
	case valueTypeBool:
		rvn = n.rb()
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs()
	case valueTypeBytes:
		rvn = n.rl()
	case valueTypeTime:
		rvn = n.rt()
	default:
		halt.errorStr2("kInterfaceNaked: unexpected valueType: ", n.v.String())
	}
	return
}

func (d *decoderCborBytes) kInterface(f *decFnInfo, rv reflect.Value) {

	isnilrv := rvIsNil(rv)

	var rvn reflect.Value

	if d.h.InterfaceReset {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			} else if !isnilrv {
				decSetNonNilRV2Zero4Intf(rv)
			}
			return
		}
	} else if isnilrv {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			}
			return
		}
	} else {

		rvn = rv.Elem()
	}

	canDecode, _ := isDecodeable(rvn)

	if !canDecode {
		rvn2 := d.oneShotAddrRV(rvn.Type(), rvn.Kind())
		rvSetDirect(rvn2, rvn)
		rvn = rvn2
	}

	d.decodeValue(rvn, nil)
	rvSetIntf(rv, rvn)
}

func (d *decoderCborBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *decoderCborBytes) kStruct(f *decFnInfo, rv reflect.Value) {
	ctyp := d.d.ContainerType()
	ti := f.ti
	var mf MissingFielder
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder)
	} else if ti.flagMissingFielderPtr {
		mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder)
	}
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var name2 []byte
		if mf != nil {
			var namearr2 [16]byte
			name2 = namearr2[:0]
		}
		var rvkencname []byte
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey()
			switch ti.keyType {
			case valueTypeString:
				rvkencname = d.d.DecodeStringAsBytes()
			case valueTypeInt:
				rvkencname = strconv.AppendInt(d.b[:0], d.d.DecodeInt64(), 10)
			case valueTypeUint:
				rvkencname = strconv.AppendUint(d.b[:0], d.d.DecodeUint64(), 10)
			case valueTypeFloat:
				rvkencname = strconv.AppendFloat(d.b[:0], d.d.DecodeFloat64(), 'f', -1, 64)
			default:
				halt.errorStr2("invalid struct key type: ", ti.keyType.String())
			}

			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else if mf != nil {

				name2 = append(name2[:0], rvkencname...)
				var f interface{}
				d.decode(&f)
				if !mf.CodecMissingField(name2, f) && d.h.ErrorIfNoField {
					d.errorStr2("no matching struct field when decoding stream map with key: ", stringView(name2))
				}
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}

		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem()
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}

		d.arrayEnd()
	} else {
		d.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderCborBytes) kSlice(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		rvbs := rvGetBytes(rv)
		if !rvCanset {

			rvbs = rvbs[:len(rvbs):len(rvbs)]
		}
		bs2 := d.decodeBytesInto(rvbs)

		if !(len(bs2) > 0 && len(bs2) == len(rvbs) && &bs2[0] == &rvbs[0]) {
			if rvCanset {
				rvSetBytes(rv, bs2)
			} else if len(rvbs) > 0 && len(bs2) > 0 {
				copy(rvbs, bs2)
			}
		}
		return
	}

	slh, containerLenS := d.decSliceHelperStart()

	if containerLenS == 0 {
		if rvCanset {
			if rvIsNil(rv) {
				rvSetDirect(rv, rvSliceZeroCap(ti.rt))
			} else {
				rvSetSliceLen(rv, 0)
			}
		}
		slh.End()
		return
	}

	rtelem0Mut := !scalarBitset.isset(ti.elemkind)
	rtelem := ti.elem

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnCborBytes

	var rvChanged bool

	var rv0 = rv
	var rv9 reflect.Value

	rvlen := rvLenSlice(rv)
	rvcap := rvCapSlice(rv)
	hasLen := containerLenS > 0
	if hasLen {
		if containerLenS > rvcap {
			oldRvlenGtZero := rvlen > 0
			rvlen1 := decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
			if rvlen1 == rvlen {
			} else if rvlen1 <= rvcap {
				if rvCanset {
					rvlen = rvlen1
					rvSetSliceLen(rv, rvlen)
				}
			} else if rvCanset {
				rvlen = rvlen1
				rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = !rvCanset
			} else {
				d.errorStr("cannot decode into non-settable slice")
			}
			if rvChanged && oldRvlenGtZero && rtelem0Mut {
				rvCopySlice(rv, rv0, rtelem)
			}
		} else if containerLenS != rvlen {
			if rvCanset {
				rvlen = containerLenS
				rvSetSliceLen(rv, rvlen)
			}
		}
	}

	var elemReset = d.h.SliceElementReset

	var j int

	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if rvCanset {
					rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
					rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
					rvcap = rvlen
					rvChanged = !rvCanset
				} else {
					d.errorStr("cannot decode into non-settable slice")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}

		if j >= rvlen {
			slh.ElemContainerState(j)

			if rvlen < rvcap {
				rvlen = rvcap
				if rvCanset {
					rvSetSliceLen(rv, rvlen)
				} else if rvChanged {
					rv = rvSlice(rv, rvlen)
				} else {
					d.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					d.onerror(errExpandSliceCannotChange)
				}
				rv, rvcap, rvCanset = rvGrowSlice(rv, f.ti, rvcap, 1)
				rvlen = rvcap
				rvChanged = !rvCanset
			}
		} else {
			slh.ElemContainerState(j)
		}
		rv9 = rvSliceIndex(rv, j, f.ti)
		if elemReset {
			rvSetZero(rv9)
		}
		d.decodeValue(rv9, fn)
	}
	if j < rvlen {
		if rvCanset {
			rvSetSliceLen(rv, j)
		} else if rvChanged {
			rv = rvSlice(rv, j)
		}

	} else if j == 0 && rvIsNil(rv) {
		if rvCanset {
			rv = rvSliceZeroCap(ti.rt)
			rvCanset = false
			rvChanged = true
		}
	}
	slh.End()

	if rvChanged {
		rvSetDirect(rv0, rv)
	}
}

func (d *decoderCborBytes) kArray(f *decFnInfo, rv reflect.Value) {

	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {

		if f.ti.elemkind != uint8(reflect.Uint8) {
			d.errorf("bytes/string in stream can decode into array of bytes, but not %v", f.ti.rt)
		}
		rvbs := rvGetArrayBytes(rv, nil)
		bs2 := d.decodeBytesInto(rvbs)
		if !byteSliceSameData(rvbs, bs2) && len(rvbs) > 0 && len(bs2) > 0 {
			copy(rvbs, bs2)
		}
		return
	}

	slh, containerLenS := d.decSliceHelperStart()

	if containerLenS == 0 {
		slh.End()
		return
	}

	rtelem := f.ti.elem
	for k := reflect.Kind(f.ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnCborBytes

	var rv9 reflect.Value

	rvlen := rv.Len()
	hasLen := containerLenS > 0
	if hasLen && containerLenS > rvlen {
		d.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
	}

	var elemReset = d.h.SliceElementReset

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {

		if j >= rvlen {
			slh.arrayCannotExpand(hasLen, rvlen, j, containerLenS)
			return
		}

		slh.ElemContainerState(j)
		rv9 = rvArrayIndex(rv, j, f.ti)
		if elemReset {
			rvSetZero(rv9)
		}

		if fn == nil {
			fn = d.fn(rtelem)
		}
		d.decodeValue(rv9, fn)
	}
	slh.End()
}

func (d *decoderCborBytes) kChan(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		d.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2 := d.d.DecodeBytes(nil)
		irv := rv2i(rv)
		ch, ok := irv.(chan<- byte)
		if !ok {
			ch = irv.(chan byte)
		}
		for _, b := range bs2 {
			ch <- b
		}
		return
	}

	var rvCanset = rv.CanSet()

	slh, containerLenS := d.decSliceHelperStart()

	if containerLenS == 0 {
		if rvCanset && rvIsNil(rv) {
			rvSetDirect(rv, reflect.MakeChan(ti.rt, 0))
		}
		slh.End()
		return
	}

	rtelem := ti.elem
	useTransient := decUseTransient && ti.elemkind != byte(reflect.Ptr) && ti.tielem.flagCanTransient

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnCborBytes

	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int
	hasLen := containerLenS > 0

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if hasLen {
					rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
				} else {
					rvlen = decDefChanCap
				}
				if rvCanset {
					rv = reflect.MakeChan(ti.rt, rvlen)
					rvChanged = true
				} else {
					d.errorStr("cannot decode into non-settable chan")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}
		slh.ElemContainerState(j)
		if rv9.IsValid() {
			rvSetZero(rv9)
		} else if decUseTransient && useTransient {
			rv9 = d.perType.TransientAddrK(ti.elem, reflect.Kind(ti.elemkind))
		} else {
			rv9 = rvZeroAddrK(ti.elem, reflect.Kind(ti.elemkind))
		}
		if !d.d.TryNil() {
			d.decodeValueNoCheckNil(rv9, fn)
		}
		rv.Send(rv9)
	}
	slh.End()

	if rvChanged {
		rvSetDirect(rv0, rv)
	}

}

func (d *decoderCborBytes) kMap(f *decFnInfo, rv reflect.Value) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	ti := f.ti
	if rvIsNil(rv) {
		rvlen := decInferLen(containerLen, d.h.MaxInitLen, int(ti.keysize+ti.elemsize))
		rvSetDirect(rv, makeMapReflect(ti.rt, rvlen))
	}

	if containerLen == 0 {
		d.mapEnd()
		return
	}

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := reflect.Kind(ti.elemkind)
	ktypeKind := reflect.Kind(ti.keykind)
	kfast := mapKeyFastKindFor(ktypeKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	vtypePtr := vtypeKind == reflect.Ptr
	ktypePtr := ktypeKind == reflect.Ptr

	vTransient := decUseTransient && !vtypePtr && ti.tielem.flagCanTransient
	kTransient := decUseTransient && !ktypePtr && ti.tikey.flagCanTransient

	var vtypeElem reflect.Type

	var keyFn, valFn *decFnCborBytes
	var ktypeLo, vtypeLo = ktype, vtype

	if ktypeKind == reflect.Ptr {
		for ktypeLo = ktype.Elem(); ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() {
		}
	}

	if vtypePtr {
		vtypeElem = vtype.Elem()
		for vtypeLo = vtypeElem; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() {
		}
	}

	rvkMut := !scalarBitset.isset(ti.keykind)
	rvvMut := !scalarBitset.isset(ti.elemkind)
	rvvCanNil := isnilBitset.isset(ti.elemkind)

	var rvk, rvkn, rvv, rvvn, rvva, rvvz reflect.Value

	var doMapGet, doMapSet bool

	if !d.h.MapValueReset {
		if rvvMut && (vtypeKind != reflect.Interface || !d.h.InterfaceReset) {
			doMapGet = true
			rvva = mapAddrLoopvarRV(vtype, vtypeKind)
		}
	}

	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId

	hasLen := containerLen > 0

	var kstrbs []byte
	var kstr2bs []byte
	var s string

	var callFnRvk bool

	fnRvk2 := func() (s string) {
		callFnRvk = false

		switch len(kstr2bs) {
		case 0:
		case 1:
			s = str256[kstr2bs[0] : kstr2bs[0]+1]
		default:
			s = d.mapKeyString(&callFnRvk, &kstrbs, &kstr2bs)
		}
		return
	}

	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		callFnRvk = false
		if j == 0 {

			if decUseTransient && vTransient && kTransient {
				rvk = d.perType.TransientAddr2K(ktype, ktypeKind)
			} else {
				rvk = rvZeroAddrK(ktype, ktypeKind)
			}
			if !rvkMut {
				rvkn = rvk
			}
			if !rvvMut {
				if decUseTransient && vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
			}
			if !ktypeIsString && keyFn == nil {
				keyFn = d.fn(ktypeLo)
			}
			if valFn == nil {
				valFn = d.fn(vtypeLo)
			}
		} else if rvkMut {
			rvSetZero(rvk)
		} else {
			rvk = rvkn
		}

		d.mapElemKey()
		if ktypeIsString {
			kstr2bs = d.d.DecodeStringAsBytes()
			rvSetString(rvk, fnRvk2())
		} else {
			d.decByteState = decByteStateNone
			d.decodeValue(rvk, keyFn)

			if ktypeIsIntf {
				if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp {
					kstr2bs = rvGetBytes(rvk2)
					rvSetIntf(rvk, rv4istr(fnRvk2()))
				}

			}
		}

		d.mapElemValue()

		if d.d.TryNil() {

			if !rvvz.IsValid() {
				rvvz = rvZeroK(vtype, vtypeKind)
			}
			if callFnRvk {
				s = d.string(kstr2bs)
				if ktypeIsString {
					rvSetString(rvk, s)
				} else {
					rvSetIntf(rvk, rv4istr(s))
				}
			}
			mapSet(rv, rvk, rvvz, kfast, visindirect, visref)
			continue
		}

		doMapSet = true

		if !rvvMut {
			rvv = rvvn
		} else if !doMapGet {
			goto NEW_RVV
		} else {
			rvv = mapGet(rv, rvk, rvva, kfast, visindirect, visref)
			if !rvv.IsValid() || (rvvCanNil && rvIsNil(rvv)) {
				goto NEW_RVV
			}
			switch vtypeKind {
			case reflect.Ptr, reflect.Map:
				doMapSet = false
			case reflect.Interface:

				rvvn = rvv.Elem()
				if k := rvvn.Kind(); (k == reflect.Ptr || k == reflect.Map) && !rvIsNil(rvvn) {
					d.decodeValueNoCheckNil(rvvn, nil)
					continue
				}

				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvSetIntf(rvvn, rvv)
				rvv = rvvn
			default:

				if decUseTransient && vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
				rvSetDirect(rvvn, rvv)
				rvv = rvvn
			}
		}
		goto DECODE_VALUE_NO_CHECK_NIL

	NEW_RVV:
		if vtypePtr {
			rvv = reflect.New(vtypeElem)
		} else if decUseTransient && vTransient {
			rvv = d.perType.TransientAddrK(vtype, vtypeKind)
		} else {
			rvv = rvZeroAddrK(vtype, vtypeKind)
		}

	DECODE_VALUE_NO_CHECK_NIL:
		d.decodeValueNoCheckNil(rvv, valFn)

		if doMapSet {
			if callFnRvk {
				s = d.string(kstr2bs)
				if ktypeIsString {
					rvSetString(rvk, s)
				} else {
					rvSetIntf(rvk, rv4istr(s))
				}
			}
			mapSet(rv, rvk, rvv, kfast, visindirect, visref)
		}
	}

	d.mapEnd()
}

type decoderCborBytes struct {
	panicHdl
	perType decPerType

	dh helperDecDriverCborBytes

	fp *fastpathDsCborBytes

	h *BasicHandle

	d cborDecDriverBytes

	decoderShared

	hh Handle

	mtid uintptr
	stid uintptr
}

func (d *decoderCborBytes) HandleName() string {
	return d.hh.Name()
}

func (d *decoderCborBytes) isBytes() bool {
	return d.bytes
}

func (d *decoderCborBytes) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()
	d.zeroCopy = d.h.ZeroCopy
	d.be = h.isBinary()
	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderShared, d).(*fastpathDsCborBytes)

	d.cbreak = d.js || d.cbor

	if d.bytes {
		d.rtidFn = &d.h.rtidFnsDecBytes
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtBytes
	} else {
		d.rtidFn = &d.h.rtidFnsDecIO
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtIO
	}

	d.reset()

}

func (d *decoderCborBytes) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
	d.decByteState = decByteStateNone
	d.depth = 0
	d.calls = 0

	d.maxdepth = decDefMaxDepth
	if d.h.MaxDepth > 0 {
		d.maxdepth = d.h.MaxDepth
	}
	d.mtid = 0
	d.stid = 0
	d.mtr = false
	d.str = false
	if d.h.MapType != nil {
		d.mtid = rt2id(d.h.MapType)
		_, d.mtr = fastpathAvIndex(d.mtid)
	}
	if d.h.SliceType != nil {
		d.stid = rt2id(d.h.SliceType)
		_, d.str = fastpathAvIndex(d.stid)
	}
}

func (d *decoderCborBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderCborBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
	return
}

func (d *decoderCborBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderCborBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderCborBytes) Decode(v interface{}) (err error) {

	if !debugging {
		defer func() {
			if x := recover(); x != nil {
				panicValToErr(d, x, &d.err)
				err = d.err
			}
		}()
	}

	d.MustDecode(v)
	return
}

func (d *decoderCborBytes) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderCborBytes) Release() {
}

func (d *decoderCborBytes) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *decoderCborBytes) nextValueBytes(start []byte) []byte {
	return d.d.nextValueBytes(start)
}

func (d *decoderCborBytes) decode(iv interface{}) {

	if iv == nil {
		d.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {

	case reflect.Value:
		if x, _ := isDecodeable(v); !x {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)
	case *string:
		*v = d.stringZC(d.d.DecodeStringAsBytes())
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *float32:
		*v = d.d.DecodeFloat32()
	case *float64:
		*v = d.d.DecodeFloat64()
	case *complex64:
		*v = complex(d.d.DecodeFloat32(), 0)
	case *complex128:
		*v = complex(d.d.DecodeFloat64(), 0)
	case *[]byte:
		*v = d.decodeBytesInto(*v)
	case []byte:

		b := d.decodeBytesInto(v[:len(v):len(v)])
		if !(len(b) > 0 && len(b) == len(v) && &b[0] == &v[0]) {
			copy(v, b)
		}
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface{}:
		d.decodeValue(rv4iptr(v), nil)

	default:

		if skipFastpathTypeSwitchInDirectCall || !d.dh.fastpathDecodeTypeSwitch(iv, d) {
			v := reflect.ValueOf(iv)
			if x, _ := isDecodeable(v); !x {
				d.haltAsNotDecodeable(v)
			}
			d.decodeValue(v, nil)
		}
	}
}

func (d *decoderCborBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderCborBytes) decodeValue(rv reflect.Value, fn *decFnCborBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *decoderCborBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnCborBytes) {

	var rvp reflect.Value
	var rvpValid bool
PTR:
	if rv.Kind() == reflect.Ptr {
		rvpValid = true
		if rvIsNil(rv) {
			rvSetDirect(rv, reflect.New(rv.Type().Elem()))
		}
		rvp = rv
		rv = rv.Elem()
		goto PTR
	}

	if fn == nil {
		fn = d.fn(rv.Type())
	}
	if fn.i.addrD {
		if rvpValid {
			rv = rvp
		} else if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else if fn.i.addrDf {
			d.errorStr("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *decoderCborBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			d.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			d.errorStr("no matching struct field found when decoding stream map with key " + rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderCborBytes) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		d.errorf("cannot expand array len during decode from %v to %v", any(sliceLen), any(streamLen))
	}
}

func (d *decoderCborBytes) haltAsNotDecodeable(rv reflect.Value) {
	if !rv.IsValid() {
		d.onerror(errCannotDecodeIntoNil)
	}

	if !rv.CanInterface() {
		d.errorf("cannot decode into a value without an interface: %v", rv)
	}
	d.errorf("cannot decode into value of kind: %v, %#v", rv.Kind(), rv2i(rv))
}

func (d *decoderCborBytes) depthIncr() {
	d.depth++
	if d.depth >= d.maxdepth {
		d.onerror(errMaxDepthExceeded)
	}
}

func (d *decoderCborBytes) depthDecr() {
	d.depth--
}

func (d *decoderCborBytes) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = zeroByteSlice
	}
	return d.d.DecodeBytes(in)
}

func (d *decoderCborBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes(zeroByteSlice)
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderCborBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.NumBytesRead(), false)
}

func (d *decoderCborBytes) NumBytesRead() int {
	return int(d.d.NumBytesRead())
}

func (d *decoderCborBytes) checkBreak() (v bool) {
	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *decoderCborBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *decoderCborBytes) mapStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerMapStart
	}
	return v
}

func (d *decoderCborBytes) mapElemKey() {
	d.d.ReadMapElemKey()
	d.c = containerMapKey
}

func (d *decoderCborBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderCborBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborBytes) arrayStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerArrayStart
	}
	return v
}

func (d *decoderCborBytes) arrayElem() {
	d.d.ReadArrayElem()
	d.c = containerArrayElem
}

func (d *decoderCborBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var s interface{}
	rv := reflect.ValueOf(v)
	rv2 := rv.Elem()
	rvk := rv2.Kind()
	if rvk == reflect.Struct || rvk == reflect.Array {
		s = ext.ConvertExt(v)
	} else {
		s = ext.ConvertExt(rv2i(rv2))
	}
	rv = reflect.ValueOf(s)

	if !rv.CanAddr() {
		rvk = rv.Kind()
		rv2 = d.oneShotAddrRV(rv.Type(), rvk)
		if rvk == reflect.Interface {
			rvSetIntf(rv2, rv)
		} else {
			rvSetDirect(rv2, rv)
		}
		rv = rv2
	}

	d.decodeValue(rv, nil)
	ext.UpdateExt(v, rv2i(rv))
}

func (d *decoderCborBytes) oneShotAddrRV(rvt reflect.Type, rvk reflect.Kind) reflect.Value {
	if decUseTransient &&
		(numBoolStrSliceBitset.isset(byte(rvk)) ||
			((rvk == reflect.Struct || rvk == reflect.Array) &&
				d.h.getTypeInfo(rt2id(rvt), rvt).flagCanTransient)) {
		return d.perType.TransientAddrK(rvt, rvk)
	}
	return rvZeroAddrK(rvt, rvk)
}

func (d *decoderCborBytes) fn(t reflect.Type) *decFnCborBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, false)
}

func (d *decoderCborBytes) fnNoExt(t reflect.Type) *decFnCborBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, true)
}

type decSliceHelperCborBytes struct {
	d     *decoderCborBytes
	ct    valueType
	Array bool
	IsNil bool
}

func (d *decoderCborBytes) decSliceHelperStart() (x decSliceHelperCborBytes, clen int) {
	x.ct = d.d.ContainerType()
	x.d = d
	switch x.ct {
	case valueTypeNil:
		x.IsNil = true
	case valueTypeArray:
		x.Array = true
		clen = d.arrayStart(d.d.ReadArrayStart())
	case valueTypeMap:
		clen = d.mapStart(d.d.ReadMapStart())
		clen += clen
	default:
		d.errorStr2("to decode into a slice, expect map/array - got ", x.ct.String())
	}
	return
}

func (x decSliceHelperCborBytes) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelperCborBytes) ElemContainerState(index int) {

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 {
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelperCborBytes) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
	x.d.arrayCannotExpand(lenv, j+1)

	x.ElemContainerState(j)
	x.d.swallow()
	j++
	for ; x.d.containerNext(j, containerLenS, hasLen); j++ {
		x.ElemContainerState(j)
		x.d.swallow()
	}
	x.End()
}

func (helperDecDriverCborBytes) newDecoderBytes(in []byte, h Handle) *decoderCborBytes {
	var c1 decoderCborBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverCborBytes) newDecoderIO(in io.Reader, h Handle) *decoderCborBytes {
	var c1 decoderCborBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverCborBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsCborBytes) (f *fastpathDCborBytes, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

type decFnCborBytes struct {
	i  decFnInfo
	fd func(*decoderCborBytes, *decFnInfo, reflect.Value)
}
type decRtidFnCborBytes struct {
	rtid uintptr
	fn   *decFnCborBytes
}

func (helperDecDriverCborBytes) decFindRtidFn(s []decRtidFnCborBytes, rtid uintptr) (i uint, fn *decFnCborBytes) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperDecDriverCborBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnCborBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnCborBytes](v))
	}
	return
}

func (dh helperDecDriverCborBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsCborBytes,
	checkExt bool) (fn *decFnCborBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverCborBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnCborBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverCborBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnCborBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnCborBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnCborBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnCborBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverCborBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborBytes) {
	fn = new(decFnCborBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderCborBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderCborBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderCborBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderCborBytes).ext
		fi.addrD = true
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fd = (*decoderCborBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderCborBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderCborBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderCborBytes).textUnmarshal
		fi.addrD = ti.flagTextUnmarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fd = fp[idx].decfn
					fi.addrD = true
					fi.addrDf = false
					if rk == reflect.Array {
						fi.addrD = false
					}
				}
			} else {

				xfe, xrt := dh.decFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf2 := xfe.decfn
					if rk == reflect.Array {
						fi.addrD = false
						fn.fd = func(d *decoderCborBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderCborBytes, xf *decFnInfo, xrv reflect.Value) {
							if xrv.Kind() == reflect.Ptr {
								xfnf2(d, xf, rvConvert(xrv, xptr2rt))
							} else {
								xfnf2(d, xf, rvConvert(xrv, xrt))
							}
						}
					}
				}
			}
		}
		if fn.fd == nil {
			switch rk {
			case reflect.Bool:
				fn.fd = (*decoderCborBytes).kBool
			case reflect.String:
				fn.fd = (*decoderCborBytes).kString
			case reflect.Int:
				fn.fd = (*decoderCborBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderCborBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderCborBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderCborBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderCborBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderCborBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderCborBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderCborBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderCborBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderCborBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderCborBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderCborBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderCborBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderCborBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderCborBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderCborBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderCborBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderCborBytes).kArray
			case reflect.Struct:
				fn.fd = (*decoderCborBytes).kStruct
			case reflect.Map:
				fn.fd = (*decoderCborBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderCborBytes).kInterface
			default:

				fn.fd = (*decoderCborBytes).kErr
			}
		}
	}
	return
}

type helperDecDriverCborBytes struct{}
type fastpathECborBytes struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderCborBytes, *encFnInfo, reflect.Value)
}
type fastpathDCborBytes struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderCborBytes, *decFnInfo, reflect.Value)
}
type fastpathEsCborBytes [56]fastpathECborBytes
type fastpathDsCborBytes [56]fastpathDCborBytes
type fastpathETCborBytes struct{}
type fastpathDTCborBytes struct{}

func (helperEncDriverCborBytes) fastpathEList() *fastpathEsCborBytes {
	var i uint = 0
	var s fastpathEsCborBytes
	fn := func(v interface{}, fe func(*encoderCborBytes, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathECborBytes{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderCborBytes).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderCborBytes).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderCborBytes).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderCborBytes).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderCborBytes).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderCborBytes).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderCborBytes).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderCborBytes).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderCborBytes).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderCborBytes).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderCborBytes).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderCborBytes).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderCborBytes).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderCborBytes).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderCborBytes).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderCborBytes).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderCborBytes).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderCborBytes).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderCborBytes).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderCborBytes).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderCborBytes).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderCborBytes).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderCborBytes).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderCborBytes).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderCborBytes).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderCborBytes).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderCborBytes).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderCborBytes).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderCborBytes).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderCborBytes).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderCborBytes).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderCborBytes).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderCborBytes).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderCborBytes).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderCborBytes).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderCborBytes).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderCborBytes).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderCborBytes).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderCborBytes).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderCborBytes).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderCborBytes).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderCborBytes).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderCborBytes).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderCborBytes).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderCborBytes).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderCborBytes).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderCborBytes).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderCborBytes).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderCborBytes).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderCborBytes).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderCborBytes).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderCborBytes).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderCborBytes).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderCborBytes).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderCborBytes).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderCborBytes).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverCborBytes) fastpathDList() *fastpathDsCborBytes {
	var i uint = 0
	var s fastpathDsCborBytes
	fn := func(v interface{}, fd func(*decoderCborBytes, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDCborBytes{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderCborBytes).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderCborBytes).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderCborBytes).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderCborBytes).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderCborBytes).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderCborBytes).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderCborBytes).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderCborBytes).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderCborBytes).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderCborBytes).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderCborBytes).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderCborBytes).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderCborBytes).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderCborBytes).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderCborBytes).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderCborBytes).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderCborBytes).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderCborBytes).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderCborBytes).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderCborBytes).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderCborBytes).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderCborBytes).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderCborBytes).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderCborBytes).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderCborBytes).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderCborBytes).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderCborBytes).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderCborBytes).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderCborBytes).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderCborBytes).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderCborBytes).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderCborBytes).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderCborBytes).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderCborBytes).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderCborBytes).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderCborBytes).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderCborBytes).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderCborBytes).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderCborBytes).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderCborBytes).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderCborBytes).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderCborBytes).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderCborBytes).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderCborBytes).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderCborBytes).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderCborBytes).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderCborBytes).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderCborBytes).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderCborBytes).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderCborBytes).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderCborBytes).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderCborBytes).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderCborBytes).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderCborBytes).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderCborBytes).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderCborBytes).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverCborBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderCborBytes) bool {
	var ft fastpathETCborBytes
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

func (e *encoderCborBytes) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceIntfV(v []interface{}, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceIntfV(v []interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceStringV(v []string, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceStringV(v []string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceBytesV(v [][]byte, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceBytesV(v [][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceFloat32V(v []float32, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceFloat32V(v []float32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceFloat64V(v []float64, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceFloat64V(v []float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceUint8V(v []uint8, e *encoderCborBytes) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETCborBytes) EncAsMapSliceUint8V(v []uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceUint64V(v []uint64, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceUint64V(v []uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceIntV(v []int, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceIntV(v []int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceInt32V(v []int32, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceInt32V(v []int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceInt64V(v []int64, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceInt64V(v []int64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceBoolV(v []bool, e *encoderCborBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborBytes) EncAsMapSliceBoolV(v []bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETCborBytes) EncMapStringIntfV(v map[string]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETCborBytes) EncMapStringStringV(v map[string]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETCborBytes) EncMapStringBytesV(v map[string][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETCborBytes) EncMapStringUint8V(v map[string]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETCborBytes) EncMapStringUint64V(v map[string]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETCborBytes) EncMapStringIntV(v map[string]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETCborBytes) EncMapStringInt32V(v map[string]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETCborBytes) EncMapStringFloat64V(v map[string]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETCborBytes) EncMapStringBoolV(v map[string]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETCborBytes) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETCborBytes) EncMapUint8StringV(v map[uint8]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETCborBytes) EncMapUint8BytesV(v map[uint8][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETCborBytes) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETCborBytes) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETCborBytes) EncMapUint8IntV(v map[uint8]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETCborBytes) EncMapUint8Int32V(v map[uint8]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETCborBytes) EncMapUint8Float64V(v map[uint8]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETCborBytes) EncMapUint8BoolV(v map[uint8]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETCborBytes) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETCborBytes) EncMapUint64StringV(v map[uint64]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETCborBytes) EncMapUint64BytesV(v map[uint64][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETCborBytes) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETCborBytes) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETCborBytes) EncMapUint64IntV(v map[uint64]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETCborBytes) EncMapUint64Int32V(v map[uint64]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETCborBytes) EncMapUint64Float64V(v map[uint64]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETCborBytes) EncMapUint64BoolV(v map[uint64]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETCborBytes) EncMapIntIntfV(v map[int]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETCborBytes) EncMapIntStringV(v map[int]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETCborBytes) EncMapIntBytesV(v map[int][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETCborBytes) EncMapIntUint8V(v map[int]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETCborBytes) EncMapIntUint64V(v map[int]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETCborBytes) EncMapIntIntV(v map[int]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETCborBytes) EncMapIntInt32V(v map[int]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETCborBytes) EncMapIntFloat64V(v map[int]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETCborBytes) EncMapIntBoolV(v map[int]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETCborBytes) EncMapInt32IntfV(v map[int32]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETCborBytes) EncMapInt32StringV(v map[int32]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETCborBytes) EncMapInt32BytesV(v map[int32][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETCborBytes) EncMapInt32Uint8V(v map[int32]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETCborBytes) EncMapInt32Uint64V(v map[int32]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETCborBytes) EncMapInt32IntV(v map[int32]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETCborBytes) EncMapInt32Int32V(v map[int32]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETCborBytes) EncMapInt32Float64V(v map[int32]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
	ft.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETCborBytes) EncMapInt32BoolV(v map[int32]bool, e *encoderCborBytes) {
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

func (helperDecDriverCborBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderCborBytes) bool {
	var ft fastpathDTCborBytes
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

func (d *decoderCborBytes) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceIntfX(vp *[]interface{}, d *decoderCborBytes) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceIntfY(v []interface{}, d *decoderCborBytes) (v2 []interface{}, changed bool) {
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
func (fastpathDTCborBytes) DecSliceIntfN(v []interface{}, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceStringX(vp *[]string, d *decoderCborBytes) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceStringY(v []string, d *decoderCborBytes) (v2 []string, changed bool) {
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
func (fastpathDTCborBytes) DecSliceStringN(v []string, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceBytesX(vp *[][]byte, d *decoderCborBytes) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceBytesY(v [][]byte, d *decoderCborBytes) (v2 [][]byte, changed bool) {
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
func (fastpathDTCborBytes) DecSliceBytesN(v [][]byte, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceFloat32X(vp *[]float32, d *decoderCborBytes) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceFloat32Y(v []float32, d *decoderCborBytes) (v2 []float32, changed bool) {
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
func (fastpathDTCborBytes) DecSliceFloat32N(v []float32, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceFloat64X(vp *[]float64, d *decoderCborBytes) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceFloat64Y(v []float64, d *decoderCborBytes) (v2 []float64, changed bool) {
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
func (fastpathDTCborBytes) DecSliceFloat64N(v []float64, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceUint8X(vp *[]uint8, d *decoderCborBytes) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceUint8Y(v []uint8, d *decoderCborBytes) (v2 []uint8, changed bool) {
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
func (fastpathDTCborBytes) DecSliceUint8N(v []uint8, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceUint64X(vp *[]uint64, d *decoderCborBytes) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceUint64Y(v []uint64, d *decoderCborBytes) (v2 []uint64, changed bool) {
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
func (fastpathDTCborBytes) DecSliceUint64N(v []uint64, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceIntX(vp *[]int, d *decoderCborBytes) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceIntY(v []int, d *decoderCborBytes) (v2 []int, changed bool) {
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
func (fastpathDTCborBytes) DecSliceIntN(v []int, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceInt32X(vp *[]int32, d *decoderCborBytes) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceInt32Y(v []int32, d *decoderCborBytes) (v2 []int32, changed bool) {
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
func (fastpathDTCborBytes) DecSliceInt32N(v []int32, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceInt64X(vp *[]int64, d *decoderCborBytes) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceInt64Y(v []int64, d *decoderCborBytes) (v2 []int64, changed bool) {
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
func (fastpathDTCborBytes) DecSliceInt64N(v []int64, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecSliceBoolX(vp *[]bool, d *decoderCborBytes) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborBytes) DecSliceBoolY(v []bool, d *decoderCborBytes) (v2 []bool, changed bool) {
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
func (fastpathDTCborBytes) DecSliceBoolN(v []bool, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringIntfX(vp *map[string]interface{}, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapStringIntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborBytes) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringStringX(vp *map[string]string, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]string, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapStringStringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringStringL(v map[string]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringBytesX(vp *map[string][]byte, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string][]byte, decInferLen(containerLen, d.h.MaxInitLen, 40))
		}
		if containerLen != 0 {
			f.DecMapStringBytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborBytes) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringUint8X(vp *map[string]uint8, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]uint8, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapStringUint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringUint64X(vp *map[string]uint64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]uint64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapStringUint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringIntX(vp *map[string]int, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]int, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapStringIntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringIntL(v map[string]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringInt32X(vp *map[string]int32, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]int32, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			f.DecMapStringInt32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringFloat64X(vp *map[string]float64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]float64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapStringFloat64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapStringBoolX(vp *map[string]bool, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]bool, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapStringBoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapUint8IntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborBytes) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8StringX(vp *map[uint8]string, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]string, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapUint8StringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8][]byte, decInferLen(containerLen, d.h.MaxInitLen, 25))
		}
		if containerLen != 0 {
			f.DecMapUint8BytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborBytes) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]uint8, decInferLen(containerLen, d.h.MaxInitLen, 2))
		}
		if containerLen != 0 {
			f.DecMapUint8Uint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]uint64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint8Uint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8IntX(vp *map[uint8]int, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]int, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint8IntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8Int32X(vp *map[uint8]int32, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]int32, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			f.DecMapUint8Int32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8Float64X(vp *map[uint8]float64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]float64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint8Float64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint8BoolX(vp *map[uint8]bool, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]bool, decInferLen(containerLen, d.h.MaxInitLen, 2))
		}
		if containerLen != 0 {
			f.DecMapUint8BoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapUint64IntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborBytes) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64StringX(vp *map[uint64]string, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapUint64StringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapUint64BytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborBytes) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint64Uint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapUint64Uint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64IntX(vp *map[uint64]int, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapUint64IntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64Int32X(vp *map[uint64]int32, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapUint64Int32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64Float64X(vp *map[uint64]float64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapUint64Float64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapUint64BoolX(vp *map[uint64]bool, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint64BoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntIntfX(vp *map[int]interface{}, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapIntIntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborBytes) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntStringX(vp *map[int]string, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapIntStringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntStringL(v map[int]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntBytesX(vp *map[int][]byte, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapIntBytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborBytes) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntUint8X(vp *map[int]uint8, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapIntUint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntUint64X(vp *map[int]uint64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapIntUint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntIntX(vp *map[int]int, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapIntIntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntIntL(v map[int]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntInt32X(vp *map[int]int32, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapIntInt32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntFloat64X(vp *map[int]float64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapIntFloat64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapIntBoolX(vp *map[int]bool, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapIntBoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			f.DecMapInt32IntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborBytes) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32StringX(vp *map[int32]string, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]string, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			f.DecMapInt32StringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32BytesX(vp *map[int32][]byte, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32][]byte, decInferLen(containerLen, d.h.MaxInitLen, 28))
		}
		if containerLen != 0 {
			f.DecMapInt32BytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderCborBytes) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborBytes) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]uint8, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			f.DecMapInt32Uint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]uint64, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapInt32Uint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32IntX(vp *map[int32]int, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]int, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapInt32IntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32Int32X(vp *map[int32]int32, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]int32, decInferLen(containerLen, d.h.MaxInitLen, 8))
		}
		if containerLen != 0 {
			f.DecMapInt32Int32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32Float64X(vp *map[int32]float64, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]float64, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapInt32Float64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (f fastpathDTCborBytes) DecMapInt32BoolX(vp *map[int32]bool, d *decoderCborBytes) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]bool, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			f.DecMapInt32BoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborBytes) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderCborBytes) {
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

type cborEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoState
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT

	h   *CborHandle
	e   *encoderShared
	w   bytesEncAppender
	enc encoderI

	b [40]byte
}

func (e *cborEncDriverBytes) EncodeNil() {
	e.w.writen1(cborBdNil)
}

func (e *cborEncDriverBytes) EncodeBool(b bool) {
	if b {
		e.w.writen1(cborBdTrue)
	} else {
		e.w.writen1(cborBdFalse)
	}
}

func (e *cborEncDriverBytes) EncodeFloat32(f float32) {
	b := math.Float32bits(f)
	if e.h.OptimumSize {
		if h := floatToHalfFloatBits(b); halfFloatToFloatBits(h) == b {
			e.w.writen1(cborBdFloat16)
			e.w.writen2(bigen.PutUint16(h))
			return
		}
	}
	e.w.writen1(cborBdFloat32)
	e.w.writen4(bigen.PutUint32(b))
}

func (e *cborEncDriverBytes) EncodeFloat64(f float64) {
	if e.h.OptimumSize {
		if f32 := float32(f); float64(f32) == f {
			e.EncodeFloat32(f32)
			return
		}
	}
	e.w.writen1(cborBdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *cborEncDriverBytes) encUint(v uint64, bd byte) {
	if v <= 0x17 {
		e.w.writen1(byte(v) + bd)
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd+0x18, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 0x19)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 0x1a)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 0x1b)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *cborEncDriverBytes) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-1-v), cborBaseNegInt)
	} else {
		e.encUint(uint64(v), cborBaseUint)
	}
}

func (e *cborEncDriverBytes) EncodeUint(v uint64) {
	e.encUint(v, cborBaseUint)
}

func (e *cborEncDriverBytes) encLen(bd byte, length int) {
	e.encUint(uint64(length), bd)
}

func (e *cborEncDriverBytes) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else if e.h.TimeRFC3339 {
		e.encUint(0, cborBaseTag)
		e.encStringBytesS(cborBaseString, stringView(t.AppendFormat(e.b[:0], time.RFC3339Nano)))
	} else {
		e.encUint(1, cborBaseTag)
		t = t.UTC().Round(time.Microsecond)
		sec, nsec := t.Unix(), uint64(t.Nanosecond())
		if nsec == 0 {
			e.EncodeInt(sec)
		} else {
			e.EncodeFloat64(float64(sec) + float64(nsec)/1e9)
		}
	}
}

func (e *cborEncDriverBytes) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	e.encUint(uint64(xtag), cborBaseTag)
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.EncodeNil()
	} else {
		e.enc.encode(v)
	}
}

func (e *cborEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encUint(uint64(re.Tag), cborBaseTag)

	if re.Value != nil {
		e.enc.encode(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *cborEncDriverBytes) WriteArrayStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteArray)
	} else {
		e.encLen(cborBaseArray, length)
	}
}

func (e *cborEncDriverBytes) WriteMapStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteMap)
	} else {
		e.encLen(cborBaseMap, length)
	}
}

func (e *cborEncDriverBytes) WriteMapEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverBytes) WriteArrayEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverBytes) EncodeString(v string) {
	bb := cborBaseString
	if e.h.StringToRaw {
		bb = cborBaseBytes
	}
	e.encStringBytesS(bb, v)
}

func (e *cborEncDriverBytes) EncodeStringBytesRaw(v []byte) {
	if v == nil {
		e.EncodeNil()
	} else {
		e.encStringBytesS(cborBaseBytes, stringView(v))
	}
}

func (e *cborEncDriverBytes) encStringBytesS(bb byte, v string) {
	if e.h.IndefiniteLength {
		if bb == cborBaseBytes {
			e.w.writen1(cborBdIndefiniteBytes)
		} else {
			e.w.writen1(cborBdIndefiniteString)
		}
		var vlen uint = uint(len(v))
		blen := vlen / 4
		if blen == 0 {
			blen = 64
		} else if blen > 1024 {
			blen = 1024
		}
		for i := uint(0); i < vlen; {
			var v2 string
			i2 := i + blen
			if i2 >= i && i2 < vlen {
				v2 = v[i:i2]
			} else {
				v2 = v[i:]
			}
			e.encLen(bb, len(v2))
			e.w.writestr(v2)
			i = i2
		}
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(bb, len(v))
		e.w.writestr(v)
	}
}

type cborDecDriverBytes struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	noBuiltInTypes

	h   *CborHandle
	d   *decoderShared
	r   bytesDecReader
	dec decoderI
	bdAndBdread
	st    bool
	bytes bool
}

func (d *cborDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *cborDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdNil || d.bd == cborBdUndefined {
		d.bdRead = false
		return true
	}
	return
}

func (d *cborDecDriverBytes) TryNil() bool {
	return d.advanceNil()
}

func (d *cborDecDriverBytes) skipTags() {
	for d.bd>>5 == cborMajorTag {
		d.decUint()
		d.bd = d.r.readn1()
	}
}

func (d *cborDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdNil {
		d.bdRead = false
		return valueTypeNil
	}
	major := d.bd >> 5
	if major == cborMajorBytes {
		return valueTypeBytes
	} else if major == cborMajorString {
		return valueTypeString
	} else if major == cborMajorArray {
		return valueTypeArray
	} else if major == cborMajorMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *cborDecDriverBytes) CheckBreak() (v bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdBreak {
		d.bdRead = false
		v = true
	}
	return
}

func (d *cborDecDriverBytes) decUint() (ui uint64) {
	v := d.bd & 0x1f
	if v <= 0x17 {
		ui = uint64(v)
	} else if v == 0x18 {
		ui = uint64(d.r.readn1())
	} else if v == 0x19 {
		ui = uint64(bigen.Uint16(d.r.readn2()))
	} else if v == 0x1a {
		ui = uint64(bigen.Uint32(d.r.readn4()))
	} else if v == 0x1b {
		ui = uint64(bigen.Uint64(d.r.readn8()))
	} else {
		halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
	}
	return
}

func (d *cborDecDriverBytes) decLen() int {
	return int(d.decUint())
}

func (d *cborDecDriverBytes) decAppendIndefiniteBytes(bs []byte, major byte) []byte {
	d.bdRead = false
	for !d.CheckBreak() {
		chunkMajor := d.bd >> 5
		if chunkMajor != major {
			halt.errorf("malformed indefinite string/bytes %x (%s); contains chunk with major type %v, expected %v",
				d.bd, cbordesc(d.bd), chunkMajor, major)
		}
		n := uint(d.decLen())
		oldLen := uint(len(bs))
		newLen := oldLen + n
		if newLen > uint(cap(bs)) {
			bs2 := make([]byte, newLen, 2*uint(cap(bs))+n)
			copy(bs2, bs)
			bs = bs2
		} else {
			bs = bs[:newLen]
		}
		d.r.readb(bs[oldLen:newLen])
		if d.h.ValidateUnicode && major == cborMajorString && !utf8.Valid(bs[oldLen:newLen]) {
			halt.errorf("indefinite-length text string contains chunk that is not a valid utf-8 sequence: 0x%x", bs[oldLen:newLen])
		}
		d.bdRead = false
	}
	d.bdRead = false
	return bs
}

func (d *cborDecDriverBytes) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.r.readn2()))))
	case cborBdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case cborBdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *cborDecDriverBytes) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd >> 5 {
	case cborMajorUint:
		ui = d.decUint()
	case cborMajorNegInt:
		ui = d.decUint()
		neg = true
	default:
		ok = false
	}
	return
}

func (d *cborDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdTrue {
		b = true
	} else if d.bd == cborBdFalse {
	} else {
		halt.errorf("not bool - %s %x/%s", msgBadDesc, d.bd, cbordesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.st {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorMap {
		halt.errorf("error reading map; got major type: %x, expected %x/%s", d.bd>>5, cborMajorMap, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverBytes) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.st {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorArray {
		halt.errorf("invalid array; got major type: %x, expect: %x/%s", d.bd>>5, cborMajorArray, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverBytes) DecodeBytes(bs []byte) (bsOut []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		d.bdRead = false
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			return d.decAppendIndefiniteBytes(d.d.b[:0], d.bd>>5)
		}
		return d.decAppendIndefiniteBytes(bs[:0], d.bd>>5)
	}
	if d.bd == cborBdIndefiniteArray {
		d.bdRead = false
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			bs = d.d.b[:0]
		} else {
			bs = bs[:0]
		}
		for !d.CheckBreak() {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}
	if d.bd>>5 == cborMajorArray {
		d.bdRead = false
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			bs = d.d.b[:]
		}
		slen := d.decLen()
		var changed bool
		if bs, changed = usableByteSlice(bs, slen); changed {
			d.d.decByteState = decByteStateNone
		}
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}
	clen := d.decLen()
	d.bdRead = false
	if d.bytes && d.h.ZeroCopy {
		d.d.decByteState = decByteStateZerocopy
		return d.r.readx(uint(clen))
	}
	if bs == nil {
		d.d.decByteState = decByteStateReuseBuf
		bs = d.d.b[:]
	}
	return decByteSlice(&d.r, clen, d.h.MaxInitLen, bs)
}

func (d *cborDecDriverBytes) DecodeStringAsBytes() (s []byte) {
	s = d.DecodeBytes(nil)
	if d.h.ValidateUnicode && !utf8.Valid(s) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", s)
	}
	return
}

func (d *cborDecDriverBytes) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
}

func (d *cborDecDriverBytes) decodeTime(xtag uint64) (t time.Time) {
	switch xtag {
	case 0:
		var err error
		t, err = time.Parse(time.RFC3339, stringView(d.DecodeStringAsBytes()))
		halt.onerror(err)
	case 1:
		f1, f2 := math.Modf(d.DecodeFloat64())
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		halt.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	}
	t = t.UTC().Round(time.Microsecond)
	return
}

func (d *cborDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	realxtag := d.decUint()
	d.bdRead = false
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = realxtag
		d.dec.decode(&re.Value)
	} else if xtag != realxtag {
		halt.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
	} else if ext == SelfExt {
		d.dec.decodeAs(rv, basetype, false)
	} else {
		d.dec.interfaceExtConvertAndDecode(rv, ext)
	}
	d.bdRead = false
}

func (d *cborDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd >> 5 {
	case cborMajorUint:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case cborMajorNegInt:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case cborMajorBytes:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString, d.h.ZeroCopy)
	case cborMajorString:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case cborMajorArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborMajorMap:
		n.v = valueTypeMap
		decodeFurther = true
	case cborMajorTag:
		n.v = valueTypeExt
		n.u = d.decUint()
		n.l = nil
		if n.u == 0 || n.u == 1 {
			d.bdRead = false
			n.v = valueTypeTime
			n.t = d.decodeTime(n.u)
		} else if d.st && d.h.getExtForTag(n.u) == nil {

			d.bdRead = false
			d.DecodeNaked()
			return
		}
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined:
			n.v = valueTypeNil
		case cborBdFalse:
			n.v = valueTypeBool
			n.b = false
		case cborBdTrue:
			n.v = valueTypeBool
			n.b = true
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			n.v = valueTypeFloat
			n.f = d.DecodeFloat64()
		default:
			halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	}
	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *cborDecDriverBytes) uintBytes() (v []byte, ui uint64) {

	switch vv := d.bd & 0x1f; vv {
	case 0x18:
		v = d.r.readx(1)
		ui = uint64(v[0])
	case 0x19:
		v = d.r.readx(2)
		ui = uint64(bigenstd.Uint16(v))
	case 0x1a:
		v = d.r.readx(4)
		ui = uint64(bigenstd.Uint32(v))
	case 0x1b:
		v = d.r.readx(8)
		ui = uint64(bigenstd.Uint64(v))
	default:
		if vv > 0x1b {
			halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
		}
		ui = uint64(vv)
	}
	return
}

func (d *cborDecDriverBytes) nextValueBytes(v0 []byte) (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	v = v0
	var h decNextValueBytesHelper

	var cursor uint
	if d.bytes {
		cursor = d.r.numread() - 1
	}
	h.append1(&v, d.bytes, d.bd)
	v = d.nextValueBytesBdReadR(v)
	d.bdRead = false

	if d.bytes {
		v = d.r.bytesReadFrom(cursor)
	}
	return
}

func (d *cborDecDriverBytes) nextValueBytesR(v0 []byte) (v []byte) {
	d.readNextBd()
	v = v0
	var h decNextValueBytesHelper
	h.append1(&v, d.bytes, d.bd)
	return d.nextValueBytesBdReadR(v)
}

func (d *cborDecDriverBytes) nextValueBytesBdReadR(v0 []byte) (v []byte) {
	v = v0
	var h decNextValueBytesHelper

	var bs []byte
	var ui uint64

	switch d.bd >> 5 {
	case cborMajorUint, cborMajorNegInt:
		bs, _ = d.uintBytes()
		h.appendN(&v, d.bytes, bs...)
	case cborMajorString, cborMajorBytes:
		if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
			for {
				d.readNextBd()
				h.append1(&v, d.bytes, d.bd)
				if d.bd == cborBdBreak {
					break
				}
				bs, ui = d.uintBytes()
				h.appendN(&v, d.bytes, bs...)
				h.appendN(&v, d.bytes, d.r.readx(uint(ui))...)
			}
		} else {
			bs, ui = d.uintBytes()
			h.appendN(&v, d.bytes, bs...)
			h.appendN(&v, d.bytes, d.r.readx(uint(ui))...)
		}
	case cborMajorArray:
		if d.bd == cborBdIndefiniteArray {
			for {
				d.readNextBd()
				h.append1(&v, d.bytes, d.bd)
				if d.bd == cborBdBreak {
					break
				}
				v = d.nextValueBytesBdReadR(v)
			}
		} else {
			bs, ui = d.uintBytes()
			h.appendN(&v, d.bytes, bs...)
			for i := uint64(0); i < ui; i++ {
				v = d.nextValueBytesR(v)
			}
		}
	case cborMajorMap:
		if d.bd == cborBdIndefiniteMap {
			for {
				d.readNextBd()
				h.append1(&v, d.bytes, d.bd)
				if d.bd == cborBdBreak {
					break
				}
				v = d.nextValueBytesBdReadR(v)
				v = d.nextValueBytesR(v)
			}
		} else {
			bs, ui = d.uintBytes()
			h.appendN(&v, d.bytes, bs...)
			for i := uint64(0); i < ui; i++ {
				v = d.nextValueBytesR(v)
				v = d.nextValueBytesR(v)
			}
		}
	case cborMajorTag:
		bs, _ = d.uintBytes()
		h.appendN(&v, d.bytes, bs...)
		v = d.nextValueBytesR(v)
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined, cborBdFalse, cborBdTrue:
		case cborBdFloat16:
			h.appendN(&v, d.bytes, d.r.readx(2)...)
		case cborBdFloat32:
			h.appendN(&v, d.bytes, d.r.readx(4)...)
		case cborBdFloat64:
			h.appendN(&v, d.bytes, d.r.readx(8)...)
		default:
			halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
	}
	return
}

func (d *cborDecDriverBytes) reset() {
	d.bdAndBdread.reset()
	d.st = d.h.SkipUnexpectedTags
}

func (d *cborEncDriverBytes) init(hh Handle, shared *encoderShared, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*CborHandle)
	d.e = shared
	if shared.bytes {
		fp = cborFpEncBytes
	} else {
		fp = cborFpEncIO
	}

	d.init2(enc)
	return
}

func (e *cborEncDriverBytes) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *cborEncDriverBytes) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *cborEncDriverBytes) writerEnd()                        { e.w.end() }

func (e *cborEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *cborEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *cborDecDriverBytes) init(hh Handle, shared *decoderShared, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*CborHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = cborFpDecBytes
	} else {
		fp = cborFpDecIO
	}

	d.init2(dec)
	return
}

func (d *cborDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *cborDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *cborDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

func (d *cborDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, cbordesc(d.bd))
}

func (d *cborDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

func (d *cborEncDriverBytes) init2(enc encoderI) {
	d.enc = enc
}

func (d *cborDecDriverBytes) init2(dec decoderI) {
	d.dec = dec
	d.d.cbor = true
}
func (e *encoderCborIO) setContainerState(cs containerState) {
	if cs != 0 {
		e.c = cs
	}
}

func (e *encoderCborIO) rawExt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *encoderCborIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderCborIO) selferMarshal(f *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderCborIO) binaryMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderCborIO) textMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderCborIO) jsonMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderCborIO) raw(f *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderCborIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderCborIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderCborIO) kBool(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderCborIO) kTime(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderCborIO) kString(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderCborIO) kFloat32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderCborIO) kFloat64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderCborIO) kComplex64(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderCborIO) kComplex128(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderCborIO) kInt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderCborIO) kInt8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderCborIO) kInt16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderCborIO) kInt32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderCborIO) kInt64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderCborIO) kUint(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderCborIO) kUint8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderCborIO) kUint16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderCborIO) kUint32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderCborIO) kUint64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderCborIO) kUintptr(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderCborIO) kErr(f *encFnInfo, rv reflect.Value) {
	e.errorf("unsupported encoding kind %s, for %#v", rv.Kind(), any(rv))
}

func (e *encoderCborIO) kSeqFn(rtelem reflect.Type) (fn *encFnCborIO) {
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}

	if rtelem.Kind() != reflect.Interface {
		fn = e.fn(rtelem)
	}
	return
}

func (e *encoderCborIO) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	if l == 0 {
		e.mapStart(0)
	} else {
		e.haltOnMbsOddLen(l)
		e.mapStart(l >> 1)
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			if j&1 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
	e.mapEnd()
}

func (e *encoderCborIO) kSliceW(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	e.arrayStart(l)
	if l > 0 {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
	e.arrayEnd()
}

func (e *encoderCborIO) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	if l == 0 {
		e.mapStart(0)
	} else {
		e.haltOnMbsOddLen(l)
		e.mapStart(l >> 1)
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			if j&1 == 0 {
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encodeValue(rv.Index(j), fn)
		}
	}
	e.mapEnd()
}

func (e *encoderCborIO) kArrayW(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	e.arrayStart(l)
	if l > 0 {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rv.Index(j), fn)
		}
	}
	e.arrayEnd()
}

func (e *encoderCborIO) kChan(f *encFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		e.errorStr("send-only channel cannot be encoded")
	}
	if !f.ti.mbs && uint8TypId == rt2id(f.ti.elem) {
		e.kSliceBytesChan(rv)
		return
	}
	rtslice := reflect.SliceOf(f.ti.elem)
	rv = chanToSlice(rv, rtslice, e.h.ChanRecvTimeout)
	ti := e.h.getTypeInfo(rt2id(rtslice), rtslice)
	if f.ti.mbs {
		e.kSliceWMbs(rv, ti)
	} else {
		e.kSliceW(rv, ti)
	}
}

func (e *encoderCborIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoderCborIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoderCborIO) kSliceBytesChan(rv reflect.Value) {

	bs0 := e.blist.peek(32, true)
	bs := bs0

	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok {
		ch = irv.(chan byte)
	}

L1:
	switch timeout := e.h.ChanRecvTimeout; {
	case timeout == 0:
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			}
		}
	case timeout > 0:
		tt := time.NewTimer(timeout)
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:

				break L1
			}
		}
	default:
		for b := range ch {
			bs = append(bs, b)
		}
	}

	e.e.EncodeStringBytesRaw(bs)
	e.blist.put(bs)
	if !byteSliceSameData(bs0, bs) {
		e.blist.put(bs0)
	}
}

func (e *encoderCborIO) kStructSfi(f *encFnInfo) []*structFieldInfo {
	if e.h.Canonical {
		return f.ti.sfi.sorted()
	}
	return f.ti.sfi.source()
}

func (e *encoderCborIO) kStructNoOmitempty(f *encFnInfo, rv reflect.Value) {
	var tisfi []*structFieldInfo
	if f.ti.toArray || e.h.StructToArray {
		tisfi = f.ti.sfi.source()
		e.arrayStart(len(tisfi))
		for _, si := range tisfi {
			e.arrayElem()
			e.encodeValue(si.path.field(rv), nil)
		}
		e.arrayEnd()
	} else {
		tisfi = e.kStructSfi(f)
		e.mapStart(len(tisfi))
		keytyp := f.ti.keyType
		for _, si := range tisfi {
			e.mapElemKey()
			e.kStructFieldKey(keytyp, si.path.encNameAsciiAlphaNum, si.encName)
			e.mapElemValue()
			e.encodeValue(si.path.field(rv), nil)
		}
		e.mapEnd()
	}
}

func (e *encoderCborIO) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
	if keyType == valueTypeString {
		if e.js && encNameAsciiAlphaNum {
			e.e.writeStringAsisDblQuoted(encName)
		} else {
			e.e.EncodeString(encName)
		}
	} else if keyType == valueTypeInt {
		e.e.EncodeInt(must.Int(strconv.ParseInt(encName, 10, 64)))
	} else if keyType == valueTypeUint {
		e.e.EncodeUint(must.Uint(strconv.ParseUint(encName, 10, 64)))
	} else if keyType == valueTypeFloat {
		e.e.EncodeFloat64(must.Float(strconv.ParseFloat(encName, 64)))
	} else {
		halt.errorStr2("invalid struct key type: ", keyType.String())
	}

}

func (e *encoderCborIO) kStruct(f *encFnInfo, rv reflect.Value) {
	var newlen int
	ti := f.ti
	toMap := !(ti.toArray || e.h.StructToArray)
	var mf map[string]interface{}
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	} else if ti.flagMissingFielderPtr {
		rv2 := e.addrRV(rv, ti.rt, ti.ptr)
		mf = rv2i(rv2).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	}
	tisfi := ti.sfi.source()
	newlen += len(tisfi)

	var fkvs = e.slist.get(newlen)[:newlen]

	recur := e.h.RecursiveEmptyCheck

	var kv sfiRv
	var j int
	if toMap {
		newlen = 0
		for _, si := range e.kStructSfi(f) {
			kv.r = si.path.field(rv)
			if si.path.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
				continue
			}
			kv.v = si
			fkvs[newlen] = kv
			newlen++
		}

		var mf2s []stringIntf
		if len(mf) > 0 {
			mf2s = make([]stringIntf, 0, len(mf))
			for k, v := range mf {
				if k == "" {
					continue
				}
				if ti.infoFieldOmitempty && isEmptyValue(reflect.ValueOf(v), e.h.TypeInfos, recur) {
					continue
				}
				mf2s = append(mf2s, stringIntf{k, v})
			}
		}

		e.mapStart(newlen + len(mf2s))

		if len(mf2s) > 0 && e.h.Canonical {
			mf2w := make([]encStructFieldObj, newlen+len(mf2s))
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				mf2w[j] = encStructFieldObj{kv.v.encName, kv.r, nil, kv.v.path.encNameAsciiAlphaNum, true}
			}
			for _, v := range mf2s {
				mf2w[j] = encStructFieldObj{v.v, reflect.Value{}, v.i, false, false}
				j++
			}
			sort.Sort((encStructFieldObjSlice)(mf2w))
			for _, v := range mf2w {
				e.mapElemKey()
				e.kStructFieldKey(ti.keyType, v.ascii, v.key)
				e.mapElemValue()
				if v.isRv {
					e.encodeValue(v.rv, nil)
				} else {
					e.encode(v.intf)
				}
			}
		} else {
			keytyp := ti.keyType
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				e.mapElemKey()
				e.kStructFieldKey(keytyp, kv.v.path.encNameAsciiAlphaNum, kv.v.encName)
				e.mapElemValue()
				e.encodeValue(kv.r, nil)
			}
			for _, v := range mf2s {
				e.mapElemKey()
				e.kStructFieldKey(keytyp, false, v.v)
				e.mapElemValue()
				e.encode(v.i)
			}
		}

		e.mapEnd()
	} else {
		newlen = len(tisfi)
		for i, si := range tisfi {
			kv.r = si.path.field(rv)

			if si.path.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
				switch kv.r.Kind() {
				case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
					kv.r = reflect.Value{}
				}
			}
			fkvs[i] = kv
		}

		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ {
			e.arrayElem()
			e.encodeValue(fkvs[j].r, nil)
		}
		e.arrayEnd()
	}

	e.slist.put(fkvs)
}

func (e *encoderCborIO) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	var keyFn, valFn *encFnCborIO

	ktypeKind := reflect.Kind(f.ti.keykind)
	vtypeKind := reflect.Kind(f.ti.elemkind)

	rtval := f.ti.elem
	rtvalkind := vtypeKind
	for rtvalkind == reflect.Ptr {
		rtval = rtval.Elem()
		rtvalkind = rtval.Kind()
	}
	if rtvalkind != reflect.Interface {
		valFn = e.fn(rtval)
	}

	var rvv = mapAddrLoopvarRV(f.ti.elem, vtypeKind)

	rtkey := f.ti.key
	var keyTypeIsString = stringTypId == rt2id(rtkey)
	if keyTypeIsString {
		keyFn = e.fn(rtkey)
	} else {
		for rtkey.Kind() == reflect.Ptr {
			rtkey = rtkey.Elem()
		}
		if rtkey.Kind() != reflect.Interface {
			keyFn = e.fn(rtkey)
		}
	}

	if e.h.Canonical {
		e.kMapCanonical(f.ti, rv, rvv, keyFn, valFn)
		e.mapEnd()
		return
	}

	var rvk = mapAddrLoopvarRV(f.ti.key, ktypeKind)

	var it mapIter
	mapRange(&it, rv, rvk, rvv, true)

	for it.Next() {
		e.mapElemKey()
		if keyTypeIsString {
			e.e.EncodeString(it.Key().String())
		} else {
			e.encodeValue(it.Key(), keyFn)
		}
		e.mapElemValue()
		e.encodeValue(it.Value(), valFn)
	}
	it.Done()

	e.mapEnd()
}

func (e *encoderCborIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnCborIO) {

	rtkey := ti.key
	rtkeydecl := rtkey.PkgPath() == "" && rtkey.Name() != ""

	mks := rv.MapKeys()
	rtkeyKind := rtkey.Kind()
	kfast := mapKeyFastKindFor(rtkeyKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	switch rtkeyKind {
	case reflect.Bool:

		if len(mks) == 2 && mks[0].Bool() {
			mks[0], mks[1] = mks[1], mks[0]
		}
		for i := range mks {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeBool(mks[i].Bool())
			} else {
				e.encodeValueNonNil(mks[i], keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mks[i], rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.String:
		mksv := make([]orderedRv[string], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.String()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeString(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		mksv := make([]orderedRv[uint64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Uint()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeUint(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		mksv := make([]orderedRv[int64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Int()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeInt(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Float32:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeFloat32(float32(mksv[i].v))
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Float64:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeFloat64(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	default:
		if rtkey == timeTyp {
			mksv := make([]timeRv, len(mks))
			for i, k := range mks {
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			}
			slices.SortFunc(mksv, cmpTimeRv)
			for i := range mksv {
				e.mapElemKey()
				e.e.EncodeTime(mksv[i].v)
				e.mapElemValue()
				e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
			}
			break
		}

		bs0 := e.blist.get(len(mks) * 16)
		mksv := bs0
		mksbv := make([]bytesRv, len(mks))

		func() {
			se := e.h.sideEncPool.Get().(encoderI)
			defer e.h.sideEncPool.Put(se)
			se.ResetBytes(&mksv)
			for i, k := range mks {
				v := &mksbv[i]
				l := len(mksv)
				se.setContainerState(containerMapKey)
				se.encode(rv2i(baseRVRV(k)))
				se.atEndOfEncode()
				se.writerEnd()
				v.r = k
				v.v = mksv[l:]
			}
		}()

		slices.SortFunc(mksbv, cmpBytesRv)
		for j := range mksbv {
			e.mapElemKey()
			e.e.writeBytesAsis(mksbv[j].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksbv[j].r, rvv, kfast, visindirect, visref), valFn)
		}
		e.blist.put(mksv)
		if !byteSliceSameData(bs0, mksv) {
			e.blist.put(bs0)
		}
	}
}

type encoderCborIO struct {
	panicHdl
	perType encPerType

	dh helperEncDriverCborIO

	fp *fastpathEsCborIO

	h *BasicHandle

	e cborEncDriverIO

	encoderShared

	hh Handle

	ci []interface{}

	slist sfiRvFreelist
}

func (e *encoderCborIO) HandleName() string {
	return e.hh.Name()
}

func (e *encoderCborIO) Release() {
}

func (e *encoderCborIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderShared, e).(*fastpathEsCborIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderCborIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderCborIO) MustEncode(v interface{}) {
	halt.onerror(e.err)
	if e.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	e.calls++
	e.encode(v)
	e.calls--
	if e.calls == 0 {
		e.e.atEndOfEncode()
		e.e.writerEnd()
	}
}

func (e *encoderCborIO) Encode(v interface{}) (err error) {

	if !debugging {
		defer func() {

			if x := recover(); x != nil {
				panicValToErr(e, x, &e.err)
				err = e.err
			}
		}()
	}

	e.MustEncode(v)
	return
}

func (e *encoderCborIO) encode(iv interface{}) {

	if iv == nil {
		e.e.EncodeNil()
		return
	}

	rv, ok := isNil(iv)
	if ok {
		e.e.EncodeNil()
		return
	}

	switch v := iv.(type) {

	case Raw:
		e.rawBytes(v)
	case reflect.Value:
		e.encodeValue(v, nil)

	case string:
		e.e.EncodeString(v)
	case bool:
		e.e.EncodeBool(v)
	case int:
		e.e.EncodeInt(int64(v))
	case int8:
		e.e.EncodeInt(int64(v))
	case int16:
		e.e.EncodeInt(int64(v))
	case int32:
		e.e.EncodeInt(int64(v))
	case int64:
		e.e.EncodeInt(v)
	case uint:
		e.e.EncodeUint(uint64(v))
	case uint8:
		e.e.EncodeUint(uint64(v))
	case uint16:
		e.e.EncodeUint(uint64(v))
	case uint32:
		e.e.EncodeUint(uint64(v))
	case uint64:
		e.e.EncodeUint(v)
	case uintptr:
		e.e.EncodeUint(uint64(v))
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)
	case complex64:
		e.encodeComplex64(v)
	case complex128:
		e.encodeComplex128(v)
	case time.Time:
		e.e.EncodeTime(v)
	case []byte:
		e.e.EncodeStringBytesRaw(v)
	case *Raw:
		e.rawBytes(*v)
	case *string:
		e.e.EncodeString(*v)
	case *bool:
		e.e.EncodeBool(*v)
	case *int:
		e.e.EncodeInt(int64(*v))
	case *int8:
		e.e.EncodeInt(int64(*v))
	case *int16:
		e.e.EncodeInt(int64(*v))
	case *int32:
		e.e.EncodeInt(int64(*v))
	case *int64:
		e.e.EncodeInt(*v)
	case *uint:
		e.e.EncodeUint(uint64(*v))
	case *uint8:
		e.e.EncodeUint(uint64(*v))
	case *uint16:
		e.e.EncodeUint(uint64(*v))
	case *uint32:
		e.e.EncodeUint(uint64(*v))
	case *uint64:
		e.e.EncodeUint(*v)
	case *uintptr:
		e.e.EncodeUint(uint64(*v))
	case *float32:
		e.e.EncodeFloat32(*v)
	case *float64:
		e.e.EncodeFloat64(*v)
	case *complex64:
		e.encodeComplex64(*v)
	case *complex128:
		e.encodeComplex128(*v)
	case *time.Time:
		e.e.EncodeTime(*v)
	case *[]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			e.e.EncodeStringBytesRaw(*v)
		}
	default:

		if skipFastpathTypeSwitchInDirectCall || !e.dh.fastpathEncodeTypeSwitch(iv, e) {
			e.encodeValue(rv, nil)
		}
	}
}

func (e *encoderCborIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderCborIO) encodeValue(rv reflect.Value, fn *encFnCborIO) {

	var sptr interface{}
	var rvp reflect.Value
	var rvpValid bool
TOP:
	switch rv.Kind() {
	case reflect.Ptr:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
		rvpValid = true
		rvp = rv
		rv = rv.Elem()
		goto TOP
	case reflect.Interface:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
		rvpValid = false
		rvp = reflect.Value{}
		rv = rv.Elem()
		goto TOP
	case reflect.Struct:
		if rvpValid && e.h.CheckCircularRef {
			sptr = rv2i(rvp)
			for _, vv := range e.ci {
				if eq4i(sptr, vv) {
					e.errorf("circular reference found: %p, %T", sptr, sptr)
				}
			}
			e.ci = append(e.ci, sptr)
		}
	case reflect.Slice, reflect.Map, reflect.Chan:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		return
	}

	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if !fn.i.addrE {

	} else if rvpValid {
		rv = rvp
	} else {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)

	if sptr != nil {
		e.ci = e.ci[:len(e.ci)-1]
	}
}

func (e *encoderCborIO) encodeValueNonNil(rv reflect.Value, fn *encFnCborIO) {
	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if fn.i.addrE {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderCborIO) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
	if rv.CanAddr() {
		return rvAddr(rv, ptrType)
	}
	if e.h.NoAddressableReadonly {
		rva = reflect.New(typ)
		rvSetDirect(rva.Elem(), rv)
		return
	}
	return rvAddr(e.perType.AddressableRO(rv), ptrType)
}

func (e *encoderCborIO) marshalUtf8(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderCborIO) marshalAsis(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderCborIO) marshalRaw(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoderCborIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		e.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderCborIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}

func (e *encoderCborIO) fn(t reflect.Type) *encFnCborIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoderCborIO) fnNoExt(t reflect.Type) *encFnCborIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

func (e *encoderCborIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderCborIO) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoderCborIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderCborIO) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoderCborIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderCborIO) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoderCborIO) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

func (e *encoderCborIO) haltOnMbsOddLen(length int) {
	if length&1 != 0 {
		e.errorInt("mapBySlice requires even slice length, but got ", int64(length))
	}
}

func (e *encoderCborIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderCborIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderCborIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderCborIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderCborIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverCborIO) newEncoderBytes(out *[]byte, h Handle) *encoderCborIO {
	var c1 encoderCborIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverCborIO) newEncoderIO(out io.Writer, h Handle) *encoderCborIO {
	var c1 encoderCborIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverCborIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsCborIO) (f *fastpathECborIO, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

type encFnCborIO struct {
	i  encFnInfo
	fe func(*encoderCborIO, *encFnInfo, reflect.Value)
}
type encRtidFnCborIO struct {
	rtid uintptr
	fn   *encFnCborIO
}

func (helperEncDriverCborIO) encFindRtidFn(s []encRtidFnCborIO, rtid uintptr) (i uint, fn *encFnCborIO) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperEncDriverCborIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnCborIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnCborIO](v))
	}
	return
}

func (dh helperEncDriverCborIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsCborIO, checkExt bool) (fn *encFnCborIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverCborIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnCborIO
	sp = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverCborIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnCborIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnCborIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnCborIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnCborIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverCborIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborIO) {
	fn = new(encFnCborIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderCborIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderCborIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderCborIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderCborIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fe = (*encoderCborIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderCborIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderCborIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderCborIO).textMarshal
		fi.addrE = ti.flagTextMarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {

			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fe = fp[idx].encfn
				}
			} else {

				xfe, xrt := dh.encFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf := xfe.encfn
					fn.fe = func(e *encoderCborIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderCborIO).kBool
			case reflect.String:

				fn.fe = (*encoderCborIO).kString
			case reflect.Int:
				fn.fe = (*encoderCborIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderCborIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderCborIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderCborIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderCborIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderCborIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderCborIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderCborIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderCborIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderCborIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderCborIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderCborIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderCborIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderCborIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderCborIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderCborIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderCborIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderCborIO).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*encoderCborIO).kStruct
				} else {
					fn.fe = (*encoderCborIO).kStructNoOmitempty
				}
			case reflect.Map:
				fn.fe = (*encoderCborIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderCborIO).kErr
			default:

				fn.fe = (*encoderCborIO).kErr
			}
		}
	}
	return
}

type helperEncDriverCborIO struct{}

func (d *decoderCborIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *decoderCborIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderCborIO) selferUnmarshal(f *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderCborIO) binaryUnmarshal(f *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	d.onerror(fnerr)
}

func (d *decoderCborIO) textUnmarshal(f *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	d.onerror(fnerr)
}

func (d *decoderCborIO) jsonUnmarshal(f *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderCborIO) jsonUnmarshalV(tm jsonUnmarshaler) {

	var bs0 = zeroByteSlice
	if !d.bytes {
		bs0 = d.blist.get(256)
	}
	bs := d.d.nextValueBytes(bs0)
	fnerr := tm.UnmarshalJSON(bs)
	if !d.bytes {
		d.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			d.blist.put(bs0)
		}
	}
	d.onerror(fnerr)
}

func (d *decoderCborIO) kErr(f *decFnInfo, rv reflect.Value) {
	d.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderCborIO) raw(f *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderCborIO) kString(f *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes()))
}

func (d *decoderCborIO) kBool(f *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderCborIO) kTime(f *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderCborIO) kFloat32(f *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderCborIO) kFloat64(f *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderCborIO) kComplex64(f *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderCborIO) kComplex128(f *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderCborIO) kInt(f *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderCborIO) kInt8(f *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderCborIO) kInt16(f *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderCborIO) kInt32(f *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderCborIO) kInt64(f *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderCborIO) kUint(f *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborIO) kUintptr(f *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborIO) kUint8(f *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderCborIO) kUint16(f *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderCborIO) kUint32(f *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderCborIO) kUint64(f *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderCborIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

	n := d.naked()
	d.d.DecodeNaked()

	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
	}
	switch n.v {
	case valueTypeMap:
		mtid := d.mtid
		if mtid == 0 {
			if d.jsms {
				mtid = mapStrIntfTypId
			} else {
				mtid = mapIntfIntfTypId
			}
		}
		if mtid == mapStrIntfTypId {
			var v2 map[string]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if mtid == mapIntfIntfTypId {
			var v2 map[interface{}]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.mtr {
			rvn = reflect.New(d.h.MapType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.MapType, reflect.Map)
			d.decodeValue(rvn, nil)
		}
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId {
			var v2 []interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.str {
			rvn = reflect.New(d.h.SliceType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.SliceType, reflect.Slice)
			d.decodeValue(rvn, nil)
		}
		if d.h.PreferArrayOverSlice {
			rvn = rvGetArray4Slice(rvn)
		}
	case valueTypeExt:
		tag, bytes := n.u, n.l
		bfn := d.h.getExtForTag(tag)
		var re = RawExt{Tag: tag}
		if bytes == nil {

			if bfn == nil {
				d.decode(&re.Value)
				rvn = rv4iptr(&re).Elem()
			} else {
				if bfn.ext == SelfExt {
					rvn = rvZeroAddrK(bfn.rt, bfn.rt.Kind())
					d.decodeValue(rvn, d.fnNoExt(bfn.rt))
				} else {
					rvn = reflect.New(bfn.rt)
					d.interfaceExtConvertAndDecode(rv2i(rvn), bfn.ext)
					rvn = rvn.Elem()
				}
			}
		} else {

			if bfn == nil {
				re.setData(bytes, false)
				rvn = rv4iptr(&re).Elem()
			} else {
				rvn = reflect.New(bfn.rt)
				if bfn.ext == SelfExt {
					sideDecode(d.h, rv2i(rvn), bytes, bfn.rt, true)
				} else {
					bfn.ext.ReadExt(rv2i(rvn), bytes)
				}
				rvn = rvn.Elem()
			}
		}

		if d.h.PreferPointerForStructOrArray && rvn.CanAddr() {
			if rk := rvn.Kind(); rk == reflect.Array || rk == reflect.Struct {
				rvn = rvn.Addr()
			}
		}
	case valueTypeNil:

	case valueTypeInt:
		rvn = n.ri()
	case valueTypeUint:
		rvn = n.ru()
	case valueTypeFloat:
		rvn = n.rf()
	case valueTypeBool:
		rvn = n.rb()
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs()
	case valueTypeBytes:
		rvn = n.rl()
	case valueTypeTime:
		rvn = n.rt()
	default:
		halt.errorStr2("kInterfaceNaked: unexpected valueType: ", n.v.String())
	}
	return
}

func (d *decoderCborIO) kInterface(f *decFnInfo, rv reflect.Value) {

	isnilrv := rvIsNil(rv)

	var rvn reflect.Value

	if d.h.InterfaceReset {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			} else if !isnilrv {
				decSetNonNilRV2Zero4Intf(rv)
			}
			return
		}
	} else if isnilrv {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			}
			return
		}
	} else {

		rvn = rv.Elem()
	}

	canDecode, _ := isDecodeable(rvn)

	if !canDecode {
		rvn2 := d.oneShotAddrRV(rvn.Type(), rvn.Kind())
		rvSetDirect(rvn2, rvn)
		rvn = rvn2
	}

	d.decodeValue(rvn, nil)
	rvSetIntf(rv, rvn)
}

func (d *decoderCborIO) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *decoderCborIO) kStruct(f *decFnInfo, rv reflect.Value) {
	ctyp := d.d.ContainerType()
	ti := f.ti
	var mf MissingFielder
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder)
	} else if ti.flagMissingFielderPtr {
		mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder)
	}
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var name2 []byte
		if mf != nil {
			var namearr2 [16]byte
			name2 = namearr2[:0]
		}
		var rvkencname []byte
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey()
			switch ti.keyType {
			case valueTypeString:
				rvkencname = d.d.DecodeStringAsBytes()
			case valueTypeInt:
				rvkencname = strconv.AppendInt(d.b[:0], d.d.DecodeInt64(), 10)
			case valueTypeUint:
				rvkencname = strconv.AppendUint(d.b[:0], d.d.DecodeUint64(), 10)
			case valueTypeFloat:
				rvkencname = strconv.AppendFloat(d.b[:0], d.d.DecodeFloat64(), 'f', -1, 64)
			default:
				halt.errorStr2("invalid struct key type: ", ti.keyType.String())
			}

			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else if mf != nil {

				name2 = append(name2[:0], rvkencname...)
				var f interface{}
				d.decode(&f)
				if !mf.CodecMissingField(name2, f) && d.h.ErrorIfNoField {
					d.errorStr2("no matching struct field when decoding stream map with key: ", stringView(name2))
				}
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}

		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem()
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}

		d.arrayEnd()
	} else {
		d.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderCborIO) kSlice(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		rvbs := rvGetBytes(rv)
		if !rvCanset {

			rvbs = rvbs[:len(rvbs):len(rvbs)]
		}
		bs2 := d.decodeBytesInto(rvbs)

		if !(len(bs2) > 0 && len(bs2) == len(rvbs) && &bs2[0] == &rvbs[0]) {
			if rvCanset {
				rvSetBytes(rv, bs2)
			} else if len(rvbs) > 0 && len(bs2) > 0 {
				copy(rvbs, bs2)
			}
		}
		return
	}

	slh, containerLenS := d.decSliceHelperStart()

	if containerLenS == 0 {
		if rvCanset {
			if rvIsNil(rv) {
				rvSetDirect(rv, rvSliceZeroCap(ti.rt))
			} else {
				rvSetSliceLen(rv, 0)
			}
		}
		slh.End()
		return
	}

	rtelem0Mut := !scalarBitset.isset(ti.elemkind)
	rtelem := ti.elem

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnCborIO

	var rvChanged bool

	var rv0 = rv
	var rv9 reflect.Value

	rvlen := rvLenSlice(rv)
	rvcap := rvCapSlice(rv)
	hasLen := containerLenS > 0
	if hasLen {
		if containerLenS > rvcap {
			oldRvlenGtZero := rvlen > 0
			rvlen1 := decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
			if rvlen1 == rvlen {
			} else if rvlen1 <= rvcap {
				if rvCanset {
					rvlen = rvlen1
					rvSetSliceLen(rv, rvlen)
				}
			} else if rvCanset {
				rvlen = rvlen1
				rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = !rvCanset
			} else {
				d.errorStr("cannot decode into non-settable slice")
			}
			if rvChanged && oldRvlenGtZero && rtelem0Mut {
				rvCopySlice(rv, rv0, rtelem)
			}
		} else if containerLenS != rvlen {
			if rvCanset {
				rvlen = containerLenS
				rvSetSliceLen(rv, rvlen)
			}
		}
	}

	var elemReset = d.h.SliceElementReset

	var j int

	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if rvCanset {
					rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
					rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
					rvcap = rvlen
					rvChanged = !rvCanset
				} else {
					d.errorStr("cannot decode into non-settable slice")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}

		if j >= rvlen {
			slh.ElemContainerState(j)

			if rvlen < rvcap {
				rvlen = rvcap
				if rvCanset {
					rvSetSliceLen(rv, rvlen)
				} else if rvChanged {
					rv = rvSlice(rv, rvlen)
				} else {
					d.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					d.onerror(errExpandSliceCannotChange)
				}
				rv, rvcap, rvCanset = rvGrowSlice(rv, f.ti, rvcap, 1)
				rvlen = rvcap
				rvChanged = !rvCanset
			}
		} else {
			slh.ElemContainerState(j)
		}
		rv9 = rvSliceIndex(rv, j, f.ti)
		if elemReset {
			rvSetZero(rv9)
		}
		d.decodeValue(rv9, fn)
	}
	if j < rvlen {
		if rvCanset {
			rvSetSliceLen(rv, j)
		} else if rvChanged {
			rv = rvSlice(rv, j)
		}

	} else if j == 0 && rvIsNil(rv) {
		if rvCanset {
			rv = rvSliceZeroCap(ti.rt)
			rvCanset = false
			rvChanged = true
		}
	}
	slh.End()

	if rvChanged {
		rvSetDirect(rv0, rv)
	}
}

func (d *decoderCborIO) kArray(f *decFnInfo, rv reflect.Value) {

	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {

		if f.ti.elemkind != uint8(reflect.Uint8) {
			d.errorf("bytes/string in stream can decode into array of bytes, but not %v", f.ti.rt)
		}
		rvbs := rvGetArrayBytes(rv, nil)
		bs2 := d.decodeBytesInto(rvbs)
		if !byteSliceSameData(rvbs, bs2) && len(rvbs) > 0 && len(bs2) > 0 {
			copy(rvbs, bs2)
		}
		return
	}

	slh, containerLenS := d.decSliceHelperStart()

	if containerLenS == 0 {
		slh.End()
		return
	}

	rtelem := f.ti.elem
	for k := reflect.Kind(f.ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnCborIO

	var rv9 reflect.Value

	rvlen := rv.Len()
	hasLen := containerLenS > 0
	if hasLen && containerLenS > rvlen {
		d.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
	}

	var elemReset = d.h.SliceElementReset

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {

		if j >= rvlen {
			slh.arrayCannotExpand(hasLen, rvlen, j, containerLenS)
			return
		}

		slh.ElemContainerState(j)
		rv9 = rvArrayIndex(rv, j, f.ti)
		if elemReset {
			rvSetZero(rv9)
		}

		if fn == nil {
			fn = d.fn(rtelem)
		}
		d.decodeValue(rv9, fn)
	}
	slh.End()
}

func (d *decoderCborIO) kChan(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		d.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2 := d.d.DecodeBytes(nil)
		irv := rv2i(rv)
		ch, ok := irv.(chan<- byte)
		if !ok {
			ch = irv.(chan byte)
		}
		for _, b := range bs2 {
			ch <- b
		}
		return
	}

	var rvCanset = rv.CanSet()

	slh, containerLenS := d.decSliceHelperStart()

	if containerLenS == 0 {
		if rvCanset && rvIsNil(rv) {
			rvSetDirect(rv, reflect.MakeChan(ti.rt, 0))
		}
		slh.End()
		return
	}

	rtelem := ti.elem
	useTransient := decUseTransient && ti.elemkind != byte(reflect.Ptr) && ti.tielem.flagCanTransient

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnCborIO

	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int
	hasLen := containerLenS > 0

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if hasLen {
					rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
				} else {
					rvlen = decDefChanCap
				}
				if rvCanset {
					rv = reflect.MakeChan(ti.rt, rvlen)
					rvChanged = true
				} else {
					d.errorStr("cannot decode into non-settable chan")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}
		slh.ElemContainerState(j)
		if rv9.IsValid() {
			rvSetZero(rv9)
		} else if decUseTransient && useTransient {
			rv9 = d.perType.TransientAddrK(ti.elem, reflect.Kind(ti.elemkind))
		} else {
			rv9 = rvZeroAddrK(ti.elem, reflect.Kind(ti.elemkind))
		}
		if !d.d.TryNil() {
			d.decodeValueNoCheckNil(rv9, fn)
		}
		rv.Send(rv9)
	}
	slh.End()

	if rvChanged {
		rvSetDirect(rv0, rv)
	}

}

func (d *decoderCborIO) kMap(f *decFnInfo, rv reflect.Value) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	ti := f.ti
	if rvIsNil(rv) {
		rvlen := decInferLen(containerLen, d.h.MaxInitLen, int(ti.keysize+ti.elemsize))
		rvSetDirect(rv, makeMapReflect(ti.rt, rvlen))
	}

	if containerLen == 0 {
		d.mapEnd()
		return
	}

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := reflect.Kind(ti.elemkind)
	ktypeKind := reflect.Kind(ti.keykind)
	kfast := mapKeyFastKindFor(ktypeKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	vtypePtr := vtypeKind == reflect.Ptr
	ktypePtr := ktypeKind == reflect.Ptr

	vTransient := decUseTransient && !vtypePtr && ti.tielem.flagCanTransient
	kTransient := decUseTransient && !ktypePtr && ti.tikey.flagCanTransient

	var vtypeElem reflect.Type

	var keyFn, valFn *decFnCborIO
	var ktypeLo, vtypeLo = ktype, vtype

	if ktypeKind == reflect.Ptr {
		for ktypeLo = ktype.Elem(); ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() {
		}
	}

	if vtypePtr {
		vtypeElem = vtype.Elem()
		for vtypeLo = vtypeElem; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() {
		}
	}

	rvkMut := !scalarBitset.isset(ti.keykind)
	rvvMut := !scalarBitset.isset(ti.elemkind)
	rvvCanNil := isnilBitset.isset(ti.elemkind)

	var rvk, rvkn, rvv, rvvn, rvva, rvvz reflect.Value

	var doMapGet, doMapSet bool

	if !d.h.MapValueReset {
		if rvvMut && (vtypeKind != reflect.Interface || !d.h.InterfaceReset) {
			doMapGet = true
			rvva = mapAddrLoopvarRV(vtype, vtypeKind)
		}
	}

	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId

	hasLen := containerLen > 0

	var kstrbs []byte
	var kstr2bs []byte
	var s string

	var callFnRvk bool

	fnRvk2 := func() (s string) {
		callFnRvk = false

		switch len(kstr2bs) {
		case 0:
		case 1:
			s = str256[kstr2bs[0] : kstr2bs[0]+1]
		default:
			s = d.mapKeyString(&callFnRvk, &kstrbs, &kstr2bs)
		}
		return
	}

	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		callFnRvk = false
		if j == 0 {

			if decUseTransient && vTransient && kTransient {
				rvk = d.perType.TransientAddr2K(ktype, ktypeKind)
			} else {
				rvk = rvZeroAddrK(ktype, ktypeKind)
			}
			if !rvkMut {
				rvkn = rvk
			}
			if !rvvMut {
				if decUseTransient && vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
			}
			if !ktypeIsString && keyFn == nil {
				keyFn = d.fn(ktypeLo)
			}
			if valFn == nil {
				valFn = d.fn(vtypeLo)
			}
		} else if rvkMut {
			rvSetZero(rvk)
		} else {
			rvk = rvkn
		}

		d.mapElemKey()
		if ktypeIsString {
			kstr2bs = d.d.DecodeStringAsBytes()
			rvSetString(rvk, fnRvk2())
		} else {
			d.decByteState = decByteStateNone
			d.decodeValue(rvk, keyFn)

			if ktypeIsIntf {
				if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp {
					kstr2bs = rvGetBytes(rvk2)
					rvSetIntf(rvk, rv4istr(fnRvk2()))
				}

			}
		}

		d.mapElemValue()

		if d.d.TryNil() {

			if !rvvz.IsValid() {
				rvvz = rvZeroK(vtype, vtypeKind)
			}
			if callFnRvk {
				s = d.string(kstr2bs)
				if ktypeIsString {
					rvSetString(rvk, s)
				} else {
					rvSetIntf(rvk, rv4istr(s))
				}
			}
			mapSet(rv, rvk, rvvz, kfast, visindirect, visref)
			continue
		}

		doMapSet = true

		if !rvvMut {
			rvv = rvvn
		} else if !doMapGet {
			goto NEW_RVV
		} else {
			rvv = mapGet(rv, rvk, rvva, kfast, visindirect, visref)
			if !rvv.IsValid() || (rvvCanNil && rvIsNil(rvv)) {
				goto NEW_RVV
			}
			switch vtypeKind {
			case reflect.Ptr, reflect.Map:
				doMapSet = false
			case reflect.Interface:

				rvvn = rvv.Elem()
				if k := rvvn.Kind(); (k == reflect.Ptr || k == reflect.Map) && !rvIsNil(rvvn) {
					d.decodeValueNoCheckNil(rvvn, nil)
					continue
				}

				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvSetIntf(rvvn, rvv)
				rvv = rvvn
			default:

				if decUseTransient && vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
				rvSetDirect(rvvn, rvv)
				rvv = rvvn
			}
		}
		goto DECODE_VALUE_NO_CHECK_NIL

	NEW_RVV:
		if vtypePtr {
			rvv = reflect.New(vtypeElem)
		} else if decUseTransient && vTransient {
			rvv = d.perType.TransientAddrK(vtype, vtypeKind)
		} else {
			rvv = rvZeroAddrK(vtype, vtypeKind)
		}

	DECODE_VALUE_NO_CHECK_NIL:
		d.decodeValueNoCheckNil(rvv, valFn)

		if doMapSet {
			if callFnRvk {
				s = d.string(kstr2bs)
				if ktypeIsString {
					rvSetString(rvk, s)
				} else {
					rvSetIntf(rvk, rv4istr(s))
				}
			}
			mapSet(rv, rvk, rvv, kfast, visindirect, visref)
		}
	}

	d.mapEnd()
}

type decoderCborIO struct {
	panicHdl
	perType decPerType

	dh helperDecDriverCborIO

	fp *fastpathDsCborIO

	h *BasicHandle

	d cborDecDriverIO

	decoderShared

	hh Handle

	mtid uintptr
	stid uintptr
}

func (d *decoderCborIO) HandleName() string {
	return d.hh.Name()
}

func (d *decoderCborIO) isBytes() bool {
	return d.bytes
}

func (d *decoderCborIO) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()
	d.zeroCopy = d.h.ZeroCopy
	d.be = h.isBinary()
	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderShared, d).(*fastpathDsCborIO)

	d.cbreak = d.js || d.cbor

	if d.bytes {
		d.rtidFn = &d.h.rtidFnsDecBytes
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtBytes
	} else {
		d.rtidFn = &d.h.rtidFnsDecIO
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtIO
	}

	d.reset()

}

func (d *decoderCborIO) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
	d.decByteState = decByteStateNone
	d.depth = 0
	d.calls = 0

	d.maxdepth = decDefMaxDepth
	if d.h.MaxDepth > 0 {
		d.maxdepth = d.h.MaxDepth
	}
	d.mtid = 0
	d.stid = 0
	d.mtr = false
	d.str = false
	if d.h.MapType != nil {
		d.mtid = rt2id(d.h.MapType)
		_, d.mtr = fastpathAvIndex(d.mtid)
	}
	if d.h.SliceType != nil {
		d.stid = rt2id(d.h.SliceType)
		_, d.str = fastpathAvIndex(d.stid)
	}
}

func (d *decoderCborIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderCborIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
	return
}

func (d *decoderCborIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderCborIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderCborIO) Decode(v interface{}) (err error) {

	if !debugging {
		defer func() {
			if x := recover(); x != nil {
				panicValToErr(d, x, &d.err)
				err = d.err
			}
		}()
	}

	d.MustDecode(v)
	return
}

func (d *decoderCborIO) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderCborIO) Release() {
}

func (d *decoderCborIO) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *decoderCborIO) nextValueBytes(start []byte) []byte {
	return d.d.nextValueBytes(start)
}

func (d *decoderCborIO) decode(iv interface{}) {

	if iv == nil {
		d.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {

	case reflect.Value:
		if x, _ := isDecodeable(v); !x {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)
	case *string:
		*v = d.stringZC(d.d.DecodeStringAsBytes())
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *float32:
		*v = d.d.DecodeFloat32()
	case *float64:
		*v = d.d.DecodeFloat64()
	case *complex64:
		*v = complex(d.d.DecodeFloat32(), 0)
	case *complex128:
		*v = complex(d.d.DecodeFloat64(), 0)
	case *[]byte:
		*v = d.decodeBytesInto(*v)
	case []byte:

		b := d.decodeBytesInto(v[:len(v):len(v)])
		if !(len(b) > 0 && len(b) == len(v) && &b[0] == &v[0]) {
			copy(v, b)
		}
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface{}:
		d.decodeValue(rv4iptr(v), nil)

	default:

		if skipFastpathTypeSwitchInDirectCall || !d.dh.fastpathDecodeTypeSwitch(iv, d) {
			v := reflect.ValueOf(iv)
			if x, _ := isDecodeable(v); !x {
				d.haltAsNotDecodeable(v)
			}
			d.decodeValue(v, nil)
		}
	}
}

func (d *decoderCborIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderCborIO) decodeValue(rv reflect.Value, fn *decFnCborIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *decoderCborIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnCborIO) {

	var rvp reflect.Value
	var rvpValid bool
PTR:
	if rv.Kind() == reflect.Ptr {
		rvpValid = true
		if rvIsNil(rv) {
			rvSetDirect(rv, reflect.New(rv.Type().Elem()))
		}
		rvp = rv
		rv = rv.Elem()
		goto PTR
	}

	if fn == nil {
		fn = d.fn(rv.Type())
	}
	if fn.i.addrD {
		if rvpValid {
			rv = rvp
		} else if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else if fn.i.addrDf {
			d.errorStr("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *decoderCborIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			d.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			d.errorStr("no matching struct field found when decoding stream map with key " + rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderCborIO) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		d.errorf("cannot expand array len during decode from %v to %v", any(sliceLen), any(streamLen))
	}
}

func (d *decoderCborIO) haltAsNotDecodeable(rv reflect.Value) {
	if !rv.IsValid() {
		d.onerror(errCannotDecodeIntoNil)
	}

	if !rv.CanInterface() {
		d.errorf("cannot decode into a value without an interface: %v", rv)
	}
	d.errorf("cannot decode into value of kind: %v, %#v", rv.Kind(), rv2i(rv))
}

func (d *decoderCborIO) depthIncr() {
	d.depth++
	if d.depth >= d.maxdepth {
		d.onerror(errMaxDepthExceeded)
	}
}

func (d *decoderCborIO) depthDecr() {
	d.depth--
}

func (d *decoderCborIO) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = zeroByteSlice
	}
	return d.d.DecodeBytes(in)
}

func (d *decoderCborIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes(zeroByteSlice)
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderCborIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.NumBytesRead(), false)
}

func (d *decoderCborIO) NumBytesRead() int {
	return int(d.d.NumBytesRead())
}

func (d *decoderCborIO) checkBreak() (v bool) {
	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *decoderCborIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *decoderCborIO) mapStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerMapStart
	}
	return v
}

func (d *decoderCborIO) mapElemKey() {
	d.d.ReadMapElemKey()
	d.c = containerMapKey
}

func (d *decoderCborIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderCborIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborIO) arrayStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerArrayStart
	}
	return v
}

func (d *decoderCborIO) arrayElem() {
	d.d.ReadArrayElem()
	d.c = containerArrayElem
}

func (d *decoderCborIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var s interface{}
	rv := reflect.ValueOf(v)
	rv2 := rv.Elem()
	rvk := rv2.Kind()
	if rvk == reflect.Struct || rvk == reflect.Array {
		s = ext.ConvertExt(v)
	} else {
		s = ext.ConvertExt(rv2i(rv2))
	}
	rv = reflect.ValueOf(s)

	if !rv.CanAddr() {
		rvk = rv.Kind()
		rv2 = d.oneShotAddrRV(rv.Type(), rvk)
		if rvk == reflect.Interface {
			rvSetIntf(rv2, rv)
		} else {
			rvSetDirect(rv2, rv)
		}
		rv = rv2
	}

	d.decodeValue(rv, nil)
	ext.UpdateExt(v, rv2i(rv))
}

func (d *decoderCborIO) oneShotAddrRV(rvt reflect.Type, rvk reflect.Kind) reflect.Value {
	if decUseTransient &&
		(numBoolStrSliceBitset.isset(byte(rvk)) ||
			((rvk == reflect.Struct || rvk == reflect.Array) &&
				d.h.getTypeInfo(rt2id(rvt), rvt).flagCanTransient)) {
		return d.perType.TransientAddrK(rvt, rvk)
	}
	return rvZeroAddrK(rvt, rvk)
}

func (d *decoderCborIO) fn(t reflect.Type) *decFnCborIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, false)
}

func (d *decoderCborIO) fnNoExt(t reflect.Type) *decFnCborIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, true)
}

type decSliceHelperCborIO struct {
	d     *decoderCborIO
	ct    valueType
	Array bool
	IsNil bool
}

func (d *decoderCborIO) decSliceHelperStart() (x decSliceHelperCborIO, clen int) {
	x.ct = d.d.ContainerType()
	x.d = d
	switch x.ct {
	case valueTypeNil:
		x.IsNil = true
	case valueTypeArray:
		x.Array = true
		clen = d.arrayStart(d.d.ReadArrayStart())
	case valueTypeMap:
		clen = d.mapStart(d.d.ReadMapStart())
		clen += clen
	default:
		d.errorStr2("to decode into a slice, expect map/array - got ", x.ct.String())
	}
	return
}

func (x decSliceHelperCborIO) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelperCborIO) ElemContainerState(index int) {

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 {
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelperCborIO) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
	x.d.arrayCannotExpand(lenv, j+1)

	x.ElemContainerState(j)
	x.d.swallow()
	j++
	for ; x.d.containerNext(j, containerLenS, hasLen); j++ {
		x.ElemContainerState(j)
		x.d.swallow()
	}
	x.End()
}

func (helperDecDriverCborIO) newDecoderBytes(in []byte, h Handle) *decoderCborIO {
	var c1 decoderCborIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverCborIO) newDecoderIO(in io.Reader, h Handle) *decoderCborIO {
	var c1 decoderCborIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverCborIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsCborIO) (f *fastpathDCborIO, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

type decFnCborIO struct {
	i  decFnInfo
	fd func(*decoderCborIO, *decFnInfo, reflect.Value)
}
type decRtidFnCborIO struct {
	rtid uintptr
	fn   *decFnCborIO
}

func (helperDecDriverCborIO) decFindRtidFn(s []decRtidFnCborIO, rtid uintptr) (i uint, fn *decFnCborIO) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperDecDriverCborIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnCborIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnCborIO](v))
	}
	return
}

func (dh helperDecDriverCborIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsCborIO,
	checkExt bool) (fn *decFnCborIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverCborIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnCborIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverCborIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnCborIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnCborIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnCborIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnCborIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverCborIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborIO) {
	fn = new(decFnCborIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderCborIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderCborIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderCborIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderCborIO).ext
		fi.addrD = true
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fd = (*decoderCborIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderCborIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderCborIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderCborIO).textUnmarshal
		fi.addrD = ti.flagTextUnmarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fd = fp[idx].decfn
					fi.addrD = true
					fi.addrDf = false
					if rk == reflect.Array {
						fi.addrD = false
					}
				}
			} else {

				xfe, xrt := dh.decFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf2 := xfe.decfn
					if rk == reflect.Array {
						fi.addrD = false
						fn.fd = func(d *decoderCborIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderCborIO, xf *decFnInfo, xrv reflect.Value) {
							if xrv.Kind() == reflect.Ptr {
								xfnf2(d, xf, rvConvert(xrv, xptr2rt))
							} else {
								xfnf2(d, xf, rvConvert(xrv, xrt))
							}
						}
					}
				}
			}
		}
		if fn.fd == nil {
			switch rk {
			case reflect.Bool:
				fn.fd = (*decoderCborIO).kBool
			case reflect.String:
				fn.fd = (*decoderCborIO).kString
			case reflect.Int:
				fn.fd = (*decoderCborIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderCborIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderCborIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderCborIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderCborIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderCborIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderCborIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderCborIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderCborIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderCborIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderCborIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderCborIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderCborIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderCborIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderCborIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderCborIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderCborIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderCborIO).kArray
			case reflect.Struct:
				fn.fd = (*decoderCborIO).kStruct
			case reflect.Map:
				fn.fd = (*decoderCborIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderCborIO).kInterface
			default:

				fn.fd = (*decoderCborIO).kErr
			}
		}
	}
	return
}

type helperDecDriverCborIO struct{}
type fastpathECborIO struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderCborIO, *encFnInfo, reflect.Value)
}
type fastpathDCborIO struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderCborIO, *decFnInfo, reflect.Value)
}
type fastpathEsCborIO [56]fastpathECborIO
type fastpathDsCborIO [56]fastpathDCborIO
type fastpathETCborIO struct{}
type fastpathDTCborIO struct{}

func (helperEncDriverCborIO) fastpathEList() *fastpathEsCborIO {
	var i uint = 0
	var s fastpathEsCborIO
	fn := func(v interface{}, fe func(*encoderCborIO, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathECborIO{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderCborIO).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderCborIO).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderCborIO).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderCborIO).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderCborIO).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderCborIO).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderCborIO).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderCborIO).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderCborIO).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderCborIO).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderCborIO).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderCborIO).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderCborIO).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderCborIO).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderCborIO).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderCborIO).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderCborIO).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderCborIO).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderCborIO).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderCborIO).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderCborIO).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderCborIO).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderCborIO).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderCborIO).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderCborIO).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderCborIO).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderCborIO).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderCborIO).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderCborIO).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderCborIO).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderCborIO).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderCborIO).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderCborIO).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderCborIO).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderCborIO).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderCborIO).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderCborIO).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderCborIO).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderCborIO).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderCborIO).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderCborIO).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderCborIO).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderCborIO).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderCborIO).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderCborIO).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderCborIO).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderCborIO).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderCborIO).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderCborIO).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderCborIO).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderCborIO).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderCborIO).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderCborIO).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderCborIO).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderCborIO).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderCborIO).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverCborIO) fastpathDList() *fastpathDsCborIO {
	var i uint = 0
	var s fastpathDsCborIO
	fn := func(v interface{}, fd func(*decoderCborIO, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDCborIO{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderCborIO).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderCborIO).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderCborIO).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderCborIO).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderCborIO).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderCborIO).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderCborIO).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderCborIO).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderCborIO).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderCborIO).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderCborIO).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderCborIO).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderCborIO).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderCborIO).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderCborIO).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderCborIO).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderCborIO).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderCborIO).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderCborIO).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderCborIO).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderCborIO).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderCborIO).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderCborIO).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderCborIO).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderCborIO).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderCborIO).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderCborIO).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderCborIO).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderCborIO).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderCborIO).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderCborIO).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderCborIO).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderCborIO).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderCborIO).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderCborIO).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderCborIO).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderCborIO).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderCborIO).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderCborIO).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderCborIO).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderCborIO).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderCborIO).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderCborIO).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderCborIO).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderCborIO).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderCborIO).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderCborIO).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderCborIO).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderCborIO).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderCborIO).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderCborIO).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderCborIO).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderCborIO).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderCborIO).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderCborIO).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderCborIO).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverCborIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderCborIO) bool {
	var ft fastpathETCborIO
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

func (e *encoderCborIO) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceIntfV(v []interface{}, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceIntfV(v []interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceStringV(v []string, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceStringV(v []string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceBytesV(v [][]byte, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceBytesV(v [][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceFloat32V(v []float32, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceFloat32V(v []float32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceFloat64V(v []float64, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceFloat64V(v []float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceUint8V(v []uint8, e *encoderCborIO) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETCborIO) EncAsMapSliceUint8V(v []uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceUint64V(v []uint64, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceUint64V(v []uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceIntV(v []int, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceIntV(v []int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceInt32V(v []int32, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceInt32V(v []int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceInt64V(v []int64, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceInt64V(v []int64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceBoolV(v []bool, e *encoderCborIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathETCborIO) EncAsMapSliceBoolV(v []bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETCborIO) EncMapStringIntfV(v map[string]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETCborIO) EncMapStringStringV(v map[string]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETCborIO) EncMapStringBytesV(v map[string][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETCborIO) EncMapStringUint8V(v map[string]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETCborIO) EncMapStringUint64V(v map[string]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETCborIO) EncMapStringIntV(v map[string]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETCborIO) EncMapStringInt32V(v map[string]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETCborIO) EncMapStringFloat64V(v map[string]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETCborIO) EncMapStringBoolV(v map[string]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETCborIO) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETCborIO) EncMapUint8StringV(v map[uint8]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETCborIO) EncMapUint8BytesV(v map[uint8][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETCborIO) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETCborIO) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETCborIO) EncMapUint8IntV(v map[uint8]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETCborIO) EncMapUint8Int32V(v map[uint8]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETCborIO) EncMapUint8Float64V(v map[uint8]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETCborIO) EncMapUint8BoolV(v map[uint8]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETCborIO) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETCborIO) EncMapUint64StringV(v map[uint64]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETCborIO) EncMapUint64BytesV(v map[uint64][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETCborIO) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETCborIO) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETCborIO) EncMapUint64IntV(v map[uint64]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETCborIO) EncMapUint64Int32V(v map[uint64]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETCborIO) EncMapUint64Float64V(v map[uint64]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETCborIO) EncMapUint64BoolV(v map[uint64]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETCborIO) EncMapIntIntfV(v map[int]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETCborIO) EncMapIntStringV(v map[int]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETCborIO) EncMapIntBytesV(v map[int][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETCborIO) EncMapIntUint8V(v map[int]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETCborIO) EncMapIntUint64V(v map[int]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETCborIO) EncMapIntIntV(v map[int]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETCborIO) EncMapIntInt32V(v map[int]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETCborIO) EncMapIntFloat64V(v map[int]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETCborIO) EncMapIntBoolV(v map[int]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETCborIO) EncMapInt32IntfV(v map[int32]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETCborIO) EncMapInt32StringV(v map[int32]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETCborIO) EncMapInt32BytesV(v map[int32][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETCborIO) EncMapInt32Uint8V(v map[int32]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETCborIO) EncMapInt32Uint64V(v map[int32]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETCborIO) EncMapInt32IntV(v map[int32]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETCborIO) EncMapInt32Int32V(v map[int32]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETCborIO) EncMapInt32Float64V(v map[int32]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
	ft.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETCborIO) EncMapInt32BoolV(v map[int32]bool, e *encoderCborIO) {
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

func (helperDecDriverCborIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderCborIO) bool {
	var ft fastpathDTCborIO
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

func (d *decoderCborIO) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceIntfX(vp *[]interface{}, d *decoderCborIO) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceIntfY(v []interface{}, d *decoderCborIO) (v2 []interface{}, changed bool) {
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
func (fastpathDTCborIO) DecSliceIntfN(v []interface{}, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceStringX(vp *[]string, d *decoderCborIO) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceStringY(v []string, d *decoderCborIO) (v2 []string, changed bool) {
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
func (fastpathDTCborIO) DecSliceStringN(v []string, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceBytesX(vp *[][]byte, d *decoderCborIO) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceBytesY(v [][]byte, d *decoderCborIO) (v2 [][]byte, changed bool) {
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
func (fastpathDTCborIO) DecSliceBytesN(v [][]byte, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceFloat32X(vp *[]float32, d *decoderCborIO) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceFloat32Y(v []float32, d *decoderCborIO) (v2 []float32, changed bool) {
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
func (fastpathDTCborIO) DecSliceFloat32N(v []float32, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceFloat64X(vp *[]float64, d *decoderCborIO) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceFloat64Y(v []float64, d *decoderCborIO) (v2 []float64, changed bool) {
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
func (fastpathDTCborIO) DecSliceFloat64N(v []float64, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceUint8X(vp *[]uint8, d *decoderCborIO) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceUint8Y(v []uint8, d *decoderCborIO) (v2 []uint8, changed bool) {
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
func (fastpathDTCborIO) DecSliceUint8N(v []uint8, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceUint64X(vp *[]uint64, d *decoderCborIO) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceUint64Y(v []uint64, d *decoderCborIO) (v2 []uint64, changed bool) {
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
func (fastpathDTCborIO) DecSliceUint64N(v []uint64, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceIntX(vp *[]int, d *decoderCborIO) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceIntY(v []int, d *decoderCborIO) (v2 []int, changed bool) {
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
func (fastpathDTCborIO) DecSliceIntN(v []int, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceInt32X(vp *[]int32, d *decoderCborIO) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceInt32Y(v []int32, d *decoderCborIO) (v2 []int32, changed bool) {
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
func (fastpathDTCborIO) DecSliceInt32N(v []int32, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceInt64X(vp *[]int64, d *decoderCborIO) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceInt64Y(v []int64, d *decoderCborIO) (v2 []int64, changed bool) {
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
func (fastpathDTCborIO) DecSliceInt64N(v []int64, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecSliceBoolX(vp *[]bool, d *decoderCborIO) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTCborIO) DecSliceBoolY(v []bool, d *decoderCborIO) (v2 []bool, changed bool) {
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
func (fastpathDTCborIO) DecSliceBoolN(v []bool, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringIntfX(vp *map[string]interface{}, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapStringIntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborIO) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringStringX(vp *map[string]string, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]string, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapStringStringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringStringL(v map[string]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringBytesX(vp *map[string][]byte, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string][]byte, decInferLen(containerLen, d.h.MaxInitLen, 40))
		}
		if containerLen != 0 {
			f.DecMapStringBytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[string][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborIO) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringUint8X(vp *map[string]uint8, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]uint8, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapStringUint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringUint64X(vp *map[string]uint64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]uint64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapStringUint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringIntX(vp *map[string]int, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]int, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapStringIntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringIntL(v map[string]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringInt32X(vp *map[string]int32, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]int32, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			f.DecMapStringInt32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringFloat64X(vp *map[string]float64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]float64, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapStringFloat64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapStringBoolX(vp *map[string]bool, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[string]bool, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapStringBoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapUint8IntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborIO) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8StringX(vp *map[uint8]string, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]string, decInferLen(containerLen, d.h.MaxInitLen, 17))
		}
		if containerLen != 0 {
			f.DecMapUint8StringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8][]byte, decInferLen(containerLen, d.h.MaxInitLen, 25))
		}
		if containerLen != 0 {
			f.DecMapUint8BytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint8][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborIO) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]uint8, decInferLen(containerLen, d.h.MaxInitLen, 2))
		}
		if containerLen != 0 {
			f.DecMapUint8Uint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]uint64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint8Uint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8IntX(vp *map[uint8]int, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]int, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint8IntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8Int32X(vp *map[uint8]int32, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]int32, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			f.DecMapUint8Int32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8Float64X(vp *map[uint8]float64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]float64, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint8Float64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint8BoolX(vp *map[uint8]bool, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint8]bool, decInferLen(containerLen, d.h.MaxInitLen, 2))
		}
		if containerLen != 0 {
			f.DecMapUint8BoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapUint64IntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborIO) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64StringX(vp *map[uint64]string, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapUint64StringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapUint64BytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[uint64][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborIO) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint64Uint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapUint64Uint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64IntX(vp *map[uint64]int, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapUint64IntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64Int32X(vp *map[uint64]int32, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapUint64Int32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64Float64X(vp *map[uint64]float64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapUint64Float64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapUint64BoolX(vp *map[uint64]bool, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[uint64]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapUint64BoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntIntfX(vp *map[int]interface{}, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapIntIntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborIO) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntStringX(vp *map[int]string, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]string, decInferLen(containerLen, d.h.MaxInitLen, 24))
		}
		if containerLen != 0 {
			f.DecMapIntStringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntStringL(v map[int]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntBytesX(vp *map[int][]byte, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int][]byte, decInferLen(containerLen, d.h.MaxInitLen, 32))
		}
		if containerLen != 0 {
			f.DecMapIntBytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborIO) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntUint8X(vp *map[int]uint8, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]uint8, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapIntUint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntUint64X(vp *map[int]uint64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]uint64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapIntUint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntIntX(vp *map[int]int, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]int, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapIntIntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntIntL(v map[int]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntInt32X(vp *map[int]int32, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]int32, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapIntInt32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntFloat64X(vp *map[int]float64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]float64, decInferLen(containerLen, d.h.MaxInitLen, 16))
		}
		if containerLen != 0 {
			f.DecMapIntFloat64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapIntBoolX(vp *map[int]bool, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int]bool, decInferLen(containerLen, d.h.MaxInitLen, 9))
		}
		if containerLen != 0 {
			f.DecMapIntBoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]interface{}, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			f.DecMapInt32IntfL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32]interface{} given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset && !d.h.InterfaceReset
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
func (d *decoderCborIO) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32StringX(vp *map[int32]string, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]string, decInferLen(containerLen, d.h.MaxInitLen, 20))
		}
		if containerLen != 0 {
			f.DecMapInt32StringL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32BytesX(vp *map[int32][]byte, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32][]byte, decInferLen(containerLen, d.h.MaxInitLen, 28))
		}
		if containerLen != 0 {
			f.DecMapInt32BytesL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderCborIO) {
	if v == nil {
		d.errorInt("cannot decode into nil map[int32][]byte given stream length: ", int64(containerLen))
		return
	}
	mapGet := v != nil && !d.h.MapValueReset
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
func (d *decoderCborIO) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]uint8, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			f.DecMapInt32Uint8L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]uint64, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapInt32Uint64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32IntX(vp *map[int32]int, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]int, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapInt32IntL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32Int32X(vp *map[int32]int32, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]int32, decInferLen(containerLen, d.h.MaxInitLen, 8))
		}
		if containerLen != 0 {
			f.DecMapInt32Int32L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32Float64X(vp *map[int32]float64, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]float64, decInferLen(containerLen, d.h.MaxInitLen, 12))
		}
		if containerLen != 0 {
			f.DecMapInt32Float64L(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (f fastpathDTCborIO) DecMapInt32BoolX(vp *map[int32]bool, d *decoderCborIO) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	if containerLen == containerLenNil {
		*vp = nil
	} else {
		if *vp == nil {
			*vp = make(map[int32]bool, decInferLen(containerLen, d.h.MaxInitLen, 5))
		}
		if containerLen != 0 {
			f.DecMapInt32BoolL(*vp, containerLen, d)
		}
		d.mapEnd()
	}
}
func (fastpathDTCborIO) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderCborIO) {
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

type cborEncDriverIO struct {
	noBuiltInTypes
	encDriverNoState
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT

	h   *CborHandle
	e   *encoderShared
	w   bufioEncWriter
	enc encoderI

	b [40]byte
}

func (e *cborEncDriverIO) EncodeNil() {
	e.w.writen1(cborBdNil)
}

func (e *cborEncDriverIO) EncodeBool(b bool) {
	if b {
		e.w.writen1(cborBdTrue)
	} else {
		e.w.writen1(cborBdFalse)
	}
}

func (e *cborEncDriverIO) EncodeFloat32(f float32) {
	b := math.Float32bits(f)
	if e.h.OptimumSize {
		if h := floatToHalfFloatBits(b); halfFloatToFloatBits(h) == b {
			e.w.writen1(cborBdFloat16)
			e.w.writen2(bigen.PutUint16(h))
			return
		}
	}
	e.w.writen1(cborBdFloat32)
	e.w.writen4(bigen.PutUint32(b))
}

func (e *cborEncDriverIO) EncodeFloat64(f float64) {
	if e.h.OptimumSize {
		if f32 := float32(f); float64(f32) == f {
			e.EncodeFloat32(f32)
			return
		}
	}
	e.w.writen1(cborBdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *cborEncDriverIO) encUint(v uint64, bd byte) {
	if v <= 0x17 {
		e.w.writen1(byte(v) + bd)
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd+0x18, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 0x19)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 0x1a)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 0x1b)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *cborEncDriverIO) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-1-v), cborBaseNegInt)
	} else {
		e.encUint(uint64(v), cborBaseUint)
	}
}

func (e *cborEncDriverIO) EncodeUint(v uint64) {
	e.encUint(v, cborBaseUint)
}

func (e *cborEncDriverIO) encLen(bd byte, length int) {
	e.encUint(uint64(length), bd)
}

func (e *cborEncDriverIO) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else if e.h.TimeRFC3339 {
		e.encUint(0, cborBaseTag)
		e.encStringBytesS(cborBaseString, stringView(t.AppendFormat(e.b[:0], time.RFC3339Nano)))
	} else {
		e.encUint(1, cborBaseTag)
		t = t.UTC().Round(time.Microsecond)
		sec, nsec := t.Unix(), uint64(t.Nanosecond())
		if nsec == 0 {
			e.EncodeInt(sec)
		} else {
			e.EncodeFloat64(float64(sec) + float64(nsec)/1e9)
		}
	}
}

func (e *cborEncDriverIO) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	e.encUint(uint64(xtag), cborBaseTag)
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.EncodeNil()
	} else {
		e.enc.encode(v)
	}
}

func (e *cborEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encUint(uint64(re.Tag), cborBaseTag)

	if re.Value != nil {
		e.enc.encode(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *cborEncDriverIO) WriteArrayStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteArray)
	} else {
		e.encLen(cborBaseArray, length)
	}
}

func (e *cborEncDriverIO) WriteMapStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteMap)
	} else {
		e.encLen(cborBaseMap, length)
	}
}

func (e *cborEncDriverIO) WriteMapEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverIO) WriteArrayEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverIO) EncodeString(v string) {
	bb := cborBaseString
	if e.h.StringToRaw {
		bb = cborBaseBytes
	}
	e.encStringBytesS(bb, v)
}

func (e *cborEncDriverIO) EncodeStringBytesRaw(v []byte) {
	if v == nil {
		e.EncodeNil()
	} else {
		e.encStringBytesS(cborBaseBytes, stringView(v))
	}
}

func (e *cborEncDriverIO) encStringBytesS(bb byte, v string) {
	if e.h.IndefiniteLength {
		if bb == cborBaseBytes {
			e.w.writen1(cborBdIndefiniteBytes)
		} else {
			e.w.writen1(cborBdIndefiniteString)
		}
		var vlen uint = uint(len(v))
		blen := vlen / 4
		if blen == 0 {
			blen = 64
		} else if blen > 1024 {
			blen = 1024
		}
		for i := uint(0); i < vlen; {
			var v2 string
			i2 := i + blen
			if i2 >= i && i2 < vlen {
				v2 = v[i:i2]
			} else {
				v2 = v[i:]
			}
			e.encLen(bb, len(v2))
			e.w.writestr(v2)
			i = i2
		}
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(bb, len(v))
		e.w.writestr(v)
	}
}

type cborDecDriverIO struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	noBuiltInTypes

	h   *CborHandle
	d   *decoderShared
	r   ioDecReader
	dec decoderI
	bdAndBdread
	st    bool
	bytes bool
}

func (d *cborDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *cborDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdNil || d.bd == cborBdUndefined {
		d.bdRead = false
		return true
	}
	return
}

func (d *cborDecDriverIO) TryNil() bool {
	return d.advanceNil()
}

func (d *cborDecDriverIO) skipTags() {
	for d.bd>>5 == cborMajorTag {
		d.decUint()
		d.bd = d.r.readn1()
	}
}

func (d *cborDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdNil {
		d.bdRead = false
		return valueTypeNil
	}
	major := d.bd >> 5
	if major == cborMajorBytes {
		return valueTypeBytes
	} else if major == cborMajorString {
		return valueTypeString
	} else if major == cborMajorArray {
		return valueTypeArray
	} else if major == cborMajorMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *cborDecDriverIO) CheckBreak() (v bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdBreak {
		d.bdRead = false
		v = true
	}
	return
}

func (d *cborDecDriverIO) decUint() (ui uint64) {
	v := d.bd & 0x1f
	if v <= 0x17 {
		ui = uint64(v)
	} else if v == 0x18 {
		ui = uint64(d.r.readn1())
	} else if v == 0x19 {
		ui = uint64(bigen.Uint16(d.r.readn2()))
	} else if v == 0x1a {
		ui = uint64(bigen.Uint32(d.r.readn4()))
	} else if v == 0x1b {
		ui = uint64(bigen.Uint64(d.r.readn8()))
	} else {
		halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
	}
	return
}

func (d *cborDecDriverIO) decLen() int {
	return int(d.decUint())
}

func (d *cborDecDriverIO) decAppendIndefiniteBytes(bs []byte, major byte) []byte {
	d.bdRead = false
	for !d.CheckBreak() {
		chunkMajor := d.bd >> 5
		if chunkMajor != major {
			halt.errorf("malformed indefinite string/bytes %x (%s); contains chunk with major type %v, expected %v",
				d.bd, cbordesc(d.bd), chunkMajor, major)
		}
		n := uint(d.decLen())
		oldLen := uint(len(bs))
		newLen := oldLen + n
		if newLen > uint(cap(bs)) {
			bs2 := make([]byte, newLen, 2*uint(cap(bs))+n)
			copy(bs2, bs)
			bs = bs2
		} else {
			bs = bs[:newLen]
		}
		d.r.readb(bs[oldLen:newLen])
		if d.h.ValidateUnicode && major == cborMajorString && !utf8.Valid(bs[oldLen:newLen]) {
			halt.errorf("indefinite-length text string contains chunk that is not a valid utf-8 sequence: 0x%x", bs[oldLen:newLen])
		}
		d.bdRead = false
	}
	d.bdRead = false
	return bs
}

func (d *cborDecDriverIO) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.r.readn2()))))
	case cborBdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case cborBdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *cborDecDriverIO) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd >> 5 {
	case cborMajorUint:
		ui = d.decUint()
	case cborMajorNegInt:
		ui = d.decUint()
		neg = true
	default:
		ok = false
	}
	return
}

func (d *cborDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdTrue {
		b = true
	} else if d.bd == cborBdFalse {
	} else {
		halt.errorf("not bool - %s %x/%s", msgBadDesc, d.bd, cbordesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.st {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorMap {
		halt.errorf("error reading map; got major type: %x, expected %x/%s", d.bd>>5, cborMajorMap, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverIO) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.st {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorArray {
		halt.errorf("invalid array; got major type: %x, expect: %x/%s", d.bd>>5, cborMajorArray, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverIO) DecodeBytes(bs []byte) (bsOut []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	if d.st {
		d.skipTags()
	}
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		d.bdRead = false
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			return d.decAppendIndefiniteBytes(d.d.b[:0], d.bd>>5)
		}
		return d.decAppendIndefiniteBytes(bs[:0], d.bd>>5)
	}
	if d.bd == cborBdIndefiniteArray {
		d.bdRead = false
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			bs = d.d.b[:0]
		} else {
			bs = bs[:0]
		}
		for !d.CheckBreak() {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}
	if d.bd>>5 == cborMajorArray {
		d.bdRead = false
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			bs = d.d.b[:]
		}
		slen := d.decLen()
		var changed bool
		if bs, changed = usableByteSlice(bs, slen); changed {
			d.d.decByteState = decByteStateNone
		}
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}
	clen := d.decLen()
	d.bdRead = false
	if d.bytes && d.h.ZeroCopy {
		d.d.decByteState = decByteStateZerocopy
		return d.r.readx(uint(clen))
	}
	if bs == nil {
		d.d.decByteState = decByteStateReuseBuf
		bs = d.d.b[:]
	}
	return decByteSlice(&d.r, clen, d.h.MaxInitLen, bs)
}

func (d *cborDecDriverIO) DecodeStringAsBytes() (s []byte) {
	s = d.DecodeBytes(nil)
	if d.h.ValidateUnicode && !utf8.Valid(s) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", s)
	}
	return
}

func (d *cborDecDriverIO) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
}

func (d *cborDecDriverIO) decodeTime(xtag uint64) (t time.Time) {
	switch xtag {
	case 0:
		var err error
		t, err = time.Parse(time.RFC3339, stringView(d.DecodeStringAsBytes()))
		halt.onerror(err)
	case 1:
		f1, f2 := math.Modf(d.DecodeFloat64())
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		halt.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	}
	t = t.UTC().Round(time.Microsecond)
	return
}

func (d *cborDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	realxtag := d.decUint()
	d.bdRead = false
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = realxtag
		d.dec.decode(&re.Value)
	} else if xtag != realxtag {
		halt.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
	} else if ext == SelfExt {
		d.dec.decodeAs(rv, basetype, false)
	} else {
		d.dec.interfaceExtConvertAndDecode(rv, ext)
	}
	d.bdRead = false
}

func (d *cborDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd >> 5 {
	case cborMajorUint:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case cborMajorNegInt:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case cborMajorBytes:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString, d.h.ZeroCopy)
	case cborMajorString:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case cborMajorArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborMajorMap:
		n.v = valueTypeMap
		decodeFurther = true
	case cborMajorTag:
		n.v = valueTypeExt
		n.u = d.decUint()
		n.l = nil
		if n.u == 0 || n.u == 1 {
			d.bdRead = false
			n.v = valueTypeTime
			n.t = d.decodeTime(n.u)
		} else if d.st && d.h.getExtForTag(n.u) == nil {

			d.bdRead = false
			d.DecodeNaked()
			return
		}
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined:
			n.v = valueTypeNil
		case cborBdFalse:
			n.v = valueTypeBool
			n.b = false
		case cborBdTrue:
			n.v = valueTypeBool
			n.b = true
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			n.v = valueTypeFloat
			n.f = d.DecodeFloat64()
		default:
			halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	}
	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *cborDecDriverIO) uintBytes() (v []byte, ui uint64) {

	switch vv := d.bd & 0x1f; vv {
	case 0x18:
		v = d.r.readx(1)
		ui = uint64(v[0])
	case 0x19:
		v = d.r.readx(2)
		ui = uint64(bigenstd.Uint16(v))
	case 0x1a:
		v = d.r.readx(4)
		ui = uint64(bigenstd.Uint32(v))
	case 0x1b:
		v = d.r.readx(8)
		ui = uint64(bigenstd.Uint64(v))
	default:
		if vv > 0x1b {
			halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
		}
		ui = uint64(vv)
	}
	return
}

func (d *cborDecDriverIO) nextValueBytes(v0 []byte) (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	v = v0
	var h decNextValueBytesHelper

	var cursor uint
	if d.bytes {
		cursor = d.r.numread() - 1
	}
	h.append1(&v, d.bytes, d.bd)
	v = d.nextValueBytesBdReadR(v)
	d.bdRead = false

	if d.bytes {
		v = d.r.bytesReadFrom(cursor)
	}
	return
}

func (d *cborDecDriverIO) nextValueBytesR(v0 []byte) (v []byte) {
	d.readNextBd()
	v = v0
	var h decNextValueBytesHelper
	h.append1(&v, d.bytes, d.bd)
	return d.nextValueBytesBdReadR(v)
}

func (d *cborDecDriverIO) nextValueBytesBdReadR(v0 []byte) (v []byte) {
	v = v0
	var h decNextValueBytesHelper

	var bs []byte
	var ui uint64

	switch d.bd >> 5 {
	case cborMajorUint, cborMajorNegInt:
		bs, _ = d.uintBytes()
		h.appendN(&v, d.bytes, bs...)
	case cborMajorString, cborMajorBytes:
		if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
			for {
				d.readNextBd()
				h.append1(&v, d.bytes, d.bd)
				if d.bd == cborBdBreak {
					break
				}
				bs, ui = d.uintBytes()
				h.appendN(&v, d.bytes, bs...)
				h.appendN(&v, d.bytes, d.r.readx(uint(ui))...)
			}
		} else {
			bs, ui = d.uintBytes()
			h.appendN(&v, d.bytes, bs...)
			h.appendN(&v, d.bytes, d.r.readx(uint(ui))...)
		}
	case cborMajorArray:
		if d.bd == cborBdIndefiniteArray {
			for {
				d.readNextBd()
				h.append1(&v, d.bytes, d.bd)
				if d.bd == cborBdBreak {
					break
				}
				v = d.nextValueBytesBdReadR(v)
			}
		} else {
			bs, ui = d.uintBytes()
			h.appendN(&v, d.bytes, bs...)
			for i := uint64(0); i < ui; i++ {
				v = d.nextValueBytesR(v)
			}
		}
	case cborMajorMap:
		if d.bd == cborBdIndefiniteMap {
			for {
				d.readNextBd()
				h.append1(&v, d.bytes, d.bd)
				if d.bd == cborBdBreak {
					break
				}
				v = d.nextValueBytesBdReadR(v)
				v = d.nextValueBytesR(v)
			}
		} else {
			bs, ui = d.uintBytes()
			h.appendN(&v, d.bytes, bs...)
			for i := uint64(0); i < ui; i++ {
				v = d.nextValueBytesR(v)
				v = d.nextValueBytesR(v)
			}
		}
	case cborMajorTag:
		bs, _ = d.uintBytes()
		h.appendN(&v, d.bytes, bs...)
		v = d.nextValueBytesR(v)
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined, cborBdFalse, cborBdTrue:
		case cborBdFloat16:
			h.appendN(&v, d.bytes, d.r.readx(2)...)
		case cborBdFloat32:
			h.appendN(&v, d.bytes, d.r.readx(4)...)
		case cborBdFloat64:
			h.appendN(&v, d.bytes, d.r.readx(8)...)
		default:
			halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
	}
	return
}

func (d *cborDecDriverIO) reset() {
	d.bdAndBdread.reset()
	d.st = d.h.SkipUnexpectedTags
}

func (d *cborEncDriverIO) init(hh Handle, shared *encoderShared, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*CborHandle)
	d.e = shared
	if shared.bytes {
		fp = cborFpEncBytes
	} else {
		fp = cborFpEncIO
	}

	d.init2(enc)
	return
}

func (e *cborEncDriverIO) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *cborEncDriverIO) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *cborEncDriverIO) writerEnd()                        { e.w.end() }

func (e *cborEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *cborEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *cborDecDriverIO) init(hh Handle, shared *decoderShared, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*CborHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = cborFpDecBytes
	} else {
		fp = cborFpDecIO
	}

	d.init2(dec)
	return
}

func (d *cborDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *cborDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *cborDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

func (d *cborDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, cbordesc(d.bd))
}

func (d *cborDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

func (d *cborEncDriverIO) init2(enc encoderI) {
	d.enc = enc
}

func (d *cborDecDriverIO) init2(dec decoderI) {
	d.dec = dec
	d.d.cbor = true
}
