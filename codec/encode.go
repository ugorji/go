// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"errors"
	"io"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
)

// defEncByteBufSize is the default size of []byte used
// for bufio buffer or []byte (when nil passed)
const defEncByteBufSize = 1 << 10 // 4:16, 6:64, 8:256, 10:1024

var errEncoderNotInitialized = errors.New("encoder not initialized")

var encBuiltinRtids []uintptr

func init() {
	for _, v := range []interface{}{
		(string)(""),
		(bool)(false),
		(int)(0),
		(int8)(0),
		(int16)(0),
		(int32)(0),
		(int64)(0),
		(uint)(0),
		(uint8)(0),
		(uint16)(0),
		(uint32)(0),
		(uint64)(0),
		(uintptr)(0),
		(float32)(0),
		(float64)(0),
		(complex64)(0),
		(complex128)(0),
		(time.Time{}),
		([]byte)(nil),
		(Raw{}),
		// (interface{})(nil),
	} {
		t := reflect.TypeOf(v)
		encBuiltinRtids = append(encBuiltinRtids, rt2id(t), rt2id(reflect.PointerTo(t)))
	}
	slices.Sort(encBuiltinRtids)
}

// encDriver abstracts the actual codec (binc vs msgpack, etc)
type encDriverI interface {
	EncodeNil()
	EncodeInt(i int64)
	EncodeUint(i uint64)
	EncodeBool(b bool)
	EncodeFloat32(f float32)
	EncodeFloat64(f float64)
	// re is never nil
	EncodeRawExt(re *RawExt)
	// ext is never nil
	EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext)
	// EncodeString using cUTF8, honor'ing StringToRaw flag
	EncodeString(v string)
	EncodeStringNoEscape4Json(v string)
	EncodeStringBytesRaw(v []byte)
	EncodeTime(time.Time)
	WriteArrayStart(length int)
	WriteArrayEnd()
	WriteMapStart(length int)
	WriteMapEnd()

	// these are no-op except for json
	encDriverContainerTracker

	// reset will reset current encoding runtime state, and cached information from the handle
	reset()

	atEndOfEncode()
	writerEnd()

	writeBytesAsis(b []byte)
	// writeStringAsisDblQuoted(v string)

	resetOutBytes(out *[]byte)
	resetOutIO(out io.Writer)

	init(h Handle, shared *encoderBase, enc encoderI) (fp interface{})

	// driverStateManager
}

type helperEncDriver[T encDriver] struct{}

type encInit2er struct{}

func (encInit2er) init2(enc encoderI) {}

type encDriverContainerTracker interface {
	WriteArrayElem()
	WriteMapElemKey()
	WriteMapElemValue()
}

type encDriverNoState struct{}

// func (encDriverNoState) captureState() interface{}  { return nil }
// func (encDriverNoState) resetState()                {}
// func (encDriverNoState) restoreState(v interface{}) {}
func (encDriverNoState) reset() {}

type encDriverNoopContainerWriter struct{}

func (encDriverNoopContainerWriter) WriteArrayStart(length int) {}
func (encDriverNoopContainerWriter) WriteArrayEnd()             {}
func (encDriverNoopContainerWriter) WriteMapStart(length int)   {}
func (encDriverNoopContainerWriter) WriteMapEnd()               {}
func (encDriverNoopContainerWriter) atEndOfEncode()             {}

// encStructFieldObj[Slice] is used for sorting when there are missing fields and canonical flag is set
type encStructFieldObj struct {
	key        string
	rv         reflect.Value
	intf       interface{}
	isRv       bool
	noEsc4json bool
	builtin    bool
}

type encStructFieldObjSlice []encStructFieldObj

func (p encStructFieldObjSlice) Len() int      { return len(p) }
func (p encStructFieldObjSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p encStructFieldObjSlice) Less(i, j int) bool {
	return p[uint(i)].key < p[uint(j)].key
}

// ----

type encFnInfo struct {
	ti    *typeInfo
	xfFn  Ext
	xfTag uint64
	addrE bool
	// addrEf bool // force: if addrE, then encode function MUST take a ptr
}

// encFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the
// code block that should be called (encapsulated in a function)
// instead of executing the checks every time.
type encFn[T encDriver] struct {
	i  encFnInfo
	fe func(*encoder[T], *encFnInfo, reflect.Value)
	// _  [1]uint64 // padding (cache-aligned)
}

type encRtidFn[T encDriver] struct {
	rtid uintptr
	fn   *encFn[T]
}

// ----

// EncodeOptions captures configuration options during encode.
type EncodeOptions struct {
	// WriterBufferSize is the size of the buffer used when writing.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	WriterBufferSize int

	// ChanRecvTimeout is the timeout used when selecting from a chan.
	//
	// Configuring this controls how we receive from a chan during the encoding process.
	//   - If ==0, we only consume the elements currently available in the chan.
	//   - if  <0, we consume until the chan is closed.
	//   - If  >0, we consume until this timeout.
	ChanRecvTimeout time.Duration

	// StructToArray specifies to encode a struct as an array, and not as a map
	StructToArray bool

	// Canonical representation means that encoding a value will always result in the same
	// sequence of bytes.
	//
	// This only affects maps, as the iteration order for maps is random.
	//
	// The implementation MAY use the natural sort order for the map keys if possible:
	//
	//     - If there is a natural sort order (ie for number, bool, string or []byte keys),
	//       then the map keys are first sorted in natural order and then written
	//       with corresponding map values to the strema.
	//     - If there is no natural sort order, then the map keys will first be
	//       encoded into []byte, and then sorted,
	//       before writing the sorted keys and the corresponding map values to the stream.
	//
	Canonical bool

	// CheckCircularRef controls whether we check for circular references
	// and error fast during an encode.
	//
	// If enabled, an error is received if a pointer to a struct
	// references itself either directly or through one of its fields (iteratively).
	//
	// This is opt-in, as there may be a performance hit to checking circular references.
	CheckCircularRef bool

	// RecursiveEmptyCheck controls how we determine whether a value is empty.
	//
	// If true, we descend into interfaces and pointers to reursively check if value is empty.
	//
	// We *might* check struct fields one by one to see if empty
	// (if we cannot directly check if a struct value is equal to its zero value).
	// If so, we honor IsZero, Comparable, IsCodecEmpty(), etc.
	// Note: This *may* make OmitEmpty more expensive due to the large number of reflect calls.
	//
	// If false, we check if the value is equal to its zero value (newly allocated state).
	RecursiveEmptyCheck bool

	// Raw controls whether we encode Raw values.
	// This is a "dangerous" option and must be explicitly set.
	// If set, we blindly encode Raw values as-is, without checking
	// if they are a correct representation of a value in that format.
	// If unset, we error out.
	Raw bool

	// StringToRaw controls how strings are encoded.
	//
	// As a go string is just an (immutable) sequence of bytes,
	// it can be encoded either as raw bytes or as a UTF string.
	//
	// By default, strings are encoded as UTF-8.
	// but can be treated as []byte during an encode.
	//
	// Note that things which we know (by definition) to be UTF-8
	// are ALWAYS encoded as UTF-8 strings.
	// These include encoding.TextMarshaler, time.Format calls, struct field names, etc.
	StringToRaw bool

	// OptimumSize controls whether we optimize for the smallest size.
	//
	// Some formats will use this flag to determine whether to encode
	// in the smallest size possible, even if it takes slightly longer.
	//
	// For example, some formats that support half-floats might check if it is possible
	// to store a float64 as a half float. Doing this check has a small performance cost,
	// but the benefit is that the encoded message will be smaller.
	OptimumSize bool

	// NoAddressableReadonly controls whether we try to force a non-addressable value
	// to be addressable so we can call a pointer method on it e.g. for types
	// that support Selfer, json.Marshaler, etc.
	//
	// Use it in the very rare occurrence that your types modify a pointer value when calling
	// an encode callback function e.g. JsonMarshal, TextMarshal, BinaryMarshal or CodecEncodeSelf.
	NoAddressableReadonly bool
}

// ---------------------------------------------

// encoderBase is shared as a field between Encoder and its encDrivers.
// This way, encDrivers need not hold a referece to the Encoder itself.
type encoderBase struct {
	perType encPerType

	h *BasicHandle

	// MARKER: these fields below should belong directly in Encoder.
	// There should not be any pointers here - just values.
	// we pack them here for space efficiency and cache-line optimization.

	rtidFn, rtidFnNoExt *atomicRtidFnSlice

	// se  encoderI
	err error

	blist bytesFreeList

	// js bool // is json encoder?
	// be bool // is binary encoder?

	bytes bool

	c containerState

	calls uint16
	seq   uint16 // sequencer (e.g. used by binc for symbols, etc)

	// ---- cpu cache line boundary
	hh Handle

	// ---- cpu cache line boundary

	// ---- writable fields during execution --- *try* to keep in sep cache line

	ci circularRefChecker

	slist sfiRvFreeList
}

func (e *encoderBase) HandleName() string {
	return e.hh.Name()
}

// Release is a no-op.
//
// Deprecated: Pooled resources are not used with an Encoder.
// This method is kept for compatibility reasons only.
func (e *encoderBase) Release() {
}

func (e *encoderBase) setContainerState(cs containerState) {
	if cs != 0 {
		e.c = cs
	}
}

// Encoder writes an object to an output stream in a supported format.
//
// Encoder is NOT safe for concurrent use i.e. a Encoder cannot be used
// concurrently in multiple goroutines.
//
// However, as Encoder could be allocation heavy to initialize, a Reset method is provided
// so its state can be reused to decode new input streams repeatedly.
// This is the idiomatic way to use.
type encoder[T encDriver] struct {
	dh helperEncDriver[T]
	fp *fastpathEs[T]
	e  T
	encoderBase
}

func (e *encoder[T]) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoder[T]) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoder[T]) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoder[T]) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoder[T]) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoder[T]) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoder[T]) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoder[T]) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoder[T]) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoder[T]) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoder[T]) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoder[T]) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoder[T]) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoder[T]) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoder[T]) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoder[T]) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoder[T]) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoder[T]) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoder[T]) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoder[T]) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoder[T]) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoder[T]) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoder[T]) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoder[T]) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoder[T]) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoder[T]) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoder[T]) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderBase) kErr(_ *encFnInfo, rv reflect.Value) {
	halt.errorf("unsupported encoding kind: %s, for %#v", rv.Kind(), any(rv))
}

func chanToSlice(rv reflect.Value, rtslice reflect.Type, timeout time.Duration) (rvcs reflect.Value) {
	rvcs = rvZeroK(rtslice, reflect.Slice)
	if timeout < 0 { // consume until close
		for {
			recv, recvOk := rv.Recv()
			if !recvOk {
				break
			}
			rvcs = reflect.Append(rvcs, recv)
		}
	} else {
		cases := make([]reflect.SelectCase, 2)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: rv}
		if timeout == 0 {
			cases[1] = reflect.SelectCase{Dir: reflect.SelectDefault}
		} else {
			tt := time.NewTimer(timeout)
			cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(tt.C)}
		}
		for {
			chosen, recv, recvOk := reflect.Select(cases)
			if chosen == 1 || !recvOk {
				break
			}
			rvcs = reflect.Append(rvcs, recv)
		}
	}
	return
}

func (e *encoder[T]) kSeqFn(rt reflect.Type) (fn *encFn[T]) {
	// if kind is reflect.Interface, do not pre-determine the encoding type,
	// because preEncodeValue may break it down to a concrete type and kInterface will bomb.
	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoder[T]) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
	var builtin bool
	var fn *encFn[T]
	var l = rvLenSlice(rv)
	if l == 0 {
		e.mapStart(0)
		goto END
	}
	e.haltOnMbsOddLen(l)
	e.mapStart(l >> 1) // e.mapStart(l / 2)
	builtin = ti.tielem.flagEncBuiltin
	if !builtin {
		fn = e.kSeqFn(ti.elem)
	}
	for j := 0; j < l; j++ {
		if j&1 == 0 { // j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		rvv := rvSliceIndex(rv, j, ti)
		if builtin {
			e.encode(rv2i(baseRVRV(rvv)))
		} else {
			e.encodeValue(rvv, fn)
		}
	}
END:
	e.mapEnd()
}

func (e *encoder[T]) kSliceW(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	e.arrayStart(l)
	if l <= 0 {
		goto END
	}
	if ti.tielem.flagEncBuiltin {
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encode(rv2i(baseRVRV(rvSliceIndex(rv, j, ti))))
		}
	} else {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
END:
	e.arrayEnd()
}

func (e *encoder[T]) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
	var builtin bool
	var fn *encFn[T]
	var l = rv.Len()
	if l == 0 {
		e.mapStart(0)
		goto END
	}
	e.haltOnMbsOddLen(l)
	e.mapStart(l >> 1) // e.mapStart(l / 2)
	builtin = ti.tielem.flagEncBuiltin
	if !builtin {
		fn = e.kSeqFn(ti.elem)
	}
	for j := 0; j < l; j++ {
		if j&1 == 0 { // j%2 == 0 {
			e.mapElemKey()
		} else {
			e.mapElemValue()
		}
		rvv := rvArrayIndex(rv, j, ti)
		if builtin {
			e.encode(rv2i(baseRVRV(rvv)))
		} else {
			e.encodeValue(rvv, fn)
		}
	}
END:
	e.mapEnd()
}

func (e *encoder[T]) kArrayW(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	e.arrayStart(l)
	if l <= 0 {
		goto END
	}
	if ti.tielem.flagEncBuiltin {
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encode(rv2i(baseRVRV(rvArrayIndex(rv, j, ti))))
		}
	} else {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rvArrayIndex(rv, j, ti), fn)
		}
	}
END:
	e.arrayEnd()
}

func (e *encoder[T]) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoder[T]) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *encoder[T]) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *encoder[T]) kSliceBytesChan(rv reflect.Value) {
	// do not use range, so that the number of elements encoded
	// does not change, and encoding does not hang waiting on someone to close chan.

	bs0 := e.blist.peek(32, true)
	bs := bs0

	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok {
		ch = irv.(chan byte)
	}

L1:
	switch timeout := e.h.ChanRecvTimeout; {
	case timeout == 0: // only consume available
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			}
		}
	case timeout > 0: // consume until timeout
		tt := time.NewTimer(timeout)
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:
				// close(tt.C)
				break L1
			}
		}
	default: // consume until close
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

func (e *encoder[T]) kStructFieldKey(keyType valueType, encName string) {
	// use if (not switch) block, so that branch prediction picks valueTypeString first
	if keyType == valueTypeString {
		e.e.EncodeString(encName)
	} else if keyType == valueTypeInt {
		e.e.EncodeInt(must.Int(strconv.ParseInt(encName, 10, 64)))
	} else if keyType == valueTypeUint {
		e.e.EncodeUint(must.Uint(strconv.ParseUint(encName, 10, 64)))
	} else if keyType == valueTypeFloat {
		e.e.EncodeFloat64(must.Float(strconv.ParseFloat(encName, 64)))
	} else {
		halt.errorStr2("invalid struct key type: ", keyType.String())
	}
	// e.dh.encStructFieldKey(e.e, encName, keyType, encNameAsciiAlphaNum, e.js)
}

func (e *encoder[T]) kStructSimple(f *encFnInfo, rv reflect.Value) {
	tisfi := f.ti.sfi.source()

	// To bypass encodeValue, we need to handle cases where
	// the field is an interface kind. To do this, we need to handle an
	// interface or a pointer to an interface differently.
	//
	// Easiest to just delegate to encodeValue.

	chkCirRef := e.h.CheckCircularRef
	var si *structFieldInfo

	// use value of chkCirRef ie if true, then send the addr of the value
	if f.ti.toArray || e.h.StructToArray { // toArray
		e.arrayStart(len(tisfi))
		for _, si = range tisfi {
			e.arrayElem()
			if si.encBuiltin {
				e.encode(rv2i(si.fieldNoAlloc(rv, true)))
			} else {
				e.encodeValue(si.fieldNoAlloc(rv, !chkCirRef), nil)
			}
		}
		e.arrayEnd()
	} else {
		if e.h.Canonical {
			tisfi = f.ti.sfi.sorted()
		}
		e.mapStart(len(tisfi))
		for _, si = range tisfi {
			e.mapElemKey()
			e.e.EncodeStringNoEscape4Json(si.encName)
			e.mapElemValue()
			if si.encBuiltin {
				e.encode(rv2i(si.fieldNoAlloc(rv, true)))
			} else {
				e.encodeValue(si.fieldNoAlloc(rv, !chkCirRef), nil)
			}
		}
		e.mapEnd()
	}
}

func (e *encoder[T]) kStruct(f *encFnInfo, rv reflect.Value) {
	ti := f.ti
	toMap := !(ti.toArray || e.h.StructToArray)
	var mf map[string]interface{}
	if ti.flagMissingFielder {
		toMap = true
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
	} else if ti.flagMissingFielderPtr {
		toMap = true
		if rv.CanAddr() {
			mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder).CodecMissingFields()
		} else {
			mf = rv2i(e.addrRV(rv, ti.rt, ti.ptr)).(MissingFielder).CodecMissingFields()
		}
	}
	newlen := len(mf)
	tisfi := ti.sfi.source()
	newlen += len(tisfi)

	var fkvs = e.slist.get(newlen)[:newlen]

	recur := e.h.RecursiveEmptyCheck
	chkCirRef := e.h.CheckCircularRef

	var kv sfiRv
	var j int
	if toMap {
		newlen = 0
		if e.h.Canonical {
			tisfi = f.ti.sfi.sorted()
		}
		for _, si := range tisfi {
			// kv.r = si.path.field(rv, false, si.encBuiltin || !chkCirRef)
			// if si.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
			// 	continue
			// }
			if si.omitEmpty {
				kv.r = si.fieldNoAlloc(rv, false) // test actual field val
				if isEmptyValue(kv.r, e.h.TypeInfos, recur) {
					continue
				}
			} else {
				kv.r = si.fieldNoAlloc(rv, si.encBuiltin || !chkCirRef)
			}
			kv.v = si
			fkvs[newlen] = kv
			newlen++
		}

		var mf2s []stringIntf
		if len(mf) != 0 {
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

		// When there are missing fields, and Canonical flag is set,
		// we cannot have the missing fields and struct fields sorted independently.
		// We have to capture them together and sort as a unit.

		if len(mf2s) != 0 && e.h.Canonical {
			mf2w := make([]encStructFieldObj, newlen+len(mf2s))
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				mf2w[j] = encStructFieldObj{kv.v.encName, kv.r, nil, true,
					!kv.v.encNameEscape4Json, kv.v.encBuiltin}
			}
			for _, v := range mf2s {
				mf2w[j] = encStructFieldObj{v.v, reflect.Value{}, v.i, false, false, false}
				j++
			}
			sort.Sort((encStructFieldObjSlice)(mf2w))
			for _, v := range mf2w {
				e.mapElemKey()
				if ti.keyType == valueTypeString && v.noEsc4json {
					e.e.EncodeStringNoEscape4Json(v.key)
				} else {
					e.kStructFieldKey(ti.keyType, v.key)
				}
				e.mapElemValue()
				if v.isRv {
					if v.builtin {
						e.encode(rv2i(baseRVRV(v.rv)))
					} else {
						e.encodeValue(v.rv, nil)
					}
				} else {
					e.encode(v.intf)
				}
			}
		} else {
			keytyp := ti.keyType
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				e.mapElemKey()
				if ti.keyType == valueTypeString && !kv.v.encNameEscape4Json {
					e.e.EncodeStringNoEscape4Json(kv.v.encName)
				} else {
					e.kStructFieldKey(keytyp, kv.v.encName)
				}
				e.mapElemValue()
				if kv.v.encBuiltin {
					e.encode(rv2i(baseRVRV(kv.r)))
				} else {
					e.encodeValue(kv.r, nil)
				}
			}
			for _, v := range mf2s {
				e.mapElemKey()
				e.kStructFieldKey(keytyp, v.v)
				e.mapElemValue()
				e.encode(v.i)
			}
		}

		e.mapEnd()
	} else {
		newlen = len(tisfi)
		for i, si := range tisfi { // use unsorted array (to match sequence in struct)
			// kv.r = si.path.field(rv, false, si.encBuiltin || !chkCirRef)
			// kv.r = si.path.field(rv, false, !si.omitEmpty || si.encBuiltin || !chkCirRef)
			if si.omitEmpty {
				// use the zero value.
				// if a reference or struct, set to nil (so you do not output too much)
				kv.r = si.fieldNoAlloc(rv, false) // test actual field val
				if isEmptyContainerValue(kv.r, e.h.TypeInfos, recur) {
					kv.r = reflect.Value{} //encode as nil
				}
			} else {
				kv.r = si.fieldNoAlloc(rv, si.encBuiltin || !chkCirRef)
			}
			kv.v = si
			fkvs[i] = kv
		}
		// encode it all
		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ {
			e.arrayElem()
			kv = fkvs[j]
			if !kv.r.IsValid() {
				e.e.EncodeNil()
			} else if kv.v.encBuiltin {
				e.encode(rv2i(baseRVRV(kv.r)))
			} else {
				e.encodeValue(kv.r, nil)
			}
		}
		e.arrayEnd()
	}

	// do not use defer. Instead, use explicit pool return at end of function.
	// defer has a cost we are trying to avoid.
	// If there is a panic and these slices are not returned, it is ok.
	e.slist.put(fkvs)
}

func (e *encoder[T]) kMap(f *encFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	// determine the underlying key and val encFn's for the map.
	// This eliminates some work which is done for each loop iteration i.e.
	// rv.Type(), ref.ValueOf(rt).Pointer(), then check map/list for fn.
	//
	// However, if kind is reflect.Interface, do not pre-determine the
	// encoding type, because preEncodeValue may break it down to
	// a concrete type and kInterface will bomb.

	var keyFn, valFn *encFn[T]

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
	var keyTypeIsString = stringTypId == rt2id(rtkey) // rtkeyid
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

	kbuiltin := f.ti.tikey.flagEncBuiltin
	vbuiltin := f.ti.tielem.flagEncBuiltin
	for it.Next() {
		rv = it.Key()
		e.mapElemKey()
		if keyTypeIsString {
			e.e.EncodeString(rvGetString(rv))
		} else if kbuiltin {
			e.encode(rv2i(baseRVRV(rv)))
		} else {
			e.encodeValue(rv, keyFn)
		}
		e.mapElemValue()
		rv = it.Value()
		if vbuiltin {
			e.encode(rv2i(baseRVRV(rv)))
		} else {
			e.encodeValue(it.Value(), valFn)
		}
	}
	it.Done()

	e.mapEnd()
}

func (e *encoder[T]) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFn[T]) {
	// The base kind of the type of the map key is sufficient for ordering.
	// We only do out of band if that kind is not ordered (number or string), bool or time.Time.
	// If the key is a predeclared type, directly call methods on encDriver e.g. EncodeString
	// but if not, call encodeValue, in case it has an extension registered or otherwise.
	rtkey := ti.key
	rtkeydecl := rtkey.PkgPath() == "" && rtkey.Name() != "" // key type is predeclared

	mks := rv.MapKeys()
	rtkeyKind := rtkey.Kind()
	kfast := mapKeyFastKindFor(rtkeyKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	switch rtkeyKind {
	case reflect.Bool:
		// though bool keys make no sense in a map, it *could* happen.
		// in that case, we MUST support it in reflection mode,
		// as that is the fallback for even codecgen and others.

		// sort the keys so that false comes before true
		// ie if 2 keys in order (true, false), then swap them
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
			v.v = rvGetString(k)
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

		// out-of-band
		// first encode each key to a []byte first, then sort them, then record
		bs0 := e.blist.get(len(mks) * 16)
		mksv := bs0
		mksbv := make([]bytesRv, len(mks))

		sideEncode(e.hh, &e.h.sideEncPool, func(se encoderI) {
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
		})

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

func (e *encoder[T]) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()
	// e.be = e.hh.isBinary()
	e.err = errEncoderNotInitialized

	// e.fp = fastpathEList[T]()
	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEs[T])

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoder[T]) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

// MustEncode is like Encode, but panics if unable to Encode.
//
// Note: This provides insight to the code location that triggered the error.
func (e *encoder[T]) MustEncode(v interface{}) {
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

// Encode writes an object into a stream.
//
// Encoding can be configured via the struct tag for the fields.
// The key (in the struct tags) that we look at is configurable.
//
// By default, we look up the "codec" key in the struct field's tags,
// and fall bak to the "json" key if "codec" is absent.
// That key in struct field's tag value is the key name,
// followed by an optional comma and options.
//
// To set an option on all fields (e.g. omitempty on all fields), you
// can create a field called _struct, and set flags on it. The options
// which can be set on _struct are:
//   - omitempty: so all fields are omitted if empty
//   - toarray: so struct is encoded as an array
//   - int: so struct key names are encoded as signed integers (instead of strings)
//   - uint: so struct key names are encoded as unsigned integers (instead of strings)
//   - float: so struct key names are encoded as floats (instead of strings)
//
// More details on these below.
//
// Struct values "usually" encode as maps. Each exported struct field is encoded unless:
//   - the field's tag is "-", OR
//   - the field is empty (empty or the zero value) and its tag specifies the "omitempty" option.
//
// When encoding as a map, the first string in the tag (before the comma)
// is the map key string to use when encoding.
// ...
// This key is typically encoded as a string.
// However, there are instances where the encoded stream has mapping keys encoded as numbers.
// For example, some cbor streams have keys as integer codes in the stream, but they should map
// to fields in a structured object. Consequently, a struct is the natural representation in code.
// For these, configure the struct to encode/decode the keys as numbers (instead of string).
// This is done with the int,uint or float option on the _struct field (see above).
//
// However, struct values may encode as arrays. This happens when:
//   - StructToArray Encode option is set, OR
//   - the tag on the _struct field sets the "toarray" option
//
// Note that omitempty is ignored when encoding struct values as arrays,
// as an entry must be encoded for each field, to maintain its position.
//
// Values with types that implement MapBySlice are encoded as stream maps.
//
// The empty values (for omitempty option) are false, 0, any nil pointer
// or interface value, and any array, slice, map, or string of length zero.
//
// Anonymous fields are encoded inline except:
//   - the struct tag specifies a replacement name (first value)
//   - the field is of an interface type
//
// Examples:
//
//	// NOTE: 'json:' can be used as struct tag key, in place 'codec:' below.
//	type MyStruct struct {
//	    _struct bool    `codec:",omitempty"`   //set omitempty for every field
//	    Field1 string   `codec:"-"`            //skip this field
//	    Field2 int      `codec:"myName"`       //Use key "myName" in encode stream
//	    Field3 int32    `codec:",omitempty"`   //use key "Field3". Omit if empty.
//	    Field4 bool     `codec:"f4,omitempty"` //use key "f4". Omit if empty.
//	    io.Reader                              //use key "Reader".
//	    MyStruct        `codec:"my1"           //use key "my1".
//	    MyStruct                               //inline it
//	    ...
//	}
//
//	type MyStruct struct {
//	    _struct bool    `codec:",toarray"`     //encode struct as an array
//	}
//
//	type MyStruct struct {
//	    _struct bool    `codec:",uint"`        //encode struct with "unsigned integer" keys
//	    Field1 string   `codec:"1"`            //encode Field1 key using: EncodeInt(1)
//	    Field2 string   `codec:"2"`            //encode Field2 key using: EncodeInt(2)
//	}
//
// The mode of encoding is based on the type of the value. When a value is seen:
//   - If a Selfer, call its CodecEncodeSelf method
//   - If an extension is registered for it, call that extension function
//   - If implements encoding.(Binary|Text|JSON)Marshaler, call Marshal(Binary|Text|JSON) method
//   - Else encode it based on its reflect.Kind
//
// Note that struct field names and keys in map[string]XXX will be treated as symbols.
// Some formats support symbols (e.g. binc) and will properly encode the string
// only once in the stream, and use a tag to refer to it thereafter.
func (e *encoder[T]) Encode(v interface{}) (err error) {
	// tried to use closure, as runtime optimizes defer with no params.
	// This seemed to be causing weird issues (like circular reference found, unexpected panic, etc).
	// Also, see https://github.com/golang/go/issues/14939#issuecomment-417836139
	if !debugging {
		defer func() {
			// if error occurred during encoding, return that error;
			// else if error occurred on end'ing (i.e. during flush), return that error.
			if x := recover(); x != nil {
				panicValToErr(e, x, &e.err)
				err = e.err
			}
		}()
	}

	e.MustEncode(v)
	return
}

func (e *encoder[T]) encode(iv interface{}) {
	// MARKER: a switch with only concrete types can be optimized.
	// consequently, we deal with nil and interfaces outside the switch.

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
	// case nil:
	// case Selfer:
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
		// we can't check non-predefined types, as they might be a Selfer or extension.
		if skipFastpathTypeSwitchInDirectCall || !e.dh.fastpathEncodeTypeSwitch(iv, e) {
			e.encodeValue(rv, nil)
		}
	}
}

// encodeValue will encode a value.
//
// Note that encodeValue will handle nil in the stream early, so that the
// subsequent calls i.e. kXXX methods, etc do not have to handle it themselves.
func (e *encoder[T]) encodeValue(rv reflect.Value, fn *encFn[T]) {
	// if a valid fn is passed, it MUST BE for the dereferenced type of rv

	// MARKER: We check if value is nil here, so that the kXXX method do not have to.

	var ciPushes int
	if e.h.CheckCircularRef {
		ciPushes = e.ci.pushRV(rv)
	}

	var rvp reflect.Value
	var rvpValid bool
TOP:
	switch rv.Kind() {
	case reflect.Ptr:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
		rvpValid = true
		rvp = rv
		rv = rv.Elem()
		goto TOP
	case reflect.Interface:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
		rvpValid = false
		rvp = reflect.Value{}
		rv = rv.Elem()
		goto TOP
	case reflect.Map, reflect.Slice, reflect.Chan:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		goto END
	}

	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if !fn.i.addrE { // typically, addrE = false, so check it first
		// keep rv same
	} else if rvpValid {
		rv = rvp
	} else if rv.CanAddr() {
		rv = rvAddr(rv, fn.i.ti.ptr)
	} else {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
END:
	if ciPushes > 0 {
		e.ci.pop(ciPushes)
	}
}

// encodeValueNonNil can encode a number, bool, or string
// OR non-nil values of kind map, slice and chan.
//
// It expects fn to be non-nil.
func (e *encoder[T]) encodeValueNonNil(rv reflect.Value, fn *encFn[T]) {
	// if fn == nil {
	// 	fn = e.fn(rv.Type())
	// }
	if fn.i.addrE { // typically, addrE = false, so check it first
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoder[T]) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

// addrRV returns a addressable value given that rv is not addressable
func (e *encoderBase) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
	// if rv.CanAddr() {
	// 	return rvAddr(rv, ptrType)
	// }
	if e.h.NoAddressableReadonly {
		rva = reflect.New(typ)
		rvSetDirect(rva.Elem(), rv)
		return
	}
	return rvAddr(e.perType.AddressableRO(rv), ptrType)
}

func (e *encoder[T]) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoder[T]) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs) // e.asis(bs)
	}
}

func (e *encoder[T]) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *encoder[T]) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderBase) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}

func (e *encoder[T]) fn(t reflect.Type) *encFn[T] {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, false)
}

func (e *encoder[T]) fnNoExt(t reflect.Type) *encFn[T] {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, true)
}

// ---- container tracker methods
// Note: We update the .c after calling the callback.
// This way, the callback can know what the last status was.

func (e *encoder[T]) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoder[T]) mapElemKey() {
	e.e.WriteMapElemKey()
	e.c = containerMapKey
}

func (e *encoder[T]) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoder[T]) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *encoder[T]) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoder[T]) arrayElem() {
	e.e.WriteArrayElem()
	e.c = containerArrayElem
}

func (e *encoder[T]) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}

// ----------

func (e *encoderBase) haltOnMbsOddLen(length int) {
	if length&1 != 0 { // similar to &1==1 or %2 == 1
		halt.errorInt("mapBySlice requires even slice length, but got ", int64(length))
	}
}

func (e *encoder[T]) writerEnd() {
	e.e.writerEnd()
}

func (e *encoder[T]) atEndOfEncode() {
	e.e.atEndOfEncode()
}

type encoderI interface {
	Encode(v interface{}) error
	MustEncode(v interface{})
	Release()
	Reset(w io.Writer)
	ResetBytes(out *[]byte)

	wrapErr(v error, err *error)
	atEndOfEncode()
	writerEnd()

	encode(v interface{})
	encodeAs(v interface{}, t reflect.Type, ext bool)

	setContainerState(cs containerState) // needed for canonical encoding via side encoder
}

var errEncNoResetBytesWithWriter = errors.New("cannot reset an Encoder which outputs to []byte with a io.Writer")
var errEncNoResetWriterWithBytes = errors.New("cannot reset an Encoder which outputs to io.Writer with a []byte")

// Reset resets the Encoder with a new output stream.
//
// This accommodates using the state of the Encoder,
// where it has "cached" information about sub-engines.
func (e *encoder[T]) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

// ResetBytes resets the Encoder with a new destination output []byte.
func (e *encoder[T]) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

// only call this iff you are sure it is a bytes encoder
func (e *encoder[T]) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

type encDriverContainerNoTrackerT struct{}

func (encDriverContainerNoTrackerT) WriteArrayElem()    {}
func (encDriverContainerNoTrackerT) WriteMapElemKey()   {}
func (encDriverContainerNoTrackerT) WriteMapElemValue() {}

type Encoder struct {
	encoderI
}

// ----

func (helperEncDriver[T]) newEncoderBytes(out *[]byte, h Handle) *encoder[T] {
	var c1 encoder[T]
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriver[T]) newEncoderIO(out io.Writer, h Handle) *encoder[T] {
	var c1 encoder[T]
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriver[T]) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEs[T]) (f *fastpathE[T], u reflect.Type) {
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

// ----

func (helperEncDriver[T]) encFindRtidFn(s []encRtidFn[T], rtid uintptr) (i uint, fn *encFn[T]) {
	// binary search. Adapted from sort/search.go. Use goto (not for loop) to allow inlining.
	var h uint // var h, i uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1 // avoid overflow when computing h // h = i + (j-i)/2
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

func (helperEncDriver[T]) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFn[T]) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFn[T]](v))
	}
	return
}

func (dh helperEncDriver[T]) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEs[T], checkExt bool) (fn *encFn[T]) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriver[T]) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEs[T],
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFn[T]) {
	rtid := rt2id(rt)
	var sp []encRtidFn[T] = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriver[T]) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEs[T],
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFn[T]) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFn[T]
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)
	// since this is an atomic load/store, we MUST use a different array each time,
	// else we have a data race when a store is happening simultaneously with a encFindRtidFn call.
	if sp == nil {
		sp = []encRtidFn[T]{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFn[T], len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFn[T]{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriver[T]) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEs[T],
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFn[T]) {
	fn = new(encFn[T])
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	// anything can be an extension except the built-in ones: time, raw and rawext.
	// ensure we check for these types, then if extension, before checking if
	// it implementes one of the pre-declared interfaces.

	// fi.addrEf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoder[T]).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoder[T]).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoder[T]).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoder[T]).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoder[T]).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoder[T]).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {
		//If JSON, we should check JSONMarshal before textMarshal
		fn.fe = (*encoder[T]).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoder[T]).textMarshal
		fi.addrE = ti.flagTextMarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			// by default (without using unsafe),
			// if an array is not addressable, converting from an array to a slice
			// requires an allocation (see helper_not_unsafe.go: func rvGetSlice4Array).
			//
			// (Non-addressable arrays mostly occur as keys/values from a map).
			//
			// However, fastpath functions are mostly for slices of numbers or strings,
			// which are small by definition and thus allocation should be fast/cheap in time.
			//
			// Consequently, the value of doing this quick allocation to elide the overhead cost of
			// non-optimized (not-unsafe) reflection is a fair price.
			var rtid2 uintptr
			if !ti.flagHasPkgPath { // un-named type (slice or mpa or array)
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key) // ti.key for arrays = reflect.SliceOf(ti.elem)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fe = fp[idx].encfn
				}
			} else { // named type (with underlying type of map or slice or array)
				// try to use mapping for underlying type
				xfe, xrt := dh.encFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf := xfe.encfn
					fn.fe = func(e *encoder[T], xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoder[T]).kBool
			case reflect.String:
				// Do not use different functions based on StringToRaw option, as that will statically
				// set the function for a string type, and if the Handle is modified thereafter,
				// behaviour is non-deterministic
				// i.e. DO NOT DO:
				//   if x.StringToRaw {
				//   	fn.fe = (*encoder[T]).kStringToRaw
				//   } else {
				//   	fn.fe = (*encoder[T]).kStringEnc
				//   }

				fn.fe = (*encoder[T]).kString
			case reflect.Int:
				fn.fe = (*encoder[T]).kInt
			case reflect.Int8:
				fn.fe = (*encoder[T]).kInt8
			case reflect.Int16:
				fn.fe = (*encoder[T]).kInt16
			case reflect.Int32:
				fn.fe = (*encoder[T]).kInt32
			case reflect.Int64:
				fn.fe = (*encoder[T]).kInt64
			case reflect.Uint:
				fn.fe = (*encoder[T]).kUint
			case reflect.Uint8:
				fn.fe = (*encoder[T]).kUint8
			case reflect.Uint16:
				fn.fe = (*encoder[T]).kUint16
			case reflect.Uint32:
				fn.fe = (*encoder[T]).kUint32
			case reflect.Uint64:
				fn.fe = (*encoder[T]).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoder[T]).kUintptr
			case reflect.Float32:
				fn.fe = (*encoder[T]).kFloat32
			case reflect.Float64:
				fn.fe = (*encoder[T]).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoder[T]).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoder[T]).kComplex128
			case reflect.Chan:
				fn.fe = (*encoder[T]).kChan
			case reflect.Slice:
				fn.fe = (*encoder[T]).kSlice
			case reflect.Array:
				fn.fe = (*encoder[T]).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoder[T]).kStructSimple
				} else {
					fn.fe = (*encoder[T]).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoder[T]).kMap
			case reflect.Interface:
				// encode: reflect.Interface are handled already by preEncodeValue
				fn.fe = (*encoder[T]).kErr
			default:
				// reflect.Ptr and reflect.Interface are handled already by preEncodeValue
				fn.fe = (*encoder[T]).kErr
			}
		}
	}
	return
}

// ----

func sideEncode(h Handle, p *sync.Pool, fn func(encoderI)) {
	var s encoderI
	if usePoolForSideEncode {
		s = p.Get().(encoderI)
		defer p.Put(s)
	} else {
		// initialization cycle error
		// s = NewEncoderBytes(nil, h).encoderI
		s = p.New().(encoderI)
	}
	fn(s)
}

func oneOffEncode(se encoderI, v interface{}, out *[]byte, basetype reflect.Type, ext bool) {
	se.ResetBytes(out)
	se.encodeAs(v, basetype, ext)
	se.atEndOfEncode()
	se.writerEnd()
	// e.sideEncoder(&bs)
	// e.sideEncode(v, basetype, 0)
}
