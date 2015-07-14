package binlog

import (
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	f, err := os.Open("mysql-bin.sakila")
	if err != nil {
		panic(err)
	}
	err = ReadBinlog(f)
	if err != nil {
		panic(err)
	}
}
