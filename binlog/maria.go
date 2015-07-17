package binlog

import (
	"fmt"
	"github.com/wenerme/myfacility/proto"
)

type MariaGtidEvent struct {
	SequenceNumber uint64
	DomainId       uint32
	Flag           MariaGtidEventFlag
}

func (p *MariaGtidEvent) Read(c proto.Reader) {
	c.Get(&p.SequenceNumber, &p.DomainId, &p.Flag)
	// reserved
	n := 6
	if p.Flag&FL_GROUP_COMMIT_ID > 0 {
		n += 2
	}
	c.Get(n, proto.IgnoreByte)
}
func (p *MariaGtidEvent) Type() EventType {
	return MARIA_GTID_EVENT
}

type MariaBinlogCheckPointEvent struct {
	BinlogFilename string
}

func (p *MariaBinlogCheckPointEvent) Read(c proto.Reader) {
	var n uint32
	c.Get(&n, &p.BinlogFilename, proto.StrVar, &n)
}
func (p *MariaBinlogCheckPointEvent) Type() EventType {
	return MARIA_BINLOG_CHECKPOINT_EVENT
}

type MariaGtidListEvent struct {
	Flag uint8
	List []MariaGtid
}

func (p *MariaGtidListEvent) Read(c proto.Reader) {
	var n uint32
	c.Get(&n)
	p.Flag = uint8(n >> 28) // higher 4 bit
	n = n & 0x1fffffff      // lower 28 bit
	for i := uint32(0); i < n; i++ {
		g := MariaGtid{}
		c.Get(&g.DomainId, &g.ServerId, &g.SequenceNumber)
		p.List = append(p.List, g)
	}
}
func (p *MariaGtidListEvent) Type() EventType {
	return MARIA_GTID_LIST_EVENT
}

type MariaGtid struct {
	DomainId       uint32
	ServerId       uint32
	SequenceNumber uint64
}

func (g MariaGtid) String() string {
	return fmt.Sprintf("%d-%d-%d", g.DomainId, g.ServerId, g.SequenceNumber)
}
