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
	"reflect"
	"strconv"
	"sync"
)

// Decoder reads and decodes RESP objects from an input stream.
type Decoder struct {
	r        io.Reader
	buf      []byte
	lastLine []byte
	off      int
	mu       *sync.Mutex
}

const (
	readLen      = 4096 // Read size.
	lineCapacity = 16   // Typical line capacity (scrach space).
)

// NewDecoder creates and returns a Decoder.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		r:        r,
		buf:      []byte{},
		lastLine: make([]byte, 0, lineCapacity),
		mu:       &sync.Mutex{},
	}
	return d
}

// Sets the initial data of the decoder.
func (self *Decoder) setData(buf []byte) {

	if self.off > 0 {
		self.off = 0
		self.buf = self.buf[:0]
	}

	self.buf = buf
	self.off = 0
}

// Reads the size of b from the buffer.
func (self *Decoder) read(b []byte) (n int, err error) {

	// Requested read size.
	lb := len(b)

	// Available read in buffer.
	lz := len(self.buf) - self.off

	// Can we read from buffer?
	if lb <= lz {
		// Yes!
		self.mu.Lock()

		// Reading from buffer...
		n = copy(b, self.buf[self.off:self.off+lb])
		// ...and advancing offset.
		self.off += lb

		self.mu.Unlock()
		return n, nil
	}

	// No, we should read from the reader.
	if lz > 0 {
		// Copying everything we have...
		self.mu.Lock()

		copy(b, self.buf[self.off:lz])

		// It's a good time to reset our buffer.
		self.buf = []byte{}
		self.off = 0

		self.mu.Unlock()
	}

	// Now let's attempt to read from the reader all that we can.
	if self.r == nil {
		return 0, ErrMissingReader
	}

	r := make([]byte, lb-lz)

	if n, err = self.r.Read(r); err != nil {
		return 0, err
	}

	self.mu.Lock()

	self.buf = append(self.buf, r[:n]...)

	n = copy(b, self.buf[self.off:self.off+n])
	self.off += n

	self.mu.Unlock()

	return n, nil
}

// Reads from the buffer until delim is found.
func (self *Decoder) readBytes(delim byte) (line []byte, err error) {
	hasRead := false

doRead:

	// Attempt to read data from buffer.
	lb := len(self.buf)

	for i := self.off; i < lb; i++ {
		if self.buf[i] == delim {
			c := make([]byte, (i+1)-self.off)
			self.mu.Lock()

			copy(c, self.buf[self.off:i+1])
			self.off = i + 1

			self.mu.Unlock()
			return c, nil
		}
	}

	// Is this our second attempt to read?
	if hasRead {
		return nil, ErrInvalidInput
	}

	// We didn't find the byte we were looking for, let's attempt to read more
	// data from the reader (if any) and try again.

	if self.r == nil {
		// Except that we don't have a reader...
		return nil, ErrInvalidInput
	}

	var n int
	self.mu.Lock()

	buf := make([]byte, readLen)

	if n, err = self.r.Read(buf); err != nil {
		self.mu.Unlock()
		return nil, err
	}

	self.buf = append(self.buf, buf[:n]...)

	self.mu.Unlock()

	hasRead = true

	goto doRead
}

// Attempts to decode the next message.
func (self *Decoder) next(out *Message) (err error) {
	// After the header, we expect a message ending with \r\n.
	if err = self.readLine(); err != nil {
		return err
	}

	out.Type = self.lastLine[0]

	line := self.lastLine[1 : len(self.lastLine)-2]

	switch out.Type {

	case StringHeader:
		out.Status = string(line)
		return

	case ErrorHeader:
		out.Error = errors.New(string(line))
		return

	case IntegerHeader:
		if out.Integer, err = strconv.ParseInt(string(line), 10, 64); err != nil {
			return err
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
			out.IsNil = true
			return
		}

		buf := make([]byte, msgLen+2)

		if _, err = self.read(buf); err != nil {
			return err
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
			out.IsNil = true
			return
		}

		out.Array = make([]*Message, arrLen)

		for i := 0; i < arrLen; i++ {
			out.Array[i] = new(Message)
			if err = self.next(out.Array[i]); err != nil {
				return err
			}
		}

		return
	}

	return ErrInvalidInput
}

// Decode attempts to decode the whole message in buffer.
func (self *Decoder) Decode(v interface{}) (err error) {
	out := new(Message)

	if err = self.next(out); err != nil {
		return err
	}

	if v == nil {
		return ErrExpectingDestination
	}

	dst := reflect.ValueOf(v)

	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		return ErrExpectingPointer
	}

	return redisMessageToType(dst.Elem(), out)
}

// Attempts to read the next line and put it on the self.lastLine space.
func (self *Decoder) readLine() (err error) {
	var chunk []byte
	var n int

	self.lastLine = self.lastLine[:0]

	// self.lastLine = nil

	// Step on every \n and check if the previous char was a \r.
	for {

		// Attempt to read until this character is found.
		if chunk, err = self.readBytes(endOfLine[1]); err != nil {
			return err
		}

		self.lastLine = append(self.lastLine, chunk...)

		// Lenght of the buffer.
		n = len(self.lastLine)

		if n < 2 {
			// Minimal read is two chars: \r\n
			return ErrInvalidInput
		}

		// The character before \n should be \r
		if self.lastLine[n-2] == endOfLine[0] {
			break
		}

	}

	return nil
}
