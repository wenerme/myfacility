// generated by stringer -output=strings.go -type=Capability,Status,Command,SessionState,ProtoType,ColumnType,CursorType; DO NOT EDIT

package proto

import "fmt"

const _Capability_name = "CLIENT_LONG_PASSWORDCLIENT_FOUND_ROWSCLIENT_LONG_FLAGCLIENT_CONNECT_WITH_DBCLIENT_NO_SCHEMACLIENT_COMPRESSCLIENT_ODBCCLIENT_LOCAL_FILESCLIENT_IGNORE_SPACECLIENT_PROTOCOL_41CLIENT_INTERACTIVECLIENT_SSLCLIENT_IGNORE_SIGPIPECLIENT_TRANSACTIONSCLIENT_RESERVEDCLIENT_SECURE_CONNECTIONCLIENT_MULTI_STATEMENTSCLIENT_MULTI_RESULTSCLIENT_PS_MULTI_RESULTSCLIENT_PLUGIN_AUTHCLIENT_CONNECT_ATTRSCLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATACLIENT_CAN_HANDLE_EXPIRED_PASSWORDSCLIENT_SSL_VERIFY_SERVER_CERTCLIENT_DEPRECATE_EOFCLIENT_REMEMBER_OPTIONSCLIENT_BASIC_FLAGSCLIENT_ALL_FLAGS"

var _Capability_map = map[Capability]string{
	1:          _Capability_name[0:20],
	2:          _Capability_name[20:37],
	4:          _Capability_name[37:53],
	8:          _Capability_name[53:75],
	16:         _Capability_name[75:91],
	32:         _Capability_name[91:106],
	64:         _Capability_name[106:117],
	128:        _Capability_name[117:135],
	256:        _Capability_name[135:154],
	512:        _Capability_name[154:172],
	1024:       _Capability_name[172:190],
	2048:       _Capability_name[190:200],
	4096:       _Capability_name[200:221],
	8192:       _Capability_name[221:240],
	16384:      _Capability_name[240:255],
	32768:      _Capability_name[255:279],
	65536:      _Capability_name[279:302],
	131072:     _Capability_name[302:322],
	262144:     _Capability_name[322:345],
	524288:     _Capability_name[345:363],
	1048576:    _Capability_name[363:383],
	2097152:    _Capability_name[383:420],
	4194304:    _Capability_name[420:455],
	8388608:    _Capability_name[455:484],
	16777216:   _Capability_name[484:504],
	2147483648: _Capability_name[504:527],
	2155870175: _Capability_name[527:545],
	2164260863: _Capability_name[545:561],
}

func (i Capability) String() string {
	if str, ok := _Capability_map[i]; ok {
		return str
	}
	return fmt.Sprintf("Capability(%d)", i)
}

const _Status_name = "SERVER_STATUS_IN_TRANSSERVER_STATUS_AUTOCOMMITSERVER_MORE_RESULTS_EXISTSSERVER_QUERY_NO_GOOD_INDEX_USEDSERVER_QUERY_NO_INDEX_USEDSERVER_STATUS_CURSOR_EXISTSSERVER_STATUS_LAST_ROW_SENTSERVER_STATUS_DB_DROPPEDSERVER_STATUS_NO_BACKSLASH_ESCAPESSERVER_STATUS_METADATA_CHANGEDSERVER_QUERY_WAS_SLOWSERVER_STATUS_CLEAR_SETSERVER_PS_OUT_PARAMSSERVER_STATUS_IN_TRANS_READONLYSERVER_SESSION_STATE_CHANGED"

var _Status_map = map[Status]string{
	1:     _Status_name[0:22],
	2:     _Status_name[22:46],
	8:     _Status_name[46:72],
	16:    _Status_name[72:103],
	32:    _Status_name[103:129],
	64:    _Status_name[129:156],
	128:   _Status_name[156:183],
	256:   _Status_name[183:207],
	512:   _Status_name[207:241],
	1024:  _Status_name[241:271],
	2048:  _Status_name[271:292],
	3576:  _Status_name[292:315],
	4096:  _Status_name[315:335],
	8192:  _Status_name[335:366],
	16384: _Status_name[366:394],
}

func (i Status) String() string {
	if str, ok := _Status_map[i]; ok {
		return str
	}
	return fmt.Sprintf("Status(%d)", i)
}

const _Command_name = "COM_SLEEPCOM_QUITCOM_INIT_DBCOM_QUERYCOM_FIELD_LISTCOM_CREATE_DBCOM_DROP_DBCOM_REFRESHCOM_SHUTDOWNCOM_STATISTICSCOM_PROCESS_INFOCOM_CONNECTCOM_PROCESS_KILLCOM_DEBUGCOM_PINGCOM_TIMECOM_DELAYED_INSERTCOM_CHANGE_USERCOM_BINLOG_DUMPCOM_TABLE_DUMPCOM_CONNECT_OUTCOM_REGISTER_SLAVECOM_STMT_PREPARECOM_STMT_EXECUTECOM_STMT_SEND_LONG_DATACOM_STMT_CLOSECOM_STMT_RESETCOM_SET_OPTIONCOM_STMT_FETCHCOM_DAEMONCOM_BINLOG_DUMP_GTIDCOM_RESET_CONNECTION"

var _Command_index = [...]uint16{0, 9, 17, 28, 37, 51, 64, 75, 86, 98, 112, 128, 139, 155, 164, 172, 180, 198, 213, 228, 242, 257, 275, 291, 307, 330, 344, 358, 372, 386, 396, 416, 436}

func (i CommandType) String() string {
	if i >= CommandType(len(_Command_index)-1) {
		return fmt.Sprintf("Command(%d)", i)
	}
	return _Command_name[_Command_index[i]:_Command_index[i+1]]
}

const _SessionState_name = "SESSION_TRACK_SYSTEM_VARIABLESSESSION_TRACK_SCHEMASESSION_TRACK_STATE_CHANGE"

var _SessionState_index = [...]uint8{0, 30, 50, 76}

func (i SessionState) String() string {
	if i >= SessionState(len(_SessionState_index)-1) {
		return fmt.Sprintf("SessionState(%d)", i)
	}
	return _SessionState_name[_SessionState_index[i]:_SessionState_index[i+1]]
}

const _ProtoType_name = "UndTypeIntInt1Int2Int3Int4Int6Int8IntEncStrEofStrNulStrEncStrVarIgnoreByte"

var _ProtoType_index = [...]uint8{0, 7, 10, 14, 18, 22, 26, 30, 34, 40, 46, 52, 58, 64, 74}

func (i ProtoType) String() string {
	if i < 0 || i >= ProtoType(len(_ProtoType_index)-1) {
		return fmt.Sprintf("ProtoType(%d)", i)
	}
	return _ProtoType_name[_ProtoType_index[i]:_ProtoType_index[i+1]]
}

const (
	_ColumnType_name_0 = "MYSQL_TYPE_DECIMALMYSQL_TYPE_TINYMYSQL_TYPE_SHORTMYSQL_TYPE_LONGMYSQL_TYPE_FLOATMYSQL_TYPE_DOUBLEMYSQL_TYPE_NULLMYSQL_TYPE_TIMESTAMPMYSQL_TYPE_LONGLONGMYSQL_TYPE_INT24MYSQL_TYPE_DATEMYSQL_TYPE_TIMEMYSQL_TYPE_DATETIMEMYSQL_TYPE_YEARMYSQL_TYPE_NEWDATEMYSQL_TYPE_VARCHARMYSQL_TYPE_BITMYSQL_TYPE_TIMESTAMP2MYSQL_TYPE_DATETIME2MYSQL_TYPE_TIME2"
	_ColumnType_name_1 = "MYSQL_TYPE_NEWDECIMALMYSQL_TYPE_ENUMMYSQL_TYPE_SETMYSQL_TYPE_TINY_BLOBMYSQL_TYPE_MEDIUM_BLOBMYSQL_TYPE_LONG_BLOBMYSQL_TYPE_BLOBMYSQL_TYPE_VAR_STRINGMYSQL_TYPE_STRINGMYSQL_TYPE_GEOMETRY"
)

var (
	_ColumnType_index_0 = [...]uint16{0, 18, 33, 49, 64, 80, 97, 112, 132, 151, 167, 182, 197, 216, 231, 249, 267, 281, 302, 322, 338}
	_ColumnType_index_1 = [...]uint8{0, 21, 36, 50, 70, 92, 112, 127, 148, 165, 184}
)

func (i ColumnType) String() string {
	switch {
	case 0 <= i && i <= 19:
		return _ColumnType_name_0[_ColumnType_index_0[i]:_ColumnType_index_0[i+1]]
	case 246 <= i && i <= 255:
		i -= 246
		return _ColumnType_name_1[_ColumnType_index_1[i]:_ColumnType_index_1[i+1]]
	default:
		return fmt.Sprintf("ColumnType(%d)", i)
	}
}

const _CursorType_name = "CURSOR_TYPE_NO_CURSORCURSOR_TYPE_READ_ONLYCURSOR_TYPE_FOR_UPDATECURSOR_TYPE_SCROLLABLE"

var _CursorType_index = [...]uint8{0, 21, 42, 64, 86}

func (i CursorType) String() string {
	if i >= CursorType(len(_CursorType_index)-1) {
		return fmt.Sprintf("CursorType(%d)", i)
	}
	return _CursorType_name[_CursorType_index[i]:_CursorType_index[i+1]]
}
