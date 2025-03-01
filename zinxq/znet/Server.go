package znet

import (
	"fmt"
	"github.com/winterqin/zinxq/utils"
	"github.com/winterqin/zinxq/ziface"
	"net"
	"reflect"
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
	s := &Server{
		Name:      utils.Config.Name,
		IPVersion: utils.Config.IPVersion,
		IP:        utils.Config.Host,
		Port:      utils.Config.TcpPort,
		Msghd:     NewMsgHandle(),
	}

	fmt.Printf("[START] Server name: %s, listen at Host: %s:%d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinxq] IPVersion: %s, MaxConnectionNum: %d,  MaxPacketSize: %d\n",
		utils.Config.IPVersion,
		utils.Config.MaxConnectionNum,
		utils.Config.MaxPacketSize)

	return s
}

func (s *Server) Start() {

	// 解析地址
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("[zinx] ResolveTCPAddr err:", err)
	}
	// 监听地址
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println("[zinx server error] ListenTCP err:", err)
	}
	var cid uint32 = 0
	// 循环处理连接
	for {
		// 阻塞等待，获取连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("[zinx server error] AcceptTCP err:", err)
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

	// 通过反射获取 router 的类型信息
	routerType := reflect.TypeOf(router)

	// 如果 router 是指针类型，需要获取其指向的类型
	if routerType.Kind() == reflect.Ptr {
		routerType = routerType.Elem()
	}

	// 获取 router 的名称
	routerName := routerType.Name()

	fmt.Println("[zinx router] AddRouter:", "ID: ", msgId, "=====> ", routerName)
}
