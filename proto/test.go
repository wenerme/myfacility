package proto
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"io/ioutil"
)

type TestFlag uint
const (
	SkipEqual TestFlag = 1 << iota
	DumpOrigin
	DumpWrite
	DumpPayload
	DumpPacket
)
func (this TestFlag) Has(c TestFlag) bool {
	return this & c != 0
}

func (this TestFlag) Remove(c TestFlag) TestFlag {
	return this & ^c
}

func (this TestFlag) Add(c TestFlag) TestFlag {
	return this | c
}
func assertCodec(data []byte, p Pack, c Capability, args...interface{}) {
	fine := false
	var t *testing.T
	flag := TestFlag(0)
	var write []byte

	proto := NewBuffer(bytes.NewBufferString(""), nil)
	proto.buf.Write(data)

	for _, arg := range args {
		if f, ok := arg.(TestFlag); ok {
			flag = flag.Add(f)
		}else if f, ok := arg.(*testing.T); ok {
			t = f
		}
	}

	if flag.Has(DumpOrigin) {
		fmt.Println("Origin data:\n", hex.Dump(data))
	}

	proto.SetCap(c)

	//	if flag.Has(DumpPayload) {
	//		fmt.Printf("Payload :\n%s\n", hex.Dump(payload))
	//	}

	defer func() {
		if !fine {
			fmt.Println("Assert Codec Failed")
			fmt.Printf("Origin data:\n%s\n", hex.Dump(data))
			fmt.Printf("Write data:\n%s\n", hex.Dump(write))
			//			fmt.Printf("Payload :\n%s\n", hex.Dump(payload))
			fmt.Printf("Packet:\n%#v\n", p)
			if t != nil {
				t.Fatal(recover())
			}
		}
	}()

	_, err := proto.RecvPacket()
	if err != nil {panic(err)}
	proto.ReadPacket(p)

	if flag.Has(DumpPacket) {
		fmt.Printf("Packet:\n%#v\n", p)
	}

	for _, t := range args {
		if f, ok := t.(func()); ok {
			f()
		}
	}

	proto.WritePacket(p)
	_, err = proto.SendPacket()
	if err != nil {panic(err)}
	write, _ = ioutil.ReadAll(proto.con)
	if flag.Has(DumpWrite) {
		fmt.Printf("Write data:\n%s\n", hex.Dump(write))
	}
	if !flag.Has(SkipEqual) && !bytes.Equal(data, write) {
		panic("Not equals")
	}
	fine = true
}