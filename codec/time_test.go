package codec

import (
	"testing"
	"time"
)

func TestTimeEncoding(t *testing.T) {
	bincHandle := new(BincHandle)
	var b []byte
	t1 := time.Now()
	if e := NewEncoderBytes(&b, bincHandle).Encode(t1); e != nil {
		panic(e)
	}
	var t2 time.Time
	if e := NewDecoderBytes(b, bincHandle).Decode(&t2); e != nil {
		panic(e)
	}
	if !t1.Equal(t2) {
		t.Errorf("%v != %v", t1, t2)
	}
	t1 = time.Time{}
	if e := NewEncoderBytes(&b, bincHandle).Encode(t1); e != nil {
		panic(e)
	}
	if e := NewDecoderBytes(b, bincHandle).Decode(&t2); e != nil {
		panic(e)
	}
	if !t1.Equal(t2) {
		t.Errorf("%v != %v", t1, t2)
	}

}
