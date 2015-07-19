package proto

import (
	"encoding/binary"
	"fmt"
	"github.com/spacemonkeygo/errors"
	"io"
	"reflect"
)

type writer struct {
	io.Writer
}

func (w *writer) Put(values ...interface{}) {
	argc := len(values)
	for i := 0; i < argc; i++ {
		v := values[i]
		// Detect next type parameter
		if i < argc-1 {
			if try, ok := values[i+1].(ProtoType); ok {
				// With proto type
				pt := try
				i++
				switch pt {
				case IgnoreByte:
					if n, ok := checkInt(v); ok {
						_, err := io.CopyN(w, ZeroReader, int64(n))
						if err != nil {
							panic(err)
						}
						continue
					}
					panic(errors.New("Ignore byte need a int size"))
				case StrVar:
					if i >= argc-1 {
						panic(errors.New("Type StrVar need a size"))
					}
					if n, ok := checkInt(values[i+1]); ok {
						i++
						switch v.(type) {
						case string:
							w.Write([]byte(v.(string))[0:n])
						case []byte:
							w.Write(v.([]byte)[0:n])
						case *string:
							w.Write([]byte(*(v.(*string)))[0:n])
						case *[]byte:
							w.Write((*(v.(*[]byte)))[0:n])
						default:
							panic(errors.New(fmt.Sprintf("Can not handle type StrVar %T(%v)", v, v)))
						}
						continue
					}
					panic(errors.New("Type StrVar need a int type size"))
				case Int:
					if i >= argc-1 {
						panic(errors.New("Type Int need a size"))
					}
					pt = UndType
					if n, ok := checkInt(values[i+1]); ok {
						i++
						switch n {
						case 1:
							pt = Int1
						case 2:
							pt = Int2
						case 3:
							pt = Int3
						case 4:
							pt = Int4
						case 6:
							pt = Int6
						case 8:
							pt = Int8
						default:
							panic(errors.New(fmt.Sprintf("Unsupport Int size %v", n)))
						}
					} else {
						panic(errors.New("Type Int need a size"))
					}
				}
				w.mustPutByType(v, pt)
				continue
			} else if try, ok := values[i+1].(reflect.Kind); ok {
				// With kind
				_ = try
				i++
				panic(errors.New("Put don't support kind"))
			}
		}

		// Normal
		switch v.(type) {
		case *string, *[]byte, string, []byte:
			w.mustPutByType(v, StrEnc)
		case *uint, uint:
			w.mustPutByType(v, IntEnc)
		case *int:
			i := uint(*v.(*int))
			w.mustPutByType(i, IntEnc)
		case int:
			i := uint(v.(int))
			w.mustPutByType(i, IntEnc)
		default:
			err := binary.Write(w, binary.LittleEndian, v)
			if err != nil {
				panic(err)
			}
		}
	}
}
func (w *writer) WriteByte(b byte) (err error) {
	if bw, ok := w.Writer.(io.ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		_, err = w.Write([]byte{b})
	}
	return
}
func (w *writer) mustPutByType(v interface{}, t ProtoType) (n int) {
	n, err := w.putByType(v, t)
	if err != nil {
		panic(err)
	}
	return n
}
func (w *writer) putByType(v interface{}, t ProtoType) (n int, err error) {
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
		bytes := mustBytes(v)
		if n, err = w.putByType(uint64(len(bytes)), IntEnc); err == nil {
			write := n
			n, err = w.Write(bytes)
			n += write
		}
	case StrEof:
		n, err = w.Write(mustBytes(v))
	case StrNul:
		n, err = w.Write(mustBytes(v))
		if err == nil {
			err = w.WriteByte(0)
			n++
		}
	default:
		goto CAN_NOT_PUT
	}
	return
CAN_NOT_PUT:
	err = errors.New(fmt.Sprintf("Can not put type %v", t))
	return
}

func mustBytes(v interface{}) []byte {
	switch v.(type) {
	case string:
		return []byte(v.(string))
	case []byte:
		return v.([]byte)
	case *string:
		return []byte(*v.(*string))
	case *[]byte:
		return *v.(*[]byte)
	default:
		panic(errors.New(fmt.Sprintf("Can not convert %T to []byte", v)))
	}
}
