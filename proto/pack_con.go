package proto
import (
	"math"
//	"golang.org/x/tools/cmd/stringer"
)

//
// If a MySQL client or server wants to send data, it:
//
// Splits the data into packets of size (224–1) bytes
//
// Prepends to each chunk a packet header
//
// Protocol::Packet
// Data between client and server is exchanged in packets of max 16MByte size.
//
// Payload
// Type	Name	Description
// int&lt;3>	payload_length	Length of the payload. The number of bytes in the packet beyond the initial 4 bytes that make up the packet header.
// int&lt;1>	sequence_id	Sequence ID
// string&lt;var>	payload	[len=payload_length] payload of the packet
// Example
// A COM_QUIT looks like this:
//
// 01 00 00 00 01
// length: 1
// sequence_id: x00
// payload: 0x01
//
// <a href=http://dev.mysql.com/doc/internals/en/mysql-packet.html>mysql-packet</a>
type Packet struct {
	//	PayloadLength Int3
	SequenceId uint64
	Payload    []byte
}

func (p *Packet)Read(c Reader) {
	var len uint
	c.Get(&len, &p.SequenceId, &p.Payload, StrVar, int(len))
}

func (p *Packet)Write(c Writer) {
	c.Put(uint(len(p.Payload)), p.SequenceId, p.Payload, StrEof)
}

//
// <pre>
// 1              [0a] protocol version
// string[NUL]    server version
// 4              connection id
// string[8]      auth-plugin-data-part-1
// 1              [00] filler
// 2              capability flags (lower 2 bytes)
// if more data in the packet:
// 1              character set
// 2              status flags
// 2              capability flags (upper 2 bytes)
// if capabilities & CLIENT_PLUGIN_AUTH {
// 1              length of auth-plugin-data
// } else {
// 1              [00]
// }
// string[10]     reserved (all [00])
// if capabilities & CLIENT_SECURE_CONNECTION {
// string[$len]   auth-plugin-data-part-2 ($len=MAX(13, length of auth-plugin-data - 8))
// if capabilities & CLIENT_PLUGIN_AUTH {
// string[NUL]    auth-plugin name
// }
//
// Fields
// protocol_version (1) -- 0x0a protocol_version
//
// server_version (string.NUL) -- human-readable server version
//
// connection_id (4) -- connection id
//
// auth_plugin_data_part_1 (string.fix_len) -- [len=8] first 8 bytes of the auth-plugin data
//
// filler_1 (1) -- 0x00
//
// capability_flag_1 (2) -- lower 2 bytes of the Protocol::CapabilityFlags (optional)
//
// character_set (1) -- default server character-set, only the lower 8-bits Protocol::CharacterSet (optional)
//
// status_flags (2) -- Protocol::StatusFlags (optional)
//
// capability_flags_2 (2) -- upper 2 bytes of the Protocol::CapabilityFlags
//
// auth_plugin_data_len (1) -- length of the combined auth_plugin_data, if auth_plugin_data_len is > 0
//
// auth_plugin_name (string.NUL) -- name of the auth_method that the auth_plugin_data belongs to
//
// <b>Note</b>
// Due to Bug#59453 the auth-plugin-name is missing the terminating NUL-char in versions prior to 5.5.10 and 5.6.2.
// </pre>
//
// <a href=http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeV10>HandshakeV10</a>
//
type Handshake struct {
	ProtocolVersion uint8
	ServerVersion   string
	ConnectionId    uint32
	Challenge1      string
	Capability      uint32
	CharacterSet    uint8
	Status          uint16
	Challenge2      string
	AuthPluginName  string
}
func (p *Handshake)Read(c Reader) {
	c.Get(
		&p.ProtocolVersion,
		&p.ServerVersion, StrNul,
		&p.ConnectionId,
		&p.Challenge1, StrVar, 8,
	)
	//  1              [00] filler
	c.SkipBytes(1)
	var t uint16
	c.Get(&t)
	p.Capability = uint32(t)
	if c.More() {
		c.Get(&p.CharacterSet, &p.Status)
		c.Get(&t)
		p.Capability |= uint32(t) << 16

		cap := Capability(p.Capability)
		var authPluginDataLen uint8
		if cap.Has(CLIENT_PLUGIN_AUTH) {
			c.Get(&authPluginDataLen)
		}else {
			c.SkipBytes(1)
		}

		//string[10]     reserved (all [00])
		c.SkipBytes(10)

		if cap.Has(CLIENT_SECURE_CONNECTION) {
			// ($len=MAX(13, length of auth-plugin-data - 8))
			// -1 to strip the last \x00 char
			c.Get(&p.Challenge2, StrVar, int(math.Max(13, float64(authPluginDataLen)-8)) - 1)
			c.SkipBytes(1)// waste the \x00 char
		}

		if cap.Has(CLIENT_PLUGIN_AUTH) {
			c.Get(&p.AuthPluginName, StrNul)
		}
	}
}

func (p *Handshake)Write(c Writer) {
	c.Put(
		&p.ProtocolVersion,
		&p.ServerVersion, StrNul,
		&p.ConnectionId,
		p.Challenge1, StrVar, 8, // len = 8
		uint8(0), // filter
		uint16(p.Capability), // lower
		&p.CharacterSet, &p.Status,
		uint16(p.Capability >> 16), // upper
	)
	cap := Capability(p.Capability)
	if cap.Has(CLIENT_PLUGIN_AUTH) {
		c.Put(uint8(len(p.Challenge2) + 8 + 1))
	}else {
		c.PutZero(1)
	}

	//string[10]     reserved (all [00])
	c.PutZero(10)

	if cap.Has(CLIENT_SECURE_CONNECTION) {
		c.Put(p.Challenge2, StrNul)
		//		c.PutZero(1)
	}

	if cap.Has(CLIENT_PLUGIN_AUTH) {
		c.Put(&p.AuthPluginName, StrNul)
	}
}
//
// Handshake Response Packet sent by 4.1+ clients supporting CLIENT_PROTOCOL_41 capability, if the server announced it in its Initial Handshake Packet. Otherwise (talking to an old server) the Protocol::HandshakeResponse320 packet has to be used.
// <pre>
//
//
// Payload
// 4              capability flags, CLIENT_PROTOCOL_41 always set
// 4              max-packet size
// 1              character set
// string[23]     reserved (all [0])
// string[NUL]    username
// if capabilities & CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA {
// lenenc-int     length of auth-response
// string[n]      auth-response
// } else if capabilities & CLIENT_SECURE_CONNECTION {
// 1              length of auth-response
// string[n]      auth-response
// } else {
// string[NUL]    auth-response
// }
// if capabilities & CLIENT_CONNECT_WITH_DB {
// string[NUL]    database
// }
// if capabilities & CLIENT_PLUGIN_AUTH {
// string[NUL]    auth plugin name
// }
// if capabilities & CLIENT_CONNECT_ATTRS {
// lenenc-int     length of all key-values
// lenenc-str     key
// lenenc-str     value
// if-more data in 'length of all key-values', more keys and value pairs
// }
//
// Fields
// capability_flags (4) -- capability flags of the client as defined in Protocol::CapabilityFlags
//
// max_packet_size (4) -- max size of a command packet that the client wants to send to the server
//
// character_set (1) -- connection's default character set as defined in Protocol::CharacterSet.
//
// username (string.fix_len) -- name of the SQL account which client wants to log in -- this string should be interpreted using the character set indicated by character set field.
//
// auth-response (string.NUL) -- opaque authentication response data generated by Authentication Method indicated by the plugin name field.
//
// database (string.NUL) -- initail database for the connection -- this string should be interpreted using the character set indicated by character set field.
//
// auth plugin name (string.NUL) -- the Authentication Method used by the client to generate auth-response value in this packet. This is an UTF-8 string.
//                                  oo.bar
// Caution
// Currently, multibyte character sets such as UCS2, UTF16 and UTF32 are not supported.
//
// Note
// If client wants to have a secure SSL connection and sets CLIENT_SSL flag it should first send the SSL Request Packet and only then, after establishing the secure connection, it should send the Handshake Response Packet.
// </pre>
//
// <a href="http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse41">HandshakeResponse41</a>
//
type HandshakeResponse struct {
	Capability     uint32
	MaxPacketSize  uint32
	CharacterSet uint8
	//string[23]     reserved (all [0])
	Username       string
	AuthResponse   []byte
	Database       string
	AuthPluginName string
	Attributes     map[string]string
}

func (p *HandshakeResponse)Read(c Reader) {
	c.Get(&p.Capability, &p.MaxPacketSize, &p.CharacterSet)
	//  string[23]     reserved (all [0])
	c.SkipBytes(23)
	c.Get(&p.Username, StrNul)
	cap := Capability(p.Capability)
	if cap.Has(CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA) {
		c.Get(&p.AuthResponse)
	}else if cap.Has(CLIENT_SECURE_CONNECTION) {
		var n uint8
		c.Get(&n, &p.AuthResponse, StrVar, int(n))
	}else {
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
		c.Get(&len)// length
		p.Attributes = make(map[string]string)
		for c.More() {
			c.Get(&k, &v)
			p.Attributes[k]=v
		}
	}
}

func (p *HandshakeResponse)Write(c Writer) {
	c.Put(&p.Capability, &p.MaxPacketSize, &p.CharacterSet)
	//  string[23]     reserved (all [0])
	c.PutZero(23)
	c.Put(&p.Username, StrNul)
	cap := Capability(p.Capability)
	if cap.Has(CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA) {
		c.Put(p.AuthResponse)
	}else if cap.Has(CLIENT_SECURE_CONNECTION) {
		c.Put(uint8(len(p.AuthResponse)), p.AuthResponse, StrEof)
	}else {
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
			l += kl + vl + bytesOfIntVar(uint64(kl))+ bytesOfIntVar(uint64(vl))
		}
		c.Put(uint(l))
		for k, v := range p.Attributes {
			c.Put(k, v)
		}
	}
}

func bytesOfIntVar(i uint64) uint {
	switch {
	case i<251:return 1
	case i<0xffff:return 3
	default: return 8
	}
}