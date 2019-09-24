package dao

import (
	"github.com/didichuxing/gatekeeper/public"
	"github.com/e421083458/gorm"
)

//ModuleConfiger module整体配置结构
type ModuleConfiger struct {
	Module []*GatewayModule `json:"module" toml:"module" validate:"required"`
}

//GetGateWayModule 根据名字获取模块
func (c *ModuleConfiger) GetGateWayModule(name string) *GatewayModule{
	for _, module := range c.Module {
		if module.Base.Name == name {
			return module
		}
	}
	return nil
}

//GatewayModule module配置
type GatewayModule struct {
	Base          *GatewayModuleBase    `json:"base" validate:"required" toml:"base"`
	MatchRule     []*GatewayMatchRule   `json:"match_rule" validate:"required"  toml:"match_rule"`
	LoadBalance   *GatewayLoadBalance   `json:"load_balance" validate:"required" toml:"load_balance"`
	AccessControl *GatewayAccessControl `json:"access_control" toml:"access_control"`
}

//GatewayModuleBase base数据表结构体
type GatewayModuleBase struct {
	ID           int64  `json:"id" toml:"-" orm:"column(id);auto" description:"自增主键"`
	LoadType     string `json:"load_type" toml:"load_type" validate:"" orm:"column(load_type);size(255)" description:"负载类型 http/tcp"`
	Name         string `json:"name" toml:"name" validate:"required" orm:"column(name);size(255)" description:"模块名"`
	ServiceName  string `json:"service_name" toml:"service_name" validate:"" orm:"column(service_name);size(255)" description:"服务名称"`
	PassAuthType int8   `json:"pass_auth_type" toml:"pass_auth_type" validate:"" orm:"column(pass_auth_type)" description:"认证传参类型"`
	FrontendAddr string `json:"frontend_addr" toml:"frontend_addr" validate:"" orm:"column(frontend_addr);size(255)" description:"前端绑定ip地址"`
}

//TableName db tablename
func (o *GatewayModuleBase) TableName() string {
	return "gateway_module_base"
}

//GetAll db getall
func (o *GatewayModuleBase) GetAll(orders... string) ([]*GatewayModuleBase, error) {
	var modules []*GatewayModuleBase
	builder := public.GormPool
	for _,order:=range orders{
		builder = builder.Order(order)
	}
	err := builder.Find(&modules).Error
	return modules, err
}

//FindByName db findbyname
func (o *GatewayModuleBase) FindByName(db *gorm.DB, name string) (*GatewayModuleBase, error) {
	var modules GatewayModuleBase
	err := db.Where("name = ?", name).First(&modules).Error
	if err == gorm.ErrRecordNotFound {
		return &modules, nil
	}
	return &modules, err
}

//FindByPort db findbyport
func (o *GatewayModuleBase) FindByPort(db *gorm.DB, port string) (*GatewayModuleBase, error) {
	var modules GatewayModuleBase
	err := db.Where("frontend_addr = ?", port).First(&modules).Error
	if err == gorm.ErrRecordNotFound {
		return &modules, nil
	}
	return &modules, err
}

//GetPk getpk
func (o *GatewayModuleBase) GetPk() int64 {
	return o.ID
}

//Save db row
func (o *GatewayModuleBase) Save(db *gorm.DB) error {
	return db.Save(o).Error
}

//Del db row
func (o *GatewayModuleBase) Del(db *gorm.DB) error {
	if err := db.Where("id = ?", o.ID).Delete(o).Error; err != nil {
		return err
	}
	return nil
}

//GatewayMatchRule match_rule数据表结构体
type GatewayMatchRule struct {
	ID         int64  `json:"id" toml:"-" orm:"column(id);auto" description:"自增主键"`
	ModuleID   int64  `json:"module_id" toml:"-" orm:"column(module_id)" description:"模块id"`
	Type       string `json:"type" toml:"type" validate:"required" orm:"column(type)" description:"匹配类型"`
	Rule       string `json:"rule" toml:"rule" validate:"required" orm:"column(rule);size(1000)" description:"规则"`
	RuleExt    string `json:"rule_ext" validate:"required" toml:"rule_ext" orm:"column(rule_ext);size(1000)" description:"拓展规则"`
	URLRewrite string `json:"url_rewrite" validate:"required" toml:"url_rewrite" orm:"column(rule_ext);size(1000)" description:"url重写"`
}

//TableName tablename
func (o *GatewayMatchRule) TableName() string {
	return "gateway_match_rule"
}

//GetAll db getall
func (o *GatewayMatchRule) GetAll() ([]*GatewayMatchRule, error) {
	var rules []*GatewayMatchRule
	err := public.GormPool.
		Find(&rules).Error
	return rules, err
}

//GetByModule db getmodule
func (o *GatewayMatchRule) GetByModule(moduleID int64) ([]*GatewayMatchRule, error) {
	var rules []*GatewayMatchRule
	err := public.GormPool.
		Where(&GatewayMatchRule{ModuleID: moduleID}).
		Find(&rules).Error
	return rules, err
}

//GetPk db get pk
func (o *GatewayMatchRule) GetPk() int64 {
	return o.ID
}

//Save db save
func (o *GatewayMatchRule) Save(db *gorm.DB) error {
	return db.Save(o).Error
}

//FindByURLPrefix db findbyprefix
func (o *GatewayMatchRule) FindByURLPrefix(db *gorm.DB, prefix string) (*GatewayMatchRule, error) {
	var rule GatewayMatchRule
	err := db.Where("rule = ?", prefix).First(&rule).Error
	if err == gorm.ErrRecordNotFound {
		return &rule, nil
	}
	return &rule, err
}

//Del db del
func (o *GatewayMatchRule) Del(db *gorm.DB) error {
	if err := db.Where("module_id = ?", o.ModuleID).Delete(o).Error; err != nil {
		return err
	}
	return nil
}

//GatewayLoadBalance load_balance数据表结构体
type GatewayLoadBalance struct {
	ID            int64  `json:"id" toml:"-" orm:"column(id);auto" description:"自增主键"`
	ModuleID      int64  `json:"module_id" toml:"-" orm:"column(module_id)"`
	CheckMethod   string `json:"check_method" validate:"required" toml:"check_method" orm:"column(check_method);size(200)" description:"检查方法"`
	CheckURL      string `json:"check_url" validate:"" toml:"check_url" orm:"column(check_url);size(500)" description:"检测url"`
	CheckTimeout  int    `json:"check_timeout" validate:"required,min=100" toml:"check_timeout" orm:"column(check_timeout);size(500)" description:"检测超时时间"`
	CheckInterval int    `json:"check_interval" validate:"required,min=100" toml:"check_interval" orm:"column(check_interval);size(500)" description:"检测url"`

	Type                string `json:"type" validate:"required" toml:"type" orm:"column(type);size(100)" description:"轮询方式"`
	IPList              string `json:"ip_list" validate:"required" toml:"ip_list" orm:"column(ip_list);size(500)" description:"ip列表"`
	WeightList          string `json:"weight_list" validate:"" toml:"weight_list" orm:"column(weight_list);size(500)" description:"ip列表"`
	ForbidList          string `json:"forbid_list" validate:"" toml:"forbid_list" orm:"column(forbid_list);size(1000)" description:"禁用 ip列表"`
	ProxyConnectTimeout int    `json:"proxy_connect_timeout" validate:"required,min=1" toml:"proxy_connect_timeout" orm:"column(proxy_connect_timeout)" description:"单位ms，连接后端超时时间"`
	ProxyHeaderTimeout  int    `json:"proxy_header_timeout" validate:"" toml:"proxy_header_timeout" orm:"column(proxy_header_timeout)" description:"单位ms，后端服务器数据回传时间"`
	ProxyBodyTimeout    int    `json:"proxy_body_timeout" validate:"" toml:"proxy_body_timeout" orm:"column(proxy_body_timeout)" description:"单位ms，后端服务器响应时间"`
	MaxIdleConn         int    `json:"max_idle_conn" validate:"" toml:"max_idle_conn" orm:"column(max_idle_conn)"`
	IdleConnTimeout     int    `json:"idle_conn_timeout" validate:"" toml:"idle_conn_timeout" orm:"column(idle_conn_timeout)" description:"keep-alived超时时间，新增"`
}

//TableName tablename
func (o *GatewayLoadBalance) TableName() string {
	return "gateway_load_balance"
}

//GetAll db getall
func (o *GatewayLoadBalance) GetAll() ([]*GatewayLoadBalance, error) {
	var rules []*GatewayLoadBalance
	err := public.GormPool.Model(&GatewayLoadBalance{}).
		Find(&rules).Error
	return rules, err
}

//GetByModule db getbymoduleid
func (o *GatewayLoadBalance) GetByModule(moduleID int64) (*GatewayLoadBalance, error) {
	var rules []*GatewayLoadBalance
	err := public.GormPool.Model(&GatewayLoadBalance{}).
		Where(&GatewayLoadBalance{ModuleID: moduleID}).
		Find(&rules).Error
	if len(rules) == 0 {
		return nil, nil
	}
	return rules[0], err
}

//GetPk getpk
func (o *GatewayLoadBalance) GetPk() int64 {
	return o.ID
}

//Save db save
func (o *GatewayLoadBalance) Save(db *gorm.DB) error {
	return db.Save(o).Error
}

//Del db del
func (o *GatewayLoadBalance) Del(db *gorm.DB) error {
	if err := db.Where("module_id = ?", o.ModuleID).Delete(o).Error; err != nil {
		return err
	}
	return nil
}

//GatewayAccessControl access_control 数据表结构体
type GatewayAccessControl struct {
	ID              int64  `json:"id" toml:"-" orm:"column(id);auto" description:"自增主键"`
	ModuleID        int64  `json:"module_id" toml:"-" orm:"column(module_id)" description:"模块id"`
	BlackList       string `json:"black_list" toml:"black_list" orm:"column(black_list);size(1000)" description:"黑名单ip"`
	WhiteList       string `json:"white_list" toml:"white_list" orm:"column(white_list);size(1000)" description:"白名单ip"`
	WhiteHostName   string `json:"white_host_name" toml:"white_host_name" orm:"column(white_host_name);size(1000)" description:"白名单主机"`
	AuthType        string `json:"auth_type" toml:"auth_type" orm:"column(auth_type);size(100)" description:"认证方法"`
	ClientFlowLimit int64  `json:"client_flow_limit" toml:"client_flow_limit" orm:"column(client_flow_limit);size(100)" description:"客户端ip限流"`
	Open            int64  `json:"open" toml:"open" orm:"column(open);size(100)" description:"是否开启权限功能"`
}

//TableName tablename
func (o *GatewayAccessControl) TableName() string {
	return "gateway_access_control"
}

//GetByModule getbymoduleid
func (o *GatewayAccessControl) GetByModule(moduleID int64) (*GatewayAccessControl, error) {
	var rules []*GatewayAccessControl
	err := public.GormPool.Model(&GatewayAccessControl{}).
		Where(&GatewayAccessControl{ModuleID: moduleID}).
		Find(&rules).Error
	if len(rules) == 0 {
		return nil, nil
	}
	return rules[0], err
}

//GetAll getall
func (o *GatewayAccessControl) GetAll() ([]*GatewayAccessControl, error) {
	var rules []*GatewayAccessControl
	err := public.GormPool.Model(&GatewayAccessControl{}).
		Find(&rules).Error
	return rules, err
}

//GetPk getpk
func (o *GatewayAccessControl) GetPk() int64 {
	return o.ID
}

//Save db save
func (o *GatewayAccessControl) Save(db *gorm.DB) error {
	return db.Save(o).Error
}

//Del db del
func (o *GatewayAccessControl) Del(db *gorm.DB) error {
	if err := db.Where("module_id = ?", o.ModuleID).Delete(o).Error; err != nil {
		return err
	}
	return nil
}
