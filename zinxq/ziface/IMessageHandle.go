package ziface

type IMessageHandle interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
}
