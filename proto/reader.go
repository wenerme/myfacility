package proto

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/spacemonkeygo/errors"
	"io"
	"io/ioutil"
	"reflect"
)

type BufReader struct {
	*bufio.Reader
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

// Handle integer case
// Better performance
func (r *BufReader) readKind(v interface{}, k reflect.Kind) (err error) {
	//	var buf []byte
	switch k {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = binary.Read(r, binary.LittleEndian, v)
	case reflect.Uint:
		r.getByType(v, IntEnc)
	case reflect.Int:
		var i uint
		r.getByType(&i, IntEnc)
		*v.(*int) = int(i)
	default:
		panic(fmt.Sprintf("Can not get %T(nil)", v))
	}
	return
}

func (r *BufReader) getByType(v interface{}, t ProtoType) (n int, err error) {
	val := reflect.ValueOf(v).Elem()
	var buf []byte
TYPE_SWITCH:
	switch t {
	case Int1:
		b, e := r.ReadByte()
		if e != nil {
			err = e
			break
		}
		n = 1
		val.SetUint(uint64(b))
	case Int2:
		buf = make([]byte, 2)
		n, err = io.ReadFull(r, buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8)
	case Int3:
		buf = make([]byte, 3)
		n, err = io.ReadFull(r, buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16)
	case Int4:
		buf = make([]byte, 4)
		n, err = io.ReadFull(r, buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24)
	case Int6:
		buf = make([]byte, 6)
		n, err = io.ReadFull(r, buf)
		if err != nil {
			break
		}
		val.SetUint(uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24 | uint64(buf[4])<<32 | uint64(buf[5])<<40)
	case Int8:
		buf = make([]byte, 8)
		n, err = io.ReadFull(r, buf)
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
		_, err = r.getByType(&size, IntEnc)
		if err != nil {
			break
		}
		bytes := make([]byte, size)
		n, err = io.ReadFull(r, bytes)
		if err != nil {
			break
		}
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
func (r *BufReader) Get(values ...interface{}) {
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
						_, err := io.CopyN(ioutil.Discard, r, int64(n))
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
						buf := make([]byte, n)
						_, err := r.Read(buf)
						if err != nil {
							panic(err)
						}
						switch v.(type) {
						case *string:
							*v.(*string) = string(buf)
						case *[]byte:
							*v.(*[]byte) = buf
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
					}
					panic(errors.New("Type Int need a size"))
				}

				r.getByType(v, pt)
				continue
			} else if try, ok := values[i+1].(reflect.Kind); ok {
				// With kind
				kind := try
				i++

				val := reflect.ValueOf(v)
				if kind == val.Kind() {
					r.readKind(v, kind)
				} else {
					neo := newKind(kind)
					r.readKind(neo, kind)
					val.Set(reflect.ValueOf(neo).Elem())
				}
				continue
			}
		}

		// Normal
		switch v.(type) {
		case *string, *[]byte:
			r.getByType(v, StrEnc)
		case *uint:
			r.getByType(v, IntEnc)
		case *int:
			var i uint
			r.getByType(&i, IntEnc)
			*v.(*int) = int(i)
		default:
			err := binary.Read(r, binary.LittleEndian, v)
			if err != nil {
				panic(err)
			}
		}
	}
}
