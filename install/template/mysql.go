package template

import (
	"bufio"
	"gatekeeper/install/tool"
	"os"
	"regexp"
	"strings"
)

var(
	TableSql map[string]string
	Tables map[string][]string
)

func InitSql() error {
	TableSql = make(map[string]string)
	Tables = make(map[string][]string)
	err := getCreateSql()
	if err != nil{
		return err
	}
	for tableName, sql := range TableSql {
		Tables[tableName] = strings.Split(sql, ";")
	}
	return nil
}

<<<<<<< HEAD
func gateway_service_info() []string  {
	sql := []string{
		"DROP TABLE IF EXISTS `gateway_service_info`;",
		"CREATE TABLE `gateway_service_info` (" +
			"`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键'," +
			"`service_type` tinyint NOT NULL DEFAULT '0' COMMENT '负载类型 0=http 1=tcp 2=grpc'," +
			"`service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称 6-128 数字字母下划线'," +
			"`service_desc` varchar(255) NOT NULL DEFAULT '' COMMENT '服务描述'," +
			"`service_port` int unsigned NOT NULL DEFAULT '0' COMMENT 'tcp/grpc端口'," +
			"`http_hosts` varchar(1000) NOT NULL DEFAULT '' COMMENT '域名信息'," +
			"`http_paths` varchar(1000) NOT NULL DEFAULT '' COMMENT '路径信息'," +
			"`http_strip_prefix` varchar(255) NOT NULL DEFAULT '' COMMENT '是否需要strip_uri'," +
			"`load_balance_strategy` varchar(255) NOT NULL DEFAULT '' COMMENT '负载策略'," +
			"`load_balance_type` varchar(255) NOT NULL DEFAULT '' COMMENT '负载类型'," +
			"`auth_type` varchar(255) NOT NULL DEFAULT '' COMMENT '鉴权类型'," +
			"`upstream_list` varchar(255) NOT NULL DEFAULT '' COMMENT '下游服务器ip和权重'," +
			"`plugin_conf` mediumtext COMMENT '插件配置'," +
			"`create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '添加时间'," +
			"`update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间'," +
			"`is_delete` tinyint DEFAULT '0' COMMENT '是否删除 1=删除'," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='网关基本信息表';",
			"LOCK TABLES `gateway_service_info` WRITE;",
			"INSERT INTO `gateway_service_info` VALUES (1,0,'test_service_name','test_service_desc',0,'','/test_service_name','0','random','upstream_config','','3','{\\\"url_rewrite\\\":{\\\"rewrite_rule\\\":\\\"\\\"},\\\"http_flow_limit\\\":{\\\"service_flow_limit_num\\\":\\\"\\\",\\\"service_flow_limit_type\\\":\\\"0\\\",\\\"clientip_flow_limit_num\\\":\\\"\\\",\\\"clientip_flow_limit_type\\\":\\\"\\\"},\\\"header_transfer\\\":{\\\"transfer_rule\\\":\\\"\\\"},\\\"http_whiteblacklist\\\":{\\\"ip_white_list\\\":\\\"\\\",\\\"ip_black_list\\\":\\\"\\\",\\\"host_white_list\\\":\\\"\\\",\\\"url_white_list\\\":\\\"\\\"},\\\"http_upstream_transport\\\":{\\\"http_upstream_connection_timeout\\\":\\\"\\\",\\\"http_upstream_header_timeout\\\":\\\"\\\"},\\\"jwt_auth\\\":{},\\\"upstream_config\\\":{\\\"upstream_list\\\":\\\"http://127.0.0.1:8081 100\\\\nhttp://127.0.0.1:8081 100\\\"}}','2021-09-20 10:55:46','2021-09-20 11:37:50',0);",
			"UNLOCK TABLES;",

=======
func getCreateSql() error{
	sqlFilePath := tool.GateKeeperPath + "/gatekeeper.sql"
	tool.LogInfo.Println("sql file path :" + sqlFilePath)
	tableName := ""
	tablePre := regexp.MustCompile(`gateway_[a-z_]*`)
	f, _ := os.Open(sqlFilePath)
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		lineSql, err := tool.ReadLine(r)
		if strings.Index(lineSql, "/*") < 0 && strings.Index(lineSql, "--") < 0 && lineSql != "" {
			if strings.Index(lineSql, "DROP TABLE IF EXISTS") >= 0 {
				tableName = tablePre.FindString(lineSql)
			}
			TableSql[tableName] += lineSql
		}
		if err != nil {
			break
		}
>>>>>>> e05d82fab91d0df1fd9b7eb3deebda7b00c90b44
	}
	return nil
}