package backend

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"

	// "coralreefci/engine/frontend"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type ArchModel struct {
	Model     *models.Model
	Conflator *conflation.Conflator
	// Benchmark        Benchmark // TODO: Struct to build
	// Scenarios        []conflation.Scenario
	// PilotScenarios   []conflation.Scenario
	// LearnedScenarios []conflation.Scenario
	// StrategyParams   StrategyParams TODO: Baseline with self evolving
	//                                       parameters: Tossing Graph?,
	//                                       Conflation Scenarios, etc.)
}

type Blender struct {
	Models []*ArchModel
	// PilotModels []*ArchModel
}

type ArchHive struct {
	Blender *Blender
	// TossingGraph       TossingGraphAlgorithm // TODO: Struct to build
	// StrategyParams     StrategyParams // TODO: Struct to build
	// AggregateBenchmark Benchmark
}

type ArchRepo struct {
	Repo   *github.Repository // NOTE: This is a deprediated field; remove.
	Hive   *ArchHive
	Client *github.Client
}

func (bs *BackendServer) NewArchRepo(repoID int) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()
	bs.Repos.Actives[repoID] = &ArchRepo{Hive: &ArchHive{Blender: &Blender{}}}
}

func (bs *BackendServer) NewClient(repoID int, token *oauth2.Token) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	oaConfig := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     ghoa.Endpoint,
		Scopes:       []string{"admin:repo_hook", "repo:status", "public_repo"}, // NOTE: Scopes may be reduced (e.g. remove hook).
	}

	oaClient := oaConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oaClient)
	bs.Repos.Actives[repoID].Client = client
}

// TODO: Instantiate the Conflator struct on the ArchRepo.

// TODO:
// Below are several potential helper methods for the ArchRepo:
// BootstrapModel() - performs preliminary training / assignments / startup
// GetModelBenchmark() TODO: Calculate AggregateBenchmark for this method
// Assign(issue github.Issue) - assign newly raised issue to contributor
// Enable()
// Disable()
// Destroy()
// InactiveDevelopers() []string
// Even though the Github API will prevent it. Heupr could try to assign an
// issue to an inactive developer. It might not be great for team morale if
// Heupr inadvertently exposed that in a UI/dashboard page.
