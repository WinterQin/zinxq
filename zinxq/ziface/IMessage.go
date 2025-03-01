package ziface

type IMessage interface {
	GetMsgID() uint32
	GetMsgLen() uint32
	GetData() []byte
	SetMsgID(uint32)
	SetMsgLen(uint32)
	SetData([]byte)
}
