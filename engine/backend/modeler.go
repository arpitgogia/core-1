package backend

import (
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
)

func (bs *BackendServer) NewModel(repoID int) error {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()
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
	bs.Repos.Actives[repoID].Hive.Blender.Models = append(bs.Repos.Actives[repoID].Hive.Blender.Models, &ArchModel{Model: &model, Conflator: &conflator})
	return nil
}
