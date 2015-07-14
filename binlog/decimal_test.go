package binlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimalToBin(t *testing.T) {
	// Example from decimal.c
	assert := assert.New(t)
	v := []byte{0x81, 0x0D, 0xFB, 0x38, 0xD2, 0x04, 0xD2}
	assert.EqualValues("1234567890.1234", toDecimal(14, 4, v).String())
	nv := []byte{0x7E, 0xF2, 0x04, 0xC7, 0x2D, 0xFB, 0x2D}
	assert.EqualValues("-1234567890.1234", toDecimal(14, 4, nv).String())
}
