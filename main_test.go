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
	"bytes"
	"errors"
	"testing"
)

var respDecoder = decoder{}

var (
	errTestFailed    = errors.New("Test failed.")
	errErrorExpected = errors.New("An error was expected.")
)

func TestReadLine(t *testing.T) {
	var test []byte
	var err error
	var n int

	if test, n, err = respDecoder.readLine([]byte("+OK\r\n")); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte("+OK"), test) == false {
		t.Fatal(errTestFailed)
	}

	if n != 5 {
		t.Fatal(errTestFailed)
	}

	if test, n, err = respDecoder.readLine([]byte("+OK")); err == nil {
		t.Fatal(errErrorExpected)
	}

	if test != nil {
		t.Fatal(errTestFailed)
	}

	if n != 0 {
		t.Fatal(errTestFailed)
	}
}

func TestDecodeString(t *testing.T) {
	var test interface{}
	var encoded []byte
	var err error

	// Simple "OK" string
	encoded = []byte("+OK\r\n")

	if test, err = respDecoder.decode(encoded); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("OK")) == false {
		t.Fatal(errTestFailed)
	}

	// Two encoded strings, must get the first one.
	encoded = []byte("+OK\r\n+NO\r\n")

	if test, err = respDecoder.decode(encoded); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("OK")) == false {
		t.Fatal(errTestFailed)
	}

	// String with a special character.
	encoded = []byte("+OK\r+NO\r\n")

	if test, err = respDecoder.decode(encoded); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("OK\r+NO")) == false {
		t.Fatal(errTestFailed)
	}
}

func TestDecodeError(t *testing.T) {
	var test interface{}
	var encoded []byte
	var err error

	// Simple "Error Message" error
	encoded = []byte("-Error Message\r\n")

	if test, err = respDecoder.decode(encoded); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(error).Error()), []byte("Error Message")) == false {
		t.Fatal(errTestFailed)
	}
}

func TestDecodeInteger(t *testing.T) {
	var test interface{}
	var encoded []byte
	var err error

	// Positive integer.
	encoded = []byte(":123\r\n")

	if test, err = respDecoder.decode(encoded); err != nil {
		t.Fatal(err)
	}

	if test.(int) != 123 {
		t.Fatal(errTestFailed)
	}

	// Negative integer.
	encoded = []byte(":-123\r\n")

	if test, err = respDecoder.decode(encoded); err != nil {
		t.Fatal(err)
	}

	if test.(int) != -123 {
		t.Fatal(errTestFailed)
	}

	// Wrong formatting
	encoded = []byte(":-12.3\r\n")

	if test, err = respDecoder.decode(encoded); err == nil {
		t.Fatal(errErrorExpected)
	}

	if test != nil {
		t.Fatal(errTestFailed)
	}

	// Wrong formatting
	encoded = []byte(":-12a3\r\n")

	if test, err = respDecoder.decode(encoded); err == nil {
		t.Fatal(errErrorExpected)
	}

	if test != nil {
		t.Fatal(errTestFailed)
	}

}

func TestDecodeBulk(t *testing.T) {
	var test interface{}
	var err error

	// "foobar" string.
	if test, err = respDecoder.decode([]byte("$6\r\nfoobar\r\n")); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("foobar")) == false {
		t.Fatal(errTestFailed)
	}

	// "foo\r\nbar" string.
	if test, err = respDecoder.decode([]byte("$8\r\nfoo\r\nbar\r\n")); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("foo\r\nbar")) == false {
		t.Fatal(errTestFailed)
	}

	// An empty string.
	if test, err = respDecoder.decode([]byte("$0\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("")) == false {
		t.Fatal(errTestFailed)
	}

	// Nil.
	if test, err = respDecoder.decode([]byte("$-1\r\n")); err != nil {
		t.Fatal(err)
	}

	if test != nil {
		t.Fatal(errTestFailed)
	}

	// UTF-8 string.
	if test, err = respDecoder.decode([]byte("$3\r\n✓\r\n")); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte(test.(string)), []byte("✓")) == false {
		t.Fatal(errTestFailed)
	}

	// Invalid.
	if test, err = respDecoder.decode([]byte("$12\r\nSmall\r\n")); err == nil {
		t.Fatal(errErrorExpected)
	}

	if test != nil {
		t.Fatal(errTestFailed)
	}

}

func TestArrayDecode(t *testing.T) {
	var test interface{}
	var err error

	// Array with zero elements.
	if test, err = respDecoder.decode([]byte("*0\r\n")); err != nil {
		t.Fatal(err)
	}

	if len(test.([]interface{})) > 0 {
		t.Fatal(errTestFailed)
	}

	// Nil.
	if test, err = respDecoder.decode([]byte("*-1\r\n")); err != nil {
		t.Fatal(err)
	}

	if test != nil {
		t.Fatal(errTestFailed)
	}
}

func TestArrayDecodeTwoElements(t *testing.T) {
	var test interface{}
	var err error

	// Array with two elements.
	if test, err = respDecoder.decode([]byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")); err != nil {
		t.Fatal(err)
	}

	res := test.([]interface{})

	if len(res) != 2 {
		t.Fatal(errTestFailed)
	}

	if res[0].(string) != "foo" {
		t.Fatal(errTestFailed)
	}

	if res[1].(string) != "bar" {
		t.Fatal(errTestFailed)
	}
}

func TestArrayDecodeThreeIntegers(t *testing.T) {
	var test interface{}
	var err error

	// Array of three integers.
	if test, err = respDecoder.decode([]byte("*3\r\n:1\r\n:2\r\n:3\r\n")); err != nil {
		t.Fatal(err)
	}

	res := test.([]interface{})

	if len(res) != 3 {
		t.Fatal(errTestFailed)
	}

	if res[0].(int) != 1 {
		t.Fatal(errTestFailed)
	}

	if res[1].(int) != 2 {
		t.Fatal(errTestFailed)
	}

	if res[2].(int) != 3 {
		t.Fatal(errTestFailed)
	}
}

func TestArrayMixed(t *testing.T) {
	var test interface{}
	var err error

	// Array of four integers and one string.
	if test, err = respDecoder.decode([]byte("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n")); err != nil {
		t.Fatal(err)
	}

	res := test.([]interface{})

	if len(res) != 5 {
		t.Fatal(errTestFailed)
	}

	if res[0].(int) != 1 {
		t.Fatal(errTestFailed)
	}

	if res[1].(int) != 2 {
		t.Fatal(errTestFailed)
	}

	if res[2].(int) != 3 {
		t.Fatal(errTestFailed)
	}

	if res[3].(int) != 4 {
		t.Fatal(errTestFailed)
	}

	if res[4].(string) != "foobar" {
		t.Fatal(errTestFailed)
	}
}

func TestArrayNested(t *testing.T) {
	var test interface{}
	var err error

	// Array of two arrays.
	if test, err = respDecoder.decode([]byte("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n")); err != nil {
		t.Fatal(err)
	}

	res := test.([]interface{})

	if len(res) != 2 {
		t.Fatal(errTestFailed)
	}

	arr1 := res[0].([]interface{})
	arr2 := res[1].([]interface{})

	if arr1[0].(int) != 1 {
		t.Fatal(errTestFailed)
	}

	if arr1[1].(int) != 2 {
		t.Fatal(errTestFailed)
	}

	if arr1[2].(int) != 3 {
		t.Fatal(errTestFailed)
	}

	if arr2[0].(string) != "Foo" {
		t.Fatal(errTestFailed)
	}
	if arr2[1].(error).Error() != "Bar" {
		t.Fatal(errTestFailed)
	}

}

func TestArrayWithNil(t *testing.T) {
	var test interface{}
	var err error

	// Array of two arrays.
	if test, err = respDecoder.decode([]byte("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n")); err != nil {
		t.Fatal(err)
	}

	res := test.([]interface{})

	if len(res) != 3 {
		t.Fatal(errTestFailed)
	}

	if res[0].(string) != "foo" {
		t.Fatal(errTestFailed)
	}
	if res[1] != nil {
		t.Fatal(errTestFailed)
	}
	if res[2].(string) != "bar" {
		t.Fatal(errTestFailed)
	}

}
