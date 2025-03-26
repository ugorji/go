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

func (e *encoderBincBytes) setContainerState(cs containerState) {
	if cs != 0 {
		e.c = cs
	}
}

func (e *encoderBincBytes) rawExt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *encoderBincBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderBincBytes) selferMarshal(f *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderBincBytes) binaryMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderBincBytes) textMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderBincBytes) jsonMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderBincBytes) raw(f *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderBincBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderBincBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderBincBytes) kBool(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderBincBytes) kTime(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderBincBytes) kString(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderBincBytes) kFloat32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderBincBytes) kFloat64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderBincBytes) kComplex64(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderBincBytes) kComplex128(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderBincBytes) kInt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderBincBytes) kInt8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderBincBytes) kInt16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderBincBytes) kInt32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderBincBytes) kInt64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderBincBytes) kUint(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderBincBytes) kUint8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderBincBytes) kUint16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderBincBytes) kUint32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderBincBytes) kUint64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderBincBytes) kUintptr(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderBincBytes) kErr(f *encFnInfo, rv reflect.Value) {
	e.errorf("unsupported encoding kind %s, for %#v", rv.Kind(), any(rv))
}

func (e *encoderBincBytes) kSeqFn(rtelem reflect.Type) (fn *encFnBincBytes) {
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}

	if rtelem.Kind() != reflect.Interface {
		fn = e.fn(rtelem)
	}
	return
}

func (e *encoderBincBytes) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincBytes) kSliceW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincBytes) kArrayW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincBytes) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoderBincBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoderBincBytes) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderBincBytes) kStructSfi(f *encFnInfo) []*structFieldInfo {
	if e.h.Canonical {
		return f.ti.sfi.sorted()
	}
	return f.ti.sfi.source()
}

func (e *encoderBincBytes) kStructNoOmitempty(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincBytes) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
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

func (e *encoderBincBytes) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincBytes) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	var keyFn, valFn *encFnBincBytes

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

func (e *encoderBincBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnBincBytes) {

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

type encoderBincBytes struct {
	panicHdl
	perType encPerType

	dh helperEncDriverBincBytes

	fp *fastpathEsBincBytes

	h *BasicHandle

	e bincEncDriverBytes

	encoderShared

	hh Handle

	ci []interface{}

	slist sfiRvFreelist
}

func (e *encoderBincBytes) HandleName() string {
	return e.hh.Name()
}

func (e *encoderBincBytes) Release() {
}

func (e *encoderBincBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderShared, e).(*fastpathEsBincBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderBincBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderBincBytes) MustEncode(v interface{}) {
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

func (e *encoderBincBytes) Encode(v interface{}) (err error) {

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

func (e *encoderBincBytes) encode(iv interface{}) {

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

func (e *encoderBincBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderBincBytes) encodeValue(rv reflect.Value, fn *encFnBincBytes) {

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

func (e *encoderBincBytes) encodeValueNonNil(rv reflect.Value, fn *encFnBincBytes) {
	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if fn.i.addrE {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderBincBytes) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
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

func (e *encoderBincBytes) marshalUtf8(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderBincBytes) marshalAsis(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderBincBytes) marshalRaw(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoderBincBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		e.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderBincBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}

func (e *encoderBincBytes) fn(t reflect.Type) *encFnBincBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoderBincBytes) fnNoExt(t reflect.Type) *encFnBincBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

func (e *encoderBincBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderBincBytes) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoderBincBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderBincBytes) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoderBincBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderBincBytes) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoderBincBytes) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

func (e *encoderBincBytes) haltOnMbsOddLen(length int) {
	if length&1 != 0 {
		e.errorInt("mapBySlice requires even slice length, but got ", int64(length))
	}
}

func (e *encoderBincBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderBincBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderBincBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderBincBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderBincBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverBincBytes) newEncoderBytes(out *[]byte, h Handle) *encoderBincBytes {
	var c1 encoderBincBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverBincBytes) newEncoderIO(out io.Writer, h Handle) *encoderBincBytes {
	var c1 encoderBincBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverBincBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsBincBytes) (f *fastpathEBincBytes, u reflect.Type) {
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

type encFnBincBytes struct {
	i  encFnInfo
	fe func(*encoderBincBytes, *encFnInfo, reflect.Value)
}
type encRtidFnBincBytes struct {
	rtid uintptr
	fn   *encFnBincBytes
}

func (helperEncDriverBincBytes) encFindRtidFn(s []encRtidFnBincBytes, rtid uintptr) (i uint, fn *encFnBincBytes) {

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

func (helperEncDriverBincBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnBincBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnBincBytes](v))
	}
	return
}

func (dh helperEncDriverBincBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsBincBytes, checkExt bool) (fn *encFnBincBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverBincBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnBincBytes
	sp = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverBincBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnBincBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnBincBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnBincBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnBincBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverBincBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincBytes) {
	fn = new(encFnBincBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderBincBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderBincBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderBincBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderBincBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fe = (*encoderBincBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderBincBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderBincBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderBincBytes).textMarshal
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
					fn.fe = func(e *encoderBincBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderBincBytes).kBool
			case reflect.String:

				fn.fe = (*encoderBincBytes).kString
			case reflect.Int:
				fn.fe = (*encoderBincBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderBincBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderBincBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderBincBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderBincBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderBincBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderBincBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderBincBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderBincBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderBincBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderBincBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderBincBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderBincBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderBincBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderBincBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderBincBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderBincBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderBincBytes).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*encoderBincBytes).kStruct
				} else {
					fn.fe = (*encoderBincBytes).kStructNoOmitempty
				}
			case reflect.Map:
				fn.fe = (*encoderBincBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderBincBytes).kErr
			default:

				fn.fe = (*encoderBincBytes).kErr
			}
		}
	}
	return
}

type helperEncDriverBincBytes struct{}

func (d *decoderBincBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *decoderBincBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderBincBytes) selferUnmarshal(f *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderBincBytes) binaryUnmarshal(f *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	d.onerror(fnerr)
}

func (d *decoderBincBytes) textUnmarshal(f *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	d.onerror(fnerr)
}

func (d *decoderBincBytes) jsonUnmarshal(f *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderBincBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

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

func (d *decoderBincBytes) kErr(f *decFnInfo, rv reflect.Value) {
	d.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderBincBytes) raw(f *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderBincBytes) kString(f *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes()))
}

func (d *decoderBincBytes) kBool(f *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderBincBytes) kTime(f *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderBincBytes) kFloat32(f *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderBincBytes) kFloat64(f *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderBincBytes) kComplex64(f *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderBincBytes) kComplex128(f *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderBincBytes) kInt(f *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderBincBytes) kInt8(f *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderBincBytes) kInt16(f *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderBincBytes) kInt32(f *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderBincBytes) kInt64(f *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderBincBytes) kUint(f *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincBytes) kUintptr(f *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincBytes) kUint8(f *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderBincBytes) kUint16(f *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderBincBytes) kUint32(f *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderBincBytes) kUint64(f *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderBincBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderBincBytes) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderBincBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *decoderBincBytes) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderBincBytes) kSlice(f *decFnInfo, rv reflect.Value) {

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

	var fn *decFnBincBytes

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

func (d *decoderBincBytes) kArray(f *decFnInfo, rv reflect.Value) {

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

	var fn *decFnBincBytes

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

func (d *decoderBincBytes) kChan(f *decFnInfo, rv reflect.Value) {

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

	var fn *decFnBincBytes

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

func (d *decoderBincBytes) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnBincBytes
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

type decoderBincBytes struct {
	panicHdl
	perType decPerType

	dh helperDecDriverBincBytes

	fp *fastpathDsBincBytes

	h *BasicHandle

	d bincDecDriverBytes

	decoderShared

	hh Handle

	mtid uintptr
	stid uintptr
}

func (d *decoderBincBytes) HandleName() string {
	return d.hh.Name()
}

func (d *decoderBincBytes) isBytes() bool {
	return d.bytes
}

func (d *decoderBincBytes) init(h Handle) {
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

	d.fp = d.d.init(h, &d.decoderShared, d).(*fastpathDsBincBytes)

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

func (d *decoderBincBytes) reset() {
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

func (d *decoderBincBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderBincBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
	return
}

func (d *decoderBincBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderBincBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderBincBytes) Decode(v interface{}) (err error) {

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

func (d *decoderBincBytes) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderBincBytes) Release() {
}

func (d *decoderBincBytes) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *decoderBincBytes) nextValueBytes(start []byte) []byte {
	return d.d.nextValueBytes(start)
}

func (d *decoderBincBytes) decode(iv interface{}) {

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

func (d *decoderBincBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderBincBytes) decodeValue(rv reflect.Value, fn *decFnBincBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *decoderBincBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnBincBytes) {

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

func (d *decoderBincBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			d.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			d.errorStr("no matching struct field found when decoding stream map with key " + rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderBincBytes) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		d.errorf("cannot expand array len during decode from %v to %v", any(sliceLen), any(streamLen))
	}
}

func (d *decoderBincBytes) haltAsNotDecodeable(rv reflect.Value) {
	if !rv.IsValid() {
		d.onerror(errCannotDecodeIntoNil)
	}

	if !rv.CanInterface() {
		d.errorf("cannot decode into a value without an interface: %v", rv)
	}
	d.errorf("cannot decode into value of kind: %v, %#v", rv.Kind(), rv2i(rv))
}

func (d *decoderBincBytes) depthIncr() {
	d.depth++
	if d.depth >= d.maxdepth {
		d.onerror(errMaxDepthExceeded)
	}
}

func (d *decoderBincBytes) depthDecr() {
	d.depth--
}

func (d *decoderBincBytes) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = zeroByteSlice
	}
	return d.d.DecodeBytes(in)
}

func (d *decoderBincBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes(zeroByteSlice)
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderBincBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.NumBytesRead(), false)
}

func (d *decoderBincBytes) NumBytesRead() int {
	return int(d.d.NumBytesRead())
}

func (d *decoderBincBytes) checkBreak() (v bool) {
	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *decoderBincBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *decoderBincBytes) mapStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerMapStart
	}
	return v
}

func (d *decoderBincBytes) mapElemKey() {
	d.d.ReadMapElemKey()
	d.c = containerMapKey
}

func (d *decoderBincBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderBincBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincBytes) arrayStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerArrayStart
	}
	return v
}

func (d *decoderBincBytes) arrayElem() {
	d.d.ReadArrayElem()
	d.c = containerArrayElem
}

func (d *decoderBincBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

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

func (d *decoderBincBytes) oneShotAddrRV(rvt reflect.Type, rvk reflect.Kind) reflect.Value {
	if decUseTransient &&
		(numBoolStrSliceBitset.isset(byte(rvk)) ||
			((rvk == reflect.Struct || rvk == reflect.Array) &&
				d.h.getTypeInfo(rt2id(rvt), rvt).flagCanTransient)) {
		return d.perType.TransientAddrK(rvt, rvk)
	}
	return rvZeroAddrK(rvt, rvk)
}

func (d *decoderBincBytes) fn(t reflect.Type) *decFnBincBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, false)
}

func (d *decoderBincBytes) fnNoExt(t reflect.Type) *decFnBincBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, true)
}

type decSliceHelperBincBytes struct {
	d     *decoderBincBytes
	ct    valueType
	Array bool
	IsNil bool
}

func (d *decoderBincBytes) decSliceHelperStart() (x decSliceHelperBincBytes, clen int) {
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

func (x decSliceHelperBincBytes) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelperBincBytes) ElemContainerState(index int) {

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 {
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelperBincBytes) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
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

func (helperDecDriverBincBytes) newDecoderBytes(in []byte, h Handle) *decoderBincBytes {
	var c1 decoderBincBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverBincBytes) newDecoderIO(in io.Reader, h Handle) *decoderBincBytes {
	var c1 decoderBincBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverBincBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsBincBytes) (f *fastpathDBincBytes, u reflect.Type) {
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

type decFnBincBytes struct {
	i  decFnInfo
	fd func(*decoderBincBytes, *decFnInfo, reflect.Value)
}
type decRtidFnBincBytes struct {
	rtid uintptr
	fn   *decFnBincBytes
}

func (helperDecDriverBincBytes) decFindRtidFn(s []decRtidFnBincBytes, rtid uintptr) (i uint, fn *decFnBincBytes) {

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

func (helperDecDriverBincBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnBincBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnBincBytes](v))
	}
	return
}

func (dh helperDecDriverBincBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsBincBytes,
	checkExt bool) (fn *decFnBincBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverBincBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnBincBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverBincBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnBincBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnBincBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnBincBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnBincBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverBincBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincBytes) {
	fn = new(decFnBincBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderBincBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderBincBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderBincBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderBincBytes).ext
		fi.addrD = true
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fd = (*decoderBincBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderBincBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderBincBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderBincBytes).textUnmarshal
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
						fn.fd = func(d *decoderBincBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderBincBytes, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderBincBytes).kBool
			case reflect.String:
				fn.fd = (*decoderBincBytes).kString
			case reflect.Int:
				fn.fd = (*decoderBincBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderBincBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderBincBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderBincBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderBincBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderBincBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderBincBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderBincBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderBincBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderBincBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderBincBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderBincBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderBincBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderBincBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderBincBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderBincBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderBincBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderBincBytes).kArray
			case reflect.Struct:
				fn.fd = (*decoderBincBytes).kStruct
			case reflect.Map:
				fn.fd = (*decoderBincBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderBincBytes).kInterface
			default:

				fn.fd = (*decoderBincBytes).kErr
			}
		}
	}
	return
}

type helperDecDriverBincBytes struct{}
type fastpathEBincBytes struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderBincBytes, *encFnInfo, reflect.Value)
}
type fastpathDBincBytes struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderBincBytes, *decFnInfo, reflect.Value)
}
type fastpathEsBincBytes [56]fastpathEBincBytes
type fastpathDsBincBytes [56]fastpathDBincBytes
type fastpathETBincBytes struct{}
type fastpathDTBincBytes struct{}

func (helperEncDriverBincBytes) fastpathEList() *fastpathEsBincBytes {
	var i uint = 0
	var s fastpathEsBincBytes
	fn := func(v interface{}, fe func(*encoderBincBytes, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathEBincBytes{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderBincBytes).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderBincBytes).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderBincBytes).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderBincBytes).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderBincBytes).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderBincBytes).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderBincBytes).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderBincBytes).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderBincBytes).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderBincBytes).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderBincBytes).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderBincBytes).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderBincBytes).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderBincBytes).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderBincBytes).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderBincBytes).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderBincBytes).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderBincBytes).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderBincBytes).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderBincBytes).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderBincBytes).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderBincBytes).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderBincBytes).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderBincBytes).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderBincBytes).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderBincBytes).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderBincBytes).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderBincBytes).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderBincBytes).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderBincBytes).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderBincBytes).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderBincBytes).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderBincBytes).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderBincBytes).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderBincBytes).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderBincBytes).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderBincBytes).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderBincBytes).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderBincBytes).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderBincBytes).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderBincBytes).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderBincBytes).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderBincBytes).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderBincBytes).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderBincBytes).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderBincBytes).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderBincBytes).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderBincBytes).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderBincBytes).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderBincBytes).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderBincBytes).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderBincBytes).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderBincBytes).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderBincBytes).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderBincBytes).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderBincBytes).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverBincBytes) fastpathDList() *fastpathDsBincBytes {
	var i uint = 0
	var s fastpathDsBincBytes
	fn := func(v interface{}, fd func(*decoderBincBytes, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDBincBytes{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderBincBytes).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderBincBytes).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderBincBytes).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderBincBytes).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderBincBytes).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderBincBytes).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderBincBytes).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderBincBytes).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderBincBytes).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderBincBytes).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderBincBytes).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderBincBytes).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderBincBytes).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderBincBytes).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderBincBytes).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderBincBytes).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderBincBytes).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderBincBytes).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderBincBytes).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderBincBytes).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderBincBytes).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderBincBytes).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderBincBytes).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderBincBytes).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderBincBytes).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderBincBytes).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderBincBytes).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderBincBytes).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderBincBytes).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderBincBytes).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderBincBytes).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderBincBytes).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderBincBytes).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderBincBytes).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderBincBytes).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderBincBytes).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderBincBytes).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderBincBytes).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderBincBytes).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderBincBytes).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderBincBytes).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderBincBytes).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderBincBytes).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderBincBytes).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderBincBytes).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderBincBytes).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderBincBytes).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderBincBytes).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderBincBytes).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderBincBytes).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderBincBytes).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderBincBytes).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderBincBytes).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderBincBytes).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderBincBytes).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderBincBytes).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverBincBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderBincBytes) bool {
	var ft fastpathETBincBytes
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

func (e *encoderBincBytes) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceIntfV(v []interface{}, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceIntfV(v []interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceStringV(v []string, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceStringV(v []string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceBytesV(v [][]byte, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceBytesV(v [][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceFloat32V(v []float32, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceFloat32V(v []float32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceFloat64V(v []float64, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceFloat64V(v []float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceUint8V(v []uint8, e *encoderBincBytes) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETBincBytes) EncAsMapSliceUint8V(v []uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceUint64V(v []uint64, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceUint64V(v []uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceIntV(v []int, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceIntV(v []int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceInt32V(v []int32, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceInt32V(v []int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceInt64V(v []int64, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceInt64V(v []int64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceBoolV(v []bool, e *encoderBincBytes) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincBytes) EncAsMapSliceBoolV(v []bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETBincBytes) EncMapStringIntfV(v map[string]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETBincBytes) EncMapStringStringV(v map[string]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETBincBytes) EncMapStringBytesV(v map[string][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETBincBytes) EncMapStringUint8V(v map[string]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETBincBytes) EncMapStringUint64V(v map[string]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETBincBytes) EncMapStringIntV(v map[string]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETBincBytes) EncMapStringInt32V(v map[string]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETBincBytes) EncMapStringFloat64V(v map[string]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETBincBytes) EncMapStringBoolV(v map[string]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETBincBytes) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETBincBytes) EncMapUint8StringV(v map[uint8]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETBincBytes) EncMapUint8BytesV(v map[uint8][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETBincBytes) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETBincBytes) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETBincBytes) EncMapUint8IntV(v map[uint8]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETBincBytes) EncMapUint8Int32V(v map[uint8]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETBincBytes) EncMapUint8Float64V(v map[uint8]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETBincBytes) EncMapUint8BoolV(v map[uint8]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETBincBytes) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETBincBytes) EncMapUint64StringV(v map[uint64]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETBincBytes) EncMapUint64BytesV(v map[uint64][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETBincBytes) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETBincBytes) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETBincBytes) EncMapUint64IntV(v map[uint64]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETBincBytes) EncMapUint64Int32V(v map[uint64]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETBincBytes) EncMapUint64Float64V(v map[uint64]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETBincBytes) EncMapUint64BoolV(v map[uint64]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETBincBytes) EncMapIntIntfV(v map[int]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETBincBytes) EncMapIntStringV(v map[int]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETBincBytes) EncMapIntBytesV(v map[int][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETBincBytes) EncMapIntUint8V(v map[int]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETBincBytes) EncMapIntUint64V(v map[int]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETBincBytes) EncMapIntIntV(v map[int]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETBincBytes) EncMapIntInt32V(v map[int]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETBincBytes) EncMapIntFloat64V(v map[int]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETBincBytes) EncMapIntBoolV(v map[int]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETBincBytes) EncMapInt32IntfV(v map[int32]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETBincBytes) EncMapInt32StringV(v map[int32]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETBincBytes) EncMapInt32BytesV(v map[int32][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETBincBytes) EncMapInt32Uint8V(v map[int32]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETBincBytes) EncMapInt32Uint64V(v map[int32]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETBincBytes) EncMapInt32IntV(v map[int32]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETBincBytes) EncMapInt32Int32V(v map[int32]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETBincBytes) EncMapInt32Float64V(v map[int32]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
	ft.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETBincBytes) EncMapInt32BoolV(v map[int32]bool, e *encoderBincBytes) {
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

func (helperDecDriverBincBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderBincBytes) bool {
	var ft fastpathDTBincBytes
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

func (d *decoderBincBytes) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceIntfX(vp *[]interface{}, d *decoderBincBytes) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceIntfY(v []interface{}, d *decoderBincBytes) (v2 []interface{}, changed bool) {
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
func (fastpathDTBincBytes) DecSliceIntfN(v []interface{}, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceStringX(vp *[]string, d *decoderBincBytes) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceStringY(v []string, d *decoderBincBytes) (v2 []string, changed bool) {
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
func (fastpathDTBincBytes) DecSliceStringN(v []string, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceBytesX(vp *[][]byte, d *decoderBincBytes) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceBytesY(v [][]byte, d *decoderBincBytes) (v2 [][]byte, changed bool) {
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
func (fastpathDTBincBytes) DecSliceBytesN(v [][]byte, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceFloat32X(vp *[]float32, d *decoderBincBytes) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceFloat32Y(v []float32, d *decoderBincBytes) (v2 []float32, changed bool) {
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
func (fastpathDTBincBytes) DecSliceFloat32N(v []float32, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceFloat64X(vp *[]float64, d *decoderBincBytes) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceFloat64Y(v []float64, d *decoderBincBytes) (v2 []float64, changed bool) {
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
func (fastpathDTBincBytes) DecSliceFloat64N(v []float64, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceUint8X(vp *[]uint8, d *decoderBincBytes) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceUint8Y(v []uint8, d *decoderBincBytes) (v2 []uint8, changed bool) {
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
func (fastpathDTBincBytes) DecSliceUint8N(v []uint8, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceUint64X(vp *[]uint64, d *decoderBincBytes) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceUint64Y(v []uint64, d *decoderBincBytes) (v2 []uint64, changed bool) {
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
func (fastpathDTBincBytes) DecSliceUint64N(v []uint64, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceIntX(vp *[]int, d *decoderBincBytes) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceIntY(v []int, d *decoderBincBytes) (v2 []int, changed bool) {
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
func (fastpathDTBincBytes) DecSliceIntN(v []int, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceInt32X(vp *[]int32, d *decoderBincBytes) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceInt32Y(v []int32, d *decoderBincBytes) (v2 []int32, changed bool) {
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
func (fastpathDTBincBytes) DecSliceInt32N(v []int32, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceInt64X(vp *[]int64, d *decoderBincBytes) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceInt64Y(v []int64, d *decoderBincBytes) (v2 []int64, changed bool) {
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
func (fastpathDTBincBytes) DecSliceInt64N(v []int64, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecSliceBoolX(vp *[]bool, d *decoderBincBytes) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincBytes) DecSliceBoolY(v []bool, d *decoderBincBytes) (v2 []bool, changed bool) {
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
func (fastpathDTBincBytes) DecSliceBoolN(v []bool, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringIntfX(vp *map[string]interface{}, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringStringX(vp *map[string]string, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringStringL(v map[string]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringBytesX(vp *map[string][]byte, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringUint8X(vp *map[string]uint8, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringUint64X(vp *map[string]uint64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringIntX(vp *map[string]int, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringIntL(v map[string]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringInt32X(vp *map[string]int32, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringFloat64X(vp *map[string]float64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapStringBoolX(vp *map[string]bool, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8StringX(vp *map[uint8]string, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8IntX(vp *map[uint8]int, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8Int32X(vp *map[uint8]int32, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8Float64X(vp *map[uint8]float64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint8BoolX(vp *map[uint8]bool, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64StringX(vp *map[uint64]string, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64IntX(vp *map[uint64]int, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64Int32X(vp *map[uint64]int32, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64Float64X(vp *map[uint64]float64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapUint64BoolX(vp *map[uint64]bool, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntIntfX(vp *map[int]interface{}, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntStringX(vp *map[int]string, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntStringL(v map[int]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntBytesX(vp *map[int][]byte, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntUint8X(vp *map[int]uint8, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntUint64X(vp *map[int]uint64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntIntX(vp *map[int]int, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntIntL(v map[int]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntInt32X(vp *map[int]int32, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntFloat64X(vp *map[int]float64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapIntBoolX(vp *map[int]bool, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32StringX(vp *map[int32]string, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32BytesX(vp *map[int32][]byte, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32IntX(vp *map[int32]int, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32Int32X(vp *map[int32]int32, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32Float64X(vp *map[int32]float64, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (f fastpathDTBincBytes) DecMapInt32BoolX(vp *map[int32]bool, d *decoderBincBytes) {
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
func (fastpathDTBincBytes) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderBincBytes) {
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

type bincEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT
	encInit2er

	h *BincHandle
	e *encoderShared
	w bytesEncAppender
	bincEncState
}

func (e *bincEncDriverBytes) EncodeNil() {
	e.w.writen1(bincBdNil)
}

func (e *bincEncDriverBytes) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else {
		bs := bincEncodeTime(t)
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	}
}

func (e *bincEncDriverBytes) EncodeBool(b bool) {
	if b {
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	}
}

func (e *bincEncDriverBytes) encSpFloat(f float64) (done bool) {
	if f == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
	} else if math.IsNaN(float64(f)) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNan)
	} else if math.IsInf(float64(f), +1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpPosInf)
	} else if math.IsInf(float64(f), -1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegInf)
	} else {
		return
	}
	return true
}

func (e *bincEncDriverBytes) EncodeFloat32(f float32) {
	if !e.encSpFloat(float64(f)) {
		e.w.writen1(bincVdFloat<<4 | bincFlBin32)
		e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
	}
}

func (e *bincEncDriverBytes) EncodeFloat64(f float64) {
	if e.encSpFloat(f) {
		return
	}
	b := bigen.PutUint64(math.Float64bits(f))
	if bincDoPrune {
		i := 7
		for ; i >= 0 && (b[i] == 0); i-- {
		}
		i++
		if i <= 6 {
			e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.w.writen1(byte(i))
			e.w.writeb(b[:i])
			return
		}
	}
	e.w.writen1(bincVdFloat<<4 | bincFlBin64)
	e.w.writen8(b)
}

func (e *bincEncDriverBytes) encIntegerPrune32(bd byte, pos bool, v uint64) {
	b := bigen.PutUint32(uint32(v))
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 3 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 3)
		e.w.writen4(b)
	}
}

func (e *bincEncDriverBytes) encIntegerPrune64(bd byte, pos bool, v uint64) {
	b := bigen.PutUint64(v)
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 7 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 7)
		e.w.writen8(b)
	}
}

func (e *bincEncDriverBytes) EncodeInt(v int64) {
	if v >= 0 {
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	} else if v == -1 {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	} else {
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	}
}

func (e *bincEncDriverBytes) EncodeUint(v uint64) {
	e.encUint(bincVdPosInt<<4, true, v)
}

func (e *bincEncDriverBytes) encUint(bd byte, pos bool, v uint64) {
	if v == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	} else if pos && v >= 1 && v <= 16 {
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd|0x0, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.encIntegerPrune32(bd, pos, v)
	} else {
		e.encIntegerPrune64(bd, pos, v)
	}
}

func (e *bincEncDriverBytes) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	var bs0, bs []byte
	if ext == SelfExt {
		bs0 = e.e.blist.get(1024)
		bs = bs0
		sideEncode(&e.h.BasicHandle, v, &bs, basetype, true)
	} else {
		bs = ext.WriteExt(v)
	}
	if bs == nil {
		e.EncodeNil()
		goto END
	}
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *bincEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *bincEncDriverBytes) encodeExtPreamble(xtag byte, length int) {
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncDriverBytes) WriteArrayStart(length int) {
	e.encLen(bincVdArray<<4, uint64(length))
}

func (e *bincEncDriverBytes) WriteMapStart(length int) {
	e.encLen(bincVdMap<<4, uint64(length))
}

func (e *bincEncDriverBytes) EncodeSymbol(v string) {

	l := len(v)
	if l == 0 {
		e.encBytesLen(cUTF8, 0)
		return
	} else if l == 1 {
		e.encBytesLen(cUTF8, 1)
		e.w.writen1(v[0])
		return
	}
	if e.m == nil {
		e.m = make(map[string]uint16, 16)
	}
	ui, ok := e.m[v]
	if ok {
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			e.w.writen2(bigen.PutUint16(ui))
		}
	} else {
		e.e.seq++
		ui = e.e.seq
		e.m[v] = ui
		var lenprec uint8
		if l <= math.MaxUint8 {

		} else if l <= math.MaxUint16 {
			lenprec = 1
		} else if int64(l) <= math.MaxUint32 {
			lenprec = 2
		} else {
			lenprec = 3
		}
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4|0x0|0x4|lenprec, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			e.w.writen2(bigen.PutUint16(ui))
		}
		if lenprec == 0 {
			e.w.writen1(byte(l))
		} else if lenprec == 1 {
			e.w.writen2(bigen.PutUint16(uint16(l)))
		} else if lenprec == 2 {
			e.w.writen4(bigen.PutUint32(uint32(l)))
		} else {
			e.w.writen8(bigen.PutUint64(uint64(l)))
		}
		e.w.writestr(v)
	}
}

func (e *bincEncDriverBytes) EncodeString(v string) {
	if e.h.StringToRaw {
		e.encLen(bincVdByteArray<<4, uint64(len(v)))
		if len(v) > 0 {
			e.w.writestr(v)
		}
		return
	}
	e.EncodeStringEnc(cUTF8, v)
}

func (e *bincEncDriverBytes) EncodeStringEnc(c charEncoding, v string) {
	if e.e.c == containerMapKey && c == cUTF8 && (e.h.AsSymbols == 1) {
		e.EncodeSymbol(v)
		return
	}
	e.encLen(bincVdString<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writestr(v)
	}
}

func (e *bincEncDriverBytes) EncodeStringBytesRaw(v []byte) {
	if v == nil {
		e.EncodeNil()
		return
	}
	e.encLen(bincVdByteArray<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writeb(v)
	}
}

func (e *bincEncDriverBytes) encBytesLen(c charEncoding, length uint64) {

	if c == cRAW {
		e.encLen(bincVdByteArray<<4, length)
	} else {
		e.encLen(bincVdString<<4, length)
	}
}

func (e *bincEncDriverBytes) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l+4))
	} else {
		e.encLenNumber(bd, l)
	}
}

func (e *bincEncDriverBytes) encLenNumber(bd byte, v uint64) {
	if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd | 0x02)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd | 0x03)
		e.w.writen8(bigen.PutUint64(uint64(v)))
	}
}

type bincDecDriverBytes struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	decInit2er
	noBuiltInTypes

	h *BincHandle
	d *decoderShared
	r bytesDecReader

	bincDecState

	bytes bool
}

func (d *bincDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
}

func (d *bincDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *bincDecDriverBytes) TryNil() bool {
	return d.advanceNil()
}

func (d *bincDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return valueTypeNil
	} else if d.vd == bincVdByteArray {
		return valueTypeBytes
	} else if d.vd == bincVdString {
		return valueTypeString
	} else if d.vd == bincVdArray {
		return valueTypeArray
	} else if d.vd == bincVdMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *bincDecDriverBytes) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.vd != bincVdTimestamp {
		halt.errorf("cannot decode time - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	t, err := bincDecodeTime(d.r.readx(uint(d.vs)))
	halt.onerror(err)
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) decFloatPruned(maxlen uint8) {
	l := d.r.readn1()
	if l > maxlen {
		halt.errorf("cannot read float - at most %v bytes used to represent float - received %v bytes", maxlen, l)
	}
	for i := l; i < maxlen; i++ {
		d.d.b[i] = 0
	}
	d.r.readb(d.d.b[0:l])
}

func (d *bincDecDriverBytes) decFloatPre32() (b [4]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn4()
	} else {
		d.decFloatPruned(4)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverBytes) decFloatPre64() (b [8]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn8()
	} else {
		d.decFloatPruned(8)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverBytes) decFloatVal() (f float64) {
	switch d.vs & 0x7 {
	case bincFlBin32:
		f = float64(math.Float32frombits(bigen.Uint32(d.decFloatPre32())))
	case bincFlBin64:
		f = math.Float64frombits(bigen.Uint64(d.decFloatPre64()))
	default:

		halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincDecDriverBytes) decUint() (v uint64) {
	switch d.vs {
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		v = uint64(bigen.Uint16(d.r.readn2()))
	case 2:
		b3 := d.r.readn3()
		var b [4]byte
		copy(b[1:], b3[:])
		v = uint64(bigen.Uint32(b))
	case 3:
		v = uint64(bigen.Uint32(d.r.readn4()))
	case 4, 5, 6:
		var b [8]byte
		lim := 7 - d.vs
		bs := d.d.b[lim:8]
		d.r.readb(bs)
		copy(b[lim:], bs)
		v = bigen.Uint64(b)
	case 7:
		v = bigen.Uint64(d.r.readn8())
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	return
}

func (d *bincDecDriverBytes) uintBytes() (bs []byte) {
	switch d.vs {
	case 0:
		bs = d.d.b[:1]
		bs[0] = d.r.readn1()
	case 1:
		bs = d.d.b[:2]
		d.r.readb(bs)
	case 2:
		bs = d.d.b[:3]
		d.r.readb(bs)
	case 3:
		bs = d.d.b[:4]
		d.r.readb(bs)
	case 4, 5, 6:
		lim := 7 - d.vs
		bs = d.d.b[lim:8]
		d.r.readb(bs)
	case 7:
		bs = d.d.b[:8]
		d.r.readb(bs)
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	return
}

func (d *bincDecDriverBytes) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdPosInt {
		ui = d.decUint()
	} else if vd == bincVdNegInt {
		ui = d.decUint()
		neg = true
	} else if vd == bincVdSmallInt {
		ui = uint64(d.vs) + 1
	} else if vd == bincVdSpecial {
		if vs == bincSpZero {

		} else if vs == bincSpNegOne {
			neg = true
			ui = 1
		} else {
			ok = false

		}
	} else {
		ok = false

	}
	return
}

func (d *bincDecDriverBytes) decFloat() (f float64, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdSpecial {
		if vs == bincSpNan {
			f = math.NaN()
		} else if vs == bincSpPosInf {
			f = math.Inf(1)
		} else if vs == bincSpZeroFloat || vs == bincSpZero {

		} else if vs == bincSpNegInf {
			f = math.Inf(-1)
		} else {
			ok = false

		}
	} else if vd == bincVdFloat {
		f = d.decFloatVal()
	} else {
		ok = false
	}
	return
}

func (d *bincDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == (bincVdSpecial | bincSpFalse) {

	} else if d.bd == (bincVdSpecial | bincSpTrue) {
		b = true
	} else {
		halt.errorf("bool - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdMap {
		halt.errorf("map - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdArray {
		halt.errorf("array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) decLen() int {
	if d.vs > 3 {
		return int(d.vs - 4)
	}
	return int(d.decLenNumber())
}

func (d *bincDecDriverBytes) decLenNumber() (v uint64) {
	if x := d.vs; x == 0 {
		v = uint64(d.r.readn1())
	} else if x == 1 {
		v = uint64(bigen.Uint16(d.r.readn2()))
	} else if x == 2 {
		v = uint64(bigen.Uint32(d.r.readn4()))
	} else {
		v = bigen.Uint64(d.r.readn8())
	}
	return
}

func (d *bincDecDriverBytes) DecodeStringAsBytes() (bs2 []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	var slen = -1
	switch d.vd {
	case bincVdString, bincVdByteArray:
		slen = d.decLen()
		if d.d.bytes {
			d.d.decByteState = decByteStateZerocopy
			bs2 = d.r.readx(uint(slen))
		} else {
			d.d.decByteState = decByteStateReuseBuf
			bs2 = decByteSlice(&d.r, slen, d.h.MaxInitLen, d.d.b[:])
		}
	case bincVdSymbol:

		var symbol uint16
		vs := d.vs
		if vs&0x8 == 0 {
			symbol = uint16(d.r.readn1())
		} else {
			symbol = uint16(bigen.Uint16(d.r.readn2()))
		}
		if d.s == nil {
			d.s = make(map[uint16][]byte, 16)
		}

		if vs&0x4 == 0 {
			bs2 = d.s[symbol]
		} else {
			switch vs & 0x3 {
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(bigen.Uint16(d.r.readn2()))
			case 2:
				slen = int(bigen.Uint32(d.r.readn4()))
			case 3:
				slen = int(bigen.Uint64(d.r.readn8()))
			}

			bs2 = decByteSlice(&d.r, slen, d.h.MaxInitLen, nil)
			d.s[symbol] = bs2
		}
	default:
		halt.errorf("string/bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if d.h.ValidateUnicode && !utf8.Valid(bs2) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", bs2)
	}

	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeBytes(bs []byte) (bsOut []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	if d.vd == bincVdArray {
		if bs == nil {
			bs = d.d.b[:]
			d.d.decByteState = decByteStateReuseBuf
		}
		slen := d.ReadArrayStart()
		var changed bool
		if bs, changed = usableByteSlice(bs, slen); changed {
			d.d.decByteState = decByteStateNone
		}
		for i := 0; i < slen; i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}
	var clen int
	if d.vd == bincVdString || d.vd == bincVdByteArray {
		clen = d.decLen()
	} else {
		halt.errorf("bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	if d.bytes && d.h.ZeroCopy {
		d.d.decByteState = decByteStateZerocopy
		return d.r.readx(uint(clen))
	}
	if bs == nil {
		bs = d.d.b[:]
		d.d.decByteState = decByteStateReuseBuf
	}
	return decByteSlice(&d.r, clen, d.h.MaxInitLen, bs)
}

func (d *bincDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if xtag > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtag)
	}
	if d.advanceNil() {
		return
	}
	xbs, realxtag1, zerocopy := d.decodeExtV(ext != nil, uint8(xtag))
	realxtag := uint64(realxtag1)
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = realxtag
		re.setData(xbs, zerocopy)
	} else if ext == SelfExt {
		sideDecode(&d.h.BasicHandle, rv, xbs, basetype, true)
	} else {
		ext.ReadExt(rv, xbs)
	}
}

func (d *bincDecDriverBytes) decodeExtV(verifyTag bool, tag byte) (xbs []byte, xtag byte, zerocopy bool) {
	if d.vd == bincVdCustomExt {
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag - got %b, expecting: %v", xtag, tag)
		}
		if d.d.bytes {
			xbs = d.r.readx(uint(l))
			zerocopy = true
		} else {
			xbs = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	} else if d.vd == bincVdByteArray {
		xbs = d.DecodeBytes(nil)
	} else {
		halt.errorf("ext expects extensions or byte array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil:
			n.v = valueTypeNil
		case bincSpFalse:
			n.v = valueTypeBool
			n.b = false
		case bincSpTrue:
			n.v = valueTypeBool
			n.b = true
		case bincSpNan:
			n.v = valueTypeFloat
			n.f = math.NaN()
		case bincSpPosInf:
			n.v = valueTypeFloat
			n.f = math.Inf(1)
		case bincSpNegInf:
			n.v = valueTypeFloat
			n.f = math.Inf(-1)
		case bincSpZeroFloat:
			n.v = valueTypeFloat
			n.f = float64(0)
		case bincSpZero:
			n.v = valueTypeUint
			n.u = uint64(0)
		case bincSpNegOne:
			n.v = valueTypeInt
			n.i = int64(-1)
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
		n.v = valueTypeUint
		n.u = uint64(int8(d.vs)) + 1
	case bincVdPosInt:
		n.v = valueTypeUint
		n.u = d.decUint()
	case bincVdNegInt:
		n.v = valueTypeInt
		n.i = -(int64(d.decUint()))
	case bincVdFloat:
		n.v = valueTypeFloat
		n.f = d.decFloatVal()
	case bincVdString:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case bincVdByteArray:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString, d.h.ZeroCopy)
	case bincVdSymbol:
		n.v = valueTypeSymbol
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case bincVdTimestamp:
		n.v = valueTypeTime
		tt, err := bincDecodeTime(d.r.readx(uint(d.vs)))
		halt.onerror(err)
		n.t = tt
	case bincVdCustomExt:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		if d.d.bytes {
			n.l = d.r.readx(uint(l))
		} else {
			n.l = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case bincVdArray:
		n.v = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if !decodeFurther {
		d.bdRead = false
	}
	if n.v == valueTypeUint && d.h.SignedInteger {
		n.v = valueTypeInt
		n.i = int64(n.u)
	}
}

func (d *bincDecDriverBytes) nextValueBytes(v0 []byte) (v []byte) {
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

func (d *bincDecDriverBytes) nextValueBytesR(v0 []byte) (v []byte) {
	d.readNextBd()
	v = v0
	var h decNextValueBytesHelper
	h.append1(&v, d.bytes, d.bd)
	return d.nextValueBytesBdReadR(v)
}

func (d *bincDecDriverBytes) nextValueBytesBdReadR(v0 []byte) (v []byte) {
	v = v0
	var h decNextValueBytesHelper

	fnLen := func(vs byte) uint {
		switch vs {
		case 0:
			x := d.r.readn1()
			h.append1(&v, d.bytes, x)
			return uint(x)
		case 1:
			x := d.r.readn2()
			h.appendN(&v, d.bytes, x[:]...)
			return uint(bigen.Uint16(x))
		case 2:
			x := d.r.readn4()
			h.appendN(&v, d.bytes, x[:]...)
			return uint(bigen.Uint32(x))
		case 3:
			x := d.r.readn8()
			h.appendN(&v, d.bytes, x[:]...)
			return uint(bigen.Uint64(x))
		default:
			return uint(vs - 4)
		}
	}

	var clen uint

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil, bincSpFalse, bincSpTrue, bincSpNan, bincSpPosInf:
		case bincSpNegInf, bincSpZeroFloat, bincSpZero, bincSpNegOne:
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
	case bincVdPosInt, bincVdNegInt:
		bs := d.uintBytes()
		h.appendN(&v, d.bytes, bs...)
	case bincVdFloat:
		fn := func(xlen byte) {
			if d.vs&0x8 != 0 {
				xlen = d.r.readn1()
				h.append1(&v, d.bytes, xlen)
				if xlen > 8 {
					halt.errorf("cannot read float - at most 8 bytes used to represent float - received %v bytes", xlen)
				}
			}
			d.r.readb(d.d.b[:xlen])
			h.appendN(&v, d.bytes, d.d.b[:xlen]...)
		}
		switch d.vs & 0x7 {
		case bincFlBin32:
			fn(4)
		case bincFlBin64:
			fn(8)
		default:
			halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdString, bincVdByteArray:
		clen = fnLen(d.vs)
		h.appendN(&v, d.bytes, d.r.readx(clen)...)
	case bincVdSymbol:
		if d.vs&0x8 == 0 {
			h.append1(&v, d.bytes, d.r.readn1())
		} else {
			h.appendN(&v, d.bytes, d.r.readx(2)...)
		}
		if d.vs&0x4 != 0 {
			clen = fnLen(d.vs & 0x3)
			h.appendN(&v, d.bytes, d.r.readx(clen)...)
		}
	case bincVdTimestamp:
		h.appendN(&v, d.bytes, d.r.readx(uint(d.vs))...)
	case bincVdCustomExt:
		clen = fnLen(d.vs)
		h.append1(&v, d.bytes, d.r.readn1())
		h.appendN(&v, d.bytes, d.r.readx(clen)...)
	case bincVdArray:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			v = d.nextValueBytesR(v)
		}
	case bincVdMap:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			v = d.nextValueBytesR(v)
			v = d.nextValueBytesR(v)
		}
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincEncDriverBytes) init(hh Handle, shared *encoderShared, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*BincHandle)
	d.e = shared
	if shared.bytes {
		fp = bincFpEncBytes
	} else {
		fp = bincFpEncIO
	}

	d.init2(enc)
	return
}

func (e *bincEncDriverBytes) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *bincEncDriverBytes) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *bincEncDriverBytes) writerEnd()                        { e.w.end() }

func (e *bincEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *bincEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *bincDecDriverBytes) init(hh Handle, shared *decoderShared, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*BincHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = bincFpDecBytes
	} else {
		fp = bincFpDecIO
	}

	d.init2(dec)
	return
}

func (d *bincDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *bincDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *bincDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

func (d *bincDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, bincdescbd(d.bd))
}

func (d *bincDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
func (e *encoderBincIO) setContainerState(cs containerState) {
	if cs != 0 {
		e.c = cs
	}
}

func (e *encoderBincIO) rawExt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *encoderBincIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderBincIO) selferMarshal(f *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderBincIO) binaryMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderBincIO) textMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderBincIO) jsonMarshal(f *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderBincIO) raw(f *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderBincIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderBincIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderBincIO) kBool(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderBincIO) kTime(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderBincIO) kString(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderBincIO) kFloat32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderBincIO) kFloat64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderBincIO) kComplex64(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderBincIO) kComplex128(f *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderBincIO) kInt(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderBincIO) kInt8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderBincIO) kInt16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderBincIO) kInt32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderBincIO) kInt64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderBincIO) kUint(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderBincIO) kUint8(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderBincIO) kUint16(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderBincIO) kUint32(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderBincIO) kUint64(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderBincIO) kUintptr(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderBincIO) kErr(f *encFnInfo, rv reflect.Value) {
	e.errorf("unsupported encoding kind %s, for %#v", rv.Kind(), any(rv))
}

func (e *encoderBincIO) kSeqFn(rtelem reflect.Type) (fn *encFnBincIO) {
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}

	if rtelem.Kind() != reflect.Interface {
		fn = e.fn(rtelem)
	}
	return
}

func (e *encoderBincIO) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincIO) kSliceW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincIO) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincIO) kArrayW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderBincIO) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoderBincIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoderBincIO) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderBincIO) kStructSfi(f *encFnInfo) []*structFieldInfo {
	if e.h.Canonical {
		return f.ti.sfi.sorted()
	}
	return f.ti.sfi.source()
}

func (e *encoderBincIO) kStructNoOmitempty(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincIO) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
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

func (e *encoderBincIO) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincIO) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	var keyFn, valFn *encFnBincIO

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

func (e *encoderBincIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnBincIO) {

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

type encoderBincIO struct {
	panicHdl
	perType encPerType

	dh helperEncDriverBincIO

	fp *fastpathEsBincIO

	h *BasicHandle

	e bincEncDriverIO

	encoderShared

	hh Handle

	ci []interface{}

	slist sfiRvFreelist
}

func (e *encoderBincIO) HandleName() string {
	return e.hh.Name()
}

func (e *encoderBincIO) Release() {
}

func (e *encoderBincIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderShared, e).(*fastpathEsBincIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderBincIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderBincIO) MustEncode(v interface{}) {
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

func (e *encoderBincIO) Encode(v interface{}) (err error) {

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

func (e *encoderBincIO) encode(iv interface{}) {

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

func (e *encoderBincIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderBincIO) encodeValue(rv reflect.Value, fn *encFnBincIO) {

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

func (e *encoderBincIO) encodeValueNonNil(rv reflect.Value, fn *encFnBincIO) {
	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if fn.i.addrE {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderBincIO) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
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

func (e *encoderBincIO) marshalUtf8(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderBincIO) marshalAsis(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderBincIO) marshalRaw(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoderBincIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		e.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderBincIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}

func (e *encoderBincIO) fn(t reflect.Type) *encFnBincIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoderBincIO) fnNoExt(t reflect.Type) *encFnBincIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

func (e *encoderBincIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderBincIO) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoderBincIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderBincIO) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoderBincIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderBincIO) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoderBincIO) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

func (e *encoderBincIO) haltOnMbsOddLen(length int) {
	if length&1 != 0 {
		e.errorInt("mapBySlice requires even slice length, but got ", int64(length))
	}
}

func (e *encoderBincIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderBincIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderBincIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderBincIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderBincIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverBincIO) newEncoderBytes(out *[]byte, h Handle) *encoderBincIO {
	var c1 encoderBincIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverBincIO) newEncoderIO(out io.Writer, h Handle) *encoderBincIO {
	var c1 encoderBincIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverBincIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsBincIO) (f *fastpathEBincIO, u reflect.Type) {
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

type encFnBincIO struct {
	i  encFnInfo
	fe func(*encoderBincIO, *encFnInfo, reflect.Value)
}
type encRtidFnBincIO struct {
	rtid uintptr
	fn   *encFnBincIO
}

func (helperEncDriverBincIO) encFindRtidFn(s []encRtidFnBincIO, rtid uintptr) (i uint, fn *encFnBincIO) {

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

func (helperEncDriverBincIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnBincIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnBincIO](v))
	}
	return
}

func (dh helperEncDriverBincIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsBincIO, checkExt bool) (fn *encFnBincIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverBincIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnBincIO
	sp = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverBincIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnBincIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnBincIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnBincIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnBincIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverBincIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincIO) {
	fn = new(encFnBincIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderBincIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderBincIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderBincIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderBincIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fe = (*encoderBincIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderBincIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderBincIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderBincIO).textMarshal
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
					fn.fe = func(e *encoderBincIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderBincIO).kBool
			case reflect.String:

				fn.fe = (*encoderBincIO).kString
			case reflect.Int:
				fn.fe = (*encoderBincIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderBincIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderBincIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderBincIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderBincIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderBincIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderBincIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderBincIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderBincIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderBincIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderBincIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderBincIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderBincIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderBincIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderBincIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderBincIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderBincIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderBincIO).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*encoderBincIO).kStruct
				} else {
					fn.fe = (*encoderBincIO).kStructNoOmitempty
				}
			case reflect.Map:
				fn.fe = (*encoderBincIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderBincIO).kErr
			default:

				fn.fe = (*encoderBincIO).kErr
			}
		}
	}
	return
}

type helperEncDriverBincIO struct{}

func (d *decoderBincIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *decoderBincIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderBincIO) selferUnmarshal(f *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderBincIO) binaryUnmarshal(f *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	d.onerror(fnerr)
}

func (d *decoderBincIO) textUnmarshal(f *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	d.onerror(fnerr)
}

func (d *decoderBincIO) jsonUnmarshal(f *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderBincIO) jsonUnmarshalV(tm jsonUnmarshaler) {

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

func (d *decoderBincIO) kErr(f *decFnInfo, rv reflect.Value) {
	d.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderBincIO) raw(f *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderBincIO) kString(f *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes()))
}

func (d *decoderBincIO) kBool(f *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderBincIO) kTime(f *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderBincIO) kFloat32(f *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderBincIO) kFloat64(f *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderBincIO) kComplex64(f *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderBincIO) kComplex128(f *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderBincIO) kInt(f *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderBincIO) kInt8(f *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderBincIO) kInt16(f *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderBincIO) kInt32(f *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderBincIO) kInt64(f *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderBincIO) kUint(f *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincIO) kUintptr(f *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincIO) kUint8(f *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderBincIO) kUint16(f *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderBincIO) kUint32(f *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderBincIO) kUint64(f *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderBincIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderBincIO) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderBincIO) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *decoderBincIO) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderBincIO) kSlice(f *decFnInfo, rv reflect.Value) {

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

	var fn *decFnBincIO

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

func (d *decoderBincIO) kArray(f *decFnInfo, rv reflect.Value) {

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

	var fn *decFnBincIO

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

func (d *decoderBincIO) kChan(f *decFnInfo, rv reflect.Value) {

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

	var fn *decFnBincIO

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

func (d *decoderBincIO) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnBincIO
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

type decoderBincIO struct {
	panicHdl
	perType decPerType

	dh helperDecDriverBincIO

	fp *fastpathDsBincIO

	h *BasicHandle

	d bincDecDriverIO

	decoderShared

	hh Handle

	mtid uintptr
	stid uintptr
}

func (d *decoderBincIO) HandleName() string {
	return d.hh.Name()
}

func (d *decoderBincIO) isBytes() bool {
	return d.bytes
}

func (d *decoderBincIO) init(h Handle) {
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

	d.fp = d.d.init(h, &d.decoderShared, d).(*fastpathDsBincIO)

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

func (d *decoderBincIO) reset() {
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

func (d *decoderBincIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderBincIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
	return
}

func (d *decoderBincIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderBincIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderBincIO) Decode(v interface{}) (err error) {

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

func (d *decoderBincIO) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderBincIO) Release() {
}

func (d *decoderBincIO) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *decoderBincIO) nextValueBytes(start []byte) []byte {
	return d.d.nextValueBytes(start)
}

func (d *decoderBincIO) decode(iv interface{}) {

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

func (d *decoderBincIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderBincIO) decodeValue(rv reflect.Value, fn *decFnBincIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *decoderBincIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnBincIO) {

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

func (d *decoderBincIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			d.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			d.errorStr("no matching struct field found when decoding stream map with key " + rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderBincIO) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		d.errorf("cannot expand array len during decode from %v to %v", any(sliceLen), any(streamLen))
	}
}

func (d *decoderBincIO) haltAsNotDecodeable(rv reflect.Value) {
	if !rv.IsValid() {
		d.onerror(errCannotDecodeIntoNil)
	}

	if !rv.CanInterface() {
		d.errorf("cannot decode into a value without an interface: %v", rv)
	}
	d.errorf("cannot decode into value of kind: %v, %#v", rv.Kind(), rv2i(rv))
}

func (d *decoderBincIO) depthIncr() {
	d.depth++
	if d.depth >= d.maxdepth {
		d.onerror(errMaxDepthExceeded)
	}
}

func (d *decoderBincIO) depthDecr() {
	d.depth--
}

func (d *decoderBincIO) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = zeroByteSlice
	}
	return d.d.DecodeBytes(in)
}

func (d *decoderBincIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes(zeroByteSlice)
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderBincIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.NumBytesRead(), false)
}

func (d *decoderBincIO) NumBytesRead() int {
	return int(d.d.NumBytesRead())
}

func (d *decoderBincIO) checkBreak() (v bool) {
	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *decoderBincIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *decoderBincIO) mapStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerMapStart
	}
	return v
}

func (d *decoderBincIO) mapElemKey() {
	d.d.ReadMapElemKey()
	d.c = containerMapKey
}

func (d *decoderBincIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderBincIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincIO) arrayStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerArrayStart
	}
	return v
}

func (d *decoderBincIO) arrayElem() {
	d.d.ReadArrayElem()
	d.c = containerArrayElem
}

func (d *decoderBincIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

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

func (d *decoderBincIO) oneShotAddrRV(rvt reflect.Type, rvk reflect.Kind) reflect.Value {
	if decUseTransient &&
		(numBoolStrSliceBitset.isset(byte(rvk)) ||
			((rvk == reflect.Struct || rvk == reflect.Array) &&
				d.h.getTypeInfo(rt2id(rvt), rvt).flagCanTransient)) {
		return d.perType.TransientAddrK(rvt, rvk)
	}
	return rvZeroAddrK(rvt, rvk)
}

func (d *decoderBincIO) fn(t reflect.Type) *decFnBincIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, false)
}

func (d *decoderBincIO) fnNoExt(t reflect.Type) *decFnBincIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, true)
}

type decSliceHelperBincIO struct {
	d     *decoderBincIO
	ct    valueType
	Array bool
	IsNil bool
}

func (d *decoderBincIO) decSliceHelperStart() (x decSliceHelperBincIO, clen int) {
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

func (x decSliceHelperBincIO) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelperBincIO) ElemContainerState(index int) {

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 {
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelperBincIO) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
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

func (helperDecDriverBincIO) newDecoderBytes(in []byte, h Handle) *decoderBincIO {
	var c1 decoderBincIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverBincIO) newDecoderIO(in io.Reader, h Handle) *decoderBincIO {
	var c1 decoderBincIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverBincIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsBincIO) (f *fastpathDBincIO, u reflect.Type) {
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

type decFnBincIO struct {
	i  decFnInfo
	fd func(*decoderBincIO, *decFnInfo, reflect.Value)
}
type decRtidFnBincIO struct {
	rtid uintptr
	fn   *decFnBincIO
}

func (helperDecDriverBincIO) decFindRtidFn(s []decRtidFnBincIO, rtid uintptr) (i uint, fn *decFnBincIO) {

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

func (helperDecDriverBincIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnBincIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnBincIO](v))
	}
	return
}

func (dh helperDecDriverBincIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsBincIO,
	checkExt bool) (fn *decFnBincIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverBincIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnBincIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverBincIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnBincIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnBincIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnBincIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnBincIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverBincIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincIO) {
	fn = new(decFnBincIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderBincIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderBincIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderBincIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderBincIO).ext
		fi.addrD = true
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fd = (*decoderBincIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderBincIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderBincIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderBincIO).textUnmarshal
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
						fn.fd = func(d *decoderBincIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderBincIO, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderBincIO).kBool
			case reflect.String:
				fn.fd = (*decoderBincIO).kString
			case reflect.Int:
				fn.fd = (*decoderBincIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderBincIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderBincIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderBincIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderBincIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderBincIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderBincIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderBincIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderBincIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderBincIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderBincIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderBincIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderBincIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderBincIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderBincIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderBincIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderBincIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderBincIO).kArray
			case reflect.Struct:
				fn.fd = (*decoderBincIO).kStruct
			case reflect.Map:
				fn.fd = (*decoderBincIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderBincIO).kInterface
			default:

				fn.fd = (*decoderBincIO).kErr
			}
		}
	}
	return
}

type helperDecDriverBincIO struct{}
type fastpathEBincIO struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderBincIO, *encFnInfo, reflect.Value)
}
type fastpathDBincIO struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderBincIO, *decFnInfo, reflect.Value)
}
type fastpathEsBincIO [56]fastpathEBincIO
type fastpathDsBincIO [56]fastpathDBincIO
type fastpathETBincIO struct{}
type fastpathDTBincIO struct{}

func (helperEncDriverBincIO) fastpathEList() *fastpathEsBincIO {
	var i uint = 0
	var s fastpathEsBincIO
	fn := func(v interface{}, fe func(*encoderBincIO, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathEBincIO{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderBincIO).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderBincIO).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderBincIO).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderBincIO).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderBincIO).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderBincIO).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderBincIO).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderBincIO).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderBincIO).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderBincIO).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderBincIO).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderBincIO).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderBincIO).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderBincIO).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderBincIO).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderBincIO).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderBincIO).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderBincIO).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderBincIO).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderBincIO).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderBincIO).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderBincIO).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderBincIO).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderBincIO).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderBincIO).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderBincIO).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderBincIO).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderBincIO).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderBincIO).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderBincIO).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderBincIO).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderBincIO).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderBincIO).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderBincIO).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderBincIO).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderBincIO).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderBincIO).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderBincIO).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderBincIO).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderBincIO).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderBincIO).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderBincIO).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderBincIO).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderBincIO).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderBincIO).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderBincIO).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderBincIO).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderBincIO).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderBincIO).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderBincIO).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderBincIO).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderBincIO).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderBincIO).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderBincIO).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderBincIO).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderBincIO).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverBincIO) fastpathDList() *fastpathDsBincIO {
	var i uint = 0
	var s fastpathDsBincIO
	fn := func(v interface{}, fd func(*decoderBincIO, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDBincIO{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderBincIO).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderBincIO).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderBincIO).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderBincIO).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderBincIO).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderBincIO).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderBincIO).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderBincIO).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderBincIO).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderBincIO).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderBincIO).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderBincIO).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderBincIO).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderBincIO).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderBincIO).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderBincIO).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderBincIO).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderBincIO).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderBincIO).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderBincIO).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderBincIO).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderBincIO).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderBincIO).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderBincIO).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderBincIO).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderBincIO).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderBincIO).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderBincIO).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderBincIO).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderBincIO).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderBincIO).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderBincIO).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderBincIO).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderBincIO).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderBincIO).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderBincIO).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderBincIO).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderBincIO).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderBincIO).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderBincIO).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderBincIO).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderBincIO).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderBincIO).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderBincIO).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderBincIO).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderBincIO).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderBincIO).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderBincIO).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderBincIO).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderBincIO).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderBincIO).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderBincIO).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderBincIO).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderBincIO).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderBincIO).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderBincIO).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverBincIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderBincIO) bool {
	var ft fastpathETBincIO
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

func (e *encoderBincIO) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceIntfV(v []interface{}, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.encode(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceIntfV(v []interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceStringV(v []string, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeString(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceStringV(v []string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceBytesV(v [][]byte, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeStringBytesRaw(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceBytesV(v [][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceFloat32V(v []float32, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat32(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceFloat32V(v []float32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceFloat64V(v []float64, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeFloat64(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceFloat64V(v []float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceUint8V(v []uint8, e *encoderBincIO) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETBincIO) EncAsMapSliceUint8V(v []uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceUint64V(v []uint64, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeUint(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceUint64V(v []uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceIntV(v []int, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceIntV(v []int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceInt32V(v []int32, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(int64(v[j]))
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceInt32V(v []int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceInt64V(v []int64, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeInt(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceInt64V(v []int64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceBoolV(v []bool, e *encoderBincIO) {
	e.arrayStart(len(v))
	for j := range v {
		e.arrayElem()
		e.e.EncodeBool(v[j])
	}
	e.arrayEnd()
}
func (fastpathETBincIO) EncAsMapSliceBoolV(v []bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETBincIO) EncMapStringIntfV(v map[string]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETBincIO) EncMapStringStringV(v map[string]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETBincIO) EncMapStringBytesV(v map[string][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETBincIO) EncMapStringUint8V(v map[string]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETBincIO) EncMapStringUint64V(v map[string]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETBincIO) EncMapStringIntV(v map[string]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETBincIO) EncMapStringInt32V(v map[string]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETBincIO) EncMapStringFloat64V(v map[string]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETBincIO) EncMapStringBoolV(v map[string]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETBincIO) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETBincIO) EncMapUint8StringV(v map[uint8]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETBincIO) EncMapUint8BytesV(v map[uint8][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETBincIO) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETBincIO) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETBincIO) EncMapUint8IntV(v map[uint8]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETBincIO) EncMapUint8Int32V(v map[uint8]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETBincIO) EncMapUint8Float64V(v map[uint8]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETBincIO) EncMapUint8BoolV(v map[uint8]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETBincIO) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETBincIO) EncMapUint64StringV(v map[uint64]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETBincIO) EncMapUint64BytesV(v map[uint64][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETBincIO) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETBincIO) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETBincIO) EncMapUint64IntV(v map[uint64]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETBincIO) EncMapUint64Int32V(v map[uint64]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETBincIO) EncMapUint64Float64V(v map[uint64]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETBincIO) EncMapUint64BoolV(v map[uint64]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETBincIO) EncMapIntIntfV(v map[int]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETBincIO) EncMapIntStringV(v map[int]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETBincIO) EncMapIntBytesV(v map[int][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETBincIO) EncMapIntUint8V(v map[int]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETBincIO) EncMapIntUint64V(v map[int]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETBincIO) EncMapIntIntV(v map[int]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETBincIO) EncMapIntInt32V(v map[int]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETBincIO) EncMapIntFloat64V(v map[int]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETBincIO) EncMapIntBoolV(v map[int]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETBincIO) EncMapInt32IntfV(v map[int32]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETBincIO) EncMapInt32StringV(v map[int32]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETBincIO) EncMapInt32BytesV(v map[int32][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETBincIO) EncMapInt32Uint8V(v map[int32]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETBincIO) EncMapInt32Uint64V(v map[int32]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETBincIO) EncMapInt32IntV(v map[int32]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETBincIO) EncMapInt32Int32V(v map[int32]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETBincIO) EncMapInt32Float64V(v map[int32]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
	ft.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETBincIO) EncMapInt32BoolV(v map[int32]bool, e *encoderBincIO) {
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

func (helperDecDriverBincIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderBincIO) bool {
	var ft fastpathDTBincIO
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

func (d *decoderBincIO) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceIntfX(vp *[]interface{}, d *decoderBincIO) {
	if v, changed := f.DecSliceIntfY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceIntfY(v []interface{}, d *decoderBincIO) (v2 []interface{}, changed bool) {
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
func (fastpathDTBincIO) DecSliceIntfN(v []interface{}, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceStringX(vp *[]string, d *decoderBincIO) {
	if v, changed := f.DecSliceStringY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceStringY(v []string, d *decoderBincIO) (v2 []string, changed bool) {
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
func (fastpathDTBincIO) DecSliceStringN(v []string, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceBytesX(vp *[][]byte, d *decoderBincIO) {
	if v, changed := f.DecSliceBytesY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceBytesY(v [][]byte, d *decoderBincIO) (v2 [][]byte, changed bool) {
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
func (fastpathDTBincIO) DecSliceBytesN(v [][]byte, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceFloat32X(vp *[]float32, d *decoderBincIO) {
	if v, changed := f.DecSliceFloat32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceFloat32Y(v []float32, d *decoderBincIO) (v2 []float32, changed bool) {
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
func (fastpathDTBincIO) DecSliceFloat32N(v []float32, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceFloat64X(vp *[]float64, d *decoderBincIO) {
	if v, changed := f.DecSliceFloat64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceFloat64Y(v []float64, d *decoderBincIO) (v2 []float64, changed bool) {
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
func (fastpathDTBincIO) DecSliceFloat64N(v []float64, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceUint8X(vp *[]uint8, d *decoderBincIO) {
	if v, changed := f.DecSliceUint8Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceUint8Y(v []uint8, d *decoderBincIO) (v2 []uint8, changed bool) {
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
func (fastpathDTBincIO) DecSliceUint8N(v []uint8, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceUint64X(vp *[]uint64, d *decoderBincIO) {
	if v, changed := f.DecSliceUint64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceUint64Y(v []uint64, d *decoderBincIO) (v2 []uint64, changed bool) {
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
func (fastpathDTBincIO) DecSliceUint64N(v []uint64, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceIntX(vp *[]int, d *decoderBincIO) {
	if v, changed := f.DecSliceIntY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceIntY(v []int, d *decoderBincIO) (v2 []int, changed bool) {
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
func (fastpathDTBincIO) DecSliceIntN(v []int, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceInt32X(vp *[]int32, d *decoderBincIO) {
	if v, changed := f.DecSliceInt32Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceInt32Y(v []int32, d *decoderBincIO) (v2 []int32, changed bool) {
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
func (fastpathDTBincIO) DecSliceInt32N(v []int32, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceInt64X(vp *[]int64, d *decoderBincIO) {
	if v, changed := f.DecSliceInt64Y(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceInt64Y(v []int64, d *decoderBincIO) (v2 []int64, changed bool) {
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
func (fastpathDTBincIO) DecSliceInt64N(v []int64, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecSliceBoolX(vp *[]bool, d *decoderBincIO) {
	if v, changed := f.DecSliceBoolY(*vp, d); changed {
		*vp = v
	}
}
func (fastpathDTBincIO) DecSliceBoolY(v []bool, d *decoderBincIO) (v2 []bool, changed bool) {
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
func (fastpathDTBincIO) DecSliceBoolN(v []bool, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringIntfX(vp *map[string]interface{}, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringStringX(vp *map[string]string, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringStringL(v map[string]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringBytesX(vp *map[string][]byte, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringUint8X(vp *map[string]uint8, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringUint64X(vp *map[string]uint64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringIntX(vp *map[string]int, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringIntL(v map[string]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringInt32X(vp *map[string]int32, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringFloat64X(vp *map[string]float64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapStringBoolX(vp *map[string]bool, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8IntfX(vp *map[uint8]interface{}, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8StringX(vp *map[uint8]string, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8BytesX(vp *map[uint8][]byte, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8Uint8X(vp *map[uint8]uint8, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8Uint64X(vp *map[uint8]uint64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8IntX(vp *map[uint8]int, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8Int32X(vp *map[uint8]int32, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8Float64X(vp *map[uint8]float64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint8BoolX(vp *map[uint8]bool, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64IntfX(vp *map[uint64]interface{}, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64StringX(vp *map[uint64]string, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64BytesX(vp *map[uint64][]byte, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64Uint8X(vp *map[uint64]uint8, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64Uint64X(vp *map[uint64]uint64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64IntX(vp *map[uint64]int, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64Int32X(vp *map[uint64]int32, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64Float64X(vp *map[uint64]float64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapUint64BoolX(vp *map[uint64]bool, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntIntfX(vp *map[int]interface{}, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntStringX(vp *map[int]string, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntStringL(v map[int]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntBytesX(vp *map[int][]byte, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntUint8X(vp *map[int]uint8, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntUint64X(vp *map[int]uint64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntIntX(vp *map[int]int, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntIntL(v map[int]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntInt32X(vp *map[int]int32, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntFloat64X(vp *map[int]float64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapIntBoolX(vp *map[int]bool, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32IntfX(vp *map[int32]interface{}, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32StringX(vp *map[int32]string, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32BytesX(vp *map[int32][]byte, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32Uint8X(vp *map[int32]uint8, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32Uint64X(vp *map[int32]uint64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32IntX(vp *map[int32]int, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32Int32X(vp *map[int32]int32, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32Float64X(vp *map[int32]float64, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (f fastpathDTBincIO) DecMapInt32BoolX(vp *map[int32]bool, d *decoderBincIO) {
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
func (fastpathDTBincIO) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderBincIO) {
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

type bincEncDriverIO struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT
	encInit2er

	h *BincHandle
	e *encoderShared
	w bufioEncWriter
	bincEncState
}

func (e *bincEncDriverIO) EncodeNil() {
	e.w.writen1(bincBdNil)
}

func (e *bincEncDriverIO) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else {
		bs := bincEncodeTime(t)
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	}
}

func (e *bincEncDriverIO) EncodeBool(b bool) {
	if b {
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	}
}

func (e *bincEncDriverIO) encSpFloat(f float64) (done bool) {
	if f == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
	} else if math.IsNaN(float64(f)) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNan)
	} else if math.IsInf(float64(f), +1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpPosInf)
	} else if math.IsInf(float64(f), -1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegInf)
	} else {
		return
	}
	return true
}

func (e *bincEncDriverIO) EncodeFloat32(f float32) {
	if !e.encSpFloat(float64(f)) {
		e.w.writen1(bincVdFloat<<4 | bincFlBin32)
		e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
	}
}

func (e *bincEncDriverIO) EncodeFloat64(f float64) {
	if e.encSpFloat(f) {
		return
	}
	b := bigen.PutUint64(math.Float64bits(f))
	if bincDoPrune {
		i := 7
		for ; i >= 0 && (b[i] == 0); i-- {
		}
		i++
		if i <= 6 {
			e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.w.writen1(byte(i))
			e.w.writeb(b[:i])
			return
		}
	}
	e.w.writen1(bincVdFloat<<4 | bincFlBin64)
	e.w.writen8(b)
}

func (e *bincEncDriverIO) encIntegerPrune32(bd byte, pos bool, v uint64) {
	b := bigen.PutUint32(uint32(v))
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 3 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 3)
		e.w.writen4(b)
	}
}

func (e *bincEncDriverIO) encIntegerPrune64(bd byte, pos bool, v uint64) {
	b := bigen.PutUint64(v)
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 7 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 7)
		e.w.writen8(b)
	}
}

func (e *bincEncDriverIO) EncodeInt(v int64) {
	if v >= 0 {
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	} else if v == -1 {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	} else {
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	}
}

func (e *bincEncDriverIO) EncodeUint(v uint64) {
	e.encUint(bincVdPosInt<<4, true, v)
}

func (e *bincEncDriverIO) encUint(bd byte, pos bool, v uint64) {
	if v == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	} else if pos && v >= 1 && v <= 16 {
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd|0x0, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.encIntegerPrune32(bd, pos, v)
	} else {
		e.encIntegerPrune64(bd, pos, v)
	}
}

func (e *bincEncDriverIO) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	var bs0, bs []byte
	if ext == SelfExt {
		bs0 = e.e.blist.get(1024)
		bs = bs0
		sideEncode(&e.h.BasicHandle, v, &bs, basetype, true)
	} else {
		bs = ext.WriteExt(v)
	}
	if bs == nil {
		e.EncodeNil()
		goto END
	}
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *bincEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *bincEncDriverIO) encodeExtPreamble(xtag byte, length int) {
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncDriverIO) WriteArrayStart(length int) {
	e.encLen(bincVdArray<<4, uint64(length))
}

func (e *bincEncDriverIO) WriteMapStart(length int) {
	e.encLen(bincVdMap<<4, uint64(length))
}

func (e *bincEncDriverIO) EncodeSymbol(v string) {

	l := len(v)
	if l == 0 {
		e.encBytesLen(cUTF8, 0)
		return
	} else if l == 1 {
		e.encBytesLen(cUTF8, 1)
		e.w.writen1(v[0])
		return
	}
	if e.m == nil {
		e.m = make(map[string]uint16, 16)
	}
	ui, ok := e.m[v]
	if ok {
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			e.w.writen2(bigen.PutUint16(ui))
		}
	} else {
		e.e.seq++
		ui = e.e.seq
		e.m[v] = ui
		var lenprec uint8
		if l <= math.MaxUint8 {

		} else if l <= math.MaxUint16 {
			lenprec = 1
		} else if int64(l) <= math.MaxUint32 {
			lenprec = 2
		} else {
			lenprec = 3
		}
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4|0x0|0x4|lenprec, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			e.w.writen2(bigen.PutUint16(ui))
		}
		if lenprec == 0 {
			e.w.writen1(byte(l))
		} else if lenprec == 1 {
			e.w.writen2(bigen.PutUint16(uint16(l)))
		} else if lenprec == 2 {
			e.w.writen4(bigen.PutUint32(uint32(l)))
		} else {
			e.w.writen8(bigen.PutUint64(uint64(l)))
		}
		e.w.writestr(v)
	}
}

func (e *bincEncDriverIO) EncodeString(v string) {
	if e.h.StringToRaw {
		e.encLen(bincVdByteArray<<4, uint64(len(v)))
		if len(v) > 0 {
			e.w.writestr(v)
		}
		return
	}
	e.EncodeStringEnc(cUTF8, v)
}

func (e *bincEncDriverIO) EncodeStringEnc(c charEncoding, v string) {
	if e.e.c == containerMapKey && c == cUTF8 && (e.h.AsSymbols == 1) {
		e.EncodeSymbol(v)
		return
	}
	e.encLen(bincVdString<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writestr(v)
	}
}

func (e *bincEncDriverIO) EncodeStringBytesRaw(v []byte) {
	if v == nil {
		e.EncodeNil()
		return
	}
	e.encLen(bincVdByteArray<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writeb(v)
	}
}

func (e *bincEncDriverIO) encBytesLen(c charEncoding, length uint64) {

	if c == cRAW {
		e.encLen(bincVdByteArray<<4, length)
	} else {
		e.encLen(bincVdString<<4, length)
	}
}

func (e *bincEncDriverIO) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l+4))
	} else {
		e.encLenNumber(bd, l)
	}
}

func (e *bincEncDriverIO) encLenNumber(bd byte, v uint64) {
	if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd | 0x02)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd | 0x03)
		e.w.writen8(bigen.PutUint64(uint64(v)))
	}
}

type bincDecDriverIO struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	decInit2er
	noBuiltInTypes

	h *BincHandle
	d *decoderShared
	r ioDecReader

	bincDecState

	bytes bool
}

func (d *bincDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
}

func (d *bincDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *bincDecDriverIO) TryNil() bool {
	return d.advanceNil()
}

func (d *bincDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return valueTypeNil
	} else if d.vd == bincVdByteArray {
		return valueTypeBytes
	} else if d.vd == bincVdString {
		return valueTypeString
	} else if d.vd == bincVdArray {
		return valueTypeArray
	} else if d.vd == bincVdMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *bincDecDriverIO) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.vd != bincVdTimestamp {
		halt.errorf("cannot decode time - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	t, err := bincDecodeTime(d.r.readx(uint(d.vs)))
	halt.onerror(err)
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) decFloatPruned(maxlen uint8) {
	l := d.r.readn1()
	if l > maxlen {
		halt.errorf("cannot read float - at most %v bytes used to represent float - received %v bytes", maxlen, l)
	}
	for i := l; i < maxlen; i++ {
		d.d.b[i] = 0
	}
	d.r.readb(d.d.b[0:l])
}

func (d *bincDecDriverIO) decFloatPre32() (b [4]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn4()
	} else {
		d.decFloatPruned(4)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverIO) decFloatPre64() (b [8]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn8()
	} else {
		d.decFloatPruned(8)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverIO) decFloatVal() (f float64) {
	switch d.vs & 0x7 {
	case bincFlBin32:
		f = float64(math.Float32frombits(bigen.Uint32(d.decFloatPre32())))
	case bincFlBin64:
		f = math.Float64frombits(bigen.Uint64(d.decFloatPre64()))
	default:

		halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincDecDriverIO) decUint() (v uint64) {
	switch d.vs {
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		v = uint64(bigen.Uint16(d.r.readn2()))
	case 2:
		b3 := d.r.readn3()
		var b [4]byte
		copy(b[1:], b3[:])
		v = uint64(bigen.Uint32(b))
	case 3:
		v = uint64(bigen.Uint32(d.r.readn4()))
	case 4, 5, 6:
		var b [8]byte
		lim := 7 - d.vs
		bs := d.d.b[lim:8]
		d.r.readb(bs)
		copy(b[lim:], bs)
		v = bigen.Uint64(b)
	case 7:
		v = bigen.Uint64(d.r.readn8())
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	return
}

func (d *bincDecDriverIO) uintBytes() (bs []byte) {
	switch d.vs {
	case 0:
		bs = d.d.b[:1]
		bs[0] = d.r.readn1()
	case 1:
		bs = d.d.b[:2]
		d.r.readb(bs)
	case 2:
		bs = d.d.b[:3]
		d.r.readb(bs)
	case 3:
		bs = d.d.b[:4]
		d.r.readb(bs)
	case 4, 5, 6:
		lim := 7 - d.vs
		bs = d.d.b[lim:8]
		d.r.readb(bs)
	case 7:
		bs = d.d.b[:8]
		d.r.readb(bs)
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	return
}

func (d *bincDecDriverIO) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdPosInt {
		ui = d.decUint()
	} else if vd == bincVdNegInt {
		ui = d.decUint()
		neg = true
	} else if vd == bincVdSmallInt {
		ui = uint64(d.vs) + 1
	} else if vd == bincVdSpecial {
		if vs == bincSpZero {

		} else if vs == bincSpNegOne {
			neg = true
			ui = 1
		} else {
			ok = false

		}
	} else {
		ok = false

	}
	return
}

func (d *bincDecDriverIO) decFloat() (f float64, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdSpecial {
		if vs == bincSpNan {
			f = math.NaN()
		} else if vs == bincSpPosInf {
			f = math.Inf(1)
		} else if vs == bincSpZeroFloat || vs == bincSpZero {

		} else if vs == bincSpNegInf {
			f = math.Inf(-1)
		} else {
			ok = false

		}
	} else if vd == bincVdFloat {
		f = d.decFloatVal()
	} else {
		ok = false
	}
	return
}

func (d *bincDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == (bincVdSpecial | bincSpFalse) {

	} else if d.bd == (bincVdSpecial | bincSpTrue) {
		b = true
	} else {
		halt.errorf("bool - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdMap {
		halt.errorf("map - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdArray {
		halt.errorf("array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) decLen() int {
	if d.vs > 3 {
		return int(d.vs - 4)
	}
	return int(d.decLenNumber())
}

func (d *bincDecDriverIO) decLenNumber() (v uint64) {
	if x := d.vs; x == 0 {
		v = uint64(d.r.readn1())
	} else if x == 1 {
		v = uint64(bigen.Uint16(d.r.readn2()))
	} else if x == 2 {
		v = uint64(bigen.Uint32(d.r.readn4()))
	} else {
		v = bigen.Uint64(d.r.readn8())
	}
	return
}

func (d *bincDecDriverIO) DecodeStringAsBytes() (bs2 []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	var slen = -1
	switch d.vd {
	case bincVdString, bincVdByteArray:
		slen = d.decLen()
		if d.d.bytes {
			d.d.decByteState = decByteStateZerocopy
			bs2 = d.r.readx(uint(slen))
		} else {
			d.d.decByteState = decByteStateReuseBuf
			bs2 = decByteSlice(&d.r, slen, d.h.MaxInitLen, d.d.b[:])
		}
	case bincVdSymbol:

		var symbol uint16
		vs := d.vs
		if vs&0x8 == 0 {
			symbol = uint16(d.r.readn1())
		} else {
			symbol = uint16(bigen.Uint16(d.r.readn2()))
		}
		if d.s == nil {
			d.s = make(map[uint16][]byte, 16)
		}

		if vs&0x4 == 0 {
			bs2 = d.s[symbol]
		} else {
			switch vs & 0x3 {
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(bigen.Uint16(d.r.readn2()))
			case 2:
				slen = int(bigen.Uint32(d.r.readn4()))
			case 3:
				slen = int(bigen.Uint64(d.r.readn8()))
			}

			bs2 = decByteSlice(&d.r, slen, d.h.MaxInitLen, nil)
			d.s[symbol] = bs2
		}
	default:
		halt.errorf("string/bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if d.h.ValidateUnicode && !utf8.Valid(bs2) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", bs2)
	}

	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeBytes(bs []byte) (bsOut []byte) {
	d.d.decByteState = decByteStateNone
	if d.advanceNil() {
		return
	}
	if d.vd == bincVdArray {
		if bs == nil {
			bs = d.d.b[:]
			d.d.decByteState = decByteStateReuseBuf
		}
		slen := d.ReadArrayStart()
		var changed bool
		if bs, changed = usableByteSlice(bs, slen); changed {
			d.d.decByteState = decByteStateNone
		}
		for i := 0; i < slen; i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		return bs
	}
	var clen int
	if d.vd == bincVdString || d.vd == bincVdByteArray {
		clen = d.decLen()
	} else {
		halt.errorf("bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	if d.bytes && d.h.ZeroCopy {
		d.d.decByteState = decByteStateZerocopy
		return d.r.readx(uint(clen))
	}
	if bs == nil {
		bs = d.d.b[:]
		d.d.decByteState = decByteStateReuseBuf
	}
	return decByteSlice(&d.r, clen, d.h.MaxInitLen, bs)
}

func (d *bincDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if xtag > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtag)
	}
	if d.advanceNil() {
		return
	}
	xbs, realxtag1, zerocopy := d.decodeExtV(ext != nil, uint8(xtag))
	realxtag := uint64(realxtag1)
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = realxtag
		re.setData(xbs, zerocopy)
	} else if ext == SelfExt {
		sideDecode(&d.h.BasicHandle, rv, xbs, basetype, true)
	} else {
		ext.ReadExt(rv, xbs)
	}
}

func (d *bincDecDriverIO) decodeExtV(verifyTag bool, tag byte) (xbs []byte, xtag byte, zerocopy bool) {
	if d.vd == bincVdCustomExt {
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag - got %b, expecting: %v", xtag, tag)
		}
		if d.d.bytes {
			xbs = d.r.readx(uint(l))
			zerocopy = true
		} else {
			xbs = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	} else if d.vd == bincVdByteArray {
		xbs = d.DecodeBytes(nil)
	} else {
		halt.errorf("ext expects extensions or byte array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil:
			n.v = valueTypeNil
		case bincSpFalse:
			n.v = valueTypeBool
			n.b = false
		case bincSpTrue:
			n.v = valueTypeBool
			n.b = true
		case bincSpNan:
			n.v = valueTypeFloat
			n.f = math.NaN()
		case bincSpPosInf:
			n.v = valueTypeFloat
			n.f = math.Inf(1)
		case bincSpNegInf:
			n.v = valueTypeFloat
			n.f = math.Inf(-1)
		case bincSpZeroFloat:
			n.v = valueTypeFloat
			n.f = float64(0)
		case bincSpZero:
			n.v = valueTypeUint
			n.u = uint64(0)
		case bincSpNegOne:
			n.v = valueTypeInt
			n.i = int64(-1)
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
		n.v = valueTypeUint
		n.u = uint64(int8(d.vs)) + 1
	case bincVdPosInt:
		n.v = valueTypeUint
		n.u = d.decUint()
	case bincVdNegInt:
		n.v = valueTypeInt
		n.i = -(int64(d.decUint()))
	case bincVdFloat:
		n.v = valueTypeFloat
		n.f = d.decFloatVal()
	case bincVdString:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case bincVdByteArray:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString, d.h.ZeroCopy)
	case bincVdSymbol:
		n.v = valueTypeSymbol
		n.s = d.d.stringZC(d.DecodeStringAsBytes())
	case bincVdTimestamp:
		n.v = valueTypeTime
		tt, err := bincDecodeTime(d.r.readx(uint(d.vs)))
		halt.onerror(err)
		n.t = tt
	case bincVdCustomExt:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		if d.d.bytes {
			n.l = d.r.readx(uint(l))
		} else {
			n.l = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case bincVdArray:
		n.v = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if !decodeFurther {
		d.bdRead = false
	}
	if n.v == valueTypeUint && d.h.SignedInteger {
		n.v = valueTypeInt
		n.i = int64(n.u)
	}
}

func (d *bincDecDriverIO) nextValueBytes(v0 []byte) (v []byte) {
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

func (d *bincDecDriverIO) nextValueBytesR(v0 []byte) (v []byte) {
	d.readNextBd()
	v = v0
	var h decNextValueBytesHelper
	h.append1(&v, d.bytes, d.bd)
	return d.nextValueBytesBdReadR(v)
}

func (d *bincDecDriverIO) nextValueBytesBdReadR(v0 []byte) (v []byte) {
	v = v0
	var h decNextValueBytesHelper

	fnLen := func(vs byte) uint {
		switch vs {
		case 0:
			x := d.r.readn1()
			h.append1(&v, d.bytes, x)
			return uint(x)
		case 1:
			x := d.r.readn2()
			h.appendN(&v, d.bytes, x[:]...)
			return uint(bigen.Uint16(x))
		case 2:
			x := d.r.readn4()
			h.appendN(&v, d.bytes, x[:]...)
			return uint(bigen.Uint32(x))
		case 3:
			x := d.r.readn8()
			h.appendN(&v, d.bytes, x[:]...)
			return uint(bigen.Uint64(x))
		default:
			return uint(vs - 4)
		}
	}

	var clen uint

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil, bincSpFalse, bincSpTrue, bincSpNan, bincSpPosInf:
		case bincSpNegInf, bincSpZeroFloat, bincSpZero, bincSpNegOne:
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
	case bincVdPosInt, bincVdNegInt:
		bs := d.uintBytes()
		h.appendN(&v, d.bytes, bs...)
	case bincVdFloat:
		fn := func(xlen byte) {
			if d.vs&0x8 != 0 {
				xlen = d.r.readn1()
				h.append1(&v, d.bytes, xlen)
				if xlen > 8 {
					halt.errorf("cannot read float - at most 8 bytes used to represent float - received %v bytes", xlen)
				}
			}
			d.r.readb(d.d.b[:xlen])
			h.appendN(&v, d.bytes, d.d.b[:xlen]...)
		}
		switch d.vs & 0x7 {
		case bincFlBin32:
			fn(4)
		case bincFlBin64:
			fn(8)
		default:
			halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdString, bincVdByteArray:
		clen = fnLen(d.vs)
		h.appendN(&v, d.bytes, d.r.readx(clen)...)
	case bincVdSymbol:
		if d.vs&0x8 == 0 {
			h.append1(&v, d.bytes, d.r.readn1())
		} else {
			h.appendN(&v, d.bytes, d.r.readx(2)...)
		}
		if d.vs&0x4 != 0 {
			clen = fnLen(d.vs & 0x3)
			h.appendN(&v, d.bytes, d.r.readx(clen)...)
		}
	case bincVdTimestamp:
		h.appendN(&v, d.bytes, d.r.readx(uint(d.vs))...)
	case bincVdCustomExt:
		clen = fnLen(d.vs)
		h.append1(&v, d.bytes, d.r.readn1())
		h.appendN(&v, d.bytes, d.r.readx(clen)...)
	case bincVdArray:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			v = d.nextValueBytesR(v)
		}
	case bincVdMap:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			v = d.nextValueBytesR(v)
			v = d.nextValueBytesR(v)
		}
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincEncDriverIO) init(hh Handle, shared *encoderShared, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*BincHandle)
	d.e = shared
	if shared.bytes {
		fp = bincFpEncBytes
	} else {
		fp = bincFpEncIO
	}

	d.init2(enc)
	return
}

func (e *bincEncDriverIO) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *bincEncDriverIO) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *bincEncDriverIO) writerEnd()                        { e.w.end() }

func (e *bincEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *bincEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *bincDecDriverIO) init(hh Handle, shared *decoderShared, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*BincHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = bincFpDecBytes
	} else {
		fp = bincFpDecIO
	}

	d.init2(dec)
	return
}

func (d *bincDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *bincDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *bincDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

func (d *bincDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, bincdescbd(d.bd))
}

func (d *bincDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
