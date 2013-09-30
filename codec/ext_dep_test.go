//+build ignore

// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// This file includes benchmarks which have dependencies on 3rdparty
// packages (bson and vmihailenco/msgpack) which must be installed locally.
//
// To run the benchmarks including these 3rdparty packages, first
//   - Uncomment first line in this file (put // // in front of it)
//   - Get those packages:
//       go get github.com/vmihailenco/msgpack
//       go get labix.org/v2/mgo/bson
//   - Run:
//       go test -bi -bench=.

import (
	vmsgpack "github.com/vmihailenco/msgpack"
	"labix.org/v2/mgo/bson"
	"testing"
)

func init() {
	benchCheckers = append(benchCheckers,
		benchChecker{"v-msgpack", fnVMsgpackEncodeFn, fnVMsgpackDecodeFn},
		benchChecker{"bson", fnBsonEncodeFn, fnBsonDecodeFn},
	)
}

func fnVMsgpackEncodeFn(ts *TestStruc) ([]byte, error) {
	return vmsgpack.Marshal(ts)
}

func fnVMsgpackDecodeFn(buf []byte, ts *TestStruc) error {
	return vmsgpack.Unmarshal(buf, ts)
}

func fnBsonEncodeFn(ts *TestStruc) ([]byte, error) {
	return bson.Marshal(ts)
}

func fnBsonDecodeFn(buf []byte, ts *TestStruc) error {
	return bson.Unmarshal(buf, ts)
}

func Benchmark__Bson_____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "bson", fnBsonEncodeFn)
}

func Benchmark__Bson_____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "bson", fnBsonEncodeFn, fnBsonDecodeFn)
}

func Benchmark__VMsgpack_Encode(b *testing.B) {
	fnBenchmarkEncode(b, "v-msgpack", fnVMsgpackEncodeFn)
}

func Benchmark__VMsgpack_Decode(b *testing.B) {
	fnBenchmarkDecode(b, "v-msgpack", fnVMsgpackEncodeFn, fnVMsgpackDecodeFn)
}

func TestMsgpackPythonGenStreams(t *testing.T) {
	doTestMsgpackPythonGenStreams(t)
}

func TestMsgpackRpcSpecGoClientToPythonSvc(t *testing.T) {
	doTestMsgpackRpcSpecGoClientToPythonSvc(t)
}

func TestMsgpackRpcSpecPythonClientToGoSvc(t *testing.T) {
	doTestMsgpackRpcSpecPythonClientToGoSvc(t)
}
