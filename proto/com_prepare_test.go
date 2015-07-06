package proto

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComPrepareExecute(t *testing.T) {
	dump := `
12 00 00 00 17 01 00 00    00 00 01 00 00 00 00 01    ................
0f 00 03 66 6f 6f                                     ...foo`
	p := &ComStmtExecute{}
	p.OK = &ComStmtPrepareOK{Params: make([]ColumnDefinition, 1)}
	assertCodec(DecodeDump(dump), p, CLIENT_BASIC_FLAGS, func() {}, t)
}
func TestComPrepareOK(t *testing.T) {
	assert := assert.New(t)
	dump := `
0c 00 00 01 00 01 00 00    00 01 00 02 00 00 00 00   ................
17 00 00 02 03 64 65 66    00 00 00 01 3f 00 0c 3f    .....def....?..?
00 00 00 00 00 fd 80 00    00 00 00 17 00 00 03 03    ................
64 65 66 00 00 00 01 3f    00 0c 3f 00 00 00 00 00    def....?..?.....
fd 80 00 00 00 00 05 00    00 04 fe 00 00 02 00 1a    ................
00 00 05 03 64 65 66 00    00 00 04 63 6f 6c 31 00    ....def....col1.
0c 3f 00 00 00 00 00 fd    80 00 1f 00 00 05 00 00    .?..............
06 fe 00 00 02 00                                     ......
`
	data := bytes.NewBuffer(DecodeDump(dump))
	b := NewBuffer(data, nil)
	b.SetCap(CLIENT_BASIC_FLAGS)
	p := ComStmtPrepareOK{}
	p.Read(b)
	b.SetSeq(1)
	p.Write(b)
	//	fmt.Printf("ORG\n%s\nWRITE\n%s\n",hex.Dump(DecodeDump(dump)), hex.Dump(data.Bytes()))
	assert.EqualValues(DecodeDump(dump), data.Bytes())
	//
	//	for i:=0; i < 255; i ++{
	//		fmt.Println(Command(i))
	//	}
}
