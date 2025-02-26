package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 尝试连接到服务器
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}
	defer conn.Close()

	for i := 0; i < 5; i++ {
		_, err = conn.Write([]byte("hello, server!"))
		if err != nil {
			fmt.Println("Write err:", err)
			return
		}

		// 读取服务器响应（如果有的话）
		buffer := make([]byte, 512)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Read err:", err)
			return
		}
		fmt.Println("Received from server:", string(buffer[:n]))
		time.Sleep(1 * time.Second)
	}
}
