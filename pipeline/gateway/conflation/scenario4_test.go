package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

var TestScenario4 = Scenario4{}

var url = "https://www.rule-of-two.com/"
var pullRequest = github.PullRequest{IssueURL: &url}
var TestWithPullRequest = &ExpandedIssue{PullRequest: CRPullRequest{pullRequest, []int{}, []CRIssue{}}}
var TestWithoutPullRequest = &ExpandedIssue{}

func TestFilter4(t *testing.T) {
	withURL := TestScenario4.Filter(TestWithPullRequest)
	if withURL != false {
		t.Error(
			"PULL REQUEST WITH ASSOCIATED ISSUES INCLUDED",
		)
	}
	withoutURL := TestScenario4.Filter(TestWithoutPullRequest)
	if withoutURL != true {
		t.Error(
			"PULL REQUEST WITHOUT ISSUES EXCLUDED",
		)
	}
}
