package controller

import (
	"github.com/didichuxing/gatekeeper/dao"
)

//APPDetailInfo app详情结构体
type APPDetailInfo struct {
	APPInfo       *dao.GatewayAPP
	DailyHourStat string //当日流量统计
	DailyHourAvg  string //当日流量统计
	DailyStatMax  int64  //当日流量统计
}

//ServiceDetailInfo 服务详情结构体
type ServiceDetailInfo struct {
	WeightList       []string //模块ip
	ModuleIPList     []string //模块ip
	ModuleIPCount    int      //模块ip数
	Module           *dao.GatewayModule
	ActiveIPList     []string //活动ip
	ActiveIPCount    int      //活动ip数
	ForbidIPList     []string //禁用ip
	ForbidIPCount    int      //禁用ip数
	AvaliableIPList  []string //可用ip
	AvaliableIPCount int      //可用ip数
	ClusterIP        string   //集群地址
	HTTPAddr         string   //http地址
	QPS              int64
	DayRequest       string
	DailyHourStat    string //当日流量统计
	DailyHourAvg     string
	DailyStatMax     int64 //当日流量统计

	//for edit
	MatchRule     string
	IPWeightList  string
	URLRewrite    string
	WhiteList     string
	BlackList     string
	WhiteHostName string
	Passport      string
	FilterRule    string
	RoutePrefix   string
}

//IsActive ip是否激活
func (u *ServiceDetailInfo) IsActive(ip string) bool {
	for _, item := range u.ActiveIPList {
		if item == ip {
			return true
		}
	}
	return false
}

//IsForbid ip是否禁用
func (u *ServiceDetailInfo) IsForbid(ip string) bool {
	for _, item := range u.ForbidIPList {
		if item == ip {
			return true
		}
	}
	return false
}

//APPListObj app列表结构体
type APPListObj struct {
	List      []APPItemObj
	ActiveURL string
}

//APPItemObj app对象结构体
type APPItemObj struct {
	*dao.GatewayAPP
	QPS int64
	QPD int64
}

//Admin admin结构体
type Admin struct {
}
