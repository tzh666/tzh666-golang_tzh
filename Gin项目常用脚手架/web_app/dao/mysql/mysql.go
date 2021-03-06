package mysql

import (
	"fmt"
	"web_app/settings"

	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		/*viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
		*/
		// 通过结构体的形式去读取配置文件数据
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.SetMaxOpenConns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.SetMaxIdleConns"))
	return
}

// 对外暴露一个关闭的方法
func Close() {
	_ = db.Close()
}
