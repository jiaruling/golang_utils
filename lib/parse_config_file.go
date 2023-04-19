package lib

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 加载配置文件并进行监听
func ParseConfig(path string, obj interface{}, isWatch bool) (err error) {
	// 读取配置文件，映射到结构体
	// 实例化viper对象
	v := viper.New()
	v.SetConfigFile(path)
	// 读取配置文件
	if err = v.ReadInConfig(); err != nil {
		return
	}
	// 反序列化为 struct 对象
	if err = v.Unmarshal(obj); err != nil {
		return
	}

	if isWatch {
		// viper的功能 -- 动态监控变化
		fmt.Println("开启配置 " + path + " 文件监听")
		go func() {
			// 开启监听功能
			v.WatchConfig()
			// 文件监听
			v.OnConfigChange(func(e fsnotify.Event) {
				// 打印变换的文件名
				fmt.Println("配置文件发生变化:", e.Name)
				_ = v.ReadInConfig() // 重新读取配置数据
				_ = v.Unmarshal(obj) // 将文件内容映射到结构体
			})
		}()
	}
	return
}

var ViperConfMap map[string]*viper.Viper

func ParseConfigViper(path, suffix string, isWatch bool) (err error) {
	// 创建一个新的 viper 实例
	v := viper.New()

	// 设置配置文件的类型
	v.SetConfigType(suffix)

	// 设置配置文件的路径
	v.SetConfigFile(path)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// 放入ViperConfMap中
	_, file := filepath.Split(path)
	fileName := strings.Split(file, ".")[0]
	if ViperConfMap == nil {
		ViperConfMap = make(map[string]*viper.Viper)
	}
	ViperConfMap[fileName] = v

	if isWatch {
		// viper的功能 -- 动态监控变化
		fmt.Println("开启配置 " + path + " 文件监听")
		go func() {
			// 开启监听功能
			v.WatchConfig()
			// 文件监听
			v.OnConfigChange(func(e fsnotify.Event) {
				// 打印变换的文件名
				fmt.Println("配置文件发生变化:", e.Name)
				err := v.ReadInConfig() // 重新读取配置数据
				if err != nil {
					return
				} else {
					_, file := filepath.Split(e.Name)
					fileName := strings.Split(file, ".")[0]
					ViperConfMap[fileName] = v
				}
			})
		}()
	}
	return
}
