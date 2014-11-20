// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// This file contains values used by tests and benchmarks.

import (
	"math"
	"time"
)

var testStrucTime = time.Date(2012, 2, 2, 2, 2, 2, 2000, time.UTC).UTC()

type AnonInTestStruc struct {
	AS        string
	AI64      int64
	AI16      int16
	AUi64     uint64
	ASslice   []string
	AI64slice []int64
}

type TestStruc struct {
	S    string
	I64  int64
	I16  int16
	Ui64 uint64
	Ui8  uint8
	B    bool
	By   byte

	Sslice    []string
	I64slice  []int64
	I16slice  []int16
	Ui64slice []uint64
	Ui8slice  []uint8
	Bslice    []bool
	Byslice   []byte

	Islice    []interface{}
	Iptrslice []*int64

	AnonInTestStruc

	//M map[interface{}]interface{}  `json:"-",bson:"-"`
	Ms    map[string]interface{}
	Msi64 map[string]int64

	Nintf      interface{} //don't set this, so we can test for nil
	T          time.Time
	Nmap       map[string]bool //don't set this, so we can test for nil
	Nslice     []byte          //don't set this, so we can test for nil
	Nint64     *int64          //don't set this, so we can test for nil
	Mtsptr     map[string]*TestStruc
	Mts        map[string]TestStruc
	Its        []*TestStruc
	Nteststruc *TestStruc
}

func newTestStruc(depth int, bench bool) (ts *TestStruc) {
	var i64a, i64b, i64c, i64d int64 = 64, 6464, 646464, 64646464

	ts = &TestStruc{
		S:    "some string",
		I64:  math.MaxInt64 * 2 / 3, // 64,
		I16:  16,
		Ui64: uint64(int64(math.MaxInt64 * 2 / 3)), // 64, //don't use MaxUint64, as bson can't write it
		Ui8:  160,
		B:    true,
		By:   5,

		Sslice:    []string{"one", "two", "three"},
		I64slice:  []int64{1, 2, 3},
		I16slice:  []int16{4, 5, 6},
		Ui64slice: []uint64{137, 138, 139},
		Ui8slice:  []uint8{210, 211, 212},
		Bslice:    []bool{true, false, true, false},
		Byslice:   []byte{13, 14, 15},

		Islice: []interface{}{"true", true, "no", false, uint64(288), float64(0.4)},

		Ms: map[string]interface{}{
			"true":     "true",
			"int64(9)": false,
		},
		Msi64: map[string]int64{
			"one": 1,
			"two": 2,
		},
		T: testStrucTime,
		AnonInTestStruc: AnonInTestStruc{
			AS:        "A-String",
			AI64:      64,
			AI16:      16,
			AUi64:     64,
			ASslice:   []string{"Aone", "Atwo", "Athree"},
			AI64slice: []int64{1, 2, 3},
		},
	}
	//For benchmarks, some things will not work.
	if !bench {
		//json and bson require string keys in maps
		//ts.M = map[interface{}]interface{}{
		//	true: "true",
		//	int8(9): false,
		//}
		//gob cannot encode nil in element in array (encodeArray: nil element)
		ts.Iptrslice = []*int64{nil, &i64a, nil, &i64b, nil, &i64c, nil, &i64d, nil}
		// ts.Iptrslice = nil
	}
	if depth > 0 {
		depth--
		if ts.Mtsptr == nil {
			ts.Mtsptr = make(map[string]*TestStruc)
		}
		if ts.Mts == nil {
			ts.Mts = make(map[string]TestStruc)
		}
		ts.Mtsptr["0"] = newTestStruc(depth, bench)
		ts.Mts["0"] = *(ts.Mtsptr["0"])
		ts.Its = append(ts.Its, ts.Mtsptr["0"])
	}
	return
}
