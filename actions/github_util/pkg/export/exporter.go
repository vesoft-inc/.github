package export

import (
	"encoding/csv"
	"os"

	"github.com/google/go-github/v70/github"
)

type Iterator interface {
	Header() []string
	Next() []string
	Rows() int
}

type prIterator struct {
	prs []*github.PullRequest
}

type issueIterator struct {
	issues []*github.Issue
}

func ExportToCsv(iter Iterator, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := iter.Header()
	if err := writer.Write(header); err != nil {
		return err
	}
	for i := 0; i < iter.Rows(); i++ {
		row := iter.Next()
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	// if err := writer.Error(); err != nil {
	// 	return err
	// }
	return nil
}

func NewPrIterator(prs []*github.PullRequest) Iterator {
	return &prIterator{prs: prs}
}

func NewIssueIterator(issues []*github.Issue) Iterator {
	return &issueIterator{issues: issues}
}

func (i *prIterator) Header() []string {
	return []string{"title", "url", "created_at", "closed_at", "merged_at", "author", "assignee"}
}

func (i *prIterator) Next() []string {
	pr := i.prs[0]
	i.prs = i.prs[1:]
	return []string{
		pr.GetTitle(),
		pr.GetHTMLURL(),
		pr.GetCreatedAt().String(),
		pr.GetClosedAt().String(),
		pr.GetMergedAt().String(),
		pr.GetUser().GetLogin(),
		pr.GetAssignee().GetLogin(),
	}
}
func (i *prIterator) Rows() int {
	return len(i.prs)
}

func (i *issueIterator) Header() []string {
	return []string{"title", "url", "created_at", "closed_at", "author", "assignee"}
}

func (i *issueIterator) Next() []string {
	issue := i.issues[0]
	i.issues = i.issues[1:]
	return []string{
		issue.GetTitle(),
		issue.GetHTMLURL(),
		issue.GetCreatedAt().String(),
		issue.GetClosedAt().String(),
		issue.GetUser().GetLogin(),
	}
}

func (i *issueIterator) Rows() int {
	return len(i.issues)
}
