package common

import (
	"fmt"
	"github.com/didi/gatekeeper/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/smartystreets/goconvey/convey"
	"net/http"
	"sync"
	"testing"
	"time"
)

//TestAppAccess 租户测试
func TestAppAccess(t *testing.T) {
	lib.Log.TagInfo(lib.NewTrace(), "_com_request_in", map[string]interface{}{
		"uri":    "uri",
		"method": "method",
	})

	config1 := service.SysConfMgr.GetModuleConfigByName("tester_appaccess")
	appConf,err:=service.SysConfMgr.GetAppConfigByAPPID("test_app")
	if err!=nil{
		t.Errorf("GetAppConfigByAPPID err%v",err)
	}

	appWhiteIPConf,err:=service.SysConfMgr.GetAppConfigByAPPID("test_app_whiteip")
	if err!=nil{
		t.Errorf("GetAppConfigByAPPID err%v",err)
	}

	checkURL := fmt.Sprintf("http://%s%s%s/get_path?app_id=%s&sign=%s", lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		config1.MatchRule[0].Rule,appConf.AppID,appConf.Secret)

	pingURL := fmt.Sprintf("http://%s%s%s/ping",
		lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		config1.MatchRule[0].Rule)

	whiteIPURL := fmt.Sprintf("http://%s%s%s/get_path?app_id=%s", lib.GetStringConf("base.cluster.cluster_ip"),
		lib.GetStringConf("base.cluster.cluster_addr"),
		config1.MatchRule[0].Rule, appWhiteIPConf.AppID)

	convey.Convey("租户测试", t, func() {
		convey.Convey("租户签名测试", func() {
			time.Sleep(1*time.Second)
			hostPath,err:= getURLContent(checkURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(hostPath, convey.ShouldEqual, "/get_path")

			code,err:=getRequestStatusCode(pingURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(code, convey.ShouldEqual, 401)
		})
		convey.Convey("白名单测试", func() {
			hostPath,err:= getURLContent(whiteIPURL)
			convey.So(err, convey.ShouldBeNil)
			convey.So(hostPath, convey.ShouldEqual, "/get_path")
		})

		convey.Convey("Qps测试", func() {
			wg := sync.WaitGroup{}
			overflowFlag := false
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if isOverFlow(100, whiteIPURL) {
						overflowFlag = true
					}
				}()
			}
			wg.Wait()
			convey.So(overflowFlag, convey.ShouldBeTrue)
		})

	})
}

func getRequestStatusCode(checkURL string) (int,error){
	client := http.Client{}
	resp, err := client.Get(checkURL)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return 0,err
	}
	resp.Body.Close()
	if err != nil {
		return 0,err
	}
	return resp.StatusCode,nil
}
