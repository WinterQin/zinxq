package znet

import (
	"fmt"
	"github.com/winterqin/zinxq/utils"
	"github.com/winterqin/zinxq/ziface"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// msg handler
	Msghd ziface.IMessageHandle
}

func InitServer() ziface.IServer {
	utils.Config.Reload()
	server := &Server{
		Name:      utils.Config.Name,
		IPVersion: utils.Config.IPVersion,
		IP:        utils.Config.Host,
		Port:      utils.Config.TcpPort,
		Msghd:     NewMsgHandle(),
	}
	return server
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at Host: %s:%d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinxq] IPVersion: %s, MaxConnectionNum: %d,  MaxPacketSize: %d\n",
		utils.Config.IPVersion,
		utils.Config.MaxConnectionNum,
		utils.Config.MaxPacketSize)
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
		delaConn := NewConnection(conn, cid, s.Msghd)
		cid++
		go delaConn.Start()
	}
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

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.Msghd.AddRouter(msgId, router)
	fmt.Println("AddRouter")
}
