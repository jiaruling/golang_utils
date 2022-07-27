# Alpine Linux 系统下开发注意事项

## 安装编译环境

``` shell
$ apk add gcc
$ apk add musl-dev
```

## 修改时区

``` shell
# url: https://blog.csdn.net/isea533/article/details/87261764

# 安装时区设置
$ apk add tzdata
# 复制上海时区
$ cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 指定为上海时区
$ echo "Asia/Shanghai" > /etc/timezone
# 验证时区
$ date -R
Fri, 22 Jul 2022 17:07:46 +0800
# 删除其它时区配置, 节省空间 【可选】
$ apk del tzdata
```

# golang 开发常用工具包功能介绍

## 链路日志

```go
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
```

## http client

```go
func TestHttp(t *testing.T) {
	// 初始化项目目录
	lib.InitDir([]string{"./logs", "./static"}, []string{"log"})

	// 初始化日志配置
	lib.InitLog("", "test", 1, 7, 7, false)

	h := lib.Client{}
	h.Request(nil, "GET", "http://192.168.0.142:8085/health",
		nil, 500, nil, lib.JSON)
}
```

## Tips

- 数据去重
- struct 转 map
- ListStruct 转 ListMap
- 生成SHA1
- StringToBytes 实现string 转换成 []byte, 不用额外的内存分配
- BytesToString 实现 []byte 转换成 string, 不需要额外的内存分配
- 拼接url
- 判断元素是否在数组/切片内
- 生成任意长度的随机字符串

## 文件操作

## 配置文件解析

## 数据操作

### MySQL

```go
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
```

### Redis

```go
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
```



