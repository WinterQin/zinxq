package znet

import (
	"errors"
	"fmt"
	"github.com/winterqin/zinxq/utils"
	"github.com/winterqin/zinxq/ziface"
	"io"
	"net"
)

type Connection struct {
	// 当前conn属于哪个server
	TcpServer ziface.IServer
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
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte //定义channel成员
}

// NewConnection 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connid uint32, msghd ziface.IMessageHandle) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connid,
		ExitChan:    make(chan bool, 1),
		isClosed:    false,
		Msghd:       msghd,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.Config.MaxMsgChanLen),
	}
	c.TcpServer.GetConnMgr().AddConnection(c)
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
				return
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := NewRequest(c, msg)

		if utils.Config.WorkerPoolSize > 0 {
			c.Msghd.SendMsgToTaskQueue(req)
		} else {
			go c.Msghd.DoMsgHandler(req)
		}

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
			//针对有缓冲channel需要些的数据处理
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break

			}
		case <-c.ExitChan:
			return
		}
	}

}

// SendBuffMsg 直接将Message数据发送给远程的TCP客户端(有缓冲)
// 添加带缓冲发送消息接口
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("[zinx] Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewMsgPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}
func (c *Connection) Send(msg ziface.IMessage) error {
	if c.isClosed {
		return errors.New("[zinx] connection closed")
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

	//==================
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
	//==================

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
	//==================
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)
	//==================
	
	c.Conn.Close()
	c.ExitChan <- true
	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().RemoveConnection(c) //删除conn从ConnManager中
	//关闭该链接全部管道
	defer close(c.ExitChan)
	defer close(c.msgChan)
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
