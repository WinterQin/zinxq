package ziface

type IConnManager interface {
	AddConnection(conn IConnection)                 //添加链接
	RemoveConnection(conn IConnection)              //删除连接
	CurConnNum() int                                //获取当前连接数量
	GetConnById(connId uint32) (IConnection, error) //利用ConnID获取链接
	ClearConn()                                     //删除并停止所有链接
}
