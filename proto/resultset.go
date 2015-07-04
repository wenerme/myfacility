package proto

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

func (p *QueryResponse) Read(proto Proto) {
	var n uint
	proto.MustRecvPacket()
	proto.Get(&n)
	p.Fields = make([]ColumnDefinition, n)
	for i := uint(0); i < n; i++ {
		c := ColumnDefinition{}
		_, err := proto.RecvReadPacket(&c)
		if err != nil {
			panic(err)
		}
		p.Fields[i] = c
	}
	if !proto.HasCap(CLIENT_DEPRECATE_EOF) {
		eof := &EOFPack{}
		proto.MustRecvPacket()
		eof.Read(proto)
	}

	for {
		proto.MustRecvPacket()
		b, err := proto.PeekByte()
		if err != nil {
			panic(err)
		}
		switch PackType(b) {
		case EOF:
			if proto.HasCap(CLIENT_DEPRECATE_EOF) {
				break
			}
			eof := &EOFPack{}
			eof.Read(proto)
			p.EOF = eof
			return
		case OK:
			ok := &OKPack{}
			ok.Read(proto)
			p.OK = ok
			return
		case ERR:
			err := &ERRPack{}
			err.Read(proto)
			p.ERR = err
			return
		}

		row := make([]*string, n)
		for i := uint(0); i < n; i++ {
			b, err := proto.PeekByte()
			if err != nil {
				panic(err)
			}
			if b == 0xfb {
				proto.SkipBytes(1)
				continue
			}
			var s string
			proto.Get(&s)
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

}
