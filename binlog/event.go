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
	Flags     BinlogEventFlag
}

func (p *EventHeader) Read(c proto.Reader) {
	var ts uint32
	c.Get(&ts, &p.EventType, &p.ServerId, &p.EventSize, &p.NextPos, &p.Flags)
	p.Timestamp = time.Unix(int64(ts), 0)
}

// The query event is used to send text querys right the binlog.
// http://dev.mysql.com/doc/internals/en/query-event.html
type QueryEvent struct {
	SlaveProxyId  uint32
	ExecutionTime uint32
	ErrorCode     uint16
	Status        []byte
	Schema        string
	Query         string
}

func (p *QueryEvent) Read(c proto.Reader) {
	var m uint8
	var n uint16
	c.Get(&p.SlaveProxyId,
		&p.ExecutionTime,
		&m,
		&p.ErrorCode,
		&n,
		&p.Status, proto.StrVar, &n,
		&p.Schema, proto.StrVar, &m,
		&p.Query)
}
func (p *QueryEvent) Type() EventType {
	return QUERY_EVENT
}
func NewEventMap() map[EventType]interface{} {
	return map[EventType]interface{}{
		UNKNOWN_EVENT:                 nil,
		START_EVENT_V3:                &StartEventV3{},
		QUERY_EVENT:                   &QueryEvent{},
		STOP_EVENT:                    StopEvent,
		ROTATE_EVENT:                  nil,
		INTVAR_EVENT:                  &IntvarEvent{},
		LOAD_EVENT:                    &LoadEvent{},
		SLAVE_EVENT:                   nil,
		CREATE_FILE_EVENT:             &CreateFileEvent{},
		APPEND_BLOCK_EVENT:            &AppendBlockEvent{},
		EXEC_LOAD_EVENT:               &ExecLoadEvent{},
		DELETE_FILE_EVENT:             &DeleteFileEvent{},
		NEW_LOAD_EVENT:                &NewLoadEvent{},
		RAND_EVENT:                    &RandEvent{},
		USER_VAR_EVENT:                &UseVarEvent{},
		FORMAT_DESCRIPTION_EVENT:      &FormatDescriptionEvent{},
		XID_EVENT:                     &XIDEvent{},
		BEGIN_LOAD_QUERY_EVENT:        &BeginLoadQueryEvent{},
		EXECUTE_LOAD_QUERY_EVENT:      &ExecuteLoadQueryEvent{},
		TABLE_MAP_EVENT:               &TableMapEvent{},
		WRITE_ROWS_EVENTv0:            nil,
		UPDATE_ROWS_EVENTv0:           nil,
		DELETE_ROWS_EVENTv0:           nil,
		WRITE_ROWS_EVENTv1:            &WriteRowsEventV1{},
		UPDATE_ROWS_EVENTv1:           nil,
		DELETE_ROWS_EVENTv1:           nil,
		INCIDENT_EVENT:                &IncidentEvent{},
		HEARTBEAT_EVENT:               HeartbeatEvent,
		IGNORABLE_EVENT:               nil,
		ROWS_QUERY_EVENT:              nil,
		WRITE_ROWS_EVENTv2:            nil,
		UPDATE_ROWS_EVENTv2:           nil,
		DELETE_ROWS_EVENTv2:           nil,
		GTID_EVENT:                    nil,
		ANONYMOUS_GTID_EVENT:          nil,
		PREVIOUS_GTIDS_EVENT:          nil,
		MARIA_ANNOTATE_ROWS_EVENT:     nil,
		MARIA_BINLOG_CHECKPOINT_EVENT: &MariaBinlogCheckPointEvent{},
		MARIA_GTID_EVENT:              &MariaGtidEvent{},
		MARIA_GTID_LIST_EVENT:         &MariaGtidListEvent{},
	}
}
