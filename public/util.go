package public

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/golang_common/log"
	"github.com/garyburd/redigo/redis"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

//MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//AuthIPList 验证ip名单
func AuthIPList(clientIP string, whiteList []string) bool {
	return InStringList(clientIP, whiteList)
}

//CheckConnPort 检查端口是否被占用
func CheckConnPort(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	ln.Close()
	return nil
}

//InStringList 数组中是否存在某值
func InStringList(t string, list []string) bool {
	for _, s := range list {
		if s == t {
			return true
		}
	}
	return false
}

//InOrPrefixStringList 字符串在string数组 或者 字符串前缀在数组中
func InOrPrefixStringList(t string, arr []string) bool {
	for _, s := range arr {
		if t == s {
			return true
		}
		if s != "" && strings.HasPrefix(t, s) {
			return true
		}
	}
	return false
}

//Substr 字符串的截取
func Substr(str string, start int64, end int64) string {
	length := int64(len(str))
	if start < 0 || start > length {
		return ""
	}
	if end < 0 {
		return ""
	}
	if end > length {
		end = length
	}
	return string(str[start:end])
}

//MapSorter map排序，按key排序
type MapSorter []MapItem

//NewMapSorter 新排序
func NewMapSorter(m map[string]string) MapSorter {
	ms := make(MapSorter, 0, len(m))
	for k, v := range m {
		ms = append(ms, MapItem{Key: k, Val: v})
	}
	sort.Sort(ms)
	return ms
}

//MapItem 排序对象
type MapItem struct {
	Key string
	Val string
}

//Len 对象长度
func (ms MapSorter) Len() int {
	return len(ms)
}

//Swap 交换位置
func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

//Less 按首字母键排序
func (ms MapSorter) Less(i, j int) bool {
	return ms[i].Key < ms[j].Key
}

//GetSign 获取签名
func GetSign(paramMap map[string]string, secret string) string {
	paramArr := NewMapSorter(paramMap)
	str := ""
	for _, v := range paramArr {
		str = str + fmt.Sprintf("%s=%s&", v.Key, url.QueryEscape(v.Val))
	}
	str = str + secret

	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)
	return md5Str[7:23]
}

//RemoteIP 获取远程IP
func RemoteIP(req *http.Request) string {
	var err error
	var remoteAddr = req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, err = net.SplitHostPort(remoteAddr)
	}
	if err != nil {
		return ""
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

//ParseGzip 解析gzip
func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	undatas, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return undatas, nil
}

//HTTPGET get请求
func HTTPGET(log *log.Logger, urlString string, urlParams url.Values, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	urlString = lib.AddGetDataToUrl(urlString, urlParams)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Warn(
			"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|err=%v",
			"_com_http_failure",
			urlString,
			float32(time.Now().UnixNano()-startTime)/1.0e9,
			"GET",
			urlParams,
			err.Error())
		return nil, nil, err
	}
	if len(header) > 0 {
		req.Header = header
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(
			"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|err=%v",
			"_com_http_failure",
			urlString,
			float32(time.Now().UnixNano()-startTime)/1.0e9,
			"GET",
			urlParams,
			err.Error())
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(
			"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|err=%v",
			"_com_http_failure",
			urlString,
			float32(time.Now().UnixNano()-startTime)/1.0e9,
			"GET",
			urlParams,
			err.Error())
		return nil, nil, err
	}
	log.Info(
		"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|result=%v",
		"_com_http_success",
		urlString,
		float32(time.Now().UnixNano()-startTime)/1.0e9,
		"GET",
		urlParams,
		string(body))
	return resp, body, nil
}

//HTTPPOST post请求
func HTTPPOST(log *log.Logger, urlString string, urlParams url.Values, msTimeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	if contextType == "" {
		contextType = "application/x-www-form-urlencoded"
	}
	req, err := http.NewRequest("POST", urlString, strings.NewReader(urlParams.Encode()))
	if len(header) > 0 {
		req.Header = header
	}
	req.Header.Set("Content-Type", contextType)
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(
			"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|err=%v",
			"_com_http_failure",
			urlString,
			float32(time.Now().UnixNano()-startTime)/1.0e9,
			"POST",
			urlParams,
			err.Error())
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(
			"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|err=%v",
			"_com_http_failure",
			urlString,
			float32(time.Now().UnixNano()-startTime)/1.0e9,
			"POST",
			urlParams,
			err.Error())
		return nil, nil, err
	}
	log.Info(
		"dltag=%v|url=%v|proc_time=%v|method=%v|args=%v|result=%v",
		"_com_http_success",
		urlString,
		float32(time.Now().UnixNano()-startTime)/1.0e9,
		"GET",
		urlParams,
		string(body))
	return resp, body, nil
}

//RedisLogDo redis单次请求
func RedisLogDo(log *log.Logger, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		log.Warn(
			"dltag=%v|method=%v|err=%v|bind=%v|proc_time=%v",
			"_com_redis_failure",
			commandName,
			err,
			args,
			fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()))
	} else {
		replyStr, _ := redis.String(reply, nil)
		log.Info(
			"dltag=%v|method=%v|bind=%v|reply=%v|proc_time=%v",
			"_com_redis_success",
			commandName,
			args,
			replyStr,
			fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()))
	}
	return reply, err
}

//RedisConfPipline redis pip请求
func RedisConfPipline(log *log.Logger, name string, pip ...func(c redis.Conn)) error {
	c, err := lib.RedisConnFactory(name)
	if err != nil {
		log.Warn(
			"dltag=%v|name=%v|err=%v",
			"_com_redis_failure",
			name,
			err)
		return err
	}
	defer c.Close()
	for _, f := range pip {
		f(c)
	}
	c.Flush()
	return nil
}

//RedisConfDo 通过配置 执行redis
func RedisConfDo(log *log.Logger, name string, commandName string, args ...interface{}) (interface{}, error) {
	c, err := lib.RedisConnFactory(name)
	if err != nil {
		log.Warn(
			"dltag=%v|method=%v|err=%v|bind=%v",
			"_com_redis_failure",
			commandName,
			err,
			args)
		return nil, err
	}
	defer c.Close()

	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		log.Warn(
			"dltag=%v|method=%v|err=%v|bind=%v|proc_time=%v",
			"_com_redis_failure",
			commandName,
			err,
			args,
			fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()))
	} else {
		replyStr, _ := redis.String(reply, nil)
		log.Info(
			"dltag=%v|method=%v|bind=%v|reply=%v|proc_time=%v",
			"_com_redis_success",
			commandName,
			args,
			replyStr,
			fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()))
	}
	return reply, err
}
