package dashboard_controller

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/didi/gatekeeper/dashboard_middleware"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)
	group.POST("/service_add", service.ServiceAdd)
	group.POST("/service_update", service.ServiceUpdate)
	group.GET("/service_plugin_config", service.ServicePluginConfig)

	group.POST("/service_add_tcp", service.ServiceAddTcp)
	group.POST("/service_update_tcp", service.ServiceUpdateTcp)
	group.POST("/service_add_grpc", service.ServiceAddGrpc)
	group.POST("/service_update_grpc", service.ServiceUpdateGrpc)
}

// ServerList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &model.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	//从db中分页读取基本信息
	serviceInfo := &model.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}

	//格式化输出信息
	outList := []model.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			dashboard_middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := ""
		if serviceDetail.Info.ServiceType == public.LoadTypeHTTP {
			host := strings.Split(serviceDetail.Info.HTTPHosts, "\n")
			paths := strings.Split(serviceDetail.Info.HTTPPaths, "\n")
			for _, v := range host {
				for _, vs := range paths {
					serviceAddr = serviceAddr + fmt.Sprintf("%s%s,", v, vs)
				}
			}
		} else {
			clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.Info.ServicePort)
		}
		upConf, err := model.GetUpstreamConfigFromString(serviceDetail.Info.UpstreamList)
		if err != nil {
			dashboard_middleware.ResponseError(c, 200, err)
			return
		}
		counter, err := handler.ServiceCounterHandler.GetCounter(public.FlowServicePrefix + listItem.ServiceName)
		if err != nil {
			dashboard_middleware.ResponseError(c, 200, err)
			return
		}
		outItem := model.ServiceListItemOutput{
			ID:          listItem.ID,
			LoadType:    listItem.ServiceType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qps:         counter.QPS,
			Qpd:         counter.TotalCount,
			TotalNode:   len(upConf.IpList),
		}
		outList = append(outList, outItem)
	}
	out := &model.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	dashboard_middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &model.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &model.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	if err := serviceInfo.Delete(c, tx, serviceInfo); err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, "")
}

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=model.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &model.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	//读取基本信息
	serviceInfo := &model.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	conf, err := model.GetPluginConfig()
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	// plugin_config的样式信息和值信息组合
	var pluginConfVal []model.PluginConfig
	loadType := ""
	if serviceInfo.ServiceType == 0 {
		loadType = "http"
		pluginConfVal = conf.HTTP
	} else if serviceInfo.ServiceType == 1 {
		loadType = "tcp"
		pluginConfVal = conf.TCP
	} else {
		loadType = "grpc"
		pluginConfVal = conf.GRCP
	}
	for _, v := range pluginConfVal {
		for i := 0; i < len(v.Items); i++ {

			v.Items[i].FieldValue = serviceDetail.PluginConf.GetPath(v.UniqueName, v.Items[i].FieldUniqueName).MustString()

		}
	}
	dashboard_middleware.ResponseSuccess(c, map[string]interface{}{
		loadType:                pluginConfVal,
		"service_name":          serviceInfo.ServiceName,
		"port":                  serviceInfo.ServicePort,
		"service_desc":          serviceInfo.ServiceDesc,
		"http_hosts":            serviceInfo.HTTPHosts,
		"http_paths":            serviceInfo.HTTPPaths,
		"need_strip_uri":        serviceInfo.HttpStripPrefix,
		"load_balance_strategy": serviceInfo.LoadBalanceStrategy,
		"auth_type":             serviceInfo.AuthType,
		"upstream_list":         serviceInfo.UpstreamList,
		"load_balance_type":     serviceInfo.LoadBalanceType,
	})
}

// ServicePluginConfig godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service_plugin_config
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=model.ServiceDetail} "success"
// @Router /service/service_plugin_config [get]
func (service *ServiceController) ServicePluginConfig(c *gin.Context) {
	conf, err := model.GetPluginConfig()
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, conf)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/service_stat [get]
func (service *ServiceController) ServiceStat(c *gin.Context) {
	params := &model.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}

	//读取基本信息
	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &model.ServiceInfo{ID: params.ID}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}

	counter, err := handler.ServiceCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2004, err)
		return
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}
	dashboard_middleware.ResponseSuccess(c, &model.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceAddHTTP godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (service *ServiceController) ServiceAdd(c *gin.Context) {
	params := &model.ServiceAddInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}
	if params.ServiceName == "" {
		dashboard_middleware.ResponseError(c, 2001, errors.New("服务名称不能为空"))
		return
	} else {
		reg, _ := regexp.MatchString(`^[0-9a-zA-Z_]{1,}$`, params.ServiceName)
		if !reg { //解释失败，返回false
			dashboard_middleware.ResponseError(c, 2001, errors.New("服务名称格式错误"))
			return
		}
	}
	if params.ServiceDesc == "" {
		dashboard_middleware.ResponseError(c, 2001, errors.New("服务描述不能为空"))
		return
	}
	if params.LoadType != 1 {
		if params.HTTPHosts == "" {
			dashboard_middleware.ResponseError(c, 2001, errors.New("服务域名不能为空"))
			return
		}
		// else {
		// reg, _ := regexp.MatchString(`^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, params.HTTPHosts)
		// if !reg { //解释失败，返回false
		// 	dashboard_middleware.ResponseError(c, 2001, errors.New("服务域名格式错误"))
		// 	return
		// }
		// }
		if params.HTTPPaths == "" {
			dashboard_middleware.ResponseError(c, 2001, errors.New("服务地址不能为空"))
			return
		} else {
			reg, _ := regexp.MatchString(`^(/[\w\-]+)+`, params.HTTPPaths)
			if !reg { //解释失败，返回false
				dashboard_middleware.ResponseError(c, 2001, errors.New("服务地址格式错误"))
				return
			}
		}
		// if params.NeedStripUri == 0 {
		// 	dashboard_middleware.ResponseError(c, 2001, errors.New("strip_url请选择是否开启"))
		// 	return
		// }
	}
	if params.LoadBalanceStrategy == "" {
		dashboard_middleware.ResponseError(c, 2001, errors.New("loadbalance策略不能为空"))
		return
	}
	if params.UpstreamList == "" {
		dashboard_middleware.ResponseError(c, 2001, errors.New("下游服务器ip和权重不能为空"))
		return
	} else {
		tmpLine := strings.Split(params.UpstreamList, "\n")
		for _, tmp := range tmpLine {
			r, _ := regexp.Compile("^(.*://)(.*?)\\s(.*?)$")
			submatch := r.FindStringSubmatch(tmp)
			if len(submatch) != 4 {
				dashboard_middleware.ResponseError(c, 2001, errors.New("下游服务器ip和权重 format error"))
				return
			}
		}
	}
	// if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
	// 	dashboard_middleware.ResponseError(c, 2004, errors.New("IP列表与权重列表数量不一致"))
	// 	return
	// }

	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	tx = tx.Begin()
	serviceInfo := &model.ServiceInfo{ServiceName: params.ServiceName}
	if _, err = serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	}

	serviceInfo = &model.ServiceInfo{HTTPHosts: params.HTTPHosts, HTTPPaths: params.HTTPPaths}
	if _, err := serviceInfo.Find(c, tx, serviceInfo); err == nil && params.LoadType != 1 {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2003, errors.New("服务路径或域名已存在"))
		return
	}

	serviceModel := &model.ServiceInfo{
		ServiceType:         params.LoadType,
		ServiceName:         params.ServiceName,
		ServiceDesc:         params.ServiceDesc,
		ServicePort:         params.Port,
		HTTPHosts:           params.HTTPHosts,
		HTTPPaths:           params.HTTPPaths,
		HttpStripPrefix:     params.NeedStripUri,
		LoadBalanceStrategy: params.LoadBalanceStrategy,
		LoadBalanceType:     params.LoadBalanceType,
		AuthType:            params.AuthType,
		UpstreamList:        params.UpstreamList,
		PluginConf:          params.PluginConf,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2005, err)
		return
	}
	tx.Commit()
	dashboard_middleware.ResponseSuccess(c, "")
}

// ServiceUpdateHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service_update_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_http [post]
func (service *ServiceController) ServiceUpdate(c *gin.Context) {
	params := &model.ServiceUpdateInput{}
	if err := params.BindValidParam(c); err != nil {
		dashboard_middleware.ResponseError(c, 2000, err)
		return
	}
	// //fmt.Println(public.Obj2Json(params))
	// if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
	// 	dashboard_middleware.ResponseError(c, 2001, errors.New("IP列表与权重列表数量不一致"))
	// 	return
	// }
	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	tx = tx.Begin()
	serviceInfo := &model.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}
	serviceInfo.UpstreamList = params.UpstreamList
	serviceInfo.AuthType = params.AuthType
	serviceInfo.HTTPHosts = params.HTTPHosts
	serviceInfo.HTTPPaths = params.HTTPPaths
	serviceInfo.LoadBalanceStrategy = params.LoadBalanceStrategy
	serviceInfo.LoadBalanceType = params.LoadBalanceType
	serviceInfo.HttpStripPrefix = params.NeedStripUri
	serviceInfo.PluginConf = params.PluginConf
	serviceInfo.ServiceName = params.ServiceName
	serviceInfo.ServiceDesc = params.ServiceDesc
	serviceInfo.ServicePort = params.Port
	if err := serviceInfo.Save(c, tx); err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2005, err)
		return
	}
	//httpRule := serviceDetail.HTTPRule
	//httpRule.NeedHttps = params.NeedHttps
	//httpRule.HTTPStripPrefix = params.HTTPStripPrefix
	//httpRule.NeedWebsocket = params.NeedWebsocket
	//httpRule.UrlRewrite = params.UrlRewrite
	//httpRule.HeaderTransfor = params.HeaderTransfor
	//if err := httpRule.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2006, err)
	//	return
	//}
	//
	//accessControl := serviceDetail.AccessControl
	//accessControl.OpenAuth = params.OpenAuth
	//accessControl.BlackList = params.BlackList
	//accessControl.WhiteList = params.WhiteList
	//accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	//accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	//if err := accessControl.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2007, err)
	//	return
	//}
	//
	//loadbalance := serviceDetail.LoadBalance
	//loadbalance.RoundType = params.RoundType
	//loadbalance.IpList = params.IpList
	//loadbalance.WeightList = params.WeightList
	//loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	//loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	//loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	//loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	//
	//loadbalance.DisfName = params.DisfName
	//loadbalance.DisfClusterName = params.DisfClusterName
	//if err := loadbalance.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2008, err)
	//	return
	//}
	tx.Commit()
	dashboard_middleware.ResponseSuccess(c, "")
}

// ServiceAddHttp godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (admin *ServiceController) ServiceAddTcp(c *gin.Context) {
	params := &model.ServiceAddTcpInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &model.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	if _, err := infoSearch.Find(c, dbPool, infoSearch); err == nil {
		dashboard_middleware.ResponseError(c, 2003, errors.New("服务名被占用，请重新输入"))
		return
	}

	////验证端口是否被占用?
	//tcpRuleSearch := &model.TcpRule{
	//	Port: params.Port,
	//}
	//if _, err := tcpRuleSearch.Find(c, dbPool, tcpRuleSearch); err == nil {
	//	dashboard_middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
	//	return
	//}
	//grpcRuleSearch := &model.GrpcRule{
	//	Port: params.Port,
	//}
	//if _, err := grpcRuleSearch.Find(c, dbPool, grpcRuleSearch); err == nil {
	//	dashboard_middleware.ResponseError(c, 2005, errors.New("服务端口被占用，请重新输入"))
	//	return
	//}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		dashboard_middleware.ResponseError(c, 2006, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := dbPool.Begin()
	info := &model.ServiceInfo{
		ServiceType: public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2007, err)
		return
	}
	//loadBalance := &model.LoadBalance{
	//	ServiceID:  info.ID,
	//	RoundType:  params.RoundType,
	//	IpList:     params.IpList,
	//	WeightList: params.WeightList,
	//	ForbidList: params.ForbidList,
	//}
	//if err := loadBalance.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2008, err)
	//	return
	//}
	//
	//httpRule := &model.TcpRule{
	//	ServiceID: info.ID,
	//	Port:      params.Port,
	//}
	//if err := httpRule.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2009, err)
	//	return
	//}
	//
	//accessControl := &model.AccessControl{
	//	ServiceID:         info.ID,
	//	OpenAuth:          params.OpenAuth,
	//	BlackList:         params.BlackList,
	//	WhiteList:         params.WhiteList,
	//	WhiteHostName:     params.WhiteHostName,
	//	ClientIPFlowLimit: params.ClientIPFlowLimit,
	//	ServiceFlowLimit:  params.ServiceFlowLimit,
	//}
	//if err := accessControl.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2010, err)
	//	return
	//}
	tx.Commit()
	dashboard_middleware.ResponseSuccess(c, "")
	return
}

// ServiceUpdateTcp godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (admin *ServiceController) ServiceUpdateTcp(c *gin.Context) {
	params := &model.ServiceUpdateTcpInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		dashboard_middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	tx := dbPool.Begin()

	service := &model.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(c, dbPool, service)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2004, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2005, err)
		return
	}

	//loadBalance := &model.LoadBalance{}
	//if detail.LoadBalance != nil {
	//	loadBalance = detail.LoadBalance
	//}
	//loadBalance.ServiceID = info.ID
	//loadBalance.RoundType = params.RoundType
	//loadBalance.IpList = params.IpList
	//loadBalance.WeightList = params.WeightList
	//loadBalance.ForbidList = params.ForbidList
	//if err := loadBalance.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2006, err)
	//	return
	//}
	//tcpRule := &model.TcpRule{}
	//if detail.TCPRule != nil {
	//	tcpRule = detail.TCPRule
	//}
	//tcpRule.ServiceID = info.ID
	//tcpRule.Port = params.Port
	//if err := tcpRule.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2005, err)
	//	return
	//}
	//accessControl := &model.AccessControl{}
	//if detail.AccessControl != nil {
	//	accessControl = detail.AccessControl
	//}
	//accessControl.ServiceID = info.ID
	//accessControl.OpenAuth = params.OpenAuth
	//accessControl.BlackList = params.BlackList
	//accessControl.WhiteList = params.WhiteList
	//accessControl.WhiteHostName = params.WhiteHostName
	//accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	//accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	//if err := accessControl.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2007, err)
	//	return
	//}

	tx.Commit()
	dashboard_middleware.ResponseSuccess(c, "")
	return
}

// ServiceAddHttp godoc
// @Summary grpc服务添加
// @Description grpc服务添加
// @Tags 服务管理
// @ID /service/service_add_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_grpc [post]
func (admin *ServiceController) ServiceAddGrpc(c *gin.Context) {
	params := &model.ServiceAddGrpcInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &model.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	if _, err := infoSearch.Find(c, dbPool, infoSearch); err == nil {
		dashboard_middleware.ResponseError(c, 2003, errors.New("服务名被占用，请重新输入"))
		return
	}

	////验证端口是否被占用?
	//tcpRuleSearch := &model.TcpRule{
	//	Port: params.Port,
	//}
	//if _, err := tcpRuleSearch.Find(c, dbPool, tcpRuleSearch); err == nil {
	//	dashboard_middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
	//	return
	//}
	//grpcRuleSearch := &model.GrpcRule{
	//	Port: params.Port,
	//}
	//if _, err := grpcRuleSearch.Find(c, dbPool, grpcRuleSearch); err == nil {
	//	dashboard_middleware.ResponseError(c, 2005, errors.New("服务端口被占用，请重新输入"))
	//	return
	//}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		dashboard_middleware.ResponseError(c, 2006, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := dbPool.Begin()
	info := &model.ServiceInfo{
		ServiceType: public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2007, err)
		return
	}

	//loadBalance := &model.LoadBalance{
	//	ServiceID:  info.ID,
	//	RoundType:  params.RoundType,
	//	IpList:     params.IpList,
	//	WeightList: params.WeightList,
	//	ForbidList: params.ForbidList,
	//}
	//if err := loadBalance.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2008, err)
	//	return
	//}
	//
	//grpcRule := &model.GrpcRule{
	//	ServiceID:      info.ID,
	//	Port:           params.Port,
	//	HeaderTransfor: params.HeaderTransfor,
	//}
	//if err := grpcRule.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2009, err)
	//	return
	//}
	//
	//accessControl := &model.AccessControl{
	//	ServiceID:         info.ID,
	//	OpenAuth:          params.OpenAuth,
	//	BlackList:         params.BlackList,
	//	WhiteList:         params.WhiteList,
	//	WhiteHostName:     params.WhiteHostName,
	//	ClientIPFlowLimit: params.ClientIPFlowLimit,
	//	ServiceFlowLimit:  params.ServiceFlowLimit,
	//}
	//if err := accessControl.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2010, err)
	//	return
	//}
	tx.Commit()
	dashboard_middleware.ResponseSuccess(c, "")
	return
}

// ServiceUpdateTcp godoc
// @Summary grpc服务更新
// @Description grpc服务更新
// @Tags 服务管理
// @ID /service/service_update_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_grpc [post]
func (admin *ServiceController) ServiceUpdateGrpc(c *gin.Context) {
	params := &model.ServiceUpdateGrpcInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		dashboard_middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	tx := dbPool.Begin()

	service := &model.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(c, dbPool, service)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2004, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		dashboard_middleware.ResponseError(c, 2005, err)
		return
	}

	//loadBalance := &model.LoadBalance{}
	//if detail.LoadBalance != nil {
	//	loadBalance = detail.LoadBalance
	//}
	//loadBalance.ServiceID = info.ID
	//loadBalance.RoundType = params.RoundType
	//loadBalance.IpList = params.IpList
	//loadBalance.WeightList = params.WeightList
	//loadBalance.ForbidList = params.ForbidList
	//if err := loadBalance.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2006, err)
	//	return
	//}
	//
	//grpcRule := &model.GrpcRule{}
	//if detail.GRPCRule != nil {
	//	grpcRule = detail.GRPCRule
	//}
	//grpcRule.ServiceID = info.ID
	////grpcRule.Port = params.Port
	//grpcRule.HeaderTransfor = params.HeaderTransfor
	//if err := grpcRule.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2007, err)
	//	return
	//}
	//
	//accessControl := &model.AccessControl{}
	//if detail.AccessControl != nil {
	//	accessControl = detail.AccessControl
	//}
	//accessControl.ServiceID = info.ID
	//accessControl.OpenAuth = params.OpenAuth
	//accessControl.BlackList = params.BlackList
	//accessControl.WhiteList = params.WhiteList
	//accessControl.WhiteHostName = params.WhiteHostName
	//accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	//accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	//if err := accessControl.Save(c, tx); err != nil {
	//	tx.Rollback()
	//	dashboard_middleware.ResponseError(c, 2008, err)
	//	return
	//}

	tx.Commit()
	dashboard_middleware.ResponseSuccess(c, "")
	return
}
