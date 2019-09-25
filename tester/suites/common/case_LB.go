package common

import (
	"fmt"
	"github.com/didi/gatekeeper/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"
)

//TestRoundRobin 负载均衡测试
func TestRoundRobin(t *testing.T) {
	//rr:=public.NewWeightedRR(public.RR_NGINX)
	//rr.Add("5001",1)
	//rr.Add("5002",2)
	//rr.Add("5003",1)
	//for i:=0;i<10;i++{
	//	fmt.Println(rr.Next())
	//}

	//公共配置
	moduleConfig := service.SysConfMgr.GetModuleConfigByName("tester_lb")
	checkURL := fmt.Sprintf("http://%s%s%s/get_host",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		moduleConfig.MatchRule[0].Rule)
	convey.Convey("负载均衡测试", t, func() {
		convey.Convey("RoundRobin测试", func() {
			rrs,err:=getRRHosts(4,checkURL)
			convey.So(err,convey.ShouldBeNil)
			convey.So(rrs,convey.ShouldResemble,[]string{
				"127.0.0.1:50002",
				"127.0.0.1:50001",
				"127.0.0.1:50003",
				"127.0.0.1:50002",
			})
		})
	})
}

func getRRHosts(times int, checkURL string) ([]string, error) {
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
		match, _ := regexp.MatchString("^[a-z0-9.]+:[0-9]+$", string(queryBody))
		if !match {
			break loop
		}
		checkHostSlice = append(checkHostSlice, string(queryBody))
	}
	return checkHostSlice, nil
}
