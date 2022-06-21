package cmd

import (
	"UserInsert/apis"
	"UserInsert/models"
	"UserInsert/readfile"
	"UserInsert/setting"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userAddCmd = &cobra.Command{
	Use:   "insert",
	Short: "insert user",
	Long:  "insert user to sc",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 执行的时候实际调用的函数
		var (
			userInfo       *models.UserInfo
			fileUserInfo   *models.FileUserInfo
			portraitLib    *models.PortraitLib
			respInsertUser []byte
			content        []string
			userID         int64
			stateCode      int
		)
		// 初始化配置文件
		setting.InitConf()

		// 初始化日志配置
		setting.InitLog()

		// 读取用户信息文件
		if content, err = readfile.ReadFile(setting.Conf.SenseCity.UserFilePath); err != nil {
			logrus.Error("Read userInfo.txt Error", err)
			return err
		}

		// 遍历切片
		for i := 0; i < len(content); i++ {
			// 反序列化到models.UserInfo
			if err = json.Unmarshal([]byte(content[i]), &fileUserInfo); err != nil {
				logrus.Error("Unmarshal to models Error", err)
				return
			}

			// 调用接口创建用户的用户的信息
			userInfo = &models.UserInfo{
				UserName: fileUserInfo.UserName,
				RealName: fileUserInfo.RealName,
				Serial:   fileUserInfo.Serial,
				Phone:    fileUserInfo.Phone,
				RoleIds:  fileUserInfo.RoleIds,
				OrgId:    fileUserInfo.OrgId,
				ImageUrl: fileUserInfo.ImageUrl,
			}

			// 调用接口创建用户
			userAddURL := fmt.Sprintf("https://%s:10220/uums/users", viper.Get("sensecity.scip").(string))
			respInsertUser, userID, err = apis.InsertUser(userAddURL, apis.GetToken(), userInfo)
			fmt.Println(string(respInsertUser))
			if userID == -1 {
				logrus.Error(" Failed to add user: ", err)
				return
			}

			// 人像库权限控制
			portraitLib = &models.PortraitLib{
				UserID:              userID,
				RealName:            fileUserInfo.RealName,
				AddTarLibSerials:    fileUserInfo.AddTarLibSerials,
				RemoveTarLibSerials: fileUserInfo.RemoveTarLibSerials,
			}

			portraitURL := fmt.Sprintf("https://%s:10220/whale-openapi/portrait-libraries/assigns", viper.Get("sensecity.scip").(string))
			// 调用接口分配用户人像库权限
			if stateCode, err = apis.PortraitLib(portraitURL, viper.Get("sensecity.token").(string), portraitLib); err != nil || stateCode != 200 {
				logrus.Error("Failed to add Portrait: ", err, stateCode)
				logrus.Info("调用接口分配用户人像库权限失败")
				return
			}

			// 打印请求接口相应
			logrus.Info("Request interface corresponding: ")

		}
		fmt.Println("执行成功")
		return
	},
}

func Execute() {
	if err := userAddCmd.Execute(); err != nil {
		logrus.Error("userAddCmd Error:", err)
	}
}
