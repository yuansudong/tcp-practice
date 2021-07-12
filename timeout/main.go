package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	now := time.Now()
	conn, err := net.Dial("tcp4", "www.google.com:443")
	if err != nil {
		fmt.Println(time.Since(now))
		fmt.Println(err.Error())
		goto end
	}
	conn.Close()
end:
	fmt.Println(time.Since(now))
}
