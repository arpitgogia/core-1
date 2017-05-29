package main

import (
	"coralreefci/analysis/cmd/endtoendtests/replay"
	"coralreefci/engine/ingestor"
	"flag"
)

func main() {

	loadArchiveFlag := flag.Bool("loadarchive", false, "load archive into the database")
	archivePathFlag := flag.String("archivepath", "", "location of github archive")
	flag.Parse()

	bufferPool := ingestor.NewPool()
	db := ingestor.Database{BufferPool: bufferPool}
	db.Open()
	bs := replay.BacktestServer{DB: &db}
	go bs.Start()

	if *loadArchiveFlag && *archivePathFlag != "" {
		bs.LoadArchive(*archivePathFlag)
	}
}
