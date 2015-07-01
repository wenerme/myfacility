package proto
import (
	"bytes"
	"reflect"
	"github.com/stretchr/testify/assert"
)


func createCodec(data []byte, p Pack, cap Capability) (*Codec, *bytes.Buffer) {
	buf := bytes.NewBuffer(make([]byte, 0))
	c := NewCodec(bytes.NewReader(data), buf)
	c.SetCapability(cap)
	if p != nil {
		c.MustReadPacket()
		c.MustReadPack(p)

		//        c.Packet.SequenceId = c.Packet.SequenceId
		c.MustWritePack(p)
		c.MustWritePacket()
		c.MustFlushWrite()
	}
	return c, buf
}
func assertData(data []byte, p Pack, cap Capability, assert *assert.Assertions) *Codec {
	c, buf := createCodec(data, p, cap)
	assert.EqualValues(data, buf.Bytes())
	return c
}

func assertValue(data []byte, p Pack, cap Capability, assert *assert.Assertions) *Codec {
	c, _ := createCodec(data, p, cap)
	np := reflect.New(reflect.TypeOf(p).Elem()).Interface().(Pack)

	r := NewPackReader(bytes.NewReader(c.Packet.Payload))
	r.SetCapability(c.Capability())
	np.Read(r)

	assert.EqualValues(p, np)

	return c
}