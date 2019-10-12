package tcpproxy

import (
	"net"
)

//TCPHandlerFunc 构建一个函数回调结构体
type TCPHandlerFunc func(c net.Conn)

//HandleConn 回调方法
func (f TCPHandlerFunc) HandleConn(c net.Conn) {
	f(c)
}

//HandlerFunc 构建一个回调target结构体
type HandlerFunc func(next Target) Target

//TCPRouter tcp路由器
type TCPRouter struct {
	*DialProxy
	prev *TCPRouter
	mw   HandlerFunc
}

//NewTCPRouter 创建一个tcp路由器
func NewTCPRouter(dialProxy *DialProxy) *TCPRouter {
	return &TCPRouter{DialProxy: dialProxy}
}

//Use 增加一个中间件
func (r *TCPRouter) Use(mws ...HandlerFunc) *TCPRouter {
	if len(mws) == 0 {
		return r
	}
	router := r
	for _, mw := range mws {
		router = router.use(mw)
	}
	return router
}

func (r *TCPRouter) use(mw HandlerFunc) *TCPRouter {
	return &TCPRouter{
		prev:      r,
		DialProxy: r.DialProxy,
		mw:        mw,
	}
}

func (r *TCPRouter) genChainHandler(handle Target) Target {
	wraphandler := &WrapHandlerEntity{
		Handler: handle,
	}
	chain := handle
	router := r
	for router.prev != nil {
		if router.mw != nil {
			chain = router.mw(wraphandler)
		}
		wraphandler = &WrapHandlerEntity{
			Handler: chain,
		}
		router = router.prev
	}
	chain = &WrapHandlerEntity{
		Handler: chain,
	}
	return chain
}

//HandleConn 请求回调
func (r *TCPRouter) HandleConn(c net.Conn) {
	chainHandler := r.genChainHandler(r.DialProxy)
	chainHandler.HandleConn(c)
}

//WrapHandlerEntity 包装实体
type WrapHandlerEntity struct {
	Handler Target
}

//HandleConn 请求回调
func (w *WrapHandlerEntity) HandleConn(c net.Conn) {
	w.Handler.HandleConn(c)
}
