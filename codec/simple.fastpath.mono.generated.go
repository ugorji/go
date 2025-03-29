//go:build !codec.notmono && !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

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
)

func (e *encoderSimpleBytes) rawExt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *encoderSimpleBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderSimpleBytes) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderSimpleBytes) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderSimpleBytes) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderSimpleBytes) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderSimpleBytes) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderSimpleBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderSimpleBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderSimpleBytes) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderSimpleBytes) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderSimpleBytes) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderSimpleBytes) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderSimpleBytes) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderSimpleBytes) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderSimpleBytes) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderSimpleBytes) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderSimpleBytes) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderSimpleBytes) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderSimpleBytes) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderSimpleBytes) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderSimpleBytes) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderSimpleBytes) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderSimpleBytes) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderSimpleBytes) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderSimpleBytes) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderSimpleBytes) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderSimpleBytes) kSeqFn(rt reflect.Type) (fn *encFnSimpleBytes) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderSimpleBytes) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleBytes) kSliceW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleBytes) kArrayW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleBytes) kChan(f *encFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		halt.errorStr("send-only channel cannot be encoded")
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

func (e *encoderSimpleBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoderSimpleBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoderSimpleBytes) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderSimpleBytes) kStructNoOmitempty(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderSimpleBytes) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
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

func (e *encoderSimpleBytes) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderSimpleBytes) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	var keyFn, valFn *encFnSimpleBytes

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

func (e *encoderSimpleBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnSimpleBytes) {

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

type encoderSimpleBytes struct {
	dh helperEncDriverSimpleBytes
	fp *fastpathEsSimpleBytes
	e  simpleEncDriverBytes
	encoderBase
}

func (e *encoderSimpleBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsSimpleBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderSimpleBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderSimpleBytes) MustEncode(v interface{}) {
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

func (e *encoderSimpleBytes) Encode(v interface{}) (err error) {

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

func (e *encoderSimpleBytes) encode(iv interface{}) {

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

func (e *encoderSimpleBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderSimpleBytes) encodeValue(rv reflect.Value, fn *encFnSimpleBytes) {

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
					halt.errorf("circular reference found: %p, %T", sptr, sptr)
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

func (e *encoderSimpleBytes) encodeValueNonNil(rv reflect.Value, fn *encFnSimpleBytes) {
	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if fn.i.addrE {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderSimpleBytes) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderSimpleBytes) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderSimpleBytes) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoderSimpleBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderSimpleBytes) fn(t reflect.Type) *encFnSimpleBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoderSimpleBytes) fnNoExt(t reflect.Type) *encFnSimpleBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

func (e *encoderSimpleBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderSimpleBytes) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoderSimpleBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderSimpleBytes) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoderSimpleBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderSimpleBytes) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoderSimpleBytes) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

func (e *encoderSimpleBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderSimpleBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderSimpleBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderSimpleBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderSimpleBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

type encFnSimpleBytes struct {
	i  encFnInfo
	fe func(*encoderSimpleBytes, *encFnInfo, reflect.Value)
}
type encRtidFnSimpleBytes struct {
	rtid uintptr
	fn   *encFnSimpleBytes
}

func (helperEncDriverSimpleBytes) newEncoderBytes(out *[]byte, h Handle) *encoderSimpleBytes {
	var c1 encoderSimpleBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverSimpleBytes) newEncoderIO(out io.Writer, h Handle) *encoderSimpleBytes {
	var c1 encoderSimpleBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverSimpleBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsSimpleBytes) (f *fastpathESimpleBytes, u reflect.Type) {
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

func (helperEncDriverSimpleBytes) encFindRtidFn(s []encRtidFnSimpleBytes, rtid uintptr) (i uint, fn *encFnSimpleBytes) {

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

func (helperEncDriverSimpleBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnSimpleBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnSimpleBytes](v))
	}
	return
}

func (dh helperEncDriverSimpleBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsSimpleBytes, checkExt bool) (fn *encFnSimpleBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverSimpleBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnSimpleBytes = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverSimpleBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnSimpleBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnSimpleBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnSimpleBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnSimpleBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverSimpleBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleBytes) {
	fn = new(encFnSimpleBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderSimpleBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderSimpleBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderSimpleBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderSimpleBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fe = (*encoderSimpleBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderSimpleBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderSimpleBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderSimpleBytes).textMarshal
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
					fn.fe = func(e *encoderSimpleBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderSimpleBytes).kBool
			case reflect.String:

				fn.fe = (*encoderSimpleBytes).kString
			case reflect.Int:
				fn.fe = (*encoderSimpleBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderSimpleBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderSimpleBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderSimpleBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderSimpleBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderSimpleBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderSimpleBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderSimpleBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderSimpleBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderSimpleBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderSimpleBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderSimpleBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderSimpleBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderSimpleBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderSimpleBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderSimpleBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderSimpleBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderSimpleBytes).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*encoderSimpleBytes).kStruct
				} else {
					fn.fe = (*encoderSimpleBytes).kStructNoOmitempty
				}
			case reflect.Map:
				fn.fe = (*encoderSimpleBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderSimpleBytes).kErr
			default:

				fn.fe = (*encoderSimpleBytes).kErr
			}
		}
	}
	return
}

type helperEncDriverSimpleBytes struct{}

func (d *decoderSimpleBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *decoderSimpleBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderSimpleBytes) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderSimpleBytes) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderSimpleBytes) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOk(d.d.DecodeStringAsBytes(nil)))
	halt.onerror(fnerr)
}

func (d *decoderSimpleBytes) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderSimpleBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

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
	halt.onerror(fnerr)
}

func (d *decoderSimpleBytes) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderSimpleBytes) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderSimpleBytes) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes(zeroByteSlice)))
}

func (d *decoderSimpleBytes) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderSimpleBytes) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderSimpleBytes) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderSimpleBytes) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderSimpleBytes) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderSimpleBytes) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderSimpleBytes) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderSimpleBytes) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderSimpleBytes) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderSimpleBytes) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderSimpleBytes) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderSimpleBytes) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleBytes) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleBytes) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderSimpleBytes) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderSimpleBytes) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderSimpleBytes) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderSimpleBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

	n := d.naked()
	d.d.DecodeNaked()

	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		halt.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
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

func (d *decoderSimpleBytes) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderSimpleBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *decoderSimpleBytes) kStruct(f *decFnInfo, rv reflect.Value) {
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
			name2 = make([]byte, 0, 16)
		}
		var rvkencname []byte
		tkt := ti.keyType
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey()

			if tkt == valueTypeString {
				rvkencname, _ = d.d.DecodeStringAsBytes(nil)
			} else if tkt == valueTypeInt {
				rvkencname = strconv.AppendInt(d.b[:0], d.d.DecodeInt64(), 10)
			} else if tkt == valueTypeUint {
				rvkencname = strconv.AppendUint(d.b[:0], d.d.DecodeUint64(), 10)
			} else if tkt == valueTypeFloat {
				rvkencname = strconv.AppendFloat(d.b[:0], d.d.DecodeFloat64(), 'f', -1, 64)
			} else {
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
					halt.errorStr2("no matching struct field when decoding stream map with key: ", stringView(name2))
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
		halt.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderSimpleBytes) kSlice(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
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

	var fn *decFnSimpleBytes

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
				halt.errorStr("cannot decode into non-settable slice")
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
					halt.errorStr("cannot decode into non-settable slice")
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
					halt.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					halt.onerror(errExpandSliceCannotChange)
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

func (d *decoderSimpleBytes) kArray(f *decFnInfo, rv reflect.Value) {

	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {

		if f.ti.elemkind != uint8(reflect.Uint8) {
			halt.errorf("bytes/string in stream can decode into array of bytes, but not %v", f.ti.rt)
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

	var fn *decFnSimpleBytes

	var rv9 reflect.Value

	rvlen := rv.Len()
	hasLen := containerLenS > 0
	if hasLen && containerLenS > rvlen {
		halt.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
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

func (d *decoderSimpleBytes) kChan(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		halt.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2, _ := d.d.DecodeBytes(nil)
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

	var fn *decFnSimpleBytes

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
					halt.errorStr("cannot decode into non-settable chan")
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

func (d *decoderSimpleBytes) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnSimpleBytes
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
	var scratchBuf bool

	fnRvk2 := func() (s string) {
		callFnRvk = false

		if len(kstr2bs) == 1 {
			s = str4byte(kstr2bs[0])
		} else if len(kstr2bs) != 0 {
			s, callFnRvk = d.mapKeyString(&kstrbs, &kstr2bs, scratchBuf)
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
			kstr2bs, scratchBuf = d.d.DecodeStringAsBytes(nil)
			rvSetString(rvk, fnRvk2())
		} else {
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

type decoderSimpleBytes struct {
	dh helperDecDriverSimpleBytes
	fp *fastpathDsSimpleBytes
	d  simpleDecDriverBytes
	decoderBase
}

func (d *decoderSimpleBytes) init(h Handle) {
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

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsSimpleBytes)

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

func (d *decoderSimpleBytes) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
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

func (d *decoderSimpleBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderSimpleBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderSimpleBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderSimpleBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderSimpleBytes) Decode(v interface{}) (err error) {

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

func (d *decoderSimpleBytes) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderSimpleBytes) Release() {
}

func (d *decoderSimpleBytes) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *decoderSimpleBytes) nextValueBytes(start []byte) []byte {
	return d.d.nextValueBytes(start)
}

func (d *decoderSimpleBytes) decode(iv interface{}) {

	if iv == nil {
		halt.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {

	case reflect.Value:
		if x, _ := isDecodeable(v); !x {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)
	case *string:
		*v = d.stringZC(d.d.DecodeStringAsBytes(nil))
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

func (d *decoderSimpleBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderSimpleBytes) decodeValue(rv reflect.Value, fn *decFnSimpleBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *decoderSimpleBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnSimpleBytes) {

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
			halt.errorStr("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *decoderSimpleBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderSimpleBytes) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = zeroByteSlice
	}
	v, _ = d.d.DecodeBytes(in)
	return
}

func (d *decoderSimpleBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes(zeroByteSlice)
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderSimpleBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderSimpleBytes) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderSimpleBytes) checkBreak() (v bool) {
	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *decoderSimpleBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *decoderSimpleBytes) mapElemKey() {
	d.d.ReadMapElemKey()
	d.c = containerMapKey
}

func (d *decoderSimpleBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderSimpleBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleBytes) arrayElem() {
	d.d.ReadArrayElem()
	d.c = containerArrayElem
}

func (d *decoderSimpleBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {
	rv := d.interfaceExtConvertAndDecodeGetRV(v, ext)
	d.decodeValue(rv, nil)
	ext.UpdateExt(v, rv2i(rv))
}

func (d *decoderSimpleBytes) fn(t reflect.Type) *decFnSimpleBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, false)
}

func (d *decoderSimpleBytes) fnNoExt(t reflect.Type) *decFnSimpleBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, true)
}

type decSliceHelperSimpleBytes struct {
	d     *decoderSimpleBytes
	ct    valueType
	Array bool
	IsNil bool
}

func (d *decoderSimpleBytes) decSliceHelperStart() (x decSliceHelperSimpleBytes, clen int) {
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
		halt.errorStr2("to decode into a slice, expect map/array - got ", x.ct.String())
	}
	return
}

func (x decSliceHelperSimpleBytes) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelperSimpleBytes) ElemContainerState(index int) {

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 {
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelperSimpleBytes) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
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

type decFnSimpleBytes struct {
	i  decFnInfo
	fd func(*decoderSimpleBytes, *decFnInfo, reflect.Value)
}
type decRtidFnSimpleBytes struct {
	rtid uintptr
	fn   *decFnSimpleBytes
}

func (helperDecDriverSimpleBytes) newDecoderBytes(in []byte, h Handle) *decoderSimpleBytes {
	var c1 decoderSimpleBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverSimpleBytes) newDecoderIO(in io.Reader, h Handle) *decoderSimpleBytes {
	var c1 decoderSimpleBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverSimpleBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsSimpleBytes) (f *fastpathDSimpleBytes, u reflect.Type) {
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

func (helperDecDriverSimpleBytes) decFindRtidFn(s []decRtidFnSimpleBytes, rtid uintptr) (i uint, fn *decFnSimpleBytes) {

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

func (helperDecDriverSimpleBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnSimpleBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnSimpleBytes](v))
	}
	return
}

func (dh helperDecDriverSimpleBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsSimpleBytes,
	checkExt bool) (fn *decFnSimpleBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverSimpleBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnSimpleBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverSimpleBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnSimpleBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnSimpleBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnSimpleBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnSimpleBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverSimpleBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleBytes) {
	fn = new(decFnSimpleBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderSimpleBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderSimpleBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderSimpleBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderSimpleBytes).ext
		fi.addrD = true
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fd = (*decoderSimpleBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderSimpleBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderSimpleBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderSimpleBytes).textUnmarshal
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
						fn.fd = func(d *decoderSimpleBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderSimpleBytes, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderSimpleBytes).kBool
			case reflect.String:
				fn.fd = (*decoderSimpleBytes).kString
			case reflect.Int:
				fn.fd = (*decoderSimpleBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderSimpleBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderSimpleBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderSimpleBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderSimpleBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderSimpleBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderSimpleBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderSimpleBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderSimpleBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderSimpleBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderSimpleBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderSimpleBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderSimpleBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderSimpleBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderSimpleBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderSimpleBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderSimpleBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderSimpleBytes).kArray
			case reflect.Struct:
				fn.fd = (*decoderSimpleBytes).kStruct
			case reflect.Map:
				fn.fd = (*decoderSimpleBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderSimpleBytes).kInterface
			default:

				fn.fd = (*decoderSimpleBytes).kErr
			}
		}
	}
	return
}

type helperDecDriverSimpleBytes struct{}
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

type simpleEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *SimpleHandle
	e *encoderBase

	w bytesEncAppender
}

func (e *simpleEncDriverBytes) EncodeNil() {
	e.w.writen1(simpleVdNil)
}

func (e *simpleEncDriverBytes) EncodeBool(b bool) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && !b {
		e.EncodeNil()
		return
	}
	if b {
		e.w.writen1(simpleVdTrue)
	} else {
		e.w.writen1(simpleVdFalse)
	}
}

func (e *simpleEncDriverBytes) EncodeFloat32(f float32) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat32)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *simpleEncDriverBytes) EncodeFloat64(f float64) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *simpleEncDriverBytes) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), simpleVdNegInt)
	} else {
		e.encUint(uint64(v), simpleVdPosInt)
	}
}

func (e *simpleEncDriverBytes) EncodeUint(v uint64) {
	e.encUint(v, simpleVdPosInt)
}

func (e *simpleEncDriverBytes) encUint(v uint64, bd uint8) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == 0 {
		e.EncodeNil()
		return
	}
	if v <= math.MaxUint8 {
		e.w.writen2(bd, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 1)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 2)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 3)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *simpleEncDriverBytes) encLen(bd byte, length int) {
	if length == 0 {
		e.w.writen1(bd)
	} else if length <= math.MaxUint8 {
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	} else if length <= math.MaxUint16 {
		e.w.writen1(bd + 2)
		e.w.writen2(bigen.PutUint16(uint16(length)))
	} else if int64(length) <= math.MaxUint32 {
		e.w.writen1(bd + 3)
		e.w.writen4(bigen.PutUint32(uint32(length)))
	} else {
		e.w.writen1(bd + 4)
		e.w.writen8(bigen.PutUint64(uint64(length)))
	}
}

func (e *simpleEncDriverBytes) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *simpleEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *simpleEncDriverBytes) encodeExtPreamble(xtag byte, length int) {
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
}

func (e *simpleEncDriverBytes) WriteArrayStart(length int) {
	e.encLen(simpleVdArray, length)
}

func (e *simpleEncDriverBytes) WriteMapStart(length int) {
	e.encLen(simpleVdMap, length)
}

func (e *simpleEncDriverBytes) EncodeString(v string) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == "" {
		e.EncodeNil()
		return
	}
	if e.h.StringToRaw {
		e.encLen(simpleVdByteArray, len(v))
	} else {
		e.encLen(simpleVdString, len(v))
	}
	e.w.writestr(v)
}

func (e *simpleEncDriverBytes) EncodeStringBytesRaw(v []byte) {

	if v == nil {
		e.EncodeNil()
		return
	}
	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
}

func (e *simpleEncDriverBytes) EncodeTime(t time.Time) {

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	v, err := t.MarshalBinary()
	halt.onerror(err)
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
}

type simpleDecDriverBytes struct {
	h *SimpleHandle
	d *decoderBase
	r bytesDecReader

	bdAndBdread
	bytes bool

	noBuiltInTypes

	decDriverNoopContainerReader
	decInit2er
}

func (d *simpleDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *simpleDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *simpleDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	switch d.bd {
	case simpleVdNil:
		d.bdRead = false
		return valueTypeNil
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		return valueTypeBytes
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		return valueTypeString
	case simpleVdArray, simpleVdArray + 1,
		simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		return valueTypeArray
	case simpleVdMap, simpleVdMap + 1,
		simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *simpleDecDriverBytes) TryNil() bool {
	return d.advanceNil()
}

func (d *simpleDecDriverBytes) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case simpleVdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *simpleDecDriverBytes) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case simpleVdPosInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case simpleVdPosInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
		neg = true
	default:
		ok = false

	}

	return
}

func (d *simpleDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == simpleVdFalse {
	} else if d.bd == simpleVdTrue {
		b = true
	} else {
		halt.errorf("cannot decode bool - %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverBytes) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverBytes) uint2Len(ui uint64) int {
	if chkOvf.Uint(ui, intBitsize) {
		halt.errorf("overflow integer: %v", ui)
	}
	return int(ui)
}

func (d *simpleDecDriverBytes) decLen() int {
	switch d.bd & 7 {
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(bigen.Uint16(d.r.readn2()))
	case 3:
		return d.uint2Len(uint64(bigen.Uint32(d.r.readn4())))
	case 4:
		return d.uint2Len(bigen.Uint64(d.r.readn8()))
	}
	halt.errorf("cannot read length: bd%%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
}

func (d *simpleDecDriverBytes) DecodeStringAsBytes(bs []byte) (s []byte, scratchBuf bool) {
	return d.DecodeBytes(bs)
}

func (d *simpleDecDriverBytes) DecodeBytes(bs []byte) (out []byte, scratchBuf bool) {
	if d.advanceNil() {
		return
	}

	if d.bd >= simpleVdArray && d.bd <= simpleVdMap+4 {
		if bs == nil {
			scratchBuf = true
			bs = d.d.b[:]
		}
		slen := d.ReadArrayStart()
		bs, _ = usableByteSlice(bs, slen)
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		out = bs
		return
	}

	clen := d.decLen()
	d.bdRead = false

	if d.bytes && d.h.ZeroCopy {
		return d.r.readx(uint(clen)), false
	}
	if bs == nil {
		scratchBuf = true
		bs = d.d.b[:]
	}
	out = decByteSlice(&d.r, clen, d.h.MaxInitLen, bs)
	return
}

func (d *simpleDecDriverBytes) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd != simpleVdTime {
		halt.errorf("invalid descriptor for time.Time - expect 0x%x, received 0x%x", simpleVdTime, d.bd)
	}
	d.bdRead = false
	clen := uint(d.r.readn1())
	b := d.r.readx(clen)
	halt.onerror((&t).UnmarshalBinary(b))
	return
}

func (d *simpleDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *simpleDecDriverBytes) decodeExtV(verifyTag bool, tag byte) (xbs []byte, xtag byte, zerocopy bool) {
	switch d.bd {
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		if d.d.bytes {
			xbs = d.r.readx(uint(l))
			zerocopy = true
		} else {
			xbs = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs, _ = d.DecodeBytes(nil)
	default:
		halt.errorf("ext - %s - expecting extensions/bytearray, got: 0x%x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd {
	case simpleVdNil:
		n.v = valueTypeNil
	case simpleVdFalse:
		n.v = valueTypeBool
		n.b = false
	case simpleVdTrue:
		n.v = valueTypeBool
		n.b = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case simpleVdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdTime:
		n.v = valueTypeTime
		n.t = d.DecodeTime()
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes(zeroByteSlice))
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		if d.d.bytes {
			n.l = d.r.readx(uint(l))
		} else {
			n.l = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2,
		simpleVdArray + 3, simpleVdArray + 4:
		n.v = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, d.bd)
	}

	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *simpleDecDriverBytes) nextValueBytes(v0 []byte) (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	v0 = append(v0, d.bd)
	d.r.startRecording(v0)
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) nextValueBytesBdReadR() {
	c := d.bd

	var length uint

	switch c {
	case simpleVdNil, simpleVdFalse, simpleVdTrue, simpleVdString, simpleVdByteArray:

	case simpleVdPosInt, simpleVdNegInt:
		d.r.readn1()
	case simpleVdPosInt + 1, simpleVdNegInt + 1:
		d.r.readx(2)
	case simpleVdPosInt + 2, simpleVdNegInt + 2, simpleVdFloat32:
		d.r.readx(4)
	case simpleVdPosInt + 3, simpleVdNegInt + 3, simpleVdFloat64:
		d.r.readx(8)
	case simpleVdTime:
		c = d.r.readn1()
		d.r.readx(uint(c))

	default:
		switch c & 7 {
		case 0:
			length = 0
		case 1:
			b := d.r.readn1()
			length = uint(b)
		case 2:
			x := d.r.readn2()
			length = uint(bigen.Uint16(x))
		case 3:
			x := d.r.readn4()
			length = uint(bigen.Uint32(x))
		case 4:
			x := d.r.readn8()
			length = uint(bigen.Uint64(x))
		}

		bExt := c >= simpleVdExt && c <= simpleVdExt+7
		bStr := c >= simpleVdString && c <= simpleVdString+7
		bByteArray := c >= simpleVdByteArray && c <= simpleVdByteArray+7
		bArray := c >= simpleVdArray && c <= simpleVdArray+7
		bMap := c >= simpleVdMap && c <= simpleVdMap+7

		if !(bExt || bStr || bByteArray || bArray || bMap) {
			halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, c)
		}

		if bExt {
			d.r.readn1()
		}

		if length == 0 {
			break
		}

		if bArray {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else if bMap {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			d.r.readx(length)
		}
	}
	return
}

func (d *simpleEncDriverBytes) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*SimpleHandle)
	d.e = shared
	if shared.bytes {
		fp = simpleFpEncBytes
	} else {
		fp = simpleFpEncIO
	}

	d.init2(enc)
	return
}

func (e *simpleEncDriverBytes) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *simpleEncDriverBytes) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *simpleEncDriverBytes) writerEnd()                        { e.w.end() }

func (e *simpleEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *simpleEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *simpleDecDriverBytes) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*SimpleHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = simpleFpDecBytes
	} else {
		fp = simpleFpDecIO
	}

	d.init2(dec)
	return
}

func (d *simpleDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *simpleDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *simpleDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

func (d *simpleDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, simpledesc(d.bd))
}

func (d *simpleDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
func (e *encoderSimpleIO) rawExt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *encoderSimpleIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderSimpleIO) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderSimpleIO) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderSimpleIO) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderSimpleIO) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderSimpleIO) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderSimpleIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderSimpleIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderSimpleIO) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderSimpleIO) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderSimpleIO) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderSimpleIO) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderSimpleIO) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderSimpleIO) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderSimpleIO) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderSimpleIO) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderSimpleIO) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderSimpleIO) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderSimpleIO) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderSimpleIO) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderSimpleIO) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderSimpleIO) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderSimpleIO) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderSimpleIO) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderSimpleIO) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderSimpleIO) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderSimpleIO) kSeqFn(rt reflect.Type) (fn *encFnSimpleIO) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderSimpleIO) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleIO) kSliceW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleIO) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleIO) kArrayW(rv reflect.Value, ti *typeInfo) {
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

func (e *encoderSimpleIO) kChan(f *encFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		halt.errorStr("send-only channel cannot be encoded")
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

func (e *encoderSimpleIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoderSimpleIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoderSimpleIO) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderSimpleIO) kStructNoOmitempty(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderSimpleIO) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
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

func (e *encoderSimpleIO) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderSimpleIO) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	var keyFn, valFn *encFnSimpleIO

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

func (e *encoderSimpleIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnSimpleIO) {

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

type encoderSimpleIO struct {
	dh helperEncDriverSimpleIO
	fp *fastpathEsSimpleIO
	e  simpleEncDriverIO
	encoderBase
}

func (e *encoderSimpleIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsSimpleIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderSimpleIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderSimpleIO) MustEncode(v interface{}) {
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

func (e *encoderSimpleIO) Encode(v interface{}) (err error) {

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

func (e *encoderSimpleIO) encode(iv interface{}) {

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

func (e *encoderSimpleIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderSimpleIO) encodeValue(rv reflect.Value, fn *encFnSimpleIO) {

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
					halt.errorf("circular reference found: %p, %T", sptr, sptr)
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

func (e *encoderSimpleIO) encodeValueNonNil(rv reflect.Value, fn *encFnSimpleIO) {
	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if fn.i.addrE {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderSimpleIO) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderSimpleIO) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderSimpleIO) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoderSimpleIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderSimpleIO) fn(t reflect.Type) *encFnSimpleIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoderSimpleIO) fnNoExt(t reflect.Type) *encFnSimpleIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

func (e *encoderSimpleIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderSimpleIO) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoderSimpleIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderSimpleIO) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoderSimpleIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderSimpleIO) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoderSimpleIO) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

func (e *encoderSimpleIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderSimpleIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderSimpleIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderSimpleIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderSimpleIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

type encFnSimpleIO struct {
	i  encFnInfo
	fe func(*encoderSimpleIO, *encFnInfo, reflect.Value)
}
type encRtidFnSimpleIO struct {
	rtid uintptr
	fn   *encFnSimpleIO
}

func (helperEncDriverSimpleIO) newEncoderBytes(out *[]byte, h Handle) *encoderSimpleIO {
	var c1 encoderSimpleIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverSimpleIO) newEncoderIO(out io.Writer, h Handle) *encoderSimpleIO {
	var c1 encoderSimpleIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverSimpleIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsSimpleIO) (f *fastpathESimpleIO, u reflect.Type) {
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

func (helperEncDriverSimpleIO) encFindRtidFn(s []encRtidFnSimpleIO, rtid uintptr) (i uint, fn *encFnSimpleIO) {

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

func (helperEncDriverSimpleIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnSimpleIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnSimpleIO](v))
	}
	return
}

func (dh helperEncDriverSimpleIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsSimpleIO, checkExt bool) (fn *encFnSimpleIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverSimpleIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnSimpleIO = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverSimpleIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnSimpleIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnSimpleIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnSimpleIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnSimpleIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverSimpleIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleIO) {
	fn = new(encFnSimpleIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderSimpleIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderSimpleIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderSimpleIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderSimpleIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fe = (*encoderSimpleIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderSimpleIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderSimpleIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderSimpleIO).textMarshal
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
					fn.fe = func(e *encoderSimpleIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderSimpleIO).kBool
			case reflect.String:

				fn.fe = (*encoderSimpleIO).kString
			case reflect.Int:
				fn.fe = (*encoderSimpleIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderSimpleIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderSimpleIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderSimpleIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderSimpleIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderSimpleIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderSimpleIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderSimpleIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderSimpleIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderSimpleIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderSimpleIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderSimpleIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderSimpleIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderSimpleIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderSimpleIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderSimpleIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderSimpleIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderSimpleIO).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*encoderSimpleIO).kStruct
				} else {
					fn.fe = (*encoderSimpleIO).kStructNoOmitempty
				}
			case reflect.Map:
				fn.fe = (*encoderSimpleIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderSimpleIO).kErr
			default:

				fn.fe = (*encoderSimpleIO).kErr
			}
		}
	}
	return
}

type helperEncDriverSimpleIO struct{}

func (d *decoderSimpleIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *decoderSimpleIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderSimpleIO) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderSimpleIO) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderSimpleIO) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOk(d.d.DecodeStringAsBytes(nil)))
	halt.onerror(fnerr)
}

func (d *decoderSimpleIO) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderSimpleIO) jsonUnmarshalV(tm jsonUnmarshaler) {

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
	halt.onerror(fnerr)
}

func (d *decoderSimpleIO) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderSimpleIO) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderSimpleIO) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes(zeroByteSlice)))
}

func (d *decoderSimpleIO) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderSimpleIO) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderSimpleIO) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderSimpleIO) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderSimpleIO) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderSimpleIO) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderSimpleIO) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderSimpleIO) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderSimpleIO) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderSimpleIO) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderSimpleIO) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderSimpleIO) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleIO) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleIO) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderSimpleIO) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderSimpleIO) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderSimpleIO) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderSimpleIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

	n := d.naked()
	d.d.DecodeNaked()

	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		halt.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
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

func (d *decoderSimpleIO) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderSimpleIO) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *decoderSimpleIO) kStruct(f *decFnInfo, rv reflect.Value) {
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
			name2 = make([]byte, 0, 16)
		}
		var rvkencname []byte
		tkt := ti.keyType
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey()

			if tkt == valueTypeString {
				rvkencname, _ = d.d.DecodeStringAsBytes(nil)
			} else if tkt == valueTypeInt {
				rvkencname = strconv.AppendInt(d.b[:0], d.d.DecodeInt64(), 10)
			} else if tkt == valueTypeUint {
				rvkencname = strconv.AppendUint(d.b[:0], d.d.DecodeUint64(), 10)
			} else if tkt == valueTypeFloat {
				rvkencname = strconv.AppendFloat(d.b[:0], d.d.DecodeFloat64(), 'f', -1, 64)
			} else {
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
					halt.errorStr2("no matching struct field when decoding stream map with key: ", stringView(name2))
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
		halt.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderSimpleIO) kSlice(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
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

	var fn *decFnSimpleIO

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
				halt.errorStr("cannot decode into non-settable slice")
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
					halt.errorStr("cannot decode into non-settable slice")
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
					halt.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					halt.onerror(errExpandSliceCannotChange)
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

func (d *decoderSimpleIO) kArray(f *decFnInfo, rv reflect.Value) {

	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {

		if f.ti.elemkind != uint8(reflect.Uint8) {
			halt.errorf("bytes/string in stream can decode into array of bytes, but not %v", f.ti.rt)
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

	var fn *decFnSimpleIO

	var rv9 reflect.Value

	rvlen := rv.Len()
	hasLen := containerLenS > 0
	if hasLen && containerLenS > rvlen {
		halt.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
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

func (d *decoderSimpleIO) kChan(f *decFnInfo, rv reflect.Value) {

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		halt.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2, _ := d.d.DecodeBytes(nil)
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

	var fn *decFnSimpleIO

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
					halt.errorStr("cannot decode into non-settable chan")
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

func (d *decoderSimpleIO) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnSimpleIO
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
	var scratchBuf bool

	fnRvk2 := func() (s string) {
		callFnRvk = false

		if len(kstr2bs) == 1 {
			s = str4byte(kstr2bs[0])
		} else if len(kstr2bs) != 0 {
			s, callFnRvk = d.mapKeyString(&kstrbs, &kstr2bs, scratchBuf)
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
			kstr2bs, scratchBuf = d.d.DecodeStringAsBytes(nil)
			rvSetString(rvk, fnRvk2())
		} else {
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

type decoderSimpleIO struct {
	dh helperDecDriverSimpleIO
	fp *fastpathDsSimpleIO
	d  simpleDecDriverIO
	decoderBase
}

func (d *decoderSimpleIO) init(h Handle) {
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

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsSimpleIO)

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

func (d *decoderSimpleIO) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
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

func (d *decoderSimpleIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderSimpleIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderSimpleIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderSimpleIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderSimpleIO) Decode(v interface{}) (err error) {

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

func (d *decoderSimpleIO) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderSimpleIO) Release() {
}

func (d *decoderSimpleIO) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *decoderSimpleIO) nextValueBytes(start []byte) []byte {
	return d.d.nextValueBytes(start)
}

func (d *decoderSimpleIO) decode(iv interface{}) {

	if iv == nil {
		halt.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {

	case reflect.Value:
		if x, _ := isDecodeable(v); !x {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)
	case *string:
		*v = d.stringZC(d.d.DecodeStringAsBytes(nil))
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

func (d *decoderSimpleIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderSimpleIO) decodeValue(rv reflect.Value, fn *decFnSimpleIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *decoderSimpleIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnSimpleIO) {

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
			halt.errorStr("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *decoderSimpleIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderSimpleIO) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = zeroByteSlice
	}
	v, _ = d.d.DecodeBytes(in)
	return
}

func (d *decoderSimpleIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes(zeroByteSlice)
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderSimpleIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderSimpleIO) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderSimpleIO) checkBreak() (v bool) {
	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *decoderSimpleIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *decoderSimpleIO) mapElemKey() {
	d.d.ReadMapElemKey()
	d.c = containerMapKey
}

func (d *decoderSimpleIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderSimpleIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleIO) arrayElem() {
	d.d.ReadArrayElem()
	d.c = containerArrayElem
}

func (d *decoderSimpleIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {
	rv := d.interfaceExtConvertAndDecodeGetRV(v, ext)
	d.decodeValue(rv, nil)
	ext.UpdateExt(v, rv2i(rv))
}

func (d *decoderSimpleIO) fn(t reflect.Type) *decFnSimpleIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, false)
}

func (d *decoderSimpleIO) fnNoExt(t reflect.Type) *decFnSimpleIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, true)
}

type decSliceHelperSimpleIO struct {
	d     *decoderSimpleIO
	ct    valueType
	Array bool
	IsNil bool
}

func (d *decoderSimpleIO) decSliceHelperStart() (x decSliceHelperSimpleIO, clen int) {
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
		halt.errorStr2("to decode into a slice, expect map/array - got ", x.ct.String())
	}
	return
}

func (x decSliceHelperSimpleIO) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelperSimpleIO) ElemContainerState(index int) {

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 {
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelperSimpleIO) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
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

type decFnSimpleIO struct {
	i  decFnInfo
	fd func(*decoderSimpleIO, *decFnInfo, reflect.Value)
}
type decRtidFnSimpleIO struct {
	rtid uintptr
	fn   *decFnSimpleIO
}

func (helperDecDriverSimpleIO) newDecoderBytes(in []byte, h Handle) *decoderSimpleIO {
	var c1 decoderSimpleIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverSimpleIO) newDecoderIO(in io.Reader, h Handle) *decoderSimpleIO {
	var c1 decoderSimpleIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverSimpleIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsSimpleIO) (f *fastpathDSimpleIO, u reflect.Type) {
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

func (helperDecDriverSimpleIO) decFindRtidFn(s []decRtidFnSimpleIO, rtid uintptr) (i uint, fn *decFnSimpleIO) {

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

func (helperDecDriverSimpleIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnSimpleIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnSimpleIO](v))
	}
	return
}

func (dh helperDecDriverSimpleIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsSimpleIO,
	checkExt bool) (fn *decFnSimpleIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverSimpleIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnSimpleIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverSimpleIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnSimpleIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnSimpleIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnSimpleIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnSimpleIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverSimpleIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleIO) {
	fn = new(decFnSimpleIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderSimpleIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderSimpleIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderSimpleIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderSimpleIO).ext
		fi.addrD = true
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {

		fn.fd = (*decoderSimpleIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderSimpleIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderSimpleIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderSimpleIO).textUnmarshal
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
						fn.fd = func(d *decoderSimpleIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderSimpleIO, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderSimpleIO).kBool
			case reflect.String:
				fn.fd = (*decoderSimpleIO).kString
			case reflect.Int:
				fn.fd = (*decoderSimpleIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderSimpleIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderSimpleIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderSimpleIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderSimpleIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderSimpleIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderSimpleIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderSimpleIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderSimpleIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderSimpleIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderSimpleIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderSimpleIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderSimpleIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderSimpleIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderSimpleIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderSimpleIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderSimpleIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderSimpleIO).kArray
			case reflect.Struct:
				fn.fd = (*decoderSimpleIO).kStruct
			case reflect.Map:
				fn.fd = (*decoderSimpleIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderSimpleIO).kInterface
			default:

				fn.fd = (*decoderSimpleIO).kErr
			}
		}
	}
	return
}

type helperDecDriverSimpleIO struct{}
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

type simpleEncDriverIO struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *SimpleHandle
	e *encoderBase

	w bufioEncWriter
}

func (e *simpleEncDriverIO) EncodeNil() {
	e.w.writen1(simpleVdNil)
}

func (e *simpleEncDriverIO) EncodeBool(b bool) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && !b {
		e.EncodeNil()
		return
	}
	if b {
		e.w.writen1(simpleVdTrue)
	} else {
		e.w.writen1(simpleVdFalse)
	}
}

func (e *simpleEncDriverIO) EncodeFloat32(f float32) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat32)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *simpleEncDriverIO) EncodeFloat64(f float64) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *simpleEncDriverIO) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), simpleVdNegInt)
	} else {
		e.encUint(uint64(v), simpleVdPosInt)
	}
}

func (e *simpleEncDriverIO) EncodeUint(v uint64) {
	e.encUint(v, simpleVdPosInt)
}

func (e *simpleEncDriverIO) encUint(v uint64, bd uint8) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == 0 {
		e.EncodeNil()
		return
	}
	if v <= math.MaxUint8 {
		e.w.writen2(bd, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 1)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 2)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 3)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *simpleEncDriverIO) encLen(bd byte, length int) {
	if length == 0 {
		e.w.writen1(bd)
	} else if length <= math.MaxUint8 {
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	} else if length <= math.MaxUint16 {
		e.w.writen1(bd + 2)
		e.w.writen2(bigen.PutUint16(uint16(length)))
	} else if int64(length) <= math.MaxUint32 {
		e.w.writen1(bd + 3)
		e.w.writen4(bigen.PutUint32(uint32(length)))
	} else {
		e.w.writen1(bd + 4)
		e.w.writen8(bigen.PutUint64(uint64(length)))
	}
}

func (e *simpleEncDriverIO) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *simpleEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *simpleEncDriverIO) encodeExtPreamble(xtag byte, length int) {
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
}

func (e *simpleEncDriverIO) WriteArrayStart(length int) {
	e.encLen(simpleVdArray, length)
}

func (e *simpleEncDriverIO) WriteMapStart(length int) {
	e.encLen(simpleVdMap, length)
}

func (e *simpleEncDriverIO) EncodeString(v string) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == "" {
		e.EncodeNil()
		return
	}
	if e.h.StringToRaw {
		e.encLen(simpleVdByteArray, len(v))
	} else {
		e.encLen(simpleVdString, len(v))
	}
	e.w.writestr(v)
}

func (e *simpleEncDriverIO) EncodeStringBytesRaw(v []byte) {

	if v == nil {
		e.EncodeNil()
		return
	}
	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
}

func (e *simpleEncDriverIO) EncodeTime(t time.Time) {

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	v, err := t.MarshalBinary()
	halt.onerror(err)
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
}

type simpleDecDriverIO struct {
	h *SimpleHandle
	d *decoderBase
	r ioDecReader

	bdAndBdread
	bytes bool

	noBuiltInTypes

	decDriverNoopContainerReader
	decInit2er
}

func (d *simpleDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *simpleDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *simpleDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	switch d.bd {
	case simpleVdNil:
		d.bdRead = false
		return valueTypeNil
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		return valueTypeBytes
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		return valueTypeString
	case simpleVdArray, simpleVdArray + 1,
		simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		return valueTypeArray
	case simpleVdMap, simpleVdMap + 1,
		simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *simpleDecDriverIO) TryNil() bool {
	return d.advanceNil()
}

func (d *simpleDecDriverIO) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case simpleVdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *simpleDecDriverIO) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case simpleVdPosInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case simpleVdPosInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
		neg = true
	default:
		ok = false

	}

	return
}

func (d *simpleDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == simpleVdFalse {
	} else if d.bd == simpleVdTrue {
		b = true
	} else {
		halt.errorf("cannot decode bool - %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverIO) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverIO) uint2Len(ui uint64) int {
	if chkOvf.Uint(ui, intBitsize) {
		halt.errorf("overflow integer: %v", ui)
	}
	return int(ui)
}

func (d *simpleDecDriverIO) decLen() int {
	switch d.bd & 7 {
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(bigen.Uint16(d.r.readn2()))
	case 3:
		return d.uint2Len(uint64(bigen.Uint32(d.r.readn4())))
	case 4:
		return d.uint2Len(bigen.Uint64(d.r.readn8()))
	}
	halt.errorf("cannot read length: bd%%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
}

func (d *simpleDecDriverIO) DecodeStringAsBytes(bs []byte) (s []byte, scratchBuf bool) {
	return d.DecodeBytes(bs)
}

func (d *simpleDecDriverIO) DecodeBytes(bs []byte) (out []byte, scratchBuf bool) {
	if d.advanceNil() {
		return
	}

	if d.bd >= simpleVdArray && d.bd <= simpleVdMap+4 {
		if bs == nil {
			scratchBuf = true
			bs = d.d.b[:]
		}
		slen := d.ReadArrayStart()
		bs, _ = usableByteSlice(bs, slen)
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		out = bs
		return
	}

	clen := d.decLen()
	d.bdRead = false

	if d.bytes && d.h.ZeroCopy {
		return d.r.readx(uint(clen)), false
	}
	if bs == nil {
		scratchBuf = true
		bs = d.d.b[:]
	}
	out = decByteSlice(&d.r, clen, d.h.MaxInitLen, bs)
	return
}

func (d *simpleDecDriverIO) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd != simpleVdTime {
		halt.errorf("invalid descriptor for time.Time - expect 0x%x, received 0x%x", simpleVdTime, d.bd)
	}
	d.bdRead = false
	clen := uint(d.r.readn1())
	b := d.r.readx(clen)
	halt.onerror((&t).UnmarshalBinary(b))
	return
}

func (d *simpleDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *simpleDecDriverIO) decodeExtV(verifyTag bool, tag byte) (xbs []byte, xtag byte, zerocopy bool) {
	switch d.bd {
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		if d.d.bytes {
			xbs = d.r.readx(uint(l))
			zerocopy = true
		} else {
			xbs = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs, _ = d.DecodeBytes(nil)
	default:
		halt.errorf("ext - %s - expecting extensions/bytearray, got: 0x%x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd {
	case simpleVdNil:
		n.v = valueTypeNil
	case simpleVdFalse:
		n.v = valueTypeBool
		n.b = false
	case simpleVdTrue:
		n.v = valueTypeBool
		n.b = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case simpleVdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdTime:
		n.v = valueTypeTime
		n.t = d.DecodeTime()
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		n.v = valueTypeString
		n.s = d.d.stringZC(d.DecodeStringAsBytes(zeroByteSlice))
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		if d.d.bytes {
			n.l = d.r.readx(uint(l))
		} else {
			n.l = decByteSlice(&d.r, l, d.h.MaxInitLen, d.d.b[:])
		}
	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2,
		simpleVdArray + 3, simpleVdArray + 4:
		n.v = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, d.bd)
	}

	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *simpleDecDriverIO) nextValueBytes(v0 []byte) (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	v0 = append(v0, d.bd)
	d.r.startRecording(v0)
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) nextValueBytesBdReadR() {
	c := d.bd

	var length uint

	switch c {
	case simpleVdNil, simpleVdFalse, simpleVdTrue, simpleVdString, simpleVdByteArray:

	case simpleVdPosInt, simpleVdNegInt:
		d.r.readn1()
	case simpleVdPosInt + 1, simpleVdNegInt + 1:
		d.r.readx(2)
	case simpleVdPosInt + 2, simpleVdNegInt + 2, simpleVdFloat32:
		d.r.readx(4)
	case simpleVdPosInt + 3, simpleVdNegInt + 3, simpleVdFloat64:
		d.r.readx(8)
	case simpleVdTime:
		c = d.r.readn1()
		d.r.readx(uint(c))

	default:
		switch c & 7 {
		case 0:
			length = 0
		case 1:
			b := d.r.readn1()
			length = uint(b)
		case 2:
			x := d.r.readn2()
			length = uint(bigen.Uint16(x))
		case 3:
			x := d.r.readn4()
			length = uint(bigen.Uint32(x))
		case 4:
			x := d.r.readn8()
			length = uint(bigen.Uint64(x))
		}

		bExt := c >= simpleVdExt && c <= simpleVdExt+7
		bStr := c >= simpleVdString && c <= simpleVdString+7
		bByteArray := c >= simpleVdByteArray && c <= simpleVdByteArray+7
		bArray := c >= simpleVdArray && c <= simpleVdArray+7
		bMap := c >= simpleVdMap && c <= simpleVdMap+7

		if !(bExt || bStr || bByteArray || bArray || bMap) {
			halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, c)
		}

		if bExt {
			d.r.readn1()
		}

		if length == 0 {
			break
		}

		if bArray {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else if bMap {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			d.r.readx(length)
		}
	}
	return
}

func (d *simpleEncDriverIO) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*SimpleHandle)
	d.e = shared
	if shared.bytes {
		fp = simpleFpEncBytes
	} else {
		fp = simpleFpEncIO
	}

	d.init2(enc)
	return
}

func (e *simpleEncDriverIO) writeBytesAsis(b []byte)           { e.w.writeb(b) }
func (e *simpleEncDriverIO) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *simpleEncDriverIO) writerEnd()                        { e.w.end() }

func (e *simpleEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *simpleEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *simpleDecDriverIO) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*SimpleHandle)
	d.bytes = shared.bytes
	d.d = shared
	if shared.bytes {
		fp = simpleFpDecBytes
	} else {
		fp = simpleFpDecIO
	}

	d.init2(dec)
	return
}

func (d *simpleDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *simpleDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *simpleDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, &d.d.blist)
}

func (d *simpleDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, simpledesc(d.bd))
}

func (d *simpleDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
