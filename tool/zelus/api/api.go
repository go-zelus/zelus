package api

import (
	apinew "github.com/go-zelus/zelus/tool/zelus/api/new"
	"github.com/spf13/cobra"
)

var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "api项目命令",
	Long:  "api项目命令",
	Run:   run,
}

func init() {
	ApiCmd.AddCommand(apinew.NewCmd)
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	}
}
