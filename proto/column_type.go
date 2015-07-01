package proto
import "errors"

type ColumnType uint8
const (
	MYSQL_TYPE_DECIMAL           ColumnType = 0x00
	MYSQL_TYPE_TINY           ColumnType = 0x01
	MYSQL_TYPE_SHORT           ColumnType = 0x02
	MYSQL_TYPE_LONG           ColumnType = 0x03
	MYSQL_TYPE_FLOAT           ColumnType = 0x04
	MYSQL_TYPE_DOUBLE           ColumnType = 0x05
	MYSQL_TYPE_NULL           ColumnType = 0x06
	MYSQL_TYPE_TIMESTAMP       ColumnType = 0x07
	MYSQL_TYPE_LONGLONG       ColumnType = 0x08
	MYSQL_TYPE_INT24           ColumnType = 0x09
	MYSQL_TYPE_DATE           ColumnType = 0x0a
	MYSQL_TYPE_TIME           ColumnType = 0x0b
	MYSQL_TYPE_DATETIME       ColumnType = 0x0c
	MYSQL_TYPE_YEAR           ColumnType = 0x0d
// internal
	MYSQL_TYPE_NEWDATE           ColumnType = 0x0e
	MYSQL_TYPE_VARCHAR           ColumnType = 0x0f
	MYSQL_TYPE_BIT               ColumnType = 0x10
// internal
	MYSQL_TYPE_TIMESTAMP2       ColumnType = 0x11
// internal
	MYSQL_TYPE_DATETIME2       ColumnType = 0x12
// internal
	MYSQL_TYPE_TIME2           ColumnType = 0x13
	MYSQL_TYPE_NEWDECIMAL       ColumnType = 0xf6
	MYSQL_TYPE_ENUM           ColumnType = 0xf7
	MYSQL_TYPE_SET            ColumnType = 0xf8
	MYSQL_TYPE_TINY_BLOB    ColumnType = 0xf9
	MYSQL_TYPE_MEDIUM_BLOB    ColumnType = 0xfa
	MYSQL_TYPE_LONG_BLOB    ColumnType = 0xfb
	MYSQL_TYPE_BLOB        ColumnType = 0xfc
	MYSQL_TYPE_VAR_STRING    ColumnType = 0xfd
	MYSQL_TYPE_STRING        ColumnType = 0xfe
	MYSQL_TYPE_GEOMETRY    ColumnType = 0xff
)

var ErrColumnTypeNoLength = errors.New("Can not get length for column type");

func (t ColumnType)Len() int {
	switch t{
	case MYSQL_TYPE_STRING: return 2
	case MYSQL_TYPE_VAR_STRING: return 2
	case MYSQL_TYPE_VARCHAR: return 2
	case MYSQL_TYPE_BLOB: return 1
	case MYSQL_TYPE_DECIMAL: return 2
	case MYSQL_TYPE_NEWDECIMAL: return 2
	case MYSQL_TYPE_DOUBLE: return 1
	case MYSQL_TYPE_FLOAT: return 1
	case MYSQL_TYPE_ENUM: return 2
	case MYSQL_TYPE_SET: return MYSQL_TYPE_ENUM.Len()
	case MYSQL_TYPE_BIT: return 0
	case MYSQL_TYPE_DATE: return 0
	case MYSQL_TYPE_DATETIME: return 0
	case MYSQL_TYPE_TIMESTAMP: return 0
	case MYSQL_TYPE_TINY: return 0
	case MYSQL_TYPE_SHORT: return 0
	case MYSQL_TYPE_INT24: return 0
	case MYSQL_TYPE_LONG: return 0
	case MYSQL_TYPE_LONGLONG: return 0
	}
	panic(ErrColumnTypeNoLength)
}