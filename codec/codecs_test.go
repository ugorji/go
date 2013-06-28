// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// Test works by using a slice of interfaces.
// It can test for encoding/decoding into/from a nil interface{}
// or passing the object to encode/decode into.
//
// There are basically 2 main tests here.
// First test internally encodes and decodes things and verifies that
// the artifact was as expected.
// Second test will use python msgpack to create a bunch of golden files,
// read those files, and compare them to what it should be. It then
// writes those files back out and compares the byte streams.
//
// Taken together, the tests are pretty extensive.

// Some hints:
// - python msgpack encodes positive numbers as uints, so use uints below
//   for positive numbers.

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type testVerifyArg int

const (
	testVerifyMapTypeSame testVerifyArg = iota
	testVerifyMapTypeStrIntf
	testVerifyMapTypeIntfIntf
)

var (
	testInitDebug   bool
	testUseIoEncDec bool
	_                           = fmt.Printf
	skipVerifyVal   interface{} = &(struct{}{})
	timeLoc                     = time.FixedZone("UTC-08:00", -8*60*60)         //time.UTC
	timeToCompare               = time.Date(2012, 2, 2, 2, 2, 2, 2000, timeLoc) //time.Time{}
	//"2012-02-02T02:02:02.000002000Z" //1328148122000002
	timeToCompareAs    interface{}   = timeToCompare.UnixNano()
	table              []interface{} // main items we encode
	tableVerify        []interface{} // we verify encoded things against this after decode
	tableTestNilVerify []interface{} // for nil interface, use this to verify (rules are different)
	tablePythonVerify  []interface{} // for verifying for python, since Python sometimes
	// will encode a float32 as float64, or large int as uint
	testRpcInt   = new(TestRpcInt)
	testMsgpackH = &MsgpackHandle{}
	testBincH    = &BincHandle{}
)

func init() {
	// delete(testDecOpts.ExtFuncs, timeTyp)
	flag.BoolVar(&testInitDebug, "tdbg", false, "Test Debug")
	flag.BoolVar(&testUseIoEncDec, "tio", false, "Use IO Reader/Writer for Marshal/Unmarshal")
	flag.Parse()
	gob.Register(new(TestStruc))
	if testInitDebug {
		ts0 := newTestStruc(2, false)
		fmt.Printf("====> depth: %v, ts: %#v\n", 2, ts0)
	}

	testMsgpackH.AddExt(byteSliceTyp, 0, testMsgpackH.BinaryEncodeExt, testMsgpackH.BinaryDecodeExt)
	testMsgpackH.AddExt(timeTyp, 1, testMsgpackH.TimeEncodeExt, testMsgpackH.TimeDecodeExt)
}

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

type TestRpcInt struct {
	i int
}

func (r *TestRpcInt) Update(n int, res *int) error      { r.i = n; *res = r.i; return nil }
func (r *TestRpcInt) Square(ignore int, res *int) error { *res = r.i * r.i; return nil }
func (r *TestRpcInt) Mult(n int, res *int) error        { *res = r.i * n; return nil }

func testVerifyVal(v interface{}, arg testVerifyArg) (v2 interface{}) {
	switch iv := v.(type) {
	case int8:
		v2 = int64(iv)
	case int16:
		v2 = int64(iv)
	case int32:
		v2 = int64(iv)
	case int64:
		v2 = int64(iv)
	case uint8:
		v2 = uint64(iv)
	case uint16:
		v2 = uint64(iv)
	case uint32:
		v2 = uint64(iv)
	case uint64:
		v2 = uint64(iv)
	case float32:
		v2 = float64(iv)
	case float64:
		v2 = float64(iv)
	case []interface{}:
		m2 := make([]interface{}, len(iv))
		for j, vj := range iv {
			m2[j] = testVerifyVal(vj, arg)
		}
		v2 = m2
	case map[string]bool:
		switch arg {
		case testVerifyMapTypeSame:
			m2 := make(map[string]bool)
			for kj, kv := range iv {
				m2[kj] = kv
			}
			v2 = m2
		case testVerifyMapTypeStrIntf:
			m2 := make(map[string]interface{})
			for kj, kv := range iv {
				m2[kj] = kv
			}
			v2 = m2
		case testVerifyMapTypeIntfIntf:
			m2 := make(map[interface{}]interface{})
			for kj, kv := range iv {
				m2[kj] = kv
			}
			v2 = m2
		}
	case map[string]interface{}:
		switch arg {
		case testVerifyMapTypeSame:
			m2 := make(map[string]interface{})
			for kj, kv := range iv {
				m2[kj] = testVerifyVal(kv, arg)
			}
			v2 = m2
		case testVerifyMapTypeStrIntf:
			m2 := make(map[string]interface{})
			for kj, kv := range iv {
				m2[kj] = testVerifyVal(kv, arg)
			}
			v2 = m2
		case testVerifyMapTypeIntfIntf:
			m2 := make(map[interface{}]interface{})
			for kj, kv := range iv {
				m2[kj] = testVerifyVal(kv, arg)
			}
			v2 = m2
		}
	case map[interface{}]interface{}:
		m2 := make(map[interface{}]interface{})
		for kj, kv := range iv {
			m2[testVerifyVal(kj, arg)] = testVerifyVal(kv, arg)
		}
		v2 = m2
	default:
		v2 = v
	}
	return
}

func init() {
	primitives := []interface{}{
		int8(-8),
		int16(-1616),
		int32(-32323232),
		int64(-6464646464646464),
		uint8(192),
		uint16(1616),
		uint32(32323232),
		uint64(6464646464646464),
		byte(192),
		float32(-3232.0),
		float64(-6464646464.0),
		float32(3232.0),
		float64(6464646464.0),
		false,
		true,
		nil,
		timeToCompare,
		"someday",
		"",
		"bytestring",
	}
	mapsAndStrucs := []interface{}{
		map[string]bool{
			"true":  true,
			"false": false,
		},
		map[string]interface{}{
			"true":         "True",
			"false":        false,
			"uint16(1616)": uint16(1616),
		},
		//add a complex combo map in here. (map has list which has map)
		//note that after the first thing, everything else should be generic.
		map[string]interface{}{
			"list": []interface{}{
				int16(1616),
				int32(32323232),
				true,
				float32(-3232.0),
				map[string]interface{}{
					"TRUE":  true,
					"FALSE": false,
				},
				[]interface{}{true, false},
			},
			"int32":        int32(32323232),
			"bool":         true,
			"LONG STRING":  "123456789012345678901234567890123456789012345678901234567890",
			"SHORT STRING": "1234567890",
		},
		map[interface{}]interface{}{
			true:     "true",
			uint8(8): false,
			"false":  uint8(0),
		},
		newTestStruc(0, false),
	}

	table = []interface{}{}
	table = append(table, primitives...)    //0-19 are primitives
	table = append(table, primitives)       //20 is a list of primitives
	table = append(table, mapsAndStrucs...) //21-24 are maps. 25 is a *struct

	tableVerify = make([]interface{}, len(table))
	tableTestNilVerify = make([]interface{}, len(table))
	tablePythonVerify = make([]interface{}, len(table))

	lp := len(primitives)
	av := tableVerify
	for i, v := range table {
		if i == lp+3 {
			av[i] = skipVerifyVal
			continue
		}
		switch v.(type) {
		case []interface{}:
			av[i] = testVerifyVal(v, testVerifyMapTypeSame)
		case map[string]interface{}:
			av[i] = testVerifyVal(v, testVerifyMapTypeSame)
		case map[interface{}]interface{}:
			av[i] = testVerifyVal(v, testVerifyMapTypeSame)
		default:
			av[i] = v
		}
	}

	av = tableTestNilVerify
	for i, v := range table {
		if i > lp+3 {
			av[i] = skipVerifyVal
			continue
		}
		av[i] = testVerifyVal(v, testVerifyMapTypeStrIntf)
	}

	av = tablePythonVerify
	for i, v := range table {
		if i > lp+3 {
			av[i] = skipVerifyVal
			continue
		}
		av[i] = testVerifyVal(v, testVerifyMapTypeStrIntf)
		//msgpack python encodes large positive numbers as unsigned ints
		switch i {
		case 16:
			av[i] = uint64(1328148122000002)
		case 20:
			//msgpack python doesn't understand time. So use number val.
			m2 := av[i].([]interface{})
			m3 := make([]interface{}, len(m2))
			copy(m3, m2)
			m3[16] = uint64(1328148122000002)
			av[i] = m3
		case 23:
			m2 := make(map[string]interface{})
			for k2, v2 := range av[i].(map[string]interface{}) {
				m2[k2] = v2
			}
			m2["int32"] = uint64(32323232)
			m3 := m2["list"].([]interface{})
			m3[0], m3[1], m3[3] = uint64(1616), uint64(32323232), float64(-3232.0)
			av[i] = m2
		}
	}

	tablePythonVerify = tablePythonVerify[:24]
}

func testUnmarshal(v interface{}, data []byte, h Handle) error {
	if testUseIoEncDec {
		return NewDecoder(bytes.NewBuffer(data), h).Decode(v)
	}
	return NewDecoderBytes(data, h).Decode(v)
}

func testMarshal(v interface{}, h Handle) (bs []byte, err error) {
	if testUseIoEncDec {
		var buf bytes.Buffer
		err = NewEncoder(&buf, h).Encode(v)
		bs = buf.Bytes()
		return
	}
	err = NewEncoderBytes(&bs, h).Encode(v)
	return
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
		Ui64slice: []uint64{7, 8, 9},
		Ui8slice:  []uint8{10, 11, 12},
		Bslice:    []bool{true, false, true, false},
		Byslice:   []byte{13, 14, 15},

		Islice: []interface{}{"true", true, "no", false, uint64(88), float64(0.4)},

		Ms: map[string]interface{}{
			"true":     "true",
			"int64(9)": false,
		},
		Msi64: map[string]int64{
			"one": 1,
			"two": 2,
		},
		T: timeToCompare,
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

// doTestCodecTableOne allows us test for different variations based on arguments passed.
func doTestCodecTableOne(t *testing.T, testNil bool, h Handle,
	vs []interface{}, vsVerify []interface{}) {
	//if testNil, then just test for when a pointer to a nil interface{} is passed. It should work.
	//Current setup allows us test (at least manually) the nil interface or typed interface.
	logT(t, "================ TestNil: %v ================\n", testNil)
	for i, v0 := range vs {
		logT(t, "..............................................")
		logT(t, "         Testing: #%d:, %T, %#v\n", i, v0, v0)
		b0, err := testMarshal(v0, h)
		if err != nil {
			logT(t, err.Error())
			failT(t)
			continue
		}
		logT(t, "         Encoded bytes: len: %v, %v\n", len(b0), b0)

		var v1 interface{}

		if testNil {
			err = testUnmarshal(&v1, b0, h)
		} else {
			if v0 != nil {
				v0rt := reflect.TypeOf(v0) // ptr
				rv1 := reflect.New(v0rt)
				err = testUnmarshal(rv1.Interface(), b0, h)
				v1 = rv1.Elem().Interface()
				// v1 = reflect.Indirect(reflect.ValueOf(v1)).Interface()
			}
		}

		logT(t, "         v1 returned: %T, %#v", v1, v1)
		// if v1 != nil {
		//	logT(t, "         v1 returned: %T, %#v", v1, v1)
		//	//we always indirect, because ptr to typed value may be passed (if not testNil)
		//	v1 = reflect.Indirect(reflect.ValueOf(v1)).Interface()
		// }
		if err != nil {
			logT(t, "-------- Error: %v. Partial return: %v", err, v1)
			failT(t)
			continue
		}
		v0check := vsVerify[i]
		if v0check == skipVerifyVal {
			logT(t, "        Nil Check skipped: Decoded: %T, %#v\n", v1, v1)
			continue
		}

		if err = deepEqual(v0check, v1); err == nil {
			logT(t, "++++++++ Before and After marshal matched\n")
		} else {
			logT(t, "-------- Before and After marshal do not match: Error: %v"+
				" ====> AGAINST: (%T) %#v, DECODED: (%T) %#v\n", err, v0check, v0check, v1, v1)
			failT(t)
		}
	}
}

func testCodecTableOne(t *testing.T, h Handle) {
	// func TestMsgpackAllExperimental(t *testing.T) {
	// dopts := testDecOpts(nil, nil, false, true, true),
	var oldWriteExt bool
	switch v := h.(type) {
	case *MsgpackHandle:
		oldWriteExt = v.WriteExt
		v.WriteExt = true
	}
	doTestCodecTableOne(t, false, h, table, tableVerify)
	//if true { panic("") }
	switch v := h.(type) {
	case *MsgpackHandle:
		v.WriteExt = oldWriteExt
	}
	// func TestMsgpackAll(t *testing.T) {

	doTestCodecTableOne(t, false, h, table[:20], tableVerify[:20])
	doTestCodecTableOne(t, false, h, table[21:], tableVerify[21:])

	// func TestMsgpackNilStringMap(t *testing.T) {
	var oldMapType reflect.Type
	switch v := h.(type) {
	case *MsgpackHandle:
		oldMapType = v.MapType
		v.MapType = mapStringIntfTyp
	case *BincHandle:
		oldMapType = v.MapType
		v.MapType = mapStringIntfTyp
	}
	//skip #16 (time.Time), and #20 ([]interface{} containing time.Time)
	doTestCodecTableOne(t, true, h, table[:16], tableTestNilVerify[:16])
	doTestCodecTableOne(t, true, h, table[17:20], tableTestNilVerify[17:20])
	doTestCodecTableOne(t, true, h, table[21:24], tableTestNilVerify[21:24])

	switch v := h.(type) {
	case *MsgpackHandle:
		v.MapType = oldMapType
	case *BincHandle:
		v.MapType = oldMapType
	}

	// func TestMsgpackNilIntf(t *testing.T) {
	doTestCodecTableOne(t, true, h, table[24:], tableTestNilVerify[24:])

	doTestCodecTableOne(t, true, h, table[17:18], tableTestNilVerify[17:18])
}

func testCodecMiscOne(t *testing.T, h Handle) {
	b, err := testMarshal(32, h)
	// Cannot do this nil one, because faster type assertion decoding will panic
	// var i *int32
	// if err = testUnmarshal(b, i, nil); err == nil {
	// 	logT(t, "------- Expecting error because we cannot unmarshal to int32 nil ptr")
	// 	t.FailNow()
	// }
	var i2 int32 = 0
	if err = testUnmarshal(&i2, b, h); err != nil {
		logT(t, "------- Cannot unmarshal to int32 ptr. Error: %v", err)
		t.FailNow()
	}
	if i2 != int32(32) {
		logT(t, "------- didn't unmarshal to 32: Received: %d", i2)
		t.FailNow()
	}

	// func TestMsgpackDecodePtr(t *testing.T) {
	ts := newTestStruc(0, false)
	b, err = testMarshal(ts, h)
	if err != nil {
		logT(t, "------- Cannot Marshal pointer to struct. Error: %v", err)
		t.FailNow()
	} else if len(b) < 40 {
		logT(t, "------- Size must be > 40. Size: %d", len(b))
		t.FailNow()
	}
	logT(t, "------- b: %v", b)
	ts2 := new(TestStruc)
	err = testUnmarshal(ts2, b, h)
	if err != nil {
		logT(t, "------- Cannot Unmarshal pointer to struct. Error: %v", err)
		t.FailNow()
	} else if ts2.I64 != math.MaxInt64*2/3 {
		logT(t, "------- Unmarshal wrong. Expect I64 = 64. Got: %v", ts2.I64)
		t.FailNow()
	}

	// func TestMsgpackIntfDecode(t *testing.T) {
	m := map[string]int{"A": 2, "B": 3}
	p := []interface{}{m}
	bs, err := testMarshal(p, h)
	if err != nil {
		logT(t, "Error marshalling p: %v, Err: %v", p, err)
		t.FailNow()
	}
	m2 := map[string]int{}
	p2 := []interface{}{m2}
	err = testUnmarshal(&p2, bs, h)
	if err != nil {
		logT(t, "Error unmarshalling into &p2: %v, Err: %v", p2, err)
		t.FailNow()
	}

	if m2["A"] != 2 || m2["B"] != 3 {
		logT(t, "m2 not as expected: expecting: %v, got: %v", m, m2)
		t.FailNow()
	}
	// log("m: %v, m2: %v, p: %v, p2: %v", m, m2, p, p2)
	if err = deepEqual(p, p2); err == nil {
		logT(t, "p and p2 match")
	} else {
		logT(t, "Not Equal: %v. p: %v, p2: %v", err, p, p2)
		t.FailNow()
	}
	if err = deepEqual(m, m2); err == nil {
		logT(t, "m and m2 match")
	} else {
		logT(t, "Not Equal: %v. m: %v, m2: %v", err, m, m2)
		t.FailNow()
	}

	// func TestMsgpackDecodeStructSubset(t *testing.T) {
	// test that we can decode a subset of the stream
	mm := map[string]interface{}{"A": 5, "B": 99, "C": 333}
	bs, err = testMarshal(mm, h)
	if err != nil {
		logT(t, "Error marshalling m: %v, Err: %v", mm, err)
		t.FailNow()
	}
	type ttt struct {
		A uint8
		C int32
	}
	var t2 ttt
	err = testUnmarshal(&t2, bs, h)
	if err != nil {
		logT(t, "Error unmarshalling into &t2: %v, Err: %v", t2, err)
		t.FailNow()
	}
	t3 := ttt{5, 333}
	if err = deepEqual(t2, t3); err != nil {
		logT(t, "Not Equal: %v. t2: %v, t3: %v", err, t2, t3)
		t.FailNow()
	}
}

func doTestRpcOne(t *testing.T, rr Rpc, h Handle, callClose, doRequest, doExit bool) {
	srv := rpc.NewServer()
	srv.Register(testRpcInt)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	// log("listener: %v", ln.Addr())
	checkErrT(t, err)
	defer ln.Close()

	// var opts *DecoderOptions
	// opts := testDecOpts
	// opts.MapType = mapStringIntfTyp
	// opts.RawToString = false
	serverExitChan := make(chan bool, 1)
	serverFn := func() {
		for {
			conn1, err1 := ln.Accept()
			if err1 != nil {
				continue
			}
			bs := make([]byte, 1)
			n1, err1 := conn1.Read(bs)
			if n1 != 1 || err1 != nil {
				conn1.Close()
				continue
			}
			var sc rpc.ServerCodec
			switch bs[0] {
			case 'R':
				sc = rr.ServerCodec(conn1, h)
			case 'X':
				serverExitChan <- true
				// <- serverExitChan
				conn1.Close()
				return // exit serverFn goroutine
			}
			if sc == nil {
				conn1.Close()
				continue
			}
			srv.ServeCodec(sc)
			// for {
			// 	if err1 = srv.ServeRequest(sc); err1 != nil {
			// 		break
			// 	}
			// }
			// if callClose {
			// 	sc.Close()
			// }
		}
	}

	clientFn := func(cc rpc.ClientCodec) {
		cl := rpc.NewClientWithCodec(cc)
		if callClose {
			defer cl.Close()
		}
		var up, sq, mult int
		// log("Calling client")
		checkErrT(t, cl.Call("TestRpcInt.Update", 5, &up))
		// log("Called TestRpcInt.Update")
		checkEqualT(t, testRpcInt.i, 5)
		checkEqualT(t, up, 5)
		checkErrT(t, cl.Call("TestRpcInt.Square", 1, &sq))
		checkEqualT(t, sq, 25)
		checkErrT(t, cl.Call("TestRpcInt.Mult", 20, &mult))
		checkEqualT(t, mult, 100)
	}

	connFn := func(req byte) (bs net.Conn) {
		// log("calling f1")
		bs, err2 := net.Dial(ln.Addr().Network(), ln.Addr().String())
		// log("f1. bs: %v, err2: %v", bs, err2)
		checkErrT(t, err2)
		n1, err2 := bs.Write([]byte{req})
		checkErrT(t, err2)
		checkEqualT(t, n1, 1)
		return
	}

	go serverFn()
	if doRequest {
		bs := connFn('R')
		cc := rr.ClientCodec(bs, h)
		clientFn(cc)
	}
	if doExit {
		bs := connFn('X')
		<-serverExitChan
		bs.Close()
		// serverExitChan <- true
	}
}

// Comprehensive testing that generates data encoded from python msgpack,
// and validates that our code can read and write it out accordingly.
// We keep this unexported here, and put actual test in ext_dep_test.go.
// This way, it can be excluded by excluding file completely.
func doTestMsgpackPythonGenStreams(t *testing.T) {
	logT(t, "TestPythonGenStreams")
	tmpdir, err := ioutil.TempDir("", "golang-msgpack-test")
	if err != nil {
		logT(t, "-------- Unable to create temp directory\n")
		t.FailNow()
	}
	defer os.RemoveAll(tmpdir)
	logT(t, "tmpdir: %v", tmpdir)
	cmd := exec.Command("python", "msgpack_test.py", "testdata", tmpdir)
	//cmd.Stdin = strings.NewReader("some input")
	//cmd.Stdout = &out
	var cmdout []byte
	if cmdout, err = cmd.CombinedOutput(); err != nil {
		logT(t, "-------- Error running python build.py. Err: %v", err)
		logT(t, "         %v", string(cmdout))
		t.FailNow()
	}

	oldMapType := testMsgpackH.MapType
	for i, v := range tablePythonVerify {
		testMsgpackH.MapType = oldMapType
		//load up the golden file based on number
		//decode it
		//compare to in-mem object
		//encode it again
		//compare to output stream
		logT(t, "..............................................")
		logT(t, "         Testing: #%d: %T, %#v\n", i, v, v)
		var bss []byte
		bss, err = ioutil.ReadFile(filepath.Join(tmpdir, strconv.Itoa(i)+".golden"))
		if err != nil {
			logT(t, "-------- Error reading golden file: %d. Err: %v", i, err)
			failT(t)
			continue
		}
		testMsgpackH.MapType = mapStringIntfTyp

		var v1 interface{}
		if err = testUnmarshal(&v1, bss, testMsgpackH); err != nil {
			logT(t, "-------- Error decoding stream: %d: Err: %v", i, err)
			failT(t)
			continue
		}
		if v == skipVerifyVal {
			continue
		}
		//no need to indirect, because we pass a nil ptr, so we already have the value
		//if v1 != nil { v1 = reflect.Indirect(reflect.ValueOf(v1)).Interface() }
		if err = deepEqual(v, v1); err == nil {
			logT(t, "++++++++ Objects match")
		} else {
			logT(t, "-------- Objects do not match: %v. Source: %T. Decoded: %T", err, v, v1)
			logT(t, "--------   AGAINST: %#v", v)
			logT(t, "--------   DECODED: %#v <====> %#v", v1, reflect.Indirect(reflect.ValueOf(v1)).Interface())
			failT(t)
		}
		bsb, err := testMarshal(v1, testMsgpackH)
		if err != nil {
			logT(t, "Error encoding to stream: %d: Err: %v", i, err)
			failT(t)
			continue
		}
		if err = deepEqual(bsb, bss); err == nil {
			logT(t, "++++++++ Bytes match")
		} else {
			logT(t, "???????? Bytes do not match. %v.", err)
			xs := "--------"
			if reflect.ValueOf(v).Kind() == reflect.Map {
				xs = "        "
				logT(t, "%s It's a map. Ok that they don't match (dependent on ordering).", xs)
			} else {
				logT(t, "%s It's not a map. They should match.", xs)
				failT(t)
			}
			logT(t, "%s   FROM_FILE: %4d] %v", xs, len(bss), bss)
			logT(t, "%s     ENCODED: %4d] %v", xs, len(bsb), bsb)
		}
	}
	testMsgpackH.MapType = oldMapType

}

func TestMsgpackCodecsTable(t *testing.T) {
	testCodecTableOne(t, testMsgpackH)
}

func TestMsgpackCodecsMisc(t *testing.T) {
	testCodecMiscOne(t, testMsgpackH)
}

func TestBincCodecsTable(t *testing.T) {
	testCodecTableOne(t, testBincH)
}

func TestBincCodecsMisc(t *testing.T) {
	testCodecMiscOne(t, testBincH)
}

func TestMsgpackRpcGo(t *testing.T) {
	doTestRpcOne(t, GoRpc, testMsgpackH, true, true, true)
}
func TestMsgpackRpcSpec(t *testing.T) {
	doTestRpcOne(t, MsgpackSpecRpc, testMsgpackH, true, true, true)
}

func TestBincRpcGo(t *testing.T) {
	doTestRpcOne(t, GoRpc, testBincH, true, true, true)
}
