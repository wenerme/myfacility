package binlog

//go:vet
//go:generate stringer -output=strings.go -type=EventType,BinlogEventFlag,MariaGtidEventFlag,RowsEventFlag

// http://dev.mysql.com/doc/internals/en/binlog-event-type.html
type EventType uint8

const (
	UNKNOWN_EVENT EventType = iota
	START_EVENT_V3
	QUERY_EVENT
	STOP_EVENT
	ROTATE_EVENT
	INTVAR_EVENT
	LOAD_EVENT
	SLAVE_EVENT
	// http://dev.mysql.com/doc/internals/en/create-file-event.html
	CREATE_FILE_EVENT
	// http://dev.mysql.com/doc/internals/en/append-block-event.html
	APPEND_BLOCK_EVENT
	EXEC_LOAD_EVENT
	DELETE_FILE_EVENT
	NEW_LOAD_EVENT
	RAND_EVENT
	USER_VAR_EVENT
	FORMAT_DESCRIPTION_EVENT
	XID_EVENT
	// http://dev.mysql.com/doc/internals/en/begin-load-query-event.html
	BEGIN_LOAD_QUERY_EVENT
	// http://dev.mysql.com/doc/internals/en/execute-load-query-event.html
	EXECUTE_LOAD_QUERY_EVENT
	TABLE_MAP_EVENT
	WRITE_ROWS_EVENTv0
	UPDATE_ROWS_EVENTv0
	DELETE_ROWS_EVENTv0
	WRITE_ROWS_EVENTv1
	UPDATE_ROWS_EVENTv1
	DELETE_ROWS_EVENTv1
	INCIDENT_EVENT
	HEARTBEAT_EVENT
	IGNORABLE_EVENT
	ROWS_QUERY_EVENT
	WRITE_ROWS_EVENTv2
	UPDATE_ROWS_EVENTv2
	DELETE_ROWS_EVENTv2
	GTID_EVENT
	ANONYMOUS_GTID_EVENT
	PREVIOUS_GTIDS_EVENT
)
const (
	MARIA_ANNOTATE_ROWS_EVENT EventType = iota + 160
	// Binlog checkpoint event. Used for XA crash recovery on the master, not used
	// in replication.
	// A binlog checkpoint event specifies a binlog file such that XA crash
	// recovery can start from that file - and it is guaranteed to find all XIDs
	// that are prepared in storage engines but not yet committed.
	MARIA_BINLOG_CHECKPOINT_EVENT
	// Gtid event. For global transaction ID, used to start a new event group,
	// instead of the old BEGIN query event, and also to mark stand-alone
	// events.
	MARIA_GTID_EVENT
	// Gtid list event. Logged at the start of every binlog, to record the
	// current replication state. This consists of the last GTID seen for
	// each replication domain.
	MARIA_GTID_LIST_EVENT
)

// IntVar type
const (
	UNSIGNED uint8 = 0x1
)

// incident-event type
// http://dev.mysql.com/doc/internals/en/incident-event.html
const (
	INCIDENT_NONE uint16 = iota
	INCIDENT_LOST_EVENTS
)

// http://dev.mysql.com/doc/internals/en/intvar-event.html
type IntVarType uint8

const (
	INVALID_INT_EVENT IntVarType = iota
	LAST_INSERT_ID_EVENT
	INSERT_ID_EVENT
)

// http://dev.mysql.com/doc/internals/en/binlog-event-flag.html
type BinlogEventFlag uint16

const (
	//gets unset in the FORMAT_DESCRIPTION_EVENT when the file gets closed to detect broken binlogs
	LOG_EVENT_BINLOG_IN_USE_F BinlogEventFlag = 1 << iota
	//unused
	LOG_EVENT_FORCED_ROTATE_F
	//event is thread specific (CREATE TEMPORARY TABLE ...)
	LOG_EVENT_THREAD_SPECIFIC_F
	//event doesn't need default database to be updated (CREATE DATABASE, ...)
	LOG_EVENT_SUPPRESS_USE_FBinlogEventFlag
	//unused
	LOG_EVENT_UPDATE_TABLE_MAP_VERSION_F
	//event is created by the slaves SQL-thread and shouldn't update the master-log pos
	LOG_EVENT_ARTIFICIAL_F
	//event is created by the slaves IO-thread when written to the relay log
	LOG_EVENT_RELAY_LOG_F
)

type MariaGtidEventFlag uint8

const (
	//FL_STANDALONE is set when there is no terminating COMMIT event.
	FL_STANDALONE MariaGtidEventFlag = 1 << iota
	//FL_GROUP_COMMIT_ID is set when event group is part of a group commit on the
	//master. Groups with same commit_id are part of the same group commit.
	FL_GROUP_COMMIT_ID
	//FL_TRANSACTIONAL is set for an event group that can be safely rolled back
	//(no MyISAM, eg.).
	FL_TRANSACTIONAL
	//FL_ALLOW_PARALLEL reflects the (negation of the) value of @@SESSION.skip_parallel_replication at the time of commit.
	FL_ALLOW_PARALLEL
	//FL_WAITED is set if a row lock wait (or other wait) is detected during the
	//execution of the transaction.
	FL_WAITED
	//FL_DDL is set for event group containing DDL.
	FL_DDL
)

// http://dev.mysql.com/doc/internals/en/rows-event.htm
type RowsEventFlag uint16

const (
	END_OF_STATEMENT RowsEventFlag = 1 << iota
	NO_FOREIGN_KEY_CHECKS
	NO_UNIQUE_KEY_CHECKS
	ROW_HAS_A_COLUMNS
)
