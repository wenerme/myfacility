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

type protoReadable interface {
	Read(proto.Reader)
}
type binlogReadable interface {
	Read(Reader)
}
type protoWritable interface {
	Write(proto.Writer)
}
type binlogWritable interface {
	Write(Writer)
}

func ReadBinlog(rd io.Reader) (err error) {
	rdBuf := bufio.NewReaderSize(rd, 19)
	c := &proto.BufReader{Reader: rdBuf}
	buf := &bytes.Buffer{}
	r := &reader{Reader: &proto.BufReader{Reader: bufio.NewReader(buf)}, context: newContext()}
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
	m := NewEventTypeMap()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("============================================")
			spew.Dump(h)
			fmt.Println(err)
			n, _ := rd.(io.Seeker).Seek(0, os.SEEK_CUR)
			fmt.Printf("Current File Seek %d\n", n)
			debug.PrintStack()
			os.Exit(1)
		}
	}()
	for {
		if !c.More() {
			fmt.Println("No more")
			os.Exit(0)
		}
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
		//		spew.Dump(buf.Bytes())
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("============================================")
					spew.Dump(h, p)
					fmt.Println(err)
					n, _ := rd.(io.Seeker).Seek(0, os.SEEK_CUR)
					fmt.Printf("Current File Seek %d\n", n)
					debug.PrintStack()
					os.Exit(1)
				}
			}()
			if _, ok := p.(protoReadable); ok {
				p.(protoReadable).Read(r)
			} else if _, ok := p.(binlogReadable); ok {
				p.(binlogReadable).Read(r)
			} else {
				panic(fmt.Sprintf("No way to read %T", p))
			}
		}()
		if h.EventType == TABLE_MAP_EVENT {
			tab := p.(*TableMapEvent)
			r.SetTableMap(*tab)
		}
		//		spew.Dump(p)
		if r.More() {
			panic("Should no more")
		}
	}
	return
}
