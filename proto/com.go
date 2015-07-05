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
var ComDaemon = comPack(COM_DAEMON)

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

// A COM_QUERY is used to send the server a text-based query that is executed immediately.
// The server replies to a COM_QUERY packet with a COM_QUERY Response.
// The length of the query-string is a taken from the packet length - 1.
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

type ComSetOption struct {
	Option uint16
}

func (p *ComSetOption) Read(c Reader) {
	c.Get(1, IgnoreByte, &p.Option)
}
func (p *ComSetOption) Write(c Writer) {
	c.Put(COM_SET_OPTION, &p.Option)
}
func (p *ComSetOption) Type() Command {
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
