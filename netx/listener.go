package netx

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

type Listener struct {
	fd      int
	sa      syscall.Sockaddr
	backlog int
}

// Listen 监听套接字。
func Listen(ip string, port int) (*Listener, error) {
	var err error
	inst := &Listener{}
	sa := &syscall.SockaddrInet4{
		Port: port,
	}
	IP := net.ParseIP(ip)
	if IP == nil {
		return nil, fmt.Errorf("ip地址有误：", ip)
	}
	copy(sa.Addr[:], IP[12:16])
	inst.sa = sa
	inst.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	log.Println("当前的文件描述符是", inst.fd)
	if err := syscall.Bind(inst.fd, inst.sa); err != nil {
		return nil, err
	}
	if err := syscall.Listen(inst.fd, 512); err != nil {
		return nil, err
	}

	return inst, nil
}

func (l *Listener) Accept() (*TCPConn, error) {
	fd, _, err := syscall.Accept(l.fd)
	if err != nil {
		return nil, err
	}
	return &TCPConn{
		fd: fd,
	}, nil
}

func (l *Listener) Close() error {
	return syscall.Close(l.fd)
}

//
// func (l *Listener) SetBacklog(backlog int) {
// 	l.backlog = backlog
// }

// SetFastopen 开启fastopen
func (l *Listener) SetFastopen() error {
	return syscall.SetsockoptInt(l.fd, syscall.IPPROTO_TCP, unix.TCP_FASTOPEN, 1)
}
