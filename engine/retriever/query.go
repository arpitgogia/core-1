package retriever

import (
    "database/sql"
    "encoding/json"
    "fmt" // TEMPORARY
    "time"

	"github.com/google/go-github/github"
)

func (d *Database) Timer() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
            issues, pullRequests, err := d.Read()
            if err != nil {
                fmt.Println(err) // TEMPORARY
            }
            fmt.Println(issues) // TEMPORARY
            fmt.Println(pullRequests) // TEMPORARY
		}
	}()
}

var issueID = 0

const ISSUE_QUERY = `SELECT id, is_pr, payload FROM github_events WHERE id > ?`

func (d *Database) Read() ([]*github.Issue, []*github.PullRequest, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, nil, err
	}
	defer results.Close()

    i, pr, err := translator(results)
    if err != nil {
        return nil, nil, err
    }
	return i, pr, nil
}

func translator(rows *sql.Rows) ([]*github.Issue, []*github.PullRequest, error) {
    i_list := []*github.Issue{}
    pr_list := []*github.PullRequest{}
    for rows.Next() {
        count := new(int)
        is_pr := new(bool)
        payload := new(string)
        if err := rows.Scan(count, is_pr, payload); err != nil {
            return nil, nil, err
        }
        // TODO: Refactor this logic into separate function.
        if !*is_pr {
            i := github.Issue{}
            if err := json.Unmarshal([]byte(*payload), i); err != nil {
                return nil, nil, err
            }
            i_list = append(i_list, &i)
            issueID = *count
        } else {
            pr := github.PullRequest{}
            if err := json.Unmarshal([]byte(*payload), pr); err != nil {
                return nil, nil, err
            }
            pr_list = append(pr_list, &pr)
            issueID = *count
        }
    }
    return i_list, pr_list, nil
}
