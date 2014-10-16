# github.com/xiam/resp

The `resp` package provides methods for encoding and decoding data using the
[RESP][1] format.

## Installation

Use `go get` to instal or upgrade (`-u`) the `resp` package:

```
go get -u github.com/xiam/resp
```

## Usage

### Encoding

`resp` provides a `Marshal()` function that creates a RESP representation of a
given value.

```go
func Marshal(v interface{}) ([]byte, error)
```

The following example converts the input string `"Foo"` into its [RESP][1]
representation: `$3\r\nFoo\r\n`.

```go
buf, err = resp.Marshal("Foo") // RESP: $3\r\nFoo\r\n

fmt.Printf("buf: %s\n", string(buf))
```

`resp` also provides an `resp.Encoder` type that you can use to write the
encoded message to the given `io.Writer`.

```go
w = bytes.NewBuffer(nil)

resp.NewEncoder(w)

err = w.Encode("Hello World!")

fmt.Printf("RESP: %s\n", w.Bytes())
```

### Decoding

`resp` also provides an `Unmarshal()` function that takes a RESP message and
creates a Go value with it.

```go
func Unmarshal(data []byte, v interface{}) error
```

Let's take the binary safe RESP encoding of the "Foo" string (`$3\r\nFoo\r\n`),
this should be decoded into the string `"Foo"`.

```go
var dest string

err = resp.Unmarshal([]byte("$3\r\nFoo\r\n"), &dest)
```

If you want to decode a stream, you can use the `resp.Decoder` type providing
an `io.Reader`:

```go
var s string

r = bytes.NewBuffer([]byte("$3\r\nFoo\r\n"))

d = resp.NewDecoder(r)

err = d.Decode(&s)
```

## License

> Copyright (c) 2014 José Carlos Nieto, https://menteslibres.net/xiam
>
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[1]: http://redis.io/topics/protocol
