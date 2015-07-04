package proto
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)


func TestQueryResponse(t *testing.T) {
	assert := assert.New(t)
	// select @@version_comment
	dump := `
00000000  01 00 00 01 01 27 00 00  02 03 64 65 66 00 00 00  |.....'....def...|
00000010  11 40 40 76 65 72 73 69  6f 6e 5f 63 6f 6d 6d 65  |.@@version_comme|
00000020  6e 74 00 0c 21 00 18 00  00 00 fd 00 00 1f 00 00  |nt..!...........|
00000030  05 00 00 03 fe 00 00 02  00 09 00 00 04 08 48 6f  |..............Ho|
00000040  6d 65 62 72 65 77 05 00  00 05 fe 00 00 02 00     |mebrew.........|
`
	_ = assert
	buf := bytes.NewBuffer(DecodeDump(dump))
	b := NewBuffer(buf, nil)
	b.SetCap(CLIENT_BASIC_FLAGS)
	rs := &QueryResponse{}
	rs.Read(b)
	fmt.Printf("%v\n", rs)
}
func TestQueryResponseLong(t *testing.T) {
	assert := assert.New(t)
	dump := `
01 00 00 01 02 54 00 00    02 03 64 65 66 12 69 6e    .....T....def.in
66 6f 72 6d 61 74 69 6f    6e 5f 73 63 68 65 6d 61    formation_schema
09 56 41 52 49 41 42 4c    45 53 09 56 41 52 49 41    .VARIABLES.VARIA
42 4c 45 53 0d 56 61 72    69 61 62 6c 65 5f 6e 61    BLES.Variable_na
6d 65 0d 56 41 52 49 41    42 4c 45 5f 4e 41 4d 45    me.VARIABLE_NAME
0c 21 00 c0 00 00 00 fd    01 00 00 00 00 4d 00 00    .!...........M..
03 03 64 65 66 12 69 6e    66 6f 72 6d 61 74 69 6f    ..def.informatio
6e 5f 73 63 68 65 6d 61    09 56 41 52 49 41 42 4c    n_schema.VARIABL
45 53 09 56 41 52 49 41    42 4c 45 53 05 56 61 6c    ES.VARIABLES.Val
75 65 0e 56 41 52 49 41    42 4c 45 5f 56 41 4c 55    ue.VARIABLE_VALU
45 0c 21 00 00 0c 00 00    fd 00 00 00 00 00 05 00    E.!.............
00 04 fe 00 00 22 00 1a    00 00 05 14 63 68 61 72    ....."......char
61 63 74 65 72 5f 73 65    74 5f 63 6c 69 65 6e 74    acter_set_client
04 75 74 66 38 1e 00 00    06 18 63 68 61 72 61 63    .utf8.....charac
74 65 72 5f 73 65 74 5f    63 6f 6e 6e 65 63 74 69    ter_set_connecti
6f 6e 04 75 74 66 38 1b    00 00 07 15 63 68 61 72    on.utf8.....char
61 63 74 65 72 5f 73 65    74 5f 72 65 73 75 6c 74    acter_set_result
73 04 75 74 66 38 1a 00    00 08 14 63 68 61 72 61    s.utf8.....chara
63 74 65 72 5f 73 65 74    5f 73 65 72 76 65 72 04    cter_set_server.
75 74 66 38 0e 00 00 09    0c 69 6e 69 74 5f 63 6f    utf8.....init_co
6e 6e 65 63 74 00 1a 00    00 0a 13 69 6e 74 65 72    nnect......inter
61 63 74 69 76 65 5f 74    69 6d 65 6f 75 74 05 32    active_timeout.2
38 38 30 30 0c 00 00 0b    07 6c 69 63 65 6e 73 65    8800.....license
03 47 50 4c 19 00 00 0c    16 6c 6f 77 65 72 5f 63    .GPL.....lower_c
61 73 65 5f 74 61 62 6c    65 5f 6e 61 6d 65 73 01    ase_table_names.
32 1b 00 00 0d 12 6d 61    78 5f 61 6c 6c 6f 77 65    2.....max_allowe
64 5f 70 61 63 6b 65 74    07 31 30 34 38 35 37 36    d_packet.1048576
18 00 00 0e 11 6e 65 74    5f 62 75 66 66 65 72 5f    .....net_buffer_
6c 65 6e 67 74 68 05 31    36 33 38 34 15 00 00 0f    length.16384....
11 6e 65 74 5f 77 72 69    74 65 5f 74 69 6d 65 6f    .net_write_timeo
75 74 02 36 30 13 00 00    10 10 71 75 65 72 79 5f    ut.60.....query_
63 61 63 68 65 5f 73 69    7a 65 01 30 14 00 00 11    cache_size.0....
10 71 75 65 72 79 5f 63    61 63 68 65 5f 74 79 70    .query_cache_typ
65 02 4f 4e 1d 00 00 12    08 73 71 6c 5f 6d 6f 64    e.ON.....sql_mod
65 13 53 54 52 49 43 54    5f 54 52 41 4e 53 5f 54    e.STRICT_TRANS_T
41 42 4c 45 53 15 00 00    13 10 73 79 73 74 65 6d    ABLES.....system
5f 74 69 6d 65 5f 7a 6f    6e 65 03 43 53 54 11 00    _time_zone.CST..
00 14 09 74 69 6d 65 5f    7a 6f 6e 65 06 53 59 53    ...time_zone.SYS
54 45 4d 1d 00 00 15 0c    74 78 5f 69 73 6f 6c 61    TEM.....tx_isola
74 69 6f 6e 0f 52 45 50    45 41 54 41 42 4c 45 2d    tion.REPEATABLE-
52 45 41 44 13 00 00 16    0c 77 61 69 74 5f 74 69    READ.....wait_ti
6d 65 6f 75 74 05 32 38    38 30 30 05 00 00 17 fe    meout.28800.....
00 00 22 00                                           ..
`
	_ = assert
	buf := bytes.NewBuffer(DecodeDump(dump))
	b := NewBuffer(buf, nil)
	b.SetCap(CLIENT_BASIC_FLAGS)
	rs := &QueryResponse{}
	rs.Read(b)
	spew.Dump(rs)
}