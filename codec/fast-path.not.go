// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build notfastpath || codec.notfastpath

package codec

import "reflect"

const fastpathEnabled = false

// The generated fast-path code is very large, and adds a few seconds to the build time.
// This causes test execution, execution of small tools which use codec, etc
// to take a long time.
//
// To mitigate, we now support the notfastpath tag.
// This tag disables fastpath during build, allowing for faster build, test execution,
// short-program runs, etc.

// type fastpathT struct{}
type fastpathE[E encDriver] struct {
	rt    reflect.Type
	encfn func(*encoder[E], *encFnInfo, reflect.Value)
}
type fastpathD[D decDriver] struct {
	rt    reflect.Type
	decfn func(*decoder[D], *decFnInfo, reflect.Value)
}
type fastpathEs[T encDriver] [0]fastpathE[T]
type fastpathDs[T decDriver] [0]fastpathD[T]

func fastpathDecodeTypeSwitch[T decDriver](iv interface{}, d *decoder[T]) bool { return false }
func fastpathEncodeTypeSwitch[T encDriver](iv interface{}, e *encoder[T]) bool { return false }

// func fastpathEncodeTypeSwitchSlice(iv interface{}, e *Encoder) bool { return false }
// func fastpathEncodeTypeSwitchMap(iv interface{}, e *Encoder) bool   { return false }

func fastpathDecodeSetZeroTypeSwitch(iv interface{}) bool { return false }

func fastpathAvIndex(rtid uintptr) (uint, bool) { return 0, false }

func fastpathEList[T encDriver]() (v *fastpathEs[T]) { return }
func fastpathDList[T decDriver]() (v *fastpathDs[T]) { return }

type fastpathRtRtid struct {
	rtid uintptr
	rt   reflect.Type
}
type fastpathARtRtid [0]fastpathRtRtid

var fastpathAvRtRtid fastpathARtRtid

//var fastpathAV fastpathA
// var fastpathTV fastpathT
