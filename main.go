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

// RESP protocol encoder/decoder.
package resp

import (
	"fmt"
	"reflect"
)

var endOfLine = []byte{'\r', '\n'}

const (
	respStringByte  = '+'
	respErrorByte   = '-'
	respIntegerByte = ':'
	respBulkByte    = '$'
	respArrayByte   = '*'
)

const (
	// Bulk Strings are used in order to represent a single binary safe string up
	// to 512 MB in length.
	bulkMessageMaxLength = 512 * 1024
)

var defaultEncoder = encoder{}
var defaultDecoder = decoder{}

// Marshal returns the RESP encoding of v. At this moment, it only works with
// string, int, []byte, nil and []interface{} types.
func Marshal(v interface{}) ([]byte, error) {
	switch t := v.(type) {
	case string:
		// Strings are not binary safe, we should use bulk type instead.
		return defaultEncoder.encode([]byte(t))
	}
	return defaultEncoder.encode(v)
}

// Unmarshal parses the RESP-encoded data and stores the result in the value
// pointed to by v. At this moment, it only works with string, int, []byte and
// []interface{} types.
func Unmarshal(data []byte, v interface{}) error {

	var out interface{}
	var err error

	dst := reflect.ValueOf(v)

	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		return ErrExpectingPointer
	}

	if out, err = defaultDecoder.decode(data); err != nil {
		return err
	}

	outV := reflect.ValueOf(out)

	// Is this a safe conversion?
	if dst.Elem().Type().Kind() != outV.Type().Kind() {
		return fmt.Errorf(ErrNotSameKind.Error(), dst.Elem().Type().Kind(), outV.Type().Kind())
	}

	dst.Elem().Set(outV)

	return nil
}
