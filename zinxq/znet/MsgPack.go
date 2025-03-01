package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/winterqin/zinxq/utils"
	"github.com/winterqin/zinxq/ziface"
)

type MsgPack struct {
}

func NewMsgPack() *MsgPack {
	return &MsgPack{}
}

func (mp *MsgPack) GetHeadLen() uint32 {
	return 8
}

func (mp *MsgPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil

}

func (mp *MsgPack) Unpack(binarydata []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewBuffer(binarydata)

	msg := &Message{}
	var msglen uint32
	var msgid uint32
	if err := binary.Read(dataBuff, binary.LittleEndian, &msglen); err != nil {
		return nil, err
	}
	msg.SetMsgLen(msglen)

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msgid); err != nil {
		return nil, err
	}
	msg.SetMsgID(msgid)

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.Config.MaxPacketSize > 0 && msg.GetMsgLen() > utils.Config.MaxPacketSize {
		return nil, errors.New("[zinx MsgPack] too large msg data recieved")
	}

	return msg, nil
}
