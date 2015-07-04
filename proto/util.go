package proto

import (
	"bytes"
	"encoding/hex"
	"github.com/op/go-logging"
	"os"
	"regexp"
	"strings"
)

var log = logging.MustGetLogger("proto")

// 初始化 Log
func init() {
	//	format := logging.MustStringFormatter("%{color}%{time:15:04:05} %{level:.4s} %{shortfunc} %{color:reset} %{message}", )
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{level:.4s} %{longfile} %{shortfunc} %{color:reset} %{message}")
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
	logging.SetLevel(logging.DEBUG, "proto")
}

// TODO Need optimize
func DecodeDump(dump string) (data []byte) {
	// 00000000  8d a6 0f 00 00 00 00 01  08 00 00 00 00 00 00 00  |................|
	// 33 34 35 36 37 38 39 30    31 32 33 34 35 36 37 38    3456789012345678
	dump = strings.TrimSpace(dump)
	lines := strings.Split(dump, "\n")
	buf := bytes.NewBufferString("")
	for _, l := range lines {
		{
			ok, err := regexp.MatchString("(?i)^[0-9a-z]{3,}", l)
			if err != nil {
				panic(err)
			}
			if ok {
				l = l[strings.Index(l, " "):]
			}
		}
		l = strings.TrimSpace(l)
		{
			reg := regexp.MustCompile(`(?i)^([0-9a-z]{2}\s+){16}`)
			tmp := reg.FindString(l)
			if tmp == "" {
				// Works for most of time
				r2 := regexp.MustCompile(`^.*\s{3,}`)
				tmp = r2.FindString(l)
			}
			l = tmp
		}
		l = strings.Replace(l, " ", "", -1)
		b, err := hex.DecodeString(l)
		if err != nil {
			panic(err)
		}
		buf.Write(b)
	}
	return buf.Bytes()
}
