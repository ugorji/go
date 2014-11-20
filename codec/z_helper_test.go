// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// All non-std package dependencies related to testing live in this file,
// so porting to different environment is easy (just update functions).
//
// Also, this file is called z_helper_test, to give a "hint" to compiler
// that its init() function should be called last. (not guaranteed by spec)

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"testing"
)

var (
	testLogToT    = true
	failNowOnFail = true
)

func init() {
	testBincHSym.AsSymbols = AsSymbolAll
	testBincHNoSym.AsSymbols = AsSymbolNone
	testInitFlags()
	flag.Parse()
	testInit()
}

var (
	testMsgpackH   = &MsgpackHandle{}
	testBincH      = &BincHandle{}
	testBincHNoSym = &BincHandle{}
	testBincHSym   = &BincHandle{}
	testSimpleH    = &SimpleHandle{}
	testCborH      = &CborHandle{}
	testJsonH      = &JsonHandle{}
)

func fnCodecEncode(ts interface{}, h Handle) (bs []byte, err error) {
	err = NewEncoderBytes(&bs, h).Encode(ts)
	return
}

func fnCodecDecode(buf []byte, ts interface{}, h Handle) error {
	return NewDecoderBytes(buf, h).Decode(ts)
}

func checkErrT(t *testing.T, err error) {
	if err != nil {
		logT(t, err.Error())
		failT(t)
	}
}

func checkEqualT(t *testing.T, v1 interface{}, v2 interface{}, desc string) (err error) {
	if err = deepEqual(v1, v2); err != nil {
		logT(t, "Not Equal: %s: %v. v1: %v, v2: %v", desc, err, v1, v2)
		failT(t)
	}
	return
}

func logT(x interface{}, format string, args ...interface{}) {
	if t, ok := x.(*testing.T); ok && t != nil && testLogToT {
		t.Logf(format, args...)
	} else if b, ok := x.(*testing.B); ok && b != nil && testLogToT {
		b.Logf(format, args...)
	} else {
		if len(format) == 0 || format[len(format)-1] != '\n' {
			format = format + "\n"
		}
		fmt.Printf(format, args...)
	}
}

func failT(t *testing.T) {
	if failNowOnFail {
		t.FailNow()
	} else {
		t.Fail()
	}
}

func deepEqual(v1, v2 interface{}) (err error) {
	if !reflect.DeepEqual(v1, v2) {
		err = errors.New("Not Match")
	}
	return
}
