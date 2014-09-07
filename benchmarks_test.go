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
	"encoding/json"
	"log"
	"testing"
)

var (
	benchmarkString       = "More than a woman."
	benchmarkBytes        = "More than a woman to me."
	benchmarkInteger      = 1234567890
	benchmarkArrayString  = []string{"fanny", "be", "tender", "with", "my", "love"}
	benchmarkArrayBytes   = [][]byte{[]byte("spicks"), []byte("and"), []byte("specks")}
	benchmarkArrayInteger = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	benchmarkArrayArray   = []interface{}{benchmarkArrayString, benchmarkArrayBytes, benchmarkArrayInteger}

	benchmarkRESPEncodedString       = mustRESPEncode(benchmarkString)
	benchmarkRESPEncodedBytes        = mustRESPEncode(benchmarkBytes)
	benchmarkRESPEncodedInteger      = mustRESPEncode(benchmarkInteger)
	benchmarkRESPEncodedArrayString  = mustRESPEncode(benchmarkArrayString)
	benchmarkRESPEncodedArrayBytes   = mustRESPEncode(benchmarkArrayBytes)
	benchmarkRESPEncodedArrayInteger = mustRESPEncode(benchmarkArrayInteger)
	benchmarkRESPEncodedArrayArray   = mustRESPEncode(benchmarkArrayArray)

	benchmarkJSONEncodedString       = mustJSONEncode(benchmarkString)
	benchmarkJSONEncodedBytes        = mustJSONEncode(benchmarkBytes)
	benchmarkJSONEncodedInteger      = mustJSONEncode(benchmarkInteger)
	benchmarkJSONEncodedArrayString  = mustJSONEncode(benchmarkArrayString)
	benchmarkJSONEncodedArrayBytes   = mustJSONEncode(benchmarkArrayBytes)
	benchmarkJSONEncodedArrayInteger = mustJSONEncode(benchmarkArrayInteger)
	benchmarkJSONEncodedArrayArray   = mustJSONEncode(benchmarkArrayArray)
)

func mustJSONEncode(i interface{}) []byte {
	b, e := json.Marshal(i)
	if e != nil {
		log.Printf("input: %v\n", i)
		panic(e)
	}
	return b
}

func mustRESPEncode(i interface{}) []byte {
	b, e := Marshal(i)
	if e != nil {
		log.Printf("input: %v\n", i)
		panic(e)
	}
	return b
}

func BenchmarkJSONMarshalString(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalString(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalBytes(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkBytes); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalBytes(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkBytes); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalInteger(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkInteger); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalInteger(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkInteger); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalArrayString(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkArrayString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalArrayString(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkArrayString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalArrayBytes(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkArrayBytes); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalArrayBytes(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkArrayBytes); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalArrayInteger(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkArrayInteger); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalArrayInteger(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkArrayInteger); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalArrayArray(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = json.Marshal(benchmarkArrayArray); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPMarshalArrayArray(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		if _, err = Marshal(benchmarkArrayArray); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalString(b *testing.B) {
	var err error
	var d string

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedString, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalString(b *testing.B) {
	var err error
	var d string

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedString, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalBytes(b *testing.B) {
	var err error
	var d []byte

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedBytes, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalBytes(b *testing.B) {
	var err error
	var d []byte

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedBytes, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalInteger(b *testing.B) {
	var err error
	var d int

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedInteger, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalInteger(b *testing.B) {
	var err error
	var d int

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedInteger, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalArrayString(b *testing.B) {
	var err error
	var d []string

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedArrayString, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalArrayString(b *testing.B) {
	var err error
	var d []string

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedArrayString, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalArrayBytes(b *testing.B) {
	var err error
	var d [][]byte

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedArrayBytes, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalArrayBytes(b *testing.B) {
	var err error
	var d [][]byte

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedArrayBytes, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalArrayInteger(b *testing.B) {
	var err error
	var d []int

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedArrayInteger, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalArrayInteger(b *testing.B) {
	var err error
	var d []int

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedArrayInteger, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalArrayArray(b *testing.B) {
	var err error
	var d []interface{}

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(benchmarkJSONEncodedArrayArray, &d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRESPUnmarshalArrayArray(b *testing.B) {
	var err error
	var d []interface{}

	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkRESPEncodedArrayArray, &d); err != nil {
			b.Fatal(err)
		}
	}
}
