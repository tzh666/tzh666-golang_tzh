package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	host string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "cmdb api",
	Long:  "cmdb program api",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("api", verbose, host)
		return nil
	},
}

func init() {
	// 把命令加到根命令下
	rootCmd.AddCommand(apiCmd)

	/*
			&db,        要绑定的自定义参数地址
			"database", 长信息
			"d",        短信息 (只能是一个字符)
			"default",  默写
		    "database"  说明信息
	*/
	// api的子命令
	apiCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "web host")

}
