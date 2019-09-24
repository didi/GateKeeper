package common

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/didichuxing/gatekeeper/service"
	"github.com/didichuxing/gatekeeper/tester/testrpc/thriftgen"
	"github.com/e421083458/golang_common/lib"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
)

//TestFlowControl Http限流测试
func TestFlowControl(t *testing.T) {
	//公共配置
	config1 := service.SysConfMgr.GetModuleConfigByName("tester_whitelist")
	url1 := fmt.Sprintf("http://%s%s%s/get_host",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		config1.MatchRule[0].Rule)
	convey.Convey("Http限流测试", t, func() {
		convey.Convey("客户端ip限流", func() {
			wg := sync.WaitGroup{}
			overflowFlag := false
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if isOverFlow(100, url1) {
						overflowFlag = true
					}
				}()
			}
			wg.Wait()
			convey.So(overflowFlag, convey.ShouldBeTrue)
		})
	})
}

//TestTCPFlowControl Tcp限流测试
func TestTCPFlowControl(t *testing.T) {
	//公共配置
	tcpConfig := service.SysConfMgr.GetModuleConfigByName("tester_tcp_whitelist")
	checkTCPAddr := fmt.Sprintf("%s%s", lib.GetStringConf("base.cluster.cluster_ip"), tcpConfig.Base.FrontendAddr)
	convey.Convey("Tcp限流测试", t, func() {
		convey.Convey("客户端ip限流", func() {
			wg := sync.WaitGroup{}
			overflowFlag := false
			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if isTCPEmptyResponse(5, checkTCPAddr) {
						overflowFlag = true
					}
				}()
			}
			wg.Wait()
			convey.So(overflowFlag, convey.ShouldBeTrue)
		})
	})
}

func isTCPEmptyResponse(times int, addr string) bool {
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println(err)
		}
	}()
	for i := 0; i < times; i++ {
		tSocket, err := thrift.NewTSocket(addr)
		if err != nil {
			log.Println("tSocket error:", err)
			continue
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, _ := transportFactory.GetTransport(tSocket)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thriftgen.NewFormatDataClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			log.Println("Error opening:", addr)
			continue
		}
		defer transport.Close()
		data := thriftgen.Data{Text: "ping"}
		if _, err := client.DoFormat(context.Background(), &data); err != nil {
			return true
		}
	}
	return false
}

func isOverFlow(times int, checkURL string) bool {
	allHostSlice := []string{}
	for i := 0; i < times; i++ {
		client := http.Client{}
		resp, err := client.Get(checkURL)
		if err != nil {
			continue
		}
		queryBody, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		c := string(queryBody)
		//log.Println(c)
		if strings.Contains(c, "QPS") {
			return true
		}
		allHostSlice = append(allHostSlice, c)
	}
	return false
}
