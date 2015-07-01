package proto
import (
	"math"
)



// General Pack interface
// Will not return error, you can use recover to capture
type Pack interface {
	WritablePack
	ReadablePack
}

type WritablePack interface {
	Write(c *PackWriter)
}

type ReadablePack interface {
	Read(c *PackReader)
}

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
//
type Packet struct {
	//	PayloadLength Int3
	SequenceId uint64
	Payload    []byte
}

func (p *Packet)Read(c *PackReader) {
	len := c.MustReadInt3()
	c.MustRead(&p.SequenceId)
	p.Payload = c.MustReadStrV(uint(len))
}

func (p *Packet)Write(c *PackWriter) {
	c.MustWriteAll(Int3(len(p.Payload)), p.SequenceId, p.Payload)
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
type HandshakeV10 struct {
	ProtocolVersion   uint8
	ServerVersion     string
	ConnectionId      uint32
	Challenge1        string
	Capability        uint32
	CharacterSet      uint64
	Status            uint16
	Challenge2        string
	AuthPluginDataLen uint64
	AuthPluginName    string
}
func (p *HandshakeV10)Read(c *Reader) {
	c.Get(&p.ProtocolVersion, &p.ServerVersion, StrNul, &p.ConnectionId, &p.Challenge1, StrVar, 8)
	//  1              [00] filler
	c.SkipBytes(1)


	if c.HasMore() {
		c.MustReadAll(&p.CharacterSet, &p.Status)
		p.Capability |= uint32(uint16(c.MustReaduint16())) << 16

		if c.HasCapability(CLIENT_PLUGIN_AUTH) {
			c.MustRead(&p.AuthPluginDataLen)
		}else {
			c.MustReaduint64()
		}

		//string[10]     reserved (all [00])
		c.MustReadStrF(10)

		if c.HasCapability(CLIENT_SECURE_CONNECTION) {
			// ($len=MAX(13, length of auth-plugin-data - 8))
			// -1 to strip the last \x00 char
			p.Challenge2 = c.MustReadStrV(uint(math.Max(13, float64(p.AuthPluginDataLen)-8)) - 1)
			c.MustReaduint64()// waste the \x00 char
		}

		if c.HasCapability(CLIENT_PLUGIN_AUTH) {
			c.MustRead(&p.AuthPluginName)
		}
	}
}

func (p *HandshakeV10)Write(c PackWriter) {
	c.MustWriteAll(&p.ProtocolVersion, &p.ServerVersion, &p.ConnectionId,
		p.Challenge1[0:8], // len = 8
		uint64(0), // filter
		uint16(p.Capability&0xffff), // lower
		&p.CharacterSet, &p.Status,
		uint16(p.Capability >> 16), // upper
	)

	if c.HasCapability(CLIENT_PLUGIN_AUTH) {
		c.MustWrite(&p.AuthPluginDataLen)
	}else {
		c.MustWriteuint64(0)
	}


	//string[10]     reserved (all [00])
	c.MustWriteNuint64(10, 0)

	if c.HasCapability(CLIENT_SECURE_CONNECTION) {
		c.MustWriteStrV(p.Challenge2)
		c.MustWriteuint64(0)
	}

	if c.HasCapability(CLIENT_PLUGIN_AUTH) {
		c.MustWrite(&p.AuthPluginName)
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
type HandshakeResponse41 struct {
	Capability     uint32
	MaxPacketSize  uint32
	CharacterSet   uint64
	//string[23]     reserved (all [0])
	Username       string
	AuthResponse   []byte
	Database       string
	AuthPluginName string
	Attributes     map[string]string
}

func (p *HandshakeResponse41)Read(c *PackReader) {
	c.MustReadAll(&p.Capability, &p.MaxPacketSize, &p.CharacterSet)
	//  string[23]     reserved (all [0])
	c.MustReadStrF(23)
	c.MustRead(&p.Username)

	if c.HasCapability(CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA) {
		l := c.MustReadIntL()
		p.AuthResponse = stringc.MustReadStrF(uint(l)))
	}else if c.HasCapability(CLIENT_SECURE_CONNECTION) {
		l := c.MustReaduint64()
		p.AuthResponse = stringc.MustReadStrF(uint(l)))
	}else {
		p.AuthResponse = stringc.MustReadStrN())
	}

	if c.HasCapability(CLIENT_CONNECT_WITH_DB) {
		c.MustRead(&p.Database)
	}
	if c.HasCapability(CLIENT_PLUGIN_AUTH) {
		c.MustRead(&p.AuthPluginName)
	}

	if c.HasCapability(CLIENT_CONNECT_ATTRS) {
		c.MustReadIntL()// length
		p.Attributes = make(map[string]string)
		for c.HasMore() {
			k, v := c.MustReadStrL(), c.MustReadStrL()
			p.Attributes[string(k)]=string(v)
		}
	}
}

func (p *HandshakeResponse41)Write(c *PackWriter) {
	c.MustWriteAll(&p.Capability, &p.MaxPacketSize, &p.CharacterSet)
	//  string[23]     reserved (all [0])
	c.MustWriteNuint64(23, 0)
	c.MustWrite(&p.Username)

	if c.HasCapability(CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA) {
		c.MustWriteIntL(uint(len(p.AuthResponse)))
		c.MustWriteStrF(string(p.AuthResponse))
	}else if c.HasCapability(CLIENT_SECURE_CONNECTION) {
		c.MustWriteuint64(uint64(len(p.AuthResponse)))
		c.MustWriteStrF(string(p.AuthResponse))
	}else {
		c.MustWriteStrN(string(p.AuthResponse))
	}

	if c.HasCapability(CLIENT_CONNECT_WITH_DB) {
		c.MustWrite(&p.Database)
	}
	if c.HasCapability(CLIENT_PLUGIN_AUTH) {
		c.MustWrite(&p.AuthPluginName)
	}

	if c.HasCapability(CLIENT_CONNECT_ATTRS) {
		l := 0
		for k, v := range p.Attributes {
			kl, vl := uint(len(k)), uint(len(v))
			l += int(kl)+int(vl)
			l += kl.Len()+vl.Len()
		}
		c.MustWriteIntL(uint(l))
		for k, v := range p.Attributes {
			c.MustWriteStrL(string(k))
			c.MustWriteStrL(string(v))
		}
	}
}