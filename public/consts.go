package public

const (
	//UserGroupPerfix 用户组权限redis key前缀
	UserGroupPerfix = "gatekeeper_upm_"

	//AccessControlAppIDTotalCallPrefix appid级别的 TotalQueryDaily
	AccessControlAppIDTotalCallPrefix     = "gatekeeper_appid_totalcall_"
	//AccessControlAppIDHourTotalCallPrefix appid hour
	AccessControlAppIDHourTotalCallPrefix = "gatekeeper_appid_hour_totalcall_"

	//ContentEncoding header相关
	ContentEncoding = "Content-Encoding"

	//RequestModuleCounterPrefix limit相关
	RequestModuleCounterPrefix     = "gatekeeper_module_counter_"
	//RequestModuleHourCounterPrefix 模块小时前缀
	RequestModuleHourCounterPrefix = "gatekeeper_module_hour_counter_"

	//AdminCookiePrefix admin相关
	AdminCookiePrefix = "admin_"
	//AdminCookieSecrit aes密钥
	AdminCookieSecrit = "1122334455667788"
	//AdminExpired 管理员登陆超时时间
	AdminExpired      = 14400

	//IPDefaultWeight 默认ip权重
	IPDefaultWeight = 50
)
