package znet

import (
	"zinxq/zinx/ziface"
)

type Request struct {
	conn ziface.IConnection
	data []byte
}

func NewRequest(conn ziface.IConnection, data []byte) *Request {
	return &Request{
		conn: conn,
		data: data,
	}
}

// GetConnection 获取当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取数据
func (r *Request) GetData() []byte {
	return r.data
}
