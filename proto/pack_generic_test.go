package proto
import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestEOF(t *testing.T) {
	// a 4.1 EOF packet with: 0 warnings, AUTOCOMMIT enabled.
	// 05 00 00 05 fe 00 00 02 00
	data := []byte{0x05, 0x00, 0x00, 0x05, 0xfe, 0x00, 0x00, 0x02, 0x00, }
	assert := assert.New(t)
	p := &EOFPack{}
	assertData(data, p, CLIENT_ALL_FLAGS, assert)
}
func TestOK(t *testing.T) {
	// OK with CLIENT_PROTOCOL_41. 0 affected rows, last-insert-id was 0, AUTOCOMMIT, 0 warnings. No further info.
	// 07 00 00 02 00 00 00 02    00 00 00
	// ...........
	data := []byte{0x07, 0x00, 0x00, 0x02, 0x00, 0x01, 0x02, 0x03, 0x00, 0x00, 0x00, }
	assert := assert.New(t)
	p := &OKPack{}
	assertData(data, p, CLIENT_PROTOCOL_41, assert)
	//    t.Logf("%#v\n",p.State())
}
func TestERR(t *testing.T) {
	// 17 00 00 01 ff 48 04 23    48 59 30 30 30 4e 6f 20       .....H.#HY000No
	// 74 61 62 6c 65 73 20 75    73 65 64                      tables used
	data := []byte{
		0x17, 0x00, 0x00, 0x01, 0xff, 0x48, 0x04, 0x23, 0x48, 0x59, 0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x20, 0x75, 0x73, 0x65, 0x64, }
	assert := assert.New(t)
	p := &ERRPack{}
	assertData(data, p, CLIENT_PROTOCOL_41, assert)
	//    t.Logf("%#v\n",p.State())
}