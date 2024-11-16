package util

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

// 保存项目根目录，即当前文件所在目录util上一级
var (
	ProjectRootPath = path.Dir(getoncurrentPath()+"/../") + "/"
)

func getoncurrentPath() string {
	_, foldername, _, _ := runtime.Caller(0) // 0代表本行代码位置，参数可为0，1，2
	return path.Dir(foldername)              // 返回本行代码位置的文件路径即目录
}

// Viper可以用来解析JSON、TOML、YAML、HCL、INI、ENV等配置文件格式
// 可以用来监听配置文件的变化（WatchConfig），无需重启程序就可读取latest
// file string就是配置文件的名称，例如key.yaml,传入key即可,配置文件保存在config，通过根目录路径获取
func CreateConfig(file string) *viper.Viper {
	config := viper.New()
	configPath := ProjectRootPath + "config/"
	config.AddConfigPath(configPath) // 读取配置文件目录
	config.SetConfigName(file)       // 通过文件名读取配置文件
	config.SetConfigType("yaml")     // 设定文件类型
	configFile := configPath + file + ".yaml"

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// os.Exit(2)
			// fmt.Errorf方式构造error
			panic(fmt.Errorf("config file not found error:%s", configFile)) // 初始化阶段错误，直接panic结束
		} else {
			panic(fmt.Errorf("read config file %s error:%s", configFile, err))
		}
	}
	return config
}
