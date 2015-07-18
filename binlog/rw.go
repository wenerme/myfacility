package binlog

import (
	"bytes"
	"fmt"
	"github.com/spacemonkeygo/errors"
	"github.com/wenerme/myfacility/proto"
	"io"
)

type Reader interface {
	proto.Reader
	TableMap(uint64) *TableMapEvent
	SetTableMap(tab TableMapEvent)
	//	Next() (Event, error)
}

type Writer interface {
	proto.Writer
	TableMap(uint64) *TableMapEvent
	SetTableMap(tab TableMapEvent)
}

var ErrBinlog = errors.NewClass("ErrBinlog")

type context struct {
	tables map[uint64]*TableMapEvent
	types  map[EventType]interface{}
}
type reader struct {
	proto.Reader
	*context
	buf    *bytes.Buffer
	bufRd  proto.Reader
	rd     proto.Reader
	header EventHeader
	event  Event
}

func newContext() *context {
	return &context{types: NewEventTypeMap(), tables: make(map[uint64]*TableMapEvent)}
}
func (c *context) TableMap(id uint64) *TableMapEvent {
	return c.tables[id]
}

func (c *context) SetTableMap(tab TableMapEvent) {
	c.tables[tab.TableId] = &tab
}
func (c *reader) Next() (*Event, error) {
	if !c.rd.More() {
		return nil, io.EOF
	}
	c.event.e = nil
	c.header.Read(c.rd)
	_, err := io.CopyN(c.buf, c.rd, int64(c.header.EventSize-19))
	if err != nil {
		return nil, ErrBinlog.Wrap(err)
	}
	if c.header.EventType == TABLE_MAP_EVENT {
		tab := c.event.Event().(*TableMapEvent)
		c.SetTableMap(*tab)
	}
	return &c.event, nil
}

func NewReader(r io.Reader) *reader {
	rd := &reader{context: newContext()}
	rd.buf = &bytes.Buffer{}
	rd.bufRd = proto.NewReader(rd.buf)
	rd.event.header = &rd.header
	rd.event.getter = rd.read
	rd.Reader = rd.bufRd
	rd.rd = proto.NewReader(r)
	return rd
}
func (c *reader) read() EventData {
	p := c.types[c.header.EventType]

	if _, ok := p.(protoReadable); ok {
		p.(protoReadable).Read(c)
	} else if _, ok := p.(binlogReadable); ok {
		p.(binlogReadable).Read(c)
	} else {
		panic(fmt.Sprintf("No way to read type %v %T", c.header.EventType, p))
	}

	return p.(EventData)
}
