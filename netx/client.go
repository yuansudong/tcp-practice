package netx

import (
	"fmt"
	"net"
	"syscall"
)

func DialTCP(ip string, port int) (*TCPConn, error) {
	IP := net.ParseIP(ip)
	if IP == nil {
		return nil, fmt.Errorf("no parse ip:%s", ip)
	}
	sa := &syscall.SockaddrInet4{
		Port: port,
	}
	copy(sa.Addr[:], IP[12:16])
	socketfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	if err := syscall.Connect(socketfd, sa); err != nil {
		return nil, err
	}

	return &TCPConn{
		fd: socketfd,
	}, nil
}
