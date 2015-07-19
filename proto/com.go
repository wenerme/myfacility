package proto

type CommandPack interface {
	Pack
	CommandType() CommandType
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
var ComDaemon = comPack(COM_DAEMON)

type comPack CommandType

func (p comPack) Read(c Proto) {
	c.Get(1, IgnoreByte)
}
func (p comPack) Write(c Proto) {
	c.Put(p)
}
func (p comPack) CommandType() CommandType {
	return CommandType(p)
}

// A COM_QUERY is used to send the server a text-based query that is executed immediately.
// The server replies to a COM_QUERY packet with a COM_QUERY Response.
// The length of the query-string is a taken from the packet length - 1.
type ComQuery struct {
	Query string
}

func (p *ComQuery) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Query, StrEof)
}
func (p *ComQuery) Write(c Proto) {
	c.Put(COM_QUERY, &p.Query, StrEof)
}
func (p *ComQuery) CommandType() CommandType {
	return COM_QUERY
}

type ComInitDb struct {
	Schema string
}

func (p *ComInitDb) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Schema, StrEof)
}
func (p *ComInitDb) Write(c Proto) {
	c.Put(p.CommandType(), &p.Schema, StrEof)
}
func (p *ComInitDb) CommandType() CommandType {
	return COM_INIT_DB
}

type ComCreateDb struct {
	Schema string
}

func (p *ComCreateDb) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Schema, StrEof)
}
func (p *ComCreateDb) Write(c Proto) {
	c.Put(p.CommandType(), &p.Schema, StrEof)
}
func (p *ComCreateDb) CommandType() CommandType {
	return COM_CREATE_DB
}

type ComSetOption struct {
	Option uint16
}

func (p *ComSetOption) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Option)
}
func (p *ComSetOption) Write(c Proto) {
	c.Put(COM_SET_OPTION, &p.Option)
}
func (p *ComSetOption) CommandType() CommandType {
	return COM_SET_OPTION
}

type ComChangeUser struct {
	Username       string
	AuthResponse   []byte
	SchemaName     string
	CharacterSet   CharacterSet
	AuthPluginName string
	Attributes     map[string]string
}

func (p *ComChangeUser) Read(c Proto) {

	c.Get(1, IgnoreByte, &p.Username, StrNul)
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

func (p *ComChangeUser) Write(c Proto) {
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
func (p *ComChangeUser) CommandType() CommandType {
	return COM_CHANGE_USER
}

type ComDropDb struct {
	Schema string
}

func (p *ComDropDb) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Schema, StrEof)
}
func (p *ComDropDb) Write(c Proto) {
	c.Put(COM_DROP_DB, &p.Schema, StrEof)
}
func (p *ComDropDb) CommandType() CommandType {
	return COM_DROP_DB
}

type ComRefresh struct {
	Subcommand uint8
}

func (p *ComRefresh) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Subcommand)
}
func (p *ComRefresh) Write(c Proto) {
	c.Put(COM_REFRESH, &p.Subcommand)
}
func (p *ComRefresh) CommandType() CommandType {
	return COM_REFRESH
}

type ComShutdown struct {
	ShutdownType uint8
}

func (p *ComShutdown) Read(c Proto) {
	c.Get(1, IgnoreByte)
	if c.More() {
		c.Get(&p.ShutdownType)
	}
}
func (p *ComShutdown) Write(c Proto) {
	// TODO Define shutdown type
	c.Put(COM_SHUTDOWN, &p.ShutdownType)
}
func (p *ComShutdown) CommandType() CommandType {
	return COM_SHUTDOWN
}

type ComProcessKill struct {
	ProcessId uint32
}

func (p *ComProcessKill) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.ProcessId)
}
func (p *ComProcessKill) Write(c Proto) {
	c.Put(COM_PROCESS_KILL, &p.ProcessId)
}
func (p *ComProcessKill) CommandType() CommandType {
	return COM_PROCESS_KILL
}

type ComFieldList struct {
	Table string
	Field string
}

func (p *ComFieldList) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Table, StrNul, &p.Field, StrEof)
}
func (p *ComFieldList) Write(c Proto) {
	c.Put(p.CommandType(), &p.Table, StrNul, &p.Field, StrEof)
}
func (p *ComFieldList) CommandType() CommandType {
	return COM_FIELD_LIST
}

type ComPack struct {
	Type CommandType
	Data []byte
}

func (p *ComPack) Read(c Proto) {
	c.Get(&p.Data, StrEof)
	p.Type = CommandType(p.Data[0])
}
func (p *ComPack) Write(c Proto) {
	c.Put(&p.Data, StrEof)
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
func NewCommandPacketMap() map[CommandType]CommandPack {
	m := map[CommandType]CommandPack{
		COM_SLEEP:               ComSleep,
		COM_QUIT:                ComQuit,
		COM_INIT_DB:             &ComInitDb{},
		COM_QUERY:               &ComQuery{},
		COM_FIELD_LIST:          &ComFieldList{},
		COM_CREATE_DB:           &ComCreateDb{},
		COM_DROP_DB:             &ComDropDb{},
		COM_REFRESH:             &ComRefresh{},
		COM_SHUTDOWN:            &ComShutdown{},
		COM_STATISTICS:          ComStatistics,
		COM_PROCESS_INFO:        ComProcessInfo,
		COM_CONNECT:             ComConnect,
		COM_PROCESS_KILL:        ComProcessInfo,
		COM_DEBUG:               ComDebug,
		COM_PING:                ComPing,
		COM_TIME:                ComTime,
		COM_DELAYED_INSERT:      ComDelayInsert,
		COM_CHANGE_USER:         &ComChangeUser{},
		COM_BINLOG_DUMP:         &ComBinlogDump{},
		COM_TABLE_DUMP:          &ComTableDump{},
		COM_CONNECT_OUT:         ComConnectOut,
		COM_REGISTER_SLAVE:      &ComRegisterSlave{},
		COM_STMT_PREPARE:        &ComStmtPrepare{},
		COM_STMT_EXECUTE:        &ComStmtExecute{},
		COM_STMT_SEND_LONG_DATA: &ComStmtSendLongData{},
		COM_STMT_CLOSE:          &ComStmtClose{},
		COM_STMT_RESET:          &ComStmtReset{},
		COM_SET_OPTION:          &ComSetOption{},
		COM_STMT_FETCH:          &ComStmtFetch{},
		COM_DAEMON:              ComDaemon,
		COM_BINLOG_DUMP_GTID:    &ComBinlogDumpGtid{},
		COM_RESET_CONNECTION:    ComResetConnection,
	}
	return m
}
