package main

import (
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8421")
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		conn.Write([]byte("hello"))
		log.Println("send", conn.LocalAddr().String())

		time.Sleep(time.Second * 1)
		data := make([]byte, 1024)
		n, _ := conn.Read(data)
		log.Println("read", conn.LocalAddr().String(), string(data[:n]))

		time.Sleep(time.Second * 5)
	}
}
