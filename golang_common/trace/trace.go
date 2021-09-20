// Package trace
package trace

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ip     string
	ctxKey key
)

const (
	//TraceId
	DIDI_HEADER_RID = "Didi-Header-Rid"
	//SpanId
	DIDI_HEADER_SPANID       = "Didi-Header-Spanid"
	DIDI_HEADER_HINT_CODE    = "Didi-Header-Hint-Code"
	DIDI_HEADER_HINT_CONTENT = "Didi-Header-Hint-Content"

	EMPRY_TRACE_ID = ""
	TRACEID        = "traceid"
	SPANID         = "spanid"
	CALLER         = "caller"
	SRC_METHOD     = "srcMethod"
	HINT_CODE      = "hintCode"
	HINT_CONTENT   = "hintContent"

	DLTagUndefined        = " _undef"
	DLTagMysqlFailed      = " _com_mysql_failure"
	DLTagMysqlSuccess     = " _com_mysql_success"
	DLTagRedisFailed      = " _com_redis_failure"
	DLTagRedisSuccess     = " _com_redis_success"
	DLTagThriftFailed     = " _com_thrift_failure"
	DLTagThriftSuccess    = " _com_thrift_success"
	DLTagHTTPSuccess      = " _com_http_success"
	DLTagHTTPFailed       = " _com_http_failure"
	DLTagBackendRPCFailed = " _com_interactive_failure"
	DLTagRequestIn        = " _com_request_in"
	DLTagRequestOut       = " _com_request_out"

	EMPTY_SPAN_ID             string = "0"
	HINT_PRESSURE_TRAFFIC            = "1"
	HINT_NORMAL_TRAFFIC              = "0"
	HINT_TRACE_SAMPLE_INIT    int64  = 0x80
	HINT_TRACE_SAMPLE         int64  = 0x40
	DEFAULT_TRACE_SAMPLE_RATE int64  = 10000
)

type Trace struct {
	TraceId      string
	SpanId       string
	Caller       string
	SrcMethod    string
	Method       string
	Host         string
	URL          string
	Params       string
	From         string
	HintCode     string
	Sampling     int
	HintContent  string
	FormatString string
}

type HintContent struct {
	Sample HintSampling `json:"Sample"`
}

type HintSampling struct {
	Rate int64 `json:"Rate"`
	Code int64 `json:"Code"`
}

type key int

func (k *key) GetCtxKey() interface{} {
	return ctxKey
}

type ICtxKey interface {
	GetCtxKey() interface{}
}

// NewCtxKey return  ctx key for get trace
func NewCtxKey() ICtxKey {
	return new(key)
}

//SetCtxTrace set context ctxKey *Trace
func SetCtxTrace(ctx context.Context, val *Trace) context.Context {
	return context.WithValue(ctx, ctxKey, val)
}

//GetCtxTrace get context *Trace
func GetCtxTrace(ctx context.Context) (*Trace, bool) {
	val, ok := ctx.Value(ctxKey).(*Trace)
	return val, ok
}

var tracePool = sync.Pool{
	New: func() interface{} {
		return new(Trace)
	},
}

func PutTrace(trace *Trace) {
	tracePool.Put(trace)
}

// New Trace，req 可为nil
func New(req *http.Request) (trace *Trace) {
	trace = tracePool.Get().(*Trace)
	trace.TraceId = trace.GetTraceId(req)
	trace.SpanId = trace.GetSpanId(req)
	trace.HintCode = trace.GetHintCode(req)
	trace.HintContent = trace.GetHintContent(req)
	trace.Method = trace.GetHttpMethod(req)
	trace.URL = trace.GetHttpURL(req)
	trace.Params = trace.GetHttpParams(req)
	trace.Host = trace.GetHttpHost(req)
	trace.From = trace.GetClientAddr(req)
	trace.FormatString = fmt.Sprintf("traceid=%s||spanid=%s||hintCode=%s||hintContent=%s||method=%s||host=%s||uri=%s||params=%s||from=%s||srcMethod=%s||caller=%s", trace.TraceId, trace.SpanId, trace.HintCode, trace.HintContent, trace.Method, trace.Host, trace.URL, trace.Params, trace.From, trace.SrcMethod, trace.Caller)
	return
}

func (tr *Trace) genFormatString(trace *Trace) string {
	return fmt.Sprintf("traceid=%s||spanid=%s||hintCode=%s||hintContent=%s||method=%s||host=%s||uri=%s||params=%s||from=%s||srcMethod=%s||caller=%s",
		trace.TraceId, trace.SpanId, trace.HintCode, trace.HintContent, trace.Method, trace.Host, trace.URL, trace.Params, trace.From, trace.SrcMethod, trace.Caller)
}

func FormatCtx(ctx context.Context) string {
	if ctx == nil {
		return "ctx_format=unset"
	}
	trace, ok := ctx.Value(ctxKey).(*Trace)
	if !ok {
		return "ctx_format=unset"
	}
	return trace.String()
}

func NewWithMap(m map[string]string) (trace *Trace) {
	trace = &Trace{}
	for k, v := range m {
		canonicalHeader := http.CanonicalHeaderKey(strings.TrimSpace(k))
		m[canonicalHeader] = v
	}
	if val, ok := m[DIDI_HEADER_RID]; ok {
		trace.TraceId = val
	} else {
		trace.genTraceId()
	}
	if val, ok := m[DIDI_HEADER_SPANID]; ok {
		trace.SpanId = val
	} else {
		trace.SpanId = trace.GenSpanId()
	}
	if val, ok := m[DIDI_HEADER_HINT_CODE]; ok {
		trace.HintCode = val
	}
	if val, ok := m[DIDI_HEADER_HINT_CONTENT]; ok {
		trace.HintContent = val
	}
	trace.FormatString = trace.genFormatString(trace)
	return
}

//AddHttpHeader 向http request header中添加trace信息
func (tr *Trace) AddHttpHeader(request *http.Request) {
	if request == nil {
		return
	}
	request.Header.Set(DIDI_HEADER_RID, tr.TraceId)
	request.Header.Set(DIDI_HEADER_SPANID, tr.GenSpanId())
	request.Header.Set(DIDI_HEADER_HINT_CODE, tr.HintCode)
	request.Header.Set(DIDI_HEADER_HINT_CONTENT, tr.HintContent)
}

//IsPressureTraffic 判断是否是压测流量
func (tr *Trace) IsPressureTraffic() bool {
	return tr.HintCode == HINT_PRESSURE_TRAFFIC
}

// GetTraceId 获取TraceId，如果不存在就生成
func (tr *Trace) GetTraceId(req *http.Request) string {
	if tr.TraceId == "" && req != nil {
		tr.TraceId = req.Header.Get(DIDI_HEADER_RID)
		if tr.TraceId == "" {
			tr.genTraceId()
		}
	}

	return tr.TraceId
}

func (tr *Trace) GetHttpMethod(req *http.Request) string {
	if tr.Method == "" && req != nil {
		tr.Method = req.Method
	}

	return tr.Method
}

func (tr *Trace) GetHttpURL(req *http.Request) string {
	if tr.URL == "" && req != nil {
		tr.URL = req.URL.Path
	}

	return tr.URL
}

func (tr *Trace) GetHttpParams(req *http.Request) string {
	if tr.Params == "" && req != nil {
		tr.Params = req.URL.Query().Encode()
	}

	return tr.Params
}

func (tr *Trace) GetHttpHost(req *http.Request) string {
	if tr.Host == "" && req != nil {
		tr.Host = req.Host
	}

	return tr.Host
}

func (tr *Trace) GetClientAddr(req *http.Request) string {
	if tr.From == "" && req != nil {
		tr.From = GetClientAddr(req)
	}

	return tr.From
}

func (tr *Trace) genTraceId() {
	tr.TraceId = GenTraceId()
}

// GetClientAddr ...
func GetClientAddr(req *http.Request) string {
	addr := req.Header.Get("X-Real-IPV6")
	if net.ParseIP(addr) != nil {
		return addr
	}
	addr = req.Header.Get("X-Real-IP")
	if net.ParseIP(addr) != nil {
		return addr
	}
	addr = req.Header.Get("X-Forwarded-For")
	if net.ParseIP(addr) != nil {
		return addr
	}
	return req.RemoteAddr
}

//GenTraceId
func GenTraceId() string {
	if ip == "" {
		ip = getLocalIP()
	}
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()
	b := bytes.Buffer{}

	b.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))
	b.WriteString(fmt.Sprintf("%x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0")

	return b.String()
}

// GetSpanId 获取TraceId，如果不存在就生成
func (tr *Trace) GetSpanId(req *http.Request) string {
	if tr.SpanId == "" && req != nil {
		tr.SpanId = req.Header.Get(DIDI_HEADER_SPANID)
		if tr.SpanId == "" {
			tr.SpanId = tr.GenSpanId()
		}
	}

	return tr.SpanId
}

// GetHintCode 获取HintCode
func (tr *Trace) GetHintCode(req *http.Request) string {
	if tr.HintCode == "" && req != nil {
		tr.HintCode = req.Header.Get(DIDI_HEADER_HINT_CODE)
	}

	return tr.HintCode
}

// GetHintContent 获取HintContent
func (tr *Trace) GetHintContent(req *http.Request) string {
	if tr.HintContent == "" && req != nil {
		tr.HintContent = req.Header.Get(DIDI_HEADER_HINT_CONTENT)
	}

	return tr.HintContent
}

// GenSpanId 生成新的spanId 用作 CSPanid
func (tr *Trace) GenSpanId() string {
	return GenSpanId()
}

// GenSpanId 生成新的spanId 用作 CSPanid
func GenSpanId() string {
	return fmt.Sprintf("%x", rand.Int63())
}

func getLocalIP() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	return ip
}

//IsTraceSampleEnabled 判断采样相关逻辑
func (tr *Trace) IsTraceSampleEnabled() (hit bool) {

	var SampleHint HintContent
	err := json.Unmarshal([]byte(tr.HintContent), &SampleHint)

	if err != nil {
		SampleHint = HintContent{
			Sample: HintSampling{
				Code: 0,
				Rate: DEFAULT_TRACE_SAMPLE_RATE,
			},
		}
	} else {
		if SampleHint.Sample.Rate == 0 {
			SampleHint.Sample.Rate = DEFAULT_TRACE_SAMPLE_RATE
		}
	}

	if HINT_TRACE_SAMPLE_INIT == (SampleHint.Sample.Code & HINT_TRACE_SAMPLE_INIT) {
		if HINT_TRACE_SAMPLE == (SampleHint.Sample.Code & HINT_TRACE_SAMPLE) {
			hit = true
		} else {
			hit = false
		}
	} else {
		rndMax := DEFAULT_TRACE_SAMPLE_RATE

		if SampleHint.Sample.Rate > 0 {
			rndMax = SampleHint.Sample.Rate
		}

		rnd := rand.Int63n(SampleHint.Sample.Rate) + 1
		if rnd == rndMax {
			hit = true
		}
	}

	return
}

func (tr *Trace) String() string {
	if len(tr.FormatString) == 0 {
		tr.FormatString = tr.genFormatString(tr)
	}

	return tr.FormatString
}
