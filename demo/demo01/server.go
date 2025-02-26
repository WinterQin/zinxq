package main

import "zinxq/zinx/znet"

func main() {
	server := znet.InitServer("winter_server", "tcp4", "0.0.0.0", 8999)
	server.RunServer()
}
