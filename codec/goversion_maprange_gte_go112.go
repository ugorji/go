// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build go1.12

// add support for reflect.Value.MapRange.

package codec

import "reflect"

func mapRange(rv reflect.Value) *reflect.MapIter {
	return rv.MapRange()
}
