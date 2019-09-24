package public

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)


//Context 对response和request方法的封装
type Context struct {
	Res        http.ResponseWriter
	Req        *http.Request
	StatusCode int
	urlValue   url.Values
	formValue  url.Values
	done       bool
}

//NewContext 构造方法
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Res:       w,
		Req:       r,
		urlValue:  nil,
		formValue: nil,
	}
}

//Exists 参数是否存在
func (c *Context) Exists(k string) bool {
	_ = c.Query("")
	formValue := map[string][]string(c.formValue)
	urlValue := map[string][]string(c.urlValue)
	if len(formValue[k]) != 0 {
		return true
	}
	if len(urlValue[k]) != 0 {
		return true
	}
	return false
}

//Query 获取参数
func (c *Context) Query(k string) string {
	if c.Method() == "GET" {
		if c.urlValue == nil {
			c.urlValue = c.Req.URL.Query()
		}
		return c.urlValue.Get(k)
	}
	if c.formValue == nil || c.urlValue == nil {
		c.Req.ParseForm()
		c.formValue = c.Req.Form
		c.urlValue = c.Req.URL.Query()
	}
	if v := c.formValue.Get(k); v != "" {
		return v
	}
	return c.urlValue.Get(k)
}

//QueryInt 获取请求参数，转换为int
func (c *Context) QueryInt(k string) (int, error) {
	sv := c.Query(k)
	return strconv.Atoi(sv)
}

//QueryInt64 获取请求参数，转换为int64
func (c *Context) QueryInt64(k string) (int64, error) {
	sv := c.Query(k)
	return strconv.ParseInt(sv, 10, 64)
}

//QueryBool 获取请求参数，转换为bool
func (c *Context) QueryBool(k string) (bool, error) {
	sv := c.Query(k)
	return strconv.ParseBool(sv)
}

//Cookie 通过name,获取cookie
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.Req.Cookie(name)
}

//Cookies 获取cookie数组
func (c *Context) Cookies() []*http.Cookie {
	return c.Req.Cookies()
}

//SetCookie 设置cookie内容
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Res, cookie)
}

//File 获取文件内容
func (c *Context) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.Req.FormFile(name)
}

//Method 获取请求方法
func (c *Context) Method() string {
	return c.Req.Method
}

//URI 获取请求方法
func (c *Context) URI() string {
	return c.Req.RequestURI
}

//Path 获取请求路径
func (c *Context) Path() string {
	return c.Req.URL.Path
}

//Schema 获取协议
func (c *Context) Schema() string {
	if c.Req.TLS != nil{
		return "https://"
	}
	return "http://"
}

//Host 获取请求域名
func (c *Context) Host() string {
	return c.Req.Host
}

//RemoteAddr 获取客户ip
func (c *Context) RemoteAddr() string {
	return c.Req.RemoteAddr
}

//IsAjaxReq 是否为ajax请求
func (c *Context) IsAjaxReq() bool {
	s := c.Req.Header.Get("HTTP_X_REQUESTED_WITH")
	s = strings.ToLower(s)
	return s == "xmlhttprequest"
}

//IsBrowser 是否为浏览器
func (c *Context) IsBrowser() bool {
	s := c.Req.Header.Get("Accept")
	return s!="*/*"
}

//AcceptJSON 是否为json请求
func (c *Context) AcceptJSON() bool {
	accept := c.Req.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

//JSON 写入response内容
func (c *Context) JSON(data interface{}) {
	if ct := c.Res.Header().Get("Content-Type"); ct == "" {
		c.Res.Header().Set("Content-Type", "application/json")
	}
	j, err := json.Marshal(data)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		c.Write([]byte(fmt.Sprintf("%s:line %d, json marshal error:%v", f, l, err)))
		return
	}
	c.Write(j)
}

//Redirect 跳转地址
func (c *Context) Redirect(location string, code ...int) {
	if len(code) != 0 {
		c.SetStatusCode(code[0])
	} else {
		c.SetStatusCode(303)
	}
	c.done = true
	http.Redirect(c.Res, c.Req, location, c.StatusCode)
}

//SetStatusCode 设置状态码
func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	c.Res.WriteHeader(code)
}

//Success 成功返回
func (c *Context) Success(data interface{}) {
	c.SetStatusCode(200)
	c.JSON(map[string]interface{}{
		"errno":  0,
		"errmsg": "",
		"data":   data,
	})
}

//Location 跳转到
func (c *Context) Location(message,url string) {
	content:="<script>alert('"+message+"');location.href='"+url+"'</script>"
	c.Write([]byte(content))
}

//Error 错误返回
func (c *Context) Error(code int, msg string) {
	c.JSON(map[string]interface{}{
		"errno":  code,
		"errmsg": msg,
		"data":   "",
	})
}

//String 写入string
func (c *Context) String(s string) {
	c.Write([]byte(s))
}

//Write 往response写入数据
func (c *Context) Write(data []byte) {
	if ct := c.Res.Header().Get("Content-Type"); ct == "" {
		c.Res.Header().Set("Content-Type", "text/plain")
	}
	if !c.done {
		if c.StatusCode == 0 {
			c.SetStatusCode(200)
		}
		c.Res.Write(data)
		c.done = true
	}
}

//Map map[string]别名
type Map map[string]interface{}
//
////Handler func别名
//type Handler func(*Context)
//
////NewHttpHandler func
//func (h Handler) NewHttpHandler() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx := NewContext(w, r)
//		h(ctx)
//	}
//}
//
////ServeHTTP func
//func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	handle := h.NewHttpHandler()
//	handle(w, r)
//}