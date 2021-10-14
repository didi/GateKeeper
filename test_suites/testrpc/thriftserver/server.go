package thriftserver

import (
	"context"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/test_suites/testrpc/thriftgen"
	"sync"
)

//TestTCPDestServer TestTCPDestServer
type TestTCPDestServer struct {
	mapTCPDestSvr map[string]*thrift.TSimpleServer
	lock          *sync.Mutex
}

//NewTestTCPDestServer NewTestTCPDestServer
func NewTestTCPDestServer() *TestTCPDestServer {
	return &TestTCPDestServer{
		mapTCPDestSvr: map[string]*thrift.TSimpleServer{},
		lock:          &sync.Mutex{},
	}
}

//Run Run
func (t *TestTCPDestServer) Run(addr string, showLog ...bool) {
	flag.Parse()
	if addr == "" {
		log.Info().Msg(lib.Purple("need addr like :8007"))
		log.Fatal()
	}
	handler := &FormatDataImpl{
		Addr: fmt.Sprintf("127.0.0.1%s", addr), //打印格式
	}
	processor := thriftgen.NewFormatDataProcessor(handler)
	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Info().Msg(lib.Purple("Error:" + err.Error()))
		log.Fatal()
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if len(showLog) == 0 || showLog[0] == true {
		log.Info().Msg(lib.Purple("start test rpcserver:" + addr))
	}
	go func() {
		server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
		t.lock.Lock()
		t.mapTCPDestSvr[addr] = server
		t.lock.Unlock()
		err := server.Serve()
		if err != nil {
		}
	}()
}

//Stop Stop
func (t *TestTCPDestServer) Stop(addr string) {
	if server, ok := t.mapTCPDestSvr[addr]; ok {
		err := server.Stop()
		if err != nil {
		}
	}
}

//FormatDataImpl FormatDataImpl
type FormatDataImpl struct {
	Addr string
}

//DoFormat DoFormat
func (fdi *FormatDataImpl) DoFormat(ctx context.Context, data *thriftgen.Data) (r *thriftgen.Data, err error) {
	var rData thriftgen.Data
	rData.Text = fdi.Addr
	return &rData, nil
}
