package proto
import "net"

type Conn interface {
	SystemVariable(string) (string, error)
	Ping() error
	Codec() Codec
	Exec(string) error
}
func NewConn(host string, port int, user string, password string) Conn {
	//    return &conn{}
	return nil
}
type ConCfg struct {
	User     string
	Password string
	Host     string
	Port     int
	Attrs    map[string]string
}
type shared struct {
	packet *Packet
}
type conn struct {
	conn  net.Conn
	codec Codec
	cfg   ConCfg
}