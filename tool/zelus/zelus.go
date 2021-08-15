package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-zelus/zelus/tool/zelus/api"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "zelus",
	Short:   "",
	Long:    "",
	Version: fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH),
}

func init() {
	rootCmd.AddCommand(api.ApiCmd)
	rootCmd.SetHelpTemplate("-h, --help 帮助文档\n")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
