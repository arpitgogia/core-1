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
        LEFT JOIN user AS creator ON issue.user_id=creator.id
        LEFT JOIN plan AS creator_plan ON creator.plan_id=creator_plan.id
        LEFT JOIN text_match AS creator_tm ON creator_tm.text_match_id=creator.text_match_id
        LEFT JOIN matches AS creator_m ON creator_tm.id=creator_m.text_match_fk
        LEFT JOIN indices AS creator_i ON creator_m.id=creator_i.match_fk
        LEFT JOIN labels ON issue.id=labels.issue_fk
        LEFT JOIN user AS assignee ON issue.assignee_id=assignee.id
        LEFT JOIN plan AS assignee_plan ON assignee.plan_id=assignee_plan.id
        LEFT JOIN text_match AS assignee_tm ON assignee_tm.text_match_id=assignee.text_match_id
        LEFT JOIN matches AS assignee_m ON assignee_tm.id=assignee_m.text_match_fk
        LEFT JOIN indices AS assignee_i ON assignee_m.id=assignee_i.match_fk
        LEFT JOIN user AS closed_by ON issue.closed_by_id=closed_by.id
        LEFT JOIN plan AS closed_by_plan ON closed_by.plan_id=closed_by_plan.id
        LEFT JOIN text_match AS closed_by_tm ON closed_by_tm.text_match_id=closed_by.text_match_id
        LEFT JOIN matches AS closed_by_m ON closed_by_tm.id=closed_by_m.text_match_fk
        LEFT JOIN indices AS closed_by_i ON closed_by_m.id=closed_by_i.match_fk
        LEFT JOIN milestones ON issue.milestone_id=milestones.id
        LEFT JOIN pull_request_links ON issue.pull_request_links_id=pull_request_links.id
        LEFT JOIN repository AS repo ON issue.repository_id=repo.id
        LEFT JOIN user AS repo_owner ON repo.user_id=repo_owner.id
        LEFT JOIN plan AS owner_plan ON repo_owner.plan_id=owner_plan.plan_id
        LEFT JOIN text_match AS repo_owner_tm ON repo_owner_tm.text_match_id=repo_owner.text_match_id
        LEFT JOIN matches AS repo_owner_m ON repo_owner_tm.id=repo_owner_m.text_match_fk
        LEFT JOIN indices AS repo_owner_i ON repo_owner_m.id=repo_owner_i.match_fk
        LEFT JOIN repository AS parent ON issue.parent_id=parent.id
        LEFT JOIN user AS parent_owner ON parent.user_id=parent_owner.id
        LEFT JOIN plan AS parent_plan ON parent_owner.plan_id=parent_plan.plan_id
        LEFT JOIN license AS parent_license ON parent.license_id=parent_license.id
        LEFT JOIN license_permissions AS parent_lp ON parent_license.permissions_id=parent_lp.id
        LEFT JOIN license_conditions AS parent_lc ON parent_license.conditions_id=parent_lc.id
        LEFT JOIN license_limitations AS parent_ll ON parent_license.limitations_id=parent_ll.id
        LEFT JOIN text_match AS parent_owner_tm ON parent_owner_tm.text_match_id=repo_owner.text_match_id
        LEFT JOIN matches AS parent_owner_m ON parent_owner_tm.id=parent_owner_m.text_match_fk
        LEFT JOIN indices AS parent_owner_i ON parent_owner_m.id=parent_owner_i.match_fk
        LEFT JOIN repository AS source ON issue.source_id=source.id
        LEFT JOIN user AS source_owner ON source.user_id=source_owner.id
        LEFT JOIN plan AS source_plan ON source_owner.plan_id=source_plan.plan_id
        LEFT JOIN license AS source_license ON source.license_id=source_license.id
        LEFT JOIN license_permissions AS source_lp ON source_license.permissions_id=source_lp.id
        LEFT JOIN license_conditions AS source_lc ON source_license.conditions_id=source_lc.id
        LEFT JOIN license_limitations AS source_ll ON source_license.limitations_id=source_ll.id
        LEFT JOIN text_match AS source_owner_tm ON source_owner_tm.text_match_id=repo_owner.text_match_id
        LEFT JOIN matches AS source_owner_m ON source_owner_tm.id=source_owner_m.text_match_fk
        LEFT JOIN indices AS source_owner_i ON source_owner_m.id=source_owner_i.match_fk
        LEFT JOIN organization AS org ON repository.organization_id=org.id
        LEFT JOIN license AS repo_license ON repo.license_id=repo_license.id
        LEFT JOIN license_permissions AS repo_lp ON repo_license.permissions_id=repo_lp.id
        LEFT JOIN license_conditions AS repo_lc ON repo_license.conditions_id=repo_lc.id
        LEFT JOIN license_limitations AS repo_ll ON repo_license.limitations_id=repo_ll.id
        LEFT JOIN text_match AS repo_tm ON repo_tm.text_match_id=assignees_member.text_match_id
        LEFT JOIN matches AS repo_m ON repo_tm.id=repo_m.text_match_fk
        LEFT JOIN indices AS repo_i ON repo_m.id=repo_i.match_fk
        LEFT JOIN reactions ON issue.reactions_id=reactions.id
        LEFT JOIN assignees ON issue.id=assignees.issue_fk
        LEFT JOIN user AS assignees_member ON assignees.user_id=assignees_member.id
        LEFT JOIN plan AS assignees_member_plan ON assignees.plan_id=assignees_member_plan.id
        LEFT JOIN text_match AS assignees_member_tm ON assignees_member_tm.text_match_id=assignees_member.text_match_id
        LEFT JOIN matches AS assignees_member_m ON assignees_member_tm.id=assignees_member_m.text_match_fk
        LEFT JOIN indices AS assignees_member_i ON assignees_member_m.id=assignees_member_i.match_fk
        LEFT JOIN text_match AS issue_tm ON issue_tm.text_match_id=issue.text_match_id
        LEFT JOIN matches AS issue_m ON issue_tm.id=issue_m.text_match_fk
        LEFT JOIN indices AS issue_i ON issue_m.id=issue_i.match_fk
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
