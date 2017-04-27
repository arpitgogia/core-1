package retriever

import (
    "database/sql"
    "time"

    "github.com/google/go-github/github"
    _ "github.com/go-sql-driver/mysql"
)

// TODO: Refactor this out into a parent file for shared access.
type Database struct {
	db *sql.DB
}

func (d *Database) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	d.db = mysql
}

func (d *Database) Close() {
	d.db.Close()
}

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
    SELECT * FROM issues
        LEFT JOIN users ON issues.id=users.issue_fk
        LEFT JOIN assignees ON issues.id=assignee.issue_fk
        LEFT JOIN closers ON issues.id=closer.issue_fk
        LEFT JOIN milestones ON isues.id=milestones.issue_fk
        LEFT JOIN pull_request_links ON issues.id=pull_request_links.issue_fk
        LEFT JOIN repos ON issues.id=repos.issue_fk
        LEFT JOIN reactions ON issues.id=reactions.issue_fk
    WHERE (id > ?);
`

// TODO: Build to support a PR query as well.
func (d *Database) Read() ([]*github.Issue, error) {
    results, err := d.db.Query(ISSUE_QUERY, issueID)
    if err != nil {
        return nil, err
    }
    defer results.Close()

    id := 0
    issues := []*github.Issue{} // NOTE: Final output slice.
    issue := new(github.Issue)

    for results.Next() {
        if err := results.Scan(&id, &issue.ID, &issue.Number, &issue.State, &issue.Locked, &issue.Title, &issue.Body, &issue.Comments, &issue.ClosedAt, &issue.CreatedAt, &issue.UpdatedAt, &issue.URL, &issue.HTMLURL,
                                // THE REST OF THE SCAN GO HERE
        ); err != nil {
            return nil, err
        }
    //     BUILD FULL ISSUE HERE AND THEN APPEND TO OUTPUT LIST
    //     issues = append(issues, issue)
        issueID = *issue.ID
    }

    // Step 0) Hit Up Issues Table and Create Issues Slice
    // Step 1)
    // For each Issue in the Issues Slice
    //    Lookup Assignee ID and Populate Assignee Struct on the Issue struct
    //    Lookup Label ID and Populate Labels Struct on the Issue struct


    // QUERY: (ALL PRIOR STUFF) AND WHERE > LAST ID
    // ^ LAST ID = 0 on startup
    // ^ Set LAST ID = max(id)

    // Query: All values for current date starting from current time - 500ms

    // Reads newest values in db (since last call)
    // Populates results into struct
    // - IssuesEvents (likely)
    // Needed logic:
    // - Handle no/many results from database query
    return issues, nil
}
