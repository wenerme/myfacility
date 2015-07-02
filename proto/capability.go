package proto


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

	CLIENT_ALL_FLAGS  Capability = CLIENT_LONG_PASSWORD| CLIENT_FOUND_ROWS| CLIENT_LONG_FLAG|
	CLIENT_CONNECT_WITH_DB| CLIENT_NO_SCHEMA| CLIENT_COMPRESS| CLIENT_ODBC|
	CLIENT_LOCAL_FILES| CLIENT_IGNORE_SPACE| CLIENT_PROTOCOL_41|
	CLIENT_INTERACTIVE| CLIENT_SSL| CLIENT_IGNORE_SIGPIPE| CLIENT_TRANSACTIONS|
	CLIENT_RESERVED| CLIENT_SECURE_CONNECTION| CLIENT_MULTI_STATEMENTS|
	CLIENT_MULTI_RESULTS| CLIENT_PS_MULTI_RESULTS| CLIENT_SSL_VERIFY_SERVER_CERT|
	CLIENT_REMEMBER_OPTIONS| CLIENT_PLUGIN_AUTH| CLIENT_CONNECT_ATTRS|
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA| CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS


	CLIENT_BASIC_FLAGS Capability = CLIENT_ALL_FLAGS & ^CLIENT_SSL & ^CLIENT_COMPRESS & ^CLIENT_SSL_VERIFY_SERVER_CERT

// <h2>Server</h2>
// can set SERVER_SESSION_STATE_CHANGED in the Status Flags and send session-state change data after a OK packet
// <h2>Client</h2>
// expects the server to send sesson-state changes after a OK packet
	CLIENT_SESSION_TRACK Capability = 0x00800000
)

