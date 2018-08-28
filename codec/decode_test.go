// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"errors"
	"reflect"
	"testing"
)

func TestDecoderWrapErr(t *testing.T) {
	var h JsonHandle
	d := NewDecoderBytes([]byte{}, &h)

	err := errors.New("test error")
	wrappedErr := d.wrapErr(err)

	expectedErr := DecodeError{Name: h.Name(), Pos: 0, Err: err}
	if !reflect.DeepEqual(expectedErr, wrappedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, wrappedErr)
	}
}
