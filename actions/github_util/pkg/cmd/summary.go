package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/controller"
)

var (
	labels         string
	severities     string
	send           bool
	dingdingSecret string
	dingdingToken  string
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "summary is a tool to calculate the assignee of the issue",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var ll, ss []string
		if labels != "" {
			ll = strings.Split(labels, ",")
		}
		if severities != "" {
			ss = strings.Split(severities, ",")
		}
		c := controller.NewSummaryController(token, ll, ss, dingdingSecret, dingdingToken)
		s, err := c.Summary(send)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().StringVarP(&labels, "labels", "l", "", "required labels")
	summaryCmd.Flags().StringVarP(&severities, "severities", "s", "severity/blocker,severity/major", "issue severities")
	summaryCmd.Flags().StringVarP(&token, "token", "k", "", "github token")
	summaryCmd.Flags().BoolVarP(&send, "send", "e", true, "send to dingding")
	summaryCmd.Flags().StringVarP(&dingdingToken, "dingding-token", "t", "", "dingding token")
	summaryCmd.Flags().StringVarP(&dingdingSecret, "dingding-secret", "c", "", "dingding secret")
}
