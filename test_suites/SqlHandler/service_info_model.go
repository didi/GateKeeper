package testsqlhandler

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/model"
)

func GetServiceStripPrefix(serviceName string) int {
	db, err := lib.GetGormPool("default")
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	task := &model.ServiceInfo{}
	if err := db.Where("service_name = ?", serviceName).Find(task).Error; err != nil {
		panic("查询数据表 , error=" + err.Error())
	}
	if task != nil {
		return task.HttpStripPrefix
	}
	return 0
}

func DeleteServiceInfo(serviceName string) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	db.Where("service_name = ?", serviceName).Delete(model.ServiceInfo{})
}

func AddServiceInfo(serviceInfo *model.ServiceInfo) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	db.Create(serviceInfo)
}
