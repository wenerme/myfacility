package proto
import (
	"testing"
	"bytes"
	"github.com/cyfdecyf/bufio"
)

func TestPacketCoder(t *testing.T) {

	// A COM_QUIT looks like this:
	// length: 1
	// sequence_id: x00
	// payload: 0x01

	data := []byte{01, 00, 00, 00, 01}
	pack := bytes.NewBuffer(make([]byte, 0))
	ReadPacket(bufio.NewReader(data), pack)
}
