package public

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
)

var (
	//GormPool gorm连接池
	GormPool *gorm.DB
)

//InitMysql 初始化mysql
func InitMysql() error {
	dbpool, err := lib.GetGormPool("default")
	if err != nil {
		return err
	}
	GormPool = dbpool
	return nil
}