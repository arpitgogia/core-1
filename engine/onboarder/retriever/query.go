package retriever

import (
	"encoding/json"
	"fmt" // TEMPORARY
	// "runtime"

	"github.com/google/go-github/github"
)

var issueID = 0

const ISSUE_QUERY = `SELECT id, is_pull, payload FROM github_events WHERE id > ?`

type RepoData struct {
	RepoID int
	Open   []*github.Issue
	Closed []*github.Issue
	Pulls  []*github.PullRequest
}

func (m *MemSQL) Read() (map[int]*RepoData, error) {
	results, err := m.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	repodata := make(map[int]*RepoData)

	for results.Next() {
		count := new(int)
		is_pull := new(bool)
		payload := new(string)
		if err := results.Scan(count, is_pull, payload); err != nil {
			return nil, err
		}
        fmt.Println(*payload) // TEMPORARY
		if !*is_pull {
			i := github.Issue{}
			// _, fn, ln := runtime.Caller(1)
			if err := json.Unmarshal([]byte(*payload), &i); err != nil {
				fmt.Println("ERROR IS HERE")
				return nil, err
			}
            fmt.Println(i) // TEMPORARY
			repodata[*i.ID].RepoID = *i.ID
			if i.ClosedAt == nil {
				if _, ok := repodata[*i.ID]; ok {
					repodata[*i.ID].Open = append(repodata[*i.ID].Open, &i)
				} else {
					repodata[*i.ID].Open = []*github.Issue{&i}
				}
			} else {
				if _, ok := repodata[*i.ID]; ok {
					repodata[*i.ID].Closed = append(repodata[*i.ID].Closed, &i)
				} else {
					repodata[*i.ID].Closed = []*github.Issue{&i}
				}
			}
		} else {
			pr := github.PullRequest{}
			if err := json.Unmarshal([]byte(*payload), &pr); err != nil {
				return nil, err
			}
			if _, ok := repodata[*pr.ID]; ok {
				repodata[*pr.ID].Pulls = append(repodata[*pr.ID].Pulls, &pr)
			} else {
				repodata[*pr.ID].Pulls = []*github.PullRequest{&pr}
			}
		}
		issueID = *count
	}
	return repodata, nil
}
