package znet

import (
	"errors"
	"fmt"
	"github.com/winterqin/zinxq/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

/*
创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) AddConnection(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.CurConnNum())
}

func (cm *ConnManager) RemoveConnection(conn ziface.IConnection) {
	//删除连接信息
	delete(cm.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.CurConnNum())
}

func (cm *ConnManager) CurConnNum() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", cm.CurConnNum())
}

func (cm *ConnManager) GetConnById(connId uint32) (ziface.IConnection, error) {
	//保护共享资源Map 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}
