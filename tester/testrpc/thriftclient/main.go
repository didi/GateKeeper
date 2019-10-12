package main

import (
	"context"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/didi/gatekeeper/tester/testrpc/thriftgen"
	"log"
	"os"
	"time"
)

func main() {
	addr := flag.String("addr", "", "input addr like 127.0.0.1:8800")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	for {
		tSocket, err := thrift.NewTSocket(*addr)
		if err != nil {
			log.Fatalln("tSocket error:", err)
		}

		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, _ := transportFactory.GetTransport(tSocket)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thriftgen.NewFormatDataClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			log.Fatalln("Error opening:", *addr)
		}
		defer transport.Close()
		data := thriftgen.Data{Text: "ping"}
		d, err := client.DoFormat(context.Background(), &data)
		if err != nil {
			fmt.Println("err:", err.Error())
		} else {
			fmt.Println("Text:", d.Text)
		}
		time.Sleep(40 * time.Millisecond)
	}
}
