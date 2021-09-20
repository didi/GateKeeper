package model

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/gin-gonic/gin"
)

const (
	PLUGIN_POSITION_HIDDEN      PluginPosition = "hidden"      //隐藏形式
	PLUGIN_POSITION_NORMAL      PluginPosition = "normal"      //文档流动形式向下展示
	PLUGIN_POSITION_AUTH        PluginPosition = "auth"        //下拉选择认证类型
	PLUGIN_POSITION_LOADBALANCE PluginPosition = "loadbalance" //下拉选择负载类型

	PLUGIN_FIELD_TYPE_INPUT    PluginFieldType = "input"    //文字
	PLUGIN_FIELD_TYPE_TEXTAREA PluginFieldType = "textarea" //文本域
	PLUGIN_FIELD_TYPE_RADIO    PluginFieldType = "radio"    //单选
	PLUGIN_FIELD_TYPE_CHECKBOX PluginFieldType = "checkbox" //多选
	PLUGIN_FIELD_TYPE_SELECT   PluginFieldType = "select"   //下拉
	PLUGIN_FIELD_TYPE_SWITCH   PluginFieldType = "switch"   //开关

	PLUGIN_FIELD_DISPLAY_INLINE PluginFieldDisplay = "inline" //内联
	PLUGIN_FIELD_DISPLAY_BLOCK  PluginFieldDisplay = "block"  //块级

	PLUGIN_FIELD_CLEAR_NONE  PluginFieldClear = "none"  //默认
	PLUGIN_FIELD_CLEAR_LEFT  PluginFieldClear = "left"  //左边
	PLUGIN_FIELD_CLEAR_RIGHT PluginFieldClear = "right" //右边
	PLUGIN_FIELD_CLEAR_BOTH  PluginFieldClear = "both"  //两边
)

type PluginPosition string
type PluginFieldType string
type PluginFieldDisplay string
type PluginFieldClear string

type HTTPProxyMiddleware interface {
	Configration()
	GetConfig() *PluginConfig
	Handler(c *gin.Context)
}

type PluginConfigs struct {
	HTTP []PluginConfig `toml:"http" json:"http"`
	TCP  []PluginConfig `toml:"tcp" json:"tcp"`
	GRCP []PluginConfig `toml:"grpc" json:"grpc"`
}

type PluginConfig struct {
	DisplayName string             `toml:"display_name" json:"display_name"` //插件名
	Sort        int                `toml:"sort" json:"sort"`                 //显示顺序
	Postion     PluginPosition     `toml:"postion" json:"postion"`           //显示位置
	UniqueName  string             `toml:"unique_name" json:"unique_name"`   //唯一名称
	Items       []PluginItemConfig `toml:"items" json:"items"`
}

type PluginItemConfig struct {
	FieldType         PluginFieldType    `toml:"field_type" json:"field_type"`                   //表单类型 文字=input,文本域=textarea,单选=radio,多选=checkbox,下拉=select,开关=switch
	FieldDisplay      PluginFieldDisplay `toml:"field_display" json:"field_display"`             //显示样式 内联=inline,块=black
	FieldClear        PluginFieldClear   `toml:"field_clear" json:"field_clear"`                 //显示清空 清空左=left
	FieldPlaceholder  string             `toml:"field_placeholder" json:"field_placeholder"`     //占位符
	FieldOption       string             `toml:"field_option" json:"field_option"`               //可选值，radio、check、select时的value
	FieldValue        string             `toml:"field_value" json:"field_value"`                 //字段值
	FieldDefaultValue string             `toml:"field_default_value" json:"field_default_value"` //字段默认值，checkbox时 111|222 表示选中多个值
	FieldUniqueName   string             `toml:"field_unique_name" json:"field_unique_name"`     //字段唯一标识
	FieldDisplayName  string             `toml:"field_display_name" json:"field_display_name"`   //字段显示名称，radio、check、select时的name
	FieldRequired     bool               `toml:"field_required" json:"field_required"`           //是否必填
	FieldValidRule    string             `toml:"field_valid_rule" json:"field_valid_rule"`       //验证规则
}

func GetPluginConfig() (*PluginConfigs, error) {
	pluginConf := &PluginConfigs{}
	err := lib.ParseConfig(lib.GetConfPath("plugin_config"), pluginConf)
	if err != nil {
		return nil, err
	}
	return pluginConf, nil
	//
	////yaml.Unmarshal()
	//return &PluginConfigs{
	//	HTTP: []PluginConfig{
	//		PluginConfig{
	//			DisplayName: "url地址重写",
	//			Sort:        6,
	//			Postion:     PLUGIN_POSITION_NORMAL,
	//			UniqueName:  "url_rewrite",
	//			Items: []PluginItemConfig{
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_TEXTAREA,
	//					FieldPlaceholder:  "格式：^/gateway/test_service(.*) $1 多条换行",
	//					FieldDefaultValue: "",
	//					FieldUniqueName:   "url_rewrite",
	//					FieldDisplayName:  "URL重写",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[\\S]+$/",
	//					FieldDisplay:      PLUGIN_FIELD_DISPLAY_BLOCK,
	//					FieldClear:        PLUGIN_FIELD_CLEAR_NONE,
	//				},
	//			},
	//		},
	//		PluginConfig{
	//			DisplayName: "限流中间件",
	//			Sort:        3,
	//			Postion:     PLUGIN_POSITION_NORMAL,
	//			UniqueName:  "http_flow_limit",
	//			Items: []PluginItemConfig{
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_INPUT,
	//					FieldOption:       "",
	//					FieldDefaultValue: "0",
	//					FieldUniqueName:   "service_flow_limit_num",
	//					FieldDisplayName:  "服务限流数",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[0-9]$/",
	//					FieldDisplay:      PLUGIN_FIELD_DISPLAY_INLINE,
	//					FieldClear:        PLUGIN_FIELD_CLEAR_NONE,
	//				},
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_SELECT,
	//					FieldOption:       "0|1|2,秒|分|时",
	//					FieldDefaultValue: "0",
	//					FieldUniqueName:   "service_flow_limit_type",
	//					FieldDisplayName:  "服务限流类型",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[0-9]$/",
	//					FieldDisplay:      PLUGIN_FIELD_DISPLAY_INLINE,
	//					FieldClear:        PLUGIN_FIELD_CLEAR_RIGHT,
	//				},
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_INPUT,
	//					FieldOption:       "",
	//					FieldDefaultValue: "0",
	//					FieldUniqueName:   "clientip_flow_limit_num",
	//					FieldDisplayName:  "客户IP限流数",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[0-9]$/",
	//					FieldDisplay:      PLUGIN_FIELD_DISPLAY_INLINE,
	//					FieldClear:        PLUGIN_FIELD_CLEAR_NONE,
	//				},
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_SELECT,
	//					FieldOption:       "0|1|2,秒|分|时",
	//					FieldDefaultValue: "0",
	//					FieldUniqueName:   "clientip_flow_limit_type",
	//					FieldDisplayName:  "客户IP限流类型",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[0-9]$/",
	//					FieldDisplay:      PLUGIN_FIELD_DISPLAY_INLINE,
	//					FieldClear:        PLUGIN_FIELD_CLEAR_RIGHT,
	//				},
	//			},
	//		},
	//		PluginConfig{
	//			DisplayName: "Header头转换",
	//			Sort:        3,
	//			Postion:     PLUGIN_POSITION_NORMAL,
	//			UniqueName:  "header_transfer_middleware",
	//			Items: []PluginItemConfig{
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_TEXTAREA,
	//					FieldPlaceholder:  "header转换支持 add(增加)/del(删除)/edit(修改) 格式：add headerName headValue",
	//					FieldDefaultValue: "",
	//					FieldUniqueName:   "service_flow_limit_num",
	//					FieldDisplayName:  "服务限流数",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[\\S]+$/",
	//					FieldDisplay:      PLUGIN_FIELD_DISPLAY_INLINE,
	//					FieldClear:        PLUGIN_FIELD_CLEAR_NONE,
	//				},
	//			},
	//		},
	//		PluginConfig{
	//			DisplayName: "白名单&黑名单",
	//			Sort:        2,
	//			Postion:     PLUGIN_POSITION_NORMAL,
	//			UniqueName:  "http_blacklist",
	//			Items: []PluginItemConfig{
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_TEXTAREA,
	//					FieldOption:       "",
	//					FieldDefaultValue: "",
	//					FieldUniqueName:   "ip_white_list",
	//					FieldDisplayName:  "IP白名单",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[\\S]+$/",
	//				},
	//				PluginItemConfig{
	//					FieldType:         PLUGIN_FIELD_TYPE_TEXTAREA,
	//					FieldOption:       "",
	//					FieldDefaultValue: "",
	//					FieldUniqueName:   "ip_black_list",
	//					FieldDisplayName:  "IP黑名单",
	//					FieldRequired:     false,
	//					FieldValidRule:    "/^[\\S]+$/",
	//				},
	//			},
	//		},
	//	},
}
