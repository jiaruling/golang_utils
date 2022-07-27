package test

import (
	"github.com/jiaruling/GolangUtil/lib"
	"testing"
)

func TestLog(t *testing.T) {
	// ----------------------------------------------
	// 初始化项目目录
	lib.InitDir([]string{"./logs", "./static"}, []string{"log"})

	// ----------------------------------------------
	// 初始化日志配置
	lib.InitLog("", "test", 1, 7, 7, false)

	log := lib.GetLog()
	for i := 0; i < 10; i++ {
		log.Debug(lib.NewTrace(), lib.DLTagRequestIn, map[string]interface{}{"method": "GET", "URL": "DEBUG", "body": "debugBody"})
		log.Info(lib.NewTrace(), lib.DLTagRequestIn, map[string]interface{}{"method": "POST", "URL": "IFNO", "body": "debugINFO"})
		log.Warn(lib.NewTrace(), lib.DLTagRequestIn, map[string]interface{}{"method": "PUT", "URL": "WARN", "body": "debugWARN"})
		log.Error(lib.NewTrace(), lib.DLTagRequestIn, map[string]interface{}{"method": "DELETE", "URL": "ERROR", "body": "debugERROE"})
	}
}
