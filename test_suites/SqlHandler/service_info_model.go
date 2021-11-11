package testsqlhandler

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/model"
	"github.com/e421083458/gorm"
)

var Db *gorm.DB

func GetServiceStripPrefix(serviceName string) int {
	db, err := lib.GetGormPool("default")
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	task := model.ServiceInfo{}
	db.Where("service_name = ?", serviceName).First(&task)
	return task.HttpStripPrefix
}

func GetServiceLoadBalanceStrategy(serviceName string) string {
	db, err := lib.GetGormPool("default")
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	task := model.ServiceInfo{}
	db.Where("service_name = ?", serviceName).First(&task)
	return task.LoadBalanceStrategy
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

//func Save(serviceInfo *model.ServiceInfo) {
//	tmp, err := lib.GetGormPool("default")
//	if err != nil {
//		panic("连接数据库失败, error=" + err.Error())
//	}
//
//	err = tmp.Save(serviceInfo).Error
//	if err != nil {
//		fmt.Println("SAVE:", err)
//	}
//}
