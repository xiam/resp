package resp

import (
	"bytes"
	"errors"
	"testing"
)

var respDecoder = decoder{}

var errTestFailed = errors.New("Test failed.")

func TestReadLine(t *testing.T) {
	var test []byte
	var err error

	encodedString := []byte("+OK\r\n")

	if test, err = respDecoder.readLine(encodedString); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal([]byte("+OK"), test) == false {
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
