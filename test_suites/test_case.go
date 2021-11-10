package suites

import (
	"encoding/json"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/test_suites/SqlHandler"
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

func TestStripPrefix(t *testing.T) {
	url1 := "http://127.0.0.1:58080/test_service_name/abc"
	columnServiceName := "test_service_name"
	convey.Convey("前缀剥离测试", t, func() {
		convey.Convey("清空测试库", func() {
			//delete all
			testsqlhandler.DeleteServiceInfo(columnServiceName)
		})
		convey.Convey("插入测试服务信息", func() {
			//insert new service
			serviceInfo := model.ServiceInfo{ID: 1, ServiceName: "test_service_name", ServiceDesc: "test_service_desc", HTTPPaths: "/test_service_name", HttpStripPrefix: 0, LoadBalanceStrategy: "round", LoadBalanceType: "default_loadbalance", UpstreamList: "http://127.0.0.1:8881 100", PluginConf: "{\"url_rewrite\":{\"rewrite_rule\":\"^/test_service(*) $1\"},\"http_flow_limit\":{\"service_flow_limit_num\":\"60\",\"service_flow_limit_type\":\"1\",\"clientip_flow_limit_num\":\"3\",\"clientip_flow_limit_type\":\"\"},\"header_transfer\":{\"header_transfer_rule\":\"add gatekeeper_power v2.0.1\"},\"http_whiteblacklist\":{\"ip_white_list\":\"\",\"url_white_list\":\"\"},\"http_upstream_transport\":{\"http_upstream_connection_timeout\":\"111\",\"http_upstream_header_timeout\":\"112\"},\"default_loadbalance\":{}}"}
			testsqlhandler.AddServiceInfo(&serviceInfo)
			fmt.Println("TEST ", testsqlhandler.GetServiceStripPrefix("test_service_name"))
		})
		convey.Convey("更新服务配置", func() {
			log.Info().Msg(lib.Purple("watching load service config from resource"))
			handler.ServiceManagerHandler.Load()
		})
		convey.Convey("访问请求", func() {
			c2, err := getURLContent(url1)
			convey.So(err, convey.ShouldBeNil)
			convey.So(c2, convey.ShouldEqual, "/test_service_name/abc")
		})

		convey.Convey("清空测试库2", func() {
			//delete all
			testsqlhandler.DeleteServiceInfo(columnServiceName)
		})
		convey.Convey("插入测试服务信息2", func() {
			//insert new service
			serviceInfo := model.ServiceInfo{ID: 1, ServiceName: "test_service_name", ServiceDesc: "test_service_desc", HTTPPaths: "/test_service_name", HttpStripPrefix: 1, LoadBalanceStrategy: "round", LoadBalanceType: "default_loadbalance", UpstreamList: "http://127.0.0.1:8881 100", PluginConf: "{\"url_rewrite\":{\"rewrite_rule\":\"^/test_service(*) $1\"},\"http_flow_limit\":{\"service_flow_limit_num\":\"60\",\"service_flow_limit_type\":\"1\",\"clientip_flow_limit_num\":\"3\",\"clientip_flow_limit_type\":\"\"},\"header_transfer\":{\"header_transfer_rule\":\"add gatekeeper_power v2.0.1\"},\"http_whiteblacklist\":{\"ip_white_list\":\"\",\"url_white_list\":\"\"},\"http_upstream_transport\":{\"http_upstream_connection_timeout\":\"111\",\"http_upstream_header_timeout\":\"112\"},\"default_loadbalance\":{}}"}
			testsqlhandler.AddServiceInfo(&serviceInfo)
			fmt.Println("TEST", testsqlhandler.GetServiceStripPrefix("test_service_name"))
		})
		convey.Convey("更新服务配置2", func() {
			log.Info().Msg(lib.Purple("watching load service config from resource"))
			handler.ServiceManagerHandler.Load()
		})
		convey.Convey("访问请求2", func() {
			c2, err := getURLContent(url1)
			convey.So(err, convey.ShouldBeNil)
			convey.So(c2, convey.ShouldEqual, "/abc")
		})
	})
}

func TestUpstreamList(t *testing.T) {
	url1 := "http://127.0.0.1:58080/test_service_name/get_host"
	columnServiceName := "test_service_name"
	convey.Convey("下游IP和权重测试", t, func() {
		convey.Convey("轮询", func() {
			testsqlhandler.DeleteServiceInfo(columnServiceName)
			serviceInfo := model.ServiceInfo{ID: 1, ServiceName: "test_service_name", ServiceDesc: "test_service_desc", HTTPPaths: "/test_service_name", HttpStripPrefix: 1, LoadBalanceStrategy: "round", LoadBalanceType: "default_loadbalance", UpstreamList: "http://127.0.0.1:8881 100\nhttp://127.0.0.1:8882 100", PluginConf: "{\"url_rewrite\":{\"rewrite_rule\":\"^/test_service(*) $1\"},\"http_flow_limit\":{\"service_flow_limit_num\":\"60\",\"service_flow_limit_type\":\"1\",\"clientip_flow_limit_num\":\"30\",\"clientip_flow_limit_type\":\"\"},\"header_transfer\":{\"header_transfer_rule\":\"add gatekeeper_power v2.0.1\"},\"http_whiteblacklist\":{\"ip_white_list\":\"\",\"url_white_list\":\"\"},\"http_upstream_transport\":{\"http_upstream_connection_timeout\":\"111\",\"http_upstream_header_timeout\":\"112\"},\"default_loadbalance\":{}}"}
			testsqlhandler.AddServiceInfo(&serviceInfo)
			//testsqlhandler.Save(&serviceInfo)
			handler.ServiceManagerHandler.Load()
			log.Info().Msg(lib.Purple("watching load service config from resource"))

			count1 := 0
			count2 := 0
			for i := 0; i < 10; i++ {
				c2, _ := getURLContent(url1)
				if c2 == "127.0.0.1:8881" {
					count1++
				}
				if c2 == "127.0.0.1:8882" {
					count2++
				}
			}

			convey.So(count1, convey.ShouldEqual, 5)
			convey.So(count2, convey.ShouldEqual, 5)
		})

		convey.Convey("权重轮询", func() {
			testsqlhandler.DeleteServiceInfo(columnServiceName)
			serviceInfo := model.ServiceInfo{ID: 1, ServiceName: "test_service_name", ServiceDesc: "test_service_desc", HTTPPaths: "/test_service_name", HttpStripPrefix: 1, LoadBalanceStrategy: "weight_round", LoadBalanceType: "default_loadbalance", UpstreamList: "http://127.0.0.1:8881 80\nhttp://127.0.0.1:8882 20", PluginConf: "{\"url_rewrite\":{\"rewrite_rule\":\"^/test_service(*) $1\"},\"http_flow_limit\":{\"service_flow_limit_num\":\"60\",\"service_flow_limit_type\":\"1\",\"clientip_flow_limit_num\":\"30\",\"clientip_flow_limit_type\":\"\"},\"header_transfer\":{\"header_transfer_rule\":\"add gatekeeper_power v2.0.1\"},\"http_whiteblacklist\":{\"ip_white_list\":\"\",\"url_white_list\":\"\"},\"http_upstream_transport\":{\"http_upstream_connection_timeout\":\"111\",\"http_upstream_header_timeout\":\"112\"},\"default_loadbalance\":{}}"}
			testsqlhandler.AddServiceInfo(&serviceInfo)
			//testsqlhandler.Save(&serviceInfo)
			fmt.Println("TEST", testsqlhandler.GetServiceLoadBalanceStrategy("test_service_name"))
			handler.ServiceManagerHandler.Load()
			log.Info().Msg(lib.Purple("watching load service config from resource"))
			count1 := 0
			count2 := 0
			for i := 0; i < 10; i++ {
				c2, _ := getURLContent(url1)
				if c2 == "127.0.0.1:8881" {
					count1++
				}
				if c2 == "127.0.0.1:8882" {
					count2++
				}
			}

			convey.So(count1, convey.ShouldEqual, 8)
			convey.So(count2, convey.ShouldEqual, 2)
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
