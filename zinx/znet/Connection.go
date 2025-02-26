package znet

import (
	"fmt"
	"io"
	"net"
	"zinxq/zinx/ziface"
)

type Connection struct {
	//当前链接的socket
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//链接处理方法
	HandlerApi ziface.HandleFunc

	//告知结束的channel
	ExitChan chan bool

	//链接的状态
	isClosed bool
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connid uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connid,
		HandlerApi: callbackApi,
		ExitChan:   make(chan bool, 1),
		isClosed:   false,
	}
	return c
}

// private
func (c *Connection) startReader() {
	fmt.Println("Reader is running......")
	defer fmt.Println("Reader is stopped, ConnID is ", c.ConnID, "RemoteAddr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端数据到buff中，目前最大为512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		fmt.Println("recv from ConnID", c.ConnID, ": ", string(buf[:cnt]))
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("recv buf err:", err)
			continue

		}

		//调用当前链接所绑定的HandleApi
		if err := c.HandlerApi(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID is :", c.ConnID, "handler err:", err)
			break
		}
	}
}

// public
// Start 启动链接
func (c *Connection) Start() {
	fmt.Println("connection start......    ConnID:", c.ConnID)
	//启动读取goroutine
	go c.startReader()
	// TODO 启动写入goroutine业务
}

// Stop 停止链接
func (c *Connection) Stop() {
	fmt.Println("connection stop......    ConnID:", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	select {
	case c.ExitChan <- true:
		c.Conn.Close()
	}
	close(c.ExitChan)
}

// GetTCPConnection 获取当前链接绑定的socket的conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前链接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetRemoteAddr 获取远程客户端的IP和端口
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据到客户端
func (c *Connection) Send(data []byte) {
	_, err := c.Conn.Write(data)
	if err != nil {
		fmt.Println("Send error:", err)
	}
}
