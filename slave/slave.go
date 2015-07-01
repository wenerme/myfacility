package slave

type Slave interface {
	Start() error
	Stop() error
	Reset() error
	Init() error
	ChangeMaster(MasterInfo) error
	SlaveStatus() (SlaveStatus, error)
}

type MasterInfo struct {
	Host     string
	User     string
	Password string
	Port     int
	LogName  string
	LogPos   uint64
}
type SlaveStatus struct {

}

type MySQL struct {
	AffectedRows     uint64
	InsertId         uint64
	ClientFlag       uint32
	ServerCapability uint32
	ProtocolVersion  uint32
	ServerStatus     uint32
	WarningCount     uint32
	Status           uint32

}

type (
	MySQLStatus uint
	ProtocolType uint

)
const (
	MYSQL_STATUS_READY MySQLStatus = iota
	MYSQL_STATUS_GET_RESULT, MYSQL_STATUS_USE_RESULT,
	MYSQL_STATUS_STATEMENT_GET_RESULT
)
const (
	MYSQL_PROTOCOL_DEFAULT ProtocolType = iota
	MYSQL_PROTOCOL_TCP, MYSQL_PROTOCOL_SOCKET,
	MYSQL_PROTOCOL_PIPE, MYSQL_PROTOCOL_MEMORY
)
