// Copyright (c) 2012-2014 José Carlos Nieto, https://menteslibres.net/xiam
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

// RESP serialization/deserialization protocol.
package resp

import (
	"bytes"
	"errors"
	"strconv"
)

var endOfLine = []byte{'\r', '\n'}

const (
	respStringByte  = '+'
	respErrorByte   = '-'
	respIntegerByte = ':'
	respBulkByte    = '$'
)

const (
	// Bulk Strings are used in order to represent a single binary safe string up
	// to 512 MB in length.
	bulkMessageMaxLength = 512 * 1024
)

type decoder struct {
}

func (self decoder) decode(in []byte) (out interface{}, err error) {

	if len(in) == 0 {
		return nil, ErrInvalidResponse
	}

	switch in[0] {
	case respStringByte:
		var line []byte
		var err error

		if line, err = self.readLine(in[1:]); err != nil {
			return nil, err
		}

		return string(line), nil
	case respErrorByte:
		var line []byte
		var err error

		if line, err = self.readLine(in[1:]); err != nil {
			return nil, err
		}

		return errors.New(string(line)), nil
	case respIntegerByte:
		var line []byte
		var err error
		var res int

		if line, err = self.readLine(in[1:]); err != nil {
			return nil, err
		}

		if res, err = strconv.Atoi(string(line)); err != nil {
			return nil, err
		}

		return res, nil

	case respBulkByte:
		// Getting string length.
		var line []byte
		var err error
		var msgLen, startOffset int

		if line, err = self.readLine(in[1:]); err != nil {
			return nil, err
		}

		if msgLen, err = strconv.Atoi(string(line)); err != nil {
			return nil, err
		}

		if msgLen > bulkMessageMaxLength {
			return nil, ErrMessageIsTooLarge
		}

		if msgLen < 0 {
			// RESP Bulk Strings can also be used in order to signal non-existence of
			// a value.
			return nil, nil
		}

		startOffset = 1 + len(line) + 2 // type + number + \r\n

		if len(in) >= (startOffset + msgLen + 2) { // message start + message length + \r\n
			out := in[startOffset : startOffset+msgLen]
			return string(out), nil
		} else {
			return nil, ErrInvalidResponse
		}
	}

	return nil, ErrInvalidResponse
}

func (self decoder) readLine(in []byte) (out []byte, err error) {
	i := bytes.Index(in, endOfLine)
	if i < 0 {
		return nil, ErrInvalidDelimiter
	}
	return in[0:i], nil
}
