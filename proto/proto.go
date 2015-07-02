package proto
import (
	"bufio"
	"encoding/binary"
	"io"
)

type (
	Pack interface {
		Read(Reader)
		Write(Writer)
	}
)

func ReadPacketTo(r *bufio.Reader, w io.Writer) (seq uint8, size int, err error) {
	var l uint32
	err = binary.Read(r, binary.LittleEndian, &l)
	if err != nil {return }
	seq = uint8(l >> 24)
	l = l << 8 >> 8
	size = int(l)
	var buf []byte = make([]byte, 16)
	for i := 0; i < size; {
		n, e := r.Read(buf)
		if n < 16 {
			w.Write(buf[0:n])
			i+=n
		}else {
			_, err = w.Write(buf)
			if e != nil {return }
			i += 16
		}

		if e == io.EOF { return }
		if e != nil {err = e; return }


	}
	return
}

func ReadCompressedPacket(r *bufio.Reader, w *bufio.Writer) (seq uint8, err error) {

	return
}


func WritePacket(seq uint8, r *bufio.Reader, w *bufio.Writer) (n int, err error) {

	return
}
func WriteCompressedPacket(seq uint8, r *bufio.Reader, w *bufio.Writer) (n int, err error) {

	return
}