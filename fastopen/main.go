package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var port int
var run string
var addr string

func main() {
	flag.IntVar(&port, "port", 9999, "--port=9999")
	flag.StringVar(&run, "run", "service", "--run=service 或者 --run=client")
	flag.StringVar(&addr, "addr", "172.81.209.185", "--addr=172.81.209.185")
	flag.Parse()

	if run == "service" {
		RunServer()
	} else {
		RunClient()
	}

}
func RunServer() {
	var serverAddr [4]byte
	IP := net.ParseIP("0.0.0.0")
	if IP == nil {
		fmt.Println("IP地址无法解析", addr)
		return
	}
	copy(serverAddr[:], IP[12:16])
	server := TFOServer{ServerAddr: serverAddr, ServerPort: port}
	err := server.Bind()
	if err != nil {
		fmt.Printf("Failed to bind socket:%s\n", err.Error())
		return
	}
	server.Accept()

}

// RunClient
func RunClient() {
	var serverAddr [4]byte
	IP := net.ParseIP(addr)
	if IP == nil {
		fmt.Println("IP地址无法解析", addr)
		return
	}

	copy(serverAddr[:], IP[12:16])
	client := New(serverAddr, port)
	if err := client.Connect(); err != nil {
		fmt.Println("初始化套接字失败：", err)
		return
	}
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := client.Read(buf)
			if err != nil || n == 0 {
				if n == 0 {
					fmt.Println("服务端关闭了套接字")
				} else {
					fmt.Println("client read:", err.Error())
				}
				return
			} else {
				fmt.Printf("%s\n", string(buf[0:n]))
			}
			time.Sleep(time.Second)
		}
	}()
	sb := []byte("client:画画的北鼻")
	for {
		n, err := client.Send(sb)
		if err != nil {
			if n == 0 {
				fmt.Println("服务端关闭了套接字")
			} else {
				fmt.Println("client write:", err.Error())
			}
			return
		}
		time.Sleep(time.Second)

	}
}
