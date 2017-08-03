package ingestor

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

type IngestorServer struct {
	Server          http.Server
	RepoInitializer RepoInitializer
}

func (i *IngestorServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != frontend.BackendSecret {
		utils.AppLog.Error("failed validating frontend-backend secret")
		return
	}
	repoInfo := r.FormValue("repos")
	// repoID, err := strconv.Atoi(string(repoInfo[0]))
	// if err != nil {
	// 	utils.AppLog.Error("converting repo ID: ", zap.Error(err))
	// 	http.Error(w, "failed converting repo ID", http.StatusForbidden)
	// 	return
	// }
	owner := string(repoInfo[1])
	repo := string(repoInfo[2])
	tokenString := r.FormValue("token")

	source := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokenString})
	token := oauth2.NewClient(oauth2.NoContext, source)
	client := *github.NewClient(token)

	isssueOpts := github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	issues := []*github.Issue{}
	for {
		gotIssues, resp, err := client.Issues.ListByRepo(context.Background(), owner, repo, &isssueOpts)
		if err != nil {
			utils.AppLog.Error("failed issue pull down: ", zap.Error(err))
			http.Error(w, "failed issue pull down", http.StatusForbidden)
			return
		}
		issues = append(issues, gotIssues...)
		if resp.NextPage == 0 {
			break
		} else {
			isssueOpts.ListOptions.Page = resp.NextPage
		}
	}

	pullsOpts := github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	pulls := []*github.PullRequest{}
	for {
		gotPulls, resp, err := client.PullRequests.List(context.Background(), owner, repo, &pullsOpts)
		if err != nil {
			utils.AppLog.Error("failed pull request pull down: ", zap.Error(err))
			http.Error(w, "failed pull request pull down", http.StatusForbidden)
			return
		}
		pulls = append(pulls, gotPulls...)
		if resp.NextPage == 0 {
			break
		} else {
			pullsOpts.ListOptions.Page = resp.NextPage
		}
	}

	bufferPool := NewPool()
	db := Database{BufferPool: bufferPool}
	db.Open()

	db.BulkInsertIssues(issues)
	db.BulkInsertPullRequests(pulls)
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	mux.HandleFunc("/activate-repos-ingestor", i.activateHandler)
	return mux
}

func (i *IngestorServer) Start() {
	i.RepoInitializer = RepoInitializer{}
	i.Server = http.Server{Addr: "127.0.0.1:8030", Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start", zap.Error(err))
	}
}

func (i *IngestorServer) Restart() error {
	db, err := bolt.Open("frontend/storage.db", 0644, nil)
	if err != nil {
		utils.AppLog.Error("failed opening bolt on ingestor restart", zap.Error(err))
		return err
	}
	defer db.Close()

	boltDB := frontend.BoltDB{DB: db}

	repos, tokens, err := boltDB.RetrieveBulk("token")
	if err != nil {
		utils.AppLog.Error("retrieve bulk tokens on ingestor restart", zap.Error(err))
	}

	allIssues := []*github.Issue{}

	for i := range tokens {
		token := oauth2.Token{}
		if err := json.Unmarshal(tokens[i], &token); err != nil {
			utils.AppLog.Error("converting tokens; ", zap.Error(err))
			return err
		}

		source := oauth2.StaticTokenSource(&token)
		oaClient := oauth2.NewClient(oauth2.NoContext, source)
		client := *github.NewClient(oaClient)

		repoID, err := strconv.Atoi(string(repos[i]))
		if err != nil {
			utils.AppLog.Error("repo id int conversion; ", zap.Error(err))
			return err
		}

		repo, _, err := client.Repositories.GetByID(context.Background(), repoID)
		if err != nil {
			utils.AppLog.Error("ingestor restart get by id; ", zap.Error(err))
			return err
		}

		owner := repo.Owner.Login
		name := repo.Name
		opts := github.IssueListByRepoOptions{
			State: "open",
			// Since: time.Time    // may be helpful if a "crash time" can be catalogued
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		issues := []*github.Issue{}
		for {
			gotIssues, resp, err := client.Issues.ListByRepo(context.Background(), *owner, *name, &opts)
			if err != nil {
				utils.AppLog.Error("restart issues pull down; ", zap.Error(err))
				return err
			}
			issues = append(issues, gotIssues...)
			if resp.NextPage == 0 {
				break
			} else {
				opts.ListOptions.Page = resp.NextPage
			}
		}

		allIssues = append(allIssues, issues...)

	}

	bufferPool := NewPool()
	memSQL := Database{BufferPool: bufferPool}
	memSQL.Open()
	defer memSQL.Close()

	memSQL.BulkInsertIssues(allIssues)

	return nil
}

func (i *IngestorServer) ContinuityCheck() {
	bufferPool := NewPool()
	db := Database{BufferPool: bufferPool}
	db.Open()

	// starts up
	// runs perpetually
	// queries on longer:
	// - time expires
	// - goroutine complete
	// query:
	// - all repos
	// - - issues of repos
	// - - pulls of repos
	// if gaps in issues/pulls IDs
	// - query GitHub w/ repo info
	// - update database w/ insert
}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
