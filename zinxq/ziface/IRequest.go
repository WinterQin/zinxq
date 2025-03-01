package ziface

type IRequest interface {
	// GetConnection 获取当前连接
	GetConnection() IConnection
	// GetData 获取数据

	GetData() []byte
	GetMsgID() uint32

	//存储上下文
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}
