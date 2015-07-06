// generated by stringer -output=strings.go -type=EventType; DO NOT EDIT

package binlog

import "fmt"

const _EventType_name = "UNKNOWN_EVENTSTART_EVENT_V3QUERY_EVENTSTOP_EVENTROTATE_EVENTINTVAR_EVENTLOAD_EVENTSLAVE_EVENTCREATE_FILE_EVENTAPPEND_BLOCK_EVENTEXEC_LOAD_EVENTDELETE_FILE_EVENTNEW_LOAD_EVENTRAND_EVENTUSER_VAR_EVENTFORMAT_DESCRIPTION_EVENTXID_EVENTBEGIN_LOAD_QUERY_EVENTEXECUTE_LOAD_QUERY_EVENTTABLE_MAP_EVENTWRITE_ROWS_EVENTv0UPDATE_ROWS_EVENTv0DELETE_ROWS_EVENTv0WRITE_ROWS_EVENTv1UPDATE_ROWS_EVENTv1DELETE_ROWS_EVENTv1INCIDENT_EVENTHEARTBEAT_EVENTIGNORABLE_EVENTROWS_QUERY_EVENTWRITE_ROWS_EVENTv2UPDATE_ROWS_EVENTv2DELETE_ROWS_EVENTv2GTID_EVENTANONYMOUS_GTID_EVENTPREVIOUS_GTIDS_EVENT"

var _EventType_index = [...]uint16{0, 13, 27, 38, 48, 60, 72, 82, 93, 110, 128, 143, 160, 174, 184, 198, 222, 231, 253, 277, 292, 310, 329, 348, 366, 385, 404, 418, 433, 448, 464, 482, 501, 520, 530, 550, 570}

func (i EventType) String() string {
	if i >= EventType(len(_EventType_index)-1) {
		return fmt.Sprintf("EventType(%d)", i)
	}
	return _EventType_name[_EventType_index[i]:_EventType_index[i+1]]
}
