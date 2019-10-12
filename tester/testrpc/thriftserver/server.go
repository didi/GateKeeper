package thriftserver

import (
	"context"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/didi/gatekeeper/tester/testrpc/thriftgen"
	"log"
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
		log.Fatal("need addr like :8007")
	}
	handler := &FormatDataImpl{
		Addr: fmt.Sprintf("127.0.0.1%s", addr), //打印格式
	}
	processor := thriftgen.NewFormatDataProcessor(handler)
	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if len(showLog) == 0 || showLog[0] == true {
		log.Println("RunTcpDestServer ", addr)
	}
	go func() {
		server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
		t.lock.Lock()
		t.mapTCPDestSvr[addr] = server
		t.lock.Unlock()
		err := server.Serve()
		if err != nil {
			//log.Printf("\nRunTcpDestServer addr:%v -err:%v ", addr, err)
		}
	}()
}

//Stop Stop
func (t *TestTCPDestServer) Stop(addr string) {
	//log.Println("t.mapTCPDestSvr",t.mapTCPDestSvr)
	if server, ok := t.mapTCPDestSvr[addr]; ok {
		err := server.Stop()
		if err != nil {
			//fmt.Printf("TestTCPDestServer.Stop -addr:%v -err:%v\n",addr,err)
		}
		//fmt.Println("TestTCPDestServer.Stop",addr)
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
