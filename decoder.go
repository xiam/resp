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

// RESP decoder. See: http://redis.io/topics/protocol
package resp

import (
	"bytes"
	"errors"
	"strconv"
)

type decoder struct {
}

func (self decoder) next(in []byte) (out interface{}, n int, err error) {

	if len(in) < 1 {
		return nil, 0, ErrInvalidInput
	}

	var line []byte

	if line, n, err = self.readLine(in[1:]); err != nil {
		return nil, 0, err
	}

	n = n + 1 // line length + type byte

	switch in[0] {

	case respStringByte:
		return string(line), n, nil

	case respErrorByte:
		return errors.New(string(line)), n, nil

	case respIntegerByte:
		var res int
		if res, err = strconv.Atoi(string(line)); err != nil {
			return nil, 0, err
		}
		return res, n, nil

	case respBulkByte:
		// Getting string length.
		var msgLen, offset int

		if msgLen, err = strconv.Atoi(string(line)); err != nil {
			return nil, 0, err
		}

		if msgLen > bulkMessageMaxLength {
			return nil, 0, ErrMessageIsTooLarge
		}

		if msgLen < 0 {
			// RESP Bulk Strings can also be used in order to signal non-existence of
			// a value.
			return nil, n, nil
		}

		offset = 1 + len(line) + 2 // type + number + \r\n

		if len(in) >= (offset + msgLen + 2) { // message start + message length + \r\n
			out := in[offset : offset+msgLen]
			return string(out), offset + msgLen + 2, nil
		} else {
			return nil, 0, ErrInvalidInput
		}
	case respArrayByte:
		// Getting string length.
		var arrLen, offset int
		var res []interface{}

		if arrLen, err = strconv.Atoi(string(line)); err != nil {
			return nil, 0, err
		}

		if arrLen < 0 {
			// The concept of Null Array exists as well, and is an alternative way to
			// specify a Null value (usually the Null Bulk String is used, but for
			// historical reasons we have two formats).
			return nil, n, nil
		}

		offset = 1 + len(line) + 2 // type + number + \r\n

		for i := 0; i < arrLen; i++ {

			if len(in) < offset {
				return nil, 0, ErrIncompleteMessage
			}

			out, n, err := self.next(in[offset:])

			if err != nil {
				return nil, 0, err
			}

			res = append(res, out)

			offset = offset + n
		}

		return res, offset, nil
	}

	return nil, -1, ErrInvalidInput
}

func (self decoder) decode(in []byte) (out interface{}, err error) {
	out, _, err = self.next(in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (self decoder) readLine(in []byte) (out []byte, n int, err error) {
	i := bytes.Index(in, endOfLine)
	if i < 0 {
		return nil, 0, ErrInvalidDelimiter
	}
	return in[0:i], i + 2, nil // header + content + \r\n
}
