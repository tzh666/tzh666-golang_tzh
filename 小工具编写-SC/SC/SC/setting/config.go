package setting

import (
	"UserInsert/models"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 单例模式
var Conf = new(models.Config)

func InitConf() {
	// 指定要读取的配置文件
	viper.SetConfigType("yaml") // 指定文件格式
	viper.SetConfigName("api")  // 文件名不用加后缀
	viper.AddConfigPath("conf") // 路径conf/api.yaml

	// 判断读取配置文件是否有误
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error("Read configuration file failed, err:", err)
		return
	}

	// 设置默认值,配置文件不存在的时候才会用这个默认值,存在就不用,默认配置文件一定要在viper.Unmarshal(Conf)之前设置
	viper.SetDefault("sensecity.userfilepath", "conf/userInfo.txt")

	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		logrus.Error("unmarshal conf failed, err:", err)
		return
	}

	// 监控配置文件变化
	viper.WatchConfig()

	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 日志记录
		logrus.Info("...配置文件已被修改...")

		// 配置文件发生变化后要同步到全局变量Conf
		if err := viper.Unmarshal(Conf); err != nil {
			logrus.Error("unmarshal conf failed, err:", err)
			return
		}

		// 到时候自动重启程序？？？
		// 有待编写
	})
}
