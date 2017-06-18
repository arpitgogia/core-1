package main

import (
	"coralreefci/analysis/cmd/endtoendtests/replay"
	"coralreefci/engine/ingestor"
	"coralreefci/engine/onboarder/signup"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"time"
)

func main() {
	runBacktestFlag := flag.Bool("runbacktest", false, "runs the end to end backtest")
	loadArchiveFlag := flag.Bool("loadarchive", false, "load archive into the database")
	archivePathFlag := flag.String("archivepath", "", "location of github archive")
	flag.Parse()

	if !*runBacktestFlag && !*loadArchiveFlag {
		fmt.Println("Usage: ./cmd --loadarchive=true --archivepath=/home/michael/Data/GithubArchive/")
		fmt.Println("Usage: ./cmd --runbacktest=true")
		return
	}

	bufferPool := ingestor.NewPool()
	db := ingestor.Database{BufferPool: bufferPool}
	db.Open()

	bs := replay.BacktestServer{DB: &db}
	go bs.Start()

	dispatcher := ingestor.Dispatcher{}
	dispatcher.Start(5)
	ingestorServer := ingestor.IngestorServer{}
	go ingestorServer.Start()

	if *loadArchiveFlag && *archivePathFlag != "" {
		bs.LoadArchive(*archivePathFlag)
	}

	if *runBacktestFlag {
		repoServer := signup.RepoServer{}
		go repoServer.Start()

		bs.AddRepo(26295345, "dotnet", "corefx")
		bs.AddRepo(724712, "rust-lang", "rust")
		bs.StreamWebhookEvents()
		time.Sleep(15 * time.Second)

		repo1 := &github.Repository{ID: github.Int(26295345), Organization: &github.Organization{Name: github.String("dotnet")}, Name: github.String("corefx")}
		repoServer.AddModel(repo1)
	}
}
