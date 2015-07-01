package proto
import (
	"io"
	"bytes"
)



type Capabilitible interface {
	HasCapability(Capability) bool
	Capability() Capability
	SetCapability(Capability)
}

// Not thread safe
// Read:
//      ReadPacket
//      ReadPack
// Write:
//      WritePack
//      WritePacket
//type Codec interface {
//    Capabilitible
//    io.Reader
//    io.Writer
//
//    // read Packet from Reader
//    MustReadPacket() *Packet
//    // read Pack from Packet
//    MustReadPack(Pack)
//
//    // write Packet to Writer
//    MustWritePacket()
//    // write Pack to Packet
//    MustWritePack(Pack) *Packet
//    MustFlushWrite()
//
//    // Packet used for read and write
//    Packet() *Packet
//    PayloadWriter() *PackWriter
//    PayloadReader() *PackReader
//    // Writer write to Packet
//    //    PackWriter() PackWriter
//    // Reader read from Packet
//    //    PackReader() PackReader
//}

type Codec struct {
	*PackWriter
	*PackReader
	*baseCodec
	Packet        *Packet
	PayloadWriter *PackWriter   // writer
	PayloadReader PackReader    // reader
	PayloadBuffer *bytes.Buffer // buffer used for write and read
}

func NewCodec(r io.Reader, w io.Writer) *Codec {
	base := newBaseCodec()
	buff := bytes.NewBuffer(make([]byte, 0))
	var pw *PackWriter
	var pr PackReader
	if r != nil {pr = newPacketReader(r, base)}
	if w != nil {pw = newPacketWriter(w, base)}
	c := &Codec{
		PackWriter: pw,
		PackReader: pr,
		baseCodec: base,
		Packet: &Packet{},
		PayloadBuffer: buff,
		PayloadWriter: newPacketWriter(buff, base),
		PayloadReader: newPacketReader(buff, base),
	}
	return c
}


func (this *Codec)MustWritePack(p WritablePack) *Packet {
	this.PayloadBuffer.Reset()
	p.Write(this.PayloadWriter)
	return this.Packet
}
func (this *Codec)MustWritePacket() {
	this.PayloadWriter.Flush()
	this.Packet.Payload = this.PayloadBuffer.Bytes()
	//    this.Packet.PayloadLength = Int3(len(this.Packet.Payload))
	this.Packet.Write(this.PackWriter)
	this.MustFlushWrite()
}
func (this *Codec)MustFlushWrite() {
	this.PackWriter.Flush()
}

func (this *Codec)MustReadPacket() *Packet {
	this.Packet.Read(this)
	this.PayloadBuffer.Reset()
	this.PayloadBuffer.Write(this.Packet.Payload)
	return this.Packet
}
func (this *Codec)MustReadPack(p ReadablePack) {
	this.PayloadBuffer.Reset()
	this.PayloadBuffer.Write(this.Packet.Payload)
	p.Read(this.PayloadReader)
}

func (this *Codec)HasCapability(c Capability) bool {
	return this.baseCodec.HasCapability(c)
}
func (this *Codec)Capability() Capability {
	return this.baseCodec.Capability()
}
func (this *Codec)SetCapability(c Capability) {
	this.baseCodec.SetCapability(c)
}