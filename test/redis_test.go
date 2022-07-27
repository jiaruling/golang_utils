package test

import (
	"context"
	"testing"

	"github.com/jiaruling/GolangUtil/lib"

	"github.com/go-redis/redis/v8"
)

func TestRedis(t *testing.T) {
	// 初始化项目目录
	lib.InitDir([]string{"./logs", "./static"}, []string{"log"})
	// 初始化日志配置
	lib.InitLog("", "redis", 1, 7, 7, false)
	log := lib.GetLog()
	// 初始化redis连接
	r := lib.NewRedis("192.168.0.142:6379", "abc123456", 2)
	r.InitRedis()
	// crud
	if rdb := lib.GetRedis(); rdb != nil {
		trace := lib.NewTrace()
		ctx := context.Background()
		if err := rdb.Set(ctx, "key", "value", 0).Err(); err != nil {
			log.Error(trace, lib.DLTagRedisFailed, map[string]interface{}{"error": err, "hint_code": 1001})
		}
		val, err := rdb.Get(ctx, "key").Result()
		if err != nil {
			log.Error(trace, lib.DLTagRedisFailed, map[string]interface{}{"error": err, "hint_code": 1002})
		}
		log.Info(trace, lib.DLTagRedisSuccess, map[string]interface{}{"key": val})
		val2, err := rdb.Get(ctx, "key2").Result()
		if err == redis.Nil {
			log.Error(trace, lib.DLTagRedisFailed, map[string]interface{}{"error": "key2 does not exist", "hint_code": 1003})
		} else if err != nil {
			log.Error(trace, lib.DLTagRedisFailed, map[string]interface{}{"error": err, "hint_code": 1004})
		} else {
			log.Info(trace, lib.DLTagRedisSuccess, map[string]interface{}{"key2": val2})
		}
	}
	// 关闭redis
	lib.DestroyRedisAll()
}
