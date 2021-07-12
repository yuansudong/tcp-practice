package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var port int
var run string
var addr string
var connCount int

func main() {
	fmt.Println("当前的进程号是:", os.Getpid())
	flag.IntVar(&connCount, "conn-count", 1, "--conn-count=5")
	flag.StringVar(&addr, "addr", "172.81.209.185", "--addr=172.81.209.185")
	flag.StringVar(&run, "run", "client", "--run=client 或者 --run=service")
	flag.IntVar(&port, "port", 12345, "--port=9999")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if run == "service" {
		RunService()
	} else {
		RunClient()
	}
}

func RunClient() {
	conns := []net.Conn{}

	for index := 0; index < 1; index++ {

		conn, err := net.Dial("tcp4", addr+":"+fmt.Sprint(port))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		fmt.Println("当前的index是：", index)
		conns = append(conns, conn)
	}
	fmt.Println("连接建立成功")
	sChan := make(chan os.Signal, 10)
	signal.Notify(sChan, syscall.SIGHUP, syscall.SIGUSR2)
	for {
		select {
		case sig := <-sChan:
			if sig == syscall.SIGUSR2 {
				fmt.Println("调用关闭函数")
				for _, conn := range conns {
					if err := conn.Close(); err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}

}

func RunService() {
	listener, err := net.Listen("tcp4", "172.17.0.15:"+fmt.Sprint(port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("端口", port)
	for {
		_, err := listener.Accept()

		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("接收了套接字。")
	}
}
