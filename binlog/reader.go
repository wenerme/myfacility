package binlog

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/wenerme/myfacility/proto"
	"io"
)

var ErrFileHeader = errors.New("Wrong binlog file header")

func ReadBinlog(reader io.Reader) (err error) {
	c := &proto.BufReader{Reader: bufio.NewReader(reader)}
	buf := &bytes.Buffer{}
	r := &proto.BufReader{Reader: bufio.NewReader(buf)}
	tmp := make([]byte, 4)
	_, err = c.Read(tmp)
	if err != nil {
		return
	}

	// http://dev.mysql.com/doc/internals/en/binlog-file-header.html
	if bytes.Compare(tmp, []byte{0xfe, 0x62, 0x69, 0x6e}) != 0 {
		err = ErrFileHeader
		return
	}

	h := &EventHeader{}
	h.Read(c)
	spew.Dump(h)
	_, err = io.CopyN(buf, c, int64(h.EventSize))
	if err != nil {
		panic(err)
	}
	f := &FormatDescriptionEvent{}
	f.Read(r)
	spew.Dump(f)
	return
}
