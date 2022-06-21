package cmd

import (
	"UserInsert/apis"
	"UserInsert/models"
	"UserInsert/setting"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getPortrait = &cobra.Command{
	Use:   "portrait",
	Short: "get portrait",
	Long:  "get sc all portrait",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 执行的时候实际调用的函数
		// 初始化配置文件
		setting.InitConf()

		// 初始化日志配置
		setting.InitLog()
		// 接口返回的人像数据
		var tortraitList models.TortraitList

		// 赋值
		portraits := models.Portraits{
			Page:          setting.Conf.Portrait.Page,
			SenseType:     setting.Conf.Portrait.SenseType,
			SageSize:      setting.Conf.Portrait.SageSize,
			TarLibSerials: setting.Conf.Portrait.TarLibSerials,
		}
		// 获取身份库信息
		respn := apis.GetPortraits(&portraits)

		// 反序列化到tortraitList
		if err := json.Unmarshal(respn, &tortraitList); err != nil {
			fmt.Println(err)
			return err
		}
		// 便利tortraitList.Data.List,然后调用接口是下载图片
		for _, v := range tortraitList.Data.List {
			imgName := fmt.Sprintf("%s.jpg", v.TargetName)
			fmt.Println(apis.DownLoadImage(v.CoverImageUrl, "./images/", imgName))
		}
		return
	},
}

func init() {
	// 把命令加到根命令下
	userAddCmd.AddCommand(getPortrait)
}
