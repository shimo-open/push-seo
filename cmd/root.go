package cmd

import (
	"log"
	"os"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/spf13/cobra"
)

// Config config路径
var config string

func init() {
	confEnv := os.Getenv("EGO_CONFIG_PATH")
	if confEnv == "" {
		confEnv = "config/default.toml"
	}
	RootCommand.PersistentFlags().StringVarP(&config, "config", "c", confEnv, "指定配置文件，默认 config/default.toml")
}

var RootCommand = &cobra.Command{
	Use: "push-seo",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Println("ConfigFile", config)
		provider, parser, tag, err := manager.NewDataSource(config, eflag.Bool("watch"))
		if err != nil {
			log.Fatal("load config fail: ", err)
		}
		if err := econf.LoadFromDataSource(provider, parser, econf.WithSquash(true), econf.WithTagName(tag)); err != nil {
			log.Fatal("data source: load config, unmarshal config err: ", err)
		}
	},
}
