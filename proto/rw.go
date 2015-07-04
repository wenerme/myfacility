package proto

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

// A reader used to read packet
type Reader interface {
	Get(...interface{})
	GetType(interface{}, ProtoType) (int, error)
	SkipBytes(int)
	More() bool
	HasCap(Capability) bool
	Peek(int) ([]byte, error)
	PeekByte() (byte, error)
	Com() Command
	SetCom(Command)
}

// A writer used to write packet
type Writer interface {
	Put(...interface{})
	PutType(interface{}, ProtoType) (int, error)
	PutZero(int)
	HasCap(Capability) bool
	Com() Command
	SetCom(Command)
}

type BufReader struct {
	*bufio.Reader
	cap Capability
	com Command
}

type BufWriter struct {
	*bufio.Writer
	cap Capability
	com Command
}

// General protocol types
// A unit for read write packet content
type ProtoType int

const (
	UndType ProtoType = iota
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
)

func (r *BufReader) SetCom(com Command) {
	r.com = com
}
func (r *BufReader) Com() Command {
	return r.com
}

func (r *BufWriter) SetCom(com Command) {
	r.com = com
}

func (r *BufWriter) Com() Command {
	return r.com
}

func (r *BufWriter) SetCap(cap Capability) {
	r.cap = cap
}
func (r *BufReader) More() bool {
	_, err := r.Peek(1)
	if err == nil {
		return true
	} else if err == io.EOF {
		return false
	} else {
		panic(err)
	}
}

func (r *BufReader) PeekByte() (b byte, err error) {
	var tmp []byte
	tmp, err = r.Peek(1)
	if err == nil {
		b = tmp[0]
	}
	return
}
func (r *BufReader) SkipBytes(n int) {
	for i := 0; i < n; i++ {
		_, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
	}
}

func (r *BufReader) Get(values ...interface{}) {
	n := len(values)
	for i := 0; i < n; i++ {
		v := values[i]
		t := UndType

		if v == nil {
			panic(fmt.Sprintf("Can not get %T(nil)", v))
		}
		val := reflect.ValueOf(v)
		if val.CanAddr() {
			panic(fmt.Sprintf("Must use a addressable value instead of %T(%v)", v, v))
		}

		if i < n-1 {
			if ty, ok := values[i+1].(ProtoType); ok {
				t = ty
				i++
			}
		}

		if t == UndType {
			// For type alias
			switch val.Elem().Kind() {
			case reflect.Uint:
				t = IntEnc
			case reflect.Uint8:
				t = Int1
			case reflect.Uint16:
				t = Int2
			case reflect.Uint32:
				t = Int4
			case reflect.Uint64:
				t = Int8
			}
		}

		if t == UndType {
			switch v.(type) {
			case *uint8:
				t = Int1
			case *uint16:
				t = Int2
			case *uint32:
				t = Int4
			case *uint64:
				t = Int8
			case *uint:
				t = IntEnc
			case *string, *[]byte:
				t = StrEnc
			default:
				panic(errors.New(fmt.Sprintf("Can not get type of %T", v)))
			}
		} else if t == StrVar {
			// need specified a size
			if i < n-1 {
				if size, ok := values[i+1].(int); ok {
					i++
					buf := make([]byte, size)
					_, err := r.Read(buf)
					if err != nil {
						panic(err)
					}
					switch v.(type) {
					case *string:
						*v.(*string) = string(buf)
					case *[]byte:
						*v.(*[]byte) = buf
					}
					continue
				} else {
					panic(errors.New("Type StrVar need a int type size"))
				}
			} else {
				panic(errors.New("Type StrVar need a size"))
			}
		}
		_, err := r.GetType(v, t)
		if err != nil {
			panic(err)
		}
	}
}
func (r *BufReader) GetType(v interface{}, t ProtoType) (n int, err error) {
	val := reflect.ValueOf(v)
	if !val.CanSet() {
		if val = val.Elem(); !val.CanSet() {
			return 0, errors.New(fmt.Sprintf("Must use a addressable value instead of %T(%v)", v, v))
		}
	}
	var buf []byte
TYPE_SWITCH:
	switch t {
	case Int1:
		{
			b, e := r.ReadByte()
			if e != nil {
				err = e
				break
			}
			n = 1
			val.SetUint(uint64(b))
		}
	case Int2:
		buf = make([]byte, 2)
		n, err = r.Read(buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8)
	case Int3:
		buf = make([]byte, 3)
		n, err = r.Read(buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16)
	case Int4:
		buf = make([]byte, 4)
		n, err = r.Read(buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24)
	case Int6:
		buf = make([]byte, 6)
		n, err = r.Read(buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24 | uint64(buf[4])<<32 | uint64(buf[5])<<40)
	case Int8:
		buf = make([]byte, 8)
		n, err = r.Read(buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24 | uint64(buf[4])<<32 | uint64(buf[5])<<40 | uint64(buf[6])<<48 | uint64(buf[7])<<56)
	case IntEnc:
		i, e := r.ReadByte()
		if e != nil {
			err = e
			break
		}
		if i <= 251 {
			n = 1
			val.SetUint(uint64(i))
			break
		}
		switch i {
		case 252:
			t = Int2
		case 253:
			t = Int3
		case 254:
			t = Int8
		}
		goto TYPE_SWITCH
	case StrEnc:
		var size uint32
		_, err = r.GetType(&size, IntEnc)
		if err != nil {
			break
		}
		bytes := make([]byte, size)
		n, err = r.Read(bytes)
		if err != nil {
			break
		}
		// How about use val.Set
		switch v.(type) {
		case *string:
			*v.(*string) = string(bytes)
		case *[]byte:
			*v.(*[]byte) = bytes
		default:
			goto CAN_NOT_GET
		}
	case StrEof:
		bytes, e := ioutil.ReadAll(r)
		if e != nil {
			err = e
			break
		}
		n = len(bytes)
		switch v.(type) {
		case *string:
			*v.(*string) = string(bytes)
		case *[]byte:
			*v.(*[]byte) = bytes
		default:
			goto CAN_NOT_GET
		}
	case StrNul:
		bytes, e := r.ReadBytes(0)
		if e != nil {
			err = e
			break
		}
		n = len(bytes)
		bytes = bytes[:n-1] // drop the nul
		switch v.(type) {
		case *string:
			*v.(*string) = string(bytes)
		case *[]byte:
			*v.(*[]byte) = bytes
		default:
			goto CAN_NOT_GET
		}

	default:
		goto CAN_NOT_GET
	}
	return
CAN_NOT_GET:
	err = errors.New(fmt.Sprintf("Can not get type %v", t))
	return
}

func (w *BufWriter) Put(values ...interface{}) {
	n := len(values)
	for i := 0; i < n; i++ {
		v := values[i]
		t := UndType
		if i < n-1 {
			if ty, ok := values[i+1].(ProtoType); ok {
				t = ty
				i++
			}
		}

		if v == nil {
			panic(fmt.Sprintf("Can not put %T(nil)", v))
		}
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			v = val.Elem().Interface()
		}

		if t == UndType {
			kind := val.Kind()
			if kind == reflect.Ptr {
				kind = val.Elem().Kind()
			}
			// For type alias
			switch kind {
			case reflect.Uint:
				t = IntEnc
			case reflect.Uint8:
				t = Int1
			case reflect.Uint16:
				t = Int2
			case reflect.Uint32:
				t = Int4
			case reflect.Uint64:
				t = Int8
			}
		}

		if t == UndType {
			switch v.(type) {
			case uint8:
				t = Int1
			case uint16:
				t = Int2
			case uint32:
				t = Int4
			case uint64:
				t = Int8
			case uint:
				t = IntEnc
			case string, []byte:
				t = StrEnc
			default:
				panic(errors.New(fmt.Sprintf("Can not get type of %T", v)))
			}
		} else if t == StrVar {
			// need specified a size
			if i < n-1 {
				if size, ok := values[i+1].(int); ok {
					i++
					switch v.(type) {
					case string:
						w.Write([]byte(v.(string))[0:size])
					case []byte:
						w.Write(v.([]byte)[0:size])
					}
					continue
				} else {
					panic(errors.New("Type StrVar need a int type size"))
				}
			} else {
				panic(errors.New("Type StrVar need a size"))
			}
		}
		_, err := w.PutType(v, t)
		if err != nil {
			panic(err)
		}
	}
}
func (w *BufWriter) PutZero(n int) {
	for i := 0; i < n; i++ {
		err := w.WriteByte(0)
		if err != nil {
			panic(err)
		}
	}
}
func (w *BufWriter) PutType(v interface{}, t ProtoType) (n int, err error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	var buf []byte
TYPE_SWITCH:
	if err != nil {
		return
	}
	switch t {
	case Int1:
		err = w.WriteByte(byte(val.Uint()))
	case Int2:
		buf = make([]byte, 2)
		u := val.Uint()
		buf[0], buf[1] = byte(u), byte(u>>8)
		n, err = w.Write(buf)
	case Int3:
		buf = make([]byte, 3)
		u := val.Uint()
		buf[0], buf[1], buf[2] = byte(u), byte(u>>8), byte(u>>16)
		n, err = w.Write(buf)
	case Int4:
		buf = make([]byte, 4)
		u := val.Uint()
		buf[0], buf[1], buf[2], buf[3] = byte(u), byte(u>>8), byte(u>>16), byte(u>>24)
		n, err = w.Write(buf)
	case Int6:
		buf = make([]byte, 6)
		u := val.Uint()
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5] = byte(u), byte(u>>8), byte(u>>16), byte(u>>24), byte(u>>32), byte(u>>40)
		n, err = w.Write(buf)
	case Int8:
		buf = make([]byte, 8)
		u := val.Uint()
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7] = byte(u), byte(u>>8), byte(u>>16), byte(u>>24), byte(u>>32), byte(u>>40), byte(u>>48), byte(u>>56)
		n, err = w.Write(buf)
	case IntEnc:
		i := val.Uint()
		if i < 251 {
			n = 1
			w.WriteByte(byte(i))
			break
		}
		n++
		switch {
		case i <= 0xffff:
			err = w.WriteByte(252)
			t = Int2
		case i <= 0xffffff:
			err = w.WriteByte(253)
			t = Int3
		default:
			err = w.WriteByte(254)
			t = Int8
		}
		goto TYPE_SWITCH
	case StrEnc:
		var bytes []byte
		switch v.(type) {
		case string:
			bytes = []byte(v.(string))
		case []byte:
			bytes = v.([]byte)
		default:
			goto CAN_NOT_PUT
		}
		if n, err = w.PutType(uint64(len(bytes)), IntEnc); err == nil {
			writeed := n
			n, err = w.Write(bytes)
			n += writeed
		}
	case StrEof, StrVar:
		switch v.(type) {
		case string:
			n, err = w.Write([]byte(v.(string)))
		case []byte:
			n, err = w.Write(v.([]byte))
		default:
			goto CAN_NOT_PUT
		}
	case StrNul:
		switch v.(type) {
		case string:
			n, err = w.Write([]byte(v.(string)))
		case []byte:
			n, err = w.Write(v.([]byte))
		default:
			goto CAN_NOT_PUT
		}
		if err == nil {
			err = w.WriteByte(0)
			n++
		}
	default:
	}
	return
CAN_NOT_PUT:
	err = errors.New(fmt.Sprintf("Can not put type %v", t))
	return
}

func (r *BufWriter) Cap() Capability {
	return r.cap
}
func (r *BufWriter) HasCap(cap Capability) bool {
	return r.cap.Has(cap)
}
func (r *BufReader) SetCap(cap Capability) {
	r.cap = cap
}
func (r *BufReader) Cap() Capability {
	return r.cap
}
func (r *BufReader) HasCap(cap Capability) bool {
	return r.cap.Has(cap)
}
