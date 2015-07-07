package binlog

import (
	"fmt"
	"github.com/wenerme/myfacility/proto"
)

//go:vet
//go:generate stringer -output=strings.go -type=EventType,BinlogEventFlag,MariaGtidEventFlag

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

// http://dev.mysql.com/doc/internals/en/incident-event.html
const (
	INCIDENT_NONE uint16 = iota
	INCIDENT_LOST_EVENTS
)

type BinlogEventFlag uint16

// http://dev.mysql.com/doc/internals/en/binlog-event-flag.html
const (
	// 	gets unset in the FORMAT_DESCRIPTION_EVENT when the file gets closed to detect broken binlogs
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

var StopEvent = eventTypePack(STOP_EVENT)
var HeartbeatEvent = eventTypePack(HEARTBEAT_EVENT)

type eventTypePack EventType

func (eventTypePack) Read(c proto.Reader) {
}
func (p eventTypePack) Write(c proto.Writer) {
}
func (p eventTypePack) Type() EventType {
	return EventType(p)
}

type IntVarType uint8

// http://dev.mysql.com/doc/internals/en/intvar-event.html
const (
	INVALID_INT_EVENT IntVarType = iota
	LAST_INSERT_ID_EVENT
	INSERT_ID_EVENT
)

type IntvarEvent struct {
	VarType IntVarType
	Value   uint64
}

func (p *IntvarEvent) Read(c proto.Reader) {
	c.Get(&p.VarType, &p.Value)
}
func (p *IntvarEvent) Type() EventType {
	return INTVAR_EVENT
}

/*
4              slave_proxy_id
4              exec_time
4              skip_lines
1              table_name_len
1              schema_len
4              num_fields

1              field_term
1              enclosed_by
1              line_term
1              line_start
1              escaped_by
1              opt_flags
1              empty_flags

string.var_len [len=1//num_fields] (array of 1-byte) field_name_lengths
string.var_len [len=sum(field_name_lengths) + num_fields] (array of nul-terminated strings) field_names
string.var_len [len=table_name_len + 1] (nul-terminated string) table_name
string.var_len [len=schema_len + 1] (nul-terminated string) schema_name
string.NUL     file_name
*/

type LoadEvent struct {
	SlaveProxyId  uint32
	ExecutionTime uint32
	SkipLines     uint32
}

func (p *LoadEvent) Read(c proto.Reader) {
	// TODO
}
func (p *LoadEvent) Type() EventType {
	return LOAD_EVENT
}

/*
4              slave_proxy_id
4              exec_time
4              skip_lines
1              table_name_len
1              schema_len
4              num_fields

1              field_term_len
string.var_len field_term
1              enclosed_by_len
string.var_len enclosed_by
1              line_term_len
string.var_len line_term
1              line_start_len
string.var_len line_start
1              escaped_by_len
string.var_len escaped_by
1              opt_flags

string.var_len [len=1//num_fields] (array of 1-byte) field_name_lengths
string.var_len [len=sum(field_name_lengths) + num_fields] (array of nul-terminated strings) field_names
string.var_len [len=table_name_len] (nul-terminated string) table_name
string.var_len [len=schema_len] (nul-terminated string) schema_name
string.EOF     file_name
*/
type NewLoadEvent struct {
	SlaveProxyId  uint32
	ExecutionTime uint32
	SkipLines     uint32
}

func (p *NewLoadEvent) Read(c proto.Reader) {
	// TODO
}
func (p *NewLoadEvent) Type() EventType {
	return NEW_LOAD_EVENT
}

type CreateFileEvent struct {
	FileId    uint32
	BlockData []byte
}

func (p *CreateFileEvent) Read(c proto.Reader) {
	c.Get(&p.FileId, &p.BlockData, proto.StrEof)
}
func (p *CreateFileEvent) Type() EventType {
	return CREATE_FILE_EVENT
}

type AppendBlockEvent CreateFileEvent

func (p *AppendBlockEvent) Type() EventType {
	return APPEND_BLOCK_EVENT
}

type BeginLoadQueryEvent CreateFileEvent

func (p *BeginLoadQueryEvent) Type() EventType {
	return BEGIN_LOAD_QUERY_EVENT
}

type DeleteFileEvent struct {
	FileId uint32
}

func (p *DeleteFileEvent) Read(c proto.Reader) {
	c.Get(&p.FileId)
}
func (p *DeleteFileEvent) Type() EventType {
	return DELETE_FILE_EVENT
}

type ExecLoadEvent struct {
	FieldId uint32
}

func (p *ExecLoadEvent) Read(c proto.Reader) {
	c.Get(&p.FieldId)
}
func (p *ExecLoadEvent) Type() EventType {
	return EXEC_LOAD_EVENT
}

type ExecuteLoadQueryEvent struct {
	SlaveProxyId     uint32
	ExecutionTime    uint32
	ErrorCode        uint16
	Status           []byte
	Schema           string
	ExecuteLoadQuery string
}

func (p *ExecuteLoadQueryEvent) Read(c proto.Reader) {
	var m uint8
	var n uint16
	c.Get(
		&p.SlaveProxyId,
		&p.ExecutionTime,
		&m,
		&p.ErrorCode,
		&n,
		&p.Status, proto.StrVar, &n,
		&p.Schema, proto.StrVar, &m,
		&p.ExecuteLoadQuery)
}
func (p *ExecuteLoadQueryEvent) Type() EventType {
	return EXECUTE_LOAD_QUERY_EVENT
}

type RandEvent struct {
	Seed1 uint64
	Seed2 uint64
}

func (p *RandEvent) Read(c proto.Reader) {
	c.Get(&p.Seed1, &p.Seed2)
}
func (p *RandEvent) Type() EventType {
	return RAND_EVENT
}

type XIDEvent struct {
	XID uint64
}

func (p *XIDEvent) Read(c proto.Reader) {
	c.Get(&p.XID)
}
func (p *XIDEvent) Type() EventType {
	return XID_EVENT
}

type UseVarEvent struct {
	Name    string
	IsNull  uint8
	VarType uint8
	Charset uint32
	Value   string
	Flags   uint8
}

func (p *UseVarEvent) Read(c proto.Reader) {
	var n uint32
	c.Get(&n, &p.Name, proto.StrVar, &n, &p.IsNull)
	if p.IsNull == 0 {
		c.Get(&p.VarType, &p.Charset, &n, &p.Value, proto.StrVar, &n)
		if c.More() {
			c.Get(&p.Flags)
		}
	}
}
func (p *UseVarEvent) Type() EventType {
	return USER_VAR_EVENT
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
	/*

	   long n = 6 + ((e.getFlags2() & MariaGtidEventData.FL_GROUP_COMMIT_ID) > 0 ? 2 : 0);
	   long skip = is.skip(n);
	   assert n == skip;
	   return e;
	*/
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
