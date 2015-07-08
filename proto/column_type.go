package proto

import "errors"

type ColumnType uint8

// http://dev.mysql.com/doc/internals/en/date-and-time-data-type-representation.html
const (
	MYSQL_TYPE_DECIMAL ColumnType = iota
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR
	// internal
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT
	// internal
	MYSQL_TYPE_TIMESTAMP2
	// internal
	MYSQL_TYPE_DATETIME2
	// internal
	MYSQL_TYPE_TIME2
	MYSQL_TYPE_NEWDECIMAL  ColumnType = 0xf6
	MYSQL_TYPE_ENUM        ColumnType = 0xf7
	MYSQL_TYPE_SET         ColumnType = 0xf8
	MYSQL_TYPE_TINY_BLOB   ColumnType = 0xf9
	MYSQL_TYPE_MEDIUM_BLOB ColumnType = 0xfa
	MYSQL_TYPE_LONG_BLOB   ColumnType = 0xfb
	MYSQL_TYPE_BLOB        ColumnType = 0xfc
	MYSQL_TYPE_VAR_STRING  ColumnType = 0xfd
	MYSQL_TYPE_STRING      ColumnType = 0xfe
	MYSQL_TYPE_GEOMETRY    ColumnType = 0xff
)

var ErrColumnTypeNoLength = errors.New("Can not get length for column type")

func (t ColumnType) Len() int {
	switch t {
	case MYSQL_TYPE_STRING:
		return 2
	case MYSQL_TYPE_VAR_STRING:
		return 2
	case MYSQL_TYPE_VARCHAR:
		return 2
	case MYSQL_TYPE_BLOB:
		return 1
	case MYSQL_TYPE_DECIMAL:
		return 2
	case MYSQL_TYPE_NEWDECIMAL:
		return 2
	case MYSQL_TYPE_DOUBLE:
		return 1
	case MYSQL_TYPE_FLOAT:
		return 1
	case MYSQL_TYPE_ENUM:
		return 2
	case MYSQL_TYPE_SET:
		return MYSQL_TYPE_ENUM.Len()
	case MYSQL_TYPE_BIT:
		return 0
	case MYSQL_TYPE_DATE:
		return 0
	case MYSQL_TYPE_DATETIME:
		return 0
	case MYSQL_TYPE_TIMESTAMP:
		return 0
	case MYSQL_TYPE_TINY:
		return 0
	case MYSQL_TYPE_SHORT:
		return 0
	case MYSQL_TYPE_INT24:
		return 0
	case MYSQL_TYPE_LONG:
		return 0
	case MYSQL_TYPE_LONGLONG:
		return 0
	}
	panic(ErrColumnTypeNoLength)
}
