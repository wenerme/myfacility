package proto

//
// A COM_QUERY is used to send the server a text-based query that is executed immediately.
// <p/>
// The server replies to a COM_QUERY packet with a COM_QUERY Response.
// <p/>
// The length of the query-string is a taken from the packet length - 1.
// <p/>
// Payload
// 1              [03] COM_QUERY
// string[EOF]    the query the server shall execute
// Fields
// command_id (1) -- 0x03 COM_QUERY
// <p/>
// query (string.EOF) -- query_text
// <p/>
// Implemented By
// mysql_query()
// <p/>
// Returns
// COM_QUERY_Response
// <p/>
// Example
// 21 00 00 00 03 73 65 6c    65 63 74 20 40 40 76 65    !....select @@ve
// 72 73 69 6f 6e 5f 63 6f    6d 6d 65 6e 74 20 6c 69    rsion_comment li
// 6d 69 74 20 31                                        mit 1
type ComQuery struct {
	Query string
}

type ComPack struct {
	Type Command
	data []byte
	buf  *Buffer
}

func (p *ComPack) Read(c Reader) {
	c.Get(&p.Type, &p.data, StrEof)
}
func (p *ComPack) Write(c Writer) {
	c.Put(&p.Type, &p.data, StrEof)
}
