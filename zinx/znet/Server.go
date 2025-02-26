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
	var cid uint32 = 0
	// 循环处理连接
	for {
		// 阻塞等待，获取连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP err:", err)
		}
		delaConn := NewConnection(conn, cid, CallBackToClient)
		cid++
		go delaConn.Start()
	}
}

func CallBackToClient(conn *net.TCPConn, buff []byte, n int) error {
	// 判断buff长度是否正好为5且前五个字符是否为"hello"
	if len(buff) > 5 && string(buff[:5]) == "hello" {
		buff = []byte("hello client!")
	}
	_, err := conn.Write(buff)
	if err != nil {
		fmt.Println("Write err:", err)
		return err
	}
	return nil
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
