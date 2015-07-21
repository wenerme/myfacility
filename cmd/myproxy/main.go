package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	. "github.com/wenerme/myfacility/proto"
	"io"
	"net"
	"os"
	"runtime/debug"
	"time"
)

const MYPROXY_VERSION = "0.0.1"

var daemon bool
var basedir string
var plugins []string
var app *cli.App
var log = logging.MustGetLogger("myproxy")

// 初始化 Log
func init() {
	//	format := logging.MustStringFormatter("%{color}%{time:15:04:05} %{level:.4s} %{shortfunc} %{color:reset} %{message}", )
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{level:.4s} %{longfile} %{shortfunc} %{color:reset} %{message}")
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
	logging.SetLevel(logging.DEBUG, "myproxy")
	flags()
}

func main() {
	app.Run(os.Args)
}

func TestClient() {
	c, err := net.Dial("tcp", "127.0.0.1:3306")
	if err != nil {
		panic(err)
	}
	_ = c
}
func run() {
	log.Info("Run myproxt at 127.0.0.1:8788")
	l, err := net.Listen("tcp", "127.0.0.1:8788")
	if err != nil {
		panic(err)
	}
	for {
		cli, err := l.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Accept %v \n", cli.RemoteAddr())
		svr, err := net.Dial("tcp", "127.0.0.1:3306")
		if err != nil {
			panic(err)
		}
		go proxy(svr, cli)
	}
}
func proxy(svr net.Conn, cli net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
			//			os.Exit(1)
		}
	}()
	cliBufReader := bufio.NewReader(cli)
	cliBufWriter := bufio.NewWriter(cli)
	cp := NewProto(bufio.NewReadWriter(cliBufReader, cliBufWriter), nil)
	sp := NewProto(svr, nil)
	_, _ = cp, sp
	phase := READ_HANDSHAKE
	var next Mode
	hs := &Handshake{}
	hsr := &HandshakeResponse{}
	com := &ComPack{}
	rs := &QueryResponse{}
	commands := NewCommandPacketMap()
	var p Pack
	var list []Pack
	var err error
	var expected, actual []byte
	_, _ = expected, actual
	for {
		switch phase {
		case READ_HANDSHAKE:
			_, err = sp.RecvReadPacket(hs)
			log.Info("Server caps:\n %v", hs.Capability.Dump())
			next = SEND_HANDSHAKE
		case SEND_HANDSHAKE:
			_, err = cp.WriteSendPacket(hs)
			next = READ_AUTH
		case READ_AUTH:
			_, err = cp.RecvReadPacket(hsr)
			cp.SetCap(hsr.Capability)
			sp.SetCap(hsr.Capability)
			log.Info("Client caps:\n %v", hsr.Capability.Dump())
			next = SEND_AUTH
		case SEND_AUTH:
			sp.SetSeq(cp.Seq())
			_, err = sp.WriteSendPacket(hsr)
			next = READ_AUTH_RESULT
		case READ_AUTH_RESULT:
			_, err = sp.RecvPacket()
			if err != nil {
				break
			}
			p, err = ReadErrOk(sp)
			next = SEND_AUTH_RESULT
		case SEND_AUTH_RESULT:
			cp.SetSeq(sp.Seq())
			_, err = cp.WriteSendPacket(p)
			// TODO check error
			next = READ_COM
		case READ_COM:
			_, err = cp.RecvReadPacket(com)
			next = SEND_COM
		case SEND_COM:
			sp.SetSeq(cp.Seq())
			_, err = sp.WriteSendPacket(com)
			log.Info("Handle command %v", com.Type)
			if log.IsEnabledFor(logging.DEBUG) {
				r := &BufReader{Reader: bufio.NewReader(bytes.NewReader(com.Data))}
				cmd := commands[com.Type]
				cmd.Read(r)
				spew.Dump(cmd)
			}
			sp.SetCom(com.Type)
			cp.SetCom(com.Type)
			switch com.Type {
			case COM_QUERY:
				next = READ_QUERY_RESULT
			case COM_FIELD_LIST:
				next = READ_FIELD_LIST
			default:
				next = READ_COM_RESULT
			}
		case READ_COM_RESULT:
			_, err = sp.RecvPacket()
			if err != nil {
				break
			}
			p, err = ReadErrOk(sp)
			next = SEND_COM_RESULT
		case SEND_COM_RESULT:
			cp.SetSeq(sp.Seq())
			_, err = cp.WriteSendPacket(p)
			next = READ_COM
		case READ_QUERY_RESULT:
			rs.Read(sp)
			next = SEND_QUERY_RESULT
		case SEND_QUERY_RESULT:
			cp.SetSeq(1)
			rs.Write(cp)
			next = READ_COM
		case READ_FIELD_LIST:
			// possible error
			list = nil
			for {
				sp.MustRecvPacket()
				if p, err = ReadErrEof(sp); err == nil {
					log.Info("Got %#v", p)
					list = append(list, p)
					break
				} else if err == ErrNotStatePack {
					err = nil
					//					b,err := sp.Peek(7)
					//					log.Info("%v %v\n%s",hex.Dump(b) , err)
					p = &ColumnDefinition{}
					p.Read(sp)
					log.Info("Got %#v", p)
					list = append(list, p)
				} else {
					break
				}
			}
			next = SEND_FIELD_LIST
		case SEND_FIELD_LIST:
			cp.SetSeq(1)
			spew.Dump(list)
			for _, p := range list {
				_, err = cp.WriteSendPacket(p)
			}
			next = READ_COM
		}
		cliBufWriter.Flush()
		phase = next
		if err, ok := p.(*ERRPack); false && ok {
			os.Stdout.Sync()
			time.Sleep(800)
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
			os.Exit(1)
		}
		if err == io.EOF {
			if _, e := cliBufReader.Peek(1); e == nil {
				panic(err)
			}
			log.Info("%s closed", cli.RemoteAddr())
			svr.Close()
			cli.Close()
			return
		}
		if err != nil {
			panic(err)
		}
	}
}

type Mode int

const (
	READ_HANDSHAKE Mode = iota + 1
	SEND_HANDSHAKE
	READ_AUTH
	SEND_AUTH
	READ_AUTH_RESULT
	SEND_AUTH_RESULT
	READ_COM
	SEND_COM
	READ_COM_RESULT
	SEND_COM_RESULT
	READ_QUERY_RESULT
	SEND_QUERY_RESULT
	READ_FIELD_LIST
	SEND_FIELD_LIST

	START_BINLOG_DUMP
)

func flags() {
	app = cli.NewApp()
	app.Name = "myproxy"
	app.Usage = "Proxy for MySQL-Server in Go"
	app.Version = MYPROXY_VERSION
	app.Authors = []cli.Author{
		{"wener", "wenermail@gmail.com"},
	}
	app.Action = func(c *cli.Context) {
		log.Info("Listening %v", c.StringSlice("proxy-address"))
		if c.Bool("no-proxy") {
			log.Debug("With no-proxyf flag, will not start proxy module")
			return
		}
		run()
	}
	// mysql-proxy Admin Options
	adminOption := []cli.Flag{
		cli.StringFlag{
			Name:  "admin-address",
			Usage: "The admin module listening host and port",
		},
		cli.StringFlag{
			Name:  "admin-lua-script",
			Usage: "Script to execute by the admin module",
		},
		cli.StringFlag{
			Name:  "admin-password",
			Usage: "Authentication password for admin module",
		},
		cli.StringFlag{
			Name:  "admin-username",
			Usage: "Authentication user name for admin module",
		},
		cli.StringSliceFlag{
			Name:  "proxy-address,P",
			Value: &cli.StringSlice{},
			Usage: "The listening proxy server host and port",
		},
	}
	// mysql-proxy Proxy Options
	proxyOption := []cli.Flag{
		cli.BoolFlag{
			Name:  "no-proxy",
			Usage: "Do not start the proxy module",
		},
		cli.StringSliceFlag{
			Name:  "proxy-backend-addresses,b",
			Value: &cli.StringSlice{},

			Usage: "The MySQL server host and port",
		},
		cli.BoolFlag{
			Name: "proxy-fix-bug-25371",
			Usage: "Enable the fix for Bug #25371 for older libmysql versions	0.8.1",
		},
		cli.StringFlag{
			Name:  "proxy-lua-script",
			Usage: "Filename for Lua script for proxy operations",
		},
		cli.BoolFlag{
			Name:  "proxy-pool-no-change-user",
			Usage: "Do not use the protocol CHANGE_USER command to reset the connection when coming from the connection pool",
		},
		cli.StringSliceFlag{
			Name: "proxy-read-only-backend-addresses,r", Value: &cli.StringSlice{},
			Usage: "The MySQL server host and port (read only)",
		},
		cli.BoolFlag{
			Name:  "proxy-skip-profiling",
			Usage: "Disable query profiling",
		},
	}
	// mysql-proxy Applications Options
	applicationOptions := []cli.Flag{
		cli.StringFlag{
			Name: "basedir",

			Usage: "The base directory prefix for paths in the configuration",
		},
		cli.BoolFlag{
			Name: "daemon",

			Usage: "Start in daemon mode",
		},
		cli.StringFlag{
			Name:  "defaults-file",
			Usage: "Read only named option file",
		},
		cli.IntFlag{
			Name: "event-threads", Value: 1,
			Usage: "Number of event-handling threads",
		},
		cli.BoolFlag{
			Name:  "keepalive",
			Usage: "Try to restart the proxy if a crash occurs",
		},
		cli.BoolFlag{
			Name:  "log-backtrace-on-crash",
			Usage: "Try to invoke the debugger and generate a backtrace on crash",
		},
		cli.StringFlag{Name: "log-file",
			Usage: "The file where error messages are logged",
		},
		// critical error warning info message debug
		cli.StringFlag{Name: "log-level", Value: "critical",
			Usage: "The logging level",
		},
		cli.BoolFlag{Name: "log-use-syslog",
			Usage: "Log errors to syslog",
		},
		cli.StringFlag{
			Name:  "lua-cpath",
			Usage: "Set the LUA_CPATH", EnvVar: "LUA_CPATH",
		},
		cli.StringFlag{
			Name:  "lua-path",
			Usage: "Set the LUA_PATH", EnvVar: "LUA_PATH",
		},
		cli.IntFlag{Name: "max-open-files",
			Usage: "The maximum number of open files to support",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Usage: "File in which to store process ID",
		},
		cli.StringFlag{
			Name:  "plugin-dir",
			Usage: "Directory containing plugin files",
		},
		cli.StringSliceFlag{
			Name: "plugins", Value: &cli.StringSlice{},
			Usage: "List of plugins to load",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "The user to use when running mysql-proxy",
		},
	}

	app.Flags = append(app.Flags, adminOption...)
	app.Flags = append(app.Flags, proxyOption...)
	app.Flags = append(app.Flags, applicationOptions...)

	app.Commands = nil
}
