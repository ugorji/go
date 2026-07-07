// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"testing"
)

// TestGH437: a []struct-with-a-map used as a map value must round-trip.
//
// Regression: isCanTransient marked any top-level slice transient without
// checking its element, so a []T whose element holds a map/pointer field reused
// the decoder transient scratch and was corrupted on decode (nil-pointer deref
// / SIGSEGV in typedmemclr). Re-fixes the regression of #367.
func TestGH437(t *testing.T) {
	type P struct {
		X string
		N int64
	}
	type Item struct {
		S string
		M map[int64][]P
	}

	in := map[string][]Item{"a": {{S: "s", M: map[int64][]P{7: {{X: "x", N: 1}}}}}}

	var h CborHandle
	var b []byte
	NewEncoderBytes(&b, &h).MustEncode(in)

	var out map[string][]Item
	NewDecoderBytes(b, &h).MustDecode(&out)

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round-trip mismatch:\n in: %#v\nout: %#v", in, out)
	}
}
