package retriever

import (
	"database/sql"
	"encoding/json"

	"github.com/google/go-github/github"
)

var issueID = 0

const ISSUE_QUERY = `SELECT id, is_pr, payload FROM github_events WHERE id > ?`

func (d *Database) Read() ([]*github.Issue, []*github.PullRequest, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, nil, err
	}
	defer results.Close()

	issues, pulls, err := translator(results)
	if err != nil {
		return nil, nil, err
	}
	return issues, pulls, nil
}

func translator(rows *sql.Rows) ([]*github.Issue, []*github.PullRequest, error) {
	issues := []*github.Issue{}
	pulls := []*github.PullRequest{}
	for rows.Next() {
		count := new(int)
		is_pr := new(bool)
		payload := new(string)
		if err := rows.Scan(count, is_pr, payload); err != nil {
			return nil, nil, err
		}
		if !*is_pr {
			i := github.Issue{}
			if err := json.Unmarshal([]byte(*payload), i); err != nil {
				return nil, nil, err
			}
			issues = append(issues, &i)
			issueID = *count
		} else {
			pr := github.PullRequest{}
			if err := json.Unmarshal([]byte(*payload), pr); err != nil {
				return nil, nil, err
			}
			pulls = append(pulls, &pr)
			issueID = *count
		}
	}
	return issues, pulls, nil
}

// TODO: Implement this refactoring (make sure nothing is being shadowed).
func parseJSON(obj *interface{}, pld string, dest *[]interface{}, cnt int) error {
	if err := json.Unmarshal([]byte(pld), obj); err != nil {
		return err
	}
	*dest = append(*dest, &pld)
	issueID = cnt
	return nil
}
