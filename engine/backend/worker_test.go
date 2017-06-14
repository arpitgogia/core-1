package backend

import "testing"

func TestNewWorker(t *testing.T) {
	testID := 1
	channel := make(chan chan *RepoData)
	testWorker := NewWorker(testID, channel)
	if testWorker.ID != testID {
		t.Error("Failure creating new worker object")
	}
}
