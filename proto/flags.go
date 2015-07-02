package proto

// http://dev.mysql.com/doc/internals/en/capability-flags.html
type Capability uint32

// http://dev.mysql.com/doc/internals/en/status-flags.html#packet-Protocol::StatusFlags
type Status uint16

func (this Status) Has(c Status) bool {
	return this & c != 0
}

func (this Status) Remove(c Status) Status {
	return this & ^c
}

func (this Status) Add(c Status) Status {
	return this | c
}


func (this Capability) Has(c Capability) bool {
	return this & c != 0
}

func (this Capability) Remove(c Capability) Capability {
	return this & ^ c
}

func (this Capability) Add(c Capability) Capability {
	return this | c
}


const (
// Is raised when a multi-statement transaction
// has been started, either explicitly, by means
// of BEGIN or COMMIT AND CHAIN, or
// implicitly, by the first transactional
// statement, when autocommit=off.
	SERVER_STATUS_IN_TRANS Status = 1
// Server in auto_commit mode
	SERVER_STATUS_AUTOCOMMIT = 2
// Multi query - next query exists
	SERVER_MORE_RESULTS_EXISTS = 8
	SERVER_QUERY_NO_GOOD_INDEX_USED = 16
	SERVER_QUERY_NO_INDEX_USED = 32
//
// The server was able to fulfill the clients request and opened a
// read-only non-scrollable cursor for a query. This flag comes
// in reply to COM_STMT_EXECUTE and COM_STMT_FETCH commands.
	SERVER_STATUS_CURSOR_EXISTS = 64
//
// This flag is sent when a read-only cursor is exhausted, in reply to
// COM_STMT_FETCH command.
	SERVER_STATUS_LAST_ROW_SENT = 128
//  A database was dropped
	SERVER_STATUS_DB_DROPPED = 256
	SERVER_STATUS_NO_BACKSLASH_ESCAPES = 512
//
// Sent to the client if after a prepared statement reprepare
// we discovered that the new statement returns a different
// number of result set columns.
	SERVER_STATUS_METADATA_CHANGED = 1024
	SERVER_QUERY_WAS_SLOW = 2048

//
// To mark ResultSet containing output parameter values.
	SERVER_PS_OUT_PARAMS = 4096

//
// Set at the same time as SERVER_STATUS_IN_TRANS if the started
// multi-statement transaction is a read-only transaction. Cleared
// when the transaction commits or aborts. Since this flag is sent
// to clients in OK and EOF packets, the flag indicates the
// transaction status at the end of command execution.
	SERVER_STATUS_IN_TRANS_READONLY = 8192

// connection state information has changed
	SERVER_SESSION_STATE_CHANGED = 0x4000
//
// Server status flags that must be cleared when starting
// execution of a new SQL statement.
// Flags from this set are only added to the
// current server status by the execution engine, but
// never removed -- the execution engine expects them
// to disappear automagically by the next command.
	SERVER_STATUS_CLEAR_SET Status = SERVER_QUERY_NO_GOOD_INDEX_USED |
	SERVER_QUERY_NO_INDEX_USED |
	SERVER_MORE_RESULTS_EXISTS |
	SERVER_STATUS_METADATA_CHANGED |
	SERVER_QUERY_WAS_SLOW |
	SERVER_STATUS_DB_DROPPED |
	SERVER_STATUS_CURSOR_EXISTS |
	SERVER_STATUS_LAST_ROW_SENT
)