package binlog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/wenerme/myfacility/proto"
	"io"
	"os"
	"runtime/debug"
)

var ErrFileHeader = errors.New("Wrong binlog file header")

func ReadBinlog(rd io.Reader) (err error) {
	rdBuf := bufio.NewReaderSize(rd, 19)
	c := &proto.BufReader{Reader: rdBuf}
	buf := &bytes.Buffer{}
	r := &reader{Reader: &proto.BufReader{Reader: bufio.NewReader(buf)}, tables: make(map[uint64]*TableMapEvent)}
	{
		// Check file header
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
	}

	h := &EventHeader{}
	m := NewEventMap()
	type readable interface {
		Read(proto.Reader)
	}
	type binlogReadable interface {
		Read(Reader)
	}
	for {
		// Make sure buffer is enough
		// FIXME If buf length is not enough, the header will not read all
		// bufio.Read only read once
		_, _ = c.Peek(19)
		h.Read(c)
		_, err = io.CopyN(buf, c, int64(h.EventSize-19))
		if err != nil {
			return
		}
		p := m[h.EventType]
		if p == nil {
			fmt.Println("Skip event ", h.EventType)
			buf.Reset()
			continue
		}
		spew.Dump(buf.Bytes())
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("============================================")
					fmt.Println(err)
					spew.Dump(h, p)
					n, _ := rd.(io.Seeker).Seek(0, os.SEEK_CUR)
					fmt.Printf("Current File Seek %d\n", n)
					debug.PrintStack()
					os.Exit(1)
				}
			}()
			if _, ok := p.(readable); ok {
				p.(readable).Read(r)
			}
			if _, ok := p.(binlogReadable); ok {
				p.(binlogReadable).Read(r)
			}
		}()
		if h.EventType == TABLE_MAP_EVENT {
			tab := p.(*TableMapEvent)
			r.SetTableMap(tab)
		}
		spew.Dump(p)
		if r.More() {
			panic("Should no more")
		}
	}
	return
}
