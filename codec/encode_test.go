// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"errors"
	"reflect"
	"testing"
)

func TestWrapErr(t *testing.T) {
	var h JsonHandle
	e := NewEncoderBytes(nil, &h)

	err := errors.New("test error")
	wrappedErr := e.wrapErr(err)

	expectedErr := EncodeError{Name: h.Name(), Err: err}
	if !reflect.DeepEqual(expectedErr, wrappedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, wrappedErr)
	}
}
