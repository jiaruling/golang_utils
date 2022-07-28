package test

import (
	"github.com/jiaruling/golang_util/lib"
	"testing"
)

func TestMysql(t *testing.T) {
	// 初始化项目目录
	lib.InitDir([]string{"./logs", "./static"}, []string{"log"})
	// 初始化日志配置
	lib.InitLog("", "mysql", 1, 7, 7, false)
	log := lib.GetLog()

	// gorm
	// 初始化mysql连接
	m := lib.NewMysqlGorm("root", "abc123456", "192.168.0.142", 3366, "gm_customer_service_system")
	m.InitMysqlGorm()

	// curd
	trace := lib.NewTrace()
	if db := lib.GetMysqlGorm(); db != nil {
		var count int64
		result := db.Table("log").Count(&count)
		if result.Error != nil {
			log.Error(trace, lib.DLTagMySqlFailed, map[string]interface{}{"hint_code": 1001, "error": result.Error})
		} else {
			log.Info(trace, lib.DLTagMySqlFailed, map[string]interface{}{"count": count, "table": "log"})
		}
	} else {
		log.Error(trace, lib.DLTagMySqlFailed, map[string]interface{}{"hint_code": 1002, "error": "获取db失败"})
	}

	// sqlx
	// 初始化mysql连接
	mx := lib.NewMysqlX("root", "abc123456", "192.168.0.142", 3366, "gm_customer_service_system")
	mx.Name = "R"
	mx.InitMysqlX()

	if db := lib.GetMysqlX("R"); db != nil {
		var count int64
		err := db.Get(&count, "select count(1) from menu;")
		if err != nil {
			log.Error(trace, lib.DLTagMySqlFailed, map[string]interface{}{"hint_code": 1003, "error": err})
		} else {
			log.Info(trace, lib.DLTagMySqlFailed, map[string]interface{}{"count": count, "table": "menu"})
		}
	} else {
		log.Error(trace, lib.DLTagMySqlFailed, map[string]interface{}{"hint_code": 1004, "error": "获取db失败"})
	}

	// 关闭数据库连接
	lib.DestroyMysqlGormAll()
	lib.DestroyMysqlXAll()
}
