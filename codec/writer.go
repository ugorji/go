// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
)

// encWriter abstracts writing to a byte array or to an io.Writer.
type encWriterI interface {
	writeb([]byte)
	writestr(string)
	writeqstr(string) // write string wrapped in quotes ie "..."
	writen1(byte)

	// add convenience functions for writing 2,4
	writen2(byte, byte)
	writen4([4]byte)
	writen8([8]byte)

	// isBytes() bool
	end()
}

// ---------------------------------------------

type bufioEncWriter struct {
	w io.Writer

	buf []byte

	n int

	b [16]byte // scratch buffer and padding (cache-aligned)
}

type bufioEncWriterM struct {
	*bufioEncWriter
}

func (z *bufioEncWriterM) Make() {
	z.bufioEncWriter = new(bufioEncWriter)
	z.w = io.Discard
}

// // flesh out these writes, to hopefully force full stenciling/monomorphization
// // unfortunately, even with this, the go generics implementation
// // still used a dynamic dispatch, and not static (calling stenciled method)
// func (z bufioEncWriterM) writeb(v []byte)    { z.bufioEncWriter.writeb(v) }
// func (z bufioEncWriterM) writestr(v string)  { z.bufioEncWriter.writestr(v) }
// func (z bufioEncWriterM) writeqstr(v string) { z.bufioEncWriter.writeqstr(v) }
// func (z bufioEncWriterM) writen1(v byte)     { z.bufioEncWriter.writen1(v) }
// func (z bufioEncWriterM) writen2(v, v2 byte) { z.bufioEncWriter.writen2(v, v2) }
// func (z bufioEncWriterM) writen4(v [4]byte)  { z.bufioEncWriter.writen4(v) }
// func (z bufioEncWriterM) writen8(v [8]byte)  { z.bufioEncWriter.writen8(v) }
// func (z bufioEncWriterM) end()               { z.bufioEncWriter.end() }

func (z *bufioEncWriter) reset(w io.Writer, bufsize int, blist *bytesFreelist) {
	z.w = w
	z.n = 0
	if bufsize <= 0 {
		bufsize = defEncByteBufSize
	}
	// bufsize must be >= 8, to accomodate writen methods (where n <= 8)
	if bufsize <= 8 {
		bufsize = 8
	}
	if cap(z.buf) < bufsize {
		if len(z.buf) > 0 && &z.buf[0] != &z.b[0] {
			blist.put(z.buf)
		}
		if len(z.b) > bufsize {
			z.buf = z.b[:]
		} else {
			z.buf = blist.get(bufsize)
		}
	}
	z.buf = z.buf[:cap(z.buf)]
}

func (z *bufioEncWriter) flushErr() (err error) {
	n, err := z.w.Write(z.buf[:z.n])
	z.n -= n
	if z.n > 0 {
		if err == nil {
			err = io.ErrShortWrite
		}
		if n > 0 {
			copy(z.buf, z.buf[n:z.n+n])
		}
	}
	return err
}

func (z *bufioEncWriter) flush() {
	halt.onerror(z.flushErr())
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
	setByteAt(z.buf, uint(z.n), '"')
	// z.buf[z.n] = '"'
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
	setByteAt(z.buf, uint(z.n), '"')
	// z.buf[z.n] = '"'
	z.n++
}

func (z *bufioEncWriter) writen1(b1 byte) {
	if 1 > len(z.buf)-z.n {
		z.flush()
	}
	setByteAt(z.buf, uint(z.n), b1)
	// z.buf[z.n] = b1
	z.n++
}
func (z *bufioEncWriter) writen2(b1, b2 byte) {
	if 2 > len(z.buf)-z.n {
		z.flush()
	}
	setByteAt(z.buf, uint(z.n+1), b2)
	setByteAt(z.buf, uint(z.n), b1)
	// z.buf[z.n+1] = b2
	// z.buf[z.n] = b1
	z.n += 2
}

func (z *bufioEncWriter) writen4(b [4]byte) {
	if 4 > len(z.buf)-z.n {
		z.flush()
	}
	// setByteAt(z.buf, uint(z.n+3), b4)
	// setByteAt(z.buf, uint(z.n+2), b3)
	// setByteAt(z.buf, uint(z.n+1), b2)
	// setByteAt(z.buf, uint(z.n), b1)
	copy(z.buf[z.n:], b[:])
	z.n += 4
}

func (z *bufioEncWriter) writen8(b [8]byte) {
	if 8 > len(z.buf)-z.n {
		z.flush()
	}
	copy(z.buf[z.n:], b[:])
	z.n += 8
}

func (z *bufioEncWriter) endErr() (err error) {
	if z.n > 0 {
		err = z.flushErr()
	}
	return
}

func (z *bufioEncWriter) end() {
	halt.onerror(z.endErr())
}

// ---------------------------------------------

var bytesEncAppenderDefOut = []byte{}

// bytesEncAppender implements encWriter and can write to an byte slice.
type bytesEncAppender struct {
	b   []byte
	out *[]byte
}

type bytesEncAppenderM struct {
	*bytesEncAppender
}

func (z *bytesEncAppenderM) Make() {
	z.bytesEncAppender = new(bytesEncAppender)
	z.out = &bytesEncAppenderDefOut
}

// // flesh out these writes, to hopefully force full stenciling/monomorphization
// // unfortunately, even with this, the go generics implementation
// // still used a dynamic dispatch, and not static (calling stenciled method)
// func (z bytesEncAppenderM) writeb(v []byte)    { z.bytesEncAppender.writeb(v) }
// func (z bytesEncAppenderM) writestr(v string)  { z.bytesEncAppender.writestr(v) }
// func (z bytesEncAppenderM) writeqstr(v string) { z.bytesEncAppender.writeqstr(v) }
// func (z bytesEncAppenderM) writen1(v byte)     { z.bytesEncAppender.writen1(v) }
// func (z bytesEncAppenderM) writen2(v, v2 byte) { z.bytesEncAppender.writen2(v, v2) }
// func (z bytesEncAppenderM) writen4(v [4]byte)  { z.bytesEncAppender.writen4(v) }
// func (z bytesEncAppenderM) writen8(v [8]byte)  { z.bytesEncAppender.writen8(v) }
// func (z bytesEncAppenderM) end()               { z.bytesEncAppender.end() }

func (z *bytesEncAppender) writeb(s []byte) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writestr(s string) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writeqstr(s string) {
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

func (z *bytesEncAppender) writen4(b [4]byte) {
	z.b = append(z.b, b[:]...)
	// z.b = append(z.b, b1, b2, b3, b4) // prevents inlining encWr.writen4
}

func (z *bytesEncAppender) writen8(b [8]byte) {
	z.b = append(z.b, b[:]...)
	// z.b = append(z.b, b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7]) // prevents inlining encWr.writen4
}

// func (z *bytesEncAppender) endErr() error {
// 	*(z.out) = z.b
// 	return nil
// }

func (z *bytesEncAppender) end() {
	*(z.out) = z.b
}

func (z *bytesEncAppender) reset(in []byte, out *[]byte) {
	z.b = in[:0]
	z.out = out
}

// --------------------------------------------------

// func (z *encWr) WriteStr(s string) {
// 	if z.bytes {
// 		z.wb.writestr(s)
// 	} else {
// 		z.wf.writestr(s)
// 	}
// }
