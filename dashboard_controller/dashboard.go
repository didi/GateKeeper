package dashboard_controller

import (
	"github.com/didi/gatekeeper/dashboard_middleware"
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/handler"
	"github.com/didi/gatekeeper/model"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	service := &DashboardController{}
	group.GET("/panel_group_data", service.PanelGroupData)
	group.GET("/flow_stat", service.FlowStat)
	group.GET("/service_stat", service.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (service *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &model.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &model.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	app := &model.App{}
	_, appNum, err := app.APPList(c, tx, &model.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	counter, err := handler.ServiceCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	out := &model.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: counter.TotalCount,
		CurrentQPS:      counter.QPS,
	}
	dashboard_middleware.ResponseSuccess(c, out)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (service *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &model.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			dashboard_middleware.ResponseError(c, 2003, errors.New("service_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &model.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	dashboard_middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (service *DashboardController) FlowStat(c *gin.Context) {
	counter, err := handler.ServiceCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
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
