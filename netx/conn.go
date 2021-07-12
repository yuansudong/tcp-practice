package netx

import "syscall"

type TCPConn struct {
	fd int
}

// SetNoDelay 用于设置是否启用Nagle算法
func (c *TCPConn) SetNoDelay(state bool) error {
	if state {
		return syscall.SetsockoptInt(c.fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	}
	return syscall.SetsockoptInt(c.fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 0)

}

// SetCrok 启用Crok算法
func (c *TCPConn) SetCrok(state bool) error {
	if state {
		return syscall.SetsockoptInt(c.fd, syscall.IPPROTO_TCP, syscall.TCP_CORK, 1)
	}
	return syscall.SetsockoptInt(c.fd, syscall.IPPROTO_TCP, syscall.TCP_CORK, 0)
}

func (c *TCPConn) SetQuickAck(state bool) error {
	if state {
		return syscall.SetsockoptInt(c.fd, syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 1)
	}
	return syscall.SetsockoptInt(c.fd, syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 0)
}

// Write 写入数据
func (c *TCPConn) Write(data []byte) (int, error) {
	return syscall.Write(c.fd, data)
}

func (c *TCPConn) Read(data []byte) (int, error) {

	return syscall.Read(c.fd, data)
}

//
func (c *TCPConn) Close() error {
	return syscall.Close(c.fd)
}
