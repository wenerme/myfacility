package proto

type CommandType uint8
const (
// internal server command
// <p/>
// Payload
// 1              [00] COM_SLEEP
// Returns
// ERR_Packet
	COM_SLEEP CommandType = iota
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

type ComPack struct {
	Type CommandType
	data []byte
	buf  *Buffer
}