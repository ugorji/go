// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build !safe
// +build !appengine
// +build go1.12
// +build !go1.13

package codec

import "unsafe"

//go:linkname mapiterelem reflect.mapitervalue
//go:noescape
func mapiterelem(it unsafe.Pointer) (elem unsafe.Pointer)
