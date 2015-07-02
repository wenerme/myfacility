package proto
import (
	"testing"
	"bytes"
	"bufio"
	"github.com/stretchr/testify/assert"
)

func TestPacketCoder(t *testing.T) {
	assert := assert.New(t)
	// A COM_QUIT looks like this:
	// length: 1
	// sequence_id: x02
	// payload: 0x01

	data := []byte{0x01, 00, 00, 0x2, 0x01}
	pack := bytes.NewBuffer(make([]byte, 0))
	r := &BufReader{Reader:bufio.NewReader(pack)}
	seq, size, err := ReadPacketTo(bufio.NewReader(bytes.NewReader(data)), pack)
	assert.NoError(err)
	assert.EqualValues(1, size)
	assert.EqualValues(2, seq)

	assert.True(r.More())
	var com uint8
	r.Get(&com)
	assert.EqualValues(1, com)
	assert.False(r.More())
}
