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
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager

	// =======================
	//新增两个hook函数原型

	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)

	// =======================

}

func InitServer() ziface.IServer {
	utils.Config.Reload()
	s := &Server{
		Name:      utils.Config.Name,
		IPVersion: utils.Config.IPVersion,
		IP:        utils.Config.Host,
		Port:      utils.Config.TcpPort,
		Msghd:     NewMsgHandle(),
		ConnMgr:   NewConnManager(), //创建ConnManager
	}

	fmt.Printf("[START] Server name: %s, listen at Host: %s:%d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinxq] IPVersion: %s, MaxConnectionNum: %d,  MaxPacketSize: %d\n",
		utils.Config.IPVersion,
		utils.Config.MaxConnectionNum,
		utils.Config.MaxPacketSize)

	return s
}

// 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) Start() {
	//0 启动worker工作池机制
	s.Msghd.StartWorkerPool()

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
			continue
		}

		if s.ConnMgr.CurConnNum() >= utils.Config.MaxConnectionNum {
			conn.Close()
			continue
		}
		delaConn := NewConnection(s, conn, cid, s.Msghd)
		cid++
		go delaConn.Start()
	}
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
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

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
