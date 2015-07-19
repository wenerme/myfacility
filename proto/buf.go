package proto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/op/go-logging"
	"io"
)

type Buffer struct {
	Reader
	Writer
	buf *bytes.Buffer
	con io.ReadWriter
	cap Capability
	com CommandType
	seq uint8
}

func NewBuffer(con io.ReadWriter, buf *bytes.Buffer) *Buffer {
	if buf == nil {
		buf = bytes.NewBufferString("")
	}
	b := &Buffer{buf: buf, con: con}
	b.Reader = NewReader(buf)
	b.Writer = NewWriter(buf)
	return b
}

func (r *Buffer) SetSeq(seq uint8) {
	r.seq = seq
}
func (r *Buffer) Seq() uint8 {
	return r.seq
}
func (r *Buffer) SetCap(cap Capability) {
	r.cap = cap
}
func (r *Buffer) Cap() Capability {
	return r.cap
}
func (r *Buffer) HasCap(cap Capability) bool {
	return r.cap.Has(cap)
}

func (b *Buffer) RecvReadPacket(p Pack) (n int, err error) {
	n, err = b.RecvPacket()
	if err != nil {
		return
	}

	p.Read(b)
	return
}
func (b *Buffer) WriteSendPacket(p Pack) (n int, err error) {
	p.Write(b)

	n, err = b.SendPacket()
	return
}
func (b *Buffer) ReadPacket(p Pack) {
	p.Read(b)
}
func (b *Buffer) WritePacket(p Pack) {
	p.Write(b)
}
func (b *Buffer) MustRecvPacket() int {
	n, err := b.RecvPacket()
	if err != nil {
		panic(err)
	}
	return n
}
func (b *Buffer) MustSendPacket() int {
	n, err := b.SendPacket()
	if err != nil {
		panic(err)
	}
	return n
}
func (b *Buffer) RecvPacket() (n int, err error) {
	var l uint32
	err = binary.Read(b.con, binary.LittleEndian, &l)
	if err != nil {
		return
	}
	b.seq = uint8(l >> 24)
	l = l << 8 >> 8
	var written int64
	written, err = io.CopyN(b.buf, b.con, int64(l))
	n = 4 + int(written)
	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Recv packet#%d (%d)\n%s", b.seq, n, hex.Dump(b.buf.Bytes()))
	}
	return
}
func (b *Buffer) SendPacket() (n int, err error) {
	var l uint32
	l = uint32(len(b.buf.Bytes()))
	l = l | (uint32(b.seq) << 24)
	err = binary.Write(b.con, binary.LittleEndian, l)
	if err != nil {
		return
	}
	n, err = b.con.Write(b.buf.Bytes())
	n += 4
	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Send packet#%d (%d)\n%s", b.seq, n, hex.Dump(b.buf.Bytes()))
	}
	b.seq++
	b.buf.Reset()
	return
}

func (r *Buffer) SetCom(com CommandType) {
	r.com = com
}
func (r *Buffer) Com() CommandType {
	return r.com
}
