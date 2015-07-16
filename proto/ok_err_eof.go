package proto

import "github.com/syndtr/goleveldb/leveldb/errors"

// http://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html
type EOFPack struct {
	Header   uint8
	Warnings uint16
	Status   Status
}

func (p *EOFPack) Read(c Proto) {
	c.Get(&p.Header)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Get(&p.Warnings, &p.Status)
	}
}
func (p *EOFPack) Write(c Proto) {
	c.Put(&p.Header)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Put(&p.Warnings, &p.Status)
	}
}

// http://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
type ERRPack struct {
	Header         uint8
	ErrorCode      uint16
	SQLStateMarker string
	SQLState       string
	ErrorMessage   string
}

func (p *ERRPack) Read(c Proto) {
	c.Get(&p.Header, &p.ErrorCode)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Get(&p.SQLStateMarker, StrVar, 1, &p.SQLState, StrVar, 5)
	}
	c.Get(&p.ErrorMessage, StrEof)
}
func (p *ERRPack) Write(c Proto) {
	c.Put(&p.Header, &p.ErrorCode)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Put(&p.SQLStateMarker, StrVar, 1, &p.SQLState, StrVar, 5)
	}
	c.Put(&p.ErrorMessage, StrEof)
}
func (p ERRPack) Error() string {
	return string(p.ErrorMessage)
}

// http://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
type OKPack struct {
	Header       uint8
	AffectedRows uint64
	LastInsertId uint64
	Status       Status
	Warnings     uint16
	Info         string
	SessionState SessionState
}

func (p *OKPack) Read(c Proto) {
	c.Get(&p.Header, &p.AffectedRows, IntEnc, &p.LastInsertId, IntEnc)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Get(&p.Warnings, &p.Status)
	} else if c.HasCap(CLIENT_TRANSACTIONS) {
		c.Get(&p.Status)
	}

	if c.HasCap(CLIENT_SESSION_TRACK) {
		c.Get(&p.Info)
		if Status(p.Status).Has(SERVER_SESSION_STATE_CHANGED) {
			c.Get(&p.SessionState)
		}
	} else {
		c.Get(&p.Info, StrEof)
	}
}
func (p *OKPack) Write(c Proto) {
	c.Put(&p.Header, &p.AffectedRows, IntEnc, &p.LastInsertId, IntEnc)
	if c.HasCap(CLIENT_PROTOCOL_41) {
		c.Put(&p.Warnings, &p.Status)
	} else if c.HasCap(CLIENT_TRANSACTIONS) {
		c.Put(&p.Status)
	}

	if c.HasCap(CLIENT_SESSION_TRACK) {
		c.Put(&p.Info)
		if Status(p.Status).Has(SERVER_SESSION_STATE_CHANGED) {
			c.Put(&p.SessionState)
		}
	} else {
		c.Put(&p.Info, StrEof)
	}
}

var ErrNotStatePack = errors.New("Not OK,ERR of EOF packet")

func ReadErrOk(proto Proto) (p Pack, err error) {
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
func ReadErrEof(proto Proto) (p Pack, err error) {
	b, err := proto.PeekByte()
	if err != nil {
		return nil, err
	}
	switch b {
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
