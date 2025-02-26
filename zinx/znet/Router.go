package znet

import (
	"zinxq/zinx/ziface"
)

// BaseRouter 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

// 这里之所以所有的方法都为空，是因为不是所有Router都需要实现这三个方法
// 如果有Router需要什么方法到时候再重写对应的方法就行了

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) CurHandle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
