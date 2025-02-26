package ziface

import "net"

// IConnection 定义链接模块的抽象层
type IConnection interface {
	// Start 启动链接
	Start()
	// Stop 停止链接
	Stop()
	// GetTCPConnection 获取当前链接绑定的socket的conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前链接的ID
	GetConnID() uint32
	// GetRemoteAddr 获取远程客户端的IP和端口
	GetRemoteAddr() net.Addr
	// Send 发送数据到客户端
	Send(data []byte)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
