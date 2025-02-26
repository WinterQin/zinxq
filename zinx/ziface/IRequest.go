package ziface

type IRequest interface {
	// GetConnection 获取当前连接
	GetConnection() IConnection
	// GetData 获取数据
	GetData() []byte
}
