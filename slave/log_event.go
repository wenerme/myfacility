package slave
import (
	. "../proto"
)
type LogEventType Int1
type LogEventFlag Int2

const (

	EVENT_INVALID_LOGGING LogEventType = iota
/*
  The event must be written to a cache and upon commit or rollback
  written to the binary log.
*/
	EVENT_NORMAL_LOGGING
/*
  The event must be written to an empty cache and immediatly written
  to the binary log without waiting for any other event.
*/
	EVENT_IMMEDIATE_LOGGING
/*
   If there is a need for different types, introduce them before this.
*/
	EVENT_CACHE_LOGGING_COUNT
)

const (

/*
   This flag only makes sense for Format_description_log_event. It is set
   when the event is written, and *reset* when a binlog file is
   closed (yes, it's the only case when MySQL modifies already written
   part of binlog).  Thus it is a reliable indicator that binlog was
   closed correctly.  (Stop_log_event is not enough, there's always a
   small chance that mysqld crashes in the middle of insert and end of
   the binlog would look like a Stop_log_event).

   This flag is used to detect a restart after a crash, and to provide
   "unbreakable" binlog. The problem is that on a crash storage engines
   rollback automatically, while binlog does not.  To solve this we use this
   flag and automatically append ROLLBACK to every non-closed binlog (append
   virtually, on reading, file itself is not changed). If this flag is found,
   mysqlbinlog simply prints "ROLLBACK" Replication master does not abort on
   binlog corruption, but takes it as EOF, and replication slave forces a
   rollback in this case.

   Note, that old binlogs does not have this flag set, so we get a
   a backward-compatible behaviour.
*/
	LOG_EVENT_BINLOG_IN_USE_F LogEventFlag = 0x1

/**
  @def LOG_EVENT_THREAD_SPECIFIC_F

  If the query depends on the thread (for example: TEMPORARY TABLE).
  Currently this is used by mysqlbinlog to know it must print
  SET @@PSEUDO_THREAD_ID=xx; before the query (it would not hurt to print it
  for every query but this would be slow).
*/
	LOG_EVENT_THREAD_SPECIFIC_F LogEventFlag = 0x4

/**
  @def LOG_EVENT_SUPPRESS_USE_F

  Suppress the generation of 'USE' statements before the actual
  statement. This flag should be set for any events that does not need
  the current database set to function correctly. Most notable cases
  are 'CREATE DATABASE' and 'DROP DATABASE'.

  This flags should only be used in exceptional circumstances, since
  it introduce a significant change in behaviour regarding the
  replication logic together with the flags --binlog-do-db and
  --replicated-do-db.
 */
	LOG_EVENT_SUPPRESS_USE_F LogEventFlag = 0x8

/*
  Note: this is a place holder for the flag
  LOG_EVENT_UPDATE_TABLE_MAP_VERSION_F (0x10), which is not used any
  more, please do not reused this value for other flags.
 */

/**
   @def LOG_EVENT_ARTIFICIAL_F
   
   Artificial events are created arbitarily and not written to binary
   log

   These events should not update the master log position when slave
   SQL thread executes them.
*/
	LOG_EVENT_ARTIFICIAL_F LogEventFlag = 0x20

/**
   @def LOG_EVENT_RELAY_LOG_F
   
   Events with this flag set are created by slave IO thread and written
   to relay log
*/
	LOG_EVENT_RELAY_LOG_F LogEventFlag = 0x40

/**
   @def LOG_EVENT_IGNORABLE_F

   For an event, 'e', carrying a type code, that a slave,
   's', does not recognize, 's' will check 'e' for
   LOG_EVENT_IGNORABLE_F, and if the flag is set, then 'e'
   is ignored. Otherwise, 's' acknowledges that it has
   found an unknown event in the relay log.
*/
	LOG_EVENT_IGNORABLE_F LogEventFlag = 0x80

/**
   @def LOG_EVENT_NO_FILTER_F

   Events with this flag are not filtered (e.g. on the current
   database) and are always written to the binary log regardless of
   filters.
*/
	LOG_EVENT_NO_FILTER_F LogEventFlag = 0x100

/**
   MTS: group of events can be marked to force its execution
   in isolation from any other Workers.
   So it's a marker for Coordinator to memorize and perform necessary
   operations in order to guarantee no interference from other Workers.
   The flag can be set ON only for an event that terminates its group.
   Typically that is done for a transaction that contains 
   a query accessing more than OVER_MAX_DBS_IN_EVENT_MTS databases.
*/
	LOG_EVENT_MTS_ISOLATE_F LogEventFlag = 0x200
)


type EventHeader struct {
	// The time when the query started, in seconds since 1970.
	Timestamp      Int4
	// LogEventType
	EventType      Int1
	// Server ID of the server that created the event.
	ServerId       Int4
	// The total size of this event, in bytes.  In other words, this
	// is the sum of the sizes of Common-Header, Post-Header, and Body.
	TotalSize      Int4
	// The position of the next event in the master binary log, in
	// bytes from the beginning of the file.  In a binlog that is not a
	// relay log, this is just the position of the next event, in bytes
	// from the beginning of the file.  In a relay log, this is
	// the position of the next event in the master's binlog.
	MasterPosition Int4
	// LogEventFlag
	Flags          Int2
}

type QueryLogEvent struct {
	/* Post Header */
	// An integer identifying the client thread that issued the
	// query.  The id is unique per server.  (Note, however, that two
	// threads on different servers may have the same slave_proxy_id.)
	// This is used when a client thread creates a temporary table local
	// to the client.  The slave_proxy_id is used to distinguish
	// temporary tables that belong to different clients.
	SlaveProxyId Int4
	// The time from when the query started to when it was logged in
	// the binlog, in seconds.
	ExecTime     Int4
	// The length of the name of the currently selected database.
	DbLen        Int1
	//Error code generated by the master.  If the master fails, the
	//slave will fail with the same error code, except for the error
	//codes ER_DB_CREATE_EXISTS == 1007 and ER_DB_DROP_EXISTS == 1008.
	ErrorCode    Int2
	//The length of the status_vars block of the Body, in bytes. See
	//@ref query_log_event_status_vars "below".
	StatusLen    Int2

	/* Body */
	//Zero or more status variables.  Each status variable consists
	//of one byte identifying the variable stored, followed by the value
	//of the variable.  The possible variables are listed separately in
	//the table @ref Table_query_log_event_status_vars "below".  MySQL
	//always writes events in the order defined below; however, it is
	//capable of reading them in any order.
	Status       StrF

	//The currently selected database, as a null-terminated string.
	//(The trailing zero is redundant since the length is already known;
	//it is db_len from Post-Header.)
	Db           StrF
	// The SQL query.
	Query        StrE
}