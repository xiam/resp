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

	if test, err = respDecoder.readLine([]byte("+OK\r\n")); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte("+OK"), test) == false {
		t.Fatal(errTestFailed)
	}

	if test, err = respDecoder.readLine([]byte("+OK")); err == nil {
		t.Fatal(errErrorExpected)
	}

	if test != nil {
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
