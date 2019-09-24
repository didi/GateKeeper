package main

import (
	"github.com/didichuxing/gatekeeper/tester/testhttp"
	"time"
)

func main()  {
	server:= testhttp.NewTestHTTPDestServer()
	server.Run(":8018")
	time.Sleep(time.Second*1000)
}
