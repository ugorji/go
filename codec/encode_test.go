// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"errors"
	"reflect"
	"testing"
)

func TestEncodeErrorCause(t *testing.T) {
	err := errors.New("test error")
	wrappedErr := EncodeError{Err: err}

	cause := wrappedErr.Cause()
	if !reflect.DeepEqual(err, cause) {
		t.Fatalf("expected %v, got %v", err, cause)
	}
}

func TestEncoderWrapErr(t *testing.T) {
	var h JsonHandle
	e := NewEncoderBytes(nil, &h)

	err := errors.New("test error")
	wrappedErr := e.wrapErr(err)

	expectedErr := EncodeError{Name: h.Name(), Err: err}
	if !reflect.DeepEqual(expectedErr, wrappedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, wrappedErr)
	}
}
