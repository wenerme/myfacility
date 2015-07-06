package binlog

import (
	"github.com/wenerme/myfacility/proto"
)

type Reader interface {
	proto.Reader
}

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

func (p *StartEventV3) Read(c Reader) {
	c.Get(&p.BinlogVersion, &p.MySQLServerVersion, proto.StrVar, 50, &p.CreateTimestamp)
}

/*



Binlog::FORMAT_DESCRIPTION_EVENT:
A format description event is the first event of a binlog for binlog-version 4. It describes how the other events are layed out.

Note
added in MySQL 5.0.0 as replacement for START_EVENT_V3

Payload
2                binlog-version
string[50]       mysql-server version
4                create timestamp
1                event header length
string[p]        event type header lengths
Fields
binlog-version (2) -- version of this binlog format.

mysql-server version (string.fix_len) -- [len=50] version of the MySQL Server that created the binlog. The string is evaluted to apply work-arounds in the slave.

create_timestamp (4) -- seconds since Unix epoch when the binlog was created

event_header_length (1) -- length of the Binlog Event Header of next events. Should always be 19.

event type header length (string.EOF) -- a array indexed by Binlog Event Type - 1 to extract the length of the event specific header.

Example
$ hexdump -v -s 4 -C relay-bin.000001
00000004  82 2d c2 4b 0f 02 00 00  00 67 00 00 00 6b 00 00  |.-.K.....g...k..|
00000014  00 00 00 04 00 35 2e 35  2e 32 2d 6d 32 00 00 00  |.....5.5.2-m2...|
00000024  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000034  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000044  00 00 00 00 00 00 00 82  2d c2 4b 13 38 0d 00 08  |........-.K.8...|
00000054  00 12 00 04 04 04 04 12  00 00 54 00 04 1a 08 00  |..........T.....|
00000064  00 00 08 08 08 02 00                              |........        |
*/
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
}
