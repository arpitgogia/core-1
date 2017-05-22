package retriever

import (
	"encoding/json"

	"github.com/google/go-github/github"
)

var issueID = 0

const ISSUE_QUERY = `SELECT id, is_pr, payload FROM github_events WHERE id > ?`

type Reader interface {
    Read() (map[int][]*github.Issue, map[int][]*github.PullRequest, map[int][]*github.Issue, error)
}

func (d *Database) Read() (map[int][]*github.Issue, map[int][]*github.PullRequest, map[int][]*github.Issue, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer results.Close()

	issues := make(map[int][]*github.Issue)
	pulls := make(map[int][]*github.PullRequest)
	open := make(map[int][]*github.Issue)

	for results.Next() {
		count := new(int)
		is_pr := new(bool)
		payload := new(string)
		if err := results.Scan(count, is_pr, payload); err != nil {
			return nil, nil, nil, err
		}
		if !*is_pr {
			i := &github.Issue{}
			if err := json.Unmarshal([]byte(*payload), i); err != nil {
				return nil, nil, nil, err
			}
			if i.ClosedAt == nil {
				if _, ok := open[*i.ID]; ok {
					open[*i.ID] = append(open[*i.ID], i)
				} else {
					open[*i.ID] = []*github.Issue{i}
				}
			} else {
				if _, ok := issues[*i.ID]; ok {
					issues[*i.ID] = append(issues[*i.ID], i)
				} else {
					issues[*i.ID] = []*github.Issue{i}
				}
			}
		} else {
			pr := &github.PullRequest{}
			if err := json.Unmarshal([]byte(*payload), pr); err != nil {
				return nil, nil, nil, err
			}
			if _, ok := pulls[*pr.ID]; ok {
				pulls[*pr.ID] = append(pulls[*pr.ID], pr)
			} else {
				pulls[*pr.ID] = []*github.PullRequest{pr}
			}
		}
		issueID = *count
	}
	return issues, pulls, open, nil
}
