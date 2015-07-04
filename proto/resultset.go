package proto

import "github.com/davecgh/go-spew/spew"

// https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnDefinition41
type ColumnDefinition struct {
	Catalog       string
	Schema        string
	Table         string
	OrgTable      string
	Name          string
	OrgName       string
	FixedLength   uint
	CharacterSet  CharacterSet
	ColumnLength  uint32
	Type          ColumnType
	Flags         uint16
	Decimals      uint8
	DefaultValues string
}

func (p *ColumnDefinition) Read(c Reader) {
	c.Get(&p.Catalog,
		&p.Schema,
		&p.Table,
		&p.OrgTable,
		&p.Name,
		&p.OrgName,
		&p.FixedLength,
		&p.CharacterSet,
		&p.ColumnLength,
		&p.Type,
		&p.Flags,
		&p.Decimals)
	c.SkipBytes(2) // filter

	//	if c.Com ==  COM_FIELD_LIST {
	//		c.Get(&p.DefaultValues)
	//	}
}
func (p *ColumnDefinition) Write(c Writer) {
	c.Put(&p.Catalog,
		&p.Schema,
		&p.Table,
		&p.OrgTable,
		&p.Name,
		&p.OrgName,
		&p.FixedLength,
		&p.CharacterSet,
		&p.ColumnLength,
		&p.Type,
		&p.Flags,
		&p.Decimals)
	c.PutZero(2) // filter

	//	if c.Com ==  COM_FIELD_LIST {
	//		c.Put(&p.DefaultValues)
	//	}
}

// https://dev.mysql.com/doc/internals/en/com-query-response.html
type QueryResponse struct {
	Fields []ColumnDefinition
	Rows   [][]*string

	EOF *EOFPack
	OK  *OKPack
	ERR *ERRPack
}
type Cell struct {
	Value string
	Col   *ColumnDefinition
}

func (p *QueryResponse) Read(c Proto) {
	var n uint
	c.MustRecvPacket()
	c.Get(&n)
	p.Fields = make([]ColumnDefinition, n)
	for i := uint(0); i < n; i++ {
		col := ColumnDefinition{}
		_, err := c.RecvReadPacket(&col)
		if err != nil {
			panic(err)
		}
		p.Fields[i] = col
	}
	if !c.HasCap(CLIENT_DEPRECATE_EOF) {
		eof := &EOFPack{}
		c.MustRecvPacket()
		eof.Read(c)
		spew.Dump(eof)
	}

	for {
		c.MustRecvPacket()
		b, err := c.PeekByte()
		if err != nil {
			panic(err)
		}
		switch PackType(b) {
		case EOF:
			if c.HasCap(CLIENT_DEPRECATE_EOF) {
				break
			}
			eof := &EOFPack{}
			eof.Read(c)
			p.EOF = eof
			return
		case OK:
			ok := &OKPack{}
			ok.Read(c)
			p.OK = ok
			return
		case ERR:
			err := &ERRPack{}
			err.Read(c)
			p.ERR = err
			return
		}

		row := make([]*string, n)
		for i := uint(0); i < n; i++ {
			b, err := c.PeekByte()
			if err != nil {
				panic(err)
			}
			if b == 0xfb {
				c.SkipBytes(1)
				continue
			}
			var s string
			c.Get(&s)
			row[i] = &s
		}
		p.Rows = append(p.Rows, row)
	}

}
func (p *QueryResponse) Write(c Proto) {
	n := uint(len(p.Fields))
	c.Put(n)
	c.MustSendPacket()
	for _, col := range p.Fields {
		col.Write(c)
		c.MustSendPacket()
	}
	if !c.HasCap(CLIENT_DEPRECATE_EOF) {
		p.EOF.Write(c)
		c.MustSendPacket()
	}
	for _, row := range p.Rows {
		for i := uint(0); i < n; i++ {
			s := row[i]
			if s == nil {
				c.Put(uint8(0xfb))
			} else {
				c.Put(s)
			}
		}
		c.MustSendPacket()
	}
	if !c.HasCap(CLIENT_DEPRECATE_EOF) {
		p.EOF.Write(c)
		c.MustSendPacket()
	} else if p.OK != nil {
		p.OK.Write(c)
		c.MustSendPacket()
	} else if p.ERR != nil {
		p.ERR.Write(c)
		c.MustSendPacket()
	}
}
