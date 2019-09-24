package dao

import (
	"github.com/didichuxing/gatekeeper/public"
	"github.com/e421083458/gorm"
)

//APPConfiger app列表配置
type APPConfiger struct {
	AuthirizedApps []*GatewayAPP `json:"authirized_apps" toml:"authirized_apps"  validate:"required"`
}

//GatewayAPP app配置
type GatewayAPP struct {
	ID              int64  `json:"id" toml:"-" orm:"column(id);auto" description:"自增主键"`
	AppID           string `json:"app_id" toml:"app_id" validate:"required"  orm:"column(app_id);" description:"租户id"`
	Name            string `json:"name" toml:"name" validate:"required"  orm:"column(name);" description:"租户名称"`
	Secret          string `json:"secret" toml:"secret" validate:"required"  orm:"column(secret);" description:"密钥"`
	Method          string `json:"method" toml:"method" validate:""  orm:"column(method);" descriId           int6ption:"请求方法"`
	Timeout         int64  `json:"timeout" toml:"timeout" orm:"column(timeout);" description:"超时时间"`
	OpenAPI         string `json:"open_api" toml:"open_api" orm:"column(open_api);" description:"接口列表，支持前缀匹配"`
	WhiteIps        string `json:"white_ips" toml:"white_ips" orm:"column(white_ips);" description:"ip白名单，支持前缀匹配"`
	CityIDs         string `json:"city_ids" toml:"city_ids" orm:"column(city_ids);" description:"city_id数据权限"`
	TotalQueryDaily int64  `json:"total_query_daily" toml:"total_query_daily" orm:"column(total_query_daily);" description:"日请求量"`
	QPS             int64  `json:"qps" toml:"qps" orm:"column(qps);" description:"qps"`
	GroupID         int64  `json:"group_id" toml:"group_id" orm:"column(group_id);" description:"数据关联id"`
}

//TableName db table
func (o *GatewayAPP) TableName() string {
	return "gateway_app"
}

//GetAll db getall
func (o *GatewayAPP) GetAll(orders... string) ([]*GatewayAPP, error) {
	var apps []*GatewayAPP
	builder := public.GormPool
	for _,order:=range orders{
		builder = builder.Order(order)
	}
	err := builder.Find(&apps).Error
	return apps, err
}

//FindByAppID db findbyappid
func (o *GatewayAPP) FindByAppID(db *gorm.DB, appID string) (*GatewayAPP, error) {
	var app GatewayAPP
	err := db.Where("app_id = ?", appID).First(&app).Error
	if err == gorm.ErrRecordNotFound {
		return &app, nil
	}
	return &app, err
}

//FindByID db findbyid
func (o *GatewayAPP) FindByID(db *gorm.DB, id int64) (*GatewayAPP, error) {
	var app GatewayAPP
	err := db.Where("id = ?", id).First(&app).Error
	if err == gorm.ErrRecordNotFound {
		return &app, nil
	}
	return &app, err
}

//GetPk get pk
func (o *GatewayAPP) GetPk() int64 {
	return o.ID
}

//Save db save
func (o *GatewayAPP) Save(db *gorm.DB) error {
	return db.Save(o).Error
}

//Del db del
func (o *GatewayAPP) Del(db *gorm.DB) error {
	if err := db.Where("id = ?", o.ID).Delete(o).Error; err != nil {
		return err
	}
	return nil
}
