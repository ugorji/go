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
	// For >= 4 elements, map outways cost of linear search (especially for reflect.Type)
	mapAccessThreshold    = 4
	binarySearchThreshold = 16
	structTagName         = "codec"
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

var (
	bigen               = binary.BigEndian
	structInfoFieldName = "_struct"

	cachedStructFieldInfos      = make(map[reflect.Type]structFieldInfos, 4)
	cachedStructFieldInfosMutex sync.RWMutex

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

	intBitsize  uint8 = uint8(reflect.TypeOf(int(0)).Bits())
	uintBitsize uint8 = uint8(reflect.TypeOf(uint(0)).Bits())
)

type encdecHandle struct {
	encHandle
	decHandle
}

func (o *encdecHandle) AddExt(
	rt reflect.Type,
	tag byte,
	encfn func(reflect.Value) ([]byte, error),
	decfn func(reflect.Value, []byte) error,
) {
	o.addEncodeExt(rt, tag, encfn)
	o.addDecodeExt(rt, tag, decfn)
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

type structFieldInfo struct {
	encName   string // encode name
	is        []int
	i         int16 // field index in struct
	omitEmpty bool
	// tag       string   // tag
	// name      string   // field name
	// encNameBs []byte   // encoded name as byte stream
	// ikind     int      // kind of the field as an int i.e. int(reflect.Kind)
}

type structFieldInfos []structFieldInfo

type sfiSortedByEncName []*structFieldInfo

func (p sfiSortedByEncName) Len() int {
	return len(p)
}

func (p sfiSortedByEncName) Less(i, j int) bool {
	return p[i].encName < p[j].encName
}

func (p sfiSortedByEncName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (sis structFieldInfos) indexForEncName(name string) int {
	sislen := len(sis)
	if sislen < binarySearchThreshold {
		// linear search. faster than binary search in my testing up to 16-field structs.
		for i := 0; i < sislen; i++ {
			if sis[i].encName == name {
				return i
			}
		}
	} else {
		// binary search. adapted from sort/search.go.
		h, i, j := 0, 0, sislen
		for i < j {
			h = i + (j-i)/2
			// i ≤ h < j
			if sis[h].encName < name {
				i = h + 1 // preserves f(i-1) == false
			} else {
				j = h // preserves f(j) == true
			}
		}
		if i < sislen && sis[i].encName == name {
			return i
		}
	}
	return -1
}

func getStructFieldInfos(rt reflect.Type) (sis structFieldInfos) {
	cachedStructFieldInfosMutex.RLock()
	sis, ok := cachedStructFieldInfos[rt]
	cachedStructFieldInfosMutex.RUnlock()
	if ok {
		return
	}

	cachedStructFieldInfosMutex.Lock()
	defer cachedStructFieldInfosMutex.Unlock()

	var siInfo *structFieldInfo
	if f, ok := rt.FieldByName(structInfoFieldName); ok {
		siInfo = parseStructFieldInfo(structInfoFieldName, f.Tag.Get(structTagName))
	}
	sisp := make([]*structFieldInfo, 0, rt.NumField())
	rgetStructFieldInfos(rt, nil, make(map[string]bool), &sisp, siInfo)
	sort.Sort(sfiSortedByEncName(sisp))

	lsis := len(sisp)
	sis = make([]structFieldInfo, lsis)
	for i := 0; i < lsis; i++ {
		sis[i] = *sisp[i]
	}
	// sis = sisp
	cachedStructFieldInfos[rt] = sis
	return
}

func rgetStructFieldInfos(rt reflect.Type, indexstack []int, fnameToHastag map[string]bool,
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
				rgetStructFieldInfos(f.Type, indexstack2, fnameToHastag, sis, siInfo)
				continue
			}
		}
		//do not let fields with same name in embedded structs override field at higher level.
		//this must be done after anonymous check, to allow anonymous field still include their child fields
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
				if s == "omitempty" {
					si.omitEmpty = true
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

//--------------------------------------------------

// // This implements the util.Codec interface
// type Codec struct {
// 	H Handle
// }

// func (x Codec) Encode(w io.Writer, v interface{}) error {
// 	return NewEncoder(w, x.H).Encode(v)
// }

// func (x Codec) EncodeBytes(out *[]byte, v interface{}) error {
// 	return NewEncoderBytes(out, x.H).Encode(v)
// }

// func (x Codec) Decode(r io.Reader, v interface{}) error {
// 	return NewDecoder(r, x.H).Decode(v)
// }

// func (x Codec) DecodeBytes(in []byte, v interface{}) error {
// 	return NewDecoderBytes(in, x.H).Decode(v)
// }
