package onboarder

import (
	"github.com/google/go-github/github"

	"coralreefci/engine/gateway"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
)

func (rs *RepoServer) AddModel(repo *github.Repository, client *github.Client) error {
	name := *repo.Name
	owner := *repo.Owner.Login
	repoID := *repo.ID
	//TODO: The comments field is not cached when using CachedGateway and will
	//      need to be fixed eventually.
	newGateway := gateway.Gateway{Client: client}
	githubIssues, err := newGateway.GetIssues(owner, name)
	if err != nil {
		// utils.Log.Error("Cannot get Issues from Gateway. ", err)
	}
	githubPulls, err := newGateway.GetPullRequests(owner, name)
	if err != nil {
		// utils.Log.Error("Cannot get PullRequests from Gateway. ", err)
	}

	context := &conflation.Context{}
	scenarios := []conflation.Scenario{&conflation.Scenario2{}}
	conflationAlgorithms := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{Context: context},
	}
	normalizer := conflation.Normalizer{Context: context}
	conflator := conflation.Conflator{
		Scenarios:            scenarios,
		ConflationAlgorithms: conflationAlgorithms,
		Normalizer:           normalizer,
		Context:              context,
	}

	issuesCopy := make([]github.Issue, len(githubIssues))
	pullsCopy := make([]github.PullRequest, len(githubPulls))

	// TODO: Evaluate this particular snippet of code as it has potential
	//       performance optimization capabilities related to the hardware
	//       level. This may ultimately live in the actual gateway.go file to
	//	     improve the actual download operations.
	for i := 0; i < len(issuesCopy); i++ {
		issuesCopy[i] = *githubIssues[i]
	}
	for i := 0; i < len(pullsCopy); i++ {
		pullsCopy[i] = *githubPulls[i]
	}

	conflator.Context.Issues = []conflation.ExpandedIssue{}
	conflator.SetIssueRequests(issuesCopy)
	conflator.SetPullRequests(pullsCopy)
	conflator.Conflate()

	trainingSet := []conflation.ExpandedIssue{}

	for i := 0; i < len(conflator.Context.Issues); i++ {
		expandedIssue := conflator.Context.Issues[i]
		if expandedIssue.Conflate {
			if expandedIssue.Issue.Assignee == nil {
				continue
			} else {
				trainingSet = append(trainingSet, conflator.Context.Issues[i])
			}
		}
	}
	// TODO: This will likely become a read from the MemSQL database which will
	//       hold the desired state of each model as a table.
	model := models.Model{Algorithm: &bhattacharya.NBModel{}}
	model.Algorithm.Learn(trainingSet)
	rs.Repos[repoID].Hive.Blender.Models = append(rs.Repos[repoID].Hive.Blender.Models, &ArchModel{Model: &model})
	return nil
}
