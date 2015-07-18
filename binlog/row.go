package binlog

/*
Row Based Replication Events
In Row Based replication the changed rows are sent to the slave which removes side-effects and makes it more reliable.
Now all statements can be sent with RBR though.
Most of the time you will see RBR and SBR side by side.
*/
import "github.com/wenerme/myfacility/proto"

/*
post-header:
    if post_header_len == 6 {
  4              table id
    } else {
  6              table id
    }
  2              flags
payload:
  1              schema name length
  string         schema name
  1              [00]
  1              table name length
  string         table name
  1              [00]
  lenenc-int     column-count
  string.var_len [length=$column-count] column-def
  lenenc-str     column-meta-def
  n              NULL-bitmask, length: (column-count + 8) / 7
*/
// http://dev.mysql.com/doc/internals/en/table-map-event.html
type TableMapEvent struct {
	TableId        uint64
	Flag           uint16
	SchemaName     string
	TableName      string
	ColumnTypes    []byte
	ColumnMetadata []uint
	NullBitMask    []byte
}

func (p *TableMapEvent) EventType() EventType {
	return TABLE_MAP_EVENT
}
func (p *TableMapEvent) Read(c proto.Reader) {
	var n uint8
	var columns uint
	// TODO is ok to ignore post_header_len ?
	// Note Can use StrNul ?
	c.Get(&p.TableId, proto.Int6,
		&p.Flag,
		&n, &p.SchemaName, proto.StrVar, &n,
		1, proto.IgnoreByte,
		&n, &p.TableName, proto.StrVar, &n,
		1, proto.IgnoreByte,
		&columns, &p.ColumnTypes, proto.StrVar, &columns)
	//	&p.ColumnMetadata
	var l uint
	c.Get(&l)
	p.ColumnMetadata = make([]uint, len(p.ColumnTypes))
	for i, t := range p.ColumnTypes {
		switch proto.ColumnType(t) {
		case proto.MYSQL_TYPE_STRING, proto.MYSQL_TYPE_VAR_STRING, proto.MYSQL_TYPE_VARCHAR:
			c.Get(&p.ColumnMetadata[i], proto.Int2)
		case proto.MYSQL_TYPE_BLOB, proto.MYSQL_TYPE_DOUBLE, proto.MYSQL_TYPE_FLOAT:
			c.Get(&p.ColumnMetadata[i], proto.Int1)
		case proto.MYSQL_TYPE_TIMESTAMP2, proto.MYSQL_TYPE_DATETIME2, proto.MYSQL_TYPE_TIME2:
			c.Get(&p.ColumnMetadata[i], proto.Int1)
		case proto.MYSQL_TYPE_DECIMAL, proto.MYSQL_TYPE_NEWDECIMAL, proto.MYSQL_TYPE_SET, proto.MYSQL_TYPE_ENUM:
			c.Get(&p.ColumnMetadata[i], proto.Int2)
		/*
			proto.MYSQL_TYPE_TIME
			proto.MYSQL_TYPE_BIT
			proto.MYSQL_TYPE_DATE
			proto.MYSQL_TYPE_DATETIME
			proto.MYSQL_TYPE_TIMESTAMP
			proto.MYSQL_TYPE_TINY
			proto.MYSQL_TYPE_SHORT
			proto.MYSQL_TYPE_INT24
			proto.MYSQL_TYPE_LONG
			proto.MYSQL_TYPE_LONGLONG
			0
		*/
		default:
			p.ColumnMetadata[i] = 0
		}
	}
	c.Get(&p.NullBitMask, proto.StrVar, (columns+7)/8)
}
func (p *TableMapEvent) Write(c proto.Writer) {
	var n uint8 = uint8(len(p.SchemaName))
	var m uint8 = uint8(len(p.TableName))
	var columns uint = uint(len(p.ColumnTypes))
	// TODO is ok to ignore post_header_len ?
	// Note Can use StrNul ?
	c.Put(&p.TableId, proto.Int6,
		&p.Flag,
		&n, &p.SchemaName, proto.StrVar, &n,
		1, proto.IgnoreByte,
		&m, &p.TableName, proto.StrVar, &m,
		1, proto.IgnoreByte,
		&columns, &p.ColumnTypes, proto.StrVar, &columns)
	//	&p.ColumnMetadata
	var l uint
	for _, t := range p.ColumnTypes {
		switch proto.ColumnType(t) {
		case proto.MYSQL_TYPE_STRING, proto.MYSQL_TYPE_VAR_STRING, proto.MYSQL_TYPE_VARCHAR:
			l += 2
		case proto.MYSQL_TYPE_BLOB, proto.MYSQL_TYPE_DOUBLE, proto.MYSQL_TYPE_FLOAT:
			l += 1
		case proto.MYSQL_TYPE_TIMESTAMP2, proto.MYSQL_TYPE_DATETIME2, proto.MYSQL_TYPE_TIME2:
			l += 1
		case proto.MYSQL_TYPE_DECIMAL, proto.MYSQL_TYPE_NEWDECIMAL, proto.MYSQL_TYPE_SET, proto.MYSQL_TYPE_ENUM:
			l += 2
		}
	}
	c.Put(l)
	for i, t := range p.ColumnTypes {
		switch proto.ColumnType(t) {
		case proto.MYSQL_TYPE_STRING, proto.MYSQL_TYPE_VAR_STRING, proto.MYSQL_TYPE_VARCHAR:
			c.Put(p.ColumnMetadata[i], proto.Int2)
		case proto.MYSQL_TYPE_BLOB, proto.MYSQL_TYPE_DOUBLE, proto.MYSQL_TYPE_FLOAT:
			c.Put(p.ColumnMetadata[i], proto.Int1)
		case proto.MYSQL_TYPE_TIMESTAMP2, proto.MYSQL_TYPE_DATETIME2, proto.MYSQL_TYPE_TIME2:
			c.Put(p.ColumnMetadata[i], proto.Int1)
		case proto.MYSQL_TYPE_DECIMAL, proto.MYSQL_TYPE_NEWDECIMAL, proto.MYSQL_TYPE_SET, proto.MYSQL_TYPE_ENUM:
			c.Put(p.ColumnMetadata[i], proto.Int2)
		/*
			proto.MYSQL_TYPE_TIME
			proto.MYSQL_TYPE_BIT
			proto.MYSQL_TYPE_DATE
			proto.MYSQL_TYPE_DATETIME
			proto.MYSQL_TYPE_TIMESTAMP
			proto.MYSQL_TYPE_TINY
			proto.MYSQL_TYPE_SHORT
			proto.MYSQL_TYPE_INT24
			proto.MYSQL_TYPE_LONG
			proto.MYSQL_TYPE_LONGLONG
			0
		*/
		default:
		}
	}
	c.Put(&p.NullBitMask, proto.StrVar, (columns+7)/8)
}

/*
header:
  if post_header_len == 6 {
4                    table id
  } else {
6                    table id
  }
2                    flags
  if version == 2 {
2                    extra-data-length
string.var_len       extra-data
  }

body:
lenenc_int           number of columns
string.var_len       columns-present-bitmap1, length: (num of columns+7)/8
  if UPDATE_ROWS_EVENTv1 or v2 {
string.var_len       columns-present-bitmap2, length: (num of columns+7)/8
  }

rows:
string.var_len       nul-bitmap, length (bits set in 'columns-present-bitmap1'+7)/8
string.var_len       value of each field as defined in table-map
  if UPDATE_ROWS_EVENTv1 or v2 {
string.var_len       nul-bitmap, length (bits set in 'columns-present-bitmap2'+7)/8
string.var_len       value of each field as defined in table-map
  }
  ... repeat rows until event-end
*/

type RowsEvent struct {
	TableId     uint64
	Flag        uint16
	ExtraData   []byte
	ColumnCount uint
	// Column included
	BeforeColumns []byte
	// Column included
	AfterColumns []byte
	Before       [][]interface{}
	After        [][]interface{}
}

type WriteRowsEventV1 RowsEvent

func (p *WriteRowsEventV1) EventType() EventType {
	return WRITE_ROWS_EVENTv1
}
func (p *WriteRowsEventV1) Read(c Reader) {
	p.Before = nil
	p.After = nil
	p.ExtraData = nil
	p.AfterColumns = nil
	p.BeforeColumns = nil
	// TODO is ok to ignore post_header_len ?
	c.Get(&p.TableId, proto.Int6, &p.Flag, &p.ColumnCount)
	c.Get(&p.AfterColumns, proto.StrVar, (p.ColumnCount+7)/8)

	tab := c.TableMap(p.TableId)
	included := bitSet{int(p.ColumnCount), p.AfterColumns}
	for {
		p.After = append(p.After, readRow(tab, included, c))
		if !c.More() {
			break
		}
	}
}

type DeleteRowsEventV1 RowsEvent

func (p *DeleteRowsEventV1) EventType() EventType {
	return DELETE_ROWS_EVENTv1
}
func (p *DeleteRowsEventV1) Read(c Reader) {
	p.Before = nil
	p.After = nil
	p.ExtraData = nil
	p.AfterColumns = nil
	p.BeforeColumns = nil
	// TODO is ok to ignore post_header_len ?
	c.Get(&p.TableId, proto.Int6, &p.Flag, &p.ColumnCount)
	c.Get(&p.BeforeColumns, proto.StrVar, (p.ColumnCount+7)/8)

	tab := c.TableMap(p.TableId)
	included := bitSet{int(p.ColumnCount), p.BeforeColumns}
	for {
		p.Before = append(p.Before, readRow(tab, included, c))
		if !c.More() {
			break
		}
	}
}

type UpdateRowsEventV1 RowsEvent

func (p *UpdateRowsEventV1) EventType() EventType {
	return UPDATE_ROWS_EVENTv1
}
func (p *UpdateRowsEventV1) Read(c proto.Reader) {
	p.Before = nil
	p.After = nil
	p.ExtraData = nil
	// TODO is ok to ignore post_header_len ?
	c.Get(&p.TableId, proto.Int6, &p.Flag)
	var n uint16
	c.Get(&n, &p.ExtraData, proto.StrVar, &n,
		&p.ColumnCount,
		&p.AfterColumns,
	)
}

const (
	RW_V_EXTRAINFO_TAG = 0
)

type RowsEventExtraData struct {
	EventType uint8
	Data      []byte
}

type RowsQueryEvent struct {
	Query string
}

func (p *RowsQueryEvent) Read(c proto.Reader) {
	c.Get(1, proto.IgnoreByte, &p.Query, proto.StrEof)
}
func (p *RowsQueryEvent) EventType() EventType {
	return ROWS_QUERY_EVENT
}
