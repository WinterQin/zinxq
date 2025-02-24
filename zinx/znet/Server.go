package znet

import (
	"fmt"
	"net"
	"zinxq/zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	// 解析地址
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
	}
	// 监听地址
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err)
	}
	// 循环处理连接
	for {
		// 阻塞等待，获取连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP err:", err)
		}
		go func() {
			defer conn.Close()
			for {
				// buffer 用于存放数据
				buff := make([]byte, 512)
				// 获取客户端发送的数据 nmsg表示数据的长度
				nmsg, err := conn.Read(buff)
				if err != nil {
					fmt.Println("Read err:", err)
					break
				}
				// 打印数据
				fmt.Println(string(buff[:nmsg]))
				// 处理数据
				newbuff := handler(buff[:nmsg])
				// 回写处理过后的数据
				_, err = conn.Write(newbuff)
				if err != nil {
					fmt.Println("Write err:", err)
				}
			}
		}()
	}
}
func handler(buff []byte) []byte {
	if string(buff[0]) == "1" && string(buff[1]) == "." {
		buff = buff[2:]
	}

	//if string(buff[0]) == "2" && string(buff[1]) == "." {
	////	stopChan <- context.
	////}
	return buff
}
func (s *Server) Stop() {

}
func (s *Server) RunServer() {
	// 创建一个可以取消的context

	go s.Start()
	//stopChan := make(chan context.Context)

	// select {} 阻塞进程
	select {}
}

func InitServer(name string, ipv string, ip string, port int) ziface.IServer {
	server := &Server{
		Name:      name,
		IPVersion: ipv,
		IP:        ip,
		Port:      port,
	}
	return server
}
