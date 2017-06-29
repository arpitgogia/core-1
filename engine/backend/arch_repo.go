package backend

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type ArchModel struct {
	Model     *models.Model
	Conflator *conflation.Conflator
}

type Blender struct {
	Models []*ArchModel
}

type ArchHive struct {
	Blender *Blender
}

type ArchRepo struct {
	Hive   *ArchHive
	Client *github.Client
}

func (bs *BackendServer) NewArchRepo(repoID int) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	ctx := &conflation.Context{}
	scn := []conflation.Scenario{&conflation.Scenario3{}}
	algo := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{
			Context: ctx,
		},
	}
	norm := conflation.Normalizer{Context: ctx}
	conf := conflation.Conflator{
		Scenarios:            scn,
		ConflationAlgorithms: algo,
		Normalizer:           norm,
		Context:              ctx,
	}
	model := ArchModel{Conflator: &conf}

	bs.Repos.Actives[repoID] = &ArchRepo{
		Hive: &ArchHive{
			Blender: &Blender{
				Models: []*ArchModel{&model},
			},
		},
	}
}

func (bs *BackendServer) NewClient(repoID int, token *oauth2.Token) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	tokenSource := oauth2.StaticTokenSource(token)
	authClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	githubClient := github.NewClient(authClient)

	bs.Repos.Actives[repoID].Client = githubClient
}
