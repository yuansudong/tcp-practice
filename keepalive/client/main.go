package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var port int
var host string
var callClose bool
var keepalive bool

func main() {
	fmt.Println("当前的PID是：", os.Getpid())
	// conn, err := net.DialTCP("tcp", nil, nil)
	flag.IntVar(&port, "port", 9999, "--port=5050")
	flag.StringVar(&host, "host", "127.0.0.1", "--host=127.0.0.1")
	flag.BoolVar(&callClose, "call-close", false, "--call-close")

	flag.BoolVar(&keepalive, "keepalive", false, "--keepalive")
	flag.Parse()

	conn, err := net.DialTCP("tcp", &net.TCPAddr{}, &net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	conn.SetKeepAlivePeriod(2 * time.Second)
	if keepalive {
		conn.SetKeepAlive(true)
	} else {
		conn.SetKeepAlive(false)
	}
	buf := make([]byte, 4096)
	fmt.Println("连接建立成功")
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			if callClose {
				conn.Close()
			}
		}
		time.Sleep(time.Second)
	}

}
