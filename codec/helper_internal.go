// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// All non-std package dependencies live in this file,
// so porting to different environment is easy (just update functions).

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

var (
	raisePanicAfterRecover = false
	debugging              = true
)

func panicValToErr(panicVal interface{}, err *error) {
	switch xerr := panicVal.(type) {
	case error:
		*err = xerr
	case string:
		*err = errors.New(xerr)
	default:
		*err = fmt.Errorf("%v", panicVal)
	}
	if raisePanicAfterRecover {
		panic(panicVal)
	}
	return
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
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
		return v.IsNil()
	}
	return false
}

func debugf(format string, args ...interface{}) {
	if debugging {
		if len(format) == 0 || format[len(format)-1] != '\n' {
			format = format + "\n"
		}
		fmt.Printf(format, args...)
	}
}

func pruneSignExt(v []byte) (n int) {
	l := len(v)
	if l < 2 {
		return
	}
	if v[0] == 0 {
		n2 := n + 1
		for v[n] == 0 && n2 < l && (v[n2]&(1<<7) == 0) {
			n++
			n2++
		}
		return
	}
	if v[0] == 0xff {
		n2 := n + 1
		for v[n] == 0xff && n2 < l && (v[n2]&(1<<7) != 0) {
			n++
			n2++
		}
		return
	}
	return
}

func implementsIntf(typ, iTyp reflect.Type) (success bool, indir int8) {
	if typ == nil {
		return
	}
	rt := typ
	// The type might be a pointer and we need to keep
	// dereferencing to the base type until we find an implementation.
	for {
		if rt.Implements(iTyp) {
			return true, indir
		}
		if p := rt; p.Kind() == reflect.Ptr {
			indir++
			if indir >= math.MaxInt8 { // insane number of indirections
				return false, 0
			}
			rt = p.Elem()
			continue
		}
		break
	}
	// No luck yet, but if this is a base type (non-pointer), the pointer might satisfy.
	if typ.Kind() != reflect.Ptr {
		// Not a pointer, but does the pointer work?
		if reflect.PtrTo(typ).Implements(iTyp) {
			return true, -1
		}
	}
	return false, 0
}
