package proto
import (
	"bufio"
	"encoding/binary"
)

func ReadPacket(r *bufio.Reader, w *bufio.Writer) (seq uint8, err error) {
	var l uint32
	err = binary.Read(r, binary.LittleEndian, &l)
	if err != nil {return }
	seq = l &0xff
	l = l >> 8
	var b byte
	for i := 0; i < l; i ++ {
		b, err = r.ReadByte()
		if err != nil {return }
		w.WriteByte(b)
	}
	return
}

func ReadCompressedPacket(r *bufio.Reader, w *bufio.Writer) (seq uint8, err error) {
	var l uint32
	err = binary.Read(r, binary.LittleEndian, &l)
	if err != nil {return }
	seq = l &0xff
	l = l >> 8
	var b byte
	for i := 0; i < l; i ++ {
		b, err = r.ReadByte()
		if err != nil {return }
		w.WriteByte(b)
	}
	return
}


func WritePacket(seq uint8, r *bufio.Reader, w *bufio.Writer) (n int, err error) {

	return
}
func WriteCompressedPacket(seq uint8, r *bufio.Reader, w *bufio.Writer) (n int, err error) {

	return
}