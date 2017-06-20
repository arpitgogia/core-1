package backend

type Worker struct {
	ID    int
	Work  chan *RepoData
	Queue chan chan *RepoData
	Repos *ActiveRepos
	Quit  chan bool
}

func (bs *BackendServer) NewWorker(workerID int, queue chan chan *RepoData) Worker {
	return Worker{
		ID:    workerID,
		Work:  make(chan *RepoData),
		Queue: queue,
		Repos: bs.Repos,
		Quit:  make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Work
			select {
			case repodata := <-w.Work:
				w.Repos.Lock()

				if w.Repos.Actives[repodata.RepoID] != nil {
					// - Generate new ArchRepo + assign
					// - Create new Client + assign
					// - Add model to sub struct in ArchRepo
				}

				if len(repodata.Open) != 0 {
					w.Repos.Actives[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetIssueRequests(repodata.Open)
				}
				if len(repodata.Closed) != 0 {
					w.Repos.Actives[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetIssueRequests(repodata.Closed)
				}
				if len(repodata.Pulls) != 0 {
					w.Repos.Actives[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetPullRequests(repodata.Pulls)
				}
				w.Repos.Actives[repodata.RepoID].Hive.Blender.Models[0].Conflator.Conflate()
				// TODO: Call Learn/Predict.
				// TODO: Add in Assignment call.
				w.Repos.Unlock()
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
