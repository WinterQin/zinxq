package ziface

type IMsgPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(binary []byte) (IMessage, error)
}
