package binlog

import (
	"github.com/wenerme/myfacility/proto"
	"strings"
)

var StopEvent = eventTypePack(STOP_EVENT)
var HeartbeatEvent = eventTypePack(HEARTBEAT_EVENT)

// A start event is the first event of a binlog for binlog-version 1 to 3.
// http://dev.mysql.com/doc/internals/en/start-event-v3.html
type StartEventV3 struct {
	// version of this binlog format.
	BinlogVersion uint16
	//  [len=50] version of the MySQL Server that created the binlog. The string is evaluted to apply work-arounds in the slave.
	MySQLServerVersion string
	// seconds since Unix epoch when the binlog was created
	CreateTimestamp uint32
}

func (p *StartEventV3) Read(c proto.Reader) {
	c.Get(&p.BinlogVersion, &p.MySQLServerVersion, proto.StrVar, 50, &p.CreateTimestamp)
}
func (p *StartEventV3) Type() EventType {
	return START_EVENT_V3
}

// A format description event is the first event of a binlog for binlog-version 4. It describes how the other events are layed out.
// Note added in MySQL 5.0.0 as replacement for START_EVENT_V3
// http://dev.mysql.com/doc/internals/en/format-description-event.html
type FormatDescriptionEvent struct {
	// version of this binlog format.
	BinlogVersion uint16
	// [len=50] version of the MySQL Server that created the binlog. The string is evaluted to apply work-arounds in the slave.
	MySQLServerVersion string
	// seconds since Unix epoch when the binlog was created
	CreateTimestamp uint32
	// length of the Binlog Event Header of next events. Should always be 19.
	EventHeaderLength uint8
	// a array indexed by Binlog Event Type - 1 to extract the length of the event specific header.
	EventTypeHeader []byte
}

func (p *FormatDescriptionEvent) Read(c proto.Reader) {
	c.Get(&p.BinlogVersion,
		&p.MySQLServerVersion, proto.StrVar, 50,
		&p.CreateTimestamp,
		&p.EventHeaderLength, // Should always be 19.
		&p.EventTypeHeader, proto.StrEof,
	)
	p.MySQLServerVersion = strings.Trim(p.MySQLServerVersion, "\000")
}
func (p *FormatDescriptionEvent) Type() EventType {
	return FORMAT_DESCRIPTION_EVENT
}

type IncidentEvent struct {
	IncidentType uint8
	Message      []byte
}

func (p *IncidentEvent) Read(c proto.Reader) {
	var n uint8
	c.Get(&p.IncidentType, &n, &p.Message, proto.StrVar, &n)
}
func (p *IncidentEvent) Type() EventType {
	return INCIDENT_EVENT
}

type RotateEvent struct {
	BinlogPos  uint64
	BinlogFile string
}

func (p *RotateEvent) Read(c proto.Reader) {
	c.Get(&p.BinlogPos, &p.BinlogFile, proto.StrEof)
}
func (p *RotateEvent) Type() EventType {
	return ROTATE_EVENT
}
