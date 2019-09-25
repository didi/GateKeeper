package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/didi/gatekeeper/dao"
	"github.com/didi/gatekeeper/public"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

//BeforeRequestAuthRegisterFuncs 验证方法列表
var BeforeRequestAuthRegisterFuncs [] func(m *dao.GatewayModule, req *http.Request, res http.ResponseWriter) (bool,error)

//ModifyResponseRegisterFuncs 过滤方法列表
var ModifyResponseRegisterFuncs [] func(m *dao.GatewayModule, req *http.Request, res *http.Response) error

//RegisterBeforeRequestAuthFunc 注册请求前验证请求方法
func RegisterBeforeRequestAuthFunc(funcs ...func(m *dao.GatewayModule, req *http.Request, res http.ResponseWriter) (bool,error)) {
	BeforeRequestAuthRegisterFuncs = append(BeforeRequestAuthRegisterFuncs, funcs...)
}

//RegisterModifyResponseFunc 注册请求后修改response方法
func RegisterModifyResponseFunc(funcs ...func(m *dao.GatewayModule, req *http.Request, res *http.Response) error) {
	ModifyResponseRegisterFuncs = append(ModifyResponseRegisterFuncs, funcs...)
}

//FilterCityData 过滤数据函数
func FilterCityData(filterURLs []string) func(m *dao.GatewayModule, req *http.Request, res *http.Response) error{
	return func(m *dao.GatewayModule, req *http.Request, res *http.Response) error {
		//获取原始请求地址
		v:=req.Context().Value(public.ContextKey("request_url"))
		requestURL,ok := v.(string)
		if !ok{
			requestURL = req.URL.Path
		}

		//获取请求内容
		payload, err := ioutil.ReadAll(res.Body)
		if err!=nil{
			return err
		}

		//验证是否匹配
		for _,matchURL:=range filterURLs{
			if matchURL==requestURL {
				//过滤规则
				filterData, err := filterJSONTreeByKey(string(payload),"data.list", "city_id", []string{"12"},)
				if err!=nil{
					return err
				}
				payload = []byte(filterData)
				break
			}
		}

		//重写请求内容
		res.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		res.ContentLength = int64(len(payload))
		res.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}
}

//ModifyResponse 修改返回内容
func ModifyResponse(m *dao.GatewayModule, req *http.Request, res *http.Response) (error) {
	for _, ff := range ModifyResponseRegisterFuncs {
		if err:=ff(m, req, res);err!=nil {
			return err
		}
	}
	return nil
}

//filterJSONTreeByKey 基于json_path过滤节点数据
func filterJSONTreeByKey(payload , jsonTree , filterKey string, filterIds []string) (string, error) {
	mapTest := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payload), &mapTest); err != nil {
		return payload, fmt.Errorf("json.Unmarshal err:%v",err)
	}
	dlRs := gjson.Get(payload, jsonTree)
	dataList := []interface{}{}
	for _, dlitem := range dlRs.Array() {
		for dKey, dValue := range dlitem.Map() {
			if dKey == filterKey && public.InStringList(dValue.String(), filterIds) {
				dlitemmap := map[string]interface{}{}
				json.Unmarshal([]byte(dlitem.String()), &dlitemmap)
				dataList = append(dataList, dlitemmap)
			}
		}
	}
	newPayload, err := sjson.Set(string(payload), jsonTree, dataList) //写入json树
	if err != nil {
		return "", fmt.Errorf("sjson.Set err:%v",err)
	}
	return newPayload, nil
}