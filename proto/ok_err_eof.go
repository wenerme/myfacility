package proto
import "github.com/syndtr/goleveldb/leveldb/errors"

//
// If CLIENT_PROTOCOL_41 is enabled, the EOF packet contains a warning count and status flags.
// <pre>
// <b>Caution</b>
// the EOF packet may appear in places where a Protocol::LengthEncodedInteger may appear. You must check whether the packet length is less than 9 to make sure that it is a EOF packet.
// * <b>Payload</b>
// Type	Name	Description
// int&lt;1>	header	[fe] EOF header
// if capabilities & CLIENT_PROTOCOL_41 {
// int&lt;2>	warnings	number of warnings
// int&lt;2>	status_flags	Status Flags
// }
// <b>Example</b>
// a 4.1 EOF packet with: 0 warnings, AUTOCOMMIT enabled.
// 05 00 00 05 fe 00 00 02 00     ..........
// <pre/>
// <a href=http://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html>EOF_Packet</a>
type EOFPack struct {
	Header   uint8
	Warnings uint16
	Status Status
}
func (p *EOFPack)Read(c Reader) {
	c.Get(&p.Header)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Get(&p.Warnings, &p.Status)
	}
}
func (p *EOFPack)Write(c Writer) {
	c.Put(&p.Header)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Put(&p.Warnings, &p.Status)
	}
}

// 
// This packet signals that an error occurred. It contains a SQL state value if CLIENT_PROTOCOL_41 is enabled.
// <pre>
// <b>Payload</b>
// Type	Name	Description
// int&lt;1>	header	[ff] header of the ERR packet
// int&lt;2>	error_code	error-code
// if capabilities & CLIENT_PROTOCOL_41 {
// string[1]	sql_state_marker	# marker of the SQL State
// string[5]	sql_state	SQL State
// }
// string&lt;EOF>	error_message	human readable error message
// <b>Example</b>
// 17 00 00 01 ff 48 04 23    48 59 30 30 30 4e 6f 20       .....H.#HY000No
// 74 61 62 6c 65 73 20 75    73 65 64                      tables used
// </pre>
// <a href=http://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html>ERR_Packet</a>
// 
type ERRPack struct {
	Header         uint8
	ErrorCode      uint16
	SQLStateMarker string
	SQLState       string
	ErrorMessage   string
}
func (p *ERRPack)Read(c Reader) {
	c.Get(&p.Header, &p.ErrorCode)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Get(&p.SQLStateMarker, StrVar, 1, &p.SQLState, StrVar, 5)
	}
	c.Get(&p.ErrorMessage, StrEof)
}
func (p *ERRPack)Write(c Writer) {
	c.Put(&p.Header, &p.ErrorCode)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Put(&p.SQLStateMarker, StrVar, 1, &p.SQLState, StrVar, 5)
	}
	c.Put(&p.ErrorMessage, StrEof)
}
func (p ERRPack)Error() string {
	return string(p.ErrorMessage)
}
// 
// An OK packet is sent from the server to the client to signal successful completion of a command.
// <p/>
// If CLIENT_PROTOCOL_41 is set, the packet contains a warning count.
// <pre>
// <b>Payload</b>
// Type	Name	Description
// int&lt;1>	header	[00] the OK packet header
// int&lt;lenenc>	affected_rows	affected rows
// int&lt;lenenc>	last_insert_id	last insert-id
// if capabilities & CLIENT_PROTOCOL_41 {
// int&lt;2>	status_flags	Status Flags
// int&lt;2>	warnings	number of warnings
// } elseif capabilities & CLIENT_TRANSACTIONS {
// int&lt;2>	status_flags	Status Flags
// }
// if capabilities & CLIENT_SESSION_TRACK {
// string&lt;lenenc>	info	human readable status information
// if status_flags & SERVER_SESSION_STATE_CHANGED {
// string&lt;lenenc>	session_state_changes	session state info
// }
// } else {
// string&lt;EOF>	info	human readable status information
// }
// <b>Example</b>
// OK with CLIENT_PROTOCOL_41. 0 affected rows, last-insert-id was 0, AUTOCOMMIT, 0 warnings. No further info.
//  * 07 00 00 02 00 00 00 02    00 00 00
// ...........
// </pre>
// http://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
type OKPack struct {
	Header       uint8
	AffectedRows uint64
	LastInsertId uint64
	Status Status
	Warnings     uint16
	Info         string
	SessionState string
}

func (p *OKPack)Read(c Reader) {
	c.Get(&p.Header, &p.AffectedRows, IntEnc, &p.LastInsertId, IntEnc)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Get(&p.Warnings, &p.Status)
	}else if c.HasCap(CLIENT_TRANSACTIONS) {
		c.Get(&p.Status)
	}

	if c.HasCap(CLIENT_SESSION_TRACK) {
		c.Get(&p.Info)
		if Status(p.Status).Has(SERVER_SESSION_STATE_CHANGED) {
			c.Get(&p.SessionState)
		}
	}else {
		c.Get(&p.Info, StrEof)
	}
}
func (p *OKPack)Write(c Writer) {
	c.Put(&p.Header, &p.AffectedRows, IntEnc, &p.LastInsertId, IntEnc)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Put(&p.Warnings, &p.Status)
	}else if c.HasCap(CLIENT_TRANSACTIONS) {
		c.Put(&p.Status)
	}

	if c.HasCap(CLIENT_SESSION_TRACK) {
		c.Put(&p.Info)
		if Status(p.Status).Has(SERVER_SESSION_STATE_CHANGED) {
			c.Put(&p.SessionState)
		}
	}else {
		c.Put(&p.Info, StrEof)
	}
}
var ErrNotStatePack = errors.New("Not OK,ERR of EOF packet")
func ReadStatePack(proto Proto) (p Pack, err error) {
	b, err := proto.PeekByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case 0:
		p = &OKPack{}
	case 0xFF:
		p = &ERRPack{}
	case 0xFE:
		p = &EOFPack{}
	default:
		return nil, ErrNotStatePack
	}
	proto.ReadPacket(p)
	return
}