// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !safe && !codec.safe && !appengine && go1.9

package codec

import "unsafe"

func encFromRtidFnSlice[E encDriver](v unsafe.Pointer) (s []encRtidFn[E]) {
	if v != nil {
		s = *(*[]encRtidFn[E])(v)
	}
	return
}

func encToRtidFnSlice[E encDriver](s *[]encRtidFn[E]) unsafe.Pointer {
	return unsafe.Pointer(s)
}

func decFromRtidFnSlice[D decDriver](v unsafe.Pointer) (s []decRtidFn[D]) {
	if v != nil {
		s = *(*[]decRtidFn[D])(v)
	}
	return
}

func decToRtidFnSlice[D decDriver](s *[]decRtidFn[D]) unsafe.Pointer {
	return unsafe.Pointer(s)
}
