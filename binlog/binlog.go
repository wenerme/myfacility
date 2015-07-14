package binlog

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/wenerme/myfacility/proto"
	"math"
	"math/big"
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
		t, meta := tab.ColumnTypes[i], tab.ColumnMetadata[i]
		_ = meta
		row[i] = readCell(proto.ColumnType(t), meta, c)
	}

	return row
}

func readCell(t proto.ColumnType, meta uint, c proto.Reader) interface{} {
	var r interface{}
	//	defer func() {
	//		log.Info("Cell %v %v %v '%v'", t, meta, l, r)
	//	}()
	switch t {
	// http://dev.mysql.com/doc/internals/en/binary-protocol-value.html
	case proto.MYSQL_TYPE_LONGLONG:
		var n uint64
		c.Get(&n)
		r = n
	case proto.MYSQL_TYPE_LONG, proto.MYSQL_TYPE_INT24:
		var n uint32
		c.Get(&n)
		r = n
	case proto.MYSQL_TYPE_SHORT, proto.MYSQL_TYPE_YEAR:
		var n uint16
		c.Get(&n)
		r = n
	case proto.MYSQL_TYPE_TINY:
		var n uint8
		c.Get(&n)
		r = n
	case proto.MYSQL_TYPE_DOUBLE:
		var n uint64
		c.Get(&n)
		r = math.Float64frombits(n)
	case proto.MYSQL_TYPE_FLOAT:
		var n uint32
		c.Get(&n)
		r = math.Float32frombits(n)
	case proto.MYSQL_TYPE_TIMESTAMP:
		var n uint32
		c.Get(&n)
		r = time.Unix(int64(n), 0)
	case proto.MYSQL_TYPE_DATETIME:
		var v uint64
		// 20060214220436
		// 8 byte
		c.Get(&v)
		p := splitDateTime(v, 100, 6)
		r = time.Date(p[5], time.Month(p[4]-1), p[3], p[2], p[1], p[0], 0, time.UTC)
	case proto.MYSQL_TYPE_DATE:
		// year 2,month 1,day 1, hour 1, minute 1, second 1, micro second 4
		var length uint8
		var year uint16
		var month, day, hour, minute, second, msecond uint8
		c.Get(&length)
		switch length {
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
			panic(fmt.Sprintf("Unkonwn type %v length %d", t, length))
		}
	case proto.MYSQL_TYPE_TIME:
		var v uint64
		c.Get(&v, proto.Int3)
		p := splitDateTime(v, 100, 3)
		r = time.Duration(p[2])*time.Hour + time.Duration(p[1])*time.Minute + time.Duration(p[0])*time.Second
	case proto.MYSQL_TYPE_VARCHAR, proto.MYSQL_TYPE_VAR_STRING:
		var length uint
		var s string
		if meta < 256 {
			c.Get(&length, proto.Int1, &s, proto.StrVar, &length)
		} else {
			c.Get(&length, proto.Int2, &s, proto.StrVar, &length)
		}
		r = s
	case proto.MYSQL_TYPE_BLOB:
		var l uint64
		var b []byte
		switch meta {
		case 1:
			c.Get(&l, proto.Int1)
		case 2:
			c.Get(&l, proto.Int2)
		case 4:
			c.Get(&l, proto.Int4)
		case 8:
			panic(fmt.Sprintf("Blob too large"))
		}
		c.Get(&b, proto.StrVar, int(l))
		r = b
	case proto.MYSQL_TYPE_NEWDECIMAL:

	default:
		panic(fmt.Sprintf("Unsupport type %s meta %d", t, meta))
	}
	return r
}

/*
   private static int[] split(long value, int divider, int length) {
       int[] result = new int[length];
       for (int i = 0; i < length - 1; i++) {
           result[i] = (int) (value % divider);
           value /= divider;
       }
       result[length - 1] = (int) value;
       return result;
   }
*/
func splitDateTime(v uint64, d int, l int) []int {
	r := make([]int, l)
	for i := 0; i < l-1; i++ {
		r[i] = int(v % uint64(d))
		v /= uint64(d)
	}
	r[l-1] = int(v)
	return r
}

var DIG_PER_DEC = 9
var DIG_TO_BYTES = []byte{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}

func determineDecimalLength(precision int, scale int) int {
	x := precision - scale
	ipDigits := x / DIG_PER_DEC
	fpDigits := scale / DIG_PER_DEC
	ipDigitsX := x - ipDigits*DIG_PER_DEC
	fpDigitsX := scale - fpDigits*DIG_PER_DEC
	return (ipDigits << 2) + int(DIG_TO_BYTES[ipDigitsX]) + (fpDigits << 2) + int(DIG_TO_BYTES[fpDigitsX])
}

// https://github.com/MariaDB/server/blob/10.1/strings/decimal.c
// decimal2bin
// bin2decimal
func toDecimal(precision int, scale int, value []byte) big.Rat {
	return big.Rat{}
}

/*
 private static BigDecimal toDecimal(int precision, int scale, byte[] value) {
        boolean positive = (value[0] & 0x80) == 0x80;
        value[0] ^= 0x80;
        if (!positive) {
            for (int i = 0; i < value.length; i++) {
                value[i] ^= 0xFF;
            }
        }
        int x = precision - scale;
        int ipDigits = x / DIG_PER_DEC;
        int ipDigitsX = x - ipDigits * DIG_PER_DEC;
        int ipSize = (ipDigits << 2) + DIG_TO_BYTES[ipDigitsX];
        int offset = DIG_TO_BYTES[ipDigitsX];
        BigDecimal ip = offset > 0 ? BigDecimal.valueOf(bigEndianInteger(value, 0, offset)) : BigDecimal.ZERO;
        for (; offset < ipSize; offset += 4) {
            int i = bigEndianInteger(value, offset, 4);
            ip = ip.movePointRight(DIG_PER_DEC).add(BigDecimal.valueOf(i));
        }
        int shift = 0;
        BigDecimal fp = BigDecimal.ZERO;
        for (; shift + DIG_PER_DEC <= scale; shift += DIG_PER_DEC, offset += 4) {
            int i = bigEndianInteger(value, offset, 4);
            fp = fp.add(BigDecimal.valueOf(i).movePointLeft(shift + DIG_PER_DEC));
        }
        if (shift < scale) {
            int i = bigEndianInteger(value, offset, DIG_TO_BYTES[scale - shift]);
            fp = fp.add(BigDecimal.valueOf(i).movePointLeft(scale));
        }
        BigDecimal result = ip.add(fp);
        return positive ? result : result.negate();
    }

*/
func bigEndianInteger(bytes []byte, offset int, length int) uint64 {
	var result uint64
	for i := offset; i < (offset + length); i++ {
		b := bytes[i]
		result = (result << 8) | uint64(b)
	}
	return result
}
