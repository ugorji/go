// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "io"

/*

// encWriter abstracts writing to a byte array or to an io.Writer.
//
//
// Deprecated: Use encWriterSwitch instead.
type encWriter interface {
	writeb([]byte)
	writestr(string)
	writeqstr(string) // write string wrapped in quotes ie "..."
	writen1(byte)
	writen2(byte, byte)
	end()
}

*/

// type ioEncWriterWriter interface {
// 	WriteByte(c byte) error
// 	WriteString(s string) (n int, err error)
// 	Write(p []byte) (n int, err error)
// }

// ---------------------------------------------

/*

type ioEncStringWriter interface {
	WriteString(s string) (n int, err error)
}

// ioEncWriter implements encWriter and can write to an io.Writer implementation
type ioEncWriter struct {
	w  io.Writer
	ww io.Writer
	bw io.ByteWriter
	sw ioEncStringWriter
	fw ioFlusher
	b  [8]byte
}

func (z *ioEncWriter) reset(w io.Writer) {
	z.w = w
	var ok bool
	if z.bw, ok = w.(io.ByteWriter); !ok {
		z.bw = z
	}
	if z.sw, ok = w.(ioEncStringWriter); !ok {
		z.sw = z
	}
	z.fw, _ = w.(ioFlusher)
	z.ww = w
}

func (z *ioEncWriter) WriteByte(b byte) (err error) {
	z.b[0] = b
	_, err = z.w.Write(z.b[:1])
	return
}

func (z *ioEncWriter) WriteString(s string) (n int, err error) {
	return z.w.Write(bytesView(s))
}

func (z *ioEncWriter) writeb(bs []byte) {
	if _, err := z.ww.Write(bs); err != nil {
		panic(err)
	}
}

func (z *ioEncWriter) writestr(s string) {
	if _, err := z.sw.WriteString(s); err != nil {
		panic(err)
	}
}

func (z *ioEncWriter) writeqstr(s string) {
	writestr("\"" + s + "\"")
}

func (z *ioEncWriter) writen1(b byte) {
	if err := z.bw.WriteByte(b); err != nil {
		panic(err)
	}
}

func (z *ioEncWriter) writen2(b1, b2 byte) {
	var err error
	if err = z.bw.WriteByte(b1); err == nil {
		if err = z.bw.WriteByte(b2); err == nil {
			return
		}
	}
	panic(err)
}

// func (z *ioEncWriter) writen5(b1, b2, b3, b4, b5 byte) {
// 	z.b[0], z.b[1], z.b[2], z.b[3], z.b[4] = b1, b2, b3, b4, b5
// 	if _, err := z.ww.Write(z.b[:5]); err != nil {
// 		panic(err)
// 	}
// }

//go:noinline - so *encWriterSwitch.XXX has the bytesEncAppender.XXX inlined
func (z *ioEncWriter) end() {
	if z.fw != nil {
		if err := z.fw.Flush(); err != nil {
			panic(err)
		}
	}
}

*/

// ---------------------------------------------

// bufioEncWriter
type bufioEncWriter struct {
	w io.Writer

	buf []byte

	n int

	// Extensions can call Encode() within a current Encode() call.
	// We need to know when the top level Encode() call returns,
	// so we can decide whether to Release() or not.
	calls uint16 // what depth in mustDecode are we in now.

	sz int // buf size
	// _ uint64 // padding (cache-aligned)

	// ---- cache line

	// write-most fields below

	// less used fields
	bytesBufPooler

	b [40]byte // scratch buffer and padding (cache-aligned)
	// a int
	// b   [4]byte
	// err
}

func (z *bufioEncWriter) reset(w io.Writer, bufsize int) {
	z.w = w
	z.n = 0
	z.calls = 0
	if bufsize <= 0 {
		bufsize = defEncByteBufSize
	}
	z.sz = bufsize
	if cap(z.buf) >= bufsize {
		z.buf = z.buf[:cap(z.buf)]
	} else if bufsize <= len(z.b) {
		z.buf = z.b[:]
	} else {
		z.buf = z.bytesBufPooler.get(bufsize)
		z.buf = z.buf[:cap(z.buf)]
		// z.buf = make([]byte, bufsize)
	}
}

func (z *bufioEncWriter) release() {
	z.buf = nil
	z.bytesBufPooler.end()
}

//go:noinline - flush only called intermittently
func (z *bufioEncWriter) flushErr() (err error) {
	n, err := z.w.Write(z.buf[:z.n])
	z.n -= n
	if z.n > 0 && err == nil {
		err = io.ErrShortWrite
	}
	if n > 0 && z.n > 0 {
		copy(z.buf, z.buf[n:z.n+n])
	}
	return err
}

func (z *bufioEncWriter) flush() {
	if err := z.flushErr(); err != nil {
		panic(err)
	}
}

func (z *bufioEncWriter) writeb(s []byte) {
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
}

func (z *bufioEncWriter) writestr(s string) {
	// z.writeb(bytesView(s)) // inlined below
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
}

func (z *bufioEncWriter) writeqstr(s string) {
	// z.writen1('"')
	// z.writestr(s)
	// z.writen1('"')

	if z.n+len(s)+2 > len(z.buf) {
		z.flush()
	}
	z.buf[z.n] = '"'
	z.n++
LOOP:
	a := len(z.buf) - z.n
	if len(s)+1 > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
	z.buf[z.n] = '"'
	z.n++
}

func (z *bufioEncWriter) writen1(b1 byte) {
	if 1 > len(z.buf)-z.n {
		z.flush()
	}
	z.buf[z.n] = b1
	z.n++
}

func (z *bufioEncWriter) writen2(b1, b2 byte) {
	if 2 > len(z.buf)-z.n {
		z.flush()
	}
	z.buf[z.n+1] = b2
	z.buf[z.n] = b1
	z.n += 2
}

func (z *bufioEncWriter) endErr() (err error) {
	if z.n > 0 {
		err = z.flushErr()
	}
	return
}

// ---------------------------------------------

// bytesEncAppender implements encWriter and can write to an byte slice.
type bytesEncAppender struct {
	b   []byte
	out *[]byte
}

func (z *bytesEncAppender) writeb(s []byte) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writestr(s string) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writeqstr(s string) {
	// z.writen1('"')
	// z.writestr(s)
	// z.writen1('"')

	z.b = append(append(append(z.b, '"'), s...), '"')

	// z.b = append(z.b, '"')
	// z.b = append(z.b, s...)
	// z.b = append(z.b, '"')
}
func (z *bytesEncAppender) writen1(b1 byte) {
	z.b = append(z.b, b1)
}
func (z *bytesEncAppender) writen2(b1, b2 byte) {
	z.b = append(z.b, b1, b2)
}
func (z *bytesEncAppender) endErr() error {
	*(z.out) = z.b
	return nil
}
func (z *bytesEncAppender) reset(in []byte, out *[]byte) {
	z.b = in[:0]
	z.out = out
}

// --------------------------------------------------

type encWriterSwitch struct {
	esep  bool // whether it has elem separators
	bytes bool // encoding to []byte
	isas  bool // whether e.as != nil
	js    bool // is json encoder?
	be    bool // is binary encoder?

	c containerState
	// _    [3]byte // padding
	// _    [2]uint64 // padding
	// _    uint64    // padding
	// wi   *ioEncWriter
	wb bytesEncAppender
	wf *bufioEncWriter
	// typ  entryType
}

func (z *encWriterSwitch) writeb(s []byte) {
	if z.bytes {
		z.wb.writeb(s)
	} else {
		z.wf.writeb(s)
	}
}
func (z *encWriterSwitch) writeqstr(s string) {
	if z.bytes {
		z.wb.writeqstr(s)
	} else {
		z.wf.writeqstr(s)
	}
}
func (z *encWriterSwitch) writestr(s string) {
	if z.bytes {
		z.wb.writestr(s)
	} else {
		z.wf.writestr(s)
	}
}
func (z *encWriterSwitch) writen1(b1 byte) {
	if z.bytes {
		z.wb.writen1(b1)
	} else {
		z.wf.writen1(b1)
	}
}
func (z *encWriterSwitch) writen2(b1, b2 byte) {
	if z.bytes {
		z.wb.writen2(b1, b2)
	} else {
		z.wf.writen2(b1, b2)
	}
}
func (z *encWriterSwitch) endErr() error {
	if z.bytes {
		return z.wb.endErr()
	}
	return z.wf.endErr()
}

func (z *encWriterSwitch) end() {
	if err := z.endErr(); err != nil {
		panic(err)
	}
}

/*

// ------------------------------------------
func (z *encWriterSwitch) writeb(s []byte) {
	switch z.typ {
	case entryTypeBytes:
		z.wb.writeb(s)
	case entryTypeIo:
		z.wi.writeb(s)
	default:
		z.wf.writeb(s)
	}
}
func (z *encWriterSwitch) writestr(s string) {
	switch z.typ {
	case entryTypeBytes:
		z.wb.writestr(s)
	case entryTypeIo:
		z.wi.writestr(s)
	default:
		z.wf.writestr(s)
	}
}
func (z *encWriterSwitch) writen1(b1 byte) {
	switch z.typ {
	case entryTypeBytes:
		z.wb.writen1(b1)
	case entryTypeIo:
		z.wi.writen1(b1)
	default:
		z.wf.writen1(b1)
	}
}
func (z *encWriterSwitch) writen2(b1, b2 byte) {
	switch z.typ {
	case entryTypeBytes:
		z.wb.writen2(b1, b2)
	case entryTypeIo:
		z.wi.writen2(b1, b2)
	default:
		z.wf.writen2(b1, b2)
	}
}
func (z *encWriterSwitch) end() {
	switch z.typ {
	case entryTypeBytes:
		z.wb.end()
	case entryTypeIo:
		z.wi.end()
	default:
		z.wf.end()
	}
}

// ------------------------------------------
func (z *encWriterSwitch) writeb(s []byte) {
	if z.bytes {
		z.wb.writeb(s)
	} else {
		z.wi.writeb(s)
	}
}
func (z *encWriterSwitch) writestr(s string) {
	if z.bytes {
		z.wb.writestr(s)
	} else {
		z.wi.writestr(s)
	}
}
func (z *encWriterSwitch) writen1(b1 byte) {
	if z.bytes {
		z.wb.writen1(b1)
	} else {
		z.wi.writen1(b1)
	}
}
func (z *encWriterSwitch) writen2(b1, b2 byte) {
	if z.bytes {
		z.wb.writen2(b1, b2)
	} else {
		z.wi.writen2(b1, b2)
	}
}
func (z *encWriterSwitch) end() {
	if z.bytes {
		z.wb.end()
	} else {
		z.wi.end()
	}
}

*/
