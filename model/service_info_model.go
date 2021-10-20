package model

import (
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/didi/gatekeeper/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceInfo struct {
	ID                  int64     `json:"id" gorm:"primary_key"`
	LoadType            int       `json:"service_type" gorm:"column:service_type" description:"服务类型 0=http 1=tcp 2=grpc"`
	ServiceName         string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc         string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	Port                int       `json:"service_port" gorm:"column:service_port" description:"服务端口(针对 tcp/grpc)"`
	HTTPHosts           string    `json:"http_hosts" gorm:"column:http_hosts" description:"http域名信息"`
	HTTPPaths           string `json:"http_paths" gorm:"column:http_paths" description:"http路径信息"`
	HTTPStripPrefix     int    `json:"http_strip_prefix" gorm:"column:http_strip_prefix" description:"http转发前剥离前缀"`
	LoadBalanceStrategy string `json:"load_balance_strategy" gorm:"column:load_balance_strategy" description:"负载策略"`
	LoadBalanceType     string    `json:"load_balance_type" gorm:"column:load_balance_type" description:"负载类型"`
	AuthType            string    `json:"auth_type" gorm:"column:auth_type" description:"鉴权类型"`
	UpstreamList        string    `json:"upstream_list" form:"upstream_list" comment:"下游服务器ip和权重"`
	PluginConf          string    `json:"plugin_conf" gorm:"column:plugin_conf" description:"插件配置"`
	CreatedAt           time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	UpdatedAt           time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete            int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) Delete(c *gin.Context, tx *gorm.DB, search *ServiceInfo) error {
	if err := tx.Where("id=?", search.ID).Delete(&ServiceInfo{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, info *ServiceInfo) (*ServiceDetail, error) {
	if info.ServiceName == "" {
		info, err := t.Find(c, tx, info)
		if err != nil {
			return nil, err
		}
		info = info
	}
	pluginConf := simplejson.New()
	if tmp, err := simplejson.NewJson([]byte(info.PluginConf)); err == nil {
		pluginConf = tmp
	}
	detail := &ServiceDetail{
		Info:       info,
		PluginConf: pluginConf,
	}
	return detail, nil
}

func (t *ServiceInfo) GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]DashServiceStatItemOutput, error) {
	list := []DashServiceStatItemOutput{}
	query := tx.SetCtx(public.GetGinTraceContext(c))
	if err := query.Table(t.TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (t *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, param *ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Count(&total)
	return list, total, nil
}

func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	query := tx.SetCtx(public.GetGinTraceContext(c))
	if search.HTTPHosts != "" {
		query = query.Where("http_hosts LIKE ? ", "%"+search.HTTPHosts+"%")
	}
	if search.HTTPPaths != "" {
		query = query.Where("http_hosts LIKE ? ", "%"+search.HTTPPaths+"%")
	}
	query = query.Where("is_delete=0").Where(search).Find(out)
	if err := query.Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
