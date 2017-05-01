package retriever

import (
	"time"

	"github.com/google/go-github/github"
)

func Timer() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			// TODO: Shit goes here.
		}
	}()
}

var issueID = 0

const ISSUE_QUERY = `
    SELECT * FROM issue
        LEFT JOIN user ON issue.user_id=user.id
        LEFT JOIN labels ON issue.id=labels.issue_fk
        LEFT JOIN user assignee ON issue.assignee_id=assignee.id
        LEFT JOIN user closedby ON issue.closed_by_id=closedby.id
        LEFT JOIN milestones ON issue.milestone_id=milestones.id
        LEFT JOIN pull_request_links ON issue.pull_request_links_id=pull_request_links.id
        LEFT JOIN repository ON issue.repository_id=repository.id
        LEFT JOIN reactions ON issue.reactions_id=reactions.id


    WHERE (id > ?);
`

/*
Tables:
- issue
- user
- plan
- permissions
- labels
- milestones
- pull_request_links
- repository
- organization
- license
- license_permissions
- license_conditions
- license_limitations
- reactions
- assignees
- text_match
- matches
- indices
*/

func (d *Database) Read() ([]*github.Issue, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	issues := []*github.Issue{}
	return issues, nil
}
