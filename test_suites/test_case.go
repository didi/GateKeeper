package suites

import (
	"encoding/json"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/handler"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGoConvey(t *testing.T) {
	url1 := "http://127.0.0.1:58080/test_service_name"
	convey.Convey("集成测试", t, func() {
		convey.Convey("访问测试1", func() {
			c, err := getURLContent(url1)
			convey.So(err, convey.ShouldBeNil)
			convey.So(c, convey.ShouldEqual, "/test_service_name")
		})
	})
}

func TestHostServiceVisit(t *testing.T) {
	//url1 := "http://www.test.com/test_service_name"
	convey.Convey("测试服务域名", t, func() {
		convey.Convey("清空测试库", func() {
			//delete all
		})
		convey.Convey("插入测试服务信息", func() {
			//insert new service
		})
		convey.Convey("更新服务配置", func() {
			log.Info().Msg(lib.Purple("watching load service config from resource"))
			handler.ServiceManagerHandler.Load()
		})
		convey.Convey("访问请求", func() {

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

func getErrmsg(c string) string {
	m := map[string]string{}
	if err := json.Unmarshal([]byte(c), m); err != nil {
		return ""
	}
	if msg, ok := m["errmsg"]; ok {
		return msg
	}
	return ""
}
