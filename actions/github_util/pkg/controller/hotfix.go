package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v70/github"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/api"
	"github.com/vesoft-inc/nebula-utils/github_util/pkg/export"
	"golang.org/x/sync/errgroup"
)

type HotfixController struct {
	client         *api.Client
	tag            string
	since          *time.Time
	branch         string
	allPr          []*github.PullRequest
	allClosedIssue []*github.Issue
	mapping        *issueMapping
}

type issueMapping struct {
	client         *api.Client
	commitPrMap    map[string]*github.PullRequest
	prMap          map[*github.PullRequest]struct{}
	issueCommitMap map[int][]string
	fixedIssues    []*github.Issue
	noIssuePr      []*github.PullRequest
}

func NewHotfixController(tag, branch string, token string) *HotfixController {
	client := api.NewClient(token, "vesoft-inc", "nebula-ng")
	return &HotfixController{
		client: client,
		tag:    tag,
		branch: branch,
	}
}

// Run is the main entry of the hotfix controller
// Get all the fixed issues from last tag, and get all issue events
// to find the commit id of the issue, then find the pull request in the branch
func (h *HotfixController) Run() error {
	// get since time
	if err := h.getSinceTime(); err != nil {
		return err
	}
	var eg errgroup.Group
	eg.Go(h.getAllClosedIssue)
	eg.Go(h.getAllPr)
	if err := eg.Wait(); err != nil {
		return err
	}
	m := &issueMapping{
		client: h.client,
	}
	if err := m.construct(h.allClosedIssue, h.allPr); err != nil {
		return err
	}
	h.mapping = m
	return nil
}

func (h *HotfixController) Export(filename string) error {
	fixedIssueFile := filename + "_fixed_issue.csv"
	noIssuePrFile := filename + "_no_issue_pr.csv"
	issueIter := export.NewIssueIterator(h.mapping.fixedIssues)
	if err := export.ExportToCsv(issueIter, fixedIssueFile); err != nil {
		return err
	}

	prIter := export.NewPrIterator(h.mapping.noIssuePr)
	if err := export.ExportToCsv(prIter, noIssuePrFile); err != nil {
		return err
	}
	return nil
}

func (h *HotfixController) getSinceTime() error {
	r, err := h.client.GetRelease(h.tag)
	if err != nil {
		return err
	}
	t := r.GetPublishedAt()
	h.since = t.GetTime()
	return nil
}

func (h *HotfixController) getAllClosedIssue() error {
	ops := &github.IssueListByRepoOptions{
		State:  "closed",
		Labels: []string{"type/bug"},
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}
	for i := 1; i < 5; i++ {
		ops.ListOptions.Page = i
		issues, err := h.client.ListIssues(ops)
		if err != nil {
			return err
		}
		for _, issue := range issues {
			if !issue.GetClosedAt().Before(*h.since) {
				h.allClosedIssue = append(h.allClosedIssue, issue)
			}
		}
	}
	return nil
}

func (h *HotfixController) getAllPr() error {
	ops := &github.PullRequestListOptions{
		State: "closed",
		Base:  h.branch,
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}
	for i := 1; i < 5; i++ {
		ops.ListOptions.Page = i
		prs, err := h.client.ListPullRequest(ops)
		if err != nil {
			return err
		}
		for _, pr := range prs {
			if pr.GetClosedAt().Before(*h.since) || pr.MergedAt == nil {
				continue
			}
			noNeedTesting := false
			for _, label := range pr.Labels {
				if label.GetName() == "no-need-testing" {
					noNeedTesting = true
					break
				}
			}
			if !noNeedTesting {
				h.allPr = append(h.allPr, pr)
			}
		}
	}
	return nil
}

func (h *HotfixController) GetFixedIssues() []*github.Issue {
	if h.mapping == nil || h.mapping.fixedIssues == nil {
		return nil
	}
	return h.mapping.fixedIssues
}

func (h *HotfixController) GetNoIssuePr() []*github.PullRequest {
	if h.mapping == nil || h.mapping.noIssuePr == nil {
		return nil
	}
	return h.mapping.noIssuePr
}

func (m *issueMapping) construct(issues []*github.Issue, prs []*github.PullRequest) error {
	m.issueCommitMap = make(map[int][]string)
	m.commitPrMap = make(map[string]*github.PullRequest)
	m.prMap = make(map[*github.PullRequest]struct{})
	for _, issue := range issues {
		events, err := m.client.ListIssueEvents(*issue.Number)
		if err != nil {
			return err
		}
		for _, event := range events {
			if event.GetEvent() == "referenced" || event.GetEvent() == "mentioned" {
				if event.CommitID == nil {
					continue
				}
				m.issueCommitMap[*issue.Number] = append(m.issueCommitMap[*issue.Number], *event.CommitID)
			}
		}
	}

	for _, pr := range prs {
		commit := pr.GetMergeCommitSHA()
		if commit == "" {
			continue
		}
		m.commitPrMap[commit] = pr
		m.prMap[pr] = struct{}{}
	}
	m.fixedIssues = make([]*github.Issue, 0)
	m.noIssuePr = make([]*github.PullRequest, 0)
	issueNumbers := make(map[int]struct{})
	for number, commits := range m.issueCommitMap {
		for _, commit := range commits {
			if pr, ok := m.commitPrMap[commit]; ok {
				issueNumbers[number] = struct{}{}
				delete(m.prMap, pr)
			}
		}
	}

	// if there's no referenced issue
	// need to check pr comments to find the issue number
	for pr := range m.prMap {
		comments, err := m.client.ListPrComments(pr.GetNumber())
		if err != nil {
			return err
		}
		for _, comment := range comments {
			body := comment.GetBody()
			issues := findIssueNumber(body, "vesoft-inc", "nebula-ng")
			if len(issues) == 0 {
				continue
			}
			for _, number := range issues {
				issueNumbers[number] = struct{}{}
			}
			delete(m.prMap, pr)
		}
	}

	for pr := range m.prMap {
		m.noIssuePr = append(m.noIssuePr, pr)
	}

	for _, issue := range issues {
		if _, ok := issueNumbers[*issue.Number]; ok {
			m.fixedIssues = append(m.fixedIssues, issue)
		}
	}

	return nil
}

var issueKeywords = []string{"fix", "close", "resolve", "fixed", "closed", "resolved"}

func findIssueNumber(body string, owner, repoName string) []int {
	var issueNumbers = make([]int, 0)
	for _, line := range strings.Split(body, "\n") {
		n := findIssueNumberInLine(line, owner, repoName)
		if n != -1 {
			issueNumbers = append(issueNumbers, n)
		}
	}
	return issueNumbers
}

func findIssueNumberInLine(line string, owner, repoName string) int {
	prefix := fmt.Sprintf("%s/%s/%s/issues/", "https://github.com", owner, repoName)
	prefix = strings.ToLower(prefix)
	line = strings.ToLower(line)
	for _, kw := range issueKeywords {
		if !strings.Contains(line, kw) {
			continue
		}
		words := strings.Fields(line)
		for _, word := range words {
			if !strings.HasPrefix(strings.ToLower(word), prefix) {
				continue
			}
			number := word[len(prefix):]
			n, err := strconv.Atoi(number)
			if err != nil {
				fmt.Println("invalid issue number ", number)
				break
			}
			return n
		}
	}
	return -1
}
