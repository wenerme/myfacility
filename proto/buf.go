package proto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/op/go-logging"
	"io"
)

// A proto implement
type buffer struct {
	Reader
	Writer
	buf  *bytes.Buffer
	conn io.ReadWriter
	cap  Capability
	com  CommandType
	seq  uint8
}

func NewProto(conn io.ReadWriter, buf *bytes.Buffer) Proto {
	if buf == nil {
		buf = bytes.NewBufferString("")
	}
	b := &buffer{buf: buf, conn: conn}
	b.Reader = NewReader(buf)
	b.Writer = NewWriter(buf)
	return b
}

func (r *buffer) SetSeq(seq uint8) {
	r.seq = seq
}
func (r *buffer) Seq() uint8 {
	return r.seq
}
func (r *buffer) SetCap(cap Capability) {
	r.cap = cap
}
func (r *buffer) Cap() Capability {
	return r.cap
}
func (r *buffer) HasCap(cap Capability) bool {
	return r.cap.Has(cap)
}

func (b *buffer) RecvReadPacket(p Pack) (n int, err error) {
	n, err = b.RecvPacket()
	if err != nil {
		return
	}

	p.Read(b)
	return
}
func (b *buffer) WriteSendPacket(p Pack) (n int, err error) {
	p.Write(b)

	n, err = b.SendPacket()
	return
}
func (b *buffer) ReadPacket(p Pack) {
	p.Read(b)
}
func (b *buffer) WritePacket(p Pack) {
	p.Write(b)
}
func (b *buffer) MustRecvPacket() int {
	n, err := b.RecvPacket()
	if err != nil {
		panic(err)
	}
	return n
}
func (b *buffer) MustSendPacket() int {
	n, err := b.SendPacket()
	if err != nil {
		panic(err)
	}
	return n
}
func (b *buffer) RecvPacket() (n int, err error) {
	var l uint32
	err = binary.Read(b.conn, binary.LittleEndian, &l)
	if err != nil {
		return
	}
	b.seq = uint8(l >> 24)
	l = l << 8 >> 8
	var written int64
	written, err = io.CopyN(b.buf, b.conn, int64(l))
	n = 4 + int(written)
	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Recv packet#%d (%d)\n%s", b.seq, n, hex.Dump(b.buf.Bytes()))
	}
	return
}
func (b *buffer) SendPacket() (n int, err error) {
	var l uint32
	l = uint32(len(b.buf.Bytes()))
	l = l | (uint32(b.seq) << 24)
	err = binary.Write(b.conn, binary.LittleEndian, l)
	if err != nil {
		return
	}
	n, err = b.conn.Write(b.buf.Bytes())
	n += 4
	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Send packet#%d (%d)\n%s", b.seq, n, hex.Dump(b.buf.Bytes()))
	}
	b.seq++
	b.buf.Reset()
	return
}

func (r *buffer) SetCom(com CommandType) {
	r.com = com
}
func (r *buffer) Com() CommandType {
	return r.com
}

func (r *buffer) Conn() io.ReadWriter {
	return r.conn
}
