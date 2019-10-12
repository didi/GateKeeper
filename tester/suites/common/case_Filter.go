package common

import (
	"fmt"
	"github.com/didi/gatekeeper/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//TestFilter 内容过滤
func TestFilter(t *testing.T) {
	lib.Log.TagInfo(lib.NewTrace(), "_com_request_in", map[string]interface{}{
		"uri":    "uri",
		"method": "method",
	})

	config1 := service.SysConfMgr.GetModuleConfigByName("tester_filter")

	checkURL := fmt.Sprintf("http://%s%s%s/goods_list", lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"), config1.MatchRule[0].Rule)

	convey.Convey("内容过滤", t, func() {
		convey.Convey("json_path 过滤测试", func() {
			c, err := getURLContent(checkURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(c, convey.ShouldEqual, "{\"errno\":0,\"errmsg\":\"\",\"data\":{\"list\":[{\"city_id\":\"12\",\"name\":\"商品1\",\"pid\":\"2018103018_9000042581625\"}]}}")
		})
	})
}
