package binlog

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/wenerme/myfacility/proto"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"testing"
)

func TestReader(t *testing.T) {
	f, err := os.Open("mysql-bin.sakila")
	if err != nil {
		panic(err)
	}
	err = ReadBinlog(f)
	if err != nil {
		panic(err)
	}
}

func TestReaderRead(t *testing.T) {
	assert := assert.New(t)
	_ = assert
	f, err := os.Open("mysql-bin.sakila")
	if err != nil {
		panic(err)
	}
	f.Seek(4, os.SEEK_SET)
	rd := NewReader(f)
	for {
		originHeader, err := rd.rd.Peek(19)
		if err != nil {
			panic(err)
		}
		eventSize := binary.LittleEndian.Uint32(originHeader[9:])
		originEvent, err := rd.rd.Peek(int(eventSize))
		originEvent = originEvent[19:]
		if err != nil {
			panic(err)
		}

		e, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		spew.Dump(e.Event())
		assertCodec(originHeader, e.Header(), t, rd.context)
		assertCodec(originEvent, e.Event(), t, rd.context)
		//		spew.Dump(e.Event())
	}
}
func writeToBytes(p interface{}, ctx *context) []byte {
	buf := &bytes.Buffer{}
	c := struct {
		proto.Writer
		*context
	}{proto.NewWriter(buf), ctx}

	if _, ok := p.(protoWritable); ok {
		p.(protoWritable).Write(c)
	} else if _, ok := p.(binlogWritable); ok {
		p.(binlogWritable).Write(c)
	} else {
		panic(fmt.Sprintf("No way to write %T", p))
	}
	return buf.Bytes()
}
func readFromBytes(p interface{}, b []byte, ctx *context) {
	buf := bytes.NewBuffer(b)
	c := struct {
		proto.Reader
		*context
	}{proto.NewReader(buf), ctx}
	if _, ok := p.(protoReadable); ok {
		p.(protoReadable).Read(c)
	} else if _, ok := p.(binlogReadable); ok {
		p.(binlogReadable).Read(c)
	} else {
		panic(fmt.Sprintf("No way to read %T", p))
	}
}

func TestEventHeader(t *testing.T) {
	p := &EventHeader{}
	assertCodec(`
00000000  5b 98 c8 51 0f 65 00 00  00 67 00 00 00 6b 00 00  |[..Q.e...g...k..|
00000010  00 01 00                                          |...|
 `, p, DumpPacket, DumpWrite, DumpOrigin)
}

type TestFlag uint

const (
	SkipEqual TestFlag = 1 << iota
	DumpOrigin
	DumpWrite
	DumpPacket
)

func (this TestFlag) Has(c TestFlag) bool {
	return this&c != 0
}

func (this TestFlag) Remove(c TestFlag) TestFlag {
	return this & ^c
}

func (this TestFlag) Add(c TestFlag) TestFlag {
	return this | c
}
func assertCodec(input interface{}, p interface{}, args ...interface{}) {
	var data []byte
	switch input := input.(type) {
	case string:
		data = proto.DecodeDump(input)
	case []byte:
		data = input
	}
	fine := false
	var t *testing.T
	flag := TestFlag(0)
	var write []byte
	var ctx *context

	for _, arg := range args {
		if f, ok := arg.(TestFlag); ok {
			flag = flag.Add(f)
		} else if f, ok := arg.(*testing.T); ok {
			t = f
		} else if f, ok := arg.(*context); ok {
			ctx = f
		}
	}

	if flag.Has(DumpOrigin) {
		fmt.Println("Origin data:\n", hex.Dump(data))
	}

	defer func() {
		if !fine {
			fmt.Println("Assert Codec Failed")
			fmt.Printf("Origin data:\n%s\n", hex.Dump(data))
			fmt.Printf("Write data:\n%s\n", hex.Dump(write))
			fmt.Println("Packet:")
			spew.Dump(p)
			fmt.Println(string(debug.Stack()))
			if t != nil {
				t.Fatal(recover())
			}
		}
	}()

	readFromBytes(p, data, ctx)

	if flag.Has(DumpPacket) {
		fmt.Printf("Packet:\n%#v\n", p)
	}

	for _, t := range args {
		if f, ok := t.(func()); ok {
			f()
		}
	}

	write = writeToBytes(p, ctx)
	if flag.Has(DumpWrite) {
		fmt.Printf("Write data:\n%s\n", hex.Dump(write))
	}
	if !flag.Has(SkipEqual) && !bytes.Equal(data, write) {
		panic("Not equals")
	}

	np := reflect.New(reflect.ValueOf(p).Elem().Type()).Interface()
	readFromBytes(np, write, ctx)
	if !flag.Has(SkipEqual) && !reflect.DeepEqual(p, np) {
		panic("Packet not equals")
	}
	fine = true
}
