package main

import (
	"UserInsert/apis"
	"UserInsert/dao/mysql"
	"UserInsert/models"
	"UserInsert/readfile"
	"UserInsert/setting"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	var (
		userInfo *models.UserInfo
		resp     []byte
		content  []string
		err      error
	)
	// 初始化配置文件
	setting.InitConf()

	// 初始化日志配置
	setting.InitLog()

	// mysql 连接
	if err = mysql.InitMySQL(); err != nil {
		logrus.Error("InitMySQL Error", err)
		return
	}

	// 读取用户信息文件
	if content, err = readfile.ReadFile(setting.Conf.SenseCity.UserFilePath); err != nil {
		logrus.Error("Read userInfo.txt Error", err)
		return
	}

	// 遍历切片
	for i := 0; i < len(content); i++ {
		// 反序列化到models.UserInfo
		if err = json.Unmarshal([]byte(content[i]), &userInfo); err != nil {
			logrus.Error("Unmarshal to models Error", err)
		}

		// 调用接口创建用户
		if resp, err = apis.PostHttpsSkip(setting.Conf.SenseCity.Url, viper.Get("sensecity.token").(string), userInfo); err != nil {
			logrus.Error(" Failed to add user: ", err)
			return
		}

		// 打印请求接口相应
		logrus.Info("Request interface corresponding: ", string(resp))
	}
	fmt.Println("任务执行完毕,详情请看日志文件")
}
