package controller

import (
	"fmt"
	"sort"

	"github.com/google/go-github/v70/github"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/api"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/notice"
)

type issue struct {
	assignee string
	count    int
}

type SummaryController struct {
	client     *api.Client
	labels     []string
	severities []string
	dingding   *notice.DingDing
}

func NewSummaryController(token string, labels []string, severities []string,
	dingdingSecret, dingdingToken string) *SummaryController {
	client := api.NewClient(token, "vesoft-inc", "nebula-ng")
	return &SummaryController{
		client:     client,
		labels:     labels,
		severities: severities,
		dingding:   notice.NewDingdinng(dingdingSecret, dingdingToken),
	}
}

func (c *SummaryController) Summary(send bool) (string, error) {
	var m map[string]int = map[string]int{}
	var issues []*github.Issue
	for _, s := range c.severities {
		ll := append(c.labels, s)
		opts := &github.IssueListByRepoOptions{
			State:  "open",
			Labels: ll,
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}
		iss, err := c.client.ListIssues(opts)
		if err != nil {
			return "", err
		}
		issues = append(issues, iss...)
	}
	githubIssues := c.processIssues(issues, m)
	var l []issue
	for k, v := range m {
		l = append(l,
			issue{assignee: k, count: v})
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].count > l[j].count
	})
	s := c.printAssignee(githubIssues, l)
	if send {
		if err := c.dingding.SendDingTalk("nebula-ng issues assignee", "## Nebula-ng issues assignee\n\n"+s); err != nil {
			return "", err
		}
	}
	return s, nil
}

func (c *SummaryController) printAssignee(issues []*github.Issue, l []issue) string {
	var s string
	s += "| Assignee | Count |\n"
	s += "| --- | --- |\n"
	for _, i := range l {
		s += fmt.Sprintf("| %s | %d |\n", i.assignee, i.count)
	}
	s += fmt.Sprintf("| Total | %d |\n", len(issues))
	return s
}

func (c *SummaryController) processIssues(issues []*github.Issue, m map[string]int) []*github.Issue {
	var githubIssues []*github.Issue
	for _, issue := range issues {
		if issue.Type == nil || issue.Type.Name == nil {
			continue
		}
		if *issue.Type.Name != "Bug" {
			continue
		}
		a := issue.Assignees
		if len(a) == 0 {
			m["unassigned"]++
		} else {
			for _, assignee := range a {
				m[*assignee.Login]++
			}
		}
		githubIssues = append(githubIssues, issue)
	}
	return githubIssues
}
