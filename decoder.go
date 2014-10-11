// Copyright (c) 2014 José Carlos Nieto, https://menteslibres.net/xiam
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

// RESP Decoder. See: http://redis.io/topics/protocol
package resp

import (
	"errors"
	"io"
	"strconv"
	"sync"
)

type Decoder struct {
	r   io.Reader
	buf []byte
	mu  *sync.Mutex
}

const minRead = 512

func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		r:   r,
		buf: []byte{},
		mu:  &sync.Mutex{},
	}
	return d
}

func (self *Decoder) read(b []byte) (n int, err error) {
	// Read from buffer.
	lb := len(b)
	lz := len(self.buf)

	if lb <= lz {
		self.mu.Lock()
		n = copy(b, self.buf[:lb])
		self.buf = self.buf[lb:]
		self.mu.Unlock()
		return n, nil
	}

	// Read from buffer
	if lz > 0 {
		self.mu.Lock()
		copy(b, self.buf[:lz])
		self.buf = []byte{}
		self.mu.Unlock()
	}

	// ...and from reader.
	r := make([]byte, lb-lz)

	if _, err = self.r.Read(r); err != nil {
		return 0, err
	}

	n = copy(b, r)

	return n, nil
}

func (self *Decoder) readBytes(delim byte) (line []byte, err error) {

	// Filling buffer
	self.mu.Lock()
	for {
		buf := make([]byte, minRead)
		n, _ := self.r.Read(buf)
		if n == 0 {
			break
		}
		self.buf = append(self.buf, buf[0:n]...)
	}
	self.mu.Unlock()

	// Looking for delim
	lb := len(self.buf)

	for i := 0; i < lb; i++ {
		if self.buf[i] == delim {
			c := make([]byte, i+1)
			self.mu.Lock()
			copy(c, self.buf[:i+1])
			self.buf = self.buf[i+1:]
			self.mu.Unlock()
			return c, nil
		}
	}

	return nil, ErrInvalidInput
}

func (self *Decoder) next() (out *Message, n int, err error) {
	//var head []byte
	var line []byte

	// After the header, we expect a message ending with \r\n.
	if line, err = self.readLine(); err != nil {
		return nil, 0, err
	}

	t := line[0]
	line = line[1:]

	switch t {

	case StringHeader:
		out = new(Message)
		out.Type = t
		out.Status = string(line)
		return

	case ErrorHeader:
		out = new(Message)
		out.Type = t
		out.Error = errors.New(string(line))
		return

	case IntegerHeader:
		out = new(Message)
		out.Type = t
		if out.Integer, err = strconv.ParseInt(string(line), 10, 64); err != nil {
			return nil, 0, err
		}
		return

	case BulkHeader:
		// Getting string length.
		var msgLen int

		if msgLen, err = strconv.Atoi(string(line)); err != nil {
			return
		}

		if msgLen > bulkMessageMaxLength {
			err = ErrMessageIsTooLarge
			return
		}

		if msgLen < 0 {
			// RESP Bulk Strings can also be used in order to signal non-existence of
			// a value.
			out = new(Message)
			out.Type = t
			out.IsNil = true
			return
		}

		out = new(Message)
		out.Type = t
		buf := make([]byte, msgLen+2)

		if _, err = self.read(buf); err != nil {
			return nil, 0, err
		}

		out.Bytes = buf[:msgLen]

		return
	case ArrayHeader:
		// Getting string length.
		var arrLen int

		if arrLen, err = strconv.Atoi(string(line)); err != nil {
			return
		}

		if arrLen < 0 {
			// The concept of Null Array exists as well, and is an alternative way to
			// specify a Null value (usually the Null Bulk String is used, but for
			// historical reasons we have two formats).
			out = new(Message)
			out.Type = t
			out.IsNil = true
			return
		}

		out = new(Message)
		out.Type = t
		out.Array = make([]*Message, arrLen)

		for i := 0; i < arrLen; i++ {

			nestedOut, nestedN, nestedErr := self.next()

			if nestedErr != nil {
				return nil, 0, nestedErr
			}

			out.Array[i] = nestedOut

			n = n + nestedN
		}

		return
	}

	err, n = ErrInvalidInput, -1

	return
}

func (self *Decoder) decode() (out *Message, err error) {
	out, _, err = self.next()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (self *Decoder) readLine() (data []byte, err error) {
	var buf []byte
	var chunk []byte

	for {

		if chunk, err = self.readBytes(endOfLine[1]); err != nil {
			return nil, err
		}

		l := len(chunk)

		if l < 2 {
			return nil, ErrInvalidInput
		}

		buf = append(buf, chunk...)

		if chunk[l-2] == '\r' && chunk[l-1] == '\n' {
			break
		}

	}

	n := len(buf)

	return buf[0 : n-2], err
}
