// Copyright (c) 2015 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package resp

import (
	"bufio"
	"bytes"
	"io"
)

// Reader reads Redis tokens from an input stream
type Reader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	d := &Reader{
		br: bufio.NewReader(r),
	}
	return d
}

// Read a line of input and its type
func (r *Reader) ReadLine() (lineType byte, line []byte, err error) {
	buf := bytes.NewBuffer(nil)
	end := endOfLine[len(endOfLine)-1]
	for !bytes.HasSuffix(buf.Bytes(), endOfLine) {
		if tmp, err := r.br.ReadBytes(end); err != nil {
			return 0, nil, err
		} else {
			buf.Write(tmp)
		}
	}
	// Line must be at least 1 byte + EOL marker
	if buf.Len() < (1 + len(endOfLine)) {
		return 0, nil, ErrInvalidInput
	}

	if lineType, err = buf.ReadByte(); err != nil {
		return 0, nil, err
	}
	buf.Truncate(buf.Len() - len(endOfLine))
	line = buf.Bytes()
	return lineType, line, nil
}

// Read a message from Redis of length n bytes (not including EOL marker)
func (r *Reader) ReadMessageBytes(n int) (buf []byte, err error) {
	bytesRemaining := n + len(endOfLine)
	buf = make([]byte, bytesRemaining)

	for {
		readStart := len(buf) - bytesRemaining
		var bytesRead int
		if bytesRead, err = r.br.Read(buf[readStart:]); err != nil {
			return nil, err
		}
		if bytesRead == bytesRemaining {
			break
		} else {
			bytesRemaining -= bytesRead
		}
	}
	// Message must terminate in EOL marker
	if !bytes.HasSuffix(buf, endOfLine) {
		return nil, ErrInvalidInput
	}

	// Truncate EOL marker from return buffer
	buf = buf[:n]
	return buf, nil
}
