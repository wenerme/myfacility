package proto
import (
	"bufio"
	"fmt"
	"errors"
	"reflect"
	"io/ioutil"
)

type Reader struct {
	*bufio.Reader
	Cap          Capability
	CharacterSet int
	Com          CommandType
}

type Writer struct {
	*bufio.Writer
	Cap          Capability
	CharacterSet int
	Com          CommandType
}

type ProtoType int
const (
	UndType ProtoType = iota
	Int1    // http://dev.mysql.com/doc/internals/en/integer.html
	Int2
	Int3
	Int4
	Int6
	Int8
	IntEnc  // int<lenenc>
	StrEof  // string<EOF>	    Protocol::RestOfPacketString
	StrNul  // string<NUL>	    Protocol::NulTerminatedString
	StrEnc  // string<lenenc>	Protocol::LengthEncodedString
	StrVar  // string<var/fix>	    Protocol::VariableLengthString:
)

func (r *Reader)SkipBytes(n int) *Reader {
	for i := 0; i < n; i ++ {
		_, err := r.ReadByte()
		if err != nil {panic(err) }
	}
	return t
}

func (r *Reader)Get(values... interface{}) *Reader {
	n := len(values)
	for i := 0; i < n; i ++ {
		v := values[i]
		t := UndType
		if i < n-1 {
			if ty, ok := values[i+1].(ProtoType); ok {
				t = ty
				i ++
			}
		}

		if t == UndType {
			switch v.(type){
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
				panic(errors.New(fmt.Sprintf("Can not get type %v", t)))
			}
		}else if t == StrVar {
			// need specified a size
			if i < n-1 {
				if size, ok := values[i+1].(int); ok {
					i ++
					buf := make([]byte, size)
					_, err := r.Read(buf)
					if err != nil {panic(err)}
					switch v.(type){
						case *string:
						*v = string(buf)
						case *[]byte:
						*v = buf
					}
					continue
				}else {
					panic(errors.New("Type StrVar need a int type size"))
				}
			}else {
				panic(errors.New("Type StrVar need a size"))
			}
		}
		_, err := r.GetType(v, t)
		if err != nil {panic(err)}
	}
	return r
}
func (r *Reader)GetType(v interface{}, t ProtoType) (n int, err error) {
	val := reflect.ValueOf(v)
	var buf []byte
	TYPE_SWITCH:
	switch t{
	case Int1: {
		b, e := r.ReadByte();
		if e != nil {err = e; break }
		n = 1
		val.SetUint(uint64(b))
	}
	case Int2:
		buf = make([]byte, 2)
		n, err= r.Read(buf)
		if err != nil {break }
		val.SetUint(uint64(buf[0]) | uint64(buf[1]) << 8)
	case Int3:
		buf = make([]byte, 3)
		n, err= r.Read(buf)
		if err != nil {break }
		val.SetUint(uint64(buf[0]) | uint64(buf[1]) << 8 | uint64(buf[2]) << 16)
	case Int4:
		buf = make([]byte, 4)
		n, err= r.Read(buf)
		if err != nil {break }
		val.SetUint(uint64(buf[0]) | uint64(buf[1]) << 8 | uint64(buf[2]) << 16 | uint64(buf[3]) << 24)
	case Int6:
		buf = make([]byte, 6)
		n, err= r.Read(buf)
		if err != nil {break }
		val.SetUint(uint64(buf[0]) | uint64(buf[1]) << 8 | uint64(buf[2]) << 16 | uint64(buf[3]) << 24 | uint64(buf[4]) << 32 | uint64(buf[5]) << 40)
	case Int8:
		buf = make([]byte, 8)
		n, err= r.Read(buf)
		if err != nil {break }
		val.SetUint(uint64(buf[0]) | uint64(buf[1]) << 8 | uint64(buf[2]) << 16 | uint64(buf[3]) << 24 | uint64(buf[4]) << 32 | uint64(buf[5]) << 40| uint64(buf[6]) << 48| uint64(buf[7]) << 56)
	case IntEnc:
		i, e := r.ReadByte()
		if e != nil {err = e; break}
		if i < 251 {
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
		_, err=r.GetType(&size, IntEnc)
		if err != nil {break}
		buf = make([]byte, size)
		n, err= r.Read(buf)
		if err != nil {break }
		val.SetString(string(buf))
	case StrEof:
		bytes, e := ioutil.ReadAll(r)
		if e != nil {err = e; break}
		n = len(bytes)
		val.SetString(string(bytes))
	case StrNul:
		bytes, e := r.ReadBytes(0)
		if e != nil {err = e; break}
		n = len(bytes)
		val.SetString(string(bytes))
	default:
		err = errors.New(fmt.Sprintf("Can not get type %v", t))
	}
	return
}




func (w *Writer)Put(values...interface{}) {
	n := len(values)
	for i := 0; i < n; i ++ {
		v := values[i]
		t := UndType
		if i < n-1 {
			if ty, ok := values[i+1].(ProtoType); ok {
				t = ty
			}
		}

		if t == UndType {
			switch v.(type){
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
				panic(errors.New(fmt.Sprintf("Can not get type %v", t)))
			}
		}
		_, err := w.PutType(v, t)
		if err != nil {panic(err)}
	}
	return w
}
func (w *Writer)PutZero(n int) *Writer {
	for i := 0; i < n; i ++ {
		err := w.WriteByte(0)
		if err != nil { panic(err) }
	}
	return w
}
func (w *Writer)PutType(v interface{}, t ProtoType) (n int, err error) {
	val := reflect.ValueOf(v)
	var buf []byte
	TYPE_SWITCH:
	switch t{
	case Int1: {
		err = w.WriteByte(byte(val.Uint()))
	}
	case Int2:
		buf = make([]byte, 2)
		u := val.Uint()
		buf[0], buf[1] = byte(u), byte(u >> 8)
		n, err= w.Write(buf)
	case Int3:
		buf = make([]byte, 3)
		u := val.Uint()
		buf[0], buf[1], buf[2] = byte(u), byte(u >> 8), byte(u >> 16)
		n, err= w.Write(buf)
	case Int4:
		buf = make([]byte, 4)
		u := val.Uint()
		buf[0], buf[1], buf[2], buf[3] = byte(u), byte(u >> 8), byte(u >> 16), byte(u >> 24)
		n, err= w.Write(buf)
	case Int6:
		buf = make([]byte, 6)
		u := val.Uint()
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5] = byte(u), byte(u >> 8), byte(u >> 16), byte(u >> 24), byte(u >> 32), byte(u >> 40)
		n, err= w.Write(buf)
	case Int8:
		buf = make([]byte, 8)
		u := val.Uint()
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7 ] = byte(u), byte(u >> 8), byte(u >> 16), byte(u >> 24), byte(u >> 32), byte(u >> 40), byte(u >> 48), byte(u >> 56)
		n, err= w.Write(buf)
	case IntEnc:
		i := val.Uint()
		if i < 251 {
			n = 1
			w.Write(byte(i))
			break
		}
		switch  {
		case i < 0xffff:
			t = Int2
		case i < 0xffffff:
			t = Int3
		default:
			t = Int8
		}
		goto TYPE_SWITCH
	case StrEnc:
	case StrEof:
	case StrNul:
	default:
		err = errors.New(fmt.Sprintf("Can not put type %v", t))
	}
	return
}