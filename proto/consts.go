package proto

//go:vet
//go:generate stringer -output=strings.go -type=Capability,Status,Command,SessionState,ProtoType,ColumnType

// Packet less than 50 will not compress
const MIN_COMPRESS_LENGTH = 50

type (
	// Packet command type
	Command uint8
	// http://dev.mysql.com/doc/internals/en/packet-OK_Packet.html#cs-sect-packet-ok-sessioninfo
	SessionState uint16

	// http://dev.mysql.com/doc/internals/en/capability-flags.html
	Capability   uint32
	CharacterSet uint16
	// http://dev.mysql.com/doc/internals/en/status-flags.html#packet-Protocol::StatusFlags
	Status uint16
)

func (this Status) Has(c Status) bool {
	return this&c != 0
}

func (this Status) Remove(c Status) Status {
	return this & ^c
}

func (this Status) Add(c Status) Status {
	return this | c
}

func (this Capability) Has(c Capability) bool {
	return this&c != 0
}

func (this Capability) Remove(c Capability) Capability {
	return this & ^c
}

func (this Capability) Add(c Capability) Capability {
	return this | c
}

const (
	// one or more system variables changed. See also: session_track_system_variables
	SESSION_TRACK_SYSTEM_VARIABLES SessionState = iota
	// schema changed. See also: session_track_schema
	SESSION_TRACK_SCHEMA
	// "track state change" changed. See also: session_track_state_change
	SESSION_TRACK_STATE_CHANGE
)

const (
	// Is raised when a multi-statement transaction
	// has been started, either explicitly, by means
	// of BEGIN or COMMIT AND CHAIN, or
	// implicitly, by the first transactional
	// statement, when autocommit=off.
	SERVER_STATUS_IN_TRANS Status = 1
	// Server in auto_commit mode
	SERVER_STATUS_AUTOCOMMIT Status = 2
	// Multi query - next query exists
	SERVER_MORE_RESULTS_EXISTS      Status = 8
	SERVER_QUERY_NO_GOOD_INDEX_USED Status = 16
	SERVER_QUERY_NO_INDEX_USED      Status = 32
	//
	// The server was able to fulfill the clients request and opened a
	// read-only non-scrollable cursor for a query. This flag comes
	// in reply to COM_STMT_EXECUTE and COM_STMT_FETCH commands.
	SERVER_STATUS_CURSOR_EXISTS Status = 64
	//
	// This flag is sent when a read-only cursor is exhausted, in reply to
	// COM_STMT_FETCH command.
	SERVER_STATUS_LAST_ROW_SENT Status = 128
	//  A database was dropped
	SERVER_STATUS_DB_DROPPED           Status = 256
	SERVER_STATUS_NO_BACKSLASH_ESCAPES Status = 512
	//
	// Sent to the client if after a prepared statement reprepare
	// we discovered that the new statement returns a different
	// number of result set columns.
	SERVER_STATUS_METADATA_CHANGED Status = 1024
	SERVER_QUERY_WAS_SLOW          Status = 2048

	//
	// To mark ResultSet containing output parameter values.
	SERVER_PS_OUT_PARAMS Status = 4096

	//
	// Set at the same time as SERVER_STATUS_IN_TRANS if the started
	// multi-statement transaction is a read-only transaction. Cleared
	// when the transaction commits or aborts. Since this flag is sent
	// to clients in OK and EOF packets, the flag indicates the
	// transaction status at the end of command execution.
	SERVER_STATUS_IN_TRANS_READONLY Status = 8192

	// connection state information has changed
	SERVER_SESSION_STATE_CHANGED Status = 0x4000
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

const (
	// internal server command
	// <p/>
	// Payload
	// 1              [00] COM_SLEEP
	// Returns
	// ERR_Packet
	COM_SLEEP Command = iota
	// Tells the server that the client wants to close the connection
	// <p/>
	// response: either a connection close or a OK_Packet
	// <p/>
	// Payload
	// 1              [01] COM_QUIT
	// Fields
	// command (1) -- 0x01 COM_QUIT
	// <p/>
	// Example
	// 01 00 00 00 01
	COM_QUIT
	// change the default schema of the connection
	// <p/>
	// Returns
	// OK_Packet or ERR_Packet
	// <p/>
	// Payload
	// 1              [02] COM_INIT_DB
	// string[EOF]    schema name
	// Fields
	// command (1) -- 0x02 COM_INIT_DB
	// <p/>
	// schema_name (string.EOF) -- name of the schema to change to
	// <p/>
	// Example
	// 05 00 00 00 02 74 65 73    74                         .....test
	COM_INIT_DB
	// A COM_QUERY is used to send the server a text-based query that is executed immediately.
	// <p/>
	// The server replies to a COM_QUERY packet with a COM_QUERY Response.
	// <p/>
	// The length of the query-string is a taken from the packet length - 1.
	// <p/>
	// Payload
	// 1              [03] COM_QUERY
	// string[EOF]    the query the server shall execute
	// Fields
	// command_id (1) -- 0x03 COM_QUERY
	// <p/>
	// query (string.EOF) -- query_text
	// <p/>
	// Implemented By
	// mysql_query()
	// <p/>
	// Returns
	// COM_QUERY_Response
	// <p/>
	// Example
	// 21 00 00 00 03 73 65 6c    65 63 74 20 40 40 76 65    !....select @@ve
	// 72 73 69 6f 6e 5f 63 6f    6d 6d 65 6e 74 20 6c 69    rsion_comment li
	// 6d 69 74 20 31                                        mit 1
	COM_QUERY
	// get the column definitions of a table
	// <p/>
	// Payload
	// 1              [04] COM_FIELD_LIST
	// string[NUL]    table
	// string[EOF]    field wildcard
	// Returns
	// COM_FIELD_LIST response
	// <p/>
	// Implemented By
	// mysql_list_fields()
	// <p/>
	// Response
	// The response to a COM_FIELD_LIST can either be a
	// <p/>
	// a ERR_Packet or
	// <p/>
	// one or more Column Definition packets and a closing EOF_Packet
	// <p/>
	// 31 00 00 01 03 64 65 66    04 74 65 73 74 09 66 69    1....def.test.fi
	// 65 6c 64 6c 69 73 74 09    66 69 65 6c 64 6c 69 73    eldlist.fieldlis
	// 74 02 69 64 02 69 64 0c    3f 00 0b 00 00 00 03 00    t.id.id.?.......
	// 00 00 00 00 fb 05 00 00    02 fe 00 00 02 00          ..............
	COM_FIELD_LIST
	// create a schema
	// <p/>
	// Payload
	// 1              [05] COM_CREATE_DB
	// string[EOF]    schema name
	// Returns
	// OK_Packet or ERR_Packet
	// <p/>
	// Example
	// 05 00 00 00 05 74 65 73    74                         .....test
	COM_CREATE_DB
	// drop a schema
	// <p/>
	// Payload
	// 1              [06] COM_DROP_DB
	// string[EOF]    schema name
	// Returns
	// OK_Packet or ERR_Packet
	// <p/>
	// Example
	// 05 00 00 00 06 74 65 73    74                         .....test
	COM_DROP_DB
	// A low-level version of several FLUSH ... and RESET ... statements.
	// <p/>
	// COM_REFRESH:
	// Call REFRESH or FLUSH statements
	// <p/>
	// Payload
	// 1              [07] COM_REFRESH
	// 1              sub_command
	// Fields
	// command (1) -- 0x07 COM_REFRESH
	// <p/>
	// sub_command (1) -- a bitmask of sub-systems to refresh
	// <p/>
	// Returns
	// OK_Packet or ERR_Packet
	COM_REFRESH
	// COM_SHUTDOWN is used to shut down the MySQL server.
	// <p/>
	// The SHUTDOWN privilege is required for this operation.
	// <p/>
	// COM_SHUTDOWN:
	// shut down the server
	// <p/>
	// Payload
	// 1              [08] COM_SHUTDOWN
	// if more data {
	// 1              shutdown type
	// }
	// Fields
	// command (1) -- 0x08 COM_SHUTDOWN
	// <p/>
	// sub_command (1) -- optional if sub_command is 0x00
	// <p/>
	// Returns
	// EOF_Packet or ERR_Packet
	// <p/>
	// Note
	// Even if several shutdown types are defined, right now only one is in use: SHUTDOWN_WAIT_ALL_BUFFERS
	COM_SHUTDOWN
	// Get a human readable string of internal statistics.
	// <p/>
	// Returns
	// string.EOF
	// <p/>
	// Payload
	// 1              [09] COM_STATISTICS
	COM_STATISTICS
	// get a list of active threads
	// <p/>
	// Returns
	// a ProtocolText::Resultset or ERR_Packet
	// <p/>
	// Payload
	// 1              [0a] COM_PROCCESS_INFO
	COM_PROCESS_INFO
	// an internal command in the server
	// <p/>
	// Payload
	// 1              [0b] COM_CONNECT
	// Returns
	// ERR_Packet
	COM_CONNECT
	// Same as KILL &lt;id>.
	// <p/>
	// ask the server to terminate a connection
	// <p/>
	// Returns
	// OK_Packet or ERR_Packet
	// <p/>
	// Payload
	// 1              [0c] COM_PROCCESS_KILL
	// 4              connection id
	COM_PROCESS_KILL
	// COM_DEBUG triggers a dump on internal debug info to stdout of the mysql-server.
	// <p/>
	// The SUPER privilege is required for this operation.
	// <p/>
	// dump debug info to stdout
	// <p/>
	// Returns
	// EOF_Packet or ERR_Packet on error
	// <p/>
	// Payload
	// 1              [0d] COM_DEBUG
	COM_DEBUG
	// check if the server is alive
	// <p/>
	// Returns
	// OK_Packet
	// <p/>
	// Payload
	// 1              [0e] COM_PING
	COM_PING
	// an internal command in the server
	// <p/>
	// Payload
	// 1              [0f] COM_TIME
	// Returns
	// ERR_Packet
	COM_TIME
	// an internal command in the server
	// <p/>
	// Payload
	// 1              [10] COM_DELAYED_INSERT
	// Returns
	// ERR_Packet
	COM_DELAYED_INSERT
	// COM_CHANGE_USER changes the user of the current connection and reset the connection state.
	// <p/>
	// user variables
	// <p/>
	// temp tables
	// <p/>
	// prepared statemants
	// <p/>
	// ... and others
	// <p/>
	// It is followed by the same states as the initial handshake.
	// <p/>
	// COM_CHANGE_USER:
	// change the user of the current connection
	// <p/>
	// Returns
	// Authentication Method Switch Request Packet or ERR_Packet
	// <p/>
	// Payload
	// 1              [11] COM_CHANGE_USER
	// string[NUL]    user
	// if capabilities & SECURE_CONNECTION {
	// 1              auth-response-len
	// string[$len]   auth-response
	// } else {
	// string[NUL]    auth-response
	// }
	// string[NUL]    schema-name
	// if more data {
	// 2              character-set
	// if capabilities & CLIENT_PLUGIN_AUTH {
	// string[NUL]    auth plugin name
	// }
	// if capabilities & CLIENT_CONNECT_ATTRS) {
	// lenenc-int     length of all key-values
	// lenenc-str     key
	// lenenc-str     value
	// if-more data in 'length of all key-values', more keys and value pairs
	// }
	// }
	// Fields
	// command (1) -- command byte
	// <p/>
	// username (string.NUL) -- user name
	// <p/>
	// auth_plugin_data_len (1) -- length of the auth_plugin_data filed
	// <p/>
	// auth_plugin_data (string.var_len) -- auth data
	// <p/>
	// schema (string.NUL) -- default schema
	// <p/>
	// character_set (2) -- new connection character set (see Protocol::CharacterSet)
	// <p/>
	// auth_plugin_name (string.NUL) -- name of the auth plugin that auth_plugin_data corresponds to
	// <p/>
	// connect_attrs_len (lenenc_int) -- length in bytes of the following block of key-value pairs
	// <p/>
	// Implemented By
	// parse_com_change_user_packet()
	// <p/>
	// character set is the connection character set as defined in Protocol::CharacterSet and is also the encoding of user and schema-name.
	COM_CHANGE_USER
	COM_BINLOG_DUMP
	COM_TABLE_DUMP
	COM_CONNECT_OUT
	COM_REGISTER_SLAVE
	// create a prepared statement
	// <p/>
	// Fields
	// command (1) -- [16] the COM_STMT_PREPARE command
	// <p/>
	// query (string.EOF) -- the query to prepare
	// <p/>
	// Example
	// 1c 00 00 00 16 53 45 4c    45 43 54 20 43 4f 4e 43    .....SELECT CONC
	// 41 54 28 3f 2c 20 3f 29    20 41 53 20 63 6f 6c 31    AT(?, ?) AS col1
	// Implemented By
	// mysqld_stmt_prepare()
	// <p/>
	// Return
	// COM_STMT_PREPARE_OK on success, ERR_Packet otherwise
	// <p/>
	// Note
	// As LOAD DATA isn't supported by COM_STMT_PREPARE yet, no Protocol::LOCAL_INFILE_Request is expected here. Compare this to COM_QUERY_Response.
	COM_STMT_PREPARE
	// COM_STMT_EXECUTE asks the server to execute a prepared statement as identified by stmt-id.
	// <pre>
	//  * It sends the values for the placeholders of the prepared statement
	// (if it contained any) in Binary Protocol Value form. The type of each
	// parameter is made up of two bytes:
	//      the type as in Protocol::ColumnType
	//      a flag byte which has the highest bit set if the type is unsigned [80]
	//  * The num-params used for this packet has to match the num_params of the COM_STMT_PREPARE_OK of the corresponding prepared statement.
	//  * The server returns a COM_STMT_EXECUTE Response.
	//  * COM_STMT_EXECUTE:
	// COM_STMT_EXECUTE
	// execute a prepared statement
	//  * direction: client -> server
	// response: COM_STMT_EXECUTE Response
	//  * payload:
	// 1              [17] COM_STMT_EXECUTE
	// 4              stmt-id
	// 1              flags
	// 4              iteration-count
	// if num-params > 0:
	// n              NULL-bitmap, length: (num-params+7)/8
	// 1              new-params-bound-flag
	// if new-params-bound-flag == 1:
	// n              type of each parameter, length: num-params * 2
	// n              value of each parameter
	//  * example:
	// 12 00 00 00 17 01 00 00    00 00 01 00 00 00 00 01    ................
	// 0f 00 03 66 6f 6f                                     ...foo
	// The iteration-count is always 1.
	//  * The flags are:
	// Flags Constant Name
	// 0x00 CURSOR_TYPE_NO_CURSOR
	// 0x01 CURSOR_TYPE_READ_ONLY
	// 0x02 CURSOR_TYPE_FOR_UPDATE
	// 0x04 CURSOR_TYPE_SCROLLABLE
	// </pre>
	// NULL-bitmap is like NULL-bitmap for the Binary Protocol Resultset Row just that it has a bit-offset of 0.
	//  * @see <a href=http://dev.mysql.com/doc/internals/en/com-stmt-execute.html>com-stmt-execute</a>
	COM_STMT_EXECUTE
	// COM_STMT_SEND_LONG_DATA sends the data for a column. Repeating to send it, appends the data to the parameter.
	// <p/>
	// No response is sent back to the client.
	// <p/>
	// COM_STMT_SEND_LONG_DATA:
	// COM_STMT_SEND_LONG_DATA
	// direction: client -> server
	// response: none
	// <p/>
	// payload:
	// 1              [18] COM_STMT_SEND_LONG_DATA
	// 4              statement-id
	// 2              param-id
	// n              data
	// COM_STMT_SEND_LONG_DATA has to be sent before COM_STMT_EXECUTE.
	COM_STMT_SEND_LONG_DATA
	// COM_STMT_CLOSE deallocates a prepared statement
	// <pre>
	// No response is sent back to the client.
	//  * COM_STMT_CLOSE:
	// COM_STMT_CLOSE
	// direction: client -> server
	// response: none
	//  * payload:
	// 1              [19] COM_STMT_CLOSE
	// 4              statement-id
	//  * example:
	// 05 00 00 00 19 01 00 00    00                         .........
	// </pre>
	//  * @see <a href=http://dev.mysql.com/doc/internals/en/com-stmt-close.html>com-stmt-close</a>
	COM_STMT_CLOSE
	// COM_STMT_RESET resets the data of a prepared statement which was accumulated with COM_STMT_SEND_LONG_DATA commands and closes the cursor if it was opened with COM_STMT_EXECUTE
	// <pre>
	// The server will send a OK_Packet if the statement could be reset, a ERR_Packet if not.
	//  * COM_STMT_RESET:
	// COM_STMT_RESET
	// direction: client -> server
	// response: OK or ERR
	//  * payload:
	// 1              [1a] COM_STMT_RESET
	// 4              statement-id
	//  * example:
	// 05 00 00 00 1a 01 00 00    00                         .........
	// </pre>
	//  * @see <a href=http://dev.mysql.com/doc/internals/en/com-stmt-reset.html>com-stmt-reset</a>
	COM_STMT_RESET
	// Allows to enable and disable: CLIENT_MULTI_STATEMENTS
	// <pre>
	// for the current connection. The option operation is one of:
	//  * Operation Constant Name
	// 0 MYSQL_OPTION_MULTI_STATEMENTS_ON
	// 1 MYSQL_OPTION_MULTI_STATEMENTS_OFF
	//  * On success it returns a EOF_Packet otherwise a ERR_Packet.
	//  * COM_SET_OPTION
	// set options for the current connection
	//  * response: EOF or ERR
	//  * payload:
	// 1              [1b] COM_SET_OPTION
	// 2              option operation
	// </pre>
	//  * @see <a href=http://dev.mysql.com/doc/internals/en/com-set-option.html>com-set-option</a>
	COM_SET_OPTION
	// Fetch rows from a existing resultset after a COM_STMT_EXECUTE.
	// <pre>
	// Payload
	// 1              [1c] COM_STMT_FETCH
	// 4              stmt-id
	// 4              num rows
	// Returns
	// a COM_STMT_FETCH response( a multi-resultset or a ERR_Packet )
	// </pre>
	COM_STMT_FETCH
	// an internal command in the server
	// <p/>
	// Payload
	// 1              [1d] COM_DAEMON
	// Returns
	// ERR_Packet
	COM_DAEMON
	COM_BINLOG_DUMP_GTID
	// Resets the session state; more lightweight than COM_CHANGE_USER because it does not close and reopen the connection, and does not re-authenticate
	// <p/>
	// Payload
	// 1              [1f] COM_RESET_CONNECTION
	// Returns
	// a ERR_Packet
	// <p/>
	// a OK_Packet
	COM_RESET_CONNECTION
)

const (
	//Use the improved version of Old Password Authentication
	//<h2>Note</h2>
	//assumed to be set since 4.1.1
	CLIENT_LONG_PASSWORD Capability = 0x00000001

	//Send found rows instead of affected rows in EOF_Packet
	CLIENT_FOUND_ROWS Capability = 0x00000002
	//
	//Longer flags in Protocol::ColumnDefinition320
	//<h2>Server</h2>
	//supports longer flags
	//<h2>Client</h2>
	//expects longer flags
	//
	CLIENT_LONG_FLAG Capability = 0x00000004
	//
	//One can specify db on connect in Handshake Response Packet
	//<h2>Server</h2>
	//supports schema-name in Handshake Response Packet
	//<h2>Client</h2>
	//Handshake Response Packet contains a schema-name
	//
	CLIENT_CONNECT_WITH_DB Capability = 0x00000008
	//
	//<h2>Server</h2>
	//Don't allow database.table.column
	//
	CLIENT_NO_SCHEMA Capability = 0x00000010
	//
	//Compression protocol supported
	//<h2>Server</h2>
	//supports compression
	//<h2>Client</h2>
	//switches to Compression compressed protocol after successful authentication
	//
	CLIENT_COMPRESS Capability = 0x00000020
	//
	//Special handling of ODBC behaviour
	//<h2>Note</h2>
	//no special behaviour since 3.22
	//
	CLIENT_ODBC Capability = 0x00000040
	//
	//Can use LOAD DATA LOCAL
	//<h2>Server</h2>
	//allows the LOCAL INFILE request of LOAD DATA|XML
	//<h2>Client</h2>
	//will handle LOCAL INFILE request
	//
	CLIENT_LOCAL_FILES Capability = 0x00000080
	//
	//<h2>Server</h2>
	//parser can ignore spaces before '('
	//<h2>Client</h2>
	//let the parser ignore spaces before '('
	//
	CLIENT_IGNORE_SPACE Capability = 0x00000100
	//
	//<h2>Server</h2>
	//supports the 4.1 protocol
	//<h2>Client</h2>
	//uses the 4.1 protocol
	//<h2>Note</h2>
	//this value was CLIENT_CHANGE_USER in 3.22, unused in 4.0
	//
	CLIENT_PROTOCOL_41 Capability = 0x00000200
	//
	//wait_timeout vs. wait_interactive_timeout
	//<h2>Server</h2>
	//supports interactive and non-interactive clients
	//<h2>Client</h2>
	//client is interactive
	//<p/>
	//See
	//mysql_real_connect()
	//
	CLIENT_INTERACTIVE Capability = 0x00000400
	//
	//<h2>Server</h2>
	//supports SSL
	//<h2>Client</h2>
	//switch to SSL after sending the capability-flags
	//
	CLIENT_SSL Capability = 0x00000800
	//
	//<h2>Client</h2>
	//Don't issue SIGPIPE if network failures (libmysqlclient only)
	//<p/>
	//See
	//mysql_real_connect()
	//
	CLIENT_IGNORE_SIGPIPE Capability = 0x00001000
	//
	//<h2>Server</h2>
	//can send status flags in EOF_Packet
	//<h2>Client</h2>
	//expects status flags in EOF_Packet
	//<h2>Note</h2>
	//this flag is optional in 3.23, but set all the time by the server since 4.0
	//
	CLIENT_TRANSACTIONS Capability = 0x00002000

	//unused
	//<h2>Note</h2>
	//Was named CLIENT_PROTOCOL_41 in 4.1.0
	//
	CLIENT_RESERVED Capability = 0x00004000
	//
	//<h2>Server</h2>
	//supports Authentication::Native41
	//<h2>Client</h2>
	//supports Authentication::Native41
	//
	CLIENT_SECURE_CONNECTION Capability = 0x00008000
	//
	//<h2>Server</h2>
	//can handle multiple statements per COM_QUERY and COM_STMT_PREPARE
	//<h2>Client</h2>
	//may send multiple statements per COM_QUERY and COM_STMT_PREPARE
	//<h2>Note</h2>
	//was named CLIENT_MULTI_QUERIES in 4.1.0, renamed later
	//<p/>
	//Requires
	//CLIENT_PROTOCOL_41
	//
	CLIENT_MULTI_STATEMENTS Capability = 0x00010000

	//<h2>Server</h2>
	//can send multiple resultsets for COM_QUERY
	//<h2>Client</h2>
	//can handle multiple resultsets for COM_QUERY
	//<p/>
	//Requires
	//CLIENT_PROTOCOL_41
	CLIENT_MULTI_RESULTS Capability = 0x00020000
	//<h2>Server</h2>
	//can send multiple resultsets for COM_STMT_EXECUTE
	//<h2>Client</h2>
	//can handle multiple resultsets for COM_STMT_EXECUTE
	//<p/>
	//Requires
	//CLIENT_PROTOCOL_41
	//
	CLIENT_PS_MULTI_RESULTS Capability = 0x00040000
	//<h2>Server</h2>
	//sends extra data in Initial Handshake Packet and supports the pluggable authentication protocol.
	//<h2>Client</h2>
	//supports auth plugins
	//<p/>
	//Requires
	//CLIENT_PROTOCOL_41
	CLIENT_PLUGIN_AUTH Capability = 0x00080000
	//
	//<h2>Server</h2>
	//allows connection attributes in Protocol::HandshakeResponse41
	//<h2>Client</h2>
	//sends connection attributes in Protocol::HandshakeResponse41
	//
	CLIENT_CONNECT_ATTRS Capability = 0x00100000
	//
	//<h2>Server</h2>
	//understands length encoded integer for auth response data in Protocol::HandshakeResponse41
	//<h2>Client</h2>
	//length of auth response data in Protocol::HandshakeResponse41 is a length encoded integer
	//<h2>Note</h2>
	//the flag was introduce in 5.6.6, but had the wrong value.
	//
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA Capability = 0x00200000
	//
	//<h2>Server</h2>
	//announces support for expired password extension
	//<h2>Client</h2>
	//can handle expired passwords
	//
	CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS Capability = 0x00400000

	//<h2>Server</h2>
	//can set SERVER_SESSION_STATE_CHANGED in the Status Flags and send session-state change data after a OK packet
	//<h2>Client</h2>
	//expects the server to send sesson-state changes after a OK packet
	//<h2>Background</h2>
	//To support CLIENT_SESSION_TRACK additional information has to be sent after all succesful commands. While the OK packet is extensible, the EOF packet is not due to the overlap of its bytes with the content of the Text Resultset Row.
	//Therefore, the EOF packet in the Text Resultset is replaced with an OK packet.
	CLIENT_SSL_VERIFY_SERVER_CERT Capability = 0x00800000

	// <h2>Server</h2>
	// can send OK after a Text Resultset
	// <h2>Client</h2>
	// expects a OK (instead of EOF) after the resultset rows of a Text Resultset.
	// <h2>Background</h2>
	// To support CLIENT_SESSION_TRACK additional information has to be sent after all succesful commands. While the OK packet is extensible,
	// the EOF packet is not due to the overlap of its bytes with the content of the Text Resultset Row.
	// Therefore, the EOF packet in the Text Resultset is replaced with an OK packet.
	CLIENT_DEPRECATE_EOF Capability = 0x01000000

	CLIENT_REMEMBER_OPTIONS Capability = 1 << 31

	CLIENT_ALL_FLAGS Capability = CLIENT_LONG_PASSWORD | CLIENT_FOUND_ROWS | CLIENT_LONG_FLAG |
		CLIENT_CONNECT_WITH_DB | CLIENT_NO_SCHEMA | CLIENT_COMPRESS | CLIENT_ODBC |
		CLIENT_LOCAL_FILES | CLIENT_IGNORE_SPACE | CLIENT_PROTOCOL_41 |
		CLIENT_INTERACTIVE | CLIENT_SSL | CLIENT_IGNORE_SIGPIPE | CLIENT_TRANSACTIONS |
		CLIENT_RESERVED | CLIENT_SECURE_CONNECTION | CLIENT_MULTI_STATEMENTS |
		CLIENT_MULTI_RESULTS | CLIENT_PS_MULTI_RESULTS | CLIENT_SSL_VERIFY_SERVER_CERT |
		CLIENT_REMEMBER_OPTIONS | CLIENT_PLUGIN_AUTH | CLIENT_CONNECT_ATTRS |
		CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA | CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS

	CLIENT_BASIC_FLAGS Capability = CLIENT_ALL_FLAGS & ^CLIENT_SSL & ^CLIENT_COMPRESS & ^CLIENT_SSL_VERIFY_SERVER_CERT

	// <h2>Server</h2>
	// can set SERVER_SESSION_STATE_CHANGED in the Status Flags and send session-state change data after a OK packet
	// <h2>Client</h2>
	// expects the server to send sesson-state changes after a OK packet
	CLIENT_SESSION_TRACK Capability = 0x00800000
)
