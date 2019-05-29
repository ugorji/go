// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "strconv"

// func parseFloat(b []byte, bitsize int) (f float64, err error) {
// 	if bitsize == 32 {
// 		return parseFloat32(b)
// 	} else {
// 		return parseFloat64(b)
// 	}
// }

func parseFloat32(b []byte) (f float32, err error) {
	return parseFloat32_custom(b)
	// return parseFloat32_strconv(b)
}

func parseFloat64(b []byte) (f float64, err error) {
	return parseFloat64_custom(b)
	// return parseFloat64_strconv(b)
}

func parseFloat32_strconv(b []byte) (f float32, err error) {
	f64, err := strconv.ParseFloat(stringView(b), 32)
	f = float32(f64)
	return
}

func parseFloat64_strconv(b []byte) (f float64, err error) {
	return strconv.ParseFloat(stringView(b), 64)
}

// ------ parseFloat custom below --------

// We assume that a lot of floating point numbers in json files will be
// those that are handwritten, and with defined precision (in terms of number
// of digits after decimal point), etc.
//
// We further assume that this ones can be written in exact format.
//
// strconv.ParseFloat has some unnecessary overhead which we can do without
// for the common case:
//
//    - expensive char-by-char check to see if underscores are in right place
//    - testing for and skipping underscores
//    - check if the string matches ignorecase +/- inf, +/- infinity, nan
//    - support for base 16 (0xFFFF...)
//
// The functions below will try a fast-path for floats which can be decoded
// without any loss of precision, meaning they:
//
//    - fits within the significand bits of the 32-bits or 64-bits
//    - exponent fits within the exponent value
//    - there is no truncation (any extra numbers are all trailing zeros)
//
// To figure out what the values are for maxMantDigits, use this idea below:
//
// 2^23 =                 838 8608 (between 10^ 6 and 10^ 7) (significand bits of uint32)
// 2^32 =             42 9496 7296 (between 10^ 9 and 10^10) (full uint32)
// 2^52 =      4503 5996 2737 0496 (between 10^15 and 10^16) (significand bits of uint64)
// 2^64 = 1844 6744 0737 0955 1616 (between 10^19 and 10^20) (full uint64)
//
// Since we only allow for up to what can comfortably fit into the significand
// ignoring the exponent, and we only try to parse iff significand fits into the

// Exact powers of 10.
var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
var float32pow10 = [...]float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}

type floatinfo struct {
	mantbits uint8
	expbits  uint8
	bias     int16

	exactPow10 int8 // Exact powers of ten are <= 10^N (32: 10, 64: 22)
	exactInts  int8 // Exact integers are <= 10^N

	maxMantDigits int8 // 10^19 fits in uint64, while 10^9 fits in uint32
}

var fi32 = floatinfo{23, 8, -127, 10, 7, 9}     // maxMantDigits = 9
var fi64 = floatinfo{52, 11, -1023, 22, 15, 19} // maxMantDigits = 19

const fMax64 = 1e15
const fMax32 = 1e7

func parseFloatErr(b []byte) error {
	return &strconv.NumError{
		Func: "ParseFloat",
		Err:  strconv.ErrSyntax,
		Num:  string(b),
	}
}

func parseFloat32_custom(b []byte) (f float32, err error) {
	mantissa, exp, i16, neg, _, bad, ok := readFloat(b, fi32)
	if bad {
		return 0, parseFloatErr(b)
	}
	// defer parseFloatDebug(b, 32, &trunc, exp, trunc, ok)
	if ok {
		f = float32(mantissa)
		if neg {
			f = -f
		}
		if exp == 0 {
		} else if exp < 0 { // int / 10^k
			f /= float32pow10[-exp]
		} else { // exp > 0
			i16 = int16(fi32.exactPow10)
			if exp > i16 {
				f *= float32pow10[exp-i16]
				exp = i16
			}
			if f < -fMax32 || f > fMax32 { // exponent may be too large - outside range
				goto FALLBACK
			}
			f *= float32pow10[exp]
		}
		return
	}
FALLBACK:
	return parseFloat32_strconv(b)
}

func parseFloat64_custom(b []byte) (f float64, err error) {
	mantissa, exp, i16, neg, _, bad, ok := readFloat(b, fi64)
	if bad {
		return 0, parseFloatErr(b)
	}
	// defer parseFloatDebug(b, 64, &trunc, exp, trunc, ok)
	if ok {
		f = float64(mantissa)
		if neg {
			f = -f
		}
		if exp == 0 {
		} else if exp < 0 { // int / 10^k
			f /= float64pow10[-exp]
		} else { // exp > 0
			i16 = int16(fi64.exactPow10)
			if exp > i16 {
				f *= float64pow10[exp-i16]
				exp = i16
			}
			if f < -fMax64 || f > fMax64 { // exponent may be too large - outside range
				goto FALLBACK
			}
			f *= float64pow10[exp]
		}
		return
	}
FALLBACK:
	return parseFloat64_strconv(b)
}

func readFloat(s []byte, y floatinfo) (mantissa uint64, exp, i16 int16, neg, trunc, bad, ok bool) {
	var i int
	if len(s) == 0 {
		bad = true
		return
	}
	switch s[0] {
	case '+':
		i++
	case '-':
		neg = true
		i++
	}

	// we considered punting early if string has length > maxMantDigits, but this doesn't account
	// for trailing 0's e.g. 700000000000000000000 can be encoded exactly as it is 7e20

	const base = 10

	// var sawdot, sawdigits, sawexp bool
	var sawdot bool
	var nd, ndMant, dp int8
L:
	for ; i < len(s); i++ {
		switch s[i] {
		case '.':
			if sawdot {
				bad = true
				return
			}
			sawdot = true
			dp = nd
		case '0':
			if nd == 0 { // ignore leading zeros
				dp--
				continue
			}
			nd++
			if ndMant < y.maxMantDigits {
				mantissa *= base
				ndMant++
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// sawdigits = true
			nd++
			if ndMant < y.maxMantDigits {
				mantissa *= base
				mantissa += uint64(s[i] - '0')
				ndMant++
			} else {
				trunc = true
				return // break L
			}
		case 'e', 'E':
			// sawexp = true
			i++
			if i == len(s) {
				break L
			}
			var eneg bool
			if s[i] == '+' {
				i++
			} else if s[i] == '-' {
				i++
				eneg = true
			}
			if i == len(s) {
				break L
			}
			// for exact match, exponent is in single or double digits (-22 to 37 for float64).
			// exit quick if exponent is more than 2 digits.
			if len(s)-i > 2 {
				return
			}
			var e int8
			for ; i < len(s); i++ {
				if s[i] < '0' || s[i] > '9' {
					bad = true
					return
				}
				e = e*base + int8(s[i]-'0')
			}
			if eneg {
				dp -= e
			} else {
				dp += e
			}
			break L
		default:
			bad = true
			return
		}
	}
	// if !sawdigits {
	// 	bad = true
	// 	return
	// }
	if !sawdot {
		dp = nd
	}

	if mantissa != 0 {
		nd = dp - ndMant
		if nd < -y.exactPow10 || nd > y.exactInts+y.exactPow10 { // cannot handle it
			return
		}
		exp = int16(nd)
	}
	ok = true && !trunc && mantissa>>y.mantbits == 0
	return
}

// func parseFloatDebug(b []byte, bitsize int, strconv *bool, exp int16, trunc, ok bool) {
// 	if false && bitsize == 64 {
// 		return
// 	}
// 	if *strconv {
// 		xdebugf("parseFloat%d: delegating: %s, exp: %d, trunc: %v, ok: %v", bitsize, b, exp, trunc, ok)
// 	} else {
// 		xdebug2f("parseFloat%d: attempting: %s, exp: %d, trunc: %v, ok: %v", bitsize, b, exp, trunc, ok)
// 	}
// }
