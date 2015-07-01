package proto
import (
	"errors"
	"encoding/binary"
)



// A ByteOrder specifies how to convert byte sequences into
// 16-, 32-, or 64-bit unsigned integers.
type PackByteOrder interface {
	binary.ByteOrder
	Uint24([]byte) uint32
	Uint40([]byte) uint64

	PutUint24([]byte, uint32)
	PutUint40([]byte, uint64)
}


type ExtendedByteOrder struct {
	PackByteOrder
}
func (this ExtendedByteOrder)MustUint(n int, b []byte) (uint64) {
	r, err := this.Uint(n, b)
	if err !=nil {
		panic(err)
	}
	return r
}
func (this ExtendedByteOrder)Uint(n int, b []byte) (uint64, error) {
	switch n{
	case 1:
		return uint64(b[0]), nil
	case 2:
		return uint64(this.Uint16(b)), nil
	case 3:
		return uint64(this.Uint24(b)), nil
	case 4:
		return uint64(this.Uint32(b)), nil
	case 6:
		return this.Uint40(b), nil
	case 8:
		return this.Uint64(b), nil
	default:
		return 0, errors.New("Can not read bytes "+string(n)+" int")
	}
}
func (this ExtendedByteOrder)MustPutUint(n int, b []byte, v uint64) {
	err := this.PutUint(n, b, v)
	if err !=nil {
		panic(err)
	}
}
func (this ExtendedByteOrder)PutUint(n int, b []byte, v uint64) (error) {
	switch n{
	case 1:
		b[0] = byte(v)
	case 2:
		this.PutUint16(b, uint16(v))
	case 3:
		this.PutUint24(b, uint32(v))
	case 4:
		this.PutUint32(b, uint32(v))
	case 6:
		this.PutUint40(b, v)
	case 8:
		this.PutUint64(b, v)
	default:
		return errors.New("Can not put bytes "+string(n)+" int")
	}
	return nil
}
type littleEndian struct {}

// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian


func (littleEndian) Uint16(b []byte) uint16 { return uint16(b[0]) | uint16(b[1])<<8 }

func (littleEndian) PutUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func (littleEndian) Uint24(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
func (littleEndian) Uint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func (littleEndian) PutUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}
func (littleEndian) PutUint24(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func (littleEndian) Uint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
	uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}
func (littleEndian) Uint40(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
	uint64(b[4])<<32
}

func (littleEndian) PutUint40(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
}
func (littleEndian) PutUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func (littleEndian) String() string { return "LittleEndian" }

func (littleEndian) GoString() string { return "binary.LittleEndian" }

type bigEndian struct {}

func (bigEndian) Uint16(b []byte) uint16 { return uint16(b[1]) | uint16(b[0])<<8 }

func (bigEndian) PutUint16(b []byte, v uint16) {
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func (bigEndian) Uint24(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func (bigEndian) Uint32(b []byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func (bigEndian) PutUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
func (bigEndian) PutUint32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func (bigEndian) Uint40(b []byte) uint64 {
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 |
	uint64(b[0])<<32
}

func (bigEndian) Uint64(b []byte) uint64 {
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
	uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func (bigEndian) PutUint40(b []byte, v uint64) {
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
}


func (bigEndian) PutUint64(b []byte, v uint64) {
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

func (bigEndian) String() string { return "BigEndian" }

func (bigEndian) GoString() string { return "binary.BigEndian" }

