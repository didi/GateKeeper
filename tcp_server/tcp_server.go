package tcp_server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed     = errors.New("tcp: Server closed")
	ErrAbortHandler     = errors.New("tcp: abort TCPHandler")
	ServerContextKey    = &contextKey{"tcp-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type TcpServer struct {
	Addr        string
	Handler     TCPHandler
	err         error
	BaseCtx     context.Context
	UpdateAt    time.Time
	ServiceName string

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.Mutex
	inShutdown int32
	doneChan   chan struct{}
	l          *onceCloseListener
}

func (s *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

func (srv *TcpServer) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	addr := srv.Addr
	if addr == "" {
		return errors.New("need addr")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{
		ln.(*net.TCPListener)})
}

func (srv *TcpServer) Close() error {
	atomic.StoreInt32(&srv.inShutdown, 1)
	close(srv.doneChan) //关闭channel
	srv.l.Close()       //执行listener关闭
	return nil
}

func (srv *TcpServer) Serve(l net.Listener) error {
	srv.l = &onceCloseListener{Listener: l}
	defer srv.l.Close() //执行listener关闭
	if srv.BaseCtx == nil {
		srv.BaseCtx = context.Background()
	}
	baseCtx := srv.BaseCtx
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			fmt.Printf("accept fail, err: %v\n", e)
			continue
		}
		c := srv.newConn(rw)
		go c.serve(ctx)
	}
	return nil
}

func (srv *TcpServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}
	// 设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return c
}

func (s *TcpServer) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func ListenAndServe(addr string, handler TCPHandler) error {
	server := &TcpServer{Addr: addr, Handler: handler, doneChan: make(chan struct{}),}
	return server.ListenAndServe()
}
