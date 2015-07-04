package proto
import (
	"testing"
	"bytes"
	"compress/zlib"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/hex"
	"strings"
	"fmt"
)


func TestCompressSeveral(t *testing.T) {
	assert := assert.New(t)
	// SELECT repeat("a", 50)
	before := `
01 00 00 01 01 25 00 00    02 03 64 65 66 00 00 00    .....%....def...
0f 72 65 70 65 61 74 28    22 61 22 2c 20 35 30 29    .repeat("a", 50)
00 0c 08 00 32 00 00 00    fd 01 00 1f 00 00 05 00    ....2...........
00 03 fe 00 00 02 00 33    00 00 04 32 61 61 61 61    .......3...2aaaa
61 61 61 61 61 61 61 61    61 61 61 61 61 61 61 61    aaaaaaaaaaaaaaaa
61 61 61 61 61 61 61 61    61 61 61 61 61 61 61 61    aaaaaaaaaaaaaaaa
61 61 61 61 61 61 61 61    61 61 61 61 61 61 05 00    aaaaaaaaaaaaaa..
00 05 fe 00 00 02 00                                  .......
`
	after := `
4a 00 00 01 77 00 00 78    9c 63 64 60 60 64 54 65    J...w..x.cd..dTe
60 60 62 4e 49 4d 63 60    60 e0 2f 4a 2d 48 4d 2c    ..bNIMc.../J-HM,
d1 50 4a 54 d2 51 30 35    d0 64 e0 e1 60 30 02 8a    .PJT.Q05.d...0..
ff 65 64 90 67 60 60 65    60 60 fe 07 54 cc 60 cc    .ed.g..e....T...
c0 c0 62 94 48 32 00 ea    67 05 eb 07 00 8d f9 1c    ..b.H2..g.......
64                                                    d
`
	_, _ = before, after

	buf := bytes.NewBufferString("")
	buf.Write(DecodeDump(after))

	r := NewCompressedReader(buf)
	b, err := ioutil.ReadAll(r)
	assert.NoError(err)
	assert.EqualValues(DecodeDump(before), b)

	// write and read again
	w := NewCompressedWriter(buf)
	_, err =w.Write(DecodeDump(before))
	assert.NoError(err)
	b, err = ioutil.ReadAll(r)
	assert.NoError(err)
	assert.EqualValues(DecodeDump(before), b)
}

func TestCompressSeveralConcat(t *testing.T) {
	assert := assert.New(t)
	// SELECT repeat("a", 50)
	before := `
01 00 00 01 01 25 00 00    02 03 64 65 66 00 00 00    .....%....def...
0f 72 65 70 65 61 74 28    22 61 22 2c 20 35 30 29    .repeat("a", 50)
00 0c 08 00 32 00 00 00    fd 01 00 1f 00 00 05 00    ....2...........
00 03 fe 00 00 02 00 33    00 00 04 32 61 61 61 61    .......3...2aaaa
61 61 61 61 61 61 61 61    61 61 61 61 61 61 61 61    aaaaaaaaaaaaaaaa
61 61 61 61 61 61 61 61    61 61 61 61 61 61 61 61    aaaaaaaaaaaaaaaa
61 61 61 61 61 61 61 61    61 61 61 61 61 61 05 00    aaaaaaaaaaaaaa..
00 05 fe 00 00 02 00                                  .......
2e 00 00 00 03 73 65 6c    65 63 74 20 22 30 31 32    .....select "012
33 34 35 36 37 38 39 30    31 32 33 34 35 36 37 38    3456789012345678
39 30 31 32 33 34 35 36    37 38 39 30 31 32 33 34    9012345678901234
35 22                                                 5"
09 00 00 00 03 53 45 4c 45 43 54 20 31                ....SELECT 1
`
	after := `
4a 00 00 01 77 00 00 78    9c 63 64 60 60 64 54 65    J...w..x.cd..dTe
60 60 62 4e 49 4d 63 60    60 e0 2f 4a 2d 48 4d 2c    ..bNIMc.../J-HM,
d1 50 4a 54 d2 51 30 35    d0 64 e0 e1 60 30 02 8a    .PJT.Q05.d...0..
ff 65 64 90 67 60 60 65    60 60 fe 07 54 cc 60 cc    .ed.g..e....T...
c0 c0 62 94 48 32 00 ea    67 05 eb 07 00 8d f9 1c    ..b.H2..g.......
64                                                    d
22 00 00 00 32 00 00 78    9c d3 63 60 60 60 2e 4e    "...2..x..c....N
cd 49 4d 2e 51 50 32 30    34 32 36 31 35 33 b7 b0    .IM.QP20426153..
c4 cd 52 02 00 0c d1 0a    6c                         ..R.....l
0d 00 00 00 00 00 00 09    00 00 00 03 53 45 4c 45    ............SELE
43 54 20 31                                           CT 1
`
	_, _ = before, after

	buf := bytes.NewBufferString("")
	buf.Write(DecodeDump(after))

	r := NewCompressedReader(buf)
	b, err := ioutil.ReadAll(r)
	fmt.Println(hex.Dump(b))
	fmt.Println(hex.Dump(DecodeDump(after)))
	fmt.Println(hex.Dump(buf.Bytes()))
	assert.NoError(err)
	assert.EqualValues(DecodeDump(before), b)


	// write and read again
	w := NewCompressedWriter(buf)
	_, err =w.Write(DecodeDump(before))
	assert.NoError(err)
	b, err = ioutil.ReadAll(r)
	assert.NoError(err)
	assert.EqualValues(DecodeDump(before), b)
}

func TestUnCompressed(t *testing.T) {
	assert := assert.New(t)
	// SELECT repeat("a", 50)
	before := `
09 00 00 00 03 53 45 4c 45 43 54 20 31                ....SELECT 1
`
	after := `
0d 00 00 00 00 00 00 09    00 00 00 03 53 45 4c 45    ............SELE
43 54 20 31                                           CT 1
`
	_, _ = before, after

	buf := bytes.NewBufferString("")
	buf.Write(DecodeDump(after))

	r := NewCompressedReader(buf)
	b, err := ioutil.ReadAll(r)
	assert.NoError(err)
	assert.EqualValues(DecodeDump(before), b)


	// write and read again
	w := NewCompressedWriter(buf)
	_, err =w.Write(DecodeDump(before))
	assert.NoError(err)
	b, err = ioutil.ReadAll(r)
	assert.NoError(err)
	assert.EqualValues(DecodeDump(before), b)
}

func testCompress(t *testing.T) {
	assert := assert.New(t)
	before := `
2e 00 00 00 03 73 65 6c    65 63 74 20 22 30 31 32    .....select "012
33 34 35 36 37 38 39 30    31 32 33 34 35 36 37 38    3456789012345678
39 30 31 32 33 34 35 36    37 38 39 30 31 32 33 34    9012345678901234
35 22                                                 5"
`
	after := `
22 00 00 00 32 00 00 78    9c d3 63 60 60 60 2e 4e    "...2..x..c....N
cd 49 4d 2e 51 50 32 30    34 32 36 31 35 33 b7 b0    .IM.QP20426153..
c4 cd 52 02 00 0c d1 0a    6c                         ..R.....l
`
	_, _ = before, after
	{
		data := bytes.NewBuffer(DecodeDump(after))
		data.Read(make([]byte, 7))
		r, err := zlib.NewReader(data)
		assert.NoError(err)
		b, err := ioutil.ReadAll(r)
		assert.NoError(err)
		assert.EqualValues(DecodeDump(before), b)
	}
	{
		var data bytes.Buffer
		w, err := zlib.NewWriterLevel(&data, zlib.BestCompression)
		assert.NoError(err)
		h, err := hex.DecodeString(strings.Replace("22 00 00 00 32 00 00", " ", "", -1))
		assert.NoError(err)
		data.Write(h)
		w.Write(DecodeDump(before))
		w.Close()// Important

		fmt.Println(hex.Dump(data.Bytes()))
		fmt.Println(hex.Dump(DecodeDump(after)))
		//		fmt.Println(hex.Dump(DecodeDump(before)))

		// 不等,因为会刷一个 00 00 FF FF 的 deflate 块边界
		//		assert.EqualValues(DecodeDump(after), data.Bytes())
		data.Read(make([]byte, 7))
		fmt.Println(hex.Dump(data.Bytes()))
		r, err := zlib.NewReader(&data)
		assert.NoError(err)
		b, err := ioutil.ReadAll(r)
		assert.NoError(err)
		r.Close()
		assert.EqualValues(DecodeDump(before), b)
	}

}
