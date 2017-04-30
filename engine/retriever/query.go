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

const ISSUE_QUERY = `
    SELECT * FROM issue
        LEFT JOIN user ON issue.user_id=user.id
        LEFT JOIN user assignee ON issue.assignee_id=assignee.id
        LEFT JOIN user closedby ON issue.closed_by_id=closedby.id
        LEFT JOIN milestones ON issue.milestone_id=milestones.id
        LEFT JOIN pull_request_links ON issue.pull_request_links_id=pull_request_links.id
        LEFT JOIN repository ON issue.repository_id=repository.id
        LEFT JOIN reactions ON issue.reactions_id=reactions.id
    WHERE (id > ?);
`

// TODO: Build to support a PR query as well.
func (d *Database) Read() ([]*github.Issue, error) {
	results, err := d.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	issues := []*github.Issue{}
	issue := new(github.Issue)
    issueID := 0    // Temporary field - "throwaway"
    userID := 0     // Temporary field - "throwaway"
    parentID := 0   // Temporary field - "throwaway"
    sourcdID := 0   // Temporary field - "throwaway"

	for results.Next() {
		if err := results.Scan(&issueID, // NOTE: Issue struct from MemSQL
                                    &issue.ID,
                                    &issue.Number,
                                    &issue.State,
                                    &issue.Locked,
                                    &issue.Title,
                                    &issue.Body,
                                    &issue.Comments,
                                    &issue.ClosedAt,
                                    &issue.CreatedAt,
                                    &issue.UpdatedAt,
                                    &issue.URL,
                                    &issue.HTMLURL,
		                        &userID, // NOTE: User field on Issue struct
                                    &issue.User.ID,
                                    &issue.User.AvatarURL,
                                    &issue.User.HTMLURL,
		                            &issue.User.GravatarID,
                                    &issue.User.Name,
                                    &issue.User.Company,
		                            &issue.User.Blog,
                                    &issue.User.Location,
                                    &issue.User.Email,
		                            &issue.User.Hireable,
                                    &issue.User.Bio,
                                    &issue.User.PublicRepos,
		                            &issue.User.PublicGists,
                                    &issue.User.Followers,
                                    &issue.User.Following,
                                    &issue.User.CreatedAt,
                                    &issue.User.UpdatedAt,
                                    &issue.User.SuspendedAt,
                                    &issue.User.Type,
                                    &issue.User.SiteAdmin,
                                    &issue.User.TotalPrivateRepos,
                                    &issue.User.OwnedPrivateRepos,
                                    &issue.User.PrivateGists,
                                    &issue.User.DiskUsage,
                                    &issue.User.Collaborators,
                                    &issue.User.Plan.Name, // NOTE: Issue User Plan info
                                        &issue.User.Plan.Space,
                                        &issue.User.Plan.Collaborators,
                                        &issue.User.Plan.PrivateRepos,
		                                &issue.User.URL,
                                        &issue.User.EventsURL,
                                        &issue.User.FollowingURL,
		                                &issue.User.FollowersURL,
                                        &issue.User.GistsURL,
                                        &issue.User.OrganizationsURL,
		                                &issue.User.ReceivedEventsURL,
                                        &issue.User.ReposURL,
                                        &issue.User.StarredURL,
		                                &issue.User.SubscriptionsURL,
		                        &userID, // NOTE: Assignee field on Issue struct
                                    &issue.Assignee.ID,
                                    &issue.Assignee.AvatarURL,
                                    &issue.Assignee.HTMLURL,
		                            &issue.Assignee.GravatarID,
                                    &issue.Assignee.Name,
                                    &issue.Assignee.Company,
		                            &issue.Assignee.Blog,
                                    &issue.Assignee.Location,
                                    &issue.Assignee.Email,
		                            &issue.Assignee.Hireable,
                                    &issue.Assignee.Bio,
                                    &issue.Assignee.PublicRepos,
		                            &issue.Assignee.PublicGists,
                                    &issue.Assignee.Followers,
                                    &issue.Assignee.Following,
		                            &issue.Assignee.CreatedAt,
                                    &issue.Assignee.UpdatedAt,
                                    &issue.Assignee.SuspendedAt,
		                            &issue.Assignee.Type,
                                    &issue.Assignee.SiteAdmin,
                                    &issue.Assignee.TotalPrivateRepos,
		                            &issue.Assignee.OwnedPrivateRepos,
                                    &issue.Assignee.PrivateGists,
		                            &issue.Assignee.DiskUsage,
                                    &issue.Assignee.Collaborators,
		                            &issue.Assignee.Plan.Name, // NOTE: Issue Assignee Plan info
                                        &issue.Assignee.Plan.Space,
		                                &issue.Assignee.Plan.Collaborators,
                                        &issue.Assignee.Plan.PrivateRepos,
		                                &issue.Assignee.URL,
                                        &issue.Assignee.EventsURL,
                                        &issue.Assignee.FollowingURL,
                                        &issue.Assignee.FollowersURL,
                                        &issue.Assignee.GistsURL,
                                        &issue.Assignee.OrganizationsURL,
		                                &issue.Assignee.ReceivedEventsURL,
                                        &issue.Assignee.ReposURL,
                                        &issue.Assignee.StarredURL,
		                                &issue.Assignee.SubscriptionsURL,
		                        &issue.Comments, // NOTE: General fields on User
                                &issue.ClosedAt,
                                &issue.CreatedAt,
                                &issue.UpdatedAt,
		                        &userID, // NOTE: ClosedBy field on Issue struct (note that this could be nil)
                                    &issue.ClosedBy.ID,
                                    &issue.ClosedBy.AvatarURL,
                                    &issue.ClosedBy.HTMLURL,
		                            &issue.ClosedBy.GravatarID,
                                    &issue.ClosedBy.Name,
                                    &issue.ClosedBy.Company,
		                            &issue.ClosedBy.Blog,
                                    &issue.ClosedBy.Location,
                                    &issue.ClosedBy.Email,
		                            &issue.ClosedBy.Hireable,
                                    &issue.ClosedBy.Bio,
                                    &issue.ClosedBy.PublicRepos,
		                            &issue.ClosedBy.PublicGists,
                                    &issue.ClosedBy.Followers,
                                    &issue.ClosedBy.Following,
		                            &issue.ClosedBy.CreatedAt,
                                    &issue.ClosedBy.UpdatedAt,
                                    &issue.ClosedBy.SuspendedAt,
		                            &issue.ClosedBy.Type,
                                    &issue.ClosedBy.SiteAdmin,
                                    &issue.ClosedBy.TotalPrivateRepos,
		                            &issue.ClosedBy.OwnedPrivateRepos,
                                    &issue.ClosedBy.PrivateGists,
		                            &issue.ClosedBy.DiskUsage,
                                    &issue.ClosedBy.Collaborators,
		                            &issue.ClosedBy.Plan.Name, // NOTE: Issue ClosedBy Plan info
                                        &issue.ClosedBy.Plan.Space,
		                                &issue.ClosedBy.Plan.Collaborators,
                                        &issue.ClosedBy.Plan.PrivateRepos,
		                                &issue.ClosedBy.URL,
                                        &issue.ClosedBy.EventsURL,
                                        &issue.ClosedBy.FollowingURL,
		                                &issue.ClosedBy.FollowersURL,
                                        &issue.ClosedBy.GistsURL,
                                        &issue.ClosedBy.OrganizationsURL,
		                                &issue.ClosedBy.ReceivedEventsURL,
                                        &issue.ClosedBy.ReposURL,
                                        &issue.ClosedBy.StarredURL,
		                                &issue.ClosedBy.SubscriptionsURL,
		                        &issue.URL, // NOTE: Link fields for Issue struct
                                    &issue.HTMLURL,
		                        &issue.Milestone.URL, // NOTE: Milestone field on Issue
                                    &issue.Milestone.HTMLURL,
                                    &issue.Milestone.LabelsURL,
		                            &issue.Milestone.ID,
                                    &issue.Milestone.Number,
                                    &issue.Milestone.State,
		                            &issue.Milestone.Title,
                                    &issue.Milestone.Description,
		                            &userID, // NOTE: Creator field on Milestone struct
                                        &issue.Milestone.Creator.ID,
                                        &issue.Milestone.Creator.AvatarURL,
                                        &issue.Milestone.Creator.HTMLURL,
		                                &issue.Milestone.Creator.GravatarID,
                                        &issue.Milestone.Creator.Name,
                                        &issue.Milestone.Creator.Company,
		                                &issue.Milestone.Creator.Blog,
                                        &issue.Milestone.Creator.Location,
                                        &issue.Milestone.Creator.Email,
		                                &issue.Milestone.Creator.Hireable,
                                        &issue.Milestone.Creator.Bio,
                                        &issue.Milestone.Creator.PublicRepos,
		                                &issue.Milestone.Creator.PublicGists,
                                        &issue.Milestone.Creator.Followers,
                                        &issue.Milestone.Creator.Following,
		                                &issue.Milestone.Creator.CreatedAt,
                                        &issue.Milestone.Creator.UpdatedAt,
                                        &issue.Milestone.Creator.SuspendedAt,
		                                &issue.Milestone.Creator.Type,
                                        &issue.Milestone.Creator.SiteAdmin,
                                        &issue.Milestone.Creator.TotalPrivateRepos,
		                                &issue.Milestone.Creator.OwnedPrivateRepos,
                                        &issue.Milestone.Creator.PrivateGists,
		                                &issue.Milestone.Creator.DiskUsage,
                                        &issue.Milestone.Creator.Collaborators,
		                                &issue.Milestone.Creator.Plan.Name, // NOTE: Milestone Creator Plan info
                                            &issue.Milestone.Creator.Plan.Space,
		                                    &issue.Milestone.Creator.Plan.Collaborators,
                                            &issue.Milestone.Creator.Plan.PrivateRepos,
	                                    &issue.Milestone.Creator.URL,
                                        &issue.Milestone.Creator.EventsURL,
                                        &issue.Milestone.Creator.FollowingURL,
		                                &issue.Milestone.Creator.FollowersURL,
                                        &issue.Milestone.Creator.GistsURL,
                                        &issue.Milestone.Creator.OrganizationsURL,
		                                &issue.Milestone.Creator.ReceivedEventsURL,
                                        &issue.Milestone.Creator.ReposURL,
                                        &issue.Milestone.Creator.StarredURL,
		                                &issue.Milestone.Creator.SubscriptionsURL,
		                            &issue.Milestone.OpenIssues,
                                    &issue.Milestone.ClosedIssues,
                                    &issue.Milestone.CreatedAt,
		                            &issue.Milestone.UpdatedAt,
                                    &issue.Milestone.ClosedAt,
                                    &issue.Milestone.DueOn,
		                        &issue.PullRequestLinks.URL, // NOTE: PullRequestLinks struct on Issue
                                    &issue.PullRequestLinks.HTMLURL,
		                            &issue.PullRequestLinks.DiffURL,
                                    &issue.PullRequestLinks.PatchURL,
		                        &issue.Repository.ID, // NOTE: Issue struct Repository field
    		                        &userID, // NOTE: Owner field on Repository struct
                                        &issue.Repository.Owner.ID,
                                        &issue.Repository.Owner.AvatarURL,
                                        &issue.Repository.Owner.HTMLURL,
    		                            &issue.Repository.Owner.GravatarID,
                                        &issue.Repository.Owner.Name,
                                        &issue.Repository.Owner.Company,
    		                            &issue.Repository.Owner.Blog,
                                        &issue.Repository.Owner.Location,
                                        &issue.Repository.Owner.Email,
    		                            &issue.Repository.Owner.Hireable,
                                        &issue.Repository.Owner.Bio,
                                        &issue.Repository.Owner.PublicRepos,
    		                            &issue.Repository.Owner.PublicGists,
                                        &issue.Repository.Owner.Followers,
                                        &issue.Repository.Owner.Following,
    		                            &issue.Repository.Owner.CreatedAt,
                                        &issue.Repository.Owner.UpdatedAt,
                                        &issue.Repository.Owner.SuspendedAt,
    		                            &issue.Repository.Owner.Type,
                                        &issue.Repository.Owner.SiteAdmin,
                                        &issue.Repository.Owner.TotalPrivateRepos,
    		                            &issue.Repository.Owner.OwnedPrivateRepos,
                                        &issue.Repository.Owner.PrivateGists,
    		                            &issue.Repository.Owner.DiskUsage,
                                        &issue.Repository.Owner.Collaborators,
    		                            &issue.Repository.Owner.Plan.Name, // NOTE: Repository Owner Plan info
                                            &issue.Repository.Owner.Plan.Space,
    		                                &issue.Repository.Owner.Plan.Collaborators,
                                            &issue.Repository.Owner.Plan.PrivateRepos,
    		                                &issue.Repository.Owner.URL,
                                            &issue.Repository.Owner.EventsURL,
                                            &issue.Repository.Owner.FollowingURL,
    		                                &issue.Repository.Owner.FollowersURL,
                                            &issue.Repository.Owner.GistsURL,
                                            &issue.Repository.Owner.OrganizationsURL,
    		                                &issue.Repository.Owner.ReceivedEventsURL,
                                            &issue.Repository.Owner.ReposURL,
                                            &issue.Repository.Owner.StarredURL,
    		                                &issue.Repository.Owner.SubscriptionsURL,
		                            &issue.Repository.Name, // NOTE: General Repository fields
                                        &issue.Repository.FullName,
                                        &issue.Repository.Description,
    		                            &issue.Repository.Homepage,
                                        &issue.Repository.DefaultBranch,
                                        &issue.Repository.MasterBranch,
    		                            &issue.Repository.CreatedAt,
                                        &issue.Repository.PushedAt,
                                        &issue.Repository.UpdatedAt,
    		                            &issue.Repository.HTMLURL,
                                        &issue.Repository.CloneURL,
                                        &issue.Repository.GitURL,
    		                            &issue.Repository.MirrorURL,
                                        &issue.Repository.SSHURL,
                                        &issue.Repository.SVNURL,
                                        &issue.Repository.Language,
                                        &issue.Repository.Fork,
                                        &issue.Repository.ForksCount,
                                        &issue.Repository.NetworkCount,
                                        &issue.Repository.OpenIssuesCount,
                                        &issue.Repository.StargazersCount,
                                        &issue.Repository.SubscribersCount,
                                        &issue.Repository.WatchersCount,
                                        &issue.Repository.Size,
                                        &issue.Repository.AutoInit,
                                    &parentID, // NOTE: This is the temporary solution to the sub-repo fields
                                    &sourcdID, // NOTE: This is the temporary solution to the sub-repo fields
                                    &issue.Repository.Organization.Login, // NOTE: Organization struct field
                                        &issue.Repository.Organization.ID,
                                        &issue.Repository.Organization.AvatarURL,
                                        &issue.Repository.Organization.HTMLURL,
                                        &issue.Repository.Organization.Name,
                                        &issue.Repository.Organization.Company,
                                        &issue.Repository.Organization.Blog,
                                        &issue.Repository.Organization.Location,
                                        &issue.Repository.Organization.Email,
                                        &issue.Repository.Organization.Description,
                                        &issue.Repository.Organization.PublicRepos,
                                        &issue.Repository.Organization.PublicGists,
                                        &issue.Repository.Organization.Followers,
                                        &issue.Repository.Organization.Following,
                                        &issue.Repository.Organization.CreatedAt,
                                        &issue.Repository.Organization.UpdatedAt,
                                        &issue.Repository.Organization.TotalPrivateRepos,
                                        &issue.Repository.Organization.OwnedPrivateRepos,
                                        &issue.Repository.Organization.PrivateGists,
                                        &issue.Repository.Organization.DiskUsage,
                                        &issue.Repository.Organization.Collaborators,
                                        &issue.Repository.Organization.BillingEmail,
                                        &issue.Repository.Organization.Type,
                                        &issue.Repository.Organization.Plan.Name, // NOTE: Plan struct for Org
                                            &issue.Repository.Organization.Plan.Space,
                                            &issue.Repository.Organization.Plan.Collaborators,
                                            &issue.Repository.Organization.Plan.PrivateRepos,
                                        &issue.Repository.Organization.URL,
                                        &issue.Repository.Organization.EventsURL,
                                        &issue.Repository.Organization.HooksURL,
                                        &issue.Repository.Organization.IssuesURL,
                                        &issue.Repository.Organization.MembersURL,
                                        &issue.Repository.Organization.PublicMembersURL,
                                        &issue.Repository.Organization.ReposURL,
                                    &issue.Repository.AllowRebaseMerge,
                                    &issue.Repository.AllowSquashMerge,
                                    &issue.Repository.AllowMergeCommit,
                                    &issue.Repository.License.Key, // NOTE: License struct
                                        &issue.Repository.License.Name,
                                        &issue.Repository.License.URL,
                                        &issue.Repository.License.SPDXID,
                                        &issue.Repository.License.HTMLURL,
                                        &issue.Repository.License.Featured,
                                        &issue.Repository.License.Description,
                                        &issue.Repository.License.Implementation,
                                        // NOTE: There are several fields here that are of
                                        // a 1:M relationship and rely on a FK which will
                                        // need separate tieup logic.
                                        &issue.Repository.License.Body,
                                    &issue.Repository.Private,
                                    &issue.Repository.HasIssues,
                                    &issue.Repository.HasWiki,
                                    &issue.Repository.HasPages,
                                    &issue.Repository.HasDownloads,
                                    &issue.Repository.LicenseTemplate,
                                    &issue.Repository.GitignoreTemplate,
                                    &issue.Repository.TeamID,
                                    &issue.Repository.URL,
                                    &issue.Repository.ArchiveURL,
                                    &issue.Repository.AssigneesURL,
                                    &issue.Repository.BlobsURL,
                                    &issue.Repository.BranchesURL,
                                    &issue.Repository.CollaboratorsURL,
                                    &issue.Repository.CommentsURL,
                                    &issue.Repository.CommitsURL,
                                    &issue.Repository.CompareURL,
                                    &issue.Repository.ContentsURL,
                                    &issue.Repository.ContributorsURL,
                                    &issue.Repository.DeploymentsURL,
                                    &issue.Repository.DownloadsURL,
                                    &issue.Repository.EventsURL,
                                    &issue.Repository.ForksURL,
                                    &issue.Repository.GitCommitsURL,
                                    &issue.Repository.GitRefsURL,
                                    &issue.Repository.GitTagsURL,
                                    &issue.Repository.HooksURL,
                                    &issue.Repository.IssueCommentURL,
                                    &issue.Repository.IssueEventsURL,
                                    &issue.Repository.IssuesURL,
                                    &issue.Repository.KeysURL,
                                    &issue.Repository.LabelsURL,
                                    &issue.Repository.LanguagesURL,
                                    &issue.Repository.MergesURL,
                                    &issue.Repository.MilestonesURL,
                                    &issue.Repository.NotificationsURL,
                                    &issue.Repository.PullsURL,
                                    &issue.Repository.ReleasesURL,
                                    &issue.Repository.StargazersURL,
                                    &issue.Repository.StatusesURL,
                                    &issue.Repository.SubscribersURL,
                                    &issue.Repository.SubscriptionURL,
                                    &issue.Repository.TagsURL,
                                    &issue.Repository.TreesURL,
                                    &issue.Repository.TeamsURL,
                                &issue.Reactions.TotalCount, // NOTE: Reactions struct on Issue struct
                                    &issue.Reactions.PlusOne,
                                    &issue.Reactions.MinusOne,
                                    &issue.Reactions.Laugh,
                                    &issue.Reactions.Confused,
                                    &issue.Reactions.Heart,
                                    &issue.Reactions.Hooray,
                                    &issue.Reactions.URL,
		); err != nil {
		    return nil, err
        }
		issues = append(issues, issue)
		issueID = *issue.ID
	}
    // TODO: Populating entire Issue struct (all 1:1 relationships).

	return issues, nil
}
