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
	"regexp"
	"strings"
	"testing"
	"time"
)

//TestCheckIPList HTTP服务发现测试
func TestCheckIPList(t *testing.T) {
	moduleConfig := service.SysConfMgr.GetModuleConfigByName("tester")
	checkURL := fmt.Sprintf("http://%s%s%s/get_host",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		moduleConfig.MatchRule[0].Rule)

	convey.Convey("HTTP服务发现测试", t, func() {

		convey.Convey("配置IP与请求IP一致", func() {
			time.Sleep(10 * time.Second)
			checkHostSlice, err := getUniqueHost(10, checkURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(checkHostSlice, convey.ShouldResemble, strings.Split(moduleConfig.LoadBalance.IPList, ","))
		})

		//convey.Convey("摘除一个服务IP时请求IP自动摘除", func() {
		//	testHTTP.Stop(":50001")
		//	testHTTP.Stop(":50002")
		//	time.Sleep(10 * time.Second)
		//	checkHostSlice, err := getUniqueHost(10, checkURL)
		//	convey.So(err, convey.ShouldBeNil)
		//	convey.So(checkHostSlice, convey.ShouldResemble, []string{"127.0.0.1:50003"})
		//	testHTTP.Run(":50001", false)
		//	testHTTP.Run(":50002", false)
		//})

		//convey.Convey("摘除全部服务IP时请求IP全部摘除", func() {
		//	testHTTP.Stop(":50001")
		//	testHTTP.Stop(":50002")
		//	testHTTP.Stop(":50003")
		//	time.Sleep(12 * time.Second)
		//	checkHostSlice2, _ := getUniqueHost(2, checkURL)
		//	convey.So(len(checkHostSlice2), convey.ShouldEqual, 0)
		//	testHTTP.Run(":50001", false)
		//	testHTTP.Run(":50002", false)
		//	testHTTP.Run(":50003", false)
		//})
	})
}

//TestTCPCheckIPList TCP服务发现测试
func TestTCPCheckIPList(t *testing.T) {
	tcpConfig := service.SysConfMgr.GetModuleConfigByName("tester_tcp")
	checkTCPAddr :=fmt.Sprintf("%s%s",lib.GetStringConf("base.cluster.cluster_ip"),tcpConfig.Base.FrontendAddr)
	convey.Convey("TCP服务发现测试", t, func() {

		convey.Convey("配置IP与请求IP一致", func() {
			time.Sleep(10 * time.Second)
			checkHostSlice, err := tcpGetUniqueHost(10, checkTCPAddr)
			convey.So(err, convey.ShouldBeNil)
			convey.So(checkHostSlice, convey.ShouldResemble, strings.Split(tcpConfig.LoadBalance.IPList, ","))
		})

		convey.Convey("摘除一个服务IP时请求IP自动摘除", func() {
			testTCP.Stop(":51001")
			testTCP.Stop(":51002") //明显是关闭时连接还存在。。。
			time.Sleep(10 * time.Second)
			checkHostSlice, err := tcpGetUniqueHost(10, checkTCPAddr)
			convey.So(err, convey.ShouldBeNil)
			convey.So(checkHostSlice, convey.ShouldResemble, []string{"127.0.0.1:51003"})
			testTCP.Run(":51001", false)
			testTCP.Run(":51002", false)
			time.Sleep(10 * time.Second)
		})

		//convey.Convey("摘除全部服务IP时请求IP全部摘除", func() {
		//	testTCP.Stop(":51001")
		//	testTCP.Stop(":51002")
		//	testTCP.Stop(":51003")
		//	time.Sleep(12 * time.Second)
		//	checkHostSlice2, err := tcpGetUniqueHost(10, checkTCPAddr)
		//	convey.So(err, convey.ShouldBeNil)
		//	convey.So(len(checkHostSlice2), convey.ShouldEqual, 0)
		//	testTCP.Run(":51001", false)
		//	testTCP.Run(":51002", false)
		//	testTCP.Run(":51003", false)
		//})
	})
}

func tcpGetUniqueHost(times int, addr string) ([]string, error) {
	checkHostSlice := []string{}
loop:
	for i := 0; i < times; i++ {
		tSocket, err := thrift.NewTSocket(addr)
		if err != nil {
			log.Println("tSocket error:", err)
			break loop
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, _ := transportFactory.GetTransport(tSocket)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thriftgen.NewFormatDataClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			log.Println("Error opening:", addr)
			break loop
		}
		defer transport.Close()
		data := thriftgen.Data{Text: "ping"}
		d, err := client.DoFormat(context.Background(), &data)
		if err != nil {
			log.Println("Error opening:", addr)
			break loop
		}
		currentIP := d.Text
		//fmt.Println("currentIP",currentIP)
		match, _ := regexp.MatchString("^[a-z0-9.]+:[0-9]+$", currentIP)
		if !match {
			break loop
		}
		for _, cip := range checkHostSlice {
			if cip == currentIP {
				break loop
			}
		}
		checkHostSlice = append(checkHostSlice, d.Text)
	}
	return checkHostSlice, nil
}

func getUniqueHost(times int, checkURL string) ([]string, error) {
	checkHostSlice := []string{}
loop:
	for i := 0; i < times; i++ {
		client := http.Client{}
		resp, err := client.Get(checkURL)
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			break loop
		}
		queryBody, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			break loop
		}
		currentIP := string(queryBody)
		match, _ := regexp.MatchString("^[a-z0-9.]+:[0-9]+$", currentIP)
		if !match {
			break loop
		}
		for _, cip := range checkHostSlice {
			if cip == currentIP {
				break loop
			}
		}
		checkHostSlice = append(checkHostSlice, string(queryBody))
	}
	return checkHostSlice, nil
}
