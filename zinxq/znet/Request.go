package znet

import (
	"github.com/winterqin/zinxq/ziface"
	"sync"
)

type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
	ctx  sync.Map // 线程安全的存储
}

func NewRequest(conn ziface.IConnection, msg ziface.IMessage) *Request {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}

// GetConnection 获取当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func (r *Request) Set(key string, value interface{}) {
	//fmt.Printf("Set: key=%s, value=%v\n", key, value)
	r.ctx.Store(key, value)
}

func (r *Request) Get(key string) (interface{}, bool) {
	value, ok := r.ctx.Load(key)
	//fmt.Printf("Get: key=%s, value=%v, ok=%v\n", key, value, ok)
	return value, ok
}
