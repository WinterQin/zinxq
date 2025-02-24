package znet

import (
	"context"
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
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP err:", err)
		}
		go func() {
			defer conn.Close()
			for {
				buff := make([]byte, 512)

				nmsg, err := conn.Read(buff)
				if err != nil {
					fmt.Println("Read err:", err)
					continue
				}
				fmt.Println(string(buff[:nmsg]))
				newbuff := handler(buff[:nmsg])
				_, err = conn.Write(newbuff)
				if err != nil {
					fmt.Println("Write err:", err)
					continue
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
	stopChan := make(chan context.Context)

	select {
	case <-stopChan:
		s.Stop()
	}
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
