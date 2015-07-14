package binlog

import (
	"github.com/wenerme/myfacility/proto"
)

type Reader interface {
	proto.Reader
	TableMap(uint64) *TableMapEvent
}

type reader struct {
	proto.Reader
	tables map[uint64]*TableMapEvent
}

func (c *reader) TableMap(id uint64) *TableMapEvent {
	return c.tables[id]
}
func (c *reader) SetTableMap(tab *TableMapEvent) {
	// make a copy
	t := *tab
	c.tables[t.TableId] = &t
}
