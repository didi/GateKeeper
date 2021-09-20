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

//APPControllerRegister admin路由注册
func APPRegister(router *gin.RouterGroup) {
	admin := APPController{}
	router.GET("/app_list", admin.APPList)
	router.GET("/app_detail", admin.APPDetail)
	router.GET("/app_stat", admin.AppStatistics)
	router.GET("/app_delete", admin.APPDelete)
	router.POST("/app_add", admin.AppAdd)
	router.POST("/app_update", admin.AppUpdate)
}

type APPController struct {
}

// APPList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
// @Router /app/app_list [get]
func (admin *APPController) APPList(c *gin.Context) {
	params := &model.APPListInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	info := &model.App{}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	list, total, err := info.APPList(c, dbPool, params)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}

	outputList := []model.APPListItemOutput{}
	for _, item := range list {
		appCounter, err := handler.ServiceCounterHandler.GetCounter(public.FlowAppPrefix + item.AppID)
		if err != nil {
			dashboard_middleware.ResponseError(c, 2004, err)
			c.Abort()
			return
		}
		outputList = append(outputList, model.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  appCounter.TotalCount,
			RealQps:  appCounter.QPS,
		})
	}
	output := model.APPListOutput{
		List:  outputList,
		Total: total,
	}
	dashboard_middleware.ResponseSuccess(c, output)
	return
}

// APPDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=model.App} "success"
// @Router /app/app_detail [get]
func (admin *APPController) APPDetail(c *gin.Context) {
	params := &model.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	search := &model.App{
		ID: params.ID,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	detail, err := search.Find(c, dbPool, search)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, detail)
	return
}

// APPDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (admin *APPController) APPDelete(c *gin.Context) {
	params := &model.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	search := &model.App{
		ID: params.ID,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	info, err := search.Find(c, dbPool, search)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	info.IsDelete = 1
	if err := info.Save(c, dbPool); err != nil {
		dashboard_middleware.ResponseError(c, 2004, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, "")
	return
}

// AppAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (admin *APPController) AppAdd(c *gin.Context) {
	params := &model.APPAddHttpInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	//验证app_id是否被占用
	search := &model.App{
		AppID: params.AppID,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	if _, err := search.Find(c, dbPool, search); err == nil {
		dashboard_middleware.ResponseError(c, 2003, errors.New("租户ID被占用，请重新输入"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	tx := dbPool
	info := &model.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := info.Save(c, tx); err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, "")
	return
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [post]
func (admin *APPController) AppUpdate(c *gin.Context) {
	params := &model.APPUpdateHttpInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}
	search := &model.App{
		ID: params.ID,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	info, err := search.Find(c, dbPool, search)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.Qps = params.Qps
	info.Qpd = params.Qpd
	if err := info.Save(c, dbPool); err != nil {
		dashboard_middleware.ResponseError(c, 2004, err)
		return
	}
	dashboard_middleware.ResponseSuccess(c, "")
	return
}

// AppStatistics godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dto.StatisticsOutput} "success"
// @Router /app/app_stat [get]
func (admin *APPController) AppStatistics(c *gin.Context) {
	params := &model.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		dashboard_middleware.ResponseError(c, 2001, err)
		return
	}

	search := &model.App{
		ID: params.ID,
	}
	dbPool, err := lib.GetGormPool("default")
	if err != nil {
		dashboard_middleware.ResponseError(c, 2002, err)
		return
	}
	detail, err := search.Find(c, dbPool, search)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2003, err)
		return
	}

	//今日流量全天小时级访问统计
	todayStat := []int64{}
	counter, err := handler.ServiceCounterHandler.GetCounter(public.FlowAppPrefix + detail.AppID)
	if err != nil {
		dashboard_middleware.ResponseError(c, 2004, err)
		c.Abort()
		return
	}
	currentTime := time.Now()
	for i := 0; i <= time.Now().In(lib.TimeLocation).Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayStat = append(todayStat, hourData)
	}

	//昨日流量全天小时级访问统计
	yesterdayStat := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayStat = append(yesterdayStat, hourData)
	}
	stat := model.StatisticsOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	}
	dashboard_middleware.ResponseSuccess(c, stat)
	return
}
