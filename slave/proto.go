package slave
import (
	"bufio"
	"bytes"
)

// Packets to bytes
type Encoder struct {
	*bytes.Buffer

}


// Bytes to Packets
type Decoder struct {
	*bytes.Reader

}

// Connection <-> Packet
type Protocol struct {
	*bufio.ReadWriter
	Compressed   bool
	Payload      []byte
	SequenceId   uint32
	ConnectionId uint32
}

type PW struct {
	*Basic
	Payload      []byte
	SequenceId   uint32
	ConnectionId uint32

}
type Basic struct {
	Payload      []byte
	SequenceId   uint32
	ConnectionId uint32
	Compressed   bool
}