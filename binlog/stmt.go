package binlog

import "github.com/wenerme/myfacility/proto"

/*
Statement Based Replication Events
Statement Based Replication or SBR sends the SQL queries a client sent to the master AS IS to the slave.
It needs extra events to mimic the client connection's state on the slave side.
*/

// The query event is used to send text querys right the binlog.
// http://dev.mysql.com/doc/internals/en/query-event.html
type QueryEvent struct {
	SlaveProxyId  uint32
	ExecutionTime uint32
	ErrorCode     uint16
	Status        []byte
	Schema        string
	Query         string
}

func (p *QueryEvent) Read(c proto.Reader) {
	var m uint8
	var n uint16
	c.Get(&p.SlaveProxyId,
		&p.ExecutionTime,
		&m,
		&p.ErrorCode,
		&n,
		&p.Status, proto.StrVar, &n,
		&p.Schema, proto.StrVar, &m,
		1, proto.IgnoreByte,
		&p.Query, proto.StrEof)
}
func (p *QueryEvent) EventType() EventType {
	return QUERY_EVENT
}

type IntvarEvent struct {
	VarType IntVarType
	Value   uint64
}

func (p *IntvarEvent) Read(c proto.Reader) {
	c.Get(&p.VarType, &p.Value)
}
func (p *IntvarEvent) EventType() EventType {
	return INTVAR_EVENT
}

type RandEvent struct {
	Seed1 uint64
	Seed2 uint64
}

func (p *RandEvent) Read(c proto.Reader) {
	c.Get(&p.Seed1, &p.Seed2)
}
func (p *RandEvent) EventType() EventType {
	return RAND_EVENT
}

type UseVarEvent struct {
	Name    string
	IsNull  uint8
	VarType uint8
	Charset uint32
	Value   string
	Flags   uint8
}

func (p *UseVarEvent) Read(c proto.Reader) {
	var n uint32
	c.Get(&n, &p.Name, proto.StrVar, &n, &p.IsNull)
	if p.IsNull == 0 {
		c.Get(&p.VarType, &p.Charset, &n, &p.Value, proto.StrVar, &n)
		if c.More() {
			c.Get(&p.Flags)
		}
	}
}
func (p *UseVarEvent) EventType() EventType {
	return USER_VAR_EVENT
}

type XIDEvent struct {
	XID uint64
}

func (p *XIDEvent) Read(c proto.Reader) {
	c.Get(&p.XID)
}
func (p *XIDEvent) EventType() EventType {
	return XID_EVENT
}
