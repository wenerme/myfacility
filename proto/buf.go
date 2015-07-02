package proto
import (
	"io"
	"bytes"
	"encoding/binary"
	"bufio"
	"encoding/hex"
	"github.com/op/go-logging"
)

type Proto interface {
	Get(...interface{})
	GetType(interface{}, ProtoType) (int, error)
	SkipBytes(int)
	More() bool

	Put(...interface{})
	PutType(interface{}, ProtoType) (int, error)
	PutZero(int)

	ReadPacket(Pack)
	WritePacket(Pack)
	RecvPacket() (int, error)
	SendPacket() (int, error)

	SetSeq(uint8)
	Seq(uint8)
}

type Buffer struct {
	*BufReader
	*BufWriter
	buf *bytes.Buffer
	con io.ReadWriter
	cap Capability
	seq uint8
}

func NewBuffer(con io.ReadWriter, buf *bytes.Buffer) *Buffer {
	if buf == nil {
		buf = bytes.NewBufferString("")
	}
	b := &Buffer{buf:buf, con:con}
	b.BufReader = &BufReader{Reader: bufio.NewReaderSize(buf, 0x400)}
	b.BufWriter = &BufWriter{Writer: bufio.NewWriterSize(buf, 0x400)}
	return b
}

func (r *Buffer)SetSeq(seq uint8) {
	r.seq = seq
}
func (r *Buffer)Seq() uint8 {
	return r.seq
}
func (r *Buffer)SetCap(cap Capability) {
	r.cap = cap
}
func (r *Buffer)Cap() Capability {
	return r.cap
}
func (r *Buffer)HasCap(cap Capability) bool {
	return r.cap.Has(cap)
}

func (b *Buffer)ReadPacket(p Pack) {
	p.Read(b)
}
func (b *Buffer)WritePacket(p Pack) {
	p.Write(b)
	b.BufWriter.Flush()
}
func (b *Buffer)MustRecvPacket() (int) {
	n, err := b.RecvPacket()
	if err != nil {panic(err)}
	return n
}
func (b *Buffer)MustSendPacket() (int) {
	n, err := b.SendPacket()
	if err != nil {panic(err)}
	return n
}
func (b *Buffer)RecvPacket() (n int, err error) {
	var l uint32
	err = binary.Read(b.con, binary.LittleEndian, &l)
	if err != nil {return }
	b.seq = uint8(l >> 24)
	l = l << 8 >> 8
	// TODO batch read
	data := make([]byte, l)
	var read int
	for {
		read, err = b.con.Read(data[n:])
		n += read
		if err != nil { return }
		if n == int(l) {break }
	}
	b.buf.Write(data)
	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Recv packet#%d (%d)\n%s", b.seq, n, hex.Dump(data))
	}
	return
}
func (b *Buffer)SendPacket() (n int, err error) {
	b.BufWriter.Flush()
	var l uint32
	l = uint32(len(b.buf.Bytes()))
	l = l | (uint32(b.seq) << 24)
	err = binary.Write(b.con, binary.LittleEndian, l)
	if err!= nil {return }
	n, err = b.con.Write(b.buf.Bytes())
	n += 4
	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Send packet#%d (%d)\n%s", b.seq, n, hex.Dump(b.buf.Bytes()))
	}
	b.buf.Reset()
	return
}