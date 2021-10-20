package template

var (
	BaseConf = baseConf()
	MysqlConf = mysqlConf()
	RedisConf = redisConf()
)

func baseConf() string  {
	baseConf := "# This is base config\n" +
		"[base]\n    " +
		"debug_mode=\"release\"\n    " +
		"time_location=\"Asia/Chongqing\"\n    " +
		"ser_name=\"gatekeeper\"\n\n" +
		"[http]\n    " +
		"addr =\":8880\"                       # 监听地址, default \":8700\"\n    " +
		"read_timeout = 10                   # 建立连接到读取请求内容超时时间\n    " +
		"write_timeout = 10                  # 读取内容到响应的超时时间\n    " +
		"max_header_bytes = 20               # 最大的header大小，二进制位长度\n    " +
		"allow_ip = [                        # 白名单ip列表\n        " +
		"\"127.0.0.1\",\n        " +
		"\"192.168.1.1\"\n    " +
		"]\n\n" +
		"[pprof]\n\n" +
		"[session]\n    " +
		"redis_server = \"#REDIS_CLIENT\"     #redis session server\n    " +
		"redis_password = \"\"\n\n" +
		"[log]\n    " +
		"on = true                          #日志开关，是否开启\n    " +
		"level = \"info\"                      #日志级别，只支持小写\n    " +
		"file_prefix = \"didi\"                #文件名前缀, 默认会加后缀.log 或者.log.wf\n    " +
		"file_dir = \"logs\"                   #生成文件目录\n    " +
		"auto_clear = false                   #是否自动清理日志\n    " +
		"clear_hours = 3                     #保留日志 n 个小时的\n    " +
		"clear_step = 3                      #清理时间从当前时间前推 n 个小时算起\n    " +
		"separate = false                     #文件是否分离，分离文件以.wf结尾\n    " +
		"disable_link = false                #是否启动link\n\n" +
		"[cluster]\n    " +
		"cluster_ip=\"127.0.0.1\"\n    " +
		"cluster_port=\"8080\"\n    " +
		"cluster_ssl_port=\"4433\"\n\n" +
		"[swagger]\n    " +
		"title=\"gatekeeper swagger API\"\n    " +
		"desc=\"This is a sample server celler server.\"\n    " +
		"host=\"\"\n    " +
		"base_path=\"\""
	return baseConf
}

func mysqlConf() string {
	mysqlConf := "[list]\n    " +
		"[list.default]\n        " +
		"driver_name = \"mysql\"\n        " +
		"data_source_name = \"#MYSQL_CLIENT\"\n        " +
		"max_open_conn = 20\n        " +
		"max_idle_conn = 10\n        " +
		"max_conn_life_time = 100"
	return mysqlConf
}

func redisConf() string {
	redisConf := "[list]\n    " +
		"[list.default]\n        " +
		"proxy_list = [\"#REDIS_CLIENT\"]\n        " +
		"conn_timeout = 500\n        " +
		"password = \"\"\n        " +
		"db = 0\n        " +
		"read_timeout = 1000\n        " +
		"write_timeout = 1000\n        " +
		"max_active = 200\n        " +
		"max_idle = 500"
	return redisConf
}