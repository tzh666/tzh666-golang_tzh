package mysql

import (
	"UserInsert/setting"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

func InitMySQL() (err error) {
	// 数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		// 通过结构体的形式去读取配置文件数据
		setting.Conf.MySQL.UserName,
		setting.Conf.MySQL.PassWord,
		setting.Conf.MySQL.Host,
		setting.Conf.MySQL.Port,
		setting.Conf.MySQL.DBName,
	)

	// 连接数据库
	if DB, err = gorm.Open("mysql", dsn); err != nil {
		logrus.Error("Open  MySQL Error:", err)
		return
	}

	// 禁用默认表名的复数形式，如果置为 true，则 `User` 的默认表名是 `user`
	DB.SingularTable(true)

	logrus.Info("连接数据库成功")
	return
}

func InitMySQLUums() (err error) {
	// 数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		// 通过结构体的形式去读取配置文件数据
		setting.Conf.MySQL.UserName,
		setting.Conf.MySQL.PassWord,
		setting.Conf.MySQL.Host,
		setting.Conf.MySQL.Port,
		"uums",
	)

	// 连接数据库
	if DB, err = gorm.Open("mysql", dsn); err != nil {
		logrus.Error("Open  MySQL Error:", err)
		return
	}

	// 禁用默认表名的复数形式，如果置为 true，则 `User` 的默认表名是 `user`
	DB.SingularTable(true)

	logrus.Info("连接数据库成功")
	return
}

// 对外暴露一个关闭的方法
func Close() {
	_ = DB.Close()
}
