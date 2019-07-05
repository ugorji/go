// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"strconv"
)

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
	// defer func() { xdebugf("strconv float32: %s, %v, err: %v", b, f, err) }()
	f64, err := strconv.ParseFloat(stringView(b), 32)
	f = float32(f64)
	return
}

func parseFloat64_strconv(b []byte) (f float64, err error) {
	// defer func() { xdebugf("strconv float64: %s, %v, err: %v", b, f, err) }()
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
// Note: we only allow for up to what can comfortably fit into the significand
// ignoring the exponent, and we only try to parse iff significand fits.

// Exact powers of 10.
const (
	thousand    = 1000
	million     = thousand * thousand
	billion     = thousand * million
	trillion    = thousand * billion
	quadrillion = thousand * trillion
	quintillion = thousand * quadrillion
)

var uint64pow10 = [...]uint64{
	1, 10, 100,
	1 * thousand, 10 * thousand, 100 * thousand,
	1 * million, 10 * million, 100 * million,
	1 * billion, 10 * billion, 100 * billion,
	1 * trillion, 10 * trillion, 100 * trillion,
	1 * quadrillion, 10 * quadrillion, 100 * quadrillion,
	1 * quintillion, 10 * quintillion,
}
var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
var float32pow10 = [...]float32{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10,
}

type floatinfo struct {
	mantbits uint8
	expbits  uint8 // (unused)
	bias     int16 // (unused)

	exactPow10 int8 // Exact powers of ten are <= 10^N (32: 10, 64: 22)
	exactInts  int8 // Exact integers are <= 10^N

	maxMantDigits int8 // 10^19 fits in uint64, while 10^9 fits in uint32
}

var fi32 = floatinfo{23, 8, -127, 10, 7, 9}
var fi64 = floatinfo{52, 11, -1023, 22, 15, 19}

var fi64u = floatinfo{64, 0, -1023, 20, 15, 19}
var fi64i = floatinfo{63, 0, -1023, 20, 15, 19}

const fMax64 = 1e15
const fMax32 = 1e7

const fBase = 10

func strconvParseErr(b []byte, fn string) error {
	return &strconv.NumError{
		Func: fn,
		Err:  strconv.ErrSyntax,
		Num:  string(b),
	}
}

func parseFloat32_reader(r readFloatResult) (f float32, fail bool) {
	// parseFloatDebug(b, 32, false, exp, trunc, ok)
	f = float32(r.mantissa)
	if r.exp != 0 {
		indx := fExpIndx(r.exp)
		if r.exp < 0 { // int / 10^k
			f /= float32pow10[indx]
		} else { // exp > 0
			if r.exp > fi32.exactPow10 {
				f *= float32pow10[r.exp-fi32.exactPow10]
				if f > fMax32 { // exponent too large - outside range
					fail = true
					return // ok = false
				}
				indx = uint8(fi32.exactPow10)
			}
			f *= float32pow10[indx]
		}
	}
	if r.neg {
		f = -f
	}
	return
}

func parseFloat32_custom(b []byte) (f float32, err error) {
	r := readFloat(b, fi32)
	if r.bad {
		return 0, strconvParseErr(b, "ParseFloat")
	}
	if r.ok {
		if f, r.bad = parseFloat32_reader(r); r.bad {
			goto FALLBACK
		}
		return
	}
FALLBACK:
	return parseFloat32_strconv(b)
}

func parseFloat64_reader(r readFloatResult) (f float64, fail bool) {
	f = float64(r.mantissa)
	if r.exp != 0 {
		indx := fExpIndx(r.exp)
		if r.exp < 0 { // int / 10^k
			f /= float64pow10[indx]
		} else { // exp > 0
			if r.exp > fi64.exactPow10 {
				f *= float64pow10[r.exp-fi64.exactPow10]
				if f > fMax64 { // exponent too large - outside range
					fail = true
					return
				}
				indx = uint8(fi64.exactPow10)
			}
			f *= float64pow10[indx]
		}
	}
	if r.neg {
		f = -f
	}
	return
}

func parseFloat64_custom(b []byte) (f float64, err error) {
	r := readFloat(b, fi64)
	if r.bad {
		return 0, strconvParseErr(b, "ParseFloat")
	}
	if r.ok {
		if f, r.bad = parseFloat64_reader(r); r.bad {
			goto FALLBACK
		}
		return
	}
FALLBACK:
	return parseFloat64_strconv(b)
}

const maxUint64 = (1<<64 - 1)
const uint64Cutoff = maxUint64/10 + 1

func parseUint64_simple(b []byte) (n uint64, ok bool) {
	for _, c := range b {
		if c < '0' || c > '9' {
			return
		}
		// unsigned integers don't overflow well on multiplication, so check cutoff here
		// e.g. (maxUint64-5)*10 doesn't overflow well ...
		if n >= uint64Cutoff {
			return
		}
		n *= 10
		n1 := n + uint64(c-'0')
		if n1 < n || n1 > maxUint64 {
			return
		}
		n = n1
	}
	ok = true
	return
}

func parseUint64_reader(r readFloatResult) (f uint64, fail bool) {
	f = r.mantissa
	if r.exp != 0 {
		indx := fExpIndx(r.exp)
		if r.exp < 0 { // int / 10^k
			f /= uint64pow10[indx]
		} else { // exp > 0
			if r.exp > fi64u.exactPow10 {
				f *= uint64pow10[r.exp-fi64u.exactPow10]
				// if f > math.MaxUint64 { // exponent too large - outside range
				// 	fail = true
				// 	return
				// }
				indx = uint8(fi64u.exactPow10)
			}
			f *= uint64pow10[indx]
		}
	}
	return
}

func parseUint64(b []byte) (f uint64, err error) {
	var ok bool
	var r readFloatResult
	if b[0] == '-' {
		goto ERROR
	}
	f, ok = parseUint64_simple(b)
	if ok {
		return
	}
	r = readFloat(b, fi64u)
	// if r.neg {
	// 	goto ERROR
	// }
	if r.okInt {
		if f, r.bad = parseUint64_reader(r); r.bad {
			goto ERROR
		}
		return
	}
ERROR:
	err = strconvParseErr(b, "ParseUint")
	return
}

func parseInt64_reader(r readFloatResult) (v int64, fail bool) {
	f := r.mantissa
	if r.exp != 0 {
		indx := fExpIndx(r.exp)
		if r.exp < 0 { // int / 10^k
			f /= uint64pow10[indx]
		} else { // exp > 0
			if r.exp > fi64i.exactPow10 {
				f *= uint64pow10[r.exp-fi64i.exactPow10]
				// if f > math.MaxInt64 { // exponent too large - outside range
				// 	fail = true
				// 	return
				// }
				indx = uint8(fi64i.exactPow10)
			}
			f *= uint64pow10[indx]
		}
	}
	if chkOvf.Uint2Int(f, r.neg) {
		fail = true
		return
	}
	if r.neg {
		v = -int64(f)
	} else {
		v = int64(f)
	}
	return
}

func parseInt64(b []byte) (v int64, err error) {
	var ok, neg bool
	var f uint64
	var r readFloatResult

	if b[0] == '-' {
		neg = true
		b = b[1:]
	}

	f, ok = parseUint64_simple(b)
	if ok {
		if chkOvf.Uint2Int(f, neg) {
			goto ERROR
		}
		if neg {
			v = -int64(f)
		} else {
			v = int64(f)
		}
		return
	}

	r = readFloat(b, fi64i)
	if r.okInt {
		r.neg = neg
		if v, r.bad = parseInt64_reader(r); r.bad {
			goto ERROR
		}
		return
	}
ERROR:
	err = strconvParseErr(b, "ParseInt")
	return
}

// parseNumber will return an integer if only composed of [-]?[0-9]+
// Else it will return a float.
func parseNumber(b []byte, z *fauxUnion, preferSignedInt bool) (err error) {
	var ok, neg bool
	var f uint64
	if b[0] == '-' {
		neg = true
		f, ok = parseUint64_simple(b[1:])
	} else {
		f, ok = parseUint64_simple(b)
	}
	if ok {
		if neg {
			z.v = valueTypeInt
			if chkOvf.Uint2Int(f, neg) {
				err = strconvParseErr(b, "ParseInt")
				return
			}
			z.i = -int64(f)
		} else if preferSignedInt {
			z.v = valueTypeInt
			if chkOvf.Uint2Int(f, neg) {
				err = strconvParseErr(b, "ParseInt")
				return
			}
			z.i = int64(f)
		} else {
			z.v = valueTypeUint
			z.u = f
		}
		return
	}

	z.v = valueTypeFloat
	z.f, err = parseFloat64_custom(b)
	return
}

func fExpIndx(v int8) uint8 {
	if v < 0 {
		return uint8(-v)
	}
	return uint8(v)
}

type readFloatResult struct {
	mantissa                                   uint64
	exp                                        int8
	neg, sawdot, sawexp, trunc, bad, okInt, ok bool
}

func readFloat(s []byte, y floatinfo) (r readFloatResult) {
	var i uint // make it uint, so that we eliminate bounds checking
	var slen = uint(len(s))
	if slen == 0 {
		r.bad = true
		return
	}

	if s[0] == '-' {
		r.neg = true
		i++
	}

	// we considered punting early if string has length > maxMantDigits, but this doesn't account
	// for trailing 0's e.g. 700000000000000000000 can be encoded exactly as it is 7e20

	var nd, ndMant, dp int8

	var c uint8
	for ; i < slen; i++ {
		c = s[i]
		if c == '.' {
			if r.sawdot {
				r.bad = true
				return
			}
			r.sawdot = true
			dp = nd
		} else if c == 'e' || c == 'E' {
			r.sawexp = true
			break
		} else if c >= '0' && c <= '9' { // !(c < '0' || c > '9') { //
			if c == '0' && nd == 0 {
				dp--
				continue
			}
			nd++
			if ndMant >= y.maxMantDigits {
				r.trunc = true
				return
			}
			// mantissa = (mantissa << 1) + (mantissa << 3) + uint64(c-'0')
			r.mantissa = r.mantissa*fBase + uint64(c-'0')
			ndMant++
		} else {
			r.bad = true
			return
		}
	}

	if !r.sawdot {
		dp = nd
	}

	if r.sawexp {
		i++
		if i < slen {
			var eneg bool
			if s[i] == '+' {
				i++
			} else if s[i] == '-' {
				i++
				eneg = true
			}
			if i < slen {
				// for exact match, exponent is 1 or 2 digits (float64: -22 to 37, float32: -1 to 17).
				// exit quick if exponent is more than 2 digits.
				if i+2 < slen {
					return
				}
				var e int8
				if s[i] < '0' || s[i] > '9' {
					r.bad = true
					return
				}
				e = int8(s[i] - '0')
				i++
				if i < slen {
					if s[i] < '0' || s[i] > '9' {
						r.bad = true
						return
					}
					e = e*fBase + int8(s[i]-'0') // (e << 1) + (e << 3) + int8(s[i]-'0')
					i++
				}
				if eneg {
					dp -= e
				} else {
					dp += e
				}
			}
		}
	}

	r.okInt = true

	if r.mantissa != 0 {
		r.exp = dp - ndMant
		if r.exp < -y.exactPow10 ||
			r.exp > y.exactInts+y.exactPow10 ||
			r.mantissa>>y.mantbits != 0 { // cannot handle it
			return
		}
	}
	r.ok = true
	return
}
