package retriever

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestNewWorker(t *testing.T) {
	testID := 1
	channel := make(chan chan github.Issue)
	testWorker := NewWorker(testID, channel)
	if testWorker.ID != testID {
		t.Error("Failure creating new worker object")
	}
}
