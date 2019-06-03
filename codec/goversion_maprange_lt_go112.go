// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build !go1.12

package codec

import "reflect"

type mapRanger struct {
	rv  reflect.Value
	mks []reflect.Value
	j   int
}

func (x *mapRanger) Next() bool {
	x.j++
	return x.j < len(x.mks)
}
func (x *mapRanger) Key() reflect.Value {
	return x.mks[x.j]
}
func (x *mapRanger) Value() reflect.Value {
	return x.rv.MapIndex(x.mks[x.j])
}

func mapRange(rv reflect.Value) (g *mapRanger) {
	return &mapRanger{rv: rv, mks: rv.MapKeys(), j: -1}
}
