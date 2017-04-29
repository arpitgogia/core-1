package retriever

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
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

const ISSUE_QUERY = `SELECT * FROM issues WHERE (id > ?);`

/*
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
// Query: All values for current date starting from current time - 500ms
*/

// TODO: Build to support a PR query as well.
func (d *Database) Read() ([]*github.Issue, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	id := 0
	issues := []*github.Issue{}
	issue := new(github.Issue)

	for results.Next() {
		if err := results.Scan(&id, &issue.ID, &issue.Number, &issue.State,
			&issue.Locked, &issue.Title, &issue.Body, &issue.Comments,
			&issue.ClosedAt, &issue.CreatedAt, &issue.UpdatedAt, &issue.URL,
			&issue.HTMLURL,
		); err != nil {
			return nil, err
		}
		/*if err := results.Scan(&id, &issue.ID, &issue.Number, &issue.State, &issue.Locked, // NOTE: Issue struct from MemSQL
		                              &issue.Title, &issue.Body, &issue.Comments, &issue.ClosedAt,
		                              &issue.CreatedAt, &issue.UpdatedAt, &issue.URL, &issue.HTMLURL,
		                          &userID, &issue.User.ID, &issue.User.AvatarURL, &issue.User.HTMLURL, // NOTE: User field on Issue struct
		                              &issue.User.GravatarID, &issue.User.Name, &issue.User.Company,
		                              &issue.User.Blog, &issue.User.Location, &issue.User.Email,
		                              &issue.User.Hireable, &issue.User.Bio, &issue.User.PublicRepos,
		                              &issue.User.PublicGists, &issue.User.Followers, &issue.User.Following,
		                              &issue.User.CreatedAt, &issue.User.UpdatedAt, &issue.User.SuspendedAt,
		                              &issue.User.Type, &issue.User.SiteAdmin, &issue.User.TotalPrivateRepos,
		                              &issue.User.OwnedPrivateRepos, &issue.User.PrivateGists,
		                              &issue.User.DiskUsage, &issue.User.Collaborators,
		                              &issue.User.Plan.Name, &issue.User.Plan.Space, // NOTE: Issue User Plan info
		                                  &issue.User.Plan.Collaborators, &issue.User.Plan.PrivateRepos,
		                                  &issue.User.URL, &issue.User.EventsURL, &issue.User.FollowingURL,
		                                  &issue.User.FollowersURL, &issue.User.GistsURL, &issue.User.OrganizationsURL,
		                                  &issue.User.ReceivedEventsURL, &issue.User.ReposURL, &issue.User.StarredURL,
		                                  &issue.User.SubscriptionsURL,
		                          &userID, &issue.Assignee.ID, &issue.Assignee.AvatarURL, &issue.Assignee.HTMLURL, // NOTE: Assignee field on Issue struct
		                              &issue.Assignee.GravatarID, &issue.Assignee.Name, &issue.Assignee.Company,
		                              &issue.Assignee.Blog, &issue.Assignee.Location, &issue.Assignee.Email,
		                              &issue.Assignee.Hireable, &issue.Assignee.Bio, &issue.Assignee.PublicRepos,
		                              &issue.Assignee.PublicGists, &issue.Assignee.Followers, &issue.Assignee.Following,
		                              &issue.Assignee.CreatedAt, &issue.Assignee.UpdatedAt, &issue.Assignee.SuspendedAt,
		                              &issue.Assignee.Type, &issue.Assignee.SiteAdmin, &issue.Assignee.TotalPrivateRepos,
		                              &issue.Assignee.OwnedPrivateRepos, &issue.Assignee.PrivateGists,
		                              &issue.Assignee.DiskUsage, &issue.Assignee.Collaborators,
		                              &issue.Assignee.Plan.Name, &issue.Assignee.Plan.Space, // NOTE: Issue Assignee Plan info
		                                  &issue.Assignee.Plan.Collaborators, &issue.Assignee.Plan.PrivateRepos,
		                                  &issue.Assignee.URL, &issue.Assignee.EventsURL, &issue.Assignee.FollowingURL,
		                                  &issue.Assignee.FollowersURL, &issue.Assignee.GistsURL, &issue.Assignee.OrganizationsURL,
		                                  &issue.Assignee.ReceivedEventsURL, &issue.Assignee.ReposURL, &issue.Assignee.StarredURL,
		                                  &issue.Assignee.SubscriptionsURL,
		                          &issue.Comments, &issue.ClosedAt, &issue.CreatedAt, &issue.UpdatedAt, // NOTE: General fields on User
		                          &userID, &issue.ClosedBy.ID, &issue.ClosedBy.AvatarURL, &issue.ClosedBy.HTMLURL, // NOTE: ClosedBy field on Issue struct
		                              &issue.ClosedBy.GravatarID, &issue.ClosedBy.Name, &issue.ClosedBy.Company,
		                              &issue.ClosedBy.Blog, &issue.ClosedBy.Location, &issue.ClosedBy.Email,
		                              &issue.ClosedBy.Hireable, &issue.ClosedBy.Bio, &issue.ClosedBy.PublicRepos,
		                              &issue.ClosedBy.PublicGists, &issue.ClosedBy.Followers, &issue.ClosedBy.Following,
		                              &issue.ClosedBy.CreatedAt, &issue.ClosedBy.UpdatedAt, &issue.ClosedBy.SuspendedAt,
		                              &issue.ClosedBy.Type, &issue.ClosedBy.SiteAdmin, &issue.ClosedBy.TotalPrivateRepos,
		                              &issue.ClosedBy.OwnedPrivateRepos, &issue.ClosedBy.PrivateGists,
		                              &issue.ClosedBy.DiskUsage, &issue.ClosedBy.Collaborators,
		                              &issue.ClosedBy.Plan.Name, &issue.ClosedBy.Plan.Space, // NOTE: Issue ClosedBy Plan info
		                                  &issue.ClosedBy.Plan.Collaborators, &issue.ClosedBy.Plan.PrivateRepos,
		                                  &issue.ClosedBy.URL, &issue.ClosedBy.EventsURL, &issue.ClosedBy.FollowingURL,
		                                  &issue.ClosedBy.FollowersURL, &issue.ClosedBy.GistsURL, &issue.ClosedBy.OrganizationsURL,
		                                  &issue.ClosedBy.ReceivedEventsURL, &issue.ClosedBy.ReposURL, &issue.ClosedBy.StarredURL,
		                                  &issue.ClosedBy.SubscriptionsURL,
		                          &issue.URL, &issue.HTMLURL, // NOTE: Link fields for Issue struct
		                          &issue.Milestone.URL, &issue.Milestone.HTMLURL, &issue.Milestone.LabelsURL, // NOTE: Milestone field on Issue
		                              &issue.Milestone.ID, &issue.Milestone.Number, &issue.Milestone.State,
		                              &issue.Milestone.Title, &issue.Milestone.Description,
		                              &userID, &issue.Milestone.Creator.ID, &issue.Milestone.Creator.AvatarURL, &issue.Milestone.Creator.HTMLURL, // NOTE: Creator field on Milestone struct
		                                  &issue.Milestone.Creator.GravatarID, &issue.Milestone.Creator.Name, &issue.Milestone.Creator.Company,
		                                  &issue.Milestone.Creator.Blog, &issue.Milestone.Creator.Location, &issue.Milestone.Creator.Email,
		                                  &issue.Milestone.Creator.Hireable, &issue.Milestone.Creator.Bio, &issue.Milestone.Creator.PublicRepos,
		                                  &issue.Milestone.Creator.PublicGists, &issue.Milestone.Creator.Followers, &issue.Milestone.Creator.Following,
		                                  &issue.Milestone.Creator.CreatedAt, &issue.Milestone.Creator.UpdatedAt, &issue.Milestone.Creator.SuspendedAt,
		                                  &issue.Milestone.Creator.Type, &issue.Milestone.Creator.SiteAdmin, &issue.Milestone.Creator.TotalPrivateRepos,
		                                  &issue.Milestone.Creator.OwnedPrivateRepos, &issue.Milestone.Creator.PrivateGists,
		                                  &issue.Milestone.Creator.DiskUsage, &issue.Milestone.Creator.Collaborators,
		                                  &issue.Milestone.Creator.Plan.Name, &issue.Milestone.Creator.Plan.Space, // NOTE: Milestone Creator Plan info
		                                      &issue.Milestone.Creator.Plan.Collaborators, &issue.Milestone.Creator.Plan.PrivateRepos,
		                                      &issue.Milestone.Creator.URL, &issue.Milestone.Creator.EventsURL, &issue.Milestone.Creator.FollowingURL,
		                                      &issue.Milestone.Creator.FollowersURL, &issue.Milestone.Creator.GistsURL, &issue.Milestone.Creator.OrganizationsURL,
		                                      &issue.Milestone.Creator.ReceivedEventsURL, &issue.Milestone.Creator.ReposURL, &issue.Milestone.Creator.StarredURL,
		                                      &issue.Milestone.Creator.SubscriptionsURL,
		                              &issue.Milestone.OpenIssues, &issue.Milestone.ClosedIssues, &issue.Milestone.CreatedAt,
		                              &issue.Milestone.UpdatedAt, &issue.Milestone.ClosedAt, &issue.Milestone.DueOn,
		                          &issue.PullRequestLinks.URL, &issue.PullRequestLinks.HTMLURL, // NOTE: PullRequestLinks struct on Issue
		                          &issue.PullRequestLinks.DiffURL, &issue.PullRequestLinks.PatchURL,
		                          &issue.Repository.ID, // NOTE: Issue struct Repository field
		                          &userID, &issue.Repository.Owner.ID, &issue.Repository.Owner.AvatarURL, &issue.Repository.Owner.HTMLURL, // NOTE: Owner field on Repository struct
		                              &issue.Repository.Owner.GravatarID, &issue.Repository.Owner.Name, &issue.Repository.Owner.Company,
		                              &issue.Repository.Owner.Blog, &issue.Repository.Owner.Location, &issue.Repository.Owner.Email,
		                              &issue.Repository.Owner.Hireable, &issue.Repository.Owner.Bio, &issue.Repository.Owner.PublicRepos,
		                              &issue.Repository.Owner.PublicGists, &issue.Repository.Owner.Followers, &issue.Repository.Owner.Following,
		                              &issue.Repository.Owner.CreatedAt, &issue.Repository.Owner.UpdatedAt, &issue.Repository.Owner.SuspendedAt,
		                              &issue.Repository.Owner.Type, &issue.Repository.Owner.SiteAdmin, &issue.Repository.Owner.TotalPrivateRepos,
		                              &issue.Repository.Owner.OwnedPrivateRepos, &issue.Repository.Owner.PrivateGists,
		                              &issue.Repository.Owner.DiskUsage, &issue.Repository.Owner.Collaborators,
		                              &issue.Repository.Owner.Plan.Name, &issue.Repository.Owner.Plan.Space, // NOTE: Repository Owner Plan info
		                                  &issue.Repository.Owner.Plan.Collaborators, &issue.Repository.Owner.Plan.PrivateRepos,
		                                  &issue.Repository.Owner.URL, &issue.Repository.Owner.EventsURL, &issue.Repository.Owner.FollowingURL,
		                                  &issue.Repository.Owner.FollowersURL, &issue.Repository.Owner.GistsURL, &issue.Repository.Owner.OrganizationsURL,
		                                  &issue.Repository.Owner.ReceivedEventsURL, &issue.Repository.Owner.ReposURL, &issue.Repository.Owner.StarredURL,
		                                  &issue.Repository.Owner.SubscriptionsURL,
		                              &issue.Repository.Name, &issue.Repository.FullName, &issue.Repository.Description,
		                              &issue.Repository.Homepage, &issue.Repository.DefaultBranch, &issue.Repository.MasterBranch,
		                              &issue.Repository.CreatedAt, &issue.Repository.PushedAt, &issue.Repository.UpdatedAt,
		                              &issue.Repository.HTMLURL, &issue.Repository.CloneURL, &issue.Repository.GitURL,
		                              &issue.Repository.MirrorURL,


		  ); err != nil {
		      return nil, err
		  }*/
		issues = append(issues, issue)
		issueID = *issue.ID
	}
	// TODO: Populating entire Issue struct (all 1:1 relationships).
	// - This requires looking into

	return issues, nil
}
