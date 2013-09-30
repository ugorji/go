// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// Contains code shared by both encode and decode.

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	// For >= mapAccessThreshold elements, map outways cost of linear search 
	//   - this was critical for reflect.Type, whose equality cost is pretty high (set to 4)
	//   - for integers, equality cost is cheap (set to 16, 32 of 64)
	// mapAccessThreshold    = 16 // 4
	
	binarySearchThreshold = 16
	structTagName         = "codec"
	
	// Support 
	//    encoding.BinaryMarshaler: MarshalBinary() (data []byte, err error)
	//    encoding.BinaryUnmarshaler: UnmarshalBinary(data []byte) error
	// This constant flag will enable or disable it.
	// 
	// Supporting this feature required a map access each time the en/decodeValue 
	// method is called to get the typeInfo and look at baseId. This caused a 
	// clear performance degradation. Some refactoring helps a portion of the loss.
	// 
	// All the band-aids we can put try to mitigate performance loss due to stack splitting:
	//    - using smaller functions to reduce func framesize
	// 
	// TODO: Look into this again later.
	supportBinaryMarshal  = true
)

type charEncoding uint8

const (
	c_RAW charEncoding = iota
	c_UTF8
	c_UTF16LE
	c_UTF16BE
	c_UTF32LE
	c_UTF32BE
)

type binaryUnmarshaler interface {
	UnmarshalBinary(data []byte) error
}

type binaryMarshaler interface {
	MarshalBinary() (data []byte, err error)
}

var (
	bigen               = binary.BigEndian
	structInfoFieldName = "_struct"

	cachedTypeInfo      = make(map[uintptr]*typeInfo, 4)
	cachedTypeInfoMutex sync.RWMutex

	nilIntfSlice     = []interface{}(nil)
	intfSliceTyp     = reflect.TypeOf(nilIntfSlice)
	intfTyp          = intfSliceTyp.Elem()
	byteSliceTyp     = reflect.TypeOf([]byte(nil))
	ptrByteSliceTyp  = reflect.TypeOf((*[]byte)(nil))
	mapStringIntfTyp = reflect.TypeOf(map[string]interface{}(nil))
	mapIntfIntfTyp   = reflect.TypeOf(map[interface{}]interface{}(nil))
	
	timeTyp          = reflect.TypeOf(time.Time{})
	ptrTimeTyp       = reflect.TypeOf((*time.Time)(nil))
	int64SliceTyp    = reflect.TypeOf([]int64(nil))
	
	timeTypId        = reflect.ValueOf(timeTyp).Pointer()
	ptrTimeTypId     = reflect.ValueOf(ptrTimeTyp).Pointer()
	byteSliceTypId   = reflect.ValueOf(byteSliceTyp).Pointer()
	
	binaryMarshalerTyp = reflect.TypeOf((*binaryMarshaler)(nil)).Elem()
	binaryUnmarshalerTyp = reflect.TypeOf((*binaryUnmarshaler)(nil)).Elem()
	
	binaryMarshalerTypId = reflect.ValueOf(binaryMarshalerTyp).Pointer()
	binaryUnmarshalerTypId = reflect.ValueOf(binaryUnmarshalerTyp).Pointer()
	
	intBitsize  uint8 = uint8(reflect.TypeOf(int(0)).Bits())
	uintBitsize uint8 = uint8(reflect.TypeOf(uint(0)).Bits())

	bsAll0x00 = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bsAll0xff = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// Handle is the interface for a specific encoding format.
// 
// Typically, a Handle is pre-configured before first time use,
// and not modified while in use. Such a pre-configured Handle
// is safe for concurrent access.
type Handle interface {
	encodeHandleI
	decodeHandleI
	newEncDriver(w encWriter) encDriver
	newDecDriver(r decReader) decDriver
}

type extTypeTagFn struct {
	rtid uintptr
	rt reflect.Type
	tag byte
	encFn func(reflect.Value) ([]byte, error)
	decFn func(reflect.Value, []byte) error
}

type extHandle map[uintptr]*extTypeTagFn

// AddExt registers an encode and decode function for a reflect.Type.
// Note that the type must be a named type, and specifically not 
// a pointer or Interface. An error is returned if that is not honored.
func (o *extHandle) AddExt(
	rt reflect.Type,
	tag byte,
	encfn func(reflect.Value) ([]byte, error),
	decfn func(reflect.Value, []byte) error,
) (err error) {
	// o is a pointer, because we may need to initialize it
	if rt.PkgPath() == "" || rt.Kind() == reflect.Interface {
		err = fmt.Errorf("codec.Handle.AddExt: Takes a named type, especially not a pointer or interface: %T", 
			reflect.Zero(rt).Interface())
		return
	}
	if o == nil {
		err = errors.New("codec.Handle.AddExt: Nil (should never happen)")
		return
	}
	rtid := reflect.ValueOf(rt).Pointer()
	if *o == nil {
		*o = make(map[uintptr]*extTypeTagFn, 4)
	}
	m := *o
	if encfn == nil || decfn == nil {
		delete(m, rtid)
	} else {
		m[rtid] = &extTypeTagFn { rtid, rt, tag, encfn, decfn }
	}
	return
}

func (o extHandle) getExt(rtid uintptr) *extTypeTagFn {
	return o[rtid]
}

func (o extHandle) getExtForTag(tag byte) *extTypeTagFn {
	for _, v := range o {
		if v.tag == tag {
			return v
		}
	}
	return nil
}

func (o extHandle) getDecodeExtForTag(tag byte) (
	rv reflect.Value, fn func(reflect.Value, []byte) error) {
	if x := o.getExtForTag(tag); x != nil {
		// ext is only registered for base
		rv = reflect.New(x.rt).Elem()
		fn = x.decFn
	}
	return
}

func (o extHandle) getDecodeExt(rtid uintptr) (tag byte, fn func(reflect.Value, []byte) error) {
	if x, ok := o[rtid]; ok {
		tag = x.tag
		fn = x.decFn
	}
	return
}

func (o extHandle) getEncodeExt(rtid uintptr) (tag byte, fn func(reflect.Value) ([]byte, error)) {
	if x, ok := o[rtid]; ok {
		tag = x.tag
		fn = x.encFn
	}
	return
}

// typeInfo keeps information about each type referenced in the encode/decode sequence.
// 
// During an encode/decode sequence, we work as below:
//   - If base is a built in type, en/decode base value
//   - If base is registered as an extension, en/decode base value
//   - If type is binary(M/Unm)arshaler, call Binary(M/Unm)arshal method
//   - Else decode appropriately based on the reflect.Kind
type typeInfo struct {
	sis       []*structFieldInfo // sorted. Used when enc/dec struct to map.
	sisp      []*structFieldInfo // unsorted. Used when enc/dec struct to array.
	// base      reflect.Type
	
	// baseId is the pointer to the base reflect.Type, after deferencing
	// the pointers. E.g. base type of ***time.Time is time.Time.
	baseId    uintptr
	baseIndir int8 // number of indirections to get to base
	m         bool // base type (T or *T) is a binaryMarshaler
	unm       bool // base type (T or *T) is a binaryUnmarshaler
	mIndir    int8 // number of indirections to get to binaryMarshaler type
	unmIndir  int8 // number of indirections to get to binaryUnmarshaler type
	toArray   bool // whether this (struct) type should be encoded as an array
}

type structFieldInfo struct {
	encName   string // encode name
	
	// only one of 'i' or 'is' can be set. If 'i' is -1, then 'is' has been set.
	
	is        []int // (recursive/embedded) field index in struct
	i         int16 // field index in struct
	omitEmpty bool  
	toArray   bool  // if field is _struct, is the toArray set?
	
	// tag       string   // tag
	// name      string   // field name
	// encNameBs []byte   // encoded name as byte stream
	// ikind     int      // kind of the field as an int i.e. int(reflect.Kind)
}

type sfiSortedByEncName []*structFieldInfo

func (p sfiSortedByEncName) Len() int           { return len(p) }
func (p sfiSortedByEncName) Less(i, j int) bool { return p[i].encName < p[j].encName }
func (p sfiSortedByEncName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (sis *typeInfo) indexForEncName(name string) int {
	//sissis := sis.sis 
	if sislen := len(sis.sis); sislen < binarySearchThreshold {
		// linear search. faster than binary search in my testing up to 16-field structs.
		// for i := 0; i < sislen; i++ {
		// 	if sis.sis[i].encName == name {
		// 		return i
		// 	}
		// }
		for i, si := range sis.sis {
			if si.encName == name {
				return i
			}
		}
	} else {
		// binary search. adapted from sort/search.go.
		h, i, j := 0, 0, sislen
		for i < j {
			h = i + (j-i)/2
			// i ≤ h < j
			if sis.sis[h].encName < name {
				i = h + 1 // preserves f(i-1) == false
			} else {
				j = h // preserves f(j) == true
			}
		}
		if i < sislen && sis.sis[i].encName == name {
			return i
		}
	}
	return -1
}

func getTypeInfo(rtid uintptr, rt reflect.Type) (sis *typeInfo) {
	var ok bool
	cachedTypeInfoMutex.RLock()
	sis, ok = cachedTypeInfo[rtid]
	cachedTypeInfoMutex.RUnlock()
	if ok {
		return
	}

	cachedTypeInfoMutex.Lock()
	defer cachedTypeInfoMutex.Unlock()
	if sis, ok = cachedTypeInfo[rtid]; ok {
		return
	}

	sis = new(typeInfo)
	
	var indir int8
	if ok, indir = implementsIntf(rt, binaryMarshalerTyp); ok {
		sis.m, sis.mIndir = true, indir
	}
	if ok, indir = implementsIntf(rt, binaryUnmarshalerTyp); ok {
		sis.unm, sis.unmIndir = true, indir
	}
	
	pt := rt
	var ptIndir int8 
	for ; pt.Kind() == reflect.Ptr; pt, ptIndir = pt.Elem(), ptIndir+1 { }
	if ptIndir == 0 {
		//sis.base = rt 
		sis.baseId = rtid
	} else {
		//sis.base = pt 
		sis.baseId = reflect.ValueOf(pt).Pointer()
		sis.baseIndir = ptIndir
	}
	
	if rt.Kind() == reflect.Struct {
		var siInfo *structFieldInfo
		if f, ok := rt.FieldByName(structInfoFieldName); ok {
			siInfo = parseStructFieldInfo(structInfoFieldName, f.Tag.Get(structTagName))
			sis.toArray = siInfo.toArray
		}
		sisp := make([]*structFieldInfo, 0, rt.NumField())
		rgetTypeInfo(rt, nil, make(map[string]bool), &sisp, siInfo)

		// // try to put all si close together
		// const tryToPutAllStructFieldInfoTogether = true
		// if tryToPutAllStructFieldInfoTogether {
		// 	sisp2 := make([]structFieldInfo, len(sisp))
		// 	for i, si := range sisp {
		// 		sisp2[i] = *si
		// 	}
		// 	for i := range sisp {
		// 		sisp[i] = &sisp2[i]
		// 	}
		// }
		
		sis.sisp = make([]*structFieldInfo, len(sisp))
		sis.sis = make([]*structFieldInfo, len(sisp))
		copy(sis.sisp, sisp)
		sort.Sort(sfiSortedByEncName(sisp))
		copy(sis.sis, sisp)
	}
	// sis = sisp
	cachedTypeInfo[rtid] = sis
	return
}

func rgetTypeInfo(rt reflect.Type, indexstack []int, fnameToHastag map[string]bool,
	sis *[]*structFieldInfo, siInfo *structFieldInfo,
) {
	for j := 0; j < rt.NumField(); j++ {
		f := rt.Field(j)
		stag := f.Tag.Get(structTagName)
		if stag == "-" {
			continue
		}
		if r1, _ := utf8.DecodeRuneInString(f.Name); r1 == utf8.RuneError || !unicode.IsUpper(r1) {
			continue
		}
		if f.Anonymous {
			//if anonymous, inline it if there is no struct tag, else treat as regular field
			if stag == "" {
				indexstack2 := append(append([]int(nil), indexstack...), j)
				rgetTypeInfo(f.Type, indexstack2, fnameToHastag, sis, siInfo)
				continue
			}
		}
		// do not let fields with same name in embedded structs override field at higher level.
		// this must be done after anonymous check, to allow anonymous field 
		// still include their child fields
		if _, ok := fnameToHastag[f.Name]; ok {
			continue
		}
		si := parseStructFieldInfo(f.Name, stag)
		// si.ikind = int(f.Type.Kind())
		if len(indexstack) == 0 {
			si.i = int16(j)
		} else {
			si.i = -1
			si.is = append(append([]int(nil), indexstack...), j)
		}

		if siInfo != nil {
			if siInfo.omitEmpty {
				si.omitEmpty = true
			}
		}
		*sis = append(*sis, si)
		fnameToHastag[f.Name] = stag != ""
	}
}

func parseStructFieldInfo(fname string, stag string) *structFieldInfo {
	if fname == "" {
		panic("parseStructFieldInfo: No Field Name")
	}
	si := structFieldInfo{
		// name: fname,
		encName: fname,
		// tag: stag,
	}

	if stag != "" {
		for i, s := range strings.Split(stag, ",") {
			if i == 0 {
				if s != "" {
					si.encName = s
				}
			} else {
				switch s {
				case "omitempty":
					si.omitEmpty = true
				case "toarray":
					si.toArray = true
				}
			}
		}
	}
	// si.encNameBs = []byte(si.encName)
	return &si
}

func panicToErr(err *error) {
	if x := recover(); x != nil {
		//debug.PrintStack()
		panicValToErr(x, err)
	}
}

func doPanic(tag string, format string, params ...interface{}) {
	params2 := make([]interface{}, len(params)+1)
	params2[0] = tag
	copy(params2[1:], params)
	panic(fmt.Errorf("%s: "+format, params2...))
}

