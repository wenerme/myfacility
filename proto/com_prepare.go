package proto

import ()

// https://dev.mysql.com/doc/internals/en/com-stmt-prepare-response.htm
type ComStmtPrepareOK struct {
	StatementId  uint32
	WarningCount uint16
	Params       []ColumnDefinition
	Columns      []ColumnDefinition
	EOF          *EOFPack
}

func (p *ComStmtPrepareOK) Read(c Proto) {
	c.MustRecvPacket()
	p.Columns = nil
	p.Params = nil
	var columns, params uint16
	c.Get(
		1, IgnoreByte,
		&p.StatementId,
		&columns, &params,
		1, IgnoreByte, /*reserved_1 (1) -- [00] filler*/
		&p.WarningCount,
	)
	if params > 0 {
		for {
			c.MustRecvPacket()
			if eof, err := ReadErrEof(c); err == nil {
				p.EOF = eof.(*EOFPack)
				break
			} else if err == ErrNotStatePack {
				col := ColumnDefinition{}
				col.Read(c)
				p.Params = append(p.Params, col)
			} else {
				panic(err)
			}
		}
	}
	if columns > 0 {
		for {
			c.MustRecvPacket()
			if eof, err := ReadErrEof(c); err == nil {
				p.EOF = eof.(*EOFPack)
				break
			} else if err == ErrNotStatePack {
				col := ColumnDefinition{}
				col.Read(c)
				p.Columns = append(p.Columns, col)
			} else {
				panic(err)
			}
		}
	}
}
func (p *ComStmtPrepareOK) Write(c Proto) {
	c.Put(
		OK, Int1,
		&p.StatementId,
		uint16(len(p.Columns)), uint16(len(p.Params)),
		1, IgnoreByte, /*reserved_1 (1) -- [00] filler*/
		&p.WarningCount,
	)
	c.MustSendPacket()
	for _, col := range p.Params {
		col.Write(c)
		c.MustSendPacket()
	}
	if len(p.Params) > 0 {
		p.EOF.Write(c)
		c.MustSendPacket()
	}
	for _, col := range p.Columns {
		col.Write(c)
		c.MustSendPacket()
	}
	if len(p.Columns) > 0 {
		p.EOF.Write(c)
		c.MustSendPacket()
	}
}

type ComStmtPrepare struct {
	Query string
}

func (p *ComStmtPrepare) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.Query, StrEof)
}
func (p *ComStmtPrepare) Write(c Proto) {
	c.Put(COM_STMT_PREPARE, p.Query, StrEof)
}
func (p *ComStmtPrepare) Type() Command {
	return COM_STMT_PREPARE
}

type ComStmtExecute struct {
	StmtId             uint32
	Flags              CursorType
	NullBitmap         []byte
	NewParamsBoundFlag uint8
	ParamsType         []uint16
	ParamsValue        []string
	OK                 *ComStmtPrepareOK
}

func (p *ComStmtExecute) Read(c Proto) {
	c.Get(1, IgnoreByte,
		&p.StmtId, &p.Flags,
		4, IgnoreByte, //4 iteration-count always 1
	)
	params := len(p.OK.Params)
	if params > 0 {
		c.Get(&p.NullBitmap, StrVar, (params+7)/8, &p.NewParamsBoundFlag)
	}
	if p.NewParamsBoundFlag == 1 {
		p.ParamsType = make([]uint16, params)
		for i := 0; i < params; i++ {
			c.Get(&p.ParamsType[i])
		}
		p.ParamsValue = make([]string, params)
		for i := 0; i < params; i++ {
			c.Get(&p.ParamsValue[i])
		}
	}
}
func (p *ComStmtExecute) Write(c Proto) {
	c.Put(COM_STMT_EXECUTE, p.StmtId, p.Flags, uint32(1))
	params := len(p.OK.Params)
	if params > 0 {
		c.Put(&p.NullBitmap, StrVar, (params+7)/8, &p.NewParamsBoundFlag)
	}
	if p.NewParamsBoundFlag == 1 {
		for i := 0; i < params; i++ {
			c.Put(&p.ParamsType[i])
		}
		for i := 0; i < params; i++ {
			c.Put(&p.ParamsValue[i])
		}
	}
}
func (p *ComStmtExecute) Type() Command {
	return COM_STMT_EXECUTE
}

type ComStmtSendLongData struct {
	StmtId  uint32
	ParamId uint16
	Data    []byte
}

func (p *ComStmtSendLongData) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.StmtId, &p.ParamId, &p.Data, StrEof)
}
func (p *ComStmtSendLongData) Write(c Proto) {
	c.Put(COM_STMT_SEND_LONG_DATA, &p.StmtId, &p.ParamId, &p.Data, StrEof)
}
func (p *ComStmtSendLongData) Type() Command {
	return COM_STMT_SEND_LONG_DATA
}

type ComStmtClose struct {
	StmtId uint32
}

func (p *ComStmtClose) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.StmtId)
}
func (p *ComStmtClose) Write(c Proto) {
	c.Put(COM_STMT_CLOSE, p.StmtId)
}
func (p *ComStmtClose) Type() Command {
	return COM_STMT_CLOSE
}

type ComStmtReset struct {
	StmtId uint32
}

func (p *ComStmtReset) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.StmtId)
}
func (p *ComStmtReset) Write(c Proto) {
	c.Put(COM_STMT_RESET, p.StmtId)
}
func (p *ComStmtReset) Type() Command {
	return COM_STMT_RESET
}

type ComStmtFetch struct {
	StmtId uint32
	Rows   uint32
}

func (p *ComStmtFetch) Read(c Proto) {
	c.Get(1, IgnoreByte, &p.StmtId, &p.Rows)
}
func (p *ComStmtFetch) Write(c Proto) {
	c.Put(COM_STMT_FETCH, p.StmtId, p.Rows)
}
func (p *ComStmtFetch) Type() Command {
	return COM_STMT_FETCH
}
