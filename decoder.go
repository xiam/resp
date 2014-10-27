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

// setData sets the initial data of the decoder.
func (d *Decoder) setData(buf []byte) {

	if d.off > 0 {
		d.off = 0
		d.buf = d.buf[:0]
	}

	d.buf = buf
	d.off = 0
}

// Reads the size of b from the buffer.
func (d *Decoder) read(b []byte) (n int, err error) {

	// Requested read size.
	lb := len(b)

	// Available read in buffer.
	lz := len(d.buf) - d.off

	// Can we read from buffer?
	if lb <= lz {
		// Yes!
		d.mu.Lock()

		// Reading from buffer...
		n = copy(b, d.buf[d.off:d.off+lb])
		// ...and advancing offset.
		d.off += lb

		d.mu.Unlock()
		return n, nil
	}

	// No, we should read from the reader.
	if lz > 0 {
		// Copying everything we have...
		d.mu.Lock()

		copy(b, d.buf[d.off:])

		// It's a good time to reset our buffer.
		d.buf = []byte{}
		d.off = 0

		d.mu.Unlock()
	}

	// Now let's attempt to read from the reader all that we can.
	if d.r == nil {
		return 0, ErrMissingReader
	}

	r := make([]byte, lb-lz)

	if n, err = d.r.Read(r); err != nil {
		return 0, err
	}

	d.mu.Lock()

	d.buf = append(d.buf, r[:n]...)

	n = copy(b, d.buf[d.off:d.off+n])
	d.off += n

	d.mu.Unlock()

	return n, nil
}

// Reads from the buffer until delim is found.
func (d *Decoder) readBytes(delim byte) (line []byte, err error) {
	hasRead := false

doRead:

	// Attempt to read data from buffer.
	lb := len(d.buf)

	for i := d.off; i < lb; i++ {
		if d.buf[i] == delim {
			c := make([]byte, (i+1)-d.off)
			d.mu.Lock()

			copy(c, d.buf[d.off:i+1])
			d.off = i + 1

			d.mu.Unlock()
			return c, nil
		}
	}

	// Is this our second attempt to read?
	if hasRead {
		return nil, ErrInvalidInput
	}

	// We didn't find the byte we were looking for, let's attempt to read more
	// data from the reader (if any) and try again.

	if d.r == nil {
		// Except that we don't have a reader...
		return nil, ErrInvalidInput
	}

	var n int
	d.mu.Lock()

	buf := make([]byte, readLen)

	if n, err = d.r.Read(buf); err != nil {
		d.mu.Unlock()
		return nil, err
	}

	d.buf = append(d.buf, buf[:n]...)

	d.mu.Unlock()

	hasRead = true

	goto doRead
}

// Attempts to decode the next message.
func (d *Decoder) next(out *Message) (err error) {
	// After the header, we expect a message ending with \r\n.
	if err = d.readLine(); err != nil {
		return err
	}

	if len(d.lastLine) < 3 {
		return ErrInvalidInput
	}

	out.Type = d.lastLine[0]

	line := d.lastLine[1 : len(d.lastLine)-2]

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

		if _, err = d.read(buf); err != nil {
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
			if err = d.next(out.Array[i]); err != nil {
				return err
			}
		}

		return
	}

	return ErrInvalidInput
}

// Decode attempts to decode the whole message in buffer.
func (d *Decoder) Decode(v interface{}) (err error) {
	out := new(Message)

	if err = d.next(out); err != nil {
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

// Attempts to read the next line and put it on the d.lastLine space.
func (d *Decoder) readLine() (err error) {
	var chunk []byte
	var n int

	d.lastLine = d.lastLine[:0]

	// d.lastLine = nil

	// Step on every \n and check if the previous char was a \r.
	for {

		// Attempt to read until this character is found.
		if chunk, err = d.readBytes(endOfLine[1]); err != nil {
			return err
		}

		d.lastLine = append(d.lastLine, chunk...)

		// Lenght of the buffer.
		n = len(d.lastLine)

		if n < 2 {
			// Minimal read is two chars: \r\n
			return ErrInvalidInput
		}

		// The character before \n should be \r
		if d.lastLine[n-2] == endOfLine[0] {
			break
		}

	}

	return nil
}
