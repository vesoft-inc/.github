package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "github_util",
	Short: "github_util is a tool to help manage github issues and pull requests",
	Long:  `github_util is a tool to help manage github issues and pull requests.`,
}
