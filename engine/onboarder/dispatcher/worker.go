package retriever

import (
	"fmt"

	"github.com/google/go-github/github"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type Worker struct {
	ID     int
	Work   chan interface{}
	Queue  chan chan interface{}
	Models map[int]models.Model
	Quit   chan bool
}

func NewWorker(id int, queue chan chan interface{}) Worker {
	return Worker{
		ID:     id,
		Work:   make(chan interface{}),
		Queue:  queue,
		Models: make(map[int]models.Model),
		Quit:   make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Work
            object := <-w.Work
            if _, ok := object.(github.Issue); ok {
                issue := object.(github.Issue)
            } else {
                pull := object.(github.PullRequest)
            }
			select {
			case issue:
				if object.ClosedAt != nil {
                    // TODO: CONFLATE + LEARN
				} else {
					// TODO: Call Predict Method
					// assignees := w.Models[*issuesEvent.Repo.ID].Algorithm.Predict(expandedIssue)
					// NOTE: This is likely where the assignment function will be called.
					// assignment.AssignContributor(assignees[0], issuesEvent, testClient())
					// HACK: using test client
				}
            case pull:
                // STUFF GOES HERE
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
