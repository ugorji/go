// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReaderI interface {
	// readx will return a view of the []byte if decoding from a []byte, OR
	// read into the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// create a new []byte and read into that
	readx(n uint) []byte

	readb([]byte)

	readxb(n, maxInitLen int, in []byte) (out []byte)

	readn1() byte
	readn2() [2]byte
	readn3() [3]byte
	readn4() [4]byte
	readn8() [8]byte
	// readn1eof() (v uint8, eof bool)

	// // read up to 8 bytes at a time
	// readn(num uint8) (v [8]byte)

	numread() uint // number of bytes read

	// skip any whitespace characters, and return the first non-matching byte
	skipWhitespace() (token byte)

	// jsonReadNum will include last read byte in first element of slice,
	// and continue numeric characters until it sees a non-numeric char
	// or EOF. If it sees a non-numeric character, it will unread that.
	jsonReadNum() []byte

	// jsonReadAsisChars will read json plain characters (anything but " or \)
	// and return a slice terminated by a non-json asis character.
	jsonReadAsisChars() []byte

	// skip will skip any byte that matches, and return the first non-matching byte
	// skip(accept *bitset256) (token byte)

	// readTo will read any byte that matches, stopping once no-longer matching.
	// readTo(accept *bitset256) (out []byte)

	// readUntil will read, only stopping once it matches the 'stop' byte (which it excludes).
	readUntil(stop byte) (out []byte)

	// only supported when reading from bytes
	// bytesReadFrom(startpos uint) []byte

	// isBytes() bool
	resetIO(r io.Reader, bufsize int, blist *bytesFreeList)

	resetBytes(in []byte)

	startRecording(v []byte)
	stopRecording() []byte
}

// ------------------------------------------------

type unreadByteStatus uint8

// unreadByteStatus goes from
// undefined (when initialized) -- (read) --> canUnread -- (unread) --> canRead ...
const (
	unreadByteUndefined unreadByteStatus = iota
	unreadByteCanRead
	unreadByteCanUnread
)

// const defBufReaderSize = 4096

// --------------------

// ioReaderByteScanner contains the io.Reader and io.ByteScanner interfaces
type ioReaderByteScanner interface {
	io.Reader
	io.ByteScanner
	// ReadByte() (byte, error)
	// UnreadByte() error
	// Read(p []byte) (n int, err error)
}

// ioReaderByteScannerImpl does a simple wrapper of a io.ByteScanner
// over a io.Reader
type ioReaderByteScannerImpl struct {
	r io.Reader

	l  byte             // last byte
	ls unreadByteStatus // last byte status

	recording bool
	bufio     bool
	b         [4]byte // tiny buffer for reading single bytes
}

func (z *ioReaderByteScannerImpl) ReadByte() (c byte, err error) {
	if z.ls == unreadByteCanRead {
		z.ls = unreadByteCanUnread
		c = z.l
	} else {
		_, err = z.Read(z.b[:1])
		c = z.b[0]
	}
	return
}

func (z *ioReaderByteScannerImpl) UnreadByte() (err error) {
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

func (z *ioReaderByteScannerImpl) Read(p []byte) (n int, err error) {
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

func (z *ioReaderByteScannerImpl) reset(r io.Reader) {
	z.r = r
	z.ls = unreadByteUndefined
	z.l = 0
}

// ioDecReader is a decReader that reads off an io.Reader.
type ioDecReader struct {
	rr ioReaderByteScannerImpl // the reader passed in, wrapped into a reader+bytescanner
	br ioReaderByteScanner     // main reader used for Read|ReadByte|UnreadByte

	blist *bytesFreeList

	n uint // num read

	bufr []byte // buffer for jsonXXX, readx, readUntil

	rec []byte
}

// var errNoBytesReadFromSupport = errors.New("bytesReadFrom() is not supported")

// func (z *ioDecReader) isBytes() bool {
// 	return false
// }

// func (z *ioDecReader) bytesReadFrom(startpos uint) []byte {
// 	halt.onerror(errNoBytesReadFromSupport)
// 	return nil
// }

func (z *ioDecReader) resetBytes(in []byte) {
	halt.errorStr("resetBytes unsupported by ioDecReader")
}

func (z *ioDecReader) resetIO(r io.Reader, bufsize int, blist *bytesFreeList) {
	z.n = 0
	z.bufr = blist.check(z.bufr, 256)
	z.blist = blist

	// - if r == nil: use eofReader
	// - else if bufsize > 0: use bufio.Reader
	//   - if we own it already and bufsize not change: reset it
	//   - else new bufio.Reader
	// - else if Reader+ByteScanner AND fixed type: use it
	// - else: use wrapper type: ioReaderByteScannerImpl

	if r == nil {
		z.br = &eofReader
		z.rr.bufio = false
		return
	}

	if bufsize > 0 {
		if z.rr.bufio {
			if bb := z.br.(*bufio.Reader); bb.Size() == bufsize {
				bb.Reset(r)
				return
			}
		}
		z.rr.bufio = true
		z.br = bufio.NewReaderSize(r, bufsize)
		return
	}

	z.rr.bufio = false
	switch bb := r.(type) {
	case *strings.Reader:
		z.br = bb
	case *bytes.Buffer:
		z.br = bb
	case *bytes.Reader:
		z.br = bb
	case *bufio.Reader:
		z.br = bb
	default:
		z.rr.reset(r)
		z.br = &z.rr
	}
	return
}

func (z *ioDecReader) numread() uint {
	return z.n
}

func (z *ioDecReader) readn1() (b uint8) {
	b, err := z.br.ReadByte()
	halt.onerror(err)
	if z.rr.recording {
		z.rec = append(z.rec, b)
	}
	z.n++
	return
}

func (z *ioDecReader) readn2() (bs [2]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn3() (bs [3]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn4() (bs [4]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn8() (bs [8]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readx(n uint) (bs []byte) {
	if n == 0 {
		return zeroByteSlice
	}
	if n < uint(cap(z.bufr)) {
		bs = z.bufr[:n]
	} else if z.blist != nil {
		z.bufr = z.blist.putGet(z.bufr, int(n))
		bs = z.bufr[:n]
	} else {
		bs = make([]byte, n)
	}
	nn, err := readFull(z.br, bs)
	if z.rr.recording {
		z.rec = append(z.rec, bs[:nn]...)
	}
	z.n += nn
	halt.onerror(err)
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	nn, err := readFull(z.br, bs)
	if z.rr.recording {
		z.rec = append(z.rec, bs[:nn]...)
	}
	z.n += nn
	halt.onerror(err)
}

func (z *ioDecReader) readxb(n, maxInitLen int, in []byte) (out []byte) {
	if n <= 0 {
		out = zeroByteSlice
	} else if cap(in) >= n {
		out = in[:n]
		z.readb(out)
	} else {
		var len2 int
		for len2 < n {
			len3 := decInferLen(n-len2, maxInitLen, 1)
			bs3 := out
			out = make([]byte, len2+len3)
			copy(out, bs3)
			z.readb(out[len2:])
			len2 += len3
		}
	}
	return
}

// func (z *ioDecReader) readn1eof() (b uint8, eof bool) {
// 	b, err := z.br.ReadByte()
// 	if err == nil {
// 		z.n++
// 	} else if err == io.EOF {
// 		eof = true
// 	} else {
// 		halt.onerror(err)
// 	}
// 	return
// }

func (z *ioDecReader) jsonReadNum() (bs []byte) {
	z.unreadn1()
	z.bufr = z.bufr[:0]
LOOP:
	// i, eof := z.readn1eof()
	i, err := z.br.ReadByte()
	if err == io.EOF {
		return z.bufr
	}
	if err != nil {
		halt.onerror(err)
	}
	if z.rr.recording {
		z.rec = append(z.rec, i)
	}
	z.n++
	if isNumberChar(i) {
		z.bufr = append(z.bufr, i)
		goto LOOP
	}
	z.unreadn1()
	return z.bufr
}

func (z *ioDecReader) jsonReadAsisChars() (bs []byte) {
	z.bufr = z.bufr[:0]
LOOP:
	i := z.readn1()
	z.bufr = append(z.bufr, i)
	if i == '"' || i == '\\' {
		return z.bufr
	}
	goto LOOP
}

func (z *ioDecReader) skipWhitespace() (token byte) {
LOOP:
	token = z.readn1()
	if isWhitespaceChar(token) {
		goto LOOP
	}
	return
}

// func (z *ioDecReader) readUntil(stop byte) []byte {
// 	z.bufr = z.bufr[:0]
// LOOP:
// 	token := z.readn1()
// 	z.bufr = append(z.bufr, token)
// 	if token == stop {
// 		return z.bufr[:len(z.bufr)-1]
// 	}
// 	goto LOOP
// }

func (z *ioDecReader) readUntil(stop byte) []byte {
	z.bufr = z.bufr[:0]
LOOP:
	token := z.readn1()
	if token == stop {
		return z.bufr
	}
	z.bufr = append(z.bufr, token)
	goto LOOP
}

func (z *ioDecReader) unreadn1() {
	err := z.br.UnreadByte()
	halt.onerror(err)
	if z.rr.recording && len(z.rec) > 0 {
		z.rec = z.rec[:len(z.rec)-1]
	}
	z.n--
}

func (z *ioDecReader) startRecording(v []byte) {
	z.rec = v
	z.rr.recording = true
}

func (z *ioDecReader) stopRecording() (v []byte) {
	v = z.rec
	z.rec = nil
	z.rr.recording = false
	return
}

// ------------------------------------

// bytesDecReader is a decReader that reads off a byte slice with zero copying
//
// Note: we do not try to convert index'ing out of bounds to an io error.
// instead, we let it bubble up to the exported Encode/Decode method
// and recover it as an io error.
//
// Every function here MUST defensively check bounds either explicitly
// or via a bounds check.
//
// see panicValToErr(...) function in helper.go.
type bytesDecReader struct {
	b []byte // data
	c uint   // cursor
	r uint   // recording cursor
}

// func (z *bytesDecReader) isBytes() bool {
// 	return true
// }

func (z *bytesDecReader) resetIO(r io.Reader, bufsize int, blist *bytesFreeList) {
	halt.errorStr("resetIO unsupported by bytesDecReader")
}

func (z *bytesDecReader) resetBytes(in []byte) {
	// it's ok to resize a nil slice, so long as it's not past 0
	z.b = in[:len(in):len(in)] // reslicing must not go past capacity
	z.c = 0
}

func (z *bytesDecReader) numread() uint {
	return z.c
}

// func (z *bytesDecReader) bytesReadFrom(startpos uint) []byte {
// 	return z.b[startpos:z.c]
// }

// Note: slicing from a non-constant start position is more expensive,
// as more computation is required to decipher the pointer start position.
// However, we do it only once, and it's better than reslicing both z.b and return value.

func (z *bytesDecReader) readx(n uint) (bs []byte) {
	// x := z.c + n
	// bs = z.b[z.c:x]
	// z.c = x
	bs = z.b[z.c : z.c+n]
	z.c += n
	return
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readx(uint(len(bs))))
}

func (z *bytesDecReader) readxb(n, maxInitLen int, in []byte) (out []byte) {
	if n <= 0 {
		return zeroByteSlice
	}
	_ = z.b[int(z.c)+n-1] // quick fail if out of bounds
	if cap(in) >= n {
		out = in[:n]
	} else {
		out = make([]byte, n)
	}
	z.readb(out)
	return
}

// MARKER: do not use this - as it calls into memmove (as the size of data to move is unknown)
// func (z *bytesDecReader) readnn(bs []byte, n uint) {
// 	x := z.c
// 	copy(bs, z.b[x:x+n])
// 	z.c += n
// }

// func (z *bytesDecReader) readn(num uint8) (bs [8]byte) {
// 	x := z.c + uint(num)
// 	copy(bs[:], z.b[z.c:x]) // slice z.b completely, so we get bounds error if past
// 	z.c = x
// 	return
// }

// func (z *bytesDecReader) readn1() uint8 {
// 	z.c++
// 	return z.b[z.c-1]
// }

// MARKER: readn{1,2,3,4,8} should throw an out of bounds error if past length.
// MARKER: readn1: explicitly ensure bounds check is done
// MARKER: readn{2,3,4,8}: ensure you slice z.b completely so we get bounds error if past end.

func (z *bytesDecReader) readn1() (v uint8) {
	v = z.b[z.c]
	z.c++
	return
}

func (z *bytesDecReader) readn2() (bs [2]byte) {
	// copy(bs[:], z.b[z.c:z.c+2])
	// bs[1] = z.b[z.c+1]
	// bs[0] = z.b[z.c]
	bs = okBytes2(z.b[z.c : z.c+2])
	z.c += 2
	return
}

func (z *bytesDecReader) readn3() (bs [3]byte) {
	// copy(bs[1:], z.b[z.c:z.c+3])
	bs = okBytes3(z.b[z.c : z.c+3])
	z.c += 3
	return
}

func (z *bytesDecReader) readn4() (bs [4]byte) {
	// copy(bs[:], z.b[z.c:z.c+4])
	bs = okBytes4(z.b[z.c : z.c+4])
	z.c += 4
	return
}

func (z *bytesDecReader) readn8() (bs [8]byte) {
	// copy(bs[:], z.b[z.c:z.c+8])
	bs = okBytes8(z.b[z.c : z.c+8])
	z.c += 8
	return
}

func (z *bytesDecReader) jsonReadNum() []byte {
	z.c-- // unread
	i := z.c
LOOP:
	// gracefully handle end of slice, as end of stream is meaningful here
	if i < uint(len(z.b)) && isNumberChar(z.b[i]) {
		i++
		goto LOOP
	}
	z.c, i = i, z.c
	// MARKER: 20230103: byteSliceOf here prevents inlining of jsonReadNum
	// return byteSliceOf(z.b, i, z.c)
	return z.b[i:z.c]
}

func (z *bytesDecReader) jsonReadAsisChars() []byte {
	i := z.c
LOOP:
	token := z.b[i]
	i++
	if token == '"' || token == '\\' {
		z.c, i = i, z.c
		return byteSliceOf(z.b, i, z.c)
		// return z.b[i:z.c]
	}
	goto LOOP
}

func (z *bytesDecReader) skipWhitespace() (token byte) {
	i := z.c
LOOP:
	token = z.b[i]
	if isWhitespaceChar(token) {
		i++
		goto LOOP
	}
	z.c = i + 1
	return
}

func (z *bytesDecReader) readUntil(stop byte) (out []byte) {
	i := z.c
LOOP:
	if z.b[i] == stop {
		out = byteSliceOf(z.b, z.c, i)
		// out = z.b[z.c:i]
		z.c = i + 1
		return
	}
	i++
	goto LOOP
}

func (z *bytesDecReader) startRecording(v []byte) {
	if z.c > 0 {
		z.r = z.c - 1
	}
}

func (z *bytesDecReader) stopRecording() (v []byte) {
	v = z.b[z.r:z.c]
	z.r = 0
	return
}

// --------------

// From out benchmarking, we see the following impact performance:
//
// - functions that are too big to inline
// - interface calls (as no inlining can occur)
//
// decRd is designed to embed a decReader, and then re-implement some of the decReader
// methods using a conditional branch.
//
// We only override the ones where the bytes version is inlined AND the wrapper method
// (containing the bytes version alongside a conditional branch) is also inlined.
//
// We use ./run.sh -z to check.
//
// Right now, only numread and "carefully crafted" readn1 can be inlined.

// func (z *decRd) numread() uint {
// 	if z.bytes {
// 		return z.rb.numread()
// 	}
// 	return z.ri.numread()
// }

// func (z *decRd) readn1() (v uint8) {
// 	if z.bytes {
// 		// return z.rb.readn1()
// 		// MARKER: calling z.rb.readn1() prevents decRd.readn1 from being inlined.
// 		// copy code, to manually inline and explicitly return here.
// 		// Keep in sync with bytesDecReader.readn1
// 		v = z.rb.b[z.rb.c]
// 		z.rb.c++
// 		return
// 	}
// 	return z.ri.readn1()
// }

// func (z *decRd) readn4() [4]byte {
// 	if z.bytes {
// 		return z.rb.readn4()
// 	}
// 	return z.ri.readn4()
// }

// func (z *decRd) readn3() [3]byte {
// 	if z.bytes {
// 		return z.rb.readn3()
// 	}
// 	return z.ri.readn3()
// }

// func (z *decRd) skipWhitespace() byte {
// 	if z.bytes {
// 		return z.rb.skipWhitespace()
// 	}
// 	return z.ri.skipWhitespace()
// }

type devNullReader struct{}

func (devNullReader) Read(p []byte) (int, error) { return 0, io.EOF }
func (devNullReader) Close() error               { return nil }
func (devNullReader) ReadByte() (byte, error)    { return 0, io.EOF }
func (devNullReader) UnreadByte() error          { return io.EOF }

func readFull(r io.Reader, bs []byte) (n uint, err error) {
	var nn int
	for n < uint(len(bs)) && err == nil {
		nn, err = r.Read(bs[n:])
		if nn > 0 {
			if err == io.EOF {
				// leave EOF for next time
				err = nil
			}
			n += uint(nn)
		}
	}
	// do not do this below - it serves no purpose
	// if n != len(bs) && err == io.EOF { err = io.ErrUnexpectedEOF }
	return
}

// var _ decReader = (*decRd)(nil)
