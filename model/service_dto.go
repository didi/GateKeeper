package model

import (
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
)

func (param *ServiceUpdateInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateInput struct {
	ID                  int64  `json:"id" form:"id" comment:"服务ID" example:"62" validate:"required,min=1"`                                //服务ID
	ServiceName         string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` //服务名称
	ServiceDesc         string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`      //服务描述
	LoadType            int    `json:"load_type" form:"load_type" comment:"负载类型 0=http 1=tcp 2=grpc"`
	HTTPHosts           string `json:"http_hosts" form:"http_hosts" comment:"域名信息" example:""`
	HTTPPaths           string `json:"http_paths" form:"http_paths" comment:"路径信息" example:""`
	NeedStripUri        string `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"` //启用strip_uri
	Port                int    `json:"port" form:"port" comment:"tcp/grpc端口"`
	LoadBalanceStrategy string `json:"load_balance_strategy" form:"load_balance_strategy" comment:"负载策略"`
	LoadBalanceType     string `json:"load_balance_type" form:"load_balance_type" comment:"负载类型"`
	AuthType            string `json:"auth_type" form:"auth_type" comment:"鉴权类型"`
	UpstreamList        string `json:"upstream_list" form:"upstream_list" comment:"下游服务器ip和权重"`
	PluginConf          string `json:"plugin_conf" form:"plugin_conf" comment:"插件配置"`
}

func (param *ServiceAddInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddInput struct {
	ServiceName         string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` //服务名称
	ServiceDesc         string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`      //服务描述
	LoadType            int    `json:"load_type" form:"load_type" comment:"负载类型 0=http 1=tcp 2=grpc"`
	HTTPHosts           string `json:"http_hosts" form:"http_hosts" comment:"域名信息" example:""`
	HTTPPaths           string `json:"http_paths" form:"http_paths" comment:"路径信息" example:""`
	NeedStripUri        string `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"` //启用strip_uri
	Port                int    `json:"port" form:"port" comment:"tcp/grpc端口"`
	LoadBalanceStrategy string `json:"load_balance_strategy" form:"load_balance_strategy" comment:"负载策略"`
	LoadBalanceType     string `json:"load_balance_type" form:"load_balance_type" comment:"负载类型"`
	AuthType            string `json:"auth_type" form:"auth_type" comment:"鉴权类型"`
	UpstreamList        string `json:"upstream_list" form:"upstream_list" comment:"下游服务器ip和权重"`
	PluginConf          string `json:"plugin_conf" form:"plugin_conf" comment:"插件配置"`
}

type ServiceDeleteInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` //服务ID
}

func (param *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceDetailInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` //服务ID
}

func (param *ServiceDetailInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` //每页条数
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     //id
	ServiceName string `json:"service_name" form:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` //服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       //类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` //服务地址
	Qps         int64  `json:"qps" form:"qps"`                   //qps
	Qpd         int64  `json:"qpd" form:"qpd"`                   //qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     //节点数
}

type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

type ServiceStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" example:"" validate:""`         //列表
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" example:"" validate:""` //列表
}

type ServiceAddGrpcInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceAddGrpcInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceUpdateGrpcInput struct {
	ID                int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceUpdateGrpcInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceAddTcpInput struct {
	ServiceName    string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc    string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header头转换" validate:"
"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceAddTcpInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceUpdateTcpInput struct {
	ID                int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceUpdateTcpInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
