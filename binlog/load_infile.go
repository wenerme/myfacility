package binlog

import "github.com/wenerme/myfacility/proto"

/*
LOAD INFILE replication
LOAD DATA|XML INFILE is a special SQL statement as it has to ship the files over to the slave too to execute the statement.
*/

/*
4              slave_proxy_id
4              exec_time
4              skip_lines
1              table_name_len
1              schema_len
4              num_fields

1              field_term
1              enclosed_by
1              line_term
1              line_start
1              escaped_by
1              opt_flags
1              empty_flags

string.var_len [len=1//num_fields] (array of 1-byte) field_name_lengths
string.var_len [len=sum(field_name_lengths) + num_fields] (array of nul-terminated strings) field_names
string.var_len [len=table_name_len + 1] (nul-terminated string) table_name
string.var_len [len=schema_len + 1] (nul-terminated string) schema_name
string.NUL     file_name
*/

type LoadEvent struct {
	SlaveProxyId  uint32
	ExecutionTime uint32
	SkipLines     uint32
}

func (p *LoadEvent) Read(c proto.Reader) {
	// TODO
}
func (p *LoadEvent) EventType() EventType {
	return LOAD_EVENT
}

/*
4              slave_proxy_id
4              exec_time
4              skip_lines
1              table_name_len
1              schema_len
4              num_fields

1              field_term_len
string.var_len field_term
1              enclosed_by_len
string.var_len enclosed_by
1              line_term_len
string.var_len line_term
1              line_start_len
string.var_len line_start
1              escaped_by_len
string.var_len escaped_by
1              opt_flags

string.var_len [len=1//num_fields] (array of 1-byte) field_name_lengths
string.var_len [len=sum(field_name_lengths) + num_fields] (array of nul-terminated strings) field_names
string.var_len [len=table_name_len] (nul-terminated string) table_name
string.var_len [len=schema_len] (nul-terminated string) schema_name
string.EOF     file_name
*/
type NewLoadEvent struct {
	SlaveProxyId  uint32
	ExecutionTime uint32
	SkipLines     uint32
}

func (p *NewLoadEvent) Read(c proto.Reader) {
	// TODO
}
func (p *NewLoadEvent) EventType() EventType {
	return NEW_LOAD_EVENT
}

type CreateFileEvent struct {
	FileId    uint32
	BlockData []byte
}

func (p *CreateFileEvent) Read(c proto.Reader) {
	c.Get(&p.FileId, &p.BlockData, proto.StrEof)
}
func (p *CreateFileEvent) EventType() EventType {
	return CREATE_FILE_EVENT
}

type AppendBlockEvent CreateFileEvent

func (p *AppendBlockEvent) EventType() EventType {
	return APPEND_BLOCK_EVENT
}

type BeginLoadQueryEvent CreateFileEvent

func (p *BeginLoadQueryEvent) EventType() EventType {
	return BEGIN_LOAD_QUERY_EVENT
}

type DeleteFileEvent struct {
	FileId uint32
}

func (p *DeleteFileEvent) Read(c proto.Reader) {
	c.Get(&p.FileId)
}
func (p *DeleteFileEvent) EventType() EventType {
	return DELETE_FILE_EVENT
}

type ExecLoadEvent struct {
	FieldId uint32
}

func (p *ExecLoadEvent) Read(c proto.Reader) {
	c.Get(&p.FieldId)
}
func (p *ExecLoadEvent) EventType() EventType {
	return EXEC_LOAD_EVENT
}

type ExecuteLoadQueryEvent struct {
	SlaveProxyId     uint32
	ExecutionTime    uint32
	ErrorCode        uint16
	Status           []byte
	Schema           string
	ExecuteLoadQuery string
}

func (p *ExecuteLoadQueryEvent) Read(c proto.Reader) {
	var m uint8
	var n uint16
	c.Get(
		&p.SlaveProxyId,
		&p.ExecutionTime,
		&m,
		&p.ErrorCode,
		&n,
		&p.Status, proto.StrVar, &n,
		&p.Schema, proto.StrVar, &m,
		&p.ExecuteLoadQuery)
}
func (p *ExecuteLoadQueryEvent) EventType() EventType {
	return EXECUTE_LOAD_QUERY_EVENT
}
