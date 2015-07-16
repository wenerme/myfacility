package proto

import (
	"fmt"
	"github.com/spacemonkeygo/errors"
	"io"
	"reflect"
)

/*
Reader Writer use reflect to get parameter type, map go type to protocol type
Also can specify protocol type explicit
*/

// A reader used to read packet
type Reader interface {
	io.Reader
	// Format
	// Get(&value,&value,&value...)
	// Get(&value,ProtoType)
	// Get(&value,StrVar,n)
	// Get(&value,StrVar,&n)
	// Get(&value,Int,n)
	// Get(&value,Int,&n)
	// Get(n, IgnoreByte)
	// Get(&value,reflect.Kind)
	Get(...interface{})
	More() bool
	Peek(int) ([]byte, error)
	PeekByte() (byte, error)
}

// A writer used to write packet
type Writer interface {
	io.Writer
	// Format
	// Pet(&value,&value,&value...)
	// Pet(value,value,value...)
	// Pet(&value,ProtoType)
	// Pet(&value,StrVar,n)
	// Pet(&value,StrVar,&n)
	// Pet(&value,Int,n)
	// Pet(&value,Int,&n)
	// Pet(n, IgnoreByte)
	Put(...interface{})
	PutZero(int)
}

// General protocol types
// A unit for read write packet content
type ProtoType int

const (
	UndType ProtoType = iota
	Int               // Must specify 1 2 3 4 5 6
	Int1              // http://dev.mysql.com/doc/internals/en/integer.html
	Int2
	Int3
	Int4
	Int6
	Int8
	IntEnc // int<lenenc>
	StrEof // string<EOF>	    Protocol::RestOfPacketString
	StrNul // string<NUL>	    Protocol::NulTerminatedString
	StrEnc // string<lenenc>	Protocol::LengthEncodedString
	StrVar // string<var/fix>	    Protocol::VariableLengthString:
	// Skip n byte for Get
	// Write n byte zero for Put
	IgnoreByte
)

type readablePack interface {
	Read(Reader)
}

func checkInt(v interface{}) (i int, ok bool) {
	ok = true
	switch v.(type) {
	case int:
		i = v.(int)
	case uint:
		i = int(v.(uint))
	case *int:
		i = *v.(*int)
	case *uint:
		i = int(*v.(*uint))
	case uint8:
		i = int(v.(uint8))
	case *uint8:
		i = int(*v.(*uint8))
	case *uint16:
		i = int(*v.(*uint16))
	case *uint32:
		i = int(*v.(*uint32))
	default:
		ok = false
	}
	return
}
func newKind(k reflect.Kind) interface{} {
	switch k {
	case reflect.Int:
		var i int
		return &i
	case reflect.Int8:
		var i int8
		return &i
	case reflect.Int16:
		var i int16
		return &i
	case reflect.Int32:
		var i int32
		return &i
	case reflect.Int64:
		var i int64
		return &i
	case reflect.Uint:
		var i uint
		return &i
	case reflect.Uint8:
		var i uint8
		return &i
	case reflect.Uint16:
		var i uint16
		return &i
	case reflect.Uint32:
		var i uint32
		return &i
	case reflect.Uint64:
		var i uint64
		return &i
	}
	panic(errors.New(fmt.Sprintf("Can not new kind %v", k)))
}
