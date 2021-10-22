-- MySQL dump 10.13  Distrib 8.0.26, for macos11.3 (x86_64)
--
-- Host: localhost    Database: gatekeeper
-- ------------------------------------------------------
-- Server version	8.0.26

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `gateway_admin`
--

DROP TABLE IF EXISTS `gateway_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_admin` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `salt` varchar(50) NOT NULL DEFAULT '' COMMENT '盐',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='管理员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_admin`
--

LOCK TABLES `gateway_admin` WRITE;
/*!40000 ALTER TABLE `gateway_admin` DISABLE KEYS */;
INSERT INTO `gateway_admin` VALUES (1,'admin','admin','2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3','2020-04-10 16:42:05','2020-04-21 06:35:08',0);
/*!40000 ALTER TABLE `gateway_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_app`
--

DROP TABLE IF EXISTS `gateway_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_app` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `app_id` varchar(255) NOT NULL DEFAULT '' COMMENT '租户id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '租户名称',
  `secret` varchar(255) NOT NULL DEFAULT '' COMMENT '密钥',
  `white_ips` varchar(1000) NOT NULL DEFAULT '' COMMENT 'ip白名单，支持前缀匹配',
  `qpd` bigint NOT NULL DEFAULT '0' COMMENT '日请求量限制',
  `qps` bigint NOT NULL DEFAULT '0' COMMENT '每秒请求量限制',
  `create_at` datetime NOT NULL COMMENT '添加时间',
  `update_at` datetime NOT NULL COMMENT '更新时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1=删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb3 COMMENT='网关租户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_app`
--

LOCK TABLES `gateway_app` WRITE;
/*!40000 ALTER TABLE `gateway_app` DISABLE KEYS */;
INSERT INTO `gateway_app` VALUES (31,'app_id_a','租户A','449441eb5e72dca9c42a12f3924ea3a2','white_ips',100000,100,'2020-04-15 20:55:02','2020-04-21 07:23:34',0),(32,'app_id_b','租户B','8d7b11ec9be0e59a36b52f32366c09cb','',20,0,'2020-04-15 21:40:52','2020-04-21 07:23:27',0),(33,'app_id','租户名称','','',0,0,'2020-04-15 22:02:23','2020-04-15 22:06:51',1),(34,'app_id45','名称','07d980f8a49347523ee1d5c1c41aec02','',0,0,'2020-04-15 22:06:38','2020-04-15 22:06:49',1);
/*!40000 ALTER TABLE `gateway_app` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gateway_service_info`
--

DROP TABLE IF EXISTS `gateway_service_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateway_service_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_type` tinyint NOT NULL DEFAULT '0' COMMENT '服务类型 0=http 1=tcp 2=grpc',
  `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称 6-128 数字字母下划线',
  `service_desc` varchar(255) NOT NULL DEFAULT '' COMMENT '服务描述',
  `service_port` int unsigned NOT NULL DEFAULT '0' COMMENT '服务端口(只针对 tcp/grpc)',
  `http_hosts` varchar(1000) NOT NULL DEFAULT '' COMMENT 'http域名信息',
  `http_paths` varchar(1000) NOT NULL DEFAULT '' COMMENT 'http路径信息',
  `http_strip_prefix` tinyint NOT NULL DEFAULT '1' COMMENT 'http转发前剥离前缀',
  `load_balance_strategy` varchar(255) NOT NULL DEFAULT '' COMMENT '负载策略',
  `load_balance_type` varchar(255) NOT NULL DEFAULT '' COMMENT '负载类型',
  `auth_type` varchar(255) NOT NULL DEFAULT '' COMMENT '鉴权类型',
  `upstream_list` varchar(255) NOT NULL DEFAULT '' COMMENT '下游服务器ip和权重',
  `plugin_conf` mediumtext COMMENT '插件配置',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '添加时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 1=删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='网关基本信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gateway_service_info`
--

LOCK TABLES `gateway_service_info` WRITE;
/*!40000 ALTER TABLE `gateway_service_info` DISABLE KEYS */;
INSERT INTO gateway_service_info VALUES (1,0,'test_service_name','test_service_desc',0,'','/test_service_name',1,'random','','','http://127.0.0.1:8881 100','{\"url_rewrite\":{\"rewrite_rule\":\"\"},\"http_flow_limit\":{\"service_flow_limit_num\":\"\",\"service_flow_limit_type\":\"0\",\"clientip_flow_limit_num\":\"\",\"clientip_flow_limit_type\":\"\"},\"header_transfer\":{\"transfer_rule\":\"\"},\"http_whiteblacklist\":{\"ip_white_list\":\"\",\"ip_black_list\":\"\",\"host_white_list\":\"\",\"url_white_list\":\"\"},\"http_upstream_transport\":{\"http_upstream_connection_timeout\":\"\",\"http_upstream_header_timeout\":\"\"},\"jwt_auth\":{},\"upstream_config\":{\"upstream_list\":\"http://127.0.0.1:8881 100\nhttp://127.0.0.1:8882 100\"}}','2021-09-20 10:55:46','2021-09-20 11:37:50',0);
/*!40000 ALTER TABLE `gateway_service_info` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-09-20 13:48:17
