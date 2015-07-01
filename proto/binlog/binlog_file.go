package proto
import (
	"bufio"
	"bytes"
)

// Read binlog header and check
// If check failed,will not affect the reader
// http://dev.mysql.com/doc/internals/en/binlog-file-header.html
func CheckBinlogFileHeader(r *bufio.Reader) (bool, error) {
	b, err := r.Peek(4)
	if err != nil {return false, err}
	// \0xfe bin
	ret := bytes.Compare(b, []byte{0xfe, 0x62, 0x69, 0x6e}) == 0
	if ret { _, err = r.Read(b[0:4])}// consume the bin log
	return ret, err
}