package controller

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v70/github"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/api"
)

type PullRequestEvent github.PullRequestEvent

type checkController struct {
	funcs    []checkFunc
	client   *api.Client
	owner    string
	repo     string
	isBugFix bool
}

type checkFunc func(pr *github.PullRequest) error

func NewCheckController(token, owner, repo string) *checkController {
	client := api.NewClient(token, owner, repo)
	ctl := &checkController{
		owner:  owner,
		repo:   repo,
		client: client,
	}
	ctl.funcs = []checkFunc{
		ctl.checkLabels,
		ctl.checkBugLink,
	}
	return ctl
}

func (c *checkController) Check(eventData []byte) error {
	// Parse the event data into a PullRequestEvent
	var event PullRequestEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return err
	}
	if event.PullRequest == nil {
		return fmt.Errorf("pull request is nil")
	}
	for _, fn := range c.funcs {
		if err := fn(event.PullRequest); err != nil {
			return err
		}
	}

	return nil
}

func (c *checkController) checkLabels(pr *github.PullRequest) error {
	bug := "pr/bugfix"
	others := []string{"pr/feature", "pr/cleanup", "pr/improvement", "pr/test", "pr/docs"}
	validMap := make(map[string]struct{})
	validMap[bug] = struct{}{}
	for _, label := range others {
		validMap[label] = struct{}{}
	}
	if pr.Labels == nil || len(pr.Labels) == 0 {
		fmt.Errorf("pull request must have one of the following labels: %v", append([]string{bug}, others...))
	}
	valid := false
	for _, label := range pr.Labels {
		if label == nil {
			continue
		}
		if _, ok := validMap[label.GetName()]; ok {
			valid = true
			if label.GetName() == bug {
				c.isBugFix = true
			}
		}
	}
	if valid {
		return nil
	}

	return fmt.Errorf("pull request must have one of the following labels: %v", append([]string{bug}, others...))
}

func (c *checkController) checkBugLink(pr *github.PullRequest) error {
	if !c.isBugFix {
		return nil
	}
	issues := findIssueNumber(*pr.Body, c.owner, c.repo)
	if len(issues) != 0 {
		fmt.Printf("found issue number %d in pull request body\n", issues[0])
		return nil
	}
	comments, err := c.client.ListPrComments(pr.GetNumber())
	if err != nil {
		return err
	}
	for _, comment := range comments {
		body := comment.GetBody()
		issues := findIssueNumber(body, c.owner, c.repo)
		if len(issues) != 0 {
			fmt.Printf("found issue number %d in comment\n", issues[0])
			return nil
		}
	}
	return fmt.Errorf("pull request must reference a bug issue, please add a link to the issue in the body or comments, e.g. fix https://github.com/vesoft-inc/nebula-ng/issues/8001")
}
