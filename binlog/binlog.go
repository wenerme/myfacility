package binlog

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/wenerme/myfacility/proto"
	"math"
	"os"
	"time"
)

var log = logging.MustGetLogger("binlog")

// 初始化 Log
func init() {
	//	format := logging.MustStringFormatter("%{color}%{time:15:04:05} %{level:.4s} %{shortfunc} %{color:reset} %{message}", )
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{level:.4s} %{longfile} %{shortfunc} %{color:reset} %{message}")
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
	logging.SetLevel(logging.DEBUG, "binlog")
}

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
	for i, t := range p.ColumnTypes {
		p.ColumnMetadata = make([]uint, len(p.ColumnTypes))
		switch ColumnType(t) {
		case MYSQL_TYPE_STRING, MYSQL_TYPE_VAR_STRING, MYSQL_TYPE_VARCHAR:
			c.Get(&p.ColumnMetadata[i], Int2)
		case MYSQL_TYPE_BLOB, MYSQL_TYPE_DOUBLE, MYSQL_TYPE_FLOAT:
			c.Get(&p.ColumnMetadata[i], Int1)
		case MYSQL_TYPE_TIMESTAMP2, MYSQL_TYPE_DATETIME2, MYSQL_TYPE_TIME2:
			c.Get(&p.ColumnMetadata[i], Int1)
		case MYSQL_TYPE_DECIMAL, MYSQL_TYPE_NEWDECIMAL, MYSQL_TYPE_SET, MYSQL_TYPE_ENUM:
			c.Get(&p.ColumnMetadata[i], Int2)
		/*
			MYSQL_TYPE_TIME
			MYSQL_TYPE_BIT
			MYSQL_TYPE_DATE
			MYSQL_TYPE_DATETIME
			MYSQL_TYPE_TIMESTAMP
			MYSQL_TYPE_TINY
			MYSQL_TYPE_SHORT
			MYSQL_TYPE_INT24
			MYSQL_TYPE_LONG
			MYSQL_TYPE_LONGLONG
			0
		*/
		default:
			p.ColumnMetadata[i] = 0
		}
	}
	c.Get(&p.NullBitMask, proto.StrVar, (columns+7)/8)
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
// http://dev.mysql.com/doc/internals/en/rows-event.htm
type RowsEventFlag uint16

const (
	END_OF_STATEMENT RowsEventFlag = 1 << iota
	NO_FOREIGN_KEY_CHECKS
	NO_UNIQUE_KEY_CHECKS
	ROW_HAS_A_COLUMNS
)

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
type UpdateRowsEventV1 RowsEvent
type BinlogReader interface {
	proto.Reader
	TableMap(uint64) TableMapEvent
}

func (p *WriteRowsEventV1) Read(c BinlogReader) {
	p.Before = nil
	p.After = nil
	p.ExtraData = nil
	p.AfterColumns = nil
	p.BeforeColumns = nil
	// TODO is ok to ignore post_header_len ?
	c.Get(&p.TableId, proto.Int6, &p.Flag, &p.ColumnCount)
	c.Get(&p.AfterColumns, proto.StrVar, (p.ColumnCount+7)/8)

	tab := c.TableMap(p.TableId)
	nulls := bitSet{}
	included := bitSet{int(p.ColumnCount), p.AfterColumns}
	for {
		c.Get(&nulls, proto.StrVar, (p.ColumnCount+7)/8)
		row := make([]interface{}, p.ColumnCount)
		for i := uint(0); i < p.ColumnCount; i++ {
			if !included.Has(int(i)) {
				continue
			}
			if nulls.Has(int(i)) {
				continue
			}
			// mysql-5.6.24 sql/log_event.cc log_event_print_value (line 1980)
			typeCode, meta := tab.ColumnTypes[i], tab.ColumnMetadata[i]
			switch proto.ColumnType(typeCode) {
			// http://dev.mysql.com/doc/internals/en/binary-protocol-value.html
			case proto.MYSQL_TYPE_LONGLONG:
				var n uint64
				c.Get(&n)
				row[i] = n
			case proto.MYSQL_TYPE_LONG, proto.MYSQL_TYPE_INT24:
				var n uint32
				c.Get(&n)
				row[i] = n
			case proto.MYSQL_TYPE_SHORT, proto.MYSQL_TYPE_YEAR:
				var n uint16
				c.Get(&n)
				row[i] = n
			case proto.MYSQL_TYPE_TINY:
				var n uint8
				c.Get(&n)
				row[i] = n
			case proto.MYSQL_TYPE_DOUBLE:
				var n uint64
				c.Get(&n)
				row[i] = math.Float64frombits(n)
			case proto.MYSQL_TYPE_FLOAT:
				var n uint32
				c.Get(&n)
				row[i] = math.Float32frombits(n)
			case proto.MYSQL_TYPE_DATE, proto.MYSQL_TYPE_DATETIME, proto.MYSQL_TYPE_TIMESTAMP:
				// year 2,month 1,day 1, hour 1, minute 1, second 1, micro second 4
				var length uint8
				var year uint16
				var month, day, hour, minute, second, msecond uint8
				c.Get(&length)
				switch length {
				case 0:
					row[i] = time.Time{}
				case 4:
					// ymd
					c.Get(&year, &month, &day)
					row[i] = time.Date(year, month, day, 0, 0, 0, 0, nil)
				case 7:
					// ymd hms
					c.Get(&year, &month, &day, &hour, &minute, &second)
					row[i] = time.Date(year, month, day, hour, minute, second, 0, nil)
				case 11:
					c.Get(&year, &month, &day, &hour, &minute, &second, &msecond)
					// TODO not sure the msecond should time 1000
					row[i] = time.Date(year, month, day, hour, minute, second, msecond*1000, nil)
				default:
					panic(fmt.Sprintf("Unkonwn type length"))
				}
			case proto.MYSQL_TYPE_TIME:
				var length uint8
				var isNeg, day, hour, minute, second, msecond uint8
				c.Get(&length)
				switch length {
				case 0, 1:
					row[i] = time.Duration(0)
				case 8:
					c.Get(&isNeg, &day, &hour, &minute, &second)
					row[i] = time.Date(0, 0, &day, hour, minute, second, 0, nil)
					if isNeg == 1 {
						day = -day
					}
					row[i] = time.Duration(int64(day)*24+int64(time.Hour))*time.Hour +
						time.Duration(minute)*time.Minute +
						time.Duration(second)*time.Second
				case 12:
					c.Get(&isNeg, &day, &hour, &minute, &second, &msecond)
					row[i] = time.Duration(int64(day)*24+int64(time.Hour))*time.Hour +
						time.Duration(minute)*time.Minute +
						time.Duration(second)*time.Second +
						time.Duration(msecond)*time.Microsecond
				default:
					panic(fmt.Sprintf("Unkonwn time length"))
				}
			case proto.MYSQL_TYPE_STRING:

			}
		}
		if !c.More() {
			break
		}
	}
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

type bitSet struct {
	len   int
	array []byte
}

func (b *bitSet) Has(pos int) bool {
	n := pos / 8
	return b.array[n]&(1<<uint8(pos%8)) > 0
}
