package proto

import (
	"math"
	//	"golang.org/x/tools/cmd/stringer"
)

// http://dev.mysql.com/doc/internals/en/mysql-packet.html
type Packet struct {
	//	PayloadLength Int3
	SequenceId uint64
	Payload    []byte
}

func (p *Packet) Read(c Proto) {
	var len uint
	c.Get(&len, &p.SequenceId, &p.Payload, StrVar, int(len))
}

func (p *Packet) Write(c Proto) {
	c.Put(uint(len(p.Payload)), p.SequenceId, p.Payload, StrEof)
}

/*
TODO Due to Bug#59453 the auth-plugin-name is missing the terminating NUL-char in versions prior to 5.5.10 and 5.6.2.
*/
// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeV10
type Handshake struct {
	ProtocolVersion uint8
	ServerVersion   string
	ConnectionId    uint32
	Challenge1      string
	Capability      Capability
	CharacterSet    uint8
	Status          uint16
	Challenge2      string
	AuthPluginName  string
}

func (p *Handshake) Read(c Proto) {
	c.Get(
		&p.ProtocolVersion,
		&p.ServerVersion, StrNul,
		&p.ConnectionId,
		&p.Challenge1, StrVar, 8,
	)
	//  1              [00] filler
	c.Get(1, IgnoreByte)
	var t uint16
	c.Get(&t)
	p.Capability = Capability(t)
	if c.More() {
		c.Get(&p.CharacterSet, &p.Status)
		c.Get(&t)
		p.Capability = p.Capability | Capability(t)<<16

		cap := Capability(p.Capability)
		var authPluginDataLen uint8
		if cap.Has(CLIENT_PLUGIN_AUTH) {
			c.Get(&authPluginDataLen)
		} else {
			c.Get(1, IgnoreByte)
		}

		//string[10]     reserved (all [00])
		c.Get(10, IgnoreByte)

		if cap.Has(CLIENT_SECURE_CONNECTION) {
			// ($len=MAX(13, length of auth-plugin-data - 8))
			// -1 to strip the last \x00 char
			c.Get(&p.Challenge2, StrVar, int(math.Max(13, float64(authPluginDataLen)-8))-1)
			c.Get(1, IgnoreByte) // waste the \x00 char
		}

		if cap.Has(CLIENT_PLUGIN_AUTH) {
			c.Get(&p.AuthPluginName, StrNul)
		}
	}
}

func (p *Handshake) Write(c Proto) {
	c.Put(
		&p.ProtocolVersion,
		&p.ServerVersion, StrNul,
		&p.ConnectionId,
		p.Challenge1, StrVar, 8, // len = 8
		uint8(0),             // filter
		uint16(p.Capability), // lower
		&p.CharacterSet, &p.Status,
		uint16(p.Capability>>16), // upper
	)
	ca := Capability(p.Capability)
	if ca.Has(CLIENT_PLUGIN_AUTH) {
		c.Put(uint8(len(p.Challenge2) + 8 + 1))
	} else {
		c.Put(1, IgnoreByte)
	}

	//string[10]     reserved (all [00])
	c.Put(10, IgnoreByte)

	if ca.Has(CLIENT_SECURE_CONNECTION) {
		c.Put(p.Challenge2, StrNul)
		//		c.PutZero(1)
	}

	if ca.Has(CLIENT_PLUGIN_AUTH) {
		c.Put(&p.AuthPluginName, StrNul)
	}
}

/*
  If client wants to have a secure SSL connection and sets CLIENT_SSL flag it should first send
  the SSL Request Packet and only then, after establishing the secure connection, it should send
  the Handshake Response Packet.
*/
// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse41
type HandshakeResponse struct {
	Capability    Capability
	MaxPacketSize uint32
	CharacterSet  CharacterSet //CharacterSet here is int1
	//string[23]     reserved (all [0])
	Username       string
	AuthResponse   []byte
	Database       string
	AuthPluginName string
	Attributes     map[string]string
}

func (p *HandshakeResponse) Read(c Proto) {
	c.Get(&p.Capability, &p.MaxPacketSize, &p.CharacterSet, Int1)
	//  string[23]     reserved (all [0])
	c.Get(23, IgnoreByte, &p.Username, StrNul)
	cap := Capability(p.Capability)
	if cap.Has(CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA) {
		c.Get(&p.AuthResponse)
	} else if cap.Has(CLIENT_SECURE_CONNECTION) {
		var n uint8
		c.Get(&n, &p.AuthResponse, StrVar, &n)
	} else {
		c.Get(&p.AuthResponse, StrNul)
	}

	if cap.Has(CLIENT_CONNECT_WITH_DB) {
		c.Get(&p.Database, StrNul)
	}
	if cap.Has(CLIENT_PLUGIN_AUTH) {
		c.Get(&p.AuthPluginName, StrNul)
	}

	if cap.Has(CLIENT_CONNECT_ATTRS) {
		var len uint
		var k, v string
		c.Get(&len) // length
		p.Attributes = make(map[string]string)
		for c.More() {
			c.Get(&k, &v)
			p.Attributes[k] = v
		}
	}
}

func (p *HandshakeResponse) Write(c Proto) {
	c.Put(&p.Capability, &p.MaxPacketSize, &p.CharacterSet, Int1)
	//  string[23]     reserved (all [0])
	c.Put(23, IgnoreByte, &p.Username, StrNul)
	cap := Capability(p.Capability)
	if cap.Has(CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA) {
		c.Put(p.AuthResponse)
	} else if cap.Has(CLIENT_SECURE_CONNECTION) {
		c.Put(uint8(len(p.AuthResponse)), p.AuthResponse, StrEof)
	} else {
		c.Put(p.AuthResponse, StrNul)
	}

	if cap.Has(CLIENT_CONNECT_WITH_DB) {
		c.Put(&p.Database, StrNul)
	}
	if cap.Has(CLIENT_PLUGIN_AUTH) {
		c.Put(&p.AuthPluginName, StrNul)
	}

	if cap.Has(CLIENT_CONNECT_ATTRS) {
		var l uint
		for k, v := range p.Attributes {
			kl, vl := uint(len(k)), uint(len(v))
			l += kl + vl + bytesOfIntVar(uint64(kl)) + bytesOfIntVar(uint64(vl))
		}
		c.Put(uint(l))
		for k, v := range p.Attributes {
			c.Put(k, v)
		}
	}
}

func bytesOfIntVar(i uint64) uint {
	switch {
	case i < 251:
		return 1
	case i < 0xffff:
		return 3
	default:
		return 8
	}
}

// SSL Connection Request Packet.
// It is like Handshake Response Packet but is truncated right before username field.
// If server supports CLIENT_SSL capability, client can send this packet to request a secure SSL connection. The CLIENT_SSL capability flag must be set inside the SSL Connection Request Packet.
// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest
type SSLRequest struct {
	Capability    Capability
	MaxPacketSize uint32
	CharacterSet  CharacterSet
}

func (p *SSLRequest) Read(c Proto) {
	c.Get(&p.Capability, &p.MaxPacketSize, &p.CharacterSet, Int1, 23, IgnoreByte)
}
func (p *SSLRequest) Write(c Proto) {
	c.Put(&p.Capability, &p.MaxPacketSize, &p.CharacterSet, Int1, 23, IgnoreByte)
}

type AuthSwitchRequest struct {
	PluginName     string
	AuthPluginData string
}

func (p *AuthSwitchRequest) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.PluginName, StrNul, &p.AuthPluginData, StrEof)
}
func (p *AuthSwitchRequest) Write(c Proto) {
	c.Put(uint8(0xFE), &p.PluginName, StrNul, &p.AuthPluginData, StrEof)
}

// Authentication Method Switch Response Packet which contains response data generated by the authenticatication method requested in Authentication Method Switch Request Packet. This data is opaque to the protocol.
type AuthSwitchResponse struct {
	AuthPluginResponse string
}

func (p *AuthSwitchResponse) Read(c Proto) {
	c.Get(&p.AuthPluginResponse, StrEof)
}
func (p *AuthSwitchResponse) Write(c Proto) {
	c.Put(&p.AuthPluginResponse, StrEof)
}

type AuthMoreData struct {
	PluginData string
}

func (p *AuthMoreData) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.PluginData, StrEof)
}
func (p *AuthMoreData) Write(c Proto) {
	c.Put(uint8(0x01), &p.PluginData, StrEof)
}
