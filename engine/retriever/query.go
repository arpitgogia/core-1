package retriever

import (
    "database/sql"
    "encoding/json"
    "time"

	"github.com/google/go-github/github"
)

func (d *Database) Timer() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			// d.Read()
            // TODO: Wrapper around
		}
	}()
}

var issueID = 0

const ISSUE_QUERY = `SELECT id, is_closed, is_pr, payload FROM github_events WHERE id > ?`

func (d *Database) Read() ([]*github.Issue, []*github.Issue, []*github.PullRequest, []*github.PullRequest, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer results.Close()

    open_i, closed_i, open_pr, closed_pr, err := translator(results)
    if err != nil {
        return nil, nil, nil, nil, err
    }
	return open_i, closed_i, open_pr, closed_pr, nil
}

func translator(rows *sql.Rows) (open_i, closed_i []*github.Issue, open_pr, closed_pr []*github.PullRequest, err error) {
    for rows.Next() {
        payload := new(string)
        is_closed := new(bool)
        is_pr := new(bool)
        count := new(int)
        if err = rows.Scan(count, is_closed, is_pr, payload); err != nil {
            return
        }
        if !*is_pr {
            i := github.Issue{}
            if err = json.Unmarshal([]byte(*payload), i); err != nil {
                return
            }
            if !*is_closed {
                open_i = append(open_i, &i)
            } else {
                closed_i = append(closed_i, &i)
            }
            issueID = *count
        } else {
            pr := github.PullRequest{}
            if err = json.Unmarshal([]byte(*payload), pr); err != nil {
                return
            }
            if !*is_closed {
                open_pr = append(open_pr, &pr)
            } else {
                closed_pr = append(closed_pr, &pr)
            }
            issueID = *count
        }
    }
    return
}
