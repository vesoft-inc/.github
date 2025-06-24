package api

import (
	"context"

	"github.com/google/go-github/v70/github"
)

type Client struct {
	client *github.Client
	repo   string
	org    string
}

func NewClient(token string, org string, repo string) *Client {
	c := github.NewClient(nil).WithAuthToken(token)
	return &Client{
		client: c,
		repo:   repo,
		org:    org,
	}
}

func (c *Client) ListIssues(ops *github.IssueListByRepoOptions) ([]*github.Issue, error) {
	issues, _, err := c.client.Issues.ListByRepo(context.Background(), c.org, c.repo, ops)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) ListPullRequest(ops *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	prs, _, err := c.client.PullRequests.List(context.Background(), c.org, c.repo, ops)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (c *Client) ListIssueEvents(number int) ([]*github.IssueEvent, error) {
	events, _, err := c.client.Issues.ListIssueEvents(context.Background(), c.org, c.repo, number, nil)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (c *Client) ListPrComments(number int) ([]*github.IssueComment, error) {
	comments, _, err := c.client.Issues.ListComments(context.Background(), c.org, c.repo, number, nil)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) GetRelease(tag string) (*github.RepositoryRelease, error) {
	release, _, err := c.client.Repositories.GetReleaseByTag(context.Background(), c.org, c.repo, tag)
	if err != nil {
		return nil, err
	}
	return release, nil
}
