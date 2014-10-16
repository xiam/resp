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

const (
	StringHeader  = '+'
	ErrorHeader   = '-'
	IntegerHeader = ':'
	BulkHeader    = '$'
	ArrayHeader   = '*'
)

type Message struct {
	Error   error
	Integer int64
	Bytes   []byte
	Status  string
	Array   []*Message
	IsNil   bool
	Type    byte
}

func (m *Message) SetStatus(s string) {
	m.Type = StringHeader
	m.Status = s
}

func (m *Message) SetError(e error) {
	m.Type = ErrorHeader
	m.Error = e
}

func (m *Message) SetInteger(i int64) {
	m.Type = IntegerHeader
	m.Integer = i
}

func (m *Message) SetBytes(b []byte) {
	m.Type = BulkHeader
	m.Bytes = b
}

func (m *Message) SetArray(a []*Message) {
	m.Type = ArrayHeader
	m.Array = a
}

func (m *Message) SetNil() {
	m.Type = 0
	m.IsNil = true
}

func (m Message) Interface() interface{} {
	switch m.Type {
	case ErrorHeader:
		return m.Error
	case IntegerHeader:
		return m.Integer
	case BulkHeader:
		return m.Bytes
	case StringHeader:
		return m.Status
	case ArrayHeader:
		return m.Array
	}
	return nil
}
