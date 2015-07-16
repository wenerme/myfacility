package proto

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestReaderWriterWithKind(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBufferString("")
	r := BufReader{bufio.NewReader(buf)}
	w := BufWriter{bufio.NewWriter(buf)}
	tests := []struct {
		data []byte
		v    interface{}
		k    reflect.Kind
	}{
		{[]byte{01, 00, 00, 00, 0, 0, 0, 0}, 1, reflect.Uint64},
		{[]byte{01, 00, 00, 00}, 1, reflect.Uint32},
		{[]byte{01, 00}, 1, reflect.Uint16},
		{[]byte{01}, 1, reflect.Uint8},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, -1, reflect.Int64},
		{[]byte{0xff, 0xff, 0xff, 0xff}, -1, reflect.Int32},
		{[]byte{0xff, 0xff}, -1, reflect.Int16},
		{[]byte{0xff}, -1, reflect.Int8},
	}

	for _, t := range tests {
		buf.Write(t.data)
		var i interface{}
		r.Get(&i, t.k)
		assert.EqualValues(i, t.v)
		w.Put(i)
		w.Flush()
		r.Get(&i, t.k)
		assert.EqualValues(i, t.v)
	}
}
