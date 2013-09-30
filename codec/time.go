// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

import (
	"time"
)

var (
	timeDigits = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

// encodeTime encodes a time.Time as a []byte, including
// information on the instant in time and UTC offset.
func encodeTime(t time.Time) []byte {
	//t := rv.Interface().(time.Time)
	tsecs, tnsecs := t.Unix(), t.Nanosecond()
	var (
		bd   byte
		btmp [8]byte
		bs   [16]byte
		i    int = 1
	)
	l := t.Location()
	if l == time.UTC {
		l = nil
	}
	if tsecs != 0 {
		bd = bd | 0x80
		bigen.PutUint64(btmp[:], uint64(tsecs))
		f := pruneSignExt(btmp[:])
		bd = bd | (byte(7-f) << 2)
		copy(bs[i:], btmp[f:])
		i = i + (8 - f)
	}
	if tnsecs != 0 {
		bd = bd | 0x40
		bigen.PutUint32(btmp[:4], uint32(tnsecs))
		f := pruneSignExt(btmp[:4])
		bd = bd | byte(3-f)
		copy(bs[i:], btmp[f:4])
		i = i + (4 - f)
	}
	if l != nil {
		bd = bd | 0x20
		// Note that Go Libs do not give access to dst flag.
		_, zoneOffset := t.Zone()
		//zoneName, zoneOffset := t.Zone()
		zoneOffset /= 60
		z := uint16(zoneOffset)
		bigen.PutUint16(btmp[:2], z)
		// clear dst flags
		bs[i] = btmp[0] & 0x3f
		bs[i+1] = btmp[1]
		i = i + 2
	}
	bs[0] = bd
	return bs[0:i]
}

// decodeTime decodes a []byte into a time.Time,
// and sets into passed reflectValue.
func decodeTime(bs []byte) (tt time.Time, err error) {
	bd := bs[0]
	var (
		tsec  int64
		tnsec uint32
		tz    uint16
		i     byte = 1
		i2    byte
		n     byte
	)
	if bd&(1<<7) != 0 {
		var btmp [8]byte
		n = ((bd >> 2) & 0x7) + 1
		i2 = i + n
		copy(btmp[8-n:], bs[i:i2])
		//if first bit of bs[i] is set, then fill btmp[0..8-n] with 0xff (ie sign extend it)
		if bs[i] & (1 << 7) != 0 {
			copy(btmp[0:8-n], bsAll0xff)
			//for j,k := byte(0), 8-n; j < k; j++ {	btmp[j] = 0xff }
		}
		i = i2
		tsec = int64(bigen.Uint64(btmp[:]))
	}
	if bd&(1<<6) != 0 {
		var btmp [4]byte
		n = (bd & 0x3) + 1
		i2 = i + n
		copy(btmp[4-n:], bs[i:i2])
		i = i2
		tnsec = bigen.Uint32(btmp[:])
	}
	if bd&(1<<5) == 0 {
		tt = time.Unix(tsec, int64(tnsec)).UTC()
		return
	}
	// In stdlib time.Parse, when a date is parsed without a zone name, it uses "" as zone name.
	// However, we need name here, so it can be shown when time is printed.
	// Zone name is in form: UTC-08:00.
	// Note that Go Libs do not give access to dst flag, so we ignore dst bits

	i2 = i + 2
	tz = bigen.Uint16(bs[i:i2])
	i = i2
	// sign extend sign bit into top 2 MSB (which were dst bits):
	if tz&(1<<13) == 0 { // positive
		tz = tz & 0x3fff //clear 2 MSBs: dst bits
	} else { // negative
		tz = tz | 0xc000 //set 2 MSBs: dst bits
		//tzname[3] = '-' (TODO: verify. this works here)
	}
	tzint := int16(tz)
	if tzint == 0 {
		tt = time.Unix(tsec, int64(tnsec)).UTC()
	} else {
		// For Go Time, do not use a descriptive timezone. 
		// It's unnecessary, and makes it harder to do a reflect.DeepEqual.
		// The Offset already tells what the offset should be, if not on UTC and unknown zone name.
		// var zoneName = timeLocUTCName(tzint)
		tt = time.Unix(tsec, int64(tnsec)).In(time.FixedZone("", int(tzint)*60))
	}
	return
}

func timeLocUTCName(tzint int16) string {
	if tzint == 0 {
		return "UTC"
	}
	var tzname = []byte("UTC+00:00")
	//tzname := fmt.Sprintf("UTC%s%02d:%02d", tzsign, tz/60, tz%60) //perf issue using Sprintf. inline below.
	//tzhr, tzmin := tz/60, tz%60 //faster if u convert to int first
	var tzhr, tzmin int16
	if tzint < 0 {
		tzname[3] = '-' // (TODO: verify. this works here)
		tzhr, tzmin = -tzint/60, (-tzint)%60
	} else {
		tzhr, tzmin = tzint/60, tzint%60
	}
	tzname[4] = timeDigits[tzhr/10]
	tzname[5] = timeDigits[tzhr%10]
	tzname[7] = timeDigits[tzmin/10]
	tzname[8] = timeDigits[tzmin%10]
	return string(tzname)
	//return time.FixedZone(string(tzname), int(tzint)*60)
}




