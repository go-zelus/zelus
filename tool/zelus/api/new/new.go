package new

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"

	"github.com/go-zelus/zelus/tool/zelus/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "创建新项目",
	Long:  "创建新项目",
	Run:   run,
}

var (
	out    string
	repo   string
	branch string
)

func init() {
	out = "."
	if repo = os.Getenv("ZELUS_API_REPO"); repo == "" {
		repo = "https://github.com/go-zelus/api-layout.git"
	}
	NewCmd.Flags().StringVarP(&out, "out", "o", out, "项目路径")
	NewCmd.Flags().StringVarP(&repo, "repo", "r", repo, "git地址")
	NewCmd.Flags().StringVarP(&branch, "branch", "b", branch, "git分支")

}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "新建项目名称",
			Help:    "新建项目名称",
		}
		survey.AskOne(prompt, &name)
		if name == "" {
			return
		}
	} else {
		name = args[0]
	}
	if out != "." {
		wd, err = filepath.Abs(out)
		if err != nil {
			fmt.Println(color.RedString(err.Error()))
			return
		}
	}
	err = New(ctx, wd, name, repo, branch)
	if err != nil {
		fmt.Println(err)
	}
}

func New(ctx context.Context, dir string, name string, layout string, branch string) error {
	to := path.Join(dir, path.Base(name))
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Println(color.RedString(fmt.Sprintf("%s 项目已存在", name)))
		override := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintln(color.RedString("是否覆盖已存在的项目?")),
			Help:    fmt.Sprintln(color.RedString("该操作将删除现有项目，重新建立新的项目")),
		}
		survey.AskOne(prompt, &override)
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	rep := util.NewRepo(layout, branch)
	if err := rep.CopyTo(ctx, to, name, []string{".git", ".github"}); err != nil {
		return err
	}
	util.Tree(to, dir)
	return nil
}
