// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !go1.9 || safe || codec.safe || appengine

package codec

func encFromRtidFnSlice[E encDriver](v interface{}) (s []encRtidFn[E]) {
	if v != nil {
		s = *(v.(*[]encRtidFn[E]))
	}
	return
}

func encToRtidFnSlice[E encDriver](s *[]encRtidFn[E]) interface{} {
	return s
}

func decFromRtidFnSlice[D decDriver](v interface{}) (s []decRtidFn[D]) {
	if v != nil {
		s = *(v.(*[]decRtidFn[D]))
	}
	return
}

func decToRtidFnSlice[D decDriver](s *[]decRtidFn[D]) interface{} {
	return s
}
