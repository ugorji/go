// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package internal

// All non-std package dependencies related to testing live in this file,
// so porting to different environment is easy (just update functions).

import (
	"errors"
	"reflect"
)

// --- these functions are used by both benchmarks and tests

var ErrDeepEqualNotMatch = errors.New("not Match")

func DeepEqual(v1, v2 interface{}) (err error) {
	if !reflect.DeepEqual(v1, v2) {
		err = ErrDeepEqualNotMatch
	}
	return
}

func ApproxDataSize(rv reflect.Value) (sum int) {
	switch rk := rv.Kind(); rk {
	case reflect.Invalid:
	case reflect.Ptr, reflect.Interface:
		sum += int(rv.Type().Size())
		sum += ApproxDataSize(rv.Elem())
	case reflect.Slice:
		sum += int(rv.Type().Size())
		for j := 0; j < rv.Len(); j++ {
			sum += ApproxDataSize(rv.Index(j))
		}
	case reflect.String:
		sum += int(rv.Type().Size())
		sum += rv.Len()
	case reflect.Map:
		sum += int(rv.Type().Size())
		for _, mk := range rv.MapKeys() {
			sum += ApproxDataSize(mk)
			sum += ApproxDataSize(rv.MapIndex(mk))
		}
	case reflect.Struct:
		//struct size already includes the full data size.
		//sum += int(rv.Type().Size())
		for j := 0; j < rv.NumField(); j++ {
			sum += ApproxDataSize(rv.Field(j))
		}
	default:
		//pure value types
		sum += int(rv.Type().Size())
	}
	return
}
