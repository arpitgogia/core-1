package retriever

import (
	"github.com/google/go-github/github"

	"coralreefci/models"
)

var Workers chan chan github.Issue

type Dispatcher struct {
	Models map[int]models.Model
}

func (d *Dispatcher) Start(count int) {
	Workers = make(chan chan github.Issue, count)
	for i := 0; i < count; i++ {
		worker := NewWorker(i+1, Workers)
		worker.Models = d.Models
		worker.Start()
	}

	go func() {
		for {
			work := <-Workload
			go func() {
				workers := <-Workers
				workers <- work
			}()
		}
	}()
}
