package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 定义一个全局的变量,到时候就通过这个全局的变量去调用自定义的结构体的内容
var Conf = new(AppConfig)

// 对应yaml的字段  (Name--转成小写对应name),如果有不符合的就用标签去指定`mapstructure:"name"`
type AppConfig struct {
	Name         string
	Mode         string
	Port         int
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int `mapstructure:"max_size"`
	MaxAge     int `mapstructure:"max_age"`
	MaxBackups int `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string `mapstructure:"dbname"`
	SetMaxOpenConns int    `mapstructure:"SetMaxOpenConns"`
	SetMaxIdleConns int    `mapstructure:"SetMaxIdleConns"`
}

type RedisConfig struct {
	Host     string
	Port     int
	DBName   int `mapstructure:"db"`
	Password string
}

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

	// 反序列化到配置结构体中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("Unmarshal failed ,err:#{err}\n")
	}

	viper.WatchConfig()
	// 自动判断配置文件是否改变
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
		// 反序列化到配置结构体中,Conf是指针所以改的跟上面的Conf是同一个东西
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("Unmarshal failed ,err:#{err}\n")
		}
		// 当然你想做其他事也可以,这里就是个钩子函数。当改变了配置文件以后就干啥干啥
	})
	return
}
