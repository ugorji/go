// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build !go1.12

package codec

import "reflect"

type mapIter struct {
	m, k, v  reflect.Value
	keys     []reflect.Value
	j        int
	kOk, vOk bool
	values   bool
}

func (t *mapIter) Next() (r bool) {
	t.j++
	r = t.j < len(t.keys)
	if r {
		if t.kOk {
			t.k.Set(t.keys[t.j])
		}
		if t.vOk {
			t.v.Set(t.rv.MapIndex(t.keys[t.j]))
		}
	}
}

func (t *mapIter) Key() reflect.Value {
	if t.kOk {
		return t.k
	}
	return t.keys[t.j]
}

func (t *mapIter) Value() (r reflect.Value) {
	if !t.values {
		return
	}
	if t.vOk {
		return t.v
	}
	return t.rv.MapIndex(t.keys[t.j])
}

func mapRange(m, k, v reflect.Value, values bool) *mapIter {
	return &mapIter{m: m, k: k, v: v,
		kOk: k.CanSet(), vOk: values && v.CanSet(),
		keys: rv.MapKeys(), j: -1,
		values: values,
	}
}

func mapIndex(m, k, v reflect.Value) (vv reflect.Value) {
	vv = m.MapIndex(k)
	if vv.IsValid() && v.CanSet() {
		v.Set(vv)
	}
	return
}

// return an addressable reflect value that can be used in mapRange and mapIndex operations.
func mapAddressableRV(t reflect.Type) (r reflect.Value) {
	return // reflect.New(t).Elem()
}
