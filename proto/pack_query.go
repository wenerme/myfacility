package proto

// lenenc_str     catalog
// lenenc_str     schema
// lenenc_str     table
// lenenc_str     org_table
// lenenc_str     name
// lenenc_str     org_name
// lenenc_int     length of fixed-length fields [0c]
// 2              character set
// 4              column length
// 1              type
// 2              flags
// 1              decimals
// 2              filler [00] [00]
// if command was COM_FIELD_LIST {
// lenenc_int     length of default-values
// string[$len]   default values
// }
// https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnDefinition41
type ColumnDefinition struct {
	Catalog       string
	Schema        string
	Table         string
	OrgTable      string
	Name          string
	OrgName       string
	FixedLength   uint
	CharacterSet  uint16
	ColumnLength  uint32
	Type          uint8
	Flags         uint16
	Decimals      uint8
	DefaultValues string
}
func (p *ColumnDefinition)Read(c *Reader) {
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
	c.SkipBytes(2)// filter

	if c.Com ==  COM_FIELD_LIST {
		c.Get(&p.DefaultValues)
	}
}
func (p *ColumnDefinition)Write(c *Writer) {
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
	c.PutZero(2)// filter

	if c.Com ==  COM_FIELD_LIST {
		c.Put(&p.DefaultValues)
	}
}


// https://dev.mysql.com/doc/internals/en/com-query-response.html
type QueryResponse struct {
	Fields []ColumnDefinition
	Rows   [][]*string
}
