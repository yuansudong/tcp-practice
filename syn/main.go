package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"runtime"
	"syscall"
	"time"
)

var port int
var run string
var addr string
var connCount int

func main() {
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

//

func RunClient() {
	arr := []net.Conn{}
	host := addr + ":" + fmt.Sprint(port)
	for i := 0; i < connCount; i++ {
		go func() {
			conn, err := net.Dial("tcp4", host)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				arr = append(arr, conn)
			}
		}()
	}
	for {
		time.Sleep(5 * time.Second)
	}
}

const LISTEN_BACKLOG = 15

func Accept1(addr [4]byte, port int) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		if err == syscall.ENOPROTOOPT {
			err = errors.New("内核不支持哇")
		}
		return err
	}
	sa := &syscall.SockaddrInet4{Addr: addr, Port: port}
	err = syscall.Bind(fd, sa)
	if err != nil {
		return err
	}
	err = syscall.Listen(fd, LISTEN_BACKLOG)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)
	for {
		time.Sleep(7200 * time.Second)
		_, _, err := syscall.Accept(fd)
		if err != nil {
			return err
		}
		fmt.Println("接收了套接字。")
		time.Sleep(15 * time.Second)
	}
}

func RunService() {
	var serverAddr [4]byte
	IP := net.ParseIP("172.17.0.15")
	if IP == nil {
		fmt.Println("IP地址无法解析", addr)
		return
	}
	copy(serverAddr[:], IP[12:16])
	if err := Accept1(serverAddr, port); err != nil {
		fmt.Println(err.Error())
	}
	// listener, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(port))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println("端口", port)
	// for {
	// 	_, err := listener.Accept()

	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	fmt.Println("接收了套接字。")
	// 	time.Sleep(15 * time.Second)
	// }
}
