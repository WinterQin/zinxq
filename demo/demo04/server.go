package main

import (
	"github.com/winterqin/zinxq/ziface"
	"github.com/winterqin/zinxq/znet"
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
	server := znet.InitServer()
	server.AddRouter(&PingRouter{})
	server.RunServer()
}
