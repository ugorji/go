package codec

import (
	"testing"
	"time"
)

func TestMsgpackLegacyTimeNotBuiltinDecode(t *testing.T) {
	// Encoded with v1.2.14 TimeNotBuiltin=true for 2026-01-01T00:00:00Z
	data := []byte{
		0x81, 0xa1, 0x54, 0xaf, 0x01, 0x00, 0x00, 0x00, 0x0e, 0xe0,
		0xe7, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff,
	}
	h := new(MsgpackHandle)
	h.TimeNotBuiltin = true
	var v struct{ T time.Time }
	if err := NewDecoderBytes(data, h).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if !v.T.Equal(want) {
		t.Fatalf("got %v want %v", v.T, want)
	}
}
