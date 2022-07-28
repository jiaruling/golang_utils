package test

import (
	"github.com/jiaruling/golang_utils/lib"
	"testing"
)

func TestHttp(t *testing.T) {
	// 初始化项目目录
	lib.InitDir([]string{"./logs", "./static"}, []string{"log"})

	// 初始化日志配置
	lib.InitLog("", "test", 1, 7, 7, false)

	h := lib.Client{}
	h.Request(nil, "GET", "http://192.168.0.142:8085/health",
		nil, 500, nil, lib.JSON)
}
