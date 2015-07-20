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
)

// Decoder reads and decodes RESP objects from an input stream.
type Decoder struct {
	r *Reader
}

// NewDecoder creates and returns a Decoder.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		r: NewReader(r),
	}
	return d
}

// Attempts to decode the next message.
func (d *Decoder) next(out *Message) (err error) {
	// After the header, we expect a message ending with \r\n.
	var line []byte
	if out.Type, line, err = d.r.ReadLine(); err != nil {
		return err
	}

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

		if out.Bytes, err = d.r.ReadMessageBytes(msgLen); err != nil {
			return
		}

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

	if err = redisMessageToType(dst.Elem(), out); err != nil {
		if out.Type == ErrorHeader {
			return errors.New(out.Error.Error())
		}
	}

	return err
}
