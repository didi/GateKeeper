package check

import (
	"fmt"
	"gatekeeper/install/template"
	"gatekeeper/install/tool"
	"io/ioutil"
	"strings"
)

var (
	ConfPath = GateKeeperPath + "/conf/dev/"
)

func InitConf() error {

	tool.LogInfo.Println("init conf start")

	err := initBase()
	if err != nil {
		return err
	}


	err = initRedis()
	if err != nil {
		return err
	}

	err = initMysql()
	if err != nil {
		return err
	}
	tool.LogInfo.Println("init conf end")
	return nil
}


func initBase() error {
	tool.LogInfo.Println("init base conf")

	fileName := ConfPath + "base.toml"
	redisClient := RedisClient.Host + ":" + RedisClient.Port
	baseConf := strings.Replace(template.BaseConf, "#REDIS_CLIENT", redisClient, 1)
	err := ioutil.WriteFile(fileName, []byte(baseConf), 0666); if err != nil{
		return err
	}
	return nil

}


func initRedis() error {
	tool.LogInfo.Println("init redis conf")

	fileName := ConfPath + "redis_map.toml"
	redisClient := RedisClient.Host + ":" + RedisClient.Port
	redisConf := strings.Replace(template.RedisConf, "#REDIS_CLIENT", redisClient, 1)
	err := ioutil.WriteFile(fileName, []byte(redisConf), 0666); if err != nil{
		return err
	}
	return nil

}


func initMysql() error {
	tool.LogInfo.Println("init mysql conf")

	fileName := ConfPath + "mysql_map.toml"
	mysqlClient := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		MysqlClient.User,
		MysqlClient.Pwd,
		MysqlClient.Host,
		MysqlClient.Port,
		MysqlClient.Database)
	mysqlConf := strings.Replace(template.MysqlConf, "#MYSQL_CLIENT", mysqlClient, 1)
	err := ioutil.WriteFile(fileName, []byte(mysqlConf), 0666); if err != nil{
		return err
	}
	return nil
}




