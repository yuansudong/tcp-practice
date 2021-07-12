package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var port int
var run string
var addr string
var connCount int
var enableNagle bool
var enableQuick bool
var enableCrok bool

// 要测试下并发写。
func main() {
	fmt.Println("当前的进程号是:", os.Getpid())
	flag.IntVar(&connCount, "conn-count", 1, "--conn-count=5")
	flag.StringVar(&addr, "addr", "172.17.0.15", "--addr=172.17.0.15")
	flag.StringVar(&run, "run", "client", "--run=client 或者 --run=service")
	flag.IntVar(&port, "port", 12345, "--port=9999")
	flag.BoolVar(&enableNagle, "nagle", false, "--nagle")
	flag.BoolVar(&enableQuick, "quick-ack", false, "--quick-ack")
	flag.BoolVar(&enableCrok, "crok", false, "--crok")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if run == "service" {
		RunService()
	} else {
		RunClient()
	}
}

func RunClient() {

	conn, err := net.DialTCP("tcp4", &net.TCPAddr{}, &net.TCPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// if err := conn.SetNoDelay(!enableNagle); err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }
	// if err := conn.SetQuickAck(enableQuick); err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }

	// if err := conn.SetCrok(enableCrok); err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }

	go func() {
		data := []byte{}
		for i := 0; i < 1; i++ {
			data = append(data, []byte("abcdef")...)
		}
		for i := 0; i < 100000; i++ {
			if _, err := conn.Write(data); err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	go func() {
		data := []byte{}
		for i := 0; i < 1; i++ {
			data = append(data, []byte("ABCDEF")...)
		}
		for i := 0; i < 100000; i++ {
			if _, err := conn.Write(data); err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	go func() {
		data := []byte{}
		for i := 0; i < 1; i++ {
			data = append(data, []byte("~!@#$%")...)
		}
		for i := 0; i < 100000; i++ {
			if _, err := conn.Write(data); err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	fmt.Println("数据发送完毕")
	sChan := make(chan os.Signal, 10)
	signal.Notify(sChan, syscall.SIGHUP, syscall.SIGUSR2)
	for {
		select {
		case sig := <-sChan:
			if sig == syscall.SIGUSR2 {
				fmt.Println("调用关闭函数")
				conn.Close()
				return
			}
		}
	}

}

func RunService() {
	listener, err := net.Listen("tcp4", addr+":"+fmt.Sprint(port))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		go ReadHandle(conn)
	}
}

// ReadHandle 用于读取
func ReadHandle(conn net.Conn) {
	buf := make([]byte, 40960)
	for {
		time.Sleep(30 * time.Second)
		n, err := conn.Read(buf)
		if n == 0 {
			err = io.EOF
		}
		if err != nil {
			conn.Close()
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(buf[:n]))
	}

}
