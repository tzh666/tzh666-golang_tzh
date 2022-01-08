package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 1、加载配置
func Init() (err error) {
	// 指定要读取的配置文件
	viper.SetConfigType("yaml")   // 指定文件格式
	viper.SetConfigName("config") // 文件名不用加后缀
	viper.AddConfigPath("./conf") // 路径

	// 判断读取配置文件是否有误
	if err := viper.ReadInConfig(); err != nil {
		//logrus.Debug(err)
		fmt.Printf("读取conf配置文件出错: %s\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
	return
}
