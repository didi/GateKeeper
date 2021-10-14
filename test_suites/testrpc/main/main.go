package main

import (
	"flag"
	"fmt"
	"github.com/didi/gatekeeper/test_suites/testrpc/thriftserver"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ports := flag.String("ports", "8218,8217", "set server port")
	flag.Parse()
	server := thriftserver.NewTestTCPDestServer()
	for _, port := range strings.Split(*ports, ",") {
		server.Run(fmt.Sprintf(":%v", port))
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
