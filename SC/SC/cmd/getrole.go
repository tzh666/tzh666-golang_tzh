package cmd

import (
	"UserInsert/apis"
	"UserInsert/dao/mysql"
	"UserInsert/setting"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getRole = &cobra.Command{
	Use:   "role",
	Short: "get role",
	Long:  "get sc all role",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 执行的时候实际调用的函数
		// 初始化配置文件
		setting.InitConf()

		// 初始化日志配置
		setting.InitLog()

		// mysql 连接
		if err = mysql.InitMySQLUums(); err != nil {
			logrus.Error("InitMySQL Error", err)
			return
		}

		defer mysql.Close()

		// 查询角色
		apis.GetRole()

		return
	},
}

func init() {
	// 把命令加到根命令下
	userAddCmd.AddCommand(getRole)
}
