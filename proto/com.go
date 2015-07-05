package proto

type CommandPack interface {
	Pack
	Type() Command
}

var ComSleep = comPack(COM_SLEEP)
var ComQuit = comPack(COM_QUIT)
var ComResetConnection = comPack(COM_RESET_CONNECTION)
var ComStatistics = comPack(COM_STATISTICS)
var ComProcessInfo = comPack(COM_PROCESS_INFO)
var ComConnect = comPack(COM_CONNECT)
var ComDebug = comPack(COM_DEBUG)
var ComPing = comPack(COM_PING)
var ComTime = comPack(COM_TIME)
var ComDelayInsert = comPack(COM_DELAYED_INSERT)
var ComConnectOut = comPack(COM_CONNECT_OUT)

type comPack Command

func (p comPack) Read(c Reader) {
	c.SkipBytes(1)
}
func (p comPack) Write(c Writer) {
	c.Put(p)
}
func (p comPack) Type() Command {
	return Command(p)
}

type stringComPack struct {
	Command Command
}

func (p *stringComPack) Read(c Reader) {
	c.SkipBytes(1)
}
func (p *stringComPack) Write(c Writer) {
	c.Put(p.Command)
}
func (p *stringComPack) Type() Command {
	return p.Command
}

//
// A COM_QUERY is used to send the server a text-based query that is executed immediately.
// <p/>
// The server replies to a COM_QUERY packet with a COM_QUERY Response.
// <p/>
// The length of the query-string is a taken from the packet length - 1.
// <p/>
// Payload
// 1              [03] COM_QUERY
// string[EOF]    the query the server shall execute
// Fields
// command_id (1) -- 0x03 COM_QUERY
// <p/>
// query (string.EOF) -- query_text
// <p/>
// Implemented By
// mysql_query()
// <p/>
// Returns
// COM_QUERY_Response
// <p/>
// Example
// 21 00 00 00 03 73 65 6c    65 63 74 20 40 40 76 65    !....select @@ve
// 72 73 69 6f 6e 5f 63 6f    6d 6d 65 6e 74 20 6c 69    rsion_comment li
// 6d 69 74 20 31                                        mit 1
type ComQuery struct {
	Query string
}

func (p *ComQuery) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.Query, StrEof)
}
func (p *ComQuery) Write(c Writer) {
	c.Put(COM_QUERY, &p.Query, StrEof)
}
func (p *ComQuery) Type() Command {
	return COM_QUERY
}

type ComInitDb struct {
	Schema string
}

func (p *ComInitDb) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.Schema, StrEof)
}
func (p *ComInitDb) Write(c Writer) {
	c.Put(COM_INIT_DB, &p.Schema, StrEof)
}
func (p *ComInitDb) Type() Command {
	return COM_INIT_DB
}

type ComChangeUser struct {
	Username       string
	AuthResponse   []byte
	SchemaName     string
	CharacterSet   CharacterSet
	AuthPluginName string
	Attributes     map[string]string
}

func (p *ComChangeUser) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.Username, StrNul)
	if c.HasCap(CLIENT_SECURE_CONNECTION) {
		var n uint8
		c.Get(&n)
		c.Get(&p.AuthResponse, StrVar, int(n))
	} else {
		c.Get(&p.AuthResponse, StrNul)
	}
	c.Get(&p.SchemaName, StrNul)
	if c.More() {
		c.Get(&p.CharacterSet)
		if c.HasCap(CLIENT_PLUGIN_AUTH) {
			c.Get(&p.AuthPluginName, StrNul)
		}

		if c.HasCap(CLIENT_CONNECT_ATTRS) {
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
}

func (p *ComChangeUser) Write(c Writer) {
	c.Put(COM_CHANGE_USER, p.Username, StrNul)
	if c.HasCap(CLIENT_SECURE_CONNECTION) {
		c.Put(uint8(len(p.AuthResponse)), p.AuthResponse, StrEof)
	} else {
		c.Put(p.AuthResponse, StrNul)
	}
	c.Put(&p.SchemaName, StrNul)
	// no character set
	if p.CharacterSet == 0 {
		return
	}
	c.Put(&p.CharacterSet)
	if c.HasCap(CLIENT_PLUGIN_AUTH) {
		c.Put(&p.AuthPluginName, StrNul)
	}
	if c.HasCap(CLIENT_CONNECT_ATTRS) {
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
func (p *ComChangeUser) Type() Command {
	return COM_CHANGE_USER
}

type ComBinlogDump struct {
	BinlogPos      uint32
	Flags          uint16
	ServerId       uint32
	BinlogFilename string
}

func (p *ComBinlogDump) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.BinlogPos, &p.Flags, &p.ServerId, &p.BinlogFilename, StrEof)
}
func (p *ComBinlogDump) Write(c Writer) {
	c.Put(COM_BINLOG_DUMP, &p.BinlogPos, &p.Flags, &p.ServerId, &p.BinlogFilename, StrEof)
}
func (p *ComBinlogDump) Type() Command {
	return COM_BINLOG_DUMP
}

/*
1              [15] COM_REGISTER_SLAVE
4              server-id
1              slaves hostname length
string[$len]   slaves hostname
1              slaves user len
string[$len]   slaves user
1              slaves password len
string[$len]   slaves password
2              slaves mysql-port
4              replication rank
4              master-id
*/
type ComRegisterSlave struct {
	ServerId        uint32
	SlaveHostname   string
	SlaveUser       string
	SlavePassword   string
	SlavePort       uint16
	ReplicationRank uint32
	MasterId        uint32
}

func (p *ComRegisterSlave) Read(c Reader) {
	c.SkipBytes(1)
	var n uint8
	c.Get(
		&p.ServerId,
		&n, &p.SlaveHostname, StrVar, n,
		&n, &p.SlaveUser, StrVar, n,
		&n, &p.SlavePassword, StrVar, n,
		&p.SlavePort,
		&p.ReplicationRank,
		&p.MasterId,
	)
}
func (p *ComRegisterSlave) Write(c Writer) {
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

type ComTableDump struct {
	Database string
	Table    string
}

func (p *ComTableDump) Read(c Reader) {
	c.SkipBytes(1)
	var n uint8
	c.Get(&n)
	c.Get(&p.Database, StrVar, int(n), &n)
	c.Get(&p.Table, StrVar, int(n))
}
func (p *ComTableDump) Write(c Writer) {
	c.Put(COM_TABLE_DUMP, uint8(len(p.Database)), p.Database, uint8(len(p.Table)), p.Table)
}
func (p *ComTableDump) Type() Command {
	return COM_TABLE_DUMP
}

type ComDropDb struct {
	Schema string
}

func (p *ComDropDb) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.Schema, StrEof)
}
func (p *ComDropDb) Write(c Writer) {
	c.Put(COM_DROP_DB, &p.Schema, StrEof)
}
func (p *ComDropDb) Type() Command {
	return COM_DROP_DB
}

type ComRefresh struct {
	Subcommand uint8
}

func (p *ComRefresh) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.Subcommand)
}
func (p *ComRefresh) Write(c Writer) {
	c.Put(COM_REFRESH, &p.Subcommand)
}
func (p *ComRefresh) Type() Command {
	return COM_REFRESH
}

type ComShutdown struct {
	ShutdownType uint8
}

func (p *ComShutdown) Read(c Reader) {
	c.SkipBytes(1)
	if c.More() {
		c.Get(&p.ShutdownType)
	}
}
func (p *ComShutdown) Write(c Writer) {
	// TODO Define shutdown type
	c.Put(COM_SHUTDOWN, &p.ShutdownType)
}
func (p *ComShutdown) Type() Command {
	return COM_SHUTDOWN
}

type ComProcessKill struct {
	ProcessId uint32
}

func (p *ComProcessKill) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.ProcessId)
}
func (p *ComProcessKill) Write(c Writer) {
	c.Put(COM_PROCESS_KILL, &p.ProcessId)
}
func (p *ComProcessKill) Type() Command {
	return COM_PROCESS_KILL
}

type ComFieldList struct {
	Table string
	Field string
}

func (p *ComFieldList) Read(c Reader) {
	c.SkipBytes(1)
	c.Get(&p.Table, StrNul, &p.Field, StrEof)
}
func (p *ComFieldList) Write(c Writer) {
	c.Put(p.Type(), &p.Table, StrNul, &p.Field, StrEof)
}
func (p *ComFieldList) Type() Command {
	return COM_FIELD_LIST
}

type ComPack struct {
	Type Command
	data []byte
	buf  *Buffer
}

func (p *ComPack) Read(c Reader) {
	c.Get(&p.data, StrEof)
	p.Type = Command(p.data[0])
}
func (p *ComPack) Write(c Writer) {
	c.Put(&p.data, StrEof)
}

func (p *ComPack) ReadPack() (com CommandPack) {
	switch p.Type {
	case COM_SLEEP:
		return ComSleep
	case COM_QUIT:
		return ComQuit
	case COM_INIT_DB:
		com = &ComInitDb{}
	case COM_QUERY:
		com = &ComQuery{}
	}
	return nil
}
