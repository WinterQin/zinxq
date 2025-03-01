package znet

import (
	"fmt"
	"github.com/winterqin/zinxq/utils"
	"github.com/winterqin/zinxq/ziface"
	"math/rand"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.Config.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.Config.WorkerPoolSize),
	}
}

func (mhd *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mhd.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("[client] api msgId = ", request.GetMsgID(), " is not FOUND!")
		request.GetConnection().Send(NotFoundMessage)
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.CurHandle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mhd *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mhd.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mhd.Apis[msgId] = router
}

func (mhd *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mhd.WorkerPoolSize); i++ {
		mhd.TaskQueue[i] = make(chan ziface.IRequest, utils.Config.MaxWorkerTaskLen)
		go mhd.StartWorkOneWorker(i, mhd.TaskQueue[i])
	}
}
func (mhd *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	randomNum := rand.Uint32()
	workerID := +randomNum % mhd.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)

	mhd.TaskQueue[workerID] <- request
}

func (mhd *MsgHandle) StartWorkOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("[zinx] Worker ID = ", workerID, " is started.")
	for {
		select {
		case request := <-taskQueue:
			mhd.DoMsgHandler(request)
		}
	}
}
