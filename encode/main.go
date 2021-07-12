package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func main() {

	by := sha1.Sum([]byte("dGhlIHNhbXBsZSBub25jZQ==258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	fmt.Println(base64.StdEncoding.EncodeToString(by[:]))
}
