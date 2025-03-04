package main

import (
	"zinxq/zinx/ziface"
	"zinxq/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) CurHandle(request ziface.IRequest) {
	conn := request.GetConnection()
	data := request.GetData()
	result := "ping...ping...ping" + string(data)
	conn.Send([]byte(result))
}

func main() {
	server := znet.InitServer("winter_server", "tcp4", "0.0.0.0", 8999)
	server.AddRouter(&PingRouter{})
	server.RunServer()
}
