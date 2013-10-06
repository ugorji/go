// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// Contains code shared by both encode and decode.

import (
	"encoding/binary"
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

	// Each Encoder or Decoder uses a cache of functions based on conditionals,
	// so that the conditionals are not run every time.
	// 
	// Either a map or a slice is used to keep track of the functions.
	// The map is more natural, but has a higher cost than a slice/array.
	// This flag (useMapForCodecCache) controls which is used.
	useMapForCodecCache = false
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
	mapStringIntfTyp = reflect.TypeOf(map[string]interface{}(nil))
	mapIntfIntfTyp   = reflect.TypeOf(map[interface{}]interface{}(nil))
	
	timeTyp          = reflect.TypeOf(time.Time{})
	int64SliceTyp    = reflect.TypeOf([]int64(nil))
	rawExtTyp        = reflect.TypeOf(RawExt{})
	
	timeTypId        = reflect.ValueOf(timeTyp).Pointer()
	byteSliceTypId   = reflect.ValueOf(byteSliceTyp).Pointer()
	rawExtTypId      = reflect.ValueOf(rawExtTyp).Pointer()
	
	binaryMarshalerTyp = reflect.TypeOf((*binaryMarshaler)(nil)).Elem()
	binaryUnmarshalerTyp = reflect.TypeOf((*binaryUnmarshaler)(nil)).Elem()
	
	binaryMarshalerTypId = reflect.ValueOf(binaryMarshalerTyp).Pointer()
	binaryUnmarshalerTypId = reflect.ValueOf(binaryUnmarshalerTyp).Pointer()
	
	intBitsize  uint8 = uint8(reflect.TypeOf(int(0)).Bits())
	uintBitsize uint8 = uint8(reflect.TypeOf(uint(0)).Bits())

	bsAll0x00 = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bsAll0xff = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// The RawExt type represents raw unprocessed extension data. 
type RawExt struct {
	Tag byte
	Data []byte
}

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

type extHandle []*extTypeTagFn

// AddExt registers an encode and decode function for a reflect.Type.
// Note that the type must be a named type, and specifically not 
// a pointer or Interface. An error is returned if that is not honored.
// 
// To Deregister an ext, call AddExt with 0 tag, nil encfn and nil decfn.
func (o *extHandle) AddExt(
	rt reflect.Type,
	tag byte,
	encfn func(reflect.Value) ([]byte, error),
	decfn func(reflect.Value, []byte) error,
) (err error) {
	// o is a pointer, because we may need to initialize it
	if rt.PkgPath() == "" || rt.Kind() == reflect.Interface {
		err = fmt.Errorf("codec.Handle.AddExt: Takes named type, especially not a pointer or interface: %T", 
			reflect.Zero(rt).Interface())
		return
	}
	
	// o cannot be nil, since it is always embedded in a Handle. 
	// if nil, let it panic.
	// if o == nil {
	// 	err = errors.New("codec.Handle.AddExt: extHandle cannot be a nil pointer.")
	// 	return
	// }
	
	rtid := reflect.ValueOf(rt).Pointer()
	for _, v := range *o {
		if v.rtid == rtid {
			v.tag, v.encFn, v.decFn = tag, encfn, decfn
			return
		}
	}
	
	*o = append(*o, &extTypeTagFn { rtid, rt, tag, encfn, decfn })
	return
}

func (o extHandle) getExt(rtid uintptr) *extTypeTagFn {
	for _, v := range o {
		if v.rtid == rtid {
			return v
		}
	}
	return nil
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
	if x := o.getExt(rtid); x != nil {
		tag = x.tag
		fn = x.decFn
	}
	return
}

func (o extHandle) getEncodeExt(rtid uintptr) (tag byte, fn func(reflect.Value) ([]byte, error)) {
	if x := o.getExt(rtid); x != nil {
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
	sfi       []*structFieldInfo // sorted. Used when enc/dec struct to map.
	sfip      []*structFieldInfo // unsorted. Used when enc/dec struct to array.
	
	rt        reflect.Type
	rtid      uintptr
	
	// baseId gives pointer to the base reflect.Type, after deferencing
	// the pointers. E.g. base type of ***time.Time is time.Time.
	base      reflect.Type
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

func (ti *typeInfo) indexForEncName(name string) int {
	//tisfi := ti.sfi 
	if sfilen := len(ti.sfi); sfilen < binarySearchThreshold {
		// linear search. faster than binary search in my testing up to 16-field structs.
		// for i := 0; i < sfilen; i++ {
		// 	if ti.sfi[i].encName == name {
		// 		return i
		// 	}
		// }
		for i, si := range ti.sfi {
			if si.encName == name {
				return i
			}
		}
	} else {
		// binary search. adapted from sort/search.go.
		h, i, j := 0, 0, sfilen
		for i < j {
			h = i + (j-i)/2
			// i ≤ h < j
			if ti.sfi[h].encName < name {
				i = h + 1 // preserves f(i-1) == false
			} else {
				j = h // preserves f(j) == true
			}
		}
		if i < sfilen && ti.sfi[i].encName == name {
			return i
		}
	}
	return -1
}

func getTypeInfo(rtid uintptr, rt reflect.Type) (pti *typeInfo) {
	var ok bool
	cachedTypeInfoMutex.RLock()
	pti, ok = cachedTypeInfo[rtid]
	cachedTypeInfoMutex.RUnlock()
	if ok {
		return
	}

	cachedTypeInfoMutex.Lock()
	defer cachedTypeInfoMutex.Unlock()
	if pti, ok = cachedTypeInfo[rtid]; ok {
		return
	}

	ti := typeInfo { rt: rt, rtid: rtid }
	pti = &ti
	
	var indir int8
	if ok, indir = implementsIntf(rt, binaryMarshalerTyp); ok {
		ti.m, ti.mIndir = true, indir
	}
	if ok, indir = implementsIntf(rt, binaryUnmarshalerTyp); ok {
		ti.unm, ti.unmIndir = true, indir
	}
	
	pt := rt
	var ptIndir int8 
	// for ; pt.Kind() == reflect.Ptr; pt, ptIndir = pt.Elem(), ptIndir+1 { }
	for pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
		ptIndir++
	}
	if ptIndir == 0 {
		ti.base = rt 
		ti.baseId = rtid
	} else {
		ti.base = pt 
		ti.baseId = reflect.ValueOf(pt).Pointer()
		ti.baseIndir = ptIndir
	}
	
	if rt.Kind() == reflect.Struct {
		var siInfo *structFieldInfo
		if f, ok := rt.FieldByName(structInfoFieldName); ok {
			siInfo = parseStructFieldInfo(structInfoFieldName, f.Tag.Get(structTagName))
			ti.toArray = siInfo.toArray
		}
		sfip := make([]*structFieldInfo, 0, rt.NumField())
		rgetTypeInfo(rt, nil, make(map[string]bool), &sfip, siInfo)

		// // try to put all si close together
		// const tryToPutAllStructFieldInfoTogether = true
		// if tryToPutAllStructFieldInfoTogether {
		// 	sfip2 := make([]structFieldInfo, len(sfip))
		// 	for i, si := range sfip {
		// 		sfip2[i] = *si
		// 	}
		// 	for i := range sfip {
		// 		sfip[i] = &sfip2[i]
		// 	}
		// }
		
		ti.sfip = make([]*structFieldInfo, len(sfip))
		ti.sfi = make([]*structFieldInfo, len(sfip))
		copy(ti.sfip, sfip)
		sort.Sort(sfiSortedByEncName(sfip))
		copy(ti.sfi, sfip)
	}
	// sfi = sfip
	cachedTypeInfo[rtid] = pti
	return
}

func rgetTypeInfo(rt reflect.Type, indexstack []int, fnameToHastag map[string]bool,
	sfi *[]*structFieldInfo, siInfo *structFieldInfo,
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
				rgetTypeInfo(f.Type, indexstack2, fnameToHastag, sfi, siInfo)
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
		*sfi = append(*sfi, si)
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

