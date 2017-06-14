package backend

import (
	"github.com/google/go-github/github"

	"testing"
)

func TestAddModel(t *testing.T) {
	BackendServer := new(BackendServer)
	id := 7
	repo := &github.Repository{ID: &id}
	if err := BackendServer.AddModel(repo); err != nil {
		t.Error("Error adding model to the BackendServer")
	}
	if len(BackendServer.Repos[id].Hive.Blender.Models) == 0 {
		t.Error("Model not added to Models slice on BackendServer")
	}
}
