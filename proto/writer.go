package proto
import (
	"io"
	"bufio"
	"errors"
)
type Flusher interface {
	Flush() error
}

type PackWriter struct {
	io.Writer
	*baseCodec
	buf []byte
}

func NewPacketWriter(w io.Writer) *PackWriter {
	return newPacketWriter(w, newBaseCodec())
}
func newPacketWriter(w io.Writer, base *baseCodec) *PackWriter {
	return &PackWriter{
		Writer: bufio.NewWriter(w),
		baseCodec:base,
		buf:make([]byte, 8),
	}
}

func (this *PackWriter)Flush() error {
	// Flush when writer is a bufio
	if w, ok := this.Writer.(Flusher); ok {
		return w.Flush()
	}
	return nil
}

func (this *PackWriter)MustWrite(data interface{}) {

}
func (this *PackWriter)MustWriteAll(data... interface{}) {
	for _, v := range data {
		this.MustWrite(v)
	}
}

func (this *PackWriter)MustWriteInt1(v Int1) {
	_, err := this.Write([]byte{byte(v)})
	if err != nil {panic(err)}
}
func (this *PackWriter)MustWriteNInt1(n int, v Int1) {
	for ; n > 0; n-=1 {
		this.MustWriteInt1(v)
	}
}
func (this *PackWriter)MustWriteInt2(v Int2) {
	this.MustWriteInt(2, uint64(v))
}
func (this *PackWriter)MustWriteInt3(v Int3) {
	this.MustWriteInt(3, uint64(v))
}
func (this *PackWriter)MustWriteInt4(v Int4) {
	this.MustWriteInt(4, uint64(v))
}
func (this *PackWriter)MustWriteInt6(v Int6) {
	this.MustWriteInt(6, uint64(v))
}
func (this *PackWriter)MustWriteInt8(v Int8) {
	this.MustWriteInt(8, uint64(v))
}
func (this *PackWriter)MustWriteInt(size int, v uint64) {
	switch size{case 1, 2, 3, 4, 6, 8: default: panic(errors.New("Can not read size "+string(size)+" int"))}
	err := this.order.PutUint(size, this.buf, v)
	if err != nil {panic(err)}
	_, err = this.Write(this.buf[0:size])
	if err != nil {panic(err)}
}
