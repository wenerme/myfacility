package binlog

import (
	"github.com/shopspring/decimal"
	"math"
)

const _DIG_PER_DEC = 9

var _DIG_TO_BYTES = []byte{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}

func determineDecimalLength(precision int, scale int) int {
	x := precision - scale
	ipDigits := x / _DIG_PER_DEC
	fpDigits := scale / _DIG_PER_DEC
	ipDigitsX := x - ipDigits*_DIG_PER_DEC
	fpDigitsX := scale - fpDigits*_DIG_PER_DEC
	return (ipDigits << 2) + int(_DIG_TO_BYTES[ipDigitsX]) + (fpDigits << 2) + int(_DIG_TO_BYTES[fpDigitsX])
}

// https://github.com/MariaDB/server/blob/10.1/strings/decimal.c
// decimal2bin
// bin2decimal
func toDecimal(precision int, scale int, value []byte) decimal.Decimal {
	positive := value[0]&0x80 == 0x80
	value[0] ^= 0x80
	if !positive {
		for i := 0; i < len(value); i++ {
			value[i] ^= 0xFF
		}
	}
	x := precision - scale
	ipDigits := x / _DIG_PER_DEC
	ipDigitsX := x - ipDigits*_DIG_PER_DEC
	ipSize := (ipDigits << 2) + int(_DIG_TO_BYTES[ipDigitsX])
	offset := int(_DIG_TO_BYTES[ipDigitsX])
	ip := decimal.Zero
	if offset > 0 {
		ip = decimal.NewFromFloat(float64(bigEndianInteger(value, 0, offset)))
	}
	for ; offset < ipSize; offset += 4 {
		i := bigEndianInteger(value, offset, 4)
		ip = ip.Mul(decimal.NewFromFloat(math.Pow10(_DIG_PER_DEC))).Add(decimal.NewFromFloat(float64(i)))
	}
	shift := 0
	fp := decimal.Zero
	for ; shift+_DIG_PER_DEC <= scale; shift, offset = shift+_DIG_PER_DEC, offset+4 {
		i := bigEndianInteger(value, offset, 4)
		fp = fp.Add(decimal.NewFromFloat(float64(i)).Div(decimal.NewFromFloat(math.Pow10(shift + _DIG_PER_DEC))))
	}

	if shift < scale {
		i := bigEndianInteger(value, offset, int(_DIG_TO_BYTES[scale-shift]))
		fp = fp.Add(decimal.NewFromFloat(float64(i)).Div(decimal.NewFromFloat(math.Pow10(scale))))
	}
	result := ip.Add(fp)
	if positive {
		return result
	}
	return result.Mul(decimal.NewFromFloat(-1))
}

func bigEndianInteger(bytes []byte, offset int, length int) uint64 {
	var result uint64
	for i := offset; i < (offset + length); i++ {
		b := bytes[i]
		result = (result << 8) | uint64(b)
	}
	return result
}
