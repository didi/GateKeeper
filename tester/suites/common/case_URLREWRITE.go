package common

import (
	"fmt"
	"github.com/didi/gatekeeper/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

//TestURLWrite Url地址重写测试
func TestURLWrite(t *testing.T) {
	//公共配置
	moduleConfig := service.SysConfMgr.GetModuleConfigByName("tester")
	checkURL := fmt.Sprintf("http://%s%s%s/get_path",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		moduleConfig.MatchRule[0].Rule)

	//公共配置
	nuwConfig := service.SysConfMgr.GetModuleConfigByName("without_urlwrite")
	nuwURL := fmt.Sprintf("http://%s%s%s/get_path",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		nuwConfig.MatchRule[0].Rule)

	convey.Convey("Url地址重写测试", t, func() {
		convey.Convey("不配置重写规则", func() {
			hostPath, err := getURLContent(nuwURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(hostPath, convey.ShouldEqual, "/gatekeeper/without_urlwrite/get_path")
		})
		convey.Convey("配置重写规则", func() {
			hostPath, err := getURLContent(checkURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(hostPath, convey.ShouldEqual, "/get_path")
		})
	})
}

func getURLContent(checkURL string) (string, error) {
	client := http.Client{}
	resp, err := client.Get(checkURL)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return "", err
	}
	queryBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	hostPath := string(queryBody)
	return hostPath, nil
}
