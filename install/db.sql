
CREATE TABLE `gateway_match_rule` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `module_id` bigint(20) NOT NULL COMMENT '模块id',
  `type` varchar(200) NOT NULL DEFAULT '' COMMENT '匹配类型',
  `rule` varchar(1000) NOT NULL DEFAULT '' COMMENT '规则',
  `rule_ext` varchar(1000) NOT NULL DEFAULT '' COMMENT '拓展规则',
  `url_rewrite` varchar(800) NOT NULL COMMENT '重写规则',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网关路由匹配表';
INSERT INTO `gateway_match_rule` VALUES (122,26,'url_prefix','/gatekeeper/test_http','','^/gatekeeper/test_http(.*) $1'),(128,27,'url_prefix','','','');

CREATE TABLE `gateway_access_control` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `module_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '模块id',
  `open` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否开启权限 1=开启',
  `black_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '黑名单ip',
  `white_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单ip',
  `white_host_name` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单主机',
  `auth_type` varchar(100) NOT NULL DEFAULT '' COMMENT '认证方法，dev强制走固定ip',
  `client_flow_limit` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端ip限流',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网关权限控制表';

INSERT INTO `gateway_access_control` VALUES (119,26,1,'','','','',0),(125,27,1,'','','','',0);

CREATE TABLE `gateway_load_balance` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `module_id` bigint(20) NOT NULL,
  `check_method` varchar(200) NOT NULL DEFAULT '' COMMENT '检查方法',
  `check_url` varchar(500) NOT NULL DEFAULT '' COMMENT '检测url',
  `check_timeout` int(10) NOT NULL DEFAULT '0' COMMENT 'check超时时间',
  `check_interval` int(11) NOT NULL DEFAULT '0' COMMENT '检查频率',
  `type` varchar(100) NOT NULL DEFAULT '' COMMENT '轮询方式',
  `ip_list` varchar(1000) NOT NULL DEFAULT '' COMMENT 'ip列表',
  `weight_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '权重列表',
  `forbid_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '禁用ip列表',
  `proxy_connect_timeout` int(11) NOT NULL DEFAULT '0' COMMENT '连接超时, 单位ms',
  `proxy_header_timeout` int(11) NOT NULL DEFAULT '0' COMMENT '获取header超时, 单位ms',
  `proxy_body_timeout` int(11) NOT NULL DEFAULT '0' COMMENT '获取body超时, 单位ms',
  `idle_conn_timeout` int(10) NOT NULL DEFAULT '0' COMMENT '链接最大空闲时间, 单位ms',
  `max_idle_conn` int(11) NOT NULL DEFAULT '0' COMMENT '最大空闲链接数',
  `max_retry_time` int(11) NOT NULL DEFAULT '0' COMMENT '重试次数，0为不重试',
  `retry_interval` int(11) NOT NULL DEFAULT '0' COMMENT '单位ms，重试间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网关负载表';

INSERT INTO `gateway_load_balance` VALUES (119,26,'httpchk','/ping',2000,5000,'round-robin','100.90.164.31:8072','50','',10001,10002,10003,10004,200,0,0),(125,27,'tcpchk','',2000,5000,'round-robin','127.0.0.1:8018','50','',10001,0,0,0,0,0,0);


CREATE TABLE `gateway_module_base` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `load_type` varchar(100) NOT NULL DEFAULT 'http' COMMENT '负载类型 http/tcp',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '模块名',
  `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称',
  `pass_auth_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '#传参类型 0=不传 1=传递全部信息 2=传递rediskey',
  `frontend_addr` varchar(100) NOT NULL DEFAULT '' COMMENT '前端绑定ip地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网关模块表';

INSERT INTO `gateway_module_base` VALUES (26,'http','test_http','test_http',2,''),(27,'tcp','test_tcp','test_tcp',2,':8900');


CREATE TABLE `gateway_app` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `app_id` varchar(255) NOT NULL DEFAULT '' COMMENT '租户id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '租户名称',
  `secret` varchar(255) NOT NULL DEFAULT '' COMMENT '密钥',
  `method` varchar(100) NOT NULL DEFAULT 'any' COMMENT '请求方法',
  `timeout` bigint(20) NOT NULL DEFAULT '0' COMMENT '超时时间',
  `open_api` text NOT NULL COMMENT '接口列表，支持前缀匹配',
  `white_ips` varchar(1000) NOT NULL DEFAULT '' COMMENT 'ip白名单，支持前缀匹配',
  `city_ids` varchar(1000) NOT NULL DEFAULT '' COMMENT 'city_id数据权限',
  `total_query_daily` bigint(20) NOT NULL DEFAULT '0' COMMENT '日请求量',
  `qps` bigint(20) NOT NULL DEFAULT '0' COMMENT 'qps',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '数据关联id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网关租户表';

INSERT INTO `gateway_app` VALUES (28,'test_app','test_app','62fda0f2212eaffd90dbf04136768c5f','any',0,'/gatekeeper','','',1000000,0,0);
