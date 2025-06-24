package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/controller"
)

var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "hotfix is a tool to help manage hotfix",
	Long:  `hotfix is a tool to help manage hotfix.`,
}

var (
	tag      string
	branch   string
	token    string
	filename string
)

var exportHotfixCmd = &cobra.Command{
	Use:   "export",
	Short: "export is a tool to help export hotfix",
	Long:  `export is a tool to help export hotfix.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Begin to export hotfix\n")
		start := time.Now()
		ctl := controller.NewHotfixController(tag, branch, token)
		if err := ctl.Run(); err != nil {
			return err
		}
		if err := ctl.Export(filename); err != nil {
			return err
		}
		fmt.Printf("Export hotfix successfully\n")
		fmt.Printf("Time cost: %v\n", time.Since(start))
		return nil
	},
}

func init() {
	hotfixCmd.AddCommand(exportHotfixCmd)
	RootCmd.AddCommand(hotfixCmd)

	exportHotfixCmd.Flags().StringVarP(&tag, "tag", "t", "", "previous tag of the hotfix")
	exportHotfixCmd.Flags().StringVarP(&branch, "branch", "b", "", "branch of the hotfix")
	exportHotfixCmd.Flags().StringVarP(&token, "token", "k", "", "github token")
	exportHotfixCmd.Flags().StringVarP(&filename, "filename", "f", "output", "output filename")
}
