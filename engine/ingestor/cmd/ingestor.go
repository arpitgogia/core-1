package main

import (
	"coralreefci/engine/gateway"
	"coralreefci/engine/ingestor"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Temp Replay Server
type BabyReplayServer struct{}

func (b *BabyReplayServer) Replay() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "23fc398670a80700b19b1ae1587825a16aa8ce57"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: client}, DiskCache: &gateway.DiskCache{}}

	githubIssues, _ := newGateway.GetIssues("dotnet", "corefx")
	githubPulls, _ := newGateway.GetPullRequests("dotnet", "corefx")

	db := ingestor.Database{}
	db.Open()

	db.BulkInsertIssues(githubIssues)
	db.BulkInsertPullRequests(githubPulls)
}

func main() {
	dispatcher := ingestor.Dispatcher{}
	dispatcher.Start(5)

	ingestorServer := ingestor.IngestorServer{}
	ingestorServer.Start()
}
