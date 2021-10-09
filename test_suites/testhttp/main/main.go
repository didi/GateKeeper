package main

import (
	"flag"
	"fmt"
	"github.com/didi/gatekeeper/test_suites/testhttp"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ports := flag.String("ports", "8018,8017", "set server port")
	flag.Parse()
	server := testhttp.NewTestHTTPDestServer()
	for _, port := range strings.Split(*ports, ",") {
		server.Run(fmt.Sprintf(":%v", port))
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
