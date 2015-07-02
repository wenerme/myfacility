package proto
import (
	"testing"
	"net"
	"bufio"
	"encoding/hex"
	"bytes"
	"fmt"
)

func TestServer(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:8788")
	if err != nil { panic(err) }
	for {
		svr, err := l.Accept()
		if err != nil { panic(err) }
		go func() {
			cli, err := net.Dial("tcp", "127.0.0.1:3306")
			if err != nil { panic(err) }
			svrReader := &BufReader{Reader: bufio.NewReader(svr)}
			svrWriter := &BufWriter{Writer: bufio.NewWriter(svr)}
			cliReader := &BufReader{Reader: bufio.NewReader(cli)}
			cliWriter := &BufWriter{Writer: bufio.NewWriter(cli)}

			{
				hs := &Handshake{}
				hs.Read(svrReader)
				fmt.Printf("%#v\n", hs)
				hs.Write(cliWriter)
			}
			{
				hr := HandshakeResponse{}
				hr.Read(cliReader)
				fmt.Printf("%#v\n", hr)
				hr.Write(svrWriter)
			}



		}()
	}
}

func TestClient(t *testing.T) {
	c, err := net.Dial("tcp", "127.0.0.1:3306")
	if err != nil { panic(err) }
	cr := &BufReader{Reader:bufio.NewReader(c)}
	p := Handshake{}
	data, n, err := cr.ReadPacket()
	print(hex.Dump(data))
	log.Debug("%v %v", n, err)
	pack := bytes.NewReader(data)
	r := &BufReader{Reader:bufio.NewReader(pack)}
	p.Read(r)
	log.Debug("%#v", p)
}