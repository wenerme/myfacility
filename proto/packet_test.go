package proto

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPacketCoder(t *testing.T) {
	assert := assert.New(t)
	// A COM_QUIT looks like this:
	// length: 1
	// sequence_id: x02
	// payload: 0x01

	data := []byte{0x01, 00, 00, 0x2, 0x01}

	buf := NewBuffer(bytes.NewBuffer(data), nil)
	n, err := buf.RecvPacket()
	assert.NoError(err)
	assert.EqualValues(5, n)
	assert.EqualValues(2, buf.Seq())

	assert.True(buf.More())
	var com uint8
	buf.Get(&com)
	assert.EqualValues(1, com)
	assert.False(buf.More())
}
