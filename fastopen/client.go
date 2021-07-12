package main

import (
	"syscall"
)

type TFOClient struct {
	sa         *syscall.SockaddrInet4
	ServerAddr [4]byte
	ServerPort int
	fd         int
}

func New(addr [4]byte, port int) *TFOClient {
	inst := new(TFOClient)
	inst.ServerAddr = addr
	inst.ServerPort = port
	inst.sa = &syscall.SockaddrInet4{Addr: inst.ServerAddr, Port: inst.ServerPort}
	return inst
}

func (c *TFOClient) Connect() (err error) {
	c.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}
	err = syscall.Sendto(c.fd, []byte("connect"), syscall.MSG_FASTOPEN, c.sa)
	return err
}

func (c *TFOClient) Close() (err error) {
	return syscall.Close(c.fd)
}

// Send发送数据
func (c *TFOClient) Send(buf []byte) (int, error) {
	return syscall.Write(c.fd, buf)
}

// Read读取数据
func (c *TFOClient) Read(buf []byte) (int, error) {
	return syscall.Read(c.fd, buf)
}
