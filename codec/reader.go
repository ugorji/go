// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "io"

/*

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
//
// Deprecated: Use decReaderSwitch instead.
type decReader interface {
	unreadn1()
	// readx will use the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// just return a view of the []byte being decoded from.
	// Ensure you call detachZeroCopyBytes later if this needs to be sent outside codec control.
	readx(n int) []byte
	readb([]byte)
	readn1() uint8
	numread() uint // number of bytes read
	track()
	stopTrack() []byte

	// skip will skip any byte that matches, and return the first non-matching byte
	skip(accept *bitset256) (token byte)
	// readTo will read any byte that matches, stopping once no-longer matching.
	readTo(accept *bitset256) (out []byte)
	// readUntil will read, only stopping once it matches the 'stop' byte.
	readUntil(stop byte) (out []byte)
}

*/

// ------------------------------------------------

type unreadByteStatus uint8

// unreadByteStatus goes from
// undefined (when initialized) -- (read) --> canUnread -- (unread) --> canRead ...
const (
	unreadByteUndefined unreadByteStatus = iota
	unreadByteCanRead
	unreadByteCanUnread
)

// func appendPool(bs []byte, b byte, bufp *bytesBufPooler) []byte {
// 	if cap(bs)-len(bs) < 1 {
// 		bs = bufp.ensureCap(len(bs)+1, bs)
// 	}
// 	bs = append(bs, b)
// 	return bs
// }

// func appendPoolMany(bs []byte, b []byte, bufp *bytesBufPooler) []byte {
// 	if cap(bs)-len(bs) < 1 {
// 		bs = bufp.ensureCap(len(bs)+1, bs)
// 	}
// 	bs = append(bs, b...)
// 	return bs
// }

// --------------------

type ioDecReaderCommon struct {
	r io.Reader // the reader passed in

	n uint // num read

	l   byte             // last byte
	ls  unreadByteStatus // last byte status
	trb bool             // tracking bytes turned on
	_   bool
	b   [4]byte // tiny buffer for reading single bytes

	tr   bytesBufPoolerPlus // buffer for tracking bytes
	bufr bytesBufPoolerPlus // buffer for readTo/readUntil
}

func (z *ioDecReaderCommon) reset(r io.Reader) {
	z.r = r
	z.ls = unreadByteUndefined
	z.l, z.n = 0, 0
	z.trb = false
	z.tr.end()
}

func (z *ioDecReaderCommon) numread() uint {
	return z.n
}

func (z *ioDecReaderCommon) track() {
	z.tr.ensureCap(1)
	z.tr.buf = z.tr.buf[:0]
	z.trb = true
}

func (z *ioDecReaderCommon) stopTrack() (bs []byte) {
	z.trb = false
	return z.tr.buf
}

func (z *ioDecReaderCommon) resetBufr() {
	z.bufr.ensureCap(1)
	z.bufr.buf = z.bufr.buf[:0]
}

func (z *ioDecReaderCommon) release() {
	z.tr.end()
	z.bufr.end()
}

// ------------------------------------------

// ioDecReader is a decReader that reads off an io.Reader.
//
// It also has a fallback implementation of ByteScanner if needed.
type ioDecReader struct {
	ioDecReaderCommon

	rr io.Reader
	br io.ByteScanner

	x [scratchByteArrayLen]byte // for: get struct field name, swallow valueTypeBytes, etc
	// _ [1]uint64                 // padding
}

func (z *ioDecReader) reset(r io.Reader) {
	z.ioDecReaderCommon.reset(r)

	var ok bool
	z.rr = r
	z.br, ok = r.(io.ByteScanner)
	if !ok {
		z.br = z
		z.rr = z
	}
}

func (z *ioDecReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	var firstByte bool
	if z.ls == unreadByteCanRead {
		z.ls = unreadByteCanUnread
		p[0] = z.l
		if len(p) == 1 {
			n = 1
			return
		}
		firstByte = true
		p = p[1:]
	}
	n, err = z.r.Read(p)
	if n > 0 {
		if err == io.EOF && n == len(p) {
			err = nil // read was successful, so postpone EOF (till next time)
		}
		z.l = p[n-1]
		z.ls = unreadByteCanUnread
	}
	if firstByte {
		n++
	}
	return
}

func (z *ioDecReader) ReadByte() (c byte, err error) {
	n, err := z.Read(z.b[:1])
	if n == 1 {
		c = z.b[0]
		if err == io.EOF {
			err = nil // read was successful, so postpone EOF (till next time)
		}
	}
	return
}

func (z *ioDecReader) UnreadByte() (err error) {
	switch z.ls {
	case unreadByteCanUnread:
		z.ls = unreadByteCanRead
	case unreadByteCanRead:
		err = errDecUnreadByteLastByteNotRead
	case unreadByteUndefined:
		err = errDecUnreadByteNothingToRead
	default:
		err = errDecUnreadByteUnknown
	}
	return
}

func (z *ioDecReader) readx(n uint) (bs []byte) {
	if n == 0 {
		return
	}
	if n < uint(len(z.x)) {
		bs = z.x[:n]
	} else {
		bs = make([]byte, n)
	}
	if _, err := decReadFull(z.rr, bs); err != nil {
		panic(err)
	}
	z.n += uint(len(bs))
	if z.trb {
		z.tr.appends(bs)
	}
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	if _, err := decReadFull(z.rr, bs); err != nil {
		panic(err)
	}
	z.n += uint(len(bs))
	if z.trb {
		z.tr.appends(bs)
	}
}

func (z *ioDecReader) readn1eof() (b uint8, eof bool) {
	b, err := z.br.ReadByte()
	if err == nil {
		z.n++
		if z.trb {
			z.tr.append(b)
		}
	} else if err == io.EOF {
		eof = true
	} else {
		panic(err)
	}
	return
}

func (z *ioDecReader) readn1() (b uint8) {
	b, err := z.br.ReadByte()
	if err == nil {
		z.n++
		if z.trb {
			z.tr.append(b)
		}
		return
	}
	panic(err)
}

func (z *ioDecReader) skip(accept *bitset256) (token byte) {
	var eof bool
	// for {
	// 	token, eof = z.readn1eof()
	// 	if eof {
	// 		return
	// 	}
	// 	if accept.isset(token) {
	// 		continue
	// 	}
	// 	return
	// }
LOOP:
	token, eof = z.readn1eof()
	if eof {
		return
	}
	if accept.isset(token) {
		goto LOOP
	}
	return
}

func (z *ioDecReader) readTo(accept *bitset256) []byte {
	// out = in

	// for {
	// 	token, eof := z.readn1eof()
	// 	if eof {
	// 		return
	// 	}
	// 	if accept.isset(token) {
	// 		out = append(out, token)
	// 	} else {
	// 		z.unreadn1()
	// 		return
	// 	}
	// }
	z.resetBufr()
LOOP:
	token, eof := z.readn1eof()
	if eof {
		return z.bufr.buf
	}
	if accept.isset(token) {
		// out = append(out, token)
		z.bufr.append(token)
		goto LOOP
	}
	z.unreadn1()
	return z.bufr.buf
}

func (z *ioDecReader) readUntil(stop byte) []byte {
	// for {
	// 	token, eof := z.readn1eof()
	// 	if eof {
	// 		panic(io.EOF)
	// 	}
	// 	out = append(out, token)
	// 	if token == stop {
	// 		return
	// 	}
	// }
	z.resetBufr()
LOOP:
	token, eof := z.readn1eof()
	if eof {
		panic(io.EOF)
	}
	z.bufr.append(token)
	if token == stop {
		return z.bufr.buf
	}
	goto LOOP
}

//go:noinline
func (z *ioDecReader) unreadn1() {
	err := z.br.UnreadByte()
	if err != nil {
		panic(err)
	}
	z.n--
	if z.trb {
		if l := len(z.tr.buf) - 1; l >= 0 {
			z.tr.buf = z.tr.buf[:l]
		}
	}
}

// ------------------------------------

type bufioDecReader struct {
	ioDecReaderCommon
	// _ uint64 // padding (cache-aligned)

	c uint // cursor
	// buf []byte

	bp bytesBufPoolerPlus

	// err error
}

func (z *bufioDecReader) reset(r io.Reader, bufsize int) {
	z.ioDecReaderCommon.reset(r)
	z.c = 0
	if cap(z.bp.buf) < bufsize {
		z.bp.get(bufsize)
		// z.buf = make([]byte, 0, bufsize)
	}
	z.bp.buf = z.bp.buf[:0]
}

func (z *bufioDecReader) release() {
	z.ioDecReaderCommon.release()
	z.bp.end()
}

func (z *bufioDecReader) readb(p []byte) {
	var n = uint(copy(p, z.bp.buf[z.c:]))
	z.n += n
	z.c += n
	if len(p) == int(n) {
		if z.trb {
			z.tr.appends(p)
		}
	} else {
		z.readbFill(p, n)
	}
}

//go:noinline - fallback when z.bp.buf is consumed
func (z *bufioDecReader) readbFill(p0 []byte, n uint) {
	// at this point, there's nothing in z.bp.buf to read (z.bp.buf is fully consumed)
	p := p0[n:]
	var n2 uint
	var err error
	if len(p) > cap(z.bp.buf) {
		n2, err = decReadFull(z.r, p)
		if err != nil {
			panic(err)
		}
		n += n2
		z.n += n2
		// always keep last byte in z.bp.buf
		z.bp.buf = z.bp.buf[:1]
		z.bp.buf[0] = p[len(p)-1]
		z.c = 1
		if z.trb {
			z.tr.appends(p0[:n])
		}
		return
	}
	// z.c is now 0, and len(p) <= cap(z.bp.buf)
LOOP:
	// for len(p) > 0 && z.err == nil {
	if len(p) > 0 {
		z.bp.buf = z.bp.buf[0:cap(z.bp.buf)]
		var n1 int
		n1, err = z.r.Read(z.bp.buf)
		n2 = uint(n1)
		if n2 == 0 && err != nil {
			panic(err)
		}
		z.bp.buf = z.bp.buf[:n2]
		n2 = uint(copy(p, z.bp.buf))
		z.c = n2
		n += n2
		z.n += n2
		p = p[n2:]
		goto LOOP
	}
	if z.c == 0 {
		z.bp.buf = z.bp.buf[:1]
		z.bp.buf[0] = p[len(p)-1]
		z.c = 1
	}
	if z.trb {
		z.tr.appends(p0[:n])
	}
}

func (z *bufioDecReader) readn1() (b byte) {
	// fast-path, so we elide calling into Read() most of the time
	if z.c < uint(len(z.bp.buf)) {
		b = z.bp.buf[z.c]
		z.c++
		z.n++
		if z.trb {
			z.tr.append(b)
		}
	} else { // meaning z.c == len(z.bp.buf) or greater ... so need to fill
		z.readbFill(z.b[:1], 0)
		b = z.b[0]
	}
	return
}

func (z *bufioDecReader) unreadn1() {
	if z.c == 0 {
		panic(errDecUnreadByteNothingToRead)
	}
	z.c--
	z.n--
	if z.trb {
		z.tr.buf = z.tr.buf[:len(z.tr.buf)-1]
	}
}

func (z *bufioDecReader) readx(n uint) (bs []byte) {
	if n == 0 {
		// return
	} else if z.c+n <= uint(len(z.bp.buf)) {
		bs = z.bp.buf[z.c : z.c+n]
		z.n += n
		z.c += n
		if z.trb {
			z.tr.appends(bs)
		}
	} else {
		bs = make([]byte, n)
		// n no longer used - can reuse
		n = uint(copy(bs, z.bp.buf[z.c:]))
		z.n += n
		z.c += n
		z.readbFill(bs, n)
	}
	return
}

// func (z *bufioDecReader) doTrack(y uint) {
// 	z.tr = append(z.tr, z.bp.buf[z.c:y]...) // cost=14???
// }

// func (z *bufioDecReader) skipLoopFn(i uint) {
// 	z.n += (i - z.c) - 1
// 	i++
// 	if z.trb {
// 		// z.tr = append(z.tr, z.bp.buf[z.c:i]...)
// 		z.doTrack(i)
// 	}
// 	z.c = i
// }

func (z *bufioDecReader) skip(accept *bitset256) (token byte) {
	// token, _ = z.search(nil, accept, 0, 1); return

	// for i := z.c; i < len(z.bp.buf); i++ {
	// 	if token = z.bp.buf[i]; !accept.isset(token) {
	// 		z.skipLoopFn(i)
	// 		return
	// 	}
	// }

	i := z.c
LOOP:
	if i < uint(len(z.bp.buf)) {
		// inline z.skipLoopFn(i) and refactor, so cost is within inline budget
		token = z.bp.buf[i]
		i++
		if accept.isset(token) {
			goto LOOP
		}
		z.n += i - 2 - z.c
		if z.trb {
			z.tr.appends(z.bp.buf[z.c:i]) // z.doTrack(i)
		}
		z.c = i
		return
	}
	return z.skipFill(accept)
}

func (z *bufioDecReader) skipFill(accept *bitset256) (token byte) {
	// defer func() { xdebugf("skipFill '%c'", token) }()
	z.n += uint(len(z.bp.buf)) - z.c
	if z.trb {
		z.tr.appends(z.bp.buf[z.c:])
	}
	var i, n2 int
	var err error
	for {
		z.c = 0
		z.bp.buf = z.bp.buf[0:cap(z.bp.buf)]
		n2, err = z.r.Read(z.bp.buf)
		if n2 == 0 && err != nil {
			panic(err)
		}
		z.bp.buf = z.bp.buf[:n2]
		for i, token = range z.bp.buf {
			if !accept.isset(token) {
				z.n += (uint(i) - z.c) - 1
				z.loopFn(uint(i + 1))
				return
			}
		}
		// for i := 0; i < n2; i++ {
		// 	if token = z.bp.buf[i]; !accept.isset(token) {
		// 		z.skipLoopFn(i)
		// 		return
		// 	}
		// }
		z.n += uint(n2)
		if z.trb {
			z.tr.appends(z.bp.buf)
		}
	}
}

// func (z *bufioDecReader) readLoopFn(i uint, out0 []byte) (out []byte) {
// 	out = appendPool(out0, z.bp.buf[z.c:i]...)
// 	z.loopFn(i)
// 	return
// }

func (z *bufioDecReader) loopFn(i uint) {
	if z.trb {
		z.tr.appends(z.bp.buf[z.c:i]) // z.doTrack(i)
	}
	z.c = i
}

// func (z *bufioDecReader) readToLoopFn(i uint, out0 []byte) (out []byte) {
// 	// out0 is never nil
// 	z.n += (i - z.c) - 1
// 	return z.readLoopFn(i, out0)
// }

func (z *bufioDecReader) readTo(accept *bitset256) (out []byte) {
	// defer func() { xdebug2f("bufio: readTo: %s", out) }()
	// _, out = z.search(in, accept, 0, 2); return

	// for i := z.c; i < len(z.bp.buf); i++ {
	// 	if !accept.isset(z.bp.buf[i]) {
	// 		return z.readToLoopFn(i, nil)
	// 	}
	// }

	i := z.c
LOOP:
	if i < uint(len(z.bp.buf)) {
		if !accept.isset(z.bp.buf[i]) {
			// return z.readToLoopFn(i, nil)
			// inline readToLoopFn here (for performance)
			z.n += (i - z.c) - 1
			out = z.bp.buf[z.c:i]
			if z.trb {
				z.tr.appends(z.bp.buf[z.c:i]) // z.doTrack(i)
			}
			z.c = i
			return
		}
		i++
		goto LOOP
	}
	return z.readToFill(accept)
}

func (z *bufioDecReader) readToFill(accept *bitset256) []byte {
	z.resetBufr()
	z.n += uint(len(z.bp.buf)) - z.c
	z.bufr.appends(z.bp.buf[z.c:])
	if z.trb {
		z.tr.appends(z.bp.buf[z.c:])
	}
	var n2 int
	var err error
	for {
		z.c = 0
		z.bp.buf = z.bp.buf[:cap(z.bp.buf)]
		n2, err = z.r.Read(z.bp.buf)
		if n2 == 0 && err != nil {
			if err == io.EOF {
				return z.bufr.buf // readTo should read until it matches or end is reached
			}
			panic(err)
		}
		z.bp.buf = z.bp.buf[:n2]
		for i, token := range z.bp.buf {
			if !accept.isset(token) {
				z.n += (uint(i) - z.c) - 1
				z.bufr.appends(z.bp.buf[z.c:i])
				z.loopFn(uint(i))
				return z.bufr.buf
			}
		}
		// for i := 0; i < n2; i++ {
		// 	if !accept.isset(z.bp.buf[i]) {
		// 		return z.readToLoopFn(i, out)
		// 	}
		// }
		z.bufr.appends(z.bp.buf)
		z.n += uint(n2)
		if z.trb {
			z.tr.appends(z.bp.buf)
		}
	}
}

// func (z *bufioDecReader) readUntilLoopFn(i uint, out0 []byte) (out []byte) {
// 	z.n += (i - z.c) - 1
// 	return z.readLoopFn(i+1, out0)
// }

func (z *bufioDecReader) readUntil(stop byte) (out []byte) {
	// defer func() { xdebug2f("bufio: readUntil: %s", out) }()
	// _, out = z.search(in, nil, stop, 4); return

	// for i := z.c; i < len(z.bp.buf); i++ {
	// 	if z.bp.buf[i] == stop {
	// 		return z.readUntilLoopFn(i, nil)
	// 	}
	// }

	i := z.c
LOOP:
	if i < uint(len(z.bp.buf)) {
		if z.bp.buf[i] == stop {
			// inline readUntilLoopFn
			// return z.readUntilLoopFn(i, nil)
			z.n += (i - z.c) - 1
			i++
			out = z.bp.buf[z.c:i]
			if z.trb {
				z.tr.appends(z.bp.buf[z.c:i]) // z.doTrack(i)
			}
			z.c = i
			return
		}
		i++
		goto LOOP
	}
	return z.readUntilFill(stop)
}

func (z *bufioDecReader) readUntilFill(stop byte) []byte {
	z.resetBufr()
	z.n += uint(len(z.bp.buf)) - z.c
	z.bufr.appends(z.bp.buf[z.c:])
	if z.trb {
		z.tr.appends(z.bp.buf[z.c:])
	}
	for {
		z.c = 0
		z.bp.buf = z.bp.buf[0:cap(z.bp.buf)]
		n1, err := z.r.Read(z.bp.buf)
		if n1 == 0 && err != nil {
			panic(err)
		}
		n2 := uint(n1)
		z.bp.buf = z.bp.buf[:n2]
		for i, token := range z.bp.buf {
			if token == stop {
				z.n += (uint(i) - z.c) - 1
				z.bufr.appends(z.bp.buf[z.c : i+1])
				z.loopFn(uint(i + 1))
				return z.bufr.buf
			}
		}
		// for i := 0; i < n2; i++ {
		// 	if z.bp.buf[i] == stop {
		// 		return z.readUntilLoopFn(i, out)
		// 	}
		// }
		z.bufr.appends(z.bp.buf)
		z.n += n2
		if z.trb {
			z.tr.appends(z.bp.buf)
		}
	}
}

// ------------------------------------

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct {
	b []byte // data
	c uint   // cursor
	t uint   // track start
	// a int    // available
}

func (z *bytesDecReader) reset(in []byte) {
	z.b = in
	// z.a = len(in)
	z.c = 0
	z.t = 0
}

func (z *bytesDecReader) numread() uint {
	return z.c
}

func (z *bytesDecReader) unreadn1() {
	if z.c == 0 || len(z.b) == 0 {
		panic(errBytesDecReaderCannotUnread)
	}
	z.c--
	// z.a++
}

func (z *bytesDecReader) readx(n uint) (bs []byte) {
	// slicing from a non-constant start position is more expensive,
	// as more computation is required to decipher the pointer start position.
	// However, we do it only once, and it's better than reslicing both z.b and return value.

	// if n <= 0 {
	// } else if z.a == 0 {
	// 	panic(io.EOF)
	// } else if n > z.a {
	// 	panic(io.ErrUnexpectedEOF)
	// } else {
	// 	c0 := z.c
	// 	z.c = c0 + n
	// 	z.a = z.a - n
	// 	bs = z.b[c0:z.c]
	// }
	// return

	if n != 0 {
		z.c += n
		if z.c > uint(len(z.b)) {
			z.c = uint(len(z.b))
			panic(io.EOF)
		}
		bs = z.b[z.c-n : z.c]
	}
	return

	// if n == 0 {
	// } else if z.c+n > uint(len(z.b)) {
	// 	z.c = uint(len(z.b))
	// 	panic(io.EOF)
	// } else {
	// 	z.c += n
	// 	bs = z.b[z.c-n : z.c]
	// }
	// return

	// if n == 0 {
	// 	return
	// }
	// if z.c == uint(len(z.b)) {
	// 	panic(io.EOF)
	// }
	// if z.c+n > uint(len(z.b)) {
	// 	panic(io.ErrUnexpectedEOF)
	// }
	// // z.a -= n
	// z.c += n
	// return z.b[z.c-n : z.c]
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readx(uint(len(bs))))
}

func (z *bytesDecReader) readn1() (v uint8) {
	if z.c == uint(len(z.b)) {
		panic(io.EOF)
	}
	v = z.b[z.c]
	z.c++
	// z.a--
	return
}

// func (z *bytesDecReader) readn1eof() (v uint8, eof bool) {
// 	if z.a == 0 {
// 		eof = true
// 		return
// 	}
// 	v = z.b[z.c]
// 	z.c++
// 	z.a--
// 	return
// }

func (z *bytesDecReader) skip(accept *bitset256) (token byte) {
	i := z.c
	// if i == len(z.b) {
	// 	goto END
	// 	// panic(io.EOF)
	// }

	// Replace loop with goto construct, so that this can be inlined
	// for i := z.c; i < blen; i++ {
	// 	if !accept.isset(z.b[i]) {
	// 		token = z.b[i]
	// 		i++
	// 		z.a -= (i - z.c)
	// 		z.c = i
	// 		return
	// 	}
	// }

	// i := z.c
LOOP:
	if i < uint(len(z.b)) {
		token = z.b[i]
		i++
		if accept.isset(token) {
			goto LOOP
		}
		// z.a -= (i - z.c)
		z.c = i
		return
	}
	// END:
	panic(io.EOF)
	// // z.a = 0
	// z.c = blen
	// return
}

func (z *bytesDecReader) readTo(accept *bitset256) (out []byte) {
	i := z.c
	if i == uint(len(z.b)) {
		panic(io.EOF)
	}

	// Replace loop with goto construct, so that this can be inlined
	// for i := z.c; i < blen; i++ {
	// 	if !accept.isset(z.b[i]) {
	// 		out = z.b[z.c:i]
	// 		z.a -= (i - z.c)
	// 		z.c = i
	// 		return
	// 	}
	// }
	// out = z.b[z.c:]
	// z.a, z.c = 0, blen
	// return

	// 	i := z.c
	// LOOP:
	// 	if i < blen {
	// 		if accept.isset(z.b[i]) {
	// 			i++
	// 			goto LOOP
	// 		}
	// 		out = z.b[z.c:i]
	// 		z.a -= (i - z.c)
	// 		z.c = i
	// 		return
	// 	}
	// 	out = z.b[z.c:]
	// 	// z.a, z.c = 0, blen
	// 	z.a = 0
	// 	z.c = blen
	// 	return

	// c := i
LOOP:
	if i < uint(len(z.b)) {
		if accept.isset(z.b[i]) {
			i++
			goto LOOP
		}
	}

	out = z.b[z.c:i]
	// z.a -= (i - z.c)
	z.c = i
	return // z.b[c:i]
	// z.c, i = i, z.c
	// return z.b[i:z.c]
}

func (z *bytesDecReader) readUntil(stop byte) (out []byte) {
	i := z.c
	// if i == len(z.b) {
	// 	panic(io.EOF)
	// }

	// Replace loop with goto construct, so that this can be inlined
	// for i := z.c; i < blen; i++ {
	// 	if z.b[i] == stop {
	// 		i++
	// 		out = z.b[z.c:i]
	// 		z.a -= (i - z.c)
	// 		z.c = i
	// 		return
	// 	}
	// }
LOOP:
	if i < uint(len(z.b)) {
		if z.b[i] == stop {
			i++
			out = z.b[z.c:i]
			// z.a -= (i - z.c)
			z.c = i
			return
		}
		i++
		goto LOOP
	}
	// z.a = 0
	// z.c = blen
	panic(io.EOF)
}

func (z *bytesDecReader) track() {
	z.t = z.c
}

func (z *bytesDecReader) stopTrack() (bs []byte) {
	return z.b[z.t:z.c]
}

// --------------

type decReaderSwitch struct {
	esep     bool // has elem separators
	mtr, str bool // whether maptype or slicetype are known types

	be   bool // is binary encoding
	js   bool // is json handle
	jsms bool // is json handle, and MapKeyAsString

	// typ   entryType
	bytes bool // is bytes reader
	bufio bool // is this a bufioDecReader?

	rb bytesDecReader
	ri *ioDecReader
	bi *bufioDecReader
}

// numread, track and stopTrack are always inlined, as they just check int fields, etc.

/*
func (z *decReaderSwitch) numread() int {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.numread()
	case entryTypeIo:
		return z.ri.numread()
	default:
		return z.bi.numread()
	}
}
func (z *decReaderSwitch) track() {
	switch z.typ {
	case entryTypeBytes:
		z.rb.track()
	case entryTypeIo:
		z.ri.track()
	default:
		z.bi.track()
	}
}
func (z *decReaderSwitch) stopTrack() []byte {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.stopTrack()
	case entryTypeIo:
		return z.ri.stopTrack()
	default:
		return z.bi.stopTrack()
	}
}

func (z *decReaderSwitch) unreadn1() {
	switch z.typ {
	case entryTypeBytes:
		z.rb.unreadn1()
	case entryTypeIo:
		z.ri.unreadn1()
	default:
		z.bi.unreadn1()
	}
}
func (z *decReaderSwitch) readx(n int) []byte {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.readx(n)
	case entryTypeIo:
		return z.ri.readx(n)
	default:
		return z.bi.readx(n)
	}
}
func (z *decReaderSwitch) readb(s []byte) {
	switch z.typ {
	case entryTypeBytes:
		z.rb.readb(s)
	case entryTypeIo:
		z.ri.readb(s)
	default:
		z.bi.readb(s)
	}
}
func (z *decReaderSwitch) readn1() uint8 {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.readn1()
	case entryTypeIo:
		return z.ri.readn1()
	default:
		return z.bi.readn1()
	}
}
func (z *decReaderSwitch) skip(accept *bitset256) (token byte) {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.skip(accept)
	case entryTypeIo:
		return z.ri.skip(accept)
	default:
		return z.bi.skip(accept)
	}
}
func (z *decReaderSwitch) readTo(accept *bitset256) (out []byte) {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.readTo(accept)
	case entryTypeIo:
		return z.ri.readTo(accept)
	default:
		return z.bi.readTo(accept)
	}
}
func (z *decReaderSwitch) readUntil(stop byte) (out []byte) {
	switch z.typ {
	case entryTypeBytes:
		return z.rb.readUntil(stop)
	case entryTypeIo:
		return z.ri.readUntil(stop)
	default:
		return z.bi.readUntil(stop)
	}
}

*/

// the if/else-if/else block is expensive to inline.
// Each node of this construct costs a lot and dominates the budget.
// Best to only do an if fast-path else block (so fast-path is inlined).
// This is irrespective of inlineExtraCallCost set in $GOROOT/src/cmd/compile/internal/gc/inl.go
//
// In decReaderSwitch methods below, we delegate all IO functions into their own methods.
// This allows for the inlining of the common path when z.bytes=true.
// Go 1.12+ supports inlining methods with up to 1 inlined function (or 2 if no other constructs).

func (z *decReaderSwitch) release() {
	if z.bytes {
	} else if z.bufio {
		z.bi.release()
	} else {
		z.ri.release()
	}
}
func (z *decReaderSwitch) numread() uint {
	if z.bytes {
		return z.rb.numread()
	} else if z.bufio {
		return z.bi.numread()
	} else {
		return z.ri.numread()
	}
}
func (z *decReaderSwitch) track() {
	if z.bytes {
		z.rb.track()
	} else if z.bufio {
		z.bi.track()
	} else {
		z.ri.track()
	}
}
func (z *decReaderSwitch) stopTrack() []byte {
	if z.bytes {
		return z.rb.stopTrack()
	} else if z.bufio {
		return z.bi.stopTrack()
	} else {
		return z.ri.stopTrack()
	}
}

// func (z *decReaderSwitch) unreadn1() {
// 	if z.bytes {
// 		z.rb.unreadn1()
// 	} else {
// 		z.unreadn1IO()
// 	}
// }
// func (z *decReaderSwitch) unreadn1IO() {
// 	if z.bufio {
// 		z.bi.unreadn1()
// 	} else {
// 		z.ri.unreadn1()
// 	}
// }

func (z *decReaderSwitch) unreadn1() {
	if z.bytes {
		z.rb.unreadn1()
	} else if z.bufio {
		z.bi.unreadn1()
	} else {
		z.ri.unreadn1() // not inlined
	}
}

func (z *decReaderSwitch) readx(n uint) []byte {
	if z.bytes {
		return z.rb.readx(n)
	}
	return z.readxIO(n)
}
func (z *decReaderSwitch) readxIO(n uint) []byte {
	if z.bufio {
		return z.bi.readx(n)
	}
	return z.ri.readx(n)
}

func (z *decReaderSwitch) readb(s []byte) {
	if z.bytes {
		z.rb.readb(s)
	} else {
		z.readbIO(s)
	}
}

//go:noinline - fallback for io, ensures z.bytes path is inlined
func (z *decReaderSwitch) readbIO(s []byte) {
	if z.bufio {
		z.bi.readb(s)
	} else {
		z.ri.readb(s)
	}
}

func (z *decReaderSwitch) readn1() uint8 {
	if z.bytes {
		return z.rb.readn1()
	}
	return z.readn1IO()
}
func (z *decReaderSwitch) readn1IO() uint8 {
	if z.bufio {
		return z.bi.readn1()
	}
	return z.ri.readn1()
}

func (z *decReaderSwitch) skip(accept *bitset256) (token byte) {
	if z.bytes {
		return z.rb.skip(accept)
	}
	return z.skipIO(accept)
}
func (z *decReaderSwitch) skipIO(accept *bitset256) (token byte) {
	if z.bufio {
		return z.bi.skip(accept)
	}
	return z.ri.skip(accept)
}

func (z *decReaderSwitch) readTo(accept *bitset256) (out []byte) {
	if z.bytes {
		return z.rb.readTo(accept)
	}
	return z.readToIO(accept)
}

//go:noinline - fallback for io, ensures z.bytes path is inlined
func (z *decReaderSwitch) readToIO(accept *bitset256) (out []byte) {
	if z.bufio {
		return z.bi.readTo(accept)
	}
	return z.ri.readTo(accept)
}
func (z *decReaderSwitch) readUntil(stop byte) (out []byte) {
	if z.bytes {
		return z.rb.readUntil(stop)
	}
	return z.readUntilIO(stop)
}

func (z *decReaderSwitch) readUntilIO(stop byte) (out []byte) {
	if z.bufio {
		return z.bi.readUntil(stop)
	}
	return z.ri.readUntil(stop)
}

// // register these here, so that staticcheck stops barfing
// var _ = (*bytesDecReader).readTo
// var _ = (*bytesDecReader).readUntil
