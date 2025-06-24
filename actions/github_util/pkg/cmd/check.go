package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/controller"
)

var (
	eventPath string
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check pr",
	RunE: func(cmd *cobra.Command, args []string) error {
		if eventPath == "" {
			eventPath = os.Getenv("GITHUB_EVENT_PATH")
		}
		if eventPath == "" {
			return fmt.Errorf("event path is required, please provide it using --event flag or set GITHUB_EVENT_PATH environment variable")
		}
		eventData, err := os.ReadFile(eventPath)
		if err != nil {
			return err
		}
		repo := os.Getenv("GITHUB_REPOSITORY")
		rr := strings.Split(repo, "/")
		if len(rr) != 2 {
			return fmt.Errorf("GITHUB_REPOSITORY must be in the format 'owner/name', got: %s", repo)
		}
		owner, name := rr[0], rr[1]
		ctl := controller.NewCheckController(token, owner, name)
		if err := ctl.Check(eventData); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&eventPath, "event", "e", "", "Path to the event file")
	checkCmd.Flags().StringVarP(&token, "token", "t", "", "GitHub token for authentication")
}
