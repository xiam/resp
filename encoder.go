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

// RESP encoder. See: http://redis.io/topics/protocol
package resp

import (
	"bytes"
)

type encoder struct {
}

var (
	encoderNil = []byte("$-1\r\n")
	digits     = []byte("0123456789")
)

const digitbuflen = 20

func intToBytes(v int) []byte {
	buf := make([]byte, digitbuflen)

	i := len(buf)

	for v >= 10 {
		i--
		buf[i] = digits[v%10]
		v = v / 10
	}

	i--
	buf[i] = digits[v%10]

	return buf[i:]
}

func writeEncoded(buf *bytes.Buffer, in interface{}) (err error) {

	switch v := in.(type) {

	case nil:
		buf.Write(encoderNil)
		return

	case string:
		buf.WriteByte(StringHeader)
		buf.WriteString(v)
		buf.Write(endOfLine)
		return

	case error:
		buf.WriteByte(ErrorHeader)
		buf.WriteString(v.Error())
		buf.Write(endOfLine)
		return

	case int:
		buf.WriteByte(IntegerHeader)
		buf.Write(intToBytes(v))
		buf.Write(endOfLine)
		return

	case []byte:
		buf.WriteByte(BulkHeader)
		buf.Write(intToBytes(len(v)))
		buf.Write(endOfLine)
		buf.Write(v)
		buf.Write(endOfLine)
		return

	case [][]byte:
		buf.WriteByte(ArrayHeader)
		buf.Write(intToBytes(len(v)))
		buf.Write(endOfLine)

		for i := range v {
			buf.WriteByte(BulkHeader)
			buf.Write(intToBytes(len(v[i])))
			buf.Write(endOfLine)
			buf.Write(v[i])
			buf.Write(endOfLine)
		}

		return
	case []string:
		buf.WriteByte(ArrayHeader)
		buf.Write(intToBytes(len(v)))
		buf.Write(endOfLine)

		for i := range v {
			buf.WriteByte(StringHeader)
			buf.WriteString(v[i])
			buf.Write(endOfLine)
		}

		return
	case []int:
		buf.WriteByte(ArrayHeader)
		buf.Write(intToBytes(len(v)))
		buf.Write(endOfLine)

		for i := range v {
			buf.WriteByte(IntegerHeader)
			buf.Write(intToBytes(v[i]))
			buf.Write(endOfLine)
		}

		return
	case []interface{}:
		buf.WriteByte(ArrayHeader)
		buf.Write(intToBytes(len(v)))
		buf.Write(endOfLine)

		for i := range v {
			if err = writeEncoded(buf, v[i]); err != nil {
				return err
			}
		}

		return
	}

	return ErrInvalidInput
}

func (self encoder) encode(in interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := writeEncoded(&buf, in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
