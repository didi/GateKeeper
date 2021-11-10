package testsqlhandler

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/didi/gatekeeper/model"
	"github.com/e421083458/gorm"
)

var Db *gorm.DB

func InitGORMHandler() {
	tmp, err := lib.GetGormPool("default")
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	Db = tmp
	log.Info().Msg("InitConfig GORMHandler Success.")
}

func UpdateServiceStripPrefix(serviceName string, stripPrefix int) {
	Db.Model(&model.ServiceInfo{}).Where("service_name = ?", serviceName).Update("http_strip_prefix", stripPrefix)
}

func GetServiceStripPrefix(serviceName string) int {
	task := model.ServiceInfo{}
	Db.Where("service_name = ?", serviceName).First(&task)
	return task.HttpStripPrefix
}

func DeleteServiceInfo(serviceName string) {
	Db.Where("service_name = ?", serviceName).Delete(model.ServiceInfo{})
}

func AddServiceInfo(serviceInfo *model.ServiceInfo) {
	Db.Create(serviceInfo)
}
