package codec

import (
	"reflect"
	"testing"
)

func TestEmptySliceEncoding(t *testing.T) {
	bincHandle := new(BincHandle)
	var b []byte
	t1 := []string{}
	if e := NewEncoderBytes(&b, bincHandle).Encode(t1); e != nil {
		panic(e)
	}
	var t2 []string
	if e := NewDecoderBytes(b, bincHandle).Decode(&t2); e != nil {
		panic(e)
	}
	if !reflect.DeepEqual(t1, t2) {
		t.Errorf("%#v != %#v", t1, t2)
	}
}
