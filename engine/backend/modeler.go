package backend

import (
	"github.com/google/go-github/github"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
)

func (bs *BackendServer) AddModel(repo *github.Repository) error {
	repoID := *repo.ID
	context := &conflation.Context{}
	scenarios := []conflation.Scenario{&conflation.Scenario2{}}
	algos := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{Context: context},
	}
	normalizer := conflation.Normalizer{Context: context}
	conflator := conflation.Conflator{
		Scenarios:            scenarios,
		ConflationAlgorithms: algos,
		Normalizer:           normalizer,
		Context:              context,
	}
	model := models.Model{Algorithm: &bhattacharya.NBModel{}}
	if bs.Repos == nil {
		bs.Repos = make(map[int]*ArchRepo)
	}
	if _, ok := bs.Repos[repoID]; !ok {
		bs.Repos[repoID] = &ArchRepo{Hive: &ArchHive{Blender: &Blender{}}}
	}
	bs.Repos[repoID].Hive.Blender.Models = append(bs.Repos[repoID].Hive.Blender.Models, &ArchModel{Model: &model, Conflator: &conflator})
	return nil
}
