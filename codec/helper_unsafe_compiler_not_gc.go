// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !safe && !codec.safe && !appengine && go1.9 && !gc
// +build !safe,!codec.safe,!appengine,go1.9,!gc

package codec

var unsafeZeroArr [1024]byte

// func unsafeNew(t reflect.Type, typ unsafe.Pointer) unsafe.Pointer {
// 	rv := reflect.New(t)
// 	return ((*unsafeReflectValue)(unsafe.Pointer(&rv))).ptr
// }
