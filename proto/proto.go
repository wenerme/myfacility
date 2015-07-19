package proto

type Pack interface {
	Read(Proto)
	Write(Proto)
}

// Protocol maintainer
type Proto interface {
	/* read */
	Get(...interface{})
	Peek(int) ([]byte, error)
	PeekByte() (byte, error)
	More() bool

	/* write */
	Put(...interface{})

	/* pack */
	ReadPacket(Pack)
	WritePacket(Pack)
	RecvPacket() (int, error)
	// Send buf as a packet
	// Will increase the sequence id
	SendPacket() (int, error)
	MustRecvPacket() int
	MustSendPacket() int
	RecvReadPacket(Pack) (int, error)
	WriteSendPacket(Pack) (int, error)

	/* ctx */
	SetSeq(uint8)
	Seq() uint8
	HasCap(Capability) bool
	Cap() Capability
	SetCap(Capability)
	Com() CommandType
	SetCom(CommandType)
}
