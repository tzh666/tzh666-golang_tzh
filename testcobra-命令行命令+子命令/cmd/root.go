package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// 定义结构体
var rootCmd = &cobra.Command{
	Use:   "cmdb",          // 程序命令
	Short: "cmdb",          // 帮助信息
	Long:  "cmdb programe", // 详细的帮助信息
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("verbose", verbose)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// 短参数名, 长参数名 --port -P,默认值,帮助信息
	// 全局用PersistentFlags
	// 非全局Flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose info")
}
