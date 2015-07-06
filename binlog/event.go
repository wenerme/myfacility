package binlog

import (
	"github.com/wenerme/myfacility/proto"
	"time"
)

type Event interface {
	Header() EventHeader
}

// The binlog event header starts each event and is either 13 or 19 bytes long,
// depending on the binlog version.
// http://dev.mysql.com/doc/internals/en/binlog-event-header.html
type EventHeader struct {
	Timestamp time.Time
	EventType EventType
	ServerId  uint32
	EventSize uint32
	NextPos   uint32
	Flags     uint16
}

func (p *EventHeader) Read(c proto.Reader) {
	var ts uint32
	c.Get(&ts, &p.EventType, &p.ServerId, &p.EventSize, &p.NextPos, &p.Flags)
	p.Timestamp = time.Unix(int64(ts), 0)
}

/*
The query event is used to send text querys right the binlog.

Post-header
4              slave_proxy_id
4              execution time
1              schema length
2              error-code
  if binlog-version ≥ 4:
2              status-vars length
Payload
string[$len]   status-vars
string[$len]   schema
1              [00]
string[EOF]    query
Fields
status_vars_length (2) -- number of bytes in the following sequence of status-vars

status_vars (string.var_len) -- [len=$status_vars_length] a sequence of status key-value pairs. The key is 1-byte, while its value is dependent on the key.
*/
// http://dev.mysql.com/doc/internals/en/query-event.html
type QueryEvent struct {
}
