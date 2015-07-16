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
func TestReaderWriterWithType(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBufferString("")
	r := BufReader{bufio.NewReader(buf)}
	w := BufWriter{bufio.NewWriter(buf)}
	tests := []struct {
		v interface{}
		t ProtoType
	}{
		{uint8(1), Int1},
		{uint16(1), Int2},
		{uint32(1), Int3},
		{uint32(1), Int4},
		{uint64(1), Int6},
		{uint64(1), Int8},
		{uint64(1), IntEnc},
		{"ABC", StrEof},
		{"DEF", StrNul},
		{"GHI", StrEnc},
		{[]byte{0xa, 0xb, 0xc}, StrEof},
		{[]byte{0xa, 0xb, 0xc}, StrNul},
		{[]byte{0xa, 0xb, 0xc}, StrEnc},
	}
	var i interface{}
	for _, t := range tests {
		w.Put(t.v, t.t)
		w.Flush()
		r.Get(&i, t.t)
		assert.EqualValues(t.v, i)
		assert.Equal(0, buf.Len())
	}
}
