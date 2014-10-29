// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"strconv"
	"testing"
)

func TestCborIndefiniteLength(t *testing.T) {
	oldMapType := testCborH.MapType
	defer func() {
		testCborH.MapType = oldMapType
	}()
	testCborH.MapType = testMapStrIntfTyp
	// var (
	// 	M1 map[string][]byte
	// 	M2 map[uint64]bool
	// 	L1 []interface{}
	// 	S1 []string
	// 	B1 []byte
	// )
	var v, vv interface{}
	// define it (v), encode it using indefinite lengths, decode it (vv), compare v to vv
	v = map[string]interface{}{
		"one-byte-key":   []byte{1, 2, 3, 4, 5, 6},
		"two-string-key": "two-value",
		"three-list-key": []interface{}{true, false, uint64(1), int64(-1)},
	}
	var buf bytes.Buffer
	// buf.Reset()
	e := NewEncoder(&buf, testCborH)
	buf.WriteByte(cborBdIndefiniteMap)
	//----
	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("one-")
	e.MustEncode("byte-")
	e.MustEncode("key")
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdIndefiniteBytes)
	e.MustEncode([]byte{1, 2, 3})
	e.MustEncode([]byte{4, 5, 6})
	buf.WriteByte(cborBdBreak)

	//----
	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("two-")
	e.MustEncode("string-")
	e.MustEncode("key")
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode([]byte("two-")) // encode as bytes, to check robustness of code
	e.MustEncode([]byte("value"))
	buf.WriteByte(cborBdBreak)

	//----
	buf.WriteByte(cborBdIndefiniteString)
	e.MustEncode("three-")
	e.MustEncode("list-")
	e.MustEncode("key")
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdIndefiniteArray)
	e.MustEncode(true)
	e.MustEncode(false)
	e.MustEncode(uint64(1))
	e.MustEncode(int64(-1))
	buf.WriteByte(cborBdBreak)

	buf.WriteByte(cborBdBreak) // close map

	NewDecoderBytes(buf.Bytes(), testCborH).MustDecode(&vv)
	if err := deepEqual(v, vv); err != nil {
		logT(t, "-------- Before and After marshal do not match: Error: %v", err)
		logT(t, "    ....... GOLDEN:  (%T) %#v", v, v)
		logT(t, "    ....... DECODED: (%T) %#v", vv, vv)
		failT(t)
	}
}

type testCborGolden struct {
	Base64    string      `json:"cbor"`
	Hex       string      `json:"hex"`
	Roundtrip bool        `json:"roundtrip"`
	Decoded   interface{} `json:"decoded"`
}

// __TestCborGoldens is disabled because it includes numbers outside the range of int64/uint64
// and it doesn't support diagnostic checking.
func __TestCborGoldens(t *testing.T) {
	// decode test-cbor-goldens.json into a list of []*testCborGolden
	// for each one,
	// - decode hex into []byte bs
	// - decode bs into interface{} v
	// - compare both using deepequal
	// - for any miss, record it
	var gs []*testCborGolden
	f, err := os.Open("test-cbor-goldens.json")
	if err != nil {
		logT(t, "error opening test-cbor-goldens.json: %v", err)
		failT(t)
	}
	d := json.NewDecoder(f)
	d.UseNumber()
	err = d.Decode(&gs)
	if err != nil {
		logT(t, "error json decoding test-cbor-goldens.json: %v", err)
		failT(t)
	}
	for i, g := range gs {
		bs, err := hex.DecodeString(g.Hex)
		if err != nil {
			logT(t, "[%v] error hex decoding %s [%v]: %v", i, g.Hex, err)
			failT(t)
		}
		var v interface{}
		NewDecoderBytes(bs, testCborH).MustDecode(&v)
		if _, ok := v.(RawExt); ok {
			continue
		}
		if x, ok := g.Decoded.(json.Number); ok {
			var doContinue bool
			js := x.String()

			switch v2 := v.(type) {
			case float64:
				xx, err := strconv.ParseFloat(js, 64)
				if err != nil {
					logT(t, "[%v] cannot parse decoded value as float64 (expect %v): %v", i, v2, err)
					failT(t)
				}
				if xx != v2 {
					logT(t, "[%v] float64 value mismatch: golden: %v, decoded: %v", i, xx, v2)
					failT(t)
				}
				doContinue = true
			case int64:
				xx, err := strconv.ParseInt(js, 10, 64)
				if err != nil {
					logT(t, "[%v] cannot parse decoded value as int64 (expect %v): %v", i, v2, err)
					failT(t)
				}
				if xx != v2 {
					logT(t, "[%v] int64 value mismatch: golden: %v, decoded: %v", i, xx, v2)
					failT(t)
				}
				doContinue = true
			case uint64:
				xx, err := strconv.ParseUint(js, 10, 64)
				if err != nil {
					logT(t, "[%v] cannot parse decoded value as uint64 (expect %v): %v", i, v2, err)
					failT(t)
				}
				if xx != v2 {
					logT(t, "[%v] uint64 value mismatch: golden: %v, decoded: %v", i, xx, v2)
					failT(t)
				}
				doContinue = true
			}
			if doContinue {
				continue
			}
		}
		if err = deepEqual(g.Decoded, v); err != nil {
			logT(t, "[%v] deepEqual error: %v", i, err)
			logT(t, "    ....... GOLDEN:  (%T) %#v", g.Decoded, g.Decoded)
			logT(t, "    ....... DECODED: (%T) %#v", v, v)
			failT(t)
		}
	}
}
