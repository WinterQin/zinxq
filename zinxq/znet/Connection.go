package znet

import (
	"errors"
	"fmt"
	"github.com/winterqin/zinxq/ziface"
	"io"
	"net"
)

type Connection struct {
	//当前链接的socket
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//告知结束的channel
	ExitChan chan bool

	//链接的状态
	isClosed bool

	// msg handler
	Msghd ziface.IMessageHandle
	// msg channel
	msgChan chan []byte
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connid uint32, msghd ziface.IMessageHandle) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connid,
		ExitChan: make(chan bool, 1),
		isClosed: false,
		Msghd:    msghd,
		msgChan:  make(chan []byte),
	}
	return c
}

// private
func (c *Connection) startReader() {
	fmt.Println("[zinx] Reader is running......")
	defer fmt.Println("[zinx] Reader is stopped, ConnID is ", c.ConnID, "RemoteAddr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()
	mp := NewMsgPack()

	for {

		HeadData := make([]byte, mp.GetHeadLen())

		_, err := io.ReadFull(c.GetTCPConnection(), HeadData)
		if err != nil {
			break
		}

		msg, err := mp.Unpack(HeadData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("[zinx connection] read msg data error ", err)
				continue
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := NewRequest(c, msg)
		go c.Msghd.DoMsgHandler(req)

	}
}

func (c *Connection) startWriter() {
	fmt.Println("[zinx] Writer is running......")
	defer fmt.Println("[zinx] Writer is stopped, ConnID is ", c.ConnID, "RemoteAddr is ", c.Conn.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("[zinx connection startWriter] Send error:", err)
			}
		case <-c.ExitChan:
			return
		}
	}

}
func (c *Connection) Send(msg ziface.IMessage) error {
	if c.isClosed {
		return errors.New("[zinx connection send] connection closed")
	}
	mp := NewMsgPack()
	data, err := mp.Pack(msg)
	_, err = c.Conn.Write(data)
	if err != nil {
		fmt.Println("pack error message id:", msg.GetMsgID())
		return errors.New("pack error message")
	}

	c.msgChan <- data
	return nil
}

// public
// Start 启动链接
func (c *Connection) Start() {
	fmt.Println("[zinx start]connection start......    ConnID:", c.ConnID)
	//启动读取goroutine
	go c.startReader()
	go c.startWriter()
	for {
		select {
		case <-c.ExitChan:
			return
		}
	}
}

// Stop 停止链接
func (c *Connection) Stop() {
	fmt.Println("[zinx] connection stop......    ConnID:", c.ConnID)
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
