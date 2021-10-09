package suites

import (
	"encoding/json"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGoConvey(t *testing.T) {
	url1 := "http://127.0.0.1:58080/test_service_name"
	convey.Convey("集成测试", t, func() {
		convey.Convey("访问测试", func() {
			c, err := getURLContent(url1)
			convey.So(err, convey.ShouldBeNil)
			convey.So(c, convey.ShouldEqual, "/test_service_name")
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
