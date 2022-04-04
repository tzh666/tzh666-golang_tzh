package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	port int
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "cmdb web",
	Long:  "cmdb program web",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("web", verbose, port)
		return nil
	},
}

func init() {
	// 把命令加到根命令下
	rootCmd.AddCommand(webCmd)

	// web的子命令 -->  go run .\main.go web -P 8081
	webCmd.Flags().IntVarP(&port, "port", "P", 8080, "web port")
}
