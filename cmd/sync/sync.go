package sync

import (
	"fmt"
	"os"

	"github.com/gotomicro/ego"
	"github.com/spf13/cobra"

	"push-seo/cmd"
	"push-seo/invoker"
	"push-seo/process"
)

// 同步指定sitmaps url结果到百度SEO
var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "sync sitmaps to baidu-seo",
	Example: "push-seo sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ego.New().Invoker(invoker.Init, Run).Run(); err != nil {
			return fmt.Errorf("run fail, %w", err)
		}
		return nil
	},
}

var (
	sitemaps []string
	dryRun   bool
)

func Run() error {
	token := os.Getenv("BAIDU_TOKEN")
	if token == "" {
		return fmt.Errorf("environment BAIDU_TOKEN can't be empty")
	}
	p := process.NewProcessor(sitemaps, token, dryRun, invoker.S3Cli, invoker.HttpCli)
	return p.Process()
}

func init() {
	syncCmd.PersistentFlags().StringSliceVarP(&sitemaps, "sitemaps", "s", []string{}, "sitemap url: https://xxx.com/sitemap.xml")
	syncCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "")
	cmd.RootCommand.AddCommand(syncCmd)
}
