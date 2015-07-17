package binlog

import (
	"encoding/binary"
	"fmt"
	"github.com/op/go-logging"
	"github.com/wenerme/myfacility/proto"
	"os"
	"reflect"
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

func readRow(tab *TableMapEvent, included bitSet, c proto.Reader) []interface{} {
	nulls := bitSet{}
	columns := len(tab.ColumnTypes)
	c.Get(&nulls.array, proto.StrVar, (columns+7)/8)
	row := make([]interface{}, columns)
	for i := 0; i < columns; i++ {
		if !included.Has(int(i)) {
			continue
		}
		if nulls.Has(int(i)) {
			continue
		}
		// mysql-5.6.24 sql/log_event.cc log_event_print_value (line 1980)
		t, meta, l := proto.ColumnType(tab.ColumnTypes[i]), tab.ColumnMetadata[i], 0
		if t == proto.MYSQL_TYPE_STRING {
			// big endian here
			meta = (meta & 0xFF << 8) | (meta >> 8)
			if meta >= 256 {
				meta0, meta1 := uint8(meta>>8), uint8(meta)
				if (meta0 & 0x30) != 0x30 {
					t = proto.ColumnType(meta0 | 0x30)
					l = int(meta1) | (((int(meta0) & 0x30) ^ 0x30) << 4)
				} else {
					// mysql-5.6.24 sql/rpl_utility.h enum_field_types (line 278)
					mt := proto.ColumnType(meta0)
					if mt == proto.MYSQL_TYPE_ENUM || mt == proto.MYSQL_TYPE_SET {
						t = mt
					}
					l = int(meta1)
				}
			} else {
				l = int(meta)
			}
		}
		row[i] = readCell(t, meta, l, c)
	}

	return row
}

func readCell(t proto.ColumnType, meta uint, length int, c proto.Reader) interface{} {
	var r interface{}
	u, u8, u16, u32, u64 := uint(0), uint8(0), uint16(0), uint32(0), uint64(0)
	_, _, _, _, _ = u, u8, u16, u32, u64
	var b []byte
	var s string
	// http://dev.mysql.com/doc/internals/en/binary-protocol-value.html
	switch t {
	case proto.MYSQL_TYPE_BIT:
		bitSetLength := (meta>>8)*8 + (meta & 0xFF)
		c.Get(&r, proto.StrVar, (bitSetLength+7)>>3)
	case proto.MYSQL_TYPE_TINY:
		c.Get(&r, reflect.Int8)
	case proto.MYSQL_TYPE_SHORT:
		c.Get(&r, reflect.Int16)
	case proto.MYSQL_TYPE_LONG:
		c.Get(&r, reflect.Int32)
	case proto.MYSQL_TYPE_LONGLONG:
		c.Get(&r, reflect.Int64)
	case proto.MYSQL_TYPE_INT24:
		c.Get(&r, proto.Int3)
	case proto.MYSQL_TYPE_FLOAT:
		c.Get(&r, reflect.Float32)
	case proto.MYSQL_TYPE_DOUBLE:
		c.Get(&r, reflect.Float64)
	case proto.MYSQL_TYPE_NEWDECIMAL:
		precision := int(meta & 0xFF)
		scale := int(meta) >> 8
		decimalLength := determineDecimalLength(precision, scale)
		var bytes []byte
		c.Get(&bytes, proto.StrVar, decimalLength)
		r = toDecimal(precision, scale, bytes)
	case proto.MYSQL_TYPE_DATE:
		// year 2,month 1,day 1, hour 1, minute 1, second 1, micro second 4
		var l uint8
		var year uint16
		var month, day, hour, minute, second, msecond uint8
		c.Get(&l)
		switch l {
		case 0:
			r = time.Time{}
		case 4:
			// ymd
			c.Get(&year, &month, &day)
			r = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, nil)
		case 7:
			// ymd hms
			c.Get(&year, &month, &day, &hour, &minute, &second)
			r = time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, nil)
		case 11:
			c.Get(&year, &month, &day, &hour, &minute, &second, &msecond)
			// TODO not sure the msecond should time 1000
			r = time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), int(msecond)*1000, nil)
		default:
			panic(fmt.Sprintf("Unkonwn type %v length %d", t, l))
		}
	case proto.MYSQL_TYPE_TIME:
		var v uint64
		c.Get(&v, proto.Int3)
		p := splitDateTime(v, 100, 3)
		r = time.Duration(p[2])*time.Hour + time.Duration(p[1])*time.Minute + time.Duration(p[0])*time.Second
	case proto.MYSQL_TYPE_TIME2:
		/*
			in big endian:

			1 bit sign (1= non-negative, 0= negative)
			1 bit unused (reserved for future extensions)
			10 bits hour (0-838)
			6 bits minute (0-59)
			6 bits second (0-59)
			= (3 bytes in total)
			+
			fractional-seconds storage (size depends on meta)
		*/
		// big endian 3byte
		b = make([]uint8, 4)
		c.Get(b[1:])
		u64 = uint64(binary.BigEndian.Uint32(b))
		r = time.Duration(extractBits(u64, 2, 10, 24))*time.Hour*24 +
			time.Duration(extractBits(u64, 12, 6, 24))*time.Minute +
			time.Duration(extractBits(u64, 18, 6, 24))*time.Second +
			time.Duration(getFractionalSeconds(int(meta), c))*time.Millisecond
	case proto.MYSQL_TYPE_TIMESTAMP:
		c.Get(&u32)
		r = time.Unix(int64(u32), 0)
	case proto.MYSQL_TYPE_TIMESTAMP2:
		c.Get(&b, proto.StrVar, 4)
		r = time.Unix(int64(binary.BigEndian.Uint32(b)), int64(getFractionalSeconds(int(meta), c)))
	case proto.MYSQL_TYPE_DATETIME:
		// 20060214220436
		// 8 byte
		c.Get(&u64)
		p := splitDateTime(u64, 100, 6)
		r = time.Date(p[5], time.Month(p[4]-1), p[3], p[2], p[1], p[0], 0, time.UTC)
	case proto.MYSQL_TYPE_DATETIME2:
		/*
		   in big endian:

		   1 bit sign (1= non-negative, 0= negative)
		   17 bits year*13+month (year 0-9999, month 0-12)
		   5 bits day (0-31)
		   5 bits hour (0-23)
		   6 bits minute (0-59)
		   6 bits second (0-59)
		   = (5 bytes in total)
		   +
		   fractional-seconds storage (size depends on meta)
		*/
		// big endian 5byte
		b = make([]uint8, 8)
		c.Get(b[3:])
		u64 := binary.BigEndian.Uint64(b)
		yearMonth := extractBits(u64, 1, 17, 40)
		r = time.Date(yearMonth/13,
			time.Month(yearMonth%13-1),
			extractBits(u64, 18, 5, 40),
			extractBits(u64, 23, 5, 40),
			extractBits(u64, 28, 6, 40),
			extractBits(u64, 34, 6, 40),
			getFractionalSeconds(int(meta), c)*1000, time.UTC)
	case proto.MYSQL_TYPE_YEAR:
		c.Get(&u8)
		r = 1900 + int(u8)
	case proto.MYSQL_TYPE_STRING:
		if length < 256 {
			c.Get(&u, proto.Int1)
		} else {
			c.Get(&u, proto.Int2)
		}
		c.Get(&s, proto.StrVar, u)
		r = s
	case proto.MYSQL_TYPE_VARCHAR, proto.MYSQL_TYPE_VAR_STRING:
		if meta < 256 {
			c.Get(&u, proto.Int1)
		} else {
			c.Get(&u, proto.Int2)
		}
		c.Get(&s, proto.StrVar, u)
		r = s
	case proto.MYSQL_TYPE_BLOB:
		c.Get(&u, proto.Int, meta,
			&r, proto.StrVar, &u)
	case proto.MYSQL_TYPE_ENUM:
		// integer
		c.Get(&r, proto.Int, length)
	case proto.MYSQL_TYPE_SET:
		// long
		c.Get(&r, proto.Int, length)
	default:
		panic(fmt.Sprintf("Unsupport column type %v meta %d", t, meta))
	}
	return r
}

func splitDateTime(v uint64, d int, l int) []int {
	r := make([]int, l)
	for i := 0; i < l-1; i++ {
		r[i] = int(v % uint64(d))
		v /= uint64(d)
	}
	r[l-1] = int(v)
	return r
}

func extractBits(v uint64, bitOffset, numberOfBits, payloadSize uint) int {
	result := v>>payloadSize - uint64(bitOffset+numberOfBits)
	return int(result & ((1 << numberOfBits) - 1))
}

func getFractionalSeconds(meta int, c proto.Reader) (fractionalSeconds int) {
	switch meta {
	case 1:
	case 2:
		c.Get(&fractionalSeconds, proto.Int1)
	case 3:
	case 4:
		c.Get(&fractionalSeconds, proto.Int2)
	case 5:
	case 6:
		c.Get(&fractionalSeconds, proto.Int3)
	default:
		return 0
	}
	if meta%2 == 1 {
		fractionalSeconds /= 10
	}
	fractionalSeconds /= 1000
	return
}
