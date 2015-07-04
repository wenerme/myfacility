package proto

type (
	Pack interface {
		Read(Reader)
		Write(Writer)
	}

	Proto interface {
		/* read */
		Get(...interface{})
		GetType(interface{}, ProtoType) (int, error)
		SkipBytes(int)
		Peek(int) ([]byte, error)
		PeekByte() (byte, error)
		More() bool

		/* write */
		Put(...interface{})
		PutType(interface{}, ProtoType) (int, error)
		PutZero(int)

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
	}
)
