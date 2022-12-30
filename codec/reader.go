// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bufio"
	"bytes"
	"io"
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface {
	// readx will return a view of the []byte if decoding from a []byte, OR
	// read into the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// create a new []byte and read into that
	readx(n uint) []byte

	readb([]byte)

	readn1() byte
	readn2() [2]byte
	// readn3 will read 3 bytes into the top-most elements of a 4-byte array
	readn3() [4]byte
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

// ioReaderByteScannerT does a simple wrapper of a io.ByteScanner
// over a io.Reader
type ioReaderByteScannerT struct {
	r io.Reader

	l  byte             // last byte
	ls unreadByteStatus // last byte status

	_ [2]byte // padding
	b [4]byte // tiny buffer for reading single bytes
}

func (z *ioReaderByteScannerT) ReadByte() (c byte, err error) {
	if z.ls == unreadByteCanRead {
		z.ls = unreadByteCanUnread
		c = z.l
	} else {
		_, err = z.Read(z.b[:1])
		c = z.b[0]
	}
	return
}

func (z *ioReaderByteScannerT) UnreadByte() (err error) {
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

func (z *ioReaderByteScannerT) Read(p []byte) (n int, err error) {
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

func (z *ioReaderByteScannerT) reset(r io.Reader) {
	z.r = r
	z.ls = unreadByteUndefined
	z.l = 0
}

// ioDecReader is a decReader that reads off an io.Reader.
type ioDecReader struct {
	rr ioReaderByteScannerT // the reader passed in, wrapped into a reader+bytescanner

	n uint // num read

	blist *bytesFreelist

	bufr []byte              // buffer for readTo/readUntil
	br   ioReaderByteScanner // main reader used for Read|ReadByte|UnreadByte
	bb   *bufio.Reader       // created internally, and reused on reset if needed

	x [64 + 48]byte // for: get struct field name, swallow valueTypeBytes, etc
}

func (z *ioDecReader) reset(r io.Reader, bufsize int, blist *bytesFreelist) {
	z.blist = blist
	z.n = 0
	z.bufr = z.blist.check(z.bufr, 256)
	z.br = nil

	var ok bool

	if bufsize <= 0 {
		z.br, ok = r.(ioReaderByteScanner)
		if !ok {
			z.rr.reset(r)
			z.br = &z.rr
		}
		return
	}

	// bufsize > 0 ...

	// if bytes.[Buffer|Reader], no value in adding extra buffer
	// if bufio.Reader, no value in extra buffer unless size changes
	switch bb := r.(type) {
	case *bytes.Buffer:
		z.br = bb
	case *bytes.Reader:
		z.br = bb
	case *bufio.Reader:
		if bb.Size() == bufsize {
			z.br = bb
		}
	}

	if z.br == nil {
		if z.bb != nil && z.bb.Size() == bufsize {
			z.bb.Reset(r)
		} else {
			z.bb = bufio.NewReaderSize(r, bufsize)
		}
		z.br = z.bb
	}
}

func (z *ioDecReader) numread() uint {
	return z.n
}

func (z *ioDecReader) readn2() (bs [2]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn3() (bs [4]byte) {
	z.readb(bs[1:])
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
		return
	}
	if n < uint(len(z.x)) {
		bs = z.x[:n]
	} else {
		bs = make([]byte, n)
	}
	_, err := readFull(z.br, bs)
	halt.onerror(err)
	z.n += uint(len(bs))
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	_, err := readFull(z.br, bs)
	halt.onerror(err)
	z.n += uint(len(bs))
}

func (z *ioDecReader) readn1() (b uint8) {
	b, err := z.br.ReadByte()
	halt.onerror(err)
	z.n++
	return
}

func (z *ioDecReader) readn1eof() (b uint8, eof bool) {
	b, err := z.br.ReadByte()
	if err == nil {
		z.n++
	} else if err == io.EOF {
		eof = true
	} else {
		halt.onerror(err)
	}
	return
}

func (z *ioDecReader) jsonReadNum() (bs []byte) {
	z.unreadn1()
	z.bufr = z.bufr[:0]
LOOP:
	i, eof := z.readn1eof()
	if eof {
		return z.bufr
	}
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

func (z *ioDecReader) readUntil(stop byte) []byte {
	z.bufr = z.bufr[:0]
LOOP:
	token := z.readn1()
	z.bufr = append(z.bufr, token)
	if token == stop {
		return z.bufr[:len(z.bufr)-1]
	}
	goto LOOP
}

func (z *ioDecReader) unreadn1() {
	err := z.br.UnreadByte()
	halt.onerror(err)
	z.n--
}

// ------------------------------------

// bytesDecReader is a decReader that reads off a byte slice with zero copying
//
// Note: we do not try to convert index'ing out of bounds to an io.EOF.
// instead, we let it bubble up to the exported Encode/Decode method
// and recover it as an io.EOF.
//
// see panicValToErr(...) function in helper.go.
type bytesDecReader struct {
	b []byte // data
	c uint   // cursor
}

func (z *bytesDecReader) reset(in []byte) {
	z.b = in[:len(in):len(in)] // reslicing must not go past capacity
	z.c = 0
}

func (z *bytesDecReader) numread() uint {
	return z.c
}

// Note: slicing from a non-constant start position is more expensive,
// as more computation is required to decipher the pointer start position.
// However, we do it only once, and it's better than reslicing both z.b and return value.

func (z *bytesDecReader) readx(n uint) (bs []byte) {
	x := z.c + n
	bs = z.b[z.c:x]
	z.c = x
	return
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readx(uint(len(bs))))
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

func (z *bytesDecReader) readn1() (v uint8) {
	v = z.b[z.c]
	z.c++
	return
}

// MARKER: for readn{2,3,4,8}, ensure you slice z.b completely so we get bounds error if past end.

func (z *bytesDecReader) readn2() (bs [2]byte) {
	// copy(bs[:], z.b[z.c:z.c+2])
	bs[1] = z.b[z.c+1]
	bs[0] = z.b[z.c]
	z.c += 2
	return
}

func (z *bytesDecReader) readn3() (bs [4]byte) {
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
	z.c--
	i := z.c
LOOP:
	if i < uint(len(z.b)) && isNumberChar(z.b[i]) {
		i++
		goto LOOP
	}
	z.c, i = i, z.c
	return z.b[i:z.c]
}

func (z *bytesDecReader) jsonReadAsisChars() []byte {
	i := z.c
LOOP:
	token := z.b[i]
	i++
	if token == '"' || token == '\\' {
		z.c, i = i, z.c
		return z.b[i:z.c]
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
		out = z.b[z.c:i]
		z.c = i + 1
		return
	}
	i++
	goto LOOP
}

// --------------

type decRd struct {
	mtr bool // is maptype a known type?
	str bool // is slicetype a known type?

	be   bool // is binary encoding
	js   bool // is json handle
	jsms bool // is json handle, and MapKeyAsString
	cbor bool // is cbor handle

	bytes bool // is bytes reader

	rb bytesDecReader
	ri *ioDecReader

	decReader
}

// From out benchmarking, we see the following in terms of performance:
//
// - interface calls
// - branch that can inline what it calls
//
// the if/else-if/else block is expensive to inline.
// Each node of this construct costs a lot and dominates the budget.
// Best to only do an if fast-path else block (so fast-path is inlined).
// This is irrespective of inlineExtraCallCost set in $GOROOT/src/cmd/compile/internal/gc/inl.go
//
// In decRd methods below, we delegate all IO functions into their own methods.
// This allows for the inlining of the common path when z.bytes=true.
// Go 1.12+ supports inlining methods with up to 1 inlined function (or 2 if no other constructs).
//
// However, up through Go 1.13, decRd's readXXX, skip and unreadXXX methods are not inlined.
// Consequently, there is no benefit to do the xxxIO methods for decRd at this time.
// Instead, we have a if/else-if/else block so that IO calls do not have to jump through
// a second unnecessary function call.
//
// If golang inlining gets better and bytesDecReader methods can be inlined,
// then we can revert to using these 2 functions so the bytesDecReader
// methods are inlined and the IO paths call out to a function.
//
// decRd is designed to embed a decReader, and then re-implement some of the decReader
// methods using a conditional branch. We only override the ones that have a bytes version
// that is small enough to be inlined. We use ./run.sh -z to check.
// Right now, only numread and readn1 can be inlined.

func (z *decRd) numread() uint {
	if z.bytes {
		return z.rb.numread()
	} else {
		return z.ri.numread()
	}
}

func (z *decRd) readn1() (v uint8) {
	if z.bytes {
		// MARKER: manually inline, else this function is not inlined.
		// Keep in sync with bytesDecReader.readn1
		// return z.rb.readn1()
		v = z.rb.b[z.rb.c]
		z.rb.c++
	} else {
		v = z.ri.readn1()
	}
	return
}

type devNullReader struct{}

func (devNullReader) Read(p []byte) (int, error) { return 0, io.EOF }
func (devNullReader) Close() error               { return nil }

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

var _ decReader = (*decRd)(nil)
