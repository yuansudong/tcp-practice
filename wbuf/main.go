package main

import (
	"flag"
	"fmt"
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

//
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

	conn, err := net.DialTCP("tcp4", &net.TCPAddr{}, &net.TCPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	conn.SetKeepAlive(false)

	fmt.Println("连接建立成功。")
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("read:", err.Error())
				os.Exit(-1)
			}
			fmt.Println(string(buf[:n]))
			time.Sleep(30 * time.Second)
		}
	}()

	// fmt.Println("连接建立成功")
	sChan := make(chan os.Signal, 10)
	signal.Notify(sChan, syscall.SIGHUP, syscall.SIGUSR2)
	for {
		select {
		case sig := <-sChan:
			if sig == syscall.SIGUSR2 {
				fmt.Println("调用发送函数")
				go func() {
					//
					file, err := os.OpenFile("/home/yuansudong/Android/Sdk.zip", os.O_RDONLY, 0644)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(-1)
					}
					defer file.Close()
					buf := make([]byte, 64)
					tn := 0
					for {
						rn, rerr := file.Read(buf)
						if rerr != nil {
							fmt.Println(err.Error())
							os.Exit(-1)
						}
						wn, werr := conn.Write(buf[:rn])
						if werr != nil {
							fmt.Println(err.Error())
							os.Exit(-1)
						}
						tn += wn
						fmt.Println("总共写:", tn)
					}
				}()
			}
		}
	}

}

func RunService() {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.ParseIP("172.17.0.15"),
		Port: port,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("端口", port)
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Println(err.Error())
		}
		conn.SetKeepAlive(false)
		fmt.Println("接收了套接字。")
		go func() {
			defer conn.Close()
			buf := make([]byte, 4096)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("read:", err.Error())
				}
				fmt.Println("读取了套接字：", n)
				time.Sleep(30 * time.Second)
			}
		}()
	}
}
