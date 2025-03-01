package znet

type Message struct {
	msgLen uint32
	data   []byte
	msgID  uint32
}

var NotFoundMessage = NewMessage(404, []byte("Request ID Not Found"))
var SuccessMessage = NewMessage(200, []byte("Success"))

func NewMessage(msgID uint32, data []byte) *Message {
	msg := &Message{data: data, msgID: msgID}
	msg.msgLen = uint32(len(data))
	return msg
}

func (m *Message) GetMsgID() uint32 {
	return m.msgID
}

func (m *Message) GetMsgLen() uint32 {
	return m.msgLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetMsgID(id uint32) {
	m.msgID = id
}

func (m *Message) SetMsgLen(msglen uint32) {
	m.msgLen = msglen
}

func (m *Message) SetData(data []byte) {
	m.data = data
}
