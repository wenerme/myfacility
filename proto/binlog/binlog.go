package proto
import "log"
import (
	. "../../proto"
)


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
	SchemaName    StrF
	TableName     StrF
	ColumnDef     StrL
	ColumnMetaDef StrL
	NullBitMask   StrV
}
func (p *TableMapEvent)Read(c PackReader) {
	p.SchemaName = c.MustReadStrF(uint(c.MustReadInt1()))
	c.MustReadInt1()// 00
	p.TableName = c.MustReadStrF(uint(c.MustReadInt1()))
	c.MustReadInt1()// 00
	p.ColumnDef = c.MustReadStrL()
	p.ColumnMetaDef = c.MustReadStrL()
	// Original is +8) /7
	p.NullBitMask = c.MustReadStrV(uint(len(p.ColumnDef)+7)/8)
}

func (p *TableMapEvent)Write(c PackWriter) {
	log.Panic("Unsupported")
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
// Only support version2
type RowsEvent struct {
	ColumnCount uint

}
type RowEvent struct {
	NullBitMask string
	Value       string
}