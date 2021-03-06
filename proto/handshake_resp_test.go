package proto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandshakeResponse41Basic(t *testing.T) {
	data := []byte{
		0x54, 0x00, 0x00, 0x01, 0x8d, 0xa6, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x01, 0x08, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x70, 0x61, 0x6d, 0x00, 0x14, 0xab, 0x09, 0xee, 0xf6, 0xbc, 0xb1, 0x32,
		0x3e, 0x61, 0x14, 0x38, 0x65, 0xc0, 0x99, 0x1d, 0x95, 0x7d, 0x75, 0xd4, 0x47, 0x74, 0x65, 0x73,
		0x74, 0x00, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70,
		0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x00,
	}
	_ = data
	c := CLIENT_BASIC_FLAGS ^ CLIENT_CONNECT_ATTRS | CLIENT_SECURE_CONNECTION
	p := &HandshakeResponse{}
	assert := assert.New(t)
	_, _, _ = c, p, assert
	assertCodec(data, p, c, func() {
		assert.EqualValues("pam", p.Username)
		assert.EqualValues("test", p.Database)
		assert.EqualValues("mysql_native_password", p.AuthPluginName)
	}, t)
}

func TestHandshakeResponse41Send(t *testing.T) {
	data := []byte{
		0x2d, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0x80, 0xff, 0xff, 0xff, 0x0, 0x21, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x72, 0x6f, 0x6f, 0x74, 0x0, 0x0, 0x74, 0x65, 0x73, 0x74, 0x0, 0x0, 0x0,
	}
	c := CLIENT_PROTOCOL_41 | CLIENT_PLUGIN_AUTH | CLIENT_SECURE_CONNECTION | CLIENT_CONNECT_WITH_DB
	p := &HandshakeResponse{}
	assert := assert.New(t)
	_ = assert
	assertCodec(data, p, c, func() {

	}, SkipEqual)
}

func TestHandshakeResponse41General(t *testing.T) {
	data := []byte{
		0xac, 0x00, 0x00, 0x01, 0x0d, 0xa6, 0x3f, 0x20, 0x00, 0x00, 0x00, 0x01, 0x21, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x72, 0x6f, 0x6f, 0x74, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74, 0x00, 0x6d,
		0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73,
		0x77, 0x6f, 0x72, 0x64, 0x00, 0x6a, 0x03, 0x5f, 0x6f, 0x73, 0x08, 0x6f, 0x73, 0x78, 0x31, 0x30,
		0x2e, 0x31, 0x30, 0x0c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
		0x08, 0x6c, 0x69, 0x62, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x04, 0x5f, 0x70, 0x69, 0x64, 0x05, 0x33,
		0x30, 0x31, 0x31, 0x32, 0x0f, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72,
		0x73, 0x69, 0x6f, 0x6e, 0x07, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x37, 0x09, 0x5f, 0x70, 0x6c,
		0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x06, 0x78, 0x38, 0x36, 0x5f, 0x36, 0x34, 0x0c, 0x70, 0x72,
		0x6f, 0x67, 0x72, 0x61, 0x6d, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x05, 0x6d, 0x79, 0x73, 0x71, 0x6c,
	}
	_ = data
	c := CLIENT_PROTOCOL_41 | CLIENT_PLUGIN_AUTH | CLIENT_SECURE_CONNECTION | CLIENT_CONNECT_WITH_DB | CLIENT_CONNECT_ATTRS
	p := &HandshakeResponse{}
	assert := assert.New(t)
	_, _, _ = c, p, assert
	attr := map[string]string{
		"_pid":            "30112",
		"_client_version": "10.0.17",
		"_platform":       "x86_64",
		"_os":             "osx10.10",
		"_client_name":    "libmysql",
		"program_name":    "mysql",
	}
	assertCodec(data, p, c, func() {
		assert.EqualValues(33, p.CharacterSet)
		assert.EqualValues(0x1000000, p.MaxPacketSize)
		assert.EqualValues("root", p.Username)
		assert.EqualValues("test", p.Database)
		assert.EqualValues("mysql_native_password", p.AuthPluginName)
		assert.EqualValues(attr, p.Attributes)
	}, SkipEqual, t)
}
