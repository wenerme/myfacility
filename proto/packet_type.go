package proto

type PackType uint8
const (
	OK CommandType = 0
	ERR CommandType = 0xFF
	EOF CommandType = 0xFE
)