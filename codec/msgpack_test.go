// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/rpc"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

// To test MsgpackSpecRpc, we test 3 scenarios:
//    - Go Client to Go RPC Service (contained within TestMsgpackRpcSpec)
//    - Go client to Python RPC Service (contained within doTestMsgpackRpcSpecGoClientToPythonSvc)
//    - Python Client to Go RPC Service (contained within doTestMsgpackRpcSpecPythonClientToGoSvc)
//
// This allows us test the different calling conventions
//    - Go Service requires only one argument
//    - Python Service allows multiple arguments

func doTestMsgpackRpcSpecGoClientToPythonSvc(t *testing.T, h Handle) {
	if testv.SkipRPCTests {
		t.Skip(testSkipRPCTestsMsg)
	}
	defer testSetup(t, &h)()

	// openPorts are between 6700 and 6800
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	openPort := strconv.FormatInt(6700+r.Int63n(99), 10)
	// openPort := "6792"
	cmd := exec.Command("python", "test.py", "rpc-server", openPort, "4")
	testCheckErr(t, cmd.Start())
	bs, err2 := net.Dial("tcp", ":"+openPort)
	maxSleepTime := 500 * time.Millisecond
	iterSleepTime := 5 * time.Millisecond
	for i := 0; i < int(maxSleepTime/iterSleepTime) && err2 != nil; i++ {
		time.Sleep(iterSleepTime) // time for python rpc server to start
		bs, err2 = net.Dial("tcp", ":"+openPort)
	}
	testCheckErr(t, err2)
	cc := MsgpackSpecRpc.ClientCodec(testReadWriteCloser(bs), h)
	cl := rpc.NewClientWithCodec(cc)
	defer cl.Close()
	var rstr string
	testCheckErr(t, cl.Call("EchoStruct", TestRpcABC{"Aa", "Bb", "Cc"}, &rstr))
	//testCheckEqual(t, rstr, "{'A': 'Aa', 'B': 'Bb', 'C': 'Cc'}")
	var mArgs MsgpackSpecRpcMultiArgs = []interface{}{"A1", "B2", "C3"}
	testCheckErr(t, cl.Call("Echo123", mArgs, &rstr))
	testCheckEqual(t, rstr, "1:A1 2:B2 3:C3", "rstr=")
	cmd.Process.Kill()
}

func doTestMsgpackRpcSpecPythonClientToGoSvc(t *testing.T, h Handle) {
	if testv.SkipRPCTests {
		t.Skip(testSkipRPCTestsMsg)
	}
	defer testSetup(t, &h)()
	// seems 10ms is not enough for test run; set high to 1 second, and client will close when done
	exitSleepDur := 1000 * time.Millisecond
	port := doTestCodecRpcOne(t, MsgpackSpecRpc, h, false, exitSleepDur)
	cmd := exec.Command("python", "test.py", "rpc-client-go-service", strconv.Itoa(port))
	var cmdout []byte
	var err error
	if cmdout, err = cmd.CombinedOutput(); err != nil {
		t.Logf("-------- Error running test.py rpc-client-go-service. Err: %v", err)
		t.Logf("         %v", string(cmdout))
		t.FailNow()
	}
	testCheckEqual(t, string(cmdout),
		fmt.Sprintf("%#v\n%#v\n", []string{"A1", "B2", "C3"}, TestRpcABC{"Aa", "Bb", "Cc"}), "cmdout=")
}

// -----------

func doTestMsgpackDecodeMapAndExtSizeMismatch(t *testing.T, h Handle) {
	if !testRecoverPanicToErr {
		t.Skip(testSkipIfNotRecoverPanicToErrMsg)
	}
	defer testSetup(t, &h)()
	fn := func(t *testing.T, b []byte, v interface{}) {
		err := NewDecoderBytes(b, h).Decode(v)
		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("expected EOF or ErrUnexpectedEOF, got %v", err)
		}
	}

	// a map claiming to have 0x10eeeeee KV pairs, but only has 1.
	var b = []byte{0xdf, 0x10, 0xee, 0xee, 0xee, 0x1, 0xa1, 0x1}
	var m1 map[int]string
	var m2 map[int][]byte
	fn(t, b, &m1)
	fn(t, b, &m2)

	// an extension claiming to have 0x7fffffff bytes, but only has 1.
	b = []byte{0xc9, 0x7f, 0xff, 0xff, 0xff, 0xda, 0x1}
	var a interface{}
	fn(t, b, &a)

	// b = []byte{0x00}
	// var s testSelferRecur
	// fn(t, b, &s)
}

func TestMsgpackCodecsTable(t *testing.T) {
	doTestCodecTableOne(t, testMsgpackH)
}

func TestMsgpackCodecsMisc(t *testing.T) {
	doTestCodecMiscOne(t, testMsgpackH)
}

func TestMsgpackCodecsEmbeddedPointer(t *testing.T) {
	doTestCodecEmbeddedPointer(t, testMsgpackH)
}

func TestMsgpackStdEncIntf(t *testing.T) {
	doTestStdEncIntf(t, testMsgpackH)
}

func TestMsgpackRaw(t *testing.T) {
	doTestRawValue(t, testMsgpackH)
}

func TestMsgpackRpcSpec(t *testing.T) {
	doTestCodecRpcOne(t, MsgpackSpecRpc, testMsgpackH, true, 0)
}

func TestMsgpackRpcGo(t *testing.T) {
	doTestCodecRpcOne(t, GoRpc, testMsgpackH, true, 0)
}

func TestMsgpackMapEncodeForCanonical(t *testing.T) {
	doTestMapEncodeForCanonical(t, testMsgpackH)
}

func TestMsgpackSwallowAndZero(t *testing.T) {
	doTestSwallowAndZero(t, testMsgpackH)
}

func TestMsgpackRawExt(t *testing.T) {
	doTestRawExt(t, testMsgpackH)
}

func TestMsgpackMapStructKey(t *testing.T) {
	doTestMapStructKey(t, testMsgpackH)
}

func TestMsgpackDecodeNilMapValue(t *testing.T) {
	doTestDecodeNilMapValue(t, testMsgpackH)
}

func TestMsgpackEmbeddedFieldPrecedence(t *testing.T) {
	doTestEmbeddedFieldPrecedence(t, testMsgpackH)
}

func TestMsgpackLargeContainerLen(t *testing.T) {
	doTestLargeContainerLen(t, testMsgpackH)
}

func TestMsgpackTime(t *testing.T) {
	doTestTime(t, testMsgpackH)
}

func TestMsgpackUintToInt(t *testing.T) {
	doTestUintToInt(t, testMsgpackH)
}

func TestMsgpackDifferentMapOrSliceType(t *testing.T) {
	doTestDifferentMapOrSliceType(t, testMsgpackH)
}

func TestMsgpackScalars(t *testing.T) {
	mh := testMsgpackH
	defer func(b bool) { mh.RawToString = b }(mh.RawToString)
	mh.RawToString = true
	doTestScalars(t, testMsgpackH)
}

func TestMsgpackOmitempty(t *testing.T) {
	doTestOmitempty(t, testMsgpackH)
}

func TestMsgpackIntfMapping(t *testing.T) {
	doTestIntfMapping(t, testMsgpackH)
}

func TestMsgpackMissingFields(t *testing.T) {
	doTestMissingFields(t, testMsgpackH)
}

func TestMsgpackMaxDepth(t *testing.T) {
	doTestMaxDepth(t, testMsgpackH)
}

func TestMsgpackSelfExt(t *testing.T) {
	doTestSelfExt(t, testMsgpackH)
}

func TestMsgpackBytesEncodedAsArray(t *testing.T) {
	doTestBytesEncodedAsArray(t, testMsgpackH)
}

func TestMsgpackStrucEncDec(t *testing.T) {
	doTestStrucEncDec(t, testMsgpackH)
}

func TestMsgpackRawToStringToRawEtc(t *testing.T) {
	doTestRawToStringToRawEtc(t, testMsgpackH)
}

func TestMsgpackStructKeyType(t *testing.T) {
	doTestStructKeyType(t, testMsgpackH)
}

func TestMsgpackPreferArrayOverSlice(t *testing.T) {
	doTestPreferArrayOverSlice(t, testMsgpackH)
}

func TestMsgpackZeroCopyBytes(t *testing.T) {
	doTestZeroCopyBytes(t, testMsgpackH)
}

func TestMsgpackNextValueBytes(t *testing.T) {
	doTestNextValueBytes(t, testMsgpackH)
}

func TestMsgpackNumbers(t *testing.T) {
	doTestNumbers(t, testMsgpackH)
}

func TestMsgpackDesc(t *testing.T) {
	m := make(map[byte]string)
	for k, v := range mpdescNames {
		m[k] = v
	}
	m[mpPosFixNumMin] = "int"
	m[mpFixStrMin] = "string|bytes"
	m[mpFixArrayMin] = "array"
	m[mpFixMapMin] = "map"
	m[mpFixExt1] = "ext"

	doTestDesc(t, testMsgpackH, m)
}

func TestMsgpackStructFieldInfoToArray(t *testing.T) {
	doTestStructFieldInfoToArray(t, testMsgpackH)
}

func TestMsgpackDecodeMapAndExtSizeMismatch(t *testing.T) {
	doTestMsgpackDecodeMapAndExtSizeMismatch(t, testMsgpackH)
}
