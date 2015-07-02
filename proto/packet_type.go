package proto

type PackType uint8
const (
	OK CommandType = 0
	ERR = 0xFF
	EOF = 0xFE
	HANDSHAKE_PACK = iota + 0xFF
	HANDSHAKE_RESPONSE_PACK
	QUERY_RESPONSE_PACK


)