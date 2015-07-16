package proto

type ComTableDump struct {
	Database string
	Table    string
}

func (p *ComTableDump) Read(c Proto) {
	var n uint8
	c.Get(1, IgnoreByte, &n, &p.Database, StrVar, &n, &n, &p.Table, StrVar, &n)
}
func (p *ComTableDump) Write(c Proto) {
	c.Put(COM_TABLE_DUMP, uint8(len(p.Database)), p.Database, uint8(len(p.Table)), p.Table)
}
func (p *ComTableDump) Type() Command {
	return COM_TABLE_DUMP
}

type ComBinlogDump struct {
	BinlogPos      uint32
	Flags          uint16
	ServerId       uint32
	BinlogFilename string
}

func (p *ComBinlogDump) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.BinlogPos, &p.Flags, &p.ServerId, &p.BinlogFilename, StrEof)
}
func (p *ComBinlogDump) Write(c Proto) {
	c.Put(COM_BINLOG_DUMP, &p.BinlogPos, &p.Flags, &p.ServerId, &p.BinlogFilename, StrEof)
}
func (p *ComBinlogDump) Type() Command {
	return COM_BINLOG_DUMP
}

/*
1              [1e] COM_BINLOG_DUMP_GTID
2              flags
4              server-id
4              binlog-filename-len
string[len]    binlog-filename
8              binlog-pos
  if flags & BINLOG_THROUGH_GTID {
4              data-size
string[len]    data
  }
*/
type ComBinlogDumpGtid struct {
	Flags          uint16
	ServerId       uint32
	BinlogFilename string
	BinlogPos      uint64
	Data           []byte
}

func (p *ComBinlogDumpGtid) Read(c Proto) {
	var n uint32
	c.Get(1, IgnoreByte,
		&p.Flags,
		&p.ServerId,
		&n, &p.BinlogFilename, StrVar, &n,
		&p.BinlogPos)
	if p.Flags&BINLOG_THROUGH_GTID > 0 {
		c.Get(&n, &p.Data, StrVar, &n)
	}
}
func (p *ComBinlogDumpGtid) Write(c Proto) {
	c.Put(COM_BINLOG_DUMP_GTID,
		&p.Flags,
		&p.ServerId,
		uint32(len(p.BinlogFilename)), &p.BinlogFilename, StrEof,
		&p.BinlogPos)
	if p.Flags&BINLOG_THROUGH_GTID > 0 {
		c.Put(uint32(len(p.Data)), &p.Data, StrEof)
	}
}
func (p *ComBinlogDumpGtid) Type() Command {
	return COM_BINLOG_DUMP_GTID
}

type ComRegisterSlave struct {
	ServerId        uint32
	SlaveHostname   string
	SlaveUser       string
	SlavePassword   string
	SlavePort       uint16
	ReplicationRank uint32
	MasterId        uint32
}

func (p *ComRegisterSlave) Read(c Proto) {
	var n uint8
	c.Get(
		1, IgnoreByte,
		&p.ServerId,
		&n, &p.SlaveHostname, StrVar, &n,
		&n, &p.SlaveUser, StrVar, &n,
		&n, &p.SlavePassword, StrVar, &n,
		&p.SlavePort,
		&p.ReplicationRank,
		&p.MasterId,
	)
}
func (p *ComRegisterSlave) Write(c Proto) {
	c.Put(
		&p.ServerId,
		uint8(len(p.SlaveHostname)), &p.SlaveHostname, StrEof,
		uint8(len(p.SlaveUser)), &p.SlaveUser, StrEof,
		uint8(len(p.SlavePassword)), &p.SlavePassword, StrEof,
		&p.SlavePort,
		&p.ReplicationRank,
		&p.MasterId,
	)
}
func (p *ComRegisterSlave) Type() Command {
	return COM_REGISTER_SLAVE
}
