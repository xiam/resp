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
)

var (
	// ErrInvalidInput is returned after any error decoding a message.
	ErrInvalidInput = errors.New(`resp: Invalid input`)

	// ErrMessageIsTooLarge is returned when a message attepts to create a buffer
	// that is considered too large.
	ErrMessageIsTooLarge = errors.New(`resp: Message is too large`)

	// ErrMissingMessageHeader is returned when the user attempts to encode a
	// message that has no header.
	ErrMissingMessageHeader = errors.New(`resp: Missing message header`)

	// ErrExpectingPointer is returned when a function expects a pointer
	// parameter.
	ErrExpectingPointer = errors.New(`resp: Expecting pointer value`)

	// ErrUnsupportedConversion is returned when the user attempts to unmarshal a
	// value into an incompatible destination type.
	ErrUnsupportedConversion = errors.New(`resp: Unsupported conversion: %s to %s`)

	// ErrMessageIsNil is returned when an user attempts to encode a nil message.
	ErrMessageIsNil = errors.New(`resp: Message is nil`)

	// ErrMissingReader is returned when the decoder needs to read additional
	// data from a writer that was not defined.
	ErrMissingReader = errors.New(`resp: Ran out of buffered data and a reader was not defined`)

	// ErrExpectingDestination is returned when a user attempts to unmarshal into
	// a nil value.
	ErrExpectingDestination = errors.New(`resp: Expecting a valid destination, but a nil value was provided`)
)
