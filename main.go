package main

import (
	"fmt"
	"os"

	"push-seo/cmd"
	_ "push-seo/cmd/sync"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}
