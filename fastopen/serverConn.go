// Interfaces for the server's establish tcp connection with a client

package main

import (
	"fmt"
	"syscall"
	"time"
)

// TFOServerConn TFO服务端套接字
type TFOServerConn struct {
	sockaddr *syscall.SockaddrInet4
	fd       int
}

// Handle 用于处理套接字
func (cxn *TFOServerConn) Handle() {
	defer cxn.Close()
	fmt.Printf("服务端接受到了套接字，远程地址: %v, 远程端口: %d\n",
		cxn.sockaddr.Addr, cxn.sockaddr.Port)
	buf := make([]byte, 1500)
	sendData := []byte("service:花花的北鼻")
	go func() {
		defer fmt.Println("服务端写套接字退出了")
		for {
			n, err := syscall.Write(cxn.fd, sendData)
			if n == 0 || err != nil {
				if n == 0 {
					fmt.Println("socket write : 远程套接字已经关闭")
				} else {
					fmt.Println("socket write : ", err)
				}

				return
			}
			time.Sleep(time.Second)
		}
	}()
	for {
		n, err := syscall.Read(cxn.fd, buf)
		if err != nil || n == 0 {
			if n == 0 {
				fmt.Println("socket read : 远程套接字已经关闭")
			} else {
				fmt.Println("socket read : ", err)
			}
			return
		}
		fmt.Printf("服务端套接字读取%d字节: %#v\n", n, string(buf[:n]))
		time.Sleep(time.Second)
	}

}

// Close 用于调用关闭
func (cxn *TFOServerConn) Close() {
	err := syscall.Shutdown(cxn.fd, syscall.SHUT_RDWR)
	if err != nil {
		fmt.Println("Shutdown：", err)
	}
	err = syscall.Close(cxn.fd)
	if err != nil {
		fmt.Println("Close", err)
	}

}
