package proto
import (
	"testing"
)

func TestHandShake10(t *testing.T) {
	data := []byte{
		//length
		0x5f, 0x00, 0x00,
		//sequence id
		0x10,
		//handshake version 10
		0x0a,
		//mysql plain text version
		0x35, 0x2e, 0x35, 0x2e, 0x33, 0x38, 0x2d, 0x30, 0x75, 0x62, 0x75, 0x6e, 0x74, 0x75, 0x30, 0x2e, 0x31,
		0x34, 0x2e, 0x30, 0x34, 0x2e, 0x31, 0x2d, 0x6c, 0x6f, 0x67, 0x00,
		//connection id
		0x05, 0x00, 0x00, 0x00,
		//auth-plugin-data-part-1 = ROw,ng;0
		0x52, 0x4f, 0x77, 0x2c, 0x6e, 0x67, 0x3b, 0x30,
		//filler
		0x00,
		//capability flags (lower 2 bytes)
		0xff, 0xf7,
		//charset
		0x08,
		//status flag
		0x02, 0x00,
		//capability flags (upper 2 bytes)
		0x0f, 0x80,
		//auth data length = 21
		0x15,
		//reserved 10 bytes
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//auth-plugin-data-part-2 = }F&):(W`Z%Gv
		0x7d, 0x46, 0x26, 0x29, 0x3a, 0x28, 0x57, 0x60, 0x5a, 0x25, 0x47, 0x76, 0x00,
		//auth-plugin name = mysql_native_password
		0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73,
		0x77, 0x6f, 0x72, 0x64, 0x00,
	}
	_ = data
	p := &Handshake{}
	assertCodec(data, p, CLIENT_BASIC_FLAGS, DumpPacket, t)
}

