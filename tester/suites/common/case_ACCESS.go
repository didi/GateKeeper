package common

import (
	"encoding/json"
	"fmt"
	"github.com/didi/gatekeeper/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//TestAccessControl 权限测试
func TestAccessControl(t *testing.T) {
	//公共配置
	config1 := service.SysConfMgr.GetModuleConfigByName("tester_whitelist")
	url1 := fmt.Sprintf("http://%s%s%s/get_host",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		config1.MatchRule[0].Rule)
	config2 := service.SysConfMgr.GetModuleConfigByName("tester_blocklist")
	url2 := fmt.Sprintf("http://%s%s%s/get_host",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		config2.MatchRule[0].Rule)

	convey.Convey("权限测试", t, func() {
		convey.Convey("黑名单ip测试", func() {
			c,err:= getURLContent(url2)
			convey.So(err,convey.ShouldBeNil)
			convey.So(getErrmsg(c),convey.ShouldNotEqual,"msg:AuthInBlackIPList")
		})
		convey.Convey("白名单ip测试", func() {
			c,err:= getURLContent(url1)
			convey.So(err,convey.ShouldBeNil)
			convey.So(c,convey.ShouldEqual,"127.0.0.1:50003")
		})
	})
}

func getErrmsg(c string)  string{
	m:=map[string]string{}
	if err:=json.Unmarshal([]byte(c),m);err!=nil{
		return ""
	}
	if msg,ok:=m["errmsg"];ok{
		return msg
	}
	return ""
}
