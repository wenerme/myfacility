package proto
import (
	"io"
	"errors"
	"bufio"
)

type PackReader struct {
	*bufio.Reader
	*baseCodec
	buf    [8]byte
	Packet *Packet
}

func NewPackReader(r io.Reader) PackReader {
	return newPacketReader(r, newBaseCodec())
}
func newPacketReader(r io.Reader, base *baseCodec) PackReader {
	return &PackReader{
		//		buf: make([]byte, 8),
		baseCodec:base,
		Reader: bufio.NewReader(r),
		Packet: &Packet{},
	}
}

func (this *PackReader)HasMore() bool {
	_, err := this.Peek(1)
	return err == nil
}

func (this *PackReader)MustReadAll(data... interface{}) {
	for _, v := range data {
		this.MustRead(v)
	}
}
func (this *PackReader)MustRead(data interface{}) {

}
func (this *PackReader)MustReadInt1() Int1 {
	_, err := this.Read(this.buf[0:1])
	if (err != nil) {
		panic(err)
	}
	return Int1(this.buf[0])
}

func (this *PackReader)MustReadInt2() Int2 {
	return Int2(this.MustReadInt(2))
}
func (this *PackReader)MustReadInt3() Int3 {
	return Int3(this.MustReadInt(3))
}
func (this *PackReader)MustReadInt4() Int4 {
	r := this.MustReadInt(4);
	return Int4(r)
}
func (this *PackReader)MustReadInt6() Int6 {
	r := this.MustReadInt(6);
	return Int6(r)
}
func (this *PackReader)MustReadInt8() Int8 {
	r := this.MustReadInt(8);
	return Int8(r)
}

func (this *PackReader)MustReadInt(size int) (uint64) {
	switch size{case 1, 2, 3, 4, 6, 8: default: panic(errors.New("Can not read size "+string(size)+" int"))}
	_, err := this.Read(this.buf[0:size])
	if err != nil {panic(err)}

	r, err := this.order.Uint(size, this.buf)
	if err != nil {panic(err)}
	return r
}
func (this *PackReader)MustPeek(n int) []byte {
	b, err := this.Peek(n)
	if err != nil {panic(err)}
	return b
}
