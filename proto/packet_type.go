package proto

type PackType uint8

const (
	OK PackType = iota
	HANDSHAKE
	HANDSHAKE_RESPONSE
	QUERY_RESPONSE
	COMMAND
	ERR PackType = 0xFF
	EOF PackType = 0xFE
)
