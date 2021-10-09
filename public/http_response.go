package public

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type ResponseCode int

//1000以下为通用码，1000以上为用户自定义码
const (
	SuccessCode ResponseCode = iota
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode
	InvalidRequestErrorCode ResponseCode = 401
	CustomizeCode           ResponseCode = 1000
	GROUPALL_SAVE_FLOWERROR ResponseCode = 2001
)

type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceId   interface{}  `json:"trace_id"`
	Stack     interface{}  `json:"stack"`
}

func ResponseError(c *gin.Context, code ResponseCode, err error) {
	traceContext := GetGinTraceContext(c)
	errMsg := fmt.Sprintf("%+v", err)

	straceMsg := ""
	tmpStack := strings.Split(errMsg, "||")
	if len(tmpStack) == 2 {
		errMsg = tmpStack[0]
		straceMsg = tmpStack[1]
	}
	strackList := strings.Split(straceMsg, "\n")
	for i, t := range strackList {
		t = strings.Replace(t, "\t", "  ", -1)
		strackList[i] = t
	}
	resp := &Response{ErrorCode: code, ErrorMsg: errMsg, Data: "", TraceId: traceContext.TraceId, Stack: strackList}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(200, err)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	traceContext := GetGinTraceContext(c)
	resp := &Response{ErrorCode: SuccessCode, ErrorMsg: "", Data: data, TraceId: traceContext.TraceId}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
