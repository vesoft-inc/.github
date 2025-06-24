package main

import (
	"os"

	"github.com/vesoft-inc/nebula-utils/github_util/pkg/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
