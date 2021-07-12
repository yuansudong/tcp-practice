package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var port int
var callClose bool
var keepalive bool

func main() {
	fmt.Println("当前的PID是：", os.Getpid())
	// conn, err := net.DialTCP("tcp", nil, nil)
	flag.IntVar(&port, "port", 9999, "--port=5050")
	flag.BoolVar(&callClose, "call-close", false, "--call-close")
	flag.BoolVar(&keepalive, "keepalive", false, "--keepalive")
	flag.Parse()
	mListener, mErr := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: port,
	})
	if mErr != nil {
		fmt.Println(mErr.Error())
		os.Exit(-1)
	}
	// mListener.SetDeadline(5 * time.Second)
	for {
		mConn, mErr := mListener.AcceptTCP()
		if mErr != nil {
			fmt.Println(mErr)
			continue
		}
		mConn.SetKeepAlivePeriod(2 * time.Second)
		//  服务端
		if !keepalive {
			mConn.SetKeepAlive(false)
		}
		go handleConn(mConn)
	}
}

func handleConn(conn *net.TCPConn) {
	//  只为验证tcp_keepalive.不处理粘包等事情
	buf := make([]byte, 4096)
	fmt.Println("有新的连接：", conn.LocalAddr().String(), conn.RemoteAddr().String())
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			if callClose {
				conn.Close()
			}
			return
		}
	}
}
