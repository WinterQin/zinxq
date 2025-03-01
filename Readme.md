# Zinx

## Zinxv0.1基础server

首先定义ziface和znet模块
ziface用于存放zinx框架的全部模块的抽象层接口类，最基本的是服务类接口IServer，定义在ziface模块中

znet模块是zinx框架中网络相关功能的实现，所有网络相关的模块都会定义在znet模块中



实现：

```bash
└── zinxq
    ├── ziface
    │   └── IServer
    └── znet
        ├── Server
```

```go
package ziface

//定义服务器接口
type IServer interface{
    //启动服务器方法
    Start()
    //停止服务器方法
    Stop()
    //开启业务服务方法
    Serve()
}
```

```go
package znet

import (
    "fmt"
    "net"
    "time"
    "zinx/ziface"
)

//iServer 接口实现，定义一个Server服务类
type Server struct {
    //服务器的名称
    Name string
    //tcp4 or other
    IPVersion string
    //服务绑定的IP地址
    IP string
    //服务绑定的端口
    Port int
}


//============== 实现 ziface.IServer 里的全部接口方法 ========

//开启网络服务
func (s *Server) Start() {
    fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

   
}

func (s *Server) Stop() {
    fmt.Println("[STOP] Zinx server , name " , s.Name)

    //TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

func (s *Server) Serve() {
    s.Start()

    //TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加


    //阻塞,否则主Go退出， listenner的go将会退出
    for {
        time.Sleep(10*time.Second)
    }
}


/*
  创建一个服务器句柄
 */
func NewServer (name string) ziface.IServer {
    s:= &Server {
        Name :name,
        IPVersion:"tcp4",
        IP:"0.0.0.0",
        Port:7777,
    }

    return s
}
```

## Zinx-V0.2-简单的连接封装与业务绑定

#### A) ziface创建iconnection.go 

```go
 
package ziface
import "net"

//定义连接接⼝口
type IConnection interface {
	//启动连接，让当前连接开始⼯工作
    Start()
	//停⽌止连接，结束当前连接状态M
	Stop()
	//从当前连接获取原始的socket TCPConn GetTCPConnection() *net.TCPConn //获取当前连接ID
	GetConnID() uint32 //获取远程客户端地址信息 RemoteAddr() net.Addr
}

//定义⼀一个统⼀一处理理链接业务的接⼝口
type HandFunc func(*net.TCPConn, []byte, int) error
```

​	该接⼝的一些基础方法，代码注释已经介绍的很清楚，这里先简单说明一个HandFunc这个函数类型， 这个是所有conn链接在处理业务的函数接口，第一参数是socket原生链接，第二个参数是客户端请求的数据，第三个参数是客户端请求的数据长度。这样，如果我们想要指定一个conn的处理业务，只要定义一个HandFunc类型的函数，然后和该链接绑定就可以了

#### B) znet 创建iconnection.go 

```go
package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//该连接的处理方法api
	handleAPI ziface.HandFunc

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}


//创建连接的方法
func NewConntion(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection{
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		handleAPI: callback_api,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}

/* 处理conn读数据的Goroutine */
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for  {
		//读取我们最大的数据到buf中

        //调用当前链接业务(这里执行的是当前conn的绑定的handle方法)
		
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <- c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}	
```

现在zinx已经初具雏形，接下来我们来继续改进它

